
tppd tx compute store ./independent_price_secret_contract/contract.wasm.gz --from faucet --fees 500000tpp --gas 2000000 --chain-id tpp --keyring-backend test -y
tppd tx compute instantiate 1 '{"estimationcount": "5"}' --label hahalololol --from faucet --fees 500000tpp --gas 2000000 --chain-id tpp --keyring-backend test -y

tppd tx compute execute cosmos18vd8fpwxzck93qlwghaj6arh4p7c5n89uzcee5 '{"send": {"amount": "46"}}' --from user1 --gas 500000 -y --fees 530tpp --chain-id tpp --keyring-backend test -y
 