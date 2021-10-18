package keeper

import (

	//"github.com/tendermint/tendermint/crypto"
	//	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/danieljdd/trst/x/item/types"
)

// Prepayment creates a buyer with a new id and update the count
func (k Keeper) Prepayment(ctx sdk.Context, msg types.MsgPrepayment) {
	// Create the buyer

	deposit := sdk.NewInt64Coin("utrst", msg.Deposit)

	buyeraddress, err := sdk.AccAddressFromBech32(msg.Buyer)
	if err != nil {
		panic(err)
	}

	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, buyeraddress, types.ModuleName, sdk.NewCoins(deposit))
	if err != nil {
		panic(err)
	}

	k.SetBuyer(ctx, msg.Itemid, msg.Buyer)

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

// Withdrawal deletes a buyer
func (k Keeper) Withdrawal(ctx sdk.Context, key []byte) {
	//store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BuyerKey))
	store := ctx.KVStore(k.storeKey)
	store.Delete(append(types.KeyPrefix(types.BuyerKey), key...))
}

// IterateBuyerItems returns all buyer items
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
