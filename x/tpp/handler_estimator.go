package tpp

import (

	//"encoding/hex"
	//"github.com/tendermint/tendermint/crypto"

	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/danieljdd/tpp/x/tpp/keeper"
	"github.com/danieljdd/tpp/x/tpp/types"
)

func handleMsgCreateEstimation(ctx sdk.Context, k keeper.Keeper, msg *types.MsgCreateEstimation) (*sdk.Result, error) {
	fmt.Printf("handling msg: %X\n", msg.Itemid)
	item := k.GetItem(ctx, msg.Itemid)
	if item.Bestestimator != "" {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "item already has an estimation")
	}
	if item.Estimationprice > 0 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "item already has an estimation")
	}

	///for production: check if estimator is item owner
	if msg.Estimator == item.Creator {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "estimator cannot be item creator")
	}

	//checks whether the estimator already estimated the item

	if k.HasEstimator(ctx, append(types.Uint64ToByte(msg.Itemid), []byte(msg.Estimator)...)) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "already estimated this item")
	}

	if msg.Deposit != item.Depositamount {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "deposit invalid")
	}
	/*
		//checks whether estimationcount will be reached
		var estimatorlistlen = strconv.Itoa(len(item.Estimatorlist) + 1)
		var estimatorlistlenhash = sha256.Sum256([]byte(estimatorlistlen + item.Seller))
		var estimatorlisthashstring = hex.EncodeToString(estimatorlistlenhash[:])
		if estimatorlisthashstring == item.Estimationcounthash {
			item.Bestestimator = "Awaiting"

			ctx.EventManager().EmitEvent(
				sdk.NewEvent(types.EventTypeItemReady, sdk.NewAttribute(types.AttributeKeyItemID, strconv.FormatUint(msg.Itemid, 10))))
		}

		var estimatorestimationhash = sha256.Sum256([]byte(strconv.FormatInt(msg.Estimation, 10) + msg.Estimator))
		var estimatorestimationhashstring = hex.EncodeToString(estimatorestimationhash[:])

		//append estimatorhash to list
		item.Estimatorestimationhashlist = append(item.Estimatorestimationhashlist, estimatorestimationhashstring)

		item.Estimatorlist = append(item.Estimatorlist, msg.Estimator)
	*/
	item.Estimatorlist = append(item.Estimatorlist, msg.Estimator)
	fmt.Printf("setting item msg: %X\n", msg.Itemid)
	k.SetItem(ctx, item)
	fmt.Printf("go to keeper item msg: %X\n", msg.Itemid)
	k.CreateEstimation(ctx, *msg)

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgUpdateLike(ctx sdk.Context, k keeper.Keeper, msg *types.MsgUpdateLike) (*sdk.Result, error) {
	var estimator = types.Estimator{
		Estimator: msg.Estimator,

		Itemid: msg.Itemid,

		Interested: msg.Interested,
	}

	// Checks that the element exists
	if !k.HasEstimator(ctx, append(types.Uint64ToByte(msg.Itemid), []byte(msg.Estimator)...)) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("like of %s doesn't exist", msg.Estimator))
	}

	// Checks if the the msg sender is the same as the current owner
	if msg.Estimator != k.GetEstimatorOwner(ctx, append(types.Uint64ToByte(msg.Itemid), []byte(msg.Estimator)...)) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.SetEstimator(ctx, estimator)

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgDeleteEstimation(ctx sdk.Context, k keeper.Keeper, msg *types.MsgDeleteEstimation) (*sdk.Result, error) {
	if !k.HasEstimator(ctx, append(types.Uint64ToByte(msg.Itemid), []byte(msg.Estimator)...)) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("estimation of %s doesn't exist", msg.Estimator))
	}
	if msg.Estimator != k.GetEstimatorOwner(ctx, append(types.Uint64ToByte(msg.Itemid), []byte(msg.Estimator)...)) {
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

	err := k.DeleteEncryptedEstimation(ctx, item, *msg)
	if err != nil {
		fmt.Printf("err executing")
		fmt.Printf("executing contract: %X\n", err)
		return nil, sdkerrors.Wrap(err, "error deleting estimation") ///panic(err)
	}

	k.DeleteEstimation(ctx, append(types.Uint64ToByte(msg.Itemid), []byte(msg.Estimator)...))

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgFlagItem(ctx sdk.Context, k keeper.Keeper, msg *types.MsgFlagItem) (*sdk.Result, error) {

	// Checks that the element exists
	if k.HasEstimator(ctx, append(types.Uint64ToByte(msg.Itemid), []byte(msg.Estimator)...)) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, fmt.Sprintf("estimation of %s does exist", msg.Estimator))
	}

	item := k.GetItem(ctx, msg.Itemid)

	if item.Transferable {
		return nil, sdkerrors.Wrap(nil, "item is already estimated")
	}

	err := k.Flag(ctx, item, *msg)
	if err != nil {
		//	fmt.Printf("err executing")
		//	fmt.Printf("executing contract: %X\n", err)
		return nil, sdkerrors.Wrap(err, "error flagging item") ///panic(err)
	}
	/*//remove item when it is flagged enough
	if int64(len(item.Estimatorlist)/2) < item.Flags+1 {

		for _, element := range item.Estimatorlist {
			//apply this to each element
			key := append(types.Uint64ToByte(msg.Itemid), []byte(element)...)
			k.DeleteEstimation(ctx, key)

		}
		item.Bestestimator = ""
		item.Lowestestimator = ""
		item.Highestestimator = ""
		item.Estimatorlist = nil
		item.Status = "Removed (Item flagged)"
		k.DeleteItem(ctx, msg.Itemid)
		k.RemoveFromItemSeller(ctx, msg.Itemid, item.Seller)
		return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
	}*/
	//item.Flags = item.Flags + 1
	//k.SetItem(ctx, item)
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}
