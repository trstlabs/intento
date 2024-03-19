package keeper_test

import (
	//"fmt"
	"fmt"

	"time"

	sdkmath "cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trstlabs/intento/x/intent/keeper"
	"github.com/trstlabs/intento/x/intent/types"
)

// BeginBlocker called every block, processes auto execution
func FakeBeginBlocker(ctx sdk.Context, k keeper.Keeper, fakeProposer sdk.ConsAddress) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	available := k.GetRelayerRewardsAvailability(ctx)
	if !available {
		k.SetRelayerRewardsAvailability(ctx, true)
	}

	autoTxs := k.GetActionsForBlock(ctx)

	timeOfBlock := ctx.BlockHeader().Time
	for _, autoTx := range autoTxs {
		autoTx = k.GetActionInfo(ctx, autoTx.ID)
		autoTxHistory, _ := k.TryGetActionHistory(ctx, autoTx.ID)
		if !k.AllowedToExecute(ctx, autoTx) {
			k.AddActionHistory(ctx, &autoTxHistory, &autoTx, timeOfBlock, sdk.Coin{}, false, nil, types.ErrActionConditions)
			autoTx.ExecTime = autoTx.ExecTime.Add(autoTx.Interval)
			k.SetActionInfo(ctx, &autoTx)
		}
		isRecurring := autoTx.ExecTime.Before(autoTx.EndTime)

		flexFee := calculateTimeBasedFlexFee(autoTx, autoTxHistory)
		fee, err := k.DistributeCoins(ctx, autoTx, flexFee, isRecurring, fakeProposer)

		k.RemoveFromActionQueue(ctx, autoTx)
		if err != nil {
			fmt.Printf("err FakeBeginBlocker DistributeCoins: %v \n", err)
			errorString := fmt.Sprintf(types.ErrActionDistribution, err.Error())
			k.AddActionHistory(ctx, &autoTxHistory, &autoTx, timeOfBlock, fee, false, nil, errorString)
		} else {
			err, executedLocally, msgResponses := k.SendAction(ctx, &autoTx)
			if err != nil {
				k.AddActionHistory(ctx, &autoTxHistory, &autoTx, ctx.BlockTime(), fee, executedLocally, msgResponses, fmt.Sprintf(types.ErrActionMsgHandling, err.Error()))
			} else {
				k.AddActionHistory(ctx, &autoTxHistory, &autoTx, ctx.BlockTime(), fee, executedLocally, msgResponses)
			}

			shouldRecur := isRecurring && (autoTx.ExecTime.Add(autoTx.Interval).Before(autoTx.EndTime) || autoTx.ExecTime.Add(autoTx.Interval) == autoTx.EndTime)
			allowedToRecur := (!autoTx.Configuration.StopOnSuccess && !autoTx.Configuration.StopOnFailure) || autoTx.Configuration.StopOnSuccess && err != nil || autoTx.Configuration.StopOnFailure && err == nil
			//fmt.Printf("%v %v\n", shouldRecur, allowedToRecur)
			if shouldRecur && allowedToRecur {

				autoTx.ExecTime = autoTx.ExecTime.Add(autoTx.Interval)
				k.InsertActionQueue(ctx, autoTx.ID, autoTx.ExecTime)
			}

		}
		k.SetActionInfo(ctx, &autoTx)
	}
}

// we may reimplement this as a configuration-based gas fee
func calculateTimeBasedFlexFee(autoTx types.ActionInfo, ActionHistory types.ActionHistory) sdkmath.Int {
	if len(ActionHistory.History) != 0 {
		prevEntry := ActionHistory.History[len(ActionHistory.History)-1].ActualExecTime
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
