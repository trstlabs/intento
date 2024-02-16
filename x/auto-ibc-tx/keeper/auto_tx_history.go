package keeper

import (
	"encoding/binary"

	"time"

	errorsmod "cosmossdk.io/errors"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trstlabs/trst/x/auto-ibc-tx/types"
)

// GetAutoTxHistory
func (k Keeper) GetAutoTxHistory(ctx sdk.Context, autoTxHistoryID uint64) types.AutoTxHistory {
	store := ctx.KVStore(k.storeKey)
	var autoTxHistory types.AutoTxHistory
	autoTxHistoryBz := store.Get(types.GetAutoTxHistoryKey(autoTxHistoryID))

	k.cdc.MustUnmarshal(autoTxHistoryBz, &autoTxHistory)
	return autoTxHistory
}

// TryGetAutoTxHistory
func (k Keeper) TryGetAutoTxHistory(ctx sdk.Context, autoTxHistoryID uint64) (types.AutoTxHistory, error) {
	store := ctx.KVStore(k.storeKey)
	var autoTxHistory types.AutoTxHistory
	autoTxHistoryBz := store.Get(types.GetAutoTxHistoryKey(autoTxHistoryID))

	err := k.cdc.Unmarshal(autoTxHistoryBz, &autoTxHistory)
	if err != nil {
		return types.AutoTxHistory{}, err
	}
	return autoTxHistory, nil
}

func (k Keeper) SetAutoTxHistory(ctx sdk.Context, autoTxId uint64, autoTxHistory *types.AutoTxHistory) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetAutoTxHistoryKey(autoTxId), k.cdc.MustMarshal(autoTxHistory))
}

func (k Keeper) importAutoTxHistory(ctx sdk.Context, autoTxHistoryId uint64, AutoTxHistory types.AutoTxHistory) error {

	store := ctx.KVStore(k.storeKey)
	key := types.GetAutoTxHistoryKey(autoTxHistoryId)
	if store.Has(key) {
		return errorsmod.Wrapf(types.ErrDuplicate, "duplicate code: %d", autoTxHistoryId)
	}
	// 0x01 | autoTxHistoryId (uint64) -> AutoTxHistory
	store.Set(key, k.cdc.MustMarshal(&AutoTxHistory))
	return nil
}

func (k Keeper) IterateAutoTxHistorys(ctx sdk.Context, cb func(uint64, types.AutoTxHistory) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.AutoTxKeyPrefix)
	iter := prefixStore.Iterator(nil, nil)
	for ; iter.Valid(); iter.Next() {
		var c types.AutoTxHistory
		k.cdc.MustUnmarshal(iter.Value(), &c)
		// cb returns true to stop early
		if cb(binary.BigEndian.Uint64(iter.Key()), c) {
			return
		}
	}
}

func (k Keeper) AddAutoTxHistory(ctx sdk.Context, autoTxHistory *types.AutoTxHistory, autoTx *types.AutoTxInfo, actualExecTime time.Time, execFee sdk.Coin, executedLocally bool, msgResponses []*cdctypes.Any, err ...string) {
	historyEntry := types.AutoTxHistoryEntry{
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
	k.SetAutoTxHistory(ctx, autoTx.TxID, autoTxHistory)

}
