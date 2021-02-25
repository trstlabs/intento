package tpp

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/danieljdd/tpp/x/tpp/keeper"
	"github.com/danieljdd/tpp/x/tpp/types"
	//"github.com/tendermint/tendermint/crypto"
)

func handleMsgCreateItem(ctx sdk.Context, k keeper.Keeper, msg *types.MsgCreateItem) (*sdk.Result, error) {
	k.CreateItem(ctx, *msg)

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgUpdateItem(ctx sdk.Context, k keeper.Keeper, msg *types.MsgUpdateItem) (*sdk.Result, error) {
	var item = types.Item{
		Creator:                     msg.Creator,
		Id:                          msg.Id,
		Title:                       msg.Title,
		Description:                 msg.Description,
		Shippingcost:                msg.Shippingcost,
		Localpickup:                 msg.Localpickup,
		Condition:                   msg.Condition,
		Shippingregion:              msg.Shippingregion,
	}

	// Checks that the element exists
	if !k.HasItem(ctx, msg.Id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %s doesn't exist", msg.Id))
	}

	// Checks if the the msg sender is the same as the current owner
	if msg.Creator != k.GetItemOwner(ctx, msg.Id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.SetItem(ctx, item)

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgDeleteItem(ctx sdk.Context, k keeper.Keeper, msg *types.MsgDeleteItem) (*sdk.Result, error) {
	if !k.HasItem(ctx, msg.Id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %s doesn't exist", msg.Id))
	}
	if msg.Creator != k.GetItemOwner(ctx, msg.Id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	item:= k.GetItem(ctx, msg.Id)
	



	if item.Status != "" && item.Buyer == "" {

		for _, element := range item.Estimatorlist {
			//apply this to each element
			key := msg.Id + "-" + element

			k.DeleteEstimator(ctx, key)
		}
	
		item.Title = "Deleted"
		item.Description = ""
		item.Shippingcost = 0
		item.Localpickup = false
		item.Estimationcounthash = ""
		item.Bestestimator = ""
		item.Lowestestimator = ""
		item.Highestestimator = ""
		item.Estimationprice = 0
		item.Estimatorlist = nil
		item.Estimatorestimationhashlist = nil
		item.Transferable = false
		item.Buyer = ""
		item.Tracking = false
		item.Comments = nil
		item.Tags = ""
		item.Condition = 0
		item.Shippingregion = ""
		

		k.SetItem(ctx, item)

	}else{

	//if estimation is made pay back all the estimators/or buyer (like handlerItemTransfer)
	//has to be a new function
	for _, element := range item.Estimatorlist {
		//apply this to each element
		key := msg.Id + "-" + element

		k.DeleteEstimator(ctx, key)

		
	}

	k.DeleteItem(ctx, msg.Id)
	}

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}
