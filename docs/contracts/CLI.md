---
order: 4
title: CLI
description: Using the CLI to instantiate an AutoExecution contracts
---

# Recurring Swap contract using the Command-line interface (CLI)

Perform a front-running resistant recurring swap using the CLI!
With this example we hope to give you a good understanding of how an AutoExecution contract works from the CLI. 

Several contracts are at play in this example
1. TIP20 token contract of input token 
2. RecurringSwap proxy contract to be instantiated
3. SwapPair contract that swaps tokens

### CodeHash

For each of these contracts you specify the CodeHash - a hash of the code - which is used to encrypt the messages and decrypt them in the secure enviroment. Code hashes are what binds a transaction to the specific contract code. Without this, it would be possible to perform a replay attack in a forked chain with a modified contract that decrypts and prints the input message. 

### Now, let's get into it!

First, upload the recurring swap contract.
This is part of the (to-be) open sourced DeFi contract bundle, which will be uploaded on [TRST Labs's Github](https://github.com/trstlabs/)

```bash
export token_addr=trust1hqrdl6wstt8qzshwc6mrumpjk9338k0l2gn4jw
export token_code_hash=642F164CDAC291AB44DBD1AB1A7F5CD946775C05351231D962B333FA0F5A2332
export swap_pair_addr=trust1xqeym28j9xgv0p93pwwt6qcxf9tdvf9zyszued
export swap_pair_code_hash=1E49EF2E8333B1776982147C271CAD7A52DFA1FCECEFC1801E00BADE3F927492

export deployer_address=$(trstd keys show -a user2  --keyring-backend test)
echo "Deployer address: '$deployer_address'"

export chain_id=trst-dev-1
export wasm_path=../build

trstd tx compute store "${wasm_path}/recurring_swap.wasm" --from user1 --keyring-backend test --chain-id "$chain_id" --gas 3000000 --fees 500utrst  -b block -y  --contract-title 'Recurring Swap - TIP20 extention' --contract-description 'This contract code is for recurring swaps. It is an extention to TIP20 contract code, and can be used in other contract codes. Contract code available at github.com/trstlabs/'  --duration 70s --interval 60s
recurring_swap_code_id=$(trstd query compute list-codes | jq '.[-1]."id"')
recurring_swap_code_hash=$(trstd query compute list-codes | jq '.[-1]."code_hash"')
echo "Stored recurring swap code: '$recurring_swap_code_id', '$recurring_swap_code_hash'"

```

Then set up the messages to pass. These are parsed as base64 messages.

* Msg: Message to send to the TIP20 token contract
* SwapMsg: Message that gets sent from TIP20 token contract to the SwapRouter contract
* AutoMsg: Message that triggers the AutoExecution

```bash
msg='{"owner": "'$deployer_address'","offer_asset": {"info":{"token":{"contract_addr":"'$token_addr'","token_code_hash":"'$token_code_hash'","viewing_key":""}}, "amount": "50000"}, "swap_pair_addr":"'$swap_pair_addr'","swap_pair_code_hash":"'$swap_pair_code_hash'"}'
msg_to_pass="$(base64 --wrap=0 <<<"$msg")"

auto_msg='{"auto_msg":{}}'
auto_msg_to_pass="$(base64 --wrap=0 <<<"$auto_msg")"
```

Then, with only one message the recurring swap is instantiated from the token contract. As it is instantiated from the TIP20 token contract, the token contract is the owner. This perserves privacy and makes it possible to give allowance right away through the [CosmWasm reply](https://docs.cosmwasm.com/docs/1.0/smart-contracts/message/submessage/) function in the TIP20 token.


```bash
export TX_HASH=$( trstd tx compute execute $(echo "$token_addr" | tr -d '"') '{"instantiate_with_allowance" : { "max_allowance": "10000000", "code_id":'$recurring_swap_code_id', "code_hash":'$recurring_swap_code_hash',  "duration":"200s","msg":"'$msg_to_pass'", "auto_msg":"'$auto_msg_to_pass'", "interval":"30s", "contract_id":"'$iter' DCA2"}}' -b block -y --amount 20000000utrst --from $deployer_name --keyring-backend test --chain-id "$chain_id"  --gas 3000000 --fees 500utrst -y -b block | 
 jq -r .txhash
)

function wait_for_tx() {
  until (trstd q tx "$1"); do
      sleep 5
  done
}

wait_for_tx "$TX_HASH" "Waiting for tx to finish on-chain..."
trstd q compute tx $TX_HASH

export recurring_swap_contract=$(trstd query compute list-contracts-by-code $recurring_swap_code_id | jq '.[-1].address')

echo Recurring Swap: "$recurring_swap_contract" | tr -d '"'
```

For this test, we used the following AutoExecution parameters:

* Interval: 30s
* Duration: 200s

No StartDurationAt specified means we will start execution in 30s.
Want to also swap right away? There's an optional message to send called "send_msg" which can be used to perform right away. Parse it like so:


```bash
swap_msg='{"swap":{}}'
swap_msg_to_pass="$(base64 --wrap=0 <<<"$swap_msg")"
send_msg='{"amount": "50000","recipient": "'$swap_pair_addr'", "recipient_code_hash" : "'$swap_pair_code_hash'", "msg": "'$swap_msg_to_pass'"}'
send_msg_to_pass="$(base64 --wrap=0 <<<"$send_msg")"
```

You should change the instantiate_with_allowance message to the following.

```bash
export TX_HASH=$( trstd tx compute execute $(echo "$token_addr" | tr -d '"') '{"instantiate_with_allowance" : { "max_allowance": "10000000", "code_id":'$recurring_swap_code_id', "code_hash":'$recurring_swap_code_hash',  "duration":"200s","msg":"'$msg_to_pass'", "auto_msg":"'$auto_msg_to_pass'", "send_msg":"'$send_msg_to_pass'","interval":"30s", "contract_id":"'$iter' DCA2"}}' -b block -y --amount 20000000utrst --from $deployer_name --keyring-backend test --chain-id "$chain_id"  --gas 3000000 --fees 500utrst -y -b block | 
 jq -r .txhash
)

wait_for_tx "$TX_HASH" "Waiting for tx to finish on-chain..."
trstd q compute tx $TX_HASH

export recurring_swap_contract=$(trstd query compute list-contracts-by-code $recurring_swap_code_id | jq '.[-1].address')

echo Recurring Swap: "$recurring_swap_contract" | tr -d '"'
```