package keeper

import (
	"fmt"
	"time"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/trstlabs/intento/x/claim/types"
)

func (k Keeper) IterateVestingQueue(ctx sdk.Context, endTime time.Time, cb func(addr sdk.AccAddress, action int32, period int32, endTime time.Time) (stop bool)) {
	iterator := k.VestingQueueIterator(ctx, endTime)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {

		addr, endTime := types.SplitVestingQueueKey(iterator.Key())

		recipientAddr, err := sdk.AccAddressFromBech32(addr)
		if err != nil {
			return
		}
		action := int32(iterator.Value()[0])
		period := int32(iterator.Value()[1])

		if cb(recipientAddr, action, period, endTime) {
			break
		}
	}
}

// VestingQueueIterator returns an sdk.Iterator for all the vesting periods in the Queue that expire by endTime
func (k Keeper) VestingQueueIterator(ctx sdk.Context, endTime time.Time) storetypes.Iterator {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	return store.Iterator(types.VestingQueuePrefix, storetypes.PrefixEndBytes(types.VestingByTimeKey(endTime))) //we check the end of the byte array for the end time
}

// InsertVestingQueue Inserts a contract into the vesting queue at endTime
func (k Keeper) InsertEntriesIntoVestingQueue(ctx sdk.Context, recipientAddr string, action byte, timeNow time.Time) error {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	p, err := k.GetParams(ctx)
	if err != nil {
		return err
	}
	//duration of 1 vesting entry for the given action
	vestDuration := p.DurationVestingPeriods[action]

	timeElapsed := ctx.BlockTime().Sub(p.AirdropStartTime)
	timeLeft := (p.DurationUntilDecay + p.DurationOfDecay) - timeElapsed
	//	fmt.Printf("timeLeft %v \n", timeLeft)
	for i := 0; i < 4; i++ {
		//	fmt.Printf("period %v \n", i)
		//exclude if vestduration is longer than timeLeft
		if vestDuration*time.Duration(i+1) > timeLeft {
			fmt.Printf("break")
			break
		}
		//fmt.Printf("duration %v \n", vestDuration*time.Duration(i+1))
		store.Set(types.VestingQueueKey(recipientAddr, timeNow.Add(vestDuration*time.Duration(i+1))), []byte{action, byte(i)})
		//		fmt.Printf("set %v \n", []byte{action, byte(i)})
	}
	return nil

}

// RemoveEntryFromVestingQueue removes a period from the vesting Queue
func (k Keeper) RemoveEntryFromVestingQueue(ctx sdk.Context, recipientAddr string, endTime time.Time) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	//fmt.Printf("removed \n")
	store.Delete(types.VestingQueueKey(recipientAddr, endTime))
}
