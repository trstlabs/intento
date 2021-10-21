import { txClient, queryClient, MissingWalletError } from './module'
// @ts-ignore
import { SpVuexError } from '@starport/vuex'

import { Code } from "./module/types/compute/genesis"
import { Contract } from "./module/types/compute/genesis"
import { Sequence } from "./module/types/compute/genesis"
import { QueryContractHistoryRequest } from "./module/types/compute/query"
import { QueryContractAddressByContractIdRequest } from "./module/types/compute/query"
import { QueryContractKeyRequest } from "./module/types/compute/query"
import { QueryContractHashRequest } from "./module/types/compute/query"
import { CodeInfoResponse } from "./module/types/compute/query"
import { QueryCodesResponse } from "./module/types/compute/query"
import { QueryContractAddressByContractIdResponse } from "./module/types/compute/query"
import { QueryContractKeyResponse } from "./module/types/compute/query"
import { QueryContractHashResponse } from "./module/types/compute/query"
import { DecryptedAnswer } from "./module/types/compute/query"
import { AccessTypeParam } from "./module/types/compute/types"
import { CodeInfo } from "./module/types/compute/types"
import { ContractInfo } from "./module/types/compute/types"
import { ContractInfoWithAddress } from "./module/types/compute/types"
import { AbsoluteTxPosition } from "./module/types/compute/types"
import { Model } from "./module/types/compute/types"


export { Code, Contract, Sequence, QueryContractHistoryRequest, QueryContractAddressByContractIdRequest, QueryContractKeyRequest, QueryContractHashRequest, CodeInfoResponse, QueryCodesResponse, QueryContractAddressByContractIdResponse, QueryContractKeyResponse, QueryContractHashResponse, DecryptedAnswer, AccessTypeParam, CodeInfo, ContractInfo, ContractInfoWithAddress, AbsoluteTxPosition, Model };

async function initTxClient(vuexGetters) {
	return await txClient(vuexGetters['common/wallet/signer'], {
		addr: vuexGetters['common/env/apiTendermint']
	})
}

async function initQueryClient(vuexGetters) {
	return await queryClient({
		addr: vuexGetters['common/env/apiCosmos']
	})
}

function mergeResults(value, next_values) {
	for (let prop of Object.keys(next_values)) {
		if (Array.isArray(next_values[prop])) {
			value[prop]=[...value[prop], ...next_values[prop]]
		}else{
			value[prop]=next_values[prop]
		}
	}
	return value
}

function getStructure(template) {
	let structure = { fields: [] }
	for (const [key, value] of Object.entries(template)) {
		let field: any = {}
		field.name = key
		field.type = typeof value
		structure.fields.push(field)
	}
	return structure
}

const getDefaultState = () => {
	return {
				ContractInfo: {},
				ContractResult: {},
				ContractsByCode: {},
				SmartContractState: {},
				Code: {},
				
				_Structure: {
						Code: getStructure(Code.fromPartial({})),
						Contract: getStructure(Contract.fromPartial({})),
						Sequence: getStructure(Sequence.fromPartial({})),
						QueryContractHistoryRequest: getStructure(QueryContractHistoryRequest.fromPartial({})),
						QueryContractAddressByContractIdRequest: getStructure(QueryContractAddressByContractIdRequest.fromPartial({})),
						QueryContractKeyRequest: getStructure(QueryContractKeyRequest.fromPartial({})),
						QueryContractHashRequest: getStructure(QueryContractHashRequest.fromPartial({})),
						CodeInfoResponse: getStructure(CodeInfoResponse.fromPartial({})),
						QueryCodesResponse: getStructure(QueryCodesResponse.fromPartial({})),
						QueryContractAddressByContractIdResponse: getStructure(QueryContractAddressByContractIdResponse.fromPartial({})),
						QueryContractKeyResponse: getStructure(QueryContractKeyResponse.fromPartial({})),
						QueryContractHashResponse: getStructure(QueryContractHashResponse.fromPartial({})),
						DecryptedAnswer: getStructure(DecryptedAnswer.fromPartial({})),
						AccessTypeParam: getStructure(AccessTypeParam.fromPartial({})),
						CodeInfo: getStructure(CodeInfo.fromPartial({})),
						ContractInfo: getStructure(ContractInfo.fromPartial({})),
						ContractInfoWithAddress: getStructure(ContractInfoWithAddress.fromPartial({})),
						AbsoluteTxPosition: getStructure(AbsoluteTxPosition.fromPartial({})),
						Model: getStructure(Model.fromPartial({})),
						
		},
		_Subscriptions: new Set(),
	}
}

// initial state
const state = getDefaultState()

export default {
	namespaced: true,
	state,
	mutations: {
		RESET_STATE(state) {
			Object.assign(state, getDefaultState())
		},
		QUERY(state, { query, key, value }) {
			state[query][JSON.stringify(key)] = value
		},
		SUBSCRIBE(state, subscription) {
			state._Subscriptions.add(subscription)
		},
		UNSUBSCRIBE(state, subscription) {
			state._Subscriptions.delete(subscription)
		}
	},
	getters: {
				getContractInfo: (state) => (params = { params: {}}) => {
					if (!(<any> params).query) {
						(<any> params).query=null
					}
			return state.ContractInfo[JSON.stringify(params)] ?? {}
		},
				getContractResult: (state) => (params = { params: {}}) => {
					if (!(<any> params).query) {
						(<any> params).query=null
					}
			return state.ContractResult[JSON.stringify(params)] ?? {}
		},
				getContractsByCode: (state) => (params = { params: {}}) => {
					if (!(<any> params).query) {
						(<any> params).query=null
					}
			return state.ContractsByCode[JSON.stringify(params)] ?? {}
		},
				getSmartContractState: (state) => (params = { params: {}}) => {
					if (!(<any> params).query) {
						(<any> params).query=null
					}
			return state.SmartContractState[JSON.stringify(params)] ?? {}
		},
				getCode: (state) => (params = { params: {}}) => {
					if (!(<any> params).query) {
						(<any> params).query=null
					}
			return state.Code[JSON.stringify(params)] ?? {}
		},
				
		getTypeStructure: (state) => (type) => {
			return state._Structure[type].fields
		}
	},
	actions: {
		init({ dispatch, rootGetters }) {
			console.log('Vuex module: trst.x.compute.v1beta1 initialized!')
			if (rootGetters['common/env/client']) {
				rootGetters['common/env/client'].on('newblock', () => {
					dispatch('StoreUpdate')
				})
			}
		},
		resetState({ commit }) {
			commit('RESET_STATE')
		},
		unsubscribe({ commit }, subscription) {
			commit('UNSUBSCRIBE', subscription)
		},
		async StoreUpdate({ state, dispatch }) {
			state._Subscriptions.forEach(async (subscription) => {
				try {
					await dispatch(subscription.action, subscription.payload)
				}catch(e) {
					throw new SpVuexError('Subscriptions: ' + e.message)
				}
			})
		},
		
		
		
		 		
		
		
		async QueryContractInfo({ commit, rootGetters, getters }, { options: { subscribe, all} = { subscribe:false, all:false}, params: {...key}, query=null }) {
			try {
				const queryClient=await initQueryClient(rootGetters)
				let value= (await queryClient.queryContractInfo( key.address)).data
				
					
				commit('QUERY', { query: 'ContractInfo', key: { params: {...key}, query}, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryContractInfo', payload: { options: { all }, params: {...key},query }})
				return getters['getContractInfo']( { params: {...key}, query}) ?? {}
			} catch (e) {
				throw new SpVuexError('QueryClient:QueryContractInfo', 'API Node Unavailable. Could not perform query: ' + e.message)
				
			}
		},
		
		
		
		
		 		
		
		
		async QueryContractResult({ commit, rootGetters, getters }, { options: { subscribe, all} = { subscribe:false, all:false}, params: {...key}, query=null }) {
			try {
				const queryClient=await initQueryClient(rootGetters)
				let value= (await queryClient.queryContractResult( key.address)).data
				
					
				commit('QUERY', { query: 'ContractResult', key: { params: {...key}, query}, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryContractResult', payload: { options: { all }, params: {...key},query }})
				return getters['getContractResult']( { params: {...key}, query}) ?? {}
			} catch (e) {
				throw new SpVuexError('QueryClient:QueryContractResult', 'API Node Unavailable. Could not perform query: ' + e.message)
				
			}
		},
		
		
		
		
		 		
		
		
		async QueryContractsByCode({ commit, rootGetters, getters }, { options: { subscribe, all} = { subscribe:false, all:false}, params: {...key}, query=null }) {
			try {
				const queryClient=await initQueryClient(rootGetters)
				let value= (await queryClient.queryContractsByCode( key.code_id)).data
				
					
				commit('QUERY', { query: 'ContractsByCode', key: { params: {...key}, query}, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryContractsByCode', payload: { options: { all }, params: {...key},query }})
				return getters['getContractsByCode']( { params: {...key}, query}) ?? {}
			} catch (e) {
				throw new SpVuexError('QueryClient:QueryContractsByCode', 'API Node Unavailable. Could not perform query: ' + e.message)
				
			}
		},
		
		
		
		
		 		
		
		
		async QuerySmartContractState({ commit, rootGetters, getters }, { options: { subscribe, all} = { subscribe:false, all:false}, params: {...key}, query=null }) {
			try {
				const queryClient=await initQueryClient(rootGetters)
				let value= (await queryClient.querySmartContractState( key.address,  key.query_data)).data
				
					
				commit('QUERY', { query: 'SmartContractState', key: { params: {...key}, query}, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QuerySmartContractState', payload: { options: { all }, params: {...key},query }})
				return getters['getSmartContractState']( { params: {...key}, query}) ?? {}
			} catch (e) {
				throw new SpVuexError('QueryClient:QuerySmartContractState', 'API Node Unavailable. Could not perform query: ' + e.message)
				
			}
		},
		
		
		
		
		 		
		
		
		async QueryCode({ commit, rootGetters, getters }, { options: { subscribe, all} = { subscribe:false, all:false}, params: {...key}, query=null }) {
			try {
				const queryClient=await initQueryClient(rootGetters)
				let value= (await queryClient.queryCode( key.code_id)).data
				
					
				commit('QUERY', { query: 'Code', key: { params: {...key}, query}, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryCode', payload: { options: { all }, params: {...key},query }})
				return getters['getCode']( { params: {...key}, query}) ?? {}
			} catch (e) {
				throw new SpVuexError('QueryClient:QueryCode', 'API Node Unavailable. Could not perform query: ' + e.message)
				
			}
		},
		
		
		async sendMsgInstantiateContract({ rootGetters }, { value, fee = [], memo = '' }) {
			try {
				const txClient=await initTxClient(rootGetters)
				const msg = await txClient.msgInstantiateContract(value)
				const result = await txClient.signAndBroadcast([msg], {fee: { amount: fee, 
	gas: "200000" }, memo})
				return result
			} catch (e) {
				if (e == MissingWalletError) {
					throw new SpVuexError('TxClient:MsgInstantiateContract:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgInstantiateContract:Send', 'Could not broadcast Tx: '+ e.message)
				}
			}
		},
		async sendMsgStoreCode({ rootGetters }, { value, fee = [], memo = '' }) {
			try {
				const txClient=await initTxClient(rootGetters)
				const msg = await txClient.msgStoreCode(value)
				const result = await txClient.signAndBroadcast([msg], {fee: { amount: fee, 
	gas: "200000" }, memo})
				return result
			} catch (e) {
				if (e == MissingWalletError) {
					throw new SpVuexError('TxClient:MsgStoreCode:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgStoreCode:Send', 'Could not broadcast Tx: '+ e.message)
				}
			}
		},
		async sendMsgExecuteContract({ rootGetters }, { value, fee = [], memo = '' }) {
			try {
				const txClient=await initTxClient(rootGetters)
				const msg = await txClient.msgExecuteContract(value)
				const result = await txClient.signAndBroadcast([msg], {fee: { amount: fee, 
	gas: "200000" }, memo})
				return result
			} catch (e) {
				if (e == MissingWalletError) {
					throw new SpVuexError('TxClient:MsgExecuteContract:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgExecuteContract:Send', 'Could not broadcast Tx: '+ e.message)
				}
			}
		},
		
		async MsgInstantiateContract({ rootGetters }, { value }) {
			try {
				const txClient=await initTxClient(rootGetters)
				const msg = await txClient.msgInstantiateContract(value)
				return msg
			} catch (e) {
				if (e == MissingWalletError) {
					throw new SpVuexError('TxClient:MsgInstantiateContract:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgInstantiateContract:Create', 'Could not create message: ' + e.message)
					
				}
			}
		},
		async MsgStoreCode({ rootGetters }, { value }) {
			try {
				const txClient=await initTxClient(rootGetters)
				const msg = await txClient.msgStoreCode(value)
				return msg
			} catch (e) {
				if (e == MissingWalletError) {
					throw new SpVuexError('TxClient:MsgStoreCode:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgStoreCode:Create', 'Could not create message: ' + e.message)
					
				}
			}
		},
		async MsgExecuteContract({ rootGetters }, { value }) {
			try {
				const txClient=await initTxClient(rootGetters)
				const msg = await txClient.msgExecuteContract(value)
				return msg
			} catch (e) {
				if (e == MissingWalletError) {
					throw new SpVuexError('TxClient:MsgExecuteContract:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgExecuteContract:Create', 'Could not create message: ' + e.message)
					
				}
			}
		},
		
	}
}
