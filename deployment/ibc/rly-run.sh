#!/bin/bash

# Define variables
TRSTDEV1_CHAIN_CONFIG="/rlyconfig/chains/trstdev-1.json"
TRSTDEV2_CHAIN_CONFIG="/rlyconfig/chains/trstdev-2.json"
LOCALJUNO_CHAIN_CONFIG="/rlyconfig/chains/juno.json"
TRSTDEV1_TRSTDEV2_PATH_CONFIG="/rlyconfig/paths/trstdev1-trstdev2.json"
TRSTDEV1_LOCALJUNO_PATH_CONFIG="/rlyconfig/paths/trstdev1-localjuno.json"
TRSTDEV1="trstdev-1"
LOCALJUNO="localjuno"
TRSTDEV2="trstdev-2"
TEST_KEY="testkey"
TRSTDEV1_TRSTDEV2_PATH="trstdev1-trstdev2"
TRSTDEV1_LOCALJUNO_PATH="trstdev1-localjuno"

# Initialize config
rly config init

# Add chains
rly chains add --file $TRSTDEV1_CHAIN_CONFIG $TRSTDEV1
rly chains add --file $TRSTDEV2_CHAIN_CONFIG $TRSTDEV2
rly chains add --file $LOCALJUNO_CHAIN_CONFIG $LOCALJUNO

# Restore keys
restore_key() {
  rly keys restore $1 $TEST_KEY "$2"
}
restore_key $TRSTDEV1 "grant rice replace explain federal release fix clever romance raise often wild taxi quarter soccer fiber love must tape steak together observe swap guitar"
restore_key $LOCALJUNO "crop staff genuine enjoy dial pact sorry bless note fall abuse more cheese clutch make ripple machine this gravity lend thank marine sell print"
restore_key $TRSTDEV2 "jelly shadow frog dirt dragon use armed praise universe win jungle close inmate rain oil canvas beauty pioneer chef soccer icon dizzy thunder meadow"

# Add paths
rly paths add $TRSTDEV1 $TRSTDEV1 $TRSTDEV1_TRSTDEV2_PATH --file $TRSTDEV1_TRSTDEV2_PATH_CONFIG
rly paths add $TRSTDEV1 testing $TRSTDEV1_LOCALJUNO_PATH --file $TRSTDEV1_LOCALJUNO_PATH_CONFIG

# Create connections
rly tx connection $TRSTDEV1_TRSTDEV2_PATH 
rly tx connection $TRSTDEV1_LOCALJUNO_PATH

# Create transfer channel 
rly tx channel $TRSTDEV1_TRSTDEV2_PATH
rly tx channel $TRSTDEV1_LOCALJUNO_PATH

# Start relayer
rly start $TRSTDEV1_TRSTDEV2_PATH $TRSTDEV1_LOCALJUNO_PATH -p events -b 100 --debug > rly.log
