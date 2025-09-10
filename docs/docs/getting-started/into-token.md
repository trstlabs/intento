---
title: The INTO Token
sidebar_position: 6
description: The INTO Token - Documentation & Whitepaper
---
# INTO Token - Documentation & Whitepaper

## Overview

INTO is the native token of the **Intento Network**, a Cosmos-based Layer 1 designed for **intent-based orchestration** across the interchain. It enables secure, programmable flows that let users and agents execute complex, self-custodial actions across chains.

INTO is not just a fee token. It is the **coordination layer** that aligns validators, relayers, builders, and users through staking, governance, and a deflationary execution model.

> For full details about the INTO token, see our [Introducing INTO: The Token Powering Interchain Flows](https://intento.zone/post/introducing-into-the-token-powering-interchain-flows/) blog post.


## Network Requirements

To operate an intent-based L1, the network must provide:

* **Security** — Validators staking $INTO to secure consensus.
* **Governance** — Token holders guiding upgrades, parameters, and community allocations.
* **Incentives** — Stakers, relayers, and contributors aligned through emissions and fees.
* **Fee System** — A mechanism to meter execution of interchain flows.

$INTO fulfills all four roles.


## Token Utility

### 1. Staking

* Secure the network by delegating to validators.
* Earn staking rewards.
* Participate in shared protocol revenue.

### 2. Flow Execution

* Pay execution fees for interchain flows.
* Save on fees when paying directly in $INTO.
* A portion of each $INTO-paid fee is **burned**, creating a deflationary loop.
* Unique **wallet fallback** feature: fees can be pulled from a user’s balance without pre-deposits.

### 3. Governance

Holders can propose and vote on:

* Fee parameters and exemptions.
* Incentive distribution (stakers, relayers, builders).
* Upgrades to flow functionality and integrations.
* Spending from the community pool.


## Flow Fee Model

* **Charged per message and execution step.**
* Complex flows cost more; simple flows cost less.
* **Condition checks and queries are free** — flows only execute when conditions are met.

Fee routing:

* **$INTO payments:** part is burned, part to community pool.
* **Non-$INTO tokens (ATOM, OSMO, ELYS):** routed entirely to treasury.


## Token Supply

### Initial Allocation (Genesis Supply: 400M $INTO)

| Category          | Amount (M) | % of Supply |
| ----------------- | ---------- | ----------- |
| Airdrop           | 90         | 22.5%       |
| Team              | 90         | 22.5%       |
| Community Pool    | 82         | 20.5%       |
| Strategic Reserve | 70         | 17.5%       |
| StreamSwap Event  | 40         | 10.0%       |
| Grant Program     | 22         | 5.5%        |
| Testnet Program   | 6          | 1.5%        |
| **Total**         | **400M**   | **100%**    |

**Liquid at Launch:** StreamSwap Event, Airdrop, and Testnet allocations.
**Long-Term Vesting:** Team, Reserve, and Grant Program.


## Vesting & Unlocks

* **Team:** Continuous vesting, 5 years.
* **Strategic Reserve:** 50% liquid at launch; 50% vests linearly over 5 years.
* **Grant Program:** Continuous vesting, 5 years.
* **Airdrop:** Unlocked via 4 staged actions over days/weeks; clawback for unclaimed.

This avoids supply shocks and ties distribution to participation.


## Emission Schedule

* **Inflation:** Starts at 10%, reduced 25% per year → near-zero after year 10.
* **Total Supply Growth:** 400M → \~559M over 20 years.
* **Distribution of New Emissions:**

  * 70% Community Pool.
  * 25% Stakers (ATOM + INTO).
  * 5% Relayers.

| Year | Supply (M) | Annual Inflation |
| ---- | ---------- | ---------------- |
| 0    | 400        | –                |
| 1    | 440        | 10.0%            |
| 3    | 470        | 6.8%             |
| 5    | 509        | 3.4%             |
| 10   | 548        | 0.7%             |
| 20   | 559        | 0.0%             |


## Airdrop Design

* **22.5% (90M)** allocated to airdrop.
* Claimable via the Intento Portal.
* **Clawback Model:** unclaimed tokens return to treasury.
* **Claim Unlock:** Users must complete flows and stake tokens to access allocation.

| Claim Rate | Tokens Distributed | Clawed Back |
| ---------- | ------------------ | ----------- |
| 20%        | 18M                | 72M         |
| 50%        | 45M                | 45M         |
| 80%        | 72M                | 18M         |

This structure ensures distribution only to active, aligned users.

### Unlock Mechanism

Recipients must complete **four meaningful on-chain actions** to unlock their allocation:

1. Orchestrate a flow on Intento.
2. Orchestrate a flow over IBC.
3. Stake tokens.
4. Participate in governance.

**Unlock model:**

* Each action → unlocks **20%** of the airdrop portion.
* Remaining portion vests over several days.
* Claiming requires staking at least **67%** of unlocked tokens.

This is not a “click-claim” airdrop. It ensures alignment through **participation and staking**.

### Decay Model

To prevent idle supply, the airdrop includes a **time-based decay mechanism**:

* **Grace Period (DurationUntilDecay):** 4 weeks after claim eligibility.

  * No penalties during this time.
* **Decay Period (DurationOfDecay):** 8 weeks after the grace period.

  * Linear reduction: unclaimed allocation decreases to 0 by the end of the period.

**Example:** If a participant waits until halfway through the decay period (\~8 weeks after launch), only \~50% of their unclaimed allocation remains.

Any tokens left unclaimed after full decay return to the **Community Pool**.

### Clawback

Unclaimed or unused allocations are not left idle. They are **clawed back** into the Community Pool to support:

* Builder grants.
* User incentives.
* Ecosystem growth campaigns.

This ensures that all $INTO either strengthens the protocol or aligns with active participants.


## Deflationary Mechanism

Every executed flow strengthens the system:

* A portion of $INTO is **burned per message**, permanently reducing supply.
* Non-$INTO fees accumulate in the treasury, governed by stakers.
* Increased flow usage → higher burn rate → tighter supply.


## Alignment Model

* **Validators:** Secured by INTO stake.
* **Relayers:** Incentivized with 5% emissions.
* **Builders:** Supported via community pool and grants.
* **Users:** Save fees and gain governance rights by holding $INTO.

No VC unlocks. No cliffs. Continuous vesting for core contributors.

## Conclusion

$INTO is the backbone of Intento: a token designed to coordinate intent execution across chains.

* **Programmable fees, deflationary design.**
* **Multi-chain execution, AI-ready orchestration.**
* **Built for builders, users, validators, and relayers alike.**

Every flow strengthens the network. Every burn aligns supply with real usage.
$INTO is not speculative fuel — it is powering the **coordination layer** for decentralized intent execution.
