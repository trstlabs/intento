---
order: 4
title: AutoMessage
description: Learn the basics of how AutoMessage enable 1-click contract self-execution
---

# AutoMessage
An AutoMessage is a message that can execute at a specified time, or recurringly with intervals. The callable function can be specified by the contract and can be anything. As this is a powerful feature, the AutoMessage is defined at instantiation, so that contract interactors can view that given contract will self-execute at a pre-defined moment.


## Process
- An AutoMessage is encrypted at contract instantiation
- A predefined execute-time is provided when storing the code
- This can only be sent to the same contract as instantiated by the creator
- Then a *callback signature* is created. This is a signature with a hash containing the address of the contract and the message, so that only the chain is able to execute the message
- When the contract is defined to execute, the AutoMessage is retrieved along with the callback signature

## Fees
To make  the AutoMessage executable, send TRST as funds. The remainder will be refunded to the account you're instantiating from. If this is a contract (e.g. TIP20) the contract accrues TRST tokens. TIP20 contracts are able to spend or burn these accrued fees. If you are instantiating directly, Trustless Hub can charge the fee directly from your available TRST balance.
At launch, Contracts with AutoMessages are incentivized, so that fees are (near) zero. Over time, incentives graduately decline.

## Security
As the AutoMessage is encrypted with the newly created contract address it can not be use my malicious nodes to expose information of other contracts.
