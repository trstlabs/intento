package keeper

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/danieljdd/tpp/x/tpp/types"
)

// GetItemCount get the total number of item
func (k Keeper) GetItemCount(ctx sdk.Context) int64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ItemCountKey))
	byteKey := types.KeyPrefix(types.ItemCountKey)
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

// SetItemCount set the total number of item
func (k Keeper) SetItemCount(ctx sdk.Context, count int64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ItemCountKey))
	byteKey := types.KeyPrefix(types.ItemCountKey)
	bz := []byte(strconv.FormatInt(count, 10))
	store.Set(byteKey, bz)
}

// CreateItem creates a item with a new id and update the count
func (k Keeper) CreateItem(ctx sdk.Context, msg types.MsgCreateItem) {
	// Create the item
	count := k.GetItemCount(ctx)

	var estimationcount = fmt.Sprint(msg.Estimationcount)
	var estimationcountHash = sha256.Sum256([]byte(estimationcount + msg.Creator))
	var estimationcountHashString = hex.EncodeToString(estimationcountHash[:])

	submitTime := ctx.BlockHeader().Time

	activePeriod := k.GetParams(ctx).MaxActivePeriod
	endTime := submitTime.Add(activePeriod)

	var item = types.Item{
		Creator:      msg.Creator,
		Seller:       msg.Creator,
		Id:           uint64(count),
		Title:        msg.Title,
		Description:  msg.Description,
		Shippingcost: msg.Shippingcost,
		Localpickup:  msg.Localpickup,
		//Estimationcounthash: estimationcountHashString,

		Tags: msg.Tags,

		Condition:      msg.Condition,
		Shippingregion: msg.Shippingregion,
		Depositamount:  msg.Depositamount,
		Submittime:     submitTime,
		Endtime:        endTime,
	}

	//k.computeKeeper.Instantiate()

	userAddress, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}

	contractAddr, err := k.computeKeeper.Instantiate(ctx, uint64(1), userAddress, msg.Estimationcount, string(count), sdk.NewCoins(sdk.NewCoin("tpp", sdk.ZeroInt())), nil)
	if err != nil {
		return nil, err
	}

	k.BindItemSeller(ctx, item.Id, msg.Creator)
	//works 100% with endtime tx.BlockHeader().Time
	k.InsertInactiveItemQueue(ctx, item.Id, endTime)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.EventTypeItemCreated, sdk.NewAttribute(types.AttributeKeyCreator, item.Creator), sdk.NewAttribute(types.AttributeKeyItemID, strconv.FormatUint(item.Id, 10))),
	)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ItemKey))
	key := append(types.KeyPrefix(types.ItemKey), types.Uint64ToByte(item.Id)...)
	value := k.cdc.MustMarshalBinaryBare(&item)
	store.Set(key, value)

	// Update item count
	k.SetItemCount(ctx, count+1)
}

// SetItem set a specific item in the store
func (k Keeper) SetItem(ctx sdk.Context, item types.Item) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ItemKey))
	b := k.cdc.MustMarshalBinaryBare(&item)
	store.Set(append(types.KeyPrefix(types.ItemKey), types.Uint64ToByte(item.Id)...), b)
}

// GetItem returns a item from its id
func (k Keeper) GetItem(ctx sdk.Context, id uint64) types.Item {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ItemKey))
	var item types.Item
	k.cdc.MustUnmarshalBinaryBare(store.Get(append(types.KeyPrefix(types.ItemKey), types.Uint64ToByte(id)...)), &item)
	return item
}

// HasItem checks if the item exists
func (k Keeper) HasItem(ctx sdk.Context, id uint64) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ItemKey))
	return store.Has(append(types.KeyPrefix(types.ItemKey), types.Uint64ToByte(id)...))

}

// GetItemOwner returns the seller of the item
func (k Keeper) GetItemOwner(ctx sdk.Context, id uint64) string {
	return k.GetItem(ctx, id).Seller
}

// DeleteItem deletes a item
func (k Keeper) DeleteItem(ctx sdk.Context, key uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ItemKey))
	store.Delete(append(types.KeyPrefix(types.ItemKey), types.Uint64ToByte(key)...))
}

// GetAllItem returns all item
func (k Keeper) GetAllItem(ctx sdk.Context) (msgs []types.Item) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ItemKey))
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefix(types.ItemKey))

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var msg types.Item
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &msg)
		msgs = append(msgs, msg)
	}

	return
}

// GetAllInactiveItems returns all inactive item
func (k Keeper) GetAllInactiveItems(ctx sdk.Context) (msgs []*types.Item) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.InactiveItemQueuePrefix)
	iterator := sdk.KVStorePrefixIterator(store, types.InactiveItemQueuePrefix)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var msg types.Item
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &msg)
		msgs = append(msgs, &msg)
	}

	return
}

// HandlePrepayment handles payment
func (k Keeper) HandlePrepayment(ctx sdk.Context, address string, coinToSend sdk.Coin) {

	userAddress, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		panic(err)
	}

	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, userAddress, sdk.NewCoins(coinToSend))
	if err != nil {
		panic(err)
	}

}

// HandleReward handles reward
func (k Keeper) HandleEstimatorReward(ctx sdk.Context, address string, coinToSend sdk.Coin) {

	userAddress, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		panic(err)
	}

	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, userAddress, sdk.NewCoins(coinToSend))
	if err != nil {
		panic(err)
	}

}

// HandleReward handles reward
func (k Keeper) HandleStakingReward(ctx sdk.Context, coinToSend sdk.Coin) {

	//distribute the same reward to the staking pool
	err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.feeCollectorName, sdk.NewCoins(coinToSend))
	if err != nil {
		panic(err)
	}
	//store.Delete(types.KeyPrefix(types.EstimatorKey + key))
}

// MintReward mints coins to a module account
func (k Keeper) MintReward(ctx sdk.Context, coinToSend sdk.Coin) {
	k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(coinToSend))

}

// BurnCoins mints coins from a module account
func (k Keeper) BurnCoins(ctx sdk.Context, coinToSend sdk.Coin) {
	k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(coinToSend))

}
