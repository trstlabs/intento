package keeper

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/trstlabs/trst/x/claim/types"
)

// IterateVestingsQueue iterates over the vesting period entries in the claims queue
// and performs a callback function
func (k Keeper) IterateVestingQueue(ctx sdk.Context, endTime time.Time, cb func(period int32) (stop bool)) {
	iterator := k.VestingQueueIterator(ctx, endTime)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {

		addr, _ := types.SplitVestingQueueKey(iterator.Key())
		recipientAddr, err := sdk.AccAddressFromBech32(addr)
		if err != nil {
			return
		}
		claimRecord, err := k.GetClaimRecord(ctx, recipientAddr)
		if err != nil {
			fmt.Printf("record for %s", addr)
			panic("Failed to get claim record ")
		}
		action := int32(iterator.Value()[0])
		period := int32(iterator.Value()[1])

		claimRecord.Status[action].VestingPeriodCompleted[period] = true
		err = k.SetClaimRecord(ctx, claimRecord)
		if err != nil {
			fmt.Printf("record for %s", addr)
			panic("Failed to set claim record ")

		}
		k.RemoveEntryFromVestingQueue(ctx, string(recipientAddr), endTime)
		if cb(period) {
			break
		}
	}
}

// VestingQueueIterator returns an sdk.Iterator for all the vesting periods in the Queue that expire by endTime
func (k Keeper) VestingQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return store.Iterator(types.VestingQueuePrefix, sdk.PrefixEndBytes(types.VestingByTimeKey(endTime))) //we check the end of the byte array for the end time
}

/*
// InsertVestingQueue Inserts a contract into the inactive vesting queue at endTime
func (k Keeper) InsertVestingQueue(ctx sdk.Context, claimableAmount sdk.Coins, recipientAddr string, endTime time.Time) {
	store := ctx.KVStore(k.storeKey)

	bz := []byte(claimableAmount.String())
	//here the key is time+contract appended (as bytes) and value is contract in bytes
	store.Set(types.VestingQueueKey(recipientAddr, endTime), bz)
}*/

// InsertVestingQueue Inserts a contract into the vesting queue at endTime
func (k Keeper) InsertEntriesIntoVestingQueue(ctx sdk.Context, recipientAddr string, action byte, timeNow time.Time) error {
	store := ctx.KVStore(k.storeKey)
	//duration of 1 vesting entry for the given action
	params, err := k.GetParams(ctx)
	if err != nil {
		return err
	}
	vestDuration := params.DurationVestingPeriods[action]
	timeElapsed := ctx.BlockTime().Sub(params.AirdropStartTime)
	timeLeft := (params.DurationUntilDecay + params.DurationOfDecay) - timeElapsed
	for i := 0; i < 4; i++ {
		fmt.Printf("period %v \n", i)
		//exclude if vestduration is longer than timeLeft
		if vestDuration*time.Duration(i+1) > timeLeft {
			fmt.Printf("break")
			break
		}
		fmt.Printf("duration %v \n", vestDuration*time.Duration(i+1))
		store.Set(types.VestingQueueKey(recipientAddr, timeNow.Add(vestDuration*time.Duration(i+1))), []byte{action, byte(i)})
		fmt.Printf("set %v \n", []byte{action, byte(i)})
	}
	return nil
	/*//here the key is time+recipientAddr appended (as bytes) and value is contract in bytes
	store.Set(types.VestingQueueKey(recipientAddr, timeNow.Add(vestDuration)), []byte{action, 0})
	store.Set(types.VestingQueueKey(recipientAddr, timeNow.Add(vestDuration*2)), []byte{action, 1})
	store.Set(types.VestingQueueKey(recipientAddr, timeNow.Add(vestDuration*3)), []byte{action, 2})
	store.Set(types.VestingQueueKey(recipientAddr, timeNow.Add(vestDuration*4)), []byte{action, 3})*/

}

// RemoveEntryFromVestingQueue removes a period from the vesting Queue
func (k Keeper) RemoveEntryFromVestingQueue(ctx sdk.Context, recipientAddr string, endTime time.Time) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.VestingQueueKey(recipientAddr, endTime))
}
