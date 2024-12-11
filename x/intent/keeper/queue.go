package keeper

import (
	"time"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trstlabs/intento/x/intent/types"
)

// IterateActionsQueue iterates over the items in the inactive action queue
// and performs a callback function
func (k Keeper) IterateActionQueue(ctx sdk.Context, execTime time.Time, cb func(action types.ActionInfo) (stop bool)) {
	iterator := k.ActionQueueIterator(ctx, execTime)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {

		actionID, _ := types.SplitActionQueueKey(iterator.Key())

		action := k.GetActionInfo(ctx, actionID)
		if cb(action) {
			break
		}
	}
}

// GetActionsForBlock returns all expiring actions for a block
func (k Keeper) GetActionsForBlock(ctx sdk.Context) (actions []types.ActionInfo) {
	k.IterateActionQueue(ctx, ctx.BlockHeader().Time, func(action types.ActionInfo) bool {
		actions = append(actions, action)
		return false
	})
	return
}

// ActionQueueIterator returns an sdk.Iterator for all the actions in the Inactive Queue that expire by execTime
func (k Keeper) ActionQueueIterator(ctx sdk.Context, execTime time.Time) storetypes.Iterator {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	return store.Iterator(types.ActionQueuePrefix, storetypes.PrefixEndBytes(types.ActionByTimeKey(execTime))) //we check the end of the bites array for the execution time
}

// InsertActionQueue Inserts a action into the action queue
func (k Keeper) InsertActionQueue(ctx sdk.Context, actionID uint64, execTime time.Time) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	bz := types.GetBytesForUint(actionID)

	//here the key is time+action appended (as bytes) and value is action in bytes
	store.Set(types.ActionQueueKey(actionID, execTime), bz)
}

// RemoveFromActionQueue removes a action from the Inactive action queue
func (k Keeper) RemoveFromActionQueue(ctx sdk.Context, action types.ActionInfo) {

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Delete(types.ActionQueueKey(action.ID, action.ExecTime))
}
