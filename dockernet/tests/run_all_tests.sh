#!/bin/bash
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

# run test files
BATS=${SCRIPT_DIR}/bats/bats-core/bin/bats
INTEGRATION_TEST_FILE=${SCRIPT_DIR}/integration_tests.bats 

# connection_id is the i of the path
CHAIN_NAME=GAIA TRANSFER_CHANNEL_NUMBER=1 CONNECTION_ID=0 HOST_CONNECTION_ID=0 HOST_TRANSFER_CHANNEL="channel-1" $BATS $INTEGRATION_TEST_FILE
#CHAIN_NAME=OSMO TRANSFER_CHANNEL_NUMBER=1 CONNECTION_ID=3 HOST_CONNECTION_ID=1 HOST_TRANSFER_CHANNEL="channel-0" $BATS $INTEGRATION_TEST_FILE
# CHAIN_NAME=HOST TRANSFER_CHANNEL_NUMBER=2 CONNECTION_ID=2 HOST_CONNECTION_ID=0 HOST_TRANSFER_CHANNEL="channel-0" $BATS $INTEGRATION_TEST_FILE