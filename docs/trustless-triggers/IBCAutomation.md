---
order: 2
title: Interchain Automation
description: How AutoTX enables automation of assets on any IBC-enabled chain
---

## Interchain Automation

Trustless Triggers are time-based transactions on IBC-compatible chains that utilize the Interchain Accounts standard. They are automated by setting an interval, duration, end time, and optional start time in a MsgSubmitAutoTx.

This automation is great for automating user actions, such as sending tokens on a recurring basis or auto-compounding TRST tokens. Developers can use this feature to automate their protocols and build solutions for end-users to automate their funds using a newly generated account on the host chain, which is controlled by Trustless Hub and the user.

Interchain Accounts are a key component of Trustless Triggers. They allow for the creation and management of accounts across different IBC-connected chains. This means that Trustless Triggers can execute transactions on other chains based on custom logic, making them extremely versatile and useful for a wide range of applications.

To use Interchain Automation with Interchain Accounts, the user must first register an interchain account. This involves creating a port ID and connection ID, which allows the user to connect their account to other chains over IBC.


A Trigger is what we call an Automated Transaction or 'AutoTX'. An AutoTX is an object containing messages that execute at a specified time, or recurringly with intervals.
AutoTX entries are scheduled at the beginning of a new block.

AutoTXs can execute Cosmos SDK blockchain transactions on Cosmos Chains such as:

- `MsgSend` for token transfers
- `MsgWithdrawDelegatorReward` for reward claiming
- `MsgExecuteContact` to execute a contract
- `MsgInstantiate` to instantiate a contract

AutoTXs can also execute transactions on other chains using an IBC protocol called Interchain Accounts. 

Using the Authz module on the host chain - the chain you want to execute at - you can grant the Interchain Account on Trustless Hub access permission to execute a specific message.

## MsgSubmitAutoTx

Submitting a MsgSubmitAutoTx takes the following input:

| Field Name        | Data Type                      | Description                                                                                       | optional |
| ----------------- | ------------------------------ | ------------------------------------------------------------------------------------------------- | -|
| `owner`           | `string`                       | The owner of the transaction                                                                      |  |
| `connection_id`   | `string`                       | The ID of the connection to use for the transaction (in YAML format)                               |✔️|
| `label`           | `string`                       | A label for the transaction                                                                       |✔️|
| `msgs`            | `repeated google.protobuf.Any` | A list of arbitrary messages to include in the transaction                                        ||
| `duration`        | `string`                       | The amount of time that the transaction code should run for                                       ||
| `start_at`        | `uint64`                       | A Unix timestamp representing the custom start time for execution (if set after block inclusion) |✔️|
| `interval`        | `string`                       | The interval between automatic message calls                                                     |✔️|
| `fee_funds`       | `repeated cosmos.base.v1beta1.Coin` | Optional funds to be used for transaction fees, limiting the amount of fees incurred | ✔️|
| `depends_on_tx_ids` | `repeated uint64`           | Optional array of transaction IDs that must be executed before the current transaction is allowed to execute | ✔️|

Comments on the optionallity of the fields

- When `Interval` is not provided, the end of the duration will be the time the AutoTX executes.
- When `FeeFunds` are not provided, fees will be deducted from the Owner account.
- When `DependsOnTxIDs` is provided, AutoTX will see if their last execution was succesfull before execution can take place, else it will fail.
- When `ConnectionID` is not provided, it is assumed that `Msgs` are local Trustless Hub chain messages.

## Automation Process

1. Register an interchain account with `MsgRegisterInterchainAccount` or `MsgRegisterInterchainAccountAndSubmitAutoTx`.
2. Submit an AutoTX using `MsgSubmitAutoTx` - if fee funds are sent along with it, a new fee address is generated
3. Chain checks if parameters are ok
4. `AutoTX` is inserted in a queue
5. In each block,scheduled AutoTXs are retreived given the current block time
6. Fees are calculated and deducted. AutoTX data is updated with information on the exact fees and execution time.
7. Transaction is sent to the host chain
8. If AutoTX is recurring, a new entry is inserted into the queue
9. Packet gets acknowledged by a relayer and the AutoTX entry is updated stating execution was succesfull. If packet times out, execution fails. 
10. Funds sent to an AutoTX-specific FeeFund account are returned to the AutoTX owner

To make packet relays succesfully, Trustless Hub allocates token incentives to relayers to acknowledge packets.

*In-depth information on how the module works can be found in the module section of our documetation.*



Overall, Trustless Triggers provide a powerful tool for automating transactions and authorized actions on IBC-compatible chains. By utilizing Interchain Accounts and Interchain Automation, developers can create highly customizable and effective solutions for end-users.


## Limitations

Trustless Triggers are a powerful tool for automating transactions and actions on IBC-compatible chains. However, they do have some limitations that should be taken into consideration when designing applications or protocols.

For complex and custom logic, Trustless Contracts should be used. While Trustless Triggers can execute transactions based on custom logic, they are not as flexible or powerful as Trustless Contracts. Trustless Contracts are specifically designed for creating protocols with built-in automation and can execute transactions on other chains based on custom logic.

Ordered IBC channels are a necessary requirement for Interchain Automation with Interchain Accounts. This means that Trustless Triggers can only be executed when the previous execution did not time out. Channels close when a chain is available but a packet times out. 

Furthermore, IBC Packets may time out, which can lead to failed transactions. This can happen due to network congestion or other reasons, and can lead to a loss of funds or other negative consequences. To reduce the risk of this happening, Trustless Triggers have a timeout equal to the interval, so that the impact is minimal.

Finally, Trustless Triggers depend on the activeness of relayers. If relayers are not actively relaying transactions, Trustless Triggers may fail or take longer to execute. It is important to keep this in mind when using Trustless Triggers, and to ensure that there are active relayers available for reliable execution. We incentivize relayers to ensure active participation of relayers.

Overall, while Trustless Triggers are a powerful tool for automating transactions and actions on IBC-compatible chains, they do have some limitations that should be taken into consideration when designing applications or protocols