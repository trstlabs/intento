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
  # random address with nil balance
  HOST_RECEIVER_ADDRESS=$(GET_VAR_VALUE ${CHAIN_NAME}_RECEIVER_ADDRESS)

  TRST_USER=$(GET_VAR_VALUE TRST_USER_ACCT)
  TRST_VAL=${TRST_VAL_PREFIX}1

  TRST_TRANFER_CHANNEL="channel-${TRANSFER_CHANNEL_NUMBER}"
  HOST_TRANSFER_CHANNEL="channel-0"

  TRANSFER_AMOUNT=5000000
  MSGSEND_AMOUNT=100000
  REDEEM_AMOUNT=10000
  PACKET_FORWARD_MSGSEND_AMOUNT=30000

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
######                TEST BASIC TRST FUNCTIONALITY                                   ######
##############################################################################################
@test "[INTEGRATION-BASIC-$CHAIN_NAME] ibc transfer updates all balances" {
  # get initial balances
  sval_trst_balance_start=$($TRST_MAIN_CMD q bank balances $(TRST_ADDRESS) --denom $TRST_DENOM | GETBAL)
  hval_trst_balance_start=$($HOST_MAIN_CMD q bank balances $HOST_USER_ADDRESS --denom $IBC_TRST_DENOM | GETBAL)
  printf "hval_trst_balance_start (sender): %s\n" "$hval_trst_balance_start"
  sval_token_balance_start=$($TRST_MAIN_CMD q bank balances $(TRST_ADDRESS) --denom $HOST_IBC_DENOM | GETBAL)
  hval_token_balance_start=$($HOST_MAIN_CMD q bank balances $HOST_USER_ADDRESS --denom $HOST_DENOM | GETBAL)
  printf "hval_token_balance_start (sender): %s\n" "$hval_token_balance_start"
  # do IBC transfer
  $TRST_MAIN_CMD tx ibc-transfer transfer transfer $TRST_TRANFER_CHANNEL $HOST_USER_ADDRESS ${TRANSFER_AMOUNT}${TRST_DENOM} --from $TRST_USER -y
  $HOST_MAIN_CMD tx ibc-transfer transfer transfer $HOST_TRANSFER_CHANNEL $(TRST_ADDRESS) ${TRANSFER_AMOUNT}${HOST_DENOM} --from $HOST_USER -y

  WAIT_FOR_BLOCK $TRST_LOGS 8

  # get new balances
  sval_trst_balance_end=$($TRST_MAIN_CMD q bank balances $(TRST_ADDRESS) --denom $TRST_DENOM | GETBAL)
  hval_trst_balance_end=$($HOST_MAIN_CMD q bank balances $HOST_USER_ADDRESS --denom $IBC_TRST_DENOM | GETBAL)
    printf "hval_trst_balance_end (sender): %s\n" "$hval_trst_balance_end"
  sval_token_balance_end=$($TRST_MAIN_CMD q bank balances $(TRST_ADDRESS) --denom $HOST_IBC_DENOM | GETBAL)
  hval_token_balance_end=$($HOST_MAIN_CMD q bank balances $HOST_USER_ADDRESS --denom $HOST_DENOM | GETBAL)
  printf "hval_token_balance_end (sender): %s\n" "$hval_token_balance_end"
  # get all TRST balance diffs
  sval_trst_balance_diff=$((sval_trst_balance_start - sval_trst_balance_end))
  hval_trst_balance_diff=$((hval_trst_balance_start - hval_trst_balance_end))
  printf "TRST balance diff (sender): %s\n" "$sval_trst_balance_diff"
  printf "TRST balance diff (receiver): %s\n" "$hval_trst_balance_diff"
  assert_equal "$sval_trst_balance_diff" "$TRANSFER_AMOUNT"
  assert_equal "$hval_trst_balance_diff" "-$TRANSFER_AMOUNT"

  # get all host balance diffs
  sval_token_balance_diff=$((sval_token_balance_start - sval_token_balance_end))
  hval_token_balance_diff=$((hval_token_balance_start - hval_token_balance_end))
  printf "Host balance diff (sender): %s\n" "$sval_token_balance_diff"
  printf "Host balance diff (receiver): %s\n" "$hval_token_balance_diff"
  assert_equal "$sval_token_balance_diff" "-$TRANSFER_AMOUNT"
  assert_equal "$hval_token_balance_diff" "$TRANSFER_AMOUNT"
}

@test "[INTEGRATION-BASIC-$CHAIN_NAME] AutoIbcTx MsgSend using AuthZ" {
  # get initial balances on host account
  user_token_balance_start=$($HOST_MAIN_CMD q bank balances $HOST_USER_ADDRESS --denom $HOST_DENOM | GETBAL)
  receiver_token_balance_start=$($HOST_MAIN_CMD q bank balances $HOST_RECEIVER_ADDRESS --denom $HOST_DENOM | GETBAL)


  # get token balance user on TRST
  user_trst_balance_start=$($TRST_MAIN_CMD q bank balances $(TRST_ADDRESS) --denom $TRST_DENOM | GETBAL)

  # build MsgRegisterAccount and retrieve trigger ICA account
  $TRST_MAIN_CMD tx autoibctx register --connection-id connection-$CONNECTION_ID --counterparty-connection-id connection-0 --from $TRST_USER -y

  sleep 20
  ICA_ADDRESS=$($TRST_MAIN_CMD q autoibctx interchainaccounts $(TRST_ADDRESS) connection-$CONNECTION_ID)
  ICA_ADDRESS=$(echo "$ICA_ADDRESS" | awk '{print $2}')
  echo "ICA ADDR: $ICA_ADDRESS"

  grant_ica=$($HOST_MAIN_CMD tx authz grant $ICA_ADDRESS generic --msg-type "/cosmos.bank.v1beta1.MsgSend" --from $HOST_USER -y)
  echo "$grant_ica"

  WAIT_FOR_BLOCK $TRST_LOGS 2

  grant=$($HOST_MAIN_CMD q authz grants-by-grantee $ICA_ADDRESS)
  echo "GRANT" "$grant"


  fund_ica=$($HOST_MAIN_CMD tx bank send $HOST_USER_ADDRESS $ICA_ADDRESS 120000$HOST_DENOM --from $HOST_USER -y)
  echo "FUND RESP" "$fund_ica"

  # build MsgSend with MSGSEND_AMOUNT
  msg_send='{
      "@type":"/cosmos.bank.v1beta1.MsgSend",
        "amount": [{
         "amount": "'$MSGSEND_AMOUNT'",
         "denom": "'$HOST_DENOM'"
        }],
      "from_address": "'$ICA_ADDRESS'",
      "to_address": "'$HOST_RECEIVER_ADDRESS'"
  }'

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
  "grantee": "$ICA_ADDRESS"
}
EOF

  # build MsgSubmitAutoTx with MsgSend, 60sec non-recurring
  msg_submit_auto_tx=$($TRST_MAIN_CMD tx autoibctx submit-auto-tx "$msg_exec_file" --label "test" --duration "60s" --connection-id connection-$CONNECTION_ID --from $TRST_USER -y)
  echo "$msg_submit_auto_tx"

  # WAIT_FOR_AUTO_TX wait for 30blocks
  WAIT_FOR_BLOCK $TRST_LOGS 10
  # Query the autoibctx to get the initial_autotxs output
  autotxs=$($TRST_MAIN_CMD q autoibctx list-auto-txs-by-owner $(TRST_ADDRESS))
  
  # Count the occurrences of 'duration:' using grep
  occurrences=$(echo "$autotxs" | grep -c 'fee_address:')
  echo "Number of occurrences: $occurrences"
  
  WAIT_FOR_BLOCK $TRST_LOGS 40 #10 blocks of 6 seconds to trigger AutoTx, 10 blocks to execute on host and call back

  # calculate difference between token balance of user before and after, should equal MSGSEND_AMOUNT
  # user_token_balance_end=$($HOST_MAIN_CMD q bank balances $HOST_USER_ADDRESS --denom $HOST_DENOM | GETBAL)
  # user_diff=$(($user_token_balance_start - $user_token_balance_end))
  # printf "Balance start: %s\n" "$user_token_balance_start"
  # printf "Balance end: %s\n" "$user_token_balance_end"
  # assert_equal "$user_diff" $MSGSEND_AMOUNT

  # calculate difference between token balance receiver before and after, should equal MSGSEND_AMOUNT
  receiver_token_balance_end=$($HOST_MAIN_CMD q bank balances $HOST_RECEIVER_ADDRESS --denom $HOST_DENOM | GETBAL)
  receiver_diff=$(( $receiver_token_balance_end - $receiver_token_balance_start))
  printf "Balance end: %s\n" "$receiver_token_balance_end"
  assert_equal "$receiver_diff" $MSGSEND_AMOUNT
}

# TODO
# test auto-tx MsgSend with user address over AuthZ grant to ICA Account

# test auto-tx MsgSend with user address over AuthZ grant to ICA Account with MsgRegisterAccountAndSubmitAutoTx ICA_ADDR parsing

# test auto-tx MsgSend from ICS20 message with Trigger Address ICA Account

# test auto-tx MsgSend from ICS20 message with user address over AuthZ grant to ICA Account

# test auto-tx MsgWithdrawalRewards, MsgSubmitProposal
