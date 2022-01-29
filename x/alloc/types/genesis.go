package types

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// this line is used by starport scaffolding # genesis/types/import
// this line is used by starport scaffolding # ibc/genesistype/import

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: Params{
			DistributionProportions: DistributionProportions{
				Staking:                     sdk.MustNewDecFromStr("0.35"), // 25%
				TrustlessContractIncentives: sdk.MustNewDecFromStr("0.40"), // 45%
				ItemIncentives:              sdk.MustNewDecFromStr("0.05"), // 45%
				DeveloperRewards:            sdk.MustNewDecFromStr("0.05"), // 25%
				CommunityPool:               sdk.MustNewDecFromStr("0.15"), // 5%
			},
			WeightedDeveloperRewardsReceivers: []WeightedAddress{{
				Address: "trust1sns5l9cvkgf4fy770nmg98e7uzet5xhhmv8njv",
				Weight:  sdk.MustNewDecFromStr("1"),
			}},
		},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	err := gs.Params.Validate()
	if err != nil {
		return err
	}
	return nil
}

// GetGenesisStateFromAppState return GenesisState
func GetGenesisStateFromAppState(cdc codec.JSONCodec, appState map[string]json.RawMessage) *GenesisState {
	var genesisState GenesisState

	if appState[ModuleName] != nil {
		cdc.MustUnmarshalJSON(appState[ModuleName], &genesisState)
	}

	return &genesisState
}
