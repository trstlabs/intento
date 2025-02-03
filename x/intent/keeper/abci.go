package keeper

import (
	"encoding/base64"
	"fmt"
	"time"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trstlabs/intento/x/intent/types"
	icqtypes "github.com/trstlabs/intento/x/interchainquery/types"
)

// HandleAction processes a single action during the block
func (k Keeper) HandleAction(ctx sdk.Context, logger log.Logger, action types.ActionInfo, timeOfBlock time.Time, queryCallback []byte) {
	var (
		errorString     = ""
		fee             = sdk.Coin{}
		executedLocally = false
		msgResponses    = []*cdctypes.Any{}
	)

	k.RemoveFromActionQueue(ctx, action)

	allowed, err := k.allowedToExecute(ctx, action)
	if !allowed {
		k.recordActionNotAllowed(ctx, &action, timeOfBlock, err)
		return
	}

	actionCtx := ctx.WithGasMeter(storetypes.NewGasMeter(types.MaxGas))
	cacheCtx, writeCtx := actionCtx.CacheContext()

	logger.Debug("action execution", "id", action.ID)

	feeAddr, feeDenom, err := k.GetFeeAccountForMinFees(cacheCtx, action, types.MaxGas)
	if err != nil {
		errorString = appendError(errorString, err.Error())
	} else if feeAddr == nil || feeDenom == "" {
		errorString = appendError(errorString, (types.ErrBalanceTooLow + feeDenom))
	}

	if errorString == "" {
		executedLocally, errorString = k.handleActionExecution(cacheCtx, &action, &msgResponses, errorString)

		fee, err = k.DistributeCoins(cacheCtx, action, feeAddr, feeDenom, ctx.BlockHeader().ProposerAddress)
		if err != nil {
			errorString = appendError(errorString, fmt.Sprintf(types.ErrActionFeeDistribution, err.Error()))
		}

	}

	k.addActionHistoryEntry(cacheCtx, &action, timeOfBlock, fee, executedLocally, msgResponses, errorString)

	writeCtx()

	if shouldRecur(action, errorString) {
		action.ExecTime = action.ExecTime.Add(action.Interval)
		k.InsertActionQueue(ctx, action.ID, action.ExecTime)
	}

	emitActionEvent(ctx, action)
	k.SetActionInfo(ctx, &action)
}

// submitInterchainQuery submits an interchain query when ICQConfig is present
func (k Keeper) SubmitInterchainQueries(ctx sdk.Context, action types.ActionInfo, logger log.Logger) {
	for i, feedbackLoop := range action.Conditions.FeedbackLoops {
		if feedbackLoop.ICQConfig == nil {
			continue
		}
		icqID := fmt.Sprintf("%s:%d:%d", types.ActionFeedbackLoopQueryKeyPrefix, action.ID, i)
		_, found := k.interchainQueryKeeper.GetQuery(ctx, icqID)
		if found {
			return //(or continue)
		}
		err := k.SubmitInterchainQuery(ctx, *feedbackLoop.ICQConfig, logger, icqID)
		if err != nil {
			k.SetActionHistoryEntry(ctx, action.ID, &types.ActionHistoryEntry{Errors: []string{fmt.Sprint("Error submitting ICQ: decoding Base64 string: ", err)}})
			return
		}
	}

	for i, comparison := range action.Conditions.Comparisons {
		if comparison.ICQConfig == nil {
			continue
		}
		icqID := fmt.Sprintf("%s:%d:%d", types.ActionFeedbackLoopQueryKeyPrefix, action.ID, i)
		_, found := k.interchainQueryKeeper.GetQuery(ctx, icqID)
		if found {
			return
		}
		err := k.SubmitInterchainQuery(ctx, *comparison.ICQConfig, logger, icqID)
		if err != nil {
			k.SetActionHistoryEntry(ctx, action.ID, &types.ActionHistoryEntry{Errors: []string{fmt.Sprint("Error submitting ICQ: decoding Base64 string: ", err)}})
			return
		}
	}

}

// submitInterchainQuery submits an interchain query when ICQConfig is present
func (k Keeper) SubmitInterchainQuery(ctx sdk.Context, icqConfig types.ICQConfig, logger log.Logger, id string) error {
	requestData, err := base64.StdEncoding.DecodeString(icqConfig.QueryKey)
	if err != nil {
		return err
	}

	query := icqtypes.Query{
		Id:              id,
		CallbackId:      ICQCallbackID_Action,
		ConnectionId:    icqConfig.ConnectionId,
		ChainId:         icqConfig.ChainId,
		TimeoutPolicy:   icqConfig.TimeoutPolicy,
		TimeoutDuration: icqConfig.TimeoutDuration,
		QueryType:       icqConfig.QueryType,
		RequestData:     requestData,
		CallbackModule:  types.ModuleName,
	}
	k.interchainQueryKeeper.SetQuery(ctx, query)
	// Log successful submission of the interchain query
	logger.Debug("interchain query submitted", "ICQ ID", id)
	return nil
}

// handleActionExecution handles the core logic of triggering an action and processing responses
func (k Keeper) handleActionExecution(ctx sdk.Context, action *types.ActionInfo, msgResponses *[]*cdctypes.Any, errorString string) (bool, string) {
	var executedLocally bool
	// Safe check to ensure conditions exist before proceeding
	if action.Conditions == nil || action.Conditions.FeedbackLoops == nil {
		// If no FeedbackLoops conditions are present, just execute all Msgs normally
		executedLocally, responses, err := k.TriggerAction(ctx, action)
		if err != nil {
			return false, appendError(errorString, fmt.Sprintf(types.ErrActionMsgHandling, err.Error()))
		}
		*msgResponses = append(*msgResponses, responses...)
		return executedLocally, errorString
	}

	// FeedbackLoops is set, handle accordingly
	executedLocally, errorString = k.handleRunFeedbackLoops(ctx, action, msgResponses, errorString)

	return executedLocally, errorString
}

// handleFeedbackLoops processes the FeedbackLoops logic when it's set
func (k Keeper) handleRunFeedbackLoops(ctx sdk.Context, action *types.ActionInfo, msgResponses *[]*cdctypes.Any, errorString string) (bool, string) {
	executedLocally := false
	actionTmp := *action
	for _, feedbackLoop := range action.Conditions.FeedbackLoops {
		// Ensure MsgsIndex is valid and avoid out-of-bound slices
		if feedbackLoop.MsgsIndex > 0 && int(feedbackLoop.MsgsIndex) <= len(action.Msgs) {
			// Split action.Msgs based on MsgsIndex (first part)
			actionTmp.Msgs = action.Msgs[:feedbackLoop.MsgsIndex]
		}

		// Trigger the first set of Msgs (before MsgsIndex)
		executedLocally, responses, err := k.TriggerAction(ctx, &actionTmp)
		if err != nil {
			return false, appendError(errorString, fmt.Sprintf(types.ErrActionMsgHandling, err.Error()))
		}
		*msgResponses = append(*msgResponses, responses...)

		// If ActionID is set, attempt to process its value based on the response of the first set
		if feedbackLoop.ActionID != 0 || executedLocally || feedbackLoop.ICQConfig != nil {
			err = k.RunFeedbackLoops(ctx, action.ID, &actionTmp.Msgs, action.Conditions)
			if err != nil {
				return false, appendError(errorString, fmt.Sprintf(types.ErrSettingActionResult, err.Error()))
			}

			// Process the remaining Msgs (after MsgsIndex)
			if feedbackLoop.MsgsIndex > 0 && int(feedbackLoop.MsgsIndex) < len(action.Msgs) {
				// Handle the remaining Msgs
				actionTmp.Msgs = action.Msgs[feedbackLoop.MsgsIndex:]

				// Trigger the remaining Msgs
				_, additionalResponses, err := k.TriggerAction(ctx, &actionTmp)
				if err != nil {
					return false, appendError(errorString, fmt.Sprintf(types.ErrActionMsgHandling, err.Error()))
				}
				*msgResponses = append(*msgResponses, additionalResponses...)
			}
		}
	}
	return executedLocally, errorString
}
func (k Keeper) allowedToExecute(ctx sdk.Context, action types.ActionInfo) (bool, error) {
	shouldRecur := action.ExecTime.Before(action.EndTime) && action.ExecTime.Add(action.Interval).Before(action.EndTime)
	conditions := action.Conditions
	if conditions == nil {
		return true, nil
	}

	// Step 1: Check comparisons
	comparisonsResult, err := k.checkComparisons(ctx, action, conditions)
	if err != nil {
		if shouldRecur {
			k.scheduleNextExecution(ctx, action)
		}
		return false, err
	}

	// Step 2: Check dependent actions
	dependenciesResult := k.checkDependentActions(ctx, conditions, &shouldRecur)

	// Combine results
	allowedToExecute := comparisonsResult && dependenciesResult

	// Step 3: Handle recurring actions if not allowed to execute
	if !allowedToExecute && shouldRecur {
		k.scheduleNextExecution(ctx, action)
	}

	return allowedToExecute, nil
}

// checkComparisons evaluates the conditions.Comparisons based on AND/OR logic.
func (k Keeper) checkComparisons(ctx sdk.Context, action types.ActionInfo, conditions *types.ExecutionConditions) (bool, error) {
	var err error = nil

	if conditions.Comparisons == nil {
		return true, nil
	}

	if conditions.UseAndForComparisons {
		// AND logic: All comparisons must evaluate to true
		for _, comparison := range conditions.Comparisons {
			isTrue, err := k.evaluateComparison(ctx, action, *comparison)
			if err != nil || !isTrue {
				return false, err
			}
		}
		return true, nil
	}

	// OR logic: At least one comparison must evaluate to true
	for _, comparison := range conditions.Comparisons {
		isTrue, err := k.evaluateComparison(ctx, action, *comparison)
		if err == nil && isTrue {
			return true, nil
		}
	}

	return false, err
}

// evaluateComparison checks a single comparison against the action history.
func (k Keeper) evaluateComparison(ctx sdk.Context, action types.ActionInfo, comparison types.Comparison) (bool, error) {
	actionID := action.ID
	if comparison.ActionID != 0 {
		actionID = comparison.ActionID
	}

	if !k.HasActionHistoryEntry(ctx, actionID) {
		return true, nil // No history means there's nothing to compare against
	}

	history, err := k.GetActionHistory(ctx, actionID)
	if err != nil {
		return false, err
	}

	if history[len(history)-1].MsgResponses == nil && history[len(history)-1].QueryResponses == nil {
		// if we should stop on failure or the condition is explicitly set to need all to be true, we return with an error here
		if action.Configuration.StopOnFailure || action.Conditions.UseAndForComparisons {
			return false, fmt.Errorf("did not make a comparison, no responses on the target history. Set SaveResponses to true to use responses for comparison")
		}
		return true, nil //default to true in case there is nothing to compare against
	}

	responses := history[len(history)-1].MsgResponses

	isTrue, err := k.CompareResponseValue(ctx, action.ID, responses, comparison)
	if err != nil {
		return false, fmt.Errorf("error in CompareResponseValue: %w", err)
	}

	return isTrue, nil
}

// checkDependentActions verifies success/failure conditions of dependent actions.
func (k Keeper) checkDependentActions(ctx sdk.Context, conditions *types.ExecutionConditions, shouldRecur *bool) bool {
	checkDependency := func(actionIDs []uint64, successCondition bool, isStop bool) bool {
		for _, actionID := range actionIDs {
			if !k.HasActionHistoryEntry(ctx, actionID) {
				continue
			}

			history, err := k.GetActionHistory(ctx, actionID)
			if err != nil || len(history) == 0 {
				return false
			}

			success := history[len(history)-1].Executed && history[len(history)-1].Errors == nil
			if success == successCondition {
				if isStop {
					// Stop means do not execute and do not recur
					*shouldRecur = false
					return false
				} else {
					// Skip means do not execute, but allow to recur
					return false
				}
			}
		}
		return true
	}

	// Stop on success
	if !checkDependency(conditions.StopOnSuccessOf, true, true) {
		return false
	}

	// Stop on failure
	if !checkDependency(conditions.StopOnFailureOf, false, true) {
		return false
	}

	// Skip on success
	if !checkDependency(conditions.SkipOnSuccessOf, true, false) {
		return false
	}

	// Skip on failure
	if !checkDependency(conditions.SkipOnFailureOf, false, false) {
		return false
	}

	return true
}

// scheduleNextExecution schedules the next execution for recurring actions.
func (k Keeper) scheduleNextExecution(ctx sdk.Context, action types.ActionInfo) {
	nextExecTime := action.ExecTime.Add(action.Interval)
	k.InsertActionQueue(ctx, action.ID, nextExecTime)
	action.ExecTime = nextExecTime
	k.SetActionInfo(ctx, &action)
}

// recordActionNotAllowed adds an action entry to the action history
func (k Keeper) recordActionNotAllowed(ctx sdk.Context, action *types.ActionInfo, timeOfBlock time.Time, errorMsg error) {
	k.Logger(ctx).Debug("action not allowed to execute", "id", action.ID)
	if errorMsg != nil {
		k.addActionHistoryEntry(ctx, action, timeOfBlock, sdk.Coin{}, false, nil, fmt.Sprintf(types.ErrActionConditions, errorMsg.Error()))
		return
	}
	k.addActionHistoryEntry(ctx, action, timeOfBlock, sdk.Coin{}, false, nil, "")
}

// shouldRecur checks whether the action should be rescheduled based on recurrence rules
func shouldRecur(action types.ActionInfo, errorString string) bool {
	isRecurring := action.ExecTime.Before(action.EndTime) && (action.ExecTime.Add(action.Interval).Before(action.EndTime) || action.ExecTime.Add(action.Interval).Equal(action.EndTime))
	allowedToRecur := (!action.Configuration.StopOnSuccess && !action.Configuration.StopOnFailure) ||
		(action.Configuration.StopOnSuccess && errorString != "") ||
		(action.Configuration.StopOnFailure && errorString == "")

	return isRecurring && allowedToRecur
}

// // appendToPriorActionHistory appends results to the prior history entry for the action
// func (k Keeper) appendToPriorActionHistory(ctx sdk.Context, action *types.ActionInfo, fee sdk.Coin, executedLocally bool, msgResponses []*cdctypes.Any, queryResponse string, errorString string) {
// 	// Fetch the last recorded action history for the action
// 	entry, found := k.getCurrentActionHistoryEntry(ctx, action.ID)
// 	if !found {
// 		return
// 	}
// 	// Append the new data to the existing history entry
// 	entry.ExecFee = entry.ExecFee.Add(fee)
// 	entry.Executed = entry.Executed || executedLocally
// 	if action.Configuration.SaveResponses {
// 		entry.MsgResponses = append(entry.MsgResponses, msgResponses...)
// 		entry.QueryResponse = queryResponse
// 		if errorString != "" {
// 			entry.Errors = append(entry.Errors, errorString)
// 		}
// 	}
// 	// Update the action history with the new appended data
// 	k.SetCurrentActionHistoryEntry(ctx, action.ID, entry)
// }

// emitActionEvent creates an event for the action execution
func emitActionEvent(ctx sdk.Context, action types.ActionInfo) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeAction,
			sdk.NewAttribute(types.AttributeKeyActionID, fmt.Sprint(action.ID)),
			sdk.NewAttribute(types.AttributeKeyActionOwner, action.Owner),
		),
	)
}

// appendError appends a new error to an existing error string
func appendError(existingError, newError string) string {
	if existingError == "" {
		return newError
	}
	return existingError + "; " + newError
}
