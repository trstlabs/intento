---
sidebar_position: 3
title: Interacting with Intento Portal
pagination_label: How to set and manage actions using Intento Portal
---

Intento Portal is a powerful dApp designed for developers and integrators. It enables advanced workflows for setting, viewing, and updating flows using conditions, comparisons, and feedback loops. It supports complex blockchain interactions and cross-chain execution logic.

### Key Features of Intento Portal

- **Conditional Automated Token Transfers**  
  Define recurring token transfers with custom conditions to multiple recipients. Useful for DAOs, payroll, and streaming payments.

- **Advanced On-Chain Time-Based Logic**  
  Build flows with feedback loops, comparisons, and time intervals that react to real-time on-chain data.

- **Cross-Chain Interoperability**  
  Execute workflows across multiple Cosmos chains using IBC and ICA (Interchain Accounts).

### Benefits

- **Efficiency and Precision**  
  Orchestrate repetitive blockchain tasks with logic-based triggers.

- **Composable Investment Strategies**  
  Reinvest, stream, or redirect rewards based on intelligent logic.

- **Dynamic Workflow Management**  
  Build composable on-chain workflows that adapt based on conditions, feedback, and state.


## Build Flows

### 1. Choose Chain and Execution Account

Go to the `Flow Builder`. First, choose the target chain for execution. Available hosted accounts will be listed — these are ICA addresses hosted by Intento which subsidizes fees.

![1](@site/docs/images/portal/build/1.png)


### 2. Build Messages

Add messages that define what you want to execute. Several examples are available for supported chains. For CosmWasm chains like Osmosis, you’ll find both contract interaction templates and DEX-related messages. For DEXes like Elys, we have several supported message types like swapping and claiming rewards.

![2](@site/docs/images/portal/build/2.png)


### 3. Configure Execution Settings

Decide how execution should behave:

- Retry until success
- Stop after an error
- Enable fallback wallet (fees can be taken from your main wallet if the flow account runs out)

![3](@site/docs/images/portal/build/3.png)


### 4. Add Conditions & Feedback Logic - Optional

Use response outputs to feed into future steps or define comparisons. Execution will only happen if all conditions return `true`.

Use these tools to create feedback loops, like:

- Re-checking balances
- Waiting for specific responses
- Responding to partial state changes

![4](@site/docs/images/portal/build/4.png)

You can find flow examples and reusable condition snippets in the dashboard or in the [ integration hub repository](https://github.com/trstlabs/intento-integration-hub)

### 5. Final Flow Settings

When you're ready, hit `Build Flow`. In the dialog:

- **Set start time** _(optional)_: when the first run begins
- **Set interval**: how often it should run
- **Set end time**: when the flow should end. You can use the calendar to select a date and time or the quick selector to select a duration.

> Unselecting the start time means the first interval (e.g., 1 hour) will be used as the initial delay.

You’ll also configure fee settings:

- Deduct from wallet (requires fallback wallet)
- Attach fee account to the flow  
  *(Ensure it holds enough $INTO. Unused fees are refunded after the final execution.)*

![dialogfullscreen](@site/docs/images/portal/build/dialogfullscreen.png)

The `Overview` section summarizes:

- Number of executions
- Start/end time
- Optional label for identifying the flow later

Click `Submit` to broadcast the flow. A notification will link to the dashboard view.


## AuthZ Permissions

Flows using hosted ICAs require `AuthZ` (authorization) from your wallet.

In Cosmos, every message has:
- a `typeUrl` (e.g., `/cosmos.bank.v1beta1.MsgSend`)
- and `value` fields (payload)

You grant `AuthZ` to an ICA using `MsgGrant` for specific message types and expiration (defaults to your flow's end time + 1 day on Intento Portal). The ICA can then use `MsgExec` to execute those messages.

✅ **Security tip**: AuthZ in Intento is scoped. ICAs can **only** execute flows you signed — reducing risk even if the grant is compromised.

:::info
It’s recommended to test with a small flow before going live, especially on mainnet.
:::


## Notes on Message Format

- Use **camelCase** in message fields (e.g., `fromAddress` not `from_address`)
- Token values are in **base units**:
  - `5 ATOM` = `5000000uatom`
  - `5 INTO` = `5000000uinto`


## Demo (June 2025)

<iframe width="560" height="315" src="https://www.youtube.com/embed/q1D9uLIh9GE" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" referrerpolicy="strict-origin-when-cross-origin" allowfullscreen></iframe>


## Future Improvements

Got ideas or feedback? Drop them on the Intento Portal  
[GitHub Repo](https://github.com/trstlabs/intento-portal) or [Twitter/X](https://twitter.com/IntentoZone) and help shape the roadmap.
