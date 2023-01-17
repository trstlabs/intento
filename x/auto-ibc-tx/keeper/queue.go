package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/trstlabs/trst/x/auto-ibc-tx/types"
)

// IterateAutoTxsQueue iterates over the items in the inactive autoTx queue
// and performs a callback function
func (k Keeper) IterateAutoTxQueue(ctx sdk.Context, execTime time.Time, cb func(autoTx types.AutoTxInfo) (stop bool)) {
	iterator := k.AutoTxQueueIterator(ctx, execTime)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {

		autoTxID, _ := types.SplitAutoTxQueueKey(iterator.Key())

		autoTx := k.GetAutoTxInfo(ctx, autoTxID)

		//fmt.Printf("info creator is:  %s \n", autoTx.AutoTxInfo.Creator)

		if cb(autoTx) {
			break
		}
	}
}

// GetAutoTxsForBlock returns all expiring autoTxs for a block
func (k Keeper) GetAutoTxsForBlock(ctx sdk.Context) (autoTxs []types.AutoTxInfo) {
	k.IterateAutoTxQueue(ctx, ctx.BlockHeader().Time, func(autoTx types.AutoTxInfo) bool {

		autoTxs = append(autoTxs, autoTx)
		return false
	})
	return
}

// AutoTxQueueIterator returns an sdk.Iterator for all the items in the Inactive Queue that expire by execTime
func (k Keeper) AutoTxQueueIterator(ctx sdk.Context, execTime time.Time) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return store.Iterator(types.AutoTxQueuePrefix, sdk.PrefixEndBytes(types.AutoTxByTimeKey(execTime))) //we check the end of the bites array for the execution time
}

// InsertAutoTxQueue Inserts a autoTx into the auto tx queue
func (k Keeper) InsertAutoTxQueue(ctx sdk.Context, autoTxID uint64, execTime time.Time) {
	store := ctx.KVStore(k.storeKey)
	bz := types.GetBytesForUint(autoTxID)

	//here the key is time+autoTx appended (as bytes) and value is autoTx in bytes
	store.Set(types.AutoTxQueueKey(autoTxID, execTime), bz)
}

// RemoveFromAutoTxQueue removes a autoTx from the Inactive Item Queue
func (k Keeper) RemoveFromAutoTxQueue(ctx sdk.Context, autoTxID uint64, execTime time.Time) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.AutoTxQueueKey(autoTxID, execTime))
}
