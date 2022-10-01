package keeper

import (
	"context"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/trstlabs/trst/x/compute/internal/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	keeper Keeper
}

func NewMsgServerImpl(k Keeper) types.MsgServer {
	return &msgServer{keeper: k}
}

func (m msgServer) StoreCode(goCtx context.Context, msg *types.MsgStoreCode) (*types.MsgStoreCodeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		sdk.NewAttribute(types.AttributeKeySigner, msg.Sender),
	))

	duration, err := time.ParseDuration(msg.DefaultDuration)

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
	}

	interval, err := time.ParseDuration(msg.DefaultInterval)

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
	}

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	codeID, err := m.keeper.Create(ctx, sender, msg.WASMByteCode, msg.Source, msg.Builder, duration, interval, msg.Title, msg.Description)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(types.AttributeKeyCodeID, fmt.Sprintf("%d", codeID)),
		),
	})

	return &types.MsgStoreCodeResponse{
		CodeID: codeID,
	}, nil
}

func (m msgServer) InstantiateContract(goCtx context.Context, msg *types.MsgInstantiateContract) (*types.MsgInstantiateContractResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	duration, err := time.ParseDuration(msg.Duration)
	p := m.keeper.GetParams(ctx)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
	}
	if duration > p.MaxContractDuration {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "contract duration must be shorter than maximum duration")
	}
	if duration != 0 && duration < p.MinContractDuration {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "contract duration must be longer than minimum duration")
	}
	interval, err := time.ParseDuration(msg.Interval)

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
	}
	if interval != 0 && interval < p.MinContractInterval {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "contract interval must be longer than minimum interval")
	}
	var startTime time.Time = ctx.BlockHeader().Time
	if msg.StartDurationAt != 0 {
		startTime = time.Unix(int64(msg.StartDurationAt), 0)
		if err != nil {
			return nil, err
		}
	}
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	contractAddr, data, err := m.keeper.Instantiate(ctx, msg.CodeID, sender, msg.Msg, msg.AutoMsg, msg.ContractId, msg.Funds, msg.CallbackSig, duration, interval, startTime)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		sdk.NewAttribute(types.AttributeKeyContractAddr, contractAddr.String()),
	))

	return &types.MsgInstantiateContractResponse{
		Address: contractAddr.String(),
		Data:    data,
	}, nil
}

func (m msgServer) ExecuteContract(goCtx context.Context, msg *types.MsgExecuteContract) (*types.MsgExecuteContractResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		sdk.NewAttribute(types.AttributeKeyContractAddr, msg.Contract),
	))
	contract, err := sdk.AccAddressFromBech32(msg.Contract)
	if err != nil {
		return nil, err
	}
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	data, err := m.keeper.Execute(ctx, contract, sender, msg.Msg, msg.Funds, msg.CallbackSig)
	if err != nil {
		return nil, err
	}

	return &types.MsgExecuteContractResponse{
		Data: data.Data,
	}, nil
}
