<template>
  <div>
    <div class="pa-4 mx-auto">
      <v-card elevation="2" class="pa-2" rounded="lg">
        <v-progress-linear
          indeterminate
          :active="loadingitem"
        ></v-progress-linear>
        <div class="pa-2 mx-auto">
          <v-row>
            <p class="mx-6 overline text-center">{{ thisitem.title }}</p>
            <v-spacer />
          <v-btn  v-if="thisitem.status == ''" text @click="removeItem()"
              ><v-icon> mdi-trash-can </v-icon></v-btn
            >
          </v-row>

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
              <v-card elevation="0">
                <div class="pl-4 overline text-center">Description</div>
                <v-card-text>
                  <div class="body-1">" {{ thisitem.description }} "</div>
                </v-card-text>
              </v-card>

      
              <v-divider class="ma-2" />
               <div class="overline mb-2 text-center">Information</div>

            <v-dialog transition="dialog-bottom-transition" max-width="300">
              <template v-slot:activator="{ on, attrs }">
                <span v-bind="attrs" v-on="on">
                  <v-chip
                    style="cursor: pointer"
                    class="ma-1 font-weight-light"
                    outlined
                    medium
                    >Condition:
                    <v-rating
                      :value="Number(thisitem.condition)"
                      readonly
                      color="primary darken-1"
                      background-color="primary lighten-1"
                      small
                      dense
                    ></v-rating>
                  </v-chip>
                </span>
              </template>
              <template v-slot:default="dialog">
                <v-card>
                  <v-toolbar color="default"
                    >Condition (provided by you)</v-toolbar
                  >
                  <v-card-text class="text-left">
                    <div class="text-p pa-2">
                      <v-icon left small> mdi-star </v-icon
                      ><v-icon left small> mdi-star-outline </v-icon
                      ><v-icon left small> mdi-star-outline </v-icon
                      ><v-icon left small> mdi-star-outline </v-icon
                      ><v-icon left small> mdi-star-outline </v-icon>
                      Bad
                    </div>
                    <div class="text-p pa-2">
                      <v-icon left small> mdi-star </v-icon
                      ><v-icon left small> mdi-star </v-icon
                      ><v-icon left small> mdi-star-outline </v-icon
                      ><v-icon left small> mdi-star-outline </v-icon
                      ><v-icon left small> mdi-star-outline </v-icon>Fixable
                    </div>
                    <div class="text-p pa-2">
                      <v-icon left small> mdi-star </v-icon
                      ><v-icon left small> mdi-star </v-icon
                      ><v-icon left small> mdi-star </v-icon
                      ><v-icon left small> mdi-star-outline </v-icon
                      ><v-icon left small> mdi-star-outline </v-icon>
                      Good
                    </div>
                    <div class="text-p pa-2">
                      <v-icon left small> mdi-star </v-icon
                      ><v-icon left small> mdi-star </v-icon
                      ><v-icon left small> mdi-star </v-icon
                      ><v-icon left small> mdi-star </v-icon
                      ><v-icon left small> mdi-star-outline </v-icon>
                      As New
                    </div>
                    <div class="text-p pa-2">
                      <v-icon left small> mdi-star </v-icon
                      ><v-icon left small> mdi-star </v-icon
                      ><v-icon left small> mdi-star </v-icon
                      ><v-icon left small> mdi-star </v-icon
                      ><v-icon left small> mdi-star </v-icon>
                      Perfect
                    </div>
                  </v-card-text>
                  <v-card-actions class="justify-end">
                    <v-btn text @click="dialog.value = false">Close</v-btn>
                  </v-card-actions>
                </v-card>
              </template>
            </v-dialog>

            <v-chip
              v-if="thisitem.localpickup"
              class="ma-1 font-weight-light"
              target="_blank"
              :href="
                'https://www.google.com/maps/search/?api=1&query=' +
                thisitem.localpickup
              "
              outlined
              ><v-icon left> mdi-map-marker-outline </v-icon> Pickup
              Location</v-chip
            >

            <v-chip
              :to="{ name: 'SearchRegion', params: { region: country } }"
              outlined
              class="ma-1 font-weight-light text-uppercase"
              v-for="country in thisitem.shippingregion"
              :key="country"
            >
              <v-icon small left> mdi-flag-variant-outline </v-icon
              >{{ country }}</v-chip
            >

            <v-chip
              :to="{ name: 'SearchTag', params: { tag: tag } }"
              outlined
              class="ma-1 font-weight-light text-capitalize"
              v-for="tag in thisitem.tags"
              :key="tag"
            >
              <v-icon small left> mdi-tag-outline </v-icon>{{ tag }}</v-chip
            >
            <v-card class="ma-1 rounded-t-xl" outlined>
              <v-list dense disabled>
                <v-subheader>About</v-subheader>
                <v-list-item-group>
                  <v-list-item >
                    <v-list-item-icon>
                      <v-icon >mdi-account-badge-outline </v-icon>
                    </v-list-item-icon>

                    <v-list-item-content>
                      <v-list-item-title class="font-weight-light"
                        ><v-row
                          ><v-col>TPP ID: </v-col>
                          <v-col>{{ thisitem.id }}</v-col></v-row
                        ></v-list-item-title
                      >
                    </v-list-item-content>
                  </v-list-item>
                  <v-list-item v-if="thisitem.creator != thisitem.seller">
                    <v-list-item-icon>
                      <v-icon> mdi-account-outline</v-icon>
                    </v-list-item-icon>

                    <v-list-item-content>
                      <v-list-item-title class="font-weight-light"
                        ><v-row
                          ><v-col
                            >Original Seller: {{ thisitem.creator }}</v-col
                          ></v-row
                        ></v-list-item-title
                      >
                    </v-list-item-content>
                    
                  </v-list-item>
                    <v-list-item  v-if="thisitem.buyer">
                    <v-list-item-icon>
                      <v-icon> mdi-shopping</v-icon>
                    </v-list-item-icon>

                    <v-list-item-content>
                      <v-list-item-title class="font-weight-light"
                        ><v-row
                          ><v-col
                            >Buyer: {{ thisitem.buyer }}</v-col
                          ></v-row
                        ></v-list-item-title
                      >
                    </v-list-item-content>
                    
                  </v-list-item>

                  <v-list-item  v-if="thisitem.status">
                    <v-list-item-icon>
                      <v-icon> mdi-tune</v-icon>
                    </v-list-item-icon>

                    <v-list-item-content>
                      <v-list-item-title class="font-weight-light"
                        ><v-row
                          ><v-col
                            >Status: {{ thisitem.status }}</v-col
                          ></v-row
                        ></v-list-item-title
                      >
                    </v-list-item-content>
                    
                  </v-list-item>

                   <v-list-item  v-if="thisitem.transferable && thisitem.buyer === '' ">
                    <v-list-item-icon>
                      <v-icon> mdi-swap-horizontal</v-icon>
                    </v-list-item-icon>

                    <v-list-item-content>
                      <v-list-item-title class="font-weight-light"
                        ><v-row
                          ><v-col
                            > <v-icon left> mdi-store </v-icon>Transferable</v-col
                          ></v-row
                        ></v-list-item-title
                      >
                    </v-list-item-content>
                    
                  </v-list-item>
                  
               
                  <v-list-item v-if="thisitem.shippingcost > 0">
                    <v-list-item-icon>
                      <v-icon> mdi-package-variant-closed </v-icon>
                    </v-list-item-icon>

                    <v-list-item-content>
                      <v-list-item-title class="font-weight-light"
                        ><v-row
                          ><v-col>Shipping Cost: </v-col>
                          <v-col
                            >{{ thisitem.shippingcost
                            }}<v-icon small right
                              >$vuetify.icons.custom</v-icon
                            ></v-col
                          ></v-row
                        ></v-list-item-title
                      >
                    </v-list-item-content>
                  </v-list-item>
                  <v-list-item v-if="thisitem.seller != thisitem.creator">
                    <v-list-item-icon>
                      <v-icon> mdi-repeat</v-icon>
                    </v-list-item-icon>

                    <v-list-item-content>
                      <v-list-item-title class="font-weight-light"
                        ><v-row
                          ><v-col>Original Price: </v-col>
                          <v-col
                            >{{ thisitem.estimationprice
                            }}<v-icon small right
                              >$vuetify.icons.custom</v-icon
                            ></v-col
                          ></v-row
                        ></v-list-item-title
                      >
                    </v-list-item-content>
                  </v-list-item>
                  <v-list-item v-else-if="thisitem.estimationprice > 0">
                    <v-list-item-icon>
                      <v-icon> mdi-check-all </v-icon>
                    </v-list-item-icon>

                    <v-list-item-content>
                      <v-list-item-title class="font-weight-light"
                        ><v-row
                          ><v-col>Estimation Price: </v-col>
                          <v-col
                            >{{ thisitem.estimationprice
                            }}<v-icon small right
                              >$vuetify.icons.custom</v-icon
                            ></v-col
                          ></v-row
                        ></v-list-item-title
                      >
                    </v-list-item-content>
                  </v-list-item>
                  <v-list-item v-if="thisitem.discount > 0">
                    <v-list-item-icon>
                      <v-icon> mdi-brightness-percent</v-icon>
                    </v-list-item-icon>

                    <v-list-item-content>
                      <v-list-item-title class="font-weight-light"
                        ><v-row
                          ><v-col>Discount: </v-col>
                          <v-col
                            >{{ thisitem.discount
                            }}<v-icon small right
                              >$vuetify.icons.custom</v-icon
                            ></v-col
                          ></v-row
                        ></v-list-item-title
                      >
                    </v-list-item-content>
                  </v-list-item>
                </v-list-item-group>
              </v-list>
            </v-card>
             <div class="overline text-center">Comments</div>
              <div v-if="commentlist">
                <div v-for="(comment, nr) in commentlist" v-bind:key="nr">
                  <v-chip color="primary" class="ma-2">{{ comment }} </v-chip>
                </div>
              </div>
              <div v-else>
                <p class="caption text-center">No comments to show right now</p>
              </div>
            </v-col>
          </v-row>
        </div>

        <v-card-actions>
          <v-btn rounded
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

        <v-expand-transition >
          <div class="pa-2 mx-auto" elevation="8" v-if="showactions">
            <v-divider></v-divider>
            <div>
              <v-stepper class="elevation-0" v-model="step" vertical>
                <v-stepper-step step="1" complete> Place Item </v-stepper-step>

                <v-stepper-step
                  :complete="
                    thisitem.bestestimator != '' || thisitem.status != '' ||thisitem.discount
                  "
                  step="2"
                >
                  Awaiting Estimations
                </v-stepper-step>

                <v-stepper-content step="2">
                 <p type="caption">
                    Awaiting estimations. Meanwhile... help
                    others by estimating other items (and earn tokens)!
                  </p>
                </v-stepper-content>

                <v-stepper-step :complete="thisitem.transferable || thisitem.status" step="3">
                   Reveal estimation
                </v-stepper-step>
                

                <v-stepper-content step="3">
               
                  <div v-if="thisitem.bestestimator != 'Awaiting'">
                    <p type="caption">
                      Wow! there is a final price. You can sell
                      {{ thisitem.title }} for {{
                        thisitem.estimationprice
                      }}
                      <v-icon small right>$vuetify.icons.custom</v-icon> tokens. By accepting, your item will be available to buy. Anyone can provide a prepayment to buy the
                      item.
                    </p>
                    <v-row>
                      <v-btn rounded
                        class="ma-4"
                        color="primary"
                        v-if="
                          !flightit &&
                          hasAddress &&
                          thisitem.bestestimator != 'Awaiting' &&
                          thisitem.transferable != true
                        "
                        @click="submititemtransferable(true, thisitem.id)"
                        ><v-icon left> mdi-checkbox-marked-circle </v-icon>
                        Accept
                        <div class="button__label" v-if="flightit">
                          <div class="button__label__icon">
                            <icon-refresh />
                          </div>
                          Placing item for $ale...
                        </div>
                      </v-btn>

                      <v-btn rounded
                        class="ma-4" 
                        color="default"
                        v-if="
                          !flightitn &&
                          hasAddress &&
                          thisitem.bestestimator != 'Awaiting' &&
                          thisitem.transferable != true
                        "
                        @click="submititemtransferable(false, thisitem.id)"
                        ><v-icon left> mdi-close </v-icon>
                        Reject
                        <div class="button__label" v-if="flightitn">
                          <div class="button__label__icon">
                            <icon-refresh />
                          </div>
                          Deleting item...
                        </div>
                      </v-btn>
                    </v-row>
                  </div>   <div v-else>
 <p type="caption">
                      The estimations came in. You can now reveal the final estimation price, after which you can accept or decline this price.
                 <v-row>  <v-btn rounded
                        class="ma-4"
                        color="primary"
                        v-if="
               
                          hasAddress &&
                          thisitem.bestestimator == 'Awaiting' 
                        "
                        @click="submitRevealEstimation()"
                        ><v-icon left> mdi-checkbox-marked-circle </v-icon>
                      
                        <div class="button__label" v-if="flightre">
                      
                          Revealing estimation...
                        </div><span v-else>Reveal</span>
                      </v-btn></v-row>
                    </p>
                    </div>
                </v-stepper-content>

                <v-stepper-step :complete="thisitem.buyer != ''" step="4">
                  Item For Sale
                </v-stepper-step>

                <v-stepper-content step="4" :complete="thisitem.status != ''">
                  <p type="subtitle"
                    >Item placed. Awaiting buyer... Tip: share your item with
                    family and friends. </p>
                
                <v-icon small>mdi-share-variant </v-icon> <input v-model="tocopy" size=50 class="mx-2 caption" type="text" ref="input" >   <v-btn text @click="copyText()">  Copy</v-btn>
                  <p
                    v-if="thisitem.shippingcost > 0 && thisitem.localpickup != '' "
                    type="caption"
                  >
                    If a buyer chooses shipping, you ship it, provide the track
                    and trace code if available, and you'll automatically get
                    your tokens. After a buyer is found and chooses local
                    pickup, the buyer can pick it up. Tip: let the buyer
                    transfer the tokens during your meetup.
                  </p>
                  <p
                    v-if="thisitem.shippingcost === 0 && thisitem.localpickup != ''"
                    type="caption"
                  >
                    After a buyer is found negotiate a meetup time and place by
                    sending a message to the buyer. Tip: let the buyer transfer
                    the tokens during your meetup.
                  </p>
                  <p
                    v-if="
                      thisitem.shippingcost > 0 &&
                      thisitem.localpickup == ''
                    "
                    type="caption"
                  >
                    After a buyer is found, find out about the address to ship
                    to by sending a message to the buyer.
                  </p>
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
                      thisitem.localpickup == '' &&
                      thisitem.buyer &&
                      thisitem.shippingcost &&
                      thisitem.status === ''
                    "
                  >
                  
                    <p type="caption">
                      Now it's time to ship the item. Provide a track and trace
                      code to the buyer if available.
                    </p>
                    <input
                      type="checkbox"
                      id="checkbox"
                      v-model="tracking"
                      v-bind:value="true"
                    />
                    <label for="checkbox">
                      I have shipped the item and provided the buyer with track
                      and trace
                    </label>
                    <v-btn rounded @click="submitItemShipping(tracking, thisitem.id)">
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
                      thisitem.localpickup != '' &&
                      thisitem.buyer &&
                      thisitem.status === ''
                    "
                  >
                    

                    <p type="caption">
                      Now its time to meet up with the buyer. 
                    </p>
                  </div>
                         <div class="justify-end">
                  <v-btn rounded
        :disabled="!this.$store.state.account.address"
        text 
        @click="createRoom"
      ><v-icon> mdi-reply</v-icon>
        Message Buyer</v-btn
      ></div>
                </v-stepper-content>

                <v-stepper-step :complete="thisitem.status != ''" step="6">
                  Complete
                </v-stepper-step>

                <v-stepper-content step="5" height="200px"
                  ><v-card></v-card>
                </v-stepper-content>
              </v-stepper>
            </div>
          </div> </v-expand-transition
      ></v-card><sign-tx v-if="submitted" :key="submitted" :fields="fields" :value="value" :msg="msg" @clicked="afterSubmit"></sign-tx>
    </div>
  </div>
</template>

<script>
import { usersRef, roomsRef, databaseRef } from "./firebase/db.js";

import ItemListSeller from "./ItemListSeller.vue";


export default {
  props: ["itemid"],
  components: { ItemListSeller },
  data() {
    return {
     


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
        shippingRules: [
          (v) =>
            !!v.length == 1 ||
            "A country is required when shipping cost is applicable",
        ],
      },
      countryCodes: ["NL", "BE", "UK", "DE", "US", "CA"],

        fields: [],
      value: {},
      msg: "",
      submitted: false,
    };
  },

  mounted() {
    this.loadingitem = true;
    const id = this.itemid;

    const imageRef = databaseRef.ref("ItemPhotoGallery/" + id + "/photos/");
    imageRef.on("value", (snapshot) => {
      const data = snapshot.val();

      if (data != null ) {
        //console.log(data[0]);
        this.photos = data;
        this.imageurl = data[0];
        this.loadingitem = false;
      }
    });
    this.loadingitem = false;
  },

  computed: {
   
    thisitem(){
      return this.$store.getters.getItemByID(this.itemid);
    },

    hasAddress() {
      return !!this.$store.state.account.address;
    },
    valid() {
      return this.thisitem.id.trim().length > 0;
    },
    tocopy(){
      return process.env.VUE_APP_URL + "/itemid=" + this.thisitem.id
    },

    commentlist() {
      //const item = this.$store.getters.getItemByID(this.itemid);

      //console.log(this.thisitem);
      return this.thisitem.comments.filter((com) => com != "") || [];
      //console.log( this.thisitem.comments.filter(i => i != ""));
    },

  
  },

  methods: {

    
    async removeItem() {
      this.loadingitem = true;
      this.flightre = true;
      const type = { type: "item" };
      const body = { id: this.thisitem.id };
      this.fields = [
        ["seller", 1, "string", "optional"],
        ["id", 2, "string", "optional"],
      ];

      this.msg = "MsgDeleteItem"

    
  this.value = {
          seller: this.$store.state.account.address,
          ...body,
        },
  
  this.submitted = true

    },


      async afterSubmit(value) {
 this.loadingitem = true;

 this.msg = ""
 this.fields = []
 this.value = {}
  if(value == true){
             await this.$store.dispatch("updateItem", this.thisitem.id)//.then(result => this.newitem = result)
        
        await this.$store.dispatch("bankBalancesGet");
         setTimeout( () => this.$router.push("/itemid="+ this.thisitem.id), 5000) }



          this.submitted = false
               this.flightre = false;
      this.loadingitem = false;
         this.flightIS = false;
        this.tracking = false;

      this.flightit = false
      this.flightitn = false

    },


 async submitRevealEstimation() {
      if (this.hasAddress) {
       this.flightre = true;
      
        this.fields = [
          ["creator", 1, "string", "optional"],
          ["itemid", 2, "string", "optional"],
        ];
        this.msg = "MsgRevealEstimation"
        // const type = { type: "item" };
        const body = { itemid: this.thisitem.id };
      
 this.value = {
          creator: this.$store.state.account.address,
          ...body,
        }
      
        this.submitted = true}
    },
    

    async getThisItem() {
      await submitrevealestimation();
      return this.thisitem();
    },
    async submititemtransferable(transferable, itemid) {
      if (this.valid && !this.flightit && this.hasAddress) {
        this.flightit = true;
        this.flightitn = true;

        this.msg = "MsgItemTransferable"
        this.fields = [
          ["seller", 1, "string", "optional"],
          ["transferable", 2, "bool", "optional"],
          ["itemid", 3, "string", "optional"],
        ];
 const body = { transferable, itemid };
        this.value = {
          seller: this.$store.state.account.address,
          ...body,
        }
    
      }  
      this.submitted = true
    },
   
    async submitItemShipping(tracking, itemid) {
      if (this.valid && !this.flightIS && this.hasAddress) {
        this.flightIS = true;

        const body = { tracking, itemid };
        this.fields = [
          ["seller", 1, "string", "optional"],
          ["tracking", 2, "bool", "optional"],
          ["itemid", 3, "string", "optional"],
        ];
     

     this.msg = "MsgItemShipping"

         this.value = {
          seller: this.$store.state.account.address,
          ...body,
        }

          this.submitted = true
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
    
     copyText() {
      const copyText = this.$refs.input;
  copyText.select();
  document.execCommand('copy');
    },


    async createRoom() {

      if (this.$store.state.user.uid) {

        const user = await usersRef.where('username', '==' , this.thisitem.buyer).get();

 let query =  roomsRef.where("users", "==", this.$store.state.user.uid, this.thisitem.buyer)
  console.log(query)
if (user && !query) {
      //await usersRef.doc(id).update({ _id: id });
      await roomsRef.add({
        users: [user.docs[0].id, this.$store.state.user.uid],
        lastUpdated: new Date(),
      });

      this.addNewRoom = false;
      this.addRoomUsername = "";
      this.fetchRooms();
     }else{
      alert("Buyer already added or buyer not found")
    }; 
       this.$router.push('/messages') } else{ alert("Sign in first (Check your Google email)")}
      
    },
  },

  
};
</script>

