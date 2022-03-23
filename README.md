# The Trustless Hub

![Welcome to the Trustless Hub](docs/images/web.png)

**The Trustless Hub (TRST)**  is an innovative smart contract hub built using the Cosmos SDK and Tendermint BFT consensus. The main token, TRST, is used to keep the Hub secure, pay for transaction fees and can be used as collateral in voting for proposals. On the Trustless Hub, smart contracts benefit from maximum programmability. Features include encrypted inputs; a private, encrypted state and a public state; automated messages and encrypted or verifiable outputs. 

The Trustless Hub aims to be the most innovative smart contract platform across all blockchain ecosystems. The Trustless Hub enables smart contracts 2.0, called Trustless Contracts. Developers will be able to use and integrate Trustless Contracts for new privacy-first services and use-cases that are currently not possible. Use-cases include NFTs, pricing and billing, yet this can be related to anything. 

There are 4 pillars to Trustless Contracts:
Private inputs
Secure execution (using Intel SGX Secure enclaves)
Verifiable or encrypted outcomes
Automatic contract results

Any CosmWasm code, like smart contracts on Terra and Secret Contracts on the Secret Network can easily be upgraded to Trustless Contract code, with minor changes.


## What is TRST?

TRST is the native token of the Trustless Hub and is used to pay for gas fees. It is used for:
Transactions
* Storing, instantiating and executing Trustless Contracts.
* Estimating and transacting with unique assets.
* Securing the hub through staking


You can stake your TRST by delegating your TRST to validators. Trough governance, you can vote on important proposals such as upgrades and parameter changes. 


## Why Trustless Contracts?
Any CosmWasm code, like smart contracts on Terra and Secret Contracts on the Secret Network can easily be transformed into a Trustless Contract.
New trustless services are now possible, like auto-ending private auctions. CosmWasm smart contracts are present in major Cosmos-based blockchains like Terra, Crypto.com, Secret Network and Irisnet, and allow for developers to easily create smart contracts in the same way across different blockchains. For example, if a developer learns to program on Terra, creating Trustless Contracts on the Trustless Hub is as simple as ABC. A developer can reuse the code, enjoy private computing and automated, verifiable outcomes.

Use cases that leverage private inputs and automatic, verifiable outputs are now possible on the Trustless Hub, such as:
- Encrypted and secured data aggregation, with verifiable outcomes
- Auto-ending secret bid auctions 
- Estimation aggregation for generating independent prices
- A private DEX with auto-swap capabilities
- Automatic private transaction relayers
- Conditional and time-dependent token transfers


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




