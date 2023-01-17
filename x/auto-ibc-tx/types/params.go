package types

import (
	"fmt"
	"time"

	yaml "gopkg.in/yaml.v2"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

const (
	// AutoTxFundsCommission percentage to distribute to community pool for leftover balances (rounded up)
	DefaultAutoTxFundsCommission int64 = 2
	// AutoTxConstantFee fee to prevent spam of auto messages, to be distributed to community pool
	DefaultAutoTxConstantFee int64 = 1_000_000 // 1trst
	// AutoTxFlexFeeMul is the denominator for the gas-dependent flex fee to prioritize auto messages in the block, to be distributed to validators
	DefaultAutoTxFlexFeeMul int64 = 100 // 100/100 = 1 = gasUsed
	// RecurringAutoTxConstantFee fee to prevent spam of auto messages, to be distributed to community pool
	DefaultRecurringAutoTxConstantFee int64 = 1_000_000 // 1trst
	// Default max period for a AutoTx that is self-executing
	DefaultMaxAutoTxDuration time.Duration = time.Hour * 24 * 366 * 10 // a little over 10 years
	// MinAutoTxDuration sets the minimum duration for a self-executing AutoTx
	DefaultMinAutoTxDuration time.Duration = time.Second * 40
	// MinAutoTxInterval sets the minimum interval self-execution
	DefaultMinAutoTxInterval time.Duration = time.Second * 20
)

// Parameter store key
var (
	KeyAutoTxFundsCommission = []byte("AutoTxFundsCommission")

	KeyAutoTxFlexFeeMul = []byte("AutoTxFlexFeeMul")

	KeyAutoTxConstantFee = []byte("AutoTxConstantFee")

	KeyRecurringAutoTxConstantFee = []byte("RecurringAutoTxConstantFee")

	KeyMaxAutoTxDuration = []byte("MaxAutoTxDuration")

	KeyMinAutoTxDuration = []byte("MinAutoTxDuration")

	KeyMinAutoTxInterval = []byte("MinAutoTxInterval")
)

// ParamSetPairs - Implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	//fmt.Print("ParamSetPairs..")
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyAutoTxFundsCommission, &p.AutoTxFundsCommission, validateAutoTxFundsCommission),
		paramtypes.NewParamSetPair(KeyAutoTxConstantFee, &p.AutoTxConstantFee, validateAutoTxConstantFee),
		paramtypes.NewParamSetPair(KeyAutoTxFlexFeeMul, &p.AutoTxFlexFeeMul, validateAutoTxFlexFeeMul),
		paramtypes.NewParamSetPair(KeyRecurringAutoTxConstantFee, &p.RecurringAutoTxConstantFee, validateRecurringAutoTxConstantFee),
		paramtypes.NewParamSetPair(KeyMaxAutoTxDuration, &p.MaxAutoTxDuration, validateAutoTxDuration),
		paramtypes.NewParamSetPair(KeyMinAutoTxDuration, &p.MinAutoTxDuration, validateAutoTxDuration),
		paramtypes.NewParamSetPair(KeyMinAutoTxInterval, &p.MinAutoTxInterval, validateAutoTxInterval),
	}
}

// NewParams creates a new Params object
func NewParams(autoTxFundsCommission int64, AutoTxConstantFee int64, AutoTxFlexFeeMul int64, RecurringAutoTxConstantFee int64, maxAutoTxDuration time.Duration, minAutoTxDuration time.Duration, minAutoTxInterval time.Duration) Params {
	//fmt.Printf("default auto-ibc-tx params. %v \n", autoTxFundsCommission)
	return Params{AutoTxFundsCommission: autoTxFundsCommission, AutoTxConstantFee: AutoTxConstantFee, AutoTxFlexFeeMul: AutoTxFlexFeeMul, RecurringAutoTxConstantFee: RecurringAutoTxConstantFee, MaxAutoTxDuration: maxAutoTxDuration, MinAutoTxDuration: minAutoTxDuration, MinAutoTxInterval: minAutoTxInterval}
}

// DefaultParams default parameters for auto-ibc-tx
func DefaultParams() Params {
	//fmt.Print("default auto-ibc-tx params..")
	return NewParams(DefaultAutoTxFundsCommission, DefaultAutoTxConstantFee, DefaultAutoTxFlexFeeMul, DefaultRecurringAutoTxConstantFee, DefaultMaxAutoTxDuration, DefaultMinAutoTxDuration, DefaultMinAutoTxInterval)
}

// Validate validates all params
func (p Params) Validate() error {
	if err := validateAutoTxFundsCommission(p.AutoTxFundsCommission); err != nil {
		return err
	}
	if err := validateAutoTxDuration(p.MaxAutoTxDuration); err != nil {
		return err
	}
	if err := validateAutoTxDuration(p.MinAutoTxDuration); err != nil {
		return err
	}
	if err := validateAutoTxInterval(p.MinAutoTxInterval); err != nil {
		return err
	}
	if err := validateAutoTxConstantFee(p.AutoTxConstantFee); err != nil {
		return err
	}
	if err := validateAutoTxFlexFeeMul(p.AutoTxFlexFeeMul); err != nil {
		return err
	}
	if err := validateRecurringAutoTxConstantFee(p.RecurringAutoTxConstantFee); err != nil {
		return err
	}

	return nil
}

func validateAutoTxFundsCommission(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v < 1 {
		return fmt.Errorf("AutoTxFundsCommission rate must be positive: %d", v)
	}

	return nil
}

func validateAutoTxDuration(i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("self-executing AutoTx period (between initiation and last self-execuion) must be longer: %d", v)
	}

	return nil
}

func validateAutoTxInterval(i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("self-executing AutoTx interval must be longer: %d", v)
	}

	return nil
}

func validateAutoTxConstantFee(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	//10_000_000 = 10TRST we do not want to go this high
	if v > 10_000_000 {
		return fmt.Errorf("AutoTxConstantFee must be lower: %T", i)
	}
	if v < 0 {
		return fmt.Errorf("AutoTxConstantFee must be 0 or higher: %d", v)
	}

	return nil
}
func validateAutoTxFlexFeeMul(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	//5000 = 50x gas cost, we do not want to go this high
	if v > 5000 {
		return fmt.Errorf("AutoTxFlexFeeMul must be lower: %T", i)
	}
	if v < 0 {
		return fmt.Errorf("AutoTxFlexFeeMul rate must be 0 or higher: %d", v)
	}

	return nil
}
func validateRecurringAutoTxConstantFee(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	//10_000_000 = 10TRST we do not want to go this high
	if v > 10_000_000 {
		return fmt.Errorf("RecurringAutoTxConstantFee must be lower: %T", i)
	}
	if v < 0 {
		return fmt.Errorf("RecurringAutoTxConstantFee rate must be 0 or higher:%d", v)
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
