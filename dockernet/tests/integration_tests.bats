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
  # random address with nil balance
  HOST_RECEIVER_ADDRESS=$(GET_VAR_VALUE ${CHAIN_NAME}_RECEIVER_ADDRESS)

  TRST_USER=$(GET_VAR_VALUE TRST_USER_ACCT)
  TRST_VAL=${TRST_VAL_PREFIX}1

  TRST_TRANFER_CHANNEL="channel-${TRANSFER_CHANNEL_NUMBER}"
  HOST_TRANSFER_CHANNEL="channel-0"

  TRANSFER_AMOUNT=5000000
  MSGSEND_AMOUNT=100000
  MSGSEND_AMOUNT_TOTAL=17280000000 #100000*120*24*60
  ICS20HOOK_AMOUNT=50000

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
  trst_user_trst_balance_start=$($TRST_MAIN_CMD q bank balances $(TRST_ADDRESS) --denom $TRST_DENOM | GETBAL)
  host_user_trst_balance_start=$($HOST_MAIN_CMD q bank balances $HOST_USER_ADDRESS --denom $IBC_TRST_DENOM | GETBAL)
  printf "host_user_trst_balance_start (sender): %s\n" "$host_user_trst_balance_start"
  trst_user_token_balance_start=$($TRST_MAIN_CMD q bank balances $(TRST_ADDRESS) --denom $HOST_IBC_DENOM | GETBAL)
  host_user_token_balance_start=$($HOST_MAIN_CMD q bank balances $HOST_USER_ADDRESS --denom $HOST_DENOM | GETBAL)
  printf "host_user_token_balance_start (sender): %s\n" "$host_user_token_balance_start"
  # do IBC transfer
  $TRST_MAIN_CMD tx ibc-transfer transfer transfer $TRST_TRANFER_CHANNEL $HOST_USER_ADDRESS ${TRANSFER_AMOUNT}${TRST_DENOM} --from $TRST_USER -y
  $HOST_MAIN_CMD_TX ibc-transfer transfer transfer $HOST_TRANSFER_CHANNEL $(TRST_ADDRESS) ${TRANSFER_AMOUNT}${HOST_DENOM} --from $HOST_USER -y

  WAIT_FOR_BLOCK $TRST_LOGS 8

  # get new balances
  trst_user_trst_balance_end=$($TRST_MAIN_CMD q bank balances $(TRST_ADDRESS) --denom $TRST_DENOM | GETBAL)
  host_user_trst_balance_end=$($HOST_MAIN_CMD q bank balances $HOST_USER_ADDRESS --denom $IBC_TRST_DENOM | GETBAL)

  trst_user_token_balance_end=$($TRST_MAIN_CMD q bank balances $(TRST_ADDRESS) --denom $HOST_IBC_DENOM | GETBAL)
  host_user_token_balance_end=$($HOST_MAIN_CMD q bank balances $HOST_USER_ADDRESS --denom $HOST_DENOM | GETBAL)

  # get all TRST balance diffs
  trst_user_trst_balance_diff=$((trst_user_trst_balance_start - trst_user_trst_balance_end))
  host_user_trst_balance_diff=$((host_user_trst_balance_start - host_user_trst_balance_end))

  assert_equal "$trst_user_trst_balance_diff" "$TRANSFER_AMOUNT"
  assert_equal "$host_user_trst_balance_diff" "-$TRANSFER_AMOUNT"

  # get all host balance diffs
  trst_user_token_balance_diff=$((trst_user_token_balance_start - trst_user_token_balance_end))
  host_user_token_balance_diff=$((host_user_token_balance_start - host_user_token_balance_end))

  assert_equal "$trst_user_token_balance_diff" "-$TRANSFER_AMOUNT"
  assert_equal "$host_user_token_balance_diff" "5002500" #TRANSFER_AMOUNT+fee
}

@test "[INTEGRATION-BASIC-$CHAIN_NAME] AutoIbcTx MsgSend using ICA" {
  # get initial balances on host account
  receiver_token_balance_start=$($HOST_MAIN_CMD q bank balances $HOST_RECEIVER_ADDRESS --denom $HOST_DENOM | GETBAL)

  # get token balance user on TRST
  user_trst_balance_start=$($TRST_MAIN_CMD q bank balances $(TRST_ADDRESS) --denom $TRST_DENOM | GETBAL)

  # build MsgRegisterAccount and retrieve trigger ICA account
  $TRST_MAIN_CMD tx autoibctx register --connection-id connection-$CONNECTION_ID --counterparty-connection-id connection-0 --from $TRST_USER -y

  sleep 20
  ICA_ADDRESS=$($TRST_MAIN_CMD q autoibctx interchainaccounts $(TRST_ADDRESS) connection-$CONNECTION_ID)
  ICA_ADDRESS=$(echo "$ICA_ADDRESS" | awk '{print $2}')

  fund_ica=$($HOST_MAIN_CMD_TX bank send $HOST_USER_ADDRESS $ICA_ADDRESS $MSGSEND_AMOUNT$HOST_DENOM --from $HOST_USER -y)

  WAIT_FOR_BLOCK $TRST_LOGS 2

  ica_token_balance_start=$($HOST_MAIN_CMD q bank balances $ICA_ADDRESS --denom $HOST_DENOM | GETBAL)

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
      "from_address": "$ICA_ADDRESS",
      "to_address": "$HOST_RECEIVER_ADDRESS"
}
EOF

  # build MsgSubmitAutoTx with MsgSend, 60sec non-recurring
  msg_submit_auto_tx=$($TRST_MAIN_CMD tx autoibctx submit-auto-tx "$msg_send_file" --label "MsgSend using ICA" --duration "60s" --connection-id connection-$CONNECTION_ID --from $TRST_USER -y)
  echo "$msg_submit_auto_tx"

  GET_AUTO_TX_ID $(TRST_ADDRESS) 8
  
  WAIT_FOR_EXECUTED_TX_BY_ID $(TRST_ADDRESS) 50

  # calculate difference between token balance user before and after, should equal MSGSEND_AMOUNT
  ica_token_balance_end=$($HOST_MAIN_CMD q bank balances $ICA_ADDRESS --denom $HOST_DENOM | GETBAL)
  ica_diff=$(($ica_token_balance_start-$ica_token_balance_end))
  assert_equal "$ica_diff" $MSGSEND_AMOUNT

  # calculate difference between token balance receiver before and after, should equal MSGSEND_AMOUNT
  receiver_token_balance_end=$($HOST_MAIN_CMD q bank balances $HOST_RECEIVER_ADDRESS --denom $HOST_DENOM | GETBAL)
  receiver_diff=$(($receiver_token_balance_end - $receiver_token_balance_start))

  assert_equal "$receiver_diff" $MSGSEND_AMOUNT 
}

@test "[INTEGRATION-BASIC-$CHAIN_NAME] AutoIbcTx MsgSend using AuthZ" {
  # get initial balances on host account
  user_token_balance_start=$($HOST_MAIN_CMD q bank balances $HOST_USER_ADDRESS --denom $HOST_DENOM | GETBAL)
  receiver_token_balance_start=$($HOST_MAIN_CMD q bank balances $HOST_RECEIVER_ADDRESS --denom $HOST_DENOM | GETBAL)

  # get token balance user on TRST
  user_trst_balance_start=$($TRST_MAIN_CMD q bank balances $(TRST_ADDRESS) --denom $TRST_DENOM | GETBAL)

  ICA_ADDRESS=$($TRST_MAIN_CMD q autoibctx interchainaccounts $(TRST_ADDRESS) connection-$CONNECTION_ID)
  ICA_ADDRESS=$(echo "$ICA_ADDRESS" | awk '{print $2}')
  echo "ICA ADDR: $ICA_ADDRESS"

  $HOST_MAIN_CMD_TX authz grant $ICA_ADDRESS generic --msg-type "/cosmos.bank.v1beta1.MsgSend" --from $HOST_USER -y
  WAIT_FOR_BLOCK $TRST_LOGS 2
  $HOST_MAIN_CMD_TX bank send $HOST_USER_ADDRESS $ICA_ADDRESS $MSGSEND_AMOUNT$HOST_DENOM --from $HOST_USER -y

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
  msg_submit_auto_tx=$($TRST_MAIN_CMD tx autoibctx submit-auto-tx "$msg_exec_file" --label "MsgSend using AuthZ" --duration "60s" --connection-id connection-$CONNECTION_ID --from $TRST_USER -y)
  echo "$msg_submit_auto_tx"

  GET_AUTO_TX_ID $(TRST_ADDRESS) 8
  
  WAIT_FOR_EXECUTED_TX_BY_ID $(TRST_ADDRESS) 50

  # calculate difference between token balance of user before and after, should equal MSGSEND_AMOUNT
  user_token_balance_end=$($HOST_MAIN_CMD q bank balances $HOST_USER_ADDRESS --denom $HOST_DENOM | GETBAL)
  user_diff=$(($user_token_balance_start - $user_token_balance_end))
  printf "Balance start: %s\n" "$user_token_balance_start"
  printf "Balance end: %s\n" "$user_token_balance_end"
  assert_equal "$user_diff" 205000 #MsgSend to ICA and MsgSend using AuthZ + host tx fees for MsgGrant,MsgSend

  # calculate difference between token balance receiver before and after, should equal MSGSEND_AMOUNT
  receiver_token_balance_end=$($HOST_MAIN_CMD q bank balances $HOST_RECEIVER_ADDRESS --denom $HOST_DENOM | GETBAL)
  receiver_diff=$(($receiver_token_balance_end - $receiver_token_balance_start))
  printf "Balance end: %s\n" "$receiver_token_balance_end"
  assert_equal "$receiver_diff" $MSGSEND_AMOUNT #from MsgSend
}

# test auto-tx MsgSend from ICS20 message with Trigger Address ICA Account with MsgSubmitAutoTx ICA_ADDR parsing
@test "[INTEGRATION-BASIC-$CHAIN_NAME] ibc ics20 transfer, create trigger and auto-parse address" {

  # get initial balances
  user_token_balance_start=$($HOST_MAIN_CMD q bank balances $HOST_USER_ADDRESS --denom $HOST_DENOM | GETBAL)
  receiver_token_balance_start=$($HOST_MAIN_CMD q bank balances $HOST_RECEIVER_ADDRESS --denom $HOST_DENOM | GETBAL)

  # do IBC transfer
  memo='{"auto_tx": {"msgs": [{"@type": "/cosmos.bank.v1beta1.MsgSend","amount": [{"amount": "'$MSGSEND_AMOUNT'","denom": "'$HOST_DENOM'"}],"from_address":"ICA_ADDR","to_address": "'$HOST_RECEIVER_ADDRESS'"}],"duration":"2880h","interval":"60s","label":"MsgSend using ICS20 hook","cid":"connection-'$CONNECTION_ID'","start_at":"0", "owner": "'$(TRST_ADDRESS)'" }}'
  $HOST_MAIN_CMD_TX ibc-transfer transfer transfer $HOST_TRANSFER_CHANNEL $(TRST_ADDRESS) ${ICS20HOOK_AMOUNT}${IBC_TRST_DENOM} --memo "$memo" --from $HOST_USER -y

  GET_AUTO_TX_ID $(TRST_ADDRESS) 8

  ICA_ADDRESS=$($TRST_MAIN_CMD q autoibctx interchainaccounts $(TRST_ADDRESS) connection-$CONNECTION_ID)
  ICA_ADDRESS=$(echo "$ICA_ADDRESS" | awk '{print $2}')
  $HOST_MAIN_CMD_TX bank send $HOST_USER_ADDRESS $ICA_ADDRESS $MSGSEND_AMOUNT_TOTAL$HOST_DENOM --from $HOST_USER -y
  
  WAIT_FOR_EXECUTED_TX_BY_ID $(TRST_ADDRESS) 50

  # calculate difference between token balance of host user before and after, should equal 2xMSGSEND_AMOUNT
  user_token_balance_end=$($HOST_MAIN_CMD q bank balances $HOST_USER_ADDRESS --denom $HOST_DENOM | GETBAL)
  user_diff=$(($user_token_balance_start - $user_token_balance_end))
  printf "Balance start: %s\n" "$user_token_balance_start"
  printf "Balance end: %s\n" "$user_token_balance_end"
  assert_equal "$user_diff" 17280005000 #MSGSEND_AMOUNT_TOTAL for all executions(10000)+2x host tx fee(2500)

  # calculate difference between token balance receiver before and after, should equal 1xMSGSEND_AMOUNT
  receiver_token_balance_end=$($HOST_MAIN_CMD q bank balances $HOST_RECEIVER_ADDRESS --denom $HOST_DENOM | GETBAL)
  receiver_diff=$(($receiver_token_balance_end - $receiver_token_balance_start))
  printf "Balance end: %s\n" "$receiver_token_balance_end"
  assert_equal "$receiver_diff" 100000 #one MsgSend received

}

# TODO
# test auto-tx MsgSend with user address over AuthZ grant to ICA Account with MsgRegisterAccountAndSubmitAutoTx ICA_ADDR parsing
# test auto-tx with other msgs like MsgWithdrawalRewards, MsgSubmitProposal
