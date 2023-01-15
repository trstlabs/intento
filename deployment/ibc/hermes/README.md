# Trustless Hub IBC setup

Two local Trustless Hub chains can communicate with each other via a Hermes relayer

## Build

```bash
docker build -f hermes.Dockerfile . --tag hermes:test
```

### Run

```bash
docker compose up
```

### Verify IBC transfers

Assuming you have a key 'a' which is not the relayer's key,
from localhost:

```bash
a_mnemonic="grant rice replace explain federal release fix clever romance raise often wild taxi quarter soccer fiber love must tape steak together observe swap guitar"

echo $a_mnemonic | trstcli keys add a --recover

trstcli add-genesis-account "$(trstcli keys show -a a)" 1000000000000000000utrst

# be on the source network (trstdev-1)
trstcli config node http://localhost:26657

# check the initial balance of a
trstcli q bank balances trust1q6k0w4cejawpkzxgqhvs4m2v6uvdzm6j2pk2jx

# transfer to the destination network
trstcli tx ibc-transfer transfer transfer channel-0 trust1he7t2wxzpmfuxfrw7qjg52vu4qljq3l56w5qqw 2utrst --from a

# check a's balance after transfer
trstcli q bank balances trust1q6k0w4cejawpkzxgqhvs4m2v6uvdzm6j2pk2jx

# switch to the destination network (trstdev-2)
trstcli config node http://localhost:36657

# check that you have an ibc-denom
trstcli q bank balances trust1ykql5ktedxkpjszj5trzu8f5dxajvgv95nuwjx # should have 1 ibc denom
```

### Interchain accounts

Message flow for interchain acccounts

```bash
# register account for address
trstd tx icamsgauth register --connection-id connection-0 --counterparty-connection-id connection-0  --keyring-backend test -y --from b --fees 600utrst

# query the ICA address
trstd q icamsgauth interchainaccounts trust1ykql5ktedxkpjszj5trzu8f5dxajvgv95nuwjx connection-0 connection-0

# query the channel, make sure channel is STATE_OPEN
trstd query ibc channel channels

# send balance to ICA on host chain (replace node and to_address here)
trstd  tx bank send trust1ykql5ktedxkpjszj5trzu8f5dxajvgv95nuwjx trust1tm8vt97s094s6egax2ms44xths0r40ftzufp93esd7sw0s3mn6dqahqh49 10000utrst --node tcp://localhost:36657 --keyring-backend test -y --from b --fees 600utrst --chain-id trstdev-2

# replace msg delegator to ICA address and submit tx
trstd tx icamsgauth submit-tx  ./msg.json --counterparty-connection-id connection-0  --keyring-backend test -y --from b --fees 600utrst --connection-id connection-0

# check balance
trstd q bank balances trust1tm8vt97s094s6egax2ms44xths0r40ftzufp93esd7sw0s3mn6dqahqh49 --node tcp://localhost:36657

# check staking delegations
trstd  q staking delegations-to trustvaloper1q6k0w4cejawpkzxgqhvs4m2v6uvdzm6jhmz5jy --node tcp://localhost:36657
```