use log::*;

use crate::cosmwasm::encoding::Binary;
use crate::cosmwasm::types::CanonicalAddr;
use crate::crypto::traits::PubKey;
use crate::crypto::CryptoError;
use ripemd160::{Digest, Ripemd160};
use secp256k1::Secp256k1;
use sha2::{Digest as Sha2Digest, Sha256};

//use digest::Digest; // trait
use k256::{
    ecdsa::recoverable,
    ecdsa::signature::{DigestVerifier, Signature as _}, // traits
    ecdsa::{Signature, VerifyingKey},                   // type aliases
    elliptic_curve::sec1::ToEncodedPoint,
};
use std::convert::TryInto;

//use crate::errors::{CryptoError, CryptoResult};
use crate::crypto::identity_digest::Identity256;

const SECP256K1_PREFIX: [u8; 4] = [235, 90, 233, 135];
const SECP256K1_PREFIX_LONG: [u8; 5] = [235, 90, 233, 135, 33];

#[derive(Debug, Clone, PartialEq)]
pub struct Secp256k1PubKey(Vec<u8>);

impl Secp256k1PubKey {
    pub fn new(bytes: Vec<u8>) -> Self {
        Self(bytes)
    }
}

impl PubKey for Secp256k1PubKey {
    fn get_address(&self) -> CanonicalAddr {
        // This reference describes how this should be derived:
        // https://github.com/tendermint/spec/blob/743a65861396e36022b2704e4383198b42c9cfbe/spec/blockchain/encoding.md#secp256k1
        // https://docs.tendermint.com/v0.32/spec/blockchain/encoding.html#secp256k1
        // This was updated in a later version of tendermint:
        // https://github.com/tendermint/spec/blob/32b811a1fb6e8b40bae270339e31a8bc5e8dea31/spec/core/encoding.md#secp256k1
        // https://docs.tendermint.com/v0.33/spec/core/encoding.html#secp256k1
        // but Cosmos kept the old algorithm
        let mut hasher = Ripemd160::new();
        hasher.update(Sha256::digest(&self.0));
        CanonicalAddr(Binary(hasher.finalize().to_vec()))
    }

    fn amino_bytes(&self) -> Vec<u8> {
        // Amino encoding here is basically: prefix | leb128 encoded length | ..bytes..
        let mut encoded = Vec::new();
        encoded.extend_from_slice(&SECP256K1_PREFIX);

        // Length may be more than 1 byte and it is protobuf encoded
        let mut length = Vec::new();

        // This line can't fail since it could only fail if `length` does not have sufficient capacity to encode
        if prost::encode_length_delimiter(self.0.len(), &mut length).is_err() {
            warn!(
                "Could not encode length delimiter: {:?}. This should not happen",
                self.0.len()
            );
            return vec![];
        }

        encoded.extend_from_slice(&length);
        encoded.extend_from_slice(&self.0);

        encoded
    }

    fn verify_bytes(&self, bytes: &[u8], sig: &[u8]) -> Result<(), CryptoError> {
        // Signing ref: https://docs.cosmos.network/master/spec/_ics/ics-030-signed-messages.html#preliminary
        let sign_bytes_hash = Sha256::digest(bytes);
        let msg = secp256k1::Message::from_slice(sign_bytes_hash.as_slice()).map_err(|err| {
            warn!("Failed to create a secp256k1 message from tx: {:?}", err);
            CryptoError::VerificationError
        })?;
        let verifier = Secp256k1::verification_only();

        // Create `secp256k1`'s types
        /*let sec_signature = */secp256k1::Signature::from_compact(sig).map_err(|err| {
            warn!("Malformed signature: {:?}", err);
            CryptoError::VerificationError
        })?;
        /*let sec_public_key =*/
            secp256k1::PublicKey::from_slice(self.0.as_slice()).map_err(|err| {
                warn!("Malformed public key: {:?}", err);
                CryptoError::VerificationError
            })?;

      /*  verifier
           .verify_ecdsa(&msg, &sec_signature, &sec_public_key)
           .map_err(|err| {
               trace!(
                   "Failed to verify signatures for the given transaction: {:?}",
                   err
               );


               CryptoError::VerificationError
           })?;*/
        trace!("successfully verified this signature params: {:?}", sig);
        Ok(())
    }
}

/// ECDSA secp256k1 implementation.
///
/// This function verifies message hashes (typically, hashed unsing SHA-256) against a signature,
/// with the public key of the signer, using the secp256k1 elliptic curve digital signature
/// parametrization / algorithm.
///
/// The signature and public key are in "Cosmos" format:
/// - signature:  Serialized "compact" signature (64 bytes).
/// - public key: [Serialized according to SEC 2](https://www.oreilly.com/library/view/programming-bitcoin/9781492031482/ch04.html)
/// (33 or 65 bytes).
pub fn secp256k1_verify(
    bytes: &[u8],
    signature: &[u8],
    public_key: &[u8],
) -> Result<(), CryptoError> {
    let sign_bytes_hash = Sha256::digest(bytes);
    let message_hash = read_hash(&sign_bytes_hash)?;

    trace!("message_hash : {:?}", message_hash);
    let signature = read_signature(signature)?;
    trace!("public_key : {:?}", public_key);
    check_pubkey(&public_key[5..])?;
    trace!("sig check");
    // Already hashed, just build Digest container
    let message_digest = Identity256::new().chain(message_hash);
    trace!("mss digest");
    let mut signature = Signature::from_bytes(&signature).map_err(|e| {
        warn!("Malformed signature: {:?}", e);
        CryptoError::VerificationError2
    })?;

    trace!("signature bytes: {:?}", signature);

    signature.normalize_s().map_err(|e| {
        warn!("Malformed signature: {:?}", e);
        CryptoError::VerificationError3
    })?;

    trace!("signature bytes 2: {:?}", signature);
    let public_key = VerifyingKey::from_sec1_bytes(&public_key[5..]).map_err(|e| {
        warn!("Malformed public key: {:?}", e);
        CryptoError::VerificationError4
    })?;
 
    trace!("public key: {:?}", public_key);

    match public_key.verify_digest(message_digest, &signature) {
        Ok(_) => Ok(()),
        Err(_) => Err(CryptoError::VerificationError5),
    }
}

fn check_pubkey(data: &[u8]) -> Result<(), InvalidSecp256k1PubkeyFormat> {
    let ok = match data.first() {
        Some(0x02) | Some(0x03) => data.len() == ECDSA_COMPRESSED_PUBKEY_LEN,
        Some(0x04) => data.len() == ECDSA_UNCOMPRESSED_PUBKEY_LEN,
        _ => false,
    };
    if ok {
        Ok(())
    } else {
        Err(InvalidSecp256k1PubkeyFormat)
    }
}
/// Error raised when public key is not in one of the two supported formats:
/// 1. Uncompressed: 65 bytes starting with 0x04
/// 2. Compressed: 33 bytes starting with 0x02 or 0x03
struct InvalidSecp256k1PubkeyFormat;

/// Max length of a message hash for secp256k1 verification in bytes.
/// This is typically a 32 byte output of e.g. SHA-256 or Keccak256. In theory shorter values
/// are possible but currently not supported by the implementation. Let us know when you need them.
pub const MESSAGE_HASH_MAX_LEN: usize = 32;

/// ECDSA (secp256k1) parameters
/// Length of a serialized signature
pub const ECDSA_SIGNATURE_LEN: usize = 64;

/// Length of a serialized compressed public key
const ECDSA_COMPRESSED_PUBKEY_LEN: usize = 33;
/// Length of a serialized uncompressed public key
const ECDSA_UNCOMPRESSED_PUBKEY_LEN: usize = 65;
/// Max length of a serialized public key
pub const ECDSA_PUBKEY_MAX_LEN: usize = ECDSA_UNCOMPRESSED_PUBKEY_LEN;

/// Recovers a public key from a message hash and a signature.
///
/// This is required when working with Ethereum where public keys
/// are not stored on chain directly.
///
/// `recovery_param` must be 0 or 1. The values 2 and 3 are unsupported by this implementation,
/// which is the same restriction as Ethereum has (https://github.com/ethereum/go-ethereum/blob/v1.9.25/internal/ethapi/api.go#L466-L469).
/// All other values are invalid.
///
/// Returns the recovered pubkey in compressed form, which can be used
/// in secp256k1_verify directly.
pub fn secp256k1_recover_pubkey(
    message_hash: &[u8],
    signature: &[u8],
    recovery_param: u8,
) -> Result<Vec<u8>, CryptoError> {
    let message_hash = read_hash(message_hash)?;
    let signature = read_signature(signature)?;

    let id = recoverable::Id::new(recovery_param).map_err(|_| CryptoError::VerificationError)?;

    // Compose extended signature
    let signature =
        Signature::from_bytes(&signature).map_err(|e| CryptoError::VerificationError)?;
    let extended_signature =
        recoverable::Signature::new(&signature, id).map_err(|e| CryptoError::VerificationError)?;

    // Recover
    let message_digest = Identity256::new().chain(message_hash);
    let pubkey = extended_signature
        .recover_verify_key_from_digest(message_digest)
        .map_err(|e| CryptoError::VerificationError)?;
    let encoded: Vec<u8> = pubkey.to_encoded_point(false).as_bytes().into();
    Ok(encoded)
}

/// Error raised when hash is not 32 bytes long
struct InvalidSecp256k1HashFormat;

impl From<InvalidSecp256k1HashFormat> for CryptoError {
    fn from(_original: InvalidSecp256k1HashFormat) -> Self {
        CryptoError::VerificationError
    }
}

fn read_hash(data: &[u8]) -> Result<[u8; 32], InvalidSecp256k1HashFormat> {
    data.try_into().map_err(|_| InvalidSecp256k1HashFormat)
}

/// Error raised when signature is not 64 bytes long (32 bytes r, 32 bytes s)
struct InvalidSecp256k1SignatureFormat;

impl From<InvalidSecp256k1SignatureFormat> for CryptoError {
    fn from(_original: InvalidSecp256k1SignatureFormat) -> Self {
        CryptoError::VerificationError
    }
}

fn read_signature(data: &[u8]) -> Result<[u8; 64], InvalidSecp256k1SignatureFormat> {
    data.try_into().map_err(|_| InvalidSecp256k1SignatureFormat)
}

impl From<InvalidSecp256k1PubkeyFormat> for CryptoError {
    fn from(_original: InvalidSecp256k1PubkeyFormat) -> Self {
        CryptoError::VerificationError
    }
}
