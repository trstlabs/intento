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

// GetFlowInfo
func (k Keeper) GetFlowInfo(ctx sdk.Context, flowID uint64) types.FlowInfo {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	var flow types.FlowInfo
	flowBz := store.Get(types.GetFlowKey(flowID))

	k.cdc.MustUnmarshal(flowBz, &flow)
	return flow
}

// TryGetFlowInfo
func (k Keeper) TryGetFlowInfo(ctx sdk.Context, flowID uint64) (types.FlowInfo, error) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	var flow types.FlowInfo
	flowBz := store.Get(types.GetFlowKey(flowID))
	if flowBz == nil {
		return types.FlowInfo{}, errorsmod.Wrapf(types.ErrNotFound, "flow")
	}
	err := k.cdc.Unmarshal(flowBz, &flow)
	if err != nil {
		return types.FlowInfo{}, err
	}

	return flow, nil
}

func (k Keeper) SetFlowInfo(ctx sdk.Context, flow *types.FlowInfo) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Set(types.GetFlowKey(flow.ID), k.cdc.MustMarshal(flow))
}

func (k Keeper) CreateFlow(ctx sdk.Context, owner sdk.AccAddress, label string, msgs []*cdctypes.Any, duration time.Duration, interval time.Duration, startAt time.Time, feeFunds sdk.Coins, configuration types.ExecutionConfiguration, HostedICAConfig types.HostedICAConfig, portID string, connectionId string, conditions types.ExecutionConditions) error {

	id := k.autoIncrementID(ctx, types.KeyLastID)
	flowAddress, err := k.createFeeAccount(ctx, id, owner, feeFunds)
	if err != nil {
		return err
	}

	endTime, execTime := k.calculateTimeAndInsertQueue(ctx, startAt, duration, id, interval)

	icaConfig := types.ICAConfig{
		PortID:       portID,
		ConnectionID: connectionId,
	}

	flow := types.FlowInfo{
		ID:              id,
		Owner:           owner.String(),
		Label:           label,
		FeeAddress:      flowAddress.String(),
		Msgs:            msgs,
		Interval:        interval,
		StartTime:       startAt,
		ExecTime:        execTime,
		EndTime:         endTime,
		ICAConfig:       &icaConfig,
		Configuration:   &configuration,
		HostedICAConfig: &HostedICAConfig,
		Conditions:      &conditions,
	}

	if err := k.SignerOk(ctx, k.cdc, flow); err != nil {
		return errorsmod.Wrap(types.ErrSignerNotOk, err.Error())
	}

	k.SetFlowInfo(ctx, &flow)
	k.addToFlowOwnerIndex(ctx, owner, id)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeFlowCreated,
			sdk.NewAttribute(types.AttributeKeyFlowID, strconv.FormatUint(id, 10)),
		))
	return nil
}

func (k Keeper) calculateTimeAndInsertQueue(ctx sdk.Context, startTime time.Time, duration time.Duration, flowID uint64, interval time.Duration) (time.Time, time.Time) {
	endTime, execTime := calculateEndAndExecTimes(ctx, startTime, duration, interval)
	k.InsertFlowQueue(ctx, flowID, execTime)

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

func (k Keeper) importFlowInfo(ctx sdk.Context, flowId uint64, flow types.FlowInfo) error {

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetFlowKey(flowId)
	if store.Has(key) {
		return errorsmod.Wrapf(types.ErrDuplicate, "duplicate code: %d", flowId)
	}
	// 0x01 | flowId (uint64) -> flow
	store.Set(key, k.cdc.MustMarshal(&flow))
	return nil
}

func (k Keeper) IterateFlowInfos(ctx sdk.Context, cb func(uint64, types.FlowInfo) bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	prefixStore := prefix.NewStore(store, types.FlowKeyPrefix)

	iter := prefixStore.Iterator(nil, nil)
	for ; iter.Valid(); iter.Next() {
		var c types.FlowInfo
		k.cdc.MustUnmarshal(iter.Value(), &c)
		// cb returns true to stop early
		if cb(binary.BigEndian.Uint64(iter.Key()), c) {
			return
		}
	}
}

// addToFlowOwnerIndex adds element to the index for flows-by-creator queries
func (k Keeper) addToFlowOwnerIndex(ctx sdk.Context, ownerAddress sdk.AccAddress, flowID uint64) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	store.Set(types.GetFlowByOwnerIndexKey(ownerAddress, flowID), []byte{})
}

// changeFlowOwnerIndex changes element to the index for flows-by-creator queries
// func (k Keeper) changeFlowOwnerIndex(ctx sdk.Context, ownerAddress, newOwnerAddress sdk.AccAddress, flowID uint64) {
// 	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

// 	store.Set(types.GetFlowByOwnerIndexKey(newOwnerAddress, flowID), []byte{})
// 	store.Delete(types.GetFlowByOwnerIndexKey(ownerAddress, flowID))
// }

// IterateFlowsByOwner iterates over all flows with given creator address in order of creation time asc.
func (k Keeper) IterateFlowsByOwner(ctx sdk.Context, owner sdk.AccAddress, cb func(address sdk.AccAddress) bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	prefixStore := prefix.NewStore(store, types.GetFlowsByOwnerPrefix(owner))
	for iter := prefixStore.Iterator(nil, nil); iter.Valid(); iter.Next() {
		key := iter.Key()
		if cb(key) {
			return
		}
	}
}

// getTmpFlowID getds tmp FlowID for a certain port and sequence. This is used to set results and timeouts.
func (k Keeper) getTmpFlowID(ctx sdk.Context, portID string, channelID string, seq uint64) uint64 {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	// Append both portID and channelID to the key
	key := append(types.TmpFlowIDLatestTX, []byte(portID)...)
	key = append(key, []byte(channelID)...)          // Append channelID after portID
	key = append(key, types.GetBytesForUint(seq)...) // Append sequence number

	flowIDBz := store.Get(key)

	return types.GetIDFromBytes(flowIDBz)
}

func (k Keeper) setTmpFlowID(ctx sdk.Context, flowID uint64, portID string, channelID string, seq uint64) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	// Append both portID and channelID to the key
	key := append(types.TmpFlowIDLatestTX, []byte(portID)...)
	key = append(key, []byte(channelID)...)          // Append channelID after portID
	key = append(key, types.GetBytesForUint(seq)...) // Append sequence number

	store.Set(key, types.GetBytesForUint(flowID))
}
