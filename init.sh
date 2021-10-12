rm -rf ~/.tpp


tppd init mainnode --chain-id=tpp


tppd keys add user1 --keyring-backend test
tppd keys add user2 --keyring-backend test
tppd keys add user3 --keyring-backend test
tppd keys add user4 --keyring-backend test
tppd keys add faucet --keyring-backend test

tppd add-genesis-account $(tppd keys show user1 -a --keyring-backend test) 10000tpp,100000000stake 
tppd add-genesis-account $(tppd keys show user2 -a --keyring-backend test) 5000tpp
tppd add-genesis-account $(tppd keys show user3 -a --keyring-backend test) 5000tpp
tppd add-genesis-account $(tppd keys show user4 -a --keyring-backend test) 5000tpp
tppd add-genesis-account $(tppd keys show faucet -a --keyring-backend test) 1000000000tpp

tppd gentx user1 100000000stake --chain-id=tpp --keyring-backend=test  --website="trustpriceprotocol.com" --security-contact="trustpriceprotocol@gmail.com"

tppd init-enclave
PUBLIC_KEY=$(tppd parse attestation_cert.der 2> /dev/null | cut -c 3-)
echo $PUBLIC_KEY
tppd init-bootstrap ./node-master-cert.der ./io-master-cert.der
echo "Collecting genesis txs..."
tppd collect-gentxs

echo "Validating genesis file..."
tppd validate-genesis

sed -i '104s/enable = false/enable = true/g' ~/.tpp/config/app.toml
sed -i 's/enabled-unsafe-cors = false/enabled-unsafe-cors = true/g' ~/.tpp/config/app.toml
sed -i 's/minimum-gas-prices = 0.0025tpp"/minimum-gas-prices = 0.0000025tpp"/g' ~/.tpp/config/app.toml

tppd start --bootstrap

