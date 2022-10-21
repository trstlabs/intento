#!/bin/bash

set -euvx

function wait_for_tx () {
    until (trstd q tx "$1" --output json)
    do
        echo "$2"
        sleep 1
    done
}

# # init the node
# rm -rf ./.sgx_secrets ~/.sgx_secrets *.der ~/*.der
# mkdir -p ./.sgx_secrets ~/.sgx_secrets

# rm -rf ~/.trstd

# #export SECRET_NETWORK_CHAIN_ID=trst_chain_1
# #export SECRET_NETWORK_KEYRING_BACKEND=test

# trstd init banana --chain-id trst_chain_1
# perl -i -pe 's/"stake"/"utrst"/g' ~/.trstd/config/genesis.json
# echo "cost member exercise evoke isolate gift cattle move bundle assume spell face balance lesson resemble orange bench surge now unhappy potato dress number acid" |
#     trstd keys add a --recover --keyring-backend test
# trstd add-genesis-account "$(trstd keys show -a --keyring-backend test a)" 1000000000000utrst
# trstd gentx a 1000000utrst --chain-id trst_chain_1 --keyring-backend test
# trstd collect-gentxs
# trstd validate-genesis

# trstd init-bootstrap node-master-cert.der io-master-cert.der
# trstd validate-genesis

# RUST_BACKTRACE=1 trstd start --bootstrap &


# export trstd_PID=$(echo $!)


# until (trstd status 2>&1 | jq -e '(.SyncInfo.latest_block_height | tonumber) > 0' &>/dev/null); do
#     echo "Waiting for chain to start..."
#     sleep 1
# done

# # trstd rest-server --laddr tcp://0.0.0.0:1337 &
# export LCD_PID=$(echo $!)
# function cleanup() {
#     kill -KILL "$trstd_PID" "$LCD_PID"
# }
# trap cleanup EXIT ERR

# store wasm code on-chain so we could later instansiate it
export STORE_TX_HASH=$(
    trstd tx compute store erc20.wasm --from a --gas 10000000 --gas-prices 0.25utrst --output json -y |
        jq -r .txhash
)

wait_for_tx "$STORE_TX_HASH" "Waiting for store to finish on-chain..."

# test storing of wasm code (this doesn't touch sgx yet)
    trstd q tx "$STORE_TX_HASH" --output json |
        jq -e '.logs[].events[].attributes[] | select(.key == "code_id" and .value == "1")'

# init the contract (ocall_init + write_db + canonicalize_address)
# a is a tendermint address (will be used in transfer: https://github.com/CosmWasm/cosmwasm-examples/blob/f5ea00a85247abae8f8cbcba301f94ef21c66087/erc20/src/contract.rs#L110)
# secret1f395p0gg67mmfd5zcqvpnp9cxnu0hg6rjep44t is just a random address
# balances are set to 108 & 53 at init
export INIT_TX_HASH=$(
    trstd tx compute instantiate 1 "{\"decimals\":10,\"initial_balances\":[{\"address\":\"$(trstd keys show a -a)\",\"amount\":\"108\"},{\"address\":\"secret1f395p0gg67mmfd5zcqvpnp9cxnu0hg6rjep44t\",\"amount\":\"53\"}],\"name\":\"ReuvenPersonalRustCoin\",\"symbol\":\"RPRC\"}" --label RPRCCoin --from a --output json -y --gas-prices 0.25utrst |
        jq -r .txhash
)

wait_for_tx "$INIT_TX_HASH" "Waiting for instantiate to finish on-chain..."

trstd q compute tx "$INIT_TX_HASH" --output json

export CONTRACT_ADDRESS=$(
    trstd q tx "$INIT_TX_HASH" --output json |
        jq -er '.logs[].events[].attributes[] | select(.key == "contract_address") | .value' |
        head -1
)

# test balances after init (ocall_query + read_db + canonicalize_address)
trstd q compute query "$CONTRACT_ADDRESS" "{\"balance\":{\"address\":\"$(trstd keys show a -a)\"}}" |
    jq -e '.balance == "108"'
trstd q compute query "$CONTRACT_ADDRESS" "{\"balance\":{\"address\":\"secret1f395p0gg67mmfd5zcqvpnp9cxnu0hg6rjep44t\"}}" |
    jq -e '.balance == "53"'

# transfer 10 balance (ocall_handle + read_db + write_db + humanize_address + canonicalize_address)
trstd tx compute execute "$CONTRACT_ADDRESS" '{"transfer":{"amount":"10","recipient":"secret1f395p0gg67mmfd5zcqvpnp9cxnu0hg6rjep44t"}}' --gas-prices 0.25utrst --from a -b block -y --output json |
    jq -r .txhash |
    xargs trstd q compute tx

# test balances after transfer (ocall_query + read_db)
trstd q compute query "$CONTRACT_ADDRESS" "{\"balance\":{\"address\":\"$(trstd keys show a -a)\"}}" |
    jq -e '.balance == "98"'
trstd q compute query "$CONTRACT_ADDRESS" "{\"balance\":{\"address\":\"secret1f395p0gg67mmfd5zcqvpnp9cxnu0hg6rjep44t\"}}" |
    jq -e '.balance == "63"'

(trstd q compute query "$CONTRACT_ADDRESS" "{\"balance\":{\"address\":\"secret1zzzzzzzzzzzzzzzzzz\"}}" || true) 2>&1 | grep -c 'canonicalize_address errored: invalid checksum'

# sleep infinity

(
    cd ./cosmwasm-js
    yarn
    cd ./packages/sdk
    yarn build
)

node ./cosmwasm/testing/cosmwasm-js-test.js

echo "All is done. Yay!"
