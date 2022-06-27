
rm -rf ~/.trst
kill -9 $(lsof -t -i:26657 -sTCP:LISTEN)
kill -9 $(lsof -t -i:1317 -sTCP:LISTEN)

trstd init FRST --chain-id=trst_chain_1

trstd prepare-genesis testnet trst_chain_1
yes exchange cabin middle shed identify soon loop vivid mutual simple sing vessel tail embody vote glide bid olive possible invite merry kitten keen nuclear | trstd keys add user1 --keyring-backend test --recover
yes impulse north bulb pistol oven fiction struggle gun season quote blush region fly sight glory glory brisk wash gate soon toddler person shield above| trstd keys add user2 --keyring-backend test --recover

echo "Adding airdrop accounts..."
trstd import-genesis-accounts-from-snapshot ./snapshot.json ./reserves.json 
echo "Getting genesis time..."


trstd add-genesis-account $(trstd keys show user1 -a --keyring-backend test) 8750000000000utrst
trstd add-genesis-account $(trstd keys show user2 -a --keyring-backend test) 8750000000000utrst

trstd gentx user1 750000000000utrst --chain-id=trst_chain_1 --keyring-backend=test  --website="trustlesshub.com" --security-contact="info@trstlabs.xyz"

echo "Collecting genesis txs..."
trstd collect-gentxs

echo "Validating genesis file..."
trstd validate-genesis

sed -i '129s/enabled-unsafe-cors = false/enabled-unsafe-cors = true/g' ~/.trst/config/app.toml


trstd start --bootstrap > init.log --log_level info


trstd tx staking delegate trustvaloper16rpg3wwxsrggxv34hj2ca5xa2gxy4jgsv3df03 200utrst  --from user2 --keyring-backend=test --fees 500utrst --chain-id=trst_chain_1
trstd tx staking delegate trustvaloper16rpg3wwxsrggxv34hj2ca5xa2gxy4jgsv3df03 33000000utrst  --from user2 --keyring-backend=test --fees 500utrst --chain-id=trst_chain_1 -y

trstd tx bank send  $(trstd keys show user2 -a --keyring-backend test) trust17xpfvakm2amg962yls6f84z3kell8c5leya4l7 69utrst --from user2 --keyring-backend=test --fees 500utrst --chain-id=trst_chain_1 -y



trstd tx bank send  $(trstd keys show user2 -a --keyring-backend test) trust1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8wz6fau 69utrst --from user2 --keyring-backend=test --fees 500utrst --chain-id=trst_chain_1 -y