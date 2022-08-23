package keeper

import (

	//"log"

	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	sdktxsigning "github.com/cosmos/cosmos-sdk/types/tx/signing"
	wasmTypes "github.com/trstlabs/trst/go-cosmwasm/types"
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

	contractInfo, codeInfo, prefixStore, err := k.contractInstance(ctx, contractAddress)
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
	res, gasUsed, execErr := k.wasmer.Execute(codeInfo.CodeHash, params, msg, prefixStore, cosmwasmAPI, querier, gasMeter(ctx), gas, verificationInfo, wasmTypes.HandleTypeExecute)

	if execErr != nil {
		return 0, sdkerrors.Wrap(types.ErrExecuteFailed, execErr.Error())
	}

	// emit all events from this contract itself
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeExecute,
		sdk.NewAttribute(types.AttributeKeyContractAddr, contractAddress.String())))

	_, err = k.handleContractResponse(ctx, contractAddress, contractInfo.IBCPortID, res, res.Messages, res.Events, res.Data, msg, verificationInfo)
	if err != nil {
		return 0, err
	}
	return gasUsed, nil

}

// DistributeCoins distributes AutoMessage fees and handles remaining contract balance
func (k Keeper) DistributeCoins(ctx sdk.Context, contract types.ContractInfoWithAddress, gas uint64, isRecurring bool) error {

	p := k.GetParams(ctx)

	//We have 2 constant fees + this gas-dependent fee.
	autoMsgFlexFee := sdk.NewCoin(types.Denom, sdk.NewIntFromUint64(gas).Mul(sdk.NewInt(100)).Quo(sdk.NewInt(p.AutoMsgFlexFeeDenom)))
	fmt.Printf("autoMsgFlexFee  %v\n", autoMsgFlexFee)
	//direct a commission of the utrst contract balance towards the community pool

	contractBalance := k.bankKeeper.GetAllBalances(ctx, contract.Address)

	//depending on the type of self-execution the constant fee may differ (gov param)
	constantFee := sdk.NewInt(p.AutoMsgConstantFee)
	if isRecurring {
		constantFee = sdk.NewInt(p.RecurringAutoMsgConstantFee)
	}

	percentageAutoMsgFundsCommission := sdk.NewDecWithPrec(p.AutoMsgFundsCommission, 2)
	amountAutoMsgFundsCommission := percentageAutoMsgFundsCommission.MulInt(contractBalance.AmountOf(types.Denom)).Ceil().TruncateInt()
	feeCoins := sdk.NewCoins(sdk.NewCoin(types.Denom, constantFee).Add(sdk.NewCoin(types.Denom, amountAutoMsgFundsCommission).Add(autoMsgFlexFee)))
	fmt.Printf("fee coins %v\n", feeCoins)

	//the contract should be funded with the fee. Iif the contract is not able to pay, the contract owner pays next in line
	err := k.distrKeeper.FundCommunityPool(ctx, feeCoins, contract.Address)
	if err != nil {
		store := ctx.KVStore(k.storeKey)
		// if a contract instantiated the contract, we do not deduct fees from it and the AutoMsg won't be written to Cache
		if !store.Has(types.GetContractEnclaveKey(contract.Owner)) {
			err := k.distrKeeper.FundCommunityPool(ctx, feeCoins, contract.ContractInfo.Owner)
			if err != nil {
				return err
			}
		}
	}

	fmt.Printf("contractBalance %v\n", contractBalance)
	//pay out the remaining balance to the contract owner after deducting fee, commision and gas cost
	contractBalance.Sort()
	toOwnerCoins, positive := contractBalance.SafeSub(feeCoins.Add(feeCoins[0]))
	if positive {
		err = k.bankKeeper.SendCoins(ctx, contract.Address, contract.ContractInfo.Owner, toOwnerCoins)
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
	k.Logger(ctx).Info("TC reward", "total", total)

	amount := total.Amount.QuoRaw(int64(len(addressList)))
	if amount.Int64() > params.MaxContractIncentive {
		amount = sdk.NewInt(params.MaxContractIncentive)
	}

	k.Logger(ctx).Info("TC reward", "amount", amount)

	for _, addr := range addressList {
		contrAddr, _ := sdk.AccAddressFromBech32(addr)
		err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, contrAddr, sdk.NewCoins(sdk.NewCoin(types.Denom, amount)))
		if err != nil {
			k.Logger(ctx).Info("TC reward", "send_err", err)
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
