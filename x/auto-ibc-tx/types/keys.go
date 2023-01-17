package types

import (
	"encoding/binary"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

const (
	ModuleName = "icamsgauth"

	StoreKey = ModuleName

	RouterKey = ModuleName

	QuerierRoute = ModuleName
)

// nolint
var (
	AutoTxKeyPrefix      = []byte{0x01}
	AutoTxStorePrefix    = []byte{0x02}
	AutoTxQueuePrefix    = []byte{0x03}
	SequenceKeyPrefix    = []byte{0x04}
	AutoTxsByOwnerPrefix = []byte{0x05}
	KeyLastTxID          = append(SequenceKeyPrefix, []byte("lastTxId")...)
	KeyLastTxAddrID      = append(SequenceKeyPrefix, []byte("lastTxAddrId")...)

/* 	TimeTimeLen          = 24 */
)

// GetAutoTxKey returns the key for the auto interchain tx
func GetAutoTxKey(autoTxID uint64) []byte {
	return append(AutoTxKeyPrefix, GetBytesForUint(autoTxID)...)
}

// GetAutoTxsByOwnerPrefix returns the autoTxs by creator prefix for the WASM autoTx instance
func GetAutoTxsByOwnerPrefix(addr sdk.AccAddress) []byte {
	bz := address.MustLengthPrefix(addr)
	return append(AutoTxsByOwnerPrefix, bz...)
}

////queue types

var lenTime = len(sdk.FormatTimeBytes(time.Now()))

// SplitAutoTxQueueKey split the listed key and returns the id and endTime
func SplitAutoTxQueueKey(key []byte) (autoTxID uint64, endTime time.Time) {
	return splitKeyWithTime(key)
}

// AutoTxByTimeKey gets the listed item queue key by endTime
func AutoTxByTimeKey(endTime time.Time) []byte {
	return append(AutoTxQueuePrefix, sdk.FormatTimeBytes(endTime)...)
}

// from the key we get the autoTx and end time
func splitKeyWithTime(key []byte) (autoTxID uint64, endTime time.Time) {

	endTime, _ = sdk.ParseTimeBytes(key[1 : 1+lenTime])

	//returns an id from bytes
	autoTxID = binary.BigEndian.Uint64(key[1+lenTime:])

	return
}

// AutoTxQueueKey returns the key with prefix for an autoTx in the Listed Item Queue
func AutoTxQueueKey(autoTxID uint64, endTime time.Time) []byte {
	return append(AutoTxByTimeKey(endTime), GetBytesForUint(autoTxID)...)
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

// GetAutoTxByOwnerSecondaryIndexKey returns the key for the second index: `<prefix><ownerAddress length><ownerAddress><autoTxID>`
func GetAutoTxByOwnerSecondaryIndexKey(bz []byte, autoTxID uint64) []byte {
	prefixBytes := GetAutoTxsByOwnerPrefix(bz)
	lenPrefixBytes := len(prefixBytes)
	r := make([]byte, lenPrefixBytes+8)

	copy(r[:lenPrefixBytes], prefixBytes)
	copy(r[lenPrefixBytes:], GetBytesForUint(autoTxID))

	return r
}
