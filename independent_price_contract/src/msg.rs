use schemars::JsonSchema;
use serde::{Deserialize, Serialize};

use cosmwasm_std::{Binary, CosmosMsg, HumanAddr, Querier, StdResult, Uint128};

use secret_toolkit::snip20::{register_receive_msg, token_info_query, transfer_msg, TokenInfo};

use crate::contract::BLOCK_SIZE;

/// Instantiation message
#[derive(Serialize, Deserialize, JsonSchema)]
pub struct InitMsg {
    /// sell contract code hash and address
    pub sell_contract: ContractInfo,
    /// estimation contract code hash and address
    pub estimation_contract: ContractInfo,
    /// amount of tokens being sold
    pub sell_amount: Uint128,
    /// minimum estimation that will be accepted
    pub minimum_estimation: Uint128,

    pub estimationcount: usize,
    /// Optional free-form description of the estimationperiod (best to avoid double quotes). As an example
    /// it could be the date the owner will likely finalize the estimationperiod, or a list of other
    /// estimationperiods for the same token, etc...
    #[serde(default)]
    pub description: Option<String>,
}

/// Handle messages
#[derive(Serialize, Deserialize, JsonSchema)]
#[serde(rename_all = "snake_case")]
pub enum HandleMsg {
    /// Receive gets called by the token contracts of the estimationperiod.  If it came from the sale token, it
    /// will consign the sent tokens.  If it came from the estimation token, it will place a estimation.  If any
    /// other address tries to call this, it will give an error message that the calling address is
    /// not a token in the estimationperiod.
    CreateEstimation {
        /// address of person or contract that sent the tokens that triggered this Receive
        sender: HumanAddr,
        /// address of the owner of the tokens sent to the estimationperiod
        from: HumanAddr,
        /// amount of tokens sent
        amount: Uint128,
        /// Optional base64 encoded message sent with the Send call -- not needed or used by this
        /// contract
        #[serde(default)]
        msg: Option<Binary>,
    },

    /// RetractEstimation will retract any active estimation the calling address has made and return the tokens
    /// that are held in escrow
    RetractEstimation {},

    /// ViewEstimation will display the amount of the active estimation made by the calling address and time the
    /// estimation was placed
   // ViewEstimation {},

    /// Finalize will close the estimationperiod
    RevealEstimation {},
  /// true if estimationperiod creator wants to keep the estimationperiod open if there are no active estimations
       // only_if_estimations: bool,
    /// If the estimationperiod holds any funds after it has closed (should never happen), this will return
    /// those funds to their owners.  Should never be needed, but included in case of unforeseen
    /// error
    ReturnAll {},
}

/// Queries
#[derive(Serialize, Deserialize, JsonSchema)]
#[serde(rename_all = "snake_case")]
pub enum QueryMsg {
    /// Displays the estimationperiod information
    EstimationPeriodInfo {},
}

/// responses to queries
#[derive(Serialize, Deserialize, Debug, JsonSchema)]
#[serde(rename_all = "snake_case")]
pub enum QueryAnswer {
    /// EstimationPeriodInfo query response
    EstimationPeriodInfo {
        /// sell token address and TokenInfo query response
       // sell_token: Token,
        /// estimation token address and TokenInfo query response
       // estimation_token: Token,
        /// amount of tokens being sold
      //  sell_amount: Uint128,
        /// minimum estimation that will be accepted
       // minimum_estimation: Uint128,
        /// Optional String description of estimationperiod
        #[serde(skip_serializing_if = "Option::is_none")]
        description: Option<String>,
        /// address of estimationperiod contract
        estimationperiod_address: HumanAddr,
        /// status of the estimationperiod can be "Accepting estimations: Tokens to be sold have(not) been
        /// consigned" or "Closed" (will also state if there are outstanding funds after estimationperiod
        /// closure
        status: String,
        /// If the estimationperiod resulted in a swap, this will state the winning estimation
        #[serde(skip_serializing_if = "Option::is_none")]
        best_estimation: Option<Uint128>,
    },
}

/*/// token's contract address and TokenInfo response
#[derive(Serialize, Deserialize, Debug, JsonSchema)]
pub struct Token {
    /// contract address of token
    pub contract_address: HumanAddr,
    /// Tokeninfo query response
    pub token_info: TokenInfo,
}*/

/// success or failure response
#[derive(Serialize, Deserialize, Debug, JsonSchema)]
pub enum ResponseStatus {
    Success,
    Failure,
}

/// Responses from handle functions
#[derive(Serialize, Deserialize, Debug, JsonSchema)]
#[serde(rename_all = "snake_case")]
pub enum HandleAnswer {
    /// response from consign attempt
    /*Consign {
        /// success or failure
        status: ResponseStatus,
        /// execution description
        message: String,
        /// Optional amount consigned
        #[serde(skip_serializing_if = "Option::is_none")]
        amount_consigned: Option<Uint128>,
        /// Optional amount that still needs to be consigned
        #[serde(skip_serializing_if = "Option::is_none")]
        amount_needed: Option<Uint128>,
        /// Optional amount of tokens returned from escrow
        #[serde(skip_serializing_if = "Option::is_none")]
        amount_returned: Option<Uint128>,
    },*/
    /// response from estimation attempt
    Estimation {
        /// success or failure
        status: ResponseStatus,
        /// execution description
        message: String,
        /// Optional amount estimation
        #[serde(skip_serializing_if = "Option::is_none")]
        amount_estimation: Option<Uint128>,
    },
    /// response from closing the estimationperiod
    RevealEstimation {
        /// success or failure
        status: ResponseStatus,
        /// execution description
        message: String,
        /// Optional amount of winning estimation
        #[serde(skip_serializing_if = "Option::is_none")]
        best_estimation: Option<Uint128>,

    },
    /// response from attempt to retract estimation
    RetractEstimation {
        /// success or failure
        status: ResponseStatus,
        /// execution description
         /// Optional amount of tokens returned from escrow
       // #[serde(skip_serializing_if = "Option::is_none")]
      //  amount_returned: Option<Uint128>,
        message: String,
       
    },
    /// generic status response
    Status {
        /// success or failure
        status: ResponseStatus,
        /// execution description
        message: String,
    },
}

/// code hash and address of a contract
#[derive(Serialize, Deserialize, JsonSchema)]
pub struct ContractInfo {
    /// contract's code hash string
    pub code_hash: String,
    /// contract's address
    pub address: HumanAddr,
}

impl ContractInfo {
    /// Returns a StdResult<CosmosMsg> used to execute Transfer
    ///
    /// # Arguments
    ///
    /// * `recipient` - address tokens are to be sent to
    /// * `amount` - Uint128 amount of tokens to send
    pub fn transfer_msg(&self, recipient: HumanAddr, amount: Uint128) -> StdResult<CosmosMsg> {
        transfer_msg(
            recipient,
            amount,
            None,
            BLOCK_SIZE,
            self.code_hash.clone(),
            self.address.clone(),
        )
    }

    /// Returns a StdResult<CosmosMsg> used to execute RegisterReceive
    ///
    /// # Arguments
    ///
    /// * `code_hash` - String holding code hash contract to be called when sent tokens
    pub fn register_receive_msg(&self, code_hash: String) -> StdResult<CosmosMsg> {
        register_receive_msg(
            code_hash,
            None,
            BLOCK_SIZE,
            self.code_hash.clone(),
            self.address.clone(),
        )
    }

    /// Returns a StdResult<TokenInfo> from performing TokenInfo query
    ///
    /// # Arguments
    ///
    /// * `querier` - a reference to the Querier dependency of the querying contract
    pub fn token_info_query<Q: Querier>(&self, querier: &Q) -> StdResult<TokenInfo> {
        token_info_query(
            querier,
            BLOCK_SIZE,
            self.code_hash.clone(),
            self.address.clone(),
        )
    }
}
