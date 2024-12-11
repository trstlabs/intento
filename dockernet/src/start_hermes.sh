#!/bin/bash

set -eu

# Ensure the Consumer ID is passed as an argument
if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <Consumer ID>"
    exit 1
fi

# Variables
CONSUMER_ID=$1
SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)
source ${SCRIPT_DIR}/../config.sh

# Paths and executables
hermes_logs=${LOGS}/hermes.log
hermes_exec=$(GET_VAR_VALUE HERMES_EXEC)

# Clear previous logs
: >$hermes_logs

echo "$(date '+%Y-%m-%d %H:%M:%S') - Starting Hermes relayer setup for Consumer ID: $CONSUMER_ID" >>$hermes_logs

# Loop through chains
for chain in ${RLY_HOST_CHAINS[@]}; do
    chain_name=$(printf "$chain" | awk '{ print tolower($0) }')
    mnemonic=$(GET_VAR_VALUE RELAYER_${chain}_MNEMONIC)
    chain_id=$(GET_VAR_VALUE ${chain}_CHAIN_ID)
    hermes_config=$STATE/hermes

    mkdir -p $hermes_config
    chmod -R 777 $STATE/hermes
    cp ${DOCKERNET_HOME}/config/hermes_config.toml $hermes_config/config.toml

    echo "$(date '+%Y-%m-%d %H:%M:%S') - Adding Hermes keys for $INTO_CHAIN_ID and $chain_id..." >>$hermes_logs
    echo $mnemonic | $hermes_exec hermes keys add --chain ${INTO_CHAIN_ID} --overwrite --mnemonic-file /dev/stdin >>$hermes_logs 2>&1
    echo $mnemonic | $hermes_exec hermes keys add --chain ${chain_id} --mnemonic-file /dev/stdin >>$hermes_logs 2>&1

    echo "$(date '+%Y-%m-%d %H:%M:%S') - Verifying balances for $chain_id..." >>$hermes_logs
    $hermes_exec hermes keys balance --chain ${chain_id} >>$hermes_logs 2>&1

    echo "$(date '+%Y-%m-%d %H:%M:%S') - Creating IBC transfer channel between $CONSUMER_ID and $chain_id..." >>$hermes_logs
    if ! $hermes_exec hermes create channel --a-chain ${INTO_CHAIN_ID} --b-chain ${chain_id} \
      --a-port transfer --b-port transfer --new-client-connection --yes >>$hermes_logs 2>&1; then
          echo "Error: Failed to create transfer channel between ${INTO_CHAIN_ID} and ${chain_id}" >>$hermes_logs
          exit 1
    fi

    echo "$(date '+%Y-%m-%d %H:%M:%S') - Creating IBC connection between $INTO_CHAIN_ID and $chain_id..." >>$hermes_logs
    $hermes_exec hermes create connection --a-chain ${INTO_CHAIN_ID} --b-chain ${chain_id} >>$hermes_logs 2>&1
done

# Special case: Create the CCV channel to GAIA
echo "$(date '+%Y-%m-%d %H:%M:%S') - Creating CCV channel to GAIA..." >>$hermes_logs
if ! $hermes_exec hermes create channel --a-chain ${INTO_CHAIN_ID} --b-chain GAIA_CHAIN_ID \
  --a-port consumer --b-port provider --order ordered --new-client-connection --channel-version 1 >>$hermes_logs 2>&1; then
      echo "Error: Failed to create CCV channel between ${INTO_CHAIN_ID} and GAIA." >>$hermes_logs
      exit 1
fi

echo "$(date '+%Y-%m-%d %H:%M:%S') - Starting Hermes relayer process..." >>$hermes_logs
$hermes_exec hermes start >>$hermes_logs 2>&1 &

echo "$(date '+%Y-%m-%d %H:%M:%S') - Hermes relayer process started successfully." >>$hermes_logs
