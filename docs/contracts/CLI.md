---
order: 6
title: CLI
description: FUsing the Command-line interface (CLI)
---

# Using the Command-line interface (CLI)


Perform a privacy-perserving recurring swap using the CLI!


First, upload the recurring swap contract.

```bash

set -xe


function wait_for_tx() {
  until (trstd q tx "$1"); do
      sleep 5
  done
}

export token_addr=trust1hqrdl6wstt8qzshwc6mrumpjk9338k0l2gn4jw
export token_code_hash=642F164CDAC291AB44DBD1AB1A7F5CD946775C05351231D962B333FA0F5A2332
export swap_pair_addr=trust1xqeym28j9xgv0p93pwwt6qcxf9tdvf9zyszued
export swap_pair_code_hash=1E49EF2E8333B1776982147C271CAD7A52DFA1FCECEFC1801E00BADE3F927492

export iter=2
export deployer_name=user2
export deployer_address=$(trstd keys show -a $deployer_name  --keyring-backend test)
echo "Deployer address: '$deployer_address'"
export test_name=user3
export test_address=$(trstd keys show -a $test_name  --keyring-backend test)

export chain_id=trst_chain_1
export wasm_path=../build

trstd tx compute store "${wasm_path}/recurring_swap.wasm" --from "$deployer_name" --keyring-backend test --chain-id "$chain_id" --gas 3000000 --fees 500utrst  -b block -y  --contract-title 'Recurring Swap - TIP20 extention' --contract-description 'This contract code is for recurring swaps. It is an extention to TIP20 contract code, and can be used in other contract codes. Contract code available at github.com/trstlabs/'  --duration 70s --interval 60s
recurring_swap_code_id=$(trstd query compute list-codes | jq '.[-1]."id"')
recurring_swap_code_hash=$(trstd query compute list-codes | jq '.[-1]."code_hash"')
echo "Stored recurring swap code: '$recurring_swap_code_id', '$recurring_swap_code_hash'"

```

Then set up the messages to pass. 

```bash


msg='{"owner": "'$deployer_address'","offer_asset": {"info":{"token":{"contract_addr":"'$token_addr'","token_code_hash":"'$token_code_hash'","viewing_key":""}}, "amount": "50000"}, "swap_pair_addr":"'$swap_pair_addr'","swap_pair_code_hash":"'$swap_pair_code_hash'"}'
msg_to_pass="$(base64 --wrap=0 <<<"$msg")"

swap_msg='{"swap":{}}'
swap_msg_to_pass="$(base64 --wrap=0 <<<"$swap_msg")"
send_msg='{"amount": "50000","recipient": "'$swap_pair_addr'", "recipient_code_hash" : "'$swap_pair_code_hash'", "msg": "'$swap_msg_to_pass'"}'
send_msg_to_pass="$(base64 --wrap=0 <<<"$send_msg")"

auto_msg='{"auto_msg":{}}'
auto_msg_to_pass="$(base64 --wrap=0 <<<"$auto_msg")"
```

Then, with only one message the recurring swap is instantiated from the token contract. As it is instantiated from the privacy-perserving TIP20 token contract, the token contract is the owner. This perserves privacy and makes it possible to give allowance right away through the [CosmWasm reply](https://docs.cosmwasm.com/docs/1.0/smart-contracts/message/submessage/) function in the TIP20 token.

```bash
export TX_HASH=$( trstd tx compute execute $(echo "$token_addr" | tr -d '"') '{"instantiate_with_allowance" : { "max_allowance": "10000000", "code_id":'$recurring_swap_code_id', "code_hash":'$recurring_swap_code_hash',  "duration":"200s","msg":"'$msg_to_pass'", "auto_msg":"'$auto_msg_to_pass'", "interval":"30s", "contract_id":"'$iter' DCA2"}}' -b block -y --amount 20000000utrst --from $deployer_name --keyring-backend test --chain-id "$chain_id"  --gas 3000000 --fees 500utrst -y -b block | 
 jq -r .txhash
)

wait_for_tx "$TX_HASH" "Waiting for tx to finish on-chain..."
trstd q compute tx $TX_HASH

export recurring_swap_contract=$(trstd query compute list-contracts-by-code $recurring_swap_code_id | jq '.[-1].address')

echo Recurring Swap: "$recurring_swap_contract" | tr -d '"'
    ```