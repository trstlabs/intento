---
sidebar_position: 9
title: Flow Submission Page
description: How to submit flows to Intento with just a button
---

### **Flow Submission Page**

for integrators. 

You can also **encode a flow directly into the page** using a `flowInput` object. This is useful for static flows—just include the `flowInput` in your frontend code or page build, and the system will handle it when the page is loaded. It’s a clean way to preload flows without runtime API calls.

let users sign up for alerts, return flow, redirect back to your site. all possible.
Supports custom theme via logo and colors.


https://portal.intento.zone/submit?


### What the /submit page expects in the URL:

* `flowInput`: your flow data as **JSON**, but **URL-encoded** (like replacing spaces and special chars with `%` codes).
* Other params:

  * `imageUrl` — a link to an image
  * `chain` — the blockchain ID
  * `bgColor` — background color in hex (like `#123abc`)
  * `theme` — either `light` or nothing

---

### What you do:

1. Turn your flow input (the object) into JSON.
2. URL-encode that JSON.
3. Add it as `flowInput=...` in the URL.
4. Add other params normally.
5. The page will decode everything and use it.

---

No need to worry about decoding — just send URL-encoded JSON for `flowInput` and plain strings for the rest. The page handles the rest.


### The `FlowInput` Structure

`FlowInput` defines the parameters required for setting up and executing an Intento Flow. This object gets encoded to proto before being sent to the Intento blockchain. Here's the breakdown of each field and what it controls:

#### Core Fields:

* **`label?: string`**
  *Optional.* A human-readable label for identifying the flow. Think of this as metadata for users or systems.

* **`msgs: string[]`**
  *Required.* An array of encoded Cosmos SDK messages (likely base64-encoded `Any` messages). These are the core actions your flow will execute.

  In the messages you can use the `Your Address` placeholder for dynamic address insertion. And with the `ICA_ADDR` placeholder the chain will insert the address of the Interchain account for the flow execution.

#### Time Control:

* **`duration: number`**
  *Required.* Total duration of the flow in seconds. Defines how long the flow will remain active or how long the streaming (DCA) will take.

* **`interval?: number`**
  *Optional.* Interval in seconds for repeating the flow. If set, the system will trigger execution every `interval` seconds until the `duration` is over. If unset, flow executes once.

* **`startTime?: number`**
  *Optional.* Unix timestamp (in seconds) for when the flow should start. If unset, the flow starts immediately.

#### Execution Control:

* **`feeFunds?: Coin`**
  *Optional.* A `Coin` object specifying extra funds to cover fees, e.g. `{ denom: 'uatom', amount: '10000' }`.

* **`configuration?: ExecutionConfiguration`**
  *Optional.* Configs like gas limits, memo, or execution flags. Imported from `intentojs`.

* **`conditions?: ExecutionConditions`**
  *Optional.* Conditions that must be met before the flow executes. For example, chain state conditions like a minimum block height or price threshold.

* **`hostedIcaConfig?: HostedICAConfig`**
  *Optional.* Configurations for Trustless Excution Agent, like the target chain, controller address, or permissions. Imported from `intentojs`.

* **`icaAddressForAuthZ?: string`**
  *Optional.* If using Authz-based flows, this specifies the address holding authorization rights (likely an ICA address).

* **`connectionId?: string`**
  *Optional.* The IBC connection ID for flows involving IBC transfers. If not set, defaults may be used.

* **`hostConnectionId?: string`**
  *Optional.* Similar to `connectionId`, but specifically for hosted flows.

#### Alerting & Notifications:

* **`email?: string`**
  *Optional.* Email address to send flow status notifications to (e.g., flow started, completed, or errored).

* **`alertType?: string`**
  *Optional.* Type of alert (e.g., `all`, `error`).

---

### Example Use Case: A Streaming Flow

Say you want to stream tokens to Osmosis over 24 hours, with actions executed every 30 minutes. You'd configure your `FlowInput` like this:

```ts
const flow: FlowInput = {
  label: 'Stream ATOM to OSMO',
  msgs: [ /* encoded IBC transfer messages here */ ],
  duration: 86400, // 24 hours in seconds
  interval: 1800, // Every 30 minutes
  feeFunds: { denom: 'uatom', amount: '5000' },
  configuration: { /* custom gas/memo settings */ },
  email: 'user@example.com', 
}
```


## Handling Responses

After submission, the page will handle the flow creation process and display appropriate success/error messages to the user. The page will automatically handle:

- Wallet connection (if not already connected)
- Transaction signing
- Error handling and user feedback
- Loading states

## Best Practices

1. **Always URL-encode parameters**: Use `encodeURIComponent()` for any dynamic values
2. **Keep URLs short**: For complex flows, use the `flow` parameter with base64-encoded JSON
3. **Handle errors gracefully**: The page will show appropriate error messages, but your integration should also handle navigation failures
4. **Test thoroughly**: Test all URL parameters in both development and production environments
5. **Respect user preferences**: Allow users to override any default theme settings

## Troubleshooting
- For theme-related issues, check the browser's console for any errors from the theme controller
- If pre-filled data isn't loading, verify that the JSON structure matches the expected `FlowInput` type

### Examples

[Autocompounding Cosmos ATOM Flow](https://portal.intento.zone/submit?flowInput=%7B%22duration%22%3A0%2C%22msgs%22%3A%5B%22%7B%5Cn++%5C%22typeUrl%5C%22%3A+%5C%22%2Fcosmos.authz.v1beta1.MsgExec%5C%22%2C%5Cn++%5C%22value%5C%22%3A+%7B%5Cn++++%5C%22grantee%5C%22%3A+%5C%22Your+Address%5C%22%2C%5Cn++++%5C%22msgs%5C%22%3A+%5B%5Cn++++++%7B%5Cn++++++++%5C%22typeUrl%5C%22%3A+%5C%22%2Fcosmos.distribution.v1beta1.MsgWithdrawDelegatorReward%5C%22%2C%5Cn++++++++%5C%22value%5C%22%3A+%7B%5Cn++++++++++%5C%22delegatorAddress%5C%22%3A+%5C%22cosmos1u7zn9sxz8s63ww8xwg8cl7xlmwkedq7a63wke7%5C%22%2C%5Cn++++++++++%5C%22validatorAddress%5C%22%3A+%5C%22cosmosvaloper19ge9c23yuj3n520xemczvkgfunsrlqfpk2add3%5C%22%5Cn++++++++%7D%5Cn++++++%7D%5Cn++++%5D%5Cn++%7D%5Cn%7D%22%2C%22%7B%5Cn++%5C%22typeUrl%5C%22%3A+%5C%22%2Fcosmos.authz.v1beta1.MsgExec%5C%22%2C%5Cn++%5C%22value%5C%22%3A+%7B%5Cn++++%5C%22grantee%5C%22%3A+%5C%22Your+Address%5C%22%2C%5Cn++++%5C%22msgs%5C%22%3A+%5B%5Cn++++++%7B%5Cn++++++++%5C%22typeUrl%5C%22%3A+%5C%22%2Fcosmos.staking.v1beta1.MsgDelegate%5C%22%2C%5Cn++++++++%5C%22value%5C%22%3A+%7B%5Cn++++++++++%5C%22delegatorAddress%5C%22%3A+%5C%22cosmos1u7zn9sxz8s63ww8xwg8cl7xlmwkedq7a63wke7%5C%22%2C%5Cn++++++++++%5C%22validatorAddress%5C%22%3A+%5C%22cosmosvaloper19ge9c23yuj3n520xemczvkgfunsrlqfpk2add3%5C%22%2C%5Cn++++++++++%5C%22amount%5C%22%3A+%7B%5Cn++++++++++++%5C%22denom%5C%22%3A+%5C%22uatom%5C%22%2C%5Cn++++++++++++%5C%22amount%5C%22%3A+%5C%2210%5C%22%5Cn++++++++++%7D%5Cn++++++++%7D%5Cn++++++%7D%5Cn++++%5D%5Cn++%7D%5Cn%7D%22%5D%2C%22conditions%22%3A%7B%22feedbackLoops%22%3A%5B%7B%22flowId%22%3A%220%22%2C%22responseIndex%22%3A0%2C%22responseKey%22%3A%22Amount.%5B0%5D%22%2C%22msgsIndex%22%3A1%2C%22msgKey%22%3A%22Amount%22%2C%22valueType%22%3A%22sdk.Coin%22%7D%5D%2C%22comparisons%22%3A%5B%7B%22flowId%22%3A%220%22%2C%22responseIndex%22%3A0%2C%22responseKey%22%3A%22Amount.%5B0%5D%22%2C%22valueType%22%3A%22sdk.Coin%22%2C%22operator%22%3A4%2C%22operand%22%3A%221uatom%22%7D%5D%2C%22stopOnSuccessOf%22%3A%5B%5D%2C%22stopOnFailureOf%22%3A%5B%5D%2C%22skipOnFailureOf%22%3A%5B%5D%2C%22skipOnSuccessOf%22%3A%5B%5D%2C%22useAndForComparisons%22%3Afalse%7D%2C%22configuration%22%3A%7B%22saveResponses%22%3Atrue%2C%22updatingDisabled%22%3Afalse%2C%22stopOnSuccess%22%3Afalse%2C%22stopOnFailure%22%3Afalse%2C%22stopOnTimeout%22%3Afalse%2C%22fallbackToOwnerBalance%22%3Atrue%7D%2C%22connectionId%22%3A%22connection-0%22%2C%22hostedIcaConfig%22%3A%7B%22agentAddress%22%3A%22into1gzakqp6uammdhhpdgcsjjqzyzayelfzn38v3q7sfgf5uacc6ltvqswckct%22%2C%22feeCoinLimit%22%3A%7B%22denom%22%3A%22uinto%22%2C%22amount%22%3A%2220%22%7D%7D%2C%22label%22%3A%22Conditional+Autocompound%22%7D&chain=GAIA&bgColor=#315faa)

[CosmWasm DCA Flow](https://portal.intento.zone/submit?flowInput=%7B%20%20%20%22msgs%22:%5B%20%20%20%20%20%22%7B%5Cn%20%20%5C%22typeUrl%5C%22:%20%5C%22/cosmos.authz.v1beta1.MsgExec%5C%22,%5Cn%20%20%5C%22value%5C%22:%20%7B%5Cn%20%20%20%20%5C%22grantee%5C%22:%20%5C%22ICA_ADDR%5C%22,%5Cn%20%20%20%20%5C%22msgs%5C%22:%20%5B%5Cn%20%20%20%20%20%20%7B%5Cn%20%20%20%20%20%20%20%20%5C%22typeUrl%5C%22:%20%5C%22/cosmwasm.wasm.v1.MsgExecuteContract%5C%22,%5Cn%20%20%20%20%20%20%20%20%5C%22value%5C%22:%20%7B%5Cn%20%20%20%20%20%20%20%20%20%20%5C%22sender%5C%22:%20%5C%22Your%20Address%5C%22,%5Cn%20%20%20%20%20%20%20%20%20%20%5C%22contract%5C%22:%20%5C%22osmo10wn49z4ncskjnmf8mq95uyfkj9kkveqx9jvxylccjs2w5lw4k6gsy4cj9l%5C%22,%5Cn%20%20%20%20%20%20%20%20%20%20%5C%22msg%5C%22:%20%7B%5Cn%20%20%20%20%20%20%20%20%20%20%20%20%5C%22subscribe%5C%22:%20%7B%5Cn%20%20%20%20%20%20%20%20%20%20%20%20%20%20%5C%22stream_id%5C%22:%2046%5Cn%20%20%20%20%20%20%20%20%20%20%20%20%7D%5Cn%20%20%20%20%20%20%20%20%20%20%7D,%5Cn%20%20%20%20%20%20%20%20%20%20%5C%22funds%5C%22:%20%5B%5Cn%20%20%20%20%20%20%20%20%20%20%20%20%7B%5Cn%20%20%20%20%20%20%20%20%20%20%20%20%20%20%5C%22denom%5C%22:%20%5C%22factory/osmo1nz7qdp7eg30sr959wvrwn9j9370h4xt6ttm0h3/ussosmo%5C%22,%5Cn%20%20%20%20%20%20%20%20%20%20%20%20%20%20%5C%22amount%5C%22:%20%5C%22100%5C%22%5Cn%20%20%20%20%20%20%20%20%20%20%20%20%7D%5Cn%20%20%20%20%20%20%20%20%20%20%5D%5Cn%20%20%20%20%20%20%20%20%7D%5Cn%20%20%20%20%20%20%7D%5Cn%20%20%20%20%5D%5Cn%20%20%7D%5Cn%7D%22%20%20%20%5D,%20%20%20%22conditions%22:%20%7B%20%20%20%20%20%22feedbackLoops%22:%20%5B%5D,%20%20%20%20%20%22comparisons%22:%20%5B%5D,%20%20%20%20%20%22stopOnSuccessOf%22:%20%5B%5D,%20%20%20%20%20%22stopOnFailureOf%22:%20%5B%5D,%20%20%20%20%20%22skipOnFailureOf%22:%20%5B%5D,%20%20%20%20%20%22skipOnSuccessOf%22:%20%5B%5D,%20%20%20%20%20%22useAndForComparisons%22:%20false%20%20%20%7D,%20%20%20%22configuration%22:%20%7B%20%20%20%20%20%22saveResponses%22:%20false,%20%20%20%20%20%22updatingDisabled%22:%20false,%20%20%20%20%20%22stopOnSuccess%22:%20false,%20%20%20%20%20%22stopOnFailure%22:%20false,%20%20%20%20%20%22stopOnTimeout%22:%20false,%20%20%20%20%20%22fallbackToOwnerBalance%22:%20true%20%20%20%7D,%20%20%20%22connectionId%22:%20%22connection-2%22,%20%20%20%22hostedIcaConfig%22:%20%7B%20%20%20%20%20%22agentAddress%22:%20%22into1p9ccttjgzh5wlewm5s55qk73j9ccjt27x00tada89sfq5t9v69rsex0977%22,%20%20%20%20%20%22feeCoinLimit%22:%20%7B%20%20%20%20%20%20%20%22denom%22:%20%22uinto%22,%20%20%20%20%20%20%20%22amount%22:%20%2250%22%20%20%20%20%20%7D%20%20%20%7D,%20%20%20%22label%22:%20%22Subscribe%20via%20hosted%20ICA%20%F0%9F%8E%AF%22%20%7D&chain=osmo-test-5&bgColor=#140739)