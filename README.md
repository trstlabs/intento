# Trustless Hub

![Welcome to Trustless Hub](docs/images/web.png)

**Trustless Hub (TRST)**  is an automation-first smart contract hub built using the Cosmos SDK and Tendermint BFT consensus. The main token, TRST, is used to keep the Hub secure, pay for automatic execution, general transactions and can be used as collateral in voting for proposals. 

Trustless Hub is built on the Cosmos SDK with IBC for secure cross-chain communication. Trustless Hub enables new services using private inputs, time-based messages, 1-click recurring execution and verifiable outcomes.

Use-cases include recurring swaps for DeFi, automatic reward claiming, instant settement of results, 1-click recurring transactions for payment and in-game battles for gaming.

Features include

 1. Encrypted inputs
 2. Private, encrypted state a public state
 3. Time-based messages 
 4. Encrypted or verifiable outputs
 5. Front-running resistance

Trustless Hub aims to be the most innovative smart contract platform across all blockchain ecosystems. Trustless Hub enables smart contracts 2.0, called Trustless Contracts. Developers will be able to use and integrate Trustless Contracts for new privacy-first services and use-cases that are currently not possible. Use-cases include NFTs, pricing and billing, yet this can be related to anything. 

There are 4 pillars to Trustless Contracts:
Private inputs
Secure execution (using Intel SGX Secure enclaves)
Verifiable or encrypted outcomes
Automatic contract results

Any CosmWasm code, like smart contracts on Terra and Secret Contracts on the Secret Network can easily be upgraded to Trustless Contract code, with minor changes.


## What is TRST?

TRST is the native token of Trustless Hub and is used to pay for gas fees. It is used for:
Transactions

* Storing, instantiating and executing Trustless Contracts.
* Automatic contract execution fees
* Securing the hub

You can stake your TRST by delegating your TRST to validators. Trough governance, you can vote on important proposals such as upgrades and parameter changes. 


## Why Trustless Contracts?

Smart contracts enable decentralized applications, where data is immutable and users can trust the data to be real and uncensored. Still, smart contracts are far from perfect as current smart contracts on major blockchains are fully public, data is stored indefinitely. Right now, developers cannot guarantee any kind of privacy when interacting with their contracts. Anyone can view user inputs. Users lose funds because smart contracts have no predefined end date. Stakeholders may not develop enough trust to interact with the contract as the data is public and the contract end date is unknown.

Yet, so far no easy way exists to build contracts with private user inputs, deterministic outcomes and automatic time-based actions to cope with these problems.

In order to truly build trustless applications for DAOs and end-users, developers should be able to program contracts that have automatic time-based functions.

Chains and CosmWasm-based apps extend their functionality by integrating with Trustless Hub over IBC. Trustless Contracts enable existing apps to have automatic execution and front-running resistance. Out of the box, recurring transactions, remittances and payroll are possible. With allowance, recipients can interact and transfer funds to a privacy-perserving balance.

## Get started

``` bash
make deb
./init.sh
```

`make deb` makes a package (linux only, with SGX enabled) with dependencies. `./init.sh` initializes and starts the node in development.

## Configure

The node can be configured with `config.yml`. To learn more see the [reference](https://github.com/tendermint/starport#documentation).

## Learn more

- [Trustless Hub website](https://trustlesshub.com/)
- [Trustless Hub wallet](https://interact.trustlesshub.com/)
- [Cosmos SDK documentation](https://docs.cosmos.network)
- [Cosmos SDK Tutorials](https://tutorials.cosmos.network)

