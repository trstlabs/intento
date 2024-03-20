package keeper

import (
	"encoding/binary"

	"time"

	errorsmod "cosmossdk.io/errors"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trstlabs/intento/x/intent/types"
)

// GetActionHistory
func (k Keeper) GetActionHistory(ctx sdk.Context, actionHistoryID uint64) types.ActionHistory {
	store := ctx.KVStore(k.storeKey)
	var actionHistory types.ActionHistory
	actionHistoryBz := store.Get(types.GetActionHistoryKey(actionHistoryID))

	k.cdc.MustUnmarshal(actionHistoryBz, &actionHistory)
	return actionHistory
}

// TryGetActionHistory
func (k Keeper) TryGetActionHistory(ctx sdk.Context, actionHistoryID uint64) (types.ActionHistory, error) {
	store := ctx.KVStore(k.storeKey)
	var actionHistory types.ActionHistory
	actionHistoryBz := store.Get(types.GetActionHistoryKey(actionHistoryID))

	err := k.cdc.Unmarshal(actionHistoryBz, &actionHistory)
	if err != nil {
		return types.ActionHistory{}, err
	}
	return actionHistory, nil
}

func (k Keeper) SetActionHistory(ctx sdk.Context, actionId uint64, actionHistory *types.ActionHistory) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetActionHistoryKey(actionId), k.cdc.MustMarshal(actionHistory))
}

func (k Keeper) importActionHistory(ctx sdk.Context, actionHistoryId uint64, ActionHistory types.ActionHistory) error {

	store := ctx.KVStore(k.storeKey)
	key := types.GetActionHistoryKey(actionHistoryId)
	if store.Has(key) {
		return errorsmod.Wrapf(types.ErrDuplicate, "duplicate code: %d", actionHistoryId)
	}
	// 0x01 | actionHistoryId (uint64) -> ActionHistory
	store.Set(key, k.cdc.MustMarshal(&ActionHistory))
	return nil
}

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

func (k Keeper) AddActionHistory(ctx sdk.Context, actionHistory *types.ActionHistory, action *types.ActionInfo, actualExecTime time.Time, execFee sdk.Coin, executedLocally bool, msgResponses []*cdctypes.Any, err ...string) {
	historyEntry := types.ActionHistoryEntry{
		ScheduledExecTime: action.ExecTime,
		ActualExecTime:    actualExecTime,
		ExecFee:           execFee,
	}

	if executedLocally {
		historyEntry.Executed = true
		historyEntry.MsgResponses = msgResponses
	}

	if len(err) == 1 && err[0] != "" {
		historyEntry.Errors = append(historyEntry.Errors, err[0])
		historyEntry.Executed = false
	}

	actionHistory.History = append(actionHistory.History, historyEntry)
	k.SetActionHistory(ctx, action.ID, actionHistory)

}
