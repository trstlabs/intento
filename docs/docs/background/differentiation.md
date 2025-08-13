---
title: Differentiation
sidebar_position: 3
---

Intento enables **intent-based flows without smart contracts**, making automated execution streamlined, self-custodial, and scalable. Unlike traditional automation that depends on bot networks or privileged smart contracts, Intento empowers users to directly own and control execution, removing unnecessary intermediaries and security risks.

### The Problem with Existing Approaches  

#### Bot Networks  

Decentralized execution via bot networks requires tasks to be registered on-chain, with bots competing to claim rewards. However, this approach has inefficiencies:

- **Network Congestion & Gas Costs** – Multiple bots executing the same transaction cause congestion and increase fees.
- **Centralization Risks** – The first-come-first-served model discourages participation from smaller operators, leading to a race-to-the-bottom effect.
- **Complex Automation Workflows** – Execution logic requires additional privileged contracts, increasing costs and security risks.

#### Custom Smart Contracts  

Privileged Smart contracts on CosmWasm chains can schedule executions via BeginBlocker/EndBlocker, but this introduces drawbacks:

- **Higher Gas Costs** – Scheduled executions require extra computational resources, increasing user fees.
- **Computational Overhead** – Smart contract-based automation incurs static costs of running inside a VM, regardless of efficiency.
- **Blockspace Limitations** – Time-based executions compete with general transactions, requiring governance oversight.
- **Security Risks** – Privileged contracts introduce attack vectors, increasing risks for users.

### How Intento Solves This  

Intento eliminates the need for bots and privileged smart contracts, enabling direct, intent-based execution:

- **User-Owned Flows** – Intento shifts power to users, enabling self-custodial, automated workflows without intermediaries.
- **Scalable & Efficient** – Execution logic runs at a fraction of the cost of traditional automation services.
- **Protocol-Neutral & Interoperable** – Works across different blockchain ecosystems without dependencies on specific contract frameworks.
- **Orchestration & Composability** – Intento enables composing multiple flows together, allowing users to create complex, automated processes across chains.
- **Direct IBC, ICA & ICQ Integration** – By integrating Inter-Blockchain Communication (IBC), Interchain Accounts (ICA), and Interchain Queries (ICQ), Intento facilitates secure, automated, and cross-chain workflows without requiring intermediary execution layers.
- **Conditions, Comparisons & Feedback Loops** – Execution logic can be condition-based, supporting real-time comparisons, iterative feedback loops, and adaptive automation.
- **Advanced Use Cases** – Supports conditional payments, streaming transactions, auto-compounding, and portfolio optimization, making it a powerful tool for DeFi and beyond.
- **Single Token Fee Abstraction** – With Trustless Agents, users can pay all fees in a single token (e.g., INTO or ATOM), simplifying automation and reducing friction.

### How Intento Differs from Anoma & Agoric  

#### Anoma (Intent-Matching vs. Intent-Execution)

Commonalities:

- **User-Centric & Flexible** – Like Anoma, Intento allows generalized intent-based automation without smart contracts.
- **Plug-and-Play Automation** – No need for additional specialized infrastructure.
- **Expands What’s Possible** – Enables applications that cannot be built purely on smart contract VMs.

Differences:

- **Intento is focused on execution, not matching** – Anoma facilitates intent discovery and matching, while Intento directly executes user-owned flows.
- **More deterministic execution** – Intento guarantees predictable execution without relying on intent-matching networks.

Examples of Anoma applications:

- Fully decentralized order book exchanges
- Decentralized Slack/Discord alternatives
- Matchmaking apps (e.g., decentralized Tinder for X)

#### Agoric (Complex vs. Simple, User-Owned Execution)  

Commonalities:

- **Orchestration & Time-Based Actions** – Both Intento and Agoric allow developers to automate and schedule cross-chain actions.
- **Cross-Chain Interoperability** – Both support multi-chain automation.

Differences:

- **Intento is simple, user-owned, and requires no smart contracts** – Agoric’s Orchestration API is powerful but requires developers to manage remote accounts, multi-block async execution, and timer-based scheduling.
- **Agoric is better suited for complex custom workflows** that require fine-tuning and for protocols that need to build a complex custom solution or want to manage and monetize orchestration infrastructure.
- **Intento enables monetization through Trustless Agents** – Protocols and integrators can offer hosted execution, which abstracts away host chain fees while keeping user flows self-custodial.

Agoric Orchestration API features:

- **Remote account control** – Create and manage accounts on remote chains.
- **Async execution over multiple blocks** – Contracts handle responses across long periods.
- **On-chain timers** – Enables scheduled execution (e.g., subscriptions).

### **Network Effects**  

- **More Adoption → Greater Demand for Execution → Increased Value Capture** – As more users adopt Intento, demand for intent execution grows, driving network utility and token value.
- **Governance-Driven Execution Parameters** – The community governs execution scalability and fee burn parameters, ensuring long-term sustainability.
- **Monetization for Integrators** – Integrators can host interchain accounts, abstracting host chain fees and allowing users to pay fees in a single token (e.g., INTO or ATOM), making automation seamless.
- **Decentralized Revenue Model** – Execution services can be monetized while maintaining full user control and self-custodial execution.
- **Permissionless Hosting** – Any integrator can run Trustless Agent, monetizing by charging fees while ensuring non-custodial execution for users.

### A New Standard for Intent-Based Execution  

Intento redefines automation by shifting from contract-based execution to direct, user-owned intent-based flows. This approach enhances security, reduces complexity, and ensures scalability without compromising decentralization

Whether for DeFi actions, cross-chain operations, or scheduled transactions, Intento delivers a faster, cheaper, and more flexible alternative. With orchestration capabilities, IBC/ICA/ICQ integration, advanced conditional logic, and monetization mechanisms, Intento unlocks new levels of automation, security, and efficiency across blockchains.
