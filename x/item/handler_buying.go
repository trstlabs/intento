package trst

import (
	"fmt"

	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/trstlabs/trst/x/item/keeper"
	"github.com/trstlabs/trst/x/item/types"
	//"github.com/tendermint/tendermint/crypto"
	//bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
)

func handleMsgPrepayment(ctx sdk.Context, k keeper.Keeper, msg *types.MsgPrepayment) (*sdk.Result, error) {
	//get item info

	item := k.GetItem(ctx, msg.Itemid)

	//check if item has a best estimator (and therefore a complete estimation)
	if item.EstimationPrice == 0 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "item does not have a price yet, cannot make prepayment")
	}

	//check if item is transferable
	if !item.Transferable {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "item not transferable, cannot make prepayment")
	}

	//check if item has a buyer already
	if item.Buyer != "" {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "item has a buyer, cannot make prepayment")
	}

	//item buyer cannot be the item creator
	if msg.Buyer == item.Creator || msg.Buyer == item.Seller {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "buyer cannot be creator/seller")
	}

	estimationPrice := item.EstimationPrice
	if item.Discount > 0 {
		estimationPrice = item.EstimationPrice - item.Discount
	}

	if item.ShippingCost > 0 && item.LocalPickup == "" {
		toPayShipping := estimationPrice + item.ShippingCost
		if toPayShipping != msg.Deposit {

			return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "deposit insufficient, cannot make prepayment")
		}

	}

	if item.ShippingCost == 0 && item.LocalPickup != "" {

		if estimationPrice != msg.Deposit {

			return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "deposit insufficient, cannot make prepayment")
		}

	}

	if item.ShippingCost > 0 && item.LocalPickup != "" {

		if estimationPrice == msg.Deposit {

			item.ShippingCost = 0

		} else {
			toPayShipping :=
				estimationPrice + item.ShippingCost

			if toPayShipping == msg.Deposit {
				item.LocalPickup = ""

			} else {

				return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "deposit insufficient, cannot make prepayment")
			}
		}

	}
	item.Buyer = msg.Buyer
	k.RemoveFromListedItemQueue(ctx, msg.Itemid, item.EndTime)
	k.SetItem(ctx, item)
	err := k.Prepayment(ctx, *msg)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.EventTypeItemPrepayment, sdk.NewAttribute(types.AttributeKeyItemID, strconv.FormatUint(msg.Itemid, 10))))
	//k.Prepayment(ctx, *msg)

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgWithdrawal(ctx sdk.Context, k keeper.Keeper, msg *types.MsgWithdrawal) (*sdk.Result, error) {
	if !k.HasBuyer(ctx, append([]byte(msg.Buyer), types.Uint64ToByte(msg.Itemid)...)) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("buyer of id %s doesn't exist", strconv.FormatUint(msg.Itemid, 10)))
	}

	//get item info
	item := k.GetItem(ctx, msg.Itemid)

	if msg.Buyer != item.Buyer {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect item buyer")
	}

	if !item.Transferable {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "item was deleted")
	}

	//if the item has a status, withdrawl upon request. if the item is not transferred, return a part of prepayment. After rating the item, the buyer gets fully refunded
	if item.Status != "" {
		item.Buyer = ""
		k.SetItem(ctx, item)
		k.DeleteBuyerKey(ctx, append([]byte(msg.Buyer), types.Uint64ToByte(msg.Itemid)...))
	} else {

		//returning the trst tokens minus a percentage that will be returned after placing a review
		percentageReturn := sdk.NewDecWithPrec(95, 2)

		bigIntEstimationPrice := sdk.NewInt(item.EstimationPrice)
		toMintAmount := percentageReturn.MulInt(bigIntEstimationPrice).TruncateInt()

		burnAmount := bigIntEstimationPrice.Sub(toMintAmount)
		k.BurnCoins(ctx, sdk.NewCoin("utrst", burnAmount))

		if item.ShippingCost > 0 {
			toMintAmount = toMintAmount.Add(sdk.NewInt(item.ShippingCost))
		}

		//minted coins (are rounded up)
		mintCoin := sdk.NewCoin("utrst", toMintAmount)
		k.SendPaymentToAccount(ctx, msg.Buyer, mintCoin)
		k.DeleteBuyerKey(ctx, append([]byte(msg.Buyer), types.Uint64ToByte(msg.Itemid)...))

		item.Status = "Withdrawal prepayment"
		item.ShippingCost = 0
		item.LocalPickup = ""
		item.EstimationCount = 0
		item.BestEstimator = ""
		item.EstimatorList = nil
		item.EstimationList = nil
		item.Transferable = false

		k.SetItem(ctx, item)
		k.DeleteBuyerKey(ctx, append([]byte(msg.Buyer), types.Uint64ToByte(msg.Itemid)...))

	}

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgItemTransfer(ctx sdk.Context, k keeper.Keeper, msg *types.MsgItemTransfer) (*sdk.Result, error) {

	//get item info
	item := k.GetItem(ctx, msg.Itemid)

	//check if item.transferable = true and therefore the seller has accepted the buyer

	if !item.Transferable {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "seller of item does not accept a transfer")
	}

	//check if item has a buyer already
	//check therefore that prepayment is done
	if item.Buyer != msg.Buyer {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "item has a different buyer")
	}
	//check therefore that prepayment is done
	if item.Status != "" {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "item already has had a transfer or transfer has been denied ")
	}
	if item.ShippingCost > 0 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "item has shippingcost")
	}

	bigIntEstimationPrice := sdk.NewInt(item.EstimationPrice - item.DepositAmount)
	if item.Discount > 0 {
		bigIntEstimationPrice = sdk.NewInt(item.EstimationPrice - item.Discount)
	}

	if item.Creator == item.Seller {

		//minted coins (are rounded up)
		maxRewardCoin := sdk.NewCoin("utrst", sdk.NewInt(item.DepositAmount))

		k.HandleBuyerReward(ctx, maxRewardCoin, sdk.AccAddress(msg.Buyer))

		item.BestEstimator = ""
		item.EstimatorList = nil
	}
	//make payment to seller
	paymentSellerCoins := sdk.NewCoin("utrst", bigIntEstimationPrice)

	k.SendPaymentToAccount(ctx, item.Seller, paymentSellerCoins)

	item.Status = "Transferred"
	k.SetItem(ctx, item)
	//k.SetBuyer(ctx, buyer)

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgItemRating(ctx sdk.Context, k keeper.Keeper, msg *types.MsgItemRating) (*sdk.Result, error) {

	//get item info
	item := k.GetItem(ctx, msg.Itemid)

	//check if msg buyer is item buyer
	if msg.Buyer != item.Buyer {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect buyer")
	}

	//check if item has a buyer already
	if item.Buyer != msg.Buyer {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "item does not belong to msg sender")
	}
	//check if the item has a status, and therefore a buyer has a reason to rate the item
	if item.Status == "" {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "item does not have a status")
	}

	if item.Status == "Withdrawal prepayment" {
		percentageReward := sdk.NewDecWithPrec(5, 2)
		bigIntEstimationPrice := sdk.NewInt(item.EstimationPrice)
		toMintAmount := percentageReward.MulInt(bigIntEstimationPrice).Ceil().TruncateInt()

		//mint the remaining payment (5%), that was burned when withdrawing payment
		mintCoin := sdk.NewCoin("utrst", toMintAmount)
		k.MintReward(ctx, mintCoin.Add(mintCoin))
		k.SendPaymentToAccount(ctx, msg.Buyer, mintCoin)
		item.EstimationPrice = 0

	} else if msg.Rating < 3 {
		//if the rating is low, we hide the buyer
		item.Buyer = ""
	}
	item.Note = msg.Note
	item.Rating = msg.Rating

	k.SetItem(ctx, item)

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}
