#!/bin/bash

set -eu
DOCKERNET_HOME=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)

STATE=$DOCKERNET_HOME/state
LOGS=$DOCKERNET_HOME/logs
UPGRADES=$DOCKERNET_HOME/upgrades
SRC=$DOCKERNET_HOME/src
PEER_PORT=26656
DOCKER_COMPOSE="docker-compose -f $DOCKERNET_HOME/docker-compose.yml"

# Logs
INTO_LOGS=$LOGS/into.log
TX_LOGS=$DOCKERNET_HOME/logs/tx.log
KEYS_LOGS=$DOCKERNET_HOME/logs/keys.log

# List of hosts enabled
HOST_CHAINS=()

# If no host zones are specified above:
#  `start-docker` defaults to just GAIA if HOST_CHAINS is empty
#  `start-docker-all` runs 3 hosts
# Available host zones:
#  - GAIA
#  - OSMO
#  - HOST (our chain enabled as a host zone)
if [[ "${ALL_HOST_CHAINS:-false}" == "true" ]]; then
  HOST_CHAINS=(GAIA OSMO) #all
  RLY_HOST_CHAINS=(GAIA OSMO)
elif [[ "${#HOST_CHAINS[@]}" == "0" ]]; then
  HOST_CHAINS=(GAIA )
  RLY_HOST_CHAINS=(GAIA) #you can add a testnet chain here (if configured)
fi

# DENOMS
INTO_DENOM="uinto"
ATOM_DENOM="uatom"
OSMO_DENOM="uosmo"
COSM_DENOM="ucosm"


IBC_INTO_DENOM='ibc/F1B5C3489F881CC56ECC12EA903EFCF5D200B4D8123852C191A88A31AC79A8E4'

IBC_GAIA_CHANNEL_0_DENOM='ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2'
IBC_GAIA_CHANNEL_1_DENOM='ibc/C4CFF46FD6DE35CA4CF4CE031E643C8FDC9BA4B99AE598E9B0ED98FE3A2319F9'
IBC_GAIA_CHANNEL_2_DENOM='ibc/9117A26BA81E29FA4F78F57DC2BD90CD3D26848101BA880445F119B22A1E254E'
IBC_GAIA_CHANNEL_3_DENOM='ibc/A4DB47A9D3CF9A068D454513891B526702455D3EF08FB9EB558C561F9DC2B701'

IBC_OSMO_CHANNEL_0_DENOM='ibc/ED07A3391A112B175915CD8FAF43A2DA8E4790EDE12566649D0C2F97716B8518'
IBC_OSMO_CHANNEL_1_DENOM='ibc/0471F1C4E7AFD3F07702BEF6DC365268D64570F7C1FDC98EA6098DD6DE59817B'
IBC_OSMO_CHANNEL_2_DENOM='ibc/13B2C536BB057AC79D5616B8EA1B9540EC1F2170718CAFF6F0083C966FFFED0B'
IBC_OSMO_CHANNEL_3_DENOM='ibc/47BD209179859CDE4A2806763D7189B6E6FE13A17880FE2B42DE1E6C1E329E23'

IBC_HOST_CHANNEL_0_DENOM='ibc/82DBA832457B89E1A344DA51761D92305F7581B7EA6C18D85037910988953C58'
IBC_HOST_CHANNEL_1_DENOM='ibc/FB7E2520A1ED6890E1632904A4ACA1B3D2883388F8E2B88F2D6A54AA15E4B49E'
IBC_HOST_CHANNEL_2_DENOM='ibc/D664DC1D38648FC4C697D9E9CF2D26369318DFE668B31F81809383A8A88CFCF4'
IBC_HOST_CHANNEL_3_DENOM='ibc/FD7AA7EB2C1D5D97A8693CCD71FFE3F5AFF12DB6756066E11E69873DE91A33EA'

# COIN TYPES
# Coin types can be found at https://github.com/satoshilabs/slips/blob/master/slip-0044.md
COSMOS_COIN_TYPE=118
#ETH_COIN_TYPE=60



# CHAIN PARAMS
MIN_GAS="0.0001000000"
BLOCK_TIME="2s"
UNBONDING_TIME="14400s"
MAX_DEPOSIT_PERIOD="30s"
VOTING_PERIOD="300s"
INITIAL_ANNUAL_PROVISIONS="10000000000000.000000000000000000"

# Tokens are denominated in the macro-unit
# (e.g. 5000000INTO implies 5000000000000uinto)
GENESIS_TOKENS=40000000
STAKE_TOKENS=30000000
ADMIN_TOKENS=10000
FAUCET_TOKENS=1000000

# faucet
# TEST_FAUCET_ACCT=faucet2
# TEST_FAUCET_MNEMONIC="wreck valve enroll onion space weasel cherry ketchup edge special certain silver"


# # CHAIN MNEMONICS
VAL_MNEMONIC_1="join secret blind dose prepare atom wrestle funny want memory spare captain empty speak logic wrestle half develop brown economy burden north example slide"
TEST_FAUCET_ACCT=faucet
TEST_FAUCET_MNEMONIC="word twist toast cloth movie predict advance crumble escape whale sail such angry muffin balcony keen move employ cook valve hurt glimpse breeze brick"


# CHAIN MNEMONICS
# VAL_MNEMONIC_1="close soup mirror crew erode defy knock trigger gather eyebrow tent farm gym gloom base lemon sleep weekend rich forget diagram hurt prize fly"
VAL_MNEMONIC_2="turkey miss hurry unable embark hospital kangaroo nuclear outside term toy fall buffalo book opinion such moral meadow wing olive camp sad metal banner"
VAL_MNEMONIC_3="tenant neck ask season exist hill churn rice convince shock modify evidence armor track army street stay light program harvest now settle feed wheat"
VAL_MNEMONIC_4="tail forward era width glory magnet knock shiver cup broken turkey upgrade cigar story agent lake transfer misery sustain fragile parrot also air document"
VAL_MNEMONIC_5="crime lumber parrot enforce chimney turtle wing iron scissors jealous indicate peace empty game host protect juice submit motor cause second picture nuclear area"
VAL_MNEMONICS=(
  "$VAL_MNEMONIC_1"
  "$VAL_MNEMONIC_2"
  "$VAL_MNEMONIC_3"
  "$VAL_MNEMONIC_4"
  "$VAL_MNEMONIC_5"
)

USER_MNEMONIC="tonight bonus finish chaos orchard plastic view nurse salad regret pause awake link bacon process core talent whale million hope luggage sauce card weasel"

# Intento
INTO_CHAIN_ID=intento-dev-1
INTO_NODE_PREFIX=into
INTO_NUM_NODES=1 #must be greater or equal GAIA_NUM_NODES
INTO_VAL_PREFIX=val
INTO_USER_ACCT=usr1
INTO_USER_ADDRESS=into1wdplq6qjh2xruc7qqagma9ya665q6qhcpse4k6
INTO_ADDRESS_PREFIX=into
INTO_DENOM=$INTO_DENOM
INTO_RPC_PORT=26657
INTO_ADMIN_ACCT=admin
INTO_ADMIN_ADDRESS=into1u20df3trc2c2zdhm8qvh2hdjx9ewh00sqnuu4e
INTO_ADMIN_MNEMONIC="tone cause tribe this switch near host damage idle fragile antique tail soda alien depth write wool they rapid unfold body scan pledge soft"
INTO_RECEIVER_ADDRESS='into1g6qdx6kdhpf000afvvpte7hp0vnpzaputyyrem'

# Binaries are contigent on whether we're doing an upgrade or not
if [[ "${UPGRADE_NAME:-}" == "" ]]; then
  INTO_BINARY="$DOCKERNET_HOME/../build/intentod"
else
  if [[ "${NEW_BINARY:-false}" == "false" ]]; then
    INTO_BINARY="$UPGRADES/binaries/intentod1"
  else
    INTO_BINARY="$UPGRADES/binaries/intentod2"
  fi
fi
INTO_MAIN_CMD="$INTO_BINARY --home $DOCKERNET_HOME/state/${INTO_NODE_PREFIX}1"

# GAIA
GAIA_CHAIN_ID=GAIA
GAIA_NODE_PREFIX=gaia
GAIA_NUM_NODES=1
GAIA_BINARY="$DOCKERNET_HOME/../build/gaiad"
GAIA_VAL_PREFIX=gval
GAIA_USER_ACCT=gusr1
GAIA_USER_ADDRESS='cosmos1wdplq6qjh2xruc7qqagma9ya665q6qhcwju3ng'
GAIA_ADDRESS_PREFIX=cosmos
GAIA_DENOM=$ATOM_DENOM
GAIA_RPC_PORT=26557
GAIA_MAIN_CMD="$GAIA_BINARY --home $DOCKERNET_HOME/state/${GAIA_NODE_PREFIX}1"
GAIA_RECEIVER_ADDRESS='cosmos1g6qdx6kdhpf000afvvpte7hp0vnpzapuyxp8uf' 

# OSMO
OSMO_CHAIN_ID=OSMO #osmo-test-5
OSMO_NODE_PREFIX=osmo
OSMO_NUM_NODES=1
OSMO_BINARY="$DOCKERNET_HOME/../build/osmosisd"
OSMO_VAL_PREFIX=oval
OSMO_USER_ACCT=ousr1
OSMO_USER_ADDRESS='osmo1wdplq6qjh2xruc7qqagma9ya665q6qhcxf0p96'
OSMO_ADDRESS_PREFIX=osmo
OSMO_DENOM=$OSMO_DENOM
OSMO_RPC_PORT=26357
OSMO_MAIN_CMD="$OSMO_BINARY --home $DOCKERNET_HOME/state/${OSMO_NODE_PREFIX}1"
OSMO_RECEIVER_ADDRESS='osmo1g6qdx6kdhpf000afvvpte7hp0vnpzapuvajh2m'

# HOST (Intento running as a host zone)
HOST_CHAIN_ID=HOST
HOST_NODE_PREFIX=host
HOST_NUM_NODES=1
HOST_BINARY="$DOCKERNET_HOME/../build/intentod"
HOST_VAL_PREFIX=hval
HOST_ADDRESS_PREFIX=into
HOST_USER_ACCT=husr1
HOST_USER_ADDRESS='into1wdplq6qjh2xruc7qqagma9ya665q6qhc80zy8t'
HOST_DENOM=$COSM_DENOM
HOST_RPC_PORT=26157
HOST_MAIN_CMD="$HOST_BINARY --home $DOCKERNET_HOME/state/${HOST_NODE_PREFIX}1"
HOST_RECEIVER_ADDRESS='into1g6qdx6kdhpf000afvvpte7hp0vnpzaputyyrem'

# EVMOS
# EVMOS_CHAIN_ID=evmos_9001-2
# EVMOS_NODE_PREFIX=evmos
# EVMOS_NUM_NODES=1
# EVMOS_BINARY="$DOCKERNET_HOME/../build/evmosd"
# EVMOS_VAL_PREFIX=eval
# EVMOS_ADDRESS_PREFIX=evmos
# EVMOS_USER_ACCT=eusr1
# EVMOS_USER_ADDRESS='TODO'
# EVMOS_DENOM=$EVMOS_DENOM
# EVMOS_RPC_PORT=26057
# EVMOS_MAIN_CMD="$EVMOS_BINARY --home $DOCKERNET_HOME/state/${EVMOS_NODE_PREFIX}1"
# EVMOS_RECEIVER_ADDRESS='evmos123z469cfejeusvk87ufrs5520wmdxmmlc7qzuw'
# EVMOS_MICRO_DENOM_UNITS="000000000000000000000000"

# RELAYER
RELAYER_GAIA_EXEC="$DOCKER_COMPOSE run --rm relayer-gaia"
RELAYER_OSMO_EXEC="$DOCKER_COMPOSE run --rm relayer-osmo"
RELAYER_HOST_EXEC="$DOCKER_COMPOSE run --rm relayer-host"
HERMES_EXEC="$DOCKER_COMPOSE run --rm hermes"

RELAYER_INTO_ACCT=rly1
RELAYER_GAIA_ACCT=rly2
RELAYER_OSMO_ACCT=rly4
RELAYER_HOST_ACCT=rly6
RELAYER_ACCTS=(
  $RELAYER_GAIA_ACCT
  $RELAYER_OSMO_ACCT
  $RELAYER_HOST_ACCT
)

RELAYER_GAIA_MNEMONIC="fiction perfect rapid steel bundle giant blade grain eagle wing cannon fever must humble dance kitchen lazy episode museum faith off notable rate flavor"
RELAYER_OSMO_MNEMONIC="giraffe few grow task opinion ahead life marble again hurry age wave creek beef force picnic couple more extra exit tenant room embody monkey"
RELAYER_HOST_MNEMONIC="renew umbrella teach spoon have razor knee sock divert inner nut between immense library inhale dog truly return run remain dune virus diamond clinic"
# RELAYER_EVMOS_MNEMONIC="science depart where tell bus ski laptop follow child bronze rebel recall brief plug razor ship degree labor human series today embody fury harvest"
RELAYER_MNEMONICS=(
  "$RELAYER_GAIA_MNEMONIC"
  "$RELAYER_OSMO_MNEMONIC"
  "$RELAYER_HOST_MNEMONIC"
)

INTO_ADDRESS() {
  # After an upgrade, the keys query can sometimes print migration info,
  # so we need to filter by valid addresses using the prefix
  $INTO_MAIN_CMD keys show ${INTO_USER_ACCT} --keyring-backend test -a | grep $INTO_ADDRESS_PREFIX
}
GAIA_ADDRESS() {
  $GAIA_MAIN_CMD keys show ${GAIA_USER_ACCT} --keyring-backend test -a
}
OSMO_ADDRESS() {
  $OSMO_MAIN_CMD keys show ${OSMO_USER_ACCT} --keyring-backend test -a
}
HOST_ADDRESS() {
  $HOST_MAIN_CMD keys show ${HOST_USER_ACCT} --keyring-backend test -a
}

CSLEEP() {
  for i in $(seq $1); do
    sleep 1
    printf "\r\t$(($1 - $i))s left..."
  done
}

GET_VAR_VALUE() {
  var_name="$1"
  echo "${!var_name}"
}

WAIT_FOR_BLOCK() {
  num_blocks="${2:-1}"
  for i in $(seq $num_blocks); do
    (tail -f -n0 $1 &) | grep -q "executed block.*height="
  done
}

FLOW_ID="1"

GET_FLOW_ID() {
  address=$1
  max_blocks=${2:-10} # Default to 10 if not specified

  # Fetch initial flow IDs
  initial_flows=($($INTO_MAIN_CMD q intent list-flows-by-owner $address | awk -v RS='  ' '$1 == "id:" {print $2}'))

  for i in $(seq $max_blocks); do
    # Fetch new IDs
   new_flows=($($INTO_MAIN_CMD q intent list-flows-by-owner $address | awk -v RS='  ' '$1 == "id:" {print $2}'))
    # Find new ID by comparing initial and new lists
    for flow_id in "${new_flows[@]}"; do
      if [[ ! " ${initial_flows[*]} " =~ " ${flow_id} " ]]; then
        echo "New Flow detected with ID: $flow_id"
       #  return $((10#$id)) # Return the ID as an integer
       FLOW_ID=$flow_id
       return 0
      fi
    done

    # Wait for the next block
    WAIT_FOR_BLOCK $INTO_LOGS

    # Optional: Handle case where no new flows are found after max_blocks
    if [[ $i -eq $max_blocks ]]; then
      echo "No new Flows found after $max_blocks blocks."
      return -1 # Indicate no new flow was found
    fi
  done
}

WAIT_FOR_EXECUTED_FLOW_BY_ID() {
  max_blocks=${1:-100} # Default if not specified

  for i in $(seq $max_blocks); do
    # Fetch transaction info for the specified id
    echo "Querying with ID: $FLOW_ID"
    FLOW_ID=$(echo $FLOW_ID | xargs)
    history=$($INTO_MAIN_CMD q intent flow-history $FLOW_ID )

     # Check if all 'executed' keys are true for the specified tx_id
    executed_count=$(echo "$history" | grep 'executed:' | wc -l)
    executed_true_count=$(echo "$history" | grep 'executed: true' | wc -l)

   if [[ "$executed_count" -gt 0 ]] && [[ "$executed_count" -eq "$executed_true_count" ]]; then
      echo "All 'executed' instances for flow ID $FLOW_ID are true."
      sleep 5
      return 0
    fi


    # Wait for the next blocks
    WAIT_FOR_BLOCK $INTO_LOGS

    # Handle case where the transaction is not executed after max_blocks
    if [[ $i -eq $max_blocks ]]; then
      echo "Flow ID $FLOW_ID not executed after $max_blocks blocks."
      return -1
    fi
  done
}

WAIT_FOR_MSG_RESPONSES_LENGTH() {
  expected_length=$1
  max_blocks=${2:-100}

  echo "Waiting for msg_responses length to be $expected_length, max blocks: $max_blocks"

  for i in $(seq $max_blocks); do
    echo "Block attempt $i/$max_blocks"
    echo "Querying flow history for ID: $FLOW_ID"
    FLOW_ID=$(echo "$FLOW_ID" | xargs)

    history=$($INTO_MAIN_CMD q intent flow-history "$FLOW_ID" 2>/dev/null)
    if [[ -z "$history" ]]; then
      echo "No flow history found for flow ID: $FLOW_ID"
      WAIT_FOR_BLOCK "$INTO_LOGS"
      continue
    fi

    # Count number of '@type' lines in msg_responses
    responses_count=$(echo "$history" | grep -o '@type' | wc -l)
    echo "Current msg_responses count: $responses_count (target: $expected_length)"

    if [[ "$responses_count" -eq "$expected_length" ]]; then
      echo "msg_responses length reached expected $expected_length."
      sleep 5
      return 0
    fi

    WAIT_FOR_BLOCK "$INTO_LOGS"
  done

  echo "Timeout: msg_responses length did not reach $expected_length after $max_blocks blocks."
  return 1
}


WAIT_FOR_UPDATING_DISABLED() {
  max_blocks=${2:-10} # Default to 10 if not specified

  for i in $(seq $max_blocks); do
    # Fetch transaction info for the specified tx_id
    FLOW_ID=$(echo $FLOW_ID | xargs)
    tx_info=$($INTO_MAIN_CMD q intent flow $FLOW_ID)

    # Check if 'updating_disabled' is present in the transaction info
    disabled=$(echo "$tx_info" | grep 'updating_disabled' | awk '{print $2}')
    if [[ "$disabled" == "true" ]]; then
      echo "flow ID info $tx_info."
      echo "flow ID $FLOW_ID updating has been disabled."
      return 0
    fi

    # Wait for the next blocks
    WAIT_FOR_BLOCK $INTO_LOGS

    # Handle case where the transaction is not executed after max_blocks
    if [[ $i -eq $max_blocks ]]; then
      echo "flow ID $FLOW_ID not executed after $max_blocks blocks."
      return -1
    fi
  done
}

GET_VAL_ADDR() {
  chain=$1
  val_index=$2

  MAIN_CMD=$(GET_VAR_VALUE ${chain}_MAIN_CMD)
  $MAIN_CMD q staking validators 2>&1 | \
  grep -A 6 "${chain}_${val_index}" | \
  grep operator_address | \
  awk -F': ' '/operator_address/ {print $2}' || echo "Operator address not found"
}

GET_ICA_ADDR() {
  connection_id="$1"
  user_address="$2"

  # $INTO_MAIN_CMD q stakeibc show-host-zone $chain_id | grep ${ica_type}_account -A 1 | grep address | awk '{print $2}'
  $INTO_MAIN_CMD query intent interchainaccounts $connection_id $user_address json | grep interchain_account_address | awk '{print $2}'
}

TRIM_TX() {
  grep -E "code:|txhash:" | sed 's/^/  /'
}

VERIFY_SUCCESS() {
  TX_RES="$1"
  
  if [[ $(echo "$TX_RES" | grep -oE "code: [0-9]+" | awk '{print $2}') -ne 0 ]]; then
    echo "Error: Failed to submit create-consumer transaction."
    exit 1
  fi
  
  echo "Transaction submitted successfully:"
  echo "$TX_RES" | grep -E "code:|txhash:" | sed 's/^/  /'
}
