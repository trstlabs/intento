---
order: 4
title: Allocation
description: Useful information regarding allocation module
---
# Allocation module

The Allocation module allocates the proportion of the inflation to the repective module accounts. 
Inflation is distributed on a per block basis.

It allocates inflation through governance parameters, like so:

```golang
	Params{
		DistributionProportions: DistributionProportions{
			Staking:                     sdk.MustNewDecFromStr("0.60"),
			RelayerIncentives:			 sdk.MustNewDecFromStr("0.10"),
			ContributorRewards:          sdk.MustNewDecFromStr("0.05"),
			CommunityPool:               sdk.MustNewDecFromStr("0.15"),
		},
		WeightedContributorRewardsReceivers: []WeightedAddress{},
	}
```


