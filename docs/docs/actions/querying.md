---
sidebar_position: 4
title: Querying actions
description: How to retreive action data
---

## Ways to Query

Retreiving `AutoTX` -related information can be through the [TriggerPortal](https://triggerportal.zone) interface, a front-end integration through `TrustlessJS`, locally through a Command-Line Interface or through a Celestia RPC endpoint.

<!--
Here's an RPC endpoint: [openrpc.intento.zone](https://openrpc.intento.zone).
A list of RPC endpoints is to-be added. -->

## Queries

The available queries are as follows:

| Query                        | Description                                                                             | Parameter                                | Returns                                   | HTTP Method | Endpoint                                        |
| ---------------------------- | --------------------------------------------------------------------------------------- | ---------------------------------------- | ----------------------------------------- | ----------- | ----------------------------------------------- |
| InterchainAccountFromAddress | Returns the interchain account for a given owner address on a specified connection pair | QueryInterchainAccountFromAddressRequest | QueryInterchainAccountFromAddressResponse | GET         | /intento/intent/v1beta1/address-to-ica             |
| AutoTX                       | Returns the auto-executing interchain account transaction for a specified ID            | QueryActionRequest                       | QueryActionResponse                       | GET         | /intento/intent/v1beta1/action/{id}               |
| AutoTXs                      | Returns all action infomration                                                          | QueryActionsRequest                      | QueryActionsResponse                      | GET         | /intento/intent/v1beta1/actions                   |
| AutoTXsForOwner              | Returns all action infomration for a given owner                                        | QueryActionsForOwnerRequest              | QueryActionsForOwnerResponse              | GET         | /intento/intent/v1beta1/actions-for-owner/{owner} |
| AutoTXHistoru                | Returns action execution history for a given action                                     | QueryActionHistoryRequest                | QueryActionHistoryResponse                | GET         | /intento/intent/v1beta1/actions-history           |
| Params                       | Returns the total set of the Intent module parameters                                | QueryParamsRequest                       | QueryParamsResponse                       | GET         | /intento/intent/v1beta1/params                     |

These proto queries provide a convenient way to interact with the Intent module and access information about automatic interchain transactions.

You can use pagination fields to narrow down the scope.

| Field                  | Type    | Description                                                                                                                                                          |
| ---------------------- | ------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| pagination.key         | string  | A value returned in PageResponse.next_key to begin querying the next page most efficiently. Only one of offset or key should be set.                                 |
| pagination.offset      | string  | A numeric offset that can be used when key is unavailable. It is less efficient than using key. Only one of offset or key should be set.                             |
| pagination.limit       | string  | The total number of results to be returned in the result page. If left empty, it will default to a value to be set by each app.                                      |
| pagination.count_total | boolean | Set to true to indicate that the result set should include a count of the total number of items available for pagination in UIs. Only respected when offset is used. |
| pagination.reverse     | boolean | Set to true if results are to be returned in the descending order.                                                                                                   |
