---
title: Modules
order: 3
---

# Overview

Intento uses general Cosmos SDK modules, modules required to run a Dymension RollApp as well as a custom Cosmos SDK module for processing intent-based actions.

The Intent module in Version 2 will enable actions to be highly configurable and conditional, whereby the actions can depend on execution results. For example, a protocol user could set up a sequence of actions such as swapping ATOM for USDC on Osmosis and then paying for a subscription that is settled on Ethereum using the Axelar bridge with General Message Passing. By enabling time-based intents, protocols and their end-users can automate complex workflows in a seamless manner.

### Fee Mechanism
    
Based on lines of code, we expect the module to be using 100,000 gas for triggering an execution. Significantly less than Gelato on EVM chains (1,000,000) and CronCat (700,000) whilst bringing trust assumptions to the minimum. This makes the module highly scalable for any specified intent.

| Gas Operation                        | enum | Gas          | Description |
|--------------------------------------|------|--------------|-------------|
| GasOperationResultContingentFlat     | 0    | 1,000        | If `[0-9]{1,};` is contingent on an execution result (StopOnSuccess, SkipOnSuccess, StopOnFailure, SkipOnFailure, StopOnSuccessOf, SkipOnSuccessOf, StopOnFailureOf, SkipOnFailureOf, UseValue) |
| GasOperationRetrieveActionByID       | 1    | 5,000        | For retrieving an `[0-9]{1,};` by ID, for `[0-9]{1,};` configurations that have StopOnSuccessOf, SkipOnSuccessOf, StopOnFailureOf, SkipOnFailureOf, Comparison, and UseValue operations. |
| GasOperationComparisonFlat           | 2    | 5,000        | Flat fee for using a comparison |
| GasOperationUseValueFlat             | 3    | 5,000        | Flat fee for using value as inputs |
| GasOperationMessagePerByte           | 4    | default to 0 | Optional, Cosmos messages are usually 100-400 characters long. CosmWasm and Interchain operations generally use more message bytes. |
| GasOperationReregisterICAAfterTimeout | 5    | 1,000        | Optional, from Intent configuration |

# Actions

Actions are specified in an Action object that contains data about when to execute, what to execute by an array of registered cosmos messages, where through an optional ICA Configuration, and how through a configuration and conditions.

## Technical Specification

### Configuration

```proto
// ExecutionConfiguration provides the execution-related configuration of the Action
message ExecutionConfiguration {
       // if true, the Action outputs are saved and can be used in condition-based logic
      bool save_msg_responses = 1;
      // if true, the Action is not updatable
      bool updating_disabled = 2;           
      // If true, will execute until we get a successful Action, if false/unset will always execute
      bool stop_on_success = 3;
      // If true, will execute until successful Action, if false/unset will always execute
      bool stop_on_failure = 4;
      // If true, owner account balance is used when trigger account funds run out
      bool fallback_to_owner_balance = 5;
      // If true, allows the Action to continue execution after an ibc channel times out (recommended)
      bool reregister_ica_after_timeout = 6 [(gogoproto.customname) = "ReregisterICAAfterTimeout"];
}
```

## Parameters

A number of action-related governance parameters can be adjusted. Parameters can be adjusted by governance to ensure that fees and rewards are appropriate. The default values are the following:

```golang
const (
 DefaultActionFundsCommission int64 = 2 //2%

 DefaultMaxActionDuration time.Duration = time.Hour * 24 * 366 * 2 // a little over 2 years

 DefaultMinActionDuration time.Duration = time.Second * 60 //1minute

 DefaultMinActionInterval time.Duration = time.Second * 60 //1minute
)

```

Gas-related parameters are also to-be governance controlled.

<!-- 
### Conditions

#### Using an event value

Let's explain conditions by examining what would be required for performing an auto-compound operation. The *action* involves withdrawing rewards and then staking these rewards again, so the rewards will compound. We would need to use the event value from the withdrawl message and use that in the new message.

-   **Attribute Extraction**: Upon a _withdraw_rewards_ event, the module needs to extract a specific attribute value associated with the key validator. This attribute contains information necessary for the next step in the autocompound process.
-   **Parsing and Message Building:** The extracted attribute value, which is expected to be in a format that can be parsed as a coin, is then used to construct a new message. This message initiates the autocompound action, effectively reinvesting the rewards back into the module.
-   **Retrieving Output from Events**: The process involves specifying _attribute_key_ is related to the values needed for the operation. For example, similar to how one would retrieve the amount from a transfer event by specifying the _attribute_key_ as amount, the module should employ a similar mechanism to fetch necessary data from the withdraw_rewards event or any other relevant event.
-   **Action ID Specification**: In scenarios where the event data needed for autocompounding is from an event triggered by another action, the optional _action_id_ may be specified. This allows the module to precisely identify and use the data from the correct event, ensuring that the auto compound action is based on the most relevant and accurate information.


```proto
// ExecutionConditions provides execution conditions for the Action
message ExecutionConditions {
    // Comparison runs before execution of Action and outputs true or false
    EventComparision event_comparison = 1;
    // Replace value with value from message or response from another Action’s latest output
    UseEventAttributeValue use_event_value = 2;
    //optional array of dependent intents that when executing succesfully, stops execution
    repeated int64 stop_on_success_of = 5;	
    //optional array of dependent intents that when not executing succesfully, stops execution
    repeated int64 stop_on_failure_of = 6;     
    //optional array of dependent intents that should be executed succesfully after their latest call before execution is allowed
    repeated int64 skip_on_failure_of = 7;
    //optional array of dependent intents that should fail after their latest call before execution is allowed
    repeated int64 skip_on_success_of = 8;
}
```

You can retrieve the output received by specifying the **attribute_key** related to the desired value. For example, for getting an event value from the typical transfer attribute.

event key is **transfer**, the **attribute_key** is **amount**. The **action_id** may be specified if the event from another intent should be used.


```proto
message UseEventAttributeValue {
  string event_key = 1;
  string action_index = 2;
  string action_key = 3;
  string attribute_key = 4;
  optional string action_id = 5;
  optional Match match = 6;
}
```


This can be troublesome when the event has multiple transfer attributes. To ensure you get the desired key you can specify a Match for the spender or recipient key with an address. The type for comparison could be specified if it should be handled differently than a string or coin.

```json
{
  "type": "transfer",
  "attributes": [
    {
      "key": "recipient",
      "value": "osmo1g7ajkk295vactngp74shkfrprvjrdwn662dg26"
    },
    {
      "key": "sender",
      "value": "osmo17a8smrhauph552zkz5864vjafz9pszpezepz68"
    },
    {
      "key": "amount",
      "value": "145506ibc/1480B8FD20AD5FCAE81EA87584D269547DD4D436843C1D20F15E00EB64743EF4"
    }
  ]
}
```

```proto
message Match {
  string attribute_key = 4;
  string match = 4;
}
```

The match iterates over the event keys when there are multiple. When the match is found, the comparison will be made.

To use the attribute value of an event as input for the new message. The action key can be a key from anywhere in the JSON parsed Cosmos SDK message. Must be readable by proto. For CosmWasm/EVM a decoder would be necessary, this can be included in a future version.


```proto
message EventComparision {
  string event_key = 1; 
  ComparisonOperator comparision_operator = 2;
  string comparison_value = 3;
  optional string attribute_key = 4; 
  optional string action_id = 5;
  optional Match match = 6;
}
```

To use the attribute value of an event as input for the new message. The action key can be a key from anywhere in the JSON parsed Cosmos SDK message. Must be readable by proto. For CosmWasm/EVM a decoder would be necessary, this can be included in a future version.
Comparison operators that can be used for various types.

```proto
enum ComparisonOperator {
  EQUAL = 0;          // Equality check (for all types)
  CONTAINS = 1;       // Contains check (for strings, arrays, etc.)
  NOT_CONTAINS = 2;   // Not contains check (for strings, arrays, etc.)
  SMALLER_THAN = 3;    // Less than check (for numeric types and coins)
  LARGER_THAN = 4;     // Greater than check (for numeric types and coins)
  GREATER_EQUAL = 5;   // Greater than or equal to check (for numeric types and coins)
  LESS_EQUAL = 6;      // Less than or equal to check (for numeric types and coins)
  STARTS_WITH = 7;     // Starts with check (for strings)
  ENDS_WITH = 8;       // Ends with check (for strings)
  NOT_EQUAL = 9;       // Not equal check (for all types)
}
```



Let’s illustrate the same example from above, comparing an event value from the typical transfer attribute. The **event_key** is **transfer**, the **attribute_key** is **amount**. We compare it with 500uosmo. The comparison operator can be any that is compatible with a numeric(coin) type. Automatically in the background, the string is parsed as coins and it will only compare when the value also is in **uosmo**.
When we expect there to be more transfer event keys, we can specify it using the match. We can match by the spender key again.

For Interchain Queries we can implement a similar structure. Due to the added complexity, in development and also in testing and auditing, we leave this out of scope but still we are excited to implement this after the grant work has been completed. With interchain queries we can allow comparisons with pool balances and oracle prices. For example Skip’s slinky oracle aggregator deployed on osmosis. With a similar structure we can look 1 level deep which is sufficient.  We can retrieve GetPriceResponse, then with a similar attribute_key we can point to price, which points to the price. We can then compare it to a comparision_value. 

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
``` -->

### Privileged Host Chain Execution

By using the Cosmos message type MsgExec of the AuthZ module in your intent, you can allow the intent address to execute any message on your behalf. This is needed for most use cases where you want to automate your own balance, such as recurring payments. 
For this it is important that the Intent address gets granted these privileges. These can be given by sending a MsgGrant with the typeUrl of the message to execute on the host chain. Front-end tools like TriggerPortal make this process easy and seamless. 
