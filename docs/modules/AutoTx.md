---
order: 4
title: AutoIbcTx
description: Automation with the interchain accounts module
---

# AutoIbcTx Module

This module contains a custom Interchain Accounts authentication module
With MsgSubmitAutoTx time-based interchain account message calls are scheduled.

This is useful for Cosmos and smart contract developers needing to automate message calls without using hard-to integrate off-chain bots

Examples include payroll (MsgSend on any chain), dollar-cost averaging (MsgSwapExactAmountIn on Osmosis), managing contract workflows (MsgExecuteContract on CosmWasm chains)

## Automate user assets

The trustless nature of on-chain automation and trustless bridging through IBC allows for a great deal of reliability. What it also brings to the table is asserting whom the message caller will be. By knowing the message caller, it is possible to create a grant and authorize the message caller to use a user's assets.  This allows for payroll services without pre-funding a contract & transfers of assets based on time. This is only available with on-chain automation.

## Fees

To encourage valuable transactions and timely execution, fees are implemented. These are designed to prevent spam, network congestion and reward validators for providing computational resources.

1. FixedFee to the community pool.  
2. FlexFee to the proposing validator


Governance can decide on how these fees are set. The FixedFee goes to the community pool to create network growth. Funds can be allocated towards ecosystem development, relayer participation, new tooling and educational content.

When funds are sent along MsgSubmitAutoTx, fees are deducted from a unique AutoTx address. If unspecified, fees are deducted from the sender account. When funds are unavailable, execution will not take place.

## Relayer Rewards

Relayers are vital to well-functioning of this module for ensuring timely, reliable execution. To incentivize relayers, rewards are minted in the mint module, allocated from the alloc module and sent to the AutoTx module. Here, relayers are incentivized. Relayer rewards are specified to different types of messages. AuthZ messages that perform authorized actions on behalf of a user such as recurring transactions and DCA strategies and WASM smart contract calls for developers automating their dApps. 

## Parameters

A number of automation-related parameters can be adjusted. Parameters can be adjusted by governance to ensure that the fees are fair. The default values are the following:

```golang

// AutoTxFundsCommission percentage to distribute to community pool for leftover balances (rounded up)

DefaultAutoTxFundsCommission int64 = 2

// AutoTxConstantFee fee to prevent spam of auto messages, to be distributed to community pool

DefaultAutoTxConstantFee int64 = 1_000_000  // 1trst

// AutoTxFlexFeeMul is the denominator for the gas-dependent flex fee to prioritize auto messages in the block, to be distributed to validators

DefaultAutoTxFlexFeeMul int64 = 100  // 100/100 = 1 = gasUsed

// RecurringAutoTxConstantFee fee to prevent spam of auto messages, to be distributed to community pool

DefaultRecurringAutoTxConstantFee int64 = 1_000_000  // 1trst

// Default max period for a AutoTx that is self-executing

DefaultMaxAutoTxDuration time.Duration = time.Hour * 24 * 366 * 10  // a little over 10 years

// MinAutoTxDuration sets the minimum duration for a self-executing AutoTx

DefaultMinAutoTxDuration time.Duration = time.Second * 60

// MinAutoTxInterval sets the minimum interval self-execution

DefaultMinAutoTxInterval time.Duration = time.Second * 60

// DefaultRelayerReward for a given autotx type (0=SDK message, 1=WASM message, 2=Osmosis message).

DefaultRelayerRewards []int64 = []int64{10_000, 15_000, 18_000}

```

## BeginBlocker

At the beginning of each block, the BeginBlocker checks if there are AutoTxs that are set for automation. The BeginBlocker is used as it can best proxy the execute time set at MsgSubmitAutoTx. 