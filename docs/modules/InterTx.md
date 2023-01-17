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

Fees are twofold, to pay validators for their services and to distribute funds back to the community.

1. FlexFee, goes to validators and their delegators for providing computational resources
2. FixedFee, goes to the community pool.  Governance can decide on how these can be used to incentivize network partners such as relayer operators, event & hackathon organizers, as well as automation tooling.

## Fee Parameters

A number of automation parameters can be adjusted. Parameters can be adjusted by governance to ensure that the fees are fair. The default values are the following:

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

DefaultMinAutoTxDuration time.Duration = time.Second * 40

// MinAutoTxInterval sets the minimum interval self-execution

DefaultMinAutoTxInterval time.Duration = time.Second * 20

```

## BeginBlocker

At the beginning of each block, the BeginBlocker checks if there are AutoTxs that are set for automation. The BeginBlocker is used as it can best proxy the execute time set at SubmitAutoTx