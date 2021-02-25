package types

import "fmt"

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		// this line is used by starport scaffolding # genesis/types/default
		EstimatorList: []*Estimator{},
		BuyerList:     []*Buyer{},
		ItemList:      []*Item{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// this line is used by starport scaffolding # genesis/types/validate
	// Check for duplicated ID in estimator
	estimatorIdMap := make(map[string]bool)

	for _, elem := range gs.EstimatorList {
		if _, ok := estimatorIdMap[(elem.Itemid + "-" + elem.Estimator)]; ok {
			return fmt.Errorf("duplicated id for estimator")
		}
		estimatorIdMap[elem.Itemid + "-" + elem.Estimator] = true
	}
	// Check for duplicated ID in buyer
	buyerIdMap := make(map[string]bool)

	for _, elem := range gs.BuyerList {
		if _, ok := buyerIdMap[elem.Itemid]; ok {
			return fmt.Errorf("duplicated id for buyer")
		}
		buyerIdMap[elem.Itemid] = true
	}
	// Check for duplicated ID in item
	itemIdMap := make(map[string]bool)

	for _, elem := range gs.ItemList {
		if _, ok := itemIdMap[elem.Id]; ok {
			return fmt.Errorf("duplicated id for item")
		}
		itemIdMap[elem.Id] = true
	}

	return nil
}
