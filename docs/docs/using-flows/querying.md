---
sidebar_position: 7
title: Querying Flows
description: How to retreive flow data
---

Retrieving `Flow`-related information can be done through multiple methods:

- **[Intento Portal](https://portal.intento.zone)** – A user-friendly interface for managing and monitoring Flows.
- **Front-end Integration via `IntentoJS`** – Enables seamless interaction with Flows in dApps.
- **Command-Line Interface (CLI)** – Local querying for developers and power users.
- **API Endpoints (RPC & LCD)** – Directly query `Flow` data through **RPC, LCD, or Swagger** endpoints.

#### **Example LCD Endpoints**

1. **Retrieve Flow Detail**

   ```
   GET /intento/intent/v1beta1/flow/{flow_id}
   ```

   Example LCD endpoint:

   ```
   https://lcd.intento.zone/intento/intent/v1beta1/flow/1
   ```

   Returns details of the specified Flow.

2. **List All Flows**
   ```
   GET /intento/intent/v1beta1/flows
   ```
   Example LCD Swagger UI:
   ```
   https://lcd.intento.zone/swagger/#get-/intento/intent/v1beta1/flows
   ```
   Fetches all active Flows.

## Queries

The available queries are as follows:

| Query                        | Description                                                                             | Parameter                                | Returns                                   | HTTP Method | Endpoint                                        |
| ---------------------------- | --------------------------------------------------------------------------------------- | ---------------------------------------- | ----------------------------------------- | ----------- | ----------------------------------------------- |
| InterchainAccountFromAddress | Returns the interchain account for a given owner address on a specified connection pair | QueryInterchainAccountFromAddressRequest | QueryInterchainAccountFromAddressResponse | GET         | /intento/intent/v1beta1/address-to-ica          |
| Flow                         | Returns the auto-executing interchain account flow for a specified ID                   | QueryFlowRequest                         | QueryFlowResponse                         | GET         | /intento/intent/v1beta1/flow/{id}               |
| Flows                        | Returns all flow infomration                                                            | QueryFlowsRequest                        | QueryFlowsResponse                        | GET         | /intento/intent/v1beta1/flows                   |
| FlowsForOwner                | Returns all flow infomration for a given owner                                          | QueryFlowsForOwnerRequest                | QueryFlowsForOwnerResponse                | GET         | /intento/intent/v1beta1/flows/{owner} |
| FlowHistory                  | Returns flow execution history for a given flow                                         | QueryFlowHistoryRequest                  | QueryFlowHistoryResponse                  | GET         | /intento/intent/v1beta1/flows-history           |
| Params                       | Returns the total set of the Intent module parameters                                   | QueryParamsRequest                       | QueryParamsResponse                       | GET         | /intento/intent/v1beta1/params                  |

These proto queries provide a convenient way to interact with the Intent module and access information about automatic interchain flows.

You can use pagination fields to narrow down the scope.

| Field                  | Type    | Description                                                                                                                                                          |
| ---------------------- | ------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| pagination.key         | string  | A value returned in PageResponse.next_key to begin querying the next page most efficiently. Only one of offset or key should be set.                                 |
| pagination.offset      | string  | A numeric offset that can be used when key is unavailable. It is less efficient than using key. Only one of offset or key should be set.                             |
| pagination.limit       | string  | The total number of results to be returned in the result page. If left empty, it will default to a value to be set by each app.                                      |
| pagination.count_total | boolean | Set to true to indicate that the result set should include a count of the total number of items available for pagination in UIs. Only respected when offset is used. |
| pagination.reverse     | boolean | Set to true if results are to be returned in the descending order.                                                                                                   |
