---
order: 2
title: CosmWasm code
description: Useful information regarding CosmWasm code
---

## Getting familiar with CosmWasm
On [Youtube](https://www.youtube.com/results?sp=mAEB&search_query=CosmWasm). there are many great tutorials on CosmWasm. 

[CosmWasm](https://docs.cosmwasm.com/docs/1.0/) provides extensive documentation


## CosmWasm on TRST
CosmWasm is a smart contract standard within the Cosmos Ecosystem. Many chains use these, including Juno, Terra, Stargaze, Irisnet, Omniflix and Secret Network. 

Code can be used to:
Send funds
Vote on proposals
Execute other contracts
Instantiate other contracts
Query other contract states


CosmWasm code fo Trustless Contracts should be designed to be:
1) Stored
2) Instantiated 
3) Executed
4) Queried

And additionally:

5) Automatically executed as predefined at instantiation 
6) Deleted as predefined at instantiation 

1,2,3,4 are similar to develop as in other smart contract platforms. Encryption and privacy are enabled by the blockchain, and the developer can carelessly use the benefits of these.
For 5, automatically executing code, an automated message should be enabled, and this can point to an existing function or to a completely seperate one than that can be executed on by the users of the contract.
You can name this AutoMessage, or anything else. Please refer to the message to be handled as AutoMessage, so people investigating the code are aware of the which part of the code to be run automatically.


![Example auto_msg on internal estimation contract](./auto_msg_example.png)
Above, an example on the AutoMessage pointing to a function on the internal estimation Trustless Contract

## Differences in contract transactions with standard CosmWasm

### Executing 
Executing code is done by encrypting the message with the code hash and sender public key, so that the message can only be executed by the code it belongs to. 

The inputs are private and are only decrypted once the message is in the Trusted Execution Environment (TEE), where the inputs are then securely handled. The TEE, that runs through Intel SGX, and is designed in such a way that no other process or application is able to view or currupt the contents.

### Instantiating 
Next to the standard InitMsg, an AutoMessage can to be sent to automatically execute code. 
The instantiation message and the automated message are encrypted and only decrypted once the message is in in the Trusted Execution Environment 

### Querying 
The contract can be queried through {RCP API URL}/compute/v1beta1/contract/{contract_address}/smart/{query_data}. 

The query inputs are encrypted, like with executing code
the result of the query is always encrypted and only viewable by the person that performs the query by decrypting the result.

In addition, like with other CosmWasm contract instances, contracts can also query other contracts.

### Contract Result
Cntract Results can be queried by anyone. It is basically a publicly viewable state of the private smart contract. The 'handleResponse' of an execution message gets saved on-chain. This is also the case for the AutoMessage. The Contract Result (last available result) is can be queried through {RCP API URL}/compute/v1beta1/contract/{contract_address}/result . 

As a developer you should keep this in mind in what you send back as information for each transaction.


