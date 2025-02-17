---
sidebar_position: 5
title: From a connected chain
description: How to setup flows from a connected chain
---

## Setting up Flows

In the previous step we showed how the flow process looks like by submitting an flow on Intento. You can do this with the [TriggerPortal](https://triggerportal.zone) interface, a IntentoJS front-end integration or locally through the CLI.

In addition, you can also submit an flow from another chain using the [ICS20 standard](https://github.com/cosmos/ibc-go/blob/main/docs/apps/transfer/messages.md).

### Interchain Accounts

Users and entities on Cosmos SDK chains may be able set up [interchain account](https://tutorials.cosmos.network/academy/3-ibc/8-ica.html) to Intento and submit the flows using `MsgSubmitFlow`. It is also easy to deploy using the ICS20 standard.

### ICS20 Standard

With Intento’s ICS20 transfer middleware, you send a transfer token memo on a chain, and Intento will convert the token transfer to an flow submission.
ICS20 is an interchain standard that enables the transfer of fungible tokens between independent blockchains. It is a protocol that defines a standard interface for token transfers across different blockchains that implement the Inter-Blockchain Communication (IBC) protocol.

Using ICS20, accounts on connected chains can create flows. This can be done by specifying flow details in the memo field of an ICS20 transfer message. Upon receiving this message, Intento's IBC hooks transforms this into a submit flow message.

This is useful for DAOs and other decentralized organizations on any connected chain. They can safely and reliably execute on Intento's connected chains. For DAOs, this gives certainty to stakeholders, whilst also reducing manual work on governance proposals.

## For DAOs

Setting up an flow on a connected chain can be particularly useful for DAOs. Using this middleware, DAOs can now automate tasks not only on their chain but also on any chain connected to Intento. What can DAO's do with this? DAOs can orchestrate periodic token swaps, payroll, payment in installments amongst other scheduled flows. These flows can be performed in one proposal, which normally require periodically voting on individual proposals. This normally requires manual flow from the proposer and DAO participants.

## ICS20 Middleware

A MsgRegisterAccountAndSubmitFlow or a MsgSubmitFlow can be derrived from the memo field in the ICS20 transfer message.

![ics20](@site/docs/images/connected_chain/from_connected_chain.png)

Our custom middleware is based on the wasmhooks implementation on [Osmosis](https://github.com/osmosis-labs/osmosis/tree/main/x/ibc-hooks).

The mechanism enabling this is a `memo` field on every ICS20 transfer packet as of [IBC v3.4.0](https://medium.com/the-interchain-foundation/moving-beyond-simple-token-transfers-d42b2b1dc29b).

ics_middleware.go is IBC middleware that parses an ICS20 transfer, and if the `memo` field is of a particular form, it creates an flow by parsing and handling a SubmitFlow message.

These are the fields for `flow` that are derived from the ICS20 message:

- **Owner**: This field is directly obtained from the ICS20 packet metadata and equals the ICS20 recipient. If unspecified, a placeholder is made from the ICS20 sender and channel.
- **Msg**: This field should be directly obtained from the ICS20 packet metadata.
- **FeeFunds**: This field is set to the amount of funds being sent over in the ICS20 packet. One detail here is that the denom in the packet is the source chains representation of the denom, this will be translated into INTO on Intento.

The constructed message for MsgSubmitFlow under the hood will look like:

```go
msg := MsgSubmitFlow{
 // If let unspecified, owner is the actor that submitted the ICS20 message and a placeholder only
 Owner: "into1-hash-of-channel-and-sender" OR packet.data.memo["flow"]["owner"],
 // Array of Msg json encoded, then transformed into a proto.message
 Msgs: packet.data.memo["flow"]["msgs"],
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
      "flow": {
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

An ICS20 packet is formatted correctly for submitting an flow if the following all hold:

- `memo` is not blank
- `memo` is valid JSON
- `memo` has at least one key, `"flow"`
- `memo["flow"]["msgs"]` is an array with valid JSON SDK message objects with a key "@type" and sdk message values
- `receiver == memo["flow"]["owner"]`. Optional, an owner can be specifed and is the address that receives remaining fee balance after execution ends.
- `memo["flow"]["cid"]`is a valid connection ID on INTO -> Destination chain, omit it for local INTO execution of the message.
- `memo["flow"]["register_ica"]` can be added, and true to register an ICA.

Fees are paid with a newly generated flow fee account.

If an ICS20 packet does not contain a memo containing "flow", a regular MsgTransfer takes place.
If an ICS20 packet is directed towards flow, and is formated incorrectly, then it returns an error.
Here’s an improved version of the example with clearer language, structure, and formatting:

---

## Example: DAO Integration

![ICS DAO](@site/docs/images/connected_chain/from_connected_chain_flow1.png)

### Overview

Using ICS20 to set up Intent-based flows is for users familiar with ICS20 transfers and Authz permissions. Our ICS20 middleware is designed to allow DAOs to submit flows, which can be done locally on Intento or on a destination chain. This can also be the source chain.

There are several caveats when setting up flows with ICS20. When automating on a destination chain for the first time, **two messages** are required to activate the flow. One flow sets up the flow and creates an Interchain Account on Intento, while the other sets permissions and sends funds to the Interchain Account on the destination chain.

In this example, we will demonstrate how a DAO can integrate with Intento by automating the payment process to a service provider.

### Scenario: DAO Payment to Service Provider

The `DAO` wishes to pay `Service Provider ABC` monthly for their services in `TOKEN1`. The DAO holds `TOKEN2` and `NTRN`.

_Service Provider ABC invoice example_

![DAODAO](@site/docs/images/connected_chain/daodao_proposal1.png)

In this case, the DAO triggers a recurring swap of `TOKEN2` for `TOKEN1` on the decentralized exchange "DEX" and automatically sends the tokens to `Service Provider ABC`. 

This setup is an ideal use case for Intento, as it involves asset movement between chains and accounts, while the DAO maintains control of the tokens. Since the flows are on-chain, the process is fully decentralized and does not require third-party trust.

The DAO can appoint an owner to manage the flow or use a placeholder account on Intento to remain in full control.

### Proposal Details

In this example, the DAO submits a proposal with the following name and description:

**Proposal Name**:

```md
[Trigger] Pay Service Provider ABC in TOKEN1
```

**Proposal Description**:

```md
Submit a flow to send TOKEN1 to "Service Provider ABC."

This flow will swap TOKEN2 for TOKEN1 on DEX "DEX" on the Destination Chain and automatically send these tokens to "Service Provider ABC." By performing this swap on a recurring basis, we achieve:

- Gradual selling pressure on TOKEN2
- Maintaining positive cash flow

This trigger on Intento automates asset workflows, ensuring liquidity and financial stability for the DAO.
```

Having sufficient liquidity ensures smooth operations and financial stability, allowing the DAO to meet its obligations and grow sustainably. Proper liquidity management helps to pay expenses, invest in growth, and strengthen the DAO's reputation in the market.

---

### 1. Submitting the Flow

![ICS DAO Flow Submission](@site/docs/images/connected_chain/from_connected_chain_flow2.png)

To submit the flow, the DAO creates a proposal with a custom message to execute. For a CosmWasm-based DAO like [DAO DAO](https://daodao.zone/) on Neutron, this message contains the ICS20 `MsgTransfer`. In the `memo`, the DAO provides flow details such as automated messages, time parameters, and a custom label.

Here’s an example of how the flow proposal message might look for a CosmWasm-based DAO:

```json
[
  {
    "stargate": {
      "typeUrl": "/ibc.applications.transfer.v1.MsgTransfer",
      "value": {
        "source_port": "transfer",
        "source_channel": "channel-to-intento",
        "token": { "denom": "ibc/....", "amount": "10" },
        "sender": "neutron1validbech32address",
        "receiver": "into1address",
        "timeout_height": "0",
        "timeout_timestamp": "0",
        "memo": {
          "flow": {
            "owner": "into1address",
            "msgs": [
              {
                "@type": "/someprefix.somemodule.someversion.sometype"
              }
            ],
            "duration": "111h",
            "interval": "11h",
            "start_at": "1677841601",
            "label": "my_label",
            "cid": "connection-0",
            "register_ica": "true"
          }
        }
      }
    }
  }
]
```

**Important Notes:**

- Ensure that `timeout_timestamp` or `timeout_height` reflects the proposal's end time, or set both to `0` if no timeout is needed.
- When creating the flow for the first time, make sure the `register_ica` field is set to `"true"`, which registers an Interchain Account on the destination chain.
  
A message within the `flow["msgs"]` array might look like this `MsgSend`:

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

Alternatively, the DAO can use a custom smart contract message like [SwapAndSendTo](https://github.com/Wasmswap/wasmswap-contracts/blob/main/src/msg.rs#:~:text=%7D%2C-,SwapAndSendTo,-%7B) in a [MsgExecuteContract](https://github.com/CosmWasm/wasmd/blob/main/proto/cosmwasm/wasm/v1/tx.proto) message to swap `TOKEN2` for `TOKEN1`.

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
  "grantee": "ICA_ADDR"
}
```

:::tip 
Write `ICA_ADDR` as a `grantee` or in other fields within the message, and Intento will parse the to-be-defined Interchain Account address.
:::

You can also use [MsgSwapExactAmountOut](https://github.com/osmosis-labs/osmosis/blob/main/proto/osmosis/gamm/v1beta1/tx.proto#:~:text=message-,MsgSwapExactAmountOut,-%7B) to swap tokens on decentralized exchanges like Osmosis or Dymension.

---

### 2. Setting Up Permissions and Funds

![ICS DAO Permissions](@site/docs/images/connected_chain/from_connected_chain_flow3.png)

The Interchain Account (ICA) must be properly funded and authorized before it can execute the flow.

#### Paying for Fees

Flow fees on Intento are automatically paid from the funds sent via ICS20. For flows executed on the destination chain, the `Flow Account` on the destination chain should be funded with the destination chain's fee token.

If needed, a separate proposal can be submitted to send tokens (via `MsgSend` or `MsgTransfer`) to the destination chain's ICA address.

:::tip
If the source and destination chains are the same, you can set up a FeeGrant for the ICA on the destination chain.
:::

#### Setting Permissions

For flows on Cosmos chains, set up an AuthZ grant using `MsgGrant`. For EVM chains, you can give the flow account `allowance` for ERC20 tokens or `approval` for NFTs.

Once funds and permissions are set up, only one flow proposal is required to trigger the flow.

---

### 3. Managing Flows

You can manage your flows on [TriggerPortal](https://triggerportal.zone/), where you can create, view, update, and control your flows.
