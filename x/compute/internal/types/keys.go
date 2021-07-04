package types

import (
	"encoding/binary"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the contract module
	ModuleName = "compute"

	// StoreKey is the string store representation
	StoreKey = ModuleName

	// TStoreKey is the string transient store representation
	TStoreKey = "transient_" + ModuleName

	// QuerierRoute is the querier route for the staking module
	QuerierRoute = ModuleName

	// RouterKey is the msg router key for the staking module
	RouterKey = ModuleName
)

const ( // event attributes
	AttributeKeyContract = "contract_address"
	AttributeKeyCodeID   = "code_id"
	AttributeKeySigner   = "signer"
)

// nolint
var (
	CodeKeyPrefix       = []byte{0x01}
	ContractKeyPrefix   = []byte{0x02}
	ContractStorePrefix = []byte{0x03}
	SequenceKeyPrefix   = []byte{0x04}
	// ContractHistoryStorePrefix = []byte{0x05}
	ContractEnclaveIdPrefix = []byte{0x06}
	ContractLabelPrefix     = []byte{0x07}

	ContractQueuePrefix = []byte{0x08}

	KeyLastCodeID     = append(SequenceKeyPrefix, []byte("lastCodeId")...)
	KeyLastInstanceID = append(SequenceKeyPrefix, []byte("lastContractId")...)
)

// GetCodeKey constructs the key for retreiving the ID for the WASM code
func GetCodeKey(codeID uint64) []byte {
	contractIDBz := sdk.Uint64ToBigEndian(codeID)
	return append(CodeKeyPrefix, contractIDBz...)
}

func decodeCodeKey(src []byte) uint64 {
	return binary.BigEndian.Uint64(src[len(CodeKeyPrefix):])
}

// GetContractAddressKey returns the key for the WASM contract instance
func GetContractAddressKey(addr sdk.AccAddress) []byte {
	return append(ContractKeyPrefix, addr...)
}

// GetContractAddressKey returns the key for the WASM contract instance
func GetContractEnclaveKey(addr sdk.AccAddress) []byte {
	return append(ContractEnclaveIdPrefix, addr...)
}

// GetContractStorePrefixKey returns the store prefix for the WASM contract instance
func GetContractStorePrefixKey(addr sdk.AccAddress) []byte {
	return append(ContractStorePrefix, addr...)
}

// GetContractStorePrefixKey returns the store prefix for the WASM contract instance
func GetContractLabelPrefix(addr string) []byte {
	return append(ContractLabelPrefix, []byte(addr)...)
}

////queue types

var lenTime = len(sdk.FormatTimeBytes(time.Now()))

// SplitContractQueueKey split the listed key and returns the id and endTime
func SplitContractQueueKey(key []byte) (contractAddr string, endTime time.Time) {
	return splitKeyWithTime(key)
}

// ContractByTimeKey gets the listed item queue key by endTime
func ContractByTimeKey(endTime time.Time) []byte {
	return append(ContractQueuePrefix, sdk.FormatTimeBytes(endTime)...)
}

//from the key we get the contract and end time
func splitKeyWithTime(key []byte) (contractAddr string, endTime time.Time) {

	/*if len(key[1:]) != 8+lenTime {
		panic(fmt.Sprintf("unexpected key length (%d ≠ %d)", len(key[1:]), lenTime+8))
	}*/

	endTime, _ = sdk.ParseTimeBytes(key[1 : 1+lenTime])
	//	if err != nil {
	//		panic(err)
	//	}
	//fmt.Printf("endTime is %s ", endTime)

	//returns an id from bytes
	contractAddr = string(key[1+lenTime:])

	//	fmt.Printf("contractAddr key is %s ", contractAddr)
	return
}

// ContractQueueKey returns the key with prefix for an contract in the Listed Item Queue
func ContractQueueKey(contractAddr string, endTime time.Time) []byte {
	return append(ContractByTimeKey(endTime), []byte(contractAddr)...)
}
