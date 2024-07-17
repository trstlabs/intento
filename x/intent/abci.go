package intent

import (
	"fmt"
	"time"

	abci "github.com/cometbft/cometbft/abci/types"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
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
		errorString := ""
		fee := sdk.Coin{}
		executedLocally := false
		msgResponses := []*cdctypes.Any{}

		allowed, err := k.AllowedToExecute(ctx, action)
		// check conditions
		if !allowed {
			k.AddActionHistory(ctx, &action, timeOfBlock, sdk.Coin{}, false, nil, fmt.Sprintf(types.ErrActionConditions, err.Error()))
			action.ExecTime = action.ExecTime.Add(action.Interval)
			k.SetActionInfo(ctx, &action)
			continue
		}

		logger.Debug("action execution", "id", action.ID)

		isRecurring := action.ExecTime.Before(action.EndTime)
		k.RemoveFromActionQueue(ctx, action)
		actionCtx := ctx.WithGasMeter(sdk.NewGasMeter(1_000_000))

		cacheCtx, writeCtx := actionCtx.CacheContext()
		feeAddr, feeDenom, err := k.GetFeeAccountForMinFees(cacheCtx, action, 1_000_000)
		if err != nil || feeAddr == nil || feeDenom == "" {
			errorString = types.ErrBalanceLow
		}
		if errorString == "" {
			err = k.UseResponseValue(cacheCtx, action.ID, &action.Msgs, action.Conditions)
			if err != nil {
				errorString = fmt.Sprintf(types.ErrActionResponseUseValue, err.Error())
			}

			if errorString == "" {
				//Handle response parsing

				if action.Conditions == nil || action.Conditions.UseResponseValue == nil || action.Conditions.UseResponseValue.MsgsIndex == 0 {
					executedLocally, msgResponses, err = k.TriggerAction(cacheCtx, &action)
					if err != nil {
						errorString = fmt.Sprintf(types.ErrActionMsgHandling, err.Error())
					}
				} else {
					actionTmp := action
					actionTmp.Msgs = action.Msgs[:action.Conditions.UseResponseValue.MsgsIndex+1]
					executedLocally, msgResponses, err = k.TriggerAction(cacheCtx, &actionTmp)
					if err != nil {
						errorString = fmt.Sprintf(types.ErrSettingActionResult + err.Error())
					}
					if errorString == "" {
						err = k.UseResponseValue(cacheCtx, action.ID, &actionTmp.Msgs, action.Conditions)
						if err != nil {
							errorString = fmt.Sprintf(types.ErrSettingActionResult + err.Error())

						} else if executedLocally {
							actionTmp.Msgs = action.Msgs[action.Conditions.UseResponseValue.MsgsIndex+1:]
							_, msgResponses2, err2 := k.TriggerAction(cacheCtx, &actionTmp)
							errorString = fmt.Sprintf(types.ErrActionMsgHandling, err2)
							msgResponses = append(msgResponses, msgResponses2...)
						}
					}
				}
			}
			fee, err = k.DistributeCoins(cacheCtx, action, feeAddr, feeDenom, isRecurring, ctx.BlockHeader().ProposerAddress)
			if err != nil {
				errorString = fmt.Sprintf(types.ErrActionFeeDistribution, err.Error())
			}
		}

		k.AddActionHistory(cacheCtx, &action, timeOfBlock, fee, executedLocally, msgResponses, errorString)
		writeCtx()
		// setting new ExecTime and adding a new entry into the queue based on interval
		shouldRecur := isRecurring && (action.ExecTime.Add(action.Interval).Before(action.EndTime) || action.ExecTime.Add(action.Interval) == action.EndTime)
		allowedToRecur := (!action.Configuration.StopOnSuccess && !action.Configuration.StopOnFailure) || action.Configuration.StopOnSuccess && err != nil || action.Configuration.StopOnFailure && err == nil

		if shouldRecur && allowedToRecur {
			action.ExecTime = action.ExecTime.Add(action.Interval)
			k.InsertActionQueue(ctx, action.ID, action.ExecTime)
		}
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeAction,
				sdk.NewAttribute(types.AttributeKeyActionID, fmt.Sprint(action.ID)),
				sdk.NewAttribute(types.AttributeKeyActionOwner, action.Owner),
			),
		)

		k.SetActionInfo(ctx, &action)
	}
}
