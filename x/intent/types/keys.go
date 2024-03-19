package types

import (
	"encoding/binary"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

const (
	ModuleName = "intent"

	StoreKey = ModuleName

	RouterKey = ModuleName

	QuerierRoute = ModuleName
)

// nolint
var (
	ActionKeyPrefix               = []byte{0x01}
	ActionHistoryPrefix           = []byte{0x02}
	ActionQueuePrefix             = []byte{0x03}
	SequenceKeyPrefix             = []byte{0x04}
	ActionsByOwnerPrefix          = []byte{0x05}
	TmpActionIDLatestTX           = []byte{0x06}
	KeyRelayerRewardsAvailability = []byte{0x07}
	ActionIbcUsageKeyPrefix       = []byte{0x08}
	KeyLastID                     = append(SequenceKeyPrefix, []byte("lastID")...)
	KeyLastTxAddrID               = append(SequenceKeyPrefix, []byte("lastTxAddrId")...)
)

// ics 20 hook
var SenderPrefix = "ibc-action-hook-intermediary"

var (
	KeyActionIncentiveForSDKTx   = 0
	KeyActionIncentiveForWasmTx  = 1
	KeyActionIncentiveForOsmoTx  = 2
	KeyActionIncentiveForAuthzTx = 3
)

// GetActionKey returns the key for the auto interchain tx
func GetActionKey(autoID uint64) []byte {
	return append(ActionKeyPrefix, GetBytesForUint(autoID)...)
}

// GetActionHistoryKey returns the key for the auto interchain tx
func GetActionHistoryKey(autoID uint64) []byte {
	return append(ActionHistoryPrefix, GetBytesForUint(autoID)...)
}

// GetActionsByOwnerPrefix returns the autoTxs by creator prefix
func GetActionsByOwnerPrefix(addr sdk.AccAddress) []byte {
	bz := address.MustLengthPrefix(addr)
	return append(ActionsByOwnerPrefix, bz...)
}

////queue types

var lenTime = len(sdk.FormatTimeBytes(time.Now()))

// SplitActionQueueKey split the listed key and returns the id and execTime
func SplitActionQueueKey(key []byte) (autoID uint64, execTime time.Time) {
	return splitKeyWithTime(key)
}

// ActionByTimeKey gets the listed item queue key by execTime
func ActionByTimeKey(execTime time.Time) []byte {
	return append(ActionQueuePrefix, sdk.FormatTimeBytes(execTime)...)
}

// from the key we get the autoTx and end time
func splitKeyWithTime(key []byte) (autoID uint64, execTime time.Time) {

	execTime, _ = sdk.ParseTimeBytes(key[1 : 1+lenTime])

	//returns an id from bytes
	autoID = binary.BigEndian.Uint64(key[1+lenTime:])

	return
}

// ActionQueueKey returns the key with prefix for an autoTx in the Listed Item Queue
func ActionQueueKey(autoID uint64, execTime time.Time) []byte {
	return append(ActionByTimeKey(execTime), GetBytesForUint(autoID)...)
}

// GetBytesForUint returns the byte representation of the itemID
func GetBytesForUint(id uint64) (idBz []byte) {
	idBz = make([]byte, 8)
	binary.BigEndian.PutUint64(idBz, id)
	return
}

/// helper functions

// GetIDFromBytes returns itemID in uint64 format from a byte array
func GetIDFromBytes(bz []byte) (id uint64) {
	return binary.BigEndian.Uint64(bz)
}

// GetActionByOwnerIndexKey returns the id: `<prefix><ownerAddress length><ownerAddress><autoID>`
func GetActionByOwnerIndexKey(bz []byte, autoID uint64) []byte {
	prefixBytes := GetActionsByOwnerPrefix(bz)
	lenPrefixBytes := len(prefixBytes)
	r := make([]byte, lenPrefixBytes+8)

	copy(r[:lenPrefixBytes], prefixBytes)
	copy(r[lenPrefixBytes:], GetBytesForUint(autoID))

	return r
}
