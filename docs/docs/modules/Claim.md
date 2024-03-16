---
order: 3
title: Claim
description: Useful information regarding allocation module
---
# Claim module

Recipients perform 4 actions
1. Perform automation of category 1
2. Perform automaton of category 2
3. Governance vote
4. Stake INTO

After each action, 20% of total elligable claims are unlocked. The remainder is unlocked following a vesting schedule of 4 vesting periods. The vesting periods end at different times per action.

Users must stake more than 67% of INTO received to submit a new claim for claimable tokens.


## Governance Parameters

```golang

var (
	DefaultClaimDenom             = "uinto"
	DefaultDurationUntilDecay     = time.Hour
	DefaultDurationOfDecay        = time.Hour * 5
	DefaultDurationVestingPeriods = []time.Duration{time.Hour, time.Hour, time.Hour, time.Hour}
)
```