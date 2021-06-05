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

func handleMsgDeleteItem(ctx sdk.Context, k keeper.Keeper, msg *types.MsgDeleteItem) (*sdk.Result, error) {
	if !k.HasItem(ctx, msg.Id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("item id %s doesn't exist", strconv.FormatUint(msg.Id, 10)))
	}

	item := k.GetItem(ctx, msg.Id)

	if item.Buyer != "" && item.Status != "" {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Item has a buyer")
	}

	if item.Status != "" && item.Buyer == "" {

		if len(item.Estimatorlist) > 0 {
			for _, element := range item.Estimatorlist {
				//apply this to each element

				key := append(types.Uint64ToByte(msg.Id), []byte(element)...)
				k.DeleteEstimation(ctx, key)
			}
		}
		// title,status and rating are kept to enhance trust
		item.Description = ""
		item.Shippingcost = 0
		item.Localpickup = ""
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
		item.Note = ""

		k.SetItem(ctx, item)

	} else {

		//if estimation is made pay back all the estimators/or buyer (like handleMsgItemTransfer)
		if len(item.Estimatorlist) > 0 {
			for _, element := range item.Estimatorlist {
				//apply this to each element
				key := append(types.Uint64ToByte(msg.Id), []byte(element)...)

				k.DeleteEstimation(ctx, key)

			}
		}
		k.RemoveFromInactiveItemQueue(ctx, msg.Id, item.Endtime)
		k.DeleteItem(ctx, msg.Id)
		k.RemoveFromItemSeller(ctx, msg.Id, msg.Seller)
	}

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgRevealEstimation(ctx sdk.Context, k keeper.Keeper, msg *types.MsgRevealEstimation) (*sdk.Result, error) {

	item := k.GetItem(ctx, msg.Itemid)

	if item.Bestestimator != "Awaiting" {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "item cannot be revealed")
	}

	var CommentList []string
	var EstimationList []int64

	for _, element := range item.Estimatorlist {
		key := append(types.Uint64ToByte(msg.Itemid), []byte(element)...)
		estimator := k.GetEstimator(ctx, key)

		//getting all of the comments into a list
		//var comment = estimator.Comment
		CommentList = append(CommentList, estimator.Comment)

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
	//var Estimationprice = sdk.NewInt(median)

	///delete item when deposit is higher than 25% of the item price (Can be altered through governance)
	if item.Depositamount > (median / 4) {
		//returns each element
		for _, element := range item.Estimatorlist {
			//apply this to each element
			key := append(types.Uint64ToByte(msg.Itemid), []byte(element)...)
			k.DeleteEstimation(ctx, key)
		}
		/*item.Bestestimator = ""
		item.Lowestestimator = ""
		item.Highestestimator = ""
		item.Estimatorlist = nil
		item.Status = "Estimators refunded"*/
		k.DeleteItem(ctx, msg.Itemid)
		k.RemoveFromItemSeller(ctx, item.Id, item.Seller)
		return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
	}

	for _, element := range item.Estimatorlist {
		//apply this to each element
		key := append(types.Uint64ToByte(msg.Itemid), []byte(element)...)
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
			item.Comments = CommentList
			k.SetItem(ctx, item)
			break
		}

	}

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgItemTransferable(ctx sdk.Context, k keeper.Keeper, msg *types.MsgItemTransferable) (*sdk.Result, error) {
	//check if item exists
	item := k.GetItem(ctx, msg.Itemid)

	//check if message creator is item seller
	if msg.Seller != k.GetItemOwner(ctx, msg.Itemid) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	//check if item has a best estimator (and therefore a complete estimation)
	if item.Bestestimator == "" {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "no estimation price yet, cannot make item transferable")
	}

	if msg.Transferable == false {

		//returns each element
		for _, element := range item.Estimatorlist {
			//apply this to each element
			key := append(types.Uint64ToByte(msg.Itemid), []byte(element)...)
			estimator := k.GetEstimator(ctx, key)

			if estimator.Estimator == item.Lowestestimator {

				k.BurnCoins(ctx, estimator.Deposit)
				k.DeleteEstimationWithoutDeposit(ctx, key)

			} else {
				k.DeleteEstimation(ctx, key)
			}

		}
		//item.TransferBool = msg.TransferBool
		//k.SetItem(ctx, item)
		//item.Bestestimator = ""
		//item.Lowestestimator = ""
		//item.Highestestimator = ""
		//item.Estimatorlist = nil

		//item has to be deleted because otherwise this function can be run again
		k.DeleteItem(ctx, msg.Itemid)
		k.RemoveFromItemSeller(ctx, item.Id, item.Seller)
		k.RemoveFromInactiveItemQueue(ctx, msg.Itemid, item.Endtime)
	} else {
		item.Transferable = msg.Transferable
		k.SetItem(ctx, item)
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.EventTypeItemTransferable, sdk.NewAttribute(types.AttributeKeyItemID, strconv.FormatUint(msg.Itemid, 10))))

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgItemShipping(ctx sdk.Context, k keeper.Keeper, msg *types.MsgItemShipping) (*sdk.Result, error) {
	//check if message seller is item seller
	if msg.Seller != k.GetItemOwner(ctx, msg.Itemid) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	//get item info
	item := k.GetItem(ctx, msg.Itemid)

	//check if item.transferable = true and therefore the seller has accepted the buyer
	if item.Transferable == false {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "item is not transferable")
	}

	//check if item has a buyer already (so that we know that prepayment is done)
	if item.Buyer == "" {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "item does not have a buyer yet")
	}
	//bonus check if item already has been transferred
	if item.Status != "" {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "item already has had a transfer or transfer has been denied ")
	}
	if item.Shippingcost < 1 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "unauthrorized, no shippingcost")
	}

	bigintestimationprice := sdk.NewInt(item.Estimationprice - item.Depositamount)
	if item.Discount > 0 {
		bigintestimationprice = sdk.NewInt(item.Estimationprice - item.Discount)
	}

	bigintshipping := sdk.NewInt(item.Shippingcost)
	if msg.Tracking == true {
		if item.Creator == item.Seller {
			//rounded down percentage for minting. Percentage may be changed through governance proposals
			/*percentageMint := sdk.NewDecWithPrec(10, 2)
			percentageReward := sdk.NewDecWithPrec(5, 2)
			//paymentSeller := percentageSeller.MulInt(bigintestimationprice)
			//roundedAmountCreaterPayout := paymentSeller.TruncateInt()


			//rounded up percentage as a reward for the estimator
			//percentageReward := sdk.NewDecWithPrec(3, 2)
			toMint := percentageMint.MulInt(bigintestimationprice)
			toMintAmount := toMint.TruncateInt()
			paymentReward := percentageReward.MulInt(bigintestimationprice)
			roundedAmountReward := paymentReward.TruncateInt()*/
			//roundedAmountRewardBestEstimator := paymentReward.TruncateInt()

			//minted coins (are rounded up)
			rewardCoins := sdk.NewCoin("tpp", sdk.NewInt(item.Depositamount))
			//paymentRewardCoins := sdk.NewCoin("tpp", roundedAmountReward)
			//paymentRewardCoinsEstimator := sdk.NewCoin("tpp", roundedAmountRewardBestEstimator)

			k.MintReward(ctx, rewardCoins)

			//for their participation in the protocol, the best estimator, a random and the stakers get rewarded.
			k.HandleEstimatorReward(ctx, item.Bestestimator, rewardCoins)

			k.HandleStakingReward(ctx, rewardCoins)

			//refund the deposits back to all of the item estimators
			for _, element := range item.Estimatorlist {
				key := append(types.Uint64ToByte(msg.Itemid), []byte(element)...)

				k.DeleteEstimation(ctx, key)
			}

			item.Bestestimator = ""
			item.Lowestestimator = ""
			item.Highestestimator = ""
			item.Estimatorlist = nil
		}
		//make payment to seller
		CreaterPayoutAndShipping := bigintestimationprice.Add(bigintshipping)
		paymentSellerCoins := sdk.NewCoin("tpp", CreaterPayoutAndShipping)

		k.HandlePrepayment(ctx, item.Seller, paymentSellerCoins)

		item.Status = "Shipped"
		k.SetItem(ctx, item)
		//k.SetBuyer(ctx, buyer)
	} else {
		repayment := bigintestimationprice.Add(bigintshipping)
		repaymentCoins := sdk.NewCoin("tpp", repayment)

		k.HandlePrepayment(ctx, item.Buyer, repaymentCoins)

		for _, element := range item.Estimatorlist {
			//apply this to each element
			key := append(types.Uint64ToByte(msg.Itemid), []byte(element)...)
			//estimator := k.GetEstimator(ctx, key)

			/*if estimator.Estimator == item.Lowestestimator {

				k.BurnCoins(ctx, estimator.Deposit)
				k.DeleteEstimationWithoutDeposit(ctx, key)

			} else {*/
			k.DeleteEstimation(ctx, key)
			//}

		}

		item.Status = "Shipping declined; buyer refunded"
		k.SetItem(ctx, item)
		//k.SetBuyer(ctx, buyer)
	}

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgItemResell(ctx sdk.Context, k keeper.Keeper, msg *types.MsgItemResell) (*sdk.Result, error) {

	// Checks that the element exists
	if !k.HasItem(ctx, msg.Itemid) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("item id %s doesn't exist", strconv.FormatUint(msg.Itemid, 10)))
	}

	item := k.GetItem(ctx, msg.Itemid)

	// Checks if the the msg sender is the same as the current buyer
	if msg.Seller != item.Buyer {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	if item.Status != "Transferred" && item.Status != "Shipped" {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Item not available to resell")
	}

	if msg.Discount > item.Estimationprice {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Discount invalid")
	}
	k.Withdrawal(ctx, append(types.Uint64ToByte(msg.Itemid), []byte(msg.Seller)...))

	item.Seller = msg.Seller
	item.Shippingcost = msg.Shippingcost
	item.Localpickup = msg.Localpickup
	item.Shippingregion = msg.Shippingregion
	item.Discount = msg.Discount
	item.Note = msg.Note
	item.Rating = 0
	item.Buyer = ""
	item.Status = ""
	item.Endtime = ctx.BlockTime().Add(types.DefaultParams().MaxActivePeriod)

	k.SetItem(ctx, item)
	k.InsertInactiveItemQueue(ctx, msg.Itemid, item.Endtime)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.EventTypeItemResellable, sdk.NewAttribute(types.AttributeKeyItemID, strconv.FormatUint(msg.Itemid, 10))))
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}
