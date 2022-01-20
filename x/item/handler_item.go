package trst

import (
	"fmt"

	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/trstlabs/trst/x/item/keeper"
	"github.com/trstlabs/trst/x/item/types"
	//"github.com/tendermint/tendermint/crypto"
)

func handleMsgCreateItem(ctx sdk.Context, k keeper.Keeper, msg *types.MsgCreateItem) (*sdk.Result, error) {

	err := k.CreateItem(ctx, *msg)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, err.Error())
	}
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

	//status but no buyer
	if item.Status != "" && item.Buyer == "" {

		/*if len(item.EstimatorList) > 0 {
			for _, element := range item.EstimatorList {
				//apply this to each element

				key := append(types.Uint64ToByte(msg.Id), []byte(element)...)
				k.DeleteEstimation(ctx, key)
			}
		}*/
		// title,status and rating are kept for record keeping
		item.Description = ""
		item.ShippingCost = 0
		item.LocalPickup = ""
		item.EstimationCount = 0
		item.BestEstimator = ""

		item.EstimationPrice = 0
		item.EstimatorList = nil
		item.EstimationList = nil
		item.Transferable = false
		item.Buyer = ""
		item.Tracking = false
		item.Comments = nil
		item.Tags = nil
		item.Condition = 0
		item.ShippingRegion = nil
		item.Note = ""
		item.Contract = ""
		item.Photos = nil

		k.SetItem(ctx, item)

	} else {

		k.RemoveFromListedItemQueue(ctx, msg.Id, item.EndTime)
		_ = k.DeleteItemContract(ctx, item.Contract)
		k.DeleteItem(ctx, msg.Id)
		k.RemoveFromItemSeller(ctx, msg.Id, msg.Seller)
	}

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgRevealEstimation(ctx sdk.Context, k keeper.Keeper, msg *types.MsgRevealEstimation) (*sdk.Result, error) {

	item := k.GetItem(ctx, msg.Itemid)

	if msg.Creator != item.Seller {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "not item seller account")
	}

	if item.EstimationCount != item.EstimationTotal {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "not enough estimators")
	}

	err := k.RevealEstimation(ctx, item, *msg)
	if err != nil {
		//	fmt.Printf("err executing")
		//	fmt.Printf("executing contract: %X\n", err)
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, err.Error()) ///panic(err)
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
	if item.IsToken {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Item is a token")
	}

	//check if item has a best estimator (and therefore a complete estimation)
	if item.BestEstimator == "" {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "no estimation price yet, cannot make item transferable")
	}

	err := k.Transferable(ctx, item, *msg)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, err.Error())
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
	if !item.Transferable {
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
	if item.ShippingCost < 1 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "unauthrorized, no shipping_cost")
	}

	bigIntEstimationPrice := sdk.NewInt(item.EstimationPrice - item.DepositAmount)
	if item.Discount > 0 {
		bigIntEstimationPrice = sdk.NewInt(item.EstimationPrice - item.Discount)
	}

	bigIntShipping := sdk.NewInt(item.ShippingCost)
	if msg.Tracking {
		if item.Creator == item.Seller {
			//rounded down percentage for minting. Percentage may be changed through governance proposals
			/*percentageMint := sdk.NewDecWithPrec(10, 2)
			percentageReward := sdk.NewDecWithPrec(5, 2)
			//paymentSeller := percentageSeller.MulInt(bigIntEstimationPrice)
			//roundedAmountCreaterPayout := paymentSeller.TruncateInt()


			//rounded up percentage as a reward for the estimator
			//percentageReward := sdk.NewDecWithPrec(3, 2)
			toMint := percentageMint.MulInt(bigIntEstimationPrice)
			toMintAmount := toMint.TruncateInt()
			paymentReward := percentageReward.MulInt(bigIntEstimationPrice)
			roundedAmountReward := paymentReward.TruncateInt()*/
			//roundedAmountRewardBestEstimator := paymentReward.TruncateInt()

			//minted coins (are rounded up)
			rewardCoins := sdk.NewCoin("utrst", sdk.NewInt(item.DepositAmount))
			//paymentRewardCoins := sdk.NewCoin("utrst", roundedAmountReward)
			//paymentRewardCoinsEstimator := sdk.NewCoin("utrst", roundedAmountRewardBestEstimator)

			k.MintReward(ctx, rewardCoins)

			//for their participation in the protocol, the best estimator, a random and the stakers get rewarded.
			//k.HandleEstimatorReward(ctx, item.BestEstimator, rewardCoins)

			k.HandleCommunityReward(ctx, rewardCoins)

			item.BestEstimator = ""
			item.EstimatorList = nil
		}
		//make payment to seller
		CreaterPayoutAndShipping := bigIntEstimationPrice.Add(bigIntShipping)
		paymentSellerCoins := sdk.NewCoin("utrst", CreaterPayoutAndShipping)

		k.HandlePrepayment(ctx, item.Seller, paymentSellerCoins)

		item.Status = "Shipped"
		k.SetItem(ctx, item)
		//k.SetBuyer(ctx, buyer)
	} else {
		repayment := bigIntEstimationPrice.Add(bigIntShipping)
		repaymentCoins := sdk.NewCoin("utrst", repayment)

		k.HandlePrepayment(ctx, item.Buyer, repaymentCoins)

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

	if msg.Discount > item.EstimationPrice {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Discount invalid")
	}
	k.Withdrawal(ctx, append(types.Uint64ToByte(msg.Itemid), []byte(msg.Seller)...))

	item.Seller = msg.Seller
	item.ShippingCost = msg.ShippingCost
	item.LocalPickup = msg.LocalPickup
	item.ShippingRegion = msg.ShippingRegion
	item.Discount = msg.Discount
	item.Note = msg.Note
	item.Rating = 0
	item.Buyer = ""
	item.Status = ""
	item.EndTime = ctx.BlockTime().Add(types.DefaultParams().MaxActivePeriod)

	k.SetItem(ctx, item)
	k.InsertListedItemQueue(ctx, msg.Itemid, item, item.EndTime)
	k.BindItemSeller(ctx, msg.Itemid, msg.Seller)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.EventTypeItemResellable, sdk.NewAttribute(types.AttributeKeyItemID, strconv.FormatUint(msg.Itemid, 10))))
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgTokenizeItem(ctx sdk.Context, k keeper.Keeper, msg *types.MsgTokenizeItem) (*sdk.Result, error) {

	// Checks that the element exists
	if !k.HasItem(ctx, msg.Id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("item id %s doesn't exist", strconv.FormatUint(msg.Id, 10)))
	}

	item := k.GetItem(ctx, msg.Id)

	// Checks if the the msg sender is the same as the current buyer or creator
	if msg.Sender != item.Buyer && msg.Sender != item.Creator {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}
	//if item was sold previously or if item is not made transferable
	if item.Status == "Transferred" || item.Status == "Shipped" {
		//Create new coin
		//Make visible in item that it it tokenized now... (bool)
		err := k.TokenizeItem(ctx, msg.Id, (msg.Sender))
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, err.Error())
		}
		item.IsToken = true

	} else if !item.Transferable && item.BestEstimator != "" {
		err := k.TokenizeItem(ctx, msg.Id, msg.Sender)

		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, err.Error())
		}
		item.IsToken = true

	} else {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Item not available to tokenize")
	}
	k.SetItem(ctx, item)
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgUnTokenizeItem(ctx sdk.Context, k keeper.Keeper, msg *types.MsgUnTokenizeItem) (*sdk.Result, error) {

	// Checks that the element exists
	if !k.HasItem(ctx, msg.Id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("item id %s doesn't exist", strconv.FormatUint(msg.Id, 10)))
	}

	item := k.GetItem(ctx, msg.Id)

	err := k.UnTokenizeItem(ctx, msg.Id, msg.Sender)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Item not available to untokenize")

	}
	item.Buyer = msg.Sender
	item.IsToken = false
	k.SetItem(ctx, item)
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}
