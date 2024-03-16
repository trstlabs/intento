---
order: 1
title: AutoTX
description: Automation with the interchain accounts module
---

# AutoTX Module

This module is used to automate processes across chains. With 'MsgSubmitAutoTx', time-based interchain account message calls are scheduled.

This is useful for Cosmos and smart contract developers wanting to automate message calls without using hard-to-integrate and unreliable off-chain bots.

Use cases include payroll, dollar-cost averaging, managing contract workflows. Respective messages for these actions are: MsgSend on any Cosmos SDK chain; MsgSwapExactAmountIn on Osmosis and MsgExecuteContract on CosmWasm chains.

Local messages to this chain can also be scheduled. This can be used to Autocompound tokens to any validator, stream local INTO tokens or even stream IBC Transfers.

## Automating user assets

The trustless nature of on-chain automation and trustless bridging through IBC have way better trust assumptions compared to bots. What it also brings to the table is asserting whom the message caller will be. By knowing the message caller it is possible to authorize the message caller to automate user assets. This is done by creating an grant in the Authz module of host chain. This allows for payroll services without pre-funding a contract & transfers of assets based on time. This is only available using on-chain automation.

### Fixed Fee

The Fixed Fee is applied per message basis. An AutoTX can have multple messages to be executed, with a maximum of 9. This number can be increased when desired.

### FlexFee Fee

The FlexFee is timedependent fee calculated for each AutoTX period on a minute basis. 

### Fund deduction

When funds are sent along MsgSubmitAutoTx, fees are deducted from the sender account. When specified, fees are deducted from a fund address. When funds are unavailable, execution will not take place. Funds left on a fund address are automatically returned to the owner after the last execution. A small commision of 2-5% is taken to the community pool.

## Relayer Rewards

Relayers are vital to well-functioning of this module for ensuring timely, reliable execution. To incentivize relayers, rewards are minted in the mint module, allocated from the alloc module and sent to the AutoTX module. For acknoledging a succesfull IBC packet containing AutoTX messages, relayers are incentivized. Relayer rewards are specified based on the category of message. The category are as follows: SDK message, WASM message and Osmosis message. AuthZ messages can perform authorized actions on behalf of a user such as recurring transactions and reward claims. On Osmosis users can perform DCA strategies and withdrawal automatically after an unbonding period ends.  WASM smart contract calls are for developers automating their dApps and users that want to automate their smart contract tasks.

## Parameters

A number of automation-related parameters can be adjusted. Parameters can be adjusted by governance to ensure that fees and rewards are fair. The default values are the following:

```golang
const (
 // AutoTXFundsCommission percentage to distribute to community pool for leftover balances (rounded up)
 DefaultAutoTxFundsCommission int64 = 2 //2%
 // AutoTXConstantFee fee to prevent spam of auto messages, to be distributed to community pool
 DefaultAutoTxConstantFee int64 = 5_000 // 0.005trst
 // AutoTXFlexFeeMul is the denominator for the gas-dependent flex fee to prioritize auto messages in the block, to be distributed to validators
 DefaultAutoTxFlexFeeMul int64 = 3 // 3% of minutes for a given period as uinto (1_000m = 20uinto)
 // RecurringAutoTxConstantFee fee to prevent spam of auto messages, to be distributed to community pool
 DefaultRecurringAutoTxConstantFee int64 = 5_000 // 0.005trst
 // Default max period for a AutoTX that is self-executing
 DefaultMaxAutoTXDuration time.Duration = time.Hour * 24 * 366 * 2 // a little over 2 years
 // MinAutoTXDuration sets the minimum duration for a self-executing AutoTX
 DefaultMinAutoTXDuration time.Duration = time.Second * 60
 // MinAutoTXInterval sets the minimum interval self-execution
 DefaultMinAutoTXInterval time.Duration = time.Second * 60
 // DefaultRelayerReward for a given AutoTX type
 DefaultRelayerReward int64 = 10_000 //0.01trst
)

## BeginBlocker

At the beginning of each block, the BeginBlocker checks if there are AutoTXs that are set for automation. The BeginBlocker is used as it can best proxy the execute time set at MsgSubmitAutoTx. 