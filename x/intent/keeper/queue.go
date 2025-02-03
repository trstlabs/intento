package keeper

import (
	"time"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trstlabs/intento/x/intent/types"
)

// IterateFlowsQueue iterates over the items in the inactive flow queue
// and performs a callback function
func (k Keeper) IterateFlowQueue(ctx sdk.Context, execTime time.Time, cb func(flow types.FlowInfo) (stop bool)) {
	iterator := k.FlowQueueIterator(ctx, execTime)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {

		flowID, _ := types.SplitFlowQueueKey(iterator.Key())

		flow := k.GetFlowInfo(ctx, flowID)
		if cb(flow) {
			break
		}
	}
}

// GetFlowsForBlock returns all expiring flows for a block
func (k Keeper) GetFlowsForBlock(ctx sdk.Context) (flows []types.FlowInfo) {
	k.IterateFlowQueue(ctx, ctx.BlockHeader().Time, func(flow types.FlowInfo) bool {
		flows = append(flows, flow)
		return false
	})
	return
}

// FlowQueueIterator returns an sdk.Iterator for all the flows in the Inactive Queue that expire by execTime
func (k Keeper) FlowQueueIterator(ctx sdk.Context, execTime time.Time) storetypes.Iterator {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	return store.Iterator(types.FlowQueuePrefix, storetypes.PrefixEndBytes(types.FlowByTimeKey(execTime))) //we check the end of the bites array for the execution time
}

// InsertFlowQueue Inserts a flow into the flow queue
func (k Keeper) InsertFlowQueue(ctx sdk.Context, flowID uint64, execTime time.Time) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	bz := types.GetBytesForUint(flowID)

	//here the key is time+flow appended (as bytes) and value is flow in bytes
	store.Set(types.FlowQueueKey(flowID, execTime), bz)
}

// RemoveFromFlowQueue removes a flow from the Inactive flow queue
func (k Keeper) RemoveFromFlowQueue(ctx sdk.Context, flow types.FlowInfo) {

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Delete(types.FlowQueueKey(flow.ID, flow.ExecTime))
}
