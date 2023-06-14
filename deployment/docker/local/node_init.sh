#!/usr/bin/env bash

set -euvo pipefail

# init the node
# rm -rf ~/.secret*
#trstd config chain-id trst-testnet
#trstd config output json
#trstd config indent true
#trstd config trust-node true
#trstd config keyring-backend test
# rm -rf ~/.trstd

mkdir -p /root/.trstd/.node
trstd config keyring-backend test
trstd config node http://bootstrap:26657
trstd config chain-id trst-pub-testnet-1

mkdir -p /root/.trstd/.node

trstd init "$(hostname)" --chain-id trst-pub-testnet-1 || true

PERSISTENT_PEERS=115aa0a629f5d70dd1d464bc7e42799e00f4edae@bootstrap:26656

sed -i 's/persistent_peers = ""/persistent_peers = "'$PERSISTENT_PEERS'"/g' ~/.trstd/config/config.toml
sed -i 's/trust_period = "168h0m0s"/trust_period = "168h"/g' ~/.trstd/config/config.toml
echo "Set persistent_peers: $PERSISTENT_PEERS"

echo "Waiting for bootstrap to start..."
sleep 20

trstd q block 1

cp /tmp/.trstd/keyring-test /root/.trstd/ -r

# MASTER_KEY="$(trstd q register node-enclave-params 2> /dev/null | cut -c 3- )"

#echo "Master key: $MASTER_KEY"

trstd init-enclave --reset

PUBLIC_KEY=$(trstd parse /opt/trustlesshub/.sgx_secrets/attestation_cert.der | cut -c 3- )

echo "Public key: $PUBLIC_KEY"

trstd parse /opt/trustlesshub/.sgx_secrets/attestation_cert.der
cat /opt/trustlesshub/.sgx_secrets/attestation_cert.der
tx_hash="$(trstd tx register auth /opt/trustlesshub/.sgx_secrets/attestation_cert.der -y --from a --gas-prices 0.25utrst | jq -r '.txhash')"

#trstd q tx "$tx_hash"
sleep 15
trstd q tx "$tx_hash"

SEED="$(trstd q register seed "$PUBLIC_KEY" | cut -c 3-)"
echo "SEED: $SEED"
#exit

trstd q register node-enclave-params

trstd configure-credentials node-master-cert.der "$SEED"

cp /tmp/.trstd/config/genesis.json /root/.trstd/config/genesis.json

trstd validate-genesis

RUST_BACKTRACE=1 trstd start --rpc.laddr tcp://0.0.0.0:26657

# ./wasmi-sgx-test.sh
