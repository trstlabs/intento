#!/bin/bash

file=~/.tppd/config/genesis.json
if [ ! -e "$file" ]
then
  # init the node
  rm -rf ~/.tppd/*
  rm -rf ~/.sgx_secrets/*
#  tppd config chain-id enigma-pub-testnet-3
#  tppd config output json
#  tppd config indent true
#  tppd config trust-node true
#  tppd config keyring-backend test
  export SECRET_NETWORK_CHAIN_ID=tppdev-1
  export SECRET_NETWORK_KEYRING_BACKEND=test
  tppd init banana --chain-id tppdev-1

  cp ~/node_key.json ~/.tppd/config/node_key.json

  perl -i -pe 's/"stake"/ "tpp"/g' ~/.tppd/config/genesis.json
  tppd keys add a
  tppd keys add b
  tppd keys add c
  tppd keys add d

  tppd add-genesis-account "$(tppd keys show -a a)" 1000000000000000000tpp
#  tppd add-genesis-account "$(tppd keys show -a b)" 1000000000000000000tpp
#  tppd add-genesis-account "$(tppd keys show -a c)" 1000000000000000000tpp
#  tppd add-genesis-account "$(tppd keys show -a d)" 1000000000000000000tpp


  tppd gentx a 1000000tpp
#  tppd gentx b 1000000tpp --keyring-backend test
#  tppd gentx c 1000000tpp --keyring-backend test
#  tppd gentx d 1000000tpp --keyring-backend test

  tppd collect-gentxs
  tppd validate-genesis

  tppd init-bootstrap
  tppd validate-genesis
fi

# sleep infinity
source /opt/sgxsdk/environment && RUST_BACKTRACE=1 tppd start --rpc.laddr tcp://0.0.0.0:26657 --bootstrap