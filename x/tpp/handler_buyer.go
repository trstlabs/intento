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
	if item.Bestestimator == "" {
		return nil, sdkerrors.Wrap(nil, "item does not have estimation yet, cannot make prepayment")
	}

	//check if item is transferable
	if item.Transferable != true {
		return nil, sdkerrors.Wrap(nil, "item  not transferable, cannot make prepayment")
	}

	//check if item has a buyer already
	if item.Buyer != "" {
		return nil, sdkerrors.Wrap(nil, "item already has a buyer, cannot make prepayment")
	}

	/*	ToPayLocal := item.EstimationPrice
		DepositCoinsLocal := sdk.NewCoins(sdk.NewCoin("token", ToPayLocal))

		ToPayShipping :=
			sdk.Int.Add(item.EstimationPrice, item.ShippingCost)
		DepositCoinsShipping := sdk.NewCoins(sdk.NewCoin("token", ToPayShipping))
	*/

	//ModuleAcct := sdk.AccAddress(crypto.AddressHash([]byte(types.ModuleName)))

	if item.Shippingcost > 0 && item.Localpickup == false {
		toPayShipping := item.Estimationprice + item.Shippingcost
		if toPayShipping != msg.Deposit {

			return nil, sdkerrors.Wrap(nil, "deposit insufficient, cannot make prepayment")
		}

		item.Buyer = msg.Buyer

		k.SetItem(ctx, item)
		k.CreateBuyer(ctx, *msg)
		//}
	}

	if item.Shippingcost == 0 && item.Localpickup == true {
		toPayLocal := item.Estimationprice
		if toPayLocal != msg.Deposit {

			return nil, sdkerrors.Wrap(nil, "deposit insufficient, cannot make prepayment")
		}

		item.Buyer = msg.Buyer

		k.SetItem(ctx, item)
		k.CreateBuyer(ctx, *msg)

	}

	if item.Shippingcost > 0 && item.Localpickup == true {
		toPayLocal := item.Estimationprice
		if toPayLocal == msg.Deposit {

			//ModuleAcct := sdk.AccAddress(crypto.AddressHash([]byte(types.ModuleName)))

			item.Shippingcost = 0
			item.Buyer = msg.Buyer

			k.SetItem(ctx, item)
			k.CreateBuyer(ctx, *msg)

		} else {
			toPayShipping :=
				item.Estimationprice + item.Shippingcost

			if toPayShipping == msg.Deposit {
				item.Localpickup = false
				item.Buyer = msg.Buyer

				k.SetItem(ctx, item)
				k.CreateBuyer(ctx, *msg)
			}
			if toPayShipping != msg.Deposit {

				return nil, sdkerrors.Wrap(nil, "deposit insufficient, cannot make prepayment")
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
	//get buyer info
	buyer := k.GetBuyer(ctx, msg.Itemid)

	//get item info
	item := k.GetItem(ctx, msg.Itemid)

	//check if item.transferable = true and therefore the creator has accepted the buyer

	if item.Transferable == false {
		return nil, sdkerrors.Wrap(nil, "creator of item does not accept a transfer")
	}

	//check if item has a buyer already
	//check therefore that prepayment is done
	if item.Buyer != msg.Buyer {
		return nil, sdkerrors.Wrap(nil, "prepayment does not belong to msg sender")
	}
	//check therefore that prepayment is done
	if item.Status != "" {
		return nil, sdkerrors.Wrap(nil, "item already has had a transfer or transfer has been denied ")
	}

	if msg.Transferable == true {
		bigintestimationprice := sdk.NewInt(item.Estimationprice)

		//rounded down percentage for minting. Percentage may be changed through governance proposals
		percentageMint := sdk.NewDecWithPrec(10, 2)
		percentageReward := sdk.NewDecWithPrec(5, 2)

		toMintAmount := percentageMint.MulInt(bigintestimationprice).TruncateInt()
		paymentReward := percentageReward.MulInt(bigintestimationprice)
		roundedAmountReward := paymentReward.Ceil().TruncateInt()

		//make payment to creator and estimator
		paymentCreatorCoins := sdk.NewCoin("tpp", bigintestimationprice)

		//minted coins (are rounded up)
		mintCoins := sdk.NewCoin("tpp", toMintAmount)
		paymentRewardCoins := sdk.NewCoin("tpp", roundedAmountReward)

		k.MintReward(ctx, mintCoins)
		k.HandlePrepayment(ctx, item.Creator, paymentCreatorCoins)

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
		item.Status = "Item transferred"
		k.SetItem(ctx, item)
		k.SetBuyer(ctx, buyer)
	}

	if msg.Transferable == false {

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
		k.SetBuyer(ctx, buyer)
	}

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}
