---
order: 4
title: Contract State
description: Learn how to integrate a contract state and query it
---

# Contract State

Contracts are private-by-default. State is encrypted. There is a public state avaliable too. 

Let's see how the public state works.
## CosmWasm Response object
At the end of a CosmWasm function, developers can define what to be added to the contract's public state in the Response object.


### Example: Instantiating TIP20 contract

- log = encrypted and viewable for sender only
- pub_db = contract state
- acc_pub_db = account_specific contract state


```rust

    let log = vec![
        log("status", "success"),
        pub_db("name", msg.name),
        pub_db("symbol", msg.symbol),
        pub_db("decimals", msg.decimals.to_string()),
        pub_db("admin", admin.as_str()),
        pub_db("total_supply", supply),
        pub_db("minter", info.sender.clone()),
        pub_db(
            "total_supply_is_public",
            init_config.public_total_supply().to_string(),
        ),//boolean indicating if total supply is viewable without viewing key
        acc_pub_db("initiated", "contract", info.sender),
    ];

    Ok(Response::default().add_attributes(log.clone()))
        
``` 

### Example: Updating admin in TIP20 contract
Here we directly add the pub_db items to the CosmWasm Response object.
![pubdb](../images/pubdb.png)


- A public state is queryable through RPC like so:
 {RCP API URL}/compute/v1beta1/contract/{contract_address}/public-state
-  An account-specific public state is queryable like so:
  {RCP API URL}/compute/v1beta1/contract/{contract_address}/public-state/{account_address}