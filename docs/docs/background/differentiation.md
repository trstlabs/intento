---
title: Differentiation
sidebar_position: 3
---

# Differentiation

Intento differs from other decentralized automation solutions in several ways.

## Bot networks

Decentralized scheduled task execution is facilitated by bot networks, addressing challenges of server-based setups and central points of failure. Tasks are registered on an on-chain smart contract, with often the first successful bot claiming a reward. Inefficiencies arise with multiple bots executing transactions simultaneously, increasing network congestion and gas costs. Additionally, critics highlight issues like rewarding agents on a first-come basis, leading to centralization as unprofitable agents exit. In addition, having arbitrary trigger addresses pose limitations, requiring additional smart contract logic for privileged execution, thereby increasing the complexity of an automation workflow.

## Privileged Smart Contracts

In CosmWasm blockchains, smart contracts can run with privileges. These can leverage BeginBlocker and EndBlocker functions for scheduling block-based execution. Several chains have integrated this. This feature may seem similar to Intento, but the system is permissioned. Privileged smart contracts demand more computational resources per execution. Running trigger engines inside smart contracts increases gas costs for users and makes such a system less flexible. Running a trigger engine inside a virtual machine has fixed computational costs, as a check for scheduled executions needs to be performed on an recurring basis. Thus, this also impacts potential fee revenue of such a solution.

With recurring executions, blockchains have to find a balance for time-based executions and general transactions in terms of blockspace usage. Hence, we also find a necessity for governance for measures like limits per block and fee settings. Intento offers a protocol-neutral, permissionless platform with integrated governance parameters, allowing the community to find the balance between blnetwork usage and fee revenue.

## A Closer Look at Gas Costs: How Intento Automation Stands Out

In crypto, the efficiency and cost of transactions, often expressed through "gas costs" play a significant role in the user experience. Our comparison between various automation services like Gelato on Ethereum, CronCat on Juno and Neutron, and Intento reveal significant enhancements in security, efficiency and cost-effectiveness.

Gelato's services on Ethereum come with a protocol logic cost of 940,000 gas, translating to about $41, with an additional 20% protocol fee that elevates the cost per execution to $49. This cost structure is further compounded by risks associated with Gelato's proxy contract. On the other hand, CronCat, operating on Juno and Neutron, offers a more efficient protocol logic, consuming around 720,000 gas for execution, of which 170,000 gas is for the execution itself and about 550,000 gas for protocol logic.

In comparison, Intento has an estimated logic cost of just about 100,000 gas based on its lines of code. This stark reduction in gas usage not only offers substantial cost savings for users, enabling more frequent automation of transactions but also broadens accessibility.

Moreover, the efficiency of Intento unlocks new use cases and significantly enhances the user experience within crypto. The USP in gas cost efficiency not only positions Intento as a highly attractive option for users but also accentuates its role in establishing the Dymension Hub as a pivotal platform for interchain coordination. 


Read more on gas cost savings in [our blog post](https://intento.zone/post/the-economics-of-modular-automation-a-comparative-gas-cost-analysis)