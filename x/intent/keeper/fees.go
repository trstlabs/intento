package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trstlabs/intento/x/intent/types"
)

// calculateFee calculates the fee based on the given gas consumption and other fee parameters.
func (k Keeper) calculateFee(feeDenom string, gasConsumed uint64, flow types.Flow, p types.Params) (sdk.Coin, sdk.Coin, error) {
	// Calculate base gas fee based on the consumed gas
	gasSmall := math.NewIntFromUint64(gasConsumed * uint64(p.FlowFlexFeeMul))
	gasFeeAmount, err := gasSmall.SafeQuo(math.NewInt(1000))
	if err != nil || !gasFeeAmount.IsPositive() {
		return sdk.Coin{}, sdk.Coin{}, errorsmod.Wrap(types.ErrUnexpectedFeeCalculation, "invalid fee calculation")
	}

	// Find the gas fee denomination
	found, coins := p.GasFeeCoins.Sort().Find(feeDenom)
	if !found {
		return sdk.Coin{}, sdk.Coin{}, errorsmod.Wrap(types.ErrNotFound, "gas fee denom not supported")
	}

	// Calculate the gas fee coin amount
	gasFeeCoin := sdk.NewCoin(feeDenom, coins.Amount.Mul(gasFeeAmount))
	if gasFeeCoin.Amount.IsZero() {
		return sdk.Coin{}, sdk.Coin{}, errorsmod.Wrap(errorsmod.ErrPanic, "calculated fee is zero, this should never happen")
	}

	// Add fixed fees if applicable (e.g., burn fee)
	totalFlowFees := gasFeeCoin
	if p.BurnFeePerMsg != 0 && feeDenom == types.Denom {
		fixedFee := math.NewInt(p.BurnFeePerMsg * int64(len(flow.Msgs)))
		fixedFeeCoin := sdk.NewCoin(feeDenom, fixedFee)
		totalFlowFees = totalFlowFees.Add(fixedFeeCoin)
		return totalFlowFees, fixedFeeCoin, nil

	}

	return totalFlowFees, sdk.Coin{}, nil
}

// DistributeCoins distributes Flow fees and handles remaining flow fee balance after last execution
func (k Keeper) DistributeCoins(ctx sdk.Context, flow types.Flow, feeAddr sdk.AccAddress, feeDenom string) (sdk.Coin, error) {
	p, err := k.GetParams(ctx)
	if err != nil {
		return sdk.Coin{}, err
	}
	k.Logger(ctx).Debug("gas", "consumed", math.NewIntFromUint64(ctx.GasMeter().GasConsumed()), "flowID", flow.ID)

	// Calculate the fee
	totalFlowFees, fixedFeeCoin, err := k.calculateFee(feeDenom, ctx.GasMeter().GasConsumed(), flow, p)
	if err != nil {
		return sdk.Coin{}, err
	}

	toCommunityPool := totalFlowFees
	if !fixedFeeCoin.IsNil() {
		toCommunityPool = totalFlowFees.Sub(fixedFeeCoin)

		err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, feeAddr, types.ModuleName, sdk.NewCoins(fixedFeeCoin))
		if err != nil {
			return sdk.Coin{}, err
		}
		err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(fixedFeeCoin))
		if err != nil {
			return sdk.Coin{}, errorsmod.Wrap(errorsmod.ErrPanic, "could not burn coins, this should never happen")
		}
		k.Logger(ctx).Debug("flow fee burn", "amount", fixedFeeCoin)

		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypeFlowFeeBurn,
				sdk.NewAttribute(types.AttributeKeyFlowID, fmt.Sprintf("%d", flow.ID)),
				sdk.NewAttribute(types.AttributeKeyFlowFeeBurnAmount, fixedFeeCoin.Amount.String()),
			),
		})

	}

	// Distribute additional funds (e.g., flow commission)
	if flow.ExecTime.Equal(flow.EndTime) && feeAddr.String() != flow.Owner {
		flowAddrBalance := k.bankKeeper.GetAllBalances(ctx, feeAddr)
		percentageFlowFundsCommission := math.LegacyNewDecWithPrec(p.FlowFundsCommission, 2)
		amountFlowFundsCommissionCoin := sdk.NewCoin(feeDenom, percentageFlowFundsCommission.MulInt(flowAddrBalance.AmountOf(feeDenom)).Ceil().TruncateInt())

		totalFlowFees = totalFlowFees.Add(amountFlowFundsCommissionCoin)
		toCommunityPool = toCommunityPool.Add(amountFlowFundsCommissionCoin)

		// Ensure that toCommunityPool does not exceed available flowAddrBalance
		if flowAddrBalance.IsAllGTE(sdk.Coins{toCommunityPool}) {
			toOwnerCoins, negative := flowAddrBalance.Sort().SafeSub(toCommunityPool)
			if !negative {
				// Continue with your normal logic here
				ownerAddr, err := sdk.AccAddressFromBech32(flow.Owner)
				if err != nil {
					return sdk.Coin{}, err
				}
				err = k.bankKeeper.SendCoins(ctx, feeAddr, ownerAddr, toOwnerCoins)
				if err != nil {
					return sdk.Coin{}, err
				}
			} else {
				return sdk.Coin{}, errorsmod.Wrap(types.ErrUnexpectedFeeCalculation, "fees exceed available balance")
			}
		} else {
			// If community pool fees exceed available balance
			return sdk.Coin{}, errorsmod.Wrap(types.ErrUnexpectedFeeCalculation, "total flow fees exceed available balance")
		}
	}

	err = k.distrKeeper.FundCommunityPool(ctx, sdk.NewCoins(toCommunityPool), feeAddr)
	if err != nil {
		return sdk.Coin{}, err
	}

	k.Logger(ctx).Debug("fee", "amount", totalFlowFees)

	return totalFlowFees, nil
}

// GetFeeAccountForMinFees checks if the flow fee address (or optionally the owner address)
// has enough balance to cover the minimum gas fee for any configured gas fee denom.
func (k Keeper) GetFeeAccountForMinFees(ctx sdk.Context, flow types.Flow, expectedGas uint64) (sdk.AccAddress, string, error) {
	p, err := k.GetParams(ctx)
	if err != nil {
		return nil, "", err
	}

	feeAddr, err := sdk.AccAddressFromBech32(flow.FeeAddress)
	if err != nil {
		return nil, "", err
	}

	feeBalances := k.bankKeeper.GetAllBalances(ctx, feeAddr)

	for _, coin := range p.GasFeeCoins {
		feeCoin, _, err := k.calculateFee(coin.Denom, expectedGas, flow, p)
		if err != nil {
			return nil, "", err
		}

		if feeBalances.AmountOf(feeCoin.Denom).GTE(feeCoin.Amount) {
			return feeAddr, feeCoin.Denom, nil
		}
	}

	// Fallback to owner balance if allowed
	if flow.Configuration != nil && flow.Configuration.FallbackToOwnerBalance {
		ownerAddr, err := sdk.AccAddressFromBech32(flow.Owner)
		if err != nil {
			return nil, "", err
		}

		ownerBalances := k.bankKeeper.GetAllBalances(ctx, ownerAddr)
		for _, coin := range p.GasFeeCoins {
			feeCoin, _, err := k.calculateFee(coin.Denom, expectedGas, flow, p)
			if err != nil {
				return nil, "", err
			}

			if ownerBalances.AmountOf(feeCoin.Denom).GTE(feeCoin.Amount) {
				return ownerAddr, feeCoin.Denom, nil
			}
		}
	}

	return nil, "", nil
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

func (k Keeper) SendFeesToHostedAdmin(ctx sdk.Context, flow types.Flow, trustlessAgent types.TrustlessAgent) error {
	// Parse addresses
	feeAddr, err := sdk.AccAddressFromBech32(flow.FeeAddress)
	if err != nil {
		return err
	}
	hostedAccAdminAddr, err := sdk.AccAddressFromBech32(trustlessAgent.FeeConfig.FeeAdmin)
	if err != nil {
		return err
	}

	supportedCoins := trustlessAgent.FeeConfig.FeeCoinsSupported.Sort()

	// Find the cheapest valid matching coin (first one that matches)
	for _, feeLimit := range flow.TrustlessAgent.FeeLimit {
		found, feeCoin := supportedCoins.Find(feeLimit.Denom)
		if !found {
			continue // skip unsupported denom
		}
		if feeCoin.Amount.GT(feeLimit.Amount) {
			continue // skip if over the configured limit
		}

		// Try to send this coin
		err = k.bankKeeper.SendCoins(ctx, feeAddr, hostedAccAdminAddr, sdk.Coins{feeCoin})
		if err == nil {
			return nil // success
		}

		// Try fallback to owner if enabled
		if flow.Configuration.FallbackToOwnerBalance {
			fallbackAddr, err := sdk.AccAddressFromBech32(flow.Owner)
			if err != nil {
				return err
			}
			return k.bankKeeper.SendCoins(ctx, fallbackAddr, hostedAccAdminAddr, sdk.Coins{feeCoin})
		}

		return err // no fallback or fallback failed
	}

	// No matching coins found within limit
	return errorsmod.Wrap(types.ErrNotFound, "no valid fee coin matched within limit")
}
