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
	ParamsKey                             = []byte{0x01}
	FlowKeyPrefix                         = []byte{0x02}
	FlowHistoryPrefix                     = []byte{0x03}
	FlowQueuePrefix                       = []byte{0x04}
	SequenceKeyPrefix                     = []byte{0x05}
	FlowsByOwnerPrefix                    = []byte{0x06}
	TmpFlowIDLatestTX                     = []byte{0x07}
	KeyRelayerRewardsAvailability         = []byte{0x08}
	FlowHistorySequencePrefix             = []byte{0x10}
	TrustlessExecutionAgentKeyPrefix      = []byte{0x11}
	TrustlessExecutionAgentsByAdminPrefix = []byte{0x12}
	FlowFeedbackLoopQueryKeyPrefix        = []byte{0x14}
	FlowComparisonQueryKeyPrefix          = []byte{0x15}
	KeyLastID                             = append(SequenceKeyPrefix, []byte("lastId")...)
	KeyLastTxAddrID                       = append(SequenceKeyPrefix, []byte("lastTxAddrId")...)
)

// ics 20 hook
var SenderPrefix = "ibc-flow-hook-intermediary"

var (
	KeyFlowIncentiveForSDKTx   = 0
	KeyFlowIncentiveForWasmTx  = 1
	KeyFlowIncentiveForOsmoTx  = 2
	KeyFlowIncentiveForAuthzTx = 3
)

// GetFlowKey returns the key for the flow
func GetFlowKey(flowID uint64) []byte {
	return append(FlowKeyPrefix, GetBytesForUint(flowID)...)
}

// GetFlowHistoryKey returns the key for the flow
func GetFlowHistoryKey(flowID uint64) []byte {
	return append(FlowHistoryPrefix, GetBytesForUint(flowID)...)
}

// GetFlowsByOwnerPrefix returns the flows by creator prefix
func GetFlowsByOwnerPrefix(addr sdk.AccAddress) []byte {
	bz := address.MustLengthPrefix(addr)
	return append(FlowsByOwnerPrefix, bz...)
}

// GetTrustlessExecutionAgentKey returns the key for the trustless excution agent
func GetTrustlessExecutionAgentKey(address string) []byte {
	return append(TrustlessExecutionAgentKeyPrefix, []byte(address)...)
}

// GetTrustlessExecutionAgentsByAdminPrefix returns the flows by creator prefix
func GetTrustlessExecutionAgentsByAdminPrefix(addr sdk.AccAddress) []byte {
	bz := address.MustLengthPrefix(addr)
	return append(TrustlessExecutionAgentsByAdminPrefix, bz...)
}

////queue types

var lenTime = len(sdk.FormatTimeBytes(time.Now()))

// SplitFlowQueueKey split the listed key and returns the id and execTime
func SplitFlowQueueKey(key []byte) (flowID uint64, execTime time.Time) {
	return splitKeyWithTime(key)
}

// FlowByTimeKey gets the listed item queue key by execTime
func FlowByTimeKey(execTime time.Time) []byte {
	return append(FlowQueuePrefix, sdk.FormatTimeBytes(execTime)...)
}

// from the key we get the flow and end time
func splitKeyWithTime(key []byte) (flowID uint64, execTime time.Time) {

	execTime, _ = sdk.ParseTimeBytes(key[1 : 1+lenTime])

	//returns an id from bytes
	flowID = binary.BigEndian.Uint64(key[1+lenTime:])

	return
}

// FlowQueueKey returns the key with prefix for an flow in the Listed Item Queue
func FlowQueueKey(flowID uint64, execTime time.Time) []byte {
	return append(FlowByTimeKey(execTime), GetBytesForUint(flowID)...)
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

// GetFlowByOwnerIndexKey returns the id: `<prefix><ownerAddress length><ownerAddress><flowID>`
func GetFlowByOwnerIndexKey(bz []byte, flowID uint64) []byte {
	prefixBytes := GetFlowsByOwnerPrefix(bz)
	lenPrefixBytes := len(prefixBytes)
	r := make([]byte, lenPrefixBytes+8)

	copy(r[:lenPrefixBytes], prefixBytes)
	copy(r[lenPrefixBytes:], GetBytesForUint(flowID))

	return r
}

// GetTrustlessExecutionAgentsByAdminIndexKey returns the id: `<prefix><adminAddress length><adminAddress><trustlessexecutionagentID>`
func GetTrustlessExecutionAgentsByAdminIndexKey(bz []byte, trustlessExecutionAgentAddress string) []byte {
	prefixBytes := GetTrustlessExecutionAgentsByAdminPrefix(bz)
	return append(prefixBytes, []byte(trustlessExecutionAgentAddress)...)
}
