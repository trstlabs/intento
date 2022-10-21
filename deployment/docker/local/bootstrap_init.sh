#!/bin/bash
set -euo pipefail

file=~/.trstd/config/genesis.json
if [ ! -e "$file" ]
then
  # init the node
  rm -rf ~/.trstd/*
  rm -rf ~/.trstcli/*
  rm -rf ~/.sgx_secrets/*
  trstcli config chain-id enigma-pub-testnet-3
  trstcli config output json
#  trstcli config indent true
#  trstcli config trust-node true
  trstcli config keyring-backend test

  trstd init banana --chain-id enigma-pub-testnet-3

  cp ~/node_key.json ~/.trstd/config/node_key.json

  perl -i -pe 's/"stake"/"utrst"/g' ~/.trstd/config/genesis.json

  perl -i -pe 's/"1814400s"/"80s"/g' ~/.trstd/config/genesis.json

  trstcli keys add a
  trstcli keys add b
  trstcli keys add c
  trstcli keys add d

  trstd add-genesis-account "$(trstcli keys show -a a)" 1000000000000000000utrst
#  trstd add-genesis-account "$(trstcli keys show -a b)" 1000000000000000000utrst
#  trstd add-genesis-account "$(trstcli keys show -a c)" 1000000000000000000utrst
#  trstd add-genesis-account "$(trstcli keys show -a d)" 1000000000000000000utrst


  trstd gentx a 1000000utrst --keyring-backend test --chain-id enigma-pub-testnet-3
  # These fail for some reason:
  # trstd gentx --name b --keyring-backend test --amount 1000000utrst
  # trstd gentx --name c --keyring-backend test --amount 1000000utrst
  # trstd gentx --name d --keyring-backend test --amount 1000000utrst

  trstd collect-gentxs
  trstd validate-genesis

  trstd init-bootstrap
  trstd validate-genesis
fi

# sleep infinity
source /opt/sgxsdk/environment && RUST_BACKTRACE=1 trstd start --rpc.laddr tcp://0.0.0.0:26657 --bootstrap
