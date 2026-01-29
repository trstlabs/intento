# Curated Proof of Stake (Governance-Gated Validator Set)

Intento implements a **Curated Proof of Stake** model built on top of the standard Cosmos SDK Proof of Stake (PoS) module.

Validator participation is **governance-gated**. While the network uses standard PoS mechanics for stake bonding, voting power, and consensus, admission to the validator set is controlled via on-chain governance.

This results in a **permissioned, economically bonded validator set**, combining authority-based admission with PoS security guarantees.

## Concept

The validator set is gated by governance:

* **Permissionless validator creation is disabled**
  Standard users cannot submit `MsgCreateValidator` transactions.

* **Validators are added via governance**
  A dedicated `ValidatorAddProposal` must be submitted and approved to admit a new validator.

This creates a curated validator set while retaining PoS-based stake weighting, slashing, and consensus behavior.

## Mechanics

### Gatekeeping

The application uses an **AnteHandler** (`GateCreateValidatorAnteHandler`) to intercept and block all `MsgCreateValidator` transactions from entering the mempool or being executed in a block.

This enforces governance-controlled admission at the protocol level.

Gatekeeping can be disabled for testnet or development environments by setting:

```
intento.poa.disable_gatekeeping = true
```

### Adding a Validator

Validators are added via a governance proposal (`ValidatorAddProposal`).

This proposal is handled by `ValidatorAdminProposalHandler`, which internally executes a privileged `MsgCreateValidator`, bypassing the AnteHandler gate.

**Proposal Fields:**

* `Title`, `Description`: Standard proposal metadata
* `Valoper`: Operator address of the new validator
* `Moniker`: Validator name
* `PubKey`: Consensus public key (Ed25519)

### Removing a Validator

Validators can be removed via `ValidatorRemoveProposal`.

This proposal forcibly unbonds the validator’s self-delegation, causing removal from the active validator set after the unbonding period.

**Proposal Fields:**

* `Valoper`: Operator address of the validator to remove

## Governance and Validator Set Management

The initial validator set (originating from ICS) forms the initial governance body.

This body controls:

* Admission of new validators
* Removal of existing validators
* Ongoing management of the curated validator set

This ensures validator participation is explicitly approved while remaining economically aligned via stake bonding and slashing.

## Configuration

To disable governance gatekeeping and allow permissionless validator creation (for testing or development), set:

```bash
--intento.poa.disable_gatekeeping=true
```

## Join the Network

For current validator sets, genesis files, and persistent peers, see the **Networks Repository**:
[https://github.com/trstlabs/networks](https://github.com/trstlabs/networks)

To apply for validator participation or ask questions, join **Discord**:
[https://discord.gg/hsVf9sYyZW](https://discord.gg/hsVf9sYyZW)
