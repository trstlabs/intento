---
order: 3
title: IBC Relayers
description: Steps to install a Trustless Hub Daemon on your local computer or dedicated server.
---

## IBC Relayers Documentation

Inter-Blockchain Communication (IBC) Relayers are crucial components of the Cosmos ecosystem. They enable the secure and efficient transfer of tokens and data between IBC-enabled chains.

To run an IBC Relayer, it is recommended to have a machine with Linux as the operating system. This can be either a dedicated server or a virtual private server (VPS). It is also important to have access to RPC endpoints or to run a light client for each chain being relayed.

One of the key benefits of IBC Relayers is that they are trustless. This means that they do not require trust in any third-party intermediaries or middlemen. Instead, IBC Relayers rely on cryptographic proofs and algorithms to ensure the integrity and security of data and transactions.

IBC is a fundamental component of many Cosmos chains, including Cosmos Hub, Terra, and Akash. It allows for cross-chain transactions and the interoperability of various blockchain applications.

Interchain Accounts are a critical component of IBC Relaying. They allow users to interact with other chains in a trustless manner, and to automate transactions using Trustless Triggers. Interchain Accounts are created by generating a new account on the host chain, which is controlled by a trusted relayer. This account can then be used to interact with other chains through IBC channels.


### Comparison of Cosmos/Relayer and Hermes Relayer

[Cosmos/Relayer](https://github.com/cosmos/relayer) (rly) and [Hermes Relayer](https://github.com/informalsystems/hermes/issues) are two popular options for running an IBC Relayer in the Cosmos ecosystem. While both offer similar functionality, there are some differences between the two. Cosmos Relayer is easy to use and preconfigurations are available. Hermes Relayer is more advanced. In `$HOME/trst/Makefile` you can find commands to run local relayers from configurations in  `$HOME/trst/deployment/ibc`.

### Cosmos/Relayer

Cosmos/Relayer is an open-source tool for running an IBC Relayer. It is maintained by the Cosmos development team and is designed to work with any IBC-enabled chain in the Cosmos ecosystem. 

### Hermes Relayer
Hermes Relayer is another open-source tool for running an IBC Relayer. It is developed and maintained by  Informal Systems, and is designed for performance. 

One of the key features of Hermes Relayer is that it is well documented. It is designed to be easy to set up and use, with a user-friendly interface and straightforward configuration options.


## Installing Cosmos/Relayer (rly) on Linux

Cosmos/Relayer (rly) is a tool for building and running IBC relayers in Cosmos-based blockchain networks. In this guide, we will walk you through the step-by-step process of installing rly on Linux.

### Prerequisites

Before installing rly, you need to ensure that your system meets the following prerequisites:

- Linux operating system
- Golang installed on your system
- Git installed on your system 

Step 1: Install Golang
First, you need to install Golang on your system. You can follow the official installation instructions for your specific operating system. After installation, you can check if Go is installed by running the following command:

```go
$ go version
```

If Go is installed on your system, you will see the version number in the output.

Step 2: Download and install rly
Once Go is installed, you can download and install rly using the following command:

```go
$ go install github.com/cosmos/relayer@latest
```

This command will download the latest version of rly from the official Github repository and install it on your system. Once the installation is complete, you can check if rly is installed by running the following command:

```go
$ rly version
If rly is installed on your system, you will see the version number in the output.
```

