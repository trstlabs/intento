package keeper

import (
	"encoding/binary"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/trstlabs/intento/x/intent/types"
)

// GetHostedAccount
func (k Keeper) GetHostedAccount(ctx sdk.Context, address string) types.HostedAccount {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	var hostedAccount types.HostedAccount
	hostedAccountBz := store.Get(types.GetHostedAccountKey(address))

	k.cdc.MustUnmarshal(hostedAccountBz, &hostedAccount)
	return hostedAccount
}

// TryGetHostedAccount
func (k Keeper) TryGetHostedAccount(ctx sdk.Context, address string) (types.HostedAccount, error) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	var hostedAccount types.HostedAccount
	hostedAccountBz := store.Get(types.GetHostedAccountKey(address))

	err := k.cdc.Unmarshal(hostedAccountBz, &hostedAccount)
	if err != nil {
		return types.HostedAccount{}, err
	}
	return hostedAccount, nil
}

func (k Keeper) SetHostedAccount(ctx sdk.Context, hostedAccount *types.HostedAccount) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Set(types.GetHostedAccountKey(hostedAccount.HostedAddress), k.cdc.MustMarshal(hostedAccount))
}

// func (k Keeper) importHostedAccount(ctx sdk.Context, address string, hostedAccount types.HostedAccount) error {

// 	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
// 	key := types.GetHostedAccountKey(address)
// 	if store.Has(key) {
// 		return errorsmod.Wrapf(types.ErrDuplicate, "duplicate address: %s", address)
// 	}
// 	// 0x01 | address (uint64) -> hostedAccount
// 	store.Set(key, k.cdc.MustMarshal(&hostedAccount))
// 	return nil
// }

func (k Keeper) IterateHostedAccounts(ctx sdk.Context, cb func(uint64, types.HostedAccount) bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	prefixStore := prefix.NewStore(store, types.HostedAccountKeyPrefix)
	iter := prefixStore.Iterator(nil, nil)
	for ; iter.Valid(); iter.Next() {
		var c types.HostedAccount
		k.cdc.MustUnmarshal(iter.Value(), &c)
		// cb returns true to stop early
		if cb(binary.BigEndian.Uint64(iter.Key()), c) {
			return
		}
	}
}

// addToHostedAccountAdminIndex adds element to the index for hostedAccounts-by-creator queries
func (k Keeper) addToHostedAccountAdminIndex(ctx sdk.Context, ownerAddress sdk.AccAddress, hostedAccountAddress string) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Set(types.GetHostedAccountsByAdminIndexKey(ownerAddress, hostedAccountAddress), []byte{})
}

// changeHostedAccountAdminIndex changes element to the index for hostedAccounts-by-creator queries
func (k Keeper) changeHostedAccountAdminIndex(ctx sdk.Context, ownerAddress, newAdminAddress sdk.AccAddress, hostedAccountAddress string) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	store.Set(types.GetHostedAccountsByAdminIndexKey(newAdminAddress, hostedAccountAddress), []byte{})
	store.Delete(types.GetHostedAccountsByAdminIndexKey(ownerAddress, hostedAccountAddress))
}

// IterateHostedAccountsByAdmin iterates over all hostedAccounts with given creator address in order of creation time asc.
func (k Keeper) IterateHostedAccountsByAdmin(ctx sdk.Context, owner sdk.AccAddress, cb func(address sdk.AccAddress) bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	prefixStore := prefix.NewStore(store, types.GetHostedAccountsByAdminPrefix(owner))

	for iter := prefixStore.Iterator(nil, nil); iter.Valid(); iter.Next() {
		key := iter.Key()
		if cb(key) {
			return
		}
	}
}
