package tpp

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"

	//"encoding/hex"
	//"github.com/tendermint/tendermint/crypto"

	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/danieljdd/tpp/x/tpp/keeper"
	"github.com/danieljdd/tpp/x/tpp/types"
)

func handleMsgCreateEstimator(ctx sdk.Context, k keeper.Keeper, msg *types.MsgCreateEstimator) (*sdk.Result, error) {

	item := k.GetItem(ctx, msg.Itemid)

	if item.Estimationprice > 0 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "item already has an estimation")
	}

	///for production: check if estimator is item owner
	if msg.Estimator == item.Creator {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "estimator cannot be item creator")
	}
	
	//checks whether the estimator already estimated the item

	if k.HasEstimator(ctx, msg.Itemid+"-"+msg.Estimator) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "already estimated this item")
	}

	if msg.Deposit != item.Depositamount {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "deposit invalid")
	}

	//checks whether estimationcount will be reached
	var estimatorlistlen = strconv.Itoa(len(item.Estimatorlist) + 1)
	var estimatorlistlenhash = sha256.Sum256([]byte(estimatorlistlen + item.Seller))
	var estimatorlisthashstring = hex.EncodeToString(estimatorlistlenhash[:])
	if estimatorlisthashstring == item.Estimationcounthash {
		item.Bestestimator = "Awaiting"
	
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(types.EventTypeItemReady, sdk.NewAttribute(types.AttributeKeyItemID, msg.Itemid)),
		)
	}

	var estimatorestimationhash = sha256.Sum256([]byte(strconv.FormatInt(msg.Estimation, 10) + msg.Estimator))
	var estimatorestimationhashstring = hex.EncodeToString(estimatorestimationhash[:])

	//append estimatorhash to list
	item.Estimatorestimationhashlist = append(item.Estimatorestimationhashlist, estimatorestimationhashstring)

	item.Estimatorlist = append(item.Estimatorlist, msg.Estimator)

	k.SetItem(ctx, item)

	k.CreateEstimator(ctx, *msg)


	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgUpdateEstimator(ctx sdk.Context, k keeper.Keeper, msg *types.MsgUpdateEstimator) (*sdk.Result, error) {
	var estimator = types.Estimator{
		Estimator: msg.Estimator,

		Itemid: msg.Itemid,

		Interested: msg.Interested,
	}

	// Checks that the element exists
	if !k.HasEstimator(ctx, msg.Itemid+"-"+msg.Estimator) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %s doesn't exist", msg.Itemid+"-"+msg.Estimator))
	}

	// Checks if the the msg sender is the same as the current owner
	if msg.Estimator != k.GetEstimatorOwner(ctx, msg.Itemid+"-"+msg.Estimator) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.SetEstimator(ctx, estimator)

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgDeleteEstimator(ctx sdk.Context, k keeper.Keeper, msg *types.MsgDeleteEstimator) (*sdk.Result, error) {
	if !k.HasEstimator(ctx, msg.Itemid+"-"+msg.Estimator) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %s doesn't exist", msg.Itemid+"-"+msg.Estimator))
	}
	if msg.Estimator != k.GetEstimatorOwner(ctx, msg.Itemid+"-"+msg.Estimator) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}
	item := k.GetItem(ctx, msg.Itemid)
//Only delete estimator when it is not lowest /highest and not transferable

if item.Status != "" {
	return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "item has a status")
}

if msg.Estimator == item.Highestestimator || msg.Estimator == item.Lowestestimator {
	return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "estimator is lowest or highest")
}

	
		
		for i, v := range item.Estimatorlist {
			if v == msg.Estimator {
				item.Estimatorlist = append(item.Estimatorlist[:i], item.Estimatorlist[i+1:]...)
				break
			}
		}
		
		k.DeleteEstimator(ctx, msg.Itemid+"-"+msg.Estimator)



	

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgCreateFlag(ctx sdk.Context, k keeper.Keeper, msg *types.MsgCreateFlag) (*sdk.Result, error) {
	
	// Checks that the element exists
	if k.HasEstimator(ctx, msg.Itemid+"-"+msg.Estimator) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, fmt.Sprintf("key %s does exist", msg.Itemid+"-"+msg.Estimator))
	}

	item := k.GetItem(ctx, msg.Itemid)

	if item.Transferable == true {
		return nil, sdkerrors.Wrap(nil, "item is already estimated")
	}


	//remove item when it is flagged enough
	if int64(len(item.Estimatorlist)/2) < item.Flags + 1 {

		for _, element := range item.Estimatorlist {
			//apply this to each element
			key := msg.Itemid + "-" + element
			k.DeleteEstimator(ctx, key)

		}
		item.Bestestimator = ""
		item.Lowestestimator = ""
		item.Highestestimator = ""
		item.Estimatorlist = nil
		item.Status = "Removed (Item flagged)"
		k.DeleteItem(ctx, msg.Itemid)
		return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
	}
	item.Flags = item.Flags + 1
	k.SetItem(ctx, item)
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}
