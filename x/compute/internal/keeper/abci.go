package keeper

import (

	//"log"

	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	sdktxsigning "github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
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
	p := types.NewEnv(ctx, contractAddress, sdk.NewCoins(sdk.NewCoin(types.Denom, sdk.ZeroInt())), contractAddress, contractKey)

	// prepare querier
	querier := QueryHandler{
		Ctx:     ctx,
		Plugins: k.queryPlugins,
	}

	gas := gasForContract(ctx)
	res, gasUsed, _, execErr := k.wasmer.Execute(codeInfo.CodeHash, p, msg, prefixStore, cosmwasmAPI, querier, gasMeter(ctx), gas, verificationInfo, wasmTypes.HandleTypeExecute)
	consumeGas(ctx, gasUsed)
	if execErr != nil {
		return 0, sdkerrors.Wrap(types.ErrExecuteFailed, execErr.Error())
	}

	// emit all events from contract itself
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeExecute,
		sdk.NewAttribute(types.AttributeKeyContractAddr, contractAddress.String())))

	_, err = k.handleContractResponse(ctx, contractAddress, contractInfo.IBCPortID, res.Attributes, res.Messages, res.Events, res.Data, msg, verificationInfo)
	if err != nil {
		return 0, err
	}
	return gasUsed, nil

}

// DistributeCoins distributes AutoMessage fees and handles remaining contract balance
func (k Keeper) DistributeCoins(ctx sdk.Context, contract types.ContractInfoWithAddress, gasUsed uint64, isRecurring bool) error {

	p := k.GetParams(ctx)

	//We have 2 constant fees + this gas-dependent fee.
	// gas dependent fee goes to validator
	//flexFeeMultiplier := sdk.NewDec(100).QuoInt64(p.AutoMsgFlexFeeMul)
	//flexFee := sdk.NewDec(int64(gasUsed)).MulTruncate(flexFeeMultiplier)

	flexFeeMultiplier := sdk.NewDec(p.AutoMsgFlexFeeMul).QuoInt64(100)
	flexFee := sdk.NewDecFromInt(sdk.NewInt(int64(gasUsed))).Mul(flexFeeMultiplier)
	//flexFeeCoin := sdk.NewCoin(types.Denom, flexFee.TruncateInt())
	//direct a commission of the utrst contract balance towards the community pool
	contractBalance := k.bankKeeper.GetAllBalances(ctx, contract.Address)

	//depending on if self-execution is recurring the constant fee may differ (gov param)
	constantFee := sdk.NewInt(p.AutoMsgConstantFee)
	if isRecurring {
		constantFee = sdk.NewInt(p.RecurringAutoMsgConstantFee)
	}

	percentageAutoMsgFundsCommission := sdk.NewDecWithPrec(p.AutoMsgFundsCommission, 2)
	amountAutoMsgFundsCommissionCoin := sdk.NewCoin(types.Denom, percentageAutoMsgFundsCommission.MulInt(contractBalance.AmountOf(types.Denom)).Ceil().TruncateInt())
	feeCommunityCoins := sdk.NewCoins(sdk.NewCoin(types.Denom, constantFee).Add(amountAutoMsgFundsCommissionCoin))
	totalExecCost := feeCommunityCoins.Add(sdk.NewCoin(types.Denom, flexFee.TruncateInt()))

	if contract.Duration >= p.MinContractDurationForIncentive {
		incentive, err := k.ContractIncentive(ctx, totalExecCost[0], contract.Address)
		if err != nil {
			fmt.Printf("err 0 %v\n", err)
			return err
		}
		contractBalance.Add(incentive)
	}

	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, contract.Address, authtypes.FeeCollectorName, sdk.NewCoins(sdk.NewCoin(types.Denom, flexFee.TruncateInt())))
	if err != nil {
		fmt.Printf("err 5 %v\n", err)
		return err
	}
	//consumeGas(ctx, flexFee.TruncateInt().Uint64())

	//the contract should be funded with the fee. Iif the contract is not able to pay, the contract owner pays next in line
	err = k.distrKeeper.FundCommunityPool(ctx, feeCommunityCoins, contract.Address)
	if err != nil {
		fmt.Printf("err %v\n", err)
		store := ctx.KVStore(k.storeKey)
		// if a contract instantiated the contract, we do not deduct fees from it and the Auto Exec Msg won't be written to Cache
		if !store.Has(types.GetContractEnclaveKey(contract.Owner)) {
			err := k.distrKeeper.FundCommunityPool(ctx, feeCommunityCoins, contract.ContractInfo.Owner)
			if err != nil {
				fmt.Printf("err 1 %v\n", err)
				return err
			}
		}
		fmt.Printf("err 2 %v\n", err)
		return err
	}

	fmt.Printf("contractBalance %v\n", contractBalance)
	fmt.Printf("feeCommunityCoins %v\n", feeCommunityCoins)

	//pay out the remaining balance to the contract owner after deducting fee, commision and gas cost
	toOwnerCoins, negative := contractBalance.Sort().SafeSub(totalExecCost)
	fmt.Printf("toOwnerCoins %v\n", toOwnerCoins)
	if !negative {
		err = k.bankKeeper.SendCoins(ctx, contract.Address, contract.ContractInfo.Owner, toOwnerCoins)
		if err != nil {
			fmt.Printf("err 3 %v\n", err)
			return err
		}

	}
	return nil
}

// ContractIncentive gives incentives to self-executing contracts
func (k Keeper) ContractIncentive(ctx sdk.Context, maxIncentive sdk.Coin, contract sdk.AccAddress) (sdk.Coin, error) {
	p := k.GetParams(ctx)

	//incentivePool := k.bankKeeper.GetBalance(ctx, k.accountKeeper.GetModuleAddress("compute"), types.Denom).Amount

	if maxIncentive.Amount.Int64() > p.MaxContractIncentive {
		maxIncentive.Amount = sdk.NewInt(p.MaxContractIncentive)
	}
	fmt.Printf("maxIncentive %v\n", maxIncentive)

	incentiveMultiplier := sdk.NewDec(p.ContractIncentiveMul).QuoInt64(100)
	fmt.Printf("incentiveMultiplier %v\n", incentiveMultiplier)
	maxIncentiveDec := sdk.NewDecFromInt(maxIncentive.Amount).Mul(incentiveMultiplier)
	fmt.Printf("maxIncentiveDec %v\n", maxIncentiveDec)
	incentiveCoin := sdk.NewCoin(maxIncentive.Denom, maxIncentiveDec.TruncateInt())
	err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, contract, sdk.NewCoins(incentiveCoin))
	if err != nil {
		return sdk.Coin{}, err
	}

	k.Logger(ctx).Info("allocated", "contract", contract, "coins", maxIncentive)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDistributedToContract,
			sdk.NewAttribute(types.AttributeKeyAddress, contract.String()),
		),
	)
	return incentiveCoin, nil
}

/*
// SetIncentiveCoins distributes compute module allocated coins to selected contracts
func (k Keeper) SetIncentiveCoins(ctx sdk.Context, addressList []string) {
	p := k.GetParams(ctx)

	total := k.bankKeeper.GetBalance(ctx, k.accountKeeper.GetModuleAddress("compute"), types.Denom)
	k.Logger(ctx).Info("TC reward", "total", total)

	amount := total.Amount.QuoRaw(int64(len(addressList)))
	if amount.Int64() > p.MaxContractIncentive {
		amount = sdk.NewInt(p.MaxContractIncentive)
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
*/
