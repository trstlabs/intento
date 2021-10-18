#!/bin/bash

set -euvx

function wait_for_tx() {
    until (./trstd q tx "$1"); do
        echo "$2"
        sleep 1
    done
}

# init the node
rm -rf ./.sgx_secrets
mkdir -p ./.sgx_secrets

rm -rf ~/.secret*

./trstd config chain-id enigma-testnet
./trstd config output json
./trstd config indent true
./trstd config trust-node true
./trstd config keyring-backend test

./trstd init banana --chain-id enigma-testnet
perl -i -pe 's/"stake"/"utrst"/g' ~/.trstd/config/genesis.json
echo "cost member exercise evoke isolate gift cattle move bundle assume spell face balance lesson resemble orange bench surge now unhappy potato dress number acid" |
    ./trstd keys add a --recover
./trstd add-genesis-account "$(./trstd keys show -a a)" 1000000000000trst
./trstd gentx --name a --keyring-backend test --amount 1000000trst
./trstd collect-gentxs
./trstd validate-genesis

./trstd init-bootstrap ./node-master-cert.der ./io-master-cert.der

./trstd validate-genesis

RUST_BACKTRACE=1 ./trstd start --bootstrap &

export trstd_PID=$(echo $!)

until (./trstd status 2>&1 | jq -e '(.sync_info.latest_block_height | tonumber) > 0' &>/dev/null); do
    echo "Waiting for chain to start..."
    sleep 1
done

./trstd rest-server --chain-id enigma-testnet --laddr tcp://0.0.0.0:1337 &
export LCD_PID=$(echo $!)
function cleanup() {
    kill -KILL "$trstd_PID" "$LCD_PID"
}
trap cleanup EXIT ERR

export STORE_TX_HASH=$(
    yes |
        ./trstd tx compute store ./x/compute/internal/keeper/testdata/test-contract/contract.wasm --from a --gas 10000000 |
        jq -r .txhash
)

wait_for_tx "$STORE_TX_HASH" "Waiting for store to finish on-chain..."

# test storing of wasm code (this doesn't touch sgx yet)
./trstd q tx "$STORE_TX_HASH" |
    jq -e '.logs[].events[].attributes[] | select(.key == "code_id" and .value == "1")'

# init the contract (ocall_init + write_db + canonicalize_address)
export INIT_TX_HASH=$(
    yes |
        ./trstd tx compute instantiate 1 '{"nop":{}}' --label baaaaaaa --from a |
        jq -r .txhash
)

wait_for_tx "$INIT_TX_HASH" "Waiting for instantiate to finish on-chain..."

./trstd q compute tx "$INIT_TX_HASH"

export CONTRACT_ADDRESS=$(
    ./trstd q tx "$INIT_TX_HASH" |
        jq -er '.logs[].events[].attributes[] | select(.key == "contract_address") | .value' | head -1
)

# exec (generate callbacks)
export EXEC_TX_HASH=$(
    yes |
        ./trstd tx compute execute --from a $CONTRACT_ADDRESS "{\"a\":{\"contract_addr\":\"$CONTRACT_ADDRESS\",\"x\":2,\"y\":3}}" |
        jq -r .txhash
)

wait_for_tx "$EXEC_TX_HASH" "Waiting for exec to finish on-chain..."

./trstd q compute tx "$EXEC_TX_HASH"

# exec (generate error inside WASM)
export EXEC_ERR_TX_HASH=$(
    yes |
        ./trstd tx compute execute --from a $CONTRACT_ADDRESS "{\"contract_error\":{\"error_type\":\"generic_err\"}}" |
        jq -r .txhash
)

wait_for_tx "$EXEC_ERR_TX_HASH" "Waiting for exec to finish on-chain..."

./trstd q compute tx "$EXEC_ERR_TX_HASH"

# exec (generate error inside WASM)
export EXEC_ERR_TX_HASH=$(
    yes |
        ./trstd tx compute execute --from a $CONTRACT_ADDRESS '{"allocate_on_heap":{"bytes":1073741824}}' |
        jq -r .txhash
)

wait_for_tx "$EXEC_ERR_TX_HASH" "Waiting for exec to finish on-chain..."

./trstd q compute tx "$EXEC_ERR_TX_HASH"
# test output data decryption
yes |
    ./trstd tx compute execute --from a "$CONTRACT_ADDRESS" '{"unicode_data":{}}' -b block |
    jq -r .txhash |
    xargs ./trstd q compute tx

# sleep infinity

(
    cd ./cosmwasm-js
    yarn
    cd ./packages/sdk
    yarn build
)

node ./cosmwasm/testing/callback-test.js
