package keeper

import (
	"context"
	"fmt"
	"strconv"
	"time"

	errorsmod "cosmossdk.io/errors"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authztypes "github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/cosmos/gogoproto/proto"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller/keeper"
	icacontrollertypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller/types"
	icatypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/types"
	"github.com/trstlabs/intento/x/intent/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	Keeper
}

// NewMsgServerImpl creates and returns a new types.MsgServer, fulfilling the intent Msg service interface
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

// RegisterAccount implements the Msg/RegisterAccount interface
func (k msgServer) RegisterAccount(goCtx context.Context, msg *types.MsgRegisterAccount) (*types.MsgRegisterAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	err := k.RegisterInterchainAccount(ctx, msg.ConnectionID, msg.Owner, msg.Version)
	if err != nil {
		return nil, err
	}
	return &types.MsgRegisterAccountResponse{}, nil
}

// SubmitTx implements the Msg/SubmitTx interface
func (k msgServer) SubmitTx(goCtx context.Context, msg *types.MsgSubmitTx) (*types.MsgSubmitTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	data, err := icatypes.SerializeCosmosTx(k.cdc, []proto.Message{msg.GetTxMsg()}, icatypes.EncodingProtobuf)
	if err != nil {
		return nil, err
	}

	packetData := icatypes.InterchainAccountPacketData{
		Type: icatypes.EXECUTE_TX,
		Data: data,
	}

	relativeTimeoutTimestamp := uint64(time.Minute.Nanoseconds())
	msgServer := icacontrollerkeeper.NewMsgServerImpl(&k.icaControllerKeeper)
	icaMsg := icacontrollertypes.NewMsgSendTx(msg.Owner, msg.ConnectionID, relativeTimeoutTimestamp, packetData)

	_, err = msgServer.SendTx(ctx, icaMsg)
	if err != nil {
		return nil, err
	}

	return &types.MsgSubmitTxResponse{}, nil
}

// SubmitFlow implements the Msg/SubmitFlow interface
func (k msgServer) SubmitFlow(goCtx context.Context, msg *types.MsgSubmitFlow) (*types.MsgSubmitFlowResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	msgOwner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrInvalidRequest, err.Error())
	}

	portID, duration, interval, startTime, configuration, conditions, trustlessAgent, err := checkAndParseFlowContent(k, msg, err, ctx)
	if err != nil {
		return nil, err
	}

	err = k.CreateFlow(ctx, msgOwner, msg.Label, msg.Msgs, duration, interval, startTime, msg.FeeFunds, configuration, trustlessAgent, portID, msg.ConnectionID, conditions)
	if err != nil {
		return nil, err
	}

	return &types.MsgSubmitFlowResponse{}, nil
}

// RegisterAccountAndSubmitFlow implements the Msg/RegisterAccountAndSubmitFlow interface
func (k msgServer) RegisterAccountAndSubmitFlow(goCtx context.Context, msg *types.MsgRegisterAccountAndSubmitFlow) (*types.MsgRegisterAccountAndSubmitFlowResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := k.RegisterInterchainAccount(ctx, msg.ConnectionID, msg.Owner, msg.Version)
	if err != nil {
		return nil, err
	}

	msgOwner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrInvalidRequest, err.Error())
	}

	portID, duration, interval, startTime, configuration, conditions, _, err := checkAndParseFlowContent(k, msg, err, ctx)
	if err != nil {
		return nil, err
	}

	err = k.CreateFlow(ctx, msgOwner, msg.Label, msg.Msgs, duration, interval, startTime, msg.FeeFunds, configuration, types.TrustlessAgentConfig{}, portID, msg.ConnectionID, conditions)
	if err != nil {
		return nil, err
	}

	return &types.MsgRegisterAccountAndSubmitFlowResponse{}, nil
}

// UpdateFlow implements the Msg/UpdateFlow interface
func (k msgServer) UpdateFlow(goCtx context.Context, msg *types.MsgUpdateFlow) (*types.MsgUpdateFlowResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	timeNowWindow := ctx.BlockTime().Add(time.Minute * 1)
	flow, err := k.TryGetflow(ctx, msg.ID)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrInvalidRequest, err.Error())
	}
	if flow.Owner != msg.Owner {
		return nil, types.ErrInvalidAddress
	}

	if flow.Configuration.UpdatingDisabled {
		return nil, errorsmod.Wrap(types.ErrInvalidRequest, "updating is disabled")
	}

	if msg.ConnectionID != "" {
		flow.SelfHostedICA.PortID, err = icatypes.NewControllerPortID(msg.Owner)
		if err != nil {
			return nil, err
		}
		flow.SelfHostedICA.ConnectionID = msg.ConnectionID
	}
	newExecTime := flow.ExecTime
	if msg.EndTime > 0 {
		endTime := time.Unix(int64(msg.EndTime), 0)
		if err != nil {
			return nil, err
		}
		if endTime.Before(timeNowWindow) {
			return nil, types.ErrInvalidTime
		}
		flow.EndTime = endTime

	}
	p, err := k.GetParams(ctx)
	if err != nil {
		panic(err)
	}

	if msg.Interval != "" {
		interval, err := time.ParseDuration(msg.Interval)
		if err != nil {
			return nil, errorsmod.Wrap(types.ErrUpdateFlow, err.Error())
		}

		if interval != 0 && interval < p.MinFlowInterval || interval > flow.EndTime.Sub(flow.StartTime) {
			return nil, errorsmod.Wrapf(types.ErrUpdateFlow, "interval: %s  must be longer than minimum interval:  %s, and execution should happen before end time", interval, p.MinFlowInterval)
		}
		flow.Interval = interval
	}

	if msg.StartAt > 0 {
		startTime := time.Unix(int64(msg.StartAt), 0)
		if startTime.Before(ctx.BlockHeader().Time.Add(time.Minute)) {
			return nil, errorsmod.Wrapf(types.ErrInvalidRequest, "custom start time: %s must be at least a minute into the future upon block submission: %s", startTime, ctx.BlockHeader().Time.Add(time.Minute))
		}
		if startTime.After(flow.EndTime) {
			return nil, errorsmod.Wrapf(types.ErrUpdateFlow, "start time: %s must be before end time", startTime)
		}
		latestEntry, err := k.GetLatestFlowHistoryEntry(ctx, flow.ID)
		if err != nil || latestEntry != nil {
			return nil, errorsmod.Wrapf(types.ErrUpdateFlow, "start time: %s must occur before first execution", startTime)
		}
		flow.StartTime = startTime
		newExecTime = startTime
	}

	if msg.Label != "" {
		flow.Label = msg.Label
	}

	if msg.Configuration != nil {
		flow.Configuration = msg.Configuration
	}

	if msg.Conditions != nil {
		err = updateConditions(flow.Conditions, msg.Msgs, flow.EndTime.Sub(timeNowWindow), flow.Interval)
		if err != nil {
			fmt.Printf("Error updating conditions: %v\n", err)
			return nil, err
		}

	}
	if len(msg.Msgs) != 0 {
		flow.Msgs = msg.Msgs
	}

	if newExecTime != flow.ExecTime {
		k.RemoveFromFlowQueue(ctx, flow)
		flow.ExecTime = newExecTime
		k.InsertFlowQueue(ctx, flow.ID, newExecTime)
	}

	flow.UpdateHistory = append(flow.UpdateHistory, ctx.BlockTime())

	if err := k.SignerOk(ctx, k.cdc, flow); err != nil {
		return nil, errorsmod.Wrap(types.ErrSignerNotOk, err.Error())
	}
	//set hosted config
	if msg.TrustlessAgent != nil {
		flow.TrustlessAgent = msg.TrustlessAgent
	}

	k.Setflow(ctx, &flow)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeFlowUpdated,
			sdk.NewAttribute(types.AttributeKeyFlowID, strconv.FormatUint(flow.ID, 10)),
			sdk.NewAttribute(types.AttributeKeyFlowOwner, flow.Owner),
		))

	return &types.MsgUpdateFlowResponse{}, nil
}

// CreateTrustlessAgent implements the Msg/CreateTrustlessAgent interface
func (k msgServer) CreateTrustlessAgent(goCtx context.Context, msg *types.MsgCreateTrustlessAgent) (*types.MsgCreateTrustlessAgentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	agentAddress, err := DeriveAgentAddress(msg.Creator, msg.ConnectionID)
	if err != nil {
		return nil, err
	}
	//register ICA
	err = k.RegisterInterchainAccount(ctx, msg.ConnectionID, agentAddress.String(), msg.Version)
	if err != nil {
		return nil, err
	}
	portID, err := icatypes.NewControllerPortID(agentAddress.String())
	if err != nil {
		return nil, err
	}
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrInvalidRequest, err.Error())
	}
	//store hosted config by address on hosted key prefix
	k.SetTrustlessAgent(ctx, &types.TrustlessAgent{AgentAddress: agentAddress.String(), FeeConfig: &types.TrustlessAgentFeeConfig{FeeAdmin: msg.Creator, FeeCoinsSupported: msg.FeeCoinsSupported}, ICAConfig: &types.ICAConfig{ConnectionID: msg.ConnectionID, PortID: portID}})
	k.addToTrustlessAgentAdminIndex(ctx, creator, agentAddress.String())
	return &types.MsgCreateTrustlessAgentResponse{Address: agentAddress.String()}, nil
}

// UpdateHosted implements the Msg/UpdateHosted interface
func (k msgServer) UpdateTrustlessAgentFeeConfig(goCtx context.Context, msg *types.MsgUpdateTrustlessAgentFeeConfig) (*types.MsgUpdateTrustlessAgentFeeConfigResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	//get hosted config by address on hosted key prefix
	trustlessAgent, err := k.TryGetTrustlessAgent(ctx, msg.AgentAddress)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrInvalidRequest, err.Error())
	}
	//check admin address
	if trustlessAgent.FeeConfig.FeeAdmin != msg.FeeAdmin {
		return nil, types.ErrInvalidAddress
	}

	agentAddress := trustlessAgent.AgentAddress

	admin := trustlessAgent.FeeConfig.FeeAdmin
	if msg.FeeConfig.FeeAdmin != "" {
		newAdminAddr, err := sdk.AccAddressFromBech32(msg.FeeConfig.FeeAdmin)
		if err != nil {
			return nil, errorsmod.Wrap(types.ErrInvalidRequest, err.Error())
		}
		currentAdminAddr, err := sdk.AccAddressFromBech32(trustlessAgent.FeeConfig.FeeAdmin)
		if err != nil {
			return nil, errorsmod.Wrap(types.ErrInvalidRequest, err.Error())
		}
		admin = msg.FeeConfig.FeeAdmin
		k.changeTrustlessAgentAdminIndex(ctx, currentAdminAddr, newAdminAddr, trustlessAgent.AgentAddress)
	}

	feeCoinsSupported := trustlessAgent.FeeConfig.FeeCoinsSupported
	if msg.FeeConfig.FeeCoinsSupported != nil {
		feeCoinsSupported = msg.FeeConfig.FeeCoinsSupported
	}

	k.SetTrustlessAgent(ctx, &types.TrustlessAgent{AgentAddress: agentAddress, FeeConfig: &types.TrustlessAgentFeeConfig{FeeAdmin: admin, FeeCoinsSupported: feeCoinsSupported}, ICAConfig: trustlessAgent.ICAConfig})

	//set hosted config by address on hosted key prefix
	return &types.MsgUpdateTrustlessAgentFeeConfigResponse{}, nil
}

// UpdateParams implements the Msg/UpdateParams interface
func (k msgServer) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate the authority
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address: %s", err)
	}

	// Verify that the authority matches the module's authority
	if msg.Authority != k.Keeper.GetAuthority() {
		return nil, errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "expected %s, got %s", k.Keeper.GetAuthority(), msg.Authority)
	}

	// Validate the parameters
	if err := msg.Params.Validate(); err != nil {
		return nil, err
	}

	// Set the new parameters
	k.Keeper.SetParams(ctx, msg.Params)

	return &types.MsgUpdateParamsResponse{}, nil
}

func checkAndParseFlowContent(
	k msgServer,
	msg sdk.Msg,
	err error,
	ctx sdk.Context,
) (string, time.Duration, time.Duration, time.Time, types.ExecutionConfiguration, types.ExecutionConditions, types.TrustlessAgentConfig, error) {
	var (
		msgOwner         string
		msgConnectionID  string
		msgDuration      string
		msgInterval      string
		msgStartAt       uint64
		msgConfiguration *types.ExecutionConfiguration = &types.ExecutionConfiguration{}
		msgConditions    *types.ExecutionConditions    = &types.ExecutionConditions{}
		TrustlessAgent   *types.TrustlessAgentConfig   = &types.TrustlessAgentConfig{}
		msgMsgs          []*cdctypes.Any               = []*cdctypes.Any{}
	)

	switch msg := msg.(type) {
	case *types.MsgSubmitFlow:
		// Existing logic for MsgSubmitFlow
		msgOwner = msg.Owner
		msgDuration = msg.Duration
		msgConnectionID = msg.ConnectionID
		msgStartAt = msg.StartAt
		msgInterval = msg.Interval
		// Use fallback if TrustlessAgent is nil
		if msg.TrustlessAgent != nil {
			TrustlessAgent = msg.TrustlessAgent
		}
		// Use fallback if msgConfiguration is nil
		if msg.Configuration != nil {
			msgConfiguration = msg.Configuration
		}
		// Use fallback if msgConditions is nil
		if msg.Conditions != nil {
			msgConditions = msg.Conditions
		}

		msgMsgs = msg.Msgs

	case *types.MsgRegisterAccountAndSubmitFlow:
		// Handle RegisterAccountAndSubmitFlow
		msgOwner = msg.Owner
		msgDuration = msg.Duration
		msgConnectionID = msg.ConnectionID
		msgStartAt = msg.StartAt
		msgInterval = msg.Interval
		if msg.Configuration != nil {
			msgConfiguration = msg.Configuration
		}
		if msg.Conditions != nil {
			msgConditions = msg.Conditions
		}

		msgMsgs = msg.Msgs

	default:
		return "", 0, 0, time.Time{}, types.ExecutionConfiguration{}, types.ExecutionConditions{}, types.TrustlessAgentConfig{}, errorsmod.Wrapf(types.ErrInvalidRequest, "unsupported intento module message type: %T", msg)
	}

	portID := ""
	if msgConnectionID != "" {
		portID, err = icatypes.NewControllerPortID(msgOwner)
		if err != nil {
			return "", 0, 0, time.Time{}, types.ExecutionConfiguration{}, types.ExecutionConditions{}, types.TrustlessAgentConfig{}, err
		}
	}

	var duration time.Duration = 0
	if msgDuration != "" {
		duration, err = time.ParseDuration(msgDuration)
		if err != nil {
			return "", 0, 0, time.Time{}, types.ExecutionConfiguration{}, types.ExecutionConditions{}, types.TrustlessAgentConfig{}, err
		}
	}

	var interval time.Duration = 0
	if msgInterval != "" {
		interval, err = time.ParseDuration(msgInterval)
		if err != nil {
			return "", 0, 0, time.Time{}, types.ExecutionConfiguration{}, types.ExecutionConditions{}, types.TrustlessAgentConfig{}, errorsmod.Wrap(types.ErrInvalidScheduling, err.Error())
		}
	}

	var startTime time.Time = ctx.BlockHeader().Time
	if msgStartAt != 0 {
		startTime = time.Unix(int64(msgStartAt), 0)
		if err != nil {
			return "", 0, 0, time.Time{}, types.ExecutionConfiguration{}, types.ExecutionConditions{}, types.TrustlessAgentConfig{}, err
		}
		if startTime.Before(ctx.BlockHeader().Time.Add(time.Minute)) {
			return "", 0, 0, time.Time{}, types.ExecutionConfiguration{}, types.ExecutionConditions{}, types.TrustlessAgentConfig{}, errorsmod.Wrapf(types.ErrInvalidScheduling, "custom start time: %s must be at least a minute into the future upon block submission: %s", startTime, ctx.BlockHeader().Time.Add(time.Minute))
		}
	}

	p, err := k.GetParams(ctx)
	if err != nil {
		panic(err)
	}
	if interval != 0 && (interval < p.MinFlowInterval || interval > duration) {
		return "", 0, 0, time.Time{}, types.ExecutionConfiguration{}, types.ExecutionConditions{}, types.TrustlessAgentConfig{}, errorsmod.Wrapf(types.ErrInvalidScheduling, "interval: %s  must be longer than minimum interval:  %s, and longer than duration: %s", interval, p.MinFlowInterval, duration)
	}
	if duration != 0 {
		if duration > p.MaxFlowDuration {
			return "", 0, 0, time.Time{}, types.ExecutionConfiguration{}, types.ExecutionConditions{}, types.TrustlessAgentConfig{}, errorsmod.Wrapf(types.ErrInvalidScheduling, "duration: %s must be shorter than maximum duration: %s", duration, p.MaxFlowDuration)
		} else if duration < p.MinFlowDuration {
			return "", 0, 0, time.Time{}, types.ExecutionConfiguration{}, types.ExecutionConditions{}, types.TrustlessAgentConfig{}, errorsmod.Wrapf(types.ErrInvalidScheduling, "duration: %s must be longer than minimum duration: %s", duration, p.MinFlowDuration)
		} else if startTime.After(ctx.BlockHeader().Time.Add(p.MaxFlowDuration)) {
			return "", 0, 0, time.Time{}, types.ExecutionConfiguration{}, types.ExecutionConditions{}, types.TrustlessAgentConfig{}, errorsmod.Wrapf(types.ErrInvalidScheduling, "start time: %s must be before current time and max duration: %s", startTime, ctx.BlockHeader().Time.Add(duration))
		}
	}

	err = updateConditions(msgConditions, msgMsgs, duration, interval)
	if err != nil {
		return "", 0, 0, time.Time{}, types.ExecutionConfiguration{}, types.ExecutionConditions{}, types.TrustlessAgentConfig{}, err
	}

	return portID, duration, interval, startTime,
		*msgConfiguration,
		*msgConditions,
		*TrustlessAgent,
		nil
}
func updateConditions(
	msgConditions *types.ExecutionConditions,
	msgMsgs []*cdctypes.Any,
	duration, interval time.Duration,
) error {
	if msgConditions == nil {
		return nil
	}

	// --- FeedbackLoops validation ---
	if msgConditions.FeedbackLoops != nil {
		for _, loop := range msgConditions.FeedbackLoops {
			if err := validateFeedbackLoop(loop, msgMsgs); err != nil {
				return err
			}
			if err := validateTimeout("FeedbackLoop", loop.ICQConfig, duration, interval); err != nil {
				return err
			}
		}
	}

	// --- Comparisons validation ---
	if msgConditions.Comparisons != nil {
		for _, cmp := range msgConditions.Comparisons {
			if int(cmp.ResponseIndex) < 0 || int(cmp.ResponseIndex) >= len(msgMsgs) {
				return errorsmod.Wrapf(
					types.ErrInvalidRequest,
					"Comparisons: response index %d out of bounds (len: %d)",
					cmp.ResponseIndex, len(msgMsgs),
				)
			}
			if err := validateTimeout("Comparison", cmp.ICQConfig, duration, interval); err != nil {
				return err
			}
		}
	}

	return nil
}

// --- Helper: FeedbackLoop validation ---
func validateFeedbackLoop(loop *types.FeedbackLoop, msgMsgs []*cdctypes.Any) error {
	if len(msgMsgs) == 0 || int(loop.MsgsIndex) < 0 {
		return errorsmod.Wrapf(types.ErrInvalidFeedbackLoop, "FeedbackLoop: empty msgMsgs or negative MsgsIndex")
	}

	if int(loop.MsgsIndex) >= len(msgMsgs) {
		msg := msgMsgs[0]
		if msg.TypeUrl == sdk.MsgTypeURL(&authztypes.MsgExec{}) {
			msgExec := &authztypes.MsgExec{}
			if err := proto.Unmarshal(msg.Value, msgExec); err != nil {
				return errorsmod.Wrapf(types.ErrInvalidFeedbackLoop, "could not unmarshal MsgExec to validate FeedbackLoop MsgsIndex")
			}
			if int(loop.MsgsIndex) >= len(msgExec.Msgs) {
				return errorsmod.Wrapf(types.ErrInvalidFeedbackLoop,
					"FeedbackLoop MsgsIndex %d out of bounds for MsgExec (len: %d)",
					loop.MsgsIndex, len(msgExec.Msgs),
				)
			}
			return errorsmod.Wrapf(types.ErrInvalidFeedbackLoop,
				"FeedbackLoop MsgsIndex %d out of bounds for msgMsgs (len: %d)",
				loop.MsgsIndex, len(msgMsgs),
			)
		}

		return errorsmod.Wrapf(types.ErrInvalidFeedbackLoop,
			"FeedbackLoop MsgsIndex %d out of bounds for msgMsgs (len: %d)",
			loop.MsgsIndex, len(msgMsgs),
		)
	}
	return nil
}

// --- Helper: Timeout validation ---
func validateTimeout(label string, icq *types.ICQConfig, duration, interval time.Duration) error {
	if icq == nil || icq.TimeoutDuration == 0 {
		return nil
	}
	if icq.TimeoutDuration > duration {
		return errorsmod.Wrapf(types.ErrInvalidICQTimeout,
			"%s TimeoutDuration (%s) exceeds flow duration (%s)", label, icq.TimeoutDuration, duration)
	}
	if interval != 0 && icq.TimeoutDuration > interval {
		return errorsmod.Wrapf(types.ErrInvalidICQTimeout,
			"%s TimeoutDuration (%s) exceeds flow interval (%s)", label, icq.TimeoutDuration, interval)
	}
	return nil
}
