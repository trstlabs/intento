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
    Ibc(IbcQuery),
    Stargate {
        /// this is the fully qualified service path used for routing,
        /// eg. custom/cosmos_sdk.x.bank.v1.Query/QueryBalance
        path: String,
        /// this is the expected protobuf message type (not any), binary encoded
        data: Binary,
    },
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

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq)]
#[serde(rename_all = "snake_case")]
pub enum DistQuery {
    /// This calls into the native bank module for all denominations.
    /// Note that this may be much more expensive than Balance and should be avoided if possible.
    /// Return value is AllBalanceResponse.
    Rewards { delegator: String },
}

#[non_exhaustive]
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
        key: Binary,
    },
     /// this queries the raw kv-store of the contract.
    /// returns the raw, unparsed data stored at that key (or `Ok(Err(StdError:NotFound{}))` if missing)
    PublicForAddr {
        contract_addr: String,
        account_addr: String,
        /// Key is the key used in the public contract's Storage
        key: Binary,
    },
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
    /// Returns all the unbonding delegations by the delegator
    UnbondingDelegations { delegator: String },
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq)]
#[serde(rename_all = "snake_case")]
pub enum IbcQuery {
    /// Gets the Port ID the current contract is bound to.
    ///
    /// Returns a `PortIdResponse`.
    PortId {},
    /// Lists all channels that are bound to a given port.
    /// If `port_id` is omitted, this list all channels bound to the contract's port.
    ///
    /// Returns a `ListChannelsResponse`.
    ListChannels { port_id: Option<String> },
    /// Lists all information for a (portID, channelID) pair.
    /// If port_id is omitted, it will default to the contract's own channel.
    /// (To save a PortId{} call)
    ///
    /// Returns a `ChannelResponse`.
    Channel {
        channel_id: String,
        port_id: Option<String>,
    },
}
