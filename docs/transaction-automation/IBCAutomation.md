---
order: 2
title: CrossChain Automation
description: How AutoTx enables automation of assets on any IBC-enabled chain
---

## CrossChain Automation

Automation is done with what we call an 'AutoTx'. An AutoTx is an object containing messages that execute at a specified time, or recurringly with intervals.
AutoTxs are scheduled at the beginning of a new block. AutoTxs can execute local Cosmos SDK blockchain transactions like MsgSend for transfers or MsgWithdrawDelegatorReward for reward claiming. AutoTxs can also execute transactions on other chains using an IBC protocol called Interchain Accounts. Using the Authz module on the host chain - the chain you want to execute at - you grant the Interchain Account on Trustless Hub access to execute the message.

## Automation using Interchain Accounts

1. Register an interchain account with MsgRegisterInterchainAccount or do MsgRegisterInterchainAccountAndSubmitAutoTx.
2. Submit an AutoTx using MsgSubmitAutoTx - if fee funds are sent along with it, a new fee address is generated
3. Chain checks if duration, interval, and messages sent are ok
4. AutoTx is inserted in a queue
5. In each block,scheduled AutoTxs are retreived given the current block time
6. Fees are calculated and deducted and the AutoTx data is updated with information on the exact fees and execution time.
7. Transaction is sent to the host chain
8. If AutoTx is recurring, a new entry is put into the queue
9. Packet gets acknowledged by a relayer and the AutoTx entry is updated stating execution was succesfull. If packet times out, execution fails. Hence, TRST provides token incentives for relayers to acknowledge packets.
10. Funds sent in a fee fund account are refunded


*In-depth information can be found in the module section of our documetation.*