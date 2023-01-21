#!/bin/bash

# Configure predefined mnemonic pharses
BINARY=rly
CHAIN_DIR=./data
CHAINID_1=trstdev-1
CHAINID_2=trstdev-2
RELAYER_DIR=./relayer
MNEMONIC_1="grant rice replace explain federal release fix clever romance raise often wild taxi quarter soccer fiber love must tape steak together observe swap guitar"
MNEMONIC_2="jelly shadow frog dirt dragon use armed praise universe win jungle close inmate rain oil canvas beauty pioneer chef soccer icon dizzy thunder meadow"

# Ensure rly is installed
if ! [ -x "$(command -v $BINARY)" ]; then
    echo "$BINARY is required to run this script..."
    echo "You can download at https://github.com/cosmos/relayer"
    exit 1
fi

echo "Starting to listen relayer..."
$BINARY start trstdev1-trstdev2 -p events -b 100 --home $CHAIN_DIR/$RELAYER_DIR
