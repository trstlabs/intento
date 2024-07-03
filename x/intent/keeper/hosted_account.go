package keeper

import (
	"encoding/binary"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/trstlabs/intento/x/intent/types"
)

// GetHostedAccount
func (k Keeper) GetHostedAccount(ctx sdk.Context, address string) types.HostedAccount {
	store := ctx.KVStore(k.storeKey)
	var action types.HostedAccount
	actionBz := store.Get(types.GetHostedAccountKey(address))

	k.cdc.MustUnmarshal(actionBz, &action)
	return action
}

// TryGetHostedAccount
func (k Keeper) TryGetHostedAccount(ctx sdk.Context, address string) (types.HostedAccount, error) {
	store := ctx.KVStore(k.storeKey)
	var action types.HostedAccount
	actionBz := store.Get(types.GetHostedAccountKey(address))

	err := k.cdc.Unmarshal(actionBz, &action)
	if err != nil {
		return types.HostedAccount{}, err
	}
	return action, nil
}

func (k Keeper) SetHostedAccount(ctx sdk.Context, hostedAccount *types.HostedAccount) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetHostedAccountKey(hostedAccount.HostedAddress), k.cdc.MustMarshal(hostedAccount))
}

// func (k Keeper) importHostedAccount(ctx sdk.Context, address string, action types.HostedAccount) error {

// 	store := ctx.KVStore(k.storeKey)
// 	key := types.GetHostedAccountKey(address)
// 	if store.Has(key) {
// 		return errorsmod.Wrapf(types.ErrDuplicate, "duplicate address: %s", address)
// 	}
// 	// 0x01 | address (uint64) -> action
// 	store.Set(key, k.cdc.MustMarshal(&action))
// 	return nil
// }

func (k Keeper) IterateHostedAccounts(ctx sdk.Context, cb func(uint64, types.HostedAccount) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.HostedAccountKeyPrefix)
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
