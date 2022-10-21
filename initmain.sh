
#!/bin/bash

DENOM=utrst
CHAIN_ID=trst_chain_1
ONE_HOUR=3600
ONE_DAY=$(($ONE_HOUR * 24))
ONE_YEAR=$(($ONE_DAY * 365))
TWO_YEARS=$(($ONE_YEAR * 2))
VALIDATOR_COINS=10000000$DENOM

rm -rf $HOME/.trstd
mkdir $HOME/opt/trustlesshub/.sgx_secrets

trstd init FRST --chain-id $CHAIN_ID
trstd prepare-genesis mainnet $CHAIN_ID


if [ "$1" == "mainnet" ]
then
    LOCKUP=ONE_YEAR
else
    LOCKUP=ONE_DAY
fi
echo "Lockup period is $LOCKUP"


trstd keys add validator

trstd add-genesis-account $(trstd keys show validator -a ) $VALIDATOR_COINS



echo "Adding airdrop accounts..."
trstd import-genesis-accounts-from-snapshot ./snapshot.json ./reserves.json 
echo "Getting genesis time..."

GENESIS_TIME=$(jq '.genesis_time' ~/.trstd/config/genesis.json | tr -d '"')
echo "Genesis time is $GENESIS_TIME"
if [[ "$OSTYPE" == "darwin"* ]]; then
    GENESIS_UNIX_TIME=$(TZ=UTC gdate "+%s" -d $GENESIS_TIME)
else
    GENESIS_UNIX_TIME=$(TZ=UTC date "+%s" -d $GENESIS_TIME)
fi


vesting_start_time=$(($GENESIS_UNIX_TIME))
vesting_end_time=$(($vesting_start_time + $LOCKUP))
vesting_end_time_two_years=$(($vesting_start_time + $TWO_YEARS))
echo "Adding vesting accounts..."
trstd add-genesis-account trust1vqan7n3hysjhamr49a3aa8keuawp9gpfdq5mts 8749990000000$DENOM
trstd add-genesis-account trust1dkf0q5u04nalrznw4fp35h5zymzddevdjdvs9t 8750000000000$DENOM \
    --vesting-amount 8750000000000$DENOM \
    --vesting-start-time $vesting_start_time \
    --vesting-end-time $vesting_end_time
trstd add-genesis-account trust1sns5l9cvkgf4fy770nmg98e7uzet5xhhmv8njv 8750000000000$DENOM \
    --vesting-amount 8750000000000$DENOM \
    --vesting-start-time $vesting_start_time \
    --vesting-end-time $vesting_end_time
trstd add-genesis-account trust1menylq7ttm3jne59lsj5xtc2all63l89lkglv3 8750000000000$DENOM \
    --vesting-amount 8750000000000$DENOM \
    --vesting-start-time $vesting_start_time \
    --vesting-end-time $vesting_end_time_two_years

echo "Gen tx ..."
trstd gentx validator 10000000utrst --chain-id=trst_chain_1 --keyring-backend=test  --website="trustlesshub.com" --security-contact="info@trstlabs.xyz"

echo "initing enclave ..."



trstd init-attestation --reset
PUBLIC_KEY=$(trstd parse attestation_cert.der 2> /dev/null | cut -c 3-)
echo $PUBLIC_KEY
trstd init-bootstrap ./node-master-cert.der ./io-master-cert.der
echo "Collecting genesis txs..."
trstd collect-gentxs

echo "Validating genesis file..."
trstd validate-genesis


trstd start --bootstrap