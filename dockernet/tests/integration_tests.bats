#!/usr/bin/env bats

load "bats/bats-support/load.bash"
load "bats/bats-assert/load.bash"

setup_file() {
  _example_run_command="CHAIN_NAME=GAIA TRANSFER_CHANNEL_NUMBER=0 bats gaia_tests.bats"
  if [[ "$CHAIN_NAME" == "" ]]; then
    echo "CHAIN_NAME variable must be set before running integration tests (e.g. $_example_run_command)" >&2
    return 1
  fi

  if [[ "$TRANSFER_CHANNEL_NUMBER" == "" ]]; then
    echo "TRANSFER_CHANNEL_NUMBER and CONNECTION_ID variable must be set before running integration tests (e.g. $_example_run_command)" >&2
    return 1
  fi

  if [[ "$CONNECTION_ID" == "" ]]; then
    echo "CONNECTION_ID variable must be set before running integration tests (e.g. $_example_run_command)" >&2
    return 1
  fi

  # set allows us to export all variables in account_vars
  set -a
  NEW_BINARY="${NEW_BINARY:-false}" source dockernet/config.sh

  HOST_CHAIN_ID=$(GET_VAR_VALUE ${CHAIN_NAME}_CHAIN_ID)
  HOST_DENOM=$(GET_VAR_VALUE ${CHAIN_NAME}_DENOM)
  HOST_IBC_DENOM=$(GET_VAR_VALUE IBC_${CHAIN_NAME}_CHANNEL_${TRANSFER_CHANNEL_NUMBER}_DENOM)
  HOST_MAIN_CMD=$(GET_VAR_VALUE ${CHAIN_NAME}_MAIN_CMD)
  HOST_USER_ADDRESS=$(GET_VAR_VALUE ${CHAIN_NAME}_USER_ADDRESS)
  HOST_USER=$(GET_VAR_VALUE ${CHAIN_NAME}_USER_ACCT)
  HOST_RECEIVER_ADDRESS=$(GET_VAR_VALUE ${CHAIN_NAME}_RECEIVER_ADDRESS) # random address with nil balance

  INTO_USER=$(GET_VAR_VALUE INTO_USER_ACCT)
  INTO_VAL=${INTO_VAL_PREFIX}1
  INTO_TRANFER_CHANNEL="channel-${TRANSFER_CHANNEL_NUMBER}"

  TRANSFER_AMOUNT=5000
  MSGSEND_AMOUNT=1000
  MSGDELEGATE_AMOUNT=100000000
  EXPECTED_FEE=2500
  HOST_MAIN_CMD_TX="$HOST_MAIN_CMD tx --fees $EXPECTED_FEE$HOST_DENOM"
  RECURRING_MSGSEND_AMOUNT_TOTAL=1728000 #100000*120*24*60
  ICS20_AMOUNT_FOR_LOCAL_GAS=5000

  GETBAL() {
    grep -oP '(?<=amount: ")[0-9]+' || echo "0"
  }
  GETSTAKE() {
    tail -n 2 | head -n 1 | grep -o -E '[0-9]+' | head -n 1
  }
  # HELPER FUNCTIONS
  DECADD() {
    echo "scale=2; $1+$2" | bc
  }
  DECMUL() {
    echo "scale=2; $1*$2" | bc
  }
  FLOOR() {
    printf "%.0f\n" $1
  }
  CEIL() {
    printf "%.0f\n" $(ADD $1 1)
  }

  set +a
}

##############################################################################################
######                              HOW TO                                              ######
##############################################################################################
# Tests are written sequentially
# Each test depends on the previous tests, and examines the chain state at a point in time

# To add a new test, take an flow then sleep for seconds / blocks / IBC_TX_WAIT_SECONDS
# Reordering existing tests could break them

##############################################################################################
######                TEST BASIC INTO FUNCTIONALITY                                   ######
##############################################################################################
@test "[INTEGRATION-BASIC-$CHAIN_NAME] ibc transfer updates all balances" {
  # get initial balances
  into_user_into_balance_start=$($INTO_MAIN_CMD q bank balance $(INTO_ADDRESS) $INTO_DENOM | GETBAL)
  host_user_into_balance_start=$($HOST_MAIN_CMD 2>&1 q bank balance $HOST_USER_ADDRESS $IBC_INTO_DENOM | GETBAL)
  into_host_user_balance_start=$($INTO_MAIN_CMD q bank balance $(INTO_ADDRESS) $HOST_IBC_DENOM | GETBAL)
  host_host_user_balance_start=$($HOST_MAIN_CMD 2>&1 q bank balance $HOST_USER_ADDRESS $HOST_DENOM | GETBAL)

  $INTO_MAIN_CMD tx ibc-transfer transfer transfer $INTO_TRANFER_CHANNEL $HOST_USER_ADDRESS ${TRANSFER_AMOUNT}${INTO_DENOM} --from $INTO_USER -y
  $HOST_MAIN_CMD_TX ibc-transfer transfer transfer $HOST_TRANSFER_CHANNEL $(INTO_ADDRESS) ${TRANSFER_AMOUNT}${HOST_DENOM} --from $HOST_USER -y

  WAIT_FOR_BLOCK $INTO_LOGS 20

  # get new balances
  into_user_into_balance_end=$($INTO_MAIN_CMD q bank balance $(INTO_ADDRESS) $INTO_DENOM | GETBAL)
  host_user_into_balance_end=$($HOST_MAIN_CMD 2>&1 q bank balance $HOST_USER_ADDRESS $IBC_INTO_DENOM | GETBAL)

  into_user_balance_end=$($INTO_MAIN_CMD q bank balance $(INTO_ADDRESS) $HOST_IBC_DENOM | GETBAL)
  host_user_balance_end=$($HOST_MAIN_CMD 2>&1 q bank balance $HOST_USER_ADDRESS $HOST_DENOM | GETBAL)

  # get all INTO balance diffs
  into_user_into_balance_diff=$((into_user_into_balance_start - into_user_into_balance_end))
  host_user_into_balance_diff=$((host_user_into_balance_start - host_user_into_balance_end))

  assert_equal "$into_user_into_balance_diff" "$TRANSFER_AMOUNT"
  assert_equal "$host_user_into_balance_diff" "-$TRANSFER_AMOUNT"

  # get all host balance diffs
  into_user_balance_diff=$((into_host_user_balance_start - into_user_balance_end))
  host_user_balance_diff=$((host_host_user_balance_start - host_user_balance_end))

  assert_equal "$into_user_balance_diff" "-$TRANSFER_AMOUNT"
  expected_diff=$(($TRANSFER_AMOUNT + $EXPECTED_FEE)) #TRANSFER_AMOUNT+fee
  assert_equal "$host_user_balance_diff" $expected_diff
}

@test "[INTEGRATION-BASIC-$CHAIN_NAME] Flow MsgTransfer" {
  host_receiver_balance_start=$($HOST_MAIN_CMD 2>&1 q bank balance $HOST_RECEIVER_ADDRESS $IBC_INTO_DENOM | GETBAL)
  user_into_balance_start=$($INTO_MAIN_CMD q bank balance $(INTO_ADDRESS) $INTO_DENOM | GETBAL)

  # Define the file path
  msg_transfer_file="msg_transfer.json"

  # build MsgSend with MSGSEND_AMOUNT
  cat <<EOF >"./$msg_transfer_file"
{
  "@type": "/ibc.applications.transfer.v1.MsgTransfer",
    "source_port": "transfer",
    "source_channel": "$INTO_TRANFER_CHANNEL",
    "token": {
      "amount": "$MSGSEND_AMOUNT",
      "denom": "$INTO_DENOM"
    },
    "sender":  "$(INTO_ADDRESS)",
    "receiver": "$HOST_RECEIVER_ADDRESS",
    "timeout_height": {
      "revision_number": "0",
      "revision_height": "0"
    },
    "timeout_timestamp": "2526374086000000000",
    "memo": "hello"
}
EOF

  msg_submit=$($INTO_MAIN_CMD tx intent submit-flow "$msg_transfer_file" --label "MsgTransfer" --duration "60s" --fallback-to-owner-balance --from $INTO_USER -y)
  echo "$msg_submit"

  GET_FLOW_ID $(INTO_ADDRESS)
  WAIT_FOR_EXECUTED_FLOW_BY_ID

  # calculate difference between token balance receiver before and after, should equal MSGSEND_AMOUNT
  host_receiver_balance_end=$($HOST_MAIN_CMD 2>&1 q bank balance $HOST_RECEIVER_ADDRESS $IBC_INTO_DENOM | GETBAL)
  receiver_diff=$(($host_receiver_balance_end - $host_receiver_balance_start))
  assert_equal "$receiver_diff" $MSGSEND_AMOUNT
}

@test "[INTEGRATION-BASIC-$CHAIN_NAME] Flow Update MsgTransfer" {
  msg_update_flow=$($INTO_MAIN_CMD tx intent update-flow $FLOW_ID --label "IBC transfer" --updating-disabled --from $INTO_USER -y)
  echo "$msg_update_flow"

  WAIT_FOR_UPDATING_DISABLED
}

@test "[INTEGRATION-BASIC-$CHAIN_NAME] Flow MsgSend using new ICA" {
  # get initial balances on host account
  host_receiver_balance_start=$($HOST_MAIN_CMD 2>&1 q bank balance $HOST_RECEIVER_ADDRESS $HOST_DENOM | GETBAL)

  # get token balance user on INTO
  user_into_balance_start=$($INTO_MAIN_CMD q bank balance $(INTO_ADDRESS) $INTO_DENOM | GETBAL)
  WAIT_FOR_BLOCK $INTO_LOGS 3
  # build MsgRegisterAccount and retrieve trigger ICA account
  msg_register_account=$($INTO_MAIN_CMD tx intent register --connection-id connection-$CONNECTION_ID --host-connection-id connection-$HOST_CONNECTION_ID --from $INTO_USER -y)
  echo $msg_register_account
  sleep 120

  ica_address=$($INTO_MAIN_CMD q intent interchainaccounts $(INTO_ADDRESS) connection-$CONNECTION_ID)
  ica_address=$(echo "$ica_address" | awk '{print $2}')

  fund_ica=$($HOST_MAIN_CMD_TX bank send $HOST_USER_ADDRESS $ica_address $MSGSEND_AMOUNT$HOST_DENOM --from $HOST_USER -y)

  WAIT_FOR_BLOCK $INTO_LOGS 2

  ica_balance_start=$($HOST_MAIN_CMD 2>&1 q bank balance $ica_address $HOST_DENOM | GETBAL)

  # Define the file path
  msg_send_file="msg_send.json"

  # build MsgSend with MSGSEND_AMOUNT
  cat <<EOF >"./$msg_send_file"
{
   "@type":"/cosmos.bank.v1beta1.MsgSend",
      "amount": [{
         "amount": "$MSGSEND_AMOUNT",
         "denom": "$HOST_DENOM"
        }],
      "from_address": "$ica_address",
      "to_address": "$HOST_RECEIVER_ADDRESS"
}
EOF

  msg_submit_flow=$($INTO_MAIN_CMD tx intent submit-flow "$msg_send_file" --label "Send from Interchain Account" --duration "60s" --connection-id connection-$CONNECTION_ID --from $INTO_USER --fallback-to-owner-balance -y)
  echo "$msg_submit_flow"

  GET_FLOW_ID $(INTO_ADDRESS)
  WAIT_FOR_EXECUTED_FLOW_BY_ID

  # calculate difference between token balance ica before and after, should equal MSGSEND_AMOUNT
  ica_balance_end=$($HOST_MAIN_CMD 2>&1 q bank balance $ica_address $HOST_DENOM | GETBAL)
  ica_diff=$(($ica_balance_start - $ica_balance_end))
  assert_equal "$ica_balance_end" 0

  # calculate difference between token balance receiver before and after, should equal MSGSEND_AMOUNT
  host_receiver_balance_end=$($HOST_MAIN_CMD 2>&1 q bank balance $HOST_RECEIVER_ADDRESS $HOST_DENOM | GETBAL)
  receiver_diff=$(($host_receiver_balance_end - $host_receiver_balance_start))
  assert_equal "$receiver_diff" $MSGSEND_AMOUNT
}

@test "[INTEGRATION-BASIC-$CHAIN_NAME] Flow MsgSend using AuthZ" {
  # get initial balances on host account
  host_user_balance_start=$($HOST_MAIN_CMD 2>&1 q bank balance $HOST_USER_ADDRESS $HOST_DENOM | GETBAL)
  host_receiver_balance_start=$($HOST_MAIN_CMD 2>&1 q bank balance $HOST_RECEIVER_ADDRESS $HOST_DENOM | GETBAL)

  # get token balance user on INTO
  user_into_balance_start=$($INTO_MAIN_CMD q bank balance $(INTO_ADDRESS) $INTO_DENOM | GETBAL)

  ica_address=$($INTO_MAIN_CMD q intent interchainaccounts $(INTO_ADDRESS) connection-$CONNECTION_ID)
  ica_address=$(echo "$ica_address" | awk '{print $2}')

  $HOST_MAIN_CMD_TX authz grant $ica_address generic --msg-type "/cosmos.bank.v1beta1.MsgSend" --from $HOST_USER -y
  WAIT_FOR_BLOCK $INTO_LOGS 3
  $HOST_MAIN_CMD_TX bank send $HOST_USER_ADDRESS $ica_address $MSGSEND_AMOUNT$HOST_DENOM --from $HOST_USER -y

  # Define the file path
  msg_exec_file="msg_exec.json"

  # Write the JSON data to the file
  cat <<EOF >"$msg_exec_file"
{
  "@type": "/cosmos.authz.v1beta1.MsgExec",
  "msgs": [
    {
      "@type": "/cosmos.bank.v1beta1.MsgSend",
      "amount": [
        {
          "amount": "$MSGSEND_AMOUNT",
          "denom": "$HOST_DENOM"
        }
      ],
      "from_address": "$HOST_USER_ADDRESS",
      "to_address": "$HOST_RECEIVER_ADDRESS"
    }
  ],
  "grantee": "$ica_address"
}
EOF

  msg_submit_flow=$($INTO_MAIN_CMD tx intent submit-flow "$msg_exec_file" --label "Send from user on host chain using AuthZ" --duration "60s" --connection-id connection-$CONNECTION_ID --from $INTO_USER --fallback-to-owner-balance -y)
  echo "$msg_submit_flow"

  GET_FLOW_ID $(INTO_ADDRESS)
  WAIT_FOR_EXECUTED_FLOW_BY_ID
  sleep 60

  # calculate difference between token balance of user before and after, should equal MSGSEND_AMOUNT
  user_balance_end=$($HOST_MAIN_CMD 2>&1 q bank balance $HOST_USER_ADDRESS $HOST_DENOM | GETBAL)
  user_diff=$(($host_user_balance_start - $user_balance_end))
  expected_diff=$(($MSGSEND_AMOUNT + $MSGSEND_AMOUNT + $EXPECTED_FEE + $EXPECTED_FEE)) #MsgSend to ICA and MsgSend using AuthZ + host tx fees for MsgGrant,MsgSend
  assert_equal "$user_diff" $expected_diff

  # calculate difference between token balance receiver before and after, should equal MSGSEND_AMOUNT
  host_receiver_balance_end=$($HOST_MAIN_CMD 2>&1 q bank balance $HOST_RECEIVER_ADDRESS $HOST_DENOM | GETBAL)
  receiver_diff=$(($host_receiver_balance_end - $host_receiver_balance_start))
  assert_equal "$receiver_diff" $MSGSEND_AMOUNT #from MsgSend
}

# test flow MsgSend from ICS20 message with Trigger Address ICA Account with MsgSubmitAutoTx ICA_ADDR parsing
@test "[INTEGRATION-BASIC-$CHAIN_NAME] ibc ics20 transfer, create trigger and auto-parse address" {
  # get initial balances
  host_user_balance_start=$($HOST_MAIN_CMD 2>&1 q bank balance $HOST_USER_ADDRESS $HOST_DENOM | GETBAL)
  host_receiver_balance_start=$($HOST_MAIN_CMD 2>&1 q bank balance $HOST_RECEIVER_ADDRESS $HOST_DENOM | GETBAL)

  # do IBC transfer
  memo='{"flow": {"msgs": [{"@type": "/cosmos.bank.v1beta1.MsgSend","amount": [{"amount": "'$MSGSEND_AMOUNT'","denom": "'$HOST_DENOM'"}],"from_address":"ICA_ADDR","to_address": "'$HOST_RECEIVER_ADDRESS'"}],"duration":"60s","label":"MsgSend submitted from ICS20 hook","cid":"connection-'$CONNECTION_ID'","host_cid":"connection-'$HOST_CONNECTION_ID'","start_at":"0", "owner": "'$(INTO_ADDRESS)'", "fallback": "true"}}'
  echo $memo
  msg_transfer=$($HOST_MAIN_CMD_TX ibc-transfer transfer transfer $HOST_TRANSFER_CHANNEL $(INTO_ADDRESS) ${ICS20_AMOUNT_FOR_LOCAL_GAS}${IBC_INTO_DENOM} --memo "$memo" --from $HOST_USER -y)
  echo $msg_transfer
  GET_FLOW_ID $(INTO_ADDRESS)

  ica_address=$($INTO_MAIN_CMD q intent interchainaccounts $(INTO_ADDRESS) connection-$CONNECTION_ID)
  ica_address=$(echo "$ica_address" | awk '{print $2}')
  $HOST_MAIN_CMD_TX bank send $HOST_USER_ADDRESS $ica_address $MSGSEND_AMOUNT$HOST_DENOM --from $HOST_USER -y

  WAIT_FOR_EXECUTED_FLOW_BY_ID
  sleep 60

  # calculate difference between token balance of host user before and after, should equal 2xMSGSEND_AMOUNT
  user_balance_end=$($HOST_MAIN_CMD 2>&1 q bank balance $HOST_USER_ADDRESS $HOST_DENOM | GETBAL)
  user_diff=$(($host_user_balance_start - $user_balance_end))
  expected_diff=$(($MSGSEND_AMOUNT + $EXPECTED_FEE + $EXPECTED_FEE)) #ICS20_MSGSEND_AMOUNT_TOTAL for all MsgSends(10000)+2x host tx fee(2500)
  assert_equal "$user_diff" $expected_diff

  # calculate difference between token balance receiver before and after, should equal 1xMSGSEND_AMOUNT
  host_receiver_balance_end=$($HOST_MAIN_CMD 2>&1 q bank balance $HOST_RECEIVER_ADDRESS $HOST_DENOM | GETBAL)
  receiver_diff=$(($host_receiver_balance_end - $host_receiver_balance_start))

  assert_equal "$receiver_diff" $MSGSEND_AMOUNT #one MsgSend received

}

@test "[INTEGRATION-BASIC-$CHAIN_NAME] Flow MsgSend using Hosted ICA" {
  # get initial balances on host account
  host_user_balance_start=$($HOST_MAIN_CMD 2>&1 q bank balance $HOST_USER_ADDRESS $HOST_DENOM | GETBAL)
  host_receiver_balance_start=$($HOST_MAIN_CMD 2>&1 q bank balance $HOST_RECEIVER_ADDRESS $HOST_DENOM | GETBAL)

  # get token balance user on INTO
  user_into_balance_start=$($INTO_MAIN_CMD q bank balance $(INTO_ADDRESS) $INTO_DENOM | GETBAL)

  # build ICA and retrieve ICA account
  msg_create_hosted_account=$($INTO_MAIN_CMD tx intent create-hosted-account --connection-id connection-$CONNECTION_ID --host-connection-id connection-$HOST_CONNECTION_ID --fee-coins-suported "10"$INTO_DENOM --from $INTO_USER --gas 250000 -y)
  echo $msg_create_hosted_account
  sleep 120

  hosted_accounts=$($INTO_MAIN_CMD q intent list-hosted-accounts --output json)

  # Use jq to filter the hosted_address based on the connection ID
  hosted_address=$(echo "$hosted_accounts" | jq -r --arg conn_id "connection-$CONNECTION_ID" '.hosted_accounts[] | select(.ica_config.connection_id == $conn_id) | .hosted_address')
  if [ -n "$hosted_address" ]; then
    # Get the interchain account address
    ica_address=$($INTO_MAIN_CMD q intent interchainaccounts "$hosted_address" connection-$CONNECTION_ID)
    ica_address=$(echo "$ica_address" | awk '{print $2}')

    echo "Interchain Account Address: $ica_address"
  else
    echo "No hosted address found for connection ID: $CONNECTION_ID"
  fi

  $HOST_MAIN_CMD_TX authz grant $ica_address generic --msg-type "/cosmos.bank.v1beta1.MsgSend" --from $HOST_USER -y
  WAIT_FOR_BLOCK $INTO_LOGS 3
  fund_ica_hosted=$($HOST_MAIN_CMD_TX bank send $HOST_USER_ADDRESS $ica_address $MSGSEND_AMOUNT$HOST_DENOM --from $HOST_USER -y)

  WAIT_FOR_BLOCK $INTO_LOGS 4

  # Define the file path
  msg_exec_file="msg_exec.json"
  cat <<EOF >"$msg_exec_file"
{
  "@type": "/cosmos.authz.v1beta1.MsgExec",
  "msgs": [
    {
      "@type": "/cosmos.bank.v1beta1.MsgSend",
      "amount": [
        {
          "amount": "$MSGSEND_AMOUNT",
          "denom": "$HOST_DENOM"
        }
      ],
      "from_address": "$HOST_USER_ADDRESS",
      "to_address": "$HOST_RECEIVER_ADDRESS"
    }
  ],
  "grantee": "$ica_address"
}
EOF

  msg_submit_flow=$($INTO_MAIN_CMD tx intent submit-flow "$msg_exec_file" --label "MsgSend using Hosted ICA and AuthZ" --duration "60s" --hosted-account $hosted_address --hosted-account-fee-limit 20$INTO_DENOM --from $INTO_USER --fallback-to-owner-balance -y)
  echo "$msg_submit_flow"

  GET_FLOW_ID $(INTO_ADDRESS)
  WAIT_FOR_EXECUTED_FLOW_BY_ID

  sleep 60
  # calculate difference between token balance of user before and after, should equal MSGSEND_AMOUNT
  user_balance_end=$($HOST_MAIN_CMD 2>&1 q bank balance $HOST_USER_ADDRESS $HOST_DENOM | GETBAL)
  user_diff=$(($host_user_balance_start - $user_balance_end))
  expected_diff=$(($MSGSEND_AMOUNT + $MSGSEND_AMOUNT + $EXPECTED_FEE + $EXPECTED_FEE)) #MsgSend to ICA and MsgSend using AuthZ + host tx fees for MsgGrant,MsgSend
  assert_equal "$user_diff" $expected_diff

  # calculate difference between token balance receiver before and after, should equal MSGSEND_AMOUNT
  host_receiver_balance_end=$($HOST_MAIN_CMD 2>&1 q bank balance $HOST_RECEIVER_ADDRESS $HOST_DENOM | GETBAL)
  receiver_diff=$(($host_receiver_balance_end - $host_receiver_balance_start))
  assert_equal "$receiver_diff" $MSGSEND_AMOUNT #from MsgSend
}

@test "[INTEGRATION-BASIC-$CHAIN_NAME] Flow MsgSend using Hosted ICA with ICQ query balance as input" {
  # get initial balances on host account
  host_user_balance_start=$($HOST_MAIN_CMD 2>&1 q bank balance $HOST_USER_ADDRESS $HOST_DENOM | GETBAL)
  host_receiver_balance_start=$($HOST_MAIN_CMD 2>&1 q bank balance $HOST_RECEIVER_ADDRESS $HOST_DENOM | GETBAL)

  # get token balance user on INTO
  user_into_balance_start=$($INTO_MAIN_CMD q bank balance $(INTO_ADDRESS) $INTO_DENOM | GETBAL)

  hosted_accounts=$($INTO_MAIN_CMD q intent list-hosted-accounts --output json)
  # Use jq to filter the hosted_address based on the connection ID
  hosted_address=$(echo "$hosted_accounts" | jq -r --arg conn_id "connection-$CONNECTION_ID" '.hosted_accounts[] | select(.ica_config.connection_id == $conn_id) | .hosted_address')
  if [ -n "$hosted_address" ]; then
    # Get the interchain account address
    ica_address=$($INTO_MAIN_CMD q intent interchainaccounts "$hosted_address" connection-$CONNECTION_ID)
    ica_address=$(echo "$ica_address" | awk '{print $2}')

    echo "Interchain Account Address: $ica_address"
  else
    echo "No hosted address found for connection ID: $CONNECTION_ID"
  fi

  # fund_ica_hosted=$($HOST_MAIN_CMD_TX bank send $HOST_USER_ADDRESS $ica_address $MSGSEND_AMOUNT$HOST_DENOM --from $HOST_USER -y)

  WAIT_FOR_BLOCK $INTO_LOGS 3

  # Define the file path
  msg_exec_file="msg_exec.json"
  cat <<EOF >"$msg_exec_file"
{
  "@type": "/cosmos.authz.v1beta1.MsgExec",
  "msgs": [
    {
      "@type": "/cosmos.bank.v1beta1.MsgSend",
      "amount": [
        {
          "amount": "$MSGSEND_AMOUNT",
          "denom": "$HOST_DENOM"
        }
      ],
      "from_address": "$HOST_USER_ADDRESS",
      "to_address": "$HOST_RECEIVER_ADDRESS"
    }
  ],
  "grantee": "$ica_address"
}
EOF

  #query INTO IBC balance
  #query_key='AhRGgNNqzbhS97+pYwK8+uF7JhF0PGliYy81MjRDNjUyMUI0NDQ4Mjc3QTBFOTgzQTVCN0U2NTZGMDREQ0UzQTZGOTAyRUZGQUM3RDMxMDBFQjQyMEYwREZF'
  query_key='AhRzQ/BoErqMPmPAB1G+lJ3WqA0C+GliYy81MjRDNjUyMUI0NDQ4Mjc3QTBFOTgzQTVCN0U2NTZGMDREQ0UzQTZGOTAyRUZGQUM3RDMxMDBFQjQyMEYwREZF'
  # echo "Query Key: $query_key"
  msg_submit_flow=$($INTO_MAIN_CMD tx intent submit-flow "$msg_exec_file" --label "ICQ and Hosted ICA" --interval "60s" --duration "120s" --hosted-account $hosted_address --hosted-account-fee-limit 20$INTO_DENOM --from $INTO_USER --fallback-to-owner-balance --conditions '{ "feedback_loops": [{"response_index":0,"response_key": "", "msgs_index":0, "msg_key":"Amount.[0].Amount","value_type": "sdk.Int", "from_icq": true, "icq_config": {"connection_id":"connection-'$CONNECTION_ID'","chain_id":"'$HOST_CHAIN_ID'","timeout_policy":2,"timeout_duration":50000000000,"query_type":"store/bank/key","query_key":"'$query_key'"}}] }' -y)
  # echo "$msg_submit_flow"

  GET_FLOW_ID $(INTO_ADDRESS)
  WAIT_FOR_EXECUTED_FLOW_BY_ID

  #sleep 40

  # # calculate difference between token balance receiver before and after, should equal MSGSEND_AMOUNT
  host_receiver_balance_end=$($HOST_MAIN_CMD 2>&1 q bank balance $HOST_RECEIVER_ADDRESS $HOST_DENOM | GETBAL)
  receiver_diff=$(($host_receiver_balance_start - $host_receiver_balance_end))
  # assert_equal "$receiver_diff" $MSGSEND_AMOUNT #from MsgSend
}

@test "[INTEGRATION-BASIC-$CHAIN_NAME] Flow Autocompound on host" {
  # call MsgDelegate
  validator_address=$(GET_VAL_ADDR $HOST_CHAIN_ID 1)
  delegate=$($HOST_MAIN_CMD_TX staking delegate $validator_address $MSGDELEGATE_AMOUNT$HOST_DENOM --from $HOST_USER -y)
  WAIT_FOR_BLOCK $INTO_LOGS 3

  hosted_accounts=$($INTO_MAIN_CMD q intent list-hosted-accounts --output json)
  # Use jq to filter the hosted_address based on the connection ID
  hosted_address=$(echo "$hosted_accounts" | jq -r --arg conn_id "connection-$CONNECTION_ID" '.hosted_accounts[] | select(.ica_config.connection_id == $conn_id) | .hosted_address')
  if [ -n "$hosted_address" ]; then
    # Get the interchain account address
    ica_address=$($INTO_MAIN_CMD q intent interchainaccounts "$hosted_address" connection-$CONNECTION_ID)
    ica_address=$(echo "$ica_address" | awk '{print $2}')

    echo "Interchain Account Address: $ica_address"
  else
    echo "No hosted address found for connection ID: $CONNECTION_ID"
  fi

  $HOST_MAIN_CMD_TX authz grant $ica_address generic --msg-type "/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward" --from $HOST_USER -y
  WAIT_FOR_BLOCK $INTO_LOGS 3
  $HOST_MAIN_CMD_TX authz grant $ica_address generic --msg-type "/cosmos.staking.v1beta1.MsgDelegate" --from $HOST_USER -y
  WAIT_FOR_BLOCK $INTO_LOGS 3

  host_user_balance_start=$($HOST_MAIN_CMD 2>&1 q bank balance $HOST_USER_ADDRESS $HOST_DENOM | GETBAL)

  msg_withdraw="MsgWithdrawDelegatorReward.json"
  cat <<EOF >"$msg_withdraw"
{
  "@type": "/cosmos.authz.v1beta1.MsgExec",
  "msgs": [
    {
      "@type": "/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward",
      "delegator_address": "$HOST_USER_ADDRESS",
      "validator_address": "$validator_address"
    }
  ],
  "grantee": "$ica_address"
}
EOF

  msg_delegate="MsgDelegate.json"
  cat <<EOF >"$msg_delegate"
{
  "@type": "/cosmos.authz.v1beta1.MsgExec",
  "msgs": [
    {
      "@type": "/cosmos.staking.v1beta1.MsgDelegate",
      "delegator_address": "$HOST_USER_ADDRESS",
      "validator_address": "$validator_address",
      "amount": {
          "amount": "10",
          "denom": "$HOST_DENOM"
      }
    }
  ],
  "grantee": "$ica_address"
}
EOF

  msg_submit_flow=$($INTO_MAIN_CMD tx intent submit-flow $msg_withdraw $msg_delegate --label "Autocompound on host chain" --duration "168h" --interval "600s" --hosted-account $hosted_address --hosted-account-fee-limit 20$INTO_DENOM --from $INTO_USER --fallback-to-owner-balance --stop-on-failure --conditions '{ "feedback_loops": [{"response_index":0,"response_key": "Amount.[0]", "msgs_index":1, "msg_key":"Amount","value_type": "sdk.Coin"}]}' -y)
  echo "$msg_submit_flow"

  GET_FLOW_ID $(INTO_ADDRESS)
  WAIT_FOR_EXECUTED_FLOW_BY_ID

  sleep 40
  # # calculate difference between token balance of user before and after, should equal MSGSEND_AMOUNT
  staking_balance_end=$($HOST_MAIN_CMD 2>&1 q staking delegation $HOST_USER_ADDRESS $validator_address $HOST_DENOM | GETBAL)
  staking_balance_diff=$(($staking_balance_end - $MSGDELEGATE_AMOUNT))

  assert_not_equal "$staking_balance_diff" 0
}

@test "[INTEGRATION-BASIC-$CHAIN_NAME] Flow Periodic MsgSend using AuthZ" {
  # get initial balances on host account
  host_user_balance_start=$($HOST_MAIN_CMD 2>&1 q bank balance $HOST_USER_ADDRESS $HOST_DENOM | GETBAL)
  host_receiver_balance_start=$($HOST_MAIN_CMD 2>&1 q bank balance $HOST_RECEIVER_ADDRESS $HOST_DENOM | GETBAL)

  # get token balance user on INTO
  user_into_balance_start=$($INTO_MAIN_CMD q bank balance $(INTO_ADDRESS) $INTO_DENOM | GETBAL)

  ica_address=$($INTO_MAIN_CMD q intent interchainaccounts $(INTO_ADDRESS) connection-$CONNECTION_ID)
  ica_address=$(echo "$ica_address" | awk '{print $2}')

  # Define the file path
  msg_exec_file="msg_exec.json"

  # Write the JSON data to the file
  cat <<EOF >"$msg_exec_file"
{
  "@type": "/cosmos.authz.v1beta1.MsgExec",
  "msgs": [
    {
      "@type": "/cosmos.bank.v1beta1.MsgSend",
      "amount": [
        {
          "amount": "$MSGSEND_AMOUNT",
          "denom": "$HOST_DENOM"
        }
      ],
      "from_address": "$HOST_USER_ADDRESS",
      "to_address": "$HOST_RECEIVER_ADDRESS"
    }
  ],
  "grantee": "$ica_address"
}
EOF

  msg_submit_flow=$($INTO_MAIN_CMD tx intent submit-flow "$msg_exec_file" --label "Recurring transfer on host chain from host user" --duration "2880h" --interval "60s" --stop-on-failure --fee-funds $RECURRING_MSGSEND_AMOUNT_TOTAL$INTO_DENOM --connection-id connection-$CONNECTION_ID --from $INTO_USER --fallback-to-owner-balance --stop-on-timeout -y)
  echo "$msg_submit_flow"

  GET_FLOW_ID $(INTO_ADDRESS)
  WAIT_FOR_EXECUTED_FLOW_BY_ID
  sleep 20
  # calculate difference between token balance receiver before and after, should equal MSGSEND_AMOUNT
  host_receiver_balance_mid=$($HOST_MAIN_CMD 2>&1 q bank balance $HOST_RECEIVER_ADDRESS $HOST_DENOM | GETBAL)
  receiver_diff=$(($host_receiver_balance_mid - $host_receiver_balance_start))
  assert_equal "$receiver_diff" $MSGSEND_AMOUNT

  sleep 60
  # WAIT_FOR_EXECUTED_FLOW_BY_ID

  # calculate difference between token balance receiver before and after, should equal MSGSEND_AMOUNT
  host_receiver_balance_end=$($HOST_MAIN_CMD 2>&1 q bank balance $HOST_RECEIVER_ADDRESS $HOST_DENOM | GETBAL)
  receiver_diff=$(($host_receiver_balance_end - $host_receiver_balance_mid))
  assert_equal "$receiver_diff" $MSGSEND_AMOUNT
}

@test "[INTEGRATION-BASIC-$CHAIN_NAME] Flow Conditional Autocompound on host" {
  # call MsgDelegate
  validator_address=$(GET_VAL_ADDR $HOST_CHAIN_ID 1)
  delegate=$($HOST_MAIN_CMD_TX staking delegate $validator_address $MSGDELEGATE_AMOUNT$HOST_DENOM --from $HOST_USER -y)
  WAIT_FOR_BLOCK $INTO_LOGS 3

  hosted_accounts=$($INTO_MAIN_CMD q intent list-hosted-accounts --output json)
  # Use jq to filter the hosted_address based on the connection ID
  hosted_address=$(echo "$hosted_accounts" | jq -r --arg conn_id "connection-$CONNECTION_ID" '.hosted_accounts[] | select(.ica_config.connection_id == $conn_id) | .hosted_address')
  if [ -n "$hosted_address" ]; then
    # Get the interchain account address
    ica_address=$($INTO_MAIN_CMD q intent interchainaccounts "$hosted_address" connection-$CONNECTION_ID)
    ica_address=$(echo "$ica_address" | awk '{print $2}')

    echo "Interchain Account Address: $ica_address"
  else
    echo "No hosted address found for connection ID: $CONNECTION_ID"
  fi

  host_user_balance_start=$($HOST_MAIN_CMD 2>&1 q bank balance $HOST_USER_ADDRESS $HOST_DENOM | GETBAL)

  msg_withdraw="MsgWithdrawDelegatorReward.json"
  cat <<EOF >"$msg_withdraw"
{
  "@type": "/cosmos.authz.v1beta1.MsgExec",
  "msgs": [
    {
      "@type": "/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward",
      "delegator_address": "$HOST_USER_ADDRESS",
      "validator_address": "$validator_address"
    }
  ],
  "grantee": "$ica_address"
}
EOF

  msg_delegate="MsgDelegate.json"
  cat <<EOF >"$msg_delegate"
{
  "@type": "/cosmos.authz.v1beta1.MsgExec",
  "msgs": [
    {
      "@type": "/cosmos.staking.v1beta1.MsgDelegate",
      "delegator_address": "$HOST_USER_ADDRESS",
      "validator_address": "$validator_address",
      "amount": {
          "amount": "10",
          "denom": "$HOST_DENOM"
      }
    }
  ],
  "grantee": "$ica_address"
}
EOF

  msg_submit_flow=$($INTO_MAIN_CMD tx intent submit-flow $msg_withdraw $msg_delegate --label "Conditional Autocompound" --duration "168h" --interval "2400s" --hosted-account $hosted_address --hosted-account-fee-limit 20$INTO_DENOM --from $INTO_USER --fallback-to-owner-balance --conditions '{ "feedback_loops": [{"response_index":0,"response_key": "Amount.[0]", "msgs_index":1, "msg_key":"Amount","value_type": "sdk.Coin"}], "comparisons": [{"response_index":0,"response_key": "Amount.[0]", "operand":"1'$HOST_DENOM'", "operator":4,"value_type": "sdk.Coin"}]}' --save-responses -y)
  echo "$msg_submit_flow"

  GET_FLOW_ID $(INTO_ADDRESS)
  WAIT_FOR_EXECUTED_FLOW_BY_ID

  # sleep 40
  # # calculate difference between token balance of user before and after, should equal MSGSEND_AMOUNT
  staking_balance_end=$($HOST_MAIN_CMD 2>&1 q staking delegation $HOST_USER_ADDRESS $validator_address $HOST_DENOM | GETBAL)
  staking_balance_diff=$(($staking_balance_end - $MSGDELEGATE_AMOUNT))

  assert_not_equal "$staking_balance_diff" 0
}