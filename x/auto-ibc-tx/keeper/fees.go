package keeper

import (

	//"log"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/trstlabs/trst/x/auto-ibc-tx/types"
)

// DistributeCoins distributes AutoTx fees and handles remaining autoTx fee balance after last execution
func (k Keeper) DistributeCoins(ctx sdk.Context, autoTxInfo types.AutoTxInfo, flexFee sdkmath.Int, isRecurring bool, proposer sdk.ConsAddress) (sdk.Coin, error) {
	cacheCtx, writeCache := ctx.CacheContext()
	p := k.GetParams(ctx)

	flexFeeMultiplier := sdk.NewDec(p.AutoTxFlexFeeMul).QuoInt64(100)

	flexFeeMulDec := sdk.NewDecFromInt(flexFee).Mul(flexFeeMultiplier)

	// //round flex to one int if it is smaller than one
	// if flexFeeMulDec.TruncateInt().IsZero() {
	// 	flexFeeMulDec = flexFeeMulDec.Ceil()
	// }
	feeAddr, err := sdk.AccAddressFromBech32(autoTxInfo.FeeAddress)
	if err != nil {
		return sdk.Coin{}, err
	}
	ownerAddr, err := sdk.AccAddressFromBech32(autoTxInfo.Owner)
	if err != nil {
		return sdk.Coin{}, err
	}

	autoTxAddrBalance := k.bankKeeper.GetAllBalances(ctx, feeAddr)

	//depending if execution is recurring the constant fee may differ (gov param)
	fixedFee := sdk.NewInt(p.AutoTxConstantFee * int64(len(autoTxInfo.Msgs)))
	if isRecurring {
		fixedFee = sdk.NewInt(p.RecurringAutoTxConstantFee * int64(len(autoTxInfo.Msgs)))
	}

	fixedFeeCommunityCoin := sdk.NewCoin(types.Denom, fixedFee)

	//if last execution, return remaining balance minus commision
	if !isRecurring && !autoTxAddrBalance.Empty() {
		percentageAutoTxFundsCommission := sdk.NewDecWithPrec(p.AutoTxFundsCommission, 2)
		amountAutoTxFundsCommissionCoin := sdk.NewCoin(types.Denom, percentageAutoTxFundsCommission.MulInt(autoTxAddrBalance.AmountOf(types.Denom)).Ceil().TruncateInt())
		fixedFeeCommunityCoin = fixedFeeCommunityCoin.Add(amountAutoTxFundsCommissionCoin)
	}
	fixedFeeCommunityCoins := sdk.NewCoins(fixedFeeCommunityCoin)
	totalAutoTxFees := fixedFeeCommunityCoin.Add(sdk.NewCoin(types.Denom, flexFeeMulDec.Ceil().TruncateInt()))

	// proposer reward
	// transfer collected fees to the distribution module account
	flexFeeCoin := sdk.NewCoin(types.Denom, flexFeeMulDec.Ceil().TruncateInt())
	if flexFeeCoin.Amount.IsZero() {
		return sdk.Coin{}, errorsmod.Wrap(sdkerrors.ErrInsufficientFee, "flexFeeCoin was zero")
	}

	proposerAddr := k.stakingKeeper.ValidatorByConsAddr(cacheCtx, proposer)

	k.distrKeeper.AllocateTokensToValidator(cacheCtx, proposerAddr, sdk.NewDecCoinsFromCoins(flexFeeCoin))

	//the trigger account should be funded with the fee amount
	err = k.bankKeeper.SendCoinsFromAccountToModule(cacheCtx, feeAddr, authtypes.FeeCollectorName, sdk.NewCoins(flexFeeCoin))
	if err != nil {
		if autoTxInfo.Configuration.FallbackToOwnerBalance {
			err := k.bankKeeper.SendCoinsFromAccountToModule(cacheCtx, ownerAddr, authtypes.FeeCollectorName, sdk.NewCoins(flexFeeCoin))
			if err != nil {
				return sdk.Coin{}, err
			}

			err = k.distrKeeper.FundCommunityPool(cacheCtx, fixedFeeCommunityCoins, ownerAddr)
			if err != nil {
				return sdk.Coin{}, err
			}
		} else {
			return sdk.Coin{}, err
		}
	} else {
		err = k.distrKeeper.FundCommunityPool(cacheCtx, fixedFeeCommunityCoins, feeAddr)
		if err != nil {
			return sdk.Coin{}, err
		}
	}
	//pay out any remaining balance to the owner after deducting fee, commision and gas
	if !isRecurring {
		toOwnerCoins, negative := autoTxAddrBalance.Sort().SafeSub(totalAutoTxFees)

		if !negative {
			err := k.bankKeeper.SendCoins(cacheCtx, feeAddr, ownerAddr, toOwnerCoins)
			if err != nil {
				return sdk.Coin{}, err
			}

		}
	}

	//we only write to the state when we know the send actions succeed
	writeCache()

	k.Logger(ctx).Debug("flex_fee", "amount", flexFeeCoin.Amount, "to", proposer.String())

	return totalAutoTxFees, nil
}
