#!/bin/bash

set -eu

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)

source ${SCRIPT_DIR}/../config.sh

into_chain_id=$(GET_VAR_VALUE INTO_CHAIN_ID)
hermes_logs=${LOGS}/hermes.log

for chain in ${RLY_HOST_CHAINS[@]}; do
    hermes_exec=$(GET_VAR_VALUE HERMES_EXEC)
    chain_name=$(printf "$chain" | awk '{ print tolower($0) }')
    mnemonic=$(GET_VAR_VALUE RELAYER_${chain}_MNEMONIC)
    chain_id=$(GET_VAR_VALUE ${chain}_CHAIN_ID)
    hermes_config=$STATE/hermes

    mkdir -p $hermes_config
    chmod -R 777 $STATE/hermes
    cp ${DOCKERNET_HOME}/config/hermes_config.toml $hermes_config/config.toml

    printf "Adding Hermes keys for intento and $chain...\n"
    echo $mnemonic | $hermes_exec hermes keys add --chain ${into_chain_id} --overwrite --mnemonic-file /dev/stdin >>$hermes_logs 2>&1
    echo $mnemonic | $hermes_exec hermes keys add --chain ${chain_id} --mnemonic-file  /dev/stdin >>$hermes_logs 2>&1


    echo "Verifying balances for $chain..."
    $hermes_exec hermes keys balance --chain ${chain_id} >>$hermes_logs 2>&1
    
    echo "Creating channel..."
    $hermes_exec hermes create channel --a-chain ${into_chain_id} --b-chain ${chain_id} --a-port transfer --b-port transfer --new-client-connection --yes >>$hermes_logs 2>&1

    echo "Creating connection..."
    $hermes_exec hermes create connection --a-chain ${into_chain_id} --b-chain ${chain_id} >>$hermes_logs 2>&1
done

echo "Starting relayer proces..."
$hermes_exec hermes start >>$hermes_logs 2>&1 &