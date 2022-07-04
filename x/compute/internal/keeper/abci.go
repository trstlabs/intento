package keeper

import (

	//"log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	sdktxsigning "github.com/cosmos/cosmos-sdk/types/tx/signing"

	"github.com/trstlabs/trst/x/compute/internal/types"
)

// SelfExecute executes the contract instance from the chain itself
func (k Keeper) SelfExecute(ctx sdk.Context, contractAddress sdk.AccAddress, msg []byte, callbackSig []byte) (uint64, error) {
	ctx.GasMeter().ConsumeGas(types.InstanceCost, "Loading compute module: execute")

	signBytes := []byte{}
	signMode := sdktxsigning.SignMode_SIGN_MODE_UNSPECIFIED
	modeInfoBytes := []byte{}
	pkBytes := []byte{}
	signerSig := []byte{}
	var err error

	// If no callback signature - we should not execute
	if callbackSig == nil {
		return 0, sdkerrors.Wrap(types.ErrExecuteFailed, "no callback sig")
	}

	verificationInfo := types.NewVerificationInfo(signBytes, signMode, modeInfoBytes, pkBytes, signerSig, callbackSig)

	codeInfo, prefixStore, err := k.contractInstance(ctx, contractAddress)
	if err != nil {
		return 0, err
	}

	store := ctx.KVStore(k.storeKey)

	contractKey := store.Get(types.GetContractEnclaveKey(contractAddress))
	if contractKey == nil {
		return 0, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "contract key not found")
	}
	params := types.NewEnv(ctx, contractAddress, sdk.NewCoins(sdk.NewCoin(types.Denom, sdk.ZeroInt())), contractAddress, contractKey)

	// prepare querier
	querier := QueryHandler{
		Ctx:     ctx,
		Plugins: k.queryPlugins,
	}

	gas := gasForContract(ctx)
	res, gasUsed, execErr := k.wasmer.Execute(codeInfo.CodeHash, params, msg, prefixStore, cosmwasmAPI, querier, gasMeter(ctx), gas, verificationInfo)

	if execErr != nil {
		return 0, sdkerrors.Wrap(types.ErrExecuteFailed, execErr.Error())
	}

	// emit all events from this contract itself
	events := types.ParseEvents(res.Log, contractAddress)
	ctx.EventManager().EmitEvents(events)

	err = k.dispatchMessages(ctx, contractAddress, res.Messages)
	if err != nil {
		return 0, err
	}
	return gasUsed, nil

}

// DeductFeesAndFundCreator handles remaining contract balance
func (k Keeper) DeductFeesAndFundCreator(ctx sdk.Context, contractAddress sdk.AccAddress, gas uint64) error {

	store := ctx.KVStore(k.storeKey)
	contractBz := store.Get(types.GetContractAddressKey(contractAddress))
	if contractBz == nil {
		return sdkerrors.Wrap(types.ErrNotFound, "contract")
	}
	var contract types.ContractInfo
	k.cdc.MustUnmarshal(contractBz, &contract)
	gasUsed := (gas / types.GasMultiplier) + 1

	//TODO multiply with a multiplier param
	gasCoin := sdk.NewCoin(types.Denom, sdk.NewIntFromUint64(gasUsed/100000))

	//direct a commission of the utrst contract balance towards the community pool
	var feeCoins sdk.Coins

	contractBalance := k.bankKeeper.GetAllBalances(ctx, contractAddress)
	p := k.GetParams(ctx)
	if contractBalance.Empty() {
		feeCoins = sdk.NewCoins(sdk.NewCoin(types.Denom, sdk.NewInt(p.AutoMsgConstantFee)).Add(gasCoin))
		err := k.distrKeeper.FundCommunityPool(ctx, feeCoins, contract.Creator)
		if err != nil {
			return err
		}
		return nil
	}
	percentageAutoMsgFundsCommission := sdk.NewDecWithPrec(p.AutoMsgFundsCommission, 2)
	amountAutoMsgFundsCommission := percentageAutoMsgFundsCommission.MulInt(contractBalance.AmountOf(types.Denom)).Ceil().TruncateInt()
	feeCoins = sdk.NewCoins(sdk.NewCoin(types.Denom, sdk.NewInt(p.AutoMsgConstantFee)).Add(sdk.NewCoin(types.Denom, amountAutoMsgFundsCommission).Add(gasCoin)))
	//if the contract is not able to pay, the contract creator pays as next in line
	err := k.distrKeeper.FundCommunityPool(ctx, feeCoins, contractAddress)
	if err != nil {
		err := k.distrKeeper.FundCommunityPool(ctx, feeCoins, contract.Creator)
		if err != nil {
			return err
		}
	}
	if contractBalance.Sub(feeCoins).AmountOf(types.Denom).IsPositive() {
		//pay out the remaining balance after deducting fee, commision and gas cost to the contract creator
		err = k.bankKeeper.SendCoins(ctx, contractAddress, contract.Creator, contractBalance.Sub(feeCoins))
		if err != nil {
			return err
		}
	}
	return nil
}

// SetIncentiveCoins distributes compute module allocated coins to selected contracts
func (k Keeper) SetIncentiveCoins(ctx sdk.Context, addressList []string) {
	params := k.GetParams(ctx)

	total := k.bankKeeper.GetBalance(ctx, k.accountKeeper.GetModuleAddress("compute"), types.Denom)
	k.Logger(ctx).Info("contract incentive", "total", total)

	amount := total.Amount.QuoRaw(int64(len(addressList)))
	if amount.Int64() > params.MaxContractIncentive {
		amount = sdk.NewInt(params.MaxContractIncentive)
	}
	k.Logger(ctx).Info("sent", "amount", amount)

	for _, addr := range addressList {
		sdkAddr, _ := sdk.AccAddressFromBech32(addr)
		err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdkAddr, sdk.NewCoins(sdk.NewCoin(types.Denom, amount)))
		if err != nil {
			k.Logger(ctx).Info("sending", "err", err)
			break
		}

		k.Logger(ctx).Info("allocated", "contract", addr, "coins", amount)
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeDistributedToContract,
				sdk.NewAttribute(types.AttributeKeyAddress, addr),
			),
		)

	}

}
