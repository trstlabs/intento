---
title: Differentiation
sidebar_position: 2
---

# Differentiation

Currently, performing intent-based actions requires a combination of smart contracts and off-chain bot networks.

For both bot networks and privileged smart contracts, developing, testing, and auditing custom smart contracts are crucial. Developers must ensure that their solutions are efficient, fair, and secure, considering factors such as network congestion, centralization risks, and computational costs.
The development of custom smart contracts can be costly. On average, the cost of developing a smart contract ranges between $5,000 and $25,000, depending on complexity and security requirements. This estimate includes the costs of designing, coding, testing, and auditing the contract. [Source: ConsenSys](https://consensys.net/blog/blockchain-development/the-cost-of-developing-a-smart-contract/) provides detailed insights into these costs. For accurate budgeting and planning, it's essential to factor in the potential expenses for ongoing maintenance and upgrades as well.

Intento makes it easy to set any action with given conditions at a fraction of the cost, and in a non-custodial manner.

## Bot Networks and Privileged Smart Contracts

### Bot Networks

Decentralized scheduled task execution through bot networks addresses challenges associated with server-based setups and central points of failure. In these networks, tasks are registered on an on-chain smart contract, and the first successful bot typically claims a reward. However, there are notable inefficiencies:

- **Network Congestion and Gas Costs:** Multiple bots executing transactions simultaneously can lead to increased congestion and higher gas costs.
- **Centralization Risks:** The first-come-first-served reward mechanism can lead to centralization, as less profitable agents may exit the network.
- **Complex Automation Workflows:** Arbitrary trigger addresses and the need for additional smart contract logic for privileged execution add complexity to the automation process.

These issues highlight the need for careful development, testing, and auditing of custom smart contracts. Developers must ensure that their smart contracts are optimized for efficiency and that the network remains decentralized and fair.

### Privileged Smart Contracts

In CosmWasm blockchains, privileged smart contracts can use BeginBlocker and EndBlocker functions for scheduling block-based executions. While this feature provides scheduling capabilities, it introduces several challenges:

- **Increased Computational Resources:** Privileged smart contracts require more computational resources per execution, leading to higher gas costs for users.
- **Fixed Computational Costs:** Running trigger engines inside smart contracts or virtual machines entails fixed costs and can impact fee revenue.
- **Balance of Blockspace Usage:** Blockchains need to balance time-based executions with general transactions, necessitating governance measures for blockspace limits and fee settings.

Compared to permissioned systems like Intento, which offers a protocol-neutral and permissionless platform with integrated governance parameters, privileged smart contracts may lack flexibility. Intento allows the community to manage the balance between blockchain network usage and fee revenue, providing a more adaptable solution for automated tasks.

## How Intento Automation Stands Out

In crypto, the efficiency and cost of transactions, often expressed through "gas costs" play a significant role in the user experience. Our comparison between various automation services like Gelato on Ethereum, CronCat on Juno and Neutron, and Intento reveal significant enhancements in security, efficiency and cost-effectiveness.

Gelato's services on Ethereum come with a protocol logic cost of 940,000 gas, translating to about $41, with an additional 20% protocol fee that elevates the cost per execution to $49. This cost structure is further compounded by risks associated with Gelato's proxy contract. On the other hand, CronCat, operating on Juno and Neutron, offers a more efficient protocol logic, consuming around 720,000 gas for execution, of which 170,000 gas is for the execution itself and about 550,000 gas for protocol logic.

In comparison, Intento has an estimated logic cost of just about 100,000 gas based on its lines of code. This stark reduction in gas usage not only offers substantial cost savings for users, enabling more frequent automation of transactions but also broadens accessibility.

The scalability unlocked by this efficiency unlocks new use cases and significantly enhances the user experience within crypto.

Read more on gas cost savings in [our blog post](https://intento.zone/post/gas-cost-in-action-processing-an-analysis/)
