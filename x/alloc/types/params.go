package types

import (
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter store keys
var (
	KeyDistributionProportions    = []byte("DistributionProportions")
	KeyContributorRewardsReceiver = []byte("ContributorRewardsReceiver")
)

// ParamTable for module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(
	distrProportions DistributionProportions,
	weightedDevRewardsReceivers []WeightedAddress,
) Params {

	return Params{
		DistributionProportions:             distrProportions,
		WeightedContributorRewardsReceivers: weightedDevRewardsReceivers,
	}
}

// default module parameters
func DefaultParams() Params {
	return Params{
		DistributionProportions: DistributionProportions{
			Staking:            sdk.MustNewDecFromStr("0.60"),
			RelayerIncentives:  sdk.MustNewDecFromStr("0.10"),
			ContributorRewards: sdk.MustNewDecFromStr("0.05"),
			CommunityPool:      sdk.MustNewDecFromStr("0.25"),
		},
		WeightedContributorRewardsReceivers: []WeightedAddress{},
	}
}

// validate params
func (p Params) Validate() error {
	if err := validateDistributionProportions(p.DistributionProportions); err != nil {
		return err
	}
	err := validateWeightedContributorRewardsReceivers(p.WeightedContributorRewardsReceivers)
	return err
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyDistributionProportions, &p.DistributionProportions, validateDistributionProportions),
		paramtypes.NewParamSetPair(
			KeyContributorRewardsReceiver, &p.WeightedContributorRewardsReceivers, validateWeightedContributorRewardsReceivers),
	}
}

func validateDistributionProportions(i interface{}) error {
	v, ok := i.(DistributionProportions)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.RelayerIncentives.IsNegative() {
		return errors.New("contract incentives distribution ratio should not be negative")
	}

	if v.Staking.IsNegative() {
		return errors.New("staking distribution ratio should not be negative")
	}

	if v.ContributorRewards.IsNegative() {
		return errors.New("developer rewards distribution ratio should not be negative")
	}

	if v.CommunityPool.IsNegative() {
		return errors.New("community pool distribution ratio should not be negative")
	}

	totalProportions := v.RelayerIncentives.Add(v.ContributorRewards)

	// at least 60% is allocated to incentives

	if !totalProportions.Sub(sdk.NewDecWithPrec(60, 2)).IsNegative() {
		return errors.New("total distributions ratio should be at least 60%")
	}

	return nil
}

func validateWeightedContributorRewardsReceivers(i interface{}) error {
	v, ok := i.([]WeightedAddress)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	// fund community pool when rewards address is empty
	if len(v) == 0 {
		return nil
	}

	weightSum := sdk.NewDec(0)
	for i, w := range v {
		// we allow address to be "" to go to community pool
		if w.Address != "" {
			_, err := sdk.AccAddressFromBech32(w.Address)
			if err != nil {
				return fmt.Errorf("invalid address at %dth", i)
			}
		}
		if !w.Weight.IsPositive() {
			return fmt.Errorf("non-positive weight at %dth", i)
		}
		if w.Weight.GT(sdk.NewDec(1)) {
			return fmt.Errorf("more than 1 weight at %dth", i)
		}
		weightSum = weightSum.Add(w.Weight)
	}

	if !weightSum.Equal(sdk.NewDec(1)) {
		return fmt.Errorf("invalid weight sum: %s", weightSum.String())
	}

	return nil
}
