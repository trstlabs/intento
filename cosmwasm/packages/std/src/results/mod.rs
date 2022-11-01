//! This module contains the messages that are sent from the contract to the VM as an execution result

mod contract_result;
mod cosmos_msg;
mod empty;
mod events;
mod query;
mod response;
mod submessages;
mod system_result;

pub use contract_result::ContractResult;
pub use cosmos_msg::{wasm_execute, wasm_instantiate,wasm_instantiate_auto, BankMsg, CosmosMsg, WasmMsg, CustomMsg};
#[cfg(feature = "staking")]
pub use cosmos_msg::{DistributionMsg, StakingMsg};
#[cfg(feature = "stargate")]
pub use cosmos_msg::{GovMsg, VoteOption};
pub use empty::Empty;
pub use events::{attr, log, log_plaintext, pub_db, acc_pub_db, acc_pub_db_bytes, pub_db_bytes,  Attribute, Event};
pub use query::QueryResponse;
pub use response::Response;
pub use submessages::{Reply, ReplyOn, SubMsg,SubMsgResponse, SubMsgResult};
pub use system_result::SystemResult;
