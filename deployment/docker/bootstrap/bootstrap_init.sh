#!/bin/bash

file=~/.trstd/config/genesis.json
if [ ! -e "$file" ]; then
  # init the node
  rm -rf ~/.trstd/*
  rm -rf /opt/trustlesshub/.sgx_secrets/*

  chain_id=${CHAINID:-supernova-1}

  mkdir -p ./.sgx_secrets
  trstd config chain-id "$chain_id"
  trstd config keyring-backend test

  trstd init banana --chain-id "$chain_id"


  cp ~/node_key.json ~/.trstd/config/node_key.json
  perl -i -pe 's/"stake"/"utrst"/g' ~/.trstd/config/genesis.json
  perl -i -pe 's/"172800000000000"/"90000000000"/g' ~/.trstd/config/genesis.json # voting period 2 days -> 90 seconds

  trstd keys add a
  trstd keys add b
  trstd keys add c
  trstd keys add d

  trstd add-genesis-account "$(trstd keys show -a a)" 1000000000000000000utrst
#  trstd add-genesis-account "$(trstd keys show -a b)" 1000000000000000000utrst



  trstd gentx a 1000000utrst --chain-id "$chain_id"
#  trstd gentx b 1000000utrst --keyring-backend test


  trstd collect-gentxs
  trstd validate-genesis

#  trstd init-enclave
  trstd init-bootstrap
#  cp new_node_seed_exchange_keypair.sealed .sgx_secrets
  trstd validate-genesis

  perl -i -pe 's/max_subscription_clients.+/max_subscription_clients = 100/' ~/.trstd/config/config.toml
  perl -i -pe 's/max_subscriptions_per_client.+/max_subscriptions_per_client = 50/' ~/.trstd/config/config.toml
fi

lcp --proxyUrl http://localhost:1317 --port 1337 --proxyPartial '' &

# sleep infinity
source /opt/sgxsdk/environment && RUST_BACKTRACE=1 trstd start --rpc.laddr tcp://0.0.0.0:26657 --bootstrap