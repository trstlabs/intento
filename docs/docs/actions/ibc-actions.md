---
sidebar_position: 1
title: Interchain Actions
pagination_label: How AutoTX enables automation of assets on any IBC-enabled chain
---

Intento can perform actions on IBC-compatible chains that utilize the Interchain Accounts standard. They are submitted by providing an interval, duration, end time, and optional start time in a `MsgSubmitAction` as well as IBC-specific settings such as the `connection_id`.

This is great for both automating actions, such as sending tokens or auto-compounding as well as for orchestrating workflows across chains. 
Developers can use this to automate their protocols and build solutions for end-users to automate their assets.

Interchain Accounts are a key component of Intento's Intents. They allow for the creation and management of accounts across different IBC-connected chains. This means that Intento's Intents can execute actions on other chains based on custom logic, making them extremely versatile and useful for a wide range of applications.

To use Interchain Actions with Interchain Accounts, the user must first register an interchain account. This involves creating a port ID and connection ID, which allows the user to connect their account to other chains over IBC.

An action is an object containing messages that are triggered at a specified time, or recurringly with intervals, given conditions.
AutoTX entries are scheduled at the beginning of a new block.

AutoTXs can execute Cosmos SDK blockchain messages on Cosmos Chains such as:

- `MsgSend` for token transfers
- `MsgSwapExactAmountIn` for token swapping on Osmosis
- `MsgWithdrawDelegatorReward` token reward claiming and auto compounding
- `MsgExecuteContract` to execute a CosmWasm contract
- `MsgInstantiateContract` to instantiate a CosmWasm contract

Intento executes messages on other chains using an IBC protocol called Interchain Accounts.

Using the Authz module on the host chain - the chain you want to execute at - you can grant the trigger on Intento permission to execute a specific message.

![IBC flow](@site/docs/images/ibc_trigger.png)

## MsgSubmitAction

Submitting a MsgSubmitAction takes the following input:

| Field Name         | Data Type                           | Description                                                                                                   | optional |
| ------------------ | ----------------------------------- | ------------------------------------------------------------------------------------------------------------- | -------- |
| `Owner`            | `string`                            | The owner of the Action                                                                                       |          |
| `Msgs`             | `repeated google.protobuf.Any`      | A list of arbitrary messages to include in the transaction                                                    |          |
| `Duration`         | `string`                            | The amount of time that the transaction code should run for                                                   |          |
| `StartAt`          | `uint64`                            | A Unix timestamp representing the custom start time for execution (if set after block inclusion)              | ✔️       |
| `Interval`         | `string`                            | The interval between automatic message calls                                                                  | ✔️       |
| `FeeFunds`         | `repeated cosmos.base.v1beta1.Coin` | Optional funds to be used for transaction fees, limiting the amount of fees incurred                          | ✔️       |
| `ConnectionID`     | `string`                            | The ID of the connection to use for the transaction (in YAML format)                                          | ✔️       |
| `HostConnectionID` | `string`                            | The ID of the host chain connection to use for the transaction (in YAML format)                               | ✔️       |
| `Configuration`    | `ExecutionConfiguration`            | Optional set of basic conditions and settings for the action                                                  | ✔️       |
| `Label`            | `string`                            | A label for the action                                                                                        | ✔️       |
| `Conditions`       | `repeated Condition`                | [V2] Powerful set of conditions for the action execution entry such as comparisons and event atribute parsing | ✔️       |

#### Optionality of the fields

- When `Interval` is not provided, the end of the duration will be the time the AutoTX executes.
- When `FeeFunds` are not provided, fees can be deducted from the Owner account by setting `FallbackToOwnerBalance` to true in `Configuration`.
- When `ConnectionID` or `HostConnectionID` are not provided, it is assumed that `Msgs` are local messages to be executed on Intento.

## Action Process

1. (Register an interchain account with `MsgRegisterAccount` or `MsgRegisterAccountAndSubmitAction` when acting on a chain for the first time)
2. Submit an AutoTX using `MsgSubmitAction` - if fee funds are sent along with it, a new fee address is generated
3. Chain checks if execution settings from `Conditions` are ok
4. `AutoTX` is inserted in a queue
5. In each block,scheduled AutoTXs are retrieved given the current block time
6. Fees are calculated and deducted. AutoTX data is updated with information on the exact fees and execution time.
7. (IBC transaction is sent to the host chain)
8. If AutoTX is recurring, a new entry is inserted into the queue
9. Packet gets acknowledged by a relayer and the AutoTX entry is updated stating execution was succesfull
10. Funds sent to an Action account are returned to the AutoTX owner

_Read more on how the module works in the [module](@site/docs/modules/index.md) section of our documetation._

<!-- ## Considerations

Intento's Intents are a powerful tool for automating actions over IBC. However, there are some limitations that should be taken into consideration when designing applications or protocols.

Ordered IBC channels are a necessary requirement for Interchain Automation with Interchain Accounts. This means that Intento's Intents can only be executed when the previous execution did not time out. Channels close when a chain is available but a packet times out. IBC Packets may time out, which can lead to failed actions. This can happen due to network congestion or other reasons, and can lead to a loss of funds or other negative consequences. To reduce the risk of this happening, Intento's Intents by default have a timeout equal to the interval, so that the impact is minimal.

Actions depend on the availability of relayers. If relayers are not actively relaying IBC transactions, Intento's Intents may fail or take longer to execute. It is important to keep this in mind when using Intento's Intents, and to ensure that there are active relayers available for reliable execution. -->