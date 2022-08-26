#!/usr/bin/env bash
set -euv

# REGISTRATION_SERVICE=
# export RPC_URL="bootstrap:26657"
# export CHAINID="trst_chain_1"
# export PERSISTENT_PEERS="115aa0a629f5d70dd1d464bc7e42799e00f4edae@bootstrap:26656"

# init the node
# rm -rf ~/.secret*

# rm -rf ~/.trstd
file=/root/.trst/config/attestation_cert.der
if [ ! -e "$file" ]
then
  rm -rf ~/.trstd/* || true

  mkdir -p /root/.trstd/.node
  # trstd config keyring-backend test
  trstd config node tcp://"$RPC_URL"
  trstd config chain-id "$CHAINID"
#  export SECRET_NETWORK_CHAIN_ID=$CHAINID
#  export SECRET_NETWORK_KEYRING_BACKEND=test
  # trstd init "$(hostname)" --chain-id enigma-testnet || true

  trstd init "$MONIKER" --chain-id "$CHAINID"

  # cp /tmp/.trstd/keyring-test /root/.trstd/ -r

  echo "Initializing chain: $CHAINID with node moniker: $(hostname)"

  sed -i 's/persistent_peers = ""/persistent_peers = "'"$PERSISTENT_PEERS"'"/g' ~/.trst/config/config.toml
  echo "Set persistent_peers: $PERSISTENT_PEERS"

  # Open RPC port to all interfaces
  perl -i -pe 's/laddr = .+?26657"/laddr = "tcp:\/\/0.0.0.0:26657"/' ~/.trst/config/config.toml

  # Open P2P port to all interfaces
  perl -i -pe 's/laddr = .+?26656"/laddr = "tcp:\/\/0.0.0.0:26656"/' ~/.trst/config/config.toml

  echo "Waiting for bootstrap to start..."
  sleep 10

  trstd init-enclave

  PUBLIC_KEY=$(trstd parse /opt/trustlesshub/.sgx_secrets/attestation_cert.der 2> /dev/null | cut -c 3- )

  echo "Public key: $(trstd parse /opt/trustlesshub/.sgx_secrets/attestation_cert.der 2> /dev/null | cut -c 3- )"

  cp /opt/trustlesshub/.sgx_secrets/attestation_cert.der /root/.trst/config/

  openssl base64 -A -in attestation_cert.der -out b64_cert
  # trstd tx register auth attestation_cert.der --from a --gas-prices 0.25utrst -y

  curl -G --data-urlencode "cert=$(cat b64_cert)" http://"$REGISTRATION_SERVICE"/register

  sleep 20

  SEED=$(trstd q register seed "$PUBLIC_KEY"  2> /dev/null | cut -c 3-)
  echo "SEED: $SEED"

  trstd q register node-enclave-params 2> /dev/null

  trstd configure-secret node-master-cert.der "$SEED"

  curl http://"$RPC_URL"/genesis | jq -r .result.genesis > /root/.trst/config/genesis.json

  echo "Downloaded genesis file from $RPC_URL "

  trstd validate-genesis

  trstd config node tcp://localhost:26657

fi
trstd start
