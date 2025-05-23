package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trstlabs/intento/x/intent/types"
)

// calculateFee calculates the fee based on the given gas consumption and other fee parameters.
func (k Keeper) calculateFee(feeDenom string, gasConsumed uint64, flow types.FlowInfo, p types.Params) (sdk.Coin, sdk.Coin, error) {
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
func (k Keeper) DistributeCoins(ctx sdk.Context, flow types.FlowInfo, feeAddr sdk.AccAddress, feeDenom string) (sdk.Coin, error) {
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
func (k Keeper) GetFeeAccountForMinFees(ctx sdk.Context, flow types.FlowInfo, expectedGas uint64) (sdk.AccAddress, string, error) {
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

func (k Keeper) SendFeesToHostedAdmin(ctx sdk.Context, flow types.FlowInfo, hostedAccount types.HostedAccount) error {
	feeAddr, err := sdk.AccAddressFromBech32(flow.FeeAddress)
	if err != nil {
		return err
	}

	hostedAccAdminAddr, err := sdk.AccAddressFromBech32(hostedAccount.HostFeeConfig.Admin)
	if err != nil {
		return err
	}
	found, feeCoin := hostedAccount.HostFeeConfig.FeeCoinsSuported.Sort().Find(flow.HostedICAConfig.FeeCoinLimit.Denom)
	if !found {
		return errorsmod.Wrap(types.ErrNotFound, "coin not in hosted config")
	}

	if feeCoin.Amount.GT(flow.HostedICAConfig.FeeCoinLimit.Amount) {
		return types.ErrHostedFeeLimit
	}

	err = k.bankKeeper.SendCoins(ctx, feeAddr, hostedAccAdminAddr, sdk.Coins{feeCoin})
	if err != nil {
		if flow.Configuration.FallbackToOwnerBalance {
			feeAddr, err = sdk.AccAddressFromBech32(flow.Owner)
			if err != nil {
				return err
			}
			err = k.bankKeeper.SendCoins(ctx, feeAddr, hostedAccAdminAddr, sdk.Coins{feeCoin})
			if err != nil {
				return err
			}
		} else {
			return err
		}

	}
	return nil
}
