---
order: 2
title: CosmWasm code
description: Useful information regarding CosmWasm code
---

## Getting familiar with CosmWasm
On [Youtube](https://www.youtube.com/results?sp=mAEB&search_query=CosmWasm). there are many great tutorials on CosmWasm. 

[CosmWasm](https://docs.cosmwasm.com/docs/1.0/) provides extensive documentation


## CosmWasm on TRST
CosmWasm is a smart contract standard within the Cosmos Ecosystem. Many chains use these, including Juno, Osmosis, Terra, Secret Network, Stargaze and Archway.  

CosmWasm code should be designed to be:
1) Stored
2) Instantiated 
3) Executed
4) Queried

And additionally for Trustless Hub they have:

5) An AutoMsg to execute 1-time or recurringly. This can defined in the CodeInfo or at Contract Instantiation

1,2,3,4 are similar to develop as in other smart contract platforms. Encryption and privacy are enabled by the blockchain, and the developer can carelessly use the benefits of these.
For 5, automatically executing code, an automated message should be enabled, and this can point to an existing function or to a completely seperate one than that can be executed on by the users of the contract.
You can name this AutoMsg, or anything else. Refer to the message to be executed as AutoMsg, so people viewing the contract code are aware of the which part of the contract can run automatically.


![Example auto_msg](./auto_msg_example.png)
Above, an example on the AutoMessage pointing to a function on a recurring swap Trustless Contract. 

#### How does this work?
After instantiating from the TIP20 token contract, the TIP20 token contract gives allowance to this contract for the max funds to swap. Trustless Hub then recurringly calls AutoMsg. 

## Differences in contract transactions with standard CosmWasm

### Executing 
Executing code is done by encrypting the message with the code hash and sender public key, so that the message can only be executed by the code it belongs to. 

The inputs are private and are only decrypted once the message is in the Trusted Execution Environment (TEE), where the inputs are then securely executed. The TEE, that runs through Intel SGX, and is designed in such a way that no other process or application is able to view or currupt the contents.

### Instantiating 
Next to the standard Msg, an AutoMessage can to be sent to automatically execute code. 
The instantiation message and the automated message are encrypted and only decrypted once the message is in in the Trusted Execution Environment 

### Querying 
The contract can be queried through {RCP API URL}/compute/v1beta1/contract/{contract_address}/smart/{query_data}. 

The query inputs eare encrypted, like with executing code
the result of the query is always encrypted and only viewable by the person that performs the query by decrypting the result.

In addition, like with other CosmWasm contract instances, contracts can also query other contracts.

### Public State
Contracts can have a public state that is cheap to query. It is also easy integrate and show users a way to view what is happening on the private by default contract. The Contract State is can be queried through {RCP API URL}/compute/v1beta1/contract/{contract_address}/public-state . 
As a developer you save outputs to the public state when you make a new Response.


