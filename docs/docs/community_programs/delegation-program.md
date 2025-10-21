---
title: Delegation Program
description: Learn how to participate in the Intento Delegation Program and earn rewards
---

# Delegation Program

### Overview
The Delegation Program allows governors to earn additional staking rewards.

> This program acts as a framework for allocating INTO to active, engaged, and impactful governance participants.


### 0) TL;DR

* **Tracks:** Governance Delegates, Ecosystem Builders, Public Goods, Research & Security.
* **Pool size:** 145M INTO (of 400M total supply).
* **Epoch:** Quarterly. Rebalance per epoch using engagement, proposal voting, usage, and integration metrics.
* **Pre-launch:** Initial delegation to all governors equally to bootstrap activity and excitement.
* **Transparency:** Reports on each rebalance epoch.


### 1) Goals & Principles

1. **Maximize governance participation:** Reward governors who actively vote, propose, and engage with the community.
2. **Encourage ecosystem integrations:** Incentivize integration with Intento Flows by supplying tools, RPC nodes, relayers, and Trustless Agents.
3. **Public goods & tooling:** Fund infrastructure that supports the network (explorers, relayers, RPC, agents).
4. **Transparency:** Deterministic scoring, clear metrics, and open dashboards.
5. **Initial excitement:** Equal pre-launch delegations to governors to signal participation and engagement.


### 2) Pools & Default Split

Total delegation budget: 145M INTO.

* **Governance Delegates — 50%**
* **Ecosystem Builders Pool — 25%**
* **Public Goods Pool — 15%**
* **Research & Security Pool — 10%**


### 3) Timeline

* **Applications Open:** 2025-10-21
* **Applications Close (Wave 0):** 2025-12-01
* **Initial Allocations Posted:** Start of Q1 2026
* **Epoch length:** Quarterly; first rebalance in the first week of Q1 2026
* **Pre-launch Delegation:** Equal split across all governors prior to Wave 0.


### 4) Track Details & Scoring

#### 4.1 Governance Delegates

**Eligibility**

* Public identity, active in proposal voting and community discussions.
* Commitment to engage with Intento-related governance and ecosystem decisions.

**Scoring (100 pts)**

* **Voting Participation (40):** % of proposals voted per epoch.
* **Proposal Engagement (30):** Number, quality, and impact of proposals submitted.
* **Community Involvement (15):** AMAs, forum engagement, office hours.

**Mechanics**

* Base allocation per delegate `A_base`, plus score-weighted variable `A_var = T * (score / Σscore)`.
* Decay for inactivity: −20% if participation drops below 50% of previous epoch.
* Rotation encouraged: max 2 consecutive terms without review.


#### 4.2 Ecosystem Builders Pool

**Eligibility**

* Teams building integrations or tooling (flows, autocompounding templates, alerting, wallets, analytics).

**Scoring (100 pts)**

* **Execution & Usage (50):** Features deployed, user metrics (TVL, volume, flows).
* **Integration Depth (25):** Ecosystem coverage (chains, ICAs, ICQ, ICS).
* **Open Source & Tools (15):** Repos, docs, example flows.
* **Alignment (10):** Use of intent orchestration patterns.


#### 4.3 Public Goods Pool

**Scope**

* RPC nodes, relayers, explorers, chain-indexers, educational content.

**Scoring (100 pts)**

* **Reliability (40):** Uptime.
* **Usage & Adoption (30):** Metrics of service usage.
* **Neutrality (10):** Open access to services.


#### 4.4 Research & Security Pool

**Scope**

* Audits, incident response, formal verification, economic security studies.

**Scoring (100 pts)**

* **Impact (50):** Published advisories, mitigations.
* **Execution (30):** Completed studies, contributions to ecosystem tooling.
* **Community Sharing (20):** Documentation, workshops, reports.


### 5) Calculation Details

* We will publish **scoring** per epoch in the mainnet validator Discord channel.

### 7) Risk & Incident Playbook

* Inactivity: apply decay, redelegate to active governors.
* Governance capture attempt: investigation; temporary cap reduction.
* MEV or tooling abuse: warning → removal if unremedied.

### 8) Reporting & Transparency

* Dashboard: live stake distribution, score breakdown, epoch countdown.
* Monthly report: allocations, changes, rationale, incidents.
* Archive: immutable snapshots (IPFS/Ceramic) per epoch.


### 9) Application Forms

Interested in joining the program? Fill out the official application form:  
 [Intento Delegation Program Application](https://docs.google.com/forms/d/e/1FAIpQLSd9Vnr58jz6h36Q_6XyL-FaB9aGsnAnpMU4F3OU762ohx5WJw/viewform?usp=dialog)


### 10) Defaults

* Epoch: quarterly
* Decay for inactivity: −20%
* Base allocations per pool configurable