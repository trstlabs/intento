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
	DefaultContractPeriod time.Duration = time.Hour * 24 * 30 // 30 days
)

// Commission to distribute to community pool for leftover balances (rounded up)
const (
	DefaultCommission int64 = 2 //time.Hour * 24 * 30 // 30 days
)

// Parameter store key
var (
	KeyMaxActiveContractPeriod = []byte("MaxActiveContractPeriod")
)

// Parameter store key
var (
	KeyCommission = []byte("Commission")
)

// ParamSetPairs - Implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	//	fmt.Print("default ParamSetPairs params..")
	return paramtypes.ParamSetPairs{

		paramtypes.NewParamSetPair(KeyMaxActiveContractPeriod, &p.MaxActiveContractPeriod, validateMaxActiveContractPeriod),
		paramtypes.NewParamSetPair(KeyCommission, &p.Commission, validateCommission),
	}
}

// NewParams creates a new ActiveParams object
func NewParams(maxActiveContractPeriod time.Duration, commission int64) Params {
	return Params{MaxActiveContractPeriod: maxActiveContractPeriod, Commission: commission}
}

// DefaultParams default parameters for Active
func DefaultParams() Params {
	//fmt.Print("default compute params..")
	return NewParams(DefaultContractPeriod, DefaultCommission)
}

// Validate validates all params
func (p Params) Validate() error {
	if err := validateMaxActiveContractPeriod(p.MaxActiveContractPeriod); err != nil {
		return err
	}

	if err := validateCommission(p.Commission); err != nil {
		return err
	}

	return nil
}

func validateMaxActiveContractPeriod(i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("active contract period must be positive: %d", v)
	}

	return nil
}

func validateCommission(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v < 1 {
		return fmt.Errorf("commission rate must be positive: %d", v)
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
