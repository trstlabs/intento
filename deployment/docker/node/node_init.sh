#!/usr/bin/env bash
set -euv

# REGISTRATION_SERVICE=
# RPC_URL=http://bootstrap:26657
# CHAINID=secret-testnet-1
# PERSISTENT_PEERS=115aa0a629f5d70dd1d464bc7e42799e00f4edae@bootstrap:26656

# init the node
# rm -rf ~/.secret*
#tppd config chain-id enigma-testnet
#tppd config output json
#tppd config indent true
#tppd config trust-node true
#tppd config keyring-backend test
# rm -rf ~/.tppd
file=/root/.tppd/config/attestation_cert.der
if [ ! -e "$file" ]
then
  rm -rf ~/.tppd/* || true

  mkdir -p /root/.tppd/.node
  export SECRET_NETWORK_CHAIN_ID=$CHAINID
  export SECRET_NETWORK_KEYRING_BACKEND=test
  # tppd init "$(hostname)" --chain-id enigma-testnet || true

  tppd init "$MONIKER" --chain-id "$CHAINID"
  echo "Initializing chain: $CHAINID with node moniker: $(hostname)"

  sed -i 's/persistent_peers = ""/persistent_peers = "'"$PERSISTENT_PEERS"'"/g' ~/.tppd/config/config.toml
  echo "Set persistent_peers: $PERSISTENT_PEERS"
  
  # Open RPC port to all interfaces
  perl -i -pe 's/laddr = .+?26657"/laddr = "tcp:\/\/0.0.0.0:26657"/' ~/.tppd/config/config.toml

  # Open P2P port to all interfaces
  perl -i -pe 's/laddr = .+?26656"/laddr = "tcp:\/\/0.0.0.0:26656"/' ~/.tppd/config/config.toml

  echo "Waiting for bootstrap to start..."
  sleep 10

  tppd init-enclave

  PUBLIC_KEY=$(tppd parse attestation_cert.der 2> /dev/null | cut -c 3- )

  echo "Public key: $(tppd parse attestation_cert.der 2> /dev/null | cut -c 3- )"

  cp attestation_cert.der /root/.tppd/config/

  openssl base64 -A -in attestation_cert.der -out b64_cert
  # tppd tx register auth attestation_cert.der --node "$RPC_URL" -y --from a
  curl -G --data-urlencode "cert=$(cat b64_cert)" http://"$REGISTRATION_SERVICE"/register

  sleep 20

  SEED=$(tppd q register seed "$PUBLIC_KEY" --node tcp://"$RPC_URL" 2> /dev/null | cut -c 3-)
  echo "SEED: $SEED"

  tppd q register tpp-enclave-params --node tcp://"$RPC_URL" 2> /dev/null

  tppd configure-secret node-master-cert.der "$SEED"

  curl http://"$RPC_URL"/genesis | jq -r .result.genesis > /root/.tppd/config/genesis.json

  echo "Downloaded genesis file from $RPC_URL "

  tppd validate-genesis

fi
tppd start
