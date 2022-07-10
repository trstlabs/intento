use schemars::JsonSchema;
use serde::{Deserialize, Serialize};

/// A full [*Cosmos SDK* event].
///
/// This version uses string attributes (similar to [*Cosmos SDK* StringEvent]),
/// which then get magically converted to bytes for Tendermint somewhere between
/// the Rust-Go interface, JSON deserialization and the `NewEvent` call in Cosmos SDK.
///
/// [*Cosmos SDK* event]: https://docs.cosmos.network/v0.42/core/events.html
/// [*Cosmos SDK* StringEvent]: https://github.com/cosmos/cosmos-sdk/blob/v0.42.5/proto/cosmos/base/abci/v1beta1/abci.proto#L56-L70
#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
#[non_exhaustive]
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

impl Event {
    /// Create a new event with the given type and an empty list of attributes.
    pub fn new(ty: impl Into<String>) -> Self {
        Event {
            ty: ty.into(),
            attributes: Vec::with_capacity(10),
        }
    }

    /// Add an attribute to the event.
    pub fn add_attribute(mut self, key: impl Into<String>, value: impl Into<String>) -> Self {
        self.attributes.push(Attribute {
            key: key.into(),
            value: value.into().as_bytes().to_vec(),
            acc_addr: None,
            pub_db: false,
            encrypted: true,
        });
        self
    }

    /// Add a plaintext attribute to the event.
    pub fn add_attribute_plaintext(
        mut self,
        key: impl Into<String>,
        value: impl Into<String>,
    ) -> Self {
        self.attributes.push(Attribute {
            key: key.into(),
            value: value.into().as_bytes().to_vec(),
            acc_addr: None,
            pub_db: false,
            encrypted: false,
        });
        self
    }

    /// Add a plaintext attribute to the event.
    pub fn store_pub_db(mut self, key: impl Into<String>, value: impl Into<String>) -> Self {
        self.attributes.push(Attribute {
            key: key.into(),
            value: value.into().as_bytes().to_vec(),
            acc_addr: None,
            pub_db: true,
            encrypted: false,
        });
        self
    }
    /// Add a plaintext attribute to the event.
    pub fn store_acc_pub_db(
        mut self,
        key: impl Into<String>,
        value: impl Into<String>,
        acc: impl Into<String>,
    ) -> Self {
        self.attributes.push(Attribute {
            key: key.into(),
            value: value.into().as_bytes().to_vec(),
            acc_addr: Some(acc.into()),
            pub_db: true,
            encrypted: false,
        });
        self
    }

    /// Bulk add attributes to the event.
    ///
    /// Anything that can be turned into an iterator and yields something
    /// that can be converted into an `Attribute` is accepted.
    pub fn add_attributes<A: Into<Attribute>>(
        mut self,
        attrs: impl IntoIterator<Item = A>,
    ) -> Self {
        self.attributes.extend(attrs.into_iter().map(A::into));
        self
    }
}

/// Return true
///
/// Only used for serde annotations
fn bool_true() -> bool {
    true
}

/// An key value pair that is used in the context of event attributes in logs
#[derive(Serialize, Deserialize, Clone, Default, Debug, PartialEq, JsonSchema)]
pub struct Attribute {
    pub key: String,
    pub value: Vec<u8>,
    pub pub_db: bool,
    pub acc_addr: Option<String>,
    pub encrypted: bool,
}

impl Attribute {
    /// Creates a new log.
    pub fn new_log(key: impl Into<String>, value: impl Into<String>) -> Self {
        let key = key.into();

        #[cfg(debug_assertions)]
        if key.starts_with('_') {
            panic!(
                "attribute key `{}` is invalid - keys starting with an underscore are reserved",
                key
            );
        }

        Self {
            key,
            value: value.into().as_bytes().to_vec(),
            pub_db: false,
            acc_addr: None,
            encrypted: true,
        }
    }

    /// Creates a new plaintext log.
    pub fn new_plaintext_log(key: impl Into<String>, value: impl Into<String>) -> Self {
        let key = key.into();

        #[cfg(debug_assertions)]
        if key.starts_with('_') {
            panic!(
                "attribute key `{}` is invalid - keys starting with an underscore are reserved",
                key
            );
        }

        Self {
            key,
            value: value.into().as_bytes().to_vec(),
            pub_db: false,
            acc_addr: None,
            encrypted: false,
        }
    }
    /// Creates a new public state Attribute.
    pub fn store_pub_db(key: impl Into<String>, value: impl Into<String>) -> Self {
        let key = key.into();

        #[cfg(debug_assertions)]
        if key.starts_with('_') {
            panic!(
                "attribute key `{}` is invalid - keys starting with an underscore are reserved",
                key
            );
        }

        Self {
            key,
            value: value.into().as_bytes().to_vec(),
            pub_db: true,
            acc_addr: None,
            encrypted: false,
        }
    }
    /// Creates a new public state Attribute for an account
    pub fn store_acc_pub_db(
        key: impl Into<String>,
        value: impl Into<String>,
        addr: impl Into<String>,
    ) -> Self {
        let key = key.into();

        #[cfg(debug_assertions)]
        if key.starts_with('_') {
            panic!(
                "attribute key `{}` is invalid - keys starting with an underscore are reserved",
                key
            );
        }

        Self {
            key,
            value: value.into().as_bytes().to_vec(),
            pub_db: true,
            acc_addr: Some(addr.into()),
            encrypted: false,
        }
    }
    /// Creates a new public state Attribute for an account
    pub fn store_pub_db_bytes(key: impl Into<String>, value: &[u8]) -> Self {
        let key = key.into();

        #[cfg(debug_assertions)]
        if key.starts_with('_') {
            panic!(
                "attribute key `{}` is invalid - keys starting with an underscore are reserved",
                key
            );
        }

        Self {
            key,
            value: value.to_vec(),
            pub_db: true,
            acc_addr: None,
            encrypted: false,
        }
    }
    /// Creates a new public state Attribute for an account
    pub fn store_acc_pub_db_bytes(
        key: impl Into<String>,
        value: &[u8],
        addr: impl Into<String>,
    ) -> Self {
        let key = key.into();

        #[cfg(debug_assertions)]
        if key.starts_with('_') {
            panic!(
                "attribute key `{}` is invalid - keys starting with an underscore are reserved",
                key
            );
        }

        Self {
            key,
            value: value.to_vec(),
            pub_db: true,
            acc_addr: Some(addr.into()),
            encrypted: false,
        }
    }
}

impl<K: Into<String>, V: Into<String>> From<(K, V)> for Attribute {
    fn from((k, v): (K, V)) -> Self {
        Attribute::new_log(k, v)
    }
}

impl<K: AsRef<str>, V: AsRef<Vec<u8>>> PartialEq<(K, V)> for Attribute {
    fn eq(&self, (k, v): &(K, V)) -> bool {
        (self.key.as_str(),  self.value.clone()) == (k.as_ref(), v.as_ref().to_vec())
    }
}

impl<K: AsRef<str>, V: AsRef<str> + std::convert::AsRef<std::vec::Vec<u8>>> PartialEq<Attribute> for (K, V) {
    fn eq(&self, attr: &Attribute) -> bool {
        attr == self
    }
}

impl<K: AsRef<str>, V: AsRef<Vec<u8>>> PartialEq<(K, V)> for &Attribute {
    fn eq(&self, (k, v): &(K, V)) -> bool {
        (self.key.as_str(), self.value.clone()) == (k.as_ref(), v.as_ref().to_vec())
    }
}

impl<K: AsRef<str>, V: AsRef<str> + std::convert::AsRef<std::vec::Vec<u8>>> PartialEq<&Attribute>
    for (K, V)
{
    fn eq(&self, attr: &&Attribute) -> bool {
        attr == self
    }
}

impl PartialEq<Attribute> for &Attribute {
    fn eq(&self, rhs: &Attribute) -> bool {
        *self == rhs
    }
}

impl PartialEq<&Attribute> for Attribute {
    fn eq(&self, rhs: &&Attribute) -> bool {
        self == *rhs
    }
}

/// Creates a new log. `Attribute::new_log` is an alias for this.
#[inline]
pub fn log(key: impl Into<String>, value: impl Into<String>) -> Attribute {
    Attribute::new_log(key, value)
}

/// Creates a new plaintext log. `Attribute::new` is an alias for this.
#[inline]
pub fn log_plaintext(key: impl Into<String>, value: impl Into<String>) -> Attribute {
    Attribute::new_plaintext_log(key, value)
}
/// A shorthand to set a public state key-value pair
#[inline]
pub fn store_pub_db(key: impl Into<String>, value: impl Into<String>) -> Attribute {
    Attribute::store_pub_db(key, value)
}
/// A shorthand to set a public state key-value pair for a given account
#[inline]
pub fn store_acc_pub_db(
    key: impl Into<String>,
    value: impl Into<String>,
    addr: impl Into<String>,
) -> Attribute {
    Attribute::store_acc_pub_db(key, value, addr)
}
/// A shorthand to set a public state key-value pair
#[inline]
pub fn store_pub_db_bytes(key: impl Into<String>, value: &[u8]) -> Attribute {
    Attribute::store_pub_db_bytes(key, value)
}
/// A shorthand to set a public state key-value pair for a given account
#[inline]
pub fn store_acc_pub_db_bytes(
    key: impl Into<String>,
    value: &[u8],
    addr: impl Into<String>,
) -> Attribute {
    Attribute::store_acc_pub_db_bytes(key, value, addr)
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::Uint128;

    #[test]
    fn event_construction() {
        let event_direct = Event {
            ty: "test".to_string(),
            attributes: vec![attr("foo", "bar"), attr("bar", "baz")],
        };
        let event_builder = Event::new("test").add_attributes(vec![("foo", "bar"), ("bar", "baz")]);

        assert_eq!(event_direct, event_builder);
    }

    #[test]
    #[should_panic]
    fn attribute_new_reserved_key_panicks() {
        Attribute::new("_invalid", "value");
    }

    #[test]
    #[should_panic]
    fn attribute_new_reserved_key_panicks2() {
        Attribute::new("_", "value");
    }

    #[test]
    fn attr_works_for_different_types() {
        let expected = ("foo", "42");

        assert_eq!(attr("foo", "42"), expected);
        assert_eq!(attr("foo", "42"), expected);
        assert_eq!(attr("foo", "42"), expected);
        assert_eq!(attr("foo", Uint128::new(42)), expected);
    }
}
