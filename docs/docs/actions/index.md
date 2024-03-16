---
title: Actions
order: 0
parent:
  title: Introduction
  order: 1
---

# Introduction

This knowledge base contains information that cover different aspects of automation of transactions on Intento.

In this module we use `Trustless Triggers` and `AutoTX` interchangeably. An AutoTX is what we call a Trustless Trigger and these are the same thing.

The AutoTX module is responsible for creating and executing automatic interchain transactions between different chains within the Cosmos ecosystem.

Automation using AutoTX module is an all-purpouse time-based automation module.



## Advantages of our approach

- Easy to integrate into dApps
- Automation of user funds using AuthZ
- No need to integrate with bots
- Removes the need for locking up tokens in Escrow or Vault smart contracts
- Schedule multiple messages into one AutoTX that execute sequentially

What's more:

- You can depend execution on other transactions
- Fee funds are refunded after execution finishes
- Relayers are incentivized for acknoledging a succesful automation packet on a host chain


### What can be automated?
<!-- TODO: highlight evm usecases -->

Trustless Triggers can automate on-chain actions by executing predefined blockchain calls in response to certain events or conditions. For example, a Trustless Trigger can be set up to automatically execute a payment function from a smart contract when a specific condition is met, such as a particular date or time, or when a specific amount of cryptocurrency is received.

Trustless Triggers are a key feature of Intento that enable automation of all kinds of function calls, such as sending transactions locally and dollar-cost averaging across IBC-enabled chains. IBC Relayers are incentivized to complete these transactions in a timely manner, and an acknowledgement is saved to ensure the transaction was successful. In case a trigger fails, it can be discarded, and if the destination chain is offline, the trigger will be discarded as well.

Trustless Triggers are a powerful tool for automating simple tasks, but they can also be used to create complex workflows by depending execution on other triggers. For example, a payment system could be built using triggers that depend on each other. A trigger could be set up to send a payment to a contractor every month. This trigger could depend on another trigger that checks if the contractor has completed the required work for the month. If this check is successful, the payment trigger is executed. If the check fails, the payment trigger is discarded.

In this way, Trustless Triggers can be used to build robust and reliable payment systems that require no manual input or verification. By chaining together triggers that depend on each other, organizations can automate the entire payment process, from checking if work has been completed to sending payments to contractors. And because the entire process is automated and decentralized, there is no need for intermediaries or trusted third parties, making the system more secure and less prone to fraud.
To prevent spam and network congestion, Intento has implemented various governance parameters for custom fees and rules. This helps to align incentives and maintain the stability and scalability of the network.

### ICS20 Middleware

ICS20 is an interchain standard that enables the transfer of fungible tokens between independent blockchains within the Cosmos ecosystem. It is a protocol that defines a standard interface for token transfers across different blockchains that implement the IBC protocol. With Intento’s ICS20 transfer middleware, you send a transfer token memo on a chain, and Intento will convert the token transfer to a trigger submission.

Using an ICS20-standard transaction, accounts on other chains can create triggers. Callers on other chains can specify trigger details in the memo field of an ICS20 message. Based on these inputs, the IBC hooks in the trigger module can build a submit message. This is useful for DAOs and other decentralized organizations on any IBC-enabled source  chain. They can now safely and reliably execute on IBC-enabled chains on a recurring basis, giving certainty to stakeholders, whilst also reducing the strain on governance proposal submitters and proposal voters.

Overall, Trustless Triggers are a powerful tool for building complex workflows across IBC-enabled chains. Whether it's automating simple transactions or building complex payment systems, Trustless Triggers can help organizations save time and money while also improving security and reliability.

