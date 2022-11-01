trstd tx compute store ~/contract/contract.wasm.gz 'Trustless Auction' 'A simple auction. Send tokens to this Trustless Contract. The person with the highest token amount will be the winning bid. When the contract ends, tokens are sent back automatically in case the minimum is not met. Code is available on github.com/trstlabs/independent_price_secret_contract/blob/master/trustless_auction/. To instantiate, the JSON fields are name and denom ' 0 --from user1 --fees 5000utrst --gas 2000000 --chain-id trst_chain_1 --keyring-backend test -y


trstd tx compute instantiate 2  '{"name": "BIG BOOK","denom": "utrst" }' --from user1 --fees 5000utrst --gas 2000000 --chain-id trst_chain_1 --keyring-backend test -y --contract_id 'For sale' --auto_msg

trstd tx compute execute trust1qxxlalvsdjd07p07y3rc5fu6ll8k4tme3pqv38 '{"bid":{}}' 2 --from user1 --fees 600utrst --gas 2000000 --chain-id trst_chain_1 --keyring-backend test -y --amount 56utrst