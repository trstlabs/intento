---
sidebar_position: 2
title: Conditions
pagination_label: Conditions
---

In this part of the docs we detail how to use `Conditions` such as `Comparisons` and `Feedback Loops` into your Intent-based Flow. With Conditions on Intento you can orchestrate conditional workflows in a structured manner. We will explore these features with examples, including an auto-compounding scenario using `MsgSend`.

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

The `Comparison` feature allows you to compare a response value with a specified operand before executing an action message. This comparison determines whether the action proceeds based on the result. You can configure up to 5 conditions per action. It evaluates the output of the last execution to decide if the current execution should proceed.

| Field            | Type                 | Description                                                                                      |
| ---------------- | -------------------- | ------------------------------------------------------------------------------------------------ |
| `action_id`      | `uint64`             | The ID of the action to fetch the latest response value from (optional).                         |
| `response_index` | `uint32`             | The index of the message response to use (optional).                                             |
| `response_key`   | `string`             | The specific response key to use (e.g., `Amount[0].Amount`, `FromAddress`, optional).            |
| `value_type`     | `string`             | The value type, such as `sdk.Int`, `sdk.Coin`, `sdk.Coins`, `string`, or other compatible types. |
| `operator`       | `ComparisonOperator` | The operator used for comparison (e.g., `==`, `!=`, `<`, `>`).                                   |
| `operand`        | `string`             | The value to compare against.                                                                    |
| `icq_config`     | `ICQConfig`          | The configuration of the Interchain Query (ICQ) to perform.                                      |

`

### Comparison Operators

The `Comparison Operator` offers several options for evaluating response values across various data types, including strings, arrays, and numeric values. These operators provide flexibility for designing precise and logical conditions in execution flows.

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
| `action_id`      | `uint64`    | The ID of the action to fetch the latest response value from (optional).                         |
| `response_index` | `uint32`    | The index of the responses to use.                                                               |
| `response_key`   | `string`    | The specific response key to use (e.g., "Amount").                                               |
| `msgs_index`     | `uint32`    | The index of the message to replace.                                                             |
| `msg_key`        | `string`    | The key of the message to replace (e.g., `Amount[0].Amount`, `FromAddress`).                     |
| `value_type`     | `string`    | The value type, such as `sdk.Int`, `sdk.Coin`, `sdk.Coins`, `string`, or other compatible types. |
| `icq_config`     | `ICQConfig` | The configuration of the Interchain Query (ICQ) to perform.                                      |

## Example: Conditional Transfers with `MsgSend`

Let's illustrate these concepts with an reward claim and send scenario.

### Scenario

1. Withdraw rewards using `MsgWithdrawDelegatorReward`.
2. Check if the withdrawn amount is greater than a threshold (e.g., 200,000 "uatom").
3. If the condition is met, transfer the amount to another account using `MsgSend`.

### Step-by-Step Example

#### 1. Define the Withdrawal Flow

First, define an action to withdraw rewards:

```proto
message MsgWithdrawDelegatorReward {
  string delegator_address = 1;
  string validator_address = 2;
}
```

#### 2. Use the Withdrawn Amount as Input

Use the `FeedbackLoop` to take the withdrawn amount as input for the send message:

```js
feedback_loops: [{
  action_id: 0 // Optional action ID of the reward withdrawal
  response_index: 0 // First response
  response_key: "amount[0].amount" // Key to extract the amount
  msgs_index: 1 // Index of the message to replace in the next action
  msg_key: "amount" // Message key to replace in MsgSend
  value_type: "sdk.Int" // Value type
}]
```


#### 3. Compare the Withdrawn Amount

Use `Comparison` to ensure the amount was greater than 200,000 "uatom":

```js
comparisons: [{
  response_index: 0 // First response
  response_key: "amount[0].amount" // Key to compare the amount
  value_type: "sdk.Int" // Value type
  comparison_operator: LARGER_THAN // Operator to check if amount is larger than
  comparison_operand: "200000" // Operand to compare against
}]
```

This will compare previous output. Feedback Loops run throughout the flow and can use message 1 outputs for constructing message 2. Comparisons on the other hand, happen at the begining of the flow and determine whether the flow should be executed or not. If you want to use  up-to-date outputs, you can configure the messages into distict Intent-based Flows and reference the withdrawal message from the send action like so:

```js
comparisons: [{
  action_id: 2 // Optional action ID of the reward withdrawal
  ...//other fields
}]

#### 4. Define the Transfer Flow

Define an action to transfer the amount using `MsgSend`:

```proto
message MsgSend {
  string from_address = 1;
  string to_address = 2;
  sdk.Coin amount = 3; // This will be replaced with the withdrawn amount
}
```

#### 5. Set Execution Conditions

Combine the conditions into `ExecutionConditions` for the transfer action:

```js
conditions {
  feedback_loops: [{
    action_id: 0 // Optional action ID of the reward withdrawal
    response_index: 0 // First response
    response_key: "amount[0].amount" // Key to extract the amount
    msgs_index: 1 // Index of the message to replace
    msg_key: "amount" // Message key to replace
    value_type: "sdk.Int" // Value type
  }],
  comparisons: [{
    action_id: 0 // Optional action ID of the reward withdrawal
    response_index: 0 // First response
    response_key: "amount[0].amount" // Key to compare the amount
    value_type: "sdk.Int" // Value type
    comparison_operator: LARGER_THAN // Operator to check if amount is larger than
    comparison_operand: "200000" // Operand to compare against
  }]
}
```

With these conditions, the `MsgSend` action message will only execute if the withdrawn amount was greater than 200,000 "uatom", and the withdrawn amount will be used as the transfer amount.
