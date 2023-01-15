#!/bin/bash

file=~/.trstd/config/genesis.json
if [ ! -e "$file" ]
then
    # init the node
    rm -rf ~/.trstd/*
    rm -rf /opt/trustlesshub/.sgx_secrets/*
    
    chain_id=${CHAINID:-trstdev-1}
    
    mkdir -p ./.sgx_secrets
    trstd config chain-id "$chain_id"
    trstd config output json
    trstd config keyring-backend test
    
    # export TRUSTLESS_HUB_CHAIN_ID=trst_chain_1
    # export TRUSTLESS_HUB_KEYRING_BACKEND=test
    trstd init banana --chain-id "$chain_id"
    
    trstd prepare-genesis testnet "$chain_id"
    
    cp ~/node_key.json ~/.trstd/config/node_key.json
    perl -i -pe 's/"stake"/"utrst"/g' ~/.trstd/config/genesis.json
    perl -i -pe 's/"172800s"/"90s"/g' ~/.trstd/config/genesis.json # voting period 2 days -> 90 seconds
    perl -i -pe 's/"1814400s"/"80s"/g' ~/.trstd/config/genesis.json # unbonding period 21 days -> 80 seconds
    
    # perl -i -pe 's/enable-unsafe-cors = false/enable-unsafe-cors = true/g' ~/.trstd/config/app.toml # enable cors
    
    a_mnemonic="grant rice replace explain federal release fix clever romance raise often wild taxi quarter soccer fiber love must tape steak together observe swap guitar"
    b_mnemonic="jelly shadow frog dirt dragon use armed praise universe win jungle close inmate rain oil canvas beauty pioneer chef soccer icon dizzy thunder meadow"
    c_mnemonic="chair love bleak wonder skirt permit say assist aunt credit roast size obtain minute throw sand usual age smart exact enough room shadow charge"
    d_mnemonic="word twist toast cloth movie predict advance crumble escape whale sail such angry muffin balcony keen move employ cook valve hurt glimpse breeze brick"
    
    echo $a_mnemonic | trstd keys add a --recover
    echo $b_mnemonic | trstd keys add b --recover
    echo $c_mnemonic | trstd keys add c --recover
    echo $d_mnemonic | trstd keys add d --recover
    
    trstd add-genesis-account "$(trstd keys show -a a)" 1000000000000000000utrst
    trstd add-genesis-account "$(trstd keys show -a b)" 1000000000000000000utrst
    trstd add-genesis-account "$(trstd keys show -a c)" 1000000000000000000utrst
    trstd add-genesis-account "$(trstd keys show -a d)" 1000000000000000000utrst
    
    
    trstd gentx a 1000000utrst --chain-id "$chain_id"
    
    trstd collect-gentxs
    trstd validate-genesis
    
    #  trstd init-enclave
    trstd init-bootstrap
    #  cp new_node_seed_exchange_keypair.sealed .sgx_secrets
    trstd validate-genesis
fi

# Setup CORS for LCD & gRPC-web
perl -i -pe 's;address = "tcp://0.0.0.0:1317";address = "tcp://0.0.0.0:1316";' .trstd/config/app.toml
# perl -i -pe 's/enable-unsafe-cors = false/enable-unsafe-cors = true/' .trstd/config/app.toml
lcp --proxyUrl http://localhost:1316 --port 1317 --proxyPartial '' &

# Setup faucet
setsid node faucet_server.js &

# Setup trstcli
cp $(which trstd) $(dirname $(which trstd))/trstcli

source /opt/sgxsdk/environment && RUST_BACKTRACE=1 LOG_LEVEL=$LOG_LEVEL trstd start --rpc.laddr tcp://0.0.0.0:26657 --bootstrap  --log_level $LOG_LEVEL

