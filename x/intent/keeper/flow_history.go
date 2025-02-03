package keeper

import (
	"encoding/binary"

	"time"

	"cosmossdk.io/math"

	"cosmossdk.io/store/prefix"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trstlabs/intento/x/intent/types"
)

func (k Keeper) GetLatestFlowHistoryEntry(ctx sdk.Context, flowId uint64) (*types.FlowHistoryEntry, error) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	prefixStore := prefix.NewStore(store, types.GetFlowHistoryKey(flowId))

	// Use a reverse prefix iterator to start from the latest entry
	iterator := prefixStore.ReverseIterator(nil, nil)
	defer iterator.Close()

	if iterator.Valid() {
		var entry types.FlowHistoryEntry
		err := k.cdc.Unmarshal(iterator.Value(), &entry)
		if err != nil {
			return nil, err
		}
		return &entry, nil
	}

	return nil, nil // Or appropriate error to indicate not found
}

// GetFlowHistory retrieves all history entries for a specific flowId.
func (k Keeper) GetFlowHistory(ctx sdk.Context, flowId uint64) ([]types.FlowHistoryEntry, error) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	prefixStore := prefix.NewStore(store, types.GetFlowHistoryKey(flowId))

	iterator := prefixStore.Iterator(nil, nil)
	defer iterator.Close()

	var historyEntries []types.FlowHistoryEntry
	for ; iterator.Valid(); iterator.Next() {
		var entry types.FlowHistoryEntry
		err := k.cdc.Unmarshal(iterator.Value(), &entry)
		if err != nil {
			return nil, err
		}
		historyEntries = append(historyEntries, entry)
	}

	return historyEntries, nil
}

// MustGetFlowHistory tries to retrieve the flow history, returning nil if it fails.
func (k Keeper) MustGetFlowHistory(ctx sdk.Context, flowId uint64) []types.FlowHistoryEntry {
	flowHistory, err := k.GetFlowHistory(ctx, flowId)
	if err != nil {
		return nil
	}
	return flowHistory
}

func (k Keeper) SetFlowHistoryEntry(ctx sdk.Context, flowId uint64, entry *types.FlowHistoryEntry) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	// Generate a unique sequence number for this entry. Alternatively, you can use a timestamp.
	// This assumes you have a function to get the next available sequence number or you can store the count somewhere.
	sequence := k.GetNextFlowHistorySequence(ctx, flowId)

	// Composite key: FlowHistoryKey + FlowID + Sequence
	key := append(types.GetFlowHistoryKey(flowId), sdk.Uint64ToBigEndian(sequence)...)

	// Store the entry
	store.Set(key, k.cdc.MustMarshal(entry))
}

func (k Keeper) GetNextFlowHistorySequence(ctx sdk.Context, flowId uint64) uint64 {
	// This is a simplified example. You need to implement the logic to get the next sequence number.
	// This could involve getting the current count from the store and incrementing it.
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	sequenceKey := append(types.FlowHistorySequencePrefix, sdk.Uint64ToBigEndian(flowId)...)
	sequenceBytes := store.Get(sequenceKey)
	var sequence uint64
	if sequenceBytes != nil {
		sequence = sdk.BigEndianToUint64(sequenceBytes)
	}
	sequence++
	store.Set(sequenceKey, sdk.Uint64ToBigEndian(sequence))
	return sequence
}

func (k Keeper) IterateFlowHistorys(ctx sdk.Context, cb func(uint64, types.FlowHistory) bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	prefixStore := prefix.NewStore(store, types.FlowKeyPrefix)
	iter := prefixStore.Iterator(nil, nil)
	for ; iter.Valid(); iter.Next() {
		var c types.FlowHistory
		k.cdc.MustUnmarshal(iter.Value(), &c)
		// cb returns true to stop early
		if cb(binary.BigEndian.Uint64(iter.Key()), c) {
			return
		}
	}
}

func (k Keeper) addFlowHistoryEntry(ctx sdk.Context, flow *types.FlowInfo, actualExecTime time.Time, execFee sdk.Coin, executedLocally bool, msgResponses []*cdctypes.Any, errorString string) {
	historyEntry := types.FlowHistoryEntry{
		ScheduledExecTime: flow.ExecTime,
		ActualExecTime:    actualExecTime,
		ExecFee:           execFee,
	}
	if flow.Configuration.SaveResponses {
		historyEntry.MsgResponses = append(historyEntry.MsgResponses, msgResponses...)
		for i, comparison := range flow.Conditions.Comparisons {
			if comparison.ICQConfig != nil {
				historyEntry.QueryResponses = append(historyEntry.QueryResponses, string(comparison.ICQConfig.Response))
				flow.Conditions.Comparisons[i].ICQConfig.Response = nil
			}
		}
		for i, feedbackLoop := range flow.Conditions.FeedbackLoops {
			if feedbackLoop.ICQConfig != nil {
				historyEntry.QueryResponses = append(historyEntry.QueryResponses, string(feedbackLoop.ICQConfig.Response))
				flow.Conditions.FeedbackLoops[i].ICQConfig.Response = nil
			}
		}

	}
	if errorString != "" {
		historyEntry.Errors = append(historyEntry.Errors, errorString)
	}
	if executedLocally {
		historyEntry.Executed = true
		// if flow.Configuration.SaveResponses {
		// 	historyEntry.MsgResponses = msgResponses
		// }
	}

	k.SetFlowHistoryEntry(ctx, flow.ID, &historyEntry)

}

func (k Keeper) SetCurrentFlowHistoryEntry(ctx sdk.Context, flowId uint64, entry *types.FlowHistoryEntry) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	sequenceKey := append(types.FlowHistorySequencePrefix, sdk.Uint64ToBigEndian(flowId)...)
	sequenceBytes := store.Get(sequenceKey)
	var sequence uint64
	if sequenceBytes != nil {
		sequence = sdk.BigEndianToUint64(sequenceBytes)
	}
	// Composite key: FlowHistoryKey + FlowID + Sequence
	key := append(types.GetFlowHistoryKey(flowId), sdk.Uint64ToBigEndian(sequence)...)

	// Store the entry
	store.Set(key, k.cdc.MustMarshal(entry))
}

func (k Keeper) HasFlowHistoryEntry(ctx sdk.Context, flowId uint64) bool {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	sequenceKey := append(types.FlowHistorySequencePrefix, sdk.Uint64ToBigEndian(flowId)...)
	return store.Has(sequenceKey)

}

// func (k Keeper) getCurrentFlowHistoryEntry(ctx sdk.Context, flowId uint64) (*types.FlowHistoryEntry, bool) {
// 	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

// 	// Retrieve the current sequence for the flowId
// 	sequenceKey := append(types.FlowHistorySequencePrefix, sdk.Uint64ToBigEndian(flowId)...)
// 	sequenceBytes := store.Get(sequenceKey)
// 	if sequenceBytes == nil {
// 		// No sequence found, so no entry exists
// 		return nil, false
// 	}

// 	// Decode the current sequence
// 	sequence := sdk.BigEndianToUint64(sequenceBytes)

// 	// Composite key: FlowHistoryKey + FlowID + Sequence (latest entry)
// 	key := append(types.GetFlowHistoryKey(flowId), sdk.Uint64ToBigEndian(sequence)...)

// 	// Fetch the current entry
// 	entryBytes := store.Get(key)
// 	if entryBytes == nil {
// 		// No entry exists at the latest sequence
// 		return nil, false
// 	}
// 	var entry types.FlowHistoryEntry
// 	k.cdc.MustUnmarshal(entryBytes, &entry)

// 	return &entry, true
// }

// we may reimplement this as a configuration-based gas fee
func (k Keeper) CalculateTimeBasedFlexFee(ctx sdk.Context, flow types.FlowInfo) math.Int {
	historyEntry, _ := k.GetLatestFlowHistoryEntry(ctx, flow.ID)

	if historyEntry != nil {
		prevEntryTime := historyEntry.ActualExecTime
		period := (flow.ExecTime.Sub(prevEntryTime))
		return math.NewInt(int64(period.Milliseconds()))
	}

	period := flow.ExecTime.Sub(flow.StartTime)
	if period.Seconds() <= 60 {
		//base fee so we do not have a zero fee
		return math.NewInt(6_000)
	}
	return math.NewInt(int64(period.Seconds() * 10))
}
