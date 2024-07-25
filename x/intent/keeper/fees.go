package keeper

import (

	//"log"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trstlabs/intento/x/intent/types"
)

// DistributeCoins distributes Action fees and handles remaining action fee balance after last execution
func (k Keeper) DistributeCoins(ctx sdk.Context, action types.ActionInfo, feeAddr sdk.AccAddress, feeDenom string, isRecurring bool, proposer sdk.ConsAddress) (sdk.Coin, error) {
	p := k.GetParams(ctx)

	k.Logger(ctx).Debug("gas", "consumed", sdk.NewIntFromUint64(ctx.GasMeter().GasConsumed()))

	gasMultipleSmall := sdk.NewIntFromUint64(ctx.GasMeter().GasConsumed() * uint64(p.ActionFlexFeeMul))
	gasMultiple := gasMultipleSmall.Quo(math.NewInt(100))
	if !gasMultiple.IsPositive() {
		return sdk.Coin{}, types.ErrIntOverflowAction
	}
	found, coins := p.GasFeeCoins.Sort().Find(feeDenom)
	if !found {
		return sdk.Coin{}, errorsmod.Wrap(types.ErrNotFound, "coin not found")
	}
	gasFeeAmount := coins.Amount.Mul(gasMultiple)

	//depending if execution is recurring the constant fee may differ (gov param)

	// proposer reward
	// transfer collected fees to the distribution module account
	flexFeeCoin := sdk.NewCoin(feeDenom, gasFeeAmount)
	if flexFeeCoin.Amount.IsZero() {
		return sdk.Coin{}, errorsmod.Wrap(errorsmod.ErrPanic, "flexFeeCoin was zero, this should never happen")
	}

	totalActionFees := flexFeeCoin
	toCommunityPool := flexFeeCoin

	//pay out any remaining balance to the owner after deducting fee, commision and gas
	if p.ActionConstantFee != 0 {
		fixedFee := sdk.NewInt(p.ActionConstantFee * int64(len(action.Msgs)))
		fixedFeeCoin := sdk.NewCoin(feeDenom, fixedFee)
		totalActionFees = totalActionFees.Add(fixedFeeCoin)

		// if feeDenom == types.Denom {
		// 	k.bankKeeper.BurnCoins(ctx, "intent", sdk.NewCoins(fixedFeeCoin))
		// } else {
		toCommunityPool = toCommunityPool.Add(fixedFeeCoin)
		//}
	}

	if !isRecurring && action.Configuration != nil && !action.Configuration.FallbackToOwnerBalance {
		actionAddrBalance := k.bankKeeper.GetAllBalances(ctx, feeAddr)

		percentageActionFundsCommission := sdk.NewDecWithPrec(p.ActionFundsCommission, 2)
		amountActionFundsCommissionCoin := sdk.NewCoin(feeDenom, percentageActionFundsCommission.MulInt(actionAddrBalance.AmountOf(feeDenom)).Ceil().TruncateInt())
		totalActionFees = totalActionFees.Add(amountActionFundsCommissionCoin)

		toOwnerCoins, negative := actionAddrBalance.Sort().SafeSub(totalActionFees)

		if !negative {
			ownerAddr, err := sdk.AccAddressFromBech32(action.Owner)
			if err != nil {
				return sdk.Coin{}, err
			}
			err = k.bankKeeper.SendCoins(ctx, feeAddr, ownerAddr, toOwnerCoins)
			if err != nil {
				return sdk.Coin{}, err
			}

		}
	}

	err := k.distrKeeper.FundCommunityPool(ctx, sdk.NewCoins(toCommunityPool), feeAddr)
	if err != nil {
		return sdk.Coin{}, err
	}
	k.Logger(ctx).Debug("fee", "amount", flexFeeCoin.Amount, "to", proposer.String())

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
		if action.Configuration.FallbackToOwnerBalance {
			feeAddr, err = sdk.AccAddressFromBech32(action.Owner)
			if err != nil {
				return err
			}
			err = k.bankKeeper.SendCoins(ctx, feeAddr, hostedAddr, sdk.Coins{feeCoin})
			if err != nil {
				return err
			}
		} else {
			return err
		}

	}
	return nil

	//nice to have: ics20 transfer to destination (needed: channelID)
}

// CheckBalanceForGasFee checks if the address has enough balance to cover the gas fee.
func (k Keeper) GetFeeAccountForMinFees(ctx sdk.Context, action types.ActionInfo, expectedGas int64) (account sdk.AccAddress, denom string, err error) {
	p := k.GetParams(ctx)

	feeAddr, err := sdk.AccAddressFromBech32(action.FeeAddress)
	if err != nil {
		return nil, "", err
	}

	actionAddrBalances := k.bankKeeper.GetAllBalances(ctx, feeAddr).Sort()
	// Calculate the required fee
	minFee := sdk.NewCoins()
	for _, coin := range p.GasFeeCoins {
		amountSmall := coin.Amount.Mul(sdk.NewInt(p.ActionFlexFeeMul * expectedGas))
		amount := amountSmall.Quo(math.NewInt(100))
		minFee = minFee.Add(sdk.NewCoin(coin.Denom, amount))
	}
	denom = GetDenomIfAnyGTE(actionAddrBalances, minFee)
	// Check if the address balance has enough coins to cover the required fee
	if denom == "" {
		if action.Configuration != nil && action.Configuration.FallbackToOwnerBalance {
			ownerAddr, err := sdk.AccAddressFromBech32(action.Owner)
			if err != nil {
				return nil, "", err
			}

			actionAddrBalances = k.bankKeeper.GetAllBalances(ctx, ownerAddr).Sort()
			denom = GetDenomIfAnyGTE(actionAddrBalances, minFee)
			if denom == "" {
				return nil, "", err
			}
			feeAddr = ownerAddr

		}

	}

	return feeAddr, denom, nil
}

func GetDenomIfAnyGTE(coins sdk.Coins, coinsB sdk.Coins) string {
	if len(coinsB) == 0 {
		return ""
	}

	for _, coin := range coins {
		amt := coinsB.AmountOf(coin.Denom)
		if coin.Amount.GTE(amt) && !amt.IsZero() {
			return coin.Denom
		}
	}

	return ""
}
