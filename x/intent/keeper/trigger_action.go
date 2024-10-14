package keeper

import (
	"fmt"
	"time"

	errorsmod "cosmossdk.io/errors"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authztypes "github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/cosmos/gogoproto/proto"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/keeper"
	icacontrollertypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/types"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	"github.com/trstlabs/intento/x/intent/types"
	icqtypes "github.com/trstlabs/intento/x/interchainquery/types"
)

func (k Keeper) TriggerAction(ctx sdk.Context, action *types.ActionInfo) (bool, []*cdctypes.Any, error) {
	// local action
	if (action.ICAConfig == nil || action.ICAConfig.ConnectionID == "") && (action.HostedConfig == nil || action.HostedConfig.HostedAddress == "") {
		txMsgs := action.GetTxMsgs(k.cdc)
		msgResponses, err := handleLocalAction(k, ctx, txMsgs, *action)
		return err == nil, msgResponses, errorsmod.Wrap(err, "could execute local action")
	}

	connectionID := action.ICAConfig.ConnectionID
	portID := action.ICAConfig.PortID
	triggerAddress := action.Owner
	//get hosted account from hosted config
	if action.HostedConfig != nil && action.HostedConfig.HostedAddress != "" {
		hostedAccount := k.GetHostedAccount(ctx, action.HostedConfig.HostedAddress)
		connectionID = hostedAccount.ICAConfig.ConnectionID
		portID = hostedAccount.ICAConfig.PortID
		triggerAddress = hostedAccount.HostedAddress
		err := k.SendFeesToHosted(ctx, *action, hostedAccount)
		if err != nil {
			return false, nil, errorsmod.Wrap(err, "could not pay hosted account")
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
		return false, nil, errorsmod.Wrap(err, "could parse and set messages")
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
	icaMsg := icacontrollertypes.NewMsgSendTx(triggerAddress, connectionID, relativeTimeoutTimestamp, packetData)

	res, err := msgServer.SendTx(ctx, icaMsg)
	if err != nil {
		return false, nil, errorsmod.Wrap(err, "could not send ICA message")
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

	}
	writeCache()
	if !action.Configuration.SaveResponses {
		msgResponses = nil
	}
	return msgResponses, nil
}

// HandleResponseAndSetActionResult sets the result of the last executed ID set at SendAction.
func (k Keeper) HandleResponseAndSetActionResult(ctx sdk.Context, portID string, channelID string, relayer sdk.AccAddress, seq uint64, msgResponses []*cdctypes.Any) error {
	id := k.getTmpActionID(ctx, portID, channelID, seq)
	if id <= 0 {
		return nil
	}
	action := k.GetActionInfo(ctx, id)

	actionHistoryEntry, newErr := k.GetLatestActionHistoryEntry(ctx, id)
	if newErr != nil {
		actionHistoryEntry.Errors = append(actionHistoryEntry.Errors, newErr.Error())
	}

	msgResponses, msgClass, err := k.HandleDeepResponses(ctx, msgResponses, relayer, action, len(actionHistoryEntry.MsgResponses))
	if err != nil {
		return err
	}

	k.UpdateActionIbcUsage(ctx, action)
	owner, err := sdk.AccAddressFromBech32(action.Owner)
	if err != nil {
		return err
	}
	// reward hooks
	if msgClass == 3 {
		k.hooks.AfterActionAuthz(ctx, owner)
	} else if msgClass == 1 {
		k.hooks.AfterActionWasm(ctx, owner)
	}

	actionHistoryEntry.Executed = true

	if action.Configuration.SaveResponses {
		actionHistoryEntry.MsgResponses = append(actionHistoryEntry.MsgResponses, msgResponses...)
	}

	k.SetCurrentActionHistoryEntry(ctx, action.ID, actionHistoryEntry)

	//trigger remaining execution
	if action.Conditions != nil && action.Conditions.UseResponseValue != nil && action.Conditions.UseResponseValue.MsgsIndex != 0 && action.Conditions.UseResponseValue.ActionID == 0 && len(actionHistoryEntry.MsgResponses)-1 < int(action.Conditions.UseResponseValue.MsgsIndex) {
		err = k.UseResponseValue(ctx, action.ID, &action.Msgs, action.Conditions, nil)
		if err != nil {
			return errorsmod.Wrap(err, types.ErrSettingActionResult)
		}
		tmpAction := action
		tmpAction.Msgs = action.Msgs[action.Conditions.UseResponseValue.MsgsIndex:]

		k.Logger(ctx).Debug("triggering msgs", "id", action.ID, "msgs", len(tmpAction.Msgs))
		_, _, err = k.TriggerAction(ctx, &tmpAction)
		if err != nil {
			actionHistoryEntry.Errors = append(actionHistoryEntry.Errors, types.ErrActionMsgHandling+err.Error())
			k.SetCurrentActionHistoryEntry(ctx, action.ID, actionHistoryEntry)
		}
		//}
	}
	k.SetActionInfo(ctx, &action)

	return nil
}

func (k Keeper) HandleDeepResponses(ctx sdk.Context, msgResponses []*cdctypes.Any, relayer sdk.AccAddress, action types.ActionInfo, previousMsgsExecuted int) ([]*cdctypes.Any, int, error) {
	var msgClass int

	for index, anyResp := range msgResponses {
		k.Logger(ctx).Debug("msg response in ICS-27 packet", "response", anyResp.GoString(), "typeURL", anyResp.GetTypeUrl())

		rewardClass := getMsgRewardType(anyResp.GetTypeUrl())
		if index == 0 && rewardClass > 0 {
			msgClass = rewardClass
			k.HandleRelayerReward(ctx, relayer, msgClass)
		}
		if anyResp.GetTypeUrl() == "/cosmos.authz.v1beta1.MsgExecResponse" {

			msgExecResponse := authztypes.MsgExecResponse{}
			err := proto.Unmarshal(anyResp.GetValue(), &msgExecResponse)
			if err != nil {
				fmt.Printf("err %v \n", err)
				k.Logger(ctx).Debug("action response", "err", err)
				return nil, 0, err
			}

			actionIndex := index + previousMsgsExecuted
			// if actionIndex < 0 {
			// 	actionIndex = 0
			// }
			msgExec := &authztypes.MsgExec{}
			if err := proto.Unmarshal(action.Msgs[actionIndex].Value, msgExec); err != nil {
				return nil, 0, err
			}

			msgResponses = []*cdctypes.Any{}

			for _, resultBytes := range msgExecResponse.Results {
				// var result sdk.Result
				// if err := proto.Unmarshal(resultBytes, &result); err != nil {
				// 	k.Logger(ctx).Debug("action result", "err", err)
				// 	fmt.Printf("err 2%v \n", err)
				// 	return nil, 0, err
				// }

				k.Logger(ctx).Debug("action result", "resultBytes", resultBytes)

				//as we do not get typeURL we have to rely in this, MsgExec is the only regisrered message that should return results
				msgRespProto, _, err := handleMsgData(&sdk.MsgData{Data: resultBytes, MsgType: msgExec.Msgs[0].TypeUrl})
				if err != nil {
					fmt.Printf("err 3%v \n", err)
					return nil, 0, err
				}
				respAny, err := cdctypes.NewAnyWithValue(msgRespProto)
				if err != nil {
					fmt.Printf("err4 %v \n", err)
					return nil, 0, err
				}

				msgResponses = append(msgResponses, respAny)

			}
		}
	}
	return msgResponses, msgClass, nil
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

// allowedToExecute checks if execution conditons are met, e.g. if dependent transactions have executed on the host chain
// insert the next entry when execution has not happend yet
func (k Keeper) allowedToExecute(ctx sdk.Context, action types.ActionInfo, queryCallback *icqtypes.Query) (bool, error) {
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
		responses := history[len(history)-1].MsgResponses
		isTrue, err := k.CompareResponseValue(ctx, action.ID, responses, *conditions.ResponseComparison, queryCallback)
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
		if len(history) != 0 {
			success := history[len(history)-1].Executed && history[len(history)-1].Errors != nil
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
		if len(history) != 0 {
			success := history[len(history)-1].Executed && history[len(history)-1].Errors != nil
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
		if len(history) != 0 {
			success := history[len(history)-1].Executed && history[len(history)-1].Errors != nil
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
		if len(history) != 0 {
			success := history[len(history)-1].Executed && history[len(history)-1].Errors != nil
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
