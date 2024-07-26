---
sidebar_position: 2
title: Conditions
pagination_label: Conditions
---

# Conditions and Comparisons

This page details how to use `ExecutionConditions` and `ResponseComparison` to create complex, conditional workflows in your application. We will explore these features with examples, including an auto-compounding scenario using `MsgSend`.

## Execution Conditions

The `ExecutionConditions` message sets the rules and dependencies for when and how an action is executed. These conditions include replacing values with previous responses, comparing response values, and controlling execution flow based on the success or failure of dependent intents.

### Message Definition

```proto
message ExecutionConditions {
  // Replace value with value from message or response from another action’s latest output
  UseResponseValue use_response_value = 2;
  // Comparison with response response value
  ResponseComparison response_comparison = 1;
  // Optional array of dependent intents that, when executed successfully, stops execution
  repeated uint64 stop_on_success_of = 5;
  // Optional array of dependent intents that, when not executed successfully, stops execution
  repeated uint64 stop_on_failure_of = 6;
  // Optional array of dependent intents that should be executed successfully after their latest call before execution is allowed
  repeated uint64 skip_on_failure_of = 7;
  // Optional array of dependent intents that should fail after their latest call before execution is allowed
  repeated uint64 skip_on_success_of = 8;
}
```

## Feedback Loops

### Using a Response Value as Input

The `UseResponseValue` message allows you to replace a value in a message with the latest response value from another action before execution. This creates feedback loops where the output of one action becomes the input for another.

```proto
message UseResponseValue {
  uint64 action_id = 1 [(gogoproto.customname) = "ActionID"]; // Action to get the latest response value from, optional
  uint32 response_index = 3;  // Index of the response
  string response_key = 2; // For example, "Amount"
  uint32 msgs_index = 4; // Index of the message to replace
  string msg_key = 5; // Key of the message to replace (e.g., Amount[0].Amount, FromAddress)
  string value_type = 6; // Can be anything from sdk.Int, sdk.Coin, sdk.Coins, string, []string, []sdk.Int
}
```

### Response Comparison

The `ResponseComparison` message allows you to compare a response value with a specified operand before the execution of an action. This comparison controls whether the action proceeds based on the outcome.

```proto
message ResponseComparison {
  uint64 action_id = 1 [(gogoproto.customname) = "ActionID"]; // Action to get the latest response value from, optional
  uint32 response_index = 2;  // Index of the response
  string response_key = 3; // e.g., Amount[0].Amount, FromAddress
  string value_type = 4; // Can be anything from sdk.Int, sdk.Coin, sdk.Coins, string, []string, []sdk.Int
  ComparisonOperator comparison_operator = 5;
  string comparison_operand = 6;
}
```

### Comparison Operators

The `ComparisonOperator` enum defines various operators for comparing values. These support different types such as strings, arrays, and numeric types.

```proto
enum ComparisonOperator {
  EQUAL = 0; // Equality check (for all types)
  CONTAINS = 1; // Contains check (for strings, arrays, etc.)
  NOT_CONTAINS = 2; // Not contains check (for strings, arrays, etc.)
  SMALLER_THAN = 3; // Less than check (for numeric types)
  LARGER_THAN = 4; // Greater than check (for numeric types)
  GREATER_EQUAL = 5; // Greater than or equal to check (for numeric types)
  LESS_EQUAL = 6; // Less than or equal to check (for numeric types)
  STARTS_WITH = 7; // Starts with check (for strings)
  ENDS_WITH = 8; // Ends with check (for strings)
  NOT_EQUAL = 9; // Not equal check (for all types)
}
```

## Example: Conditional Transfers with `MsgSend`

Let's illustrate these concepts with an reward claim and send scenario.

### Scenario

1. Withdraw rewards using `MsgWithdrawDelegatorReward`.
2. Check if the withdrawn amount is greater than a threshold (e.g., 200,000 "uatom").
3. If the condition is met, transfer the amount to another account using `MsgSend`.

### Step-by-Step Example

#### 1. Define the Withdrawal Action

First, define an action to withdraw rewards:

```proto
message MsgWithdrawDelegatorReward {
  string delegator_address = 1;
  string validator_address = 2;
}
```

#### 2. Use the Withdrawn Amount as Input

Use the `UseResponseValue` to take the withdrawn amount as input for the next action:

```proto
UseResponseValue {
  action_id: 1 // Action ID of the reward withdrawal
  response_index: 0 // First response
  response_key: "amount[0].amount" // Key to extract the amount
  msgs_index: 1 // Index of the message to replace in the next action
  msg_key: "amount" // Message key to replace in MsgSend
  value_type: "sdk.Int" // Value type
}
```

#### 3. Compare the Withdrawn Amount

Use `ResponseComparison` to ensure the amount is greater than 200,000 "uatom":

```proto
ResponseComparison {
  action_id: 1 // Action ID of the reward withdrawal
  response_index: 0 // First response
  response_key: "amount[0].amount" // Key to compare the amount
  value_type: "sdk.Int" // Value type
  comparison_operator: LARGER_THAN // Operator to check if amount is larger than
  comparison_operand: "200000" // Operand to compare against
}
```

#### 4. Define the Transfer Action

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

```proto
ExecutionConditions {
  use_response_value: {
    action_id: 1 // Action ID of the reward withdrawal
    response_index: 0 // First response
    response_key: "amount[0].amount" // Key to extract the amount
    msgs_index: 1 // Index of the message to replace
    msg_key: "amount" // Message key to replace
    value_type: "sdk.Int" // Value type
  }
  response_comparison: {
    action_id: 1 // Action ID of the reward withdrawal
    response_index: 0 // First response
    response_key: "amount[0].amount" // Key to compare the amount
    value_type: "sdk.Int" // Value type
    comparison_operator: LARGER_THAN // Operator to check if amount is larger than
    comparison_operand: "200000" // Operand to compare against
  }
}
```

With these conditions, the `MsgSend` action will only execute if the withdrawn amount is greater than 200,000 "uatom", and the withdrawn amount will be used as the transfer amount.