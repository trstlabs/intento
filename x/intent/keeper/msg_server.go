package keeper

import (
	"context"
	"fmt"
	"time"

	errorsmod "cosmossdk.io/errors"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

	portID, duration, interval, startTime, configuration, conditions, HostedICAConfig, err := checkAndParseFlowContent(k, msg, err, ctx)
	if err != nil {
		return nil, err
	}

	err = k.CreateFlow(ctx, msgOwner, msg.Label, msg.Msgs, duration, interval, startTime, msg.FeeFunds, configuration, HostedICAConfig, portID, msg.ConnectionID, conditions)
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

	err = k.CreateFlow(ctx, msgOwner, msg.Label, msg.Msgs, duration, interval, startTime, msg.FeeFunds, configuration, types.HostedICAConfig{}, portID, msg.ConnectionID, conditions)
	if err != nil {
		return nil, err
	}

	return &types.MsgRegisterAccountAndSubmitFlowResponse{}, nil
}

// UpdateFlow implements the Msg/UpdateFlow interface
func (k msgServer) UpdateFlow(goCtx context.Context, msg *types.MsgUpdateFlow) (*types.MsgUpdateFlowResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	timeNowWindow := ctx.BlockTime().Add(time.Minute * 2)
	flow, err := k.TryGetFlowInfo(ctx, msg.ID)
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
		flow.ICAConfig.PortID, err = icatypes.NewControllerPortID(msg.Owner)
		if err != nil {
			return nil, err
		}
		flow.ICAConfig.ConnectionID = msg.ConnectionID
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

		if flow.Interval != 0 && msg.Interval != "" {
			newExecTime = endTime
		}
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
	if msg.HostedICAConfig != nil {
		flow.HostedICAConfig = msg.HostedICAConfig
	}

	k.SetFlowInfo(ctx, &flow)

	return &types.MsgUpdateFlowResponse{}, nil
}

// CreateHostedAccount implements the Msg/CreateHostedAccount interface
func (k msgServer) CreateHostedAccount(goCtx context.Context, msg *types.MsgCreateHostedAccount) (*types.MsgCreateHostedAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	hostedAddress, err := DeriveHostedAddress(msg.Creator, msg.ConnectionID)
	if err != nil {
		return nil, err
	}
	//register ICA
	err = k.RegisterInterchainAccount(ctx, msg.ConnectionID, hostedAddress.String(), msg.Version)
	if err != nil {
		return nil, err
	}
	portID, err := icatypes.NewControllerPortID(hostedAddress.String())
	if err != nil {
		return nil, err
	}
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrInvalidRequest, err.Error())
	}
	//store hosted config by address on hosted key prefix
	k.SetHostedAccount(ctx, &types.HostedAccount{HostedAddress: hostedAddress.String(), HostFeeConfig: &types.HostFeeConfig{Admin: msg.Creator, FeeCoinsSuported: msg.FeeCoinsSuported}, ICAConfig: &types.ICAConfig{ConnectionID: msg.ConnectionID, PortID: portID}})
	k.addToHostedAccountAdminIndex(ctx, creator, hostedAddress.String())
	return &types.MsgCreateHostedAccountResponse{Address: hostedAddress.String()}, nil
}

// UpdateHosted implements the Msg/UpdateHosted interface
func (k msgServer) UpdateHostedAccount(goCtx context.Context, msg *types.MsgUpdateHostedAccount) (*types.MsgUpdateHostedAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	//get hosted config by address on hosted key prefix
	hostedAccount, err := k.TryGetHostedAccount(ctx, msg.HostedAddress)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrInvalidRequest, err.Error())
	}
	//check admin address
	if hostedAccount.HostFeeConfig.Admin != msg.Admin {
		return nil, types.ErrInvalidAddress
	}

	hostedAddress := hostedAccount.HostedAddress

	admin := hostedAccount.HostFeeConfig.Admin
	if msg.HostFeeConfig.Admin != "" {
		newAdminAddr, err := sdk.AccAddressFromBech32(msg.HostFeeConfig.Admin)
		if err != nil {
			return nil, errorsmod.Wrap(types.ErrInvalidRequest, err.Error())
		}
		currentAdminAddr, err := sdk.AccAddressFromBech32(hostedAccount.HostFeeConfig.Admin)
		if err != nil {
			return nil, errorsmod.Wrap(types.ErrInvalidRequest, err.Error())
		}
		admin = msg.HostFeeConfig.Admin
		k.changeHostedAccountAdminIndex(ctx, currentAdminAddr, newAdminAddr, hostedAccount.HostedAddress)
	}

	feeCoinsSupported := hostedAccount.HostFeeConfig.FeeCoinsSuported
	if msg.HostFeeConfig.FeeCoinsSuported != nil {
		feeCoinsSupported = msg.HostFeeConfig.FeeCoinsSuported
	}

	k.SetHostedAccount(ctx, &types.HostedAccount{HostedAddress: hostedAddress, HostFeeConfig: &types.HostFeeConfig{Admin: admin, FeeCoinsSuported: feeCoinsSupported}, ICAConfig: hostedAccount.ICAConfig})

	//set hosted config by address on hosted key prefix
	return &types.MsgUpdateHostedAccountResponse{}, nil
}

func checkAndParseFlowContent(
	k msgServer,
	msg sdk.Msg,
	err error,
	ctx sdk.Context,
) (string, time.Duration, time.Duration, time.Time, types.ExecutionConfiguration, types.ExecutionConditions, types.HostedICAConfig, error) {
	var (
		msgOwner         string
		msgConnectionID  string
		msgDuration      string
		msgInterval      string
		msgStartAt       uint64
		msgConfiguration *types.ExecutionConfiguration = &types.ExecutionConfiguration{}
		msgConditions    *types.ExecutionConditions    = &types.ExecutionConditions{}
		HostedICAConfig  *types.HostedICAConfig        = &types.HostedICAConfig{}
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
		// Use fallback if HostedICAConfig is nil
		if msg.HostedICAConfig != nil {
			HostedICAConfig = msg.HostedICAConfig
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
		return "", 0, 0, time.Time{}, types.ExecutionConfiguration{}, types.ExecutionConditions{}, types.HostedICAConfig{}, fmt.Errorf("unsupported message type: %T", msg)
	}

	portID := ""
	if msgConnectionID != "" {
		portID, err = icatypes.NewControllerPortID(msgOwner)
		if err != nil {
			return "", 0, 0, time.Time{}, types.ExecutionConfiguration{}, types.ExecutionConditions{}, types.HostedICAConfig{}, err
		}
	}

	var duration time.Duration = 0
	if msgDuration != "" {
		duration, err = time.ParseDuration(msgDuration)
		if err != nil {
			return "", 0, 0, time.Time{}, types.ExecutionConfiguration{}, types.ExecutionConditions{}, types.HostedICAConfig{}, err
		}
	}

	var interval time.Duration = 0
	if msgInterval != "" {
		interval, err = time.ParseDuration(msgInterval)
		if err != nil {
			return "", 0, 0, time.Time{}, types.ExecutionConfiguration{}, types.ExecutionConditions{}, types.HostedICAConfig{}, errorsmod.Wrap(types.ErrInvalidRequest, err.Error())
		}
	}

	var startTime time.Time = ctx.BlockHeader().Time
	if msgStartAt != 0 {
		startTime = time.Unix(int64(msgStartAt), 0)
		if err != nil {
			return "", 0, 0, time.Time{}, types.ExecutionConfiguration{}, types.ExecutionConditions{}, types.HostedICAConfig{}, err
		}
		if startTime.Before(ctx.BlockHeader().Time.Add(time.Minute)) {
			return "", 0, 0, time.Time{}, types.ExecutionConfiguration{}, types.ExecutionConditions{}, types.HostedICAConfig{}, errorsmod.Wrapf(types.ErrInvalidRequest, "custom start time: %s must be at least a minute into the future upon block submission: %s", startTime, ctx.BlockHeader().Time.Add(time.Minute))
		}
	}

	p, err := k.GetParams(ctx)
	if err != nil {
		panic(err)
	}
	if interval != 0 && (interval < p.MinFlowInterval || interval > duration) {
		return "", 0, 0, time.Time{}, types.ExecutionConfiguration{}, types.ExecutionConditions{}, types.HostedICAConfig{}, errorsmod.Wrapf(types.ErrInvalidRequest, "interval: %s  must be longer than minimum interval:  %s, and longer than duration: %s", interval, p.MinFlowInterval, duration)
	}
	if duration != 0 {
		if duration > p.MaxFlowDuration {
			return "", 0, 0, time.Time{}, types.ExecutionConfiguration{}, types.ExecutionConditions{}, types.HostedICAConfig{}, errorsmod.Wrapf(types.ErrInvalidRequest, "duration: %s must be shorter than maximum duration: %s", duration, p.MaxFlowDuration)
		} else if duration < p.MinFlowDuration {
			return "", 0, 0, time.Time{}, types.ExecutionConfiguration{}, types.ExecutionConditions{}, types.HostedICAConfig{}, errorsmod.Wrapf(types.ErrInvalidRequest, "duration: %s must be longer than minimum duration: %s", duration, p.MinFlowDuration)
		} else if startTime.After(ctx.BlockHeader().Time.Add(p.MaxFlowDuration)) {
			return "", 0, 0, time.Time{}, types.ExecutionConfiguration{}, types.ExecutionConditions{}, types.HostedICAConfig{}, errorsmod.Wrapf(types.ErrInvalidRequest, "start time: %s must be before current time and max duration: %s", startTime, ctx.BlockHeader().Time.Add(duration))
		}
	}

	err = updateConditions(msgConditions, msgMsgs, duration, interval)
	if err != nil {
		return "", 0, 0, time.Time{}, types.ExecutionConfiguration{}, types.ExecutionConditions{}, types.HostedICAConfig{}, err
	}

	return portID, duration, interval, startTime,
		*msgConfiguration,
		*msgConditions,
		*HostedICAConfig,
		nil
}

func updateConditions(
	msgConditions *types.ExecutionConditions,
	msgMsgs []*cdctypes.Any,
	duration, interval time.Duration,
) error {
	if msgConditions != nil && msgConditions.FeedbackLoops != nil {
		for _, feedbackLoop := range msgConditions.FeedbackLoops {
			// Validate MsgsIndex for FeedbackLoops
			if int(feedbackLoop.MsgsIndex) >= len(msgMsgs) {
				if msgMsgs[0].TypeUrl == sdk.MsgTypeURL(&authztypes.MsgExec{}) {
					msgExec := &authztypes.MsgExec{}
					if err := proto.Unmarshal(msgMsgs[0].Value, msgExec); err != nil {
						return errorsmod.Wrapf(types.ErrInvalidRequest, "msg exec could not unmarshal, cannot check conditions")
					}

					if int(feedbackLoop.MsgsIndex) >= len(msgExec.Msgs) {
						return errorsmod.Wrapf(types.ErrInvalidRequest, "msgs index: %v must be shorter than length msgExec msgs array: %s", feedbackLoop.MsgsIndex, msgExec.Msgs)
					} else {
						return errorsmod.Wrapf(types.ErrInvalidRequest, "msgs index: %v must be shorter than length msgs array: %s", feedbackLoop.MsgsIndex, msgMsgs)
					}
				}
			}
			// Validate ICQConfig TimeoutDuration for FeedbackLoops
			if feedbackLoop.ICQConfig != nil {
				if feedbackLoop.ICQConfig.TimeoutDuration != 0 {
					if feedbackLoop.ICQConfig.TimeoutDuration > duration || (interval != 0 && feedbackLoop.ICQConfig.TimeoutDuration > interval) {
						return errorsmod.Wrapf(types.ErrInvalidRequest, "TimeoutDuration must be shorter than the flow interval or duration")
					}
				}
			}
		}
	}

	if msgConditions != nil && msgConditions.Comparisons != nil {
		for _, comparison := range msgConditions.Comparisons {
			// Validate ResponseIndex for Comparisons
			if int(comparison.ResponseIndex) >= len(msgMsgs) {
				return errorsmod.Wrapf(types.ErrInvalidRequest, "response index: %v must be shorter than length msgs array: %s", comparison.ResponseIndex, msgMsgs)
			}
			// Validate ICQConfig TimeoutDuration for Comparisons
			if comparison.ICQConfig != nil {
				if comparison.ICQConfig.TimeoutDuration != 0 {
					if comparison.ICQConfig.TimeoutDuration > duration || (interval != 0 && comparison.ICQConfig.TimeoutDuration > interval) {
						return errorsmod.Wrapf(types.ErrInvalidRequest, "TimeoutDuration must be shorter than the flow interval or duration")
					}
				}
			}
		}
	}

	return nil
}
