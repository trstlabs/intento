
trstd tx compute store ./independent_price_secret_contract/contract.wasm.gz --from faucet --fees 500000trst --gas 2000000 --chain-id trst --keyring-backend test -y
trstd tx compute instantiate 1 '{"estimationcount": "5"}' --label hahalololol --from faucet --fees 500000trst --gas 2000000 --chain-id trst --keyring-backend test -y

trstd tx compute execute cosmos18vd8fpwxzck93qlwghaj6arh4p7c5n89uzcee5 '{"send": {"amount": "46"}}' --from user1 --gas 500000 -y --fees 530trst --chain-id trst --keyring-backend test -y
 