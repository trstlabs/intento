---
order: 3
title: Claim
description: Useful information regarding allocation module
---
# Claim module

Recipients perform 4 actions
1. Perform RecurringSwap (DCA) strategy
2. Perform AutoSwap (DCA) strategy
3. Governance vote
4. Stake TRST

After each action, 20% of total elligable claims are unlocked. The remainder is unlocked following a vesting schedule of 4 vesting periods. The vesting periods end at different times per action.

Users must stake more than 67% of TRST received to submit a new claim for claimable tokens.


## Params
```golang

var (
	DefaultClaimDenom             = "utrst"
	DefaultDurationUntilDecay     = time.Hour
	DefaultDurationOfDecay        = time.Hour * 5
	DefaultDurationVestingPeriods = []time.Duration{time.Hour, time.Hour, time.Hour, time.Hour}
)
```