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
	"github.com/trstlabs/intento/x/intent/types"
)

func (k Keeper) TriggerFlow(ctx sdk.Context, flow *types.FlowInfo) (bool, []*cdctypes.Any, error) {
	// local flow
	if (flow.ICAConfig == nil || flow.ICAConfig.ConnectionID == "") && (flow.HostedConfig == nil || flow.HostedConfig.HostedAddress == "") {
		txMsgs := flow.GetTxMsgs(k.cdc)
		msgResponses, err := handleLocalFlow(k, ctx, txMsgs, *flow)
		return err == nil, msgResponses, errorsmod.Wrap(err, "could execute local flow")
	}

	connectionID := flow.ICAConfig.ConnectionID
	portID := flow.ICAConfig.PortID
	triggerAddress := flow.Owner
	//get hosted account from hosted config
	if flow.HostedConfig != nil && flow.HostedConfig.HostedAddress != "" {
		hostedAccount := k.GetHostedAccount(ctx, flow.HostedConfig.HostedAddress)
		connectionID = hostedAccount.ICAConfig.ConnectionID
		portID = hostedAccount.ICAConfig.PortID
		triggerAddress = hostedAccount.HostedAddress
		err := k.SendFeesToHosted(ctx, *flow, hostedAccount)
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
	k.setTmpFlowID(ctx, flow.ID, portID, channelID, res.Sequence)
	return false, nil, nil
}

func handleLocalFlow(k Keeper, ctx sdk.Context, txMsgs []sdk.Msg, flow types.FlowInfo) ([]*cdctypes.Any, error) {
	// CacheContext returns a new context with the multi-store branched into a cached storage object
	// writeCache is called only if all msgs succeed, performing state transitions atomically
	var msgResponses []*cdctypes.Any

	cacheCtx, writeCache := ctx.CacheContext()
	for index, msg := range txMsgs {
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

		signers, _, err := k.cdc.GetMsgV1Signers(msg)
		if err != nil {
			return nil, err
		}
		for _, acct := range signers {
			if sdk.AccAddress(acct).String() != flow.Owner {
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
	if !flow.Configuration.SaveResponses {
		msgResponses = nil
	}
	return msgResponses, nil
}

// HandleResponseAndSetFlowResult sets the result of the last executed ID set at SendFlow.
func (k Keeper) HandleResponseAndSetFlowResult(ctx sdk.Context, portID string, channelID string, relayer sdk.AccAddress, seq uint64, msgResponses []*cdctypes.Any) error {
	id := k.getTmpFlowID(ctx, portID, channelID, seq)
	if id <= 0 {
		return nil
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

	if len(flow.Conditions.FeedbackLoops) != 0 {
		for _, feedbackLoop := range flow.Conditions.FeedbackLoops {
			// Validate MsgsIndex and FlowID
			if feedbackLoop.MsgsIndex == 0 || feedbackLoop.FlowID != 0 {
				continue // Skip invalid FeedbackLoops or if FlowID is set to non-default
			}

			// Ensure MsgsIndex is within bounds
			if len(flowHistoryEntry.MsgResponses)-1 < int(feedbackLoop.MsgsIndex) {
				continue // Skip if MsgsIndex exceeds available responses
			}

			// Trigger remaining execution for the valid FeedbackLoop
			tmpFlow := flow
			tmpFlow.Msgs = flow.Msgs[feedbackLoop.MsgsIndex:]

			if err := triggerRemainingMsgs(k, ctx, tmpFlow, flowHistoryEntry); err != nil {
				return err // Return on the first encountered error
			}
		}
	}

	k.SetCurrentFlowHistoryEntry(ctx, flow.ID, flowHistoryEntry)
	return nil

}

func triggerRemainingMsgs(k Keeper, ctx sdk.Context, flow types.FlowInfo, flowHistoryEntry *types.FlowHistoryEntry) error {
	var errorString = ""

	allowed, err := k.allowedToExecute(ctx, flow)
	if !allowed {
		k.recordFlowNotAllowed(ctx, &flow, ctx.BlockTime(), err)

	}

	flowCtx := ctx.WithGasMeter(storetypes.NewGasMeter(types.MaxGas))
	cacheCtx, writeCtx := flowCtx.CacheContext()
	k.Logger(ctx).Debug("continuing msg execution", "id", flow.ID)

	feeAddr, feeDenom, err := k.GetFeeAccountForMinFees(cacheCtx, flow, types.MaxGas)
	if err != nil {
		errorString = appendError(errorString, err.Error())
	} else if feeAddr == nil || feeDenom == "" {
		errorString = appendError(errorString, (types.ErrBalanceTooLow + feeDenom))
	}

	err = k.RunFeedbackLoops(cacheCtx, flow.ID, &flow.Msgs, flow.Conditions)
	if err != nil {
		return errorsmod.Wrap(err, fmt.Sprintf(types.ErrSettingFlowResult, err))
	}

	k.Logger(ctx).Debug("triggering msgs", "id", flow.ID, "msgs", len(flow.Msgs))
	_, _, err = k.TriggerFlow(cacheCtx, &flow)
	if err != nil {
		errorString = appendError(errorString, fmt.Sprintf(types.ErrFlowMsgHandling, err.Error()))
	}
	fee, err := k.DistributeCoins(cacheCtx, flow, feeAddr, feeDenom, ctx.BlockHeader().ProposerAddress)
	if err != nil {
		errorString = appendError(errorString, fmt.Sprintf(types.ErrFlowFeeDistribution, err.Error()))
	}
	flowHistoryEntry.ExecFee = flowHistoryEntry.ExecFee.Add(fee)

	if errorString != "" {
		flowHistoryEntry.Executed = false
		flowHistoryEntry.Errors = append(flowHistoryEntry.Errors, types.ErrFlowMsgHandling+err.Error())

	}
	k.SetCurrentFlowHistoryEntry(cacheCtx, flow.ID, flowHistoryEntry)
	writeCtx()
	return nil
}

func (k Keeper) HandleDeepResponses(ctx sdk.Context, msgResponses []*cdctypes.Any, relayer sdk.AccAddress, flow types.FlowInfo, previousMsgsExecuted int) ([]*cdctypes.Any, int, error) {
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
				k.Logger(ctx).Debug("handling deep flow response unmarshalling", "err", err)
				return nil, 0, err
			}

			flowIndex := index + previousMsgsExecuted
			if flowIndex >= len(flow.Msgs) {
				return nil, 0, errorsmod.Wrapf(types.ErrMsgResponsesHandling, "expected more message responses")
			}
			msgExec := &authztypes.MsgExec{}
			if err := proto.Unmarshal(flow.Msgs[flowIndex].Value, msgExec); err != nil {
				return nil, 0, err
			}

			msgResponses = []*cdctypes.Any{}

			for _, resultBytes := range msgExecResponse.Results {
				var msgResponse = cdctypes.Any{}
				if err := proto.Unmarshal(resultBytes, &msgResponse); err == nil {
					typeUrl := msgResponse.GetTypeUrl()

					if typeUrl != "" && strings.Contains(typeUrl, "Msg") {
						// _, err := k.interfaceRegistry.Resolve(typeUrl)
						// if err == nil {
						k.Logger(ctx).Debug("parsing response authz v0.52+", "msgResponse", msgResponse)
						msgResponses = append(msgResponses, &msgResponse)
						continue
						//}
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
				respAny, err := cdctypes.NewAnyWithValue(msgRespProto)
				if err != nil {
					return nil, 0, err
				}

				msgResponses = append(msgResponses, respAny)

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
	if flow.Configuration.ReregisterICAAfterTimeout {
		flow := k.GetFlowInfo(ctx, id)
		metadataString := icatypes.NewDefaultMetadataString(flow.ICAConfig.ConnectionID, flow.ICAConfig.HostConnectionID)
		err := k.RegisterInterchainAccount(ctx, flow.ICAConfig.ConnectionID, flow.Owner, metadataString)
		if err != nil {
			return err
		}
	} else {
		k.RemoveFromFlowQueue(ctx, flow)
	}
	k.Logger(ctx).Debug("flow packet timed out", "flow_id", id)

	flowHistoryEntry, err := k.GetLatestFlowHistoryEntry(ctx, id)
	if err != nil {
		return err
	}

	flowHistoryEntry.TimedOut = true
	k.SetCurrentFlowHistoryEntry(ctx, id, flowHistoryEntry)

	return nil
}

// SetFlowOnTimeout sets the flow timeout result to the flow
func (k Keeper) SetFlowError(ctx sdk.Context, sourcePort string, channelID string, seq uint64, err string) {
	id := k.getTmpFlowID(ctx, sourcePort, channelID, seq)
	if id <= 0 {
		return
	}

	k.Logger(ctx).Debug("flow", "id", id, "error", err)

	flowHistoryEntry, newErr := k.GetLatestFlowHistoryEntry(ctx, id)
	if newErr != nil {
		flowHistoryEntry.Errors = append(flowHistoryEntry.Errors, newErr.Error())
	}

	flowHistoryEntry.Errors = append(flowHistoryEntry.Errors, err)
	k.SetCurrentFlowHistoryEntry(ctx, id, flowHistoryEntry)
}
