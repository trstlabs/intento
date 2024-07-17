package keeper

import (
	"strings"
	"time"

	errorsmod "cosmossdk.io/errors"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	icacontrollerkeeper "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/keeper"
	icacontrollertypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/types"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	"github.com/trstlabs/intento/x/intent/types"
)

func (k Keeper) TriggerAction(ctx sdk.Context, action *types.ActionInfo) (bool, []*cdctypes.Any, error) {
	// local action
	if (action.ICAConfig == nil || action.ICAConfig.ConnectionID == "") && (action.HostedConfig == nil || action.HostedConfig.HostedAddress == "") {
		txMsgs := action.GetTxMsgs(k.cdc)
		msgResponses, err := handleLocalAction(k, ctx, txMsgs, *action)
		return err == nil, msgResponses, err
	}

	connectionID := action.ICAConfig.ConnectionID
	portID := action.ICAConfig.PortID

	//get hosted account from hosted config
	if action.HostedConfig.HostedAddress != "" {
		hostedAccount := k.GetHostedAccount(ctx, action.HostedConfig.HostedAddress)
		connectionID = hostedAccount.ICAConfig.ConnectionID
		portID = hostedAccount.ICAConfig.PortID
		err := k.SendFeesToHosted(ctx, *action, hostedAccount)
		if err != nil {
			return false, nil, err
		}
	}

	//check channel is active
	channelID, found := k.icaControllerKeeper.GetActiveChannelID(ctx, connectionID, portID)
	if !found {
		return false, nil, icatypes.ErrActiveChannelNotFound
	}

	//if a message contains "ICA_ADDR" string, the ICA address for the action is retrieved and parsed
	txMsgs, err := k.parseAndSetMsgs(ctx, action, connectionID, portID)
	if err != nil {
		return false, nil, err
	}
	data, err := icatypes.SerializeCosmosTx(k.cdc, txMsgs)
	if err != nil {
		return false, nil, err
	}
	packetData := icatypes.InterchainAccountPacketData{
		Type: icatypes.EXECUTE_TX,
		Data: data,
	}

	relativeTimeoutTimestamp := uint64(time.Minute.Nanoseconds())

	msgServer := icacontrollerkeeper.NewMsgServerImpl(&k.icaControllerKeeper)
	icaMsg := icacontrollertypes.NewMsgSendTx(action.Owner, connectionID, relativeTimeoutTimestamp, packetData)

	res, err := msgServer.SendTx(ctx, icaMsg)
	if err != nil {
		return false, nil, err
	}

	k.Logger(ctx).Debug("action", "ibc_sequence", res.Sequence)
	k.setTmpActionID(ctx, action.ID, portID, channelID, res.Sequence)
	return false, nil, nil
}

func handleLocalAction(k Keeper, ctx sdk.Context, txMsgs []sdk.Msg, action types.ActionInfo) ([]*cdctypes.Any, error) {
	// CacheContext returns a new context with the multi-store branched into a cached storage object
	// writeCache is called only if all msgs succeed, performing state transitions atomically
	var msgResponses []*cdctypes.Any

	cacheCtx, writeCache := ctx.CacheContext()
	for index, msg := range txMsgs {
		if action.Msgs[index].TypeUrl == "/ibc.applications.transfer.v1.MsgTransfer" {
			transferMsg, err := types.GetTransferMsg(k.cdc, action.Msgs[index])
			if err != nil {
				return nil, err
			}
			_, err = k.transferKeeper.Transfer(sdk.WrapSDKContext(ctx), &transferMsg)
			if err != nil {
				return nil, err
			}
			continue
		}

		handler := k.msgRouter.Handler(msg)
		for _, acct := range msg.GetSigners() {
			if acct.String() != action.Owner {
				return nil, errorsmod.Wrap(types.ErrUnauthorized, "owner doesn't have permission to send this message")
			}
		}

		res, err := handler(cacheCtx, msg)
		if err != nil {
			return nil, err
		}

		msgResponses = append(msgResponses, res.MsgResponses...)
		//autocompound example
		// if sdk.MsgTypeURL(msg) == "/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward" {
		// 	validator := ""
		// 	amount := sdk.NewCoin(types.Denom, sdk.ZeroInt())
		// 	for _, ev := range res.Events {
		// 		if ev.Type == distrtypes.EventTypeWithdrawRewards {
		// 			for _, attr := range ev.Attributes {
		// 				if string(attr.Key) == distrtypes.AttributeKeyValidator {
		// 					validator = string(attr.Value)
		// 				}
		// 				if string(attr.Key) == sdk.AttributeKeyAmount {
		// 					amount, err = sdk.ParseCoinNormalized(string(attr.Value))
		// 					if err != nil {
		// 						return nil, err
		// 					}
		// 				}
		// 			}

		// 			msgDelegate := stakingtypes.MsgDelegate{DelegatorAddress: action.Owner, ValidatorAddress: validator, Amount: amount}
		// 			handler := k.msgRouter.Handler(&msgDelegate)
		// 			_, err = handler(cacheCtx, &msgDelegate)
		// 			if err != nil {
		// 				return nil, err
		// 			}
		// 		}
		// 	}

		// }

	}
	writeCache()
	if !action.Configuration.SaveMsgResponses {
		msgResponses = nil
	}
	return msgResponses, nil
}

// HandleResponseAndSetActionResult sets the result of the last executed ID set at SendAction.
func (k Keeper) HandleResponseAndSetActionResult(ctx sdk.Context, portID string, channelID string, rewardType int, seq uint64, msgResponses []*cdctypes.Any) error {
	id := k.getTmpActionID(ctx, portID, channelID, seq)
	if id <= 0 {
		return nil
	}

	k.Logger(ctx).Debug("action", "executed", "on host")

	action := k.GetActionInfo(ctx, id)

	k.UpdateActionIbcUsage(ctx, action)
	owner, err := sdk.AccAddressFromBech32(action.Owner)
	if err != nil {
		return err
	}
	//airdrop reward hooks
	if rewardType == 3 {
		k.hooks.AfterActionAuthz(ctx, owner)
	} else if rewardType == 1 {
		k.hooks.AfterActionWasm(ctx, owner)
	}

	actionHistoryEntry, newErr := k.GetLatestActionHistoryEntry(ctx, id)
	if newErr != nil {
		actionHistoryEntry.Errors = append(actionHistoryEntry.Errors, newErr.Error())
	}

	actionHistoryEntry.Executed = true

	if action.Configuration.SaveMsgResponses {
		actionHistoryEntry.MsgResponses = msgResponses
	}

	//trigger remaining executions
	if action.Conditions != nil && action.Conditions.UseResponseValue != nil && action.Conditions.UseResponseValue.MsgsIndex != 0 {
		if len(msgResponses) > int(action.Conditions.UseResponseValue.MsgsIndex) && strings.Contains(msgResponses[action.Conditions.UseResponseValue.MsgsIndex].TypeUrl, action.Msgs[action.Conditions.UseResponseValue.MsgsIndex].TypeUrl) {
			tmpAction := action
			tmpAction.Msgs = action.Msgs[action.Conditions.UseResponseValue.MsgsIndex+1:]
			_, _, err = k.TriggerAction(ctx, &tmpAction)
			if err != nil {
				actionHistoryEntry.Errors = append(actionHistoryEntry.Errors, types.ErrSettingActionResult+err.Error())
			}
		}
	}
	k.SetActionInfo(ctx, &action)

	k.SetCurrentActionHistoryEntry(ctx, action.ID, actionHistoryEntry)

	return nil
}

// SetActionOnTimeout sets the action timeout result to the action

func (k Keeper) SetActionOnTimeout(ctx sdk.Context, sourcePort string, channelID string, seq uint64) error {
	id := k.getTmpActionID(ctx, sourcePort, channelID, seq)
	if id <= 0 {
		return nil
	}
	action := k.GetActionInfo(ctx, id)
	if action.Configuration.ReregisterICAAfterTimeout {
		action := k.GetActionInfo(ctx, id)
		metadataString := icatypes.NewDefaultMetadataString(action.ICAConfig.ConnectionID, action.ICAConfig.HostConnectionID)
		err := k.RegisterInterchainAccount(ctx, action.ICAConfig.ConnectionID, action.Owner, metadataString)
		if err != nil {
			return err
		}
	} else {
		k.RemoveFromActionQueue(ctx, action)
	}
	k.Logger(ctx).Debug("action packet timed out", "action_id", id)

	actionHistoryEntry, err := k.GetLatestActionHistoryEntry(ctx, id)
	if err != nil {
		return err
	}

	actionHistoryEntry.TimedOut = true
	k.SetCurrentActionHistoryEntry(ctx, id, actionHistoryEntry)

	return nil
}

// SetActionOnTimeout sets the action timeout result to the action
func (k Keeper) SetActionError(ctx sdk.Context, sourcePort string, channelID string, seq uint64, err string) {
	id := k.getTmpActionID(ctx, sourcePort, channelID, seq)
	if id <= 0 {
		return
	}

	k.Logger(ctx).Debug("action", "id", id, "error", err)

	actionHistoryEntry, newErr := k.GetLatestActionHistoryEntry(ctx, id)
	if newErr != nil {
		actionHistoryEntry.Errors = append(actionHistoryEntry.Errors, newErr.Error())
	}

	actionHistoryEntry.Errors = append(actionHistoryEntry.Errors, err)
	k.SetCurrentActionHistoryEntry(ctx, id, actionHistoryEntry)
}

// AllowedToExecute checks if execution conditons are met, e.g. if dependent transactions have executed on the host chain
// insert the next entry when execution has not happend yet
func (k Keeper) AllowedToExecute(ctx sdk.Context, action types.ActionInfo) (bool, error) {
	allowedToExecute := true
	shouldRecur := action.ExecTime.Before(action.EndTime) && action.ExecTime.Add(action.Interval).Before(action.EndTime)
	conditions := action.Conditions
	if conditions == nil {
		return true, nil
	}
	if conditions.ResponseComparison != nil {
		history, err := k.GetActionHistory(ctx, conditions.ResponseComparison.ActionID)
		if err != nil {
			allowedToExecute = false
			return allowedToExecute, err
		}
		responses := history.History[len(history.History)-1].MsgResponses
		isTrue, err := k.CompareResponseValue(ctx, action.ID, responses, *conditions.ResponseComparison)
		if !isTrue {
			allowedToExecute = false
			return allowedToExecute, err
		}
	}
	//check if dependent tx executions succeeded
	for _, actionId := range conditions.StopOnSuccessOf {
		history, err := k.GetActionHistory(ctx, actionId)
		if err != nil {
			allowedToExecute = false
		}
		if len(history.History) != 0 {
			success := history.History[len(history.History)-1].Executed && history.History[len(history.History)-1].Errors != nil
			if !success {
				allowedToExecute = false
				shouldRecur = false
			}
		}
	}

	//check if dependent tx executions failed
	for _, actionId := range conditions.StopOnFailureOf {
		history, err := k.GetActionHistory(ctx, actionId)
		if err != nil {
			allowedToExecute = false
		}
		if len(history.History) != 0 {
			success := history.History[len(history.History)-1].Executed && history.History[len(history.History)-1].Errors != nil
			if success {
				allowedToExecute = false
				shouldRecur = false
			}
		}
	}

	//check if dependent tx executions succeeded
	for _, actionId := range conditions.SkipOnFailureOf {
		history, err := k.GetActionHistory(ctx, actionId)
		if err != nil {
			allowedToExecute = false
		}
		if len(history.History) != 0 {
			success := history.History[len(history.History)-1].Executed && history.History[len(history.History)-1].Errors != nil
			if !success {
				allowedToExecute = false
			}
		}
	}

	//check if dependent tx executions failed
	for _, actionId := range conditions.SkipOnSuccessOf {
		history, err := k.GetActionHistory(ctx, actionId)
		if err != nil {
			allowedToExecute = false
		}
		if len(history.History) != 0 {
			success := history.History[len(history.History)-1].Executed && history.History[len(history.History)-1].Errors != nil
			if success {
				allowedToExecute = false
			}
		}
	}

	//if not allowed to execute, remove entry
	if !allowedToExecute {
		k.RemoveFromActionQueue(ctx, action)
		//insert the next entry given a recurring tx
		if shouldRecur {
			// adding next execTime and a new entry into the queue based on interval
			k.InsertActionQueue(ctx, action.ID, action.ExecTime.Add(action.Interval))
		}
	}
	return allowedToExecute, nil

}
