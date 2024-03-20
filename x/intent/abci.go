package intent

import (
	"fmt"
	"time"

	sdkmath "cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trstlabs/intento/x/intent/keeper"
	"github.com/trstlabs/intento/x/intent/types"
)

// BeginBlocker called every block, processes auto execution
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	available := k.GetRelayerRewardsAvailability(ctx)
	if !available {
		k.SetRelayerRewardsAvailability(ctx, true)
	}

	logger := k.Logger(ctx)
	actions := k.GetActionsForBlock(ctx)

	timeOfBlock := ctx.BlockHeader().Time
	for _, action := range actions {
		actionHistory, _ := k.TryGetActionHistory(ctx, action.ID)
		// check dependent txs
		if !k.AllowedToExecute(ctx, action) {
			k.AddActionHistory(ctx, &actionHistory, &action, timeOfBlock, sdk.Coin{}, false, nil, types.ErrActionConditions)
			action.ExecTime = action.ExecTime.Add(action.Interval)
			k.SetActionInfo(ctx, &action)
			continue
		}

		logger.Debug("action execution", "id", action.ID)

		isRecurring := action.ExecTime.Before(action.EndTime)

		flexFee := calculateTimeBasedFlexFee(action, actionHistory)
		fee, err := k.DistributeCoins(ctx, action, flexFee, isRecurring, req.Header.ProposerAddress)

		k.RemoveFromActionQueue(ctx, action)
		if err != nil {
			errorString := fmt.Sprintf(types.ErrActionDistribution, err.Error())
			k.AddActionHistory(ctx, &actionHistory, &action, timeOfBlock, fee, false, nil, errorString)
		} else {
			err, executedLocally, msgResponses := k.SendAction(ctx, &action)
			if err != nil {
				k.AddActionHistory(ctx, &actionHistory, &action, ctx.BlockTime(), fee, executedLocally, msgResponses, fmt.Sprintf(types.ErrActionMsgHandling, err.Error()))
			} else {
				k.AddActionHistory(ctx, &actionHistory, &action, ctx.BlockTime(), fee, executedLocally, msgResponses)
			}

			// setting new ExecTime and adding a new entry into the queue based on interval
			shouldRecur := isRecurring && (action.ExecTime.Add(action.Interval).Before(action.EndTime) || action.ExecTime.Add(action.Interval) == action.EndTime)
			allowedToRecur := (!action.Configuration.StopOnSuccess && !action.Configuration.StopOnFailure) || action.Configuration.StopOnSuccess && err != nil || action.Configuration.StopOnFailure && err == nil

			if shouldRecur && allowedToRecur {
				//fmt.Printf("action will recur: %v \n", action.ID)
				action.ExecTime = action.ExecTime.Add(action.Interval)
				k.InsertActionQueue(ctx, action.ID, action.ExecTime)
			}
		}
		k.SetActionInfo(ctx, &action)
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeAction,
				sdk.NewAttribute(types.AttributeKeyActionID, fmt.Sprint(action.ID)),
				sdk.NewAttribute(types.AttributeKeyActionOwner, action.Owner),
			),
		)
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
		return sdk.NewInt(60_000)
	}
	return sdk.NewInt(int64(period.Milliseconds()))
}
