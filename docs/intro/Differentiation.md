---
order: 2
title: Differentiation
description: Built-in Automation is Smart Contracts 2.0
---

## On-Chain Automation

Because Trustless Hub performs automation on-chain, the caller of the execution is the validator set. Because of this, transfers of assets can be automated.
This is the major difference that Trustless Hub provides compared to other solutions. With Transaction Automation using Interchain Accounts, the caller is a predictable address generated on-chain.

| Trustless Triggers                                | Trustless Contracts                                                       | Bot Networks (Gelato, CronCat)                                     |
|---------------------------------------------------|---------------------------------------------------------------------------|--------------------------------------------------------------------|
| Time-based transactions                           | Triggered by a wide range of events                                       | Automation of transactions are executed by a third party addresses |
| Automates financial transactions                  | Can create protocols with built-in automation                             | Can be used to create trading bots and other automated tools       |
| Easy for end-users to set up                      | Requires technical expertise to develop and deploy                        | May require technical expertise to set up and integrate with bots  |
| Highly customizable and flexible                  | Adaptable to a wide range of use cases                                    | Can be customized and optimized for specific use cases             |
| Enhances security and control over funds          | Offers advanced security features with code-based contracts               | May pose security risks if not implemented correctly               |
| Reduces transaction fees and increases efficiency | Can reduce the need for intermediaries and third-party payment processors | Can help reduce transaction costs and increase trading efficiency  |


## Self-Executing contracts

Trustless Contracts can have custom logic with CosmWasm. Effient and secure. 
On Trustless Hub it is possible to build trustless cross-chain dApps over IBC. 
In the block creation process on our chain, automatic execution has priority over general transactions. This allows developers and end users to experience time-based execution in a reliable manner. Execution Fees go to the community, fostering developement and growth of the ecosystem. Funds and fees get refunded after execution. This is not possble in bot networks where rewards go to bot operators. 

On bot networks your trade and information may be read beforehand and used against you. With Trustless Hub's smart contracts, encrypted inputs ensure no MEV takes place before and during execution.

Trustless Hub use CosmWasm, Cosmos's leading smart contract VM. As automation is built-in, automating function calls is simple. Contracts on Trustless Hub can also execute with IBC-enabled contracts on other chains.  

### Key differences with bots
![differentiation](./../images/exec.png)