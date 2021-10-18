#!/usr/bin/env bash

set -euv

# init the node
# rm -rf ~/.secret*
#trstd config chain-id enigma-testnet
#trstd config output json
#trstd config indent true
#trstd config trust-node true
#trstd config keyring-backend test
rm -rf ~/.trstd

mkdir -p /root/.trstd/.node

trstd init "$(hostname)" --chain-id enigma-testnet || true

PERSISTENT_PEERS=''

sed -i 's/persistent_peers = ""/persistent_peers = "'$PERSISTENT_PEERS'"/g' ~/.trstd/config/config.toml
echo "Set persistent_peers: $PERSISTENT_PEERS"

echo "Waiting for bootstrap to start..."
sleep 20

# MASTER_KEY="$(trstd q register trst-enclave-params --node http://bootstrap:26657 2> /dev/null | cut -c 3- )"

#echo "Master key: $MASTER_KEY"

trstd init-enclave

PUBLIC_KEY=$(trstd parse attestation_cert.der 2> /dev/null | cut -c 3- )

echo "Public key: $(trstd parse attestation_cert.der 2> /dev/null | cut -c 3- )"

trstd tx register auth attestation_cert.der --node http://bootstrap:26657 -y --from a

sleep 10

SEED=$(trstd q register seed "$PUBLIC_KEY" --node http://bootstrap:26657 2> /dev/null | cut -c 3-)
echo "SEED: $SEED"

trstd q register trst-enclave-params --node http://bootstrap:26657 2> /dev/null

trstd configure-secret node-master-cert.der "$SEED"

cp /tmp/.trstd/config/genesis.json /root/.trstd/config/genesis.json

trstd validate-genesis

RUST_BACKTRACE=1 trstd start &

./wasmi-sgx-test.sh