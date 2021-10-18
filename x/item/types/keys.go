package types

import (
	"encoding/binary"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "item"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_capability"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

const (
	ItemKey      = "Item-value-"
	ItemCountKey = "Item-count-"
)

const (
	BuyerKey = "Buyer-value-"
	//BuyerCountKey = "Buyer-count-"
)

const (
	EstimatorKey      = "Estimator-value-"
	EstimatorCountKey = "Estimator-count-"
)

var ListedItemQueuePrefix = []byte{0x02}
var ItemSellerPrefix = []byte{0x03}

var lenTime = len(sdk.FormatTimeBytes(time.Now()))

// SplitInactiveProposalQueueKey split the listed key and returns the id and endTime
func SplitListedItemQueueKey(key []byte) (itemid uint64, endTime time.Time) {
	return splitKeyWithTime(key)
}

// InactiveProposalByTimeKey gets the listed item queue key by endTime
func ListedItemByTimeKey(endTime time.Time) []byte {
	return append(ListedItemQueuePrefix, sdk.FormatTimeBytes(endTime)...)
}

//from the key we get the itemid and end time
func splitKeyWithTime(key []byte) (itemid uint64, endTime time.Time) {
	if len(key[1:]) != 8+lenTime {
		panic(fmt.Sprintf("unexpected key length (%d ≠ %d)", len(key[1:]), lenTime+8))
	}

	endTime, err := sdk.ParseTimeBytes(key[1 : 1+lenTime])
	if err != nil {
		panic(err)
	}

	//eturns an id from bytes
	itemid = GetItemIDFromBytes(key[1+lenTime:])
	return
}

// InactiveProposalQueueKey returns the key with prefix for an itemid in the Listed Item Queue
func ListedItemQueueKey(itemid uint64, endTime time.Time) []byte {
	return append(ListedItemByTimeKey(endTime), Uint64ToByte(itemid)...)
}

//----seller functions

// ItemSellerKey returns the key with prefix for an itemid in seller
func ItemSellerKey(itemid uint64, seller string) []byte {
	return append(ItemSellerBySellerKey(seller), Uint64ToByte(itemid)...)
}

// ItemSellerBySellerKey
func ItemSellerBySellerKey(seller string) []byte {
	return append(ItemSellerPrefix, []byte(seller)...)
}

// IteItemBuyerKeyBuyerKey returns the key with prefix for an itemid in seller
func ItemBuyerKey(itemid uint64, buyer string) []byte {
	return append(ItemBuyerByBuyerKey(buyer), Uint64ToByte(itemid)...)
}

// ItemBuyerByBuyerKey
func ItemBuyerByBuyerKey(buyer string) []byte {
	return append(KeyPrefix(BuyerKey), []byte(buyer)...)
}

/// helper functions

// Uint64ToByte - marshals uint64 to a bigendian byte slice so it can be sorted
func Uint64ToByte(i uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, i)
	return b
}

// GetItemIDBytes returns the byte representation of the itemID
func GetItemIDBytes(itemID uint64) (ItemIDBz []byte) {
	ItemIDBz = make([]byte, 8)
	binary.BigEndian.PutUint64(ItemIDBz, itemID)
	return
}

// GetItemIDFromBytes returns itemID in uint64 format from a byte array
func GetItemIDFromBytes(bz []byte) (itemID uint64) {
	return binary.BigEndian.Uint64(bz)
}
