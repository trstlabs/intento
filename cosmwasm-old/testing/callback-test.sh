#!/bin/bash

set -euvx

function wait_for_tx() {
    until (./tppd q tx "$1"); do
        echo "$2"
        sleep 1
    done
}

# init the node
rm -rf ./.sgx_secrets
mkdir -p ./.sgx_secrets

rm -rf ~/.secret*

./tppd config chain-id enigma-testnet
./tppd config output json
./tppd config indent true
./tppd config trust-node true
./tppd config keyring-backend test

./tppd init banana --chain-id enigma-testnet
perl -i -pe 's/"stake"/"tpp"/g' ~/.tppd/config/genesis.json
echo "cost member exercise evoke isolate gift cattle move bundle assume spell face balance lesson resemble orange bench surge now unhappy potato dress number acid" |
    ./tppd keys add a --recover
./tppd add-genesis-account "$(./tppd keys show -a a)" 1000000000000tpp
./tppd gentx --name a --keyring-backend test --amount 1000000tpp
./tppd collect-gentxs
./tppd validate-genesis

./tppd init-bootstrap ./node-master-cert.der ./io-master-cert.der

./tppd validate-genesis

RUST_BACKTRACE=1 ./tppd start --bootstrap &

export tppd_PID=$(echo $!)

until (./tppd status 2>&1 | jq -e '(.sync_info.latest_block_height | tonumber) > 0' &>/dev/null); do
    echo "Waiting for chain to start..."
    sleep 1
done

./tppd rest-server --chain-id enigma-testnet --laddr tcp://0.0.0.0:1337 &
export LCD_PID=$(echo $!)
function cleanup() {
    kill -KILL "$tppd_PID" "$LCD_PID"
}
trap cleanup EXIT ERR

export STORE_TX_HASH=$(
    yes |
        ./tppd tx compute store ./x/compute/internal/keeper/testdata/test-contract/contract.wasm --from a --gas 10000000 |
        jq -r .txhash
)

wait_for_tx "$STORE_TX_HASH" "Waiting for store to finish on-chain..."

# test storing of wasm code (this doesn't touch sgx yet)
./tppd q tx "$STORE_TX_HASH" |
    jq -e '.logs[].events[].attributes[] | select(.key == "code_id" and .value == "1")'

# init the contract (ocall_init + write_db + canonicalize_address)
export INIT_TX_HASH=$(
    yes |
        ./tppd tx compute instantiate 1 '{"nop":{}}' --label baaaaaaa --from a |
        jq -r .txhash
)

wait_for_tx "$INIT_TX_HASH" "Waiting for instantiate to finish on-chain..."

./tppd q compute tx "$INIT_TX_HASH"

export CONTRACT_ADDRESS=$(
    ./tppd q tx "$INIT_TX_HASH" |
        jq -er '.logs[].events[].attributes[] | select(.key == "contract_address") | .value' | head -1
)

# exec (generate callbacks)
export EXEC_TX_HASH=$(
    yes |
        ./tppd tx compute execute --from a $CONTRACT_ADDRESS "{\"a\":{\"contract_addr\":\"$CONTRACT_ADDRESS\",\"x\":2,\"y\":3}}" |
        jq -r .txhash
)

wait_for_tx "$EXEC_TX_HASH" "Waiting for exec to finish on-chain..."

./tppd q compute tx "$EXEC_TX_HASH"

# exec (generate error inside WASM)
export EXEC_ERR_TX_HASH=$(
    yes |
        ./tppd tx compute execute --from a $CONTRACT_ADDRESS "{\"contract_error\":{\"error_type\":\"generic_err\"}}" |
        jq -r .txhash
)

wait_for_tx "$EXEC_ERR_TX_HASH" "Waiting for exec to finish on-chain..."

./tppd q compute tx "$EXEC_ERR_TX_HASH"

# exec (generate error inside WASM)
export EXEC_ERR_TX_HASH=$(
    yes |
        ./tppd tx compute execute --from a $CONTRACT_ADDRESS '{"allocate_on_heap":{"bytes":1073741824}}' |
        jq -r .txhash
)

wait_for_tx "$EXEC_ERR_TX_HASH" "Waiting for exec to finish on-chain..."

./tppd q compute tx "$EXEC_ERR_TX_HASH"
# test output data decryption
yes |
    ./tppd tx compute execute --from a "$CONTRACT_ADDRESS" '{"unicode_data":{}}' -b block |
    jq -r .txhash |
    xargs ./tppd q compute tx

# sleep infinity

(
    cd ./cosmwasm-js
    yarn
    cd ./packages/sdk
    yarn build
)

node ./cosmwasm/testing/callback-test.js
