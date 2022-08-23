use cosmwasm_std::{
    entry_point, log, to_binary, Api, Binary, Coin, CosmosMsg, Env, Response,
    HumanAddr, DepsMut, MessageInfo, StdError, Querier, StdResult, Storage, VoteOption,
};

use crate::msg::{ExecuteMsg, InstantiateMsg};

#[entry_point]
pub fn instantiate(
    deps: DepsMut,
    _env: Env,
    _info: MessageInfo,
    msg: InstantiateMsg,
) -> StdResult<Response> {
    Ok(Response::default())
}

#[entry_point]
pub fn execute(
    deps: DepsMut,
    env: Env,
    info: MessageInfo,
    msg: ExecuteMsg,
) -> Result<Response,StdError> {
    match msg {
        ExecuteMsg::Rewards { address } => execute_query_rewards(deps, env, address),
    }
}

pub fn execute_query_rewards(
    deps: DepsMut,
    env: Env,
    address: String,
) -> StdResult<Response> {
    let query = DistQuery::Rewards {
        delegator: address.clone(),
    };

    let mut query_rewards =
        deps.querier
            .query(&query.into())
            .unwrap_or_else(|_| RewardsResponse {
                rewards: vec![],
                total: vec![],
            });

    let active_proposal = query_rewards
        .total
        .pop()
        .unwrap_or_else(|| Coin {
            denom: "stake".to_string(),
            amount: Default::default(),
        })
        .amount
        .0 as u64;

    Ok(Response.new().set_data(Binary::from(active_proposal.to_be_bytes().to_vec())))
    
}
