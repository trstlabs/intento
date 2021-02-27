rm -rf ~/.tppd


tppd init test --chain-id=tppd

tppd config output json
tppd config indent true
tppd config trust-node true
tppd config chain-id trustitemstest
tppd config keyring-backend test

tppd keys add user1
tppd keys add user2
tppd keys add user3

tppd add-genesis-account $(tppd keys show user1 -a) 1000token,100000000stake
tppd add-genesis-account $(tppd keys show user2 -a) 500token
tppd add-genesis-account $(tppd keys show user3 -a) 500token

tppd gentx --name user1 --keyring-backend test

echo "Collecting genesis txs..."
tppd collect-gentxs

echo "Validating genesis file..."
tppd validate-genesis