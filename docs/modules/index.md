---
title: Introduction
order: 0
parent:
  order: 2
  title: Cosmos SDK Modules
---

# Introduction

Trustless Hub uses a number of custom modules.

## List of Modules

Here are a number of modules that are used in the TRST chain:

- auto-tx - automating transactions containing messages. These can be executed locally and on other chains using Interchain Accounts
- claims - Module managing the claiming process for the Rainbow Airdrop. Rainbow airdrop creates long-term allignment for airdrop recipients. 
- alloc - Distribution of the inflation rewards to the respective module accounts (mint, autotx, compute, community pool). Alloc allocates inflation minted to module addresses following a governance-defined issuance schedule
- compute - Module that binds to the CosmWasm VM that runs in Intel SGX secure enclaves. It also creates callback signatures needed for secure recurring transactions.
- registration - Module that registers keys related to the enclave with Intel SGX as well as on-chain. 
- mint - Module that regulates the inflation of TRST on the Hub. Mint mints inflation accroding to a 'fourting' schedule. Inflation is reduced by 25% yearly. Inflation graduatly decreases over time. 