rm -rf ~/.trst


trstd init mainnode --chain-id=trst


trstd keys add user1 --keyring-backend test
trstd keys add user2 --keyring-backend test
trstd keys add user3 --keyring-backend test
trstd keys add user4 --keyring-backend test
trstd keys add faucet --keyring-backend test

trstd add-genesis-account $(trstd keys show user1 -a --keyring-backend test) 10000trst,100000000stake 
trstd add-genesis-account $(trstd keys show user2 -a --keyring-backend test) 5000trst
trstd add-genesis-account $(trstd keys show user3 -a --keyring-backend test) 5000trst
trstd add-genesis-account $(trstd keys show user4 -a --keyring-backend test) 5000trst
trstd add-genesis-account $(trstd keys show faucet -a --keyring-backend test) 1000000000trst

trstd gentx user1 100000000stake --chain-id=trst --keyring-backend=test  --website="trustlesshub.com" --security-contact="trustlesshub@gmail.com"

trstd init-enclave --reset
PUBLIC_KEY=$(trstd parse attestation_cert.der 2> /dev/null | cut -c 3-)
echo $PUBLIC_KEY
trstd init-bootstrap ./node-master-cert.der ./io-master-cert.der
echo "Collecting genesis txs..."
trstd collect-gentxs

echo "Validating genesis file..."
trstd validate-genesis

sed -i '104s/enable = false/enable = true/g' ~/.trst/config/app.toml
sed -i 's/enabled-unsafe-cors = false/enabled-unsafe-cors = true/g' ~/.trst/config/app.toml
sed -i 's/minimum-gas-prices = 0.0025trst"/minimum-gas-prices = 0.0000025trst"/g' ~/.trst/config/app.toml
sed -i 's/cors_allowed_origins = []/cors_allowed_origins = ["*"]/g' ~/.trst/config/config.toml
trstd start --bootstrap

