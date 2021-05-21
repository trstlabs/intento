
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
    sellerItemList: [],
    creatorItemList: [],
    estimatorItemList: [],
    buyerItemList: [],
    InterestedItemList: [],
    buyItemList: [],
    toEstimateList: [],
    sellerActionList: [],
    tagList: [],
    buySellerList: [],
    locationList: [],
    regionList: [],
    user: null,
    sentTransactions: {},
    receivedTransactions: {}
    //user: { uid: "B1Xk6qliE2ceNJN6HsoCk2MQO2K2"},
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
  

    setSellerItemList(state, payload) {
      //state.SellerItemList.push(payload);
      state.sellerItemList = payload;
    },
    setCreatorItemList(state, payload) {
    
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

    setSellerActionList(state, payload) {
      state.sellerActionList = payload;
    },
    setTagList(state, payload) {
      state.tagList = payload;
    },
    setBuySellerItemList(state, payload) {
      state.buySellerList = payload;
    },
  },
  actions: {
    async init({ dispatch, state }) {
      await dispatch("chainIdFetch");
      const type = { type: "item" };
      await dispatch("entityFetch", type )

      //await axios.get(process.env.VUE_APP_API + '/cosmos/tx/v1beta1/txs?events=transfer.sender%3D%27' + address + '%27').data;
    
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



    async setSellerItemList({ commit, state }, input) {
    
     if (!!input) { const rs = state.data.item.filter(item => item.seller === input
      ) || [];  
     // console.log("LIST")
    // console.log(rs)
      commit("setSellerItemList", rs);
     }
    },

    async setCreatorItemList({ commit, state }, input) {
      if (!!input) { const rs = state.data.item.filter(item => item.creator === input
       ) || [];
      // console.log(rs)
       commit("setCreatorItemList", rs);
      }
     },

   
    async setEstimatorItemList({ commit, state }, input) {
      
      if (input) { 
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
       const rs = state.data.item.filter(item => !item.buyer && item.transferable === true && item.localpickup != ""
      );
      
      commit("setBuyItemList", rs);
    },
    async updateBuyItemList({ commit, state }, input) {
      if (input != "") {  let rs = state.data.item.filter(item => !item.buyer && item.transferable === true && item.title.toLowerCase().includes(input)
      )
      commit("updateBuyItemList", rs);}else{const rs = state.data.item.filter(item => !item.buyer && item.transferable === true
        ) || [];
        commit("setBuyItemList", rs)

      }
    },

    async filterBuyItemList({ commit }, input) {

      commit("updateBuyItemList", input);
    },

    async tagBuyItemList({ commit, state }, input) {
      if (!!input) { 
      const rs = state.buyItemList.filter(item =>  item.tags.find(tags => tags.includes(input)) && item.transferable === true)
        ;
        if (rs == []){
          const rs = state.data.item.filter(item => !item.buyer && item.transferable === true && item.tags.find(tags => tags.includes(input)) && item.transferable === true)
        }

      commit("updateBuyItemList", rs);}
    },

    async locationBuyItemList({ commit, state }, input) {

      if (!!input) { 
      const rs = state.buyItemList.filter(item => item.shippingregion.find(loc => loc.toLowerCase()).includes(input.toLowerCase()))
        ;

      commit("updateBuyItemList", rs);}
    },

    async priceMinBuyItemList({ commit, state }, input) {
      
      if (!!input) { 
      const rs = state.buyItemList.filter(item =>  (Number(item.estimationprice) > input))
        ;
   
      commit("updateBuyItemList", rs);}
    },

    async priceMaxBuyItemList({ commit, state }, input) {
      
      if (!!input) { 
      const rs = state.buyItemList.filter(item =>  (Number(item.estimationprice) < input))
        ;
   
      commit("updateBuyItemList", rs);}
    },

   

    async tagToEstimateList({ commit, state }, input) {
      if (!!input) {  //const A = state.data.item.filter(item => !item.buyer && item.tags.find(tags => tags.includes(input)) && item.transferable === false)
        ;
      //const B = state.estimatorItemList;

      //const rs = A.filter(a => !B.map(b => b.itemid).includes(a.id));
      let rs = state.toEstimateList.filter(item => item.tags.find(tag => tag == input) )
        ;
      //console.log(input);
      //console.log(rs);
      //console.log(A);
      //console.log(B);
      commit("setToEstimateList", rs);}
    },

    async regionToEstimateList({ commit, state }, input) {
      if (!!input) { let rs = state.toEstimateList.filter(item => item.shippingregion.find(region => region == input) )
        ;
     
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
      let merged = [].concat.apply([], rs);
      let frequency = {};
      merged.forEach(function (value) {  frequency[value.toLowerCase()] = 0; });

      let uniques = merged.filter(function (value) {
        return ++frequency[value.toLowerCase()] == 1;
      });

      let sorted = uniques.sort(function (a, b) {
        return frequency[b] - frequency[a];
      });

      //console.log(merged)
      /*console.log(rs)
      console.log(merged)
      console.log(uniques)
      console.log(sorted)
*/
      //console.log("TEST",sorted);

      commit("setTagList", sorted);
    },
    async setSortedLocationList({ commit, state }) {
      
      const rs = state.buyItemList.map(item => item.shippingregion);

      let merged = [].concat.apply([], rs);
      let frequency = {};
      merged.forEach(function (value) { frequency[value] = 0; });

      let uniques = merged.filter(function (value) {
        return ++frequency[value] == 1;
      });
let uppercase = uniques.map(tag => tag.toUpperCase())
      let sorted = uppercase.sort(function (a, b) {
        return frequency[b] - frequency[a];
      });
      
      if (sorted[0]) {
        commit("set", { key: 'locationList', value: sorted } );
      }else{
        //console.log(merged)
        commit("set", { key: 'locationList', value: merged } );
      }
     /* console.log(rs)
    console.log(merged)
      console.log(uniques)
      console.log(sorted)*/
      
    },

    async setToEstimateRegions({ commit, state }) {
//console.log("test")
      const rs = state.toEstimateList.map(item => item.shippingregion);
      let merged = [].concat.apply([], rs);
      let frequency = {};
      merged.forEach(function (value) { frequency[value.toLowerCase()] = 0; });

      let uniques = merged.filter(function (value) {
        return ++frequency[value.toLowerCase()] == 1;
      });

      let sorted = uniques.sort(function (a, b) {
        return frequency[b] - frequency[a];
      });

      if (sorted[0]) {
        commit("set", { key: 'regionList', value: sorted } );
      }else{
        //console.log(merged)
        commit("set", { key: 'regionList', value: merged } );
      }  
    },
   

    async setToEstimateList({ commit, state }) {
      const A = state.data.item.filter(item => item.estimationprice < 1 && item.status == '' && item.bestestimator == '');
      const B = state.estimatorItemList;
     /* const rs = A.filter(a => !B.map(b => b.itemid).includes(a.id));*/



      //console.log(A.filter(a => !B.map(b=>b.id).includes(a.id)));

      //where the items are not in estimator list
      const D = A.filter(a => !B.map(b => b.itemid).includes(a.id));

      const E = state.sellerItemList;

      const rs = D.filter(d => !E.map(e => e.id).includes(d.id));
      
      /*console.log(A);
      console.log(B);
      console.log(D);
      console.log(E);
      console.log(rs);*/
      commit("setToEstimateList", rs);


    },

    async setInterestedItemList({ commit, state }, input) {
      if (!!input) {  const rs = state.data.estimator.filter(estimator => estimator.estimator === input && estimator.interested
      );


      commit("setInterestedItemList", rs);}
    },
    async setSellerActionList({ commit, state }, input) {
      if (!!input) { 

      const toAccept = state.data.item.filter(item => item.seller == input && item.estimationprice > 0 && !item.buyer && !item.transferable
      );
      const toShip = state.data.item.filter(item => !item.buyer && item.seller === input && item.localpickup == '' && !item.tracking
      );
      //console.log(input);
      //console.log(state.account.address);
      //console.log(state.client.senderAddress);



      //console.log(toAccept);
      //console.log(toShip);
      toAccept.concat(toShip);
      //console.log(toAccept);

      commit("setSellerActionList", toAccept);}
    },

    async setBuySellerItemList({ commit, state }, input) {
      if (!!input) { 
      const rs = state.data.item.filter(item => item.seller === input) || [];

      commit("setBuySellerItemList", rs);}
    },
async updateItem({ commit, state }, input) {
      console.log(input)
      const url = `${process.env.VUE_APP_API}/${process.env.VUE_APP_PATH.replace(/\./g, '/')}/${"item/"+ input}`;
      console.log(url)
      axios.get(url).then(result => {
        let itemindex = state.data.item.findIndex(item => item.id == result.data.Item.id)
      //state.data.item.map(item => updated.id === item.id || item);
      state.data.item[itemindex] = result.data.Item

      console.log(result.data.Item)
      console.log(itemindex)
      console.log(state.data.item[itemindex])
      return result.data.Item

    }, error => {
        console.error("Got nothing from node")
    })


},
async setTransactions({ commit, state }, address) {


  try {
    let sent = (await axios.get(process.env.VUE_APP_API + '/cosmos/tx/v1beta1/txs?events=transfer.sender%3D%27' + address + '%27')).data;

    //let sentTransactions = JSON.stringify(sent.result)


    let received = (await axios.get(process.env.VUE_APP_API + '/cosmos/tx/v1beta1/txs?events=transfer.recipient%3D%27' + address+ '%27')).data;

   // let receivedTransactions = JSON.stringify(received.result)MsgItemTransfer
    //console.log(received)
   // console.log(receivedTransactions)
   // console.log(receivedTransactions)
   // console.log(sentTransactions)
    commit("set", { key: 'sentTransactions', value: sent });
    commit("set", { key: 'receivedTransactions', value: received });
  }


  catch (e) {
    //console.error(new SpVuexError('QueryClient:ServiceGetTxsEvent', 'API Node Unavailable. Could not perform query.'));
    console.log("ERROR" + e)

  }
},
async setEvent({ commit }, {type, attribute, value}){

  console.log(attribute)
  console.log(value)
        try {
          let event = (await axios.get(process.env.VUE_APP_API + '/cosmos/tx/v1beta1/txs?events=' + type + '.' + attribute + '%3D%27' + value + '%27')).data;
  console.log(event)
          commit('entitySet', { type, body: event })
        }
  
  
        catch (e) {
          //console.error(new SpVuexError('QueryClient:ServiceGetTxsEvent', 'API Node Unavailable. Could not perform query.'));
          console.log("ERROR" + e)
  
        }
      },
  
    },
    getters: {
      account: state => state.account, bankBalances: state => state.bankBalances, getSellerItemList: state => state.sellerItemList, getEstimatorItemList: state => state.estimatorItemList, getBuyerItemList: state => state.buyerItemList, getBuyItemList: state => state.buyItemList, getInterestedItemList: state => state.InterestedItemList, getItemByID: state => id => state.data.item.find((item) => item.id === id), getToEstimateList: state => state.toEstimateList, getSellerActionList: state => state.sellerActionList, getTagList: state => state.tagList, getLocationList: state => state.locationList, getRegionList: state => state.regionList, getBuySellerList: state => state.buySellerList, getReceivedTransactions: state => state.receivedTransactions, getSentTransactions: state => state.sentTransactions,getEvent: state => t => state.data[t],
  
    }
  
  
  });
  