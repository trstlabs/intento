---
sidebar_position: 3
title: Creating actions from a connected chain
description: How to setup automation from a connected chain
---

## Setting up Automation

In the previous step we showed how the action process looks like by submitting an AutoTX on Intento. You can do this with the [TriggerPortal](https://triggerportal.zone) interface, a TrustlessJS front-end integration or locally through the CLI.

In addition, you can also submit an AutoTX from another chain using the [ICS20 standard](https://github.com/cosmos/ibc-go/blob/main/docs/apps/transfer/messages.md).

### Interchain Accounts

Users and entities on Cosmos SDK chains may be able set up [interchain account](https://tutorials.cosmos.network/academy/3-ibc/8-ica.html) to Intento and submit the actions using `MsgSubmitAction`. It is also easy to deploy using the ICS20 standard.

### ICS20 Standard

With Intento’s ICS20 transfer middleware, you send a transfer token memo on a chain, and Intento will convert the token transfer to an action submission.
ICS20 is an interchain standard that enables the transfer of fungible tokens between independent blockchains. It is a protocol that defines a standard interface for token transfers across different blockchains that implement the Inter-Blockchain Communication (IBC) protocol.

Using ICS20, accounts on connected chains can create actions. Sebders specify action details in the memo field of the ICS20 message. Based on these inputs, the IBC hooks in the intent module build a submit message.

This is useful for DAOs and other decentralized organizations on any connected chain. They can safely and reliably execute on Intento's connected chains. For DAOs, this gives certainty to stakeholders, whilst also reducing manual work on governance proposals.

## Automation for DAOs

Setting up an action on a connected chain can be particularly useful for DAOs. Using this middleware, DAOs can now automate tasks not only on their chain but also on any chain connected to Intento. What can DAO's do with this? DAOs can automate periodic token swaps, payroll, payment in installments amongst other scheduled actions. These actions can be performed in one proposal, which normally require periodically voting on individual proposals. This normally requires manual action from the proposer and DAO participants.

## ICS20 Middleware

A MsgRegisterAccountAndSubmitAction or a MsgSubmitAction can be derrived from the memo field in the ICS20 transfer message.

![ics20](@site/docs/images/connected_chain/from_connected_chain.png)

Our custom middleware is based on the wasmhooks implementation on [Osmosis](https://github.com/osmosis-labs/osmosis/tree/main/x/ibc-hooks).

The mechanism enabling this is a `memo` field on every ICS20 transfer packet as of [IBC v3.4.0](https://medium.com/the-interchain-foundation/moving-beyond-simple-token-transfers-d42b2b1dc29b).

ics_middleware.go is IBC middleware that parses an ICS20 transfer, and if the `memo` field is of a particular form, it creates an action by parsing and handling a SubmitAction message.

These are the fields for `action` that are derived from the ICS20 message:

- **Owner**: This field is directly obtained from the ICS20 packet metadata and equals the ICS20 recipient. If unspecified, a placeholder is made from the ICS20 sender and channel.
- **Msg**: This field should be directly obtained from the ICS20 packet metadata.
- **FeeFunds**: This field is set to the amount of funds being sent over in the ICS20 packet. One detail here is that the denom in the packet is the source chains representation of the denom, this will be translated into INTO on Intento.

The constructed message for MsgSubmitAction under the hood will look like:

```go
msg := MsgSubmitAction{
 // If let unspecified, owner is the actor that submitted the ICS20 message and a placeholder only
 Owner: "into1-hash-of-channel-and-sender" OR packet.data.memo["action"]["owner"],
 // Array of Msg json encoded, then transformed into a proto.message
 Msgs: packet.data.memo["action"]["msgs"],
 // Funds coins that are transferred to the owner
 FeeFunds: sdk.NewCoin{Denom: ibc.ConvertSenderDenomToLocalDenom(packet.data.Denom), Amount: packet.data.Amount}

 // other fields
}
```

### ICS20 packet structure

So given the details above, we propogate the implied ICS20 packet data structure.
ICS20 is JSON native, so we use JSON for the memo format.

```json
{
  //... other ibc fields that we don't care about
  "data": {
    "denom": "INTO denom on counterparty chain (e.g. ibc/abc...)",
    "amount": "1000", //for execution fees
    "sender": "...",
    "receiver": "A INTO addr prefixed with into1",
    "memo": {
      "action": {
        "owner": "into1address", //owner is optional
        "msgs": [
          {
            "@type": "/cosmos.somemodule.v1beta1.sometype"
            //message values in JSON format
          }
        ],
        "duration": "111h",
        "start_at": "11h",
        "label": "my_label",
        "interval": "11h", //optional
        "cid": "connection-0", //connection ID is optional, omit or leave blank in case local INTO message.
        "cp_cid": "connection-0", //counterparty connection ID is optional and is only needed to register ICA.
        "register_ica": "false", //optional, set to true to register interchain account
        //////configuration,optional
        "save_responses": "true", //save message responses of Cosmos SDK v0.46+ chain output, defaults to false
        "update_disabled": "true", //optional, disables the owner's ability to update the config, defaults to false
        "stop_on_success": "true", //optional, defaults to false
        "stop_on_fail": "true" //optional, defaults to false
      }
    }
  }
}
```

An ICS20 packet is formatted correctly for submitting an action if the following all hold:

- `memo` is not blank
- `memo` is valid JSON
- `memo` has at least one key, `"action"`
- `memo["action"]["msgs"]` is an array with valid JSON SDK message objects with a key "@type" and sdk message values
- `receiver == memo["action"]["owner"]`. Optional, an owner can be specifed and is the address that receives remaining fee balance after execution ends.
- `memo["action"]["cid"]`is a valid connection ID on INTO -> Destination chain, omit it for local INTO execution of the message.
- `memo["action"]["register_ica"]` can be added, and true to register an ICA.

Fees are paid with a newly generated and AutoTX specific fee account.

If an ICS20 packet does not contain a memo containing "action", a regular MsgTransfer takes place.
If an ICS20 packet is directed towards AutoTX, and is formated incorrectly, then it returns an error.

## Example: DAO Integration

![icsdao](@site/docs/images/connected_chain/from_connected_chain_flow1.png)

Using ICS20 to set up actions is for advanced users that are familiar with terms such as ICS20 transfers and Authz permissions.

Our ICS20 middleware is developped to allow DAOs to perform time-based actions. These can be executed locally on Intento, a destination chain which can also be the source chain.

There are several caveats to using ICS20 to set up an action. When automating on a destination chain for the first time, not one but two transactions is required to activate the trigger. One transaction is to set the action and create an Interchain Account at Intento and another is to set permissions/send funds to the Interchain Account on the destination chain.

To demonstrate how DAOs can integrate with Intento's Intents we explain the process using an example.

`DAO` wants to pay `Service provider ABC` for their services that they provide on a monthly basis in `TOKEN1` and holds `TOKEN2` and `NTRN`.

_Service provider ABC invoice example_
![DAODAO](@site/docs/images/ics20/daodao_proposal1.png)

In this example a `DAO` on DAODAO triggers a swap of TOKEN2 for TOKEN1 on dex "DEX" on a recurring basis.
Then it can automatically send these tokens to `Service provider ABC`.
This is a neat use case for Intento's Intents as it requires movement from assets between chains and accounts. The tokens to send can remain in control of the DAO.
As action triggers are on-chain, the process can be done completely without trusted third parties.

The DAO can appoint an owner that can manage the action or use a placeholder account on Intento to remain in full control.

For this invoice example we have the following proposal name and description:

_Example Proposal name_

```md
[Trigger] Pay service provider ABC in TOKEN1
```

_Example Description_

```md
Submit an action to send TOKEN1 to "service provider ABC"

This will swap on a recurring basis TOKEN2 for TOKEN1 on dex "DEX" on Destination Chain and send these tokens to "service provider ABC"
This will swap TOKEN1 for TOKEN2 on a recurring basis so that:

- the selling pressure on TOKEN2 is gradual
- we maintain a positive cash flow

This Trigger on Intento allows us to automate asset workflows. It allows our DAO to maintain positive cashflow and have better liquidity.
```

Having adequate liquidity is crucial to ensure smooth operations and financial stability. It enables you to fulfill your financial obligations to suppliers, lenders, and other stakeholders, which strengthens your credibility and reputation in the marketplace. With sufficient liquidity, you can easily pay for bills and expenses, invest in growth opportunities, and effectively manage your DAO's resources. Thus, prioritizing and managing liquidity effectively can help you achieve long-term success and sustainable growth.

### 1. Submitting AutoTX

![icsdao](@site/docs/images/connected_chain/from_connected_chain_flow2.png)

Submit a proposal with a custom message to execute. For a CosmWasm-based DAO like [DAO DAO](https://daodao.zone/) DAOs on Neutron, This message contains the ICS20 message`MsgTransfer`. In the `memo` you provide the `Action` details such as the messages to automate, time parameters and set a custom label.

Below is how the proposal action message could look like for a CosmWasm-based DAO:

```json
[
  {
    "stargate": {
      "typeUrl": "/ibc.applications.transfer.v1.MsgTransfer",
      "value": {
        "source_port": "transfer",
        "source_channel": "channel-to-intento",
        //token should be IBC-denominated INTO, which is used to pay for fees
        "token": { "denom": "ibc/....", "amount": "10" },
        // sender - a bech32 address on source chain
        "sender": "neutron1validbech32address",
        // the recipient address on Intento
        // is replaced with a newly generated address when action["owner"] in memo is blank
        "receiver": "into1address",
        // Timeout height relative to the current block height.
        // The timeout is disabled when set to 0.
        "timeout_height": "0",
        // Timeout timestamp in absolute nanoseconds since unix epoch.
        // The timeout is disabled when set to 0.
        "timeout_timestamp": "0",
        "memo": {
          "action": {
            //owner is optional and should equal recipient
            "owner": "into1address",
            "msgs": [
              {
                "@type": "/someprefix.somemodule.someversion.sometype"
                //add message values
                //in JSON format
              }
            ],
            //duration from now or from start_at. s=seconds, m=minutes, h=hours
            "duration": "111h",
            //interval is optional. s=seconds, m=minutes, h=hours
            "interval": "11h",
            //start_at is optional, in nanoseconds since UNIX epoch
            "start_at": "1677841601",
            //label is optional
            "label": "my_label",
            //connection_id is optional, omit or leave blank in case local Intento message.
            "cid": "connection-0",
            //register_ica is optional, set to true to register interchain account
            "register_ica": "true"
          }
        }
      }
    }
  }
]
```

Here it is important that the `timeout_timestamp` or `timeout_height` takes into account the proposal's end time. Alternatively both of the timeout values can be set to `0`.

Further it is important that the field `register_ica` is set to `"true"` the first time you create an AutoTX. This will register an Interchain Account from Intento on the destination chain.

A message in the `action["msgs"]` array can be like the following MsgSend:

```json
{
  "@type": "/cosmos.bank.v1beta1.MsgSend",
  "amount": [
    {
      "amount": "70",
      "denom": "stake"
    }
  ],
  "from_address": "ICA_ADDR",
  "to_address": "some_destination_chain_address"
}
```

For the DAO service provider example, we can have a [SwapAndSendTo](https://github.com/Wasmswap/wasmswap-contracts/blob/main/src/msg.rs#:~:text=%7D%2C-,SwapAndSendTo,-%7B) in a [MsgExecuteContract](https://github.com/CosmWasm/wasmd/blob/main/proto/cosmwasm/wasm/v1/tx.proto) message to swap TOKEN2 for TOKEN1 locally on `Neutron`. We can specify the out recipient the recipient's address.

```json
{
  "@type": "/cosmos.authz.v1beta1.MsgExec",
  "msgs": [
    {
      "@type": "/cosmwasm.wasm.v1.MsgExecuteContract",
      "msg": {
        "swap_and_send_to": {
          "input_token": "TOKEN2",
          "min_token": "500",
          "recipient": "neutron1_address"
        }
      },
      "sender": "neutron1_address_dao",
      "contract": "neutron1_address_swap_contract",
      "funds": []
    }
  ],
  "grantee": "ICA_ADDR" //TIP: leave "ICA_ADDR" as is
}
```

:::tip write `ICA_ADDR` as a `grantee` or any other field in the message and Intento will parse the to-be defined Interchain Account address. This is transformed when the first execution happens.
:::

We can also perform [MsgSwapExactAmountOut](https://github.com/osmosis-labs/osmosis/blob/main/proto/osmosis/gamm/v1beta1/tx.proto#:~:text=message-,MsgSwapExactAmountOut,-%7B) to `Osmosis` or `Dymension` given that there are funds on the DAO's Intento-Osmosis Interchain Account on Osmosis. To send funds for execution, the DAO can make a proposal for `MsgTransfer` that funds the Intento Interchain Account. The DAO then controls an action on Intento that controls the swap.

### 2. Setting Up Permissions and Funds

![example2](@site/docs/images/connected_chain/from_connected_chain_flow3.png)

<!-- ![proposal2](@site/docs/images/ics20/daodao_proposal2.png)
*Proposal 2 can be the following and should contain the ICA_ADDR, which can found by querying the interchain account address.*  -->

The ICA_ADDR should also have some [funds to execute](#paying-for-fees) and have the proper [permissions](#setting-permissions) set up

:::tip write ICA_ADDR as a `sender` or any other field in an AutoTX and Intento will parse the to-be defined Interchain Account address.
:::

#### Paying for fees

Action fees on Intento are paid from the funds sent in the ICS20 message, so that you do not need to worry about it.

**Action on Intento**

If you want to execute actions locally on Intento, set `action["connection_id"]` to `""`. The funds sent with the ICS20 transfer will be used to pay for Action fees.

**With execution over IBC**

To execute actions on the `destination chain`, the `Action Account` on the destination chain should be funded with `destination chain`'s fee token.
The action can be funded by sending tokens to the account on `destination chain` before the first trigger. <!-- Fees can  granting the ICA to use fees from another account using a Feegrant. -->

Alternatively, you can submit an additional proposal to send tokens ( `MsgSend` or `MsgTransfer`) to `destination_chain` with the `destination_chain_ica_address` as the receiver.
For this, a message needs to be sent.

<!-- 
Here's an example:

```json
[
  {
    "stargate": {
      "typeUrl": "/ibc.applications.transfer.v1.MsgTransfer",
      "value": {
        "source_port": "transfer",
        "source_channel": "channel_source_chain__destination_chain",
        "token": { "denom": "ibc/....", "amount": "10" },
        "sender": "dao_address",
        "receiver": "destination_chain_ica_address",
        "timeout_height": "0",
        "timeout_timestamp": "0",
        "cid": "connection-0123"
      }
    }
  }
]
``` 
-->

:::tip Tip: if the destination chain is the same as the source chain, the DAO can create a FeeGrant for the `destination_chain_ica_address`.
:::

#### Setting permissions

If your action is on a Cosmos chain, you can set an AuthZ grant through `MsgGrant`. If you action is going to be on an EVM chain you can give the action account `allowance` for an ERC20 token or `approval` in case of an NFT.

<!--
#### Setting permissions

Alternatively, you can create a FeeGrant for the `destination_chain_ica_address` when the DAO is present on destination chain. You can also allow the ICA address to execute transactions for the DAO using `AuthZ`. [Here](https://blog.cosmos.network/secret-powers-what-are-the-authz-and-fee-grant-modules-c57d0e808794) is an article explaining the power of both `FeeGrant` and `AuthZ`.

Initiating an action for the first time locally on Neutron can be done as follows:

 1. Proposal 1 Message 1: IC20 token tranfser: Submit AutoTX on Intento.  In the `action` memo, `msgs` should contain AuthZ's MsgExec with `MsgExecuteContract`.
 2. Proposal 2 Message 1: Allow the DAO's NTRN tokens to be used for fees by `destination_chain_ica_address` with a FeeGrant.
  Proposal 2 Message 2: AuthZ grant to ICA to execute `MsgExecuteContract`.

 For a DAO swapping tokens on a recurring basis on Osmosis for the first time, 2 proposals should be made.

 1. Proposal 1 Message 1: IC20 token tranfser: Submit AutoTX on Intento. In the `action` memo, `msgs` should contain AuthZ's MsgExec with `MsgSwapExactAmountOut`
 2. Proposal 2 Message 1: Send OSMO tokens to using `MsgTransfer` to `destination_chain_ica_address`.
  Proposal 2 Message 2: AuthZ grant to ICA to execute `MsgSwapExactAmountOut`.

<!-- FeeGrant typeURL:  `/cosmos.feegrant.v1beta1.Grant`
AuthZ typeURL: `/cosmos.authz.v1beta1.MsgGrant` -->

when the funds and grants are in place from a previous action on that particular Interchain Account only one proposal action containing the `action` memo has to be specified.

![example4](@site/docs/images/connected_chain/from_connected_chain_flow4.png)

### 3. Managing Actions

On [TriggerPortal](https://triggerportal.zone/), you can create, view , update and manage your actions.
