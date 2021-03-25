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
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("item %s does exist", msg.Itemid))
	}

	///for production: check if estimator is item owner
	//
	//

	//checks whether the estimator already estimated the item
	//[for testing this is may be disabled]
	if k.HasEstimator(ctx, msg.Itemid+"-"+msg.Estimator) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %s does exist", msg.Itemid+"-"+msg.Estimator))
	}

	if msg.Deposit != item.Depositamount {
		return nil, sdkerrors.Wrap(nil, "deposit invalid")
	}

	//checks whether estimationcount has been reached
	var estimatorlistlen = strconv.Itoa(len(item.Estimatorlist))
	var estimatorlistlenhash = sha256.Sum256([]byte(estimatorlistlen + item.Seller))
	var estimatorlisthashstring = hex.EncodeToString(estimatorlistlenhash[:])
	if estimatorlisthashstring == item.Estimationcounthash {
		return nil, sdkerrors.Wrap(nil, "final estimation has already been made, estimation can not be added")
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

	k.DeleteEstimator(ctx, msg.Itemid+"-"+msg.Estimator)

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgCreateFlag(ctx sdk.Context, k keeper.Keeper, msg *types.MsgCreateFlag) (*sdk.Result, error) {
	var estimator = types.Estimator{
		Estimator: msg.Estimator,

		Itemid: msg.Itemid,
	}

	// Checks that the element exists
	if k.HasEstimator(ctx, msg.Itemid+"-"+msg.Estimator) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, fmt.Sprintf("key %s does exist", msg.Itemid+"-"+msg.Estimator))
	}

	item := k.GetItem(ctx, msg.Itemid)

	if item.Transferable == true {
		return nil, sdkerrors.Wrap(nil, "item is already estimated")
	}

	if msg.Flag == true {
		item.Flags = item.Flags + 1
		//test
		k.SetEstimator(ctx, estimator)
		k.SetItem(ctx, item)
	}

	k.SetEstimator(ctx, estimator)

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}
