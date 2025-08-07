package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/trstlabs/intento/x/claim/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) ClaimClaimable(goCtx context.Context, msg *types.MsgClaimClaimable) (*types.MsgClaimClaimableResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	claimer, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	err = k.ClaimClaimableForAddr(ctx, claimer)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeClaim),
		),
	)

	return &types.MsgClaimClaimableResponse{}, nil
}

func (m msgServer) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_, err := sdk.AccAddressFromBech32(msg.GetAuthority())
	if err != nil {
		return nil, err
	}

	if msg.GetAuthority() != m.Keeper.GetAuthority() {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "sender address is not authorized address to update module params")
	}

	err = msg.GetParams().Validate() // need to explicitly validate as x/gov invokes this msg and it does not validate
	if err != nil {
		return nil, err
	}

	err = m.SetParams(ctx, msg.GetParams())
	if err != nil {
		return nil, err
	}

	return &types.MsgUpdateParamsResponse{}, nil
}
