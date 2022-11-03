
rm -rf ~/.trstd
mkdir $HOME/opt/trustlesshub/.sgx_secrets

CONFIG_PATH="$HOME/.trstd/config/config.toml"

kill -9 $(lsof -t -i:26657 -sTCP:LISTEN)
kill -9 $(lsof -t -i:1317 -sTCP:LxISTEN)

trstd init FRST --chain-id=trst_chain_1

trstd prepare-genesis testnet trst_chain_1
yes exchange cabin middle shed identify soon loop vivid mutual simple sing vessel tail embody vote glide bid olive possible invite merry kitten keen nuclear | trstd keys add user1 --keyring-backend test --recover

yes comic broom zone grass reject apology erupt chef wish add actor damage deputy hip aware connect addict excite poem arrive since bird couple artwork | trstd keys add user3 --keyring-backend test --recover

yes impulse north bulb pistol oven fiction struggle gun season quote blush region fly sight glory glory brisk wash gate soon toddler person shield above| trstd keys add user2 --keyring-backend test --recover

yes kiwi obtain scrub aunt female shoulder dune shove budget salt mechanic plug beef right pact economy swear flash update wild change puppy hurdle power | trstd keys add user4 --keyring-backend test --recover

yes orchard thing tooth dismiss seat couple define atom antenna language fuel wrist napkin tired undo toddler virus cherry shock mimic toss rifle predict crisp |trstd keys add faucet --keyring-backend test --recover



trstd add-genesis-account $(trstd keys show user1 -a --keyring-backend test) 8750000000000utrst
trstd add-genesis-account $(trstd keys show user2 -a --keyring-backend test) 8750000000000utrst
trstd add-genesis-account $(trstd keys show user3 -a --keyring-backend test) 8750000000000utrst
trstd add-genesis-account $(trstd keys show user4 -a --keyring-backend test) 8750000000000utrst

trstd gentx user1 750000000000utrst --chain-id=trst_chain_1 --keyring-backend=test  --website="trustlesshub.com" --security-contact="info@trstlabs.xyz"


trstd init-attestation
PUBLIC_KEY=$(trstd parse attestation_cert.der 2> /dev/null | cut -c 3-)
echo $PUBLIC_KEY
trstd init-bootstrap ./node-master-cert.der ./io-master-cert.der
echo "Collecting genesis txs..."
trstd collect-gentxs

echo "Validating genesis file..."
trstd validate-genesis

sed -i '384s/timeout_prevote = "1s"/timeout_prevote = "200ms"/g'  ~/.trstd/config/config.toml
sed -i '388s/timeout_precommit = "1s"/timeout_precommit = "200ms"/g'  ~/.trstd/config/config.toml
sed -i '394s/timeout_commit = "5s"/timeout_commit = "200ms"/g'  ~/.trstd/config/config.toml

sed -i '133s/enabled-unsafe-cors = false/enabled-unsafe-cors = true/g' ~/.trstd/config/app.toml

trstd config broadcast-mode block
trstd config keyring-backend test
trstd config chain-id trst_chain_1

trstd start --bootstrap > init.log --log_level info

