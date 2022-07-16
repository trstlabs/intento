

use serde::{Deserialize, Serialize};

use super::addresses::Addr;
use super::coins::Coin;
use super::timestamp::Timestamp;

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq)]
pub struct FullEnv {
    pub block: BlockInfo,
    pub message: MessageInfo,
    pub contract: ContractInfo,
    pub contract_key: Option<String>,
}
