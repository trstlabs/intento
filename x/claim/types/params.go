package types

import (
	"time"
)

var (
	DefaultClaimDenom             = "uinto"
	DefaultDurationUntilDecay     = time.Hour
	DefaultDurationOfDecay        = time.Hour * 5
	DefaultDurationVestingPeriods = []time.Duration{time.Hour, time.Hour, time.Hour, time.Hour}
)
