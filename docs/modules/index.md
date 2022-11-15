---
title: Introduction
order: 0
parent:
  order: 3
  title: Cosmos SDK Modules
---

# Introduction

Trustless Hub uses a number of custom modules.

## List of Modules

Here are a number of modules that is used in the TRST chain:

- claims - Module managing the claiming process for the Rainbow Airdrop. Rainbow airdrop creates long-term allignment for airdrop recipients. 
- alloc - Distribution of the inflation rewards to the respective  accounts (Staking, compute, community pool, contributors). Alloc allocates inflation minted to addresses and modules, following a governance-defined schedule
- compute - Module that binds to the CosmWasm VM that runs in Intel SGX secure enclaves. It also creates callback signatures needed for secure recurring transactions.
- registration - Module that registers keys related to the enclave with Intel SGX as well as on-chain. 
- mint - Module that regulates the inflation of TRST on the Hub. Mint mints inflation accroding to a 'thirdning' schedule. Inflation graduatly decreases over time. 