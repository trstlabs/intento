package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	//"github.com/tendermint/tendermint/crypto"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/danieljdd/tpp/x/tpp/types"
)

// GetBuyerCount get the total number of buyer
func (k Keeper) GetBuyerCount(ctx sdk.Context) int64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BuyerCountKey))
	byteKey := types.KeyPrefix(types.BuyerCountKey)
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	count, err := strconv.ParseInt(string(bz), 10, 64)
	if err != nil {
		// Panic because the count should be always formattable to int64
		panic("cannot decode count")
	}

	return count
}

// SetBuyerCount set the total number of buyer
func (k Keeper) SetBuyerCount(ctx sdk.Context, count int64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BuyerCountKey))
	byteKey := types.KeyPrefix(types.BuyerCountKey)
	bz := []byte(strconv.FormatInt(count, 10))
	store.Set(byteKey, bz)
}

// CreateBuyer creates a buyer with a new id and update the count
func (k Keeper) CreateBuyer(ctx sdk.Context, msg types.MsgCreateBuyer) {
	// Create the buyer
	count := k.GetBuyerCount(ctx)
	deposit := sdk.NewInt64Coin("tpp", msg.Deposit)

	var buyer = types.Buyer{
		Buyer: msg.Buyer,

		Itemid: msg.Itemid,

		Deposit: deposit,
	}

	buyeraddress, err := sdk.AccAddressFromBech32(msg.Buyer)
	if err != nil {
		panic(err)
	}
	//moduleAcct := sdk.AccAddress(crypto.AddressHash([]byte(types.ModuleName)))

	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, buyeraddress, types.ModuleName, sdk.NewCoins(deposit))
	//sdkError := bankkeeper.keeper.SendCoinsFromAccountToModule(ctx, buyer, ModuleAcct, depositCoinsShipping)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BuyerKey))
	key := types.KeyPrefix(types.BuyerKey + buyer.Itemid)
	value := k.cdc.MustMarshalBinaryBare(&buyer)
	store.Set(key, value)

	// Update buyer count
	k.SetBuyerCount(ctx, count+1)
}

// SetBuyer set a specific buyer in the store
func (k Keeper) SetBuyer(ctx sdk.Context, buyer types.Buyer) {

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BuyerKey))
	b := k.cdc.MustMarshalBinaryBare(&buyer)
	store.Set(types.KeyPrefix(types.BuyerKey+buyer.Itemid), b)
}

// GetBuyer returns a buyer from its id
func (k Keeper) GetBuyer(ctx sdk.Context, key string) types.Buyer {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BuyerKey))
	var buyer types.Buyer
	k.cdc.MustUnmarshalBinaryBare(store.Get(types.KeyPrefix(types.BuyerKey+key)), &buyer)
	return buyer
}

// HasBuyer checks if the buyer exists
func (k Keeper) HasBuyer(ctx sdk.Context, itemid string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BuyerKey))
	return store.Has(types.KeyPrefix(types.BuyerKey + itemid))
}

// GetBuyerOwner returns the creator of the buyer
func (k Keeper) GetBuyerOwner(ctx sdk.Context, key string) string {
	return k.GetBuyer(ctx, key).Buyer
}

// DeleteBuyer deletes a buyer
func (k Keeper) DeleteBuyer(ctx sdk.Context, key string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BuyerKey))
	store.Delete(types.KeyPrefix(types.BuyerKey + key))
}

// GetAllBuyer returns all buyer
func (k Keeper) GetAllBuyer(ctx sdk.Context) (msgs []types.Buyer) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BuyerKey))
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefix(types.BuyerKey))

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var msg types.Buyer
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &msg)
		msgs = append(msgs, msg)
	}

	return
}
