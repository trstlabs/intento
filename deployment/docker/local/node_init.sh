#!/usr/bin/env bash

set -euvo pipefail

# init the node
# rm -rf ~/.secret*
#trstcli config chain-id enigma-testnet
#trstcli config output json
#trstcli config indent true
#trstcli config trust-node true
#trstcli config keyring-backend test
# rm -rf ~/.trstd

mkdir -p /root/.trstd/.node
trstd config keyring-backend test
trstd config node http://bootstrap:26657
trstd config chain-id enigma-pub-testnet-3

mkdir -p /root/.trstd/.node

trstd init "$(hostname)" --chain-id enigma-pub-testnet-3 || true

PERSISTENT_PEERS=115aa0a629f5d70dd1d464bc7e42799e00f4edae@bootstrap:26656

sed -i 's/persistent_peers = ""/persistent_peers = "'$PERSISTENT_PEERS'"/g' ~/.trst/config/config.toml
sed -i 's/trust_period = "168h0m0s"/trust_period = "168h"/g' ~/.trst/config/config.toml
echo "Set persistent_peers: $PERSISTENT_PEERS"

echo "Waiting for bootstrap to start..."
sleep 20

trstcli q block 1

cp /tmp/.trstd/keyring-test /root/.trstd/ -r

# MASTER_KEY="$(trstcli q register node-enclave-params 2> /dev/null | cut -c 3- )"

#echo "Master key: $MASTER_KEY"

trstd init-enclave --reset

PUBLIC_KEY=$(trstd parse /opt/trustlesshub/.sgx_secrets/attestation_cert.der | cut -c 3- )

echo "Public key: $PUBLIC_KEY"

trstd parse /opt/trustlesshub/.sgx_secrets/attestation_cert.der
cat /opt/trustlesshub/.sgx_secrets/attestation_cert.der
tx_hash="$(trstcli tx register auth /opt/trustlesshub/.sgx_secrets/attestation_cert.der -y --from a --gas-prices 0.25utrst | jq -r '.txhash')"

#trstcli q tx "$tx_hash"
sleep 15
trstcli q tx "$tx_hash"

SEED="$(trstcli q register seed "$PUBLIC_KEY" | cut -c 3-)"
echo "SEED: $SEED"
#exit

trstcli q register node-enclave-params

trstd configure-secret node-master-cert.der "$SEED"

cp /tmp/.trst/config/genesis.json /root/.trst/config/genesis.json

trstd validate-genesis

RUST_BACKTRACE=1 trstd start --rpc.laddr tcp://0.0.0.0:26657

# ./wasmi-sgx-test.sh
