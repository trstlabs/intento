---
sidebar_position: 1
title: Governance-based Fees
description: Fee paremeters for Intent-Based Flows
---

The Intent module enables flows to be highly configurable and conditional, whereby the flows can depend on execution results. For example, a protocol user could set up a sequence of flows such as swapping ATOM for USDC on Osmosis and then paying for a subscription that is settled on Ethereum using the Axelar bridge with General Message Passing. By enabling these user intents, protocols and their end-users can automate complex workflows in a seamless manner.

### Gas operations

Based on lines of code, we expect the module to be using 80,000 to 100,000 gas for triggering up to 10 executions. Significantly less than via bots and custom smart contracts (700,000-1,000,000) whilst bringing trust assumptions to the minimum. This makes the module highly scalable for any specified intent. To manage network congestion, make the chain scalable, and provide value for token holders, protocol fees can be set and adjusted over time by token holders via chain governance.

## Flow Governance Proposal Parameters

The Flow Governance Proposal Parameters define the rules and constraints governing the execution of Flows within the network. These parameters are unique as they are set through on-chain governance, ensuring transparency and adaptability over time. By optimizing key economic and operational aspects, these parameters allow the network to scale efficiently without becoming congested.

A notable feature of these parameters is the ability to use multiple tokens beyond the native denomination for transaction fees. This enhances user experience and aligns incentives with ATOM, promoting broader ecosystem participation and interoperability.

| Parameter               | Type                       | Description                                                                                                                                             | Example Value                          |
| ----------------------- | -------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------- |
| `flow_funds_commission` | `int64`                    | Commission rate to return remaining funds in flow fee account after final execution (e.g., 100 = 1X, 250 = 250)                                         | `2` (2%)                               |
| `flow_flex_fee_mul`     | `int64`                    | Multiplier to the flexible flow gas fee (e.g., 100 = 1X, 250 = 250)                                                                                     | `250` (2.5X)                           |
| `burn_fee_per_msg`      | `int64`                    | Fixed burn fee per message execution to burn native denom                                                                                               | `10_000` (0.01uinto)                   |
| `gas_fee_coins`         | `repeated Coin`            | Array of denoms that can be used for fee payment together with an amount                                                                                | `[1uinto, 0.05ibc/chain_channel_hash]` |
| `max_flow_duration`     | `google.protobuf.Duration` | Maximum period for self-executing Flow                                                                                                                  | `263520h` (a little over 3 years)      |
| `min_flow_duration`     | `google.protobuf.Duration` | Minimum period for self-executing Flow                                                                                                                  | `1m` (1 minute)                        |
| `min_flow_interval`     | `google.protobuf.Duration` | Minimum interval for self-executing Flow                                                                                                                | `1m` (1 minute)                        |
| `relayer_rewards`       | `repeated int64`           | Relayer rewards in uinto for each message type (0 = SDK, 1 = Wasm, 2 = Osmo). Rewards are in uinto and topped up in the module account by alloc module. | `[10_000, 15_000, 18_000, 22_000]`     |
