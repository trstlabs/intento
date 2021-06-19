rm -rf ~/.tpp


tppd init mainnode --chain-id=tpp


tppd keys add user1 --keyring-backend test
tppd keys add user2 --keyring-backend test
tppd keys add user3 --keyring-backend test
tppd keys add user4 --keyring-backend test
tppd keys add faucet --keyring-backend test

tppd add-genesis-account $(tppd keys show user1 -a --keyring-backend test) 1000tpp,100000000stake 
tppd add-genesis-account $(tppd keys show user2 -a --keyring-backend test) 500tpp
tppd add-genesis-account $(tppd keys show user3 -a --keyring-backend test) 500tpp
tppd add-genesis-account $(tppd keys show user4 -a --keyring-backend test) 500tpp
tppd add-genesis-account $(tppd keys show faucet -a --keyring-backend test) 100000000tpp

tppd gentx user1 100000000stake --chain-id=tpp --keyring-backend=test  --website="trustpriceprotocol.com" --security-contact="trustpriceprotocol@gmail.com"

echo "Collecting genesis txs..."
tppd collect-gentxs

echo "Validating genesis file..."
tppd validate-genesis

tppd start --bootstrap