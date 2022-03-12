<!--
order: 0
-->

# List of Modules

Here are a number of modules that is used in the TRST chain:

- [claims](claims/spec/README.md) - Module managing the claiming process for the vested fairdrop.
- [alloc](alloc/spec/README.md) - Distribution of the inflation rewards to the respective  accounts (Staking, compute, community pool, devs).
- [compute](compute) - Module that binds to the CosmWasm VM that runs in Intel SGX secure enclaves
- [registration](registration) - Module that registers keys related to the enclave with Intel SGX as well as on-chain. 
- [mint](mint) - Module that regulates the inflation of TRST on the Hub

- [item](item) - **disabled and will not be part of the main chain** Module for on-chain NFTs that get their price from aggregating independent estimations by estimators. 