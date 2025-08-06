package keeper

import (
	"encoding/binary"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/trstlabs/intento/x/intent/types"
)

// GetTrustlessAgent
func (k Keeper) GetTrustlessAgent(ctx sdk.Context, address string) types.TrustlessAgent {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	var trustlessExecutionAgent types.TrustlessAgent
	trustlessExecutionAgentBz := store.Get(types.GetTrustlessAgentKey(address))

	k.cdc.MustUnmarshal(trustlessExecutionAgentBz, &trustlessExecutionAgent)
	return trustlessExecutionAgent
}

// TryGetTrustlessAgent
func (k Keeper) TryGetTrustlessAgent(ctx sdk.Context, address string) (types.TrustlessAgent, error) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	var trustlessExecutionAgent types.TrustlessAgent
	trustlessExecutionAgentBz := store.Get(types.GetTrustlessAgentKey(address))

	err := k.cdc.Unmarshal(trustlessExecutionAgentBz, &trustlessExecutionAgent)
	if err != nil {
		return types.TrustlessAgent{}, err
	}
	return trustlessExecutionAgent, nil
}

func (k Keeper) SetTrustlessAgent(ctx sdk.Context, trustlessExecutionAgent *types.TrustlessAgent) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Set(types.GetTrustlessAgentKey(trustlessExecutionAgent.AgentAddress), k.cdc.MustMarshal(trustlessExecutionAgent))
}

// func (k Keeper) importTrustlessAgent(ctx sdk.Context, address string, trustlessExecutionAgent types.TrustlessAgent) error {

// 	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
// 	key := types.GetTrustlessAgentKey(address)
// 	if store.Has(key) {
// 		return errorsmod.Wrapf(types.ErrDuplicate, "duplicate address: %s", address)
// 	}
// 	// 0x01 | address (uint64) -> trustlessExecutionAgent
// 	store.Set(key, k.cdc.MustMarshal(&trustlessExecutionAgent))
// 	return nil
// }

func (k Keeper) IterateTrustlessAgents(ctx sdk.Context, cb func(uint64, types.TrustlessAgent) bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	prefixStore := prefix.NewStore(store, types.TrustlessAgentKeyPrefix)
	iter := prefixStore.Iterator(nil, nil)
	for ; iter.Valid(); iter.Next() {
		var c types.TrustlessAgent
		k.cdc.MustUnmarshal(iter.Value(), &c)
		// cb returns true to stop early
		if cb(binary.BigEndian.Uint64(iter.Key()), c) {
			return
		}
	}
}

// addToTrustlessAgentAdminIndex adds element to the index for trustlessExecutionAgents-by-creator queries
func (k Keeper) addToTrustlessAgentAdminIndex(ctx sdk.Context, ownerAddress sdk.AccAddress, trustlessExecutionAgentAddress string) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Set(types.GetTrustlessAgentsByAdminIndexKey(ownerAddress, trustlessExecutionAgentAddress), []byte{})
}

// changeTrustlessAgentAdminIndex changes element to the index for trustlessExecutionAgents-by-creator queries
func (k Keeper) changeTrustlessAgentAdminIndex(ctx sdk.Context, ownerAddress, newAdminAddress sdk.AccAddress, trustlessExecutionAgentAddress string) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	store.Set(types.GetTrustlessAgentsByAdminIndexKey(newAdminAddress, trustlessExecutionAgentAddress), []byte{})
	store.Delete(types.GetTrustlessAgentsByAdminIndexKey(ownerAddress, trustlessExecutionAgentAddress))
}

// IterateTrustlessAgentsByAdmin iterates over all trustlessExecutionAgents with given creator address in order of creation time asc.
func (k Keeper) IterateTrustlessAgentsByAdmin(ctx sdk.Context, owner sdk.AccAddress, cb func(address sdk.AccAddress) bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	prefixStore := prefix.NewStore(store, types.GetTrustlessAgentsByAdminPrefix(owner))

	for iter := prefixStore.Iterator(nil, nil); iter.Valid(); iter.Next() {
		key := iter.Key()
		if cb(key) {
			return
		}
	}
}
