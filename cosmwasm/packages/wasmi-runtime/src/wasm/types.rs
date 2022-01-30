use log::*;

use enclave_ffi_types::EnclaveError;
use protobuf::Message;
use serde::{Deserialize, Serialize};

use crate::cosmwasm::{
    encoding::Binary,
    types::{CanonicalAddr, HumanAddr},
};
use crate::crypto::{
    multisig::MultisigThresholdPubKey, secp256k1::Secp256k1PubKey, traits::PubKey, AESKey,
    CryptoError, Ed25519PublicKey, SIVEncryptable, HASH_SIZE,
};
use crate::proto;

use super::io::calc_encryption_key;
use crate::cosmwasm::coins::Coin;
use crate::cosmwasm::math::Uint128;

use crate::wasm::contract_validation::calc_contract_hash;

pub type IoNonce = [u8; 32];

pub struct ContractCode<'code> {
    code: &'code [u8],
    hash: [u8; HASH_SIZE],
}

impl<'code> ContractCode<'code> {
    pub fn new(code: &'code [u8]) -> Self {
        let hash = calc_contract_hash(code);
        Self { code, hash }
    }

    pub fn code(&self) -> &[u8] {
        self.code
    }

    pub fn hash(&self) -> [u8; HASH_SIZE] {
        self.hash
    }
}

#[derive(Serialize, Deserialize, PartialEq, Debug)]
pub struct SecretMessage {
    pub nonce: IoNonce,
    pub user_public_key: Ed25519PublicKey,
    pub msg: Vec<u8>,
}

impl SecretMessage {
    pub fn encrypt_in_place(&mut self) -> Result<(), EnclaveError> {
        self.msg = self
            .encryption_key()
            .encrypt_siv(self.msg.as_slice(), None)
            .map_err(|err| {
                error!("got an error while trying to encrypt the msg: {:?}", err);
                EnclaveError::EncryptionError
            })?;

        Ok(())
    }

    pub fn decrypt(&self) -> Result<Vec<u8>, EnclaveError> {
        let key = self.encryption_key();

        // pass
        let msg = key.decrypt_siv(self.msg.as_slice(), None).map_err(|err| {
            error!("got an error while trying to decrypt the msg: {:?}", err);
            EnclaveError::DecryptionError
        })?;

        Ok(msg)
    }

    pub fn encryption_key(&self) -> AESKey {
        calc_encryption_key(&self.nonce, &self.user_public_key)
    }

    pub fn from_base64(
        msg_b64: String,
        nonce: IoNonce,
        user_public_key: Ed25519PublicKey,
    ) -> Result<Self, EnclaveError> {
        let msg = base64::decode(&msg_b64.to_owned().into_bytes()).map_err(|err| {
            error!(
                "got an error while trying to decode msg to next contract as base64 {:?}: {:?}",
                msg_b64, err
            );
            EnclaveError::FailedToDeserialize
        })?;

        Ok(SecretMessage {
            msg,
            nonce,
            user_public_key,
        })
    }

    pub fn from_slice(msg: &[u8]) -> Result<Self, EnclaveError> {
        // 32 bytes of nonce
        // 32 bytes of 25519 compressed public key
        // 16+ bytes of encrypted data

        if msg.len() < 82 {
            error!(
                "Encrypted message length {:?} is too short. Cannot parse",
                msg.len()
            );
            return Err(EnclaveError::DecryptionError);
        };

        let mut nonce = [0u8; 32];
        nonce.copy_from_slice(&msg[0..32]);

        let mut user_pubkey = [0u8; 32];
        user_pubkey.copy_from_slice(&msg[32..64]);

        debug!(
            "SecretMessage::from_slice nonce = {:?} pubkey = {:?}",
            nonce,
            hex::encode(user_pubkey)
        );

        Ok(SecretMessage {
            nonce,
            user_public_key: user_pubkey,
            msg: msg[64..].to_vec(),
        })
    }

    pub fn to_vec(&self) -> Vec<u8> {
        let mut packed_msg: Vec<u8> = self.nonce.to_vec();
        packed_msg.extend_from_slice(&self.user_public_key);
        packed_msg.extend_from_slice(self.msg.as_slice());
        packed_msg
    }
}

#[derive(PartialEq, Clone, Debug)]
pub enum CosmosPubKey {
    Secp256k1(Secp256k1PubKey),
    Multisig(MultisigThresholdPubKey),
}

/// `"/"` + `proto::crypto::multisig::LegacyAminoPubKey::descriptor_static().full_name()`
const TYPE_URL_MULTISIG_LEGACY_AMINO_PUBKEY: &str = "/cosmos.crypto.multisig.LegacyAminoPubKey";
/// `"/"` + `proto::crypto::secp256k1::PubKey::descriptor_static().full_name()`
const TYPE_URL_SECP256K1_PUBKEY: &str = "/cosmos.crypto.secp256k1.PubKey";

impl CosmosPubKey {
    pub fn from_proto(public_key: &protobuf::well_known_types::Any) -> Result<Self, CryptoError> {
        let public_key_parser = match public_key.type_url.as_str() {
            TYPE_URL_SECP256K1_PUBKEY => Self::secp256k1_from_proto,
            TYPE_URL_MULTISIG_LEGACY_AMINO_PUBKEY => Self::multisig_legacy_amino_from_proto,
            _ => {
                warn!("found public key of unsupported type: {:?}", public_key);
                return Err(CryptoError::ParsingError);
            }
        };

        public_key_parser(&public_key.value)
    }

    fn secp256k1_from_proto(public_key_bytes: &[u8]) -> Result<Self, CryptoError> {
        use proto::crypto::secp256k1::PubKey;
        let pub_key = PubKey::parse_from_bytes(public_key_bytes).map_err(|_err| {
            warn!(
                "Could not parse secp256k1 public key from these bytes: {}",
                Binary(public_key_bytes.to_vec())
            );
            CryptoError::ParsingError
        })?;
        Ok(CosmosPubKey::Secp256k1(Secp256k1PubKey::new(pub_key.key)))
    }

    fn multisig_legacy_amino_from_proto(public_key_bytes: &[u8]) -> Result<Self, CryptoError> {
        use proto::crypto::multisig::LegacyAminoPubKey;
        let multisig_key =
            LegacyAminoPubKey::parse_from_bytes(public_key_bytes).map_err(|_err| {
                warn!(
                    "Could not parse multisig public key from these bytes: {}",
                    Binary(public_key_bytes.to_vec())
                );
                CryptoError::ParsingError
            })?;
        let mut pubkeys = vec![];
        for public_key in &multisig_key.public_keys {
            pubkeys.push(CosmosPubKey::from_proto(public_key)?);
        }
        Ok(CosmosPubKey::Multisig(MultisigThresholdPubKey::new(
            multisig_key.threshold,
            pubkeys,
        )))
    }
}

impl PubKey for CosmosPubKey {
    fn get_address(&self) -> CanonicalAddr {
        match self {
            CosmosPubKey::Secp256k1(pubkey) => pubkey.get_address(),
            CosmosPubKey::Multisig(pubkey) => pubkey.get_address(),
        }
    }

    fn amino_bytes(&self) -> Vec<u8> {
        match self {
            CosmosPubKey::Secp256k1(pubkey) => pubkey.amino_bytes(),
            CosmosPubKey::Multisig(pubkey) => pubkey.amino_bytes(),
        }
    }

    fn verify_bytes(&self, bytes: &[u8], sig: &[u8]) -> Result<(), CryptoError> {
        match self {
            CosmosPubKey::Secp256k1(pubkey) => pubkey.verify_bytes(bytes, sig),
            CosmosPubKey::Multisig(pubkey) => pubkey.verify_bytes(bytes, sig),
        }
    }
}

// This type is a copy of the `proto::tx::signing::SignMode` allowing us
// to create a Deserialize impl for it without touching the autogenerated type.
// See: https://serde.rs/remote-derive.html
#[allow(non_camel_case_types)]
#[derive(Deserialize)]
#[serde(remote = "proto::tx::signing::SignMode")]
pub enum SignModeDef {
    SIGN_MODE_UNSPECIFIED = 0,
    SIGN_MODE_DIRECT = 1,
    SIGN_MODE_TEXTUAL = 2,
    SIGN_MODE_LEGACY_AMINO_JSON = 127,
}

// This is called `VerificationInfo` on the Go side
#[derive(Deserialize, Clone, Debug, PartialEq)]
pub struct SigInfo {
    pub sign_bytes: Binary,
    #[serde(with = "SignModeDef")]
    pub sign_mode: proto::tx::signing::SignMode,
    pub mode_info: Binary,
    pub public_key: Binary,
    pub signature: Binary,
    pub callback_sig: Option<Binary>,
}

// Should be in sync with https://github.com/cosmos/cosmos-sdk/blob/v0.38.3/x/auth/types/stdtx.go#L216
#[derive(Deserialize, Clone, Default, Debug, PartialEq)]
pub struct StdSignDoc {
    pub account_number: String,
    pub chain_id: String,
    pub memo: String,
    pub msgs: Vec<StdCosmWasmMsg>,
    pub sequence: String,
}

#[derive(Debug)]
pub struct SignDoc {
    pub body: TxBody,
    pub auth_info: AuthInfo,
    pub chain_id: String,
    pub account_number: u64,
}

impl SignDoc {
    pub fn from_bytes(bytes: &[u8]) -> Result<Self, EnclaveError> {
        let raw_sign_doc = proto::tx::tx::SignDoc::parse_from_bytes(bytes).map_err(|err| {
            warn!(
                "got an error while trying to deserialize sign doc bytes from protobuf: {}: {}",
                err,
                Binary(bytes.into()),
            );
            EnclaveError::FailedToDeserialize
        })?;

        let body = TxBody::from_bytes(&raw_sign_doc.body_bytes)?;
        let auth_info = AuthInfo::from_bytes(&raw_sign_doc.auth_info_bytes)?;

        Ok(Self {
            body,
            auth_info,
            chain_id: raw_sign_doc.chain_id,
            account_number: raw_sign_doc.account_number,
        })
    }
}

#[derive(Debug)]
pub struct TxBody {
    pub messages: Vec<CosmWasmMsg>,
    // Leaving this here for discoverability. We can use this, but don't verify it today.
    memo: (),
    timeout_height: (),
}

impl TxBody {
    pub fn from_bytes(bytes: &[u8]) -> Result<Self, EnclaveError> {
        let tx_body = proto::tx::tx::TxBody::parse_from_bytes(bytes).map_err(|err| {
            warn!(
                "got an error while trying to deserialize cosmos message body bytes from protobuf: {}: {}",
                err,
                Binary(bytes.into()),
            );
            EnclaveError::FailedToDeserialize
        })?;

        let messages = tx_body
            .messages
            .into_iter()
            .map(|any| CosmWasmMsg::from_bytes(&any.value))
            .collect::<Result<Vec<_>, _>>()?;

        Ok(TxBody {
            messages,
            memo: (),
            timeout_height: (),
        })
    }
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq)]
#[serde(rename_all = "snake_case", tag = "type", content = "value")]
pub enum StdCosmWasmMsg {
    #[serde(alias = "compute/MsgExecuteContract")]
    Execute {
        sender: HumanAddr,
        contract: HumanAddr,
        /// msg is the json-encoded HandleMsg struct (as raw Binary)
        msg: String,
        sent_funds: Vec<Coin>,
        callback_sig: Option<Vec<u8>>,
    },
    #[serde(alias = "compute/MsgInstantiateContract")]
    Instantiate {
        sender: HumanAddr,
        code_id: String,
        init_msg: String,
        auto_msg: String,
        init_funds: Vec<Coin>,
        contract_id: String,
        callback_sig: Option<Vec<u8>>,
    },
    #[serde(alias = "trstlabs.trst.trst.MsgCreateItem")]
    CreateItem {
    // message fields
     creator: HumanAddr,
    title: String,
    description: String,
    shipping_cost: i64,
    local_pickup: String,
    estimation_count: i64,
    tags: Vec<String>,// ::protobuf::RepeatedField<::std::string::String>,//std::string::String,//Vec<String>,
    condition: i64,
    shipping_region: Vec<String>,// ::protobuf::RepeatedField<::std::string::String>,//std::string::String,//Vec<String>,
    deposit_amount: i64,
    init_msg: String,
    auto_msg: String,
    photos: Vec<String>,
    token_uri: String,// ::protobuf::RepeatedField<::std::string::String>,//std::string::String,//Vec<String>,
},
#[serde(alias = "trstlabs.trst.trst.MsgCreateEstimation")]
CreateEstimation {
// message fields
estimator: HumanAddr,
estimate_msg: String,
itemid: u64,
deposit: i64,
interested: bool,
},

#[serde(alias = "trstlabs.trst.trst.MsgFlagItem")]
FlagItem {
// message fields
estimator: HumanAddr,
itemid: u64,
flag_msg: String,
},

#[serde(alias = "trstlabs.trst.trst.MsgItemTransferable")]
ItemTransferable {
// message fields
seller: HumanAddr,
transferable_msg: String,
itemid: u64,

},

#[serde(alias = "trstlabs.trst.trst.MsgRevealEstimation")]
RevealEstimation {
// message fields
creator: HumanAddr,
itemid: u64,
reveal_msg: String,
}

}

impl StdCosmWasmMsg {
    pub(crate) fn into_cosmwasm_msg(self) -> Result<CosmWasmMsg, EnclaveError> {
        match self {
            Self::Execute {
                sender,
                contract,
                msg,
                sent_funds,
                callback_sig,
            } => {
                let sender = CanonicalAddr::from_human(&sender).map_err(|err| {
                    warn!("failed to turn human addr to canonical addr when parsing CosmWasmMsg: {:?}", err);
                    EnclaveError::FailedToDeserialize
                })?;
                let msg = Binary::from_base64(&msg).map_err(|err| {
                    warn!(
                        "failed to parse base64 msg when parsing CosmWasmMsg: {:?}",
                        err
                    );
                    EnclaveError::FailedToDeserialize
                })?;
                let msg = msg.0;
                Ok(CosmWasmMsg::Execute {
                    sender,
                    contract,
                    msg,
                    sent_funds,
                    callback_sig,
                })
            }
            Self::Instantiate {
                sender,
                init_msg,
                auto_msg,
                init_funds,
                contract_id,
                callback_sig,
                code_id: _,
            } => {
                let sender = CanonicalAddr::from_human(&sender).map_err(|err| {
                    warn!("failed to turn human addr to canonical addr when parsing CosmWasmMsg: {:?}", err);
                    EnclaveError::FailedToDeserialize
                })?;
                let init_msg = Binary::from_base64(&init_msg).map_err(|err| {
                    warn!(
                        "failed to parse base64 init_msg when parsing CosmWasmMsg: {:?}",
                        err
                    );
                    EnclaveError::FailedToDeserialize
                })?;
                let auto_msg = Binary::from_base64(&auto_msg).map_err(|err| {
                    warn!(
                        "failed to parse base64 auto_msg when parsing last CosmWasmMsg: {:?}",
                        err
                    );
                    EnclaveError::FailedToDeserialize
                })?;
                let init_msg = init_msg.0;
                let auto_msg = auto_msg.0;
                Ok(CosmWasmMsg::Instantiate {
                    sender,
                    init_msg,
                    auto_msg, 
                    init_funds,
                    contract_id,
                    callback_sig,
                })
            }
            Self::CreateItem {
                creator,
                title,
                description,
                shipping_cost,
                local_pickup,
                estimation_count,
                tags,
                condition,
                shipping_region,
                deposit_amount,
                init_msg,
                auto_msg,
                photos,
                token_uri
            } => {
               let creator = CanonicalAddr::from_human(&creator).map_err(|err| {
                    warn!("failed to turn human addr to canonical addr when parsing CosmWasmMsg: {:?}", err);
                    EnclaveError::FailedToDeserialize
                })?;

                let init_msg = Binary::from_base64(&init_msg).map_err(|err| {
                    warn!(
                        "failed to parse base64 msg when parsing CosmWasmMsg: {:?}",
                        err
                    );
                    EnclaveError::FailedToDeserialize
                })?;
                let init_msg = init_msg.0;

                let auto_msg = Binary::from_base64(&auto_msg).map_err(|err| {
                    warn!(
                        "failed to parse base64 msg when parsing CosmWasmMsg: {:?}",
                        err
                    );
                    EnclaveError::FailedToDeserialize
                })?;
                let auto_msg = auto_msg.0;
            
                Ok(CosmWasmMsg::CreateItem {
                creator,
                title,
                description,
                shipping_cost,
                local_pickup,
                estimation_count,
                tags,
                condition,
                shipping_region,
                deposit_amount,
                init_msg,
                auto_msg,
                photos,
                token_uri
                })
            }
            Self::FlagItem {
                estimator,
                itemid,
                flag_msg,
            } => {
               let estimator = CanonicalAddr::from_human(&estimator).map_err(|err| {
                    warn!("failed to turn human addr to canonical addr when parsing CosmWasmMsg: {:?}", err);
                    EnclaveError::FailedToDeserialize
                })?;

                let flag_msg = Binary::from_base64(&flag_msg).map_err(|err| {
                    warn!(
                        "failed to parse base64 msg when parsing CosmWasmMsg: {:?}",
                        err
                    );
                    EnclaveError::FailedToDeserialize
                })?;
                let flag_msg = flag_msg.0;
            
                Ok(CosmWasmMsg::FlagItem {
                    estimator,
                    itemid,
                    flag_msg,
                })
            },
            Self::ItemTransferable {
                seller,
                transferable_msg,
                itemid,
             
            } => {
               let seller = CanonicalAddr::from_human(&seller).map_err(|err| {
                    warn!("failed to turn human addr to canonical addr when parsing CosmWasmMsg: {:?}", err);
                    EnclaveError::FailedToDeserialize
                })?;

                let transferable_msg = Binary::from_base64(&transferable_msg).map_err(|err| {
                    warn!(
                        "failed to parse base64 msg when parsing CosmWasmMsg: {:?}",
                        err
                    );
                    EnclaveError::FailedToDeserialize
                })?;
                let transferable_msg = transferable_msg.0;
            
                Ok(CosmWasmMsg::ItemTransferable {
                    seller,
                    itemid,
                    transferable_msg,
                })
            }
             Self::CreateItem {
                creator,
                title,
                description,
                shipping_cost,
                local_pickup,
                estimation_count,
                tags,
                condition,
                shipping_region,
                deposit_amount,
                init_msg,
                auto_msg,
                photos,
                token_uri,
            } => {
               let creator = CanonicalAddr::from_human(&creator).map_err(|err| {
                    warn!("failed to turn human addr to canonical addr when parsing CosmWasmMsg: {:?}", err);
                    EnclaveError::FailedToDeserialize
                })?;

                let init_msg = Binary::from_base64(&init_msg).map_err(|err| {
                    warn!(
                        "failed to parse base64 msg when parsing CosmWasmMsg: {:?}",
                        err
                    );
                    EnclaveError::FailedToDeserialize
                })?;
                let init_msg = init_msg.0;
            
                let auto_msg = Binary::from_base64(&auto_msg).map_err(|err| {
                    warn!(
                        "failed to parse base64 msg when parsing CosmWasmMsg: {:?}",
                        err
                    );
                    EnclaveError::FailedToDeserialize
                })?;
                let auto_msg = auto_msg.0;
            

                Ok(CosmWasmMsg::CreateItem {
                creator,
                title,
                description,
                shipping_cost,
                local_pickup,
                estimation_count,
                tags,
                condition,
                shipping_region,
                deposit_amount,
                init_msg,
                auto_msg,
                photos,
                token_uri
                })
            }
            Self::RevealEstimation {
                creator,
                itemid,
                reveal_msg,
            } => {
               let creator = CanonicalAddr::from_human(&creator).map_err(|err| {
                    warn!("failed to turn human addr to canonical addr when parsing CosmWasmMsg: {:?}", err);
                    EnclaveError::FailedToDeserialize
                })?;

                let reveal_msg = Binary::from_base64(&reveal_msg).map_err(|err| {
                    warn!(
                        "failed to parse base64 msg when parsing CosmWasmMsg: {:?}",
                        err
                    );
                    EnclaveError::FailedToDeserialize
                })?;
                let reveal_msg = reveal_msg.0;
            
                Ok(CosmWasmMsg::RevealEstimation {
                    creator,
                    itemid,
                    reveal_msg,
                })
            }
            Self::FlagItem {
                estimator,
                itemid,
                flag_msg,
            } => {
               let estimator = CanonicalAddr::from_human(&estimator).map_err(|err| {
                    warn!("failed to turn human addr to canonical addr when parsing CosmWasmMsg: {:?}", err);
                    EnclaveError::FailedToDeserialize
                })?;

                let flag_msg = Binary::from_base64(&flag_msg).map_err(|err| {
                    warn!(
                        "failed to parse base64 msg when parsing CosmWasmMsg: {:?}",
                        err
                    );
                    EnclaveError::FailedToDeserialize
                })?;
                let flag_msg = flag_msg.0;
            
                Ok(CosmWasmMsg::FlagItem {
                    estimator,
                    itemid,
                    flag_msg,
                })
            }
            Self::ItemTransferable {
                seller,
                transferable_msg,
                itemid,
          
            } => {
               let seller = CanonicalAddr::from_human(&seller).map_err(|err| {
                    warn!("failed to turn human addr to canonical addr when parsing CosmWasmMsg: {:?}", err);
                    EnclaveError::FailedToDeserialize
                })?;

                let transferable_msg = Binary::from_base64(&transferable_msg).map_err(|err| {
                    warn!(
                        "failed to parse base64 msg when parsing CosmWasmMsg: {:?}",
                        err
                    );
                    EnclaveError::FailedToDeserialize
                })?;
                let transferable_msg = transferable_msg.0;
            
                Ok(CosmWasmMsg::ItemTransferable {
                    seller,
                    transferable_msg,
                    itemid,
                })
            }
            Self::CreateEstimation {
                estimator,
                estimate_msg,
                itemid,
                deposit,
                interested,
            } => {
               let estimator = CanonicalAddr::from_human(&estimator).map_err(|err| {
                    warn!("failed to turn human addr to canonical addr when parsing CosmWasmMsg: {:?}", err);
                    EnclaveError::FailedToDeserialize
                })?;

                let estimate_msg = Binary::from_base64(&estimate_msg).map_err(|err| {
                    warn!(
                        "failed to parse base64 msg when parsing CosmWasmMsg: {:?}",
                        err
                    );
                    EnclaveError::FailedToDeserialize
                })?;
                let estimate_msg = estimate_msg.0;
            
                Ok(CosmWasmMsg::CreateEstimation {
                    estimator,
                    estimate_msg,
                    itemid,
                    deposit,
                    interested,
                   
                })
            }
        }
    }
}


#[derive(Debug)]
pub enum CosmWasmMsg {
    Execute {
        sender: CanonicalAddr,
        contract: HumanAddr,
        msg: Vec<u8>,
        sent_funds: Vec<Coin>,
        callback_sig: Option<Vec<u8>>,
    },
    Instantiate {
        sender: CanonicalAddr,
        init_msg: Vec<u8>,
        auto_msg: Vec<u8>,
        init_funds: Vec<Coin>,
        contract_id: String,
        callback_sig: Option<Vec<u8>>,
    },
    CreateItem {
        creator: CanonicalAddr,
        title: String,
        description:String,
        shipping_cost: i64,
        local_pickup:String,
        estimation_count: i64,
        tags: Vec<String>,//V ::protobuf::RepeatedField<::std::string::String>,//std::string::String,
        condition: i64,
        shipping_region: Vec<String>,//V ::protobuf::RepeatedField<::std::string::String>,//std::string::String,
        deposit_amount: i64,
        init_msg: Vec<u8>,
        auto_msg: Vec<u8>,
        photos:  Vec<String>,
        token_uri:  String,//V ::protobuf::RepeatedField<::std::string::String>,//std::string::String,
    },
    CreateEstimation {
        estimator: CanonicalAddr,
        estimate_msg:  Vec<u8>,
        itemid: u64,
        deposit: i64,
        interested:bool,
    },
    FlagItem {
        estimator: CanonicalAddr,
        itemid: u64,
        flag_msg:  Vec<u8>,
    },
    ItemTransferable {
        seller: CanonicalAddr,
        transferable_msg:  Vec<u8>,
        itemid: u64,
    },
    RevealEstimation {
        creator: CanonicalAddr,
        itemid: u64,
        reveal_msg:  Vec<u8>,
    },
    
    Other,
}

impl CosmWasmMsg {
    pub fn from_bytes(bytes: &[u8]) -> Result<Self, EnclaveError> {
        Self::try_parse_reveal_estimation(bytes)
            .or_else(|_| Self::try_parse_instantiate(bytes)).or_else(|_| Self::try_parse_create_item(bytes)).or_else(|_| Self::try_parse_create_estimation(bytes)).or_else(|_| Self::try_parse_execute(bytes)).or_else(|_| Self::try_parse_flag_item(bytes))
            .or_else(|_| {
                warn!(
                    "got an error while trying to deserialize cosmwasm message bytes from protobuf: {}",
                    Binary(bytes.into())
                );
                Ok(CosmWasmMsg::Other)
            })
    }

    fn try_parse_instantiate(bytes: &[u8]) -> Result<Self, EnclaveError> {
        use proto::cosmwasm::msg::MsgInstantiateContract;

        let raw_msg = MsgInstantiateContract::parse_from_bytes(bytes)
            .map_err(|_| EnclaveError::FailedToDeserialize)?;

        trace!(
            "try_parse_instantiate sender: len={} val={:?}",
            raw_msg.sender.len(),
            raw_msg.sender
        );
        let sender = CanonicalAddr::from_human(&HumanAddr(raw_msg.sender)).map_err(|err| {
            warn!("failed to turn human addr to canonical addr when try_parse_create_item CosmWasmMsg: {:?}", err);
            EnclaveError::FailedToDeserialize
        })?;

        let init_funds = Self::parse_funds(raw_msg.init_funds)?;

        let callback_sig = Some(raw_msg.callback_sig);

        Ok(CosmWasmMsg::Instantiate {
            sender: sender,
            init_msg: raw_msg.init_msg,
            auto_msg: raw_msg.auto_msg,
            init_funds,
            contract_id: raw_msg.contract_id,
            callback_sig,
        })
    }

    fn try_parse_create_item(bytes: &[u8]) -> Result<Self, EnclaveError> {
        use proto::item::item::MsgCreateItem;

        let raw_msg = MsgCreateItem::parse_from_bytes(bytes)
            .map_err(|_| EnclaveError::FailedToDeserialize)?;

        trace!(
            "try_parse_create_item creator: len={} val={:?}",
            raw_msg.creator.len(),
            raw_msg.creator
        );

        let sender = CanonicalAddr::from_human(&HumanAddr(raw_msg.creator)).map_err(|err| {
            warn!("failed to turn human addr to canonical addr when try_parse_create_item CosmWasmMsg: {:?}", err);
            EnclaveError::FailedToDeserialize
        })?;

        let tags = Self::parse_vec(raw_msg.tags)?;
        let shipping_region = Self::parse_vec(raw_msg.shipping_region)?;

        let photos = Self::parse_vec(raw_msg.photos)?;


        Ok(CosmWasmMsg::CreateItem {
            creator: sender,
            title: raw_msg.title,
            description:raw_msg.description,
            shipping_cost: raw_msg.shipping_cost,
            local_pickup:raw_msg.local_pickup,
            estimation_count: raw_msg.estimation_count,
            tags: tags,
            condition: raw_msg.condition,
            shipping_region: shipping_region,
            deposit_amount: raw_msg.deposit_amount,
            init_msg: raw_msg.init_msg,
            auto_msg: raw_msg.auto_msg,
            photos: photos,
            token_uri: raw_msg.token_uri,
        })
    }



    fn try_parse_execute(bytes: &[u8]) -> Result<Self, EnclaveError> {
        use proto::cosmwasm::msg::MsgExecuteContract;

        let raw_msg = MsgExecuteContract::parse_from_bytes(bytes)
            .map_err(|_| EnclaveError::FailedToDeserialize)?;

        trace!(
            "try_parse_execute sender: len={} val={:?}",
            raw_msg.sender.len(),
            raw_msg.sender
        );
        trace!(
            "try_parse_execute contract: len={} val={:?}",
            raw_msg.contract.len(),
            raw_msg.contract
        );
        let sender = CanonicalAddr::from_human(&HumanAddr(raw_msg.sender)).map_err(|err| {
            warn!("failed to turn human addr to canonical addr when try_parse_create_item CosmWasmMsg: {:?}", err);
            EnclaveError::FailedToDeserialize
        })?;
     

        if raw_msg.contract.clone().len() == 0  {
          /*  warn!(
                "Contract address was empty: {}",
                raw_msg.contract.len(),
            );*/
            EnclaveError::FailedToDeserialize;
        };
       
        let sent_funds = Self::parse_funds(raw_msg.sent_funds)?;

        let callback_sig = Some(raw_msg.callback_sig);

        Ok(CosmWasmMsg::Execute {
            sender: sender,
            contract: HumanAddr(raw_msg.contract),
            msg: raw_msg.msg,
            sent_funds,
            callback_sig,
        })
    }

    fn try_parse_create_estimation(bytes: &[u8]) -> Result<Self, EnclaveError> {
        use proto::item::estimator::MsgCreateEstimation;

        let raw_msg = MsgCreateEstimation::parse_from_bytes(bytes)
            .map_err(|_| EnclaveError::FailedToDeserialize)?;

        trace!(
            "try_parse_create_estimation sender: len={} val={:?}",
            raw_msg.estimator.len(),
            raw_msg.estimator
        );
        let sender = CanonicalAddr::from_human(&HumanAddr(raw_msg.estimator)).map_err(|err| {
            warn!("failed to turn human addr to canonical addr when try_parse_create_item CosmWasmMsg: {:?}", err);
            EnclaveError::FailedToDeserialize
        })?;

      //  let deposit_funds = Self::parse_funds(raw_msg.deposit)?;

        Ok(CosmWasmMsg::CreateEstimation {
            estimator: sender,
            estimate_msg: raw_msg.estimate_msg,
            itemid: raw_msg.itemid,
            deposit:  raw_msg.deposit,
            interested: raw_msg.interested,
        })
    }
    fn try_parse_flag_item(bytes: &[u8]) -> Result<Self, EnclaveError> {
        use proto::item::estimator::MsgFlagItem;

        let raw_msg = MsgFlagItem::parse_from_bytes(bytes)
            .map_err(|_| EnclaveError::FailedToDeserialize)?;

        trace!(
            "try_parse_flag_item sender: len={} val={:?}",
            raw_msg.estimator.len(),
            raw_msg.estimator
        );
        let sender = CanonicalAddr::from_human(&HumanAddr(raw_msg.estimator)).map_err(|err| {
            warn!("failed to turn human addr to canonical addr when try_parse_flag_item CosmWasmMsg: {:?}", err);
            EnclaveError::FailedToDeserialize
        })?;

        Ok(CosmWasmMsg::FlagItem {
            estimator:sender,
            itemid: raw_msg.itemid,
            flag_msg: raw_msg.flag_msg,
        })
    }
    fn try_parse_item_transferable(bytes: &[u8]) -> Result<Self, EnclaveError> {
        use proto::item::item::MsgItemTransferable;

        let raw_msg = MsgItemTransferable::parse_from_bytes(bytes)
            .map_err(|_| EnclaveError::FailedToDeserialize)?;

        trace!(
            "try_parse_flag_item sender: len={} val={:?}",
            raw_msg.seller.len(),
            raw_msg.seller
        );
        let sender = CanonicalAddr::from_human(&HumanAddr(raw_msg.seller)).map_err(|err| {
            warn!("failed to turn human addr to canonical addr when try_parse_flag_item CosmWasmMsg: {:?}", err);
            EnclaveError::FailedToDeserialize
        })?;

        Ok(CosmWasmMsg::ItemTransferable {
            seller:sender,
            transferable_msg: raw_msg.transferable_msg,
            itemid: raw_msg.itemid,
           
        })
    }
    fn try_parse_reveal_estimation(bytes: &[u8]) -> Result<Self, EnclaveError> {
        use proto::item::item::MsgRevealEstimation;

        let raw_msg = MsgRevealEstimation::parse_from_bytes(bytes)
            .map_err(|_| EnclaveError::FailedToDeserialize)?;

        trace!(
            "try_parse_reveal_estimation sender: len={} val={:?}",
            raw_msg.creator.len(),
            raw_msg.creator
        );
        let sender = CanonicalAddr::from_human(&HumanAddr(raw_msg.creator)).map_err(|err| {
            warn!("failed to turn human addr to canonical addr when try_parse_reveal_estimation CosmWasmMsg: {:?}", err);
            EnclaveError::FailedToDeserialize
        })?;

        Ok(CosmWasmMsg::RevealEstimation {
            creator: sender,
            itemid: raw_msg.itemid,
            reveal_msg: raw_msg.reveal_msg,
        })
    }
    pub fn parse_funds(
        raw_init_funds: protobuf::RepeatedField<proto::base::coin::Coin>,
    ) -> Result<Vec<Coin>, EnclaveError> {
        let mut init_funds = Vec::with_capacity(raw_init_funds.len());
        for raw_coin in raw_init_funds {
            let amount: u128 = raw_coin.amount.parse().map_err(|_err| {
                warn!(
                    "instantiate message funds were not a numeric string: {:?}",
                    raw_coin.amount,
                );
                EnclaveError::FailedToDeserialize
            })?;
            let coin = Coin {
                amount: Uint128(amount),
                denom: raw_coin.denom,
            };
            init_funds.push(coin);
        }

        Ok(init_funds)
    }

    fn parse_vec(
        raw_vec: protobuf::RepeatedField<::std::string::String>,
    ) -> Result<Vec<String>, EnclaveError> {
        let mut vec = Vec::with_capacity(raw_vec.len());
        for raw_string in raw_vec {
            vec.push(raw_string)
        }

        Ok(vec)
    }

    pub fn sender(&self) -> Option<&CanonicalAddr> {
        match self {
            CosmWasmMsg::Execute { sender, .. } | CosmWasmMsg::Instantiate { sender, .. } | CosmWasmMsg::CreateItem { creator: sender, .. }| CosmWasmMsg::RevealEstimation { creator: sender, .. } | CosmWasmMsg::FlagItem { estimator: sender, .. }| CosmWasmMsg::ItemTransferable { seller: sender, .. }| CosmWasmMsg::CreateEstimation { estimator: sender, .. }=> {
                Some(sender)
            }
            CosmWasmMsg::Other => None,
        }
    }
}

#[derive(Debug)]
pub struct AuthInfo {
    pub signer_infos: Vec<SignerInfo>,
    // Leaving this here for discoverability. We can use this, but don't verify it today.
    fee: (),
}

impl AuthInfo {
    pub fn from_bytes(bytes: &[u8]) -> Result<Self, EnclaveError> {
        let raw_auth_info = proto::tx::tx::AuthInfo::parse_from_bytes(&bytes).map_err(|err| {
            warn!("Could not parse AuthInfo from protobuf bytes: {:?}", err);
            EnclaveError::FailedToDeserialize
        })?;

        let mut signer_infos = vec![];
        for raw_signer_info in raw_auth_info.signer_infos {
            let signer_info = SignerInfo::from_proto(raw_signer_info)?;
            signer_infos.push(signer_info);
        }

        if signer_infos.is_empty() {
            warn!("No signature information provided for this TX. signer_infos empty");
            return Err(EnclaveError::FailedToDeserialize);
        }

        Ok(Self {
            signer_infos,
            fee: (),
        })
    }

    pub fn sender_public_key(&self, sender: &CanonicalAddr) -> Option<&CosmosPubKey> {
        self.signer_infos
            .iter()
            .find(|signer_info| &signer_info.public_key.get_address() == sender)
            .map(|si| &si.public_key)
    }
}

#[derive(Debug)]
pub struct SignerInfo {
    pub public_key: CosmosPubKey,
    pub sequence: u64,
}

impl SignerInfo {
    pub fn from_proto(raw_signer_info: proto::tx::tx::SignerInfo) -> Result<Self, EnclaveError> {
        if !raw_signer_info.has_public_key() {
            warn!("One of the provided signers had no public key");
            return Err(EnclaveError::FailedToDeserialize);
        }

        // unwraps valid after checks above
        let any_public_key = raw_signer_info.public_key.get_ref();

        let public_key = CosmosPubKey::from_proto(any_public_key)
            .map_err(|_| EnclaveError::FailedToDeserialize)?;

        let signer_info = Self {
            public_key,
            sequence: raw_signer_info.sequence,
        };
        Ok(signer_info)
    }
}

#[cfg(feature = "test")]
pub mod tests {

    use super::*;
    // use crate::crypto::{AESKey, SIVEncryptable, Seed, KEY_MANAGER};

    // todo: fix test vectors to actually work
    pub fn test_new_from_slice() {
        let nonce = [0u8; 32];
        let user_public_key = [0u8; 32];
        let msg = "{\"ok\": \"{\"balance\": \"108\"}\"}";

        let mut slice = nonce.to_vec();
        slice.extend_from_slice(&user_public_key);
        slice.extend_from_slice(msg.as_bytes());

        let secret_msg = SecretMessage {
            nonce,
            user_public_key,
            msg: msg.as_bytes().to_vec(),
        };

        let msg_from_slice = SecretMessage::from_slice(&slice).unwrap();

        assert_eq!(secret_msg, msg_from_slice);
    }

    // This is commented out because it's trying to modify KEY_MANAGER which is immutable.
    // // todo: fix test vectors to actually work
    // pub fn test_msg_decrypt() {
    //     let seed = Seed::new().unwrap();
    //
    //     KEY_MANAGER
    //         .set_consensus_seed(seed)
    //         .expect("Failed to set seed");
    //
    //     let nonce = [0u8; 32];
    //     let user_public_key = [0u8; 32];
    //
    //     let msg = "{\"ok\": \"{\"balance\": \"108\"}\"}";
    //     let key = calc_encryption_key(&nonce, &user_public_key);
    //
    //     let encrypted_msg = key.encrypt_siv(msg.as_bytes(), &[&[]]);
    //
    //     let secret_msg = SecretMessage {
    //         nonce,
    //         user_public_key,
    //         msg: encrypted_msg,
    //     };
    //
    //     let decrypted_msg = secret_msg.decrypt()?;
    //
    //     assert_eq!(decrypted_msg, msg)
    // }
}
