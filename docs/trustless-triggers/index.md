---
order: 0
parent:
  title: Trustless Triggers
  order: 1
---

# Introduction

This knowledge base contains information that cover different aspects of automation of transactions on Trustless Hub.

In this module we use `Trustless Triggers` and `AutoTX` interchangeably. An AutoTX is what we call a Trustless Trigger and these are the same thing.

The AutoTX module is responsible for creating and executing automatic interchain transactions between different chains within the Cosmos ecosystem.

Transaction automation using AutoTX module is a relatively simple time-based automation module as opposed to our advanced contract automation that can handle arbitrary messages. 

## Advantages of our approach

+ Easy to integrate into dApps
+ Automation of user funds using AuthZ
+ Schedule multiple messages into one AutoTX that execute after eachother
+ You can depend execution on other transactions
+ No need to integrate with bots
+ Fee funds are  refunded after execution finishes
+ Relayers are incentivized for acknoledging a succesful automation packet on a host chain