<!--
order: 0
-->

# List of Modules

Here are a number of customized Cosmos SDK modules that are used in TRST:

- [auto-ibc-tx](auto-ibc-tx/README.md) - Module for Intent-Centric Automation

- [claim](claim/spec/README.md) - Module managing the claiming process for the vested airdrop.
- [alloc](alloc/README.md) - Distribution of the inflation rewards to the respective  accounts (Staking, Relayer Rewards, Community Pool).
- [mint](mint) - Module that regulates the inflation of TRST on the Hub

## Disabled

- [item](item) - **disabled and will not be part of the main chain** Module for on-chain NFTs that get their price from aggregating independent estimations by estimators. 
- [compute](compute) - Module that binds to the CosmWasm VM that runs in Intel SGX secure enclaves
- [registration](registration) - Module that registers keys related to the enclave with Intel SGX as well as on-chain. 