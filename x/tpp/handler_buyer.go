package tpp

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/danieljdd/tpp/x/tpp/keeper"
	"github.com/danieljdd/tpp/x/tpp/types"
	//"github.com/tendermint/tendermint/crypto"
	//bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
)

func handleMsgCreateBuyer(ctx sdk.Context, k keeper.Keeper, msg *types.MsgCreateBuyer) (*sdk.Result, error) {
	//get item info

	item := k.GetItem(ctx, msg.Itemid)

	//check if item has a best estimator (and therefore a complete estimation)
	if item.Estimationprice == 0 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "item does not have estimation yet, cannot make prepayment")
	}

	//check if item is transferable
	if item.Transferable != true {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "item  not transferable, cannot make prepayment")
	}

	//check if item has a buyer already
	if item.Buyer != "" {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "item already has a buyer, cannot make prepayment")
	}

	//item buyer cannot be the item creator
	if msg.Buyer == item.Creator || msg.Buyer == item.Seller  {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Buyer cannot be creator/seller")
	}

	estimationPrice := item.Estimationprice
	if item.Discount > 0 {
		estimationPrice = item.Estimationprice - item.Discount
	}
	
	if item.Shippingcost > 0 && item.Localpickup == false {
		toPayShipping := estimationPrice + item.Shippingcost
		if toPayShipping != msg.Deposit {

			return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "deposit insufficient, cannot make prepayment")
		}

		item.Buyer = msg.Buyer

		k.SetItem(ctx, item)
		k.CreateBuyer(ctx, *msg)
		//}
	}

	if item.Shippingcost == 0 && item.Localpickup == true {
		
		if estimationPrice != msg.Deposit {

			return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "deposit insufficient, cannot make prepayment")
		}

		item.Buyer = msg.Buyer

		k.SetItem(ctx, item)
		k.CreateBuyer(ctx, *msg)

	}

	if item.Shippingcost > 0 && item.Localpickup == true {
	
		if estimationPrice == msg.Deposit {

			//ModuleAcct := sdk.AccAddress(crypto.AddressHash([]byte(types.ModuleName)))

			item.Shippingcost = 0
			item.Buyer = msg.Buyer

			k.SetItem(ctx, item)
			k.CreateBuyer(ctx, *msg)

		} else {
			toPayShipping :=
			estimationPrice + item.Shippingcost

			if toPayShipping == msg.Deposit {
				item.Localpickup = false
				item.Buyer = msg.Buyer

				k.SetItem(ctx, item)
				k.CreateBuyer(ctx, *msg)
			}
			if toPayShipping != msg.Deposit {

				return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "deposit insufficient, cannot make prepayment")
			}
		}

	}

	//k.CreateBuyer(ctx, *msg)

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgUpdateBuyer(ctx sdk.Context, k keeper.Keeper, msg *types.MsgUpdateBuyer) (*sdk.Result, error) {
	deposit := sdk.NewInt64Coin("tpp", msg.Deposit)

	var buyer = types.Buyer{
		Buyer: msg.Buyer,
		Itemid:       msg.Itemid,
		Transferable: msg.Transferable,
		Deposit:      deposit,
	}

	// Checks that the element exists
	if !k.HasBuyer(ctx, msg.Itemid) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %s doesn't exist", msg.Itemid))
	}

	// Checks if the the msg sender is the same as the current owner
	if msg.Buyer != k.GetBuyerOwner(ctx, msg.Itemid) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.SetBuyer(ctx, buyer)

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgDeleteBuyer(ctx sdk.Context, k keeper.Keeper, msg *types.MsgDeleteBuyer) (*sdk.Result, error) {
	if !k.HasBuyer(ctx, msg.Itemid) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %s doesn't exist", msg.Itemid))
	}
	if msg.Buyer != k.GetBuyerOwner(ctx, msg.Itemid) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.DeleteBuyer(ctx, msg.Itemid)

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgItemTransfer(ctx sdk.Context, k keeper.Keeper, msg *types.MsgItemTransfer) (*sdk.Result, error) {
	//check if message creator is item creator
	if msg.Buyer != k.GetBuyerOwner(ctx, msg.Itemid) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	//get item info
	item := k.GetItem(ctx, msg.Itemid)

	//check if item.transferable = true and therefore the seller has accepted the buyer

	if item.Transferable == false {
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

	if msg.Transferable == true {
		
		bigintestimationprice := sdk.NewInt(item.Estimationprice)
		if item.Discount > 0 {
			bigintestimationprice = sdk.NewInt(item.Estimationprice - item.Discount)
		}

		if (item.Creator == item.Seller) {
		//rounded down percentage for minting. Percentage may be changed through governance proposals
		percentageMint := sdk.NewDecWithPrec(10, 2)
		percentageReward := sdk.NewDecWithPrec(5, 2)


		toMintAmount := percentageMint.MulInt(bigintestimationprice).TruncateInt()
		paymentReward := percentageReward.MulInt(bigintestimationprice)
		roundedAmountReward := paymentReward.TruncateInt()

		
		//minted coins (are rounded up)
		mintCoins := sdk.NewCoin("tpp", toMintAmount)
		paymentRewardCoins := sdk.NewCoin("tpp", roundedAmountReward)

		k.MintReward(ctx, mintCoins)
		

		//for their participation in the protocol, the best estimator and the buyer get rewarded.
		k.HandlePrepayment(ctx, item.Bestestimator, paymentRewardCoins)
		k.HandlePrepayment(ctx, item.Buyer, paymentRewardCoins)

		//refund the deposits back to all of the item estimators
		for _, element := range item.Estimatorlist {
			key := msg.Itemid + "-" + element

			k.DeleteEstimator(ctx, key)
		}

		item.Bestestimator = ""
		item.Lowestestimator = ""
		item.Highestestimator = ""
		item.Estimatorlist = nil
	}
		//make payment to seller
		paymentSellerCoins := sdk.NewCoin("tpp", bigintestimationprice)

		k.HandlePrepayment(ctx, item.Seller, paymentSellerCoins)

		item.Status = "Transferred"
		k.SetItem(ctx, item)
		//k.SetBuyer(ctx, buyer)
	}

	if msg.Transferable == false {
		buyer := k.GetBuyer(ctx, msg.Itemid)
		k.HandlePrepayment(ctx, item.Buyer, buyer.Deposit)

		for _, element := range item.Estimatorlist {
			//apply this to each element
			key := msg.Itemid + "-" + element
			estimator := k.GetEstimator(ctx, key)

			if estimator.Estimator == item.Highestestimator {
				
					k.BurnCoins(ctx, estimator.Deposit)
			

			} else {
				k.DeleteEstimator(ctx, key)
			}

		}

		item.Bestestimator = ""
		item.Lowestestimator = ""
		item.Highestestimator = ""
		item.Estimatorlist = nil
		item.Status = "Item transfer declined"
		k.SetItem(ctx, item)
		//k.SetBuyer(ctx, buyer)
	}

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgItemThank(ctx sdk.Context, k keeper.Keeper, msg *types.MsgItemThank) (*sdk.Result, error) {
	//check if message seller is item seller
	if msg.Buyer != k.GetBuyerOwner(ctx, msg.Itemid) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}
	//get buyer info
	//buyer := k.GetBuyer(ctx, msg.Itemid)

	//get item info
	item := k.GetItem(ctx, msg.Itemid)

	//check if item.transferable = false

	if item.Transferable == true {
		return nil, sdkerrors.Wrap(nil, "seller of item does not accept a transfer")
	}

	//check if item has a buyer already
	if item.Buyer != msg.Buyer {
		return nil, sdkerrors.Wrap(nil, "Item does not belong to msg sender")
	}
	//check therefore that prepayment is done
	if item.Status != "Shipped" || item.Status != "Transferred" {
		return nil, sdkerrors.Wrap(nil, "item not transferred yet")
	}
	item.Thank = msg.Thank
	
	k.SetItem(ctx, item)
	
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}
