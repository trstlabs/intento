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

func (k Keeper) TriggerFlow(ctx sdk.Context, flow *types.Flow) (int64, []*cdctypes.Any, error) {
	// local flow
	if (flow.SelfHostedICA == nil || flow.SelfHostedICA.ConnectionID == "") && (flow.TrustlessAgent == nil || flow.TrustlessAgent.AgentAddress == "") {
		txMsgs := flow.GetTxMsgs(k.cdc)
		msgResponses, err := handleLocalFlow(k, ctx, txMsgs, *flow)
		if err != nil {
			return 0, msgResponses, errorsmod.Wrap(err, "could not execute local flow")
		}
		//indicates flow was executed locally
		return -1, msgResponses, nil
	}

	connectionID := flow.SelfHostedICA.ConnectionID
	portID := flow.SelfHostedICA.PortID
	triggerAddress := flow.Owner
	//get trustless agent from hosted config
	if flow.TrustlessAgent != nil && flow.TrustlessAgent.AgentAddress != "" {
		trustlessAgent := k.GetTrustlessAgent(ctx, flow.TrustlessAgent.AgentAddress)
		if trustlessAgent.AgentAddress == "" || trustlessAgent.ICAConfig == nil {
			return 0, nil, errorsmod.Wrapf(types.ErrInvalidTrustlessAgent, "trustless agent or ICAConfig is nil for address %s", flow.TrustlessAgent.AgentAddress)
		}
		connectionID = trustlessAgent.ICAConfig.ConnectionID
		portID = trustlessAgent.ICAConfig.PortID
		triggerAddress = trustlessAgent.AgentAddress
		if trustlessAgent.FeeConfig != nil && trustlessAgent.FeeConfig.FeeCoinsSupported != nil {
			err := k.SendFeesToTrustlessAgentFeeAdmin(ctx, *flow, trustlessAgent)
			if err != nil {
				return 0, nil, errorsmod.Wrap(err, "could not pay trustless agent")
			}
		}

	}

	//check channel is active
	channelID, found := k.icaControllerKeeper.GetActiveChannelID(ctx, connectionID, portID)
	if !found {
		return 0, nil, icatypes.ErrActiveChannelNotFound
	}

	//if a message contains "ICA_ADDR" string, the ICA address for the flow is retrieved and parsed
	txMsgs, err := k.parseAndSetMsgs(ctx, flow, connectionID, portID)
	if err != nil {
		return 0, nil, errorsmod.Wrap(err, "could parse and set messages")
	}
	data, err := icatypes.SerializeCosmosTx(k.cdc, txMsgs, icatypes.EncodingProtobuf)
	if err != nil {
		return 0, nil, err
	}
	packetData := icatypes.InterchainAccountPacketData{
		Type: icatypes.EXECUTE_TX,
		Data: data,
	}
	const maxTimeout = time.Hour

	timeout := maxTimeout
	if flow.Interval > 0 && flow.Interval < maxTimeout {
		timeout = flow.Interval
	}
	relativeTimeoutTimestamp := uint64(timeout.Nanoseconds())
	msgServer := icacontrollerkeeper.NewMsgServerImpl(&k.icaControllerKeeper)
	icaMsg := icacontrollertypes.NewMsgSendTx(triggerAddress, connectionID, relativeTimeoutTimestamp, packetData)

	res, err := msgServer.SendTx(ctx, icaMsg)
	if err != nil {
		return 0, nil, errorsmod.Wrap(err, "could not send ICA message")
	}

	k.Logger(ctx).Debug("flow", "ibc_sequence", res.Sequence, "message", flow.GetTxMsgs(k.cdc)[0].String())
	k.SetTmpFlowID(ctx, flow.ID, portID, channelID, res.Sequence)

	return int64(res.Sequence), nil, nil
}

func handleLocalFlow(k Keeper, ctx sdk.Context, txMsgs []sdk.Msg, flow types.Flow) ([]*cdctypes.Any, error) {
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
		} else {
			flowOwnerAddr, err := sdk.AccAddressFromBech32(flow.Owner)
			if err != nil {
				return nil, err
			}
			k.hooks.AfterActionLocal(cacheCtx, flowOwnerAddr)
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

// HandleResponseAndSetFlowResult sets the result of the initiated interchain flow.
func (k Keeper) HandleResponseAndSetFlowResult(ctx sdk.Context, portID string, channelID string, relayer sdk.AccAddress, seq uint64, msgResponses []*cdctypes.Any) error {
	k.Logger(ctx).Debug("HandleResponseAndSetFlowResult:", "portID", portID, "channelID", channelID, "seq", seq)
	id := k.getTmpFlowID(ctx, portID, channelID, seq)
	if id <= 0 {
		return fmt.Errorf("flow not found")
	}
	flow := k.GetFlow(ctx, id)

	flowHistoryEntry, newErr := k.GetLatestFlowHistoryEntry(ctx, id)
	if newErr != nil {
		flowHistoryEntry.Errors = append(flowHistoryEntry.Errors, newErr.Error())
	}

	if len(flowHistoryEntry.MsgResponses) == len(flow.Msgs) {
		k.Logger(ctx).Debug("all messages executed, skipping deep response handling")
		return nil
	}

	msgResponses, _, err := k.HandleDeepResponses(ctx, msgResponses, relayer, flow, len(flowHistoryEntry.MsgResponses))
	if err != nil {
		return err
	}

	owner, err := sdk.AccAddressFromBech32(flow.Owner)
	if err != nil {
		return err
	}

	k.hooks.AfterActionICA(ctx, owner)

	flowHistoryEntry.Executed = true
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeFlowMsgResponse,
			sdk.NewAttribute(types.AttributeKeyFlowID, fmt.Sprint(flow.ID)),
			sdk.NewAttribute(types.AttributeKeyFlowOwner, flow.Owner),
		),
	)
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
			flowHistoryEntryCheck := flowHistoryEntry
			if feedbackLoop.FlowID != 0 {
				flowHistoryEntryCheck, err = k.GetLatestFlowHistoryEntry(ctx, feedbackLoop.FlowID)
				if err != nil {
					k.Logger(ctx).Debug("error getting latest flow history entry", "error", err)

					return errorsmod.Wrapf(err, "error getting latest flow history entry for flow ID %d", feedbackLoop.FlowID)
				}
			}
			// Skip if this isn't a feedback loop for the current response
			if int(feedbackLoop.ResponseIndex) != len(flowHistoryEntryCheck.MsgResponses)-1 {
				k.Logger(ctx).Debug("skipping feedback loop - wrong response index",
					"expected", len(flowHistoryEntry.MsgResponses)-1,
					"got", feedbackLoop.ResponseIndex,
				)
				continue
			}

			// Validate MsgsIndex and FlowID
			if feedbackLoop.MsgsIndex == 0 && feedbackLoop.FlowID == 0 {
				k.Logger(ctx).Debug("skipping feedback loop - invalid MsgsIndex or FlowID")
				continue // Skip invalid FeedbackLoops or if FlowID is set to non-default
			}

			// Only process messages that come after the current message
			if feedbackLoop.FlowID == 0 && int(feedbackLoop.MsgsIndex) >= len(flow.Msgs) {
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

			k.Logger(ctx).Debug("execute batch of messages",
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

func executeMessageBatch(k Keeper, ctx sdk.Context, flow types.Flow, nextMsgs []*cdctypes.Any, flowHistoryEntry *types.FlowHistoryEntry) error {
	var errorString = ""

	allowed, err := k.allowedToExecute(ctx, flow)
	if err != nil {
		k.Logger(ctx).Error("error allow execution at batch execution of flow", "error", err, "id", flow.ID)
		return errorsmod.Wrap(err, fmt.Sprintf(types.ErrSettingFlowResult, err))
	}
	if !allowed {
		k.Logger(ctx).Debug("flow not allowed to execute", "id", flow.ID)
		return nil
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
		// Only try to distribute fees if we have a valid fee address and denom
		feeAddr, feeDenom, feeErr := k.GetFeeAccountForMinFees(cacheCtx, flow, types.MaxGas)
		if feeErr != nil {
			errorString = appendError(errorString, feeErr.Error())
		} else if feeAddr == nil || feeDenom == "" {
			errorString = appendError(errorString, types.ErrBalanceTooLow)
		}

		k.Logger(ctx).Debug("triggering next msgs", "id", flow.ID, "msgs", len(nextMsgs))
		flowCopy := flow
		flowCopy.Msgs = nextMsgs
		ibcSequence, _, err := k.TriggerFlow(cacheCtx, &flowCopy)
		if err != nil {
			errorString = appendError(errorString, fmt.Sprintf(types.ErrFlowMsgHandling, err.Error()))
		}
		if ibcSequence > 0 {
			flowHistoryEntry.PacketSequences = append(flowHistoryEntry.PacketSequences, uint64(ibcSequence))
		}
		fee, distErr := k.DistributeCoins(cacheCtx, flow, feeAddr, feeDenom)
		if distErr != nil {
			errorString = appendError(errorString, fmt.Sprintf(types.ErrFlowFeeDistribution, distErr.Error()))
		} else if len(flowHistoryEntry.ExecFee) == 0 || feeDenom != flowHistoryEntry.ExecFee[0].Denom {
			// If the fee denom is different, reset the exec fee to the new fee. TODO: Use Coins for ExecFee to handle this better
			flowHistoryEntry.ExecFee = sdk.NewCoins(sdk.NewCoin(feeDenom, fee.Amount))
			k.Logger(ctx).Debug("execFee reset to new denom")
		} else {
			flowHistoryEntry.ExecFee = sdk.NewCoins(flowHistoryEntry.ExecFee[0].Add(fee))
		}
		k.Logger(ctx).Debug("execFee", flowHistoryEntry.ExecFee)
		if errorString != "" {
			flowHistoryEntry.Executed = false
			flowHistoryEntry.Errors = append(flowHistoryEntry.Errors, errorString)
		}

		k.SetCurrentFlowHistoryEntry(cacheCtx, flow.ID, flowHistoryEntry)
		k.SetFlow(ctx, &flow)
		writeCtx()
	}

	return nil
}

// HandleDeepResponses checks responses and unwraps authz responses from the IBC packet
func (k Keeper) HandleDeepResponses(ctx sdk.Context, msgResponses []*cdctypes.Any, relayer sdk.AccAddress, flow types.Flow, previousMsgsExecuted int) ([]*cdctypes.Any, int, error) {
	var msgClass int

	if previousMsgsExecuted == len(flow.Msgs) {
		return nil, 0, errorsmod.Wrapf(types.ErrMsgResponsesHandling, "no more messages to execute for flow %d, got %d responses for flow messages %d", flow.ID, len(msgResponses), len(flow.Msgs))
	}
	for index, anyResp := range msgResponses {
		k.Logger(ctx).Debug("msg response in ICS-27 packet", "response", anyResp.GoString(), "typeURL", anyResp.GetTypeUrl())

		if entry, ok := msg_registry.MsgRegistry[anyResp.GetTypeUrl()]; ok {
			if index == 0 && entry.RewardType > 0 {
				connectionID := flow.SelfHostedICA.ConnectionID
				if flow.TrustlessAgent != nil {
					if flow.TrustlessAgent.AgentAddress != "" {
						connectionID = k.GetTrustlessAgent(ctx, flow.TrustlessAgent.AgentAddress).ICAConfig.ConnectionID
					}
				}
				msgClass = entry.RewardType
				k.HandleRelayerReward(ctx, relayer, msgClass, connectionID)
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
				return nil, 0, errorsmod.Wrapf(types.ErrMsgResponsesHandling, "expected more message responses for flow %d, got %d responses for flow messages %d", flow.ID, len(msgExecResponse.Results), len(flow.Msgs))
			}
			msgExec := &authztypes.MsgExec{}
			if err := proto.Unmarshal(flow.Msgs[flowIndex].Value, msgExec); err != nil {
				return nil, 0, errorsmod.Wrapf(types.ErrMsgResponsesHandling, "failed to unmarshal MsgExec: %s", err.Error())
			}

			// Instead of resetting msgResponses, collect inner responses
			var innerResponses []*cdctypes.Any

			if len(msgExecResponse.Results) != len(msgExec.Msgs) {
				return nil, 0, errorsmod.Wrapf(types.ErrMsgResponsesHandling, "number of results (%d) does not match number of messages (%d) in MsgExec", len(msgExecResponse.Results), len(msgExec.Msgs))
			}

			for i, resultBytes := range msgExecResponse.Results {

				var msgResponse = cdctypes.Any{}
				if err := proto.Unmarshal(resultBytes, &msgResponse); err == nil {
					typeUrl := msgResponse.GetTypeUrl()
					if typeUrl != "" && strings.Contains(typeUrl, "Msg") {
						_, err := k.interfaceRegistry.Resolve(typeUrl)
						if err == nil {
							k.Logger(ctx).Debug("parsing response authz v0.52+", "response", msgResponse.GoString(), "typeURL", msgResponse.GetTypeUrl())
							innerResponses = append(innerResponses, &msgResponse)
							continue
						}
					}
				}
				k.Logger(ctx).Debug("parsing response authz v0.52-", "response", msgResponse.GoString(), "typeURL", msgResponse.GetTypeUrl())

				// fallback: handle as MsgData
				msgRespProto, _, err := handleMsgData(&sdk.MsgData{Data: resultBytes, MsgType: msgExec.Msgs[i].TypeUrl})
				if err != nil {
					return nil, 0, err
				}
				// Only create Any if we have a non-nil response
				if msgRespProto != nil {
					respAny, err := cdctypes.NewAnyWithValue(msgRespProto)
					if err != nil {
						return nil, 0, err
					}
					innerResponses = append(innerResponses, respAny)
				}
			}

			// Replace the MsgExecResponse with its inner responses
			msgResponses = append(msgResponses[:index], append(innerResponses, msgResponses[index+1:]...)...)
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
	flow := k.GetFlow(ctx, id)
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
func (k Keeper) SetFlowError(ctx sdk.Context, sourcePort, channelID string, seq uint64, errString string) {
	id := k.getTmpFlowID(ctx, sourcePort, channelID, seq)
	if id == 0 {
		return
	}

	k.Logger(ctx).Debug("flow", "id", id, "error", errString)

	flowHistoryEntry, err := k.GetLatestFlowHistoryEntry(ctx, id)
	if err != nil {
		flowHistoryEntry = &types.FlowHistoryEntry{Errors: []string{err.Error()}}
	}

	flow, err := k.TryGetFlow(ctx, id)
	if err != nil {
		flowHistoryEntry.Errors = append(flowHistoryEntry.Errors, err.Error())
	}

	if errString != "" {
		//only append error if it is not already in the list
		alreadyExists := false
		for _, e := range flowHistoryEntry.Errors {
			if e == errString {
				alreadyExists = true
				break
			}
		}
		if flowHistoryEntry.Errors == nil || !alreadyExists {
			flowHistoryEntry.Errors = append(flowHistoryEntry.Errors, errString)
		}
		flowHistoryEntry.Executed = false
		k.SetCurrentFlowHistoryEntry(ctx, id, flowHistoryEntry)

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeFlowError,
				sdk.NewAttribute(types.AttributeKeyFlowID, fmt.Sprint(id)),
				sdk.NewAttribute(types.AttributeKeyFlowOwner, flow.Owner),
				sdk.NewAttribute(types.AttributeKeyError, errString),
			),
		)
	}
}
