#!/bin/bash

# init the node
rm -rf ~/.trstd*
trstd config chain-id trst_chain_1
trstd config output json
trstd config indent true
trstd config trust-node true
trstd config keyring-backend test

trstd init banana --chain-id trst_chain_1

cp ~/node_key.json ~/.trstd/config/node_key.json

perl -i -pe 's/"stake"/"utrst"/g' ~/.trstd/config/genesis.json
trstd keys add a

trstd add-genesis-account "$(trstd keys show -a a)" 1000000000000utrst
trstd gentx --name a --keyring-backend test --amount 1000000utrst
trstd collect-gentxs
trstd validate-genesis

trstd init-bootstrap
trstd validate-genesis

sed -i 's/persistent_peers = ""/persistent_peers = "'"$PERSISTENT_PEERS"'"/g' ~/.trstd/config/config.toml

sed -i '104s/enable = false/enable = true/g' ~/.trstd/config/app.toml

sed -i 's/enabled-unsafe-cors = false/enabled-unsafe-cors = true/g' ~/.trstd/config/app.toml


source /opt/sgxsdk/environment && RUST_BACKTRACE=1 trstd start --rpc.laddr tcp://0.0.0.0:26657 --bootstrap --log_level $LOG_LEVEL