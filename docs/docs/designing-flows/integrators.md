---
title: Integrators
sidebar_position: 1
description: Empower your users with smart flows that enhance security, efficiency, and flexibility
---

Imagine a world where your blockchain actions are seamless, automated, and entirely under your control. No more manual rebalancing, forgotten governance votes, or missed staking rewards. Whether you’re managing liquidity on a DEX, automating wallet functions, or streamlining governance processes, you need a system that works for you—without constant oversight.

This guide is for integrators who want to leverage interchain execution to bring efficiency and innovation to their users.
By designing smart execution flows, you empower your users with flows that enhance security, efficiency, and flexibility while keeping them in control.

#### Differences in Interchain Execution

Integrators can choose from multiple interchain execution options based on their requirements:

- **Self-Hosted Interchain Accounts:** Fully self-custodial approach where integrators maintain complete execution control.
- **Trustless Agent:** Allows execution of interchain actions while managing fees.
- **Scalable Managed Solution:** A streamlined, managed approach for executing interchain actions at scale.

#### Message Types

Understanding message types is crucial for designing effective execution flows:

- **Local Messages:** Standard messages sent directly by users.
- **Authz MsgExec Messages:** Enables one account to execute messages on behalf of another.
- **ICA (Interchain Accounts) Messages:** Messages sent via IBC (Inter-Blockchain Communication) through an interchain account.

##### Hosted vs. Self-Hosted ICA Messages

- **Hosted ICA Messages** can only call `MsgExec`.
- **Self-Hosted ICA Messages** provide full execution control.

#### Using Conditions in Flows

Conditions help define execution logic and handle different scenarios:

- **Stoplights:** Stop execution on failure to prevent unwanted state changes.
- **Comparisons:** Execute actions based on specific conditions or thresholds.
- **Feedback Loops:** Continuously adapt execution based on prior results.
- **Interchain Queries (ICQ):** Query a blockchain's state to make informed execution decisions.
