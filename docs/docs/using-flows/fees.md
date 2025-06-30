---
sidebar_position: 7
title: Flow Fees
description: How Flow fees are calculated and examples
---

## Fees & Execution Costs

Flow fees consist of two components:

### 1. Gas Fee (always charged)

A dynamic fee calculated as:

```
gasUsed * flowFlexFeeMul
```

* The **denom** used to pay this fee is selected in the `gasFeeCoins` field.
* The actual cost depends on gas usage and is scaled based on the selected token.

### 2. Burn Fee (only when paying in $INTO)

If the selected `gasFeeCoins` includes **INTO**, a **fixed burn fee** per message is added:

```
burnFee = burnFeePerMsg * messageCount
```

* This burn is **only** applied when using `INTO` as the gas token.
* It is not charged if you're paying fees in ATOM, OSMO, or other tokens.

> TL;DR:
>
> * Pay in INTO → gas fee + burn fee
> * Pay in ATOM/OSMO → only gas fee

---

## Hosted Account Fee (extra)

When using a **hosted interchain account**, there's an additional fee charged **per execution**. This is set by the fee admin of the hosted account.

In Intento Portal, we **automatically add a fee coin limit** during flow submission to cover this hosted fee — including a buffer to avoid underfunded execution errors.

To check the latest hosted fee programmatically:

```
GET /intento/intent/v1beta1/hosted-account/{address}
```

**LCD URL:**
[https://lcd.intento.zone/swagger/#get-/intento/intent/v1beta1/hosted-account/-address-](https://lcd.intento.zone/swagger/#get-/intento/intent/v1beta1/hosted-account/-address-)

This endpoint returns the `fee_coins_supported` array of supported fee coins for that specific hosted account.

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
