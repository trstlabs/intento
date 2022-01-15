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
	DefaultPeriod time.Duration = time.Minute * 2 //time.Minute //time.Hour * 24 * 30 // 30 days
)

// Parameter store key
var (
	KeyMaxActivePeriod = []byte("MaxActivePeriod")
)

/*
// ParamKeyTable - Key declaration for parameters
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable(
		paramtypes.NewParamSetPair(ParamStoreKeyActiveParams, ActiveParams{}, validateActiveParams),

	)
}



// ParamKeyTable - Key declaration for parameters
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable(
		paramtypes.NewParamSetPair(KeyMaxActivePeriod, Params{}, Validate),

	)
}

*/

// ParamSetPairs - Implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyMaxActivePeriod, &p.MaxActivePeriod, validateMaxActivePeriod),
	}
}

// NewParams creates a new ActiveParams object
func NewParams(maxActivePeriod time.Duration) Params {
	return Params{MaxActivePeriod: maxActivePeriod}
}

// DefaultParams default parameters for Active
func DefaultParams() Params {

	return NewParams(DefaultPeriod)
}

// Validate validates all params
func (p Params) Validate() error {
	if err := validateMaxActivePeriod(p.MaxActivePeriod); err != nil {
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

/*
func Validate(i interface{}) error {
	v, ok := i.(Params)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.MaxActivePeriod <= 0 {
		return fmt.Errorf("maximum active period must be positive: %d", v.MaxActivePeriod)
	}

	return nil
}

// Params returns all of the  params
type Params struct {

	ActiveParams ActiveParams `json:"active_params" yaml:"active_params"`
}

func (p Params) String() string {
	return  p.Params.String()
}*/
// String implements the stringer interface for Params
func (p Params) String() string {
	out, err := yaml.Marshal(p)
	if err != nil {
		return ""
	}
	return string(out)
}
