rm -rf ~/.trst


trstd init mainnode --chain-id=trst_chain_1

yes exchange cabin middle shed identify soon loop vivid mutual simple sing vessel tail embody vote glide bid olive possible invite merry kitten keen nuclear | trstd keys add user1 --keyring-backend test --recover

yes comic broom zone grass reject apology erupt chef wish add actor damage deputy hip aware connect addict excite poem arrive since bird couple artwork | trstd keys add user3 --keyring-backend test --recover

yes clutch title peasant erosion ugly sorry wear fringe question fuel allow silk outdoor sing oppose juice clown thunder aware parade bind health pizza balcony | trstd keys add user2 --keyring-backend test --recover

yes kiwi obtain scrub aunt female shoulder dune shove budget salt mechanic plug beef right pact economy swear flash update wild change puppy hurdle power | trstd keys add user4 --keyring-backend test --recover

yes orchard thing tooth dismiss seat couple define atom antenna language fuel wrist napkin tired undo toddler virus cherry shock mimic toss rifle predict crisp |trstd keys add faucet --keyring-backend test --recover


trstd add-genesis-account $(trstd keys show user1 -a --keyring-backend test) 10000000000000000utrst,100000000000stake  --vesting-amount 200000utrst  --vesting-end-time 1638485671
trstd add-genesis-account $(trstd keys show user2 -a --keyring-backend test) 10000000000000000utrst
trstd add-genesis-account $(trstd keys show user3 -a --keyring-backend test) 10000000000000000utrst
trstd add-genesis-account $(trstd keys show user4 -a --keyring-backend test) 10000000000000000utrst
trstd add-genesis-account $(trstd keys show faucet -a --keyring-backend test) 10000000000000000utrst

trstd gentx user1 100000000000stake --chain-id=trst_chain_1 --keyring-backend=test  --website="trustlesshub.com" --security-contact="trustlesshub@gmail.com"

trstd init-enclave --reset
PUBLIC_KEY=$(trstd parse attestation_cert.der 2> /dev/null | cut -c 3-)
echo $PUBLIC_KEY
trstd init-bootstrap ./node-master-cert.der ./io-master-cert.der
echo "Collecting genesis txs..."
trstd collect-gentxs

echo "Validating genesis file..."
trstd validate-genesis

sed -i '104s/enable = false/enable = true/g' ~/.trst/config/app.toml
sed -i 's/cors_allowed_origins = []/cors_allowed_origins = ["*"]/g' ~/.trst/config/config.toml
trstd start --bootstrap

