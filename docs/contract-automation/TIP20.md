---
order: 6
title: TIP20
description: Learn how TIP20 - Trustless Hub's token standard works
---

# TIP20

## Token standard with maximum programmability
*CosmWasm token standard which can instantiate with allowance, has privacy, can be backed 1-1 by a chain-native asset*

This is a privacy-first token implementation on Trustless Hub, TIP20. It can be backed by
an IBC-native coin like TRST, ATOM, OSMO, JUNO ect. and has a fixed 1-to-1 exchange ratio
with it.

Using it is simple - you deposit the token of the contract (e.g. TRST) into the contract, and you get pToken
(e.g. pTRST), which you can then use with the ERC-20-like functionality that
the contract provides including: sending/receiving/allowance and withdrawing
back to TRST.

In terms of privacy the deposit & withdrawals are public, as they are
transactions on-chain. The rest of the functionality can be private (so no one can
see if you send pTRST and to whom, and receiving pTRST can also be hidden).

You can execute and instantiate other contracts. To the chain, only your address has interacted with this privacy-perserving token.

## Usage examples

check out the schema provided for the JSON messages to send.

## Play with it on testnet

The deployed pTRST contract address on the testnet is
`trst...` and label `ptrst`

## Troubleshooting

All transactions are encrypted, so if you want to see the error returned by a
failed transaction, you need to use the command

```
trstd q compute tx <TX_HASH>
```

## Notes

* Addresses of private-TRST accounts are the same as of their respective trust
  account.
* The exchange ratio is fixed at 1-to-1.
* The total supply will always equal the amount of
  TRST locked in the contract, which can be seen in the explorer and in the contract public state.
* Messages can be padded with a multiple of 256 bytes by convention to maximize
  privacy.
