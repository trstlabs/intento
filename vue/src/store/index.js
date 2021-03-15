
import Vue from "vue";
import Vuex from "vuex";
import axios from "axios";
import app from "./app.js";


//import cosmos from "@tendermint/vue/src/store/cosmos.js";
//import { assert } from "@cosmjs/utils";




Vue.use(Vuex);

const API = process.env.VUE_APP_API
//const RPC = process.env.VUE_APP_RPC
//const API = "http://localhost:1317";
//const API = "https://node.trustpriceprotocol.com"
//const API = "http://localhost:1317"
//const ADDRESS_PREFIX = 'cosmos';
const PATH = process.env.VUE_APP_PATH

//const RPC = 'https://cli.trustpriceprotocol.com'
//const RPC = 'http://localhost:26657'

export default new Vuex.Store({
  state: {
    app,
    account: {},
    chain_id: "",
    data: {},
    client: null,
    wallet: null,
    errorMessage: "",
    newitemID: {},
    bankBalances: [],

    creatorItemList: [],
    estimatorItemList: [],
    buyerItemList: [],
    InterestedItemList: [],
    buyItemList: [],
    toEstimateList: [],
    creatorActionList: [],
    tagList: [],
    sellerList: [],
  },

  mutations: {
    set(state, { key, value }) {
      state[key] = value
    },
    accountUpdate(state, { account }) {
      state.account = account;
    },
    chainIdSet(state, { chain_id }) {
      state.chain_id = chain_id;
    },
    entitySet(state, { type, body }) {
      const updated = {};
      updated[type] = body;
      state.data = { ...state.data, ...updated };
    },
  

    setCreatorItemList(state, payload) {
      //state.CreatorItemList.push(payload);
      state.creatorItemList = payload;
    },
    setEstimatorItemList(state, payload) {
      state.estimatorItemList = payload;
    },
    setBuyerItemList(state, payload) {
      state.buyerItemList = payload;
    },
    setBuyItemList(state, payload) {
      state.buyItemList = payload;
    },

    updateBuyItemList(state, payload) {
      //console.log(payload);
      state.buyItemList = payload;
    },

    setInterestedItemList(state, payload) {
      state.InterestedItemList = payload;
    },

    setToEstimateList(state, payload) {
      state.toEstimateList = payload;
    },

    setCreatorActionList(state, payload) {
      state.creatorActionList = payload;
    },
    setTagList(state, payload) {
      state.tagList = payload;
    },
    setSellerItemList(state, payload) {
      state.sellerList = payload;
    },
  },
  actions: {
    async init({ dispatch, state }) {
      await dispatch("chainIdFetch");
      const type = { type: "item" };
      await dispatch("entityFetch", type )
      //let { type } = state.app.types.find( type  => { type == "item"})
      
     /* state.app.types.forEach(({ type }) => {
        if (type == "estimator/flag" || type == "item/reveal" || type == "item/transferable" || type == "item/transfer" || type == "item/shipping") { return; }
        dispatch("entityFetch", { type });
      });*/

      //await dispatch('accountSignInTry');

    },

    async chainIdFetch({ commit }) {
      const node_info = (await axios.get(`${API}/node_info`)).data.node_info;
      commit("chainIdSet", { chain_id: node_info.network });
    },
    /*async accountSignInTry({ state, dispatch }) {
      const mnemonic = localStorage.getItem('mnemonic')
      const privkey = localStorage.getItem('privkey')
      if (mnemonic) {
        await dispatch('accountSignIn', { mnemonic })}
        if (privkey) {
          await dispatch('torusSignIn',  privkey )}
        if (mnemonic || privkey) {
        let type = { type: "estimator" };
       await dispatch("entityFetch", type )
       let buyer = { type: "buyer" };
      await dispatch("entityFetch", buyer )
        await dispatch("setEstimatorItemList", state.account.address);
        await dispatch("setToEstimateList", state.account.address);
        await dispatch("setCreatorActionList", state.account.address);
        await dispatch("setSortedTagList");
        await dispatch("setCreatorItemList");
        await dispatch("setBuyerItemList", state.account.address);

        await dispatch("setInterestedItemList", state.account.address);
        //$emit('signedIn');
      }
    },

*/
    async bankBalancesGet({ commit, state },) {
      //const API = rootState.cosmos.env.env.API
      const { address } = state.account
      //console.log({ address })
      const url = `${API}/bank/balances/${address}`
      const value = (await axios.get(url)).data.result
      //console.log(value)
      commit('set', { key: 'bankBalances', value: value })
    },


    async accountSignOut() {
      localStorage.removeItem('mnemonic')
      localStorage.removeItem('privkey')
      window.location.reload()
    },

    async entityFetch({ commit }, { type }) {
      //const { chain_id } = state;
      const url = `${API}/${PATH.replace(/\./g, '/')}/${type}`;
      const body = (await axios.get(url)).data
      const uppercase = type.charAt(0).toUpperCase() + type.slice(1)
      if (body && body[uppercase]) {
        commit('entitySet', { type, body: body[uppercase] })
      }


    },


    async accountUpdate({ state, commit }) {
      const url = `${API}/auth/accounts/${state.account.address}`;
      const acc = (await axios.get(url)).data;
      const account = acc.result.value;
      commit("accountUpdate", { account });
    },




    ////__________________//// ////__________________//// ////__________________//// ////__________________//// ////__________________////



    async setCreatorItemList({ commit, state }, input) {
     if (!!input) { const rs = state.data.item.filter(item => item.creator === input
      ) || [];
      commit("setCreatorItemList", rs);
     }
    },

    async setEstimatorItemList({ commit, state }, input) {
      if (!!input) { 
      const rse = state.data.estimator.filter(estimator => estimator.estimator === input
      ) || [];


      commit("setEstimatorItemList", rse);}
    },

    async setBuyerItemList({ commit, state }, input) {
      if (!!input) {  const rs = state.data.item.filter(item => item.buyer === input
      )
      commit("setBuyerItemList", rs);}
    },

    async setBuyItemList({ commit, state }) {
  
      const rs = state.data.item.filter(item => !item.buyer && item.transferable === true
      ) || [];

      commit("setBuyItemList", rs);
    },

    async setLocalBuyItemList({ commit, state }) {
       const rs = state.data.item.filter(item => !item.buyer && item.transferable === true && item.localpickup === true
      );
      commit("setBuyItemList", rs);
    },
    async updateBuyItemList({ commit, state }, input) {
      if (!!input) {  const rs = state.data.item.filter(item => !item.buyer && item.transferable === true && item.title.toLowerCase().includes(input)
      );
      commit("updateBuyItemList", rs);}
    },

    async filterBuyItemList({ commit }, input) {

      commit("updateBuyItemList", input);
    },

    async tagBuyItemList({ commit, state }, input) {
      if (!!input) { 
      const rs = state.data.item.filter(item => !item.buyer && item.tags.find(tags => tags.includes(input)) && item.transferable === true)
        ;

      commit("updateBuyItemList", rs);}
    },

    async tagToEstimateList({ commit, state }, input) {
      if (!!input) {  const A = state.data.item.filter(item => !item.buyer && item.tags.find(tags => tags.includes(input)) && item.transferable === false)
        ;
      const B = state.estimatorItemList;

      const rs = A.filter(a => !B.map(b => b.itemid).includes(a.id));
      //console.log(input);
      //console.log(rs);
      //console.log(A);
      //console.log(B);
      commit("setToEstimateList", rs);}
    },

    async setSortedTagList({ commit, state }) {
      const rs = state.data.item.map(item => item.tags);
      //console.log("TEST", rs);
      /*var filtered = rs.filter(function (el) {
        return el != null;
      });
      var merged = [].concat.apply([], filtered);*/
      //console.log("TEST",merged);
      var merged = [].concat.apply([], rs);
      var frequency = {};
      merged.forEach(function (value) { frequency[value.toLowerCase()] = 0; });

      var uniques = merged.filter(function (value) {
        return ++frequency[value] == 1;
      });

      var sorted = uniques.sort(function (a, b) {
        return frequency[b] - frequency[a];
      });


      //console.log("TEST",sorted);

      commit("setTagList", sorted);
    },

    async setToEstimateList({ commit, state }) {
      const A = state.data.item.filter(item => !item.bestestimator && item.status == '');
      const B = state.estimatorItemList;
      //const rsEIL = state.data.estimator.filter(estimator => estimator.estimator === state.client.anyValidAddress);
      console.log(A);
      console.log(B);
      //console.log(A.filter(a => !B.map(b=>b.id).includes(a.id)));
      const rs = A.filter(a => !B.map(b => b.itemid).includes(a.id));

      commit("setToEstimateList", rs);


    },

    async setInterestedItemList({ commit, state }, input) {
      if (!!input) {  const rs = state.data.estimator.filter(estimator => estimator.estimator === input && estimator.interested
      );


      commit("setInterestedItemList", rs);}
    },
    async setCreatorActionList({ commit, state }, input) {
      if (!!input) { 

      const toAccept = state.data.item.filter(item => item.creator == input && item.estimationprice > 0 && !item.buyer && !item.transferable
      );
      const toShip = state.data.item.filter(item => !item.buyer && item.creator === input && !item.localpickup && !item.tracking
      );
      //console.log(input);
      //console.log(state.account.address);
      //console.log(state.client.senderAddress);



      //console.log(toAccept);
      //console.log(toShip);
      toAccept.concat(toShip);
      //console.log(toAccept);

      commit("setCreatorActionList", toAccept);}
    },

    async setSellerItemList({ commit, state }, input) {
      if (!!input) { 
      const rs = state.data.item.filter(item => item.creator === input) || [];

      commit("setSellerItemList", rs);}
    },
  },
  getters: {
    account: state => state.account, bankBalances: state => state.bankBalances, getCreatorItemList: state => state.creatorItemList, getEstimatorItemList: state => state.estimatorItemList, getBuyerItemList: state => state.buyerItemList, getBuyItemList: state => state.buyItemList, getInterestedItemList: state => state.InterestedItemList, getItemByID: state => id => state.data.item.find((item) => item.id === id), getToEstimateList: state => state.toEstimateList, getCreatorActionList: state => state.creatorActionList, getTagList: state => state.tagList, getSellerList: state => state.sellerList,

  }


});
