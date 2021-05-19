package types


import (
	
	
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

)


const (
	// ModuleName defines the module name
	ModuleName = "tpp"

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
	BuyerKey      = "Buyer-value-"
	BuyerCountKey = "Buyer-count-"
)

const (
	EstimatorKey      = "Estimator-value-"
	EstimatorCountKey = "Estimator-count-"
)

var InactiveItemQueuePrefix = []byte{0x02}


var lenTime = len(sdk.FormatTimeBytes(time.Now()))

// SplitInactiveProposalQueueKey split the inactive key and returns the id and endTime
func SplitInactiveItemQueueKey(key []byte) (itemid string, endTime time.Time) {
	return splitKeyWithTime(key)
}


// InactiveProposalByTimeKey gets the inactive proposal queue key by endTime
func InactiveItemByTimeKey(endTime time.Time) []byte {
	return append(InactiveItemQueuePrefix, sdk.FormatTimeBytes(endTime)...)
}//////we have prefix and end time only? not sufficient??



//from the key we get the itemid and end time
func splitKeyWithTime(key []byte) (itemid string, endTime time.Time) {
//	if len(key[1:]) != 22+lenTime {
//		panic(fmt.Sprintf("unexpected key length (%d ≠ %d)", len(key[1:]), lenTime+22))
//	}

	endTime, err := sdk.ParseTimeBytes(key[1 : 1+lenTime])
	if err != nil {
		panic(err)
	}

	//not sure about this, in gov returns an id from bytes
	itemid = string(key[1+lenTime:])
	return
}


// InactiveProposalQueueKey returns the key with prefix for an itemid in the inactiveProposalQueue
func InactiveItemQueueKey(itemid string, endTime time.Time) []byte {
	return append(InactiveItemByTimeKey(endTime), []byte(itemid)...)
}