---
title: Differentiation
sidebar_position: 3
---

# Differentiation

Intento differs from other decentralized automation solutions in several ways.

## Bot networks

Bot networks offer decentralized scheduled execution of tasks. Setting up a bot on a server is hard-to-implement, and poses a central point of failure. Hence, decentralized networks of bots such as Gelato and Chainlink Keepers on Ethereum and CronCat on CosmWasm chains aim to solve this problem. Tasks are registered in an on-chain smart contract. The bot that is first to execute a task successfully can claim a reward. This solves some of the problems for enabling trust-minimized scheduled execution. However,  opponents of this method identify several topics of discussion. The most common critique is on rewarding agents. By rewarding the agent on a first come basis, agents with better resources will win each time. As unprofitable agents exit the network, centralization occurs over time. Moreover, by printing fee tokens for rewards, the value of the token erodes over time, as can be seen in the price action of reward tokens. 

A further downside not often mentioned is inefficiency. Transactions will be signed and executed by several bots at once. This is expensive, as fees will have to be paid by the bot network participants. Further, this increases network congestion. In times of network growth, this may lead to increased gas cost and reduced scalability of the bot network and the underlying general purpose smart contract blockchain.
We also find some profound limitations for developers to this approach over our on-chain execution network. One of these is that the agent address that executes the task is arbitrary. Thus, it is not possible to automate assets based on a message grant by granting the executor to execute on an account’s behalf. Assets can be transferred to the bot network smart contract or to an intermediate contract, extracting them away from the user’s wallet. To liquify these assets, these assets would have to first have to be withdrawn from the smart contract. This adds complexity and poses smart contract risks, especially for when intermediate contracts are not programmed to account for arbitrary contract callers. Locking up assets for the complete duration can significantly increase liquidity risks of the task scheduler. As mentioned previously, the executor in Intento is the chain itself, allowing for rapid and secure development of time-based tasks in an easier, more efficient manner with minimal externality risks.

## Privileged smart contracts

In CosmWasm implementations it is possible for privileged smart contacts to call functions on a block basis based on the BeginBlocker or EndBlocker function. Production blockchains where this is integrated  are Injective and Tgrade. On Injective, CosmWasm contracts can interact with the on-chain orderbook, which uses block-based execution for trades. In addition, privileged smart contacts, granted by the team or chain governance can use the chain’s BeginBlocker or EndBlocker hooks. These are blockchain functions that run in the beginning and ending of the block, respectively, and allow the chain/contract to call the contract. On a high level, this approach is highly similar to our Trustless Contracts solution, as these also depend on the BeginBlocker or EndBlocker functions. Next to our solution being permissionless, execution parameters can be set easily at contract instantiation. In theory, a similar on-chain solution to Intento's Intents may be built into a privileged smart contract. However, this poses several challenges and risks, to the underlying smart contract blockchain and their users in the long term. 

Blockchains and their communities indirectly make a trade-off between blockspace availability for time-based calls and general purpose transactions. To combat unnecessary network congestion and spam, some measures would have to be taken at the smart contract and/or chain level. One of these measures is allowing only privileged smart contracts to execute. Other measures include a transaction limit or a fee mechanism. These mechanisms should be actively governed based on network conditions. Inadequate regulation endangers both protocol managers and end users in performing general purpose transactions. Excessive regulation may reduce decentralization and adoption.
By granting only selected decentralized applications privileges, centralization may occur. As more privileged decentralized applications are built, approving privileges to new decentralized applications with minor enhancements becomes harder to justify. Furthermore, having to request permission makes the network less attractive for new projects to choose the network as their preferred platform. 

Intento provides a neutral and permissionless platform for developing time-based automation applications and has these mechanisms built right into the chain’s governance. Governance of trigger parameters fits into Cosmos’s App-Chain thesis as time-based programmability is the primary application of the chain. In practice, this means that sophisticated trade-offs can be made where the interest of to-be scheduled automations is safeguarded, and where developers and users do not experience excessive fees or restrictions. 

## On-Chain Automation

Because Intento  automation on-chain, the caller of the execution is the validator set. Because of this, transfers of assets can be automated.
This is the major difference that Intento provides compared to other solutions. With Transaction Automation using Interchain Accounts, the caller is a predictable address generated on-chain.

| Intento's Intents                                | Trustless Contracts                                                       | Bot Networks (Gelato, CronCat)                                     |
|---------------------------------------------------|---------------------------------------------------------------------------|--------------------------------------------------------------------|
| Time-based actions with a predictable address                          | Triggers triggered by a the contract that supports a wide range of events                                       | Triggers are executed by a third party addresses |
| Automates financial transactions                  | Can create protocols with built-in automation                             | Can be used to create trading bots and other automated tools       |
| Easy for end-users to set up                      | Requires technical expertise to develop and deploy                        | May require technical expertise to set up and integrate with bots  |
| Highly customizable and flexible                  | Adaptable to a wide range of use cases                                    | Can be customized and optimized for specific use cases             |
| Enhances security and control over funds          |  smart contract risks              |  Smart contract risks and may pose additional security risks if not implemented correctly               |
| Reduces transaction fees and increases efficiency | Can reduce the need for intermediaries and third-party payment processors | No relay time to consider as bots are off-chain |


## A Closer Look at Gas Costs: How Intento Automation Stands Out

In crypto, the efficiency and cost of transactions, often expressed through "gas costs" play a significant role in the user experience. Our comparison between various automation services like Gelato on Ethereum, CronCat on Juno and Neutron, and Intento reveal significant enhancements in security, efficiency and cost-effectiveness.

Gelato's services on Ethereum come with a protocol logic cost of 940,000 gas, translating to about $41, with an additional 20% protocol fee that elevates the cost per execution to $49. This cost structure is further compounded by risks associated with Gelato's proxy contract. On the other hand, CronCat, operating on Juno and Neutron, offers a more efficient protocol logic, consuming around 720,000 gas for execution, of which 170,000 gas is for the execution itself and about 550,000 gas for protocol logic.

In comparison, Intento has an estimated logic cost of just about 100,000 gas based on its lines of code. This stark reduction in gas usage not only offers substantial cost savings for users, enabling more frequent automation of transactions but also broadens accessibility.

Moreover, the efficiency of Intento unlocks new use cases and significantly enhances the user experience within crypto. The USP in gas cost efficiency not only positions Intento as a highly attractive option for users but also accentuates its role in establishing the Dymension Hub as a pivotal platform for interchain coordination. 


Read more on gas cost savings in [our blog post](https://intento.zone/post/the-economics-of-modular-automation-a-comparative-gas-cost-analysis)