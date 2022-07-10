
use serde::{Deserialize, Serialize};

use crate::encoding::Binary;

#[non_exhaustive]
#[derive(Serialize, Deserialize, Clone, Debug, PartialEq)]
#[serde(rename_all = "snake_case")]
pub enum WasmQuery {
    /// this queries the public API of another contract at a known address (with known ABI)
    /// return value is whatever the contract returns (caller should know)
    Private {
        contract_addr: String,
        /// callback_code_hash is the hex encoded hash of the code. This is used by trst to harden against replaying the contract
        /// It is used to bind the request to a destination contract in a stronger way than just the contract address which can be faked
        callback_code_hash: String,
        /// msg is the json-encoded QueryMsg struct
        msg: Binary,
    },
    /// this queries the raw kv-store of the contract.
    /// returns the raw, unparsed data stored at that key (or `Ok(Err(StdError:NotFound{}))` if missing)
    Public {
        contract_addr: String,
        /// callback_code_hash is the hex encoded hash of the code. This is used by trst to harden against replaying the contract
        /// It is used to bind the request to a destination contract in a stronger way than just the contract address which can be faked
        //callback_code_hash: String,
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
