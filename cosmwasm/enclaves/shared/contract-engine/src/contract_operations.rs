use log::*;

use enclave_ffi_types::{Ctx, EnclaveError};
use std::str::from_utf8;
use enclave_cosmos_types::types::{ContractCode, SigInfo, MsgInfo};
use enclave_cosmwasm_types::encoding::Binary;
use enclave_cosmwasm_types::types::{CanonicalAddr, HumanAddr, Env, WasmMsg};
use enclave_crypto::Ed25519PublicKey;
use enclave_utils::coalesce;
use std::convert::TryInto;
use crate::external::results::{HandleSuccess, InitSuccess, QuerySuccess, CallbackSigSuccess};

use super::contract_validation::{
    extract_contract_key, generate_encryption_key, validate_contract_key, validate_msg,
    verify_params, ContractKey,
};
use super::gas::WasmCosts;
use super::io::{create_callback_signature,  encrypt_output, copy_into_array, encrypt_msg};
use super::module_cache::create_module_instance;
use super::types::{IoNonce, ContractMessage};
use super::wasm::{ContractInstance, ContractOperation, Engine};
use crate::const_callback_sig_addresses::{COMMUNITY_POOL_ADDR};
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

    let mut parsed_env: Env = serde_json::from_slice(env).map_err(|err| {
        warn!(
            "got an error while trying to deserialize env input bytes into json {:?}: {}",
            String::from_utf8_lossy(&env),
            err
        );
        EnclaveError::FailedToDeserialize
    })?;
    parsed_env.contract_code_hash = hex::encode(contract_code.hash());

    let canonical_contract_address = CanonicalAddr::from_human(&parsed_env.contract.address).map_err(|err| {
        warn!(
            "got an error while trying to deserialize parsed_env.contract.address from bech32 string to bytes {:?}: {}",
            parsed_env.contract.address, err
        );
        EnclaveError::FailedToDeserialize
    })?;
    let contract_key = generate_encryption_key(
        &parsed_env,
        contract_code.hash(),
        &(canonical_contract_address.0).0,
    )?;
    trace!("Init: Contract Key: {:?}", hex::encode(contract_key));

    let parsed_sig_info: SigInfo = serde_json::from_slice(sig_info).map_err(|err| {
        warn!(
            "got an error while trying to deserialize sig_info input bytes into json {:?}: {}",
            String::from_utf8_lossy(&sig_info),
            err
        );
        EnclaveError::FailedToDeserialize
    })?;

    trace!("Init input before decryption: {:?}", base64::encode(&msg));
    let contract_msg = ContractMessage::from_slice(msg)?;

    verify_params(&parsed_sig_info, &parsed_env, &contract_msg)?;

    let decrypted_msg = contract_msg.decrypt()?;

    let validated_msg = validate_msg(&decrypted_msg, contract_code.hash())?;

    trace!(
        "Init input after decryption: {:?}",
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

    let new_env = serde_json::to_vec(&parsed_env).map_err(|err| {
        warn!(
            "got an error while trying to serialize parsed_env into bytes {:?}: {}",
            parsed_env, err
        );
        EnclaveError::FailedToSerialize
    })?;

    let env_ptr = engine.write_to_memory(&new_env)?;
    let msg_ptr = engine.write_to_memory(&validated_msg)?;

    let auto_msg = ContractMessage::from_slice(auto_msg)?;

    let sent_funds = [];
    trace!("sent_funds {:?}...", sent_funds);
    let sig = create_callback_signature(&canonical_contract_address, &auto_msg, &sent_funds);

    let array = copy_into_array(&sig[..]);
   // trace!("array length...{:?}", array.len());
    let callback_sig: [u8; 32] = array;
    // This wrapper is used to coalesce all errors in this block to one object
    // so we can `.map_err()` in one place for all of them
    let output = coalesce!(EnclaveError, {
        let vec_ptr = engine.init(env_ptr, msg_ptr)?;
        let output = engine.extract_vector(vec_ptr)?;
        // TODO: copy cosmwasm's structures to enclave
        // TODO: ref: https://github.com/CosmWasm/cosmwasm/blob/b971c037a773bf6a5f5d08a88485113d9b9e8e7b/packages/std/src/init_handle.rs#L129
        // TODO: ref: https://github.com/CosmWasm/cosmwasm/blob/b971c037a773bf6a5f5d08a88485113d9b9e8e7b/packages/std/src/query.rs#L13
        let output = encrypt_output(
            output,
            contract_msg.nonce,
            contract_msg.user_public_key,
            &canonical_contract_address,
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

pub fn handle(
    context: Ctx,
    gas_limit: u64,
    used_gas: &mut u64,
    contract: &[u8],
    env: &[u8],
    msg: &[u8],
    sig_info: &[u8],
) -> Result<HandleSuccess, EnclaveError> {
    let contract_code = ContractCode::new(contract);

    let mut parsed_env: Env = serde_json::from_slice(env).map_err(|err| {
        warn!(
            "got an error while trying to deserialize env input bytes into json {:?}: {}",
            env, err
        );
        EnclaveError::FailedToDeserialize
    })?;
    parsed_env.contract_code_hash = hex::encode(contract_code.hash());

    let canonical_contract_address = CanonicalAddr::from_human(&parsed_env.contract.address).map_err(|err| {
        warn!(
            "got an error while trying to deserialize parsed_env.contract.address from bech32 string to bytes {:?}: {}",
            parsed_env.contract.address, err
        );
        EnclaveError::FailedToDeserialize
    })?;

    let contract_key = extract_contract_key(&parsed_env)?;

    if !validate_contract_key(&contract_key, &canonical_contract_address, &contract_code) {
        warn!("got an error while trying to deserialize output bytes");
        return Err(EnclaveError::FailedContractAuthentication);
    }

    trace!("handle parsed_env: {:?}", parsed_env);

    let parsed_sig_info: SigInfo = serde_json::from_slice(sig_info).map_err(|err| {
        warn!(
            "got an error while trying to deserialize env input bytes into json {:?}: {}",
            String::from_utf8_lossy(&sig_info),
            err
        );
        EnclaveError::FailedToDeserialize
    })?;

    trace!("Handle input before decryption: {:?}", base64::encode(&msg));
    let contract_msg = ContractMessage::from_slice(msg)?;

    // Verify env parameters against the signed tx
    verify_params(&parsed_sig_info, &parsed_env, &contract_msg)?;

    let contract_msg = ContractMessage::from_slice(msg)?;
    let decrypted_msg = contract_msg.decrypt()?;

    let validated_msg = validate_msg(&decrypted_msg, contract_code.hash())?;

    trace!(
        "Handle input afer decryption: {:?}",
        String::from_utf8_lossy(&validated_msg)
    );

    trace!("Successfully authenticated the contract!");

    trace!("Handle: Contract Key: {:?}", hex::encode(contract_key));

    let mut engine = start_engine(
        context,
        gas_limit,
        contract_code,
        &contract_key,
        ContractOperation::Handle,
        contract_msg.nonce,
        contract_msg.user_public_key,
    )?;

    let new_env = serde_json::to_vec(&parsed_env).map_err(|err| {
        warn!(
            "got an error while trying to serialize parsed_env into bytes {:?}: {}",
            parsed_env, err
        );
        EnclaveError::FailedToSerialize
    })?;

    let env_ptr = engine.write_to_memory(&new_env)?;
    let msg_ptr = engine.write_to_memory(&validated_msg)?;

    // This wrapper is used to coalesce all errors in this block to one object
    // so we can `.map_err()` in one place for all of them
    let output = coalesce!(EnclaveError, {
        let vec_ptr = engine.handle(env_ptr, msg_ptr)?;

        let output = engine.extract_vector(vec_ptr)?;

        debug!(
            "(2) nonce just before encrypt_output: nonce = {:?} pubkey = {:?}",
            contract_msg.nonce, contract_msg.user_public_key
        );
        let output = encrypt_output(
            output,
            contract_msg.nonce,
            contract_msg.user_public_key,
            &canonical_contract_address,
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

    let mut parsed_env: Env = serde_json::from_slice(env).map_err(|err| {
        warn!(
            "query got an error while trying to deserialize env input bytes into json {:?}: {}",
            env, err
        );
        EnclaveError::FailedToDeserialize
    })?;
    parsed_env.contract_code_hash = hex::encode(contract_code.hash());

    trace!("query env_v010: {:?}", parsed_env);

    let canonical_contract_address = CanonicalAddr::from_human(&parsed_env.contract.address).map_err(|err| {
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
    let validated_msg = validate_msg(&decrypted_msg, contract_code.hash())?;

    let mut engine = start_engine(
        context,
        gas_limit,
        contract_code,
        &contract_key,
        ContractOperation::Query,
        contract_msg.nonce,
        contract_msg.user_public_key,
    )?;

    let msg_ptr = engine.write_to_memory(&validated_msg)?;

    // This wrapper is used to coalesce all errors in this block to one object
    // so we can `.map_err()` in one place for all of them
    let output = coalesce!(EnclaveError, {
        let vec_ptr = engine.query(msg_ptr)?;

        let output = engine.extract_vector(vec_ptr)?;

        let output = encrypt_output(
            output,
            contract_msg.nonce,
            contract_msg.user_public_key,
            &CanonicalAddr(Binary(Vec::new())), // Not used for queries (can't init a new contract from a query)
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
    );

    Ok(Engine::new(contract_instance, module))
}



pub fn create_callback_sig(
    msg: &[u8], //message with args
    msg_info: &[u8] //code hash and funds
) -> Result<CallbackSigSuccess, EnclaveError> {

    let parsed_msg_info: MsgInfo = serde_json::from_slice(msg_info).map_err(|err| {
        warn!(
            "got an error while trying to deserialize msg input bytes into json {:?}: {}",
            String::from_utf8_lossy(&msg_info),
            err
        );
        EnclaveError::FailedToDeserialize
    })?;              
   
    let send_as_addr = CanonicalAddr::from_human(&HumanAddr(COMMUNITY_POOL_ADDR.to_string())).map_err(|err| {
        warn!("failed to turn human addr to canonical addr when create_callback_sig: {:?}", err);
        EnclaveError::FailedToDeserialize
    })?;

    let mut nonce_placeholder = [0u8; 32];
    nonce_placeholder.copy_from_slice(&msg[0..32]);
    let mut pubkey_placeholder = [0u8; 32];
    pubkey_placeholder.copy_from_slice(&msg[32..64]);

    let mut msg_callback = Binary::from(msg);
    let sig = encrypt_msg(&mut msg_callback, nonce_placeholder, pubkey_placeholder, &send_as_addr, hex::encode(&parsed_msg_info.code_hash.as_slice()), parsed_msg_info.funds).map_err(|err| {
        warn!(
            "got an error while trying to encrypt wasm_msg into encrypted message {:?}: {}",
            msg,
            err
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

/// Look at the first 8 bytes of the input and reinterpret them as a u64
fn read_be_u64(input: &[u8]) -> u64 {
    assert!(input.len() >= std::mem::size_of::<u64>());
    u64::from_be_bytes(input[0..std::mem::size_of::<u64>()].try_into().unwrap())
}