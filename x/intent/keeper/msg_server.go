package keeper

import (
	"context"
	"time"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authztypes "github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/cosmos/gogoproto/proto"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/keeper"
	icacontrollertypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/types"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
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
	err := k.RegisterInterchainAccount(ctx, msg.ConnectionId, msg.Owner, msg.Version)
	if err != nil {
		return nil, err
	}
	return &types.MsgRegisterAccountResponse{}, nil
}

// SubmitTx implements the Msg/SubmitTx interface
func (k msgServer) SubmitTx(goCtx context.Context, msg *types.MsgSubmitTx) (*types.MsgSubmitTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	data, err := icatypes.SerializeCosmosTx(k.cdc, []proto.Message{msg.GetTxMsg()})
	if err != nil {
		return nil, err
	}

	packetData := icatypes.InterchainAccountPacketData{
		Type: icatypes.EXECUTE_TX,
		Data: data,
	}

	relativeTimeoutTimestamp := uint64(time.Minute.Nanoseconds())
	msgServer := icacontrollerkeeper.NewMsgServerImpl(&k.icaControllerKeeper)
	icaMsg := icacontrollertypes.NewMsgSendTx(msg.Owner, msg.ConnectionId, relativeTimeoutTimestamp, packetData)

	_, err = msgServer.SendTx(ctx, icaMsg)
	if err != nil {
		return nil, err
	}

	// //store 0 as action id as a regular submit is not action
	// k.setTmpActionID(ctx, 0, portID, "", sequence)
	return &types.MsgSubmitTxResponse{}, nil
}

// SubmitAction implements the Msg/SubmitAction interface
func (k msgServer) SubmitAction(goCtx context.Context, msg *types.MsgSubmitAction) (*types.MsgSubmitActionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	msgOwner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrInvalidRequest, err.Error())
	}
	portID := ""
	if msg.ConnectionId != "" {
		portID, err = icatypes.NewControllerPortID(msg.Owner)
		if err != nil {
			return nil, err
		}
	}

	var duration time.Duration = 0
	if msg.Duration != "" {
		duration, err = time.ParseDuration(msg.Duration)
		if err != nil {
			return nil, err
		}
	}

	var interval time.Duration = 0
	if msg.Interval != "" {
		interval, err = time.ParseDuration(msg.Interval)
		if err != nil {
			return nil, errorsmod.Wrap(types.ErrInvalidRequest, err.Error())
		}
	}

	var startTime time.Time = ctx.BlockHeader().Time
	if msg.StartAt != 0 {
		startTime = time.Unix(int64(msg.StartAt), 0)
		if err != nil {
			return nil, err
		}
		if startTime.Before(ctx.BlockHeader().Time.Add(time.Minute)) {
			return nil, errorsmod.Wrapf(types.ErrInvalidRequest, "custom start time: %s must be at least a minute into the future upon block submission: %s", startTime, ctx.BlockHeader().Time.Add(time.Minute))
		}
	}

	p := k.GetParams(ctx)
	if interval != 0 && (interval < p.MinActionInterval || interval > duration) {
		return nil, errorsmod.Wrapf(types.ErrInvalidRequest, "interval: %s  must be longer than minimum interval:  %s, and longer than duration: %s", interval, p.MinActionInterval, duration)
	}
	if duration != 0 {
		if duration > p.MaxActionDuration {
			return nil, errorsmod.Wrapf(types.ErrInvalidRequest, "duration: %s must be shorter than maximum duration: %s", duration, p.MaxActionDuration)
		} else if duration < p.MinActionDuration {
			return nil, errorsmod.Wrapf(types.ErrInvalidRequest, "duration: %s must be longer than minimum duration: %s", duration, p.MinActionDuration)
		} else if startTime.After(ctx.BlockHeader().Time.Add(p.MaxActionDuration)) {
			return nil, errorsmod.Wrapf(types.ErrInvalidRequest, "start time: %s must be before current time and max duration: %s", startTime, ctx.BlockHeader().Time.Add(duration))
		}
	}
	configuration := types.ExecutionConfiguration{}
	if msg.Configuration != nil {
		configuration = *msg.Configuration
	}

	conditions := types.ExecutionConditions{}
	if msg.Conditions != nil {
		if msg.Conditions.UseResponseValue != nil {
			if int(msg.Conditions.UseResponseValue.MsgsIndex) >= len(msg.Msgs) {
				if msg.Msgs[0].TypeUrl == sdk.MsgTypeURL(&authztypes.MsgExec{}) {
					msgExec := &authztypes.MsgExec{}
					if err := proto.Unmarshal(msg.Msgs[0].Value, msgExec); err != nil {
						return nil, errorsmod.Wrapf(types.ErrInvalidRequest, "msg exec could not unmarshal, can not check conditions")
					}

					if int(msg.Conditions.UseResponseValue.MsgsIndex) >= len(msgExec.Msgs) {
						return nil, errorsmod.Wrapf(types.ErrInvalidRequest, "msgs index: %v must be shorter than length msgExec msgs array: %s", msg.Conditions.UseResponseValue.MsgsIndex, msgExec.Msgs)

					} else {
						return nil, errorsmod.Wrapf(types.ErrInvalidRequest, "msgs index: %v must be shorter than length msgs array: %s", msg.Conditions.UseResponseValue.MsgsIndex, msg.Msgs)
					}
				}
			}
		}
		if msg.Conditions.ResponseComparison != nil {
			if int(msg.Conditions.ResponseComparison.ResponseIndex) >= len(msg.Msgs) {
				return nil, errorsmod.Wrapf(types.ErrInvalidRequest, "response index: %v must be shorter than length msgs array: %s", msg.Conditions.ResponseComparison.ResponseIndex, msg.Msgs)

			}
		}

		if msg.Conditions.ICQConfig != nil {
			if msg.Conditions.ICQConfig.TimeoutDuration != 0 {
				if msg.Conditions.ICQConfig.TimeoutDuration > duration || msg.Conditions.ICQConfig.TimeoutDuration > interval {
					return nil, errorsmod.Wrapf(types.ErrInvalidRequest, "TimeoutDuration must be shorter than the action interval or duration")
				}
			}
			if msg.Conditions.UseResponseValue == nil || !msg.Conditions.UseResponseValue.FromICQ {
				if msg.Conditions.ResponseComparison == nil || !msg.Conditions.ResponseComparison.FromICQ {
					return nil, errorsmod.Wrapf(types.ErrInvalidRequest, "Query must be used in UseResponseValue or ResponseComparison")
				}
			}

		}
		conditions = *msg.Conditions
	}

	hostedConfig := types.HostedConfig{}
	if msg.HostedConfig != nil {
		hostedConfig = *msg.HostedConfig
	}

	err = k.CreateAction(ctx, msgOwner, msg.Label, msg.Msgs, duration, interval, startTime, msg.FeeFunds, configuration, hostedConfig, portID, msg.ConnectionId, msg.HostConnectionId, conditions)
	if err != nil {
		return nil, err
	}

	return &types.MsgSubmitActionResponse{}, nil
}

// RegisterAccountAndSubmitAction implements the Msg/RegisterAccountAndSubmitAction interface
func (k msgServer) RegisterAccountAndSubmitAction(goCtx context.Context, msg *types.MsgRegisterAccountAndSubmitAction) (*types.MsgRegisterAccountAndSubmitActionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := k.RegisterInterchainAccount(ctx, msg.ConnectionId, msg.Owner, msg.Version)
	if err != nil {
		return nil, err
	}

	portID, err := icatypes.NewControllerPortID(msg.Owner)
	if err != nil {
		return nil, err
	}

	var duration time.Duration = 0
	if msg.Duration != "" {
		duration, err = time.ParseDuration(msg.Duration)
		if err != nil {
			return nil, err
		}
	}

	var interval time.Duration = 0
	if msg.Interval != "" {
		interval, err = time.ParseDuration(msg.Interval)
		if err != nil {
			return nil, errorsmod.Wrap(types.ErrInvalidRequest, err.Error())
		}
	}

	var startTime time.Time = ctx.BlockHeader().Time
	if msg.StartAt != 0 {
		startTime = time.Unix(int64(msg.StartAt), 0)
		if err != nil {
			return nil, err
		}
		if startTime.Before(ctx.BlockHeader().Time.Add(time.Minute)) {
			return nil, errorsmod.Wrapf(types.ErrInvalidRequest, "custom start time: %s must be at least a minute into the future upon block submission: %s", startTime, ctx.BlockHeader().Time.Add(time.Minute))
		}
	}

	msgOwner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrInvalidRequest, err.Error())
	}

	p := k.GetParams(ctx)
	if interval != 0 && (interval < p.MinActionInterval || interval > duration) {
		return nil, errorsmod.Wrapf(types.ErrInvalidRequest, "interval: %s  must be longer than minimum interval:  %s, and longer than duration: %s", interval, p.MinActionInterval, duration)
	}
	if duration != 0 {
		if duration > p.MaxActionDuration {
			return nil, errorsmod.Wrapf(types.ErrInvalidRequest, "duration: %s must be shorter than maximum duration: %s", duration, p.MaxActionDuration)
		} else if duration < p.MinActionDuration {
			return nil, errorsmod.Wrapf(types.ErrInvalidRequest, "duration: %s must be longer than minimum duration: %s", duration, p.MinActionDuration)
		} else if startTime.After(ctx.BlockHeader().Time.Add(p.MaxActionDuration)) {
			return nil, errorsmod.Wrapf(types.ErrInvalidRequest, "start time: %s must be before current time and maximum duration: %s", startTime, ctx.BlockHeader().Time.Add(duration))
		}
	}
	configuration := types.ExecutionConfiguration{}
	if msg.Configuration != nil {
		configuration = *msg.Configuration
	}

	conditions := types.ExecutionConditions{}
	if msg.Conditions != nil {
		conditions = *msg.Conditions
	}
	err = k.CreateAction(ctx, msgOwner, msg.Label, msg.Msgs, duration, interval, startTime, msg.FeeFunds, configuration, types.HostedConfig{}, portID, msg.ConnectionId, msg.HostConnectionId, conditions)
	if err != nil {
		return nil, err
	}

	return &types.MsgRegisterAccountAndSubmitActionResponse{}, nil
}

// UpdateAction implements the Msg/UpdateAction interface
func (k msgServer) UpdateAction(goCtx context.Context, msg *types.MsgUpdateAction) (*types.MsgUpdateActionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	action, err := k.TryGetActionInfo(ctx, msg.ID)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrInvalidRequest, err.Error())
	}
	if action.Owner != msg.Owner {
		return nil, types.ErrInvalidAddress
	}

	if action.Configuration.UpdatingDisabled {
		return nil, errorsmod.Wrap(types.ErrInvalidRequest, "updating is disabled")
	}

	if msg.ConnectionId != "" {
		action.ICAConfig.PortID, err = icatypes.NewControllerPortID(msg.Owner)
		if err != nil {
			return nil, err
		}
		action.ICAConfig.ConnectionID = msg.ConnectionId
	}
	newExecTime := action.ExecTime
	if msg.EndTime > 0 {
		endTime := time.Unix(int64(msg.EndTime), 0)
		if err != nil {
			return nil, err
		}
		if endTime.Before(ctx.BlockTime().Add(time.Minute * 2)) {
			return nil, types.ErrInvalidTime
		}
		action.EndTime = endTime

		if action.Interval != 0 && msg.Interval != "" {
			newExecTime = endTime
		}
	}
	p := k.GetParams(ctx)

	if msg.Interval != "" {
		interval, err := time.ParseDuration(msg.Interval)
		if err != nil {
			return nil, errorsmod.Wrap(types.ErrUpdateAction, err.Error())
		}

		if interval != 0 && interval < p.MinActionInterval || interval > action.EndTime.Sub(action.StartTime) {
			return nil, errorsmod.Wrapf(types.ErrUpdateAction, "interval: %s  must be longer than minimum interval:  %s, and execution should happen before end time", interval, p.MinActionInterval)
		}
		action.Interval = interval
	}

	if msg.StartAt > 0 {
		startTime := time.Unix(int64(msg.StartAt), 0)
		if err != nil {
			return nil, err
		}
		if startTime.Before(ctx.BlockHeader().Time.Add(time.Minute)) {
			return nil, errorsmod.Wrapf(types.ErrInvalidRequest, "custom start time: %s must be at least a minute into the future upon block submission: %s", startTime, ctx.BlockHeader().Time.Add(time.Minute))
		}
		if startTime.After(action.EndTime) {
			return nil, errorsmod.Wrapf(types.ErrUpdateAction, "start time: %s must be before end time", startTime)
		}
		latestEntry, err := k.GetLatestActionHistoryEntry(ctx, action.ID)
		if err != nil || latestEntry != nil {
			return nil, errorsmod.Wrapf(types.ErrUpdateAction, "start time: %s must occur before first execution", startTime)
		}
		action.StartTime = startTime
		newExecTime = startTime
	}

	if msg.Label != "" {
		action.Label = msg.Label
	}

	if msg.Configuration != nil {
		action.Configuration = msg.Configuration
	}

	if msg.Conditions != nil {
		if msg.Conditions.UseResponseValue != nil {
			if int(msg.Conditions.UseResponseValue.MsgsIndex) >= len(msg.Msgs) {
				if msg.Msgs[0].TypeUrl == sdk.MsgTypeURL(&authztypes.MsgExec{}) {
					msgExec := &authztypes.MsgExec{}
					if err := proto.Unmarshal(msg.Msgs[0].Value, msgExec); err != nil {
						return nil, errorsmod.Wrapf(types.ErrInvalidRequest, "msg exec could not unmarshal")
					}

					if int(msg.Conditions.UseResponseValue.MsgsIndex) >= len(msgExec.Msgs) {
						return nil, errorsmod.Wrapf(types.ErrInvalidRequest, "msgs index: %v must be shorter than length msgExec msgs array: %s", msg.Conditions.UseResponseValue.MsgsIndex, msgExec.Msgs)

					} else {
						return nil, errorsmod.Wrapf(types.ErrInvalidRequest, "msgs index: %v must be shorter than length msgs array: %s", msg.Conditions.UseResponseValue.MsgsIndex, msg.Msgs)
					}
				}
			}
		}
		action.Conditions = msg.Conditions
	}
	if len(msg.Msgs) != 0 {
		action.Msgs = msg.Msgs
	}

	if newExecTime != action.ExecTime {
		k.RemoveFromActionQueue(ctx, action)
		action.ExecTime = newExecTime
		k.InsertActionQueue(ctx, action.ID, newExecTime)
	}

	action.UpdateHistory = append(action.UpdateHistory, ctx.BlockTime())

	if !action.ActionAuthzSignerOk(k.cdc) {
		return nil, errorsmod.Wrapf(types.ErrUpdateAction, "action signer: %s is not message signer", action.Owner)
	}

	//set hosted config
	if msg.HostedConfig != nil {
		action.HostedConfig = msg.HostedConfig
	}

	k.SetActionInfo(ctx, &action)

	return &types.MsgUpdateActionResponse{}, nil
}

// CreateHostedAccount implements the Msg/CreateHostedAccount interface
func (k msgServer) CreateHostedAccount(goCtx context.Context, msg *types.MsgCreateHostedAccount) (*types.MsgCreateHostedAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	hostedAddress, err := DeriveHostedAddress(msg.Creator, msg.ConnectionId)
	if err != nil {
		return nil, err
	}
	//register ICA
	err = k.RegisterInterchainAccount(ctx, msg.ConnectionId, hostedAddress.String(), msg.Version)
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
	k.SetHostedAccount(ctx, &types.HostedAccount{HostedAddress: hostedAddress.String(), HostFeeConfig: &types.HostFeeConfig{Admin: msg.Creator, FeeCoinsSuported: msg.FeeCoinsSuported}, ICAConfig: &types.ICAConfig{ConnectionID: msg.ConnectionId, HostConnectionID: msg.HostConnectionId, PortID: portID}})
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
	connectionID := hostedAccount.ICAConfig.ConnectionID
	hostConnectionID := hostedAccount.ICAConfig.HostConnectionID
	if msg.ConnectionId != "" && msg.HostConnectionId != "" {
		// hostedAddressAcc, err := DeriveHostedAddress(hostedAccount.HostedAddress, msg.ConnectionId)
		// if err != nil {
		// 	return nil, err
		// }
		// hostedAddress = hostedAddressAcc.String()
		// portID, err := icatypes.NewControllerPortID(hostedAddress)
		// if err != nil {
		// 	return nil, err
		// }
		connectionID = msg.ConnectionId
		hostConnectionID = msg.HostConnectionId
	}

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

	k.SetHostedAccount(ctx, &types.HostedAccount{HostedAddress: hostedAddress, HostFeeConfig: &types.HostFeeConfig{Admin: admin, FeeCoinsSuported: feeCoinsSupported}, ICAConfig: &types.ICAConfig{ConnectionID: connectionID, HostConnectionID: hostConnectionID, PortID: hostedAccount.ICAConfig.PortID}})

	//set hosted config by address on hosted key prefix
	return &types.MsgUpdateHostedAccountResponse{}, nil
}
