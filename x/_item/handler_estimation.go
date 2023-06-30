package trst

import (

	//"encoding/hex"
	//"github.com/cometbft/cometbft/crypto"

	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/trstlabs/trst/x/item/keeper"
	"github.com/trstlabs/trst/x/item/types"
)

func handleMsgCreateEstimation(ctx sdk.Context, k keeper.Keeper, msg *types.MsgCreateEstimation) (*sdk.Result, error) {
	//fmt.Printf("handling msg: %X\n", msg.Itemid)
	item := k.GetItem(ctx, msg.Itemid)
	if item.Id == 0 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "item not found")
	}
	if item.Id != msg.Itemid {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "item not found")
	}
	if item.Estimation.BestEstimator != "" {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "item already has an estimation")
	}
	if item.Estimation.EstimationPrice > 0 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "item already has an estimation")
	}

	///for production: check if estimator is item buyer
	if msg.Estimator == item.Creator {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "estimator cannot be item creator")
	}

	if msg.Deposit != item.Estimation.DepositAmount {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "deposit invalid")
	}

	item.Estimation.EstimatorList = append(item.Estimation.EstimatorList, msg.Estimator)
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

	var estimationInfo = types.EstimationInfo{
		Itemid:     msg.Itemid,
		Interested: msg.Interested,
	}

	err := k.UpdateEstimationInfo(ctx, estimationInfo, msg.Estimator)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "estimation info not found")
	}

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgDeleteEstimation(ctx sdk.Context, k keeper.Keeper, msg *types.MsgDeleteEstimation) (*sdk.Result, error) {

	item := k.GetItem(ctx, msg.Itemid)
	//Only delete estimator when it is not lowest /highest and not transferable

	if item.Status != "" {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "item has a status")
	}

	for i, v := range item.Estimation.EstimatorList {
		if v == msg.Estimator {
			item.Estimation.EstimatorList = append(item.Estimation.EstimatorList[:i], item.Estimation.EstimatorList[i+1:]...)
			break
		} else {
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "estimator not found in item info")
		}
	}

	err := k.DeleteEncryptedEstimation(ctx, item, *msg)
	if err != nil {
		//	fmt.Printf("err executing")
		//	fmt.Printf("executing contract: %X\n", err)
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, err.Error()) ///panic(err)
	}

	//k.DeleteEstimation(ctx, append(types.Uint64ToByte(msg.Itemid), []byte(msg.Estimator)...))

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgFlagItem(ctx sdk.Context, k keeper.Keeper, msg *types.MsgFlagItem) (*sdk.Result, error) {

	item := k.GetItem(ctx, msg.Itemid)

	if item.Estimation.EstimationPrice > 0 {
		return nil, sdkerrors.Wrap(nil, "estimation period has ended")
	}

	err := k.Flag(ctx, item, *msg)
	if err != nil {
		//	fmt.Printf("err executing")
		//	fmt.Printf("executing contract: %X\n", err)
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, err.Error()) ///panic(err)
	}
	/*//remove item when it is flagged enough
	if int64(len(item.Estimation.EstimatorList)/2) < item.Flags+1 {

		for _, element := range item.Estimation.EstimatorList {
			//apply this to each element
			key := append(types.Uint64ToByte(msg.Itemid), []byte(element)...)
			k.DeleteEstimation(ctx, key)

		}
		item.Estimation.BestEstimator = ""
		item.Lowestestimator = ""
		item.Highestestimator = ""
		item.Estimation.EstimatorList = nil
		item.Status = "Removed (Item flagged)"
		k.DeleteItem(ctx, msg.Itemid)
		k.RemoveFromSellerItems(ctx, msg.Itemid, item.Transfer.Seller)
		return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
	}*/
	//item.Flags = item.Flags + 1
	//k.SetItem(ctx, item)
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}
