// use serde::de::DeserializeOwned;
// use std::fmt;

// use cosmwasm_std::{Env, HandleResult, InitResult, MigrateResult, QueryResult};

// use crate::errors::{VmError, VmResult};
use crate::errors::VmResult;
/*
use crate::instance::{Func, Instance};
*/
use crate::instance::Instance;
// use crate::serde::{from_slice, to_vec};
use crate::traits::{Api, Querier, Storage};
// use schemars::JsonSchema;



/// Calls Wasm export "init" and returns raw data from the contract.
/// The result is length limited to prevent abuse but otherwise unchecked.
pub fn call_init_raw<S: Storage + 'static, A: Api + 'static, Q: Querier + 'static>(
    instance: &mut Instance<S, A, Q>,
    env: &[u8],
    msg: &[u8],
    auto_msg: &[u8],
    sig_info: &[u8],
) -> VmResult<Vec<u8>> {
    instance.set_storage_readonly(false);
    /*
    call_raw(instance, "init", &[env, msg], MAX_LENGTH_INIT)
    */
    instance.call_init(env, msg, auto_msg, sig_info)
}

/// Calls Wasm export "handle" and returns raw data from the contract.
/// The result is length limited to prevent abuse but otherwise unchecked.
pub fn call_handle_raw<S: Storage + 'static, A: Api + 'static, Q: Querier + 'static>(
    instance: &mut Instance<S, A, Q>,
    env: &[u8],
    msg: &[u8],
    sig_info: &[u8],
) -> VmResult<Vec<u8>> {
    instance.set_storage_readonly(false);
    /*
    call_raw(instance, "handle", &[env, msg], MAX_LENGTH_HANDLE)
    */
    instance.call_handle(env, msg, sig_info)
}

/// Calls Wasm export "migrate" and returns raw data from the contract.
/// The result is length limited to prevent abuse but otherwise unchecked.
pub fn call_migrate_raw<S: Storage + 'static, A: Api + 'static, Q: Querier + 'static>(
    instance: &mut Instance<S, A, Q>,
    env: &[u8],
    msg: &[u8],
) -> VmResult<Vec<u8>> {
    instance.set_storage_readonly(false);
    /*
    call_raw(instance, "migrate", &[env, msg], MAX_LENGTH_MIGRATE)
    */
    instance.call_migrate(env, msg)
}

/// Calls Wasm export "query" and returns raw data from the contract.
/// The result is length limited to prevent abuse but otherwise unchecked.
pub fn call_query_raw<S: Storage + 'static, A: Api + 'static, Q: Querier + 'static>(
    instance: &mut Instance<S, A, Q>,
    msg: &[u8],
) -> VmResult<Vec<u8>> {
    instance.set_storage_readonly(true);
    /*
    call_raw(instance, "query", &[msg], MAX_LENGTH_QUERY)
    */
    instance.call_query(msg)
}

#[cfg(not(feature = "default-enclave"))]
fn call_raw<S: Storage + 'static, A: Api + 'static, Q: Querier + 'static>(
    instance: &mut Instance<S, A, Q>,
    name: &str,
    args: &[&[u8]],
    result_max_length: usize,
) -> VmResult<Vec<u8>> {
    let mut arg_region_ptrs = Vec::<u32>::with_capacity(args.len());
    for arg in args {
        let region_ptr = instance.allocate(arg.len())?;
        instance.write_memory(region_ptr, arg)?;
        arg_region_ptrs.push(region_ptr);
    }

    let res_region_ptr = match args.len() {
        1 => {
            let func: Func<u32, u32> = instance.func(name)?;
            func.call(arg_region_ptrs[0])?
        }
        2 => {
            let func: Func<(u32, u32), u32> = instance.func(name)?;
            func.call(arg_region_ptrs[0], arg_region_ptrs[1])?
        }
        _ => panic!("call_raw called with unsupported number of arguments"),
    };

    let data = instance.read_memory(res_region_ptr, result_max_length)?;
    // free return value in wasm (arguments were freed in wasm code)
    instance.deallocate(res_region_ptr)?;
    Ok(data)
}
