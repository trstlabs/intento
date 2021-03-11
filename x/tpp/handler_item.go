package tpp

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"sort"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/danieljdd/tpp/x/tpp/keeper"
	"github.com/danieljdd/tpp/x/tpp/types"
	//"github.com/tendermint/tendermint/crypto"
)

func handleMsgCreateItem(ctx sdk.Context, k keeper.Keeper, msg *types.MsgCreateItem) (*sdk.Result, error) {
	k.CreateItem(ctx, *msg)

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgUpdateItem(ctx sdk.Context, k keeper.Keeper, msg *types.MsgUpdateItem) (*sdk.Result, error) {
	var item = types.Item{
		Creator: msg.Creator,
		Id:      msg.Id,

		Shippingcost: msg.Shippingcost,
		Localpickup:  msg.Localpickup,

		Shippingregion: msg.Shippingregion,
	}

	// Checks that the element exists
	if !k.HasItem(ctx, msg.Id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %s doesn't exist", msg.Id))
	}

	// Checks if the the msg sender is the same as the current owner
	if msg.Creator != k.GetItemOwner(ctx, msg.Id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.SetItem(ctx, item)

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgDeleteItem(ctx sdk.Context, k keeper.Keeper, msg *types.MsgDeleteItem) (*sdk.Result, error) {
	if !k.HasItem(ctx, msg.Id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %s doesn't exist", msg.Id))
	}
	if msg.Creator != k.GetItemOwner(ctx, msg.Id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	item := k.GetItem(ctx, msg.Id)

	if item.Status != "" && item.Buyer == "" {

		if len(item.Estimatorlist) > 0 {
			for _, element := range item.Estimatorlist {
				//apply this to each element
				key := msg.Id + "-" + element

				k.DeleteEstimator(ctx, key)
			}
		}

		item.Title = "Deleted"
		item.Description = ""
		item.Shippingcost = 0
		item.Localpickup = false
		item.Estimationcounthash = ""
		item.Bestestimator = ""
		item.Lowestestimator = ""
		item.Highestestimator = ""
		item.Estimationprice = 0
		item.Estimatorlist = nil
		item.Estimatorestimationhashlist = nil
		item.Transferable = false
		item.Buyer = ""
		item.Tracking = false
		item.Comments = nil
		item.Tags = nil
		item.Condition = 0
		item.Shippingregion = nil

		k.SetItem(ctx, item)

	} else {

		//if estimation is made pay back all the estimators/or buyer (like handlerItemTransfer)
		if len(item.Estimatorlist) > 0 {
			for _, element := range item.Estimatorlist {
				//apply this to each element
				key := msg.Id + "-" + element

				k.DeleteEstimator(ctx, key)

			}
		}
		k.DeleteItem(ctx, msg.Id)
	}

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgRevealEstimation(ctx sdk.Context, k keeper.Keeper, msg *types.MsgRevealEstimation) (*sdk.Result, error) {

	item := k.GetItem(ctx, msg.Itemid)

	if item.Estimatorlist == nil {
		return nil, sdkerrors.Wrap(nil, "Item does not have estimators ")
	}

	//[optional idea]can maybe replace estimatorlist hash with estimatorestimationlist hash to reduce computation
	var estimatorlistlen = len(item.Estimatorlist)
	var estimatorlistlenstring = strconv.Itoa(estimatorlistlen)
	var estimatorlistlenhash = sha256.Sum256([]byte(estimatorlistlenstring + item.Creator))
	var estimatorlisthashstring = hex.EncodeToString(estimatorlistlenhash[:])

	if estimatorlisthashstring != item.Estimationcounthash {
		return nil, sdkerrors.Wrap(nil, "Estimation count has not been reached, final estimation can not be made")
	}

	if item.Bestestimator != "" {
		return nil, sdkerrors.Wrap(nil, "item already has an estimation price")
	}

	//remove item when it is flagged a lot w/ at least three estimators
	if (int64(estimatorlistlen)/2) < item.Flags && estimatorlistlen > 3 {

		for _, element := range item.Estimatorlist {
			//apply this to each element

			key := msg.Itemid + "-" + element
			k.DeleteEstimator(ctx, key)

		}
		item.Bestestimator = ""
		item.Lowestestimator = ""
		item.Highestestimator = ""
		item.Estimatorlist = nil
		item.Status = "Removed (Reason: reported too often)"
		k.DeleteItem(ctx, msg.Itemid)

	}

	var Commentlist []string
	var EstimationList []int64

	for _, element := range item.Estimatorlist {
		key := msg.Itemid + "-" + element
		estimator := k.GetEstimator(ctx, key)

		//getting all of the comments into a list
		//var comment = estimator.Comment
		Commentlist = append(Commentlist, estimator.Comment)

		//append estimation to estimationlist
		//estimation := estimator.Estimation.Int64()
		EstimationList = append(EstimationList, estimator.Estimation)

		//create median
		//medianIndex := int64(math.Floor(float64(len(EstimationList))-1.0) / 2)
		sortedList := make([]int64, len(EstimationList))
		copy(sortedList, EstimationList)
		sort.Slice(sortedList, func(i, j int) bool { return sortedList[i] < sortedList[j] })

		//median := sortedList[medianIndex]

		//update the highest and lowest estimator
		if estimator.Estimation == sortedList[0] {
			item.Lowestestimator = estimator.Estimator
		}
		if estimator.Estimation == sortedList[(len(sortedList)-1)] {
			item.Highestestimator = estimator.Estimator
		}

	}

	//create median
	medianIndex := int64(math.Floor(float64(len(EstimationList))-1.0) / 2)
	sortedList := make([]int64, len(EstimationList))
	copy(sortedList, EstimationList)
	sort.Slice(sortedList, func(i, j int) bool { return sortedList[i] < sortedList[j] })
	median := sortedList[medianIndex]
	//var EstimationPrice = sdk.NewInt(median)

	for _, element := range item.Estimatorlist {
		//apply this to each element
		key := msg.Itemid + "-" + element
		estimator := k.GetEstimator(ctx, key)

		//finding out if the creator of the estimation belongs to the best estimated price
		var estimatorestimation = []byte(strconv.FormatInt(median, 10) + estimator.Estimator)
		var estimatorestimationhash = sha256.Sum256(estimatorestimation)
		var estimatorestimationhashstring = hex.EncodeToString(estimatorestimationhash[:])

		//assigns revealer of the best estimation to the item
		_, found := types.Find(item.Estimatorestimationhashlist, estimatorestimationhashstring)
		if found == true {
			item.Bestestimator = estimator.Estimator
			item.Estimationprice = median
			item.Comments = Commentlist
			k.SetItem(ctx, item)
		}

	}

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgItemTransferable(ctx sdk.Context, k keeper.Keeper, msg *types.MsgItemTransferable) (*sdk.Result, error) {
	//check if item exists
	item := k.GetItem(ctx, msg.Itemid)

	//check if message creator is item creator
	if msg.Creator != k.GetItemOwner(ctx, msg.Itemid) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	//check if item has a best estimator (and therefore a complete estimation)
	if item.Bestestimator == "" {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "item does not have estimation yet, cannot make item transferable")
	}

	if msg.Transferable == false {

		//returns each element
		for _, element := range item.Estimatorlist {
			//apply this to each element
			key := msg.Itemid + "-" + element
			estimator := k.GetEstimator(ctx, key)

			if estimator.Estimator == item.Lowestestimator {
				if item.Depositamount < (item.Estimationprice / 4) {
					k.BurnCoins(ctx, estimator.Deposit)
				} else {
					k.DeleteEstimator(ctx, key)
				}

			} else {
				k.DeleteEstimator(ctx, key)
			}

		}
		//item.TransferBool = msg.TransferBool
		//k.SetItem(ctx, item)
		item.Bestestimator = ""
		item.Lowestestimator = ""
		item.Highestestimator = ""
		item.Estimatorlist = nil

		//item has to be deleted because otherwise this function can be run again and again causing an uneven distribution and a drain of the moduleacc
		k.DeleteItem(ctx, msg.Itemid)

	} else {
		item.Transferable = msg.Transferable
		k.SetItem(ctx, item)
	}

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgItemShipping(ctx sdk.Context, k keeper.Keeper, msg *types.MsgItemShipping) (*sdk.Result, error) {
	//check if message creator is item creator
	if msg.Creator != k.GetItemOwner(ctx, msg.Itemid) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}
	//get buyer info
	buyer := k.GetBuyer(ctx, msg.Itemid)

	//get item info
	item := k.GetItem(ctx, msg.Itemid)

	//check if item.transferable = true and therefore the creator has accepted the buyer
	////[to do] in case this is false the prepayment will  be returned, item.buyer will be gone. this shall be in another function
	if item.Transferable == false {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "creator of item does not accept a transfer")
	}

	//check if item has a buyer already (so that we know that prepayment is done)
	if item.Buyer == "" {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "item does not have a buyer yet")
	}
	//bonus check if item already has been transferred
	if item.Status != "" {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "item already has had a transfer or transfer has been denied ")
	}

	bigintestimationprice := sdk.NewInt(item.Estimationprice)
	bigintshipping := sdk.NewInt(item.Shippingcost)
	if msg.Tracking == true {

		//rounded down percentage for minting. Percentage may be changed through governance proposals
		percentageMint := sdk.NewDecWithPrec(10, 2)
		percentageReward := sdk.NewDecWithPrec(5, 2)
		//paymentCreator := percentageCreator.MulInt(bigintestimationprice)
		//roundedAmountCreaterPayout := paymentCreator.TruncateInt()
		CreaterPayoutAndShipping := bigintestimationprice.Add(bigintshipping)
		//rounded up percentage as a reward for the estimator
		//percentageReward := sdk.NewDecWithPrec(3, 2)
		toMint := percentageMint.MulInt(bigintestimationprice)
		toMintAmount := toMint.TruncateInt()
		paymentReward := percentageReward.MulInt(bigintestimationprice)
		roundedAmountReward := paymentReward.TruncateInt()
		//roundedAmountRewardBestEstimator := paymentReward.TruncateInt()

		//make payment to creator and estimator
		paymentCreatorCoins := sdk.NewCoin("tpp", CreaterPayoutAndShipping)

		//minted coins (are rounded up)
		mintCoins := sdk.NewCoin("tpp", toMintAmount)
		paymentRewardCoins := sdk.NewCoin("tpp", roundedAmountReward)
		//paymentRewardCoinsEstimator := sdk.NewCoin("tpp", roundedAmountRewardBestEstimator)

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

		item.Status = "Shipped"
		k.SetItem(ctx, item)
		k.SetBuyer(ctx, buyer)
	}

	if msg.Tracking == false {
		repayment := bigintestimationprice.Add(bigintshipping)
		repaymentCoins := sdk.NewCoin("tpp", repayment)

		k.HandlePrepayment(ctx, item.Buyer, repaymentCoins)

		for _, element := range item.Estimatorlist {
			//apply this to each element
			key := msg.Itemid + "-" + element
			estimator := k.GetEstimator(ctx, key)

			if estimator.Estimator == item.Lowestestimator {
				if item.Depositamount < (item.Estimationprice / 4) {
					k.BurnCoins(ctx, estimator.Deposit)
				} else {
					k.DeleteEstimator(ctx, key)
				}

			} else {
				k.DeleteEstimator(ctx, key)
			}

		}

		item.Status = "Shipping declined; buyer refunded"
		k.SetItem(ctx, item)
		k.SetBuyer(ctx, buyer)
	}

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}
