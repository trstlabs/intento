#!/bin/bash

set -eu
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

source $SCRIPT_DIR/../config.sh

CHAIN="$1"
KEYS_LOGS=$DOCKERNET_HOME/logs/keys.log

CHAIN_ID=$(GET_VAR_VALUE    ${CHAIN}_CHAIN_ID)
BINARY=$(GET_VAR_VALUE      ${CHAIN}_BINARY)
MAIN_CMD=$(GET_VAR_VALUE    ${CHAIN}_MAIN_CMD)
DENOM=$(GET_VAR_VALUE       ${CHAIN}_DENOM)
RPC_PORT=$(GET_VAR_VALUE    ${CHAIN}_RPC_PORT)
NUM_NODES=$(GET_VAR_VALUE   ${CHAIN}_NUM_NODES)
NODE_PREFIX=$(GET_VAR_VALUE ${CHAIN}_NODE_PREFIX)
VAL_PREFIX=$(GET_VAR_VALUE  ${CHAIN}_VAL_PREFIX)

# THe host zone can optionally specify additional the micro-denom granularity
# If they don't specify the ${CHAIN}_MICRO_DENOM_UNITS variable,
# EXTRA_MICRO_DENOM_UNITS will include 6 0's
MICRO_DENOM_UNITS_VAR_NAME=${CHAIN}_MICRO_DENOM_UNITS
MICRO_DENOM_UNITS="${!MICRO_DENOM_UNITS_VAR_NAME:-000000}"

GENESIS_TOKENS=${GENESIS_TOKENS}${MICRO_DENOM_UNITS}
STAKE_TOKENS=${STAKE_TOKENS}${MICRO_DENOM_UNITS}
ADMIN_TOKENS=${ADMIN_TOKENS}${MICRO_DENOM_UNITS}
FAUCET_TOKENS=${FAUCET_TOKENS}${MICRO_DENOM_UNITS}

set_into_genesis() {
    genesis_config=$1
    
    # update params
    jq '.app_state.claim.claim_records[0].address = "into1wdplq6qjh2xruc7qqagma9ya665q6qhcpse4k6"' $genesis_config > json.tmp && mv json.tmp $genesis_config
    jq '.app_state.claim.claim_records[0].maximum_claimable_amount = {"amount":"10000","denom":"uinto"}' $genesis_config > json.tmp && mv json.tmp $genesis_config
    jq '.app_state.claim.claim_records[0].status[0].action_completed = false' $genesis_config > json.tmp && mv json.tmp $genesis_config
    jq '.app_state.claim.claim_records[0].status[0].vesting_periods_completed = [false,false,false,false]' $genesis_config > json.tmp && mv json.tmp $genesis_config
    jq '.app_state.claim.claim_records[0].status[0].vesting_periods_claimed = [false,false,false,false]' $genesis_config > json.tmp && mv json.tmp $genesis_config
    
    jq '.app_state.claim.params.duration_vesting_periods = ["40s","50s","60s","70s"]' $genesis_config > json.tmp && mv json.tmp $genesis_config
    
    jq '.app_state.staking.params.unbonding_time = $newVal' --arg newVal "$UNBONDING_TIME" $genesis_config > json.tmp && mv json.tmp $genesis_config
    jq '.app_state.gov.params.max_deposit_period = $newVal' --arg newVal "$MAX_DEPOSIT_PERIOD" $genesis_config > json.tmp && mv json.tmp $genesis_config
    jq '.app_state.gov.params.voting_period = $newVal' --arg newVal "$VOTING_PERIOD" $genesis_config > json.tmp && mv json.tmp $genesis_config
    
    # enable intento as an interchain accounts controller
    jq "del(.app_state.interchain_accounts)" $genesis_config > json.tmp && mv json.tmp $genesis_config
    interchain_accts=$(cat $DOCKERNET_HOME/config/ica_controller.json)
    jq ".app_state += $interchain_accts" $genesis_config > json.tmp && mv json.tmp $genesis_config
    
}



#MIN_ATOM_AMOUNT="[{"denom":"uatom","amount":"0"}]"
set_host_genesis() {
    genesis_config=$1
    #epoch
    HOST_DAY_EPOCH_DURATION="600s"
    HOST_HOUR_EPOCH_DURATION="600s"
    HOST_WEEK_EPOCH_DURATION="600s"
    HOST_MINT_EPOCH_DURATION="600s"
    # Shorten epochs and unbonding time
    jq '(.app_state.epochs.epochs[]? | select(.identifier=="day") ).duration = $epochLen' --arg epochLen $HOST_DAY_EPOCH_DURATION $genesis_config > json.tmp && mv json.tmp $genesis_config
    jq '(.app_state.epochs.epochs[]? | select(.identifier=="hour") ).duration = $epochLen' --arg epochLen $HOST_HOUR_EPOCH_DURATION $genesis_config > json.tmp && mv json.tmp $genesis_config
    jq '(.app_state.epochs.epochs[]? | select(.identifier=="week") ).duration = $epochLen' --arg epochLen $HOST_WEEK_EPOCH_DURATION $genesis_config > json.tmp && mv json.tmp $genesis_config
    jq '(.app_state.epochs.epochs[]? | select(.identifier=="mint") ).duration = $epochLen' --arg epochLen $HOST_MINT_EPOCH_DURATION $genesis_config > json.tmp && mv json.tmp $genesis_config
    
    jq '.app_state.staking.params.unbonding_time = $newVal' --arg newVal "$UNBONDING_TIME" $genesis_config > json.tmp && mv json.tmp $genesis_config
    
    # # Shorten voting period (we need both of these to support both versions of the SDK)
    # if [[ "$(jq '.app_state.gov | has("voting_params")' $genesis_config)" == "true" ]]; then
    #     jq '.app_state.gov.voting_params.voting_period = $newVal' --arg newVal "$VOTING_PERIOD" $genesis_config > json.tmp && mv json.tmp $genesis_config
    # fi
    # if [[ "$(jq '.app_state.gov | has("params")' $genesis_config)" == "true" ]]; then
    
    # fi
    
    # Set the mint start time to the genesis time if the chain configures inflation at the block level (e.g. stars)
    # also reduce the number of initial annual provisions so the inflation rate is not too high
    genesis_time=$(jq .genesis_time $genesis_config | tr -d '"')
    jq 'if .app_state.mint.params.start_time? then .app_state.mint.params.start_time=$newVal else . end' --arg newVal "$genesis_time" $genesis_config > json.tmp && mv json.tmp $genesis_config
    jq 'if .app_state.mint.params.initial_annual_provisions? then .app_state.mint.params.initial_annual_provisions=$newVal else . end' --arg newVal "$INITIAL_ANNUAL_PROVISIONS" $genesis_config > json.tmp && mv json.tmp $genesis_config
    
    # Add interchain accounts to the genesis set
    jq "del(.app_state.interchain_accounts)" $genesis_config > json.tmp && mv json.tmp $genesis_config
    interchain_accts=$(cat $DOCKERNET_HOME/config/ica_host.json)
    jq ".app_state += $interchain_accts" $genesis_config > json.tmp && mv json.tmp $genesis_config
    
    # Slightly harshen slashing parameters (if 5 blocks are missed, the validator will be slashed)
    # This makes it easier to test updating weights after a host zone validator is slashed
    sed -i -E 's|"signed_blocks_window": "100"|"signed_blocks_window": "10"|g' $genesis_config
    sed -i -E 's|"downtime_jail_duration": "600s"|"downtime_jail_duration": "10s"|g' $genesis_config
    sed -i -E 's|"slash_fraction_downtime": "0.010000000000000000"|"slash_fraction_downtime": "0.050000000000000000"|g' $genesis_config
    
    # LSM params
    if [[ "$CHAIN" == "GAIA" ]]; then
        # LSM Params
        LSM_VALIDATOR_BOND_FACTOR="250"
        LSM_GLOBAL_LIQUID_STAKING_CAP="0.25"
        LSM_VALIDATOR_LIQUID_STAKING_CAP="0.50"
        UPGRADE_TIMEOUT="2526374086000000000"
        jq '.app_state.gov.params.voting_period = $newVal' --arg newVal "$VOTING_PERIOD" $genesis_config > json.tmp && mv json.tmp $genesis_config
        jq '.app_state.staking.params.min_commission_rate = $newVal' --arg newVal "0.010000000000000000" $genesis_config > json.tmp && mv json.tmp $genesis_config
        jq '.app_state.staking.params.validator_bond_factor = $newVal' --arg newVal "$LSM_VALIDATOR_BOND_FACTOR" $genesis_config > json.tmp && mv json.tmp $genesis_config
        jq '.app_state.staking.params.validator_liquid_staking_cap = $newVal' --arg newVal "$LSM_VALIDATOR_LIQUID_STAKING_CAP" $genesis_config > json.tmp && mv json.tmp $genesis_config
        jq '.app_state.staking.params.global_liquid_staking_cap = $newVal' --arg newVal "$LSM_GLOBAL_LIQUID_STAKING_CAP" $genesis_config > json.tmp && mv json.tmp $genesis_config
        jq '.app_state.ibc.channel_genesis.params.upgrade_timeout.timestamp = $newVal' --arg newVal "$UPGRADE_TIMEOUT" $genesis_config > json.tmp && mv json.tmp $genesis_config
        jq '.app_state.provider.params.max_provider_consensus_validators = $newVal' --arg newVal "180" $genesis_config > json.tmp && mv json.tmp $genesis_config
        jq '.app_state.provider.params.blocks_per_epoch = $newVal' --arg newVal "3" $genesis_config > json.tmp && mv json.tmp $genesis_config
        
        # jq '.app_state.feemarket.params.min_base_gas_price = $newVal' --arg newVal "0.0025" $genesis_config > json.tmp && mv json.tmp $genesis_config
        # jq '.app_state.feemarket.params.fee_denom = $newVal' --arg newVal "uatom" $genesis_config > json.tmp && mv json.tmp $genesis_config
        # jq '.app_state.feemarket.params.enabled = false' $genesis_config > json.tmp && mv json.tmp $genesis_config
        #  jq '.app_state.feemarket.params.window = 1' $genesis_config > json.tmp && mv json.tmp $genesis_config
        jq '.app_state.feemarket.state.base_gas_price = $newVal' --arg newVal "0.002" $genesis_config > json.tmp && mv json.tmp $genesis_config
        jq '.app_state.feemarket.state.window = ["1"]' $genesis_config > json.tmp && mv json.tmp $genesis_config
        jq '.app_state.feemarket.state.learning_rate = "0.1"' $genesis_config > json.tmp && mv json.tmp $genesis_config
        jq '.app_state.feemarket.params = {
    alpha: "0.1",
    beta: "0.1",
    gamma: "0.1",
    delta: "0.1",
    min_base_gas_price: "0.002",
    min_learning_rate: "0.1",
    max_learning_rate: "0.1",
    max_block_utilization: "10",
    window: "1",
    enabled: false,
    fee_denom: "uatom"
        }' $genesis_config > json.tmp && mv json.tmp $genesis_config
        
        jq "del(.app_state.provider.mature_unbonding_ops)" $genesis_config > genesis.tmp && mv genesis.tmp $genesis_config
        jq "del(.app_state.provider.unbonding_ops)" $genesis_config > genesis.tmp && mv genesis.tmp $genesis_config
        jq "del(.app_state.provider.consumer_addition_proposals)" $genesis_config > genesis.tmp && mv genesis.tmp $genesis_config
        jq "del(.app_state.provider.consumer_removal_proposals)" $genesis_config > genesis.tmp && mv genesis.tmp $genesis_config
        jq "del(.app_state.provider.consumer_addrs_to_prune)" $genesis_config > genesis.tmp && mv genesis.tmp $genesis_config
        
        jq "del(.app_state.provider.params.max_throttled_packets)" $genesis_config > genesis.tmp && mv genesis.tmp $genesis_config
        jq "del(.app_state.provider.params.init_timeout_period)" $genesis_config > genesis.tmp && mv genesis.tmp $genesis_config
        jq "del(.app_state.provider.params.vsc_timeout_period)" $genesis_config > genesis.tmp && mv genesis.tmp $genesis_config
    fi
    
    
    # Add poolmanager params
    if [[ "$CHAIN" == "OSMO" ]]; then
        jq "del(.app_state.poolincentives.pool_to_gauges)" $genesis_config > genesis.tmp && mv genesis.tmp $genesis_config
        jq "del(.app_state.bank.supply_offsets)" $genesis_config > genesis.tmp && mv genesis.tmp $genesis_config
        jq "del(.app_state.staking.params.min_self_delegation)" $genesis_config > genesis.tmp && mv genesis.tmp $genesis_config
        jq "del(.app_state.gov.deposit_params.min_expedited_deposit)" $genesis_config > genesis.tmp && mv genesis.tmp $genesis_config
        jq "del(.app_state.gov.deposit_params.min_initial_deposit_ratio)" $genesis_config > genesis.tmp && mv genesis.tmp $genesis_config
        jq "del(.app_state.gov.voting_params.proposal_voting_periods)" $genesis_config > genesis.tmp && mv genesis.tmp $genesis_config
        jq "del(.app_state.gov.voting_params.expedited_voting_period)" $genesis_config > genesis.tmp && mv genesis.tmp $genesis_config
        jq "del(.app_state.gov.tally_params.expedited_threshold)" $genesis_config > genesis.tmp && mv genesis.tmp $genesis_config
        jq "del(.app_state.gov.tally_params.expedited_quorum)" $genesis_config > genesis.tmp && mv genesis.tmp $genesis_config
        jq "del(.app_state.wasm.params.code_upload_access.address)" $genesis_config > genesis.tmp && mv genesis.tmp $genesis_config
        
        #jq '.app_state.concentratedliquidity.authorized_uptimes = ["60s"]'  $genesis_config > genesis.tmp && mv genesis.tmp $genesis_config
        jq '.app_state.incentives.params.internal_uptime = "1m0s"'  $genesis_config > genesis.tmp && mv genesis.tmp $genesis_config
        # Set default governance params if not present
        jq '.app_state.gov.params = {
    "min_deposit": [{"denom": "stake", "amount": "10000000"}],
    "max_deposit_period": "172800s",
    "voting_period": "432000s",
    "quorum": "0.334",
    "threshold": "0.5",
    "veto_threshold": "0.334",
    "min_initial_deposit_ratio": "0.2",
    "expedited_voting_period": "86400s",
    "expedited_threshold": "0.67",
    "expedited_min_deposit": [{"denom": "stake", "amount": "100000000"}],
    "burn_vote_quorum": false,
    "burn_proposal_deposit_prevote": false,
    "burn_vote_veto": false
        }'  $genesis_config > genesis.tmp && mv genesis.tmp $genesis_config
        
        jq '.app_state.poolmanager = {
    "next_pool_id": 1,
    "params": {
      "pool_creation_fee": [{"denom": "uosmo", "amount": "1000000000"}],
      "taker_fee_params": {
        "default_taker_fee": "0.000000000000000000",
        "osmo_taker_fee_distribution": {
          "staking_rewards": "1.000000000000000000",
          "community_pool": "0.000000000000000000"
        },
        "non_osmo_taker_fee_distribution": {
          "staking_rewards": "0.670000000000000000",
          "community_pool": "0.330000000000000000"
        },
        "admin_addresses": [],
        "community_pool_denom_to_swap_non_whitelisted_assets_to": "ibc/D189335C6E4A68B513C10AB227BF1C1D38C746766278BA3EEB4FB14124F1D858",
        "reduced_fee_whitelist": []
      },
      "authorized_quote_denoms": [
        "uosmo",
        "ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2",
        "ibc/0CD3A0285E1341859B5E86B6AB7682F023D03E97607CCC1DC95706411D866DF7",
        "ibc/D189335C6E4A68B513C10AB227BF1C1D38C746766278BA3EEB4FB14124F1D858"
      ]
        },
    "pool_routes": [],
    "taker_fees_tracker": {
        "taker_fees_to_stakers": [
            {"denom": "uosmo", "amount": "0"}
        ],
        "taker_fees_to_community_pool": [
            {"denom": "uosmo", "amount": "0"}
        ],
        "height_accounting_starts_from": 0
    },
    "pool_volumes": [],
    "denom_pair_taker_fee_store": []
        }'  $genesis_config > genesis.tmp && mv genesis.tmp $genesis_config
        
        jq '.app_state.protorev = {
    "params": {
        "enabled": true,
        "admin": "osmo1wdplq6qjh2xruc7qqagma9ya665q6qhcxf0p96"
    },
    "token_pair_arb_routes": [],
    "base_denoms": [
        {"denom": "uosmo", "step_size":"1"}
    ],
    "days_since_module_genesis": 0,
    "developer_fees": [
        {"denom": "uosmo", "amount": "0"}
    ],
    "latest_block_height": 0,
    "developer_address": "osmo1wdplq6qjh2xruc7qqagma9ya665q6qhcxf0p96",
    "max_pool_points_per_block": 200,
    "max_pool_points_per_tx": 50,
    "point_count_for_block": 0,
    "profits": [
        {"denom": "uosmo", "amount": "0"}
    ],
    "info_by_pool_type": {
        "stable": {
            "weight": 100
        },
        "balancer": {
            "weight": 100
        },
        "concentrated": {
            "weight": 100,
            "max_ticks_crossed": 6
        },
        "cosmwasm": {
            "weight_maps": [
                {
                    "weight": 100,
                    "contract_address": "osmo1wdplq6qjh2xruc7qqagma9ya665q6qhcxf0p96"
                }
            ]
        }
    },
    "cyclic_arb_tracker": {
        "cyclic_arb": [
            {
                "denom": "uosmo",
                "amount": "0"
            }
        ],
        "height_accounting_starts_from": 0
    }
        }' $genesis_config > genesis.tmp && mv genesis.tmp $genesis_config
        
        
        jq '.app_state.incentives.lockable_durations = ["1s", "120s", "180s","240s"]'  $genesis_config > genesis.tmp && mv genesis.tmp $genesis_config
        jq '.app_state.wasm.params.code_upload_access.permission = "Everybody"'  $genesis_config > genesis.tmp && mv genesis.tmp $genesis_config
    fi
}


# set_consumer_genesis() {
#     genesis_config=$1

#     # add consumer genesis
#     home_directories=""
#     for (( i=1; i <= $NUM_NODES; i++ )); do
#         home_directories+="${STATE}/${NODE_PREFIX}${i},"
#     done

#     $MAIN_CMD add-consumer-section --validator-home-directories $home_directories
#     jq '.app_state.ccvconsumer.params.unbonding_period = $newVal' --arg newVal "$UNBONDING_TIME" $genesis_config > json.tmp && mv json.tmp $genesis_config
# }


MAIN_ID=1 # Node responsible for genesis and persistent_peers
MAIN_NODE_NAME=""
MAIN_NODE_ID=""
MAIN_CONFIG=""
MAIN_GENESIS=""
echo "Initializing $CHAIN chain..."
for (( i=1; i <= $NUM_NODES; i++ )); do
    # Node names will be of the form: "trst1"
    node_name="${NODE_PREFIX}${i}"
    # Moniker is of the form: INTO_1
    moniker=$(printf "${NODE_PREFIX}_${i}" | awk '{ print toupper($0) }')
    
    # Create a state directory for the current node and initialize the chain
    mkdir -p $STATE/$node_name
    
    # If the chains commands are run only from docker, grab the command from the config
    # Otherwise, if they're run locally, append the home directory
    if [[ $BINARY == docker-compose* ]]; then
        cmd=$BINARY
    else
        cmd="$BINARY --home ${STATE}/$node_name"
    fi
    
    # Initialize the chain
    $cmd init $moniker --chain-id $CHAIN_ID --overwrite &> /dev/null
    chmod -R 777 $STATE/$node_name
    
    if [[ "$CHAIN" == "INTO"  ]]; then
        TEST_FILES_DIR=$DOCKERNET_HOME/tests/test_files
        $cmd prepare-genesis testnet $CHAIN_ID
        if [ -d "$TEST_FILES_DIR" ]; then
            $cmd export-snapshot $TEST_FILES_DIR/active-users.csv $TEST_FILES_DIR/nft1.csv $TEST_FILES_DIR/nft2.csv $TEST_FILES_DIR/snapshot_output.json --nft-weight-1 20 --nft-weight 10 --user-weight 5
            $cmd import-genesis-accounts-from-snapshot $TEST_FILES_DIR/snapshot_output.json $TEST_FILES_DIR/non-airdrop-accounts.json --airdrop-amount=90_000_000_000_000
        fi
    fi
    # Update node networking configuration
    config_toml="${STATE}/${node_name}/config/config.toml"
    client_toml="${STATE}/${node_name}/config/client.toml"
    app_toml="${STATE}/${node_name}/config/app.toml"
    genesis_json="${STATE}/${node_name}/config/genesis.json"
    
    
    # sed -i -E "s|cors_allowed_origins = \[\]|cors_allowed_origins = [\"\*\"]|g" $config_toml
    sed -i -E "s|127.0.0.1|0.0.0.0|g" $config_toml
    sed -i -E "s|timeout_commit = \"5s\"|timeout_commit = \"${BLOCK_TIME}\"|g" $config_toml
    sed -i -E "s|timeout_commit = \"2s\"|timeout_commit = \"${BLOCK_TIME}\"|g" $config_toml
    sed -i -E "s|timeout_commit = \"500ms\"|timeout_commit = \"${BLOCK_TIME}\"|g" $config_toml
    sed -i -E "s|timeout_propose = \"3s\"|timeout_propose = \"${BLOCK_TIME}\"|g" $config_toml
    sed -i -E 's|timeout_propose = "3s"|timeout_propose = "1s"|g' $config_toml
    sed -i -E 's|timeout_propose = "1.8s"|timeout_propose = "1s"|g' $config_toml
    
    sed -i -E "s|timeout_propose = \"2s\"|timeout_propose = \"1s\"|g" $config_toml
    sed -i -E "s|prometheus = false|prometheus = true|g" $config_toml
    
    sed -i -E "s|minimum-gas-prices = \".*\"|minimum-gas-prices = \"0${DENOM}\"|g" $app_toml
    sed -i -E '/\[api\]/,/^enable = .*$/ s/^enable = .*$/enable = true/' $app_toml
    sed -i -E 's|unsafe-cors = .*|unsafe-cors = true|g' $app_toml
    sed -i -E 's|swagger = .*|swagger = true|g' $app_toml
    sed -i -E "s|snapshot-interval = 0|snapshot-interval = 300|g" $app_toml
    sed -i -E 's|localhost|0.0.0.0|g' $app_toml
    
    sed -i -E "s|chain-id = \"\"|chain-id = \"${CHAIN_ID}\"|g" $client_toml
    sed -i -E "s|keyring-backend = \"os\"|keyring-backend = \"test\"|g" $client_toml
    sed -i -E "s|node = \".*\"|node = \"tcp://localhost:$RPC_PORT\"|g" $client_toml
    
    sed -i -E "s|\"stake\"|\"${DENOM}\"|g" $genesis_json
    sed -i -E "s|\"aphoton\"|\"${DENOM}\"|g" $genesis_json # ethermint default
    
    # Get the endpoint and node ID
    node_id=$($cmd tendermint show-node-id)@$node_name:$PEER_PORT
    echo "Node #$i ID: $node_id"
    
    # add a validator account
    val_acct="${VAL_PREFIX}${i}"
    val_mnemonic="${VAL_MNEMONICS[((i-1))]}"
    echo "$val_mnemonic" | $cmd keys add $val_acct --recover --keyring-backend=test >> $KEYS_LOGS 2>&1
    val_addr=$($cmd keys show $val_acct --keyring-backend test -a | tr -cd '[:alnum:]._-')
    # Add this account to the current node
    if [ "$CHAIN" == "GAIA" ]; then
        $cmd genesis add-genesis-account ${val_addr} ${GENESIS_TOKENS}${DENOM}
        $cmd genesis gentx $val_acct ${STAKE_TOKENS}${DENOM} --chain-id $CHAIN_ID --keyring-backend test &> /dev/null
    else
        $cmd add-genesis-account ${val_addr} ${GENESIS_TOKENS}${DENOM}
        $cmd gentx $val_acct ${STAKE_TOKENS}${DENOM} --chain-id $CHAIN_ID --keyring-backend test &> /dev/null
    fi
    # actually set this account as a validator on the current node
    
    
    # Cleanup from seds
    rm -rf ${client_toml}-E
    rm -rf ${genesis_json}-E
    rm -rf ${app_toml}-E
    
    if [ $i -eq $MAIN_ID ]; then
        MAIN_NODE_NAME=$node_name
        MAIN_NODE_ID=$node_id
        MAIN_CONFIG=$config_toml
        MAIN_GENESIS=$genesis_json
    else
        # also add this account and it's genesis tx to the main node
        $MAIN_CMD add-genesis-account ${val_addr} ${GENESIS_TOKENS}${DENOM}
        cp ${STATE}/${node_name}/config/gentx/*.json ${STATE}/${MAIN_NODE_NAME}/config/gentx/
        
        # and add each validator's keys to the first state directory
        echo "$val_mnemonic" | $MAIN_CMD keys add $val_acct --recover --keyring-backend=test &> /dev/null
    fi
done

if [ "$CHAIN" == "INTO" ]; then
    # Add the into admin account
    echo "$INTO_ADMIN_MNEMONIC" | $MAIN_CMD keys add $INTO_ADMIN_ACCT --recover --keyring-backend=test >> $KEYS_LOGS 2>&1
    INTO_ADMIN_ADDRESS=$($MAIN_CMD keys show $INTO_ADMIN_ACCT --keyring-backend test -a)
    $MAIN_CMD add-genesis-account ${INTO_ADMIN_ADDRESS} ${ADMIN_TOKENS}${DENOM}
    
    echo "$TEST_FAUCET_MNEMONIC" | $MAIN_CMD keys add $TEST_FAUCET_ACCT --recover --keyring-backend=test >> $KEYS_LOGS 2>&1
    INTO_TEST_FAUCET_ADDRESS=$($MAIN_CMD keys show $TEST_FAUCET_ACCT --keyring-backend test -a)
    echo "FAUCET Address: " $INTO_TEST_FAUCET_ADDRESS
    $MAIN_CMD add-genesis-account ${INTO_TEST_FAUCET_ADDRESS} ${FAUCET_TOKENS}${DENOM}
    
    # Add a user account
    USER_ACCT_VAR=${CHAIN}_USER_ACCT
    USER_ACCT=${!USER_ACCT_VAR}
    echo $USER_MNEMONIC | $MAIN_CMD keys add $USER_ACCT --recover --keyring-backend=test >> $KEYS_LOGS 2>&1
    USER_ADDRESS=$($MAIN_CMD keys show $USER_ACCT --keyring-backend test -a | tr -cd '[:alnum:]._-')
    echo "USER Address: " $USER_ADDRESS
    $MAIN_CMD add-genesis-account ${USER_ADDRESS} ${GENESIS_TOKENS}${DENOM}
    
    # Add relayer accounts
    for i in "${!RELAYER_ACCTS[@]}"; do
        RELAYER_ACCT="${RELAYER_ACCTS[i]}"
        RELAYER_MNEMONIC="${RELAYER_MNEMONICS[i]}"
        
        echo "$RELAYER_MNEMONIC" | $MAIN_CMD keys add $RELAYER_ACCT --recover --keyring-backend=test >> $KEYS_LOGS 2>&1
        RELAYER_ADDRESS=$($MAIN_CMD keys show $RELAYER_ACCT --keyring-backend test -a)
        $MAIN_CMD add-genesis-account ${RELAYER_ADDRESS} ${GENESIS_TOKENS}${DENOM}
    done
    
    
else
    # Add a user account
    USER_ACCT_VAR=${CHAIN}_USER_ACCT
    USER_ACCT=${!USER_ACCT_VAR}
    echo $USER_MNEMONIC | $MAIN_CMD keys add $USER_ACCT --recover --keyring-backend=test >> $KEYS_LOGS 2>&1
    USER_ADDRESS=$($MAIN_CMD keys show $USER_ACCT --keyring-backend test -a | tr -cd '[:alnum:]._-')
    echo "USER Address: " $USER_ADDRESS
    if [ "$CHAIN" == "GAIA" ]; then
        $MAIN_CMD genesis add-genesis-account ${USER_ADDRESS} ${GENESIS_TOKENS}${DENOM}
    else
        $MAIN_CMD add-genesis-account ${USER_ADDRESS} ${GENESIS_TOKENS}${DENOM}
    fi
    
    echo "$TEST_FAUCET_MNEMONIC" | $MAIN_CMD keys add $TEST_FAUCET_ACCT --recover --keyring-backend=test >> $KEYS_LOGS 2>&1
    TEST_FAUCET_ADDRESS=$($MAIN_CMD keys show $TEST_FAUCET_ACCT --keyring-backend test -a)
    echo "FAUCET Address: " $TEST_FAUCET_ADDRESS
    
    if [ "$CHAIN" == "GAIA" ]; then
        $MAIN_CMD genesis add-genesis-account ${TEST_FAUCET_ADDRESS} ${FAUCET_TOKENS}${DENOM}
    else
        $MAIN_CMD add-genesis-account ${TEST_FAUCET_ADDRESS} ${FAUCET_TOKENS}${DENOM}
    fi
    # Add a relayer account
    RELAYER_ACCT=$(GET_VAR_VALUE RELAYER_${CHAIN}_ACCT)
    RELAYER_MNEMONIC=$(GET_VAR_VALUE RELAYER_${CHAIN}_MNEMONIC)
    
    echo "$RELAYER_MNEMONIC" | $MAIN_CMD keys add $RELAYER_ACCT --recover --keyring-backend=test >> $KEYS_LOGS 2>&1
    RELAYER_ADDRESS=$($MAIN_CMD keys show $RELAYER_ACCT --keyring-backend test -a | tr -cd '[:alnum:]._-')
    
    if [ "$CHAIN" == "GAIA" ]; then
        $MAIN_CMD genesis add-genesis-account ${RELAYER_ADDRESS} ${GENESIS_TOKENS}${DENOM}
    else
        $MAIN_CMD add-genesis-account ${RELAYER_ADDRESS} ${GENESIS_TOKENS}${DENOM}
    fi
fi



if [ "$CHAIN" == "GAIA" ]; then
    $MAIN_CMD genesis collect-gentxs &> /dev/null
else
    # now we process gentx txs on the main node
    $MAIN_CMD collect-gentxs &> /dev/null
fi
# wipe out the seeds and persistent peers for the main node (these are incorrectly autogenerated for each validator during collect-gentxs)
sed -i -E "s|persistent_peers = .*|persistent_peers = \"\"|g" $MAIN_CONFIG
sed -i -E "s|seeds = .*|seeds = \"\"|g" $MAIN_CONFIG

# update chain-specific settings
if [ "$CHAIN" == "INTO" ]; then
    #sed -i -E "s|log_level = \"info\"|log_level = \"debug\"|g" $MAIN_CONFIG
    sed -i -E "s|timeout_commit = \"5s\"|timeout_commit = \"500ms\"|g" $MAIN_CONFIG
    sed -i -E "s|timeout_propose = \"3s\"|timeout_propose = \"1s\"|g" $MAIN_CONFIG
    set_into_genesis $MAIN_GENESIS
else
    #sed -i -E "s|log_level = \"info\"|log_level = \"debug\"|g" $MAIN_CONFIG
    set_host_genesis $MAIN_GENESIS
fi


# update consumer genesis for binary chains
# if [[ "$CHAIN" == "INTO" || "$CHAIN" == "HOST" ]]; then
#     set_consumer_genesis $MAIN_GENESIS
# fi


# for all peer nodes....
for (( i=2; i <= $NUM_NODES; i++ )); do
    node_name="${NODE_PREFIX}${i}"
    config_toml="${STATE}/${node_name}/config/config.toml"
    genesis_json="${STATE}/${node_name}/config/genesis.json"
    
    # add the main node as a persistent peer
    sed -i -E "s|persistent_peers = .*|persistent_peers = \"${MAIN_NODE_ID}\"|g" $config_toml
    # copy the main node's genesis to the peer nodes to ensure they all have the same genesis
    cp $MAIN_GENESIS $genesis_json
    
    rm -rf ${config_toml}-E
done

# Cleanup from seds
rm -rf ${MAIN_CONFIG}-E
rm -rf ${MAIN_GENESIS}-E