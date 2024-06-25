package intent

import (
	"fmt"
	"time"

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

		// check conditions
		if !k.AllowedToExecute(ctx, action) {
			k.AddActionHistory(ctx, &action, timeOfBlock, sdk.Coin{}, false, nil, types.ErrActionConditions)
			action.ExecTime = action.ExecTime.Add(action.Interval)
			k.SetActionInfo(ctx, &action)
			continue
		}

		logger.Debug("action execution", "id", action.ID)

		isRecurring := action.ExecTime.Before(action.EndTime)

		flexFee := k.CalculateTimeBasedFlexFee(ctx, action)
		fee, err := k.DistributeCoins(ctx, action, flexFee, isRecurring, req.Header.ProposerAddress)

		k.RemoveFromActionQueue(ctx, action)
		if err != nil {
			errorString := fmt.Sprintf(types.ErrActionFeeDistribution, err.Error())
			k.AddActionHistory(ctx, &action, timeOfBlock, fee, false, nil, errorString)
		} else {
			executedLocally, msgResponses, err := k.TriggerAction(ctx, &action)
			if err != nil {
				k.AddActionHistory(ctx, &action, ctx.BlockTime(), fee, executedLocally, msgResponses, fmt.Sprintf(types.ErrActionMsgHandling, err.Error()))
			} else {
				k.AddActionHistory(ctx, &action, ctx.BlockTime(), fee, executedLocally, msgResponses)
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
