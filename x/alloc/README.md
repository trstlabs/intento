<!--
order: 0
-->

# Allocation module

The Allocation module allocates the proportion of the inflation to the repective module accounts. 
It is currently distributed on a per block basis, but will be moved to an epoch basis (once this module is added to the SDK), and we will import it directly from there.

It allocates inflation through the following default params:

```golang
	Params{
		DistributionProportions: DistributionProportions{
			Staking:                     sdk.MustNewDecFromStr("0.55"),
			TrustlessContractIncentives: sdk.MustNewDecFromStr("0.25"),
			DeveloperRewards:            sdk.MustNewDecFromStr("0.05"),
			CommunityPool:               sdk.MustNewDecFromStr("0.15"),
		},
		WeightedDeveloperRewardsReceivers: []WeightedAddress{},
	}
```
