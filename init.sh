rm -rf ~/.tpp


tppd init mainnode --chain-id=tpp


tppd keys add user1 --keyring-backend os --recoverclient announce finger antenna obey calm enable affair faculty gorilla shield cattle bicycle lizard mango churn deposit gather anger gown prevent outside scheme solve
tppd keys add user2 --keyring-backend os
tppd keys add user3 --keyring-backend os
tppd keys add user4 --keyring-backend os
tppd keys add faucet --keyring-backend os

tppd add-genesis-account $(tppd keys show user1 -a) 1000tpp,100000000stake  --vesting-amount=200tpp --vesting-end-time="1640257934"
tppd add-genesis-account $(tppd keys show user2 -a) 500tpp
tppd add-genesis-account $(tppd keys show user3 -a) 500tpp
tppd add-genesis-account $(tppd keys show user4 -a) 500tpp
tppd add-genesis-account $(tppd keys show faucet -a) 100000000tpp

tppd gentx user1 100000000stake --chain-id=tpp-test-1 --keyring-backend=os  --website="trustpriceprotocol.com" --security-contact="trustpriceprotocol@gmail.com"

echo "Collecting genesis txs..."
tppd collect-gentxs

echo "Validating genesis file..."
tppd validate-genesis

