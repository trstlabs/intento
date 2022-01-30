package trst

import (
	"fmt"

	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/trstlabs/trst/x/item/keeper"
	"github.com/trstlabs/trst/x/item/types"
	//"github.com/tendermint/tendermint/crypto"
)

func handleMsgCreateItem(ctx sdk.Context, k keeper.Keeper, msg *types.MsgCreateItem) (*sdk.Result, error) {

	err := k.CreateItem(ctx, *msg)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, err.Error())
	}
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgDeleteItem(ctx sdk.Context, k keeper.Keeper, msg *types.MsgDeleteItem) (*sdk.Result, error) {
	if !k.HasItem(ctx, msg.Id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("item id %s doesn't exist", strconv.FormatUint(msg.Id, 10)))
	}

	item := k.GetItem(ctx, msg.Id)

	if item.Transfer.Buyer != "" && item.Status != "" {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Item has a buyer")
	}

	//status but no buyer
	if item.Status != "" && item.Transfer.Buyer == "" {

		// title,status and rating are kept for record keeping
		item.Description = ""

		item.Estimation = &types.Estimation{}
		item.Transfer.ShippingCost = 0
		item.Transfer.LocalPickup = ""
		item.Transfer.Buyer = ""
		item.Transfer.Tracking = false
		item.Transfer.ShippingRegion = nil
		item.Transfer.Note = ""
		item.Properties.Photos = nil
		item.Properties.Tags = nil
		item.Properties.Condition = 0
		item.Properties.Transferable = false

		k.SetItem(ctx, item)

	} else {

		k.RemoveFromListedItemQueue(ctx, msg.Id, item.ListingDuration.EndTime)
		_ = k.DeleteItemContract(ctx, item.Estimation.Contract)
		k.DeleteItem(ctx, msg.Id)
		k.RemoveFromSellerItems(ctx, msg.Id, msg.Seller)
	}

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgRevealEstimation(ctx sdk.Context, k keeper.Keeper, msg *types.MsgRevealEstimation) (*sdk.Result, error) {

	item := k.GetItem(ctx, msg.Itemid)

	if msg.Creator != item.Transfer.Seller {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "not item seller")
	}

	if item.Estimation.EstimationCount != item.Estimation.EstimationTotal {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "not enough estimators")
	}

	err := k.RevealEstimation(ctx, item, *msg)
	if err != nil {
		//	fmt.Printf("err executing")
		//	fmt.Printf("executing contract: %X\n", err)
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, err.Error()) ///panic(err)
	}

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgItemTransferable(ctx sdk.Context, k keeper.Keeper, msg *types.MsgItemTransferable) (*sdk.Result, error) {
	//check if item exists
	item := k.GetItem(ctx, msg.Itemid)

	//check if message sender is item seller
	if item.Properties.EstimationOnly {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "item is estimation only")
	}

	//check if message sender is item buyer
	if msg.Seller != item.Transfer.Seller {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect seller")
	}

	//check if item is a token
	if item.Properties.IsToken {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Item is a token")
	}

	//check if item has a best estimator (and therefore a complete estimation)
	if item.Estimation.BestEstimator == "" {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "no estimation price yet, cannot make item transferable")
	}

	err := k.SetTransferable(ctx, item, *msg)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, err.Error())
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.EventTypeItemTransferable, sdk.NewAttribute(types.AttributeKeyItemID, strconv.FormatUint(msg.Itemid, 10))))

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgItemShipping(ctx sdk.Context, k keeper.Keeper, msg *types.MsgItemShipping) (*sdk.Result, error) {

	//get item info
	item := k.GetItem(ctx, msg.Itemid)

	//check if message seller is item seller
	if msg.Seller != item.Transfer.Seller {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect seller")
	}

	//check if item.transferable = true and therefore the seller has accepted the buyer
	if !item.Properties.Transferable {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "item is not transferable")
	}

	//check if item has a buyer already (so that we know that prepayment is done)
	if item.Transfer.Buyer == "" {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "item does not have a buyer yet")
	}
	//bonus check if item already has been transferred
	if item.Status != "" {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "item already has had a transfer or transfer has been denied ")
	}
	if len(item.Transfer.ShippingRegion) == 0 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "unauthrorized, no shipping region")
	}

	bigIntEstimationPrice := sdk.NewInt(item.Estimation.EstimationPrice - item.Transfer.ShippingCost)
	if item.Transfer.Discount > 0 {
		bigIntEstimationPrice = sdk.NewInt(item.Estimation.EstimationPrice - item.Transfer.Discount)
	}

	bigIntShipping := sdk.NewInt(item.Transfer.ShippingCost)
	if msg.Tracking {
		if item.Creator == item.Transfer.Seller {

			maxRewardCoin := sdk.NewCoin("utrst", sdk.NewInt(item.Transfer.ShippingCost))
			k.HandleBuyerReward(ctx, maxRewardCoin, sdk.AccAddress(item.Transfer.Buyer))

			item.Estimation.BestEstimator = ""
			item.Estimation.EstimatorList = nil
		}
		//make payment to seller
		CreaterPayoutAndShipping := bigIntEstimationPrice.Add(bigIntShipping)
		paymentSellerCoin := sdk.NewCoin("utrst", CreaterPayoutAndShipping)

		k.SendPaymentToAccount(ctx, item.Transfer.Seller, paymentSellerCoin)
		k.RemoveFromListedItemQueue(ctx, item.Id, item.ListingDuration.EndTime)
		item.Status = "Shipped"
		k.SetItem(ctx, item)
		//k.SetBuyer(ctx, buyer)
	} else {
		repayment := bigIntEstimationPrice.Add(bigIntShipping)
		repaymentCoin := sdk.NewCoin("utrst", repayment)

		k.SendPaymentToAccount(ctx, item.Transfer.Buyer, repaymentCoin)

		item.Status = "Shipping declined; buyer refunded"
		k.SetItem(ctx, item)
		//k.SetBuyer(ctx, buyer)
	}

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgItemResell(ctx sdk.Context, k keeper.Keeper, msg *types.MsgItemResell) (*sdk.Result, error) {

	// Checks that the element exists
	if !k.HasItem(ctx, msg.Itemid) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("item id %s doesn't exist", strconv.FormatUint(msg.Itemid, 10)))
	}

	item := k.GetItem(ctx, msg.Itemid)

	// Checks if the the msg sender is the same as the current buyer
	if msg.Seller != item.Transfer.Buyer {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect buyer")
	}

	if item.Status != "Transferred" && item.Status != "Shipped" {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Item not available to resell")
	}

	if msg.Discount > item.Estimation.EstimationPrice {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Discount invalid")
	}
	k.DeleteBuyerKey(ctx, append(types.Uint64ToByte(msg.Itemid), []byte(msg.Seller)...))

	item.Transfer.Seller = msg.Seller
	item.Transfer.ShippingCost = msg.ShippingCost
	item.Transfer.LocalPickup = msg.LocalPickup
	item.Transfer.ShippingRegion = msg.ShippingRegion
	item.Transfer.Discount = msg.Discount
	item.Transfer.Note = msg.Note
	item.Transfer.Rating = 0
	item.Transfer.Buyer = ""
	item.Status = ""
	item.ListingDuration.EndTime = ctx.BlockTime().Add(types.DefaultParams().MaxActivePeriod)

	k.SetItem(ctx, item)
	k.InsertListedItemQueue(ctx, item)
	k.BindItemToSellerItems(ctx, msg.Itemid, msg.Seller)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.EventTypeItemResellable, sdk.NewAttribute(types.AttributeKeyItemID, strconv.FormatUint(msg.Itemid, 10))))
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgTokenizeItem(ctx sdk.Context, k keeper.Keeper, msg *types.MsgTokenizeItem) (*sdk.Result, error) {

	// Checks that the element exists
	if !k.HasItem(ctx, msg.Id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("item id %s doesn't exist", strconv.FormatUint(msg.Id, 10)))
	}

	item := k.GetItem(ctx, msg.Id)

	// Checks if the the msg sender is the same as the current buyer or creator
	if msg.Sender != item.Transfer.Buyer && msg.Sender != item.Creator {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect buyer")
	}
	//if item was not transferred previously
	if item.Status != "Transferred" && item.Status != "Shipped" {
		//and if item is made transferable or it has no best estimator
		if item.Properties.Transferable || item.Estimation.BestEstimator == "" {
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "item must have an estimation and not be transferable ")
		}
	}
	//Create new coin
	//Set item to tokenized
	err := k.TokenizeItem(ctx, msg.Id, (msg.Sender))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, err.Error())
	}
	item.Properties.IsToken = true
	k.RemoveFromListedItemQueue(ctx, item.Id, item.ListingDuration.EndTime)
	k.SetItem(ctx, item)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.EventTypeItemTokenized, sdk.NewAttribute(types.AttributeKeyItemID, strconv.FormatUint(msg.Id, 10))))
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgUnTokenizeItem(ctx sdk.Context, k keeper.Keeper, msg *types.MsgUnTokenizeItem) (*sdk.Result, error) {

	// Checks that the element exists
	if !k.HasItem(ctx, msg.Id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("item id %s doesn't exist", strconv.FormatUint(msg.Id, 10)))
	}

	item := k.GetItem(ctx, msg.Id)

	err := k.UnTokenizeItem(ctx, msg.Id, msg.Sender)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Item not available to untokenize")

	}
	item.Transfer.Buyer = msg.Sender
	item.Properties.IsToken = false
	k.SetItem(ctx, item)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.EventTypeItemUnTokenized, sdk.NewAttribute(types.AttributeKeyItemID, strconv.FormatUint(msg.Id, 10))))
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}
