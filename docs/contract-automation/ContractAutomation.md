---
order: 3
title: Contract Automation
description: How AutoMessage enable contract self-execution and enables 1-click user experiences for time-based actions
---

## Contract Automation with AutoMessage

Automation in self-executing contracts are performed using what we call AutoMessages. An AutoMessage is a message that executes at a specified time, or recurringly with intervals. The callable function is specified in the contract. As this is a powerful feature, the 'AutoMsg' must be defined at instantiation of a contract. Contract interactors can view that given contract will self-execute at pre-defined time(s).

## Process

- AutoMessage is encrypted at contract instantiation
- A predefined execute-time is provided when storing the code
- This can only be sent to the same contract as instantiated by the creator
- Then a *callback signature* is created. This is a signature with a hash containing the address of the contract and the message, so that only the chain is able to execute the message
- When the contract is set to execute according to the predefined execution schedule, the AutoMessage is retrieved along with the callback signature. 

## Fees

Auto execution messages have a fee, we direct a portion of these to validators. Governance can alter the fee.

The fee consists of:

1. Gas-dependent Flex Fee (Goes to validators)

2. Constant Fee (Goes to community pool)

3. Commission (Goes to community pool)

To make the AutoMessage executable, TRST is sent at instantiation of the proxy contract. You don't have to worry about sending too much, as remaining tokens are refunded automatically. The community pool may take a commission. The remainder will be refunded to the account that is set as the owner. If you are instantiating directly from your account, Trustless Hub can also charge the fee directly from your available TRST balance.
At launch, Contracts with AutoMessages are incentivized, so that fees are (near) zero. Over time, incentives graduately decline.

## Security

As the AutoMessage is encrypted with the newly created contract address it can not be use my malicious nodes to expose information of other contracts.

## Proxy-contract

A contract that extends the functionality of an address. Through giving this contract allowance/approval, the proxy contract acts on your behalf. TIP20 tokens can instantiate a proxy contract with approval in just one message. Trustless Hub uses proxy contracts as it is the safest way to run self execution. There are no third parties. As you or a TIP20 token instantiates the contract, your tokens remain in safe hands.

## Contracts always execute

Should it be the case that a block is full, the automatic execution does not get deleted but rather will take place in the following block. Should the chain halt, AutoExecution will execute in the following block(s). For recurring execution, the next execution time is based on the previous expected execution time + duration. This means that recurring execution always runs following the foreseen execution pattern.

When a contract wants to skip the execution if time differs significantly, custom time-out logic can be implemented easily into the CosmWasm contract code. Block height and time is available during any type of execution, including auto-execution. For our DCA and RecurrentlySend contracts our implementation is to always execute, even when actual block time differs from the foreseen execution time.

TRST Labs will continue to find the optimal way to implement a fair mechanism that users will love. Do feel free to share your thoughts on our design. We want to find the best fee structure design that will result in a high transaction volume and onboard the most users. 