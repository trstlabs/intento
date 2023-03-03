---
order: 3
title: Interchain Setup
description: How to instantiate automation from another chain
---

## Setting up Automation 

In the previous step we showed how the AutoTx process looks like. We assumed you want to start automation by submitting an AutoTx on Trustless Hub. This can be through the [TriggerPørtal](triggerportal.netlify.app) interface, a TrustlessJS front-end integration or locally through a light client Command-Line Interface.

![msgflow](../images/msgflow.png)

However, You can also submit an AutoTX from another chain using the [ICS20 standard](https://github.com/cosmos/ibc-go/blob/main/docs/apps/transfer/messages.md). This is done through ICS20 transfer middleware. ICS20 is an interchain standard that enables the transfer of fungible tokens between independent blockchains within the Cosmos ecosystem. It is a protocol that defines a standard interface for token transfers across different blockchains that implement the Inter-Blockchain Communication (IBC) protocol.


![ics20](../images/ics20msgflow.png)

## Using ICS20 Middleware

By specifying AutoTx details in the memo, TRST will submit an AutoTx using the inputs provided. A MsgRegisterAccountAndSubmitAutoTx or a MsgSubmitAutoTx can be derrived from the memo field in the ICS20 transfer message. This is useful for users, entities and contract callers on other chains. This way DAOs and other entities can create autotxs using an ICS20-standard transaction.

Our custom middleware is loosely based on the wasmhooks implementation on [Osmosis](https://github.com/osmosis-labs/osmosis/tree/main/x/ibc-hooks).

The mechanism enabling this is a `memo` field on every ICS20 transfer packet as of [IBC v3.4.0](https://medium.com/the-interchain-foundation/moving-beyond-simple-token-transfers-d42b2b1dc29b).
ics_middleware.go is IBC middleware that parses an ICS20 transfer, and if the `memo` field is of a particular form, creates a trigger by parsing and handling a SubmitAutoTx message.

We now detail the field s format for `auto_tx`.

* Sender: We cannot trust the sender of an IBC packet, the counterparty chain has full ability to lie about it.
We cannot risk this sender being confused for a particular user or module address.

* Owner: This field is directly obtained from the ICS-20 packet metadata and equals the ICS20 recipient. If unspecified, a placeholder is made from the ICS20 sender and channel.
* Msg: This field should be directly obtained from the ICS-20 packet metadata.
* Funds: This field is set to the amount of funds being sent over in the ICS 20 packet. One detail is that the denom in the packet is the counterparty chains representation of the denom, so we have to translate it to TRST's representation.

So our constructed message for MsgSubmitAutoTx will contain the following:

```go
msg := MsgSubmitAutoTx{
 // If let unspecified, owner is the actor that submitted the ICS20 message and a placeholder only
 Owner: packet.data.memo["auto_tx"]["owner"] OR "trust1-hash-of-channel-and-sender",
 // Array of Msg json encoded, then transformed into a proto.message
 Msgs: packet.data.memo["auto_tx"]["msgs"],
 // Funds coins that are transferred to the owner
 FeeFunds: sdk.NewCoin{Denom: ibc.ConvertSenderDenomToLocalDenom(packet.data.Denom), Amount: packet.data.Amount}
```

### ICS20 packet structure

So given the details above, we propogate the implied ICS20 packet data structure.
ICS20 is JSON native, so we use JSON for the memo format.

```json
{
    //... other ibc fields that we don't care about
    "data":{
       "denom": "denom on counterparty chain (e.g. uatom)",
        "amount": "1000",
        "sender": "...", // ignored
        "receiver": "A TRST addr prefixed with trust1",
         "memo": {
           "auto_tx": {
            "owner": "trust1address", //optional
              "msgs": [{
                "@type":"/cosmos.somemodule.v1beta1.sometype",
                //message values in JSON format
            }],
            "duration":"111h",
            "interval":"11h",
            "start_at":"11h",
            "label":"my_label",
            "connection_id":"connection-0", //optional, omit or leave blank in case local TRST message.
            "register_ica": "false"//optional, set to true to register interchain account
        },
        //"version":""//optional, will attempt to register account when filled (this will never override any existing ICA address)
    }
}
}
```

An ICS20 packet is formatted correctly for submitting an auto_tx if the following all hold:

* `memo` is not blank
* `memo` is valid JSON
* `memo` has at least one key, with value `"auto_tx"`
* `memo["auto_tx"]` has exactly the entries mentioned above
* `memo["auto_tx"]["msgs"]` is an array with valid JSON SDK message objects with a key "@type" and sdk message values
* `receiver == memo["auto_tx"]["owner"]`. Optional, owner is the address that receives remaining fee balance after execution ends.
* `memo["auto_tx"]["connection_id"]`is a valid connectionID on TRST -> Destination chain, or blank/empty for local TRST execution of the message.
* `memo["auto_tx"]["register_ica"]` can be added, and true to register an ICA.

Fees are paid with a newly generated and AutoTx specific fee account.

If an ICS20 packet does not contain a memo containing "auto_tx", a regular MsgTransfer takes place.
If an ICS20 packet is directed towards autoTx, and is formated incorrectly, then it returns an error.

## DAO Integration

Submit a proposal with a custom message to execute. For a CosmWasm-based DAO like DAO DAO DAOs on Juno, the message could look like the following:
Here it is important that the timeout takes into account the proposal's end time. Alternatively it can be set to zero.

```json
[
  {
    "stargate": {
      "typeUrl": "/ibc.applications.transfer.v1.MsgTransfer",
      "value": {
        "source_port": "transfer",
        "source_channel": "channel-123",
        "token": "IBC/XZY", //should be TRST, which will pay for fees
        "sender": "juno1validbech32address",
        // the recipient address on Trustless Hub
        "receiver": "trust1address",//will be omitted when auto_tx["owner"] in memo is blank
        // Timeout height relative to the current block height.
        // The timeout is disabled when set to 0.
        "timeout_height": "",
        // Timeout timestamp in absolute nanoseconds since unix epoch.
         // The timeout is disabled when set to 0.
        "timeout_timestamp": "0",
        "memo": {
           "auto_tx": {
            "owner": "trust1address", //optional
              "msgs": [{
                "@type":"/someprefix.somemodule.someversion.sometype",
                //message values in JSON format
            }],
            "duration":"111h",
            "interval":"11h",
            "start_at":"11h",
            "label":"my_label",
            "connection_id":"connection-0", //optional, omit or leave blank in case local TRST message.
            "register_ica": "false"//optional, set to true to register interchain account
           },
        }
      }
    }
  }
]
```
