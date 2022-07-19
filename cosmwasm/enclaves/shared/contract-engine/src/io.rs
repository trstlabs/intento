/// This contains all the user-facing functions. In these functions we will be using
/// the consensus_io_exchange_keypair and a user-generated key to create a symmetric key
/// that is unique to the user and the enclave
///
use super::types::{ContractMessage, IoNonce};
use enclave_cosmwasm_types::encoding::Binary;
//use enclave_cosmwasm_types::math::Uint128;
use enclave_cosmwasm_types::addresses::{CanonicalAddr, HumanAddr};
use enclave_cosmwasm_types::coins::Coin;
use enclave_cosmwasm_types::results::{
    CosmosMsg, Event, Reply, ReplyOn, Response, SubMsgResponse, SubMsgResult, WasmMsg,
    REPLY_ENCRYPTION_MAGIC_BYTES,
};
use enclave_ffi_types::EnclaveError;
use std::convert::TryInto;

use enclave_crypto::{AESKey, Ed25519PublicKey, Kdf, SIVEncryptable, KEY_MANAGER};

use log::*;
use serde::{Deserialize, Serialize};
use serde_json::json;
use serde_json::Value;
use sha2::Digest;

/// The internal_reply_enclave_sig is being passed with the reply (Only if the reply is wasm reply)
/// This is used by the receiver of the reply to:
/// a. Verify the sender (Cotnract address)
/// b. Authenticate the reply.
#[derive(Serialize, Deserialize, Clone, Debug, PartialEq)]
#[serde(untagged)]
enum WasmOutput {
    Err {
        #[serde(rename = "error")]
        err: Value,
        internal_msg_id: Option<Binary>,
        internal_reply_enclave_sig: Option<Binary>,
    },
    QueryOk {
        #[serde(rename = "ok")]
        ok: String,
    },
    Ok {
        #[serde(rename = "ok")]
        ok: Response,
        internal_reply_enclave_sig: Option<Binary>,
        internal_msg_id: Option<Binary>,
    },
}

pub fn calc_encryption_key(nonce: &IoNonce, user_public_key: &Ed25519PublicKey) -> AESKey {
    let enclave_io_key = KEY_MANAGER.get_consensus_io_exchange_keypair().unwrap();

    let tx_encryption_ikm = enclave_io_key.diffie_hellman(user_public_key);

    let tx_encryption_key = AESKey::new_from_slice(&tx_encryption_ikm).derive_key_from_this(nonce);

    trace!("rust tx_encryption_key {:?}", tx_encryption_key.get());

    tx_encryption_key
}

fn encrypt_serializable<T>(
    key: &AESKey,
    val: &T,
    reply_params: &Option<(Vec<u8>, u64)>,
) -> Result<String, EnclaveError>
where
    T: ?Sized + Serialize,
{
    let serialized: String = serde_json::to_string(val).map_err(|err| {
        debug!("got an error while trying to encrypt output error {}", err);
        EnclaveError::EncryptionError
    })?;

    let trimmed = serialized.trim_start_matches('"').trim_end_matches('"');

    encrypt_preserialized_string(key, trimmed, reply_params)
}

// use this to encrypt a vec value
fn encrypt_vec(
    key: &AESKey,
    val: Vec<u8>,
    reply_params: &Option<(Vec<u8>, u64)>,
) -> Result<Vec<u8>, EnclaveError> {
    let serialized = match reply_params {
        Some((reply_recipient_contract_hash, _)) => {
            let mut ser = vec![];
            ser.extend_from_slice(&reply_recipient_contract_hash);
            ser.extend_from_slice(&val);
            ser
        }
        None => val.to_vec(),
    };
    let encrypted_data = key.encrypt_siv(&serialized, None).map_err(|err| {
        debug!(
            "got an error while trying to encrypt binary output error {:?}: {}",
            err, err
        );
        EnclaveError::EncryptionError
    })?;

    Ok(encrypted_data.as_slice().to_vec())
}

// use this to encrypt a String that has already been serialized.  When that is the case, if
// encrypt_serializable is called instead, it will get double serialized, and any escaped
// characters will be double escaped
fn encrypt_preserialized_string(
    key: &AESKey,
    val: &str,
    reply_params: &Option<(Vec<u8>, u64)>,
) -> Result<String, EnclaveError> {
    let serialized = match reply_params {
        Some((reply_recipient_contract_hash, _)) => {
            let mut ser = vec![];
            ser.extend_from_slice(&reply_recipient_contract_hash);
            ser.extend_from_slice(val.as_bytes());
            ser
        }
        None => val.as_bytes().to_vec(),
    };
    let encrypted_data = key
        .encrypt_siv(serialized.as_slice(), None)
        .map_err(|err| {
            debug!(
                "got an error while trying to encrypt output error {:?}: {}",
                err, err
            );
            EnclaveError::EncryptionError
        })?;

    Ok(b64_encode(encrypted_data.as_slice()))
}

fn b64_encode(data: &[u8]) -> String {
    base64::encode(data)
}

pub fn encrypt_output(
    output: Vec<u8>,
    nonce: IoNonce,
    user_public_key: Ed25519PublicKey,
    contract_addr: &CanonicalAddr,
    contract_hash: &String,
    reply_params: Option<(Vec<u8>, u64)>,
    sender_addr: &CanonicalAddr,
) -> Result<Vec<u8>, EnclaveError> {
    // When encrypting an output we might encrypt an output that is a reply to a caller contract (Via "Reply" endpoint).
    // Therefore if reply_recipient_contract_hash is not "None" we append it to any encrypted data besided submessages that are irrelevant for replies.
    // More info in: https://github.com/CosmWasm/cosmwasm/blob/v1.0.0/packages/std/src/results/submessages.rs#L192-L198
    let encryption_key = calc_encryption_key(&nonce, &user_public_key);
    trace!(
        "Output before encryption: {:?}",
        String::from_utf8_lossy(&output)
    );

    let mut output: WasmOutput = serde_json::from_slice(&output).map_err(|err| {
        warn!("got an error while trying to deserialize output bytes into json");
        trace!("output: {:?} error: {:?}", output, err);
        EnclaveError::FailedToDeserialize
    })?;

    match &mut output {
        WasmOutput::Err {
            err,
            internal_reply_enclave_sig,
            internal_msg_id,
        } => {
            let encrypted_err = encrypt_serializable(&encryption_key, err, &reply_params)?;

            // Putting the error inside a 'generic_err' envelope, so we can encrypt the error itself
            *err = json!({"generic_err":{"msg":encrypted_err}});

            let msg_id = match reply_params {
                Some(ref r) => {
                    let encrypted_id = Binary::from_base64(&encrypt_preserialized_string(
                        &encryption_key,
                        &r.1.to_string(),
                        &reply_params,
                    )?)?;

                    Some(encrypted_id)
                }
                None => None,
            };

            *internal_msg_id = msg_id.clone();

            *internal_reply_enclave_sig = match reply_params {
                Some(_) => {
                    let reply = Reply {
                        id: msg_id.unwrap(),
                        result: SubMsgResult::Err(encrypted_err),
                    };
                    let reply_as_vec = serde_json::to_vec(&reply).map_err(|err| {
                        warn!(
                            "got an error while trying to serialize reply into bytes for internal_reply_enclave_sig  {:?}: {}",
                            reply, err
                        );
                        EnclaveError::FailedToSerialize
                    })?;
                    let tmp_contract_msg = ContractMessage {
                        nonce,
                        user_public_key,
                        msg: reply_as_vec.clone(),
                    };

                    Some(Binary::from(
                        create_callback_signature(sender_addr, &tmp_contract_msg, &[]).as_slice(),
                    ))
                }
                None => None, // Not a reply, we don't need enclave sig
            }
        }
        WasmOutput::QueryOk { ok } => {
            *ok = encrypt_serializable(&encryption_key, ok, &reply_params)?;
        }
        WasmOutput::Ok {
            ok,
            internal_reply_enclave_sig,
            internal_msg_id,
        } => {
            for sub_msg in &mut ok.messages {
                if let CosmosMsg::Wasm(wasm_msg) = &mut sub_msg.msg {
                    encrypt_wasm_msg(
                        wasm_msg,
                        &sub_msg.reply_on,
                        sub_msg.id,
                        nonce,
                        user_public_key,
                        contract_addr,
                        contract_hash,
                    )?;

                    // The ID can be extracted from the encrypted wasm msg
                    // We don't encrypt it here to remain with the same type (u64)
                    sub_msg.id = 0;
                }
            }

            // v1: The logs that will be emitted as part of a "wasm" event.
            for log in ok.attributes.iter_mut().filter(|log| log.encrypted) {
                log.key = encrypt_preserialized_string(&encryption_key, &log.key, &reply_params)?;
                log.value = encrypt_vec(&encryption_key, log.value.clone(), &reply_params)
                    .map_err(|err| {
                        debug!(
                            "got an error while trying to encrypt vec value {:?}: {}",
                            &log.value, err
                        );
                        EnclaveError::FailedToDeserialize
                    })?;
            }

            // v1: Extra, custom events separate from the main wasm one. These will have "wasm-"" prepended to the type.
            for event in ok.events.iter_mut() {
                for log in event.attributes.iter_mut().filter(|log| log.encrypted) {
                    log.key =
                        encrypt_preserialized_string(&encryption_key, &log.key, &reply_params)?;
                    log.value = encrypt_vec(&encryption_key, log.value.clone(), &reply_params)
                        .map_err(|err| {
                            debug!(
                                "got an error while trying to encrypt vec value {:?}: {}",
                                &log.value, err
                            );
                            EnclaveError::FailedToDeserialize
                        })?;
                }
            }

            if let Some(data) = &mut ok.data {
                *data = Binary::from_base64(&encrypt_serializable(
                    &encryption_key,
                    data,
                    &reply_params,
                )?)?;
            }

            let msg_id = match reply_params {
                Some(ref r) => {
                    let encrypted_id = Binary::from_base64(&encrypt_preserialized_string(
                        &encryption_key,
                        &r.1.to_string(),
                        &reply_params,
                    )?)?;

                    Some(encrypted_id)
                }
                None => None,
            };

            *internal_msg_id = msg_id.clone();

            *internal_reply_enclave_sig = match reply_params {
                Some(_) => {
                    let reply = Reply {
                        id: msg_id.unwrap(),
                        result: SubMsgResult::Ok(SubMsgResponse {
                            events: ok.events.clone(),
                            data: ok.data.clone(),
                        }),
                    };
                    let reply_as_vec = serde_json::to_vec(&reply).map_err(|err| {
                        warn!(
                            "got an error while trying to serialize reply into bytes for internal_reply_enclave_sig  {:?}: {}",
                            reply, err
                        );
                        EnclaveError::FailedToSerialize
                    })?;
                    let tmp_contract_msg = ContractMessage {
                        nonce,
                        user_public_key,
                        msg: reply_as_vec.clone(),
                    };

                    Some(Binary::from(
                        create_callback_signature(sender_addr, &tmp_contract_msg, &[]).as_slice(),
                    ))
                }
                None => None, // Not a reply, we don't need enclave sig
            }
        }
    };

    trace!("WasmOutput: {:?}", output);

    let encrypted_output = serde_json::to_vec(&output).map_err(|err| {
        debug!(
            "got an error while trying to serialize output json into bytes {:?}: {}",
            output, err
        );
        EnclaveError::FailedToSerialize
    })?;

    Ok(encrypted_output)
}

pub fn encrypt_msg(
    msg: &mut Binary,
    nonce: IoNonce,
    user_public_key: Ed25519PublicKey,
    send_as_addr: &CanonicalAddr,
    code_hash: String,
    funds: Vec<Coin>,
) -> Result<Vec<u8>, EnclaveError> {
    let mut hash_appended_msg = code_hash.as_bytes().to_vec();
    hash_appended_msg.extend_from_slice(msg.as_slice());

    let mut msg_to_pass = ContractMessage::from_base64(
        Binary(hash_appended_msg).to_base64(),
        nonce,
        user_public_key,
    )?;
    msg_to_pass.encrypt_in_place()?;
    let callback_sig_bytes = create_callback_signature(send_as_addr, &msg_to_pass, &funds);
    *msg = Binary::from(msg_to_pass.to_vec().as_slice());
    Ok(callback_sig_bytes)
}

fn encrypt_wasm_msg(
    wasm_msg: &mut WasmMsg,
    reply_on: &ReplyOn,
    msg_id: u64, // In every submessage there is a field called "id", currently used only by "reply".
    nonce: IoNonce,
    user_public_key: Ed25519PublicKey,
    contract_addr: &CanonicalAddr,
    reply_recipient_contract_hash: &String,
) -> Result<(), EnclaveError> {
    match wasm_msg {
        WasmMsg::Execute {
            msg,
            code_hash,
            callback_sig,
            funds,
            ..
        } => {
            let mut hash_appended_msg = code_hash.as_bytes().to_vec();
            if *reply_on != ReplyOn::Never {
                hash_appended_msg.extend_from_slice(REPLY_ENCRYPTION_MAGIC_BYTES);
                hash_appended_msg.extend_from_slice(&msg_id.to_be_bytes());
                hash_appended_msg.extend_from_slice(reply_recipient_contract_hash.as_bytes());
            };
            let mut msg_to_pass = ContractMessage::from_base64(
                Binary(hash_appended_msg).to_base64(),
                nonce,
                user_public_key,
            )?;

            msg_to_pass.encrypt_in_place()?;
            *msg = Binary::from(msg_to_pass.to_vec().as_slice());
            *callback_sig = Some(create_callback_signature(
                contract_addr,
                &msg_to_pass,
                &funds,
            ));
        }
        WasmMsg::Instantiate {
            msg,
            code_hash,
            auto_msg,
            callback_sig,
            funds,
            ..
        } => {
            let mut hash_appended_msg = code_hash.as_bytes().to_vec();
            if *reply_on != ReplyOn::Never {
                hash_appended_msg.extend_from_slice(REPLY_ENCRYPTION_MAGIC_BYTES);
                hash_appended_msg.extend_from_slice(&msg_id.to_be_bytes());
                hash_appended_msg.extend_from_slice(reply_recipient_contract_hash.as_bytes());
            };
            let mut hash_appended_auto_msg = code_hash.as_bytes().to_vec();
            let mut msg_to_pass = ContractMessage::from_base64(
                Binary(hash_appended_msg).to_base64(),
                nonce,
                user_public_key,
            )?;
            msg_to_pass.encrypt_in_place()?;
            *msg = Binary::from(msg_to_pass.to_vec().as_slice());
            *callback_sig = Some(create_callback_signature(
                contract_addr,
                &msg_to_pass,
                &funds,
            ));

            if auto_msg.is_some() {
                let auto_msg_unwrap = auto_msg.clone().unwrap();
                hash_appended_auto_msg.extend_from_slice(auto_msg_unwrap.as_slice());
                let mut auto_msg_to_pass = ContractMessage::from_base64(
                    Binary(hash_appended_auto_msg).to_base64(),
                    nonce,
                    user_public_key,
                )?;
                auto_msg_to_pass.encrypt_in_place()?;
                *auto_msg = Some(Binary::from(auto_msg_to_pass.to_vec().as_slice()));
            }
        }
    }

    Ok(())
}

pub fn create_callback_signature(
    sender_addr: &CanonicalAddr,
    msg_to_sign: &ContractMessage,
    funds_to_send: &[Coin],
) -> Vec<u8> {
    trace!(
        "callback sig sender_addr: {:?}",
        &HumanAddr::from_canonical(sender_addr).or(Err(EnclaveError::FailedToSerialize)),
    );
    trace!("callback sig msg_to_sign: {:?}", msg_to_sign);
    // Hash(Enclave_secret | sender(current contract) | msg_to_pass | funds)
    let mut callback_sig_bytes = KEY_MANAGER
        .get_consensus_callback_secret()
        .unwrap()
        .get()
        .to_vec();

    callback_sig_bytes.extend(sender_addr.as_slice());
    callback_sig_bytes.extend(msg_to_sign.msg.as_slice());
    callback_sig_bytes.extend(serde_json::to_vec(funds_to_send).unwrap());

    sha2::Sha256::digest(callback_sig_bytes.as_slice()).to_vec()
}

pub fn copy_into_array(slice: &[u8]) -> [u8; 32] {
    slice.try_into().expect("slice with incorrect length")
}
