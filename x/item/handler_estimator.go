package trst

import (

	//"encoding/hex"
	//"github.com/tendermint/tendermint/crypto"

	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/trstlabs/trst/x/item/keeper"
	"github.com/trstlabs/trst/x/item/types"
)

func handleMsgCreateEstimation(ctx sdk.Context, k keeper.Keeper, msg *types.MsgCreateEstimation) (*sdk.Result, error) {
	//fmt.Printf("handling msg: %X\n", msg.Itemid)
	item := k.GetItem(ctx, msg.Itemid)
	if item.Id != msg.Itemid {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "item not found")
	}
	if item.BestEstimator != "" {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "item already has an estimation")
	}
	if item.EstimationPrice > 0 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "item already has an estimation")
	}

	///for production: check if estimator is item owner
	if msg.Estimator == item.Creator {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "estimator cannot be item creator")
	}

	if msg.Deposit != item.DepositAmount {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "deposit invalid")
	}

	item.EstimatorList = append(item.EstimatorList, msg.Estimator)
	//fmt.Printf("setting item msg: %X\n", msg.Itemid)
	k.SetItem(ctx, item)
	//fmt.Printf("go to keeper item msg: %X\n", msg.Itemid)
	err := k.CreateEstimation(ctx, *msg)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, err.Error())
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.EventTypeItemReady, sdk.NewAttribute(types.AttributeKeyItemID, strconv.FormatUint(msg.Itemid, 10))))

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgUpdateLike(ctx sdk.Context, k keeper.Keeper, msg *types.MsgUpdateLike) (*sdk.Result, error) {
	var estimator = types.Estimator{
		Estimator: msg.Estimator,

		Itemid: msg.Itemid,

		Interested: msg.Interested,
	}

	// Checks that the element exists
	if !k.HasEstimationInfo(ctx, append(types.Uint64ToByte(msg.Itemid), []byte(msg.Estimator)...)) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("like of %s doesn't exist", msg.Estimator))
	}

	// Checks if the the msg sender is the same as the current owner
	if msg.Estimator != k.GetEstimator(ctx, append(types.Uint64ToByte(msg.Itemid), []byte(msg.Estimator)...)) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.SetEstimationInfo(ctx, estimator)

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgDeleteEstimation(ctx sdk.Context, k keeper.Keeper, msg *types.MsgDeleteEstimation) (*sdk.Result, error) {

	if msg.Estimator != k.GetEstimator(ctx, append(types.Uint64ToByte(msg.Itemid), []byte(msg.Estimator)...)) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}
	item := k.GetItem(ctx, msg.Itemid)
	//Only delete estimator when it is not lowest /highest and not transferable

	if item.Status != "" {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "item has a status")
	}

	for i, v := range item.EstimatorList {
		if v == msg.Estimator {
			item.EstimatorList = append(item.EstimatorList[:i], item.EstimatorList[i+1:]...)
			break
		}
	}

	err := k.DeleteEncryptedEstimation(ctx, item, *msg)
	if err != nil {
		//	fmt.Printf("err executing")
		//	fmt.Printf("executing contract: %X\n", err)
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, err.Error()) ///panic(err)
	}

	k.DeleteEstimation(ctx, append(types.Uint64ToByte(msg.Itemid), []byte(msg.Estimator)...))

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgFlagItem(ctx sdk.Context, k keeper.Keeper, msg *types.MsgFlagItem) (*sdk.Result, error) {

	item := k.GetItem(ctx, msg.Itemid)

	if item.Transferable {
		return nil, sdkerrors.Wrap(nil, "item is already estimated")
	}

	err := k.Flag(ctx, item, *msg)
	if err != nil {
		//	fmt.Printf("err executing")
		//	fmt.Printf("executing contract: %X\n", err)
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, err.Error()) ///panic(err)
	}
	/*//remove item when it is flagged enough
	if int64(len(item.EstimatorList)/2) < item.Flags+1 {

		for _, element := range item.EstimatorList {
			//apply this to each element
			key := append(types.Uint64ToByte(msg.Itemid), []byte(element)...)
			k.DeleteEstimation(ctx, key)

		}
		item.BestEstimator = ""
		item.Lowestestimator = ""
		item.Highestestimator = ""
		item.EstimatorList = nil
		item.Status = "Removed (Item flagged)"
		k.DeleteItem(ctx, msg.Itemid)
		k.RemoveFromItemSeller(ctx, msg.Itemid, item.Seller)
		return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
	}*/
	//item.Flags = item.Flags + 1
	//k.SetItem(ctx, item)
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}
