---
sidebar_position: 4
title: Interchain Query
description: Interchain Query Feature for Intent-Based Actions
---

The **Interchain Query (ICQ) feature** enables seamless cross-chain interactions by querying the state of one blockchain from another. This feature allows developers to create intent-based actions that trigger specific behaviors based on the queried data. By utilizing interchain queries, you can automate decision-making processes across multiple blockchains, providing an efficient way to orchestrate complex, multi-chain applications.

### How It Works

Interchain queries allow you to access the key-value store of a different blockchain by providing a specific key to query. IBC relayers submit query responses. This queried state can then be used to determine what actions should be triggered on your own chain, based on predefined conditions. For instance, you can check the balance of a particular account on one blockchain and then execute a corresponding action on another chain, such as staking, transferring, or adjusting governance proposals.

### Supported Types and Proto Interfaces

You can query various data types from the state of other chains. There are two primary methods to accomplish this:

1.  **Supported Types:** You can use one of the [supported types](../module/supported_types.md), which are predefined and commonly used data types within the Cosmos ecosystem.
2.  **Registered Proto Interfaces:** Alternatively, you can utilize the registered protocol buffer interfaces that adhere to Cosmos SDK standards. These provide a more flexible way to query and interpret the state, allowing you to work with complex data structures.

### Feedback Loops and Comparisons

Once the queried data is retrieved, it can be used for [**comparisons**](./conditions.md#comparison-operators) or to establish [**feedback loops**](./conditions.md#feedback-loops). These operate similarly to the response mechanisms available with message-based transactions. For example, if a queried balance exceeds a certain threshold, you can trigger an action to stake the excess funds. Likewise, if a validator’s status on one chain changes, you can automatically adjust delegations or governance votes on another chain.

This capability unlocks numerous possibilities for cross-chain automation, simplifying multi-chain dApp logic and empowering developers to build more dynamic and responsive applications in the interchain ecosystem.

## User Stories

#### 1. Balance Check for Friday Action

**As a** blockchain participant, I want the system to check my ATOM balance every Friday and perform a predefined action if my balance exceeds 100.

#### 2. Action Based on DAO Vote (Negative Outcome)

**As a** DAO member, I want the system to monitor DAO vote outcomes every Wednesday and trigger a specific action if the result is negative.

#### 3. Validator Set Monitoring

**As a** validator, I want the system to check my validator status daily for one year and automatically unbond my stake if I am removed from the set.

#### 4. Automatic Staking if ATOM Balance is Above Threshold

**As a** token holder, I want the system to automatically stake any ATOM balance exceeding 100 every Friday to maximize rewards.

#### 5. Alert for Validator Removal After One Year

**As a** blockchain participant, I want the system to monitor my validator status for one year and alert me if my validator is out of the set, triggering an unbonding action.

#### 6. DAO Vote Rejection and Action Trigger

**As a** project lead in a DAO, I want the system to trigger a review process if a vote is rejected on Wednesday to adjust the proposal or take alternative steps.

#### 7. Balance-Based Action Adjustment

**As a** token holder, I want the system to scale actions based on my ATOM balance, staking or transferring the amount that exceeds 100 ATOM every Friday.

#### 8. Emergency Action for Validator Removal

**As a** validator, I want the system to automatically unbond my stake and transfer my funds if I am removed from the validator set after one year to avoid penalties.

#### 9. Weekly Balance Check for Alternative Action

**As a** DeFi participant, I want the system to check my ATOM balance every Friday and trigger borrowing or selling assets if my balance is below 100 ATOM.

#### 10. DAO Vote and Treasury Allocation

**As a** DAO member, I want the system to monitor negative vote outcomes every Wednesday and automatically reallocate funds or adjust the treasury in line with governance rules.

## Integrating ICQs

You can use the ICQ feature by attaching an ICQConfig into the Conditions object and setting fromICQ to true in the UseValueResponse or ResponseComparison objects.
In the ICQConfig you specify what to query and where and how to handle a timeout scenario.
In the UseValueResponse or ResponseComparison you may additionally specify a field to look into, if the queried value is an object.

```protobuf

//config for using interchain queries
message ICQConfig {
  string connection_id = 1;
  string chain_id = 2;
  intento.interchainquery.v1.TimeoutPolicy timeout_policy = 3;
  google.protobuf.Duration timeout_duration = 4 [
    (gogoproto.nullable) = false,
    (gogoproto.stdduration) = true
  ];
  string query_type = 5; // e.g. store/bank/key store/staking/key
  string query_key = 6; //key in the store to query e.g. stakingtypes.GetValidatorKey(validatorAddressBz) idea: abstract this in TP. See x/interchainquery/types/keys.go
}
For example, query_type can be store/bank/key store/staking/key
  The query_key is the key in the store to query such as stakingtypes.GetValidatorKey(validatorAddressBz) We abstract this away and provide examples in the TriggerPortal frondend. See x/interchainquery/types/keys.go for more examples.

```

## x/interchainqueries Msgs

SubmitQueryResponse is used to return the query responses.

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

## x/interchainqueries Queries

Query PendingQueries lists all queries that have been requested (i.e. emitted but have not had a response submitted yet

```protobuf
message QueryPendingQueriesRequest {}
```
