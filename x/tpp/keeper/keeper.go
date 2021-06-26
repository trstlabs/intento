package keeper

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	wasm "github.com/danieljdd/tpp/go-cosmwasm"
	"github.com/danieljdd/tpp/x/tpp/types"
	"github.com/tendermint/tendermint/libs/log"
)

type (
	Keeper struct {
		cdc              codec.Marshaler
		storeKey         sdk.StoreKey
		memKey           sdk.StoreKey
		paramSpace       paramtypes.Subspace
		accountKeeper    types.AccountKeeper
		bankKeeper       types.BankKeeper
		feeCollectorName string
		wasmer           wasm.Wasmer
		computeKeeper    types.ComputeKeeper
	}
)

func NewKeeper(cdc codec.Marshaler, storeKey, memKey sdk.StoreKey, paramSpace paramtypes.Subspace, ak types.AccountKeeper, bk types.BankKeeper, feeCollectorName string, homeDir string, wasmConfig types.WasmConfig, supportedFeatures string, ck types.ComputeKeeper) *Keeper {

	wasmer, err := wasm.NewWasmer(filepath.Join(homeDir, "wasm"), supportedFeatures, wasmConfig.CacheSize)
	if err != nil {
		panic(err)
	}

	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(ParamKeyTable())
	}

	// ensure reward pool module account is set
	if addr := ak.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	return &Keeper{cdc: cdc,
		storeKey:         storeKey,
		memKey:           memKey,
		paramSpace:       paramSpace,
		bankKeeper:       bk,
		accountKeeper:    ak,
		feeCollectorName: feeCollectorName,
		wasmer:           *wasmer,
		computeKeeper:    ck,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// IterateInactiveItemsQueue iterates over the items in the inactive item queue
// and performs a callback function
func (k Keeper) IterateInactiveItemsQueue(ctx sdk.Context, endTime time.Time, cb func(item types.Item) (stop bool)) {
	iterator := k.InactiveItemQueueIterator(ctx, endTime)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		//get the itemid from endTime (key)
		itemID, _ := types.SplitInactiveItemQueueKey(iterator.Key())
		item := k.GetItem(ctx, itemID)
		if cb(item) {
			break
		}
	}
}

// InactiveItemQueueIterator returns an sdk.Iterator for all the items in the Inactive Queue that expire by endTime
func (k Keeper) InactiveItemQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return store.Iterator(types.InactiveItemQueuePrefix, sdk.PrefixEndBytes(types.InactiveItemByTimeKey(endTime))) //we check the end of the bites array for the end time
}

// InsertInactiveItemQueue Inserts a itemid into the inactive item queue at endTime
func (k Keeper) InsertInactiveItemQueue(ctx sdk.Context, itemid uint64, endTime time.Time) {
	store := ctx.KVStore(k.storeKey)
	bz := types.Uint64ToByte(itemid)

	//here the key is time+itemid appended (as bytes) and value is itemid in bytes
	store.Set(types.InactiveItemQueueKey(itemid, endTime), bz)
}

// RemoveFromInactiveItemQueue removes a itemid from the Inactive Item Queue
func (k Keeper) RemoveFromInactiveItemQueue(ctx sdk.Context, itemid uint64, endTime time.Time) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.InactiveItemQueueKey(itemid, endTime))
}

/////Seller functions

// BindItemSeller binds a itemid with the seller address
func (k Keeper) BindItemSeller(ctx sdk.Context, itemid uint64, seller string) {
	store := ctx.KVStore(k.storeKey)
	bz := types.Uint64ToByte(itemid)

	//here the key is seller+itemid appended (as bytes) and value is itemid in bytes
	store.Set(types.ItemSellerKey(itemid, seller), bz)
}

// RemoveFromInactiveItemQueue removes a itemid from the seller
func (k Keeper) RemoveFromItemSeller(ctx sdk.Context, itemid uint64, seller string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.ItemSellerKey(itemid, seller))
}

// GetAllSellerItems returns all seller items
func (k Keeper) IterateItems(ctx sdk.Context, seller string, cb func(item types.Item) (stop bool)) {
	//iterator := k.ItemSellerIterator(ctx, seller)
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.ItemSellerBySellerKey(seller))

	defer iterator.Close()
	//store := ctx.KVStore(k.storeKey)
	for ; iterator.Valid(); iterator.Next() {

		//var item types.Item
		//get the itemid from endTime (key)
		itemID := store.Get(iterator.Key())
		item := k.GetItem(ctx, types.GetItemIDFromBytes(itemID))
		//items = append(items, &item)
		if cb(item) {
			break
		}

	}
	//return //items
}

// GetAllSellerItems returns all seller items on chain based on the seller
func (k Keeper) GetAllSellerItems(ctx sdk.Context, seller string) (items []*types.Item) {
	k.IterateItems(ctx, seller, func(item types.Item) bool {
		items = append(items, &item)
		return false
	})
	return
}
