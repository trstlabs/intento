<template>
  <div>
    <div class="pa-4 mx-auto">
      <v-card elevation="2" rounded="lg" >
        <v-progress-linear
          indeterminate
          :active="loadingitem"
        ></v-progress-linear>
        <div class="pa-2 mx-auto">
          

          <v-row> <p class="mx-6 overline text-center"> {{ thisitem.title }} </p><v-spacer /> <!--<v-btn   fab outlined
      
      small
      @click="setItem()"><v-icon >
        mdi-marker
      </v-icon></v-btn>--><v-btn text @click="removeItem()"><v-icon >
        mdi-trash-can
      </v-icon></v-btn> </v-row>
          
            
            <div class="ma-2" elevation="8">
            <v-carousel
              
              height="400"
              hide-delimiter-background
              show-arrows-on-hover
            >
              <v-carousel-item
                v-for="(photo, i) in photos"
                :key="i"
                :src="photo"
              >
              </v-carousel-item>
            </v-carousel>
          </div>
       
          <v-row align="start">
            <v-col cols="12">
              <v-card elevation="0" >  <div class="pl-4 overline text-center">Description</div> <v-card-text>
    
     
  <div class="body-1 "> "
           {{thisitem.description }} "
         </div> </v-card-text> </v-card>

             

 <!--<div v-for="comment in thisitem.comments" v-bind:key="comment" >
<v-text-field v-if="comment != ''" class="mt-2"
            :value="comment"
            label="Comment"
            auto-grow
            outlined
            readonly
    >
     </v-text-field>

</div> -->
<v-divider class="ma-2"/>
 <v-chip
      class="ma-1 caption"
      label
      outlined
      medium
    >
    <v-icon left>
        mdi-account-badge-outline
      </v-icon>
      Identifier: {{ thisitem.id }}
    </v-chip>

<v-chip v-if="thisitem.localpickup"
      class="ma-1 caption"
      label
      outlined
      medium
    >
    <v-icon left>
       mdi-map-marker-outline
      </v-icon>
      Local pickup available
    </v-chip>
       <v-chip
      class="ma-1 caption"
      label
      outlined
      medium

    >
    <v-icon left>
        mdi-star-outline
      </v-icon>
      Condition: {{thisitem.condition}}/5
    </v-chip>
    
          <v-chip v-if="thisitem.shippingcost"
      class="ma-1 caption"
      label
      outlined
      medium
    >
    <v-icon left>
        mdi-package-variant
      </v-icon>
      Shipping available
    </v-chip>

    <v-chip v-if="thisitem.shippingcost"
      class="ma-1 caption"
      label
      outlined
      medium
    >
    <v-icon left>
        mdi-package-variant-closed
      </v-icon>
      Shipping cost: ${{thisitem.shippingcost}} TPP
    </v-chip>
 <v-chip outlined medium label class="ma-1 caption"
            v-for="country in thisitem.shippingregion" :key="country"
          > <v-icon small left>
        mdi-flag-variant-outline
      </v-icon>{{ country }}</v-chip>

    

    <v-chip v-if="thisitem.bestestimator"
      class="ma-1 caption"
      label 
      outlined
      medium
    >
    <v-icon left>
        mdi-check-all
      </v-icon>
      Estimation Price: ${{thisitem.estimationprice}} TPP
    </v-chip>

    

<v-chip v-if="thisitem.transferable"
      class="ma-1 caption"
      label
      outlined
      medium
    >
    <v-icon left>
        mdi-swap-horizontal
      </v-icon>
      Transferable
    </v-chip>

    <v-chip  v-if="thisitem.buyer"
      class="ma-1 caption"
      label
      outlined
      medium
    >
    <v-icon left >
        mdi-cart-outline
      </v-icon>
      Buyer: {{thisitem.buyer}}
    </v-chip>

    <v-chip v-if="thisitem.status"
      class="ma-1 caption"
      label
      outlined
      medium
    >
    <v-icon left>
        mdi-clock-time-three-outline
      </v-icon>
      Status: {{ thisitem.status }}
    </v-chip>
             

              <v-chip v-if="(thisitem.buyer === '' && thisitem.transferable === true)"
      class="ma-1 caption"
      label
      outlined
      medium
    >
    <v-icon left>
        mdi-store
      </v-icon>
      Currently for sale
    </v-chip>
   
    <v-chip outlined medium label class="ma-1 caption"
            v-for="itemtags in thisitem.tags" :key="itemtags"
          > <v-icon small left>
        mdi-tag-outline
      </v-icon>{{ itemtags }}</v-chip>

    
        

<v-divider class="ma-2"/>
      
  <div class="overline text-center"> Comments </div> 
     <div v-if="thisitem.comments">
       <div  v-for="(comment, nr) in commentlist" v-bind:key="nr">
    <v-chip  color="primary" class="ma-2 "
        
    >{{ comment }}
     </v-chip>
</div>
     
     </div>
     <div v-if="!thisitem.comments">
<p  class="caption text-center"> No comments to show right now </p> </div>
     </v-col>
     </v-row>
<!--<v-divider class="ma-2"/>

           <v-row>    
           <v-btn class="pa-2 mt-2"
        
        text
        icon
        @click="shippingcost = 0"
      >
        <v-icon > {{shippingcost === 0 ? 'mdi-package-variant' : 'mdi-package-variant-closed'}} </v-icon>
      </v-btn>

                <v-slider class="pa-2 mt-2"
                  hint="Set to 0 tokens no for shipping"
                  
                  thumb-label
                  label="Shipping cost"
                  suffix="tokens"
                  :persistent-hint="shippingcost != 0"
                  
                  placeholder="Shipping cost"
                  :thumb-size="70"
                  v-model="shippingcost"
                  
                ><template v-slot:thumb-label="item">
            {{ item.value }} tokens
          </template> </v-slider>
</v-row>  <v-row  v-if="shippingcost"> <v-col>  <v-row>    
           <v-btn class="pa-2"
        
        text
        icon
        @click="localpickup = !localpickup"
      >
        <v-icon > {{localpickup ? 'mdi-map-marker' : 'mdi-map-marker-off'}} </v-icon>
      </v-btn>


      <v-switch class="ml-2" 
      v-model="localpickup"
      inset
      label="Local pickup"
      
      
    ></v-switch>  

    

               
                </v-row></v-col><v-col> <v-select
                 prepend-icon="mdi-earth"
                 hint="Leave blank for no shipping location"
                 :persistent-hint="selectedCountries == 0"
                
          v-model="selectedCountries"
          :items="countryCodes"
         :rules="rules.shippingRules"
          label="Ships to"
          deletable-chips
          multiple
          chips
          
        > </v-select> 
</v-col></v-row>

               <v-row v-if="shippingcost == 0 "> 
                
           <v-btn class="pa-2"
        
        text
        icon
        @click="localpickup = !localpickup"
      >
        <v-icon > {{localpickup ? 'mdi-map-marker' : 'mdi-map-marker-off'}} </v-icon>
      </v-btn>


      <v-switch class="ml-2" 
      v-model="localpickup"
      inset
      label="Local pickup"
      
      
    ></v-switch>  

    

               
                </v-row>-->
                
          <!--<v-divider class="ma-4"/>  
             
             
             
            

            
          
          <div class="text-center" v-if="thisitem.bestestimator != '' && getThisItem">
            <v-chip outlined>
              Final Estimation Price is revealed! 
            </v-chip>
          </div>-->
        </div>

        <v-card-actions>
          <v-btn
            color="blue"
            text
            @click="(showactions = !showactions), createStep()"
          >
            Actions
          </v-btn>

          <v-spacer></v-spacer>

          <v-btn icon @click="(showactions = !showactions), createStep()">
            <v-icon>{{
              showactions ? "mdi-chevron-up" : "mdi-chevron-down"
            }}</v-icon>
          </v-btn>
        </v-card-actions>

        <v-expand-transition>
          <div class="pa-2 mx-auto" elevation="8" v-if="showactions">
            <v-divider></v-divider>
            <div>
              <v-stepper class="elevation-0" v-model="step" vertical>
                <v-stepper-step step="1" complete> Place Item </v-stepper-step>

                <v-stepper-step
                  :complete="thisitem.bestestimator != '' || thisitem.estimationprice > 0"
                  step="2"
                >
                  Awaiting Estimation
                </v-stepper-step>

                <v-stepper-content step="2">
                 
                    <app-text type="subtitle">
                      Awaiting estimators to estimate the
                      item. Meanwhile... help others by estimating other items
                      (and earn tokens)!
                    </app-text>
                    
                  
                </v-stepper-content>

                <v-stepper-step  :complete="thisitem.transferable" step="3" >
                  Accept Estimation
                </v-stepper-step>

                <v-stepper-content step="3">
                  
                    <div>
                      
                        <app-text type="subtitle">
                          Wow! there is an estimation. You can sell {{thisitem.title}} for ${{thisitem.estimationprice}} TPP tokens. By accepting your item wil directly be able to
                          be purchased. Anyone can provide a prepayment to buy
                          the item. 
                        </app-text>
<v-row>
                        <v-btn class="ma-4" color="primary"
                          v-if="
                            !flightit &&
                            hasAddress &&
                            thisitem.bestestimator != '' &&
                            thisitem.transferable != true
                          "
                          @click="submititemtransferable(true, thisitem.id)"
                        ><v-icon left>
         mdi-checkbox-marked-circle
      </v-icon>
                          Accept 
                          <div class="button__label" v-if="flightit">
                            <div class="button__label__icon">
                              <icon-refresh />
                            </div>
                            Placing item for $ale...
                          </div>
                        </v-btn>
                      
                        <v-btn class="ma-4" color="default"
                          v-if="
                            !flightitn &&
                            hasAddress &&
                            thisitem.bestestimator != '' &&
                            thisitem.transferable != true
                          "
                          @click="submititemtransferable(false, thisitem.id)"
                        ><v-icon left>
         mdi-cancel
      </v-icon>
                          Reject 
                          <div class="button__label" v-if="flightitn">
                            <div class="button__label__icon">
                              <icon-refresh />
                            </div>
                            Deleting item...
                          </div>
                        </v-btn>
                        
                     
                      </v-row>
                    </div>
                 
                </v-stepper-content>

                <v-stepper-step :complete="thisitem.buyer != ''" step="4">
                  Item For Sale
                </v-stepper-step>

                <v-stepper-content step="4" :complete="thisitem.status != ''">
                  <app-text type="subtitle"
                    >Item placed. Awaiting buyer... Tip: share your item with family and friends.
                  </app-text>
                  <app-text
                    v-if="thisitem.shippingcost > 0 && thisitem.localpickup"
                    type="caption"
                  >
                    If a buyer chooses shipping, you ship it,
                    provide the track and trace code if available, and you'll automatically get
                    your tokens. After a buyer is found and chooses local
                    pickup, the buyer can pick it up. Tip:
                    let the buyer transfer the tokens during your meetup.
                  </app-text>
                  <app-text
                    v-if="thisitem.shippingcost === 0 && thisitem.localpickup"
                    type="caption"
                  >
                    After a buyer is found negotiate a meetup time and place by sending a message to the buyer.
                    Tip: let the buyer transfer the tokens during your meetup.
                  </app-text>
                  <app-text
                    v-if="
                      thisitem.shippingcost > 0 && thisitem.localpickup === false
                    "
                    type="caption"
                  >
                    After a buyer is found, find out about the address to ship to by sending a message to the buyer.
                  </app-text>
                </v-stepper-content>
                <v-stepper-step :complete="thisitem.status != ''" step="5">
                  Item Transfer
                </v-stepper-step>

                <v-stepper-content step="5">
                  <div
                    class="pa-8 mx-lg-auto"
                    v-if="
                      !!valid &&
                      !flightIS &&
                      hasAddress &&
                      thisitem.localpickup != true &&
                      thisitem.buyer &&
                      thisitem.shippingcost &&
                      thisitem.status === ''
                    "
                  >
                    
                    
                      <app-text type="caption"> 
                        Now it's time to ship the item. Provide a track and
                        trace code to the buyer if available.
                      </app-text>
                      <input
                        type="checkbox"
                        id="checkbox"
                        v-model="tracking"
                        v-bind:value="true"
                      />
                      <label for="checkbox"
                        > I have shipped the item and provided the buyer with
                        track and trace
                      </label>
                      <v-btn @click="submitItemShipping(tracking, thisitem.id)">
                        Receive tokens
                        <div class="button__label" v-if="flightIS">
                          <div class="button__label__icon">
                            <icon-refresh />
                          </div>
                          Collecting tokens...
                        </div>
                      </v-btn>
                    
                  </div>
                  <div
                    class="pa-8 mx-lg-auto"
                    v-if="
                      !!valid &&
                      !flightIS &&
                      hasAddress &&
                      thisitem.localpickup &&
                      thisitem.buyer &&
                      thisitem.status === ''
                    "
                  >
                    <v-divider class="pa-1" ></v-divider>
                   
                      <app-text type="caption">
                        Now its the time to meet up with the buyer! To complete the item
                        transfer, make sure the buyer sends the tokens at the
                        pick-up.
                      </app-text>
                   
                  </div>
                </v-stepper-content>

                <v-stepper-step :complete="thisitem.status != ''" step="6">
                  Done
                </v-stepper-step>

                <v-stepper-content step="5" height="200px"
                  ><v-card></v-card>
                </v-stepper-content>
              </v-stepper>
            </div>
          </div> </v-expand-transition
      ></v-card>
    </div>
  </div>
</template>

<script>
import { databaseRef } from './firebase/db';

import ItemListCreator from "./ItemListCreator.vue";
import { SigningStargateClient, assertIsBroadcastTxSuccess } from "@cosmjs/stargate";
import {  Registry } from '@cosmjs/proto-signing/';
import { Type, Field } from 'protobufjs';

export default {
  props: ["itemid"],
  components: { ItemListCreator },
  data() {
    return {
  
      shippingcost: "0",
      localpickup: false,
      selectedCountries: [],
      flightre: false,
      flightit: false,
      flightitn: false,
      flightIS: false,

      loadingitem: false,
      
      showactions: false,
      transferbool: false,
      tracking: false,
      photos: [],
      imageurl: "",
      step: 2,
      rules: {

         shippingRules:  [ 
          (v) => !!v.length == 1 || "A country is required when shipping cost is applicable",

        ], 
      },
      countryCodes:["NL", "BE", "UK", "DE", "US","CA"]
    };
  },

  mounted() {
    this.loadingitem = true;
    const id = this.itemid;


    const imageRef = databaseRef.ref("ItemPhotoGallery/" + id);
    imageRef.on("value", (snapshot) => {
      const data = snapshot.val();

      if (data != null && data.photo != null) {
        //console.log(data.photo);
        this.photos = data;
        this.imageurl = data.photo;
        this.loadingitem = false;
      }
    });
    this.loadingitem = false;
  },

  computed: {
    thisitem() {
      this.loadingitem = true;
      return this.$store.getters.getItemByID(this.itemid);
      this.loadingitem = false; 
    },

    hasAddress() {
      return !!this.$store.state.account.address;
    },
    valid() {
      return this.thisitem.id.trim().length > 0;
    },

    commentlist() {
      //const item = this.$store.getters.getItemByID(this.itemid);
      
      //console.log(this.thisitem);
      return this.thisitem.comments.filter(com => com != '') || [];
      //console.log( this.thisitem.comments.filter(i => i != ""));
    },
  },

  methods: {

     async removeItem() {
      
        this.loadingitem = true;
        this.flightre = true;
        const type = { type: "item" };
        const body = { id: this.thisitem.id };
      const fields = [
        ["creator", 1,'string', "optional"],                         
        ["id",2,'string', "optional"],
      ];
       
       this.itemdeleteSubmit({ ...type, body,fields })
     
        this.flightre = false;
        this.loadingitem = false;
    
    
      
    },

async itemdeleteSubmit({ body, fields }) {
      const wallet = this.$store.state.wallet
      const typeUrl = `/${process.env.VUE_APP_PATH}.MsgDeleteItem`;
      let MsgCreate = new Type(`MsgDeleteItem`);
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
          creator: this.$store.state.account.address,
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
   
    async getThisItem() {
      await submitrevealestimation();
      return this.thisitem();
      
    },
    async submititemtransferable(transferable, itemid) {
      if (this.valid && !this.flightit && this.hasAddress) {
     
        this.flightit = true;
        this.flightitn = true; 
         const fields = [
        ["creator", 1,'string', "optional"],
        ["transferable",2,'bool', "optional"],
         [ "itemid", 3,'string', "optional"] ,                                                    
     
      ];
        const body = { transferable, itemid };
        await this.transferableSubmit({  body, fields});
        
   
      }
    },
    async transferableSubmit( { body, fields }) {
      const wallet = this.$store.state.wallet
      const typeUrl = `/${process.env.VUE_APP_PATH}.MsgItemTransferable`;
      let MsgCreate = new Type(`MsgItemTransferable`);
      const registry = new Registry([[typeUrl, MsgCreate]]);

      fields.forEach(f => {
        MsgCreate = MsgCreate.add(new Field(f[0], f[1], f[2], f[3]))
      })

      const client = await SigningStargateClient.connectWithSigner(
        process.env.VUE_APP_RPC,
        wallet,
        { registry }
      );
      //console.log("TEST" + client)
      const msg = {
        typeUrl,
        value: {
          creator: this.$store.state.account.address,
          ...body
        }
      };

      const fee = {
        amount: [{ amount: '0', denom: 'tpp' }],
        gas: '200000'
      };

      const result = await client.signAndBroadcast(this.$store.state.account.address, [msg], fee);
      assertIsBroadcastTxSuccess(result);
      alert(" Placed! ");
    },

    async submitItemShipping(tracking, itemid) {
      if (this.valid && !this.flightIS && this.hasAddress) {
        this.flightIS = true;

        const body = { tracking, itemid };
        const fields = [
        ["creator", 1,'string', "optional"],
         [ "tracking", 2,'bool', "optional"] ,                                                    
        ["itemid",3,'string', "optional"],
    
      ];
        await this.shippingSubmit({ body, fields });
        
        this.flightIS = false;
        this.tracking = false;

      }
    },

     async shippingSubmit({ body, fields }) {
      const wallet = this.$store.state.wallet
      const typeUrl = `/${VUE_APP_PATH}.MsgItemShipping`;
      let MsgCreate = new Type(`MsgItemShipping`);
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
          creator: this.$store.state.account.address,
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
      alert("Transaction sent");

    },
    createStep() {
      if (this.thisitem.buyer != "") {
        this.step = 5;
      } else if (this.thisitem.transferable === true) {
        this.step = 4;
      } else if (this.thisitem.bestestimator != "") {
        this.step = 3;
      }
    },
  },
};
</script>

