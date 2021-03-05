
import Vue from "vue";
import Vuex from "vuex";
import axios from "axios";
import app from "./app.js";

//import cosmos from "@tendermint/vue/src/store/cosmos.js";
import { assert } from "@cosmjs/utils";
import { assertIsBroadcastTxSuccess, makeCosmoshubPath, coin } from '@cosmjs/launchpad'
import { SigningStargateClient } from "@cosmjs/stargate";
import { DirectSecp256k1HdWallet } from '@cosmjs/proto-signing';
import { Type, Field } from 'protobufjs';
import { Registry } from '@cosmjs/proto-signing';


Vue.use(Vuex);


const API = "http://localhost:1317";
//const API = "https://node.trustpriceprotocol.com"
const ADDRESS_PREFIX = 'cosmos';
const PATH = 'danieljdd.tpp.tpp'

const RPC = 'http://localhost:26657'

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
    clientUpdate(state, { client }) {
      state.client = client;
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
     
    
      state.app.types.forEach(({ type }) => {
        if (type == "estimator/flag" || type == "item/reveal" || type == "item/transferable" || type == "item/transfer" || type == "item/shipping") { return; }
        dispatch("entityFetch", { type });
      });

      await dispatch('accountSignInTry');

    },

    async chainIdFetch({ commit }) {
      const node_info = (await axios.get(`${API}/node_info`)).data.node_info;
      commit("chainIdSet", { chain_id: node_info.network });
    },
    async accountSignInTry({ state, dispatch }) {
      const mnemonic = localStorage.getItem('mnemonic')
      if (mnemonic) {
        await dispatch('accountSignIn', { mnemonic })
        await  dispatch("setEstimatorItemList", state.account.address);
        await  dispatch("setToEstimateList", state.account.address);
        await  dispatch("setCreatorActionList", state.account.address);
        await   dispatch("setSortedTagList");
        await    dispatch("setCreatorItemList");
        await    dispatch("setBuyerItemList", state.account.address);

        await dispatch("setInterestedItemList", state.account.address);
      //$emit('signedIn');
      }
    },
    async accountSignIn(
      { commit, dispatch },
      { mnemonic }
    ) {
      //const { API, RPC, ADDR_PREFIX } = rootState.cosmos.env.env
      const wallet = await DirectSecp256k1HdWallet.fromMnemonic(
        mnemonic,
        makeCosmoshubPath(0),
        ADDRESS_PREFIX
      )
  
      //console.log("fdgadagfgfd")
      localStorage.setItem('mnemonic', mnemonic)
      const { address } = wallet
      const url = `${API}/auth/accounts/${address}`
      const acc = (await axios.get(url)).data
      const account = acc.result.value
      commit('set', { key: 'wallet', value: wallet })
      commit('set', { key: 'account', value: account })
      //console.log("fdgadagfgfd" + SigningStargateClient.connectWithSigner());
      ////onsole.log(RPC)
      const client = await SigningStargateClient.connectWithSigner(RPC, wallet, {});
      commit('set', { key: 'client', value: client })
      //console.log(client)
      try {
        await dispatch('bankBalancesGet')
      } catch {
        console.log('Error in getting a bank balance.')
      }
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


    async entitySubmit({ state, dispatch }, { type, fields, body }) {
      const mnemonic = localStorage.getItem('mnemonic')
      const wallet = await DirectSecp256k1HdWallet.fromMnemonic(
        mnemonic,
        makeCosmoshubPath(0),
        ADDRESS_PREFIX
      )

      const type2 = type.charAt(0).toUpperCase() + type.slice(1)
      const typeUrl = `/${PATH}.MsgCreate${type2}`;
      let MsgCreate = new Type(`MsgCreate${type2}`);
      const registry = new Registry([[typeUrl, MsgCreate]]);
      fields.forEach(f => {
        MsgCreate = MsgCreate.add(new Field(f[0], f[1], f[2], f[3]))
      })
      //console.log(registry );
      const [firstAccount] = await wallet.getAccounts();
      //console.log("creator" + state.wallet.address);
      const client = await SigningStargateClient.connectWithSigner(
        RPC,
        wallet,
        { registry }
      );
      //console.log("TEST" + client)
      const msg = {
        typeUrl,
        value: {
          creator: state.account.address,
          ...body
        }
      };
      const fee = {
        amount: [{ amount: '0', denom: 'tpp' }],
        gas: '200000'
      };
       const result = await client.signAndBroadcast(firstAccount.address, [msg], fee);
        assertIsBroadcastTxSuccess(result);

        console.log(result)

      try {
        //const path = "danieljdd.tpp.tpp".replace(/\./g, '/')
        
        //console.log(data)
        //console.log(firstAccount.address, [msg], fee);
        await dispatch('entityFetch', {
          type: type
        //  path: path
        }
        )
      } catch (e) {
        console.log(e)
      }

    },


    async itemSubmit({ state, commit, dispatch }, { type, fields, body }) {
      
      const wallet = state.wallet

      console.log("TESTwallet" + wallet )
      const type2 = type.charAt(0).toUpperCase() + type.slice(1)
      const typeUrl = `/${PATH}.MsgCreate${type2}`;
      let MsgCreate = new Type(`MsgCreate${type2}`);
      const registry = new Registry([[typeUrl, MsgCreate]]);
      fields.forEach(f => {
        MsgCreate = MsgCreate.add(new Field(f[0], f[1], f[2], f[3]))
      })

      const [firstAccount] = await wallet.getAccounts();

      const client = await SigningStargateClient.connectWithSigner(
        RPC,
        wallet,
        { registry }
      );

      const msg = {
        typeUrl,
        value: {
          creator: state.account.address,
          ...body
        }
      };

      const fee = {
        amount: [{ amount: '0', denom: 'tpp' }],
        gas: '200000'
      };
      await dispatch('entityFetch', {
        type: type
      })
      await dispatch("setCreatorItemList", state.account.address)
      let creatoritems = state.creatorItemList || []
      console.log(creatoritems)

      try {
        const result = await client.signAndBroadcast(firstAccount.address, [msg], fee);
        assertIsBroadcastTxSuccess(result);
        await dispatch('entityFetch', {
          type: type
        })
        await dispatch("setCreatorItemList", state.account.address)
        let newcreatoritems = state.creatorItemList

        let len = (creatoritems.length)
        console.log((newcreatoritems[len].id))
        commit('set', { key: 'newitemID', value: (newcreatoritems[len].id) })
      } catch (e) {
        console.log(e)
      }

    },

    async estimationSubmit({ state, dispatch }, { type, body }) {
      const mnemonic = localStorage.getItem('mnemonic')
      const wallet = await DirectSecp256k1HdWallet.fromMnemonic(
        mnemonic,
        makeCosmoshubPath(0),
        ADDRESS_PREFIX
      )
      const type2 = type.charAt(0).toUpperCase() + type.slice(1)
      const typeUrl = `/${PATH}.MsgCreate${type2}`;
      let MsgCreate = new Type(`MsgCreate${type2}`);
      const registry = new Registry([[typeUrl, MsgCreate]]);
      const fields = [
        ["estimator", 1,'string', "optional"],
         [ "estimation", 2,'int64', "optional"] ,                                                    
        ["itemid",3,'string', "optional"],
       ["deposit", 4, "int64", "optional"],
        ["interested",5,'bool', "optional"],
        ["comment",6,'string', "optional" ],  
      ];
console.log(fields)
fields.forEach(f => {
  MsgCreate = MsgCreate.add(new Field(f[0], f[1], f[2], f[3]))
})
console.log(MsgCreate)
      //console.log(registry );
      //console.log("creator" + state.wallet.address);
      const client = await SigningStargateClient.connectWithSigner(
        RPC,
        wallet,
        { registry }
      );
      //console.log("TEST" + client)
      const msg = {
        typeUrl,
        value: {
          estimator: state.account.address,
          deposit: 5,
          ...body
        }
      };
      
      console.log(msg)
      const fee = {
        amount: [{ amount: '0', denom: 'tpp' }],
        gas: '200000'
      };
       const result = await client.signAndBroadcast(state.account.address, [msg], fee);
        assertIsBroadcastTxSuccess(result);

        console.log(result)

      try {
    
        await dispatch('entityFetch', {
          type: type
        //  path: path
        }
        )
      } catch (e) {
        console.log(e)
      }

    },


    async revealSubmit({ state }, { body, fields }) {
      const mnemonic = localStorage.getItem('mnemonic')
      const wallet = await DirectSecp256k1HdWallet.fromMnemonic(
        mnemonic,
        makeCosmoshubPath(0),
        ADDRESS_PREFIX
      )
     
      const typeUrl = `/${PATH}.MsgRevealEstimation`;
      let MsgCreate = new Type(`MsgRevealEstimation`);
      const registry = new Registry([[typeUrl, MsgCreate]]);

fields.forEach(f => {
  MsgCreate = MsgCreate.add(new Field(f[0], f[1], f[2], f[3]))
})

      const client = await SigningStargateClient.connectWithSigner(
        RPC,
        wallet,
        { registry }
      );
      //console.log("TEST" + client)
      const msg = {
        typeUrl,
        value: {
          creator: state.account.address,
          ...body
        }
      };
      
      console.log(msg)
      const fee = {
        amount: [{ amount: '0', denom: 'tpp' }],
        gas: '200000'
      };
       //const result = await client.signAndBroadcast(state.account.address, [msg], fee);
        //assertIsBroadcastTxSuccess(result);
        await client.signAndBroadcast(state.account.address, [msg], fee);
       
      

    },

    async transferableSubmit({ state }, { body, fields }) {
      const mnemonic = localStorage.getItem('mnemonic')
      const wallet = await DirectSecp256k1HdWallet.fromMnemonic(
        mnemonic,
        makeCosmoshubPath(0),
        ADDRESS_PREFIX
      )
     
      const typeUrl = `/${PATH}.MsgItemTransferable`;
      let MsgCreate = new Type(`MsgItemTransferable`);
      const registry = new Registry([[typeUrl, MsgCreate]]);

fields.forEach(f => {
  MsgCreate = MsgCreate.add(new Field(f[0], f[1], f[2], f[3]))
})

      const client = await SigningStargateClient.connectWithSigner(
        RPC,
        wallet,
        { registry }
      );
      //console.log("TEST" + client)
      const msg = {
        typeUrl,
        value: {
          creator: state.account.address,
          ...body
        }
      };
      
      console.log(msg)
      const fee = {
        amount: [{ amount: '0', denom: 'tpp' }],
        gas: '200000'
      };
        
       const result = await client.signAndBroadcast(state.account.address, [msg], fee);
        assertIsBroadcastTxSuccess(result);
        alert(" Placed! ");
      

    },

    async transferSubmit({ state }, { body, fields }) {
      const mnemonic = localStorage.getItem('mnemonic')
      const wallet = await DirectSecp256k1HdWallet.fromMnemonic(
        mnemonic,
        makeCosmoshubPath(0),
        ADDRESS_PREFIX
      )
     
      const typeUrl = `/${PATH}.MsgItemTransfer`;
      let MsgCreate = new Type(`MsgItemTransfer`);
      const registry = new Registry([[typeUrl, MsgCreate]]);
console.log(fields)
fields.forEach(f => {
  MsgCreate = MsgCreate.add(new Field(f[0], f[1], f[2], f[3]))
})

      const client = await SigningStargateClient.connectWithSigner(
        RPC,
        wallet,
        { registry }
      );
      //console.log("TEST" + client)
      const msg = {
        typeUrl,
        value: {
          buyer: state.account.address,
          ...body
        }
      };
      const fee = {
        amount: [{ amount: '0', denom: 'tpp' }],
        gas: '200000'
      };

        const result = await client.signAndBroadcast(state.account.address, [msg], fee);
        assertIsBroadcastTxSuccess(result);
        alert("Transaction sent");

    },

    async paySubmit({ state }, { body, fields }) {
      const mnemonic = localStorage.getItem('mnemonic')
      const wallet = await DirectSecp256k1HdWallet.fromMnemonic(
        mnemonic,
        makeCosmoshubPath(0),
        ADDRESS_PREFIX
      )
     
      const typeUrl = `/${PATH}.MsgCreateBuyer`;
      let MsgCreate = new Type(`MsgCreateBuyer`);
      const registry = new Registry([[typeUrl, MsgCreate]]);
console.log(fields)
fields.forEach(f => {
  MsgCreate = MsgCreate.add(new Field(f[0], f[1], f[2], f[3]))
})

      const client = await SigningStargateClient.connectWithSigner(
        RPC,
        wallet,
        { registry }
      );
      //console.log("TEST" + client)
      const msg = {
        typeUrl,
        value: {
          buyer: state.account.address,
          ...body
        }
      };
      
      console.log(msg)
      const fee = {
        amount: [{ amount: '0', denom: 'tpp' }],
        gas: '200000'
      };
       //const result = await client.signAndBroadcast(state.account.address, [msg], fee);
        //assertIsBroadcastTxSuccess(result);
        const result = await client.signAndBroadcast(state.account.address, [msg], fee);
        assertIsBroadcastTxSuccess(result);
        alert("Transaction sent");

    },

   

    async shippingSubmit({ state }, { body, fields }) {
      const mnemonic = localStorage.getItem('mnemonic')
      const wallet = await DirectSecp256k1HdWallet.fromMnemonic(
        mnemonic,
        makeCosmoshubPath(0),
        ADDRESS_PREFIX
      )
     
      const typeUrl = `/${PATH}.MsgItemShipping`;
      let MsgCreate = new Type(`MsgItemShipping`);
      const registry = new Registry([[typeUrl, MsgCreate]]);
console.log(fields)
fields.forEach(f => {
  MsgCreate = MsgCreate.add(new Field(f[0], f[1], f[2], f[3]))
})

      const client = await SigningStargateClient.connectWithSigner(
        RPC,
        wallet,
        { registry }
      );
    
      const msg = {
        typeUrl,
        value: {
          creator: state.account.address,
          ...body
        }
      };
      
      console.log(msg)
      const fee = {
        amount: [{ amount: '0', denom: 'tpp' }],
        gas: '200000'
      };
    
        const result = await client.signAndBroadcast(state.account.address, [msg], fee);
        assertIsBroadcastTxSuccess(result);
        alert("Transaction sent");

    },


    //for a delete request [cors error in development]
    async entityDelete({ state }, { type, body }) {
      const { chain_id } = state;
      const creator = state.account.address;
      const base_req = { chain_id, from: creator };
      const req = { base_req, creator, ...body };
      //const headers = { 'Authorization': 'token', 'content-type': 'text/plain', 'Access-Control-Allow-Origin': '*',  'Access-Control-Allow-Methods': 'DELETE'};
      //const { data } = await axios.request(`${API}/${chain_id}/${type}`, req, 'delete')
      console.log("req is=" + req)


      /* let headers = {
         'Access-Control-Allow-Origin': '*',
         'Access-Control-Allow-Methods': 'GET, POST ,PUT ,DELETE ,OPTIONS',
         'Access-Control-Allow-Headers':
           'Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With'
       }*/






      const { data } = await axios.delete(`${API}/${chain_id}/${type}`, req)
        .catch(error => {
          this.errorMessage = error.message;
          console.error("There was an error!", error)
        });


    },

    //for a put request [cors error in development]
    async entitySet({ state }, { type, body }) {
      const { chain_id } = state;
      const creator = state.account.address;
      const base_req = { chain_id, from: creator };
      const req = { base_req, creator, ...body };

      const { data } = await axios.put(`${API}/${chain_id}/${type}`, req)

    },




    async setCreatorItemList({ commit, state }, input) {
      const rs = state.data.item.filter(item => item.creator === input
      ) || [];
      commit("setCreatorItemList", rs);

    },

    async setEstimatorItemList({ commit, state }, input) {

      const rse = state.data.estimator.filter(estimator => estimator.estimator === input
      ) || [];


      commit("setEstimatorItemList", rse);
    },

    async setBuyerItemList({ commit, state }, input) {
      const rs = state.data.item.filter(item => item.buyer === input
      );
      commit("setBuyerItemList", rs);
    },

    async setBuyItemList({ commit, state }) {

      const rs = state.data.item.filter(item => !item.buyer && item.transferable === true
      ) || [];

      commit("setBuyItemList", rs);
    },

    async setLocalBuyItemList({ commit, state }) {
      const rs = state.data.item.filter(item => !item.buyer&& item.transferable === true && item.localpickup === true
      );
      commit("setBuyItemList", rs);
    },
    async updateBuyItemList({ commit, state }, input) {
      const rs = state.data.item.filter(item => !item.buyer && item.transferable === true && item.title.toLowerCase().includes(input)
      );
      commit("updateBuyItemList", rs);
    },

    async filterBuyItemList({ commit }, input) {

      commit("updateBuyItemList", input);
    },

    async tagBuyItemList({ commit, state }, input) {

      const rs = state.data.item.filter(item => !item.buyer && item.tags.find(tags => tags.includes(input)) && item.transferable === true)
        ;

      commit("updateBuyItemList", rs);
    },

    async tagToEstimateList({ commit, state }, input) {
      const A = state.data.item.filter(item => !item.buyer && item.tags.find(tags => tags.includes(input)) && item.transferable === false)
        ;

      //const test2 = state.data.item.filter(item => !item.buyer && item.tags.strconv().includes(input)  && item.transferable === false) 
      //;


      //const test2 = state.data.item.forEach(filter(item.tags => item.tags.includes(input));
      //const test = state.data.item.tags.filter(tags => tags.includes(input));
      //console.log(test2);

      const B = state.estimatorItemList;

      const rs = A.filter(a => !B.map(b => b.itemid).includes(a.id));
      //console.log(input);
      //console.log(rs);
      //console.log(A);
      //console.log(B);
      commit("setToEstimateList", rs);
    },

    /*async setTagList({commit, state}) {
      //console.log("asd");
      //console.log(state.data.item);
      const rs = state.data.item.map(item => item.tags);
      //console.log(rs);
      var merged = [].concat.apply([], rs);
      //console.log(merged);
      commit("setTagList", merged);
    },*/

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
      const A = state.data.item.filter(item => !item.bestestimator);
      const B = state.estimatorItemList;
      //const rsEIL = state.data.estimator.filter(estimator => estimator.estimator === state.client.anyValidAddress);
      console.log(A);
      console.log(B);
      //console.log(A.filter(a => !B.map(b=>b.id).includes(a.id)));
      const rs = A.filter(a => !B.map(b => b.itemid).includes(a.id));

      commit("setToEstimateList", rs);


    },

    async setInterestedItemList({ commit, state }, input) {
      const rs = state.data.estimator.filter(estimator => estimator.estimator === input && estimator.interested
      );


      commit("setInterestedItemList", rs);
    },
    async setCreatorActionList({ commit, state }, input) {


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

      commit("setCreatorActionList", toAccept);
    },

    async setSellerItemList({ commit, state }, input) {

      const rs = state.data.item.filter(item => item.creator === input) || [];

      commit("setSellerItemList", rs);
    },
  },
  getters: {
    account: state => state.account, bankBalances: state => state.bankBalances, getCreatorItemList:  state => state.creatorItemList, getEstimatorItemList: state => state.estimatorItemList, getBuyerItemList: state => state.buyerItemList, getBuyItemList: state => state.buyItemList, getInterestedItemList: state => state.InterestedItemList, getItemByID: state => id => state.data.item.find((item) => item.id === id), getToEstimateList: state => state.toEstimateList, getCreatorActionList: state => state.creatorActionList, getTagList: state => state.tagList, getSellerList: state => state.sellerList,

  }


});
