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
		case *types.MsgCreateEstimation:
			return handleMsgCreateEstimation(ctx, k, msg)

		case *types.MsgUpdateLike:
			return handleMsgUpdateLike(ctx, k, msg)

		case *types.MsgDeleteEstimation:
			return handleMsgDeleteEstimation(ctx, k, msg)

		case *types.MsgPrepayment:
			return handleMsgPrepayment(ctx, k, msg)

		case *types.MsgWithdrawal:
			return handleMsgWithdrawal(ctx, k, msg)

		case *types.MsgCreateItem:
			return handleMsgCreateItem(ctx, k, msg)

		case *types.MsgDeleteItem:
			return handleMsgDeleteItem(ctx, k, msg)

		case *types.MsgRevealEstimation:
			return handleMsgRevealEstimation(ctx, k, msg)

		case *types.MsgItemTransferable:
			return handleMsgItemTransferable(ctx, k, msg)

		case *types.MsgItemShipping:
			return handleMsgItemShipping(ctx, k, msg)

		case *types.MsgItemTransfer:
			return handleMsgItemTransfer(ctx, k, msg)

		case *types.MsgFlagItem:
			return handleMsgFlagItem(ctx, k, msg)
		case *types.MsgItemRating:
			return handleMsgItemRating(ctx, k, msg)
		case *types.MsgItemResell:
			return handleMsgItemResell(ctx, k, msg)
		case *types.MsgTokenizeItem:
			return handleMsgTokenizeItem(ctx, k, msg)
		case *types.MsgUnTokenizeItem:
			return handleMsgUnTokenizeItem(ctx, k, msg)

		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}
