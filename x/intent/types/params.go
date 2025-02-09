package types

import (
	"fmt"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	yaml "gopkg.in/yaml.v2"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	// FlowFundsCommission percentage to distribute to community pool for leftover balances (rounded up)
	DefaultFlowFundsCommission int64 = 2 //2%
	// BurnFeePerMsg fee to prevent spam of auto messages, to be distributed to community pool
	DefaultBurnFeePerMsg int64 = 5_000 // 0.005trst
	// FlowFlexFeeMul is the denominator for the gas fee
	DefaultFlowFlexFeeMul int64 = 10 // in %
	// GasFeeCoins fee to prevent spam of auto messages, to be distributed to community pool
	DefaultGasFeeCoins sdk.Coins = sdk.NewCoins(sdk.NewCoin(Denom, math.NewInt(1))) // 1uinto
	// Default max period for a Flow that is self-executing
	DefaultMaxFlowDuration time.Duration = time.Hour * 24 * 366 * 10 // a little over 2 years
	// MinFlowDuration sets the minimum duration for a Flow
	DefaultMinFlowDuration time.Duration = time.Second * 60
	// MinFlowInterval sets the minimum interval self-execution
	DefaultMinFlowInterval time.Duration = time.Second * 60
	// DefaultRelayerReward for a given flow type
	DefaultRelayerReward int64 = 10_000 //0.01trst

)

// Parameter store key
var (
	KeyFlowFundsCommission = []byte("FlowFundsCommission")
	KeyFlowFlexFeeMul      = []byte("FlowFlexFeeMul")
	KeyBurnFeePerMsg       = []byte("BurnFeePerMsg")
	KeyGasFeeCoins         = []byte("GasFeeCoins")
	KeyMaxFlowDuration     = []byte("MaxFlowDuration")
	KeyMinFlowDuration     = []byte("MinFlowDuration")
	KeyMinFlowInterval     = []byte("MinFlowInterval")
	KeyRelayerRewards      = []byte("RelayerRewards")
)

// ParamSetPairs - Implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	//fmt.Print("ParamSetPairs..")
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyFlowFundsCommission, &p.FlowFundsCommission, validateFlowFundsCommission),
		paramtypes.NewParamSetPair(KeyBurnFeePerMsg, &p.BurnFeePerMsg, validateBurnFeePerMsg),
		paramtypes.NewParamSetPair(KeyFlowFlexFeeMul, &p.FlowFlexFeeMul, validateFlowFlexFeeMul),
		paramtypes.NewParamSetPair(KeyGasFeeCoins, &p.GasFeeCoins, validateGasFeeCoins),
		paramtypes.NewParamSetPair(KeyMaxFlowDuration, &p.MaxFlowDuration, validateFlowDuration),
		paramtypes.NewParamSetPair(KeyMinFlowDuration, &p.MinFlowDuration, validateFlowDuration),
		paramtypes.NewParamSetPair(KeyMinFlowInterval, &p.MinFlowInterval, validateFlowInterval),
		paramtypes.NewParamSetPair(KeyRelayerRewards, &p.RelayerRewards, validateRelayerRewards),
	}
}

// NewParams creates a new Params object
func NewParams(flowFundsCommission int64, BurnFeePerMsg int64, FlowFlexFeeMul int64, GasFeeCoins sdk.Coins, maxFlowDuration time.Duration, minFlowDuration time.Duration, minFlowInterval time.Duration, relayerRewards []int64) Params {
	//fmt.Printf("default intent params. %v \n", flowFundsCommission)
	return Params{FlowFundsCommission: flowFundsCommission, BurnFeePerMsg: BurnFeePerMsg, FlowFlexFeeMul: FlowFlexFeeMul, GasFeeCoins: GasFeeCoins, MaxFlowDuration: maxFlowDuration, MinFlowDuration: minFlowDuration, MinFlowInterval: minFlowInterval, RelayerRewards: relayerRewards}
}

// DefaultParams default parameters for intent
func DefaultParams() Params {
	//fmt.Print("default intent params..")
	return NewParams(DefaultFlowFundsCommission, DefaultBurnFeePerMsg, DefaultFlowFlexFeeMul, DefaultGasFeeCoins, DefaultMaxFlowDuration, DefaultMinFlowDuration, DefaultMinFlowInterval, []int64{DefaultRelayerReward, DefaultRelayerReward, DefaultRelayerReward, DefaultRelayerReward})
}

// Validate validates all params
func (p Params) Validate() error {
	if err := validateFlowFundsCommission(p.FlowFundsCommission); err != nil {
		return err
	}
	if err := validateFlowDuration(p.MaxFlowDuration); err != nil {
		return err
	}
	if err := validateFlowDuration(p.MinFlowDuration); err != nil {
		return err
	}
	if err := validateFlowInterval(p.MinFlowInterval); err != nil {
		return err
	}
	if err := validateBurnFeePerMsg(p.BurnFeePerMsg); err != nil {
		return err
	}
	if err := validateFlowFlexFeeMul(p.FlowFlexFeeMul); err != nil {
		return err
	}
	if err := validateGasFeeCoins(p.GasFeeCoins); err != nil {
		return err
	}
	if err := validateRelayerRewards(p.RelayerRewards); err != nil {
		return err
	}

	return nil
}

func validateFlowFundsCommission(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v < 1 {
		return fmt.Errorf("FlowFundsCommission rate must be positive: %d", v)
	}

	return nil
}

func validateFlowDuration(i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("Flow period (between initiation and last self-execuion) must be longer: %d", v)
	}

	return nil
}

func validateFlowInterval(i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("Flow interval must be longer: %d", v)
	}

	return nil
}

func validateBurnFeePerMsg(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	//10_000_000 = 10INTO we do not want to go this high
	if v > 10_000_000 {
		return fmt.Errorf("BurnFeePerMsg must be lower: %T", i)
	}
	if v < 0 {
		return fmt.Errorf("BurnFeePerMsg must be 0 or higher: %d", v)
	}

	return nil
}
func validateFlowFlexFeeMul(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	//5000 = 50x gas cost, we do not want to go this high
	if v > 5000 {
		return fmt.Errorf("FlowFlexFeeMul must be lower: %T", i)
	}
	if v < 0 {
		return fmt.Errorf("FlowFlexFeeMul rate must be 0 or higher: %d", v)
	}

	return nil
}
func validateGasFeeCoins(i interface{}) error {
	v, ok := i.(sdk.Coins)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return v.Validate()
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
