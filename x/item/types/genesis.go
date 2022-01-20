package types

import (
	"fmt"
)

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		// this line is used by starport scaffolding # genesis/types/default
		ProfileList: []*Profile{},
		//BuyerList:     []*Buyer{},
		ItemList: []*Item{},
		Params:   DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// this line is used by starport scaffolding # genesis/types/validate

	// Check for duplicated ID in item
	itemIdMap := make(map[uint64]bool)

	for _, elem := range gs.ItemList {
		if _, ok := itemIdMap[elem.Id]; ok {
			return fmt.Errorf("duplicated id for item")
		}
		itemIdMap[elem.Id] = true

	}

	// Check for duplicated ID in estimator profiles
	itemIdEstimatorMap := make(map[uint64]bool)

	for _, elem := range gs.ProfileList {
		for _, elem := range elem.Estimations {

			if _, ok := itemIdEstimatorMap[elem.Itemid]; ok {
				return fmt.Errorf("duplicated estimation in info")
			}

			itemIdEstimatorMap[elem.Itemid] = true
		}
	}
	/*
		err := gs.Params.ValidateParams()
		if err != nil {
			return err
		}
	*/
	return nil
}
