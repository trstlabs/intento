package tpp

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/danieljdd/tpp/x/tpp/keeper"
	"github.com/danieljdd/tpp/x/tpp/types"
)

// NewHandler ...
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		// this line is used by starport scaffolding # 1
		case *types.MsgCreateEstimator:
			return handleMsgCreateEstimator(ctx, k, msg)

		case *types.MsgUpdateEstimator:
			return handleMsgUpdateEstimator(ctx, k, msg)

		case *types.MsgDeleteEstimator:
			return handleMsgDeleteEstimator(ctx, k, msg)

		case *types.MsgCreateBuyer:
			return handleMsgCreateBuyer(ctx, k, msg)

		case *types.MsgUpdateBuyer:
			return handleMsgUpdateBuyer(ctx, k, msg)

		case *types.MsgDeleteBuyer:
			return handleMsgDeleteBuyer(ctx, k, msg)

		case *types.MsgCreateItem:
			return handleMsgCreateItem(ctx, k, msg)

		case *types.MsgUpdateItem:
			return handleMsgUpdateItem(ctx, k, msg)

		case *types.MsgDeleteItem:
			return handleMsgDeleteItem(ctx, k, msg)

		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}
