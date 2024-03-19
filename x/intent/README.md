
# Intent Module

This module is used to automate processes across chains. It is essentially a custom Interchain Accounts authentication module
With 'MsgSubmitAction', time-based interchain account message calls are scheduled.

This is useful for Cosmos and smart contract developers wanting to automate message calls without using hard-to-integrate and unreliable off-chain bots.

Use cases include payroll, dollar-cost averaging, managing contract workflows. Respective messages for these actions are: MsgSend on any Cosmos SDK chain; MsgSwapExactAmountIn on Osmosis and MsgExecuteContract on CosmWasm chains.

Local messages to this chain can also be scheduled. This can be used to Autocompound tokens to any validator, stream local INTO tokens or even stream IBC Transfers.

## Automating user assets

The trustless nature of on-chain automation and trustless bridging through IBC have way better trust assumptions compared to bots. What it also brings to the table is asserting whom the message caller will be. By knowing the message caller it is possible to authorize the message caller to automate user assets. This is done by creating an grant in the Authz module of host chain. This allows for payroll services without pre-funding a contract & transfers of assets based on time. This is only available using on-chain automation.

## Fees

To encourage valuable transactions and timely execution, fees are implemented. These are designed to prevent spam, network congestion and reward validators for providing computational resources.

1. FixedFee to the community pool.  
2. FlexFee to the proposing validator

Governance can decide on how these fees are set. The FixedFee goes to the community pool to create network growth. Funds can be allocated towards ecosystem development, new tooling, relayer participation, and educational content.

### Fixed Fee

The Fixed Fee is applied per message basis. An Action can have multple messages to be executed, with a maximum of 9. This number can be increased when desired.

### FlexFee Fee

The FlexFee is timedependent fee calculated for each Action period on a minute basis.

### Refunds

When funds are sent along MsgSubmitAction, fees are deducted from the sender account. When specified, fees are deducted from a fund address. When funds are unavailable, execution will not take place. Funds left on a fund address are automatically returned to the owner after the last execution. A small community pool commision is applied.

## Relayer Rewards

Relayers are vital to well-functioning of this module for ensuring timely, reliable execution. To incentivize relayers, rewards are minted in the mint module, allocated from the alloc module and sent to the Intent module. For acknoledging a succesfull IBC packet containing Action messages, relayers are incentivized. Relayer rewards are specified based on the category of message. The category are as follows: SDK message, WASM message and Osmosis message. AuthZ messages can perform authorized actions on behalf of a user such as recurring transactions and reward claims. On Osmosis users can perform DCA strategies and withdrawal automatically after an unbonding period ends.  WASM smart contract calls are for developers automating their dApps and users that want to automate their smart contract tasks.

## Parameters

A number of automation-related parameters can be adjusted. Parameters can be adjusted by governance to ensure that fees and rewards are fair. The default values are the following:

```golang
const (
 // ActionFundsCommission percentage to distribute to community pool for leftover balances (rounded up)
 DefaultActionFundsCommission int64 = 2 //2%
 // ActionConstantFee fee to prevent spam of auto messages, to be distributed to community pool
 DefaultActionConstantFee int64 = 5_000 // 0.005trst
 // ActionFlexFeeMul is the denominator for the gas-dependent flex fee to prioritize auto messages in the block, to be distributed to validators
 DefaultActionFlexFeeMul int64 = 3 // 3% of minutes for a given period as uinto (1_000m = 20uinto)
 // RecurringActionConstantFee fee to prevent spam of auto messages, to be distributed to community pool
 DefaultRecurringActionConstantFee int64 = 5_000 // 0.005trst
 // Default max period for a Action that is self-executing
 DefaultMaxActionDuration time.Duration = time.Hour * 24 * 366 * 2 // a little over 2 years
 // MinActionDuration sets the minimum duration for a self-executing Action
 DefaultMinActionDuration time.Duration = time.Second * 60
 // MinActionInterval sets the minimum interval self-execution
 DefaultMinActionInterval time.Duration = time.Second * 60
 // DefaultRelayerReward for a given action type
 DefaultRelayerReward int64 = 10_000 //0.01trst
)

```

## BeginBlocker

At the beginning of each block, the BeginBlocker checks if there are Actions that are set for automation. The BeginBlocker is used as it can best proxy the execute time set at MsgSubmitAction

## ICS20 Transfer Middleware

This module contains ICS20 transfer middleware. A MsgRegisterAccountAndSubmitAction or a MsgSubmitAction can be derrived from the memo field. This is useful for contract callers on other chains. This way DAOs and other entities can create actions using an ICS20-standard transaction. This middleware is based on the wasm_hooks implementation on [Osmosis](https://github.com/osmosis-labs/osmosis/tree/main/x/ibc-hooks).

This is useful for a variety of use cases.

The mechanism enabling this is a `memo` field on every ICS20 transfer packet as of [IBC v3.4.0](https://medium.com/the-interchain-foundation/moving-beyond-simple-token-transfers-d42b2b1dc29b).
ics_middleware.go is IBC middleware that parses an ICS20 transfer, and if the `memo` field is of a particular form, creates a trigger by parsing and handling a SubmitAction message. 

We now detail the field s format for `action`.

* Sender: We cannot trust the sender of an IBC packet, the counterparty chain has full ability to lie about it.
We cannot risk this sender being confused for a particular user or module address. 

* Owner: This field is directly obtained from the ICS-20 packet metadata and should equal the ICS20 transfer recipient. If unspecified, a placeholder is made from the ICS20 sender and channel.
* Msg: This field should be directly obtained from the ICS-20 packet metadata.
* Funds: This field is set to the amount of funds being sent over in the ICS 20 packet. One detail is that the denom in the packet is the counterparty chains representation of the denom, so we have to translate it to Intento's representation.

So our constructed message for MsgSubmitAction will contain the following:

```go
msg := MsgSubmitAction{
 // If let unspecified, owner is the actor that submitted the ICS20 message and a placeholder only
 Owner: packet.data.memo["action"]["owner"] OR "into1-hash-of-channel-and-sender",
 // Array of Msg json encoded, then transformed into a proto.message
 Msgs: packet.data.memo["action"]["msgs"],
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
        "receiver": "into1address",
         "memo": {
           "action": {
            "owner": "into1address ", //optional
              "msgs": [{
                "@type":"/cosmos.somemodule.v1beta1.sometype",
                //message values in JSON format
            }],
            "duration":"111h",
            "interval":"11h",
            "start_at":"11h",
            "label":"my_label",
            "cid":"connection-0", //optional, omit or leave blank in case local INTO message.
            "register_ica": "false"//optional, set to true to register interchain account
        },
        //"version":""//optional, will attempt to register account when filled (this will never override any existing ICA address)
    }
}
}
```

An ICS20 packet is formatted correctly for submitting an action if the following all hold:

* `memo` is not blank
* `memo` is valid JSON
* `memo` has at least one key, with value `"action"`
* `memo["action"]` has exactly the entries mentioned above
* `memo["action"]["msgs"]` is an array with valid JSON SDK message objects with a key "@type" and sdk message values
* `receiver == memo["action"]["owner"]`. Optional, owner should equal receiver and is the address that receives remaining fee balance after execution ends. 
* `memo["action"]["connection_id"]`is a valid connectionID on INTO -> Destination chain, or blank/empty for local INTO execution of the message.
* `memo["action"]["register_ica"]` can be added, and true to register an ICA. 

Fees are paid with a newly generated and Action specific fee account.

If an ICS20 packet does not contain a memo containing "action", a regular MsgTransfer takes place.
If an ICS20 packet is directed towards autoTx, and is formated incorrectly, then it returns an error.
