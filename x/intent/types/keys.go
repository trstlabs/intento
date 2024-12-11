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
	// ParamsKey stores the module params
	ParamsKey                     = []byte{0x01}
	ActionKeyPrefix               = []byte{0x02}
	ActionHistoryPrefix           = []byte{0x03}
	ActionQueuePrefix             = []byte{0x04}
	SequenceKeyPrefix             = []byte{0x05}
	ActionsByOwnerPrefix          = []byte{0x06}
	TmpActionIDLatestTX           = []byte{0x07}
	KeyRelayerRewardsAvailability = []byte{0x08}
	ActionIbcUsageKeyPrefix       = []byte{0x09}
	ActionHistorySequencePrefix   = []byte{0x10}
	HostedAccountKeyPrefix        = []byte{0x11}
	HostedAccountsByAdminPrefix   = []byte{0x12}
	KeyLastID                     = append(SequenceKeyPrefix, []byte("lastId")...)
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

// GetActionKey returns the key for the action
func GetActionKey(actionID uint64) []byte {
	return append(ActionKeyPrefix, GetBytesForUint(actionID)...)
}

// GetActionHistoryKey returns the key for the action
func GetActionHistoryKey(actionID uint64) []byte {
	return append(ActionHistoryPrefix, GetBytesForUint(actionID)...)
}

// GetActionsByOwnerPrefix returns the actions by creator prefix
func GetActionsByOwnerPrefix(addr sdk.AccAddress) []byte {
	bz := address.MustLengthPrefix(addr)
	return append(ActionsByOwnerPrefix, bz...)
}

// GetHostedAccountKey returns the key for the hosted account
func GetHostedAccountKey(address string) []byte {
	return append(HostedAccountKeyPrefix, []byte(address)...)
}

// GetHostedAccountsByAdminPrefix returns the actions by creator prefix
func GetHostedAccountsByAdminPrefix(addr sdk.AccAddress) []byte {
	bz := address.MustLengthPrefix(addr)
	return append(HostedAccountsByAdminPrefix, bz...)
}

////queue types

var lenTime = len(sdk.FormatTimeBytes(time.Now()))

// SplitActionQueueKey split the listed key and returns the id and execTime
func SplitActionQueueKey(key []byte) (actionID uint64, execTime time.Time) {
	return splitKeyWithTime(key)
}

// ActionByTimeKey gets the listed item queue key by execTime
func ActionByTimeKey(execTime time.Time) []byte {
	return append(ActionQueuePrefix, sdk.FormatTimeBytes(execTime)...)
}

// from the key we get the action and end time
func splitKeyWithTime(key []byte) (actionID uint64, execTime time.Time) {

	execTime, _ = sdk.ParseTimeBytes(key[1 : 1+lenTime])

	//returns an id from bytes
	actionID = binary.BigEndian.Uint64(key[1+lenTime:])

	return
}

// ActionQueueKey returns the key with prefix for an action in the Listed Item Queue
func ActionQueueKey(actionID uint64, execTime time.Time) []byte {
	return append(ActionByTimeKey(execTime), GetBytesForUint(actionID)...)
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

// GetActionByOwnerIndexKey returns the id: `<prefix><ownerAddress length><ownerAddress><actionID>`
func GetActionByOwnerIndexKey(bz []byte, actionID uint64) []byte {
	prefixBytes := GetActionsByOwnerPrefix(bz)
	lenPrefixBytes := len(prefixBytes)
	r := make([]byte, lenPrefixBytes+8)

	copy(r[:lenPrefixBytes], prefixBytes)
	copy(r[lenPrefixBytes:], GetBytesForUint(actionID))

	return r
}

// GetHostedAccountsByAdminIndexKey returns the id: `<prefix><adminAddress length><adminAddress><hostedaccountID>`
func GetHostedAccountsByAdminIndexKey(bz []byte, hostedAccountAddress string) []byte {
	prefixBytes := GetHostedAccountsByAdminPrefix(bz)
	return append(prefixBytes, []byte(hostedAccountAddress)...)
}
