---
sidebar_position: 1
title: Technical overview
description: A technical overview of the Intento blockchain, including its architecture, consensus mechanism, block structure, network topology, automation solutions, and data availability.
---


Intento's intent-based flow framework has been meticulously designed to execute transactions based on defined schedules, leveraging the blockchain's inherent security. This framework, devoid of reliance on external agents or smart contracts, utilizes custom BeginBlocker functions for flow executions. The integration with the IBC Interchain Accounts standard, permit the Intento chain to execute transactions across IBC-enabled chains without moving the assets out of the user's control.

![architecture](@site/docs/images/architecture.png)

Intento’s execution mechanism queues flows, checking them at the beginning of each block for their scheduled execution time. In the event of a blockchain halt, the system is designed to resume queued executions in subsequent blocks, ensuring reliability and continuity.
With Intento you can use the power of IBC for your user intents. You can use Interchain Queries (ICQ) and use their responses for comparisons and build feedback loops. Or use Interchain Accounts (ICA) to execute actions on connected chains. Below are just some of the examples of how flows can look like.

![Example flows](@site/docs/images/example_flows.png)

## Leveraging Interchain Security

Intento is designed to provide a decentralized, scalable, and secure environment for executing intent-based action flows across multiple blockchains. A key component of Intento’s architecture is **Interchain Security**, which enables the protocol to inherit security from established validator sets while maintaining its autonomy and flexibility. ICS allows Intento to optimize security, decentralization, and network efficiency as it scales.

### Why ICS?

ICS is integral to Intento’s architecture for several reasons:

1. **Decentralized Security** – By leveraging ICS, Intento benefits from the Cosmos Hub’s highly decentralized validator set, reducing the risk of centralization and single points of failure.
2. **Economic Security** – Rather than bootstrapping its own validator network, Intento inherits the economic security of the Cosmos Hub’s staked ATOM, making attacks more costly and difficult.
3. **Scalability & Modularity** – ICS allows Intento to focus on execution and orchestration of action flows while relying on an external, secure validator set for consensus.
4. **Permissionless Growth** – With ICS v2 introducing partial set security, Intento can progressively decentralize its validator set while maintaining a robust security foundation.

### How Intento Uses ICS

Intento operates as a consumer chain under ICS, meaning its consensus mechanism is secured by Cosmos Hub validators who validate Intento’s transactions in exchange for rewards. This setup ensures:

- A high level of security from day one.
- A streamlined consensus mechanism without requiring Intento to maintain its own validator set.
- Efficient block finalization and execution of action flows.

#### Action Flows and ICS

Intento enables intent-based action execution across IBC-connected chains, which requires highly reliable and timely processing. By utilizing ICS, Intento ensures:

- **Secure orchestration** of workflows such as cross-chain portfolio management, autocompounding staking rewards, and scheduled transactions.
- **Minimal latency** in processing conditional triggers, governance actions, and smart contract executions.
- **Resilient infrastructure** that can handle high volumes of interchain queries and state updates efficiently.

## CometBFT and Time Management

CometBFT, with its proposer-based timestamp mechanism, ensures a consistent and secure timestamping system for block creation. This approach mitigates risks associated with inaccurate timestamps, maintaining the blockchain's integrity. The adoption of precision and delay parameters among validators facilitates a synchronized agreement on the block timestamps, crucial for the orderly function of the blockchain.

## Conclusion

Intento’s architecture enables secure, scalable, and efficient execution of decentralized workflows. By integrating ICS and IBC, Intento provides a next-generation solution for cross-chain orchestration while maintaining self-custodial security. Intento is set to scale alongside the Cosmos ecosystem, ensuring a robust and future-proof infrastructure for intent-based action flows.
