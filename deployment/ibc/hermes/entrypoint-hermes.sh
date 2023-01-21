#! /bin/bash

echo "veryfing balances"
hermes keys balance --chain trstdev-1
hermes keys balance --chain trstdev-2
hermes --config /hermes-user/alternative-config.toml keys balance --chain trstdev-1
hermes --config /hermes-user/alternative-config.toml keys balance --chain trstdev-2

# echo "creating channel"
# hermes create channel --a-chain trstdev-1 --b-chain trstdev-2 --a-port transfer --b-port transfer --new-client-connection --yes

echo "creating connection"
hermes create connection --a-chain trstdev-1 --b-chain trstdev-2

echo "relaying forever"
hermes start
