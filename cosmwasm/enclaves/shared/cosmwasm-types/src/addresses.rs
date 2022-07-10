#![allow(deprecated)]

use serde::{Deserialize, Serialize};
use std::fmt;
use crate::consts::BECH32_PREFIX_ACC_ADDR;
use crate::encoding::Binary;
use bech32::{FromBase32, ToBase32};
/// A human readable address.
///
/// In Cosmos, this is typically bech32 encoded. But for multi-chain smart contracts no
/// assumptions should be made other than being UTF-8 encoded and of reasonable length.
///
/// This type represents a validated address. It can be created in the following ways
/// 1. Use `Addr::unchecked(input)`
/// 2. Use `let checked: Addr = deps.api.addr_validate(input)?`
/// 3. Use `let checked: Addr = deps.api.addr_humanize(canonical_addr)?`
/// 4. Deserialize from JSON. This must only be done from JSON that was validated before
///    such as a contract's state. `Addr` must not be used in messages sent by the user
///    because this would result in unvalidated instances.
///
/// This type is immutable. If you really need to mutate it (Really? Are you sure?), create
/// a mutable copy using `let mut mutable = Addr::to_string()` and operate on that `String`
/// instance.
#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq, PartialOrd, Ord, Hash)]
pub struct Addr(pub String);

impl Addr {
    /// Creates a new `Addr` instance from the given input without checking the validity
    /// of the input. Since `Addr` must always contain valid addresses, the caller is
    /// responsible for ensuring the input is valid.
    ///
    /// Use this in cases where the address was validated before or in test code.
    /// If you see this in contract code, it should most likely be replaced with
    /// `let checked: Addr = deps.api.addr_humanize(canonical_addr)?`.
    ///
    /// ## Examples
    ///
    /// ```
    /// # use cosmwasm_std::{Addr};
    /// let address = Addr::unchecked("foobar");
    /// assert_eq!(address, "foobar");
    /// ```
    pub fn unchecked(input: impl Into<String>) -> Addr {
        Addr(input.into())
    }

    #[inline]
    pub fn as_str(&self) -> &str {
        self.0.as_str()
    }

    /// Returns the UTF-8 encoded address string as a byte array.
    ///
    /// This is equivalent to `address.as_str().as_bytes()`.
    #[inline]
    pub fn as_bytes(&self) -> &[u8] {
        self.0.as_bytes()
    }

    /// Utility for explicit conversion to `String`.
    #[inline]
    pub fn into_string(self) -> String {
        self.0
    }
}

impl fmt::Display for Addr {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "{}", &self.0)
    }
}

impl AsRef<str> for Addr {
    #[inline]
    fn as_ref(&self) -> &str {
        self.as_str()
    }
}

/// Implement `Addr == &str`
impl PartialEq<&str> for Addr {
    fn eq(&self, rhs: &&str) -> bool {
        self.0 == *rhs
    }
}

/// Implement `&str == Addr`
impl PartialEq<Addr> for &str {
    fn eq(&self, rhs: &Addr) -> bool {
        *self == rhs.0
    }
}

/// Implement `Addr == String`
impl PartialEq<String> for Addr {
    fn eq(&self, rhs: &String) -> bool {
        &self.0 == rhs
    }
}

/// Implement `String == Addr`
impl PartialEq<Addr> for String {
    fn eq(&self, rhs: &Addr) -> bool {
        self == &rhs.0
    }
}

impl From<Addr> for String {
    fn from(addr: Addr) -> Self {
        addr.0
    }
}

impl From<&Addr> for String {
    fn from(addr: &Addr) -> Self {
        addr.0.clone()
    }
}

//TODO rm HumanAddr
#[derive(Serialize, Deserialize, Clone, Default, Debug, PartialEq)]
pub struct HumanAddr(pub String);


impl HumanAddr {
    pub fn as_str(&self) -> &str {
        &self.0
    }
    pub fn len(&self) -> usize {
        self.0.len()
    }
    pub fn is_empty(&self) -> bool {
        self.0.is_empty()
    }
    pub fn from_canonical(canonical_addr: &CanonicalAddr) -> Result<Self, bech32::Error> {
        let human_addr_str = bech32::encode(
            BECH32_PREFIX_ACC_ADDR,
            canonical_addr.as_slice().to_base32(),
        )?;

        Ok(HumanAddr(human_addr_str))
    }
}

impl fmt::Display for HumanAddr {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "{}", &self.0)
    }
}

impl From<&str> for HumanAddr {
    fn from(addr: &str) -> Self {
        HumanAddr(addr.to_string())
    }
}

impl From<&HumanAddr> for HumanAddr {
    fn from(addr: &HumanAddr) -> Self {
        HumanAddr(addr.0.to_string())
    }
}

#[derive(Serialize, Deserialize, Clone, Default, Debug, PartialEq)]
pub struct CanonicalAddr(pub Binary);

impl CanonicalAddr {
    pub fn as_slice(&self) -> &[u8] {
        &self.0.as_slice()
    }
    pub fn len(&self) -> usize {
        self.0.len()
    }
    pub fn is_empty(&self) -> bool {
        self.0.is_empty()
    }
    pub fn from_human(human_addr: &HumanAddr) -> Result<Self, bech32::Error> {
        let (decoded_prefix, data) = bech32::decode(human_addr.as_str())?;
        let canonical = Vec::<u8>::from_base32(&data)?;

        Ok(CanonicalAddr(Binary(canonical)))
    }
    pub fn from_addr(val: &Addr) -> Result<Self, bech32::Error> {
        let (decoded_prefix, data) = bech32::decode(val.as_str())?;
        let canonical = Vec::<u8>::from_base32(&data)?;

        Ok(CanonicalAddr(Binary(canonical)))
    }
    pub fn from_str(val: &str) -> Result<Self, bech32::Error> {
        let (decoded_prefix, data) = bech32::decode(val)?;
        let canonical = Vec::<u8>::from_base32(&data)?;

        Ok(CanonicalAddr(Binary(canonical)))
    }

    pub fn from_vec(vec: Vec<u8>) -> Self {
        Self(Binary(vec))
    }
}

impl fmt::Display for CanonicalAddr {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        self.0.fmt(f)
    }
}