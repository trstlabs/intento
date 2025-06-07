---
sidebar_position: 7
title: Flow Fees
description: How Flow fees are calculated and examples
---

# Flow Fees: Calculation and Examples

## Overview

Flow fees are made up of two components:

1. **Gas Fee**:
   A dynamic fee based on gas usage (`gasUsed`) and the multiplier `flowFlexFeeMul`. The gas fee is scaled based on the selected **denom** in `gasFeeCoins`.

2. **Burn Fee**:
   A fixed fee per message (`burnFeePerMsg`), applied only when using `INTO` as the payment denom.

By default, **INTO** is the standard denom for paying Flow fees. If you select another gas denom (e.g., `ATOM`, `OSMO`), no burn fee applies. Only the gas fee.

---

## Payment Source

When you **create a Flow**, the system attempts to charge your wallet directly. If you create a Flow **via MsgTransfer from another chain** (e.g., IBC from Cosmos Hub), the funds sent along are **deposited to your account** on the Flow chain. The Flow fee is then deducted from your account. This **wallet fallback** ensures you can pay for flows even when initiating from another chain.

We aim to support **ATOM**, **OSMO**, and other integrator tokens in our gas configuration, making it easier to pay for flows.

---

## Worked Examples

### Autocompound Flow

**Scenario:**

- 2 messages per Flow
- Runs once per week for 1 year (52 runs)
- Uses `INTO` as the gas denom

**Fee Parameters:**

- `gasUsed = 63700`
- `flowFlexFeeMul = 2`
- `burnFeePerMsg = 10000`
- `INTO` gas price = `30` per unit

**Per Flow Fee:**

- Gas fee units: `(63700 * 2) / 1000 = 127.4`
- Gas fee in `INTO` (microdenom): `127.4 * 30 = 3822`
- Burn fee in `INTO` (microdenom): `10000 * 2 = 20000`
- Total Flow fee (microdenom): `23822`
- Converted to `INTO`: `0.023822`

**Burned `INTO` per Flow:**
`0.02 INTO` (approx, from the burn fee)

**Total Fees over 52 runs:**

- Total Flow fee: `1.2387 INTO`
- Total Burned: `1.04 INTO` (`0.02 * 52`)

---

### Token Stream Flow

**Scenario:**

- 1 message per Flow
- Runs 10 times over 3 days
- Uses `INTO` as the gas denom

**Per Flow Fee:**

- Gas fee units: `127.4`
- Gas fee in `INTO`: `127.4 * 30 = 3822`
- Burn fee in `INTO`: `10000 * 1 = 10000`
- Total Flow fee (microdenom): `13822`
- Converted to `INTO`: `0.013822`

**Burned `INTO` per Flow:**
`0.01 INTO` (from burn fee)

**Total Fees over 10 runs:**

- Total Flow fee: `0.13822 INTO`
- Total Burned: `0.1 INTO` (`0.01 * 10`)

---

### Token Stream with IBC Denom

**Scenario:**

- 1 message per Flow
- Runs 10 times over 3 days
- Uses `ibc/...` denom (`5` per unit) e.g. ATOM

**Per Flow Fee:**

- Gas fee units: `127.4`
- Gas fee in denom: `127.4 * 5 = 637`
- Burn fee: Not applied for IBC denoms
- Total Flow fee (microdenom): `637`
- Converted to denom: `0.000637`

**Burned `INTO`:**
`0` (no burn fee for non-INTO denoms)

**Total Fees over 10 runs:**

- Total Flow fee: `0.00637` denom
- Total Burned: `0 INTO`

---

## Summary Table

| Flow Type    | Messages | Runs | Denom  | Fee per Flow | Total Fee | Burned `INTO` |
| ------------ | -------- | ---- | ------ | ------------ | --------- | ------------- |
| Autocompound | 2        | 52   | `INTO` | \~0.0238     | \~1.2387  | \~1.04        |
| Token Stream | 1        | 10   | `INTO` | \~0.0138     | \~0.13822 | \~0.1         |
| Token Stream | 1        | 10   | `ATOM` | \~0.000637   | \~0.00637 | 0             |
