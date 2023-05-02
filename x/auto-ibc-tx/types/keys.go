package types

import (
	"encoding/binary"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

const (
	ModuleName = "autoibctx"

	StoreKey = ModuleName

	RouterKey = ModuleName

	QuerierRoute = ModuleName
)

// nolint
var (
	AutoTxKeyPrefix               = []byte{0x01}
	AutoTxStorePrefix             = []byte{0x02}
	AutoTxQueuePrefix             = []byte{0x03}
	SequenceKeyPrefix             = []byte{0x04}
	AutoTxsByOwnerPrefix          = []byte{0x05}
	TmpAutoTxIDLatestTX           = []byte{0x06}
	KeyRelayerRewardsAvailability = []byte{0x07}
	AutoTxIbcUsageKeyPrefix       = []byte{0x08}
	KeyLastTxID                   = append(SequenceKeyPrefix, []byte("lastTxId")...)
	KeyLastTxAddrID               = append(SequenceKeyPrefix, []byte("lastTxAddrId")...)
)

// ics 20 hook
var SenderPrefix = "ibc-auto-tx-hook-intermediary"

var (
	KeyAutoTxIncentiveForSDKTx   = 0
	KeyAutoTxIncentiveForWasmTx  = 1
	KeyAutoTxIncentiveForOsmoTx  = 2
	KeyAutoTxIncentiveForAuthzTx = 3
)

// GetAutoTxKey returns the key for the auto interchain tx
func GetAutoTxKey(autoTxID uint64) []byte {
	return append(AutoTxKeyPrefix, GetBytesForUint(autoTxID)...)
}

// GetAutoTxsByOwnerPrefix returns the autoTxs by creator prefix
func GetAutoTxsByOwnerPrefix(addr sdk.AccAddress) []byte {
	bz := address.MustLengthPrefix(addr)
	return append(AutoTxsByOwnerPrefix, bz...)
}

////queue types

var lenTime = len(sdk.FormatTimeBytes(time.Now()))

// SplitAutoTxQueueKey split the listed key and returns the id and execTime
func SplitAutoTxQueueKey(key []byte) (autoTxID uint64, execTime time.Time) {
	return splitKeyWithTime(key)
}

// AutoTxByTimeKey gets the listed item queue key by execTime
func AutoTxByTimeKey(execTime time.Time) []byte {
	return append(AutoTxQueuePrefix, sdk.FormatTimeBytes(execTime)...)
}

// from the key we get the autoTx and end time
func splitKeyWithTime(key []byte) (autoTxID uint64, execTime time.Time) {

	execTime, _ = sdk.ParseTimeBytes(key[1 : 1+lenTime])

	//returns an id from bytes
	autoTxID = binary.BigEndian.Uint64(key[1+lenTime:])

	return
}

// AutoTxQueueKey returns the key with prefix for an autoTx in the Listed Item Queue
func AutoTxQueueKey(autoTxID uint64, execTime time.Time) []byte {
	return append(AutoTxByTimeKey(execTime), GetBytesForUint(autoTxID)...)
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

// GetAutoTxByOwnerIndexKey returns the id: `<prefix><ownerAddress length><ownerAddress><autoTxID>`
func GetAutoTxByOwnerIndexKey(bz []byte, autoTxID uint64) []byte {
	prefixBytes := GetAutoTxsByOwnerPrefix(bz)
	lenPrefixBytes := len(prefixBytes)
	r := make([]byte, lenPrefixBytes+8)

	copy(r[:lenPrefixBytes], prefixBytes)
	copy(r[lenPrefixBytes:], GetBytesForUint(autoTxID))

	return r
}
