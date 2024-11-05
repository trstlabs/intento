package keeper

import (
	"encoding/binary"

	"time"

	sdkmath "cosmossdk.io/math"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trstlabs/intento/x/intent/types"
)

func (k Keeper) GetLatestActionHistoryEntry(ctx sdk.Context, actionId uint64) (*types.ActionHistoryEntry, error) {
	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, types.GetActionHistoryKey(actionId))

	// Use a reverse prefix iterator to start from the latest entry
	iterator := prefixStore.ReverseIterator(nil, nil)
	defer iterator.Close()

	if iterator.Valid() {
		var entry types.ActionHistoryEntry
		err := k.cdc.Unmarshal(iterator.Value(), &entry)
		if err != nil {
			return nil, err
		}
		return &entry, nil
	}

	return nil, nil // Or appropriate error to indicate not found
}

// GetActionHistory retrieves all history entries for a specific actionId.
func (k Keeper) GetActionHistory(ctx sdk.Context, actionId uint64) ([]types.ActionHistoryEntry, error) {
	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, types.GetActionHistoryKey(actionId))

	iterator := prefixStore.Iterator(nil, nil)
	defer iterator.Close()

	var historyEntries []types.ActionHistoryEntry
	for ; iterator.Valid(); iterator.Next() {
		var entry types.ActionHistoryEntry
		err := k.cdc.Unmarshal(iterator.Value(), &entry)
		if err != nil {
			return nil, err
		}
		historyEntries = append(historyEntries, entry)
	}

	return historyEntries, nil
}

// MustGetActionHistory tries to retrieve the auto transaction history, returning nil if it fails.
func (k Keeper) MustGetActionHistory(ctx sdk.Context, actionId uint64) []types.ActionHistoryEntry {
	actionHistory, err := k.GetActionHistory(ctx, actionId)
	if err != nil {
		return nil
	}
	return actionHistory
}

func (k Keeper) SetActionHistoryEntry(ctx sdk.Context, actionId uint64, entry *types.ActionHistoryEntry) {
	store := ctx.KVStore(k.storeKey)

	// Generate a unique sequence number for this entry. Alternatively, you can use a timestamp.
	// This assumes you have a function to get the next available sequence number or you can store the count somewhere.
	sequence := k.GetNextActionHistorySequence(ctx, actionId)

	// Composite key: ActionHistoryKey + ActionId + Sequence
	key := append(types.GetActionHistoryKey(actionId), sdk.Uint64ToBigEndian(sequence)...)

	// Store the entry
	store.Set(key, k.cdc.MustMarshal(entry))
}

func (k Keeper) GetNextActionHistorySequence(ctx sdk.Context, actionId uint64) uint64 {
	// This is a simplified example. You need to implement the logic to get the next sequence number.
	// This could involve getting the current count from the store and incrementing it.
	store := ctx.KVStore(k.storeKey)
	sequenceKey := append(types.ActionHistorySequencePrefix, sdk.Uint64ToBigEndian(actionId)...)
	sequenceBytes := store.Get(sequenceKey)
	var sequence uint64
	if sequenceBytes != nil {
		sequence = sdk.BigEndianToUint64(sequenceBytes)
	}
	sequence++
	store.Set(sequenceKey, sdk.Uint64ToBigEndian(sequence))
	return sequence
}

// func (k Keeper) importActionHistory(ctx sdk.Context, actionHistoryId uint64, ActionHistory types.ActionHistory) error {

// 	store := ctx.KVStore(k.storeKey)
// 	key := types.GetActionHistoryKey(actionHistoryId)
// 	if store.Has(key) {
// 		return errorsmod.Wrapf(types.ErrDuplicate, "duplicate code: %d", actionHistoryId)
// 	}
// 	// 0x01 | actionHistoryId (uint64) -> ActionHistory
// 	store.Set(key, k.cdc.MustMarshal(&ActionHistory))
// 	return nil
// }

func (k Keeper) IterateActionHistorys(ctx sdk.Context, cb func(uint64, types.ActionHistory) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.ActionKeyPrefix)
	iter := prefixStore.Iterator(nil, nil)
	for ; iter.Valid(); iter.Next() {
		var c types.ActionHistory
		k.cdc.MustUnmarshal(iter.Value(), &c)
		// cb returns true to stop early
		if cb(binary.BigEndian.Uint64(iter.Key()), c) {
			return
		}
	}
}

func (k Keeper) addActionHistory(ctx sdk.Context, action *types.ActionInfo, actualExecTime time.Time, execFee sdk.Coin, executedLocally bool, msgResponses []*cdctypes.Any, queryResponse string, errorString string) {
	historyEntry := types.ActionHistoryEntry{
		ScheduledExecTime: action.ExecTime,
		ActualExecTime:    actualExecTime,
		ExecFee:           execFee,
	}
	if action.Configuration.SaveResponses {
		historyEntry.MsgResponses = append(historyEntry.MsgResponses, msgResponses...)
		historyEntry.QueryResponse = queryResponse

	}
	if errorString != "" {
		historyEntry.Errors = append(historyEntry.Errors, errorString)
	}
	if executedLocally {
		historyEntry.Executed = true
		if action.Configuration.SaveResponses {
			historyEntry.MsgResponses = msgResponses
		}
	}

	k.SetActionHistoryEntry(ctx, action.ID, &historyEntry)

}

func (k Keeper) SetCurrentActionHistoryEntry(ctx sdk.Context, actionId uint64, entry *types.ActionHistoryEntry) {
	store := ctx.KVStore(k.storeKey)
	sequenceKey := append(types.ActionHistorySequencePrefix, sdk.Uint64ToBigEndian(actionId)...)
	sequenceBytes := store.Get(sequenceKey)
	var sequence uint64
	if sequenceBytes != nil {
		sequence = sdk.BigEndianToUint64(sequenceBytes)
	}
	// Composite key: ActionHistoryKey + ActionId + Sequence
	key := append(types.GetActionHistoryKey(actionId), sdk.Uint64ToBigEndian(sequence)...)

	// Store the entry
	store.Set(key, k.cdc.MustMarshal(entry))
}

func (k Keeper) getCurrentActionHistoryEntry(ctx sdk.Context, actionId uint64) (*types.ActionHistoryEntry, bool) {
	store := ctx.KVStore(k.storeKey)

	// Retrieve the current sequence for the actionId
	sequenceKey := append(types.ActionHistorySequencePrefix, sdk.Uint64ToBigEndian(actionId)...)
	sequenceBytes := store.Get(sequenceKey)
	if sequenceBytes == nil {
		// No sequence found, so no entry exists
		return nil, false
	}

	// Decode the current sequence
	sequence := sdk.BigEndianToUint64(sequenceBytes)

	// Composite key: ActionHistoryKey + ActionId + Sequence (latest entry)
	key := append(types.GetActionHistoryKey(actionId), sdk.Uint64ToBigEndian(sequence)...)

	// Fetch the current entry
	entryBytes := store.Get(key)
	if entryBytes == nil {
		// No entry exists at the latest sequence
		return nil, false
	}
	var entry types.ActionHistoryEntry
	k.cdc.MustUnmarshal(entryBytes, &entry)

	return &entry, true
}

// we may reimplement this as a configuration-based gas fee
func (k Keeper) CalculateTimeBasedFlexFee(ctx sdk.Context, action types.ActionInfo) sdkmath.Int {
	historyEntry, _ := k.GetLatestActionHistoryEntry(ctx, action.ID)

	if historyEntry != nil {
		prevEntryTime := historyEntry.ActualExecTime
		period := (action.ExecTime.Sub(prevEntryTime))
		return sdk.NewInt(int64(period.Milliseconds()))
	}

	period := action.ExecTime.Sub(action.StartTime)
	if period.Seconds() <= 60 {
		//base fee so we do not have a zero fee
		return sdk.NewInt(6_000)
	}
	return sdk.NewInt(int64(period.Seconds() * 10))
}
