package keeper

import (
	"encoding/binary"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/trstlabs/intento/x/intent/types"
)

// GetTrustlessExecutionAgent
func (k Keeper) GetTrustlessExecutionAgent(ctx sdk.Context, address string) types.TrustlessExecutionAgent {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	var trustlessExecutionAgent types.TrustlessExecutionAgent
	trustlessExecutionAgentBz := store.Get(types.GetTrustlessExecutionAgentKey(address))

	k.cdc.MustUnmarshal(trustlessExecutionAgentBz, &trustlessExecutionAgent)
	return trustlessExecutionAgent
}

// TryGetTrustlessExecutionAgent
func (k Keeper) TryGetTrustlessExecutionAgent(ctx sdk.Context, address string) (types.TrustlessExecutionAgent, error) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	var trustlessExecutionAgent types.TrustlessExecutionAgent
	trustlessExecutionAgentBz := store.Get(types.GetTrustlessExecutionAgentKey(address))

	err := k.cdc.Unmarshal(trustlessExecutionAgentBz, &trustlessExecutionAgent)
	if err != nil {
		return types.TrustlessExecutionAgent{}, err
	}
	return trustlessExecutionAgent, nil
}

func (k Keeper) SetTrustlessExecutionAgent(ctx sdk.Context, trustlessExecutionAgent *types.TrustlessExecutionAgent) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Set(types.GetTrustlessExecutionAgentKey(trustlessExecutionAgent.AgentAddress), k.cdc.MustMarshal(trustlessExecutionAgent))
}

// func (k Keeper) importTrustlessExecutionAgent(ctx sdk.Context, address string, trustlessExecutionAgent types.TrustlessExecutionAgent) error {

// 	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
// 	key := types.GetTrustlessExecutionAgentKey(address)
// 	if store.Has(key) {
// 		return errorsmod.Wrapf(types.ErrDuplicate, "duplicate address: %s", address)
// 	}
// 	// 0x01 | address (uint64) -> trustlessExecutionAgent
// 	store.Set(key, k.cdc.MustMarshal(&trustlessExecutionAgent))
// 	return nil
// }

func (k Keeper) IterateTrustlessExecutionAgents(ctx sdk.Context, cb func(uint64, types.TrustlessExecutionAgent) bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	prefixStore := prefix.NewStore(store, types.TrustlessExecutionAgentKeyPrefix)
	iter := prefixStore.Iterator(nil, nil)
	for ; iter.Valid(); iter.Next() {
		var c types.TrustlessExecutionAgent
		k.cdc.MustUnmarshal(iter.Value(), &c)
		// cb returns true to stop early
		if cb(binary.BigEndian.Uint64(iter.Key()), c) {
			return
		}
	}
}

// addToTrustlessExecutionAgentAdminIndex adds element to the index for trustlessExecutionAgents-by-creator queries
func (k Keeper) addToTrustlessExecutionAgentAdminIndex(ctx sdk.Context, ownerAddress sdk.AccAddress, trustlessExecutionAgentAddress string) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Set(types.GetTrustlessExecutionAgentsByAdminIndexKey(ownerAddress, trustlessExecutionAgentAddress), []byte{})
}

// changeTrustlessExecutionAgentAdminIndex changes element to the index for trustlessExecutionAgents-by-creator queries
func (k Keeper) changeTrustlessExecutionAgentAdminIndex(ctx sdk.Context, ownerAddress, newAdminAddress sdk.AccAddress, trustlessExecutionAgentAddress string) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	store.Set(types.GetTrustlessExecutionAgentsByAdminIndexKey(newAdminAddress, trustlessExecutionAgentAddress), []byte{})
	store.Delete(types.GetTrustlessExecutionAgentsByAdminIndexKey(ownerAddress, trustlessExecutionAgentAddress))
}

// IterateTrustlessExecutionAgentsByAdmin iterates over all trustlessExecutionAgents with given creator address in order of creation time asc.
func (k Keeper) IterateTrustlessExecutionAgentsByAdmin(ctx sdk.Context, owner sdk.AccAddress, cb func(address sdk.AccAddress) bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	prefixStore := prefix.NewStore(store, types.GetTrustlessExecutionAgentsByAdminPrefix(owner))

	for iter := prefixStore.Iterator(nil, nil); iter.Valid(); iter.Next() {
		key := iter.Key()
		if cb(key) {
			return
		}
	}
}
