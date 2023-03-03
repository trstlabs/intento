---
order: 3
title: Querying AutoTxs
description: How to retreive TrustlessTrigger data
---

## Ways to Query

Retreiving `AutoTx` info can be through the [TriggerPørtal](triggerportal.netlify.app) interface, a TrustlessJS front-end integration or locally through a light client Command-Line Interface.

A list of RPC endpoints are to-be added. Here's one for now: [openrpc.trustlesshub.com](openrpc.trustlesshub.com)

## GRPC Queries

The proto queries define the gRPC querier service for interacting with the AutoIbcTx module. 

The available queries are as follows:

| Query | Description | Parameter | Returns | HTTP Method | Endpoint |
|-------|-------------|-----------|---------|-------------|----------|
| InterchainAccountFromAddress | Returns the interchain account for a given owner address on a specified connection pair | QueryInterchainAccountFromAddressRequest | QueryInterchainAccountFromAddressResponse | GET | /auto-ibc-tx/v1beta1/address-to-ica |
| AutoTx | Returns the auto-executing interchain account transaction for a specified ID | QueryAutoTxRequest | QueryAutoTxResponse | GET | /auto-ibc-tx/v1beta1/auto-tx/{id} |
| AutoTxs | Returns all auto-executing interchain account messages | QueryAutoTxsRequest | QueryAutoTxsResponse | GET | /auto-ibc-tx/v1beta1/auto-txs |
| AutoTxsForOwner | Returns all auto-executing interchain account messages for a given owner | QueryAutoTxsForOwnerRequest | QueryAutoTxsForOwnerResponse | GET | /auto-ibc-tx/v1beta1/auto-txs-for-owner/{owner} |
| Params | Returns the total set of AutoIbcTx parameters | QueryParamsRequest | QueryParamsResponse | GET | /auto-ibc-tx/v1beta1/params |


These proto queries provide a convenient way to interact with the AutoIbcTx module and access information about automatic interchain transactions.
