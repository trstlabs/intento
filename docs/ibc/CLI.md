---
order: 3
title: CLI
description: How to use the CLI and hermes to transact using IBC
---


## IBC setup

Two local Trustless Hub chains can communicate with each other via a Hermes relayer

:::tip You can always switch between chains using `trstd config node RPC_ENDPOINT`.
:::

## Build

```bash
docker build -f hermes.Dockerfile . --tag hermes:test
```

### Run

```bash
docker compose up
```

### ICS20 Transfers

Assuming you have a key 'a' which is not the relayer's key,
from localhost:

```bash
a_mnemonic="grant rice replace explain federal release fix clever romance raise often wild taxi quarter soccer fiber love must tape steak together observe swap guitar"

echo $a_mnemonic | trstd keys add a --recover

trstd add-genesis-account "$(trstd keys show -a a)" 100000000000000utrst

# be on the source network (trstdev-1)
trstd config node http://localhost:26657

# check the initial balance of a
trstd q bank balances trust1q6k0w4cejawpkzxgqhvs4m2v6uvdzm6j2pk2jx

# transfer to the destination network
trstd tx ibc-transfer transfer transfer channel-0 trust1ykql5ktedxkpjszj5trzu8f5dxajvgv95nuwjx 2utrst --from a

# check a's balance after transfer
trstd q bank balances trust1q6k0w4cejawpkzxgqhvs4m2v6uvdzm6j2pk2jx

# switch to the destination network (trstdev-2)
trstd config node http://localhost:36657

# check that you have an ibc-denom
trstd q bank balances trust1ykql5ktedxkpjszj5trzu8f5dxajvgv95nuwjx # should have 1 ibc denom
```

### Interchain Accounts

Message flow for using interchain accounts to directly controll an interchain account on the host chain.

```bash

# register account for address after channel is open
 trstd tx autoibctx register --connection-id connection-0 --keyring-backend test -y --from b --fees 600utrst

# wait for the interchain account to initialize and query the ICA address
trstd q autoibctx interchainaccounts trust1ykql5ktedxkpjszj5trzu8f5dxajvgv95nuwjx connection-0

# Register an interchain account on behalf of WALLET_1 where chain test-2 is the interchain accounts host
icad tx autoibctx register --from $WALLET_1 --connection-id connection-0 --chain-id test-1 --home ./data/test-1 --node tcp://localhost:16657 --keyring-backend test -y

# Query the address of the interchain account
export ICA_ADDR=$(icad query autoibctx interchainaccounts connection-0 $WALLET_1 --home ./data/test-1 --node tcp://localhost:16657 -o json | jq -r '.interchain_account_address') && echo $ICA_ADDR //trust12gxmzpucje8aflw2vz45rv8x4nyaaj3rp8vjh03dulehkdl5fu6s93ewkp

# send balance to ICA on host chain to provide an initial balance to execute transactions (replace node and to_address here)
trstd  tx bank send trust1ykql5ktedxkpjszj5trzu8f5dxajvgv95nuwjx $ICA_ADDR 10000utrst --node tcp://localhost:36657 --keyring-backend test -y --from b --fees 600utrst --chain-id trstdev-2

# replace delegator_address to ICA address and submit tx
trstd tx autoibctx submit-tx  '{
    "@type": "/cosmos.staking.v1beta1.MsgDelegate",
    "amount": {
        "amount": "70",
        "denom": "utrst"
    },
    "delegator_address": "trust12gxmzpucje8aflw2vz45rv8x4nyaaj3rp8vjh03dulehkdl5fu6s93ewkp",
    "validator_address": "trustvaloper1q6k0w4cejawpkzxgqhvs4m2v6uvdzm6jhmz5jy"
}' --keyring-backend test -y --from b --fees 600utrst --connection-id connection-0
 
# check balance (should be 1000-70=9931)
trstd q bank balances $ICA_ADDR  --node tcp://localhost:36657

# (optional) check staking delegations
trstd q staking delegations-to trustvaloper1q6k0w4cejawpkzxgqhvs4m2v6uvdzm6jhmz5jy --node tcp://localhost:36657

# (optional) you can query the events related to the packet for packet and message info. The command takes a channel id and packet sequence
trstd q interchain-accounts host packet-events channel-0 1 --node tcp://localhost:36657 

# (debug) query to check the channel, make sure channel is STATE_OPEN 
trstd query ibc channel channels
```

### AutoTXs with Interchain Accounts

Message flow is similar to interchain acccounts, but with `submit-auto-tx` instead of `submit-tx`. It takes several parameters like duration and interval as flags. Duration is required parameter.

```bash

# replace delegator_address to ICA address and submit tx
trstd tx autoibctx submit-auto-tx  '{
    "@type": "/cosmos.staking.v1beta1.MsgDelegate",
    "amount": {
        "amount": "70",
        "denom": "utrst"
    },
    "delegator_address": "trust12gxmzpucje8aflw2vz45rv8x4nyaaj3rp8vjh03dulehkdl5fu6s93ewkp",
    "validator_address": "trustvaloper1q6k0w4cejawpkzxgqhvs4m2v6uvdzm6jhmz5jy"
}' --duration 60s --keyring-backend test -y --from b --fees 600utrst --connection-id connection-0

trstd tx autoibctx submit-auto-tx  '{
    "@type":"/cosmos.bank.v1beta1.MsgSend",
    "amount": [{
        "amount": "70",
        "denom": "utrst"
    }],
    "from_address": "trust12gxmzpucje8aflw2vz45rv8x4nyaaj3rp8vjh03dulehkdl5fu6s93ewkp",
    "to_address": "trust1ykql5ktedxkpjszj5trzu8f5dxajvgv95nuwjx"
}' --duration 16h --interval 60s --keyring-backend test -y --from b --fees 600utrst --connection-id connection-0 --retries 2

# query the AutoTXs to see if it worked. After a time-based execution the auto-tx history should update
trstd q autoibctx list-auto-txs-by-owner trust1ykql5ktedxkpjszj5trzu8f5dxajvgv95nuwjx

# to get more info you can query the events related to the packet for packet and message info
trstd q interchain-accounts host packet-events channel-0 1 --node tcp://localhost:36657 

```

### Trustless Automation of user-based tasks AutoTXs with Interchain Accounts and AuthZ

Message flow is similar to interchain acccounts but using AuthZ. It is required to create a grant on the host chain. In the messages below, the grantee is the interchain account address.

```bash


trstd tx authz grant $ICA_ADDR generic --msg-type /cosmos.staking.v1beta1.MsgDelegate --keyring-backend test -y --from b --fees 600utrst  --node tcp://localhost:36657  --chain-id trstdev-2

trstd tx autoibctx submit-auto-tx  '{
    "@type":"/cosmos.authz.v1beta1.MsgExec",
    "msgs": [{
    "@type": "/cosmos.staking.v1beta1.MsgDelegate",
    "amount": {
        "amount": "70",
        "denom": "utrst"
    },
    "delegator_address": "trust1ykql5ktedxkpjszj5trzu8f5dxajvgv95nuwjx",
    "validator_address": "trustvaloper1q6k0w4cejawpkzxgqhvs4m2v6uvdzm6jhmz5jy"
}],
    "grantee": "trust12gxmzpucje8aflw2vz45rv8x4nyaaj3rp8vjh03dulehkdl5fu6s93ewkp"
}' --duration 4h --interval 60s --keyring-backend test -y --from b --fees 600utrst --connection-id connection-0



trstd tx authz grant $ICA_ADDR generic --msg-type /cosmos.staking.v1beta1.MsgUndelegate --keyring-backend test -y --from b --fees 600utrst  --node tcp://localhost:36657  --chain-id trstdev-2

trstd tx autoibctx submit-auto-tx  '{
    "@type":"/cosmos.authz.v1beta1.MsgExec",
    "msgs": [{
   "@type": "/cosmos.staking.v1beta1.MsgUndelegate",
    "amount": {
        "amount": "70",
        "denom": "utrst"
    },
    "delegator_address": "trust1ykql5ktedxkpjszj5trzu8f5dxajvgv95nuwjx",
    "validator_address": "trustvaloper1q6k0w4cejawpkzxgqhvs4m2v6uvdzm6jhmz5jy"
}],
    "grantee": "trust12gxmzpucje8aflw2vz45rv8x4nyaaj3rp8vjh03dulehkdl5fu6s93ewkp"
}' --duration 4h --interval 60s --keyring-backend test -y --from b --fees 600utrst --connection-id connection-0

```


### Design

If host chain is offline then AutoTXs will resume afterwards. 

If packet times out when there are no relayers, channel will be closed and `AutoTX` will stop execution.

We check if the message is valid using `ValidateBasic()` in `SubmitAutoTx`. This is because if message returns an error at basic validation on the host chain, the packet won't be included. If a next packet arrives before timeout, the packet sequence will mismatch. The ICA channel will timeout. A new Interchain account address must be generated to resume using ICA. By checking for validity we prevent these issues from arising.

If a message retuns an error after basic validation, there won't be an acknoledgement. `Executed` will be `false` for that given execution entry. The ICA remains active and AutoTX resumes.