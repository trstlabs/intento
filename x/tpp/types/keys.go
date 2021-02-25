package types

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
