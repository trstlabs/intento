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
  HOST_MAIN_CMD_TX="$HOST_MAIN_CMD tx --fees 2500$HOST_DENOM"
  HOST_USER_ADDRESS=$(GET_VAR_VALUE ${CHAIN_NAME}_USER_ADDRESS)
  HOST_USER=$(GET_VAR_VALUE ${CHAIN_NAME}_USER_ACCT)
  HOST_TRANSFER_CHANNEL="channel-0"
  HOST_RECEIVER_ADDRESS=$(GET_VAR_VALUE ${CHAIN_NAME}_RECEIVER_ADDRESS) # random address with nil balance

  INTO_USER=$(GET_VAR_VALUE INTO_USER_ACCT)
  INTO_VAL=${INTO_VAL_PREFIX}1
  INTO_TRANFER_CHANNEL="channel-${TRANSFER_CHANNEL_NUMBER}"

  TRANSFER_AMOUNT=50000000
  MSGSEND_AMOUNT=100000
  EXPECTED_FEE=2500
  ICS20_MSGSEND_AMOUNT_TOTAL=200000          #100000*2
  RECURRING_MSGSEND_AMOUNT_TOTAL=17280000000 #100000*120*24*60
  ICS20_AMOUNT_FOR_LOCAL_GAS=50000
  #PACKET_FORWARD_MSGSEND_AMOUNT=30000

  GETBAL() {
    head -n 1 | grep -o -E '[0-9]+' || "0"
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

# To add a new test, take an action then sleep for seconds / blocks / IBC_TX_WAIT_SECONDS
# Reordering existing tests could break them

##############################################################################################
######                TEST BASIC INTO FUNCTIONALITY                                   ######
##############################################################################################
@test "[INTEGRATION-BASIC-$CHAIN_NAME] ibc transfer updates all balances" {
  # get initial balances
  into_user_into_balance_start=$($INTO_MAIN_CMD q bank balances $(INTO_ADDRESS) --denom $INTO_DENOM | GETBAL)
  host_user_into_balance_start=$($HOST_MAIN_CMD q bank balances $HOST_USER_ADDRESS --denom $IBC_INTO_DENOM | GETBAL)

  into_host_user_balance_start=$($INTO_MAIN_CMD q bank balances $(INTO_ADDRESS) --denom $HOST_IBC_DENOM | GETBAL)
  host_host_user_balance_start=$($HOST_MAIN_CMD q bank balances $HOST_USER_ADDRESS --denom $HOST_DENOM | GETBAL)

  $INTO_MAIN_CMD tx ibc-transfer transfer transfer $INTO_TRANFER_CHANNEL $HOST_USER_ADDRESS ${TRANSFER_AMOUNT}${INTO_DENOM} --from $INTO_USER -y
  $HOST_MAIN_CMD_TX ibc-transfer transfer transfer $HOST_TRANSFER_CHANNEL $(INTO_ADDRESS) ${TRANSFER_AMOUNT}${HOST_DENOM} --from $HOST_USER -y

  WAIT_FOR_BLOCK $INTO_LOGS 8

  # get new balances
  into_user_into_balance_end=$($INTO_MAIN_CMD q bank balances $(INTO_ADDRESS) --denom $INTO_DENOM | GETBAL)
  host_user_into_balance_end=$($HOST_MAIN_CMD q bank balances $HOST_USER_ADDRESS --denom $IBC_INTO_DENOM | GETBAL)

  into_user_balance_end=$($INTO_MAIN_CMD q bank balances $(INTO_ADDRESS) --denom $HOST_IBC_DENOM | GETBAL)
  host_user_balance_end=$($HOST_MAIN_CMD q bank balances $HOST_USER_ADDRESS --denom $HOST_DENOM | GETBAL)

  # get all INTO balance diffs
  into_user_into_balance_diff=$((into_user_into_balance_start - into_user_into_balance_end))
  host_user_into_balance_diff=$((host_user_into_balance_start - host_user_into_balance_end))

  assert_equal "$into_user_into_balance_diff" "$TRANSFER_AMOUNT"
  assert_equal "$host_user_into_balance_diff" "-$TRANSFER_AMOUNT"

  # get all host balance diffs
  into_user_balance_diff=$((into_host_user_balance_start - into_user_balance_end))
  host_user_balance_diff=$((host_host_user_balance_start - host_user_balance_end))

  assert_equal "$into_user_balance_diff" "-$TRANSFER_AMOUNT"
  assert_equal "$host_user_balance_diff" "50002500" #TRANSFER_AMOUNT+fee
}

@test "[INTEGRATION-BASIC-$CHAIN_NAME] Action MsgTransfer" {
  host_receiver_balance_start=$($HOST_MAIN_CMD q bank balances $HOST_RECEIVER_ADDRESS --denom $IBC_INTO_DENOM | GETBAL)
  user_into_balance_start=$($INTO_MAIN_CMD q bank balances $(INTO_ADDRESS) --denom $INTO_DENOM | GETBAL)

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

  msg_submit=$($INTO_MAIN_CMD tx intent submit-action "$msg_transfer_file" --label "MsgTransfer" --duration "60s" --fallback-to-owner-balance --from $INTO_USER -y)
  echo "$msg_submit"

  GET_ACTION_ID $(INTO_ADDRESS)
  WAIT_FOR_EXECUTED_ACTION_BY_ID

  # calculate difference between token balance receiver before and after, should equal MSGSEND_AMOUNT
  host_receiver_balance_end=$($HOST_MAIN_CMD q bank balances $HOST_RECEIVER_ADDRESS --denom $IBC_INTO_DENOM | GETBAL)
  receiver_diff=$(($host_receiver_balance_end - $host_receiver_balance_start))
  assert_equal "$receiver_diff" $MSGSEND_AMOUNT
}

@test "[INTEGRATION-BASIC-$CHAIN_NAME] Action MsgSend" {
  receiver_into_balance_start=$($INTO_MAIN_CMD q bank balances $INTO_RECEIVER_ADDRESS --denom $INTO_DENOM | GETBAL)
  user_into_balance_start=$($INTO_MAIN_CMD q bank balances $(INTO_ADDRESS) --denom $INTO_DENOM | GETBAL)

  # Define the file path
  msg_send_file="msg_send.json"

  # build MsgSend with MSGSEND_AMOUNT
  cat <<EOF >"./$msg_send_file"
{
   "@type":"/cosmos.bank.v1beta1.MsgSend",
      "amount": [{
         "amount": "$MSGSEND_AMOUNT",
         "denom": "$INTO_DENOM"
        }],
      "from_address": "$(INTO_ADDRESS)",
      "to_address": "$INTO_RECEIVER_ADDRESS"
}
EOF

  msg_submit_action=$($INTO_MAIN_CMD tx intent submit-action "$msg_send_file" --label "MsgSend 12345" --duration "60s" --fallback-to-owner-balance --from $INTO_USER -y)
  echo "$msg_submit_action"

  GET_ACTION_ID $(INTO_ADDRESS)
  WAIT_FOR_EXECUTED_ACTION_BY_ID

  # calculate difference between token balance receiver before and after, should equal MSGSEND_AMOUNT
  receiver_into_balance_end=$($INTO_MAIN_CMD q bank balances $INTO_RECEIVER_ADDRESS --denom $INTO_DENOM | GETBAL)
  receiver_diff=$(($receiver_into_balance_end - $receiver_into_balance_start))
  assert_equal "$receiver_diff" $MSGSEND_AMOUNT
}

@test "[INTEGRATION-BASIC-$CHAIN_NAME] Action Update MsgSend" {
  msg_update_action=$($INTO_MAIN_CMD tx intent update-action $ACTION_ID --label "MsgSend" --updating-disabled --from $INTO_USER -y)
  echo "$msg_update_action"

  WAIT_FOR_UPDATING_DISABLED
  sleep 10
}

@test "[INTEGRATION-BASIC-$CHAIN_NAME] Action MsgSend using new ICA" {
  # get initial balances on host account
  host_receiver_balance_start=$($HOST_MAIN_CMD q bank balances $HOST_RECEIVER_ADDRESS --denom $HOST_DENOM | GETBAL)

  # get token balance user on INTO
  user_into_balance_start=$($INTO_MAIN_CMD q bank balances $(INTO_ADDRESS) --denom $INTO_DENOM | GETBAL)

  # build MsgRegisterAccount and retrieve trigger ICA account
  $INTO_MAIN_CMD tx intent register --connection-id connection-$CONNECTION_ID --host-connection-id connection-0 --from $INTO_USER -y

  sleep 40

  ica_address=$($INTO_MAIN_CMD q intent interchainaccounts $(INTO_ADDRESS) connection-$CONNECTION_ID)
  ica_address=$(echo "$ica_address" | awk '{print $2}')

  fund_ica=$($HOST_MAIN_CMD_TX bank send $HOST_USER_ADDRESS $ica_address $MSGSEND_AMOUNT$HOST_DENOM --from $HOST_USER -y)

  WAIT_FOR_BLOCK $INTO_LOGS 2

  ica_balance_start=$($HOST_MAIN_CMD q bank balances $ica_address --denom $HOST_DENOM | GETBAL)

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

  msg_submit_action=$($INTO_MAIN_CMD tx intent submit-action "$msg_send_file" --label "MsgSend from ICA" --duration "60s" --connection-id connection-$CONNECTION_ID --host-connection-id connection-0 --from $INTO_USER --fallback-to-owner-balance -y)
  echo "$msg_submit_action"

  GET_ACTION_ID $(INTO_ADDRESS)
  WAIT_FOR_EXECUTED_ACTION_BY_ID

  # calculate difference between token balance user before and after, should equal MSGSEND_AMOUNT
  ica_balance_end=$($HOST_MAIN_CMD q bank balances $ica_address --denom $HOST_DENOM | GETBAL)
  ica_diff=$(($ica_balance_start - $ica_balance_end))
  assert_equal "$ica_balance_end" 0

  # calculate difference between token balance receiver before and after, should equal MSGSEND_AMOUNT
  host_receiver_balance_end=$($HOST_MAIN_CMD q bank balances $HOST_RECEIVER_ADDRESS --denom $HOST_DENOM | GETBAL)
  receiver_diff=$(($host_receiver_balance_end - $host_receiver_balance_start))
  assert_equal "$receiver_diff" $MSGSEND_AMOUNT
}

@test "[INTEGRATION-BASIC-$CHAIN_NAME] Action MsgSend using AuthZ" {
  # get initial balances on host account
  host_user_balance_start=$($HOST_MAIN_CMD q bank balances $HOST_USER_ADDRESS --denom $HOST_DENOM | GETBAL)
  host_receiver_balance_start=$($HOST_MAIN_CMD q bank balances $HOST_RECEIVER_ADDRESS --denom $HOST_DENOM | GETBAL)

  # get token balance user on INTO
  user_into_balance_start=$($INTO_MAIN_CMD q bank balances $(INTO_ADDRESS) --denom $INTO_DENOM | GETBAL)

  ica_address=$($INTO_MAIN_CMD q intent interchainaccounts $(INTO_ADDRESS) connection-$CONNECTION_ID)
  ica_address=$(echo "$ica_address" | awk '{print $2}')

  $HOST_MAIN_CMD_TX authz grant $ica_address generic --msg-type "/cosmos.bank.v1beta1.MsgSend" --from $HOST_USER -y
  WAIT_FOR_BLOCK $INTO_LOGS 2
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

  msg_submit_action=$($INTO_MAIN_CMD tx intent submit-action "$msg_exec_file" --label "MsgSend from user on host chain using AuthZ" --duration "60s" --connection-id connection-$CONNECTION_ID --host-connection-id connection-0 --from $INTO_USER --fallback-to-owner-balance -y)
  echo "$msg_submit_action"

  GET_ACTION_ID $(INTO_ADDRESS)
  WAIT_FOR_EXECUTED_ACTION_BY_ID

  # calculate difference between token balance of user before and after, should equal MSGSEND_AMOUNT
  user_balance_end=$($HOST_MAIN_CMD q bank balances $HOST_USER_ADDRESS --denom $HOST_DENOM | GETBAL)
  user_diff=$(($host_user_balance_start - $user_balance_end))
  expected_diff=$(($MSGSEND_AMOUNT + $MSGSEND_AMOUNT + $EXPECTED_FEE + $EXPECTED_FEE))  #MsgSend to ICA and MsgSend using AuthZ + host tx fees for MsgGrant,MsgSend
  assert_equal "$user_diff" $expected_diff

  # calculate difference between token balance receiver before and after, should equal MSGSEND_AMOUNT
  host_receiver_balance_end=$($HOST_MAIN_CMD q bank balances $HOST_RECEIVER_ADDRESS --denom $HOST_DENOM | GETBAL)
  receiver_diff=$(($host_receiver_balance_end - $host_receiver_balance_start))
  assert_equal "$receiver_diff" $MSGSEND_AMOUNT #from MsgSend
}

# test action MsgSend from ICS20 message with Trigger Address ICA Account with MsgSubmitAutoTx ICA_ADDR parsing
@test "[INTEGRATION-BASIC-$CHAIN_NAME] ibc ics20 transfer, create trigger and auto-parse address" {
  # get initial balances
  host_user_balance_start=$($HOST_MAIN_CMD q bank balances $HOST_USER_ADDRESS --denom $HOST_DENOM | GETBAL)
  host_receiver_balance_start=$($HOST_MAIN_CMD q bank balances $HOST_RECEIVER_ADDRESS --denom $HOST_DENOM | GETBAL)

  # do IBC transfer
  memo='{"action": {"msgs": [{"@type": "/cosmos.bank.v1beta1.MsgSend","amount": [{"amount": "'$MSGSEND_AMOUNT'","denom": "'$HOST_DENOM'"}],"from_address":"ICA_ADDR","to_address": "'$HOST_RECEIVER_ADDRESS'"}],"duration":"60s","label":"MsgSend submitted from ICS20 hook","cid":"connection-'$CONNECTION_ID'","host_cid":"connection-0","start_at":"0", "owner": "'$(INTO_ADDRESS)'"}}'
  $HOST_MAIN_CMD_TX ibc-transfer transfer transfer $HOST_TRANSFER_CHANNEL $(INTO_ADDRESS) ${ICS20_AMOUNT_FOR_LOCAL_GAS}${IBC_INTO_DENOM} --memo "$memo" --from $HOST_USER -y

  GET_ACTION_ID $(INTO_ADDRESS)

  ica_address=$($INTO_MAIN_CMD q intent interchainaccounts $(INTO_ADDRESS) connection-$CONNECTION_ID)
  ica_address=$(echo "$ica_address" | awk '{print $2}')
  $HOST_MAIN_CMD_TX bank send $HOST_USER_ADDRESS $ica_address $MSGSEND_AMOUNT$HOST_DENOM --from $HOST_USER -y

  WAIT_FOR_EXECUTED_ACTION_BY_ID

  # calculate difference between token balance of host user before and after, should equal 2xMSGSEND_AMOUNT
  user_balance_end=$($HOST_MAIN_CMD q bank balances $HOST_USER_ADDRESS --denom $HOST_DENOM | GETBAL)
  user_diff=$(($host_user_balance_start - $user_balance_end))
  expected_diff=$(($MSGSEND_AMOUNT + $EXPECTED_FEE + $EXPECTED_FEE)) #ICS20_MSGSEND_AMOUNT_TOTAL for all MsgSends(10000)+2x host tx fee(2500)
  assert_equal "$user_diff" $expected_diff

  # calculate difference between token balance receiver before and after, should equal 1xMSGSEND_AMOUNT
  host_receiver_balance_end=$($HOST_MAIN_CMD q bank balances $HOST_RECEIVER_ADDRESS --denom $HOST_DENOM | GETBAL)
  receiver_diff=$(($host_receiver_balance_end - $host_receiver_balance_start))
  assert_equal "$receiver_diff" $MSGSEND_AMOUNT #one MsgSend received

}

@test "[INTEGRATION-BASIC-$CHAIN_NAME] Action Periodic MsgSend using AuthZ" {
  # get initial balances on host account
  host_user_balance_start=$($HOST_MAIN_CMD q bank balances $HOST_USER_ADDRESS --denom $HOST_DENOM | GETBAL)
  host_receiver_balance_start=$($HOST_MAIN_CMD q bank balances $HOST_RECEIVER_ADDRESS --denom $HOST_DENOM | GETBAL)

  # get token balance user on INTO
  user_into_balance_start=$($INTO_MAIN_CMD q bank balances $(INTO_ADDRESS) --denom $INTO_DENOM | GETBAL)

  ica_address=$($INTO_MAIN_CMD q intent interchainaccounts $(INTO_ADDRESS) connection-$CONNECTION_ID)
  ica_address=$(echo "$ica_address" | awk '{print $2}')

  # $HOST_MAIN_CMD_TX authz grant $ica_address generic --msg-type "/cosmos.bank.v1beta1.MsgSend" --from $HOST_USER -y
  # WAIT_FOR_BLOCK $INTO_LOGS 2
  # $HOST_MAIN_CMD_TX bank send $HOST_USER_ADDRESS $ica_address $MSGSEND_AMOUNT$HOST_DENOM --from $HOST_USER -y

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

  msg_submit_action=$($INTO_MAIN_CMD tx intent submit-action "$msg_exec_file" --label "Recurring transfer on host chain from host user" --duration "2880h" --interval "120s" --fee-funds $RECURRING_MSGSEND_AMOUNT_TOTAL$INTO_DENOM --connection-id connection-$CONNECTION_ID --host-connection-id connection-0 --from $INTO_USER --fallback-to-owner-balance --reregister_ica_after_timeout -y)
  echo "$msg_submit_action"

  GET_ACTION_ID $(INTO_ADDRESS)
  WAIT_FOR_EXECUTED_ACTION_BY_ID
  sleep 10
  # calculate difference between token balance receiver before and after, should equal MSGSEND_AMOUNT
  host_receiver_balance_mid=$($HOST_MAIN_CMD q bank balances $HOST_RECEIVER_ADDRESS --denom $HOST_DENOM | GETBAL)
  receiver_diff=$(($host_receiver_balance_mid - $host_receiver_balance_start))
  assert_equal "$receiver_diff" $MSGSEND_AMOUNT

  sleep 120
  # WAIT_FOR_EXECUTED_ACTION_BY_ID

  # calculate difference between token balance receiver before and after, should equal MSGSEND_AMOUNT
  host_receiver_balance_end=$($HOST_MAIN_CMD q bank balances $HOST_RECEIVER_ADDRESS --denom $HOST_DENOM | GETBAL)
  receiver_diff=$(($host_receiver_balance_end - $host_receiver_balance_mid))
  assert_equal "$receiver_diff" $MSGSEND_AMOUNT
}

# TODO
# test action ICS20 MsgSend with user address over AuthZ grant to ICA Account with MsgRegisterAccountAndSubmitAutoTx ICA_ADDR parsing
# test action with other msgs like MsgWithdrawalRewards, MsgSubmitProposal
