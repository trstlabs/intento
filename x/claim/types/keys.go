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

	// ActionKey defines the store key to store user accomplished actions
	ActionKey = "action"

	MemStoreKey = "mem_claim"
)

// ParamsKey stores the module params
var ParamsKey = []byte{0x01}

var ClaimsPortions = int64(5)

// nolint
var (
	VestingStorePrefix = []byte{0x01}
	VestingQueuePrefix = []byte{0x02}
)

var lenTime = len(sdk.FormatTimeBytes(time.Now()))

// VestingByTimeKey = prefix | time
func VestingByTimeKey(endTime time.Time) []byte {
	return append(VestingQueuePrefix, sdk.FormatTimeBytes(endTime)...)
}

// VestingQueueKey = prefix | time | addr | action | period
func VestingQueueKey(addr string, endTime time.Time, action byte, period byte) []byte {
	key := VestingByTimeKey(endTime)
	key = append(key, []byte(addr)...) // store Bech32 string
	key = append(key, action)
	key = append(key, period)
	return key
}

// SplitVestingQueueKey reverses VestingQueueKey
func SplitVestingQueueKey(key []byte) (addr string, action int32, period int32, endTime time.Time) {
	// strip prefix
	key = key[1:]

	// time part
	endTime, _ = sdk.ParseTimeBytes(key[:len(sdk.FormatTimeBytes(time.Now()))])
	rest := key[len(sdk.FormatTimeBytes(time.Now())):]

	// last two bytes are action + period
	if len(rest) < 2 {
		panic("invalid vesting queue key: too short")
	}
	action = int32(rest[len(rest)-2])
	period = int32(rest[len(rest)-1])

	// everything else is the address string
	addr = string(rest[:len(rest)-2])
	return
}
