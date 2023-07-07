<!--
order: 0
-->

# Allocation module

The Allocation module allocates the proportion of the inflation to the repective module accounts. 
It is currently distributed on a per block basis.

It allocates inflation through the following default params:

```golang
	Params{
		DistributionProportions: DistributionProportions{
			Staking:                     sdk.MustNewDecFromStr("0.60"),
			RelayerIncentives: 			 sdk.MustNewDecFromStr("0.10"),
			ContributorRewards:          sdk.MustNewDecFromStr("0.05"),
			CommunityPool:               sdk.MustNewDecFromStr("0.15"),
		},
		WeightedContributorRewardsReceivers: []WeightedAddress{},
	}
```
