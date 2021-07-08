#!/bin/bash

# init the node
rm -rf ~/.tpp*
tppd config chain-id enigma-testnet
tppd config output json
tppd config indent true
tppd config trust-node true
tppd config keyring-backend test

tppd init banana --chain-id enigma-testnet

cp ~/node_key.json ~/.tppd/config/node_key.json

perl -i -pe 's/"stake"/"uscrt"/g' ~/.tppd/config/genesis.json
tppd keys add a

tppd add-genesis-account "$(tppd keys show -a a)" 1000000000000uscrt
tppd gentx --name a --keyring-backend test --amount 1000000uscrt
tppd collect-gentxs
tppd validate-genesis

tppd init-bootstrap
tppd validate-genesis

sed -i 's/persistent_peers = ""/persistent_peers = "'"$PERSISTENT_PEERS"'"/g' ~/.tpp/config/config.toml



sed -i '104s/enable = false/enable = true/g' ~/.tpp/config/app.toml

sed -i 's/enabled-unsafe-cors = false/enabled-unsafe-cors = true/g' ~/.tpp/config/app.toml


source /opt/sgxsdk/environment && RUST_BACKTRACE=1 tppd start --rpc.laddr tcp://0.0.0.0:26657 --bootstrap