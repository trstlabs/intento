package keeper

import (
	"fmt"
	"strings"
	"time"

	errorsmod "cosmossdk.io/errors"
	storetypes "cosmossdk.io/store/types"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authztypes "github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/cosmos/gogoproto/proto"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller/keeper"
	icacontrollertypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller/types"
	icatypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/types"
	"github.com/trstlabs/intento/x/intent/msg_registry"
	"github.com/trstlabs/intento/x/intent/types"
)

func (k Keeper) TriggerFlow(ctx sdk.Context, flow *types.FlowInfo) (bool, []*cdctypes.Any, error) {
	// local flow
	if (flow.ICAConfig == nil || flow.ICAConfig.ConnectionID == "") && (flow.HostedICAConfig == nil || flow.HostedICAConfig.HostedAddress == "") {
		txMsgs := flow.GetTxMsgs(k.cdc)
		msgResponses, err := handleLocalFlow(k, ctx, txMsgs, *flow)
		return err == nil, msgResponses, errorsmod.Wrap(err, "could execute local flow")
	}

	connectionID := flow.ICAConfig.ConnectionID
	portID := flow.ICAConfig.PortID
	triggerAddress := flow.Owner
	//get hosted account from hosted config
	if flow.HostedICAConfig != nil && flow.HostedICAConfig.HostedAddress != "" {
		hostedAccount := k.GetHostedAccount(ctx, flow.HostedICAConfig.HostedAddress)
		if hostedAccount.HostedAddress == "" || hostedAccount.ICAConfig == nil {
			return false, nil, errorsmod.Wrapf(types.ErrInvalidHostedAccount, "hosted account or ICAConfig is nil for address %s", flow.HostedICAConfig.HostedAddress)
		}
		connectionID = hostedAccount.ICAConfig.ConnectionID
		portID = hostedAccount.ICAConfig.PortID
		triggerAddress = hostedAccount.HostedAddress
		err := k.SendFeesToHostedAdmin(ctx, *flow, hostedAccount)
		if err != nil {
			return false, nil, errorsmod.Wrap(err, "could not pay hosted account")
		}

	}

	//check channel is active
	channelID, found := k.icaControllerKeeper.GetActiveChannelID(ctx, connectionID, portID)
	if !found {
		return false, nil, icatypes.ErrActiveChannelNotFound
	}

	//if a message contains "ICA_ADDR" string, the ICA address for the flow is retrieved and parsed
	txMsgs, err := k.parseAndSetMsgs(ctx, flow, connectionID, portID)
	if err != nil {
		return false, nil, errorsmod.Wrap(err, "could parse and set messages")
	}
	data, err := icatypes.SerializeCosmosTx(k.cdc, txMsgs, icatypes.EncodingProtobuf)
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

	k.Logger(ctx).Debug("flow", "ibc_sequence", res.Sequence)
	k.SetTmpFlowID(ctx, flow.ID, portID, channelID, res.Sequence)
	return false, nil, nil
}

func handleLocalFlow(k Keeper, ctx sdk.Context, txMsgs []sdk.Msg, flow types.FlowInfo) ([]*cdctypes.Any, error) {
	// CacheContext returns a new context with the multi-store branched into a cached storage object
	// writeCache is called only if all msgs succeed, performing state transitions atomically
	var msgResponses []*cdctypes.Any

	cacheCtx, writeCache := ctx.CacheContext()
	for index, msg := range txMsgs {

		signers, _, err := k.cdc.GetMsgV1Signers(msg)
		if err != nil {
			return nil, err
		}
		for _, acct := range signers {
			if sdk.AccAddress(acct).String() != flow.Owner {
				return nil, errorsmod.Wrap(types.ErrUnauthorized, "owner doesn't have permission to send this message")
			}
		}

		if flow.Msgs[index].TypeUrl == "/ibc.applications.transfer.v1.MsgTransfer" {
			transferMsg, err := types.GetTransferMsg(k.cdc, flow.Msgs[index])
			if err != nil {
				return nil, err
			}
			_, err = k.transferKeeper.Transfer(ctx, &transferMsg)
			if err != nil {
				return nil, err
			}
			continue
		}

		handler := k.msgRouter.Handler(msg)
		res, err := handler(cacheCtx, msg)
		if err != nil {
			return nil, err
		}

		msgResponses = append(msgResponses, res.MsgResponses...)

	}
	writeCache()
	if !flow.Configuration.SaveResponses {
		msgResponses = nil
	}
	return msgResponses, nil
}

// HandleResponseAndSetFlowResult sets the result of the last executed ID set at SendFlow.
func (k Keeper) HandleResponseAndSetFlowResult(ctx sdk.Context, portID string, channelID string, relayer sdk.AccAddress, seq uint64, msgResponses []*cdctypes.Any) error {
	k.Logger(ctx).Debug("HandleResponseAndSetFlowResult: portID=%s, channelID=%s, seq=%d\n", portID, channelID, seq)
	id := k.getTmpFlowID(ctx, portID, channelID, seq)
	if id <= 0 {
		return fmt.Errorf("flow not found")
	}
	flow := k.GetFlowInfo(ctx, id)

	flowHistoryEntry, newErr := k.GetLatestFlowHistoryEntry(ctx, id)
	if newErr != nil {
		flowHistoryEntry.Errors = append(flowHistoryEntry.Errors, newErr.Error())
	}

	msgResponses, msgClass, err := k.HandleDeepResponses(ctx, msgResponses, relayer, flow, len(flowHistoryEntry.MsgResponses))
	if err != nil {
		return err
	}

	owner, err := sdk.AccAddressFromBech32(flow.Owner)
	if err != nil {
		return err
	}
	// reward hooks
	if msgClass == 3 {
		k.hooks.AfterActionLocal(ctx, owner)
	} else if msgClass == 1 {
		k.hooks.AfterActionICA(ctx, owner)
	}

	flowHistoryEntry.Executed = true

	if flow.Configuration.SaveResponses {
		flowHistoryEntry.MsgResponses = append(flowHistoryEntry.MsgResponses, msgResponses...)
	}
	k.SetCurrentFlowHistoryEntry(ctx, flow.ID, flowHistoryEntry)
	if len(flow.Conditions.FeedbackLoops) != 0 {
		// Only process feedback loops that match the current response index
		for i, feedbackLoop := range flow.Conditions.FeedbackLoops {
			k.Logger(ctx).Debug("checking feedback loop",
				"index", i,
				"responseIndex", feedbackLoop.ResponseIndex,
				"msgIndex", feedbackLoop.MsgsIndex,
			)

			// Skip if this isn't a feedback loop for the current response
			if int(feedbackLoop.ResponseIndex) != len(flowHistoryEntry.MsgResponses)-1 {
				k.Logger(ctx).Debug("skipping feedback loop - wrong response index",
					"expected", len(flowHistoryEntry.MsgResponses)-1,
					"got", feedbackLoop.ResponseIndex,
				)
				continue
			}

			// Validate MsgsIndex and FlowID
			if feedbackLoop.MsgsIndex == 0 || feedbackLoop.FlowID != 0 {
				k.Logger(ctx).Debug("skipping feedback loop - invalid MsgsIndex or FlowID")
				continue // Skip invalid FeedbackLoops or if FlowID is set to non-default
			}

			// Only process messages that come after the current message
			if int(feedbackLoop.MsgsIndex) >= len(flow.Msgs) {
				k.Logger(ctx).Debug("skipping feedback loop - MsgsIndex out of bounds")
				continue // Skip if MsgsIndex is out of bounds
			}

			// Only include messages from the target message index onwards
			// tmpFlowMsgs := make([]*cdctypes.Any, len(flow.Msgs[feedbackLoop.MsgsIndex:]))
			// Find the next message that needs a response for subsequent feedback
			nextStopIndex := len(flow.Msgs)
			for _, nextLoop := range flow.Conditions.FeedbackLoops {
				if int(nextLoop.MsgsIndex) > int(feedbackLoop.MsgsIndex) &&
					int(nextLoop.MsgsIndex) < nextStopIndex {
					nextStopIndex = int(nextLoop.MsgsIndex)
				}
			}

			// Execute only up to the next feedback point
			messagesToExecute := flow.Msgs[feedbackLoop.MsgsIndex:nextStopIndex]

			k.Logger(ctx).Debug("processing feedback loop",
				"responseIndex", feedbackLoop.ResponseIndex,
				"msgIndex", feedbackLoop.MsgsIndex,
				"nextMsgs", len(messagesToExecute),
			)

			if err := executeMessageBatch(k, ctx, flow, messagesToExecute, flowHistoryEntry); err != nil {
				return err // Return on the first encountered error
			}
		}
	}

	return nil

}

func executeMessageBatch(k Keeper, ctx sdk.Context, flow types.FlowInfo, nextMsgs []*cdctypes.Any, flowHistoryEntry *types.FlowHistoryEntry) error {
	var errorString = ""

	allowed, err := k.allowedToExecute(ctx, flow)
	if !allowed {
		k.recordFlowNotAllowed(ctx, &flow, ctx.BlockTime(), err)
	}

	flowCtx := ctx.WithGasMeter(storetypes.NewGasMeter(types.MaxGas))
	cacheCtx, writeCtx := flowCtx.CacheContext()
	k.Logger(ctx).Debug("continuing msg execution", "id", flow.ID, "next_msgs", len(nextMsgs))

	// Only run feedback loops if we have conditions and messages to process
	if flow.Conditions != nil && len(nextMsgs) > 0 {
		err = k.RunFeedbackLoops(cacheCtx, flow.ID, &flow.Msgs, flow.Conditions)
		if err != nil {
			k.Logger(ctx).Error("error running feedback loops", "error", err)
			return errorsmod.Wrap(err, fmt.Sprintf(types.ErrSettingFlowResult, err))
		}

		// After running feedback loops, check if we still have messages to process
		if len(nextMsgs) == 0 {
			k.Logger(ctx).Debug("no messages to process after feedback loops")
			return nil
		}

		k.Logger(ctx).Debug("triggering next msgs", "id", flow.ID, "msgs", len(nextMsgs))
		_, _, err = k.TriggerFlow(cacheCtx, &flow)
		if err != nil {
			errorString = appendError(errorString, fmt.Sprintf(types.ErrFlowMsgHandling, err.Error()))
		}

		// Only try to distribute fees if we have a valid fee address and denom
		feeAddr, feeDenom, feeErr := k.GetFeeAccountForMinFees(cacheCtx, flow, types.MaxGas)
		if feeErr != nil {
			errorString = appendError(errorString, feeErr.Error())
		} else if feeAddr == nil || feeDenom == "" || flowHistoryEntry.ExecFee.IsZero() || flowHistoryEntry.ExecFee.Denom != feeDenom {
			fmt.Print("feeDenom", feeDenom, "flowHistoryEntry.ExecFee", flowHistoryEntry.ExecFee)
			errorString = appendError(errorString, types.ErrBalanceTooLow)
		} else {
			k.Logger(ctx).Debug("feeDenom", feeDenom, "flowHistoryEntry.ExecFee", flowHistoryEntry.ExecFee)
			fee, distErr := k.DistributeCoins(cacheCtx, flow, feeAddr, feeDenom)
			if distErr != nil {
				errorString = appendError(errorString, fmt.Sprintf(types.ErrFlowFeeDistribution, distErr.Error()))
			} else {
				flowHistoryEntry.ExecFee = flowHistoryEntry.ExecFee.Add(fee)
			}
		}

		if errorString != "" {
			flowHistoryEntry.Executed = false
			flowHistoryEntry.Errors = append(flowHistoryEntry.Errors, errorString)
		}

		k.SetCurrentFlowHistoryEntry(cacheCtx, flow.ID, flowHistoryEntry)
		k.SetFlowInfo(ctx, &flow)
		writeCtx()
	}

	return nil
}

func (k Keeper) HandleDeepResponses(ctx sdk.Context, msgResponses []*cdctypes.Any, relayer sdk.AccAddress, flow types.FlowInfo, previousMsgsExecuted int) ([]*cdctypes.Any, int, error) {
	var msgClass int

	for index, anyResp := range msgResponses {
		k.Logger(ctx).Debug("msg response in ICS-27 packet", "response", anyResp.GoString(), "typeURL", anyResp.GetTypeUrl())

		if entry, ok := msg_registry.MsgRegistry[anyResp.GetTypeUrl()]; ok {
			if index == 0 && entry.RewardType > 0 {
				msgClass = entry.RewardType
				k.HandleRelayerReward(ctx, relayer, msgClass)
			}
		}
		if anyResp.GetTypeUrl() == "/cosmos.authz.v1beta1.MsgExecResponse" {

			msgExecResponse := authztypes.MsgExecResponse{}
			err := proto.Unmarshal(anyResp.GetValue(), &msgExecResponse)
			if err != nil {
				k.Logger(ctx).Debug("handling deep flow response unmarshalling", "err", err)
				return nil, 0, err
			}
			flowIndex := index + previousMsgsExecuted
			if flowIndex >= len(flow.Msgs) {
				return nil, 0, errorsmod.Wrapf(types.ErrMsgResponsesHandling, "expected more message responses")
			}
			msgExec := &authztypes.MsgExec{}
			if err := proto.Unmarshal(flow.Msgs[flowIndex].Value, msgExec); err != nil {
				return nil, 0, errorsmod.Wrapf(types.ErrMsgResponsesHandling, "failed to unmarshal MsgExec: %s", err.Error())
			}

			msgResponses = []*cdctypes.Any{}

			for _, resultBytes := range msgExecResponse.Results {
				if len(resultBytes) == 0 {
					continue // Skip empty results
				}

				var msgResponse = cdctypes.Any{}
				if err := proto.Unmarshal(resultBytes, &msgResponse); err == nil {
					typeUrl := msgResponse.GetTypeUrl()

					if typeUrl != "" && strings.Contains(typeUrl, "Msg") {
						_, err := k.interfaceRegistry.Resolve(typeUrl)
						if err == nil {
							k.Logger(ctx).Debug("parsing response authz v0.52+", "msgResponse", msgResponse)
							msgResponses = append(msgResponses, &msgResponse)
							continue
						}
					}

				}
				// in v0.50.8 we were writing msgResponse.Data in that [][]byte and no marshalled anys
				// https://github.com/cosmos/cosmos-sdk/blob/v0.50.8/x/authz/keeper/keeper.go#L166-L186
				//	k.Logger(ctx).Debug("flow result", "resultBytes", resultBytes)
				//as we do not get typeURL (until cosmos 0.52 and is not possible in 51) we have to rely in this, MsgExec is the only regisrered message that should return results
				msgRespProto, _, err := handleMsgData(&sdk.MsgData{Data: resultBytes, MsgType: msgExec.Msgs[0].TypeUrl})
				if err != nil {
					return nil, 0, err
				}
				// Only create Any if we have a non-nil response
				if msgRespProto != nil {
					respAny, err := cdctypes.NewAnyWithValue(msgRespProto)
					if err != nil {
						return nil, 0, err
					}
					msgResponses = append(msgResponses, respAny)
				}
			}
		}
	}
	return msgResponses, msgClass, nil
}

// SetFlowOnTimeout sets the flow timeout result to the flow

func (k Keeper) SetFlowOnTimeout(ctx sdk.Context, sourcePort string, channelID string, seq uint64) error {
	id := k.getTmpFlowID(ctx, sourcePort, channelID, seq)
	if id <= 0 {
		return nil
	}
	flow := k.GetFlowInfo(ctx, id)
	if flow.Configuration.StopOnTimeout {
		k.RemoveFromFlowQueue(ctx, flow)
	}
	k.Logger(ctx).Debug("flow packet timed out", "flow_id", id)

	flowHistoryEntry, err := k.GetLatestFlowHistoryEntry(ctx, id)
	if err != nil {
		return err
	}

	flowHistoryEntry.TimedOut = true
	flowHistoryEntry.Executed = false
	k.SetCurrentFlowHistoryEntry(ctx, id, flowHistoryEntry)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeFlowTimedOut,
			sdk.NewAttribute(types.AttributeKeyFlowID, fmt.Sprint(flow.ID)),
			sdk.NewAttribute(types.AttributeKeyFlowOwner, flow.Owner),
		),
	)
	return nil
}

// SetFlowError sets the flow error result and emits an error event
func (k Keeper) SetFlowError(ctx sdk.Context, sourcePort, channelID string, seq uint64, errStr string) {
	id := k.getTmpFlowID(ctx, sourcePort, channelID, seq)
	if id == 0 {
		return
	}

	k.Logger(ctx).Debug("flow", "id", id, "error", errStr)

	flowHistoryEntry, err := k.GetLatestFlowHistoryEntry(ctx, id)
	if err != nil {
		flowHistoryEntry = &types.FlowHistoryEntry{Errors: []string{err.Error()}}
	}

	flow, err := k.TryGetFlowInfo(ctx, id)
	if err != nil {
		flowHistoryEntry.Errors = append(flowHistoryEntry.Errors, err.Error())
	}

	if errStr != "" {
		flowHistoryEntry.Errors = append(flowHistoryEntry.Errors, errStr)
		flowHistoryEntry.Executed = false
		k.SetCurrentFlowHistoryEntry(ctx, id, flowHistoryEntry)

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeFlowError,
				sdk.NewAttribute(types.AttributeKeyFlowID, fmt.Sprint(id)),
				sdk.NewAttribute(types.AttributeKeyFlowOwner, flow.Owner),
				sdk.NewAttribute(types.AttributeKeyError, errStr),
			),
		)
	}
}
