#!/bin/sh

set -o errexit

CHAINID=$1
GENACCT=$2

if [ -z "$1" ]; then
  echo "Need to input chain id..."
  exit 1
fi

rm -rf ~/.trstd

# Build genesis file incl account for passed address
coins="10000000000utrst,100000000000stake"
trstd init --chain-id $CHAINID $CHAINID
trstd keys add validator --keyring-backend="test"
trstd add-genesis-account $(trstd keys show validator -a --keyring-backend="test") $coins

if [ ! -z "$2" ]; then
  trstd add-genesis-account $GENACCT $coins
fi

trstd gentx validator 5000000000utrst --keyring-backend="test" --chain-id $CHAINID
trstd collect-gentxs

# Set proper defaults and change ports
sed -i 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:26657"#g' ~/.trstd/config/config.toml
sed -i 's/timeout_commit = "5s"/timeout_commit = "1s"/g' ~/.trstd/config/config.toml
sed -i 's/timeout_propose = "3s"/timeout_propose = "1s"/g' ~/.trstd/config/config.toml
sed -i 's/index_all_keys = false/index_all_keys = true/g' ~/.trstd/config/config.toml
perl -i -pe 's/"stake"/ "utrst"/g' ~/.trstd/config/genesis.json

# Start the trstd
trstd start --pruning=nothing --bootstrap