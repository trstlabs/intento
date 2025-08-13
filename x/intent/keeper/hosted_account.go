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
	var trustlessAgent types.TrustlessAgent
	trustlessAgentBz := store.Get(types.GetTrustlessAgentKey(address))

	k.cdc.MustUnmarshal(trustlessAgentBz, &trustlessAgent)
	return trustlessAgent
}

// TryGetTrustlessAgent
func (k Keeper) TryGetTrustlessAgent(ctx sdk.Context, address string) (types.TrustlessAgent, error) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	var trustlessAgent types.TrustlessAgent
	trustlessAgentBz := store.Get(types.GetTrustlessAgentKey(address))

	err := k.cdc.Unmarshal(trustlessAgentBz, &trustlessAgent)
	if err != nil {
		return types.TrustlessAgent{}, err
	}
	return trustlessAgent, nil
}

func (k Keeper) SetTrustlessAgent(ctx sdk.Context, trustlessAgent *types.TrustlessAgent) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Set(types.GetTrustlessAgentKey(trustlessAgent.AgentAddress), k.cdc.MustMarshal(trustlessAgent))
}

// func (k Keeper) importTrustlessAgent(ctx sdk.Context, address string, trustlessAgent types.TrustlessAgent) error {

// 	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
// 	key := types.GetTrustlessAgentKey(address)
// 	if store.Has(key) {
// 		return errorsmod.Wrapf(types.ErrDuplicate, "duplicate address: %s", address)
// 	}
// 	// 0x01 | address (uint64) -> trustlessAgent
// 	store.Set(key, k.cdc.MustMarshal(&trustlessAgent))
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

// addToTrustlessAgentAdminIndex adds element to the index for trustlessAgents-by-creator queries
func (k Keeper) addToTrustlessAgentAdminIndex(ctx sdk.Context, ownerAddress sdk.AccAddress, trustlessAgentAddress string) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Set(types.GetTrustlessAgentsByAdminIndexKey(ownerAddress, trustlessAgentAddress), []byte{})
}

// changeTrustlessAgentAdminIndex changes element to the index for trustlessAgents-by-creator queries
func (k Keeper) changeTrustlessAgentAdminIndex(ctx sdk.Context, ownerAddress, newAdminAddress sdk.AccAddress, trustlessAgentAddress string) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	store.Set(types.GetTrustlessAgentsByAdminIndexKey(newAdminAddress, trustlessAgentAddress), []byte{})
	store.Delete(types.GetTrustlessAgentsByAdminIndexKey(ownerAddress, trustlessAgentAddress))
}

// IterateTrustlessAgentsByAdmin iterates over all trustlessAgents with given creator address in order of creation time asc.
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
