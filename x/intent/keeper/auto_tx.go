package keeper

import (
	"encoding/binary"
	"strconv"
	"time"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/store/prefix"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/trstlabs/intento/x/intent/types"
)

// GetActionInfo
func (k Keeper) GetActionInfo(ctx sdk.Context, actionID uint64) types.ActionInfo {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	var action types.ActionInfo
	actionBz := store.Get(types.GetActionKey(actionID))

	k.cdc.MustUnmarshal(actionBz, &action)
	return action
}

// TryGetActionInfo
func (k Keeper) TryGetActionInfo(ctx sdk.Context, actionID uint64) (types.ActionInfo, error) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	var action types.ActionInfo
	actionBz := store.Get(types.GetActionKey(actionID))
	if actionBz == nil {
		return types.ActionInfo{}, errorsmod.Wrapf(types.ErrNotFound, "action")
	}
	err := k.cdc.Unmarshal(actionBz, &action)
	if err != nil {
		return types.ActionInfo{}, err
	}

	return action, nil
}

func (k Keeper) SetActionInfo(ctx sdk.Context, action *types.ActionInfo) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Set(types.GetActionKey(action.ID), k.cdc.MustMarshal(action))
}

func (k Keeper) CreateAction(ctx sdk.Context, owner sdk.AccAddress, label string, msgs []*cdctypes.Any, duration time.Duration, interval time.Duration, startAt time.Time, feeFunds sdk.Coins, configuration types.ExecutionConfiguration, hostedConfig types.HostedConfig, portID string, connectionId string, hostConnectionId string, conditions types.ExecutionConditions) error {

	id := k.autoIncrementID(ctx, types.KeyLastID)
	actionAddress, err := k.createFeeAccount(ctx, id, owner, feeFunds)
	if err != nil {
		return err
	}

	endTime, execTime := k.calculateTimeAndInsertQueue(ctx, startAt, duration, id, interval)

	icaConfig := types.ICAConfig{
		PortID:           portID,
		ConnectionID:     connectionId,
		HostConnectionID: hostConnectionId,
	}

	action := types.ActionInfo{
		ID:            id,
		Owner:         owner.String(),
		Label:         label,
		FeeAddress:    actionAddress.String(),
		Msgs:          msgs,
		Interval:      interval,
		StartTime:     startAt,
		ExecTime:      execTime,
		EndTime:       endTime,
		ICAConfig:     &icaConfig,
		Configuration: &configuration,
		HostedConfig:  &hostedConfig,
		Conditions:    &conditions,
	}

	if !action.ActionAuthzSignerOk(k.cdc) {
		return errorsmod.Wrapf(types.ErrAuthzSigner, "action id: %v", id)
	}
	k.SetActionInfo(ctx, &action)
	k.addToActionOwnerIndex(ctx, owner, id)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeAction,
			sdk.NewAttribute(types.AttributeKeyActionID, strconv.FormatUint(id, 10)),
		))
	return nil
}

func (k Keeper) calculateTimeAndInsertQueue(ctx sdk.Context, startTime time.Time, duration time.Duration, actionID uint64, interval time.Duration) (time.Time, time.Time) {
	endTime, execTime := calculateEndAndExecTimes(ctx, startTime, duration, interval)
	k.InsertActionQueue(ctx, actionID, execTime)

	return endTime, execTime
}

func calculateEndAndExecTimes(ctx sdk.Context, startTime time.Time, duration time.Duration, interval time.Duration) (time.Time, time.Time) {
	endTime := startTime.Add(duration)
	execTime := calculateExecTime(ctx, duration, interval, startTime)

	return endTime, execTime
}

func calculateExecTime(ctx sdk.Context, duration, interval time.Duration, startTime time.Time) time.Time {
	if startTime.After(ctx.BlockTime()) {
		return startTime
	}
	if interval != 0 {
		return startTime.Add(interval)
	}
	return startTime.Add(duration)

}

// peekAutoIncrementID reads the current value without incrementing it.
func (k Keeper) peekAutoIncrementID(ctx sdk.Context, lastIDKey []byte) uint64 {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	bz := store.Get(lastIDKey)
	id := uint64(1)
	if bz != nil {
		id = binary.BigEndian.Uint64(bz)
	}
	return id
}

func (k Keeper) importAutoIncrementID(ctx sdk.Context, lastIDKey []byte, val uint64) error {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	if store.Has(lastIDKey) {
		return errorsmod.Wrapf(types.ErrDuplicate, "autoincrement id: %s", string(lastIDKey))
	}
	bz := sdk.Uint64ToBigEndian(val)
	store.Set(lastIDKey, bz)
	return nil
}

func (k Keeper) importActionInfo(ctx sdk.Context, actionId uint64, action types.ActionInfo) error {

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetActionKey(actionId)
	if store.Has(key) {
		return errorsmod.Wrapf(types.ErrDuplicate, "duplicate code: %d", actionId)
	}
	// 0x01 | actionId (uint64) -> action
	store.Set(key, k.cdc.MustMarshal(&action))
	return nil
}

func (k Keeper) IterateActionInfos(ctx sdk.Context, cb func(uint64, types.ActionInfo) bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	prefixStore := prefix.NewStore(store, types.ActionKeyPrefix)

	iter := prefixStore.Iterator(nil, nil)
	for ; iter.Valid(); iter.Next() {
		var c types.ActionInfo
		k.cdc.MustUnmarshal(iter.Value(), &c)
		// cb returns true to stop early
		if cb(binary.BigEndian.Uint64(iter.Key()), c) {
			return
		}
	}
}

// addToActionOwnerIndex adds element to the index for actions-by-creator queries
func (k Keeper) addToActionOwnerIndex(ctx sdk.Context, ownerAddress sdk.AccAddress, actionID uint64) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	store.Set(types.GetActionByOwnerIndexKey(ownerAddress, actionID), []byte{})
}

// changeActionOwnerIndex changes element to the index for actions-by-creator queries
// func (k Keeper) changeActionOwnerIndex(ctx sdk.Context, ownerAddress, newOwnerAddress sdk.AccAddress, actionID uint64) {
// 	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

// 	store.Set(types.GetActionByOwnerIndexKey(newOwnerAddress, actionID), []byte{})
// 	store.Delete(types.GetActionByOwnerIndexKey(ownerAddress, actionID))
// }

// IterateActionsByOwner iterates over all actions with given creator address in order of creation time asc.
func (k Keeper) IterateActionsByOwner(ctx sdk.Context, owner sdk.AccAddress, cb func(address sdk.AccAddress) bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	prefixStore := prefix.NewStore(store, types.GetActionsByOwnerPrefix(owner))
	for iter := prefixStore.Iterator(nil, nil); iter.Valid(); iter.Next() {
		key := iter.Key()
		if cb(key) {
			return
		}
	}
}

// getTmpActionID getds tmp ActionId for a certain port and sequence. This is used to set results and timeouts.
func (k Keeper) getTmpActionID(ctx sdk.Context, portID string, channelID string, seq uint64) uint64 {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	// Append both portID and channelID to the key
	key := append(types.TmpActionIDLatestTX, []byte(portID)...)
	key = append(key, []byte(channelID)...)          // Append channelID after portID
	key = append(key, types.GetBytesForUint(seq)...) // Append sequence number

	actionIDBz := store.Get(key)

	return types.GetIDFromBytes(actionIDBz)
}

func (k Keeper) setTmpActionID(ctx sdk.Context, actionID uint64, portID string, channelID string, seq uint64) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	// Append both portID and channelID to the key
	key := append(types.TmpActionIDLatestTX, []byte(portID)...)
	key = append(key, []byte(channelID)...)          // Append channelID after portID
	key = append(key, types.GetBytesForUint(seq)...) // Append sequence number

	store.Set(key, types.GetBytesForUint(actionID))
}
