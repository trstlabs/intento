---
sidebar_position: 1
title: Introduction
pagination_label: Interchain actions. Any action, anywhere
---

Intents contain actions and predefined conditions, which determine how actions are processed. In crypto, we want to process actions with no counterparty risk and in a permissionless and decentralized manner. This has so far only been possible by directly sending transactions to the destination chain. With Intento, you will be able specify your intent to perform any action, anywhere at any moment given any conditions.

Intent-based actions are submitted by providing an interval, duration, end time, and optional start time in a `MsgSubmitAction`, IBC-specific settings such as the `ConnectionID` or a `HostedAccount`, and conditions such as comparisions and feedback loops.

An action in technical terms is an object containing messages that are triggered at a specified time, or recurringly with intervals, with given conditions.

Action entries are scheduled at the beginning of a new block.

Interchain Accounts are a key component of Intento. They allow for the secure management of assets across different IBC-connected chains. This means that Intento's Intents can execute actions on other chains based on custom logic, making them extremely versatile and useful for a wide range of applications. With this, you can execute Cosmos SDK blockchain messages such as:

- `MsgSend` for token transfers
- `MsgSwapExactAmountIn` for token swapping on Osmosis
- `MsgWithdrawDelegatorReward` token reward claiming and auto compounding
- `MsgExecuteContract` to execute a CosmWasm contract
- `MsgInstantiateContract` to instantiate a CosmWasm contract

#### Approaches for Executing Messages on Other Chains

Intento can execute messages on other chains using several approaches:

1. **ICS20 Transfers with a Memo**
   - Easy to set up on available chains by using packet forwarding
   - Memo field actions have limited support by chains
2. **Hosted Interchain Accounts**
   - Easy to set up and manage
   - Host chain fees are managed by an admin
   - You configure a fee limit
3. **Self-Hosted Interchain Accounts**
   - Full control over the account
   - You have to manage host chain fee balances yourself

To use Self-Hosted Interchain Accounts, you first register an interchain account. This involves creating a port ID and connection ID, which allows you to connect their account to other chains over IBC. You additionaly have to send funds for fees on the host chain.

Using the Authz module on the host chain - the chain you want to execute at - you can grant the trigger on Intento permission to execute a specific message.

![IBC flow](@site/docs/images/ibc_trigger.png)

## MsgSubmitAction

Submitting an action with MsgSubmitAction can be done with the following input:

| Field Name              | Data Type                           | Description                                                                                      | optional |
| ----------------------- | ----------------------------------- | ------------------------------------------------------------------------------------------------ | -------- |
| `Owner`                 | `string`                            | The owner of the action                                                                          |          |
| `Msgs`                  | `repeated google.protobuf.Any`      | A list of arbitrary messages to include in the transaction                                       |          |
| `Duration`              | `string`                            | The amount of time that the transaction code should run for                                      |          |
| `Label`                 | `string`                            | A label for the action                                                                           | ✔️       |
| `StartAt`               | `uint64`                            | A Unix timestamp representing the custom start time for execution (if set after block inclusion) | ✔️       |
| `Interval`              | `string`                            | The interval between automatic message calls                                                     | ✔️       |
| `FeeFunds`              | `repeated cosmos.base.v1beta1.Coin` | Optional funds to be used for transaction fees, limiting the amount of fees incurred             | ✔️       |
| `ConnectionID`          | `string`                            | The ID of the connection to use for a self-hosted ICA                                            | ✔️       |
| `HostConnectionID`      | `string`                            | The ID of the host chain connection to use for a self-hosted ICA                                 | ✔️       |
| `HostedAccount`         | `string`                            | Hosted ICA account that executes on a host chain on your behalf                                  | ✔️       |
| `HostedAccountFeeLimit` | `cosmos.base.v1beta1.Coin`          | A limit of the fees a hosted account can charge per action execution                             | ✔️       |
| `Configuration`         | `ExecutionConfiguration`            | Optional set of basic conditions and settings for the action                                     | ✔️       |
| `Conditions`            | `repeated Condition`                | Powerful set of conditions for the action execution entry such as comparisons and feedback loops | ✔️       |

#### Notes

- When `Interval` is not provided, the end of the duration will be the time the action executes.
- When `FeeFunds` are not provided, fees can be deducted from the Owner account by setting `FallbackToOwnerBalance` to true in `Configuration`.
- When `ConnectionID`,`HostConnectionID` and `HostedAccount` are not provided, it is assumed that `Msgs` are local messages to be executed on Intento.
- `HostedAccount` requires `HostedAccountFeeLimit`

## Action History

The `Action History` captures the timeline and details of every executed action and it's messages, providing transparency and accountability. By tracking scheduled times, actual execution times, outcomes, and errors, actions can be monitored and acted upon.

### Anatomy of an Action History Entry

| **Field**              | **Description**                                                                                                   |
|------------------------|-------------------------------------------------------------------------------------------------------------------|
| `scheduled_exec_time`  | Specifies the time at which the action was originally scheduled for execution.                                    |
| `actual_exec_time`     | The time at which the action was actually triggered. This provides insight into delays or on-time performance.     |
| `exec_fee`             | The fee incurred for executing the action. This is represented as a `Coin` object and is non-nullable.            |
| `executed`             | A boolean flag indicating whether all messages within the action were executed, regardless of the outcome.        |
| `timed_out`            | A boolean flag that specifies if the action timed out during execution over IBC.                                  |
| `errors`               | A list of error messages encountered during the execution. If empty, the execution was successful.                |
| `msg_responses`        | A collection of responses for the executed messages. This field remains empty if `save_responses` is set to false.|
| `query_responses`      | Contains responses to any queries performed. This field is also empty if `save_responses` is false.               |

## Execution Process

A high level overview of the process from action submission to triggering messages across chains.

1. Submit an action using `MsgSubmitAction` - if fee funds are sent along with it, a new fee address is generated
2. Chain checks if execution settings from `Conditions` are ok.
3. An Intento `Action` is inserted in a queue
4. In each block, scheduled actions are retrieved given the current block time
5. If set, comparisons are made with prior outputs and Interchain Queries and return true or false.
6. Feedback Loops are checked. Message and and Interchain Query responses are inserted into messages.
7. IBC transaction is sent and executed to the host chain
8. If action is recurring, a new entry is inserted into the queue
9. IBC Packet gets acknowledged by a relayer and the action entry is updated
10. Remaining funds sent to an action account are returned to the action owner

<!-- _Continue reading on on how intent module works in the [module](@site/docs/module/index.md) section of our documetation._ -->

<!-- ## Considerations

Intento's Intents are a powerful tool for automating actions over IBC. However, there are some limitations that should be taken into consideration when designing applications or protocols.

Ordered IBC channels are a necessary requirement for Interchain Automation with Interchain Accounts. This means that Intento's Intents can only be executed when the previous execution did not time out. Channels close when a chain is available but a packet times out. IBC Packets may time out, which can lead to failed actions. This can happen due to network congestion or other reasons, and can lead to a loss of funds or other negative consequences. To reduce the risk of this happening, Intento's Intents by default have a timeout equal to the interval, so that the impact is minimal.

Actions depend on the availability of relayers. If relayers are not actively relaying IBC transactions, Intento's Intents may fail or take longer to execute. It is important to keep this in mind when using Intento's Intents, and to ensure that there are active relayers available for reliable execution. -->
