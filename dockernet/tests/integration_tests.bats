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

  HOST_USR_ADDRESS=$(${CHAIN_NAME}_ADDRESS)
  HOST_RECEIVER_ADDRESS=$(GET_VAR_VALUE ${CHAIN_NAME}_RECEIVER_ADDRESS)

  TRST_USR=$(${CHAIN_NAME}_USR_ACCT)
  # HOST_USR="$(GET_VAR_VALUE ${CHAIN_NAME}_VAL_PREFIX)1"
  # TRST_VAL=${TRST_VAL_PREFIX}1

  TRST_TRANFER_CHANNEL="channel-${TRANSFER_CHANNEL_NUMBER}"
  HOST_TRANSFER_CHANNEL="channel-0"

  TRANSFER_AMOUNT=5000000
  MSGSEND_AMOUNT=1000000
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

# ##############################################################################################
# ######                              SETUP TESTS                                         ######
# ##############################################################################################
# # confirm host zone is registered
# @test "[INTEGRATION-BASIC-$CHAIN_NAME] host zones successfully registered" {
#   run $TRST_MAIN_CMD q stakeibc show-host-zone $HOST_CHAIN_ID
#   assert_line "  host_denom: $HOST_DENOM"
#   assert_line "  chain_id: $HOST_CHAIN_ID"
#   assert_line "  transfer_channel_id: channel-$TRANSFER_CHANNEL_NUMBER"
#   refute_line '  delegation_account: null'
#   refute_line '  fee_account: null'
#   refute_line '  redemption_account: null'
#   refute_line '  withdrawal_account: null'
#   assert_line '  unbonding_frequency: "1"'
# }

##############################################################################################
######                TEST BASIC TRST FUNCTIONALITY                                   ######
##############################################################################################

@test "[INTEGRATION-BASIC-$CHAIN_NAME] ibc transfer updates all balances" {
  # get initial balances
  sval_trst_balance_start=$($TRST_MAIN_CMD q bank balances $(TRST_ADDRESS) --denom $TRST_DENOM | GETBAL)
  hval_trst_balance_start=$($HOST_MAIN_CMD q bank balances $HOST_USR_ADDRESS --denom $IBC_TRST_DENOM | GETBAL)
  sval_token_balance_start=$($TRST_MAIN_CMD q bank balances $(TRST_ADDRESS) --denom $HOST_IBC_DENOM | GETBAL)
  hval_token_balance_start=$($HOST_MAIN_CMD q bank balances $HOST_USR_ADDRESS --denom $HOST_DENOM | GETBAL)

  # do IBC transfer
  $TRST_MAIN_CMD tx ibc-transfer transfer transfer $TRST_TRANFER_CHANNEL $HOST_USR_ADDRESS ${TRANSFER_AMOUNT}${TRST_DENOM} --from $TRST_VAL -y
  $HOST_MAIN_CMD tx ibc-transfer transfer transfer $HOST_TRANSFER_CHANNEL $(TRST_ADDRESS) ${TRANSFER_AMOUNT}${HOST_DENOM} --from $HOST_USR -y

  WAIT_FOR_BLOCK $TRST_LOGS 8

  # get new balances
  sval_trst_balance_end=$($TRST_MAIN_CMD q bank balances $(TRST_ADDRESS) --denom $TRST_DENOM | GETBAL)
  hval_trst_balance_end=$($HOST_MAIN_CMD q bank balances $HOST_USR_ADDRESS --denom $IBC_TRST_DENOM | GETBAL)
  sval_token_balance_end=$($TRST_MAIN_CMD q bank balances $(TRST_ADDRESS) --denom $HOST_IBC_DENOM | GETBAL)
  hval_token_balance_end=$($HOST_MAIN_CMD q bank balances $HOST_USR_ADDRESS --denom $HOST_DENOM | GETBAL)

  # get all TRST balance diffs
  sval_trst_balance_diff=$(($sval_trst_balance_start - $sval_trst_balance_end))
  hval_trst_balance_diff=$(($hval_trst_balance_start - $hval_trst_balance_end))
  assert_equal "$sval_trst_balance_diff" "$TRANSFER_AMOUNT"
  assert_equal "$hval_trst_balance_diff" "-$TRANSFER_AMOUNT"

  # get all host balance diffs
  sval_token_balance_diff=$(($sval_token_balance_start - $sval_token_balance_end))
  hval_token_balance_diff=$(($hval_token_balance_start - $hval_token_balance_end))
  assert_equal "$sval_token_balance_diff" "-$TRANSFER_AMOUNT"
  assert_equal "$hval_token_balance_diff" "$TRANSFER_AMOUNT"
}

# test auto-tx MsgSend with Trigger Address ICA Account
# get token balance user1 on $CHAIN_NAME
# get token balance user2 on $CHAIN_NAME
# get token balance user1 on TRST
# build MsgRegisterAccount
# build MsgSend with MSGSEND_AMOUNT
# build MsgSubmitAutoTx with MsgSend, 60sec non-recurring

# WAIT_FOR_AUTO_TX_TRIGGER wait for 60sec
# calculate difference between token balance user1 before and after, should equal MSGSEND_AMOUNT
# calculate difference between token balance user2 before and after, should equal MSGSEND_AMOUNT

@test "[INTEGRATION-BASIC-$CHAIN_NAME] liquid stake mint and transfer" {
  # get initial balances on host account
  user_token_balance_start=$($HOST_MAIN_CMD q bank balances $(HOST_USR_ADDRESS) --denom $HOST_DENOM | GETBAL)
  receiver_token_balance_start=$($HOST_MAIN_CMD q bank balances $(HOST_RECEIVER_ADDRESS) --denom $HOST_DENOM | GETBAL)

  # get token balance user on TRST
  user_trst_balance_start=$($TRST_MAIN_CMD q bank balances $(TRST_ADDRESS) --denom $TRST_DENOM | GETBAL)

  # build MsgRegisterAccount and retrieve trigger ICA account
  trigger_ica_account=$($TRST_MAIN_CMD tx autoibctx register --connection-id $CONNECTION_ID --from $TRST_USR | grep interchain_account_address)
  echo "$trigger_ica_account"

  grant=$($HOST_MAIN_CMD tx authz grant $trigger_ica_account generic --msg-type /cosmos.bank.v1beta1.MsgSend --from $(HOST_USR_ADDRESS))
  echo "$grant"

  # build MsgSend with MSGSEND_AMOUNT
  # msg_send='{
  #   "@type":"/cosmos.bank.v1beta1.MsgSend",
  #   "amount": [{
  #       "amount": "$MSGSEND_AMOUNT",
  #       "denom": "utrst"
  #   }],
  #   "from_address": "$(TRST_ADDRESS)",
  #   "to_address": "$(HOST_RECEIVER_ADDRESS)"
  # }'

  msg_exec="{
    "@type":"/cosmos.authz.v1beta1.MsgExec",
    "msgs": [{
      "@type":"/cosmos.bank.v1beta1.MsgSend",
        "amount": [{
         "amount": "$MSGSEND_AMOUNT",
         "denom": "utrst"
        }],
      "from_address": "$(TRST_ADDRESS)",
      "to_address": "$(HOST_RECEIVER_ADDRESS)"
    }],
    "grantee": "$trigger_ica_account"
  }"
  echo "$msg_exec"
  #"$HOST_MAIN_CMD tx bank send $(TRST_ADDRESS) $(HOST_RECEIVER_ADDRESS) $MSGSEND_AMOUNT$TRST_DENOM--from $TRST_USR"

  # build MsgSubmitAutoTx with MsgSend, 60sec non-recurring
  msg_submit_auto_tx=$($TRST_MAIN_CMD tx autoibctx submit-auto-tx $msg_send --label "test" --duration "60s" --connection-id $CONNECTION_ID --from $TRST_USR)
  echo "$msg_submit_auto_tx"
  # WAIT_FOR_AUTO_TX_TRIGGER wait for 30blocks
  WAIT_FOR_AUTO_TX

  # calculate difference between token balance receiver before and after, should equal MSGSEND_AMOUNT
  balance_end=$($HOST_MAIN_CMD q bank balances $(HOST_RECEIVER_ADDRESS) --denom $HOST_DENOM | GETBAL)
  diff=$(($balance_end - $receiver_token_balance_start))
  assert_equal "$diff" $MSGSEND_AMOUNT
}

# TODO
# test auto-tx MsgSend with user address over AuthZ grant to ICA Account

# test auto-tx MsgSend with user address over AuthZ grant to ICA Account with MsgRegisterAccountAndSubmitAutoTx ICA_ADDR parsing

# test auto-tx MsgSend from ICS20 message with Trigger Address ICA Account

# test auto-tx MsgSend from ICS20 message with user address over AuthZ grant to ICA Account


# test auto-tx MsgWithdrawalRewards, MsgSubmitProposal

# @test "[INTEGRATION-BASIC-$CHAIN_NAME] liquid stake mint and transfer" {
#   # get initial balances on trst account
#   token_balance_start=$($TRST_MAIN_CMD q bank balances $(TRST_ADDRESS) --denom $HOST_IBC_DENOM | GETBAL)
#   sttoken_balance_start=$($TRST_MAIN_CMD q bank balances $(TRST_ADDRESS) --denom st$HOST_DENOM | GETBAL)

#   user_address=$($HOST_MAIN_CMD keys show ${HOST_USER_KEY_NAME} --keyring-backend test -a)

#   # get initial ICA accound balance
#   delegation_address=$(GET_ICA_ADDR $HOST_CHAIN_ID delegation)
#   delegation_ica_balance_start=$($HOST_MAIN_CMD q bank balances $delegation_address --denom $HOST_DENOM | GETBAL)

#   # liquid stake
#   $TRST_MAIN_CMD tx stakeibc liquid-stake $MSGSEND_AMOUNT $HOST_DENOM --from $TRST_VAL -y

#   # wait for the stTokens to get minted
#   WAIT_FOR_BALANCE_CHANGE TRST $(TRST_ADDRESS) st$HOST_DENOM

#   # make sure IBC_DENOM went down
#   token_balance_end=$($TRST_MAIN_CMD q bank balances $(TRST_ADDRESS) --denom $HOST_IBC_DENOM | GETBAL)
#   token_balance_diff=$(($token_balance_start - $token_balance_end))
#   assert_equal "$token_balance_diff" $MSGSEND_AMOUNT

#   # make sure stToken went up
#   sttoken_balance_end=$($TRST_MAIN_CMD q bank balances $(TRST_ADDRESS) --denom st$HOST_DENOM | GETBAL)
#   sttoken_balance_diff=$(($sttoken_balance_end - $sttoken_balance_start))
#   assert_equal "$sttoken_balance_diff" $MSGSEND_AMOUNT

#   # Wait for the transfer to complete
#   WAIT_FOR_BALANCE_CHANGE $CHAIN_NAME $delegation_address $HOST_DENOM

#   # get the new delegation ICA balance
#   delegation_ica_balance_end=$($HOST_MAIN_CMD q bank balances $delegation_address --denom $HOST_DENOM | GETBAL)
#   diff=$(($delegation_ica_balance_end - $delegation_ica_balance_start))
#   assert_equal "$diff" $MSGSEND_AMOUNT
# }

# @test "[INTEGRATION-BASIC-$CHAIN_NAME] packet forwarding automatically liquid stakes" {
#   memo='{ "autopilot": { "receiver": "'"$(TRST_ADDRESS)"'",  "stakeibc": { "action": "LiquidStake" } } }'

#   # get initial balances
#   sttoken_balance_start=$($TRST_MAIN_CMD q bank balances $(TRST_ADDRESS) --denom st$HOST_DENOM | GETBAL)

#   # Send the IBC transfer with the JSON memo
#   transfer_msg_prefix="$HOST_MAIN_CMD tx ibc-transfer transfer transfer $HOST_TRANSFER_CHANNEL"
#   if [[ "$CHAIN_NAME" == "GAIA" ]]; then
#     # For GAIA (ibc-v3), pass the memo into the receiver field
#     $transfer_msg_prefix "$memo" ${PACKET_FORWARD_MSGSEND_AMOUNT}${HOST_DENOM} --from $HOST_USR -y
#   elif [[ "$CHAIN_NAME" == "HOST" ]]; then
#     # For HOST (ibc-v5), pass an address for a receiver and the memo in the --memo field
#     $transfer_msg_prefix $(TRST_ADDRESS) ${PACKET_FORWARD_MSGSEND_AMOUNT}${HOST_DENOM} --memo "$memo" --from $HOST_USR -y
#   else
#     # For all other hosts, skip this test
#     skip "Packet forward liquid stake test is only run on GAIA and HOST"
#   fi

#   # Wait for the transfer to complete
#   WAIT_FOR_BALANCE_CHANGE TRST $(TRST_ADDRESS) st$HOST_DENOM

#   # make sure stATOM balance increased
#   sttoken_balance_end=$($TRST_MAIN_CMD q bank balances $(TRST_ADDRESS) --denom st$HOST_DENOM | GETBAL)
#   sttoken_balance_diff=$(($sttoken_balance_end - $sttoken_balance_start))
#   assert_equal "$sttoken_balance_diff" "$PACKET_FORWARD_MSGSEND_AMOUNT"
# }

# # check that tokens on the host are staked
# @test "[INTEGRATION-BASIC-$CHAIN_NAME] tokens on $CHAIN_NAME were staked" {
#   # wait for another epoch to pass so that tokens are staked
#   WAIT_FOR_STRING $TRST_LOGS "\[DELEGATION\] success on $HOST_CHAIN_ID"
#   WAIT_FOR_BLOCK $TRST_LOGS 4

#   # check staked tokens
#   NEW_STAKE=$($HOST_MAIN_CMD q staking delegation $(GET_ICA_ADDR $HOST_CHAIN_ID delegation) $(GET_VAL_ADDR $CHAIN_NAME 1) | GETSTAKE)
#   stake_diff=$(($NEW_STAKE > 0))
#   assert_equal "$stake_diff" "1"
# }

# # check that redemptions and claims work
# @test "[INTEGRATION-BASIC-$CHAIN_NAME] redemption works" {
#   # get initial balance of redemption ICA
#   redemption_ica_balance_start=$($HOST_MAIN_CMD q bank balances $(GET_ICA_ADDR $HOST_CHAIN_ID redemption) --denom $HOST_DENOM | GETBAL)

#   # call redeem-stake
#   $TRST_MAIN_CMD tx stakeibc redeem-stake $REDEEM_AMOUNT $HOST_CHAIN_ID $HOST_RECEIVER_ADDRESS \
#     --from $TRST_VAL --keyring-backend test --chain-id $TRST_CHAIN_ID -y

#   WAIT_FOR_STRING $TRST_LOGS "\[REDEMPTION] completed on $HOST_CHAIN_ID"
#   WAIT_FOR_BLOCK $TRST_LOGS 2

#   # check that the tokens were transferred to the redemption account
#   redemption_ica_balance_end=$($HOST_MAIN_CMD q bank balances $(GET_ICA_ADDR $HOST_CHAIN_ID redemption) --denom $HOST_DENOM | GETBAL)
#   diff_positive=$(($redemption_ica_balance_end > $redemption_ica_balance_start))
#   assert_equal "$diff_positive" "1"
# }

# @test "[INTEGRATION-BASIC-$CHAIN_NAME] claimed tokens are properly distributed" {
#   # get balance before claim
#   start_balance=$($HOST_MAIN_CMD q bank balances $HOST_RECEIVER_ADDRESS --denom $HOST_DENOM | GETBAL)

#   # grab the epoch number for the first deposit record in the list od DRs
#   EPOCH=$($TRST_MAIN_CMD q records list-user-redemption-record | grep -Fiw 'epoch_number' | head -n 1 | grep -o -E '[0-9]+')

#   # claim the record (send to trst address)
#   $TRST_MAIN_CMD tx stakeibc claim-undelegated-tokens $HOST_CHAIN_ID $EPOCH $(TRST_ADDRESS) \
#     --from $TRST_VAL --keyring-backend test --chain-id $TRST_CHAIN_ID -y

#   WAIT_FOR_STRING $TRST_LOGS "\[CLAIM\] success on $HOST_CHAIN_ID"
#   WAIT_FOR_BLOCK $TRST_LOGS 2

#   # check that the tokens were transferred to the sender account
#   end_balance=$($HOST_MAIN_CMD q bank balances $HOST_RECEIVER_ADDRESS --denom $HOST_DENOM | GETBAL)

#   # check that the undelegated tokens were transfered to the sender account
#   diff_positive=$(($end_balance > $start_balance))
#   assert_equal "$diff_positive" "1"
# }

# # check that a second liquid staking call kicks off reinvestment
# @test "[INTEGRATION-BASIC-$CHAIN_NAME] rewards are being reinvested, exchange rate updating" {
#   # check that the exchange rate has increased (i.e. redemption rate is greater than 1)
#   MULT=1000000
#   redemption_rate=$($TRST_MAIN_CMD q stakeibc show-host-zone $HOST_CHAIN_ID | grep -Fiw 'redemption_rate' | grep -Eo '[+-]?[0-9]+([.][0-9]+)?')
#   redemption_rate_increased=$(($(FLOOR $(DECMUL $redemption_rate $MULT)) > $(FLOOR $(DECMUL 1.00000000000000000 $MULT))))
#   assert_equal "$redemption_rate_increased" "1"
# }

# # rewards have been collected and distributed to trst stakers
# @test "[INTEGRATION-BASIC-$CHAIN_NAME] rewards are being distributed to stakers" {
#   # collect the 2nd validator's outstanding rewards
#   val_address=$($TRST_MAIN_CMD keys show ${TRST_VAL_PREFIX}2 --keyring-backend test -a)
#   $TRST_MAIN_CMD tx distribution withdraw-all-rewards --from ${TRST_VAL_PREFIX}2 -y
#   WAIT_FOR_BLOCK $TRST_LOGS 2

#   # confirm they've recieved stTokens
#   sttoken_balance=$($TRST_MAIN_CMD q bank balances $val_address --denom st$HOST_DENOM | GETBAL)
#   rewards_accumulated=$(($sttoken_balance > 0))
#   assert_equal "$rewards_accumulated" "1"
# }
