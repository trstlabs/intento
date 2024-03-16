# Claims module

## Abstract

The INTO claims module has users begin with a portion of their total airdrop allocation,
and then be able to automatically claim higher percentages as they perform certain tasks on-chain.
Furthermore, these claimable assets 'expire' if not claimed.
Users have two months (`DurationUntilDecay`) to claim their full airdrop amount.
After two months, the reward amount available will decline over 4 months (`DurationOfDecay`) in real time, until it hits `0%` at 6 months from launch (`DurationUntilDecay + DurationOfDecay`).

After 6 months from launch, all unclaimed tokens get sent to the community pool.

## Contents

1. **[Concept](01_concepts.md)**
2. **[State](02_state.md)**
3. **[Events](03_events.md)**
4. **[Keeper](04_keeper.md)**  
5. **[React Hooks](05_react_hooks.md)**  
6. **[Queries](06_queries.md)**  
7. **[Params](07_params.md)**

## Genesis State

## Actions

All accounts start out with 20% of their entire airdrop allocation.

There are 4 types of actions, each of which release another 20% of the airdrop allocation.
The 4 actions are as follows:

```golang
ActionComputeInstantiate Action = 0
ActionComputeExecute          Action = 1
ActionInterchainAccountCall          Action = 2 (TBD)
ActionDelegateStake Action = 3
```


Performing each of these activities will unlock 20% of the claimable amount for that action. The remaining vests according to a vesting schedule, and will be released in 4 phases, which each unlock 20%. This process will ensure that airdrop recipients are active users of the Hub, and that liquidity does not automatically flood the market. By performing these specific actions, users will get in touch with the 4 pillars that make Trustless Contract unique.

These actions are monitored by registring claim **hooks** to the governance, staking, module, and lockup modules.
This means that when you perform an action, the claims module will immediately unlock those coins if they are applicable.
These actions can be performed in any order.

The code is structured by separating out a segment of the tokens as "claimable", indexed by each action type.
So if Alice delegates tokens, the claims module will move the 20% of the claimables associated with staking to her liquid balance.
If she delegates again, there will not be additional tokens given, as the relevant action has already been performed.
Every action must be performed to claim the full amount.

## ClaimRecords

A claim record is a struct that contains data about the claims process of each airdrop recipient.

It contains an address, the initial claimable airdrop amount, and an array of bools representing 
whether each action has been completed. The position in the array refers to enum number of the action.

So for example, `[false, true, true, false]` means that `ActionComputeInstantiate` and `ActionComputeExecute` are completed.

```golang
type ClaimRecord struct {
    // address of claim user
    Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty" yaml:"address"`

    // total initial claimable amount for the user
    InitialClaimableAmount sdk.Coins
    
    // true if action is completed
    // index of bool in array refers to action enum #
    ActionCompleted []bool
}

```


## Params

The airdrop logic has 4 parameters:

```golang
type Params struct {
    // Time that marks the beginning of the airdrop disbursal,
    // should be set to chain launch time.
    AirdropStartTime   time.Time
    DurationUntilDecay time.Duration
    DurationOfDecay    time.Duration
    // denom of claimable asset
    ClaimDenom string
}
```
