package keeper

import (

	//"log"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/trstlabs/intento/x/intent/types"
)

// DistributeCoins distributes Action fees and handles remaining action fee balance after last execution
func (k Keeper) DistributeCoins(ctx sdk.Context, action types.ActionInfo, flexFee sdkmath.Int, isRecurring bool, proposer sdk.ConsAddress) (sdk.Coin, error) {
	cacheCtx, writeCache := ctx.CacheContext()
	p := k.GetParams(ctx)

	flexFeeMultiplier := sdk.NewDec(p.ActionFlexFeeMul).QuoInt64(100)

	flexFeeMulDec := sdk.NewDecFromInt(flexFee).Mul(flexFeeMultiplier)

	// //round flex to one int if it is smaller than one
	// if flexFeeMulDec.TruncateInt().IsZero() {
	// 	flexFeeMulDec = flexFeeMulDec.Ceil()
	// }
	feeAddr, err := sdk.AccAddressFromBech32(action.FeeAddress)
	if err != nil {
		return sdk.Coin{}, err
	}
	ownerAddr, err := sdk.AccAddressFromBech32(action.Owner)
	if err != nil {
		return sdk.Coin{}, err
	}

	actionAddrBalance := k.bankKeeper.GetAllBalances(ctx, feeAddr)

	//depending if execution is recurring the constant fee may differ (gov param)
	fixedFee := sdk.NewInt(p.ActionConstantFee * int64(len(action.Msgs)))
	if isRecurring {
		fixedFee = sdk.NewInt(p.RecurringActionConstantFee * int64(len(action.Msgs)))
	}

	fixedFeeCommunityCoin := sdk.NewCoin(types.Denom, fixedFee)

	//if last execution, return remaining balance minus commision
	if !isRecurring && !actionAddrBalance.Empty() {
		percentageActionFundsCommission := sdk.NewDecWithPrec(p.ActionFundsCommission, 2)
		amountActionFundsCommissionCoin := sdk.NewCoin(types.Denom, percentageActionFundsCommission.MulInt(actionAddrBalance.AmountOf(types.Denom)).Ceil().TruncateInt())
		fixedFeeCommunityCoin = fixedFeeCommunityCoin.Add(amountActionFundsCommissionCoin)
	}
	fixedFeeCommunityCoins := sdk.NewCoins(fixedFeeCommunityCoin)
	totalActionFees := fixedFeeCommunityCoin.Add(sdk.NewCoin(types.Denom, flexFeeMulDec.Ceil().TruncateInt()))

	// proposer reward
	// transfer collected fees to the distribution module account
	flexFeeCoin := sdk.NewCoin(types.Denom, flexFeeMulDec.Ceil().TruncateInt())
	if flexFeeCoin.Amount.IsZero() {
		return sdk.Coin{}, errorsmod.Wrap(errorsmod.ErrPanic, "flexFeeCoin was zero, this should never happen")
	}

	proposerAddr := k.stakingKeeper.ValidatorByConsAddr(cacheCtx, proposer)

	k.distrKeeper.AllocateTokensToValidator(cacheCtx, proposerAddr, sdk.NewDecCoinsFromCoins(flexFeeCoin))

	//the trigger account should be funded with the fee amount
	err = k.bankKeeper.SendCoinsFromAccountToModule(cacheCtx, feeAddr, authtypes.FeeCollectorName, sdk.NewCoins(flexFeeCoin))
	if err != nil {
		if action.Configuration.FallbackToOwnerBalance {
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
		toOwnerCoins, negative := actionAddrBalance.Sort().SafeSub(totalActionFees)

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

	return totalActionFees, nil
}

func (k Keeper) SendFeesToHosted(ctx sdk.Context, action types.ActionInfo, hostedAccount types.HostedAccount) error {
	feeAddr, err := sdk.AccAddressFromBech32(action.FeeAddress)
	if err != nil {
		return err
	}

	hostedAddr, err := sdk.AccAddressFromBech32(hostedAccount.HostedAddress)
	if err != nil {
		return err
	}
	found, feeCoin := hostedAccount.HostFeeConfig.FeeCoinsSuported.Sort().Find(action.HostedConfig.FeeCoinLimit.Denom)
	if !found {
		return errorsmod.Wrap(types.ErrNotFound, "coin not in hosted config")
	}

	if feeCoin.Amount.GT(action.HostedConfig.FeeCoinLimit.Amount) {
		return types.ErrHostedFeeLimit
	}

	err = k.bankKeeper.SendCoins(ctx, feeAddr, hostedAddr, sdk.Coins{feeCoin})
	if err != nil {
		return err
	}
	return nil

	//nice to have: ics20 transfer to destination (needed: channelID)
}
