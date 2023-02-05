package keeper

import (

	//"log"

	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/trstlabs/trst/x/auto-ibc-tx/types"
)

// DistributeCoins distributes AutoTx fees and handles remaining autoTx fee balance after last execution
func (k Keeper) DistributeCoins(ctx sdk.Context, autoTxInfo types.AutoTxInfo, flexFee sdk.Int, isRecurring bool, proposer sdk.ConsAddress) (sdk.Coin, error) {

	p := k.GetParams(ctx)
	fmt.Printf(" flexFee: %v \n", flexFee)
	flexFeeMultiplier := sdk.NewDec(p.AutoTxFlexFeeMul).QuoInt64(100)
	fmt.Printf(" flexFeeMul: %v \n", flexFeeMultiplier)
	flexFeeMulDec := sdk.NewDecFromInt(flexFee).Mul(flexFeeMultiplier)

	//direct a commission of the utrst autoTxInfo balance towards the community pool
	autoTxInfoBalance := k.bankKeeper.GetAllBalances(ctx, autoTxInfo.Address)

	//depending on if self-execution is recurring the constant fee may differ (gov param)
	constantFee := sdk.NewInt(p.AutoTxConstantFee)
	if isRecurring {
		constantFee = sdk.NewInt(p.RecurringAutoTxConstantFee)
	}
	communityCoins := sdk.NewCoins(sdk.NewCoin(types.Denom, constantFee))

	if !isRecurring && !autoTxInfoBalance.Empty() {
		percentageAutoTxFundsCommission := sdk.NewDecWithPrec(p.AutoTxFundsCommission, 2)
		amountAutoTxFundsCommissionCoin := sdk.NewCoin(types.Denom, percentageAutoTxFundsCommission.MulInt(autoTxInfoBalance.AmountOf(types.Denom)).Ceil().TruncateInt())
		communityCoins = communityCoins.Add(amountAutoTxFundsCommissionCoin)
	}

	totalAutoTxFees := communityCoins.Add(sdk.NewCoin(types.Denom, flexFeeMulDec.TruncateInt()))
	fmt.Printf(" totalAutoTxFees: %v \n", totalAutoTxFees)
	if !isRecurring {
		//pay out the remaining balance to the autoTxInfo owner after deducting fee, commision and gas cost
		toOwnerCoins, negative := autoTxInfoBalance.Sort().SafeSub(totalAutoTxFees)
		//fmt.Printf("toOwnerCoins %v\n", toOwnerCoins)
		if !negative {
			err := k.bankKeeper.SendCoins(ctx, autoTxInfo.Address, autoTxInfo.Owner, toOwnerCoins)
			if err != nil {
				return sdk.Coin{}, err
			}

		}
	}

	// proposer reward
	// transfer collected fees to the distribution module account
	flexFeeCoin := sdk.NewCoin(types.Denom, flexFeeMulDec.TruncateInt())
	if flexFeeCoin.Amount.IsZero() {
		return sdk.Coin{}, sdkerrors.Wrap(sdkerrors.ErrInsufficientFee, "flexFeeCoin was zero")
	}

	proposerAddr := k.stakingKeeper.ValidatorByConsAddr(ctx, proposer)
	//fmt.Printf("allocating flexFeeCoin :%s \n", flexFeeCoin.Amount)
	//fmt.Printf("proposer :%s \n", proposer.String())

	k.Logger(ctx).Debug("auto_tx_flex_fee", "flexFeeCoin", flexFeeCoin.Amount, "to_proposer", proposer.String())

	k.distrKeeper.AllocateTokensToValidator(ctx, proposerAddr, sdk.NewDecCoinsFromCoins(flexFeeCoin))

	//the autoTxInfo should be funded with the fee. Iif the autoTxInfo is not able to pay, the autoTxInfo owner pays next in line
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, autoTxInfo.Address, authtypes.FeeCollectorName, sdk.NewCoins(flexFeeCoin))
	if err != nil {
		err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, autoTxInfo.Owner, authtypes.FeeCollectorName, sdk.NewCoins(flexFeeCoin))
		if err != nil {
			return sdk.Coin{}, err
		}
		err = k.distrKeeper.FundCommunityPool(ctx, communityCoins, autoTxInfo.Owner)
		if err != nil {
			return sdk.Coin{}, err
		}

	} else {
		err = k.distrKeeper.FundCommunityPool(ctx, communityCoins, autoTxInfo.Address)
		return sdk.Coin{}, err
	}

	return totalAutoTxFees[0], nil
}

/*
// DistributeCoins distributes AutoTx fees and handles remaining autoTx fee balance
func (k Keeper) DistributeCoins(ctx sdk.Context, autoTxInfo types.AutoTxInfo, flexFee uint64, isRecurring bool, proposer sdk.ConsAddress) (sdk.Coin, error) {
	p := k.GetParams(ctx)

	flexFeeMultiplier := sdk.NewDec(p.AutoTxFlexFeeMul).QuoInt64(100)
	flexFeeMul := sdk.NewDecFromInt(sdk.NewInt(int64(flexFee))).Mul(flexFeeMultiplier)

	//direct a commission of the utrst autoTxInfo balance towards the community pool
	autoTxInfoBalance := k.bankKeeper.GetAllBalances(ctx, autoTxInfo.Address)

	//depending on if self-execution is recurring the constant fee may differ (gov param)
	constantFee := sdk.NewInt(p.AutoTxConstantFee)
	if isRecurring {
		constantFee = sdk.NewInt(p.RecurringAutoTxConstantFee)
	}
	communityCoins := sdk.NewCoins(sdk.NewCoin(types.Denom, constantFee))

	if !isRecurring && !autoTxInfoBalance.Empty() {
		percentageAutoTxFundsCommission := sdk.NewDecWithPrec(p.AutoTxFundsCommission, 2)
		amountAutoTxFundsCommissionCoin := sdk.NewCoin(types.Denom, percentageAutoTxFundsCommission.MulInt(autoTxInfoBalance.AmountOf(types.Denom)).Ceil().TruncateInt())
		communityCoins = communityCoins.Add(amountAutoTxFundsCommissionCoin)
	}

	totalAutoTxFees := communityCoins.Add(sdk.NewCoin(types.Denom, flexFeeMul.TruncateInt()))

	if !isRecurring {
		//pay out the remaining balance to the autoTxInfo owner after deducting fee, commision and gas cost
		toOwnerCoins, negative := autoTxInfoBalance.Sort().SafeSub(totalAutoTxFees)

		if !negative {
			err := k.bankKeeper.SendCoins(ctx, autoTxInfo.Address, autoTxInfo.Owner, toOwnerCoins)
			if err != nil {
				return sdk.Coin{}, err
			}
		}
	}

	// proposer reward
	// transfer collected fees to the distribution module account
	flexFeeCoin := sdk.NewCoin(types.Denom, flexFeeMul.TruncateInt())
	if flexFeeCoin.Amount.IsZero() {
		return sdk.Coin{}, sdkerrors.Wrap(sdkerrors.ErrInsufficientFee, "flexFeeCoin was zero")
	}

	proposerAddr := k.stakingKeeper.ValidatorByConsAddr(ctx, proposer)
	k.Logger(ctx).Debug("auto_tx_flex_fee", "flexFeeCoin", flexFeeCoin.Amount, "to_proposer", proposer.String())

	k.distrKeeper.AllocateTokensToValidator(ctx, proposerAddr, sdk.NewDecCoinsFromCoins(flexFeeCoin))

	//the autoTxInfo should be funded with the fee. Iif the autoTxInfo is not able to pay, the autoTxInfo owner pays next in line
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, autoTxInfo.Address, authtypes.FeeCollectorName, sdk.NewCoins(flexFeeCoin))
	if err != nil {
		err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, autoTxInfo.Owner, authtypes.FeeCollectorName, sdk.NewCoins(flexFeeCoin))
		if err != nil {
			return sdk.Coin{}, err
		}
		err = k.distrKeeper.FundCommunityPool(ctx, communityCoins, autoTxInfo.Owner)
		if err != nil {
			return sdk.Coin{}, err
		}

	} else {
		err = k.distrKeeper.FundCommunityPool(ctx, communityCoins, autoTxInfo.Address)
		return sdk.Coin{}, err
	}

	return totalAutoTxFees[0], nil
}
*/
