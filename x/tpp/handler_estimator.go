package tpp

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"github.com/tendermint/tendermint/crypto"

	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/danieljdd/tpp/x/tpp/keeper"
	"github.com/danieljdd/tpp/x/tpp/types"
)

func handleMsgCreateEstimator(ctx sdk.Context, k keeper.Keeper, msg *types.MsgCreateEstimator) (*sdk.Result, error) {

	item := k.GetItem(ctx, msg.Itemid)

	///for production: check if estimator is item owner
	//
	//

	//checks whether the estimator already estimated the item
	//[for testing this is may be disabled]
	if k.HasEstimator(ctx, msg.Itemid + "-" + msg.Estimator) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %s does exist", msg.Itemid + "-" + msg.Estimator))
	}

	//checks whether estimationcount has been reached
	var estimatorlistlen = strconv.Itoa(len(item.Estimatorlist))
	var estimatorlistlenhash = sha256.Sum256([]byte(estimatorlistlen))
	var estimatorlisthashstring = hex.EncodeToString(estimatorlistlenhash[:])
	if estimatorlisthashstring == item.Estimationcounthash {
		return nil, sdkerrors.Wrap(nil, "final estimation has already been made, estimation can not be added")
	}
	estimatoraddress, err := sdk.AccAddressFromBech32(msg.Estimator)
	if err != nil {
		sdkerrors.Wrap(err, "not an address")
	}

	moduleAcct := sdk.AccAddress(crypto.AddressHash([]byte(types.ModuleName)))
	sdkError := k.BankKeeper.SendCoinsFromAccountToModule(ctx, estimatoraddress, moduleAcct.String(), sdk.NewCoins(msg.Deposit))
	if sdkError != nil {
		return nil, sdkError
	}



	

	k.CreateEstimator(ctx,  *msg)

	//append estimatorhash to list
	item.Estimatorestimationhashlist = append(item.Estimatorestimationhashlist, msg.Estimatorestimationhash)

	item.Estimatorlist = append(item.Estimatorlist, msg.Estimator)
	
	k.SetItem(ctx, item)


	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgUpdateEstimator(ctx sdk.Context, k keeper.Keeper, msg *types.MsgUpdateEstimator) (*sdk.Result, error) {
	var estimator = types.Estimator{
		Estimator:                 msg.Estimator,
		
	
		Itemid:                  msg.Itemid,
	
		Interested:              msg.Interested,

	}

	// Checks that the element exists
	if !k.HasEstimator(ctx, msg.Itemid + "-" + msg.Estimator) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %s doesn't exist", msg.Itemid + "-" + msg.Estimator))
	}

	// Checks if the the msg sender is the same as the current owner
	if msg.Estimator != k.GetEstimatorOwner(ctx, msg.Itemid + "-" + msg.Estimator) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.SetEstimator(ctx, estimator)

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgDeleteEstimator(ctx sdk.Context, k keeper.Keeper, msg *types.MsgDeleteEstimator) (*sdk.Result, error) {
	if !k.HasEstimator(ctx, msg.Itemid + "-" + msg.Estimator) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %s doesn't exist", msg.Itemid + "-" + msg.Estimator))
	}
	if msg.Estimator != k.GetEstimatorOwner(ctx, msg.Itemid + "-" + msg.Estimator) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.DeleteEstimator(ctx, msg.Itemid + "-" + msg.Estimator)

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}
