use cosmwasm_std::{
    log, to_binary, Api, CanonicalAddr, Env, Extern, HandleResponse, HandleResult, HumanAddr,
    InitResponse, InitResult, Querier, QueryResult, StdError, Storage, Uint128,
};

use std::collections::HashSet;

use serde_json_wasm as serde_json;

use secret_toolkit::utils::{pad_handle_result};

use crate::msg::{
    HandleAnswer, HandleMsg, InitMsg, ResponseStatus,
    ResponseStatus::{Failure, Success},
   // Token,
};
use crate::state::{load, may_load, remove, save, State};

use chrono::NaiveDateTime;

/// storage key for estimationperiod state
pub const CONFIG_KEY: &[u8] = b"config";

/// pad handle responses and log attributes to blocks of 256 bytes to prevent leaking info based on
/// response size
pub const BLOCK_SIZE: usize = 256;

////////////////////////////////////// Init ///////////////////////////////////////
/// Returns InitResult
///
/// Initializes the estimationperiod state and registers Receive function with sell and estimation
/// token contracts
///
/// # Arguments
///
/// * `deps` - mutable reference to Extern containing all the contract's external dependencies
/// * `env` - Env of contract's environment
/// * `msg` - InitMsg passed in with the instantiation message
pub fn init<S: Storage, A: Api, Q: Querier>(
    deps: &mut Extern<S, A, Q>,
    env: Env,
    msg: InitMsg,
) -> InitResult {
    if msg.sell_amount == Uint128(0) {
        return Err(StdError::generic_err("Sell amount must be greater than 0"));
    }
    if msg.sell_contract.address == msg.estimation_contract.address {
        return Err(StdError::generic_err(
            "Sell contract and estimation contract must be different",
        ));
    }
    let state = State {
        estimationperiod_addr: env.contract.address,
        seller: env.message.sender,

      sell_contract: msg.sell_contract,
      estimation_contract: msg.estimation_contract,
       // sell_amount: msg.sell_amount.u128(),
       // minimum_estimation: msg.minimum_estimation.u128(),
       estimationcount: msg.estimationcount,
        estimators:  vec![],
        estimation_list:  vec![],
        is_completed: false,
       // tokens_consigned: false,
       // description: msg.description,
        estimation_price: 0,
       // lowestesitmator: 
    };

    save(&mut deps.storage, CONFIG_KEY, &state)?;
  //  Ok(Response::default())

    // register receive with the estimation/sell token contracts
    Ok(InitResponse {
        messages: vec![
            state
                .sell_contract
                .register_receive_msg(env.contract_code_hash.clone())?,
            state
                .estimation_contract
                .register_receive_msg(env.contract_code_hash)?,
        ],
        log: vec![],
    })
}

///////////////////////////////////// Handle //////////////////////////////////////
/// Returns HandleResult
///
/// # Arguments
///
/// * `deps` - mutable reference to Extern containing all the contract's external dependencies
/// * `env` - Env of contract's environment
/// * `msg` - HandleMsg passed in with the execute message
pub fn handle<S: Storage, A: Api, Q: Querier>(
    deps: &mut Extern<S, A, Q>,
    env: Env,
    msg: HandleMsg,
) -> HandleResult {
    let response = match msg {
        HandleMsg::RetractEstimation { .. } => try_retract(deps, env.message.sender),
        HandleMsg::RevealEstimation { .. } => try_finalize(deps, env),
        HandleMsg::ReturnAll { .. } => try_finalize(deps, env),
        HandleMsg::CreateEstimation { from, amount, .. } => try_create_estimation(deps, env, from, amount),
    //    HandleMsg::ViewEstimation { .. } => try_view_estimation(deps, &env.message.sender),
    };
    pad_handle_result(response, BLOCK_SIZE)
}

/// Returns HandleResult
///
/// process the Receive message sent after either estimation or sell token contract sent tokens to
/// estimationperiod escrow
///
/// # Arguments
///
/// * `deps` - mutable reference to Extern containing all the contract's external dependencies
/// * `env` - Env of contract's environment
/// * `from` - address of owner of tokens sent to escrow
/// * `amount` - Uint128 amount sent to escrow
fn try_create_estimation<S: Storage, A: Api, Q: Querier>(
    deps: &mut Extern<S, A, Q>,
    env: Env,
    from: HumanAddr,
    amount: Uint128,
) -> HandleResult {
    let mut state: State = load(&deps.storage, CONFIG_KEY)?;

    if env.message.sender != state.seller && state.estimators.iter().len()  < state.estimationcount {
        try_estimate(deps, env, from, amount, &mut state)
    } else {
        let message = format!(
            "Address: {} is from seller or estimation count is reached",
            env.message.sender
        );
        let resp = serde_json::to_string(&HandleAnswer::Status {
            status: Failure,
            message,
        })
        .unwrap();

        return Ok(HandleResponse {
            messages: vec![],
            log: vec![log("response", resp)],
            data: None,
        });
    }
}

/// Returns HandleResult
///
/// process the estimation attempt
///
/// # Arguments
///
/// * `deps` - mutable reference to Extern containing all the contract's external dependencies
/// * `env` - Env of contract's environment
/// * `estimator` - address of owner of tokens sent to escrow
/// * `amount` - Uint128 amount sent to escrow
/// * `state` - mutable reference to estimationperiod state
fn try_estimate<S: Storage, A: Api, Q: Querier>(
    deps: &mut Extern<S, A, Q>,
    env: Env,
    estimator: HumanAddr,
    amount: Uint128,
    state: &mut State,
) -> HandleResult {
    // if estimationperiod is over, send the tokens back
    if state.is_completed {
        let message = String::from("Estimation period has ended. Deposit tokens have been returned");

        let resp = serde_json::to_string(&HandleAnswer::Estimation {
            status: Failure,
            message,
         //   previous_estimation: None,
            amount_estimation: None,
         
        })
        .unwrap();

        return Ok(HandleResponse {
            messages: vec![state.estimation_contract.transfer_msg(estimator, amount)?],
            log: vec![log("response", resp)],
            data: None,
        });
    }
    // don't accept a 0 estimation
    else if amount == Uint128(0) {
        let message = String::from("Estimation must be greater than 0");

        let resp = serde_json::to_string(&HandleAnswer::Estimation {
            status: Failure,
            message,
           // previous_estimation: None,
           amount_estimation: None,
          
        })
        .unwrap();

        return Ok(HandleResponse {
            messages: vec![],
            log: vec![log("response", resp)],
            data: None,
        });
    }
   
  //  let mut return_amount: Option<Uint128> = None;
    //let estimator_raw = &deps.api.canonical_address(&estimator)?;

    // if there is no active estimation from this address
    if !state.estimators.contains(&estimator.to_string()) {
      //  let estimation: Option<Estimation> = may_load(&deps.storage, estimator_raw.as_slice())?;
      
        // insert in list of estimators and save
        state.estimators.push(estimator.to_string());
        state.estimation_list.push(amount.u128());
        save(&mut deps.storage, CONFIG_KEY, &state)?;
    }
    
  //  save(&mut deps.storage, estimator_raw.as_slice())?;

    let message = String::from("estimation accepted");
    let cos_msg = Vec::new();

    let resp = serde_json::to_string(&HandleAnswer::Estimation {
        status: Success,
        message,
     //   previous_estimation: None,
        amount_estimation: Some(amount),
      //  amount_returned: return_amount,
    })
    .unwrap();

    return Ok(HandleResponse {
        messages: cos_msg,
        log: vec![log("response", resp)],
        data: None,
    })
}

/// Returns HandleResult
///
/// attempt to retract current estimation
///
/// # Arguments
///
/// * `deps` - mutable reference to Extern containing all the contract's external dependencies
/// * `estimator` - address of estimator
fn try_retract<S: Storage, A: Api, Q: Querier>(
    deps: &mut Extern<S, A, Q>,
    estimator: HumanAddr,
) -> HandleResult {
    let mut state: State = load(&deps.storage, CONFIG_KEY)?;

    //let estimator_raw = &deps.api.canonical_address(&estimator)?;
    let mut cos_msg = Vec::new();
   // let sent: Option<Uint128>;
    let mut log_msg = String::new();
    let status: ResponseStatus;
   // let mut index_to_remove: usize;
   let mut index_to_remove = Option::None;
    //let index = state.estimators.iter().position(|&x| x == estimator.clone().to_string()).unwrap();

    for (index, element) in state.estimators.iter().enumerate() {
        if element.to_owned() == estimator.clone().to_string() {
            index_to_remove = Some(index);
            break;
        }
    }

    match index_to_remove {
        None => {
            status = Failure;
           // sent = None;
            log_msg.push_str(&format!("No active estimation for address: {}", estimator))
    
        
    }  // Don't do any
        Some(i) => {
    // if there was a active estimation from this address, remove the estimation and return tokens
   // if state.estimators.contains(&estimator.to_string())   {


       let estimation: u128;

    
       estimation = state.estimation_list[i];
        //let estimation: Option<Estimation> = may_load(&deps.storage, estimator_raw.as_slice())?;
       // if let Some(old_estimation) = estimation {
           // remove(&mut deps.storage, estimator.to_string());
            state.estimators.remove(i);
            state.estimation_list.remove(i);
            save(&mut deps.storage, CONFIG_KEY, &state)?;
            cos_msg.push(
                state
                    .estimation_contract
                    .transfer_msg(estimator, Uint128(estimation))?,
            );
            status = Success;
           // sent = Some(Uint128(estimation));
            log_msg.push_str("Estimation retracted.  ");
       // } else {
         //   status = Failure;
           // sent = None;
            //log_msg.push_str(&format!("No active estimation for address: {}", estimator));
   //     }
    // no active estimation found
    }}  return Ok(HandleResponse {
    messages: cos_msg,
    log: vec![],
    data: Some(to_binary(&HandleAnswer::RetractEstimation {
        status,
        message: log_msg,
       // amount_returned: sent,
    })?),
});  }

/// Returns HandleResult
///
/// closes the estimationperiod and sends all the tokens in escrow to where they belong
///
/// # Arguments
///
/// * `deps` - mutable reference to Extern containing all the contract's external dependencies
/// * `env` - Env of contract's environment
/// * `only_if_estimations` - true if estimationperiod should stay open if there are no estimations
/// * `return_all` - true if being called from the return_all fallback plan
fn try_finalize<S: Storage, A: Api, Q: Querier>(
    deps: &mut Extern<S, A, Q>,
    env: Env,
 
   //return_all: bool,
) -> HandleResult {
    let mut state: State = load(&deps.storage, CONFIG_KEY)?;
 let best_estimation: Option<Uint128> = None;
    // can only do a return_all if the estimationperiod is closed
    if  !state.is_completed {
        return Ok(HandleResponse {
            messages: vec![],
            log: vec![],
            data: Some(to_binary(&HandleAnswer::RevealEstimation {
                status: Failure,
                message: String::from(
                    "return_all can only be executed after the estimation period has ended",
                ),
                best_estimation: None,
               // amount_returned: None,
            })?),
        });
    }
    // if not the estimationperiod owner, can't finalize, but you can return_all
    if env.message.sender != state.seller {
        return Ok(HandleResponse {
            messages: vec![],
            log: vec![],
            data: Some(to_binary(&HandleAnswer::RevealEstimation {
                status: Failure,
                message: String::from("Only creator can finalize the estimation"),
                best_estimation:  None,
                //amount_returned: None,
            })?),
        });
    };
    // if there are no active estimations, and owner only wants to close if estimations
    if !state.is_completed  && state.estimators.iter().len() > state.estimationcount {
        return Ok(HandleResponse {
            messages: vec![],
            log: vec![],
            data: Some(to_binary(&HandleAnswer::RevealEstimation {
                status: Failure,
                message: String::from("Did not close because there are no active estimations"),
                best_estimation:  None,
              //  amount_returned: None,
            })?),
        });
    };

    state.estimation_list.sort();
    let mid = state.estimationcount + 1 / 2;
    let mut cos_msg = Vec::new();


  //  for (index, estimator) in &state.estimators {

    //}




    /*
    let mut cos_msg = Vec::new();
    let mut update_state = false;
    let mut winning_amount: Option<Uint128> = None;
    let mut amount_returned: Option<Uint128> = None;

    let no_estimations = state.estimators.is_empty();
   
        // load all the estimations
        struct OwnedEstimation {
            pub estimator: CanonicalAddr,
            pub estimation: Estimation,
        }
        let mut estimation_list: Vec<OwnedEstimation> = Vec::new();
        for estimator in &state.estimators {
            let estimation: Option<Estimation> = may_load(&deps.storage, estimator.as_slice())?;
            if let Some(found_estimation) = estimation {
                estimation_list.push(OwnedEstimation {
                    estimator: CanonicalAddr::from(estimator.as_slice()),
                    estimation: found_estimation,
                });
            }
        }
        // closing an estimationperiod that has been fully consigned
        if state.tokens_consigned && !state.is_completed {
            estimation_list.sort_by(|a, b| {
                a.estimation
                    .amount
                    .cmp(&b.estimation.amount)
                    .then(b.estimation.timestamp.cmp(&a.estimation.timestamp))
            });
            // if there was a winner, swap the tokens
            if let Some(best_estimation) = estimation_list.pop() {
                cos_msg.push(
                    state
                        .estimation_contract
                        .transfer_msg(state.seller.clone(), Uint128(best_estimation.estimation.amount))?,
                );
                cos_msg.push(state.sell_contract.transfer_msg(
                    deps.api.human_address(&best_estimation.estimator)?,
                    Uint128(state.sell_amount),
                )?);
                state.currently_consigned = 0;
                update_state = true;
                winning_amount = Some(Uint128(best_estimation.estimation.amount));
                state.best_estimation = best_estimation.estimation.amount;
                remove(&mut deps.storage, &best_estimation.estimator.as_slice());
                state
                    .estimators
                    .remove(&best_estimation.estimator.as_slice().to_vec());
            }
        }
        // loops through all remaining estimations to return them to the estimators
        for losing_estimation in &estimation_list {
            cos_msg.push(state.estimation_contract.transfer_msg(
                deps.api.human_address(&losing_estimation.estimator)?,
                Uint128(losing_estimation.estimation.amount),
            )?);
            remove(&mut deps.storage, &losing_estimation.estimator.as_slice());
            update_state = true;
            state.estimators.remove(&losing_estimation.estimator.as_slice().to_vec());
        }

    // return any tokens that have been consigned to the estimationperiod owner (can happen if owner
    // finalized the estimationperiod before consigning the full sale amount or if there were no estimations)
    if state.currently_consigned > 0 {
        cos_msg.push(
            state
                .sell_contract
                .transfer_msg(state.seller.clone(), Uint128(state.currently_consigned))?,
        );
        if !return_all {
            amount_returned = Some(Uint128(state.currently_consigned));
        }
        state.currently_consigned = 0;
        update_state = true;
    }*/

    cos_msg.push( state
        .sell_contract   
        .transfer_msg(
       state.seller.clone(), 
      
       Uint128(state.estimation_list[mid]))?);


    // mark that estimationperiod had ended
    if !state.is_completed {
        state.is_completed = true;
       // update_state = true;
    }
   // if update_state {
        save(&mut deps.storage, CONFIG_KEY, &state)?;
   // }

    let log_msg = /*if winning_amount.is_some() {
        "Sale finalized.  You have been sent the winning estimation tokens".to_string()
    } else if amount_returned.is_some() {
        let cause = if !state.tokens_consigned {
            " because you did not consign the full sale amount"
        } else if no_estimations {
            " because there were no active estimations"
        } else {
            ""
        };
        format!(
            "EstimationPeriod closed.  You have been returned the consigned tokens{}",
            cause
        )
    } else if return_all {
        "Outstanding funds have been returned".to_string()
    } else */{
        "Estimation Period has been closed".to_string()
    };
    return Ok(HandleResponse {
        messages: cos_msg,
        log: vec![],
        data: Some(to_binary(&HandleAnswer::RevealEstimation {
            status: Success,
            message: log_msg,
            best_estimation: Some(Uint128(state.estimation_list[mid])),
            //amount_returned,
        })?),
    });
}
/*
/////////////////////////////////////// Query /////////////////////////////////////
/// Returns QueryResult
///
/// # Arguments
///
/// * `deps` - reference to Extern containing all the contract's external dependencies
/// * `msg` - QueryMsg passed in with the query call
pub fn query<S: Storage, A: Api, Q: Querier>(deps: &Extern<S, A, Q>, msg: QueryMsg) -> QueryResult {
    let response = match msg {
        QueryMsg::EstimationPeriodInfo { .. } => try_query_info(deps),
    };
    pad_query_result(response, BLOCK_SIZE)
}

/// Returns QueryResult
///
/// # Arguments
///
/// * `deps` - reference to Extern containing all the contract's external dependencies
fn try_query_info<S: Storage, A: Api, Q: Querier>(deps: &Extern<S, A, Q>) -> QueryResult {
    let state: State = load(&deps.storage, CONFIG_KEY)?;

    // get sell token info
    let sell_token_info = state.sell_contract.token_info_query(&deps.querier)?;
    // get estimation token info
    let estimation_token_info = state.estimation_contract.token_info_query(&deps.querier)?;

    // build status string
    let status = if state.is_completed {
        let locked = if !state.estimators.is_empty() || state.currently_consigned > 0 {
            ", but found outstanding balances.  Please run either retract_estimation to \
                retrieve your non-winning estimation, or return_all to return all outstanding estimations/\
                consignment."
        } else {
            ""
        };
        format!("Closed{}", locked)
    } else {
        let consign = if !state.tokens_consigned { " NOT" } else { "" };
        format!(
            "Accepting estimations: Token(s) to be sold have{} been consigned to the estimationperiod",
            consign
        )
    };

    let best_estimation = if state.best_estimation == 0 {
        None
    } else {
        Some(Uint128(state.best_estimation))
    };

    to_binary(&QueryAnswer::EstimationPeriodInfo {
        sell_token: Token {
            contract_address: state.sell_contract.address,
            token_info: sell_token_info,
        },
        estimation_token: Token {
            contract_address: state.estimation_contract.address,
            token_info: estimation_token_info,
        },
        sell_amount: Uint128(state.sell_amount),
        minimum_estimation: Uint128(state.minimum_estimation),
        description: state.description,
        estimationperiod_address: state.estimationperiod_addr,
        status,
        best_estimation,
    })
}
*/