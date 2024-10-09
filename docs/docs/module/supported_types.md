---
sidebar_position: 1
title: Supported Types
description: ICQ Supported Types for Interchain Queries
---

When performing interchain queries, data is retrieved from the **key-value store (KV store)** of the relevant blockchain. Below is a comprehensive list of supported types from Cosmos SDK, Osmosis, CosmWasm, and Ethermint, along with their proto definitions, KV store paths, and query paths.

### Cosmos SDK Types

1. **Account Balances**
   - **Proto Type:** `cosmos.bank.v1beta1.Balance`
   - **KV Store Path:** `store/bank/key`
   - **Query Path:** `/bank/balances/{address}`

2. **Validator Information**
   - **Proto Type:** `cosmos.staking.v1beta1.Validator`
   - **KV Store Path:** `store/staking/key/validators`
   - **Query Path:** `/staking/validators/{validator_address}`

3. **Delegations**
   - **Proto Type:** `cosmos.staking.v1beta1.Delegation`
   - **KV Store Path:** `store/staking/key/delegations`
   - **Query Path:** `/staking/delegations/{delegator_address}`

4. **Governance Proposals**
   - **Proto Type:** `cosmos.gov.v1beta1.Proposal`
   - **KV Store Path:** `store/gov/proposals`
   - **Query Path:** `/gov/proposals/{proposal_id}`

5. **Voting Power**
   - **Proto Type:** `tendermint.abci.VoteInfo`
   - **KV Store Path:** `store/consensus/key/voting_power`
   - **Query Path:** `/validatorsets/{height}`

6. **Transaction History**
   - **Proto Type:** `cosmos.tx.v1beta1.Tx`
   - **KV Store Path:** `store/tx/history`
   - **Query Path:** `/txs/{tx_hash}`

7. **Unbonding Delegations**
   - **Proto Type:** `cosmos.staking.v1beta1.UnbondingDelegation`
   - **KV Store Path:** `store/staking/key/unbonding_delegations`
   - **Query Path:** `/staking/unbonding_delegations/{delegator_address}`

8. **Community Pool**
   - **Proto Type:** `cosmos.distribution.v1beta1.CommunityPool`
   - **KV Store Path:** `store/distribution/key/community_pool`
   - **Query Path:** `/distribution/community_pool`

### Osmosis Types

1. **Pool Information**
   - **Proto Type:** `osmosis.gamm.v1beta1.Pool`
   - **KV Store Path:** `store/gamm/key/pools`
   - **Query Path:** `/osmosis/gamm/v1beta1/pools/{pool_id}`

2. **Liquidity Incentives**
   - **Proto Type:** `osmosis.incentives.v1beta1.Incentive`
   - **KV Store Path:** `store/incentives/key`
   - **Query Path:** `/osmosis/incentives/v1beta1/incentives/{gauge_id}`

3. **Superfluid Staking**
   - **Proto Type:** `osmosis.superfluid.v1beta1.SuperfluidAsset`
   - **KV Store Path:** `store/superfluid/key/assets`
   - **Query Path:** `/osmosis/superfluid/v1beta1/assets`

4. **Incentivized Gauges**
   - **Proto Type:** `osmosis.incentives.v1beta1.Gauge`
   - **KV Store Path:** `store/incentives/key/gauges`
   - **Query Path:** `/osmosis/incentives/v1beta1/gauges/{gauge_id}`

5. **TWAP**
   - **KV Store Path:** `store/twap/recent_twap|{pool_id}|{denom1}|{denom2}`
   - **Historical TWAP Pool Index:** `store/twap/historical_pool_index|{pool_id}|{denom1}|{denom2}|{time}`


### CosmWasm Types

1. **Smart Contract Code**
   - **Proto Type:** `cosmwasm.wasm.v1.CodeInfo`
   - **KV Store Path:** `store/wasm/0x01/{code_id}` (CodeKeyPrefix)
   - **Query Path:** `/wasm/codes/{code_id}`
   - **Description:** Stores compiled contract code deployed on-chain, identified by code IDs.

2. **Smart Contract Instance (Contract Info)**
   - **Proto Type:** `cosmwasm.wasm.v1.ContractInfo`
   - **KV Store Path:** `store/wasm/0x02/{contract_address}` (ContractKeyPrefix)
   - **Query Path:** `/wasm/contracts/{contract_address}/info`
   - **Description:** Stores metadata about the contract, including instantiator, admin, and the code it runs.

3. **Smart Contract State**
   - **Proto Type:** N/A (Custom Query)
   - **KV Store Path:** `store/wasm/0x03/{contract_address}/{key}` (ContractStorePrefix)
   - **Query Path:** `/wasm/contracts/{contract_address}/state/{key}`
   - **Description:** Stores the state of CosmWasm contracts. Contracts store arbitrary data based on their internal logic and this data is stored under this key.

4. **Contract Code History**
   - **Proto Type:** `cosmwasm.wasm.v1.CodeHistoryEntry`
   - **KV Store Path:** `store/wasm/0x05/{contract_address}/{position}` (ContractCodeHistoryElementPrefix)
   - **Query Path:** `/wasm/contracts/{contract_address}/history`
   - **Description:** Tracks code history changes (such as migrations) of a contract.

5. **Pinned Contract Codes**
   - **Proto Type:** N/A
   - **KV Store Path:** `store/wasm/0x07/{code_id}` (PinnedCodeIndexPrefix)
   - **Query Path:** `/wasm/pinned_codes/{code_id}`
   - **Description:** Stores contracts that are pinned in memory, meaning they are cached for faster execution.

6. **Contracts by Creator**
   - **Proto Type:** N/A
   - **KV Store Path:** `store/wasm/0x09/{creator_address}/{created_time}/{contract_address}` (ContractsByCreatorPrefix)
   - **Query Path:** N/A
   - **Description:** Index of contracts created by a specific address, useful for querying all contracts created by an entity.

### Ethermint Types

Ethermint provides an EVM (Ethereum Virtual Machine) compatible layer within Cosmos, making it possible to run Ethereum smart contracts. Some relevant types include:

### Ethermint (EVM) Types

1. **EVM Code (Contract Bytecode)**
   - **Proto Type:** N/A (Stored as raw bytecode)
   - **KV Store Path:** `store/evm/0x01/{contract_address}` (KeyPrefixCode)
   - **Query Path:** `/evm/code/{contract_address}`
   - **Description:** Stores the raw bytecode of Ethereum-compatible smart contracts deployed on the Ethermint chain. Each contract is identified by its address.

2. **EVM Storage (Contract Storage)**
   - **Proto Type:** N/A (Key-Value mapping for contract state)
   - **KV Store Path:** `store/evm/0x02/{contract_address}/{storage_key}` (KeyPrefixStorage)
   - **Query Path:** `/evm/storage/{contract_address}/{storage_key}`
   - **Description:** Stores the state of Ethereum-compatible contracts. This is a mapping of storage keys to values for each contract, organized by the contract's address.

3. **EVM Module Parameters**
   - **Proto Type:** `ethermint.evm.v1.Params`
   - **KV Store Path:** `store/evm/0x03` (KeyPrefixParams)
   - **Query Path:** `/evm/params`
   - **Description:** Stores the configuration parameters for the Ethermint EVM module, including gas costs, state parameters, and transaction settings.

4. **EVM Account State**
   - **Proto Type:** N/A (Stored as raw EVM state data)
   - **KV Store Path:** `store/evm/0x02/{contract_address}/{storage_key}` (AddressStoragePrefix)
   - **Query Path:** `/evm/account/{contract_address}`
   - **Description:** Stores the complete state of an EVM account, including balances, nonces, and storage information.


---

For further details on these specific paths and proto types, refer to the official documentation for [Cosmos SDK](https://docs.cosmos.network), [Osmosis](https://docs.osmosis.zone), [CosmWasm](https://docs.cosmwasm.com), and [Ethermint](https://docs.ethermint.zone).
