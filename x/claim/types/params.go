package types

import (
	fmt "fmt"
	"time"
)

var (
	DefaultClaimDenom             = "uinto"
	DefaultDurationUntilDecay     = time.Hour
	DefaultDurationOfDecay        = time.Hour * 5
	DefaultDurationVestingPeriods = []time.Duration{time.Hour, time.Hour, time.Hour, time.Hour}
)

// Validate checks that the Params fields are valid
func (p Params) Validate() error {
	if p.ClaimDenom == "" {
		return fmt.Errorf("claim denom cannot be empty")
	}
	if p.DurationUntilDecay < 0 {
		return fmt.Errorf("duration until decay cannot be negative")
	}
	if p.DurationOfDecay < 0 {
		return fmt.Errorf("duration of decay cannot be negative")
	}
	if len(p.DurationVestingPeriods) == 0 {
		return fmt.Errorf("duration vesting periods cannot be empty")
	}
	for i, d := range p.DurationVestingPeriods {
		if d <= 0 {
			return fmt.Errorf("duration vesting period at index %d must be positive", i)
		}
	}
	return nil
}
