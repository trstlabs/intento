use log::*;

use enclave_ffi_types::{Ctx, EnclaveError};

use crate::external::results::{CallbackSigSuccess, HandleSuccess, InitSuccess, QuerySuccess};

use enclave_cosmos_types::types::{ContractCode, HandleType, MsgInfo, SigInfo};

use enclave_cosmwasm_types::addresses::{Addr, CanonicalAddr};
use enclave_cosmwasm_types::coins::Coin;
use enclave_cosmwasm_types::encoding::Binary;
use enclave_cosmwasm_types::results::{ Event, Reply, SubMsgResponse, SubMsgResult};
use enclave_cosmwasm_types::types::{BlockInfo, ContractInfo, MessageInfo};
use enclave_cosmwasm_types::types::{Env, FullEnv};
//use enclave_cosmwasm_types::timestamp::Timestamp;

use enclave_crypto::Ed25519PublicKey;
use enclave_utils::coalesce;
use super::parse_msg::parse_message;
use super::contract_validation::{
    extract_contract_key, generate_encryption_key, validate_contract_key, validate_msg,
    verify_params, ContractKey, ReplyParams, ValidatedMessage,
};
use super::gas::WasmCosts;
use super::io::{copy_into_array, create_callback_signature, encrypt_msg, encrypt_output};
use super::module_cache::create_module_instance;
use super::types::{ContractMessage, IoNonce};
use super::wasm::{ContractInstance, ContractOperation, Engine};
use crate::const_callback_sig_addresses::COMMUNITY_POOL_ADDR;


/*
Each contract is compiled with these functions already implemented in wasm:
fn cosmwasm_api_0_6() -> i32;  // Seems unused, but we should support it anyways
fn allocate(size: usize) -> *mut c_void;
fn deallocate(pointer: *mut c_void);
fn init(env_ptr: *mut c_void, msg_ptr: *mut c_void) -> *mut c_void
fn handle(env_ptr: *mut c_void, msg_ptr: *mut c_void) -> *mut c_void
fn query(msg_ptr: *mut c_void) -> *mut c_void

Re `init`, `handle` and `query`: We need to pass `env` & `msg`
down to the wasm implementations, but because they are buffers
we need to allocate memory regions inside the VM's instance and copy
`env` & `msg` into those memory regions inside the VM's instance.
*/
#[allow(clippy::too_many_arguments)]
pub fn init(
    context: Ctx,       // need to pass this to read_db & write_db
    gas_limit: u64,     // gas limit for this execution
    used_gas: &mut u64, // out-parameter for gas used in execution
    contract: &[u8],    // contract code wasm bytes
    env: &[u8],         // blockchain state
    msg: &[u8],         // can contain function calls and args
    auto_msg: &[u8],    // can contain auto function calls and args. we create and return sigature.
    sig_info: &[u8],    // info about signature verification
) -> Result<InitSuccess, EnclaveError> {
    let contract_code = ContractCode::new(contract);
    let mut parsed_env: FullEnv = serde_json::from_slice(env).map_err(|err| {
        warn!(
            "got an error while trying to deserialize env input bytes into json {:?}: {}",
            String::from_utf8_lossy(&env),
            err
        );
        EnclaveError::FailedToDeserialize
    })?;
    parsed_env.contract.code_hash = hex::encode(contract_code.hash());

    let canonical_contract_address = CanonicalAddr::from_addr(&parsed_env.contract.address).map_err(|err| {
        warn!(
            "got an error while trying to deserialize parsed_env.contract.address from bech32 string to bytes {:?}: {}",
            parsed_env.contract.address, err
        );
        EnclaveError::FailedToDeserialize
    })?;

    let canonical_sender_address = CanonicalAddr::from_addr(&parsed_env.message.sender).map_err(|err| {
        warn!(
            "init got an error while trying to deserialize parsed_env.message.sender from bech32 string to bytes {:?}: {}",
            parsed_env.message.sender, err
        );
        EnclaveError::FailedToDeserialize
    })?;

    let contract_key = generate_encryption_key(
        &parsed_env,
        contract_code.hash(),
        &(canonical_contract_address.0).0,
    )?;
    trace!("init contract key: {:?}", hex::encode(contract_key));

    let parsed_sig_info: SigInfo = serde_json::from_slice(sig_info).map_err(|err| {
        warn!(
            "init got an error while trying to deserialize env input bytes into json {:?}: {}",
            String::from_utf8_lossy(&sig_info),
            err
        );
        EnclaveError::FailedToDeserialize
    })?;

    trace!("Init input before decryption: {:?}", base64::encode(&msg));
    let contract_msg = ContractMessage::from_slice(msg)?;

    verify_params(&parsed_sig_info, &parsed_env, &contract_msg)?;

    let decrypted_msg = contract_msg.decrypt()?;

    let ValidatedMessage {
        validated_msg,
        reply_params,
    } = validate_msg(&decrypted_msg, &contract_code.hash(), None, None)?;

    trace!(
        "init input after decryption: {:?}",
        String::from_utf8_lossy(&validated_msg)
    );

    let mut engine = start_engine(
        context,
        gas_limit,
        contract_code,
        &contract_key,
        ContractOperation::Init,
        contract_msg.nonce,
        contract_msg.user_public_key,
    )?;

    let (contract_env_bytes, contract_msg_info_bytes) = parse_msg_info_bytes(&mut parsed_env)?;

    let env_ptr = engine.write_to_memory(&contract_env_bytes)?;
    let msg_info_ptr = engine.write_to_memory(&contract_msg_info_bytes)?;
    let msg_ptr = engine.write_to_memory(&validated_msg)?;
    let auto_msg = ContractMessage::from_slice(auto_msg)?;

    let funds = [];
    trace!("funds {:?}...", funds);
    let sig = create_callback_signature(&canonical_contract_address, &auto_msg, &funds);

    let array = copy_into_array(&sig[..]);
    let callback_sig: [u8; 32] = array;

    // This wrapper is used to coalesce all errors in this block to one object
    // so we can `.map_err()` in one place for all of them
    let output = coalesce!(EnclaveError, {
        let vec_ptr = engine.init(env_ptr, msg_info_ptr, msg_ptr)?;
        let output = engine.extract_vector(vec_ptr)?;
        // TODO: copy cosmwasm's structures to enclave
        // TODO: ref: https://github.com/CosmWasm/cosmwasm/blob/b971c037a773bf6a5f5d08a88485113d9b9e8e7b/packages/std/src/init_handle.rs#L129
        // TODO: ref: https://github.com/CosmWasm/cosmwasm/blob/b971c037a773bf6a5f5d08a88485113d9b9e8e7b/packages/std/src/query.rs#L13
        let output = encrypt_output(
            output,
            &contract_msg,
            &canonical_contract_address,
            &parsed_env.contract.code_hash,
            reply_params,
            &canonical_sender_address,
        )?;

        Ok(output)
    })
    .map_err(|err| {
        *used_gas = engine.gas_used();
        err
    })?;

    *used_gas = engine.gas_used();
    // todo: can move the key to somewhere in the output message if we want

    Ok(InitSuccess {
        output,
        contract_key,
        callback_sig,
    })
}

pub struct TaggedBool {
    b: bool,
}

impl From<bool> for TaggedBool {
    fn from(b: bool) -> Self {
        TaggedBool { b }
    }
}

impl Into<bool> for TaggedBool {
    fn into(self) -> bool {
        self.b
    }
}

pub struct ParsedMessage {
    pub should_validate_sig_info: bool,
    pub was_msg_encrypted: bool,
    pub contract_msg: ContractMessage,
    pub decrypted_msg: Vec<u8>,
    pub contract_hash_for_validation: Option<Vec<u8>>,
}

pub fn reduct_custom_events(reply: &mut Reply) {
    reply.result = match &reply.result {
        SubMsgResult::Ok(r) => {
            let events: Vec<Event> = Default::default();
            /*  let filtered_types = vec![
              //  "execute".to_string(),
              //  "instantiate".to_string(),
                "wasm".to_string(),
            ];
            let filtered_attributes = vec!["contract_address".to_string(), "code_id".to_string()];
            for ev in r.events.iter() {
                if filtered_types.contains(&ev.ty) {
                    let mut new_ev = Event {
                        ty: ev.ty.clone(),
                        attributes: vec![],
                    };

                    for attr in &ev.attributes {
                        if !filtered_attributes.contains(&attr.key) {
                            new_ev.attributes.push(attr.clone());
                        }
                    }

                    if new_ev.attributes.len() > 0 {
                        events.push(new_ev);
                    }

                }
            }*/

            SubMsgResult::Ok(SubMsgResponse {
                events,
                data: r.data.clone(),
            })
        }
        SubMsgResult::Err(_) => reply.result.clone(),
    };
}

#[allow(clippy::too_many_arguments)]
pub fn handle(
    context: Ctx,
    gas_limit: u64,
    used_gas: &mut u64,
    contract: &[u8],
    env: &[u8],
    msg: &[u8],
    sig_info: &[u8],
    handle_type: u8,
) -> Result<HandleSuccess, EnclaveError> {
    let contract_code = ContractCode::new(contract);

    let mut parsed_env: FullEnv = serde_json::from_slice(env).map_err(|err| {
        warn!(
            "got an error while trying to deserialize env input bytes into json {:?}: {}",
            env, err
        );
        EnclaveError::FailedToDeserialize
    })?;
    parsed_env.contract.code_hash = hex::encode(contract_code.hash());

    let canonical_contract_address = CanonicalAddr::from_addr(&&parsed_env.contract.address).map_err(|err| {
        warn!(
            "got an error while trying to deserialize parsed_env.contract.address from bech32 string to bytes {:?}: {}",
            parsed_env.contract.address, err
        );
        EnclaveError::FailedToDeserialize
    })?;

    let canonical_sender_address = match to_canonical(&parsed_env.message.sender) {
        Ok(can) => can,
        Err(_) => CanonicalAddr::from_vec(vec![]),
    };

    let contract_key = extract_contract_key(&parsed_env)?;

    if !validate_contract_key(&contract_key, &canonical_contract_address, &contract_code) {
        warn!("got an error while trying to deserialize output bytes");
        return Err(EnclaveError::FailedContractAuthentication);
    }

    let parsed_sig_info: SigInfo = serde_json::from_slice(sig_info).map_err(|err| {
        warn!(
            "handle got an error while trying to deserialize sig info input bytes into json {:?}: {}",
            String::from_utf8_lossy(&sig_info),
            err
        );
        EnclaveError::FailedToDeserialize
    })?;

    // The flow of handle is used for multiple messages (such ash Handle, Reply)
    // When the message is handle, we expect it always to be encrypted while in Reply for example it might be plaintext
    let parsed_handle_type = HandleType::try_from(handle_type)?;

    let ParsedMessage {
        should_validate_sig_info,
        was_msg_encrypted,
        contract_msg, // params to be verified with reducted events. Should equal callback sig.
        decrypted_msg, //to be validated. Complete message.
        contract_hash_for_validation,
    } = parse_message(msg, &parsed_sig_info, &parsed_handle_type)?;

    trace!(
        "handle input after decryption: {:?}",
        String::from_utf8_lossy(&decrypted_msg)
    );

    // There is no signature to verify when the input isn't signed.
    // Receiving unsigned messages is only possible in Handle. (Init tx are always signed)
    // All of these functions go through handle but the data isn't signed:
    //  Reply (that is not WASM reply)
    if should_validate_sig_info {
        // Verify env parameters against the signed tx
        verify_params(&parsed_sig_info, &parsed_env, &contract_msg)?;
    }

    let mut validated_msg = decrypted_msg.clone();
    let mut reply_params: Option<Vec<ReplyParams>> = None;
    if was_msg_encrypted {
        let x = validate_msg(
            &decrypted_msg,
            &contract_code.hash(),
            contract_hash_for_validation,
            Some(parsed_handle_type.clone()),
        )?;
        validated_msg = x.validated_msg;
        reply_params = x.reply_params;
    }

    trace!("Successfully authenticated the contract!");

    trace!("Handle: Contract Key: {:?}", hex::encode(contract_key));

    // Although the operation here is not always handle it is irrelevant in this case
    // because it only helps to decide whether to check floating points or not
    // In this case we want to do the same as in Handle both for Reply and for others so we can always pass "Handle".
    let mut engine = start_engine(
        context,
        gas_limit,
        contract_code,
        &contract_key,
        ContractOperation::Handle,
        contract_msg.nonce,
        contract_msg.user_public_key,
    )?;

    let (contract_env_bytes, contract_msg_info_bytes) = parse_msg_info_bytes(&mut parsed_env)?;

    let env_ptr = engine.write_to_memory(&contract_env_bytes)?;
    let msg_info_ptr = engine.write_to_memory(&contract_msg_info_bytes)?;
    let msg_ptr = engine.write_to_memory(&validated_msg)?;

    // This wrapper is used to coalesce all errors in this block to one object
    // so we can `.map_err()` in one place for all of them
    let output = coalesce!(EnclaveError, {
        let vec_ptr = engine.handle(env_ptr, msg_info_ptr, msg_ptr, parsed_handle_type)?;

        let output = engine.extract_vector(vec_ptr)?;

        debug!(
            "(2) nonce just before encrypt_output: nonce = {:?} pubkey = {:?}",
            contract_msg.nonce, contract_msg.user_public_key
        );

        let output = encrypt_output(
            output,
            &contract_msg,
            &canonical_contract_address,
            &parsed_env.contract.code_hash,
            reply_params,
            &canonical_sender_address,
        )?;

        Ok(output)
    })
    .map_err(|err| {
        *used_gas = engine.gas_used();
        err
    })?;

    *used_gas = engine.gas_used();
    Ok(HandleSuccess { output })
}

pub fn query(
    context: Ctx,
    gas_limit: u64,
    used_gas: &mut u64,
    contract: &[u8],
    env: &[u8],
    msg: &[u8],
) -> Result<QuerySuccess, EnclaveError> {
    let contract_code = ContractCode::new(contract);

    let mut parsed_env: FullEnv = serde_json::from_slice(env).map_err(|err| {
        warn!(
            "query got an error while trying to deserialize env input bytes into json {:?}: {}",
            env, err
        );
        EnclaveError::FailedToDeserialize
    })?;
    parsed_env.contract.code_hash = hex::encode(contract_code.hash());

    trace!("query env: {:?}", parsed_env);

    let canonical_contract_address = CanonicalAddr::from_addr(&parsed_env.contract.address).map_err(|err| {
        warn!(
            "got an error while trying to deserialize parsed_env.contract.address from bech32 string to bytes {:?}: {}",
            parsed_env.contract.address, err
        );
        EnclaveError::FailedToDeserialize
    })?;

    let contract_key = extract_contract_key(&parsed_env)?;

    if !validate_contract_key(&contract_key, &canonical_contract_address, &contract_code) {
        warn!("query got an error while trying to validate contract key");
        return Err(EnclaveError::FailedContractAuthentication);
    }

    trace!("successfully authenticated the contract!");
    trace!("query contract key: {:?}", hex::encode(contract_key));

    trace!("query input before decryption: {:?}", base64::encode(&msg));
    let contract_msg = ContractMessage::from_slice(msg)?;
    let decrypted_msg = contract_msg.decrypt()?;
    trace!(
        "query input afer decryption: {:?}",
        String::from_utf8_lossy(&decrypted_msg)
    );
    let ValidatedMessage { validated_msg, .. } =
        validate_msg(&decrypted_msg, &contract_code.hash(), None, None)?;

    let mut engine = start_engine(
        context,
        gas_limit,
        contract_code,
        &contract_key,
        ContractOperation::Query,
        contract_msg.nonce,
        contract_msg.user_public_key,
    )?;

    let (contract_env_bytes, _ /* no msg_info in query */) = parse_msg_info_bytes(&mut parsed_env)?;

    let env_ptr = engine.write_to_memory(&contract_env_bytes)?;
    let msg_ptr = engine.write_to_memory(&validated_msg)?;

    // This wrapper is used to coalesce all errors in this block to one object
    // so we can `.map_err()` in one place for all of them
    let output = coalesce!(EnclaveError, {
        let vec_ptr = engine.query(env_ptr, msg_ptr)?;

        let output = engine.extract_vector(vec_ptr)?;

        let output = encrypt_output(
            output,
            &contract_msg,
            &CanonicalAddr(Binary(Vec::new())), // Not used for queries (can't init a new contract from a query)
            &"".to_string(), // Not used for queries (can't call a sub-message from a query),
            None,            // Not used for queries (Query response is not replied to the caller),
            &CanonicalAddr(Binary(Vec::new())), // Not used for queries (used only for replies)
        )?;
        Ok(output)
    })
    .map_err(|err| {
        *used_gas = engine.gas_used();
        err
    })?;

    *used_gas = engine.gas_used();
    Ok(QuerySuccess { output })
}

fn start_engine(
    context: Ctx,
    gas_limit: u64,
    contract_code: ContractCode,
    contract_key: &ContractKey,
    operation: ContractOperation,
    nonce: IoNonce,
    user_public_key: Ed25519PublicKey,
) -> Result<Engine, EnclaveError> {
    let module = create_module_instance(contract_code, operation)?;

    // Set the gas costs for wasm op-codes (there is an inline stack_height limit in WasmCosts)
    let wasm_costs = WasmCosts::default();

    let contract_instance = ContractInstance::new(
        context,
        module.clone(),
        gas_limit,
        wasm_costs,
        *contract_key,
        operation,
        nonce,
        user_public_key,
    )?;

    Ok(Engine::new(contract_instance, module))
}

fn parse_msg_info_bytes(env: &mut FullEnv) -> Result<(Vec<u8>, Vec<u8>), EnclaveError> {
    let new_env = Env {
        block: BlockInfo {
            height: env.block.height,
            time: env.block.time,
            chain_id: env.block.chain_id.clone(),
        },
        contract: ContractInfo {
            address: Addr(env.contract.address.0.clone()),
            code_hash: env.contract.code_hash.clone(),
        },
    };

    let env_bytes = serde_json::to_vec(&new_env).map_err(|err| {
        warn!(
            "got an error while trying to serialize env (CosmWasm v1) into bytes {:?}: {}",
            env, err
        );
        EnclaveError::FailedToSerialize
    })?;

    let msg_info = MessageInfo {
        sender: Addr(env.message.sender.clone().to_string()),
        funds: env
            .message
            .funds
            .iter()
            .map(|coin| Coin::new(coin.amount.u128(), coin.denom.clone()))
            .collect::<Vec<enclave_cosmwasm_types::coins::Coin>>(),
    };

    let msg_info_bytes = serde_json::to_vec(&msg_info).map_err(|err| {
        warn!(
            "got an error while trying to serialize msg_info (CosmWasm v1) into bytes {:?}: {}",
            msg_info, err
        );
        EnclaveError::FailedToSerialize
    })?;

    Ok((env_bytes, msg_info_bytes))
}

pub fn create_callback_sig(
    msg: &[u8],      //message with args
    msg_info: &[u8], //code hash and funds
) -> Result<CallbackSigSuccess, EnclaveError> {
    let parsed_msg_info: MsgInfo = serde_json::from_slice(msg_info).map_err(|err| {
        warn!(
            "got an error while trying to deserialize msg input bytes into json {:?}: {}",
            String::from_utf8_lossy(&msg_info),
            err
        );
        EnclaveError::FailedToDeserialize
    })?;

    let send_as_addr =
        CanonicalAddr::from_addr(&Addr(COMMUNITY_POOL_ADDR.to_string())).map_err(|err| {
            warn!(
                "failed to turn human addr to canonical addr when create_callback_sig: {:?}",
                err
            );
            EnclaveError::FailedToDeserialize
        })?;

    let nonce_placeholder = [0u8; 32];
    let pubkey_placeholder = [0u8; 32];

    let mut msg_callback = Binary::from(msg);
    let sig = encrypt_msg(
        &mut msg_callback,
        nonce_placeholder,
        pubkey_placeholder,
        &send_as_addr,
        hex::encode(&parsed_msg_info.code_hash.as_slice()),
        parsed_msg_info.funds,
    )
    .map_err(|err| {
        warn!(
            "got an error while trying to encrypt wasm_msg into encrypted message {:?}: {}",
            msg, err
        );
        EnclaveError::FailedToDeserialize
    })?;

    let array = copy_into_array(&sig[..]);
    let callback_sig: [u8; 32] = array;
    Ok(CallbackSigSuccess {
        callback_sig,
        encrypted_msg: msg_callback.as_slice().to_vec(),
    })
}

fn to_canonical(contract_address: &Addr) -> Result<CanonicalAddr, EnclaveError> {
    CanonicalAddr::from_addr(contract_address).map_err(|err| {
        warn!(
            "error while trying to deserialize address from bech32 string to bytes {:?}: {}",
            contract_address, err
        );
        EnclaveError::FailedToDeserialize
    })
}