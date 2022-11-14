use std::sync::SgxRwLock;

use lazy_static::lazy_static;
use log::*;
use lru::LruCache;

use enclave_ffi_types::EnclaveError;

use enclave_cosmos_types::types::ContractCode;
use enclave_crypto::HASH_SIZE;

use super::{gas, validation};
use crate::gas::WasmCosts;

use crate::cosmwasm_config::ContractOperation;

pub struct Code {
    pub code: Vec<u8>,

}

impl Code {
    pub fn new(code: Vec<u8>) -> Self {
        Self { code }
    }
}

lazy_static! {
    static ref MODULE_CACHE: SgxRwLock<LruCache<[u8; HASH_SIZE], Code>> =
        SgxRwLock::new(LruCache::new(0));
}

pub fn configure_module_cache(cap: usize) {
    debug!("configuring module cache: {}", cap);
    MODULE_CACHE.write().unwrap().resize(cap)
}

pub fn create_module_instance(
    contract_code: &ContractCode,
    gas_costs: &WasmCosts,
    operation: ContractOperation,
) -> Result<Code, EnclaveError> {
    debug!("fetching module from cache");
    let cache = MODULE_CACHE.read().unwrap();

    // If the cache is disabled, don't try to use it and just compile the module.
    if cache.cap() == 0 {
        debug!("cache is disabled, building module");
        return analyze_module(contract_code, gas_costs, operation);
    }
    debug!("cache is enabled");

    // Try to fetch a cached instance
    let mut code = None;

    debug!("peeking in cache");
    let peek_result = cache.peek(&contract_code.hash());
    if let Some(Code {
        code: cached_code,
    }) = peek_result
    {
        debug!("found instance in cache!");
        code = Some(cached_code.clone());
    }

    drop(cache); // Release read lock

    // if we couldn't find the code in the cache, analyze it now
    if code.is_none() {
        debug!("code not found in cache! analyzing now");
        let versioned_code = analyze_module(contract_code, gas_costs, operation)?;
        code = Some(versioned_code.code);
    }

    // If we analyzed the code in the previous step, insert it to the LRU cache
    debug!("updating cache");
    let mut cache = MODULE_CACHE.write().unwrap();
    if let Some(code) = code.clone() {
        debug!("storing code in cache");
        cache.put(contract_code.hash(), Code::new(code));
    } else {
        // Touch the cache to update the LRU value
        debug!("updating LRU without storing anything");
        cache.get(&contract_code.hash());
    }

    let code = code.unwrap();

    debug!("returning built instance");
    Ok(Code::new(code))
}

pub fn analyze_module(
    contract_code: &ContractCode,
    gas_costs: &WasmCosts,
    operation: ContractOperation,
) -> Result<Code, EnclaveError> {
    let mut module = walrus::ModuleConfig::new()
        .generate_producers_section(false)
        .parse(contract_code.code())
        .map_err(|_| EnclaveError::InvalidWasm)?;

    for import in module.imports.iter() {
        trace!("import {:?}", import)
    }
    for export in module.exports.iter() {
        trace!("export {:?}", export)
    }

    use walrus::Export;
    let exports = module.exports.iter();
  
   
    drop(exports);

    validation::validate_memory(&mut module)?;

    if let ContractOperation::Init = operation {
        if module.has_floats() {
            debug!("contract was found to contain floating point operations");
            return Err(EnclaveError::WasmModuleWithFP);
        }
    }

    gas::add_metering(&mut module, gas_costs);

    let code = module.emit_wasm();

    Ok(Code::new(code))
}