package types

import (
	"fmt"
	"time"

	yaml "gopkg.in/yaml.v2"

	//sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Default period for active
const (
	DefaultMaxActivePeriod       time.Duration = time.Hour * 24 * 30 // 30 days
	DefaultEstimatorCreatorRatio int64         = 50                  //range from 0 - 100
	DefaultMaxBuyerReward        int64         = 500000000           //amount of utrst

//MaxSameCreator 5
//MaxCreatorRatio 20%
)

// Parameter store key
var (
	KeyMaxActivePeriod          = []byte("MaxActivePeriod")
	KeyMaxEstimatorCreatorRatio = []byte("MaxEstimatorCreatorRatio")
	KeyMaxBuyerReward           = []byte("MaxBuyerReward")
)

// ParamSetPairs - Implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyMaxActivePeriod, &p.MaxActivePeriod, validateMaxActivePeriod),
		paramtypes.NewParamSetPair(KeyMaxEstimatorCreatorRatio, &p.MaxEstimatorCreatorRatio, validateMaxEstimatorCreatorRatio),
		paramtypes.NewParamSetPair(KeyMaxBuyerReward, &p.MaxBuyerReward, validateMaxBuyerReward),
	}
}

// NewParams creates a new ActiveParams object
func NewParams(maxActivePeriod time.Duration, maxEstimatorCreatorRatio int64, maxBuyerReward int64) Params {
	return Params{MaxActivePeriod: maxActivePeriod, MaxEstimatorCreatorRatio: maxEstimatorCreatorRatio, MaxBuyerReward: maxBuyerReward}
}

// DefaultParams default parameters for Active
func DefaultParams() Params {

	return NewParams(DefaultMaxActivePeriod, DefaultEstimatorCreatorRatio, DefaultMaxBuyerReward)
}

// Validate validates all params
func (p Params) Validate() error {
	if err := validateMaxActivePeriod(p.MaxActivePeriod); err != nil {
		return err
	}
	if err := validateMaxEstimatorCreatorRatio(p.MaxEstimatorCreatorRatio); err != nil {
		return err
	}
	if err := validateMaxBuyerReward(p.MaxBuyerReward); err != nil {
		return err
	}

	return nil
}

func validateMaxActivePeriod(i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("active must be positive: %d", v)
	}

	return nil
}

func validateMaxEstimatorCreatorRatio(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if 0 > v || v > 100 {
		return fmt.Errorf("ratio must be within 0-100: %d", v)
	}

	return nil
}
func validateMaxBuyerReward(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v < 0 {
		return fmt.Errorf("must be positive: %d", v)
	}

	return nil
}

// String implements the stringer interface for Params
func (p Params) String() string {
	out, err := yaml.Marshal(p)
	if err != nil {
		return ""
	}
	return string(out)
}
