---
sidebar_position: 4
title: Technical overview
description: A technical overview of the Intento blockchain, including its architecture, consensus mechanism, block structure, network topology, automation solutions, and data availability.
---

# Technical overview

This section offers a detailed technical overview, the modular stack including Celestia as a Data Availability (DA) layer. We'll cover the architecture, consensus mechanism, block structure, network topology, automation solutions, and other details.

## Architecture

Intento, leveraging the flexible and modular Cosmos SDK, has expanded its interoperability and scalability through strategic integrations
Celestia, acts as a dedicated Data Availability layer, ensuring that data related to transactions and contracts is reliably available and verifiable by anyone. This layer significantly bolsters the security and integrity of the blockchain, providing a trustless environment for users.

![architecture](@site/docs/images/architecture.png)
_Figure 1: High-level overview of execution at the beginning of a new block._

## Non-custodial execution

Intento's flow framework has been meticulously designed to execute transactions based on defined schedules, leveraging the blockchain's inherent security. This framework, devoid of reliance on external agents or smart contracts, utilizes custom BeginBlocker functions for time-based executions. The integration with the IBC Interchain Accounts standard, permit the Intento chain to execute transactions across IBC-enabled chains without moving the assets out of the user's control

Intento’s execution mechanism queues triggers and contracts, checking them at the beginning of each block for their scheduled execution time. In the event of a blockchain halt, the system is designed to resume queued executions in subsequent blocks, ensuring reliability and continuity.

## CometBFT and Time Management

CometBFT, with its proposer-based timestamp mechanism, ensures a consistent and secure timestamping system for block creation. This approach mitigates risks associated with inaccurate timestamps, maintaining the blockchain's integrity. The adoption of precision and delay parameters among validators facilitates a synchronized agreement on the block timestamps, crucial for the orderly function of the blockchain.

