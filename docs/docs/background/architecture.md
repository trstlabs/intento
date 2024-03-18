---
sidebar_position: 4
title: Technical overview
description: A technical overview of the Intento blockchain, including its architecture, consensus mechanism, block structure, network topology, automation solutions, Dymension Hub and Celestia.
---

# Technical overview

This section offers a detailed technical overview, the modular stack including Dymension Hub as a settlement layer and Celestia as a Data Availability (DA) layer. We'll cover the architecture, consensus mechanism, block structure, network topology, automation solutions, and other pertinent details.

##  Architecture

Intento, leveraging the flexible and modular Cosmos SDK, has expanded its interoperability and scalability through strategic integrations. By utilizing Dymension Hub and Celestia, Intento not only benefits from fast and secure consensus with CometBFT (previously Tendermint) but also enjoys a robust settlement and DA framework. 

Dymension Hub serves as a pivotal settlement layer, facilitating seamless and efficient transactions across multiple blockchains. Its Rollup Development Kit (RDK) enables Intento to efficiently build and deploy smart contracts and applications, enhancing the developer experience and broadening the blockchain's capabilities.

Celestia, on the other hand, acts as a dedicated Data Availability layer, ensuring that data related to transactions and contracts is reliably available and verifiable by anyone. This layer significantly bolsters the security and integrity of the blockchain, providing a trustless environment for users.

The adoption of eIBC (escrow Inter-Blockchain Communication) for connecting to Dymension Hub enables Intento to send and receive IBC packets to all chains linked to the Hub. This feature amplifies Intento’s interoperability, allowing assets and information to flow effortlessly across various blockchain ecosystems.

## Integrating with Dymension and Celestia

The strategic use of Dymension Hub and Celestia introduces a new dimension of scalability, interoperability, and data availability to Intento. These integrations not only enhance the blockchain's foundational aspects but also extend its reach and functionality across the broader blockchain ecosystem. Developers and users alike stand to benefit from a more efficient, secure, and versatile blockchain environment, setting a new standard for blockchain innovation and application.

![architecture](@site/docs/images/architecture.png)
*Figure 1: High-level overview of execution at the beginning of a new block.*


## On-chain Automation and Execution

Intento's on-chain automation framework has been meticulously designed to execute transactions based on defined schedules, leveraging the blockchain's inherent security. This framework, devoid of reliance on external agents or smart contracts, utilizes custom BeginBlocker functions for time-based executions. The integration with Dymension Hub and the use of IBC, alongside the Interchain Accounts standard, permit automated transactions across IBC-enabled chains without moving the assets out of the wallet or chain.

Intento’s execution mechanism queues triggers and contracts, checking them at the beginning of each block for their scheduled execution time. In the event of a blockchain halt, the system is designed to resume queued executions in subsequent blocks, ensuring reliability and continuity.

## CometBFT and Time Management

CometBFT, with its proposer-based timestamp mechanism, ensures a consistent and secure timestamping system for block creation. This approach mitigates risks associated with inaccurate timestamps, maintaining the blockchain's integrity. The adoption of precision and delay parameters among validators facilitates a synchronized agreement on the block timestamps, crucial for the orderly function of the blockchain.

By integrating these advanced layers, Intento solidifies its position as a forward-thinking blockchain, ready to tackle the demands of modern decentralized applications and systems.