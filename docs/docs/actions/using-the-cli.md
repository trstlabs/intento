---
order: 4
title: Using the CLI
description: How to set up actions using the CLI
---

### Installing the Intento deamon

Download `intentod` from [the main repo](https://github.com/trstlabs/intento) with `wget https://github.com/trstlabs/intento.git` and perform `make deb` to build a binary. Then install the binary with `sudo apt install ./Intento_VERSION_amd64.deb -y`. Test if it is installed correctly by running `intentod`.

You can configure the RPC endpoint with `intentod config node https://endpoint:port` or start the chain locally by running the `./inittest.sh` script.

## Local Action execution

You can perform local execution, for example for recurring payments, using the following format:

```bash

intentod tx intent submit-action  '{
    "@type":"/cosmos.bank.v1beta1.MsgSend",
    "amount": [{
        "amount": "70",
        "denom": "uinto"
    }],
    "from_address": "into12gxmzpucje8aflw2vz45rv8x4nyaaj3rp8vjh03dulehkdl5fu6s93ewkp",
    "to_address": "into1ykql5ktedxkpjszj5trzu8f5dxajvgv95nuwjx"
}' --duration 16h --interval 60s --keyring-backend test -y --from b --fees 600uinto
```

Fees are denominated in `uinto`. Similar to other Cosmos chains, it is a 6-decimal denomination of `INTO`.
`intentod tx intent submit-action` takes one argument, a JSON-encoded Cosmos message. `@type` can be retrieved from the destination chain's `/proto` directory.
For parsing a message to automate with the CLI, use snake case. `fromAddress` will be `from_address`.
:::warning the Cosmos message type should be a registered type on Intento
:::


This "type url" is based of the proto package name, in this example it is `cosmos.bank.v1beta1` and the name of the message is `MsgSend`. TRST Labs has registerd types of CosmWasm types and the Osmosis types. Any message should be registered before it can be encoded correctly on Intento.

`--duration` and `--interval` flags must contain either `s`,`m`or `h`, which stand for seconds, minutes and hours. It can be a combination of all of these. Like so: `48h30m30s`.
`--start_at` can be specified to create a custom start time for the first execution. It is in nanoseconds since UNIX epoch such as `1678199141`.

When an Action should execute once, `duration` will be the time between now and the execution time.
If the Action is recurring, `duration` will be the time execution will be unavailable.

## With Interchain Accounts

Message flow is similar to the above. A *connection_id* should be specified. You can find all available connections with `intentod q ibc connection connections`.

### What is a connection and how does it differ from a channel?

A connection is the direct link between two chains, while a channel is a specific path within that connection that enables communication between two specific modules on different chains. Using a connection, multiple channels can be created. Both connections and channels are usually different for on source chain and on the destination chain.

Channels have ports such as `icacontroller-into1...` or `transfer`. By using a connection, Intento creates an IBC port ID for your address and the destination chain.

This source port is denoted as `icacontroller-into1...`. For an Interchain Account host, the destination port is `host`.
A connection is needed to create the channels and ports.

### Register and Submit Action

When you submit an Action for the first time using interchain Accounts, you should perfrom `register-ica-and-submit-action` instead of `submit-action`.

Alternatively you can `intentod tx intent register` and specify the `--connection-id` flag.

You can submit an Action with the following command:

```bash
intentod tx intent submit-action  '{
    "@type":"/cosmos.bank.v1beta1.MsgSend",
    "amount": [{
        "amount": "70",
        "denom": "uinto"
    }],
    "from_address": "ICA_ADDR",
    "to_address": "into1ykql5ktedxkpjszj5trzu8f5dxajvgv95nuwjx"
}' --duration 16h --interval 60s --keyring-backend test -y --from b --fees 600uinto --connection-id connection-0 --retries 2
```

:::tip use ICA_ADDR as a `from_address` or any other field in an AutoTX and Intento will parse the to-be defined Interchain Account Address.
:::

```bash
# query the AutoTXs to see if it worked. After a time-based execution the action history should update
intentod query intent list-actions-by-owner $WALLET_1
```

You can get the Interchain Account address with the following command:

```bash
export ICA_ADDR=$(intentod q intent interchainaccounts $WALLET_1 connection-0 -o json | jq -r '.interchain_account_address') && echo $ICA_ADDR
```

:::tip to get more info you can query the events related to the packet for packet and message info `intentod q interchain-accounts host packet-events channel-0 1 --node tcp://localhost:36657`
:::

## With AuthZ

Message flow is similar to interchain acccounts. For using AuthZ, it is required to create a grant on the host chain. In the messages below, the grantee is the interchain account address.

For this we use another local Intento chain. On a host chain you can create a grant with a command like the following:

```bash
intentod tx authz grant $ICA_ADDR generic --msg-type /cosmos.staking.v1beta1.MsgDelegate --keyring-backend test -y --from b --fees 600uinto  --node tcp://localhost:36657  --chain-id intentodev-2
```

Hereafter, on Intento you can execute the message with the ICA as the grantee and the Msgs array in the MsgExec should include the `--msg-type` from before.

```bash
intentod tx intent submit-action  '{
    "@type":"/cosmos.authz.v1beta1.MsgExec",
    "msgs": [{
    "@type": "/cosmos.staking.v1beta1.MsgDelegate",
    "amount": {
        "amount": "70",
        "denom": "uinto"
    },
    "delegator_address": "into1ykql5ktedxkpjszj5trzu8f5dxajvgv95nuwjx",
    "validator_address": "trustvaloper1q6k0w4cejawpkzxgqhvs4m2v6uvdzm6jhmz5jy"
}],
    "grantee": "ICA_ADDR"
}' --duration 4h --interval 60s --keyring-backend test -y --from b --fees 600uinto --connection-id connection-0
```


## Configuration

You can set and the following flags:

| Flag                           | Description                                                                               | Example Value     |
| ------------------------------ | ----------------------------------------------------------------------------------------- | ----------------- |
| `connection-id`                | Identifier for the connection end on the controller chain.                                | `connnection-123` |
| `host-connection-id`           | Identifier for the controller chain channel version.                                      | `connnection-456` |
| `label`                        | Custom label for the Action, such as transaction type or operation name. Optional.        | `AutoTransfer`    |
| `duration`                     | Duration for which the Action remains active. Optional.                                   | `48h`             |
| `interval`                     | Custom interval between Action executions. Optional.                                      | `2h`              |
| `start-at`                     | Custom start time for the Action in UNIX time format. Optional.                           | `1625097600`      |
| `fee-funds`                    | Coins sent to limit the fees incurred during transaction execution. Optional.             | `100atom`         |
| `end-at`                       | Custom end time for the Action in UNIX time format. Optional.                             | `1625184000`      |
| `updating-disabled`            | Disables future updates to the Action configuration.                                      | `false`           |
| `save-msg-responses`           | Saves message responses to transaction history for Cosmos SDK v0.46+ chains only.         | `true`            |
| `fallback-to-owner-balance`    | Uses owner's balance as fallback for transaction fees if `fee-funds` are insufficient.    | `true`            |
| `stop-on-success`              | Stops execution of the Action after a successful transaction.                             | `true`            |
| `stop-on-failure`              | Stops execution of the Action after a failed transaction.                                 | `true`            |
| `stop-on-success-of`           | Stops execution if a specified action succeeds. Optional and requires custom logic.       | `23,58`           |
| `stop-on-failure-of`           | Stops execution if a specified action fails. Optional and requires custom logic.          | `4536,234`        |
| `skip-on-success-of`           | Skips the next execution if a specified action succeeds. Optional, requires custom logic. | `234,234`         |
| `skip-on-failure-of`           | Skips the next execution if a specified action fails. Optional, requires custom logic.    | ``3456,12`        |
| `reregister_ica_after_timeout` | Allows Action to continue execution after an IBC channel timeout. Recommended.            | `true`            |
