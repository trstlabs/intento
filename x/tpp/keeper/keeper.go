package keeper

import (
	"fmt"
	"time"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/danieljdd/tpp/x/tpp/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

type (
	Keeper struct {
		cdc      codec.Marshaler
		storeKey sdk.StoreKey
		memKey   sdk.StoreKey
		paramSpace paramtypes.Subspace
		accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	}
)

func NewKeeper(cdc codec.Marshaler, storeKey, memKey sdk.StoreKey, paramSpace paramtypes.Subspace, ak types.AccountKeeper, bk types.BankKeeper) *Keeper {

	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(ParamKeyTable())
	}

	// ensure reward pool module account is set
	if addr := ak.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}



	return &Keeper{cdc:      cdc,
		storeKey: storeKey,
		memKey:   memKey,
		paramSpace: paramSpace,
		bankKeeper:    bk,
		accountKeeper: ak,
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
	return store.Iterator(types.InactiveItemQueuePrefix, sdk.PrefixEndBytes(types.InactiveItemByTimeKey(endTime)))//we check the end of the bites array for the end time
}

// InsertInactiveItemQueue Inserts a itemid into the inactive item queue at endTime
func (k Keeper) InsertInactiveItemQueue(ctx sdk.Context, itemid string, endTime time.Time) {
	store := ctx.KVStore(k.storeKey)
	bz := []byte(itemid)

	//here the key is time+itemid appended (as bytes) and value is itemid in bytes
	store.Set(types.InactiveItemQueueKey(itemid, endTime), bz)
}

// RemoveFromInactiveItemQueue removes a itemid from the Inactive Item Queue
func (k Keeper) RemoveFromInactiveItemQueue(ctx sdk.Context, itemid string, endTime time.Time) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.InactiveItemQueueKey(itemid, endTime))
}