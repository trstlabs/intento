---
title: Highlight - Autocompound
order: 2
description: How Wallets Can Enhance UX By Offer Smart Auto-Compounding
---

Imagine a user, Emma, who stakes her ATOM tokens using her favorite wallet, but she constantly has to manually claim and restake her rewards. She wants her staking rewards to grow, but keeping track of them and optimizing when to reinvest is tedious.

What if Emma’s wallet could handle this for her intelligently? Instead of just automating restaking at fixed intervals, her wallet could integrate **intent-based auto-compounding**, ensuring that her rewards are reinvested under optimal conditions. For example:

- **Skip compounding if claimed tokens are below a certain threshold**
- **Withdraw rewards instead of compounding if the token price hits a take-profit target**
- **Auto-compound only if inflation is below 10%**
- **Dynamically adjust reinvestment intervals based on reward size**

By offering a smart auto-compounding solution, wallets can provide users with better staking experiences, increase token lock-up rates, and add value to their platforms. Here are different ways this can be implemented:

### Approach 1 - Conditional Autocompounding via User-Owned Self-Trustless agent

![1](@site/docs/images/designing-flows/auto-compound/1.png)

This approach gives users full control over their auto-compounding while keeping everything within their own wallet.

**Advantages:**

- User-owned flow with complete control.
- Users can choose any validator.
- Customizable duration and interval.
- No reliance on third-party services.
- Users can set personal conditions for when to compound, ensuring optimal rewards reinvestment.

**Disadvantages:**

- Users must manage the interchain account themselves.
- Users pay host chain fees directly.
- Requires technical knowledge to maintain.

### Approach 2 - Conditional Autocompounding via Trustless agent

![2](@site/docs/images/designing-flows/auto-compound/2.png)

A Trustless Agent with fees managed by the wallet provider, making auto-compounding easier while ensuring a seamless staking experience.

**Advantages:**

- Users can choose any validator.
- Simplifies auto-compounding by removing the need for self-management.
- Fees can be paid in a single token (INTO, ATOM).
- Wallet providers can generate revenue from managing Trustless Agents.
- Predefined smart conditions ensure efficient auto-compounding.

**Disadvantages:**

- Users must trust the wallet provider for execution.
- Some flexibility is lost compared to self-hosted interchain accounts.
- Still requires fees, but they are abstracted.

### Approach 3 - Autocompounding via One Managed Flow

![3](@site/docs/images/designing-flows/auto-compound/3.png)

For comparison, a fully managed approach where the wallet provider handles bulk execution of compounding actions for all users. This is similar to how auto-compounding is usually done, except for it being on-chain.

**Advantages:**

- No need for custom smart contracts or maintenance of scripts.
- Cost-efficient: executes actions in bulk, optimizing gas fees.
- Fees can be easily abstracted from end users.
- Ensures high staking participation with minimal user effort.

**Disadvantages:**

- One flow means the same conditions for all users.
- Difficult to monetize effectively.
- Users depend on the wallet provider for execution integrity.

By integrating smart auto-compounding, wallets can help users like Emma maximize their staking returns while removing the hassle of manual compounding. This leads to higher staking participation, increased token security, and a better user experience—all while strengthening the wallet’s value proposition.
