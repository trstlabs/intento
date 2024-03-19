package types

import (
	"fmt"
	"time"

	yaml "gopkg.in/yaml.v2"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

const (
	// ActionFundsCommission percentage to distribute to community pool for leftover balances (rounded up)
	DefaultActionFundsCommission int64 = 2 //2%
	// ActionConstantFee fee to prevent spam of auto messages, to be distributed to community pool
	DefaultActionConstantFee int64 = 5_000 // 0.005trst
	// ActionFlexFeeMul is the denominator for the gas-dependent flex fee to prioritize auto messages in the block, to be distributed to validators
	DefaultActionFlexFeeMul int64 = 3 // 3% of minutes for a given period as uinto (1_000m = 20uinto)
	// RecurringActionConstantFee fee to prevent spam of auto messages, to be distributed to community pool
	DefaultRecurringActionConstantFee int64 = 5_000 // 0.005trst
	// Default max period for a Action that is self-executing
	DefaultMaxActionDuration time.Duration = time.Hour * 24 * 366 * 2 // a little over 2 years
	// MinActionDuration sets the minimum duration for a self-executing Action
	DefaultMinActionDuration time.Duration = time.Second * 60
	// MinActionInterval sets the minimum interval self-execution
	DefaultMinActionInterval time.Duration = time.Second * 60
	// DefaultRelayerReward for a given action type
	DefaultRelayerReward int64 = 10_000 //0.01trst
)

// Parameter store key
var (
	KeyActionFundsCommission      = []byte("ActionFundsCommission")
	KeyActionFlexFeeMul           = []byte("ActionFlexFeeMul")
	KeyActionConstantFee          = []byte("ActionConstantFee")
	KeyRecurringActionConstantFee = []byte("RecurringActionConstantFee")
	KeyMaxActionDuration          = []byte("MaxActionDuration")
	KeyMinActionDuration          = []byte("MinActionDuration")
	KeyMinActionInterval          = []byte("MinActionInterval")
	KeyRelayerRewards             = []byte("RelayerRewards")
)

// ParamSetPairs - Implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	//fmt.Print("ParamSetPairs..")
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyActionFundsCommission, &p.ActionFundsCommission, validateActionFundsCommission),
		paramtypes.NewParamSetPair(KeyActionConstantFee, &p.ActionConstantFee, validateActionConstantFee),
		paramtypes.NewParamSetPair(KeyActionFlexFeeMul, &p.ActionFlexFeeMul, validateActionFlexFeeMul),
		paramtypes.NewParamSetPair(KeyRecurringActionConstantFee, &p.RecurringActionConstantFee, validateRecurringActionConstantFee),
		paramtypes.NewParamSetPair(KeyMaxActionDuration, &p.MaxActionDuration, validateActionDuration),
		paramtypes.NewParamSetPair(KeyMinActionDuration, &p.MinActionDuration, validateActionDuration),
		paramtypes.NewParamSetPair(KeyMinActionInterval, &p.MinActionInterval, validateActionInterval),
		paramtypes.NewParamSetPair(KeyRelayerRewards, &p.RelayerRewards, validateRelayerRewards),
	}
}

// NewParams creates a new Params object
func NewParams(autoTxFundsCommission int64, ActionConstantFee int64, ActionFlexFeeMul int64, RecurringActionConstantFee int64, maxActionDuration time.Duration, minActionDuration time.Duration, minActionInterval time.Duration, relayerRewards []int64) Params {
	//fmt.Printf("default intent params. %v \n", autoTxFundsCommission)
	return Params{ActionFundsCommission: autoTxFundsCommission, ActionConstantFee: ActionConstantFee, ActionFlexFeeMul: ActionFlexFeeMul, RecurringActionConstantFee: RecurringActionConstantFee, MaxActionDuration: maxActionDuration, MinActionDuration: minActionDuration, MinActionInterval: minActionInterval, RelayerRewards: relayerRewards}
}

// DefaultParams default parameters for intent
func DefaultParams() Params {
	//fmt.Print("default intent params..")
	return NewParams(DefaultActionFundsCommission, DefaultActionConstantFee, DefaultActionFlexFeeMul, DefaultRecurringActionConstantFee, DefaultMaxActionDuration, DefaultMinActionDuration, DefaultMinActionInterval, []int64{DefaultRelayerReward, DefaultRelayerReward, DefaultRelayerReward, DefaultRelayerReward})
}

// Validate validates all params
func (p Params) Validate() error {
	if err := validateActionFundsCommission(p.ActionFundsCommission); err != nil {
		return err
	}
	if err := validateActionDuration(p.MaxActionDuration); err != nil {
		return err
	}
	if err := validateActionDuration(p.MinActionDuration); err != nil {
		return err
	}
	if err := validateActionInterval(p.MinActionInterval); err != nil {
		return err
	}
	if err := validateActionConstantFee(p.ActionConstantFee); err != nil {
		return err
	}
	if err := validateActionFlexFeeMul(p.ActionFlexFeeMul); err != nil {
		return err
	}
	if err := validateRecurringActionConstantFee(p.RecurringActionConstantFee); err != nil {
		return err
	}
	if err := validateRelayerRewards(p.RelayerRewards); err != nil {
		return err
	}

	return nil
}

func validateActionFundsCommission(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v < 1 {
		return fmt.Errorf("ActionFundsCommission rate must be positive: %d", v)
	}

	return nil
}

func validateActionDuration(i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("self-executing Action period (between initiation and last self-execuion) must be longer: %d", v)
	}

	return nil
}

func validateActionInterval(i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("self-executing Action interval must be longer: %d", v)
	}

	return nil
}

func validateActionConstantFee(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	//10_000_000 = 10INTO we do not want to go this high
	if v > 10_000_000 {
		return fmt.Errorf("ActionConstantFee must be lower: %T", i)
	}
	if v < 0 {
		return fmt.Errorf("ActionConstantFee must be 0 or higher: %d", v)
	}

	return nil
}
func validateActionFlexFeeMul(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	//5000 = 50x gas cost, we do not want to go this high
	if v > 5000 {
		return fmt.Errorf("ActionFlexFeeMul must be lower: %T", i)
	}
	if v < 0 {
		return fmt.Errorf("ActionFlexFeeMul rate must be 0 or higher: %d", v)
	}

	return nil
}
func validateRecurringActionConstantFee(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	//10_000_000 = 10INTO we do not want to go this high
	if v > 10_000_000 {
		return fmt.Errorf("RecurringActionConstantFee must be lower: %T", i)
	}
	if v < 0 {
		return fmt.Errorf("RecurringActionConstantFee rate must be 0 or higher:%d", v)
	}

	return nil
}

func validateRelayerRewards(i interface{}) error {
	list, ok := i.([]int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	for i, v := range list {
		// 10_000_000 = 10INTO we do not want to go this high
		if v > 10_000_000 {
			return fmt.Errorf("RelayerReward for message must be lower: %T", i)
		}

		if i > 3 {
			return fmt.Errorf("only 4 types of incentives supported for now: %d", v)
		}
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
