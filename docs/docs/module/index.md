---
title: Module
order: 3
---

# Overview

Intento uses general Cosmos SDK modules, modules required to run a rollup as well as a custom Cosmos SDK module for processing intent-based actions.

The Intent module enables actions to be highly configurable and conditional, whereby the actions can depend on execution results. For example, a protocol user could set up a sequence of actions such as swapping ATOM for USDC on Osmosis and then paying for a subscription that is settled on Ethereum using the Axelar bridge with General Message Passing. By enabling these user intents, protocols and their end-users can automate complex workflows in a seamless manner.

## Technical Specification

### Gas operations

Based on lines of code, we expect the module to be using 80,000 to 100,000 gas for triggering up to 10 executions. Significantly less than Gelato on EVM chains (1,000,000) and CronCat (700,000) whilst bringing trust assumptions to the minimum. This makes the module highly scalable for any specified intent.

# Actions

Actions are specified in an action object that contains data about when to execute, what to execute by an array of registered cosmos messages, where through an optional ICA Configuration, and how through a configuration and conditions.

### Configuration

```proto
// ExecutionConfiguration provides the execution-related configuration of the action
message ExecutionConfiguration {
       // if true, the action outputs are saved and can be used in condition-based logic
      bool save_msg_responses = 1;
      // if true, the action is not updatable
      bool updating_disabled = 2;
      // If true, will execute until we get a successful Action, if false/unset will always execute
      bool stop_on_success = 3;
      // If true, will execute until successful Action, if false/unset will always execute
      bool stop_on_failure = 4;
      // If true, owner account balance is used when trigger account funds run out
      bool fallback_to_owner_balance = 5;
      // If true, allows the action to continue execution after an ibc channel times out (recommended)
      bool reregister_ica_after_timeout = 6 [(gogoproto.customname) = "ReregisterICAAfterTimeout"];
}
```

## Parameters

A number of action-related governance parameters can be adjusted. Parameters can be adjusted by governance to ensure that fees and rewards are appropriate. The default values are the following:

```golang
const (
 ActionFundsCommission int64 = 2 //2%

 ActionConstantFee int64 = 0 // e.g. 5_000 in default denom

 GasFeeCoins sdk.Coins = sdk.NewCoins(sdk.NewCoin(Denom, sdk.NewInt(1))) // 1uinto

 MinActionDuration time.Duration = time.Second * 60 //1minute

 MinActionInterval time.Duration = time.Second * 60 //1minute

MaxActionDuration time.Duration = time.Hour * 24 * 366 * 10 // a little over 10 years
)
```

## Execution Conditions

The `ExecutionConditions` message defines a set of rules and dependencies that determine when and how an action is executed. It includes conditions for using values from previous responses, comparing response values, and controlling execution flow based on the success or failure of dependent intents.

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

### Feedback Loops

#### Using a Response Value as Input

The `UseResponseValue` message allows you to replace a value in a message with the latest response value from another action before execution. This is useful for creating feedback loops where the output of one action is used as input for another.

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

For example, consider a response type `/cosmos.distribution.v1beta1.MsgWithdrawDelegatorRewardResponse`, where the output might look something like this:

```json
{
  "amount": [{ "denom": "uatom", "amount": "211816" }]
}
```

You can use `UseResponseValue` to take the "amount" field from this response and use it as an input in another action.

#### Response Comparison

The `ResponseComparison` message allows you to compare a response value with a specified operand before the execution of an action. This comparison can be used to control whether the action proceeds based on the outcome of the comparison.

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

#### Comparison Operators

The `ComparisonOperator` enum defines various operators that can be used to compare values. These operators support different types such as strings, arrays, and numeric types.

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

### Example

Let's illustrate the above concepts with an example:

Suppose you have an action that withdraws rewards and you want to use the withdrawn amount as input for another action, but only if the amount is greater than 200,000 "uatom".

1. **Using Response Value:**

   ```proto
   UseResponseValue {
     action_id: 12345 // Action ID of the reward withdrawal
     response_index: 0 // First response
     response_key: "amount[0].amount" // Key to extract the amount
     msgs_index: 0 // Index of the message to replace
     msg_key: "Amount" // Message key to replace
     value_type: "sdk.Int" // Value type
   }
   ```

2. **Response Comparison:**

   ```proto
   ResponseComparison {
     action_id: 12345 // Action ID of the reward withdrawal
     response_index: 0 // First response
     response_key: "amount[0].amount" // Key to compare the amount
     value_type: "sdk.Int" // Value type
     comparison_operator: LARGER_THAN // Operator to check if amount is larger than
     comparison_operand: "200000" // Operand to compare against
   }
   ```

With these conditions, the subsequent action will only execute if the withdrawn amount is greater than 200,000 "uatom", and the amount will be used as an input for this action.

<!-- For Interchain Queries we can implement a similar structure. Due to the added complexity, in development and also in testing and auditing, we leave this out of scope but still we are excited to implement this after the grant work has been completed. With interchain queries we can allow comparisons with pool balances and oracle prices. For example Skip’s slinky oracle aggregator deployed on osmosis. With a similar structure we can look 1 level deep which is sufficient. We can retrieve GetPriceResponse, then with a similar attribute_key we can point to price, which points to the price. We can then compare it to a comparision_value. -->

<!--
### Creating Intents

```proto
message MsgCreateIntent {
	//Set of actions
}
```

Intents are a collection of actions to be processed. As a prerequisite for submitting intents over IBC, the intent creator should have interchain accounts registered for the host chains.

```proto
message MsgRegisterICASAndSubmitIntent {
	//Set of actions
	//IBC version
}
```

### Privileged Host Chain Execution

By using the Cosmos message type MsgExec of the AuthZ module in your intent, you can allow the intent address to execute any message on your behalf. This is needed for most use cases where you want to automate your own balance, such as recurring payments.
For this it is important that the Intent address gets granted these privileges. These can be given by sending a MsgGrant with the typeUrl of the message to execute on the host chain. Front-end tools like TriggerPortal make this process easy and seamless.  -->
