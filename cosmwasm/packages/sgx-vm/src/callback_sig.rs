//use sgx_types::{sgx_status_t, SgxResult};
//use crate::errors::VmResult;
use std::mem::MaybeUninit;
use enclave_ffi_types::{CallbackSigResult, UserSpaceBuffer};
//use std::{self};
//use log::*;
use sgx_types::*;

use crate::enclave::ENCLAVE_DOORBELL;
//use crate::errors::{EnclaveError, VmResult};



extern "C" {
    pub fn ecall_create_callback_sig(
        eid: sgx_enclave_id_t,
        retval: *mut CallbackSigResult,
        msg: *const u8,
        msg_len: u32,
        msg_info: *const u8,
        msg_info_len: u32,

    ) -> sgx_status_t;
}


/// create_callback_sig_raw creates callback sig for the given message given the community pool address that is hardcoded
pub fn create_callback_sig_raw(
    msg: &[u8],
    msg_info: &[u8]

) -> SgxResult<Vec<u8>> {
   // Bind the token to a local variable to ensure its
        // destructor runs in the end of the function
        let enclave_access_token = ENCLAVE_DOORBELL
        .get_access(false) // This can never be recursive
        .ok_or(sgx_status_t::SGX_ERROR_BUSY)?;
    let enclave = (*enclave_access_token)?;
    let eid = enclave.geteid();
    let mut callback_sig_result = MaybeUninit::<CallbackSigResult>::uninit();
    //let mut retval = CallbackSigResult;
  //  let mut callback_sig = [0u8; 32];
    let status = unsafe {
        ecall_create_callback_sig(
            eid,
            callback_sig_result.as_mut_ptr(),
            msg.as_ptr(),
            msg.len() as u32,
            msg_info.as_ptr(),
            msg_info.len() as u32,

        )
    };

  /*  if status != sgx_status_t::SGX_SUCCESS {
        return Err(status);
    }
   /* if retval != CallbackSigResult::Success {
        return Ok(Err(retval));
    }*/
    if callback_sig.is_empty() {
        error!("Got empty callback sig from encryption");
        return Err(sgx_status_t::SGX_ERROR_UNEXPECTED);
    }
    Ok(Ok(callback_sig))*/

    match status {
        sgx_status_t::SGX_SUCCESS => {
            let callback_sig_result = unsafe { callback_sig_result.assume_init() };
           callback_sig_result_to_sgx_result(callback_sig_result)
        },
        _ => { Err(sgx_status_t::SGX_ERROR_UNEXPECTED)}
    }
}

pub fn callback_sig_result_to_sgx_result(other: CallbackSigResult) -> SgxResult<Vec<u8>> {
    match other {
        CallbackSigResult::Success {
            callback_sig,
            encrypted_msg,
        } => { let mut callback_sig_vec = callback_sig.to_vec();
            //let encrypted_msg_vec = encrypted_msg.to_vec();
            let encrypted_msg_vec =  unsafe { recover_buffer(encrypted_msg) }.unwrap_or_else(Vec::new);
            callback_sig_vec.extend_from_slice(&encrypted_msg_vec);
            Ok(callback_sig_vec)
        },
        CallbackSigResult::Failure { .. } => Ok(Vec::new())
    }
}


/// Take a pointer as returned by `ocall_allocate` and recover the Vec<u8> inside of it.
pub unsafe fn recover_buffer(ptr: UserSpaceBuffer) -> Option<Vec<u8>> {
    if ptr.ptr.is_null() {
        return None;
    }
    let boxed_vector = Box::from_raw(ptr.ptr as *mut Vec<u8>);
    Some(*boxed_vector)
}