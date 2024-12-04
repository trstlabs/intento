package types

import (
	"encoding/json"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
)

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default alloc genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: Params{
			DistributionProportions: DistributionProportions{
				RelayerIncentives: math.LegacyNewDecWithPrec(10, 2),
				DeveloperRewards:  math.LegacyNewDecWithPrec(5, 2),
				CommunityPool:     math.LegacyNewDecWithPrec(30, 2),
			},
			WeightedDeveloperRewardsReceivers: []WeightedAddress{},
		},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	return gs.Params.Validate()
}

// GetGenesisStateFromAppState return GenesisState
func GetGenesisStateFromAppState(cdc codec.JSONCodec, appState map[string]json.RawMessage) *GenesisState {
	var genesisState GenesisState

	if appState[ModuleName] != nil {
		cdc.MustUnmarshalJSON(appState[ModuleName], &genesisState)
	}

	return &genesisState
}
