package types

import (
	"fmt"
	"time"

	yaml "gopkg.in/yaml.v2"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

const (
	// AutoMsgFundsCommission percentage to distribute to community pool for leftover balances (rounded up)
	DefaultAutoMsgFundsCommission int64 = 2

	// AutoMsgConstantFee fee to prevent spam of auto messages, to be distributed to community pool
	DefaultAutoMsgConstantFee int64 = 1000000 // 1utrst

	// ContinuousAutoMsgConstantFee fee to prevent spam of auto messages, to be distributed to community pool
	DefaultContinuousAutoMsgConstantFee int64 = 1000000 // 1utrst

	// Default max period for a contract that is self-executing
	DefaultMaxContractDuration time.Duration = time.Hour * 24 * 366 // 366 days
	// MinContractDuration sets the minimum duration for a self-executing contract

	DefaultMinContractDuration time.Duration = time.Second * 45 // time.Hour * 24 // 1 day
	// MinContractDurationForIncentive to distribute reward to contracts we want to incentivize
	DefaultMinContractDurationForIncentive time.Duration = time.Hour * 24 // time.Hour * 24 // 1 day

	// DefaultMaxContractIncentive max amount of utrst coins to give to a contract as incentive
	DefaultMaxContractIncentive int64 = 500000000 // 500utrst

	// MinContractBalanceForIncentive minimum balance required to be elligable for an incentive
	DefaultMinContractBalanceForIncentive int64 = 50000000 // 50utrst
)

// Parameter store key
var (
	KeyAutoMsgFundsCommission = []byte("AutoMsgFundsCommission")

	KeyAutoMsgConstantFee = []byte("AutoMsgConstantFee")

	KeyContinuousAutoMsgConstantFee = []byte("ContinuousAutoMsgConstantFee")

	KeyMaxContractDuration = []byte("MaxContractDuration")

	KeyMinContractDuration = []byte("MinContractDuration")

	KeyMinContractDurationForIncentive = []byte("MinContractDurationForIncentive")

	KeyMaxContractIncentive = []byte("MaxContractIncentive")

	KeyMinContractBalanceForIncentive = []byte("MinContractBalanceForIncentive")
)

// ParamSetPairs - Implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	//	fmt.Print("default ParamSetPairs params..")
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyAutoMsgFundsCommission, &p.AutoMsgFundsCommission, validateAutoMsgFundsCommission),
		paramtypes.NewParamSetPair(KeyAutoMsgConstantFee, &p.AutoMsgConstantFee, validateAutoMsgConstantFee),
		paramtypes.NewParamSetPair(KeyContinuousAutoMsgConstantFee, &p.ContinuousAutoMsgConstantFee, validateContinuousAutoMsgConstantFee),
		paramtypes.NewParamSetPair(KeyMaxContractDuration, &p.MaxContractDuration, validateContractDuration),
		paramtypes.NewParamSetPair(KeyMinContractDuration, &p.MinContractDuration, validateContractDuration),
		paramtypes.NewParamSetPair(KeyMinContractDurationForIncentive, &p.MinContractDurationForIncentive, validateMinContractDurationForIncentive),
		paramtypes.NewParamSetPair(KeyMaxContractIncentive, &p.MaxContractIncentive, validateMaxContractIncentive),
		paramtypes.NewParamSetPair(KeyMinContractBalanceForIncentive, &p.MinContractBalanceForIncentive, validateMinContractBalanceForIncentive),
	}
}

// NewParams creates a new Params object
func NewParams(autoMsgFundsCommission int64, autoMsgConstantFee int64, continuousAutoMsgConstantFee int64, maxContractDuration time.Duration, minContractDuration time.Duration, minContractDurationForIncentive time.Duration, maxContractIncentive int64, minContractBalanceForIncentive int64) Params {
	return Params{AutoMsgFundsCommission: autoMsgFundsCommission, AutoMsgConstantFee: autoMsgConstantFee, ContinuousAutoMsgConstantFee: continuousAutoMsgConstantFee, MaxContractDuration: maxContractDuration, MinContractDuration: minContractDuration, MinContractDurationForIncentive: minContractDurationForIncentive, MaxContractIncentive: maxContractIncentive, MinContractBalanceForIncentive: minContractBalanceForIncentive}
}

// DefaultParams default parameters for compute
func DefaultParams() Params {
	//fmt.Print("default compute params..")
	return NewParams(DefaultAutoMsgFundsCommission, DefaultAutoMsgConstantFee, DefaultContinuousAutoMsgConstantFee, DefaultMaxContractDuration, DefaultMinContractDuration, DefaultMinContractDurationForIncentive, DefaultMaxContractIncentive, DefaultMinContractBalanceForIncentive)
}

// Validate validates all params
func (p Params) Validate() error {
	if err := validateContractDuration(p.MaxContractDuration); err != nil {
		return err
	}
	if err := validateContractDuration(p.MinContractDuration); err != nil {
		return err
	}

	if err := validateAutoMsgFundsCommission(p.AutoMsgFundsCommission); err != nil {
		return err
	}

	return nil
}

func validateContractDuration(i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("self-executing contract period between initiation and self-execuion must be longer: %d", v)
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

func validateAutoMsgFundsCommission(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v < 1 {
		return fmt.Errorf("AutoMsgFundsCommission rate must be positive: %d", v)
	}

	return nil
}
func validateAutoMsgConstantFee(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v < 0 {
		return fmt.Errorf("AutoMsgConstantFee rate must be 0 or higher: %d", v)
	}

	return nil
}
func validateContinuousAutoMsgConstantFee(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v < 1 {
		return fmt.Errorf("AutoMsgConstantFee rate must be 0 or higher:%d", v)
	}

	return nil
}
func validateMaxContractIncentive(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v < 1 {
		return fmt.Errorf("AutoMsgFundsCommission rate must be positive: %d", v)
	}

	return nil
}
func validateMinContractBalanceForIncentive(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v < 1 {
		return fmt.Errorf("AutoMsgFundsCommission rate must be positive: %d", v)
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
