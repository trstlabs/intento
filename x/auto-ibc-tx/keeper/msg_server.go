package keeper

import (
	"context"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	icatypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/types"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v3/modules/core/24-host"
	"github.com/trstlabs/trst/x/auto-ibc-tx/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	Keeper
}

// NewMsgServerImpl creates and returns a new types.MsgServer, fulfilling the auto-ibc-tx Msg service interface
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

// RegisterAccount implements the Msg/RegisterAccount interface
func (k msgServer) RegisterAccount(goCtx context.Context, msg *types.MsgRegisterAccount) (*types.MsgRegisterAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	err := k.RegisterInterchainAccount(ctx, msg.ConnectionId, msg.Owner)
	if err != nil {
		return nil, err
	}
	return &types.MsgRegisterAccountResponse{}, nil
}

// SubmitTx implements the Msg/SubmitTx interface
func (k msgServer) SubmitTx(goCtx context.Context, msg *types.MsgSubmitTx) (*types.MsgSubmitTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	portID, err := icatypes.NewControllerPortID(msg.Owner)
	if err != nil {
		return nil, err
	}

	channelID, found := k.icaControllerKeeper.GetActiveChannelID(ctx, msg.ConnectionId, portID)
	if !found {
		return nil, sdkerrors.Wrapf(icatypes.ErrActiveChannelNotFound, "failed to retrieve active channel for port %s", portID)
	}

	chanCap, found := k.scopedKeeper.GetCapability(ctx, host.ChannelCapabilityPath(portID, channelID))
	if !found {
		return nil, sdkerrors.Wrap(channeltypes.ErrChannelCapabilityNotFound, "module does not own channel capability")
	}

	data, err := icatypes.SerializeCosmosTx(k.cdc, []sdk.Msg{msg.GetTxMsg()})
	if err != nil {
		return nil, err
	}

	packetData := icatypes.InterchainAccountPacketData{
		Type: icatypes.EXECUTE_TX,
		Data: data,
	}
	timeoutTimestamp := time.Now().Add(time.Minute).UnixNano()
	sequence, err := k.icaControllerKeeper.SendTx(ctx, chanCap, msg.ConnectionId, portID, packetData, uint64(timeoutTimestamp))
	if err != nil {
		return nil, err
	}
	//store 0 as autoTx id as a regular submit is not autoTx
	k.setTmpAutoTxID(ctx, 0, portID, sequence)
	return &types.MsgSubmitTxResponse{}, nil
}

// SubmitAutoTx implements the Msg/SubmitAutoTx interface
func (k msgServer) SubmitAutoTx(goCtx context.Context, msg *types.MsgSubmitAutoTx) (*types.MsgSubmitAutoTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	portID, err := icatypes.NewControllerPortID(msg.Owner)
	if err != nil {
		return nil, err
	}

	// check if the msg contains valid inputs
	err = msg.GetTxMsg().ValidateBasic()
	if err != nil {
		return nil, err
	}

	data, err := icatypes.SerializeCosmosTx(k.cdc, []sdk.Msg{msg.GetTxMsg()})
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
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
		}
	}

	var startTime time.Time = ctx.BlockHeader().Time
	if msg.StartAt != 0 {
		startTime = time.Unix(int64(msg.StartAt), 0)
		if err != nil {
			return nil, err
		}
	}

	msgOwner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
	}

	p := k.GetParams(ctx)
	if interval != 0 && interval < p.MinAutoTxInterval && interval > duration {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "AutoTx interval: %s  must be longer than minimum interval:  %s, and longer than duration: %s", interval, p.MinAutoTxInterval, duration)

	}
	if duration != 0 {
		if duration > p.MaxAutoTxDuration {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "AutoTx duration: %s must be shorter than maximum duration: %s", duration, p.MaxAutoTxDuration)
		}
		if duration < p.MinAutoTxDuration {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "AutoTx duration: %s must be longer than minimum duration: %s", duration, p.MinAutoTxDuration)
		}
		if startTime.After(ctx.BlockHeader().Time.Add(duration)) {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "AutoTx start time: %s must be before AutoTx end time : %s", startTime, ctx.BlockHeader().Time.Add(duration))
		}

	}
	if len(msg.DependsOnTxIds) >= 10 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "AutoTx must depend on less than 10 autoTxIDs")
	}
	if msg.Retries > 5 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "AutoTx can retry for a maximum of 5 times")
	}

	err = k.CreateAutoTx(ctx, msgOwner, portID, data, msg.ConnectionId, duration, interval, startTime, msg.FeeFunds, msg.Retries, msg.DependsOnTxIds)
	if err != nil {
		return nil, err
	}

	return &types.MsgSubmitAutoTxResponse{}, nil
}

// SubmitAutoTx implements the Msg/SubmitAutoTx interface
func (k msgServer) RegisterAccountAndSubmitAutoTx(goCtx context.Context, msg *types.MsgRegisterAccountAndSubmitAutoTx) (*types.MsgRegisterAccountAndSubmitAutoTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := k.RegisterInterchainAccount(ctx, msg.ConnectionId, msg.Owner)
	if err != nil {
		return nil, err
	}

	portID, err := icatypes.NewControllerPortID(msg.Owner)
	if err != nil {
		return nil, err
	}

	data, err := icatypes.SerializeCosmosTx(k.cdc, []sdk.Msg{msg.GetTxMsg()})
	if err != nil {
		return nil, err
	}

	// check if the msg contains valid inputs
	err = msg.GetTxMsg().ValidateBasic()
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
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
		}
	}

	var startTime time.Time = ctx.BlockHeader().Time
	if msg.StartAt != 0 {
		startTime = time.Unix(int64(msg.StartAt), 0)
		if err != nil {
			return nil, err
		}
	}

	msgOwner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
	}

	p := k.GetParams(ctx)
	if interval != 0 && interval < p.MinAutoTxInterval && interval > duration {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "AutoTx interval: %s  must be longer than minimum interval:  %s, and longer than duration: %s", interval, p.MinAutoTxInterval, duration)

	}
	if duration != 0 {
		if duration > p.MaxAutoTxDuration {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "AutoTx duration: %s must be shorter than maximum duration: %s", duration, p.MaxAutoTxDuration)
		}
		if duration < p.MinAutoTxDuration {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "AutoTx duration: %s must be longer than minimum duration: %s", duration, p.MinAutoTxDuration)
		}
		if startTime.After(ctx.BlockHeader().Time.Add(duration)) {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "AutoTx start time: %s must be before AutoTx end time : %s", startTime, ctx.BlockHeader().Time.Add(duration))
		}

	}
	if len(msg.DependsOnTxIds) >= 10 {

		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "AutoTx must depend on less than 10 autoTxIDs")

	}
	if msg.Retries > 5 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "AutoTx can retry for a maximum of 5 times")
	}

	err = k.CreateAutoTx(ctx, msgOwner, portID, data, msg.ConnectionId, duration, interval, startTime, msg.FeeFunds, msg.Retries, msg.DependsOnTxIds)
	if err != nil {
		return nil, err
	}

	return &types.MsgRegisterAccountAndSubmitAutoTxResponse{}, nil
}
