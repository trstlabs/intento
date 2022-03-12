package keeper

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/trstlabs/trst/x/claim/types"
)

// IterateVestingsQueue iterates over the claims in the succesfull claims queue
// and performs a callback function
func (k Keeper) IterateVestingQueue(ctx sdk.Context, endTime time.Time, cb func(coins sdk.Coins) (stop bool)) {
	iterator := k.VestingQueueIterator(ctx, endTime)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {

		addr, _ := types.SplitVestingQueueKey(iterator.Key())
		claimerAddr, err := sdk.AccAddressFromBech32(addr)
		if err != nil {
			return
		}

		stringCoins := string(iterator.Value())
		coins, err := sdk.ParseCoinsNormalized(stringCoins)
		if err != nil {
			fmt.Printf("coins string %s", stringCoins)
			panic("Failed to parse coins")
		}
		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, claimerAddr, coins)
		if err != nil {
			return
		}
		k.RemoveFromVestingQueue(ctx, string(claimerAddr), endTime)
		if cb(coins) {
			break
		}
	}
}

// VestingQueueIterator returns an sdk.Iterator for all the items in the Inactive Queue that expire by endTime
func (k Keeper) VestingQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return store.Iterator(types.VestingQueuePrefix, sdk.PrefixEndBytes(types.VestingByTimeKey(endTime))) //we check the end of the bites array for the end time
}

/*
// InsertVestingQueue Inserts a contract into the inactive item queue at endTime
func (k Keeper) InsertVestingQueue(ctx sdk.Context, claimableAmount sdk.Coins, contractAddr string, endTime time.Time) {
	store := ctx.KVStore(k.storeKey)

	bz := []byte(claimableAmount.String())
	//here the key is time+contract appended (as bytes) and value is contract in bytes
	store.Set(types.VestingQueueKey(contractAddr, endTime), bz)
}*/

// InsertVestingQueue Inserts a contract into the inactive item queue at endTime
func (k Keeper) InsertEntriesIntoVestingQueue(ctx sdk.Context, claimableAmount sdk.Coins, contractAddr string, duration time.Duration, time time.Time) {
	store := ctx.KVStore(k.storeKey)
	bz := []byte(claimableAmount.String())
	//here the key is time+contract appended (as bytes) and value is contract in bytes
	store.Set(types.VestingQueueKey(contractAddr, time.Add(duration)), bz)
	store.Set(types.VestingQueueKey(contractAddr, time.Add(duration*2)), bz)
	store.Set(types.VestingQueueKey(contractAddr, time.Add(duration*3)), bz)
	store.Set(types.VestingQueueKey(contractAddr, time.Add(duration*4)), bz)
}

// RemoveFromVestingQueue removes a contract from the Inactive Item Queue
func (k Keeper) RemoveFromVestingQueue(ctx sdk.Context, contractAddr string, endTime time.Time) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.VestingQueueKey(contractAddr, endTime))
}
