<template>
  <div>
    <div class="pa-4 mx-auto">
      <v-card elevation="2" rounded="lg" >
        <v-progress-linear
          indeterminate
          :active="loadingitem"
        ></v-progress-linear>
        <div class="pa-2 mx-auto">
          

          <v-row> <p class="pa-2 h3 font-weight-medium "> {{ thisitem.title }} </p><v-spacer /> <v-btn   fab outlined
      
      small
      @click="setItem()"><v-icon >
        mdi-marker
      </v-icon></v-btn><v-btn text @click="removeItem()"><v-icon >
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
      class="ma-1"
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
      class="ma-1"
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
      class="ma-1"
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
      class="ma-1"
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
      class="ma-1"
      label
      outlined
      medium
    >
    <v-icon left>
        mdi-package-variant-closed
      </v-icon>
      Shipping cost: ${{thisitem.shippingcost}} TPP
    </v-chip>

    <v-chip v-if="thisitem.bestestimator"
      class="ma-1"
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
      class="ma-1"
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
      class="ma-1"
      label
      outlined
      medium
    >
    <v-icon left >
        mdi-cart
      </v-icon>
      Buyer: {{thisitem.buyer}}
    </v-chip>

    <v-chip v-if="thisitem.status"
      class="ma-1"
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
      class="ma-1"
      label
      outlined
      medium
    >
    <v-icon left>
        mdi-store
      </v-icon>
      This item is for sale currently
    </v-chip>
   
    <v-chip outlined medium label class="ma-1"
            v-for="itemtag in thisitem.tags" :key="itemtag"
          > <v-icon small left>
        mdi-tag-outline
      </v-icon>{{ itemtag }}</v-chip>
        

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
<v-divider class="ma-2"/>

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

    

               
                </v-row>
                
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
                  :complete="thisitem.bestestimator != ''"
                  step="2"
                >
                  Awaiting Estimation
                </v-stepper-step>

                <v-stepper-content step="2">
                 
                    <p>
                      Awaiting estimators to estimate the
                      item. Meanwhile... help others by estimating other items
                      (and earn tokens)!
                    </p>
                    <!-- <v-btn
                      v-if="
                        thisitem.bestestimator === '' &&
                        getThisItem &&
                        !flightre &&
                        thisitem.transferable != true &&
                        thisitem.buyer === ''
                      "
                      @click="submitrevealestimation(thisitem.id), getThisItem"
                    >
                      Reveal Item Estimation Price
                      <div class="button__label" v-if="flightre">
                        <div class="button__label__icon">
                          <icon-refresh />
                        </div>
                        Sending transaction...
                      </div>
                    </v-btn>-->
                  
                </v-stepper-content>

                <v-stepper-step :complete="thisitem.transferable" step="3">
                  Accept Estimation
                </v-stepper-step>

                <v-stepper-content step="3">
                  
                    <div>
                      
                        <app-text type="p">
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

                <v-stepper-content step="4">
                  <app-text type="p"
                    >Awaiting buyer... Share your item!
                  </app-text>
                  <app-text
                    v-if="thisitem.shippingcost > 0 && thisitem.localpickup"
                    type="p"
                  >
                    After a buyer is found and chooses shipping, you ship it,
                    provide the track and trace and you'll automatically get
                    your tokens! After a buyer is found and chooses local
                    pickup, the buyer can pick it up at your convienience. Tip:
                    let the buyer transfer the tokens during your meetup.
                  </app-text>
                  <app-text
                    v-if="thisitem.shippingcost === 0 && thisitem.localpickup"
                    type="p"
                  >
                    After a buyer is found, negotiate a meetup time and place.
                    Tip: let the buyer transfer the tokens during your meetup.
                  </app-text>
                  <app-text
                    v-if="
                      thisitem.shippingcost > 0 && thisitem.localpickup === false
                    "
                    type="p"
                  >
                    After a buyer is found, negotiate a meetup time and place.
                    Tip: let the buyer transfer the tokens during your meetup.
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
                    
                    
                      <app-text>
                        Now its the time to ship the item! provide the track and
                        trace and you'll automatically get your tokens!
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
                   
                      <app-text>
                        Now its the time to meet up with the buyer! To complete the item
                        transfer, make sure the buyer sends the tokens at the
                        pick-up.
                      </app-text>
                   
                  </div>
                </v-stepper-content>

                <v-stepper-step :complete="step > 6" step="6">
                  Done!
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

import ItemListCreator from "./ItemListCreator.vue";
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
    const db = firebase.database();

    const imageRef = db.ref("ItemPhotoGallery/" + id);
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
    async submitrevealestimation(itemid) {
      if (this.valid && !this.flightre && this.hasAddress) {
        this.loadingitem = true;
        this.flightre = true;
        const type = { type: "item/reveal" };
        const body = { itemid };
        await this.$store.dispatch("entitySubmit", { ...type, body });
        //await this.$store.dispatch("entityFetch", type);
        //await this.$store.dispatch("accountUpdate");
        this.flightre = false;
        this.loadingitem = false;
        //this.deposit = "";
        alert("Transaction sent");
      }
    },
     async removeItem() {
      
        this.loadingitem = true;
        this.flightre = true;
        const type = { type: "item/delete" };
        const body = { id: this.thisitem.id };
      
       
        await this.$store.dispatch("entitySubmit", { ...type, body });
     
        this.flightre = false;
        this.loadingitem = false;
    
        alert("Transaction sent");
      
    },

     async setItem() {
      
        this.loadingitem = true;
        this.flightre = true;
        const type = { type: "item/set" };
        const body = { id: this.thisitem.id, shippingcost: this.shippingcost.toString(), localpickup: this.localpickup, shippingregion: this.shippingregion };
      
       
        await this.$store.dispatch("entitySubmit", { ...type, body });
     
        this.flightre = false;
        this.loadingitem = false;
    
        alert("Transaction sent");
      
    },

    async getThisItem() {
      await submitrevealestimation();
      return this.thisitem();
      
    },
    async submititemtransferable(transferbool, itemid) {
      if (this.valid && !this.flightit && this.hasAddress) {
       /* if (transferbool === true) {
          this.flightit = true;
        }
        if (transferbool === false) {
          this.flightitn = true;
        }*/
        this.flightit = true;
        this.flightitn = true; 
        const type = { type: "item/transferable" };
        const body = { transferbool, itemid };
        await this.$store.dispatch("entitySubmit", { ...type, body });
        //await this.$store.dispatch("entityFetch", type);
        //await this.$store.dispatch("accountUpdate");

         const thisitemcheck = await this.$store.getters.getItemByID(this.thisitem.id);
        
        console.log(thisitemcheck);
        if (thisitemcheck.transferbool) {
          alert("Item Placed ");
        }
        
      






        //this.flightit = false;
        //this.flightitn = false;
        //this.deposit = "";
        alert("Transaction sent");
      }
    },
    async submitItemShipping(tracking, itemid) {
      if (this.valid && !this.flightIS && this.hasAddress) {
        this.flightIS = true;
        const type = { type: "item/shipping" };
        const body = { tracking, itemid };
        await this.$store.dispatch("entitySubmit", { ...type, body });
        this.flightIS = false;
        this.tracking = false;
        alert("Transaction sent");
      }
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

<style scoped>
button {
  background: none;
  border: none;
  color: rgba(0, 125, 255);
  padding: 0;
  font-size: inherit;
  font-weight: 800;
  font-family: inherit;
  text-transform: uppercase;
  margin-top: 0.5rem;
  cursor: pointer;
  transition: opacity 0.1s;
  letter-spacing: 0.03em;
  transition: color 0.25s;
  display: inline-flex;
  align-items: center;
}
.item {
  box-shadow: inset 0 0 0 1px rgba(0, 0, 0, 0.1);
  margin-bottom: 1rem;
  padding: 1rem;
  border-radius: 0.5rem;
  overflow: hidden;
}
.item__field {
  display: grid;
  line-height: 1.5;
  grid-template-columns: 15% 1fr;
  grid-template-rows: 1fr;
  word-break: break-all;
}
.item__field__key {
  color: rgba(0, 0, 0, 0.25);
  word-break: keep-all;
  overflow: hidden;
}
button:focus {
  opacity: 0.85;
  outline: none;
}

.button__label {
  display: inline-flex;
  align-items: center;
}
.button__label__icon {
  height: 1em;
  width: 1em;
  margin: 0 0.5em 0 0.5em;
  fill: rgba(0, 0, 0, 0.25);
  animation: rotate linear 4s infinite;
}
.button.button__valid__false {
  color: rgba(0, 0, 0, 0.25);
  cursor: not-allowed;
}
.card__empty {
  margin-bottom: 1rem;
  border: 1px dashed rgba(0, 0, 0, 0.1);
  padding: 1rem;
  display: flex;
  align-items: center;
  justify-content: center;
  box-sizing: border-box;
  border-radius: 8px;
  color: rgba(0, 0, 0, 0.25);
  text-align: center;
  min-height: 8rem;
}
@keyframes rotate {
  from {
    transform: rotate(0);
  }
  to {
    transform: rotate(-360deg);
  }
}
@media screen and (max-width: 980px) {
  .narrow {
    padding: 0;
  }
}
</style>
