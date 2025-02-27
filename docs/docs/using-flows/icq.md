---
sidebar_position: 4
title: Integrating Interchain Queries
description: Interchain Queries for conditions in Intent-Based Flows
---

The **Interchain Query (ICQ) feature** enables seamless cross-chain interactions by querying the state of one blockchain from another. This feature allows developers to create intent-based actions that trigger specific behaviors based on the queried data. By utilizing interchain queries, you can automate decision-making processes across multiple blockchains, providing an efficient way to orchestrate complex, multi-chain applications.

### How It Works

Interchain queries allow you to access the key-value store of a different blockchain by providing a specific key to query. IBC relayers submit query responses. This queried state can then be used to determine what actions should be triggered, based on predefined conditions. For instance, you can check the balance of a particular account on one blockchain and then execute a corresponding action on another chain, such as staking, transferring, or adjusting governance proposals.

### Supported Types and Proto Interfaces

You can query various data types from the state of other chains. There are two primary methods to accomplish this:

1. **Supported Types:** You can use one of the [supported types](./../module/supported_types.md), which are predefined and commonly used data types within the Cosmos ecosystem.
2. **Registered Proto Interfaces:** Alternatively, you can utilize the registered protocol buffer interfaces that adhere to Cosmos SDK standards. These provide a more flexible way to query and interpret the state, allowing you to work with complex data structures.

### Feedback Loops and Comparisons

Once the queried data is retrieved, it can be used for [**comparisons**](./conditions.md#comparison-operators) or to establish [**feedback loops**](./conditions.md#feedback-loops). For example, if a queried balance exceeds a certain threshold, you can trigger an action to stake the excess funds. Likewise, if a validator’s status on one chain changes, you can automatically adjust delegations or governance votes on another chain.

This capability unlocks numerous possibilities for cross-chain workflows, simplifying multi-chain dApp logic and empowering developers to build more dynamic and responsive applications in the interchain ecosystem.

## Integrating ICQs

You can use the ICQ feature by attaching an `ICQConfig` into [**comparisons**](./conditions.md#comparisons) and [**feedback loops**](./conditions.md#feedback-loops). In the `ICQConfig`, you specify what to query, where to query it, and how to handle a timeout scenario.

| Field              | Type                                       | Description                                                                               |
| ------------------ | ------------------------------------------ | ----------------------------------------------------------------------------------------- |
| `connection_id`    | `string`                                   | The ID of the connection to use for the interchain query.                                 |
| `chain_id`         | `string`                                   | The ID of the blockchain to query.                                                        |
| `timeout_policy`   | `intento.interchainquery.v1.TimeoutPolicy` | The policy to apply when a timeout occurs.                                                |
| `timeout_duration` | `google.protobuf.Duration`                 | The duration to wait before a timeout is triggered.                                       |
| `query_type`       | `string`                                   | The type of query to perform (e.g., `store/bank/key`, `store/staking/key`).               |
| `query_key`        | `string`                                   | The key in the store to query (e.g., `stakingtypes.GetValidatorKey(validatorAddressBz)`). |

For example, the `query_type` can be `store/bank/key` or `store/staking/key`. The `query_key` is the key in the store to query, such as `stakingtypes.GetValidatorKey(validatorAddressBz)`. The generation of query keys is abstracted in the TriggerPortal frontend.

```proto
// Config for using interchain queries
message ICQConfig {
  string connection_id = 1;
  string chain_id = 2;
  intento.interchainquery.v1.TimeoutPolicy timeout_policy = 3;
  google.protobuf.Duration timeout_duration = 4 [
    (gogoproto.nullable) = false,
    (gogoproto.stdduration) = true
  ];
  string query_type = 5;
  string query_key = 6;
}
```

If SaveResponses in the Flow Configuration is set to true, query responses are added to the Flow History. Check out the [**Supported Types**](./../module/supported_types.md) page or the TriggerPortal Flow Builder for some example queries.

Here’s a tutorial for integrating `ICQConfig` for querying balances and adding `connectionId`, `hostConnectionId`, and `HostedConfig` in `submitFlow`.

---

## Conditional Transfers with Intent-Based Flows

In this tutorial, we will explore how to use **Intent-Based Flows** to automate a token transfer, ensuring that it only executes if a certain balance condition is met. This is particularly useful for scenarios where funds should only be moved when a specific threshold is reached.

We'll achieve this by:

1. Querying the account balance using **Interchain Queries (ICQ)**.
2. Checking if the balance exceeds `200,000 uatom`.
3. Only then executing the **MsgSend** message.

---

### 1: Defining Execution Configuration

First, we need to set up the execution behavior of our flow.

```typescript
import {
  Coin,
  msgRegistry,
  Registry,
  Conditions,
  Comparison,
  ICQConfig,
  HostedConfig,
} from "intentojs";

const config: ExecutionConfiguration = {
  saveResponses: false,
  updatingDisabled: false,
  stopOnFailure: true,
  stopOnSuccess: false,
  fallbackToOwnerBalance: true,
  reregisterIcaAfterTimeout: true,
};
```

This ensures that if a condition fails, the execution stops, but if it succeeds, it continues.

---

### 2: Setting Up an Interchain Query (ICQ)

Before transferring funds, we need to **query the account balance** to determine if it meets the required threshold.

```typescript
const queryKey = createBankBalanceQueryKey("cosmos1delegatoraddress", "uatom");

const icqConfig: ICQConfig = {
  connectionId: "connection-123",
  chainId: "host-chain-1",
  timeoutPolicy: 2,
  timeoutDuration: 50000000000, // 50 seconds
  queryType: "store/bank/key",
  queryKey: queryKey,
  response: new Uint8Array(), // Will be populated with the ICQ response
};

const createBankBalanceQueryKey = (address: string, denom: string): string => {
  try {
    const { words } = bech32.decode(address);
    const addressBytes = new Uint8Array(bech32.fromWords(words));

    // Prefix (0x02) and address length
    const prefix = new Uint8Array([0x02, addressBytes.length]);

    // Convert denom to bytes
    const denomBytes = new TextEncoder().encode(denom);

    // Concatenate all parts into a single Uint8Array
    const queryData = new Uint8Array([
      ...prefix,
      ...addressBytes,
      ...denomBytes,
    ]);

    // Convert Uint8Array to Base64 for the query key
    return btoa(String.fromCharCode(...queryData));
  } catch (error) {
    console.error("Error creating query key:", error);
    return "";
  }
}
```

This configuration tells the system to query the balance from the **store/bank/key** module and wait up to 50 seconds for a response.

---

### Step 3: Implementing a Feedback Loop

A **feedback loop** will use the queried balance and dynamically insert it into the `MsgSend` message.

```typescript
const feedbackLoop: FeedbackLoop = {
  flowId: BigInt(0), // Balance query flow
  responseIndex: 0, // First response
  responseKey: "amount.[0].amount", // Extract balance amount
  valueType: "sdk.Int",
  msgsIndex: 1, // Index of MsgSend in the array
  msgKey: "amount", // Replace amount field in MsgSend
  icqConfig: icqConfig, // Uses ICQConfig for balance query
};
```

---

### 4: Adding a Conditional Check

We want to **only send funds if the balance exceeds 200,000 uatom**. We define a `Comparison` object to enforce this rule.

```typescript
const comparison: Comparison = {
  flowId: BigInt(0), // Balance query flow
  responseIndex: 0,
  responseKey: "amount.[0].amount",
  valueType: "sdk.Int",
  operator: 4, // LARGER_THAN
  operand: "200000uatom",
  icqConfig: icqConfig, // Uses ICQConfig for validation
};
```

---

### 5: Setting Execution Conditions

Now, we define conditions to ensure that:

- The **comparison rule** is met before executing the transfer.
- The **feedback loop** dynamically updates the `MsgSend` message.

```typescript
const initConditions: ExecutionConditions = {
  stopOnSuccessOf: [],
  stopOnFailureOf: [],
  skipOnFailureOf: [],
  skipOnSuccessOf: [],
  feedbackLoops: [feedbackLoop],
  comparisons: [comparison],
  useAndForComparisons: false,
};
```

---

### 6: Constructing the Messages

#### Message for Balance Query

The balance query is performed automatically via `ICQConfig`, so we do not need to explicitly add a query message.

#### Message for Conditional Transfer

```typescript
const msgSend = cosmos.bank.v1beta1.MessageComposer.withTypeUrl.send({
  fromAddress: "cosmos1delegatoraddress",
  toAddress: "cosmos1recipientaddress",
  amount: [{ denom: "uatom", amount: "0" }], // Will be replaced dynamically
});
```

---

### 7: Submitting the Intent-Based Flow

To submit the flow, we also include **hosted account configuration** using `HostedConfig`.

```typescript
const hostedConfig: HostedConfig = {
  hostedAddress: "cosmos1hostedaddress",
  feeCoinLimit: { denom: "uatom", amount: "100000" },
};

const msgSubmitFlow =
  intento.intent.v1beta1.MessageComposer.withTypeUrl.submitFlow({
    label: "Balance Query and Send Flow",
    owner: "into1wdplq6qjh2xruc7qqagma9ya665q6qhcpse4k6",
    msgs: [msgSend],
    duration: "1440h",
    interval: "600s",
    startAt: "1739781618",
    feeFunds: [{ denom: "uinto", amount: "5000000" }],
    configuration: config,
    hostedConfig: hostedConfig, // Config for hosted account
  });
```

---

### 8: Signing and Broadcasting the Transaction

Finally, we sign and send the transaction.

```typescript
client.signAndBroadcast(owner, [msgSubmitFlow], {
  amount: [],
  gas: "300000",
});
```

---

Final Thoughts

With this setup, our **Intent-Based Flow** automatically:

✅ Queries the account balance using **ICQ**.  
✅ Only proceeds with the transfer if the balance **exceeds 200,000 uatom**.  
✅ Dynamically inserts the correct amount into the `MsgSend` message.

This ensures efficient, automated, and secure fund transfers without requiring manual intervention. 🚀

Would you like to see this example with **dynamic fee handling** or **multi-asset transfers**? Let us know!

## ICQ module details - x/interchainqueries

SubmitQueryResponse is used to return the query responses. It is used by IBC Relayers. As we use the Stride implementation, it is the same as with Stride.

```protobuf
message MsgSubmitQueryResponse {
  string chain_id = 1;
  string query_id = 2;
  bytes result = 3;
  tendermint.crypto.ProofOps proof_ops = 4;
  int64 height = 5;
  string from_address = 6;
}
```

Query PendingQueries lists all queries that have been requested (i.e. emitted but have not had a response submitted yet)

```protobuf
message QueryPendingQueriesRequest {}
```
