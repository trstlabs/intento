---
sidebar_position: 3
title: Interacting with Intento Portal
pagination_label: How to set and manage actions using Intento Portal
---

Intento Portal is a powerful dApp designed for developers, and integrators, offering an advanced solution for setting, viewing and updating flows with robust conditions, comparisons, and feedback loops. It allows users to define sophisticated workflows for various blockchain actions.

### Key Features of Intento Portal

- **Conditional Automated Token Transfers**

  Effortlessly set up token transfers with conditional logic to multiple recipients at regular intervals. Ideal for organisations or individuals managing recurring payments, Intento Portal ensures precise control over the distribution process based on predefined conditions.

- **Sophisticated On-Chain Time-Based Flows**
  Intento Portal allows users to set up flows incorporating feedback loops and comparisons to ensure actions are performed under specific conditions. This capability supports intricate workflows that react dynamically to on-chain events.

### Benefits for Developers and Integrators

- **Efficiency and Precision**
  - Save time and ensure accuracy by automating repetitive tasks with condition-based actions and feedback mechanisms.
- **Cross-Chain Interoperability**
  - Seamlessly integrate and build processes across multiple blockchains, leveraging IBC capabilities.
- **Enhanced Investment Strategies**
  - Automatically reinvest rewards based on intelligent conditions, optimizing growth without manual effort.
- **Dynamic Workflow Management**
  - Create sophisticated, buildd workflows that respond to real-time on-chain data and conditions, enhancing operational efficiency.

Intento Portal is your comprehensive tool for automating and managing blockchain interactions with advanced flows. Whether you are sending recurring payments, scheduling cross-chain messages, or optimizing token compounding, Intento Portal equips you with the tools to execute with precision and reliability.

## Build Flows

### Getting started

1. Head to the `Flow Builder` page in the menu to start building flows.
   First, specify on what chain to execute on.
   ![1](@site/docs/images/triggerportal/build/1.png)
   Hosted accounts that are available for you to use are displayed here, with relevant information such as the address that will execute on the host chain and the fee coins supported and the current fees. Hosted accounts by Intento Portal will be subsidized.
   When you have an interchain account address registered, it will pop up here too. The interchain account should be funded on host chain. You can set this in the flows dialog.
   You can register an interchain account by clicking `Set ICA`.
2. Second, build messages.
   There are several message examples available. For a CosmWasm-supported chain such as `Neutron` there are examples available to interact with smart contracts. For `Osmosis` there are also next to the smart contract examples examples for the DEX.
   ![2](@site/docs/images/triggerportal/build/2.png)
3. Thirdly, configure execution details. You can execute until you run into an error for example. Or, maybe you want to keep trying until it succeeds. You can use a wallet fallback so fee funds can be deducted from your main balance automatically.
   ![3](@site/docs/images/triggerportal/build/3.png)
4. Next, you can set your conditions. You can use response outputs as inputs to create feedback loops, or create comparisons that must turn true for execution to take place.
   ![4](@site/docs/images/triggerportal/build/4.png)
   You can find inspiration for this in the dashboard!
   ![dialogfullscreen](@site/docs/images/triggerportal/build/Copy.png)
5. Build Dialog

When selecting `Build Flow`, a dialog will pop up. In this dialog you can specify the duration, interval and a start time.

By clicking on the selected interval or start time, you can unselect it. By unselecting start time, the first interval, in this case `1 hour`, will be used as the first time the action gets processed.

![dialogfullscreen](@site/docs/images/triggerportal/build/dialogfullscreen.png)

You can specify whether to deduct from your account or create a fee account and attatch funds to the flow. If you attatch funds to the flow without a wallet fallback, be sure that it has suffient INTO balance to pay for fees at the moment of execution. Fees are returned after the final execution.

In the `Overview` section, you can view how many times execution will take place, when it starts and ends. You can specify a `label` to name your flow, this is optional. You can retreive your flows in the dashboard.

You can now click on `Submit` . An alert will pop up which you can use to navigate to your flow.

### AuthZ permissions

Cosmos SDK messages contain `values` and a `typeUrl`. Message `values` are what you send, the `typeUrl` specifies what module and what version to send it to, along with the proper function.

AuthZ grants are permissions you can grant to an external address to execute messages on your behalf. This can possibly be dangerous when given to a third party. However, with Intento the permission is given to an address that can only execute flows where you were the signer from. There are multiple checks in place for this. This eliminates the risks that arise when granting another account approval. You can grant an ICA with `MsgGrant` with the type of message that is allowed and an `expiration`, and allow the ICA to execute using `MsgExec`. The default expiration on Intento Portal is 1 year.

:::info Double check what you are doing and it is recommended to test first if you are doing it for the first time.
:::

### Notes

When setting a message, you should be aware that keys are in camelCase. So, use `fromAddress` instead of `from_address`. You should also be aware that sending tokens is usually denominated with smaller decimals. For most Cosmos chains, native tokens are denominated in 6 decimals, and have a `u` in front. For 5 INTO you can specify `5000000uinto`. For 5 ATOM you specify `5000000uatom`.

<!-- ## Autocompound Staking Rewards

You can stake INTO tokens to secure the network and earn staking rewards. Staking rewards can be compounded to earn additonal tokens.
Autocompound is a feature that automatically restakes earned rewards back to the validator, compounding earnings over time.

![autocompound](@site/docs/images/triggerportal/build/autocompound.png)

There are several terms used in autocompounding staked tokens.

`Nominal APR` refers to the annual percentage rate that doesn't take into account compounding interest. It's the simple staking reward rate over the course of a year.

`RealTime APR` refers to the annual percentage rate that is calculated and updated in real-time base based on the current block time.

`APY` stands for Annual Percentage Yield and represents the effective annual rate of return of staked INTO tokens that is compounded over the course of a year. In the case of Weekly Compound APY, the rewards are calculated and added to the staking balance every week.

Using the actions dialog you can specify the interval of the autocompound. Your strategy should take into account execution fees which are estimated under `Execution Settings`. -->

## Demo (June 2025)

<iframe width="560" height="315" src="https://www.youtube.com/embed/q1D9uLIh9GE" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" referrerpolicy="strict-origin-when-cross-origin" allowfullscreen></iframe>

## Future Improvements

Do you have an interesting feature in mind? Mention it on the Intento Portal [GitHub repository](https://github.com/trstlabs/intento-portal) or [X/Twitter](https://twitter.com/IntentoZone) and it may get added to the roadmap.
