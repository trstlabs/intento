---
order: 2
title: Running a Validator
description: Steps to install a Trustless Hub Daemon on your local computer or dedicated server.
---

# Running a Validator

Before setting up a validator node, make sure to have completed the [Client](./trstd.md) guide.

If you plan to use a Key Management System (KMS), please complete the following steps:  [Using a KMS](https://hub.cosmos.network/main/validators/kms/kms.html).


## What is a Validator?

Validators play a crucial role in the Cosmos ecosystem by committing new blocks to the blockchain through a voting process. Validators risk having their stake slashed if they are unavailable or sign blocks at the same height. To protect your node from denial-of-service (DDoS) attacks and ensure high availability, we recommend familiarizing yourself with the Sentry Node Architecture.


::: danger Warning
If you want to become a validator for `mainnet`, you should learn more about [security](./security.md).
:::

The following instructions assume that you have already set up a full-node and are synchronized to the latest block height.

## Create Your Validator

To create a new validator, you need to stake tokens using your `trustvalconspub`. Find your validator's public key by running the following command:

```bash
trstd tendermint show-validator
```

To create your validator, just use the following command:

::: warning
Don't use more `utrst` than you posess!
:::

```bash
trstd tx staking create-validator \
  --amount=1000000utrst \
  --pubkey=$(trstd tendermint show-validator) \
  --moniker="choose a moniker" \
  --chain-id=<chain_id> \
  --commission-rate="0.10" \
  --commission-max-rate="0.20" \
  --commission-max-change-rate="0.01" \
  --min-self-delegation="1000000" \
  --gas="auto" \
  --gas-prices="0.0025utrst" \
  --from=<key_name>
```

::: tip
When specifying commission parameters, the `commission-max-change-rate` measures the percentage point change over the `commission-rate`. For example, a change from 1% to 2% represents a 100% rate increase, but only 1 percentage point.
:::

::: tip
`Min-self-delegation` is a stritly positive integer that represents the minimum amount of self-delegated voting power your validator must always have. A `min-self-delegation` of `1000000` means your validator will never have a self-delegation lower than `1 TRST`
:::

Initially, you may not have enough TRST to be part of the active set of validators. Users can delegate to inactive validators (those outside of the active set) using the Keplr web app. You can confirm your validator's inclusion in the validator set by using a third-party explorer like Mintscan.

## Edit Validator Description

You can edit your validator's public description. This info is to identify your validator, and will be relied on by delegators to decide which validators to stake to. Make sure to provide input for every flag below. If a flag is not included in the command the field will default to empty (`--moniker` defaults to the machine name) if the field has never been set or remain the same if it has been set in the past.

The <key_name> specifies which validator you are editing. If you choose to not include some of the flags below, remember that the --from flag **must** be included to identify the validator to update.

The `--identity` can be used as to verify identity with systems like Keybase or UPort. When using Keybase, `--identity` should be populated with a 16-digit string that is generated with a [keybase.io](https://keybase.io) account. It's a cryptographically secure method of verifying your identity across multiple online networks. The Keybase API allows us to retrieve your Keybase avatar. This is how you can add a logo to your validator profile.

```bash
trstd tx staking edit-validator
  --moniker="choose a moniker" \
  --website="https://trstlabs.xyz" \
  --identity=6A0D65E29A4CBC8E \
  --details="To infinity and beyond!" \
  --chain-id=<chain_id> \
  --gas="auto" \
  --gas-prices="0.0025utrst" \
  --from=<key_name> \
  --commission-rate="0.10"
```

::: danger Warning
Please note that some parameters such as `commission-max-rate` and `commission-max-change-rate` cannot be changed once your validator is up and running.
:::

**Note**: The `commission-rate` value must adhere to the following rules:

- Must be between 0 and the validator's `commission-max-rate`
- Must not exceed the validator's `commission-max-change-rate` which is maximum
  % point change rate **per day**. In other words, a validator can only change
  its commission once per day and within `commission-max-change-rate` bounds.

## View Validator Description

View the validator's information with this command:

```bash
trstd query staking validator <account_cosmos>
```

## Track Validator Signing Information

In order to keep track of a validator's signatures in the past you can do so by using the `signing-info` command:

```bash
trstd query slashing signing-info <validator-pubkey>\
  --chain-id=<chain_id>
```

## Unjail Validator

When a validator is "jailed" for downtime, you must submit an `Unjail` transaction from the operator account in order to be able to get block proposer rewards again (depends on the zone fee distribution).

```bash
trstd tx slashing unjail \
 --from=<key_name> \
 --chain-id=<chain_id>
```

## Confirm Your Validator is Running

Your validator is active if the following command returns anything:

```bash
trstd query tendermint-validator-set | grep "$(trstd tendermint show-address)"
```

You should now see your validator in one of the explorers. You are looking for the `bech32` encoded `address` in the `~/.trst/config/priv_validator.json` file.

## Halting Your Validator

When attempting to perform routine maintenance or planning for an upcoming coordinated upgrade, it can be useful to have your validator systematically and gracefully halt. You can achieve this by either setting the `halt-height` to the height at which you want your node to shutdown or by passing the `--halt-height` flag to `trstd`. The node will shutdown with a zero exit code at that given height after committing
the block.

## Advanced configuration

You can find more advanced information about running a node or a validator on the [CometBFT Core documentation](https://docs.cometbft.com/v0.34/core/validators).
