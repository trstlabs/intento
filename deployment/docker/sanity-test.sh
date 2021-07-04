#!/bin/bash

set -euvx

OUTPUT_FORMAT='json'

function wait_for_tx () {
    until (tppd q tx "$1" --output $OUTPUT_FORMAT)
    do
        echo "$2"
        sleep 1
    done
}

# init the node
rm -rf ./.sgx_secrets
mkdir -p ./.sgx_secrets

rm -rf ~/.secret*
#mkdir -p ~/.tppd/config
#echo 'chain-id="secret-sanity"
#output="json"
#indent=true
#trust-node=true
#keyring-backend="test"' > ~/.tppd/config/config.toml

export SECRET_NETWORK_CHAIN_ID=tppdev-1
export SECRET_NETWORK_KEYRING_BACKEND=test

tppd init banana --chain-id tppdev-1
perl -i -pe 's/"stake"/"uscrt"/g' ~/.tppd/config/genesis.json
echo "cost member exercise evoke isolate gift cattle move bundle assume spell face balance lesson resemble orange bench surge now unhappy potato dress number acid" |
    tppd keys add a --recover
tppd add-genesis-account "$(tppd keys show -a a)" 1000000000000uscrt
tppd gentx a 1000000uscrt
tppd collect-gentxs
tppd validate-genesis

tppd init-bootstrap # ./node-master-cert.der ./io-master-cert.der

tppd validate-genesis

RUST_BACKTRACE=1 tppd start --bootstrap &

export tppd_PID=$(echo $!)

until (tppd status 2>&1 | jq -e '(.SyncInfo.latest_block_height | tonumber) > 0' &>/dev/null); do
    echo "Waiting for chain to start..."
    sleep 1
done

tppd rest-server --laddr tcp://0.0.0.0:1337 &
export LCD_PID=$(echo $!)
function cleanup() {
    kill -KILL "$tppd_PID" "$LCD_PID"
}
trap cleanup EXIT ERR

# store wasm code on-chain so we could later instansiate it
export STORE_TX_HASH=$(
    yes |
    tppd tx compute store erc20.wasm --from a --gas 1200000 --gas-prices 0.25uscrt --output $OUTPUT_FORMAT |
    jq -r .txhash
)

wait_for_tx "$STORE_TX_HASH" "Waiting for store to finish on-chain..."

# test storing of wasm code (this doesn't touch sgx yet)
tppd q tx "$STORE_TX_HASH" --output $OUTPUT_FORMAT |
    jq -e '.logs[].events[].attributes[] | select(.key == "code_id" and .value == "1")'

# init the contract (ocall_init + write_db + canonicalize_address)
# a is a tendermint address (will be used in transfer: https://github.com/CosmWasm/cosmwasm-examples/blob/f5ea00a85247abae8f8cbcba301f94ef21c66087/erc20/src/contract.rs#L110)
# secret1f395p0gg67mmfd5zcqvpnp9cxnu0hg6rjep44t is just a random address
# balances are set to 108 & 53 at init
export INIT_TX_HASH=$(
    yes |
    tppd tx compute instantiate 1 '{
        "decimals":10,
        "initial_balances":[
            {"address":"'"$(tppd keys show a -a)"'","amount":"108"},
            {"address":"secret1f395p0gg67mmfd5zcqvpnp9cxnu0hg6rjep44t","amount":"53"}
        ],
        "name":"ReuvenPersonalRustCoin",
        "symbol":"RPRC"
    }' --label RPRCCoin --gas 1000000 --gas-prices 0.25uscrt --from a --output $OUTPUT_FORMAT |
    jq -r .txhash
)

wait_for_tx "$INIT_TX_HASH" "Waiting for instantiate to finish on-chain..."

tppd q compute tx "$INIT_TX_HASH"

export CONTRACT_ADDRESS=$(
    tppd q tx "$INIT_TX_HASH" |
    jq -er '.logs[].events[].attributes[] | select(.key == "contract_address") | .value'
)

# test balances after init (ocall_query + read_db + canonicalize_address)
tppd q compute query "$CONTRACT_ADDRESS" "{\"balance\":{\"address\":\"$(tppd keys show a -a)\"}}" |
    jq -e '.balance == "108"'
tppd q compute query "$CONTRACT_ADDRESS" "{\"balance\":{\"address\":\"secret1f395p0gg67mmfd5zcqvpnp9cxnu0hg6rjep44t\"}}" |
    jq -e '.balance == "53"'

# transfer 10 balance (ocall_handle + read_db + write_db + humanize_address + canonicalize_address)
yes |
    tppd tx compute execute --from a "$CONTRACT_ADDRESS" '{"transfer":{"amount":"10","recipient":"secret1f395p0gg67mmfd5zcqvpnp9cxnu0hg6rjep44t"}}' -b block |
    jq -r .txhash |
    xargs tppd q compute tx

# test balances after transfer (ocall_query + read_db)
tppd q compute query "$CONTRACT_ADDRESS" "{\"balance\":{\"address\":\"$(tppd keys show a -a)\"}}" |
    jq -e '.balance == "98"'
tppd q compute query "$CONTRACT_ADDRESS" "{\"balance\":{\"address\":\"secret1f395p0gg67mmfd5zcqvpnp9cxnu0hg6rjep44t\"}}" |
    jq -e '.balance == "63"'

(tppd q compute query "$CONTRACT_ADDRESS" "{\"balance\":{\"address\":\"secret1zzzzzzzzzzzzzzzzzz\"}}" || true) 2>&1 | grep -c 'canonicalize_address errored: invalid checksum'

# sleep infinity

(
    cd ./cosmwasm-js
    yarn
    cd ./packages/sdk
    yarn build
)

node ./cosmwasm/testing/cosmwasm-js-test.js

echo "All is done. Yay!"
