# The Trustless Hub

![Welcome to the Trustless Hub](docs/images/web.jpeg)

**The Trustless Hub (TRST)** is a smart contract hub built using the Cosmos SDK and Tendermint BFT consensus. The main token, TRST, is used to keep the Hub secure, pay for transaction fees and can be used as collateral in voting for proposals. The Trustless Hubs run CosmWasm contracts on Intel SGX Secure Enclaves, and provides flexibility in contracts.

On the Trustless Hub, smart contracts with private inputs and encrypted or verifiable outputs are possible. In addition, we introduce automated messages, these enables new DeFi-related services.

The Trustless Hub aims to be the most innovative smart contract platform across all blockchains. The Trustless Hub enables smart contracts 2.0, called Trustless Contracts. Developers will be able to use and integrate Trustless Contracts for new services and use-cases that are currently not possible. Use-cases include NFTs, pricing and billing, yet this can be related to anything.
There are 4 pillars to Trustless Contracts:
1. Private inputs
2. Secure execution (using Intel SGX Secure enclaves)
3. Verifiable or encrypted outcomes
4. Automatic contract results

Any CosmWasm code, like smart contracts on Terra and Secret Contracts on the Secret Network can easily be upgraded to Trustless Contract code, with minor changes.

## What is TRST?

TRST is the native token of the Trustless Hub and is used to pay for gas fees. It is used for:

* Transactions
* Storing, instantiating and executing Trustless Contracts. 
* Estimating and transacting with unique assets. 

You can also stake your TRSt by delegating your TRST to validators. Trough governance, you can vote on important proposals such as upgrades and parameter changes. 


## Why Trustless Contracts?
The Trustless Hub aims to be the most innovative smart contract platform across all blockchains. The Trustless Hub enables smart contracts 2.0, called Trustless Contracts. Developers will be able to use and integrate Trustless Contracts for services and use-cases that are currently not possible. Many of use-cases are related to NFTs, pricing and billing, yet this can be for anything.


There are 4 pillars to Trustless Contracts:
1. Private inputs
2. Secure execution (using Intel SGX Secure enclaves)
3. Automatic result generation and finalization
4. Verifiable outcomes

Any CosmWasm code, like a smart contracts on Terra and Secret Contracts on the Secret Network can easily be transformed into a Trustless Contract.

New trustless services are now possible, like auto-ending private auctions. CosmWasm smart contracts are present in major Cosmos-based blockchains like Terra, Crypto.com, Secret Network and Irisnet, and allow for developers to easily create smart contracts in the same way across different blockchains. For example, if a developer learns to program on Terra, creating Trustless Contracts on the Trustless Hub is as simple as ABC. A developer can reuse the code, enjoy private computing or add time-dependent functions. Cryptocurrencies like Terra’s stablecoin UST can be transferred onto any Trustless Contract using IBC, making the switch a breeze. If you’re unfamiliar with CosmWasm, there are many [tutorials out there](https://www.youtube.com/results?search_query=CosmWasm).
By deploying [Cosmwasm](https://cosmwasm.com/) code, you can use Trustless Contracts for:

Use cases that leverage private inputs and automatic, verifiable outcomes are now possible on the Trustless Hub, such as:
- Trustless data aggregation
- Auto-ending secret bid auctions 
- Estimation aggregation for generating independent prices
- A private DEX which can auto-swap back tokens
- Automatic private transaction relayers
- Conditional token and NFT transfers


## Get started

```
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




