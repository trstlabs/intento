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
func (k Keeper) DistributeCoins(ctx sdk.Context, contract types.ContractInfoWithAddress, gasUsed uint64, isRecurring bool, isLastExec bool, proposer sdk.ConsAddress) (sdk.Coin, error) {

	p := k.GetParams(ctx)

	flexFeeMultiplier := sdk.NewDec(p.AutoMsgFlexFeeMul).QuoInt64(100)
	flexFee := sdk.NewDecFromInt(sdk.NewInt(int64(gasUsed))).Mul(flexFeeMultiplier)

	//direct a commission of the utrst contract balance towards the community pool
	contractBalance := k.bankKeeper.GetAllBalances(ctx, contract.Address)

	//depending on if self-execution is recurring the constant fee may differ (gov param)
	constantFee := sdk.NewInt(p.AutoMsgConstantFee)
	if isRecurring {
		constantFee = sdk.NewInt(p.RecurringAutoMsgConstantFee)
	}
	communityCoins := sdk.NewCoins(sdk.NewCoin(types.Denom, constantFee))
	if isLastExec {
		percentageAutoMsgFundsCommission := sdk.NewDecWithPrec(p.AutoMsgFundsCommission, 2)
		amountAutoMsgFundsCommissionCoin := sdk.NewCoin(types.Denom, percentageAutoMsgFundsCommission.MulInt(contractBalance.AmountOf(types.Denom)).Ceil().TruncateInt())
		communityCoins = communityCoins.Add(amountAutoMsgFundsCommissionCoin)
	}

	totalExecCost := communityCoins.Add(sdk.NewCoin(types.Denom, flexFee.TruncateInt()))

	if contract.Duration >= p.MinContractDurationForIncentive {
		fmt.Printf("contr bal%s\n", contractBalance)
		incentive, err := k.ContractIncentive(ctx, totalExecCost[0], contract.Address)
		if err != nil {
			return sdk.Coin{}, err
		}
		contractBalance = contractBalance.Add(incentive)
	}

	if isLastExec {
		//pay out the remaining balance to the contract owner after deducting fee, commision and gas cost
		toOwnerCoins, negative := contractBalance.Sort().SafeSub(totalExecCost)
		//fmt.Printf("toOwnerCoins %v\n", toOwnerCoins)
		if !negative {
			err := k.bankKeeper.SendCoins(ctx, contract.Address, contract.ContractInfo.Owner, toOwnerCoins)
			if err != nil {
				return sdk.Coin{}, err
			}

		}
	}

	// proposer reward
	// transfer collected fees to the distribution module account
	flexFeeCoin := sdk.NewCoin(types.Denom, flexFee.TruncateInt())
	if flexFeeCoin.Amount.IsZero() {
		return sdk.Coin{}, sdkerrors.Wrap(sdkerrors.ErrInsufficientFee, "flexFeeCoin was zero")
	}
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, contract.Address, authtypes.FeeCollectorName, sdk.NewCoins(flexFeeCoin))
	if err != nil {
		return sdk.Coin{}, err
	}
	proposerAddr := k.stakingKeeper.ValidatorByConsAddr(ctx, proposer)
	fmt.Printf("-allocating flexFeeCoin :%s \n", flexFeeCoin.Amount)
	fmt.Printf("-asdfproposer :%s \n", proposer.String())
	k.Logger(ctx).Info("allocating", "flexFeeCoin", flexFeeCoin.Amount, "proposer", proposer.String())
	k.distrKeeper.AllocateTokensToValidator(ctx, proposerAddr, sdk.NewDecCoinsFromCoins(flexFeeCoin))

	//the contract should be funded with the fee. Iif the contract is not able to pay, the contract owner pays next in line
	err = k.distrKeeper.FundCommunityPool(ctx, communityCoins, contract.Address)
	if err != nil {
		store := ctx.KVStore(k.storeKey)
		// unless a contract instantiated the contract, we deduct fees so execution can be written
		if !store.Has(types.GetContractEnclaveKey(contract.Creator)) {
			err := k.distrKeeper.FundCommunityPool(ctx, communityCoins, contract.ContractInfo.Owner)
			if err != nil {
				return sdk.Coin{}, err
			}
		}
		return sdk.Coin{}, err
	}

	return totalExecCost[0], nil
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
