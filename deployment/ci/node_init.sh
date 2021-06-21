#!/usr/bin/env bash

set -euv

# init the node
# rm -rf ~/.secret*
#tppd config chain-id enigma-testnet
#tppd config output json
#tppd config indent true
#tppd config trust-node true
#tppd config keyring-backend test
rm -rf ~/.tppd

mkdir -p /root/.tppd/.node

tppd init "$(hostname)" --chain-id enigma-testnet || true

PERSISTENT_PEERS=''

sed -i 's/persistent_peers = ""/persistent_peers = "'$PERSISTENT_PEERS'"/g' ~/.tppd/config/config.toml
echo "Set persistent_peers: $PERSISTENT_PEERS"

echo "Waiting for bootstrap to start..."
sleep 20

# MASTER_KEY="$(tppd q register tpp-enclave-params --node http://bootstrap:26657 2> /dev/null | cut -c 3- )"

#echo "Master key: $MASTER_KEY"

tppd init-enclave

PUBLIC_KEY=$(tppd parse attestation_cert.der 2> /dev/null | cut -c 3- )

echo "Public key: $(tppd parse attestation_cert.der 2> /dev/null | cut -c 3- )"

tppd tx register auth attestation_cert.der --node http://bootstrap:26657 -y --from a

sleep 10

SEED=$(tppd q register seed "$PUBLIC_KEY" --node http://bootstrap:26657 2> /dev/null | cut -c 3-)
echo "SEED: $SEED"

tppd q register tpp-enclave-params --node http://bootstrap:26657 2> /dev/null

tppd configure-secret node-master-cert.der "$SEED"

cp /tmp/.tppd/config/genesis.json /root/.tppd/config/genesis.json

tppd validate-genesis

RUST_BACKTRACE=1 tppd start &

./wasmi-sgx-test.sh