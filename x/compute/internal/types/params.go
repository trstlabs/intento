package types

import (
	"fmt"
	"time"

	yaml "gopkg.in/yaml.v2"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Default max period for a contract that is active
const (
	DefaultMaxContractPeriod time.Duration = time.Hour * 24 * 30 // 30 days

	// Commission to distribute to community pool for leftover balances (rounded up)

	DefaultCommission int64 = 2

	// MinContractDurationForIncentive to distribute reward to contract

	DefaultMinContractDurationForIncentive time.Duration = time.Second // time.Hour * 24 // 1 day

	// DefaultMaxContractIncentive max amount of utrst coins to give to a contract as incentive

	DefaultMaxContractIncentive int64 = 500000000 // 500utrst

	// MinContractBalanceForIncentive minimum balance required to be elligable for an incentive

	DefaultMinContractBalanceForIncentive int64 = 50000000 // 50utrst
)

// Parameter store key
var (
	KeyMaxActiveContractPeriod = []byte("MaxActiveContractPeriod")

	KeyMinContractDurationForIncentive = []byte("MinContractDurationForIncentive")

	KeyCommission = []byte("Commission")

	KeyMaxContractIncentive = []byte("MaxContractIncentive")

	KeyMinContractBalanceForIncentive = []byte("MinContractBalanceForIncentive")
)

// ParamSetPairs - Implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	//	fmt.Print("default ParamSetPairs params..")
	return paramtypes.ParamSetPairs{

		paramtypes.NewParamSetPair(KeyMaxActiveContractPeriod, &p.MaxActiveContractPeriod, validateMaxActiveContractPeriod),
		paramtypes.NewParamSetPair(KeyCommission, &p.Commission, validateCommission),
		paramtypes.NewParamSetPair(KeyMinContractDurationForIncentive, &p.MinContractDurationForIncentive, validateMinContractDurationForIncentive),
		paramtypes.NewParamSetPair(KeyMaxContractIncentive, &p.MaxContractIncentive, validateMaxContractIncentive),
		paramtypes.NewParamSetPair(KeyMinContractBalanceForIncentive, &p.MinContractBalanceForIncentive, validateMinContractBalanceForIncentive),
	}
}

// NewParams creates a new ActiveParams object
func NewParams(maxActiveContractPeriod time.Duration, commission int64, minContractDurationForIncentive time.Duration, maxContractIncentive int64, minContractBalanceForIncentive int64) Params {
	return Params{MaxActiveContractPeriod: maxActiveContractPeriod, Commission: commission, MinContractDurationForIncentive: minContractDurationForIncentive, MaxContractIncentive: maxContractIncentive, MinContractBalanceForIncentive: minContractBalanceForIncentive}
}

// DefaultParams default parameters for Active
func DefaultParams() Params {
	//fmt.Print("default compute params..")
	return NewParams(DefaultMaxContractPeriod, DefaultCommission, DefaultMinContractDurationForIncentive, DefaultMaxContractIncentive, DefaultMinContractBalanceForIncentive)
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

func validateMinContractDurationForIncentive(i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("min contract for reward duration must be positive: %d", v)
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
func validateMaxContractIncentive(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v < 1 {
		return fmt.Errorf("commission rate must be positive: %d", v)
	}

	return nil
}
func validateMinContractBalanceForIncentive(i interface{}) error {
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
