---
order: 2
title: CrossChain Automation
description: How AutoTx enables automation of assets on any IBC-enabled chain
---

## CrossChain Automation

Automation is done with what we call an 'AutoTx'. An AutoTx is an object containing messages that execute at a specified time, or recurringly with intervals.
AutoTxs are scheduled at the beginning of a new block. AutoTxs can execute local Cosmos SDK blockchain transactions like MsgSend for transfers or MsgWithdrawalRewards for reward claiming and can also execute transactions on other chains using an IBC protocol called Interchain Accounts.

##  Automation using Interchain Accounts

1. Register an interchain account with MsgRegisterInterchainAccount or do MsgRegisterInterchainAccountAndSubmitAutoTx to do both in one message.
2. SubmitAutoTx is submitted - if fee funds are sent along with it, a new fee address is generated
3. Chain checks if duration, interval, and messages sent are ok
4. AutoTx is inserted in a queue
5. In each block, AutoTxs to be scheduled are retreived given the current block time
6. Fees are calculated and deducted, AutoTx entry is set with information on the exact fees and execution time.
7.  Message is sent to the host chain
8. If AutoTx is recurring, a new entry is put into the queue
9.  Packet gets acknowledged by a relayer and the AutoTx entry is updated stating execution was succesfull. If packet times out, execution fails. Hence, TRST provides token incentives for relayers to acknowledge packets.
10. Funds sent in a fee fund account are refunded


*In-depth information can be found in the module section of our documetation.*