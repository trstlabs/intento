package keeper

import (
	"encoding/base64"
	"fmt"
	"strconv"
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

	allowed, err := k.allowedToExecute(ctx, action, queryCallback)
	if !allowed {
		k.recordActionNotAllowed(ctx, &action, timeOfBlock, fmt.Sprintf(types.ErrActionConditions, err.Error()))
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
		executedLocally, errorString = k.handleActionExecution(cacheCtx, &action, &msgResponses, errorString, queryCallback)

		fee, err = k.DistributeCoins(cacheCtx, action, feeAddr, feeDenom, ctx.BlockHeader().ProposerAddress)
		if err != nil {
			errorString = appendError(errorString, fmt.Sprintf(types.ErrActionFeeDistribution, err.Error()))
		}

	}

	k.addActionHistory(cacheCtx, &action, timeOfBlock, fee, executedLocally, msgResponses, string(queryCallback), errorString)

	writeCtx()

	if shouldRecur(action, errorString) {
		action.ExecTime = action.ExecTime.Add(action.Interval)
		k.InsertActionQueue(ctx, action.ID, action.ExecTime)
	}

	emitActionEvent(ctx, action)
	k.SetActionInfo(ctx, &action)
}

// submitInterchainQuery submits an interchain query when ICQConfig is present
func (k Keeper) SubmitInterchainQuery(ctx sdk.Context, action types.ActionInfo, logger log.Logger) {
	_, found := k.interchainQueryKeeper.GetQuery(ctx, strconv.FormatUint(action.ID, 10))
	if found {
		return
	}
	requestData, err := base64.StdEncoding.DecodeString(action.Conditions.ICQConfig.QueryKey)
	if err != nil {
		k.SetActionHistoryEntry(ctx, action.ID, &types.ActionHistoryEntry{Errors: []string{fmt.Sprint("Error submitting ICQ: decoding Base64 string: ", err)}})
		k.RemoveFromActionQueue(ctx, action)
		return
	}
	id := strconv.FormatUint(action.ID, 10)
	// Submit the interchain query using the configuration provided in action.Conditions.ICQConfig
	query := icqtypes.Query{
		Id:              id,
		CallbackId:      ICQCallbackID_Action,
		ConnectionId:    action.Conditions.ICQConfig.ConnectionId,
		ChainId:         action.Conditions.ICQConfig.ChainId,
		TimeoutPolicy:   action.Conditions.ICQConfig.TimeoutPolicy,
		TimeoutDuration: action.Conditions.ICQConfig.TimeoutDuration,
		QueryType:       action.Conditions.ICQConfig.QueryType,
		RequestData:     requestData,
		CallbackModule:  types.ModuleName,
	}
	k.interchainQueryKeeper.SetQuery(ctx, query)
	// Log successful submission of the interchain query
	logger.Debug("interchain query submitted", "actionID", action.ID, "ICQ ID", id)
}

// handleActionExecution handles the core logic of triggering an action and processing responses
func (k Keeper) handleActionExecution(ctx sdk.Context, action *types.ActionInfo, msgResponses *[]*cdctypes.Any, errorString string, queryCallback []byte) (bool, string) {
	var executedLocally bool
	// Safe check to ensure conditions exist before proceeding
	if action.Conditions == nil || action.Conditions.UseResponseValue == nil {
		// If no UseResponseValue conditions are present, just execute all Msgs normally
		executedLocally, responses, err := k.TriggerAction(ctx, action)
		if err != nil {
			return false, appendError(errorString, fmt.Sprintf(types.ErrActionMsgHandling, err.Error()))
		}
		*msgResponses = append(*msgResponses, responses...)
		return executedLocally, errorString
	}

	// UseResponseValue is set, handle accordingly
	executedLocally, errorString = k.handleUseResponseValue(ctx, action, msgResponses, errorString, queryCallback)

	return executedLocally, errorString
}

// handleUseResponseValue processes the UseResponseValue logic when it's set
func (k Keeper) handleUseResponseValue(ctx sdk.Context, action *types.ActionInfo, msgResponses *[]*cdctypes.Any, errorString string, queryCallback []byte) (bool, string) {
	actionTmp := *action

	// Ensure MsgsIndex is valid and avoid out-of-bound slices
	if action.Conditions.UseResponseValue.MsgsIndex > 0 && int(action.Conditions.UseResponseValue.MsgsIndex) <= len(action.Msgs) {
		// Split action.Msgs based on MsgsIndex (first part)
		actionTmp.Msgs = action.Msgs[:action.Conditions.UseResponseValue.MsgsIndex]
	}

	// Trigger the first set of Msgs (before MsgsIndex)
	executedLocally, responses, err := k.TriggerAction(ctx, &actionTmp)
	if err != nil {
		return false, appendError(errorString, fmt.Sprintf(types.ErrActionMsgHandling, err.Error()))
	}
	*msgResponses = append(*msgResponses, responses...)

	// If ActionID is set, attempt to process its value based on the response of the first set
	if action.Conditions.UseResponseValue.ActionID != 0 || executedLocally || queryCallback != nil {
		err = k.UseResponseValue(ctx, action.ID, &actionTmp.Msgs, action.Conditions, queryCallback)
		if err != nil {
			return false, appendError(errorString, fmt.Sprintf(types.ErrSettingActionResult, err.Error()))
		}

		// Process the remaining Msgs (after MsgsIndex)
		if action.Conditions.UseResponseValue.MsgsIndex > 0 && int(action.Conditions.UseResponseValue.MsgsIndex) < len(action.Msgs) {
			// Handle the remaining Msgs
			actionTmp.Msgs = action.Msgs[action.Conditions.UseResponseValue.MsgsIndex:]

			// Trigger the remaining Msgs
			_, additionalResponses, err := k.TriggerAction(ctx, &actionTmp)
			if err != nil {
				return false, appendError(errorString, fmt.Sprintf(types.ErrActionMsgHandling, err.Error()))
			}
			*msgResponses = append(*msgResponses, additionalResponses...)
		}
	}

	return executedLocally, errorString
}

// recordActionNotAllowed adds an action entry to the action history
func (k Keeper) recordActionNotAllowed(ctx sdk.Context, action *types.ActionInfo, timeOfBlock time.Time, errorMsg string) {
	k.Logger(ctx).Debug("action not allowed to execute", "id", action.ID)
	k.addActionHistory(ctx, action, timeOfBlock, sdk.Coin{}, false, nil, "", errorMsg)
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
