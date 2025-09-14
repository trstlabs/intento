package keeper

import (
	"time"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/trstlabs/intento/x/claim/types"
)

// InsertEntriesIntoVestingQueue inserts up to 4 vesting entries
func (k Keeper) InsertEntriesIntoVestingQueue(ctx sdk.Context, recipientAddr string, action byte, now time.Time) error {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	p, err := k.GetParams(ctx)
	if err != nil {
		return err
	}
	vestDuration := p.DurationVestingPeriods[action]

	timeElapsed := ctx.BlockTime().Sub(p.AirdropStartTime)
	timeLeft := (p.DurationUntilDecay + p.DurationOfDecay) - timeElapsed

	for i := 0; i < 4; i++ {
		if vestDuration*time.Duration(i+1) > timeLeft {
			break
		}
		endTime := now.Add(vestDuration * time.Duration(i+1))
		key := types.VestingQueueKey(recipientAddr, endTime, action, byte(i))
		store.Set(key, []byte{1})
	}
	return nil
}

// VestingQueueIterator returns iterator for all entries up to execTime
func (k Keeper) VestingQueueIterator(ctx sdk.Context, execTime time.Time) storetypes.Iterator {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	return store.Iterator(
		types.VestingQueuePrefix,
		storetypes.PrefixEndBytes(types.VestingByTimeKey(execTime)),
	)
}

// IterateVestingQueue runs cb for each due vesting entry
func (k Keeper) IterateVestingQueue(ctx sdk.Context, execTime time.Time, cb func(addr sdk.AccAddress, action int32, period int32, endTime time.Time) bool) {
	iter := k.VestingQueueIterator(ctx, execTime)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		addrStr, action, period, endTime := types.SplitVestingQueueKey(iter.Key())
		addr, err := sdk.AccAddressFromBech32(addrStr)
		if err != nil {
			panic(err)
		}
		if cb(addr, action, period, endTime) {
			break
		}
	}
}

// RemoveEntryFromVestingQueue deletes a vesting entry
func (k Keeper) RemoveEntryFromVestingQueue(ctx sdk.Context, addr string, endTime time.Time, action byte, period byte) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.VestingQueueKey(addr, endTime, action, period)
	store.Delete(key)
}
