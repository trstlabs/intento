#!/bin/bash
set -euo pipefail

file=~/.trstd/config/genesis.json
if [ ! -e "$file" ]
then
    # init the node
    rm -rf ~/.trstd/*
    rm -rf ~/.sgx_secrets/*
    trstd config chain-id trst-pub-testnet-1
    trstd config output json
    trstd config keyring-backend test
    
    trstd init FRST --chain-id trst-pub-testnet-1
    trstd prepare-genesis testnet trst-pub-testnet-1
    cp ~/node_key.json ~/.trstd/config/node_key.json
    
    perl -i -pe 's/"stake"/"utrst"/g' ~/.trstd/config/genesis.json
    
    perl -i -pe 's/"1814400s"/"80s"/g' ~/.trstd/config/genesis.json
    
    trstd keys add a
    trstd keys add b
    trstd keys add c
    trstd keys add d
    
    trstd add-genesis-account "$(trstd keys show -a a)" 100000000000000utrst
    trstd gentx a 1000000utrst --keyring-backend test --chain-id trst-pub-testnet-1
    
    
    trstd collect-gentxs
    trstd validate-genesis
    
    trstd init-bootstrap
    trstd validate-genesis
fi

# sleep infinity
source /opt/sgxsdk/environment && RUST_BACKTRACE=1 trstd start --rpc.laddr tcp://0.0.0.0:26657 --bootstrap > init.log
