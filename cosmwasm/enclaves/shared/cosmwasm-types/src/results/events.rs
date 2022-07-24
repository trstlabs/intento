
use serde::{Deserialize, Serialize};

/// A full [*Cosmos SDK* event].
///
/// This version uses string attributes (similar to [*Cosmos SDK* StringEvent]),
/// which then get magically converted to bytes for Tendermint somewhere between
/// the Rust-Go interface, JSON deserialization and the `NewEvent` call in Cosmos SDK.
///
/// [*Cosmos SDK* event]: https://docs.cosmos.network/v0.42/core/events.html
/// [*Cosmos SDK* StringEvent]: https://github.com/cosmos/cosmos-sdk/blob/v0.42.5/proto/cosmos/base/abci/v1beta1/abci.proto#L56-L70
#[derive(Serialize, Deserialize, Clone, Debug, PartialEq)]
pub struct Event {
    /// The event type. This is renamed to "ty" because "type" is reserved in Rust. This sucks, we know.
    #[serde(rename = "type")]
    pub ty: String,
    /// The attributes to be included in the event.
    ///
    /// You can learn more about these from [*Cosmos SDK* docs].
    ///
    /// [*Cosmos SDK* docs]: https://docs.cosmos.network/v0.42/core/events.html
    pub attributes: Vec<Attribute>,
}
/// An key value pair that is used in the context of event attributes in logs
#[derive(Serialize, Deserialize, Clone, Debug, PartialEq)]
pub struct Attribute {
    pub key: String,
    #[serde(with = "serde_bytes")]
    pub value: Vec<u8>,
    #[serde(skip_deserializing)]
    pub pub_db: bool,
    #[serde(skip_deserializing)]
    pub acc_addr: Option<String>,
    #[serde(skip_deserializing)]
    pub encrypted: bool,
}

impl Attribute {
    /// helper.
    pub fn to_kv(&self) -> Self {

        Self {
            key: self.key.clone(),
            value: self.value.clone(),
            pub_db: false,
            acc_addr: None,
            encrypted: false,
        }
    }
}


