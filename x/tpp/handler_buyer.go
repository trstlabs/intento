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
		return nil, sdkerrors.Wrap(nil,"item  not transferable, cannot make prepayment")
	}

	//check if item has a buyer already
	if item.Buyer != "" {
		return nil, sdkerrors.Wrap(nil,"item already has a buyer, cannot make prepayment")
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
		depositCoinsShipping := sdk.NewInt64Coin("tpp", toPayShipping)
		equal := depositCoinsShipping.IsEqual(msg.Deposit)
		if equal == false {
			return nil, sdkerrors.Wrap(nil, "deposit insufficient, cannot make prepayment")
		}
		
		item.Buyer = msg.Buyer

		k.SetItem(ctx, item)
		k.CreateBuyer(ctx, *msg)
		//}
	}

	if item.Shippingcost == 0 && item.Localpickup == true {
		toPayLocal := item.Estimationprice
		depositCoinsLocal := sdk.NewInt64Coin("tpp", toPayLocal)
		equallocal := depositCoinsLocal.IsEqual(msg.Deposit)
		if equallocal == false {
			return nil, sdkerrors.Wrap(nil, "deposit insufficient, cannot make prepayment")
		}
	
		item.Buyer = msg.Buyer
	
		k.SetItem(ctx, item)
		k.CreateBuyer(ctx, *msg)

	}

	if item.Shippingcost > 0 && item.Localpickup == true {
		toPayLocal := item.Estimationprice
		depositCoinsLocal := sdk.NewInt64Coin("tpp", toPayLocal)
		equallocal := depositCoinsLocal.IsEqual(msg.Deposit)
		if equallocal == true {
			//ModuleAcct := sdk.AccAddress(crypto.AddressHash([]byte(types.ModuleName)))
			
			item.Shippingcost = 0
			item.Buyer = msg.Buyer
		
			k.SetItem(ctx, item)
			k.CreateBuyer(ctx, *msg)

		} else {
		toPayShipping :=
			item.Estimationprice + item.Shippingcost
		depositCoinsShipping := sdk.NewInt64Coin("tpp", toPayShipping)
		equalshipping := depositCoinsShipping.IsEqual(msg.Deposit)
		if equalshipping == true {
			item.Localpickup = false
			item.Buyer = msg.Buyer
		
			k.SetItem(ctx, item)
			k.CreateBuyer(ctx, *msg)
		} }
		
	} 
	
		//k.CreateBuyer(ctx, *msg)

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}


func handleMsgUpdateBuyer(ctx sdk.Context, k keeper.Keeper, msg *types.MsgUpdateBuyer) (*sdk.Result, error) {
	var buyer = types.Buyer{
		Buyer:      msg.Buyer,

		Itemid:       msg.Itemid,
		Transferable: msg.Transferable,
		Deposit:      msg.Deposit,
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
