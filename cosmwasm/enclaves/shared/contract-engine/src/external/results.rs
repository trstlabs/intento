use sgx_types::sgx_status_t;

use enclave_ffi_types::{
    CallbackSigResult, EnclaveError, HandleResult, InitResult, QueryResult, UntrustedVmError,
    UserSpaceBuffer,
};

use crate::external::ocalls::ocall_allocate;

/// This struct is returned from module initialization.
pub struct InitSuccess {
    /// The output of the execution
    pub output: Vec<u8>,
    /// The contract_key of this contract.
    pub contract_key: [u8; 64],
    /// The callback_sig for this contract.
    pub callback_sig: [u8; 32],
}

pub fn result_init_success_to_initresult(result: Result<InitSuccess, EnclaveError>) -> InitResult {
    match result {
        Ok(InitSuccess {
            output,
            contract_key,
            callback_sig,
        }) => {
            let user_buffer = unsafe {
                let mut user_buffer = std::mem::MaybeUninit::<UserSpaceBuffer>::uninit();
                match ocall_allocate(user_buffer.as_mut_ptr(), output.as_ptr(), output.len()) {
                    sgx_status_t::SGX_SUCCESS => { /* continue */ }
                    _ => {
                        return InitResult::Failure {
                            err: EnclaveError::FailedOcall {
                                vm_error: UntrustedVmError::default(),
                            },
                        }
                    }
                }
                user_buffer.assume_init()
            };
            InitResult::Success {
                output: user_buffer,
                contract_key,
                callback_sig,
            }
        }
        Err(err) => InitResult::Failure { err },
    }
}

/// This struct is returned from a handle method.
pub struct HandleSuccess {
    /// The output of the execution
    pub output: Vec<u8>,
}

pub fn result_handle_success_to_handleresult(
    result: Result<HandleSuccess, EnclaveError>,
) -> HandleResult {
    match result {
        Ok(HandleSuccess { output }) => {
            let user_buffer = unsafe {
                let mut user_buffer = std::mem::MaybeUninit::<UserSpaceBuffer>::uninit();
                match ocall_allocate(user_buffer.as_mut_ptr(), output.as_ptr(), output.len()) {
                    sgx_status_t::SGX_SUCCESS => { /* continue */ }
                    _ => {
                        return HandleResult::Failure {
                            err: EnclaveError::FailedOcall {
                                vm_error: UntrustedVmError::default(),
                            },
                        }
                    }
                }
                user_buffer.assume_init()
            };
            HandleResult::Success {
                output: user_buffer,
            }
        }
        Err(err) => HandleResult::Failure { err },
    }
}

/// This struct is returned from a query method.
pub struct QuerySuccess {
    /// The output of the query
    pub output: Vec<u8>,
}

pub fn result_query_success_to_queryresult(
    result: Result<QuerySuccess, EnclaveError>,
) -> QueryResult {
    match result {
        Ok(QuerySuccess { output }) => {
            let user_buffer = unsafe {
                let mut user_buffer = std::mem::MaybeUninit::<UserSpaceBuffer>::uninit();
                match ocall_allocate(user_buffer.as_mut_ptr(), output.as_ptr(), output.len()) {
                    sgx_status_t::SGX_SUCCESS => { /* continue */ }
                    _ => {
                        return QueryResult::Failure {
                            err: EnclaveError::FailedOcall {
                                vm_error: UntrustedVmError::default(),
                            },
                        }
                    }
                }
                user_buffer.assume_init()
            };
            QueryResult::Success {
                output: user_buffer,
            }
        }
        Err(err) => QueryResult::Failure { err },
    }
}

/// This struct is returned from a create callback method.
pub struct CallbackSigSuccess {
    /// The output of the callback sig creation
    // pub output: Vec<u8>,
    pub callback_sig: [u8; 32],
    pub encrypted_msg: Vec<u8>,
}

pub fn result_callback_sig_success_to_callbackresult(
    result: Result<CallbackSigSuccess, EnclaveError>,
) -> CallbackSigResult {
    match result {
        Ok(CallbackSigSuccess {
            callback_sig,
            encrypted_msg,
        }) => {
            let user_buffer = unsafe {
                let mut user_buffer = std::mem::MaybeUninit::<UserSpaceBuffer>::uninit();
                match ocall_allocate(
                    user_buffer.as_mut_ptr(),
                    encrypted_msg.as_ptr(),
                    encrypted_msg.len(),
                ) {
                    sgx_status_t::SGX_SUCCESS => { /* continue */ }
                    _ => {
                        return CallbackSigResult::Failure {
                            err: EnclaveError::FailedOcall {
                                vm_error: UntrustedVmError::default(),
                            },
                        }
                    }
                }
                user_buffer.assume_init()
            };
            CallbackSigResult::Success {
                callback_sig,
                encrypted_msg: user_buffer,
            }
        }
        Err(err) => CallbackSigResult::Failure { err },
    }
}
