trstd tx trst create-item 'Rolex Submariner 1997 Gray' 'Rolex Submariner in good condition, it has no visible scratches and still works great. Bought in 1998, model year is 1997. It is the gray edition' 5 '40.741895,-73.989308' 3 watch,submariner,rolex 5 nl 3 "" --from=user1 -y --chain-id=trst_chain_1 --keyring-backend test --fees 150trst

trstd tx trst create-estimation 145 3 1 'Great Photos!' 0 --from=user4 -y --chain-id=trst_chain_1 --keyring-backend test --fees 150trst



trstd q trst list-item

trstd tx trst create-flag 0 --from user2  --keyring-backend test --fees 150trst --chain-id=trst_chain_1 -y

trstd tx trst delete-estimator 0 --from user4  --keyring-backend test --fees 150trst --chain-id=trst_chain_1 -y

trstd keys export user1 --keyring-backend test --unarmored-hex --unsafe


trstd tx compute store contract.wasm.gz --from faucet --fees 500000trst --gas 2000000 --chain-id trst --keyring-backend test -y

trstd tx compute instantiate 2 '{"estimationcount": "3"}' --label ddd --from faucet --fees 500000trst --gas 2000000 --amount 50trst  --chain-id trst --keyring-backend test -y
 trstd q compute list-contract-by-code 2

 trstd q account cosmos10pyejy66429refv3g35g2t7am0was7yacjc2l4

trstd query bank balances cosmos1qxxlalvsdjd07p07y3rc5fu6ll8k4tmecu7e9y

./testitems.sh



trstd tx compute instantiate 2 '{"estimationcount": "3"}' --label sdd --from faucet --fees 500000trst --gas 2000000 --amount 50trst  --chain-id trst --keyring-backend test -y