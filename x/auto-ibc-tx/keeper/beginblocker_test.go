package keeper_test

import (
	//"fmt"
	"fmt"

	"time"

	sdkmath "cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trstlabs/trst/x/auto-ibc-tx/keeper"
	"github.com/trstlabs/trst/x/auto-ibc-tx/types"
)

// BeginBlocker called every block, processes auto execution
func FakeBeginBlocker(ctx sdk.Context, k keeper.Keeper, fakeProposer sdk.ConsAddress) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	available := k.GetRelayerRewardsAvailability(ctx)
	if !available {
		k.SetRelayerRewardsAvailability(ctx, true)
	}

	autoTxs := k.GetAutoTxsForBlock(ctx)

	timeOfBlock := ctx.BlockHeader().Time
	for _, autoTx := range autoTxs {
		autoTx = k.GetAutoTxInfo(ctx, autoTx.TxID)
		autoTxHistory, _ := k.TryGetAutoTxHistory(ctx, autoTx.TxID)
		if !k.AllowedToExecute(ctx, autoTx) {
			k.AddAutoTxHistory(ctx, &autoTxHistory, &autoTx, timeOfBlock, sdk.Coin{}, false, nil, types.ErrAutoTxConditions)
			autoTx.ExecTime = autoTx.ExecTime.Add(autoTx.Interval)
			k.SetAutoTxInfo(ctx, &autoTx)
		}
		isRecurring := autoTx.ExecTime.Before(autoTx.EndTime)

		flexFee := calculateTimeBasedFlexFee(autoTx, autoTxHistory)
		fee, err := k.DistributeCoins(ctx, autoTx, flexFee, isRecurring, fakeProposer)

		k.RemoveFromAutoTxQueue(ctx, autoTx)
		if err != nil {
			fmt.Printf("err FakeBeginBlocker DistributeCoins: %v \n", err)
			errorString := fmt.Sprintf(types.ErrAutoTxDistribution, err.Error())
			k.AddAutoTxHistory(ctx, &autoTxHistory, &autoTx, timeOfBlock, fee, false, nil, errorString)
		} else {
			err, executedLocally, msgResponses := k.SendAutoTx(ctx, &autoTx)
			if err != nil {
				k.AddAutoTxHistory(ctx, &autoTxHistory, &autoTx, ctx.BlockTime(), fee, executedLocally, msgResponses, fmt.Sprintf(types.ErrAutoTxMsgHandling, err.Error()))
			} else {
				k.AddAutoTxHistory(ctx, &autoTxHistory, &autoTx, ctx.BlockTime(), fee, executedLocally, msgResponses)
			}

			shouldRecur := isRecurring && (autoTx.ExecTime.Add(autoTx.Interval).Before(autoTx.EndTime) || autoTx.ExecTime.Add(autoTx.Interval) == autoTx.EndTime)
			allowedToRecur := (!autoTx.Configuration.StopOnSuccess && !autoTx.Configuration.StopOnFailure) || autoTx.Configuration.StopOnSuccess && err != nil || autoTx.Configuration.StopOnFailure && err == nil
			//fmt.Printf("%v %v\n", shouldRecur, allowedToRecur)
			if shouldRecur && allowedToRecur {

				autoTx.ExecTime = autoTx.ExecTime.Add(autoTx.Interval)
				k.InsertAutoTxQueue(ctx, autoTx.TxID, autoTx.ExecTime)
			}

		}
		k.SetAutoTxInfo(ctx, &autoTx)
	}
}

// we may reimplement this as a configuration-based gas fee
func calculateTimeBasedFlexFee(autoTx types.AutoTxInfo, AutoTxHistory types.AutoTxHistory) sdkmath.Int {
	if len(AutoTxHistory.History) != 0 {
		prevEntry := AutoTxHistory.History[len(AutoTxHistory.History)-1].ActualExecTime
		period := (autoTx.ExecTime.Sub(prevEntry))
		return sdk.NewInt(int64(period.Milliseconds()))
	}

	period := autoTx.ExecTime.Sub(autoTx.StartTime)
	if period.Seconds() <= 60 {
		//base fee so we do not have a zero fee
		return sdk.NewInt(6_000)
	}
	return sdk.NewInt(int64(period.Seconds() * 10))
}
