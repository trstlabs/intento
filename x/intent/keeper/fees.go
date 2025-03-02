package keeper

import (

	//"log"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trstlabs/intento/x/intent/types"
)

// DistributeCoins distributes Flow fees and handles remaining flow fee balance after last execution
func (k Keeper) DistributeCoins(ctx sdk.Context, flow types.FlowInfo, feeAddr sdk.AccAddress, feeDenom string) (sdk.Coin, error) {
	p, err := k.GetParams(ctx)
	if err != nil {
		panic(err)
	}
	k.Logger(ctx).Debug("gas", "consumed", math.NewIntFromUint64(ctx.GasMeter().GasConsumed()), "flowID", flow.ID)

	gasSmall := math.NewIntFromUint64(ctx.GasMeter().GasConsumed() * uint64(p.FlowFlexFeeMul))

	gas, err := gasSmall.SafeQuo(math.NewInt(1000))
	if err != nil {
		return sdk.Coin{}, errorsmod.Wrap(types.ErrUnexpectedFeeCalculation, "could not substract fee")
	}

	if !gas.IsPositive() {
		return sdk.Coin{}, types.ErrIntOverflowFlow
	}
	found, coins := p.GasFeeCoins.Sort().Find(feeDenom)
	if !found {
		return sdk.Coin{}, errorsmod.Wrap(types.ErrNotFound, "gas fee denom not supported")
	}
	gasFeeAmount := coins.Amount.Mul(gas)

	flexFeeCoin := sdk.NewCoin(feeDenom, gasFeeAmount)
	if flexFeeCoin.Amount.IsZero() {
		return sdk.Coin{}, errorsmod.Wrap(errorsmod.ErrPanic, "flexFeeCoin was zero, this should never happen")
	}

	totalFlowFees := flexFeeCoin
	toCommunityPool := flexFeeCoin

	if p.BurnFeePerMsg != 0 && feeDenom == types.Denom {
		fixedFee := math.NewInt(p.BurnFeePerMsg * int64(len(flow.Msgs)))
		fixedFeeCoin := sdk.NewCoin(feeDenom, fixedFee)
		totalFlowFees = totalFlowFees.Add(fixedFeeCoin)

		if feeDenom == types.Denom {
			err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, feeAddr, types.ModuleName, sdk.NewCoins(fixedFeeCoin))
			if err != nil {
				return sdk.Coin{}, err
			}

			err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(fixedFeeCoin))
			if err != nil {
				return sdk.Coin{}, errorsmod.Wrap(errorsmod.ErrPanic, "could not burn coins, this should never happen")
			}
		}
	}
	if flow.ExecTime.Equal(flow.EndTime) {
		if feeAddr.String() != flow.Owner {
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
	}

	err = k.distrKeeper.FundCommunityPool(ctx, sdk.NewCoins(toCommunityPool), feeAddr)
	if err != nil {
		return sdk.Coin{}, err
	}

	k.Logger(ctx).Debug("fee", "amount", totalFlowFees)

	return totalFlowFees, nil
}

func (k Keeper) SendFeesToHosted(ctx sdk.Context, flow types.FlowInfo, hostedAccount types.HostedAccount) error {
	feeAddr, err := sdk.AccAddressFromBech32(flow.FeeAddress)
	if err != nil {
		return err
	}

	hostedAddr, err := sdk.AccAddressFromBech32(hostedAccount.HostedAddress)
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

	err = k.bankKeeper.SendCoins(ctx, feeAddr, hostedAddr, sdk.Coins{feeCoin})
	if err != nil {
		if flow.Configuration.FallbackToOwnerBalance {
			feeAddr, err = sdk.AccAddressFromBech32(flow.Owner)
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
func (k Keeper) GetFeeAccountForMinFees(ctx sdk.Context, flow types.FlowInfo, expectedGas uint64) (account sdk.AccAddress, denom string, err error) {
	p, err := k.GetParams(ctx)
	if err != nil {
		panic(err)
	}

	feeAddr, err := sdk.AccAddressFromBech32(flow.FeeAddress)
	if err != nil {
		return nil, "", err
	}

	flowAddrBalances := k.bankKeeper.GetAllBalances(ctx, feeAddr).Sort()
	// Calculate the required fee
	minFee := sdk.NewCoins()
	for _, coin := range p.GasFeeCoins {
		amountSmall := coin.Amount.Mul(math.NewInt(p.FlowFlexFeeMul * int64(expectedGas)))
		amount, err := amountSmall.SafeQuo(math.NewInt(1000))
		if err != nil {
			return nil, "", errorsmod.Wrap(types.ErrIntOverflowFlow, "could not calculated expected fee")
		}
		minFee = minFee.Add(sdk.NewCoin(coin.Denom, amount))
	}
	denom = GetDenomIfAnyGTE(flowAddrBalances, minFee)
	// Check if the address balance has enough coins to cover the required fee
	if denom == "" {
		if flow.Configuration != nil && flow.Configuration.FallbackToOwnerBalance {
			ownerAddr, err := sdk.AccAddressFromBech32(flow.Owner)
			if err != nil {
				return nil, "", err
			}

			flowAddrBalances = k.bankKeeper.GetAllBalances(ctx, ownerAddr).Sort()
			if flowAddrBalances.IsZero() {
				return nil, "", errorsmod.Wrap(types.ErrNotFound, "flow owner bank balance is zero")
			}
			denom = GetDenomIfAnyGTE(flowAddrBalances, minFee)
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
