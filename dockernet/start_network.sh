#!/bin/bash

set -eu 
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
source ${SCRIPT_DIR}/config.sh

# Confirm binaries are present
for chain in INTO ${HOST_CHAINS[@]}; do
    binary_path=$(GET_VAR_VALUE ${chain}_BINARY)
    binary_path=$(realpath "$binary_path")
    if [[ ! -e "$binary_path" ]]; then
        echo "ERROR: Binary for $chain does not exist"
        echo "It should be present at $binary_path"
        echo "To build the binary, ensure submodules are updated and pass the host chain flag as a build argument (e.g. 'make start-dockernet build=g')"
        exit 1
    fi
done

# cleanup any stale state
rm -rf $STATE $LOGS 
mkdir -p $STATE
mkdir -p $LOGS


# If we're testing an upgrade, setup cosmovisor
if [[ "${UPGRADE_NAME:-}" != "" ]]; then
    printf "\n>>> UPGRADE ENABLED! ($UPGRADE_NAME)\n\n"
    
    # Update binary #2 with the binary that was just compiled
    mkdir -p $UPGRADES/binaries
    rm -f $UPGRADES/binaries/intentod2
    cp $DOCKERNET_HOME/../build/intentod $UPGRADES/binaries/intentod2

    # Build a cosmovisor image with the old binary and replace the into docker image with a new one
    #  that has both binaries and is running cosmovisor
    # The reason for having a separate cosmovisor image is so we can cache the building of cosmovisor and the old binary
    echo "Building Cosmovisor..."
    docker build \
        -t intento:cosmovisor \
        --build-arg old_commit_hash=$UPGRADE_OLD_VERSION \
        --build-arg into_admin_address=$INTO_ADMIN_ADDRESS \
        -f $UPGRADES/Dockerfile.cosmovisor .

    echo "Re-Building INTO with Upgrade Support..."
    docker build \
        -t intento:into \
        --build-arg upgrade_name=$UPGRADE_NAME \
        -f $UPGRADES/Dockerfile.into .

    echo "Done"
fi


# Initialize the state for each chain
for chain in INTO ${HOST_CHAINS[@]}; do
    echo "Initializing $chain..."
    bash $SRC/init_chain.sh $chain
done


echo "Start the chain..."
bash $SRC/start_chain.sh 

echo "Start relayers and the transfer channels..."
bash $SRC/start_relayers.sh 
