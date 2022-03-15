package keeper

import (

	//"github.com/tendermint/tendermint/crypto"
	//	"strconv"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/trstlabs/trst/x/item/types"
)

// Prepayment creates a buyer with a new id and update the count
func (k Keeper) Prepayment(ctx sdk.Context, msg types.MsgPrepayment) (err error) {
	//get item info

	item := k.GetItem(ctx, msg.Itemid)

	//check if item has a best estimator (and therefore a complete estimation)
	if item.Estimation.EstimationPrice == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "item does not have a price yet, cannot make prepayment")
	}

	//check if item is transferable
	if !item.Properties.Transferable {
		return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "item not transferable, cannot make prepayment")
	}

	//check if item has a buyer already
	if item.Transfer.Buyer != "" {
		return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "item has a buyer, cannot make prepayment")
	}

	//item buyer cannot be the item creatorc
	if msg.Buyer == item.Creator || msg.Buyer == item.Transfer.Seller {
		return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "buyer cannot be creator/seller")
	}

	estimationPrice := item.Estimation.EstimationPrice
	if item.Transfer.Discount > 0 {
		estimationPrice = item.Estimation.EstimationPrice - item.Transfer.Discount
	}

	if item.Transfer.ShippingCost > 0 && item.Transfer.Location == "" {
		toPayShipping := estimationPrice + item.Transfer.ShippingCost
		if toPayShipping != msg.Deposit {

			return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "deposit insufficient, cannot make prepayment")
		}

	}

	if item.Transfer.ShippingCost == 0 && item.Transfer.Location != "" {

		if estimationPrice != msg.Deposit {

			return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "deposit insufficient, cannot make prepayment")
		}

	}

	if item.Transfer.ShippingCost > 0 && item.Transfer.Location != "" {

		if estimationPrice == msg.Deposit {

			item.Transfer.ShippingCost = 0

		} else {
			toPayShipping :=
				estimationPrice + item.Transfer.ShippingCost

			if toPayShipping == msg.Deposit {
				item.Transfer.Location = ""

			} else {

				return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "deposit insufficient, cannot make prepayment")
			}
		}

	}
	item.Transfer.Buyer = msg.Buyer
	k.RemoveFromListedItemQueue(ctx, msg.Itemid, item.ListingDuration.EndTime)
	k.SetItem(ctx, item)

	deposit := sdk.NewInt64Coin("utrst", msg.Deposit)

	buyeraddress, err := sdk.AccAddressFromBech32(msg.Buyer)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "address not found")
	}

	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, buyeraddress, types.ModuleName, sdk.NewCoins(deposit))
	if err != nil {
		//panic(err)
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "insufficent balance")
	}

	k.SetBuyer(ctx, msg.Itemid, msg.Buyer)

	return nil

}

// SetBuyer set a specific buyer in the store
func (k Keeper) SetBuyer(ctx sdk.Context, itemid uint64, buyer string) {
	store := ctx.KVStore(k.storeKey)

	store.Set(types.ItemBuyerKey(itemid, buyer), types.Uint64ToByte(itemid))
	//store.Set(append(append(types.KeyPrefix(types.BuyerKey), []byte(buyer)...), types.Uint64ToByte(itemid)...), types.Uint64ToByte(itemid))

}

// HasBuyer checks if the buyer exists
func (k Keeper) HasBuyer(ctx sdk.Context, key []byte) bool {
	//store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BuyerKey))
	store := ctx.KVStore(k.storeKey)
	return store.Has(append(types.KeyPrefix(types.BuyerKey), key...))
}

// DeleteBuyerKey deletes a buyer
func (k Keeper) DeleteBuyerKey(ctx sdk.Context, key []byte) {
	//store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BuyerKey))
	store := ctx.KVStore(k.storeKey)
	store.Delete(append(types.KeyPrefix(types.BuyerKey), key...))
}

// IterateBuyerItems returns all items from an buyer
func (k Keeper) IterateBuyerItems(ctx sdk.Context, buyer string, cb func(item types.Item) (stop bool)) {

	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.ItemBuyerByBuyerKey(buyer))

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {

		itemID := store.Get(iterator.Key())
		item := k.GetItem(ctx, types.GetItemIDFromBytes(itemID))

		if cb(item) {
			break
		}
	}

}

// GetAllBuyerItems returns all items based on a buyer
func (k Keeper) GetAllBuyerItems(ctx sdk.Context, buyer string) (items []*types.Item) {
	k.IterateBuyerItems(ctx, buyer, func(item types.Item) bool {
		items = append(items, &item)
		return false
	})
	return
}
