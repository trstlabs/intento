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
func (k Keeper) GetActionHistory(ctx sdk.Context, autoTxHistoryID uint64) types.ActionHistory {
	store := ctx.KVStore(k.storeKey)
	var autoTxHistory types.ActionHistory
	autoTxHistoryBz := store.Get(types.GetActionHistoryKey(autoTxHistoryID))

	k.cdc.MustUnmarshal(autoTxHistoryBz, &autoTxHistory)
	return autoTxHistory
}

// TryGetActionHistory
func (k Keeper) TryGetActionHistory(ctx sdk.Context, autoTxHistoryID uint64) (types.ActionHistory, error) {
	store := ctx.KVStore(k.storeKey)
	var autoTxHistory types.ActionHistory
	autoTxHistoryBz := store.Get(types.GetActionHistoryKey(autoTxHistoryID))

	err := k.cdc.Unmarshal(autoTxHistoryBz, &autoTxHistory)
	if err != nil {
		return types.ActionHistory{}, err
	}
	return autoTxHistory, nil
}

func (k Keeper) SetActionHistory(ctx sdk.Context, actionId uint64, autoTxHistory *types.ActionHistory) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetActionHistoryKey(actionId), k.cdc.MustMarshal(autoTxHistory))
}

func (k Keeper) importActionHistory(ctx sdk.Context, autoTxHistoryId uint64, ActionHistory types.ActionHistory) error {

	store := ctx.KVStore(k.storeKey)
	key := types.GetActionHistoryKey(autoTxHistoryId)
	if store.Has(key) {
		return errorsmod.Wrapf(types.ErrDuplicate, "duplicate code: %d", autoTxHistoryId)
	}
	// 0x01 | autoTxHistoryId (uint64) -> ActionHistory
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

func (k Keeper) AddActionHistory(ctx sdk.Context, autoTxHistory *types.ActionHistory, autoTx *types.ActionInfo, actualExecTime time.Time, execFee sdk.Coin, executedLocally bool, msgResponses []*cdctypes.Any, err ...string) {
	historyEntry := types.ActionHistoryEntry{
		ScheduledExecTime: autoTx.ExecTime,
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

	autoTxHistory.History = append(autoTxHistory.History, historyEntry)
	k.SetActionHistory(ctx, autoTx.ID, autoTxHistory)

}
