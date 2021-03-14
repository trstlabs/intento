<template>
  <div>
    <div class="pa-2 mx-auto">
      <v-card elevation="2" rounded="lg">
        <v-progress-linear
          indeterminate
          :active="loadingitem"
        ></v-progress-linear>
        <div class="pa-2 mx-auto" elevation="8">
         <v-row> <p class="pa-2 h3 font-weight-medium "> {{ thisitem.title }} </p><v-spacer /><v-btn text @click="removeItem()"><v-icon >
        mdi-trash-can
      </v-icon></v-btn> </v-row>
          
            <v-divider></v-divider>
         
          <v-row align="start">
            <v-col cols="8">
                <v-chip
      class="mt-1"
      label
      outlined
      medium
    >
    <v-icon left>
        mdi-account-badge
      </v-icon>
      Identifier: {{ thisitem.id }}
    </v-chip>

    <v-chip v-if="thisitem.bestestimator"
      class="mt-1"
      label
      outlined
      medium
    >
    <v-icon left>
        mdi-check-all
      </v-icon>
      Best Estimation: $ {{thisitem.estimationprice}} tokens
    </v-chip>

 <v-card elevation="0" >  <div class="pa-2 overline">Description</div> <v-card-text>
    
     
  <div class="body-1 "> "
           {{thisitem.description }} "
         </div> </v-card-text> </v-card>

              <app-text class="mt-1" v-if="thisitem.bestestimator === userAddress" type="p"
                >You are the best estimator. Check your balance. If the item is
                transferred, you will be rewarded tokens.
              </app-text>
              <app-text  class="mt-1"  v-if="thisitem.lowestestimator === userAddress" type="p"
                >
                <v-icon left>
        mdi-account-arrow-left
      </v-icon>
                
                
                You are the lowest estimator. If the item owner does not accept
                the estimation price, you  lose the deposit.
              </app-text>
              <app-text class="mt-1" 
                v-if="thisitem.highestestimator === userAddress"
                type="p"
                >
                <v-icon left>
        mdi-account-arrow-right
      </v-icon>
      You are the highest estimator. If the item is not transferred,
                you lose the deposit.
              </app-text>
            </v-col>

            <v-col cols="4">
              <div v-if="imageurl" class="d-flex flex-row-reverse text-center">
                <v-avatar class="ma-2 rounded" size="125" tile>
                  <v-img :src="imageurl"></v-img>
                </v-avatar>
              </div>
            </v-col>
          </v-row>
          
        </div>
      </v-card>
    </div>
  </div>
</template>

<script>
import { databaseRef } from './firebase/db';
import ItemListEstimator from "./ItemListEstimator.vue";
import { SigningStargateClient, assertIsBroadcastTxSuccess } from "@cosmjs/stargate";
import {  Registry } from '@cosmjs/proto-signing/';
import { Type, Field } from 'protobufjs';


export default {
  props: ["itemid"],
  components: { ItemListEstimator },
  data() {
    return {
      loadingitem: true,

      photos: [],
      imageurl: "",
    };
  },

  mounted() {
    this.loadingitem = true;
    const id = this.itemid;

    const imageRef = databaseRef.ref("ItemPhotoGallery/" + id);
    imageRef.on("value", (snapshot) => {
      const data = snapshot.val();

      if (data != null && data.photo != null) {
        console.log(data.photo);
        this.photos = data;
        this.imageurl = data.photo;
        this.loadingitem = false;
      }
    });
    this.loadingitem = false;
  },

  computed: {
    thisitem() {
      //this.loadingitem = true;
      return this.$store.getters.getItemByID(this.itemid);
      this.loadingitem = false;
    },
    hasAddress() {
      return !!this.$store.state.account.address;
    },

    userAddress() {
      return this.$store.state.account.address;
    },
    valid() {
      return this.thisitem.id.trim().length > 0;
    },
  },

  methods: {
    async removeItem() {
      
        this.loadingitem = true;
        this.flightre = true;
        const type = { type: "estimator" };
        const body = { itemid: this.itemid };
       const fields = [
        ["estimator", 1,'string', "optional"],                         
        ["itemid",2,'string', "optional"],
      ];
       
        await this.estimatordeleteSubmit({  body, fields });
     
        this.flightre = false;
        this.loadingitem = false;
    

      
    },
    async estimatordeleteSubmit({ body, fields }) {
      const wallet = this.$store.state.wallet
      const typeUrl = `/${process.env.VUE_APP_PATH}.MsgDeleteEstimator`;
      let MsgCreate = new Type(`MsgDeleteEstimator`);
      const registry = new Registry([[typeUrl, MsgCreate]]);
      console.log(fields)
      fields.forEach(f => {
        MsgCreate = MsgCreate.add(new Field(f[0], f[1], f[2], f[3]))
      })

      const client = await SigningStargateClient.connectWithSigner(
        process.env.VUE_APP_RPC,
        wallet,
        { registry }
      );

      const msg = {
        typeUrl,
        value: {
          estimator: this.$store.state.account.address,
          ...body
        }
      };

      console.log(msg)
      const fee = {
        amount: [{ amount: '0', denom: 'tpp' }],
        gas: '200000'
      };

      const result = await client.signAndBroadcast(this.$store.state.account.address, [msg], fee);
      assertIsBroadcastTxSuccess(result);
      alert("Delete request sent");

    },
  },
};
</script>
