
#!/bin/bash

set -eu

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
source ${SCRIPT_DIR}/../config.sh

# Paths and constants
PROVIDER_HOME="$DOCKERNET_HOME/state/${GAIA_NODE_PREFIX}1"
CONSUMER_HOME_PREFIX="$DOCKERNET_HOME/state/${INTO_NODE_PREFIX}"
CONSUMER_HOME="${CONSUMER_HOME_PREFIX}1"
PROVIDER_BINARY=$GAIA_BINARY
PROVIDER_CHAIN_ID=$GAIA_CHAIN_ID
PROVIDER_RPC_ADDR="localhost:$GAIA_RPC_PORT"
DENOM=$ATOM_DENOM
PROVIDER_CMD="$PROVIDER_BINARY --home $PROVIDER_HOME --fees 200000uatom"
PROVIDER_CMD_Q="$PROVIDER_BINARY --home $PROVIDER_HOME q"
SOVEREIGN_CHAIN_ID=$INTO_CHAIN_ID
SOVEREIGN_HOME="$DOCKERNET_HOME/state/sovereign"


# Build consumer chain proposal file
tee $PROVIDER_HOME/consumer-create.json<<EOF
{
  "chain_id": "$SOVEREIGN_CHAIN_ID",
  "metadata": {
	"name": "consumer chain",
	"description": "no description",
	"metadata": "no metadata"
  },
  "initialization_parameters": {
	"initial_height": {
  	"revision_number": 1,
  	"revision_height": 1
	},
	"genesis_hash": "",
	"binary_hash": "",
	"spawn_time": null,
  "unbonding_period": 3456000000000000,
  "ccv_timeout_period": 4838400000000000,
  "transfer_timeout_period": 43200000000000,
  "consumer_redistribution_fraction": "1.0",
  "blocks_per_distribution_transmission": 5000,
  "historical_entries": 10000
},
  "power_shaping_parameters": {
	"top_N": 0,
	"validators_power_cap": 0,
	"validator_set_cap": 0,
	"allowlist": [],
	"denylist": [],
	"min_stake": 0,
	"allow_inactive_vals": true
  }
}
EOF

# Step 1: Submit create-consumer transaction
echo "Submitting create-consumer transaction..."
TX_RES=$($PROVIDER_CMD tx provider create-consumer $PROVIDER_HOME/consumer-create.json \
    --chain-id $PROVIDER_CHAIN_ID --node tcp://$PROVIDER_RPC_ADDR \
--from ${GAIA_VAL_PREFIX}1 --keyring-backend test -y --log_format json)

# echo $TX_RES
sleep 5
# # Verify success
# if [[ $(echo $TX_RES | jq -r '.code') -ne 0 ]]; then
#   echo "Error: Failed to submit create-consumer transaction."
#   exit 1
# fi

# Step 2: Fetch the consumer_id
echo "Fetching consumer_id..."
CONSUMER_ID=$($PROVIDER_CMD_Q provider list-consumer-chains --node tcp://$PROVIDER_RPC_ADDR -o json | jq -r '.chains[-1].consumer_id')

if [[ -z "$CONSUMER_ID" ]]; then
    echo "Error: Unable to fetch consumer_id."
    exit 1
fi
echo "Consumer ID: $CONSUMER_ID"
sleep 5
TX_RES=$($PROVIDER_CMD tx provider opt-in $CONSUMER_ID \
    --chain-id $PROVIDER_CHAIN_ID --node tcp://$PROVIDER_RPC_ADDR \
--from ${GAIA_VAL_PREFIX}1 --keyring-backend test -y --log_format json)

# echo $TX_RES
sleep 5

# Get current local time in the correct format (local timezone with the required Z07:00 format)
LAUNCH_DATE=$(date +"%Y-%m-%dT%H:%M:%S.000000%:z")

# Step 3: Submit update-consumer transaction
echo "Submitting update-consumer transaction..."
tee ${PROVIDER_HOME}/update-consumer.json <<EOF
{
  "chain_id": "$CONSUMER_ID",
  "consumer_id": "$CONSUMER_ID",
  "metadata": {
    "name": "consumer chain",
    "description": "no description",
    "metadata": "no metadata"
  },
  "initialization_parameters": {
    "initial_height": {
      "revision_number": 1,
      "revision_height": 1
    },
    "genesis_hash": "",
    "binary_hash": "",
    "spawn_time": "$LAUNCH_DATE",
    "unbonding_period": 3456000000000000,
    "ccv_timeout_period": 4838400000000000,
    "transfer_timeout_period": 43200000000000,
    "consumer_redistribution_fraction": "1.0",
    "blocks_per_distribution_transmission": 5000,
    "historical_entries": 10000
},
  "power_shaping_parameters": {
    "top_N": 0,
    "validators_power_cap": 0,
    "validator_set_cap": 0,
    "allowlist": [],
    "denylist": [],
    "min_stake": 0,
    "allow_inactive_vals": true
  }
}
EOF

TX_RES=$($PROVIDER_CMD tx provider update-consumer ${PROVIDER_HOME}/update-consumer.json \
    --chain-id $PROVIDER_CHAIN_ID --node tcp://$PROVIDER_RPC_ADDR \
--from ${GAIA_VAL_PREFIX}1 --keyring-backend test -y --log_format json)

echo $TX_RES
sleep 10
# Verify success
# if [[ $(echo $TX_RES | jq -r '.code') -ne 0 ]]; then
#   echo "Error: Failed to submit update-consumer transaction."
#   exit 1
# fi
# echo "Successfully updated consumer chain."
#sleep 5
# Step 4: Announce launch
# echo "Announce the launch on Discord and submit a PR to the Cosmos testnets repository."

# Step 5: Update genesis file after chain launch
echo "Updating genesis file with CCV state..."

$PROVIDER_CMD_Q provider consumer-genesis $CONSUMER_ID -o json --node tcp://$PROVIDER_RPC_ADDR > ccv-state.json
jq '.params.reward_denoms |= ["'$INTO_DENOM'"]' ccv-state.json > ccv-denom.json
jq '.params.provider_reward_denoms |= ["'$GAIA_DENOM'"]' ccv-denom.json > ccv-provider-denom.json
jq -s '.[0].app_state.ccvconsumer = .[1] | .[0]' ccv-state.json ccv-provider-denom.json > consumer-genesis-updated.json
echo "Genesis file updated."



echo "Consumer chain setup complete! ID: $CONSUMER_ID "


#!/bin/bash

set -eu
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

source ${SCRIPT_DIR}/../config.sh


log_file=$DOCKERNET_HOME/logs/${INTO_NODE_PREFIX}.log

# Step 6: Starting INTO chain
echo "Starting $SOVEREIGN_CHAIN_ID chain"
nodes_names=$(i=1; while [ $i -le $INTO_NUM_NODES ]; do printf "%s " ${INTO_NODE_PREFIX}${i}; i=$(($i + 1)); done;)
$DOCKER_COMPOSE up -d $nodes_names

$DOCKER_COMPOSE logs -f ${INTO_NODE_PREFIX}1 | sed -r -u "s/\x1B\[([0-9]{1,3}(;[0-9]{1,2})?)?[mGK]//g" > $log_file 2>&1 &

printf "Waiting for $SOVEREIGN_CHAIN_ID to start..."


log_file=$DOCKERNET_HOME/logs/${INTO_NODE_PREFIX}.log

( tail -f -n0 $log_file & ) | grep -q "finalizing commit of block"
echo "Done"

sleep 5

# Step 7: Create IBC connections and channels
echo "Creating IBC connections and channels..."

bash $SRC/start_hermes.sh $CONSUMER_ID
