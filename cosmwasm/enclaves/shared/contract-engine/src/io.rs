/// This contains all the user-facing functions. In these functions we will be using
/// the consensus_io_exchange_keypair and a user-generated key to create a symmetric key
/// that is unique to the user and the enclave
///
use super::types::{ContractMessage, IoNonce};
use crate::contract_validation::ReplyParams;
use enclave_cosmwasm_types::encoding::Binary;
use enclave_cosmwasm_types::math::Uint128;
use enclave_cosmwasm_types::addresses::{CanonicalAddr, HumanAddr};
use enclave_cosmwasm_types::coins::Coin;
use enclave_cosmwasm_types::results::{
    CosmosMsg, Reply, ReplyOn, Response, SubMsgResponse, SubMsgResult, WasmMsg,
    REPLY_ENCRYPTION_MAGIC_BYTES,Attribute,
};
use enclave_ffi_types::EnclaveError;
use std::convert::TryInto;

use enclave_crypto::{AESKey, Ed25519PublicKey, Kdf, SIVEncryptable, KEY_MANAGER};
use log::*;
use serde::{Deserialize, Serialize};
use serde_json::Value;
use serde_json::json;
use enclave_cosmwasm_types::ibc::{IbcBasicResponse, IbcReceiveResponse,IbcChannelOpenResponse};
use sha2::Digest;

/// The internal_reply_enclave_sig is being passed with the reply (Only if the reply is wasm reply)
/// This is used by the receiver of the reply to:
/// a. Verify the sender (Cotnract address)
/// b. Authenticate the reply.
#[derive(Serialize, Deserialize, Clone, Debug, PartialEq)]
#[serde(untagged)]
pub enum WasmOutput {
    Err {
        #[serde(rename = "Err")]
        err: Value,
        internal_msg_id: Option<Binary>,
        internal_reply_enclave_sig: Option<Binary>,
    },
    QueryOk {
       #[serde(rename = "Ok")]
        ok: String,
    },
    Ok {
        #[serde(rename = "Ok")]
        ok: Response,
        internal_reply_enclave_sig: Option<Binary>,
        internal_msg_id: Option<Binary>,
    },
    OkIBCBasic {
        #[serde(rename = "Ok")]
        ok: IbcBasicResponse,
    },
    OkIBCPacketReceive {
        #[serde(rename = "Ok")]
        ok: IbcReceiveResponse,
    },
    OkIBCOpenChannel {
        #[serde(rename = "Ok")]
        ok: IbcChannelOpenResponse,
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
    reply_params: &Option<Vec<ReplyParams>>,
    should_append_all_reply_params: bool,
) -> Result<String, EnclaveError>
where
    T: ?Sized + Serialize,
{
    let serialized: String = serde_json::to_string(val).map_err(|err| {
        debug!("got an error while trying to encrypt output error {}", err);
        EnclaveError::EncryptionError
    })?;

    let trimmed = serialized.trim_start_matches('"').trim_end_matches('"');

    encrypt_preserialized_string(key, trimmed, reply_params, should_append_all_reply_params)
}

// use this to encrypt a vec value
fn encrypt_vec(key: &AESKey, val: Vec<u8>) -> Result<Vec<u8>, EnclaveError> {
    let encrypted_data = key.encrypt_siv(&val, None).map_err(|err| {
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
    reply_params: &Option<Vec<ReplyParams>>,
    should_append_all_reply_params: bool,
) -> Result<String, EnclaveError> {
    let serialized = match reply_params {
        Some(v) => {
            let mut ser = vec![];
            ser.extend_from_slice(&v[0].recipient_contract_hash);
            if should_append_all_reply_params {
                for item in v.iter().skip(1) {
                    ser.extend_from_slice(REPLY_ENCRYPTION_MAGIC_BYTES);
                    ser.extend_from_slice(&item.sub_msg_id.to_be_bytes());
                    ser.extend_from_slice(item.recipient_contract_hash.as_slice());
                }
            }
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
    contract_msg: &ContractMessage,
    contract_addr: &CanonicalAddr,
    contract_hash: &str,
    reply_params: Option<Vec<ReplyParams>>,
    sender_addr: &CanonicalAddr,
) -> Result<Vec<u8>, EnclaveError> {
    // When encrypting an output we might encrypt an output that is a reply to a caller contract (Via "Reply" endpoint).
    // Therefore if reply_recipient_contract_hash is not "None" we append it to any encrypted data besided submessages that are irrelevant for replies.
    // More info in: https://github.com/CosmWasm/cosmwasm/blob/v1.0.0/packages/std/src/results/submessages.rs#L192-L198
    let encryption_key = calc_encryption_key(&contract_msg.nonce, &contract_msg.user_public_key);
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
            let encrypted_err = encrypt_serializable(&encryption_key, err, &reply_params, false)?;
            //trace!("output error: {:?}", encrypted_err);
            // Putting the error inside a 'generic_err' envelope, so we can encrypt the error itself
            *err = json!({"generic_err":{"msg":encrypted_err}});

            let msg_id = match reply_params {
                Some(ref r) => {
                    let encrypted_id = Binary::from_base64(&encrypt_preserialized_string(
                        &encryption_key,
                        &r[0].sub_msg_id.to_string(),
                        &reply_params,
                        true
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
                        was_orig_msg_encrypted: true,
                        is_encrypted: true,
                    };
                    let reply_as_vec = serde_json::to_vec(&reply).map_err(|err| {
                        warn!(
                            "got an error while trying to serialize reply into bytes for internal_reply_enclave_sig  {:?}: {}",
                            reply, err
                        );
                        EnclaveError::FailedToSerialize
                    })?;
                    let tmp_contract_msg = ContractMessage {
                        nonce: contract_msg.nonce,
                        user_public_key: contract_msg.user_public_key,
                        msg: reply_as_vec,
                    };

                    Some(Binary::from(
                        create_callback_signature(sender_addr, &tmp_contract_msg, &[]).as_slice(),
                    ))
                }
                None => None, // Not a reply, we don't need enclave sig
            }
        }
        WasmOutput::QueryOk { ok } => {
            *ok = encrypt_serializable(&encryption_key, ok, &reply_params, false)?;
        }
        WasmOutput::Ok {
            ok,
            internal_reply_enclave_sig,
            internal_msg_id,
        } => {
            for sub_msg in &mut ok.messages {
                trace!("submsg id: {:?}", sub_msg.id.clone());
                if let CosmosMsg::Wasm(wasm_msg) = &mut sub_msg.msg {
                    encrypt_wasm_msg(
                        wasm_msg,
                        &sub_msg.reply_on,
                        sub_msg.id,
                        contract_msg.nonce,
                        contract_msg.user_public_key,
                        contract_addr,
                        contract_hash,
                        &reply_params,
                    )?;

                    // The ID can be extracted from the encrypted wasm msg
                    // We don't encrypt it here to remain with the same type (u64)
                    sub_msg.id = 0;
                    trace!("encpypted submsg: {:?}", &wasm_msg);
                }
                sub_msg.was_msg_encrypted = true;
            }

            // v1: The logs that will be emitted as part of a "wasm" event.
            for log in ok.attributes.iter_mut().filter(|log| log.encrypted) {
                log.key = encrypt_preserialized_string(&encryption_key, &log.key, &None, false)?;
                log.value = encrypt_vec(&encryption_key, log.value.clone()).map_err(|err| {
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
                        encrypt_preserialized_string(&encryption_key, &log.key, &None, false)?;
                    log.value = encrypt_vec(&encryption_key, log.value.clone()).map_err(|err| {
                        debug!(
                            "got an error while trying to encrypt vec value {:?}: {}",
                            &log.value, err
                        );
                        EnclaveError::FailedToDeserialize
                    })?;
                }
            }

            if let Some(data) = &mut ok.data {
                trace!("reply data: {:?}", data);
                *data = Binary::from_base64(&encrypt_serializable(
                    &encryption_key,
                    data,
                    &reply_params,
                    false,
                )?)?;
            }

            let msg_id = match reply_params {
                Some(ref r) => {
                    let encrypted_id = Binary::from_base64(&encrypt_preserialized_string(
                        &encryption_key,
                        &r[0].sub_msg_id.to_string(),
                        &reply_params,
                        true,
                    )?)?;

                    Some(encrypted_id)
                }
                None => None,
            };

            *internal_msg_id = msg_id.clone();

            *internal_reply_enclave_sig = match reply_params {
                Some(_) => {
                  /*  let mut attributes: Vec<Attribute> = vec![];
                    for a in ok.attributes.clone(){
                        attributes.push(a.to_kv());
                    };
                     let events = match ok.attributes.len() {
                            0 => vec![],
                            _ => vec![Event {
                                ty: "wasm".to_string(),
                                attributes,
                            }],
                        };
                    let reply = Reply {
                        id: msg_id.unwrap(),
                        result: SubMsgResult::Ok(SubMsgResponse {
                            events,
                            data: ok.data.clone(),
                        }),
                    };*/
                    let reply = Reply {
                        id: msg_id.unwrap(),
                        result: SubMsgResult::Ok(SubMsgResponse {
                            events: vec![],
                            data: ok.data.clone(),
                        }),
                        was_orig_msg_encrypted: true,
                        is_encrypted: true,
                    };
                    trace!("reply : {:?}", reply);
                    let reply_as_vec = serde_json::to_vec(&reply).map_err(|err| {
                        warn!(
                            "got an error while trying to serialize reply into bytes for internal_reply_enclave_sig  {:?}: {}",
                            reply, err
                        );
                        EnclaveError::FailedToSerialize
                    })?;
                    trace!("reply_as_vec: {:?}", reply_as_vec);
                    let tmp_contract_msg = ContractMessage {
                        nonce: contract_msg.nonce,
                        user_public_key: contract_msg.user_public_key,
                        msg: reply_as_vec,
                    };
         
                    Some(Binary::from(
                        create_callback_signature(sender_addr, &tmp_contract_msg, &[]).as_slice(),
                    ))
                }
                None => None, // Not a reply, we don't need enclave sig
            }
        }
        WasmOutput::OkIBCPacketReceive { ok } => {
            for sub_msg in &mut ok.messages {
                if let CosmosMsg::Wasm(wasm_msg) = &mut sub_msg.msg {
                    match wasm_msg {
                        WasmMsg::Execute {
                            callback_sig,
                            msg,
                            funds,
                            ..
                        }
                        | WasmMsg::Instantiate {
                            callback_sig,
                            msg,
                            funds,
                            ..
                        } | WasmMsg::InstantiateAuto {
                            callback_sig,
                            msg,
                            funds,
                            ..
                        } => {
                            let msg_to_sign = ContractMessage {
                                nonce: [0; 32],
                                user_public_key: [0; 32],
                                msg: msg.as_slice().to_vec(),
                            };
                            *callback_sig = Some(create_callback_signature(
                                contract_addr,
                                &msg_to_sign,
                                &funds
                                    .iter()
                                    .map(|coin| Coin {
                                        denom: coin.denom.clone(),
                                        amount: Uint128::new(coin.amount.u128()),
                                    })
                                    .collect::<Vec<Coin>>()[..],
                            ));
                        }
                    }
                }
            }
        }
        WasmOutput::OkIBCOpenChannel { ok: _ } => { }
        WasmOutput::OkIBCBasic { ok:  _} =>  { 
        }
       /* WasmOutput::OkIBCPacketReceive { ok } =>  { 
        }*/
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

#[allow(clippy::too_many_arguments)]
fn encrypt_wasm_msg(
    wasm_msg: &mut WasmMsg,
    reply_on: &ReplyOn,
    msg_id: u64, // In every submessage there is a field called "id", currently used only by "reply".
    nonce: IoNonce,
    user_public_key: Ed25519PublicKey,
    contract_addr: &CanonicalAddr,
    reply_recipient_contract_hash: &str,
    reply_params: &Option<Vec<ReplyParams>>,
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
            }
            hash_appended_msg.extend_from_slice(msg.as_slice());

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
                funds,
            ));
        }
        WasmMsg::Instantiate {
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
            }

            if let Some(r) = reply_params {
                for param in r.iter() {
                    hash_appended_msg
                        .extend_from_slice(REPLY_ENCRYPTION_MAGIC_BYTES);
                    hash_appended_msg.extend_from_slice(&param.sub_msg_id.to_be_bytes());
                    hash_appended_msg.extend_from_slice(param.recipient_contract_hash.as_slice());
                }
            }
            
            hash_appended_msg.extend_from_slice(msg.as_slice());

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
                funds,
            ));
        }
        WasmMsg::InstantiateAuto {
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
            }

            if let Some(r) = reply_params {
                for param in r.iter() {
                    hash_appended_msg
                        .extend_from_slice(REPLY_ENCRYPTION_MAGIC_BYTES);
                    hash_appended_msg.extend_from_slice(&param.sub_msg_id.to_be_bytes());
                    hash_appended_msg.extend_from_slice(param.recipient_contract_hash.as_slice());
                }
            }
            
            hash_appended_msg.extend_from_slice(msg.as_slice());

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
                funds,
            ));


            if auto_msg.is_some() {
                let auto_msg_unwrap = auto_msg.clone().unwrap();
                let mut hash_appended_auto_msg = code_hash.as_bytes().to_vec(); 
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
    // Hash(Enclave_contract | sender(current contract) | msg_to_pass | funds)
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

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq)]
pub struct IBCOpenChannelOutput {
    #[serde(rename = "ok")]
    pub ok: Option<String>,
    #[serde(rename = "Err")]
    pub err: Option<Value>,
}


pub fn manipulate_callback_sig_for_plaintext(
    contract_addr: &CanonicalAddr,
    output: Vec<u8>,
) -> Result<WasmOutput, EnclaveError> {
    let mut raw_output: WasmOutput = serde_json::from_slice(&output).map_err(|err| {
        warn!("got an error while trying to deserialize output bytes into json");
        trace!("output: {:?} error: {:?}", output, err);
        EnclaveError::FailedToDeserialize
    })?;

    match &mut raw_output {
        WasmOutput::Ok { ok, .. } => {
            for sub_msg in &mut ok.messages {
                if let CosmosMsg::Wasm(wasm_msg) = &mut sub_msg.msg {
                    match wasm_msg {
                        WasmMsg::Execute {
                            callback_sig,
                            msg,
                            funds,
                            ..
                        }
                        | WasmMsg::Instantiate {
                            callback_sig,
                            msg,
                            funds,
                            ..
                        } | WasmMsg::InstantiateAuto {
                            callback_sig,
                            msg,
                            funds,
                            ..
                        } => {
                            let msg_to_sign = ContractMessage {
                                nonce: [0; 32],
                                user_public_key: [0; 32],
                                msg: msg.as_slice().to_vec(),
                            };
                            *callback_sig = Some(create_callback_signature(
                                contract_addr,
                                &msg_to_sign,
                                &funds
                                    .iter()
                                    .map(|coin| Coin {
                                        denom: coin.denom.clone(),
                                        amount: Uint128::new(coin.amount.u128()),
                                    })
                                    .collect::<Vec<Coin>>()[..],
                            ));
                        }
                    }
                }
            }
        }
        WasmOutput::OkIBCPacketReceive { ok } => {
            for sub_msg in &mut ok.messages {
                if let CosmosMsg::Wasm(wasm_msg) = &mut sub_msg.msg {
                    match wasm_msg {
                        WasmMsg::Execute {
                            callback_sig,
                            msg,
                            funds,
                            ..
                        }
                        | WasmMsg::Instantiate {
                            callback_sig,
                            msg,
                            funds,
                            ..
                        } | WasmMsg::InstantiateAuto {
                            callback_sig,
                            msg,
                            funds,
                            ..
                        } => {
                            let msg_to_sign = ContractMessage {
                                nonce: [0; 32],
                                user_public_key: [0; 32],
                                msg: msg.as_slice().to_vec(),
                            };
                            *callback_sig = Some(create_callback_signature(
                                contract_addr,
                                &msg_to_sign,
                                &funds
                                    .iter()
                                    .map(|coin| Coin {
                                        denom: coin.denom.clone(),
                                        amount: Uint128::new(coin.amount.u128()),
                                    })
                                    .collect::<Vec<Coin>>()[..],
                            ));
                            
                        }
                    }
                }
            }
        }
        _ => {}
    }

    Ok(raw_output)
}

pub fn set_attributes_to_plaintext(attributes: &mut Vec<Attribute>) {
    for attr in attributes {
        attr.encrypted = false;
    }
}

pub fn set_all_logs_to_plaintext(raw_output: &mut WasmOutput) {
    match raw_output {
        WasmOutput::Ok { ok, .. } => {
            set_attributes_to_plaintext(&mut ok.attributes);
            for ev in &mut ok.events {
                set_attributes_to_plaintext(&mut ev.attributes);
            }
        }
        WasmOutput::OkIBCPacketReceive { ok } => {
            set_attributes_to_plaintext(&mut ok.attributes);
            for ev in &mut ok.events {
                set_attributes_to_plaintext(&mut ev.attributes);
            }
        }
        _ => {}
    }
}