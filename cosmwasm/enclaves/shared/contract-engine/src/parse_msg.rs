use cosmos_proto::tx::signing::SignMode;
use enclave_cosmos_types::types::{HandleType, SigInfo};
use enclave_cosmwasm_types::encoding::Binary;
use enclave_cosmwasm_types::ibc::IbcPacketReceiveMsg;

use enclave_crypto::HASH_SIZE;
use enclave_ffi_types::EnclaveError;
use log::*;
const HEX_ENCODED_HASH_SIZE: usize = HASH_SIZE * 2;
use super::contract_operations::{redact_custom_events, ParsedMessage};
use super::types::ContractMessage;
use enclave_cosmwasm_types::results::{DecryptedReply, Reply, SubMsgResponse, SubMsgResult};

// Parse the message that was passed to handle (Based on the assumption that it might be a reply or IBC as well)
pub fn parse_message(
    message: &[u8],
    sig_info: &SigInfo,
    handle_type: &HandleType,
) -> Result<ParsedMessage, EnclaveError> {
    return match handle_type {
        HandleType::HANDLE_TYPE_EXECUTE => match try_get_decrypted_msg(message) {
            Some(decrypted_contract_msg) => {
                trace!(
                    "execute input before decryption: {:?}",
                    base64::encode(&message)
                );

                Ok(ParsedMessage {
                    should_validate_sig_info: true,
                    was_msg_encrypted: true,
                    should_encrypt_output: true,
                    contract_msg: decrypted_contract_msg.contract_msg,
                    decrypted_msg: decrypted_contract_msg.decrypted_msg,
                    data_for_validation: None,
                })
            }
            None => {
                trace!(
                    "execute input was plaintext: {:?}",
                    base64::encode(&message)
                );

                let contract_msg = ContractMessage {
                    nonce: [0; 32],
                    user_public_key: [0; 32],
                    msg: message.into(),
                };

                let decrypted_msg = contract_msg.msg.clone();

                Ok(ParsedMessage {
                    should_validate_sig_info: true,
                    was_msg_encrypted: false,
                    should_encrypt_output: false,
                    contract_msg,
                    decrypted_msg,
                    data_for_validation: None,
                })
            }
        },

        HandleType::HANDLE_TYPE_REPLY => {
            let orig_contract_msg = ContractMessage::from_slice(message)?;
            let mut parsed_reply: Reply =
                serde_json::from_slice(&orig_contract_msg.msg).map_err(|err| {
                    warn!(
                    "reply got an error while trying to deserialize reply bytes into json {:?}: {}",
                    String::from_utf8_lossy(&orig_contract_msg.msg.clone()),
                    err
                );
                    EnclaveError::FailedToDeserialize
                })?;
            trace!("reply msg {:?}", parsed_reply);
            if !parsed_reply.is_encrypted {
                trace!(
                    "reply input is not encrypted: {:?}",
                    base64::encode(&message)
                );

                let msg_id =
                    String::from_utf8(parsed_reply.id.as_slice().to_vec()).map_err(|err| {
                        warn!(
                            "Failed to parse message id as string {:?}: {}",
                            parsed_reply.id.as_slice().to_vec(),
                            err
                        );
                        EnclaveError::FailedToDeserialize
                    })?;

                let msg_id_as_num = match msg_id.parse::<u64>() {
                    Ok(m) => m,
                    Err(err) => {
                        warn!("Failed to parse message id as number {}: {}", msg_id, err);
                        return Err(EnclaveError::FailedToDeserialize);
                    }
                };

                let decrypted_reply = DecryptedReply {
                    id: msg_id_as_num,
                    result: parsed_reply.result.clone(),
                };

                redact_custom_events(&mut parsed_reply);
                let serialized_encrypted_reply : Vec<u8> = serde_json::to_vec(&parsed_reply).map_err(|err| {
                        warn!(
                            "got an error while trying to serialize encrypted reply into bytes {:?}: {}",
                            parsed_reply, err
                        );
                        EnclaveError::FailedToSerialize
                    })?;

                let reply_contract_msg = ContractMessage {
                    nonce: orig_contract_msg.nonce,
                    user_public_key: orig_contract_msg.user_public_key,
                    msg: serialized_encrypted_reply,
                };

                let serialized_reply: Vec<u8> = serde_json::to_vec(&decrypted_reply).map_err(|err| {
                        warn!(
                            "got an error while trying to serialize decrypted reply into bytes {:?}: {}",
                            decrypted_reply, err
                        );
                        EnclaveError::FailedToSerialize
                    })?;

                return Ok(ParsedMessage {
                    should_validate_sig_info: false,
                    was_msg_encrypted: false,
                    should_encrypt_output: parsed_reply.was_orig_msg_encrypted,
                    contract_msg: reply_contract_msg,
                    decrypted_msg: serialized_reply,
                    data_for_validation: None,
                });
            }
            // Here we are sure the reply is OK because only OK is encrypted
            trace!(
                "reply input before decryption: {:?}",
                base64::encode(&message)
            );
            let mut parsed_encrypted_reply: Reply = serde_json::from_slice(
                &orig_contract_msg.msg.as_slice().to_vec(),
            )
            .map_err(|err| {
                warn!(
            "reply got an error while trying to deserialize msg input bytes into json {:?}: {}",
            String::from_utf8_lossy(&orig_contract_msg.msg),
            err
            );
                EnclaveError::FailedToDeserialize
            })?;

            match parsed_encrypted_reply.result.clone() {
                SubMsgResult::Ok(response) => {
                    let decrypted_msg_data = match response.data {
                        Some(data) => {
                            trace!(
                                "reply data before decryption: {:?}",
                                &data.as_slice().to_vec()
                            );
                            /*let tmp_contract_msg_data = ContractMessage {
                                nonce: orig_contract_msg.nonce,
                                user_public_key: orig_contract_msg.user_public_key,
                                msg: data.as_slice().to_vec(),
                            };*/
                            let data_msg = ContractMessage::from_slice(data.as_slice())?;
                            let decrypted_msg = data_msg.decrypt()?;
                            trace!(
                                "data_msg input afer decryption: {:?}",
                                String::from_utf8_lossy(&decrypted_msg)
                            );
                            trace!(
                                "data_msg binary afer decryption: {:?}",
                                Binary(decrypted_msg.clone())
                            );
                            /*  trace!(
                                "reply data after decryption: {:?}",
                                Binary(
                                    tmp_contract_msg_data.decrypt()?[HEX_ENCODED_HASH_SIZE..].to_vec(),
                                )
                            );*/
                            Some(Binary(decrypted_msg))
                        }
                        None => None,
                    };

                    let tmp_contract_msg_id = ContractMessage {
                        nonce: orig_contract_msg.nonce,
                        user_public_key: orig_contract_msg.user_public_key,
                        msg: parsed_encrypted_reply.id.as_slice().to_vec(),
                    };

                    let tmp_decrypted_msg_id = tmp_contract_msg_id.decrypt()?;

                    // Now we need to create synthetic ContractMessage to fit the API in "handle"
                    let result = SubMsgResult::Ok(SubMsgResponse {
                        events: response.events,
                        data: decrypted_msg_data,
                    });

                    let msg_id =
                        String::from_utf8(tmp_decrypted_msg_id[HEX_ENCODED_HASH_SIZE..].to_vec())
                            .map_err(|err| {
                            warn!(
                                "Failed to parse message id as string {:?}: {}",
                                tmp_decrypted_msg_id[HEX_ENCODED_HASH_SIZE..].to_vec(),
                                err
                            );
                            EnclaveError::FailedToDeserialize
                        })?;

                    let msg_id_as_num = match msg_id.parse::<u64>() {
                        Ok(m) => m,
                        Err(err) => {
                            warn!("Failed to parse message id as number {}: {}", msg_id, err);
                            return Err(EnclaveError::FailedToDeserialize);
                        }
                    };

                    let decrypted_reply = DecryptedReply {
                        id: msg_id_as_num,
                        result,
                    };

                    let decrypted_reply_as_vec =
                        serde_json::to_vec(&decrypted_reply).map_err(|err| {
                            warn!(
                                "got an error while trying to serialize reply into bytes {:?}: {}",
                                decrypted_reply, err
                            );
                            EnclaveError::FailedToSerialize
                        })?;
                    info!("reply {:?}", &parsed_encrypted_reply);
                    redact_custom_events(&mut parsed_encrypted_reply);
                    info!("redact_custom_events {:?}", &parsed_encrypted_reply);
                    //  let msg_for_sig = parse_data(&parsed_encrypted_reply);
                    // info!("msg_for_sig {}", msg_for_sig.clone());
                    let serialized_encrypted_reply : Vec<u8> = serde_json::to_vec(&parsed_encrypted_reply).map_err(|err| {
                    warn!(
                        "got an error while trying to serialize encrypted reply into bytes {:?}: {}",
                        parsed_encrypted_reply, err
                    );
                    EnclaveError::FailedToSerialize
                })?;

                    let reply_contract_msg = ContractMessage {
                        nonce: orig_contract_msg.nonce,
                        user_public_key: orig_contract_msg.user_public_key,
                        msg: serialized_encrypted_reply,
                    };

                    return Ok(ParsedMessage {
                        should_validate_sig_info: true,
                        was_msg_encrypted: true,
                        should_encrypt_output: true,
                        contract_msg: reply_contract_msg,
                        decrypted_msg: decrypted_reply_as_vec,
                        data_for_validation: Some(
                            tmp_decrypted_msg_id[..HEX_ENCODED_HASH_SIZE].to_vec(),
                        ),
                    });
                }
                SubMsgResult::Err(response) => {
                    let contract_msg = ContractMessage {
                        nonce: orig_contract_msg.nonce,
                        user_public_key: orig_contract_msg.user_public_key,
                        msg: base64::decode(response.clone()).map_err(|err| {
                            warn!(
                                "got an error while trying to serialize err reply from base64 {:?}: {}",
                                    response, err
                            );
                            EnclaveError::FailedToSerialize
                        })?
                    };

                    let decrypted_error = contract_msg.decrypt()?;

                    let tmp_contract_msg_id = ContractMessage {
                        nonce: orig_contract_msg.nonce,
                        user_public_key: orig_contract_msg.user_public_key,
                        msg: parsed_encrypted_reply.id.as_slice().to_vec(),
                    };

                    let tmp_decrypted_msg_id = tmp_contract_msg_id.decrypt()?;

                    // Now we need to create synthetic ContractMessage to fit the API in "handle"
                    let result = SubMsgResult::Err(
                        String::from_utf8(decrypted_error[HEX_ENCODED_HASH_SIZE..].to_vec())
                            .map_err(|err| {
                                warn!(
                                    "Failed to parse error as string {:?}: {}",
                                    decrypted_error[HEX_ENCODED_HASH_SIZE..].to_vec(),
                                    err
                                );
                                EnclaveError::FailedToDeserialize
                            })?,
                    );

                    let msg_id =
                        String::from_utf8(tmp_decrypted_msg_id[HEX_ENCODED_HASH_SIZE..].to_vec())
                            .map_err(|err| {
                            warn!(
                                "Failed to parse message id as string {:?}: {}",
                                tmp_decrypted_msg_id[HEX_ENCODED_HASH_SIZE..].to_vec(),
                                err
                            );
                            EnclaveError::FailedToDeserialize
                        })?;

                    let msg_id_as_num = match msg_id.parse::<u64>() {
                        Ok(m) => m,
                        Err(err) => {
                            warn!("Failed to parse message id as number {}: {}", msg_id, err);
                            return Err(EnclaveError::FailedToDeserialize);
                        }
                    };

                    let decrypted_reply = DecryptedReply {
                        id: msg_id_as_num,
                        result,
                    };

                    let decrypted_reply_as_vec =
                        serde_json::to_vec(&decrypted_reply).map_err(|err| {
                            warn!(
                                "got an error while trying to serialize reply into bytes {:?}: {}",
                                decrypted_reply, err
                            );
                            EnclaveError::FailedToSerialize
                        })?;

                    let serialized_encrypted_reply : Vec<u8> = serde_json::to_vec(&parsed_encrypted_reply).map_err(|err| {
                    warn!(
                        "got an error while trying to serialize encrypted reply into bytes {:?}: {}",
                        parsed_encrypted_reply, err
                    );
                    EnclaveError::FailedToSerialize
                })?;

                    let reply_contract_msg = ContractMessage {
                        nonce: orig_contract_msg.nonce,
                        user_public_key: orig_contract_msg.user_public_key,
                        msg: serialized_encrypted_reply,
                    };

                    return Ok(ParsedMessage {
                        should_validate_sig_info: true,
                        was_msg_encrypted: true,
                        should_encrypt_output: true,
                        contract_msg: reply_contract_msg,
                        decrypted_msg: decrypted_reply_as_vec,
                        data_for_validation: Some(
                            tmp_decrypted_msg_id[..HEX_ENCODED_HASH_SIZE].to_vec(),
                        ),
                    });
                }
            }
        }

        HandleType::HANDLE_TYPE_IBC_CHANNEL_OPEN
        | HandleType::HANDLE_TYPE_IBC_CHANNEL_CONNECT
        | HandleType::HANDLE_TYPE_IBC_CHANNEL_CLOSE
        | HandleType::HANDLE_TYPE_IBC_PACKET_ACK
        | HandleType::HANDLE_TYPE_IBC_PACKET_TIMEOUT => {
            trace!(
                "parsing {} msg (Should always be plaintext): {:?}",
                HandleType::to_export_name(&handle_type),
                base64::encode(&message)
            );

            let contract_msg = ContractMessage {
                nonce: [0; 32],
                user_public_key: [0; 32],
                msg: message.into(),
            };

            let decrypted_msg = contract_msg.msg.clone();

            Ok(ParsedMessage {
                should_validate_sig_info: false,
                was_msg_encrypted: false,
                should_encrypt_output: false,
                contract_msg: contract_msg,
                decrypted_msg,
                data_for_validation: None,
            })
        }
        HandleType::HANDLE_TYPE_IBC_PACKET_RECEIVE => {
            // TODO: Maybe mark whether the message was encrypted or not.
            let mut parsed_encrypted_ibc_packet: IbcPacketReceiveMsg =
                    serde_json::from_slice(&message.to_vec()).map_err(|err| {
                        warn!(
                "Got an error while trying to deserialize input bytes msg into IbcPacketReceiveMsg message {:?}: {}",
                String::from_utf8_lossy(&message),
                err
            );
                        EnclaveError::FailedToDeserialize
                    })?;

            let tmp_contract_data =
                parse_contract_msg(parsed_encrypted_ibc_packet.packet.data.as_slice());
            let mut was_msg_encrypted = false;
            let mut orig_msg = tmp_contract_data;

            match orig_msg.decrypt() {
                Ok(decrypted_msg) => {
                    // IBC packet was encrypted

                    trace!(
                        "ibc_packet_receive data before decryption: {:?}",
                        base64::encode(&message)
                    );

                    parsed_encrypted_ibc_packet.packet.data = decrypted_msg.as_slice().into();
                    was_msg_encrypted = true;
                }
                Err(_) => {
                    // assume data is not encrypted

                    trace!(
                        "ibc_packet_receive data was plaintext: {:?}",
                        base64::encode(&message)
                    );
                    orig_msg = ContractMessage {
                        nonce: [0; 32],
                        user_public_key: [0; 32],
                        msg: message.into(),
                    };
                }
            }
            Ok(ParsedMessage {
                    should_validate_sig_info: false,
                    was_msg_encrypted,
                    should_encrypt_output: was_msg_encrypted,
                    contract_msg: orig_msg,
                    decrypted_msg: serde_json::to_vec(&parsed_encrypted_ibc_packet).map_err(|err| {
                        warn!(
                            "got an error while trying to serialize IbcPacketReceive msg into bytes {:?}: {}",
                            parsed_encrypted_ibc_packet, err
                        );
                        EnclaveError::FailedToSerialize
                    })?,
                    data_for_validation: None,
                })
        }
    };
}
pub fn parse_contract_msg(message: &[u8]) -> ContractMessage {
    match ContractMessage::from_slice(message) {
        Ok(orig_msg) => orig_msg,
        Err(_) => {
            trace!(
                "Msg is not ContractMessage (probably plaintext): {:?}",
                base64::encode(&message)
            );

            ContractMessage {
                nonce: [0; 32],
                user_public_key: [0; 32],
                msg: message.into(),
            }
        }
    }
}

pub fn is_ibc_msg(handle_type: HandleType) -> bool {
    match handle_type {
        HandleType::HANDLE_TYPE_EXECUTE | HandleType::HANDLE_TYPE_REPLY => false,
        HandleType::HANDLE_TYPE_IBC_CHANNEL_OPEN
        | HandleType::HANDLE_TYPE_IBC_CHANNEL_CONNECT
        | HandleType::HANDLE_TYPE_IBC_CHANNEL_CLOSE
        | HandleType::HANDLE_TYPE_IBC_PACKET_RECEIVE
        | HandleType::HANDLE_TYPE_IBC_PACKET_ACK
        | HandleType::HANDLE_TYPE_IBC_PACKET_TIMEOUT => true,
    }
}

pub fn try_get_decrypted_msg(message: &[u8]) -> Option<DecryptedContractMessage> {
    let contract_msg = parse_contract_msg(message);
    match contract_msg.decrypt() {
        Ok(decrypted_msg) => Some(DecryptedContractMessage {
            contract_msg,
            decrypted_msg,
        }),
        Err(_) => None,
    }
}

pub struct DecryptedContractMessage {
    pub contract_msg: ContractMessage,
    pub decrypted_msg: Vec<u8>,
}
