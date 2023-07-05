package keeper

import (
	"context"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/gogoproto/proto"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
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
	err := k.RegisterInterchainAccount(ctx, msg.ConnectionId, msg.Owner, msg.Version)
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

	data, err := icatypes.SerializeCosmosTx(k.cdc, []proto.Message{msg.GetTxMsg()})
	if err != nil {
		return nil, err
	}

	packetData := icatypes.InterchainAccountPacketData{
		Type: icatypes.EXECUTE_TX,
		Data: data,
	}
	timeoutTimestamp := ctx.BlockTime().Add(time.Minute).UnixNano()
	sequence, err := k.icaControllerKeeper.SendTx(ctx, nil, msg.ConnectionId, portID, packetData, uint64(timeoutTimestamp))
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

	msgOwner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
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
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
		}
	}

	var startTime time.Time = ctx.BlockHeader().Time
	if msg.StartAt != 0 {
		startTime = time.Unix(int64(msg.StartAt), 0)
		if err != nil {
			return nil, err
		}
		if startTime.Before(ctx.BlockHeader().Time.Add(time.Minute)) {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "custom start time: %s must be at least a minute into the future upon block submission: %s", startTime, ctx.BlockHeader().Time.Add(time.Minute))
		}
	}

	p := k.GetParams(ctx)
	if interval != 0 && (interval < p.MinAutoTxInterval || interval > duration) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "interval: %s  must be longer than minimum interval:  %s, and longer than duration: %s", interval, p.MinAutoTxInterval, duration)
	}
	if duration != 0 {
		if duration > p.MaxAutoTxDuration {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "duration: %s must be shorter than maximum duration: %s", duration, p.MaxAutoTxDuration)
		} else if duration < p.MinAutoTxDuration {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "duration: %s must be longer than minimum duration: %s", duration, p.MinAutoTxDuration)
		} else if startTime.After(ctx.BlockHeader().Time.Add(p.MaxAutoTxDuration)) {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "start time: %s must be before current time and max duration: %s", startTime, ctx.BlockHeader().Time.Add(duration))
		}
	}
	if len(msg.DependsOnTxIds) >= 10 || len(msg.Msgs) >= 10 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "must depend on less than 10 autoTxIDs and have less than 10 messages")
	}
	// if msg.Retries > 5 {
	// 	return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "can retry for a maximum of 5 times")
	// }

	err = k.CreateAutoTx(ctx, msgOwner, msg.Label, portID, msg.Msgs, msg.ConnectionId, duration, interval, startTime, msg.FeeFunds, msg.DependsOnTxIds)
	if err != nil {
		return nil, err
	}

	return &types.MsgSubmitAutoTxResponse{}, nil
}

// RegisterAccountAndSubmitAutoTx implements the Msg/RegisterAccountAndSubmitAutoTx interface
func (k msgServer) RegisterAccountAndSubmitAutoTx(goCtx context.Context, msg *types.MsgRegisterAccountAndSubmitAutoTx) (*types.MsgRegisterAccountAndSubmitAutoTxResponse, error) {
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
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
		}
	}

	var startTime time.Time = ctx.BlockHeader().Time
	if msg.StartAt != 0 {
		startTime = time.Unix(int64(msg.StartAt), 0)
		if err != nil {
			return nil, err
		}
		if startTime.Before(ctx.BlockHeader().Time.Add(time.Minute)) {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "custom start time: %s must be at least a minute into the future upon block submission: %s", startTime, ctx.BlockHeader().Time.Add(time.Minute))
		}
	}

	msgOwner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	p := k.GetParams(ctx)
	if interval != 0 && (interval < p.MinAutoTxInterval || interval > duration) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "interval: %s  must be longer than minimum interval:  %s, and longer than duration: %s", interval, p.MinAutoTxInterval, duration)
	}
	if duration != 0 {
		if duration > p.MaxAutoTxDuration {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "duration: %s must be shorter than maximum duration: %s", duration, p.MaxAutoTxDuration)
		} else if duration < p.MinAutoTxDuration {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "duration: %s must be longer than minimum duration: %s", duration, p.MinAutoTxDuration)
		} else if startTime.After(ctx.BlockHeader().Time.Add(p.MaxAutoTxDuration)) {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "start time: %s must be before current time and maximum duration: %s", startTime, ctx.BlockHeader().Time.Add(duration))
		}
	}
	if len(msg.DependsOnTxIds) >= 10 || len(msg.Msgs) >= 10 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "must depend on less than 10 autoTxIDs and have less than 10 messages")
	}
	// if msg.Retries > 5 {
	// 	return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "can retry for a maximum of 5 times")
	// }

	err = k.CreateAutoTx(ctx, msgOwner, msg.Label, portID, msg.Msgs, msg.ConnectionId, duration, interval, startTime, msg.FeeFunds, msg.DependsOnTxIds)
	if err != nil {
		return nil, err
	}

	return &types.MsgRegisterAccountAndSubmitAutoTxResponse{}, nil
}

// UpdateAutoTx implements the Msg/UpdateAutoTx interface
func (k msgServer) UpdateAutoTx(goCtx context.Context, msg *types.MsgUpdateAutoTx) (*types.MsgUpdateAutoTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	autoTx, err := k.TryGetAutoTxInfo(ctx, msg.TxId)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}
	if autoTx.Owner != msg.Owner {
		return nil, sdkerrors.ErrInvalidAddress
	}

	if msg.ConnectionId != "" {
		autoTx.PortID, err = icatypes.NewControllerPortID(msg.Owner)
		if err != nil {
			return nil, err
		}
		autoTx.ConnectionID = msg.ConnectionId
	}
	newExecTime := autoTx.ExecTime
	if msg.EndTime > 0 {
		endTime := time.Unix(int64(msg.EndTime), 0)
		if err != nil {
			return nil, err
		}
		if endTime.Before(ctx.BlockTime().Add(time.Minute * 2)) {
			return nil, types.ErrInvalidTime
		}
		autoTx.EndTime = endTime

		if autoTx.Interval != 0 && msg.Interval != "" {
			newExecTime = endTime
		}
	}
	p := k.GetParams(ctx)

	//var interval time.Duration = 0
	if msg.Interval != "" {
		interval, err := time.ParseDuration(msg.Interval)
		if err != nil {
			return nil, sdkerrors.Wrap(types.ErrUpdateAutoTx, err.Error())
		}

		if interval != 0 && interval < p.MinAutoTxInterval || interval > autoTx.EndTime.Sub(autoTx.StartTime) {
			return nil, sdkerrors.Wrapf(types.ErrUpdateAutoTx, "interval: %s  must be longer than minimum interval:  %s, and execution should happen before end time", interval, p.MinAutoTxInterval)
		}
		autoTx.Interval = interval
		//newExecTime := interval
	}

	if msg.StartAt > 0 {
		startTime := time.Unix(int64(msg.StartAt), 0)
		if err != nil {
			return nil, err
		}
		if startTime.Before(ctx.BlockHeader().Time.Add(time.Minute)) {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "custom start time: %s must be at least a minute into the future upon block submission: %s", startTime, ctx.BlockHeader().Time.Add(time.Minute))
		}
		if startTime.After(autoTx.EndTime) {
			return nil, sdkerrors.Wrapf(types.ErrUpdateAutoTx, "start time: %s must be before end time", startTime)
		}
		// if startTime.After(autoTx.ExecTime) {
		// 	return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "start time: %s must be before next AutoTx exec time", startTime)
		// }
		if len(autoTx.AutoTxHistory) != 0 {
			return nil, sdkerrors.Wrapf(types.ErrUpdateAutoTx, "start time: %s must occur before first execution", startTime)
		}
		autoTx.StartTime = startTime
		newExecTime = startTime
	}
	if len(msg.DependsOnTxIds) >= 10 || len(msg.Msgs) >= 10 {
		return nil, sdkerrors.Wrap(types.ErrUpdateAutoTx, "must depend on less than 10 autoTxIDs and have less than 10 messages")
	}

	if msg.Label != "" {
		autoTx.Label = msg.Label
	}

	if len(msg.DependsOnTxIds) != 0 {
		autoTx.DependsOnTxIds = msg.DependsOnTxIds
	}

	if len(msg.Msgs) != 0 {
		autoTx.Msgs = msg.Msgs
	}

	if newExecTime != autoTx.ExecTime {
		k.RemoveFromAutoTxQueue(ctx, autoTx)
		autoTx.ExecTime = newExecTime
		k.InsertAutoTxQueue(ctx, autoTx.TxID, newExecTime)
	}

	autoTx.UpdateHistory = append(autoTx.UpdateHistory, ctx.BlockTime())

	k.SetAutoTxInfo(ctx, &autoTx)

	return &types.MsgUpdateAutoTxResponse{}, nil
}
