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

	actions := k.GetActionsForBlock(ctx)

	timeOfBlock := ctx.BlockHeader().Time
	for _, action := range actions {
		action = k.GetActionInfo(ctx, action.ID)
		actionHistory, _ := k.TryGetActionHistory(ctx, action.ID)
		if !k.AllowedToExecute(ctx, action) {
			k.AddActionHistory(ctx, &actionHistory, &action, timeOfBlock, sdk.Coin{}, false, nil, types.ErrActionConditions)
			action.ExecTime = action.ExecTime.Add(action.Interval)
			k.SetActionInfo(ctx, &action)
		}
		isRecurring := action.ExecTime.Before(action.EndTime)

		flexFee := calculateTimeBasedFlexFee(action, actionHistory)
		fee, err := k.DistributeCoins(ctx, action, flexFee, isRecurring, fakeProposer)

		k.RemoveFromActionQueue(ctx, action)
		if err != nil {
			fmt.Printf("err FakeBeginBlocker DistributeCoins: %v \n", err)
			errorString := fmt.Sprintf(types.ErrActionDistribution, err.Error())
			k.AddActionHistory(ctx, &actionHistory, &action, timeOfBlock, fee, false, nil, errorString)
		} else {
			err, executedLocally, msgResponses := k.SendAction(ctx, &action)
			if err != nil {
				k.AddActionHistory(ctx, &actionHistory, &action, ctx.BlockTime(), fee, executedLocally, msgResponses, fmt.Sprintf(types.ErrActionMsgHandling, err.Error()))
			} else {
				k.AddActionHistory(ctx, &actionHistory, &action, ctx.BlockTime(), fee, executedLocally, msgResponses)
			}

			shouldRecur := isRecurring && (action.ExecTime.Add(action.Interval).Before(action.EndTime) || action.ExecTime.Add(action.Interval) == action.EndTime)
			allowedToRecur := (!action.Configuration.StopOnSuccess && !action.Configuration.StopOnFailure) || action.Configuration.StopOnSuccess && err != nil || action.Configuration.StopOnFailure && err == nil
			//fmt.Printf("%v %v\n", shouldRecur, allowedToRecur)
			if shouldRecur && allowedToRecur {

				action.ExecTime = action.ExecTime.Add(action.Interval)
				k.InsertActionQueue(ctx, action.ID, action.ExecTime)
			}

		}
		k.SetActionInfo(ctx, &action)
	}
}

// we may reimplement this as a configuration-based gas fee
func calculateTimeBasedFlexFee(action types.ActionInfo, ActionHistory types.ActionHistory) sdkmath.Int {
	if len(ActionHistory.History) != 0 {
		prevEntry := ActionHistory.History[len(ActionHistory.History)-1].ActualExecTime
		period := (action.ExecTime.Sub(prevEntry))
		return sdk.NewInt(int64(period.Milliseconds()))
	}

	period := action.ExecTime.Sub(action.StartTime)
	if period.Seconds() <= 60 {
		//base fee so we do not have a zero fee
		return sdk.NewInt(6_000)
	}
	return sdk.NewInt(int64(period.Seconds() * 10))
}
