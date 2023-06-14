# Trustless Hub

![Welcome to Trustless Hub](docs/images/web.png)

**Trustless Hub** contains custom Cosmos SDK modules for us to develop, test and launch an automation-first blockchain hub built using the Cosmos SDK and Tendermint BFT consensus. 

The main token, TRST, is used to keep the Hub secure, pay for automatic execution, general transactions, and can be used as collateral in voting for proposals.

Trustless Hub is built on the Cosmos SDK with IBC for secure cross-chain communication. Trustless Hub enables new trustless services, empowering users by building non-custodial on-chain and permissionless solutions. It allows every individual, including minorities and people in oppressive regimes, to use blockchain services in a seamless manner.

Use-cases of Trustless Hub include intent-centric triggers for blockchain calls, recurring swaps for DeFi, automatic reward claiming, instant settlement of results, and one-click recurring transactions for payments. It opens up possibilities for innovative and user-centric applications.

Key Features:

1. Intent-centric triggers for blockchain calls
2. Non-custodial and permissionless solutions
3. Time-based blockchain calls and recurring execution
4. Automation over IBC for trustless account automation on other blockchains
5. Integration with AuthZ for secure and permissioned automation

## Use Cases

* Scheduling of Payment (This includes Streaming of transfers, Payroll, Subscriptions, Payment in Installments)
* Auto Compounding of assets
* Automating on other blockchains over Interchain Accounts and AuthZ
  * CosmWasm smart contract calls on CosmWasm enabled blockchains
  * EVM smart contract calls on IBC-connected EVM chains such as Evmos 
  * EVM smart contract calls over Axelar using a MsgTransfer with EVM parsed payload specified in the memo field
  * DeFi swaps on e.g. Osmosis, in the intent, the swapped amount can be combined with a MsgSend to send the funds to a (new) recipient on any chain.
  * ...
* Scheduling of Governance Proposals


Trustless Hub aims to be the most innovative blockchain platform across all ecosystems. By enabling intent-centric automation and building trustless solutions, Trustless Hub empowers users to create new decentralized applications and services that were previously not possible. It revolutionizes the way transactions are executed, bringing greater control and accessibility to users.

## What is TRST?

TRST is the native token of Trustless Hub and is used to pay for gas fees. It is used for:
Transactions

* Storing, instantiating and executing Trustless Contracts.
* Intent-centric automation fees
* Securing the hub

You can stake your TRST by delegating your TRST to validators. Trough governance, you can vote on important proposals such as upgrades and parameter changes. 

## Why Trustless Services?

Trustless services are crucial for building applications that prioritize user privacy, security, and autonomy. Traditional smart contract platforms lack the ability to offer private inputs, deterministic outcomes, and automatic time-based actions. This limits developers' capabilities in creating truly trustless and user-centric applications.

Trustless Hub addresses these limitations by providing intent-centric automation and building non-custodial, on-chain, and permissionless solutions. It empowers developers to program trustless services with automatic time-based functions, ensuring privacy, verifiability, and efficiency.

By integrating with Trustless Hub over IBC, other chains and applications based on the Cosmos SDK can extend their functionality and benefit from automation and trustless account automation. This opens up new possibilities for recurring transactions, remittances, payroll systems, and more, while preserving privacy and ensuring trustlessness.

Trustless Hub also offers integration with AuthZ, enabling secure and permissioned automation. Developers can leverage AuthZ to define access control policies for automated actions, ensuring that only authorized parties can execute specific functions.

## Configure

The node can be configured with `config.yml`. To learn more see the [reference](https://github.com/tendermint/starport#documentation).

## Learn more

[Trustless Hub documentation](https://docs.trustlesshub.com)

Other useful links

[Trustless Hub website](https://trustlesshub.com/)
[TRST Labs website](https://info.trstlabs.xyz/)
[TriggerPortal - one stop automation tool for the interchain](https://info.trstlabs.xyz/)
[Trustless Hub wallet interface](https://interact.trustlesshub.com/)
[Cosmos SDK documentation](https://docs.cosmos.network)
[Cosmos SDK Tutorials](https://tutorials.cosmos.network)
