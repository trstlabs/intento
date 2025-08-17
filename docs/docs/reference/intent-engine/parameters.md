---
title: Engine Parameters
description: Configuration parameters and fee settings for the Intent Engine
sidebar_position: 1
sidebar_class_name: sidebar-item-reference
---

The Intent engine enables highly configurable intent-based flows. Conditional, whereby flows can depend on execution results. The module typically uses just 80,000 to 100,000 gas for triggering up to 10 executions. Significantly less than via bots and custom smart contracts (700,000-1,000,000) whilst bringing trust assumptions to the minimum. This makes the engine highly scalable, capable of handling millions of flows per hour. To manage network congestion, ensure the chain stays scalable, and provide value for token holders, fees are set and adjusted over time by token holders via chain governance.

## Flow Governance Parameters

The Flow Governance Parameters define the rules and constraints governing the execution of Flows within the network. These parameters are unique as they are set through on-chain governance, ensuring transparency and adaptability over time. By optimizing key economic and operational aspects, these parameters allow the network to scale efficiently without becoming congested.

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
| `connection_relayer_rewards`       | `repeated ConnectionRelayerReward`           | 'Per connection, relayer rewards in uinto for each message type (0 = Low Gas, 1 = Medium Gas, 2 = High Gas, 3 = Authz). Rewards are in uinto and topped up in the module account by alloc module. | `[{connection_id: "connection-1", relayer_rewards: [10_000, 15_000, 18_000, 22_000]}]`     |
