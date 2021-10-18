package trst

import (
	"fmt"

	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/danieljdd/trst/x/trst/keeper"
	"github.com/danieljdd/trst/x/trst/types"
	//"github.com/tendermint/tendermint/crypto"
	//bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
)

func handleMsgPrepayment(ctx sdk.Context, k keeper.Keeper, msg *types.MsgPrepayment) (*sdk.Result, error) {
	//get item info

	item := k.GetItem(ctx, msg.Itemid)

	//check if item has a best estimator (and therefore a complete estimation)
	if item.Estimationprice == 0 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "item does not have estimation yet, cannot make prepayment")
	}

	//check if item is transferable
	if !item.Transferable {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "item  not transferable, cannot make prepayment")
	}

	//check if item has a buyer already
	if item.Buyer != "" {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "item already has a buyer, cannot make prepayment")
	}

	//item buyer cannot be the item creator
	if msg.Buyer == item.Creator || msg.Buyer == item.Seller {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Buyer cannot be creator/seller")
	}

	estimationPrice := item.Estimationprice
	if item.Discount > 0 {
		estimationPrice = item.Estimationprice - item.Discount
	}

	if item.Shippingcost > 0 && item.Localpickup == "" {
		toPayShipping := estimationPrice + item.Shippingcost
		if toPayShipping != msg.Deposit {

			return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "deposit insufficient, cannot make prepayment")
		}

		item.Buyer = msg.Buyer
		k.RemoveFromListedItemQueue(ctx, msg.Itemid, item.Endtime)
		k.SetItem(ctx, item)
		k.Prepayment(ctx, *msg)
		//}
	}

	if item.Shippingcost == 0 && item.Localpickup != "" {

		if estimationPrice != msg.Deposit {

			return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "deposit insufficient, cannot make prepayment")
		}

		item.Buyer = msg.Buyer
		k.RemoveFromListedItemQueue(ctx, msg.Itemid, item.Endtime)
		k.SetItem(ctx, item)
		k.Prepayment(ctx, *msg)

	}

	if item.Shippingcost > 0 && item.Localpickup != "" {

		if estimationPrice == msg.Deposit {

			//ModuleAcct := sdk.AccAddress(crypto.AddressHash([]byte(types.ModuleName)))

			item.Shippingcost = 0
			item.Buyer = msg.Buyer
			k.RemoveFromListedItemQueue(ctx, msg.Itemid, item.Endtime)
			k.SetItem(ctx, item)
			k.Prepayment(ctx, *msg)

		} else {
			toPayShipping :=
				estimationPrice + item.Shippingcost

			if toPayShipping == msg.Deposit {
				item.Localpickup = ""
				item.Buyer = msg.Buyer
				k.RemoveFromListedItemQueue(ctx, msg.Itemid, item.Endtime)
				k.SetItem(ctx, item)
				k.Prepayment(ctx, *msg)
			}
			if toPayShipping != msg.Deposit {

				return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "deposit insufficient, cannot make prepayment")
			}
		}

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
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	if !item.Transferable {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "item was deleted")
	}

	//if the item has a status delete the buyer upon request, if the item is not transferred, return a part of prepayment. After rating the item, the buyer gets fully refunded
	if item.Status != "" {
		item.Buyer = ""
		k.SetItem(ctx, item)
		k.Withdrawal(ctx, append([]byte(msg.Buyer), types.Uint64ToByte(msg.Itemid)...))
	} else {

		//returning the trst tokens
		percentageReturn := sdk.NewDecWithPrec(95, 2)

		bigIntEstimationPrice := sdk.NewInt(item.Estimationprice)
		toMintAmount := percentageReturn.MulInt(bigIntEstimationPrice).TruncateInt()

		burnAmount := bigIntEstimationPrice.Sub(toMintAmount)
		k.BurnCoins(ctx, sdk.NewCoin("utrst", burnAmount))

		if item.Shippingcost > 0 {
			toMintAmount = toMintAmount.Add(sdk.NewInt(item.Shippingcost))
		}

		//minted coins (are rounded up)
		mintCoins := sdk.NewCoin("utrst", toMintAmount)
		k.HandlePrepayment(ctx, msg.Buyer, mintCoins)
		k.Withdrawal(ctx, append([]byte(msg.Buyer), types.Uint64ToByte(msg.Itemid)...))

		for _, element := range item.Estimatorlist {
			//apply this to each element
			key := append(types.Uint64ToByte(msg.Itemid), []byte(element)...)

			//	estimator := k.GetEstimator(ctx, key)

			//	if estimator.Estimator == item.Highestestimator {

			//	k.BurnCoins(ctx, estimator.Deposit)
			//	k.DeleteEstimationWithoutDeposit(ctx, key)

			//} else {
			k.DeleteEstimation(ctx, key)
			//	}

		}

		item.Status = "Withdrawal prepayment"
		item.Shippingcost = 0
		item.Localpickup = ""
		item.Estimationcount = 0
		item.Bestestimator = ""
		item.Lowestestimator = ""
		item.Highestestimator = ""

		item.Estimatorlist = nil
		item.Estimationlist = nil
		item.Transferable = false

		k.SetItem(ctx, item)
		k.Withdrawal(ctx, append([]byte(msg.Buyer), types.Uint64ToByte(msg.Itemid)...))

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
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "prepayment does not belong to msg sender")
	}
	//check therefore that prepayment is done
	if item.Status != "" {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "item already has had a transfer or transfer has been denied ")
	}
	if item.Shippingcost > 0 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "item has shippingcost")
	}

	bigIntEstimationPrice := sdk.NewInt(item.Estimationprice - item.Depositamount)
	if item.Discount > 0 {
		bigIntEstimationPrice = sdk.NewInt(item.Estimationprice - item.Discount)
	}

	if item.Creator == item.Seller {

		//minted coins (are rounded up)
		rewardCoins := sdk.NewCoin("utrst", sdk.NewInt(item.Depositamount))

		k.MintReward(ctx, rewardCoins)

		//for their participation in the protocol, the best estimator, a random and the stakers get rewarded.
		k.HandleEstimatorReward(ctx, item.Bestestimator, rewardCoins)

		k.HandleStakingReward(ctx, rewardCoins)

		//refund the deposits back to all of the item estimators
		for _, element := range item.Estimatorlist {
			key := append(types.Uint64ToByte(msg.Itemid), []byte(element)...)

			k.DeleteEstimation(ctx, key)
		}

		item.Bestestimator = ""
		item.Lowestestimator = ""
		item.Highestestimator = ""
		item.Estimatorlist = nil
	}
	//make payment to seller
	paymentSellerCoins := sdk.NewCoin("utrst", bigIntEstimationPrice)

	k.HandlePrepayment(ctx, item.Seller, paymentSellerCoins)

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
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
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
		bigIntEstimationPrice := sdk.NewInt(item.Estimationprice)
		toMintAmount := percentageReward.MulInt(bigIntEstimationPrice).Ceil().TruncateInt()

		//minted coins (are rounded up)
		mintCoins := sdk.NewCoin("utrst", toMintAmount)
		k.MintReward(ctx, mintCoins.Add(mintCoins))
		k.HandlePrepayment(ctx, msg.Buyer, mintCoins)
		item.Estimationprice = 0

	} else if msg.Rating < 3 {
		item.Buyer = ""
	}
	item.Note = msg.Note
	item.Rating = msg.Rating

	k.SetItem(ctx, item)

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}
