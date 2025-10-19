package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/runtime"
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
	if !fixedFeeCoin.IsNil() && fixedFeeCoin.Denom == types.Denom {
		toCommunityPool = totalFlowFees.Sub(fixedFeeCoin)

		err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, feeAddr, types.ModuleName, sdk.NewCoins(fixedFeeCoin))
		if err != nil {
			return sdk.Coin{}, errorsmod.Wrap(types.ErrUnexpectedFeeCalculation, "could not send coins to module: "+err.Error())
		}
		err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(fixedFeeCoin))
		if err != nil {
			return sdk.Coin{}, errorsmod.Wrap(types.ErrUnexpectedFeeCalculation, "could not burn coins: "+err.Error())
		}

		// Track total burnt coins
		k.addToTotalBurnt(ctx, fixedFeeCoin)

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

// addToTotalBurnt adds the given coin to the total burnt coins
func (k Keeper) addToTotalBurnt(ctx sdk.Context, coin sdk.Coin) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	fmt.Println("coin", coin)
	var totalBurnt sdk.Coin
	bz := store.Get(types.TotalBurntKey)
	if bz != nil {
		k.cdc.MustUnmarshal(bz, &totalBurnt)
	}

	if totalBurnt.Denom == "" {
		totalBurnt = coin
	} else {
		totalBurnt = totalBurnt.Add(coin)
	}

	store.Set(types.TotalBurntKey, k.cdc.MustMarshal(&totalBurnt))
}

func (k Keeper) GetTotalBurnt(ctx sdk.Context) sdk.Coin {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	var totalBurnt sdk.Coin
	bz := store.Get(types.TotalBurntKey)
	if bz != nil {
		k.cdc.MustUnmarshal(bz, &totalBurnt)
	}
	return totalBurnt
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

	// Reorder FeeLimit so "uinto" comes first if available
	gasFeeCoins := p.GasFeeCoins
	for i, coin := range gasFeeCoins {
		if coin.Denom == types.Denom {
			gasFeeCoins[0], gasFeeCoins[i] = gasFeeCoins[i], gasFeeCoins[0]
			break
		}
	}

	for _, coin := range gasFeeCoins {
		feeCoin, _, err := k.calculateFee(coin.Denom, expectedGas, flow, p)
		if err != nil {
			return nil, "", err
		}

		if feeBalances.AmountOf(feeCoin.Denom).GTE(feeCoin.Amount) {
			return feeAddr, feeCoin.Denom, nil
		}
	}

	// Fallback to owner balance if allowed
	if flow.Configuration != nil && flow.Configuration.WalletFallback {
		ownerAddr, err := sdk.AccAddressFromBech32(flow.Owner)
		if err != nil {
			return nil, "", err
		}

		ownerBalances := k.bankKeeper.GetAllBalances(ctx, ownerAddr)
		for _, coin := range gasFeeCoins {
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

func (k Keeper) SendFeesToTrustlessAgentFeeAdmin(ctx sdk.Context, flow types.Flow, trustlessAgent types.TrustlessAgent) error {
	if flow.TrustlessAgent.FeeLimit == nil {
		return nil
	}

	feeAddr, err := sdk.AccAddressFromBech32(flow.FeeAddress)
	if err != nil {
		return err
	}
	trustlessAgentAdminAddr, err := sdk.AccAddressFromBech32(trustlessAgent.FeeConfig.FeeAdmin)
	if err != nil {
		return err
	}
	supportedCoins := trustlessAgent.FeeConfig.FeeCoinsSupported

	// Reorder FeeLimit so "uinto" comes first if available
	feeLimits := flow.TrustlessAgent.FeeLimit
	for i, coin := range feeLimits {
		if coin.Denom == types.Denom {
			feeLimits[0], feeLimits[i] = feeLimits[i], feeLimits[0]
			break
		}
	}

	// Try each fee coin in order
	for _, feeLimit := range feeLimits {
		found, feeCoin := supportedCoins.Find(feeLimit.Denom)
		if supportedCoins.Len() != 0 && !found {
			continue // skip unsupported denom
		}
		if feeCoin.Amount.GT(feeLimit.Amount) {
			continue // skip if over the configured limit
		}

		// Check balance before attempting transfer
		balance := k.bankKeeper.GetBalance(ctx, feeAddr, feeCoin.Denom)
		if balance.IsLT(feeCoin) {
			// Try fallback address if enabled and insufficient balance
			if flow.Configuration.WalletFallback {
				fallbackAddr, err := sdk.AccAddressFromBech32(flow.Owner)
				if err == nil {
					fallbackBalance := k.bankKeeper.GetBalance(ctx, fallbackAddr, feeCoin.Denom)
					if fallbackBalance.IsGTE(feeCoin) {
						return k.bankKeeper.SendCoins(ctx, fallbackAddr, trustlessAgentAdminAddr, sdk.Coins{feeCoin})
					}
				}
			}
			continue // try next fee option
		}

		// If we get here, we have sufficient balance to make the transfer
		err = k.bankKeeper.SendCoins(ctx, feeAddr, trustlessAgentAdminAddr, sdk.Coins{feeCoin})
		if err == nil {
			return nil // success
		}
	}

	// No matching coins found with sufficient balance within limit
	return errorsmod.Wrap(types.ErrNotFound, "no valid fee coin with sufficient balance found within limit")
}
