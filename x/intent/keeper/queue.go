package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/trstlabs/intento/x/intent/types"
)

// IterateActionsQueue iterates over the items in the inactive autoTx queue
// and performs a callback function
func (k Keeper) IterateActionQueue(ctx sdk.Context, execTime time.Time, cb func(autoTx types.ActionInfo) (stop bool)) {
	iterator := k.ActionQueueIterator(ctx, execTime)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {

		autoID, _ := types.SplitActionQueueKey(iterator.Key())
		//fmt.Printf("autoTx id is:  %v \n", autoID)
		autoTx := k.GetActionInfo(ctx, autoID)

		if cb(autoTx) {
			break
		}
	}
}

// GetActionsForBlock returns all expiring autoTxs for a block
func (k Keeper) GetActionsForBlock(ctx sdk.Context) (autoTxs []types.ActionInfo) {
	k.IterateActionQueue(ctx, ctx.BlockHeader().Time, func(autoTx types.ActionInfo) bool {
		autoTxs = append(autoTxs, autoTx)
		return false
	})
	return
}

// ActionQueueIterator returns an sdk.Iterator for all the items in the Inactive Queue that expire by execTime
func (k Keeper) ActionQueueIterator(ctx sdk.Context, execTime time.Time) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return store.Iterator(types.ActionQueuePrefix, sdk.PrefixEndBytes(types.ActionByTimeKey(execTime))) //we check the end of the bites array for the execution time
}

// InsertActionQueue Inserts a autoTx into the auto tx queue
func (k Keeper) InsertActionQueue(ctx sdk.Context, autoID uint64, execTime time.Time) {
	store := ctx.KVStore(k.storeKey)
	bz := types.GetBytesForUint(autoID)

	//here the key is time+autoTx appended (as bytes) and value is autoTx in bytes
	store.Set(types.ActionQueueKey(autoID, execTime), bz)
}

// RemoveFromActionQueue removes a autoTx from the Inactive Item Queue
func (k Keeper) RemoveFromActionQueue(ctx sdk.Context, autoTx types.ActionInfo) {

	store := ctx.KVStore(k.storeKey)
	store.Delete(types.ActionQueueKey(autoTx.ID, autoTx.ExecTime))
}
