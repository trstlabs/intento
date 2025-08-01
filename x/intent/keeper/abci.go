package keeper

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trstlabs/intento/x/intent/types"
	icqtypes "github.com/trstlabs/intento/x/interchainquery/types"
)

// HandleFlow processes a single flow during the block
func (k Keeper) HandleFlow(ctx sdk.Context, logger log.Logger, flow types.FlowInfo, timeOfBlock time.Time, queryCallback []byte) {
	var (
		errorString     = ""
		fee             = sdk.Coin{}
		executedLocally = false
		msgResponses    = []*cdctypes.Any{}
	)

	k.RemoveFromFlowQueue(ctx, flow)

	allowed, err := k.allowedToExecute(ctx, flow)
	if !allowed {
		k.recordFlowNotAllowed(ctx, &flow, timeOfBlock, err)
		return
	}

	flowCtx := ctx.WithGasMeter(storetypes.NewGasMeter(types.MaxGas))
	cacheCtx, writeCtx := flowCtx.CacheContext()

	logger.Debug("flow execution", "id", flow.ID)

	feeAddr, feeDenom, err := k.GetFeeAccountForMinFees(cacheCtx, flow, types.MaxGas)
	if err != nil {
		errorString = appendError(errorString, err.Error())
	} else if feeAddr == nil || feeDenom == "" {
		errorString = appendError(errorString, types.ErrBalanceTooLow)
	}

	if errorString == "" {
		executedLocally, errorString = k.handleFlowExecution(cacheCtx, &flow, &msgResponses, errorString)

		fee, err = k.DistributeCoins(cacheCtx, flow, feeAddr, feeDenom)
		if err != nil {
			errorString = appendError(errorString, fmt.Sprintf(types.ErrFlowFeeDistribution, err.Error()))
		}

	}

	k.addFlowHistoryEntry(cacheCtx, &flow, timeOfBlock, fee, executedLocally, msgResponses, errorString)

	writeCtx()

	if k.shouldRecur(ctx, flow, errorString) {
		flow.ExecTime = flow.ExecTime.Add(flow.Interval)
		k.InsertFlowQueue(ctx, flow.ID, flow.ExecTime)
	}

	emitFlowEvent(ctx, flow)
	k.SetFlowInfo(ctx, &flow)
}

// submitInterchainQuery submits an interchain query when ICQConfig is present
func (k Keeper) SubmitInterchainQueries(ctx sdk.Context, flow types.FlowInfo, logger log.Logger) {
	for i, feedbackLoop := range flow.Conditions.FeedbackLoops {
		if feedbackLoop.ICQConfig == nil {
			continue
		}
		icqID := fmt.Sprintf("%s:%d:%d", types.FlowFeedbackLoopQueryKeyPrefix, flow.ID, i)
		_, found := k.interchainQueryKeeper.GetQuery(ctx, icqID)
		if found {
			return //(or continue)
		}
		err := k.SubmitInterchainQuery(ctx, *feedbackLoop.ICQConfig, logger, icqID)
		if err != nil {
			k.SetFlowHistoryEntry(ctx, flow.ID, &types.FlowHistoryEntry{Errors: []string{fmt.Sprint("Error submitting ICQ: decoding Base64 string: ", err)}})
			return
		}
	}

	for i, comparison := range flow.Conditions.Comparisons {
		if comparison.ICQConfig == nil {
			continue
		}
		icqID := fmt.Sprintf("%s:%d:%d", types.FlowComparisonQueryKeyPrefix, flow.ID, i)
		_, found := k.interchainQueryKeeper.GetQuery(ctx, icqID)
		if found {
			return
		}
		err := k.SubmitInterchainQuery(ctx, *comparison.ICQConfig, logger, icqID)
		if err != nil {
			k.SetFlowHistoryEntry(ctx, flow.ID, &types.FlowHistoryEntry{Errors: []string{fmt.Sprint("Error submitting ICQ: decoding Base64 string: ", err)}})
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
		CallbackId:      ICQCallbackID_Flow,
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

// handleFlowExecution handles the core logic of triggering an flow and processing responses
func (k Keeper) handleFlowExecution(ctx sdk.Context, flow *types.FlowInfo, msgResponses *[]*cdctypes.Any, errorString string) (bool, string) {
	var executedLocally bool
	// Safe check to ensure conditions exist before proceeding
	if flow.Conditions == nil || flow.Conditions.FeedbackLoops == nil {
		// If no FeedbackLoops conditions are present, just execute all Msgs normally
		executedLocally, responses, err := k.TriggerFlow(ctx, flow)
		if err != nil {
			return false, appendError(errorString, fmt.Sprintf(types.ErrFlowMsgHandling, err.Error()))
		}
		*msgResponses = append(*msgResponses, responses...)
		return executedLocally, errorString
	}

	// FeedbackLoops is set, handle accordingly
	executedLocally, errorString = k.handleRunFeedbackLoops(ctx, flow, msgResponses, errorString)

	return executedLocally, errorString
}

func (k Keeper) handleRunFeedbackLoops(ctx sdk.Context, flow *types.FlowInfo, msgResponses *[]*cdctypes.Any, errorString string) (bool, string) {
	executedLocally := false

	// Group feedback loops by MsgsIndex to process them together
	feedbackGroups := make(map[uint32][]*types.FeedbackLoop)
	for _, loop := range flow.Conditions.FeedbackLoops {
		targetFlow := flow
		if loop.FlowID != 0 && loop.FlowID != flow.ID {
			otherFlow, err := k.TryGetFlowInfo(ctx, loop.FlowID)
			if err != nil {
				errorString = appendError(errorString, fmt.Sprintf("invalid feedback loop FlowID: %d (not found)", loop.FlowID))
				continue
			}
			targetFlow = &otherFlow
		}
		if int(loop.MsgsIndex) >= len(targetFlow.Msgs) {
			errorString = appendError(errorString, fmt.Sprintf("invalid feedback loop MsgsIndex: %d (only %d messages in flow %d)", loop.MsgsIndex, len(targetFlow.Msgs), targetFlow.ID))
			continue
		}
		feedbackGroups[loop.MsgsIndex] = append(feedbackGroups[loop.MsgsIndex], loop)
	}

	allResponses := []*cdctypes.Any{}

	for i := 0; i < len(flow.Msgs); {
		loops, hasFeedback := feedbackGroups[uint32(i)]

		if hasFeedback {
			preFeedbackMsgs := flow.Msgs[:i]

			if len(preFeedbackMsgs) > 0 {
				flowCopy := *flow
				flowCopy.Msgs = preFeedbackMsgs

				localExec, responses, err := k.TriggerFlow(ctx, &flowCopy)
				if err != nil {
					return false, appendError(errorString, fmt.Sprintf(types.ErrFlowMsgHandling, err.Error()))
				}
				executedLocally = executedLocally || localExec
				allResponses = append(allResponses, responses...)
				*msgResponses = append(*msgResponses, responses...)
			}

			remainingMsgs := make([]*cdctypes.Any, len(flow.Msgs[i:]))
			copy(remainingMsgs, flow.Msgs[i:])

			for _, loop := range loops {
				if int(loop.ResponseIndex) < len(allResponses) {
					modifiedMsgs := make([]*cdctypes.Any, len(remainingMsgs))
					copy(modifiedMsgs, remainingMsgs)

					err := k.RunFeedbackLoops(ctx, flow.ID, &modifiedMsgs, &types.ExecutionConditions{
						FeedbackLoops: []*types.FeedbackLoop{loop},
					})
					if err != nil {
						return false, appendError(errorString, fmt.Sprintf(types.ErrSettingFlowResult, err.Error()))
					}

					if !msgsEqual(remainingMsgs, modifiedMsgs) {

						remainingMsgs = modifiedMsgs
					}

				}
			}

			// If there are remaining messages after processing feedback loops, execute them if they are local
			if len(remainingMsgs) > 0 && executedLocally {
				flowCopy := *flow
				flowCopy.Msgs = remainingMsgs
				localExec, responses, err := k.TriggerFlow(ctx, &flowCopy)
				if err != nil {
					return false, appendError(errorString, fmt.Sprintf(types.ErrFlowMsgHandling, err.Error()))
				}
				executedLocally = executedLocally || localExec
				*msgResponses = append(*msgResponses, responses...)
			}
			break
		}

		i++
	}

	return executedLocally, errorString
}

func msgsEqual(a, b []*cdctypes.Any) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !a[i].Equal(b[i]) {
			return false
		}
	}
	return true
}

func (k Keeper) allowedToExecute(ctx sdk.Context, flow types.FlowInfo) (bool, error) {
	shouldRecur := flow.ExecTime.Before(flow.EndTime) && flow.ExecTime.Add(flow.Interval).Before(flow.EndTime)
	conditions := flow.Conditions
	if conditions == nil {
		return true, nil
	}

	// Step 1: Check comparisons
	comparisonsResult, err := k.checkComparisons(ctx, flow, conditions)
	if err != nil {
		if shouldRecur {
			k.scheduleNextExecution(ctx, flow)
		}
		return false, err
	}

	// Step 2: Check dependent flows
	dependenciesResult := k.checkDependentFlows(ctx, conditions, &shouldRecur)

	// Combine results
	allowedToExecute := comparisonsResult && dependenciesResult

	// Step 3: Handle recurring flows if not allowed to execute
	if !allowedToExecute && shouldRecur {
		k.scheduleNextExecution(ctx, flow)
	}

	return allowedToExecute, nil
}

// checkComparisons evaluates the conditions.Comparisons based on AND/OR logic.
func (k Keeper) checkComparisons(ctx sdk.Context, flow types.FlowInfo, conditions *types.ExecutionConditions) (bool, error) {
	var err error = nil

	if conditions.Comparisons == nil {
		return true, nil
	}

	if conditions.UseAndForComparisons {
		// AND logic: All comparisons must evaluate to true
		for _, comparison := range conditions.Comparisons {
			isTrue, err := k.evaluateComparison(ctx, flow, *comparison)
			if err != nil || !isTrue {
				return false, err
			}
		}
		return true, nil
	}

	// OR logic: At least one comparison must evaluate to true
	for _, comparison := range conditions.Comparisons {
		isTrue, err := k.evaluateComparison(ctx, flow, *comparison)
		if err == nil && isTrue {
			return true, nil
		}
	}

	return false, err
}

// evaluateComparison checks a single comparison against the flow history.
func (k Keeper) evaluateComparison(ctx sdk.Context, flow types.FlowInfo, comparison types.Comparison) (bool, error) {
	flowID := flow.ID
	if comparison.FlowID != 0 {
		flowID = comparison.FlowID
	}

	if !k.HasFlowHistoryEntry(ctx, flowID) {
		return true, nil // No history means there's nothing to compare against
	}

	history, err := k.GetFlowHistory(ctx, flowID)
	if err != nil {
		return false, err
	}

	if history[len(history)-1].MsgResponses == nil && history[len(history)-1].QueryResponses == nil {
		// if we should stop on failure or the condition is explicitly set to need all to be true, we return with an error here
		if flow.Configuration.StopOnFailure || flow.Conditions.UseAndForComparisons {
			return false, fmt.Errorf("did not make a comparison, no responses on the target history. Set SaveResponses to true to use responses for comparison")
		}
		return true, nil //default to true in case there is nothing to compare against
	}

	responses := history[len(history)-1].MsgResponses

	isTrue, err := k.CompareResponseValue(ctx, flow.ID, responses, comparison)
	if err != nil {
		return false, fmt.Errorf("error in CompareResponseValue: %w", err)
	}

	return isTrue, nil
}

// checkDependentFlows verifies success/failure conditions of dependent flows.
func (k Keeper) checkDependentFlows(ctx sdk.Context, conditions *types.ExecutionConditions, shouldRecur *bool) bool {
	checkDependency := func(flowIDs []uint64, successCondition bool, isStop bool) bool {
		for _, flowID := range flowIDs {
			if !k.HasFlowHistoryEntry(ctx, flowID) {
				continue
			}

			history, err := k.GetFlowHistory(ctx, flowID)
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

// scheduleNextExecution schedules the next execution for recurring flows.
func (k Keeper) scheduleNextExecution(ctx sdk.Context, flow types.FlowInfo) {
	nextExecTime := flow.ExecTime.Add(flow.Interval)
	k.InsertFlowQueue(ctx, flow.ID, nextExecTime)
	flow.ExecTime = nextExecTime
	k.SetFlowInfo(ctx, &flow)
}

// recordFlowNotAllowed adds an flow entry to the flow history
func (k Keeper) recordFlowNotAllowed(ctx sdk.Context, flow *types.FlowInfo, timeOfBlock time.Time, errorMsg error) {
	k.Logger(ctx).Debug("flow not allowed to execute", "id", flow.ID)
	if errorMsg != nil {
		k.addFlowHistoryEntry(ctx, flow, timeOfBlock, sdk.Coin{}, false, nil, fmt.Sprintf(types.ErrFlowConditions, errorMsg.Error()))
		return
	}
	k.addFlowHistoryEntry(ctx, flow, timeOfBlock, sdk.Coin{}, false, nil, "")
}

// shouldRecur checks whether the flow should be rescheduled based on recurrence rules
func (k Keeper) shouldRecur(ctx sdk.Context, flow types.FlowInfo, errorString string) bool {
	if strings.Contains(errorString, types.ErrBalanceTooLow) {
		return false
	}

	// Check for errors in the latest flow history entry
	hasHistoryError := false
	if k.HasFlowHistoryEntry(ctx, flow.ID) {
		history, err := k.GetFlowHistory(ctx, flow.ID)
		if err == nil && len(history) > 0 {
			latest := history[len(history)-1]
			for _, errStr := range latest.Errors {
				if errStr != "" {
					hasHistoryError = true
					break
				}
			}
		}
	}

	isRecurring := flow.ExecTime.Before(flow.EndTime) && (flow.ExecTime.Add(flow.Interval).Before(flow.EndTime) || flow.ExecTime.Add(flow.Interval).Equal(flow.EndTime))
	allowedToRecur := (!flow.Configuration.StopOnSuccess && !flow.Configuration.StopOnFailure) ||
		(flow.Configuration.StopOnSuccess && (errorString != "" || hasHistoryError)) ||
		(flow.Configuration.StopOnFailure && (errorString == "" && !hasHistoryError))

	return isRecurring && allowedToRecur
}

// // appendToPriorFlowHistory appends results to the prior history entry for the flow
// func (k Keeper) appendToPriorFlowHistory(ctx sdk.Context, flow *types.FlowInfo, fee sdk.Coin, executedLocally bool, msgResponses []*cdctypes.Any, queryResponse string, errorString string) {
// 	// Fetch the last recorded flow history for the flow
// 	entry, found := k.getCurrentFlowHistoryEntry(ctx, flow.ID)
// 	if !found {
// 		return
// 	}
// 	// Append the new data to the existing history entry
// 	entry.ExecFee = entry.ExecFee.Add(fee)
// 	entry.Executed = entry.Executed || executedLocally
// 	if flow.Configuration.SaveResponses {
// 		entry.MsgResponses = append(entry.MsgResponses, msgResponses...)
// 		entry.QueryResponse = queryResponse
// 		if errorString != "" {
// 			entry.Errors = append(entry.Errors, errorString)
// 		}
// 	}
// 	// Update the flow history with the new appended data
// 	k.SetCurrentFlowHistoryEntry(ctx, flow.ID, entry)
// }

// emitFlowEvent creates an event for the flow execution
func emitFlowEvent(ctx sdk.Context, flow types.FlowInfo) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeFlowTriggered,
			sdk.NewAttribute(types.AttributeKeyFlowID, fmt.Sprint(flow.ID)),
			sdk.NewAttribute(types.AttributeKeyFlowOwner, flow.Owner),
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
