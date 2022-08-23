package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "claim"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// ClaimRecordsStorePrefix defines the store prefix for the claim records
	ClaimRecordsStorePrefix = "claimrecords"

	// ParamsKey defines the store key for claim module parameters
	ParamsKey = "params"

	// ActionKey defines the store key to store user accomplished actions
	ActionKey = "action"

	MemStoreKey = "mem_claim"
)

// nolint
var (
	VestingStorePrefix = []byte{0x01}
	VestingQueuePrefix = []byte{0x02}
)

// GetVestingStorePrefixKey returns the store prefix for the WASM contract instance
func GetVestingStorePrefixKey(addr sdk.AccAddress) []byte {
	return append(VestingStorePrefix, addr...)
}

////queue types

var lenTime = len(sdk.FormatTimeBytes(time.Now()))

// SplitVestingQueueKey split the listed key and returns the id and endTime
func SplitVestingQueueKey(key []byte) (vestingAddr string, endTime time.Time) {
	return splitKeyWithTime(key)
}

// VestingByTimeKey gets the listed claim queue key by endTime
func VestingByTimeKey(endTime time.Time) []byte {
	return append(VestingQueuePrefix, sdk.FormatTimeBytes(endTime)...)
}

//from the key we get the claim and end time
func splitKeyWithTime(key []byte) (vestingAddr string, endTime time.Time) {

	/*if len(key[1:]) != 8+lenTime {
		panic(fmt.Sprintf("unexpected key length (%d ≠ %d)", len(key[1:]), lenTime+8))
	}*/

	endTime, _ = sdk.ParseTimeBytes(key[1 : 1+lenTime])
	//fmt.Printf("time %s \n", endTime)
	//returns an id from bytes
	vestingAddr = string(key[1+lenTime:])

	//	fmt.Printf("vestingAddr key is %s ", vestingAddr)
	return
}

// VestingQueueKey returns the key with prefix for an contract in the Listed claim Queue
func VestingQueueKey(vestingAddr string, endTime time.Time) []byte {
	return append(VestingByTimeKey(endTime), []byte(vestingAddr)...)
}
