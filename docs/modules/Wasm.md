---
order: 1
title: Wasm
description: Useful information regarding compute module
---

# Wasm Module

This is a brief overview of the functionality of Trustless Hub's CosmWasm engine. It is optimized to support auto execution and privacy by default in contracts. Next to the default private state, contracts have a public state. This state is easily accessable through RPC endpoints.
The Wasm module is suppports CosmWasm v1 contracts optimized for Trustless Hub.

## Auto Execution Fees

Auto execution messages have a fee. Governance can update the fee.


The fee consists of:

1. Gas-dependent FlexFee (Goes to block proposer)

2. Constant Fee (Goes to community pool)

3. Commission (Goes to community pool)

These parameters can be updated by governance. These fees may differ between 1-time execution and recurring execution.

The chain directs a flexible portion of these to validators called FexFee. With this part of the fee, validators are incentivized to include auto execution in their blocks over general transactions. This ensures you that auto execution will happen at pre-defined times.

The gas-dependent fee can be increased and decreased with the AutoMsgFlexFeeMul governance param and the Constant Fee is altered through the AutoMsgConstantFee and  RecurringAutoMsgConstantFee params.

Upon instantiaiton auto execution fees are transferred to the contract. The instantiator gets the contract balance refunded automatically after execution in case a contract owner is set. 

At the end of each block, the BeginBlock checks if there are contracts that are to be executed. These contracts can be incentivized.
MinContractDurationForIncentive,MaxContractIncentive, MinContractBalanceForIncentive can be adjusted by governane to ensure that the contract incentives are distributed in a fair manner. These incentives can not exceed the total cost of auto execution.


## Parameters

Through governance, a number of compute parameters can be adjusted. The default values are the following:
```golang

const (
	// AutoMsgFundsCommission percentage to distribute to community pool for leftover balances (rounded up)
	DefaultAutoMsgFundsCommission int64 = 2

	// AutoMsgConstantFee fee to prevent spam of auto messages, to be distributed to community pool
	DefaultAutoMsgConstantFee int64 = 1000000 // 1utrst

	// AutoMsgFlexFeeMul is a multiplier for the gas-dependent flex fee to prevent spam of auto messages, to be distributed to community pool
	DefaultAutoMsgFlexFeeMul int64 = 100

	// RecurringAutoMsgConstantFee fee to prevent spam of auto messages, to be distributed to community pool
	DefaultRecurringAutoMsgConstantFee int64 = 1000000 // 1utrst

	// Default max period for a contract that is self-executing
	DefaultMaxContractDuration time.Duration = time.Hour * 24 * 366 // 366 days
	// MinContractDuration sets the minimum duration for a self-executing contract
	DefaultMinContractDuration time.Duration = time.Second * 45
	// MinContractInterval sets the minimum interval self-execution
	DefaultMinContractInterval time.Duration = time.Second * 20
	// MinContractDurationForIncentive to distribute reward to contracts we want to incentivize
	DefaultMinContractDurationForIncentive time.Duration = time.Hour * 24 // time.Hour * 24 // 1 day

	// DefaultMaxContractIncentive max amount of utrst coins to give to a contract as incentive
	DefaultMaxContractIncentive int64 = 500000000 // 500utrst

	// MinContractBalanceForIncentive minimum balance required to be elligable for an incentive
	DefaultMinContractBalanceForIncentive int64 = 50000000 // 50utrst
)
```

## Configuration

You can add the following section to `config/app.toml`. Below is shown with defaults:

```toml
[wasm]
# This is the maximum sdk gas (wasm and storage) that we allow for any x/compute "smart" queries
query_gas_limit = 300000
# This is the number of wasm vm instances we keep cached in memory for speed-up
# Warning: this is currently unstable and may lead to crashes, best to keep for 0 unless testing locally
lru_size = 0
```

## Events

A number of events are returned to allow good indexing of the transactions from smart contracts.

Every call to Instantiate or Execute will be tagged with the info on the contract that was executed and who executed it.
It should look something like this (with different addresses). The module is always `wasm`, and `code_id` is only present
when Instantiating a contract, so you can subscribe to new instances, it is omitted on Execute. There is also an `action` tag
which is auto-added by the Cosmos SDK and has a value of either `store-code`, `instantiate` or `execute` depending on which message
was sent:

```json
{
  "Type": "message",
  "Attr": [
    {
      "key": "module",
      "value": "wasm"
    },
    {
      "key": "action",
      "value": "instantiate"
    },
    {
      "key": "signer",
      "value": "trust1vx8knpllrj7n963p9ttd80w47kpacrhuts497x"
    },
    {
      "key": "code_id",
      "value": "1"
    },
    {
      "key": "contract_address",
      "value": "trust18vd8fpwxzck93qlwghaj6arh4p7c5n894lxvdh"
    }
  ]
}
```

Finally, the contract itself can emit custom events.
We add a `contract_address` attribute that contains the actual contract that emitted that event.
Event logs are private by default.
Here is an example from the escrow contract successfully releasing funds to the destination address:

```json
{
  "Type": "wasm",
  "Attr": [
    {
      "key": "contract_address",
      "value": "trust18vd8fpwxzck93qlwghaj6arh4p7c5n894lxvdh"
    },
    {
      "key": "action",
      "value": "release"
    },
    {
      "key": "destination",
      "value": "trust14k7v7ms4jxkk2etmg9gljxjm4ru3qjdugfsflq"
    }
  ]
}
```
