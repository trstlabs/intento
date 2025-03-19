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

    # echo $hermes_exec hermes --json query clients --host-chain ${INTO_CHAIN_ID}
    # Check if client already exists
    existing_client=$($hermes_exec hermes --json query clients --host-chain ${INTO_CHAIN_ID} | jq -r '.clients[] | select(.client_id.client_id == "07-tendermint-0") | .client_id.client_id' 2>/dev/null || true)

    if [ -n "$existing_client" ]; then
        echo $existing_client
        echo "$(date '+%Y-%m-%d %H:%M:%S') - Reusing existing client: $existing_client" >>$hermes_logs
        client_id=$existing_client
    else
        echo "$(date '+%Y-%m-%d %H:%M:%S') - Creating a new client..." >>$hermes_logs
        client_id="07-tendermint-0"
        #client_id=$($hermes_exec hermes --json create client --host-chain ${INTO_CHAIN_ID} --reference-chain ${chain_id} | jq -r '.result.client_id' 2>/dev/null || echo "fallback-client-id")
    fi
    #echo $client_id
    echo "$(date '+%Y-%m-%d %H:%M:%S') - Creating IBC connection between $INTO_CHAIN_ID and $chain_id..." >>$hermes_logs
    connection_output=$($hermes_exec hermes --json create connection --a-chain ${INTO_CHAIN_ID} --a-client ${client_id} --b-client ${client_id})

    echo "$connection_output" >>$hermes_logs
    a_side_connection_id=$(echo "$connection_output" | jq -r '.result.a_side.connection_id' | grep -v '^null$' | tail -n 1)
    if [ -z "$a_side_connection_id" ]; then
        echo "Error: Connection ID is empty or null" >>$hermes_logs
        exit 1
    fi
    echo $a_side_connection_id
    echo "$(date '+%Y-%m-%d %H:%M:%S') - Creating IBC channel with connection ID: $a_side_connection_id..." >>$hermes_logs
    channel_output=$($hermes_exec hermes create channel --a-chain ${INTO_CHAIN_ID} --a-connection $a_side_connection_id --a-port consumer --b-port provider --order ordered --channel-version 1)
    echo "$channel_output" >>$hermes_logs
done

echo "$(date '+%Y-%m-%d %H:%M:%S') - Starting Hermes relayer process..." >>$hermes_logs
$hermes_exec hermes start >>$hermes_logs 2>&1 &

echo "$(date '+%Y-%m-%d %H:%M:%S') - Hermes relayer process started successfully." >>$hermes_logs
