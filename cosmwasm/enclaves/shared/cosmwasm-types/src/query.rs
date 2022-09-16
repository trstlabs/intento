//! must keep this file in sync with cosmwasm/packages/std/src/query.rs

use serde::{Deserialize, Serialize};

use super::coins::Coin;
use super::encoding::Binary;


#[derive(Serialize, Deserialize, Clone, Debug, PartialEq)]
#[serde(rename_all = "snake_case")]
pub enum QueryRequest {
    Bank(BankQuery),
    Custom(serde_json::Value),
    Staking(StakingQuery),
    Wasm(WasmQuery),
    Dist(DistQuery),
    Mint(MintQuery),
    Gov(GovQuery),
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq)]
#[serde(rename_all = "snake_case")]
pub enum MintQuery {
    /// This calls into the native bank module for all denominations.
    /// Note that this may be much more expensive than Balance and should be avoided if possible.
    /// Return value is AllBalanceResponse.
    Inflation {},
    BondedRatio {},
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq)]
#[serde(rename_all = "snake_case")]
pub enum BankQuery {
    /// This calls into the native bank module for one denomination
    /// Return value is BalanceResponse
    Balance { address: String, denom: String },
    /// This calls into the native bank module for all denominations.
    /// Note that this may be much more expensive than Balance and should be avoided if possible.
    /// Return value is AllBalanceResponse.
    AllBalances { address: String },
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq)]
#[serde(rename_all = "snake_case")]
pub enum GovQuery {
    /// Returns all the currently active proposals. Might be useful to filter out invalid votes, and trigger
    /// in-contract voting periods
    Proposals {},
}

/// ProposalsResponse is data format returned from GovQuery::Proposals query
#[derive(Serialize, Deserialize, Clone, Debug, PartialEq)]
#[serde(rename_all = "snake_case")]
pub struct ProposalsResponse {
    pub proposals: Vec<Proposal>,
}

/// ProposalsResponse is data format returned from GovQuery::Proposals query
#[derive(Serialize, Deserialize, Clone, Debug, PartialEq)]
#[serde(rename_all = "snake_case")]
pub struct Proposal {
    pub id: u64,
    /// Time of the block where MinDeposit was reached. -1 if MinDeposit is not reached
    pub voting_start_time: u64,
    /// Time that the VotingPeriod for this proposal will end and votes will be tallied
    pub voting_end_time: u64,
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq)]
#[serde(rename_all = "snake_case")]
pub enum DistQuery {
    /// This calls into the native bank module for all denominations.
    /// Note that this may be much more expensive than Balance and should be avoided if possible.
    /// Return value is AllBalanceResponse.
    Rewards { delegator: String },
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq)]
#[serde(rename_all = "snake_case")]
pub enum WasmQuery {
    /// this queries the public API of another contract at a known address (with known ABI)
    /// return value is whatever the contract returns (caller should know)
    Private {
        contract_addr: String,
        /// code_hash is the hex encoded hash of the code. This is used by trst to harden against replaying the contract
        /// It is used to bind the request to a destination contract in a stronger way than just the contract address which can be faked
        code_hash: String,
        /// msg is the json-encoded QueryMsg struct
        msg: Binary,
    },
    /// this queries the raw kv-store of the contract.
    /// returns the raw, unparsed data stored at that key (or `Ok(Err(StdError:NotFound{}))` if missing)
    Public {
        contract_addr: String,
        /// code_hash is the hex encoded hash of the code. This is used by trst to harden against replaying the contract
        /// It is used to bind the request to a destination contract in a stronger way than just the contract address which can be faked
        //code_hash: String,
        /// Key is the key used in the public contract's Storage
        key: String,
    },
     /// this queries the raw kv-store of the contract.
    /// returns the raw, unparsed data stored at that key (or `Ok(Err(StdError:NotFound{}))` if missing)
    PublicForAddr {
        contract_addr: String,
        account_addr: String,
        /// Key is the key used in the public contract's Storage
        key: String,
    },
}

impl From<GovQuery> for QueryRequest {
    fn from(msg: GovQuery) -> Self {
        QueryRequest::Gov(msg)
    }
}

impl From<MintQuery> for QueryRequest {
    fn from(msg: MintQuery) -> Self {
        QueryRequest::Mint(msg)
    }
}

impl From<DistQuery> for QueryRequest {
    fn from(msg: DistQuery) -> Self {
        QueryRequest::Dist(msg)
    }
}

impl From<BankQuery> for QueryRequest {
    fn from(msg: BankQuery) -> Self {
        QueryRequest::Bank(msg)
    }
}

#[cfg(feature = "staking")]
impl From<StakingQuery> for QueryRequest {
    fn from(msg: StakingQuery) -> Self {
        QueryRequest::Staking(msg)
    }
}

impl From<WasmQuery> for QueryRequest {
    fn from(msg: WasmQuery) -> Self {
        QueryRequest::Wasm(msg)
    }
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq)]
#[serde(rename_all = "snake_case")]
pub enum StakingQuery {
    /// Returns the denomination that can be bonded (if there are multiple native tokens on the chain)
    BondedDenom {},
    /// AllDelegations will return all delegations by the delegator
    AllDelegations { delegator: String },
    /// Delegation will return more detailed info on a particular
    /// delegation, defined by delegator/validator pair
    Delegation {
        delegator: String,
        validator: String,
    },
    /// Returns all registered Validators on the system
    Validators {},
     /// The query response type is `ValidatorResponse`.
    Validator {
        /// The validator's address (e.g. (e.g. cosmosvaloper1...))
        address: String,
    },
    /// Returns all the unbonding delegations by the delegator
    UnbondingDelegations { delegator: String },
}

/// Delegation is basic (cheap to query) data about a delegation
#[derive(Serialize, Deserialize, Clone, Debug, PartialEq)]
pub struct Delegation {
    pub delegator: String,
    pub validator: String,
    /// How much we have locked in the delegation
    pub amount: Coin,
}

impl From<FullDelegation> for Delegation {
    fn from(full: FullDelegation) -> Self {
        Delegation {
            delegator: full.delegator,
            validator: full.validator,
            amount: full.amount,
        }
    }
}

/// UnbondingDelegationsResponse is data format returned from StakingRequest::UnbondingDelegations query
#[derive(Serialize, Deserialize, Clone, Debug, PartialEq)]
#[serde(rename_all = "snake_case")]
pub struct UnbondingDelegationsResponse {
    pub delegations: Vec<Delegation>,
}

/// FullDelegation is all the info on the delegation, some (like accumulated_reward and can_redelegate)
/// is expensive to query
#[derive(Serialize, Deserialize, Clone, Debug, PartialEq)]
pub struct FullDelegation {
    pub delegator: String,
    pub validator: String,
    /// How much we have locked in the delegation
    pub amount: Coin,
    /// can_redelegate captures how much can be immediately redelegated.
    /// 0 is no redelegation and can_redelegate == amount is redelegate all
    /// but there are many places between the two
    pub can_redelegate: Coin,
    /// How much we can currently withdraw
    pub accumulated_rewards: Coin,
}

/// Delegation is basic (cheap to query) data about a delegation
#[derive(Serialize, Deserialize, Clone, Debug, PartialEq)]
pub struct RewardsResponse {
    pub rewards: Vec<ValidatorRewards>,
    pub total: Vec<Coin>,
}

/// Delegation is basic (cheap to query) data about a delegation
#[derive(Serialize, Deserialize, Clone, Debug, PartialEq)]
pub struct ValidatorRewards {
    pub validator_address: String,
    pub reward: Vec<Coin>,
}

/// Inflation response
#[derive(Serialize, Deserialize, Clone, Debug, PartialEq)]
pub struct InflationResponse {
    pub inflation_rate: String,
}

/// Bonded Ratio response
#[derive(Serialize, Deserialize, Clone, Debug, PartialEq)]
pub struct BondedRatioResponse {
    pub bonded_ratio: String,
}