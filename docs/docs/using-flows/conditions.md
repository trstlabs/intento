---
sidebar_position: 2
title: Conditions
pagination_label: Conditions
---

In this part of the docs we detail how to use `Conditions` such as `Comparisons` and `Feedback Loops` into your Intent-based Flow. With Conditions on Intento you can orchestrate conditional workflows in a structured manner. We will explore these features with tutorials, including an auto-compounding scenario using `MsgSend`.

## Conditions

`Execution Conditions` define the rules that determine when and how actions are executed. By combining feedback loops, comparisons, and dependent intent logic, these conditions enable precise control over execution flows.

Imagine a financial ecosystem where every action depends on intricate conditions to proceed. Feedback loops ensure that the system adapts in real-time by feeding outputs from one action as inputs to another. Meanwhile, comparisons act as gatekeepers, evaluating critical data points to decide whether an action should continue. Dependent intent arrays bring context-aware decision-making, stopping or skipping execution based on the success or failure of other operations. And finally, the flexibility to combine comparisons with AND or OR logic provides a fine-grained control mechanism for complex workflows.

### Key Elements of Execution Conditions

| Field                     | Type                    | Description                                                                                                    |
| ------------------------- | ----------------------- | -------------------------------------------------------------------------------------------------------------- |
| `feedback_loops`          | `repeated FeedbackLoop` | A list of feedback loops that dynamically replace message values with the latest response values.              |
| `comparisons`             | `repeated Comparison`   | A list of comparisons to evaluate response values and decide whether execution should proceed.                 |
| `stop_on_success_of`      | `repeated uint64`       | An array of dependent intents; execution halts if any of these execute successfully.                           |
| `stop_on_failure_of`      | `repeated uint64`       | An array of dependent intents; execution halts if any of these fail to execute.                                |
| `skip_on_failure_of`      | `repeated uint64`       | An array of dependent intents that must execute successfully after their latest call for execution to proceed. |
| `skip_on_success_of`      | `repeated uint64`       | An array of dependent intents that must fail after their latest call for execution to proceed.                 |
| `use_and_for_comparisons` | `bool`                  | Determines whether comparisons are combined with AND (true) or OR (false).                                     |

## Comparisons

The `Comparison` feature allows you to compare a response value with a specified operand before executing an action message. This comparison determines whether the action proceeds based on the result. You can configure up to 5 comparisons per flow. It evaluates the output of the last execution to decide if the current execution should proceed.

| Field            | Type                 | Description                                                                                      |
| ---------------- | -------------------- | ------------------------------------------------------------------------------------------------ |
| `flow_id`        | `uint64`             | The ID of the action to fetch the latest response value from (optional).                         |
| `response_index` | `uint32`             | The index of the message response to use (optional).                                             |
| `response_key`   | `string`             | The specific response key to use (e.g., `Amount[0].Amount`, `FromAddress`, optional).            |
| `value_type`     | `string`             | The value type, such as `sdk.Int`, `sdk.Coin`, `sdk.Coins`, `string`, or other compatible types. |
| `operator`       | `ComparisonOperator` | The operator used for comparison (e.g., `==`, `!=`, `<`, `>`).                                   |
| `operand`        | `string`             | The value to compare against.                                                                    |
| `icq_config`     | `ICQConfig`          | The configuration of the Interchain Queryto perform.                                             |

`

### Comparison Operators

The `Comparison Operator` offers several options for evaluating response values across various data types, including strings, arrays, and numeric values. These operators provide flexibility for designing precise and logical conditions in flows.

Imagine you're managing a workflow where decisions must hinge on dynamic and evolving conditions. The `ComparisonOperator` enables fine-grained control to ensure actions proceed only when certain criteria are met. Whether you're verifying numerical thresholds, checking the presence of an item in a list, or ensuring a string starts or ends with specific characters, these operators have you covered.

| Operator        | Description                                                    | Supported Types |
| --------------- | -------------------------------------------------------------- | --------------- |
| `EQUAL`         | Checks for equality between two values.                        | All types       |
| `CONTAINS`      | Verifies if a value is present within a string or array.       | Strings, arrays |
| `NOT_CONTAINS`  | Confirms that a value is not present within a string or array. | Strings, arrays |
| `SMALLER_THAN`  | Evaluates if a value is less than another.                     | Numeric types   |
| `LARGER_THAN`   | Evaluates if a value is greater than another.                  | Numeric types   |
| `GREATER_EQUAL` | Checks if a value is greater than or equal to another.         | Numeric types   |
| `LESS_EQUAL`    | Checks if a value is less than or equal to another.            | Numeric types   |
| `STARTS_WITH`   | Verifies that a string begins with a specified prefix.         | Strings         |
| `ENDS_WITH`     | Verifies that a string ends with a specified suffix.           | Strings         |
| `NOT_EQUAL`     | Checks for inequality between two values.                      | All types       |

These operators ensure that comparisons align with your workflow's logical requirements, offering clarity and precision for any conditional execution scenario.
By leveraging the appropriate comparison operators, you can build intelligent, responsive systems capable of handling even the most complex decision-making requirements.

## Feedback Loops

The `Feedback Loops` feature enables you to replace a value in a message with the latest response value from another action before execution. This creates a feedback loop where the output of one action or Interchain Query becomes the input for another. You can configure up to 5 feedback loops per action.

| Field            | Type        | Description                                                                                      |
| ---------------- | ----------- | ------------------------------------------------------------------------------------------------ |
| `flow_id`        | `uint64`    | The ID of the action to fetch the latest response value from (optional).                         |
| `response_index` | `uint32`    | The index of the responses to use.                                                               |
| `response_key`   | `string`    | The specific response key to use (e.g., "Amount").                                               |
| `msgs_index`     | `uint32`    | The index of the message to replace.                                                             |
| `msg_key`        | `string`    | The key of the message to replace (e.g., `Amount[0].Amount`, `FromAddress`).                     |
| `value_type`     | `string`    | The value type, such as `sdk.Int`, `sdk.Coin`, `sdk.Coins`, `string`, or other compatible types. |
| `icq_config`     | `ICQConfig` | The configuration of the Interchain Queryto perform.                                             |

---

## Tutorial: Conditional Transfers with `MsgSend`

This tutorial demonstrates a reward claim and transfer scenario using intent-based flows.

### Scenario

1. Withdraw staking rewards using `MsgWithdrawDelegatorReward`.
2. Check if the withdrawn amount is greater than 200,000 `uatom`.
3. If the condition is met, transfer the amount to another account using `MsgSend`.

---

### 1. Define the Withdrawal Rewards Message

The withdrawal message is structured as follows:

```ts
const msgWithdrawReward =
  cosmos.distribution.v1beta1.MessageComposer.withTypeUrl.withdrawDelegatorReward(
    {
      delegatorAddress: "cosmos1delegatoraddress",
      validatorAddress: "cosmos1validatoraddress",
    }
  );
```

---

### 2. Define Feedback Loop to Use Withdrawn Amount

The withdrawn amount will be used as input for `MsgSend`:

```ts
const feedbackLoop: FeedbackLoop = {
  flowId: BigInt(0), // Reward withdrawal flow
  responseIndex: 0, // First response index
  responseKey: "amount.[0]", // Extract the withdrawn amount
  valueType: "sdk.Coin", // Value type for replacement
  msgsIndex: 1, // Message index to modify
  msgKey: "amount.[0]", // Key in MsgSend to replace
  icqConfig: undefined,
};
```

---

### 3. Compare the Withdrawn Amount

A `Comparison` condition ensures the withdrawn amount is above the threshold:

```ts
const comparison: Comparison = {
  flowId: BigInt(0), // Reward withdrawal flow
  responseIndex: 0, // First response index
  responseKey: "amount.[0].amount", // Extracted amount key
  valueType: "sdk.Int", // Value type
  operator: 4, // LARGER_THAN
  operand: "200000", // Threshold
  icqConfig: undefined,
};
```

---

### 4. Define the Transfer Flow (`MsgSend`)

If the condition is met, the amount is transferred:

```ts
const msgSend = cosmos.bank.v1beta1.MessageComposer.withTypeUrl.send({
  fromAddress: "cosmos1delegatoraddress",
  toAddress: "cosmos1recipientaddress",
  amount: [{ denom: "uatom", amount: "0" }], // Replaced by feedback loop
});
```

---

### 5. Set Execution Conditions

Execution conditions ensure that the transfer occurs only if the withdrawn amount meets the criteria:

```ts
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

### 6. Submit the Intent-Based Flow

The Intent-based Flow includes both withdrawal and transfer actions:

```ts
const msgSubmitFlow =
  intento.intent.v1beta1.MessageComposer.withTypeUrl.submitFlow({
    label: "Reward Claim and Send Flow",
    owner: "into1wdplq6qjh2xruc7qqagma9ya665q6qhcpse4k6",
    msgs: [msgWithdrawReward, msgSend],
    duration: "1440h",
    interval: "600s",
    feeFunds: [{ denom: "uinto", amount: "5000000" }],
    configuration: config,
    connectionId: "connection-12",
    hostConnectionId: "connection-345",
  });
```

---

### 7. Sign and Broadcast the Transaction

Finally, the transaction is signed and broadcasted:

```ts
client.signAndBroadcast(owner, [msgSubmitFlow], {
  amount: [],
  gas: "300000",
});
```

## Tutorial: Automated Delegation After Staking Rewards Withdrawal

This tutorial follows a similar structure to the first one, but instead of transferring the withdrawn amount to another account, it automatically delegates the rewards to a validator.

### Scenario\*\*

1. Withdraw staking rewards using `MsgWithdrawDelegatorReward`. _(Same as the first tutorial)_
2. Check if the withdrawn amount is greater than 150,000 `uatom`.
3. If the condition is met, delegate the amount to a validator using `MsgDelegate`.

---

### 1. Define Withdrawal Message

This step remains unchanged from the first tutorial:

```ts
const msgWithdrawReward =
  cosmos.distribution.v1beta1.MessageComposer.withTypeUrl.withdrawDelegatorReward(
    {
      delegatorAddress: "cosmos1delegatoraddress",
      validatorAddress: "cosmos1validatoraddress",
    }
  );
```

---

### 2. Define Feedback Loop for Delegation

Instead of using `MsgSend`, the feedback loop will now replace the delegation amount:

```ts
const feedbackLoopDelegation: FeedbackLoop = {
  flowId: BigInt(0), // Reward withdrawal flow
  responseIndex: 0, // First response index
  responseKey: "amount.[0].amount", // Extract the withdrawn amount
  valueType: "sdk.Int", // Type of value
  msgsIndex: 1, // Index in message array to modify
  msgKey: "amount.amount", // Key in MsgDelegate to replace
  icqConfig: undefined,
};
```

---

### 3. Define Comparison for Delegation Threshold

The threshold for delegation is set to 150,000 `uatom`:

```ts
const comparisonDelegation: Comparison = {
  flowId: BigInt(0), // Reward withdrawal flow
  responseIndex: 0, // First response index
  responseKey: "amount.[0]",
  valueType: "sdk.Coin",
  operator: 4, // LARGER_THAN
  operand: "150000uatom", // Threshold
  icqConfig: undefined,
};
```

---

### 4. Define MsgDelegate for Automated Staking

If the condition is met, the withdrawn rewards will be delegated:

```ts
const msgDelegate = cosmos.staking.v1beta1.MessageComposer.withTypeUrl.delegate(
  {
    delegatorAddress: "cosmos1delegatoraddress",
    validatorAddress: "cosmos1validatoraddress",
    amount: { denom: "uatom", amount: "0" }, // Will be replaced by feedback loop
  }
);
```

---

### 5. Set Execution Conditions for Delegation

Execution conditions reference the new feedback loop and comparison:

```ts
const initConditionsDelegation: ExecutionConditions = {
  stopOnSuccessOf: [],
  stopOnFailureOf: [],
  skipOnFailureOf: [],
  skipOnSuccessOf: [],
  feedbackLoops: [feedbackLoopDelegation],
  comparisons: [comparisonDelegation],
  useAndForComparisons: false,
};
```

---

### 6. Submit the Intent-Based Flow for Delegation

This step mirrors the submission process from the first tutorial, with `MsgDelegate` replacing `MsgSend`:

```ts
const msgSubmitFlowDelegation =
  intento.intent.v1beta1.MessageComposer.withTypeUrl.submitFlow({
    label: "Reward Claim and Delegate Flow",
    owner: "into1wdplq6qjh2xruc7qqagma9ya665q6qhcpse4k6",
    msgs: [msgWithdrawReward, msgDelegate],
    duration: "1440h",
    interval: "600s",
    feeFunds: [{ denom: "uinto", amount: "5000000" }],
    configuration: config,
    connectionId: "connection-12",
    hostConnectionId: "connection-345",
  });
```

Just like before, the transaction is signed and broadcasted.
