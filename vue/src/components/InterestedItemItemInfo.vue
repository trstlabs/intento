<template>
  <div class="pa-2 mx-auto"  >
    <v-card elevation="2" rounded="lg">
      <v-progress-linear
        indeterminate
        :active="loadingitem"
      ></v-progress-linear>

      <div class="pa-2 mx-auto">
       
            <p class="pa-2 overline"> {{ thisitem.title }} </p>
          
            
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
              <span >  <div class="pl-4 overline text-center">Description</div> <v-card-text>
    
     
  <div class="caption ">
           {{thisitem.description }}
         </div> </v-card-text> </span>
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
                   <v-list-item>
                    <v-list-item-icon>
                      <v-icon> mdi-clock</v-icon>
                    </v-list-item-icon>

                    <v-list-item-content>
                      <v-list-item-title class="font-weight-light"
                        ><v-row
                          ><v-col>Created on: </v-col>
                          <v-col>{{ getFmtTime(thisitem.submittime) }}</v-col></v-row
                        ></v-list-item-title
                      >
                    </v-list-item-content>
                  </v-list-item>
                   <v-list-item v-if="!thisitem.buyer && thisitem.status == ''">
                    <v-list-item-icon>
                      <v-icon> mdi-progress-clock</v-icon>
                    </v-list-item-icon>

                    <v-list-item-content>
                      <v-list-item-title class="font-weight-light"
                        ><v-row
                          ><v-col>Expires on: </v-col>
                          <v-col>{{ getFmtTime(thisitem.endtime) }}</v-col></v-row
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
              </v-list></v-card>
            </v-col>
        </v-row>
      </div>
      <v-card-actions>
        <v-btn
          color="blue"
          text
          @click="(showinfo = !showinfo), getItemPhotos()"
        >
          Actions
        </v-btn>

        <v-spacer></v-spacer>

        <v-btn icon @click="(showinfo = !showinfo), getItemPhotos()">
          <v-icon>{{
            showinfo ? "mdi-chevron-up" : "mdi-chevron-down"
          }}</v-icon>
        </v-btn>
      </v-card-actions>

      <v-expand-transition>
        <div>
          <div class="pa-2 mx-auto" elevation="8" v-if="showinfo">
            <div>
             
      
             

              <v-row> <v-col>
            <v-btn rounded block color="primary"
              v-if="thisitem.localpickup != ''"
              @click="submitLP(itemid), getThisItem"
            >
              Buy Item
              <div class="button__label" v-if="flight">
                <div class="button__label__icon">
                  <icon-refresh />
                </div>
                Sending transaction...
              </div>
            </v-btn>
            </v-col><v-col>
            <v-btn rounded block color="primary"
              v-if="thisitem.shippingcost"
              @click="submitSP(itemid), getThisItem"
            >
              Buy item + shipping
              <div class="button__label" v-if="flight">
                <div class="button__label__icon">
                  <icon-refresh />
                </div>
                Sending transaction...
              </div>
            </v-btn>
            </v-col><v-col>
            <v-btn block color="warning"
              rounded
              @click="submitInterest(itemid), getThisItem"
            >
              Unlike Item
              <div class="button__label" v-if="flight">
                <div class="button__label__icon">
                  <icon-refresh />
                </div>
                Unliking item...
              </div>
            </v-btn>
            </v-col></v-row>

            <div v-if="thisitem.buyer != ''">
              <p>Item buyer is {{ thisitem.buyer }}</p>
            </div>
            <div>
              <!-- <router-link
                :to="{ name: 'BuyItemDetails', params: { id: itemid } }"
                >Full Details (loads new page)
              </router-link> -->
            </div>
            </div>
          </div>
        </div>
      </v-expand-transition>
    </v-card><sign-tx v-if="submitted" :key="submitted" :fields="fields" :value="value" :msg="msg" @clicked="afterSubmit"></sign-tx>
  </div>
</template>

<script>
import { databaseRef } from './firebase/db';
import ItemListInterested from "./ItemListInterested.vue";
import dayjs from 'dayjs'

export default {
  props: ["itemid"],
  components: { ItemListInterested },
  data() {
    return {
      //itemid: this.item.id,
      //make sure deposit is number+token before sending tx
      amount: "",
      flight: false,
      flightLP: false,
      flightSP: false,
      showinfo: false,
      imageurl: "",
      loadingitem: true,
      photos: [],
      
        fields: [],
      value: {},
      msg: "",
      submitted: false,
    };
  },

  mounted() {
    this.loadingitem = true;

    const imageRef = databaseRef.ref("ItemPhotoGallery/" + this.itemid + "/photos/");
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
    thisitem() {
      //console.log(this.itemid)
      return this.$store.getters.getItemByID(this.itemid);
    },
   
    hasAddress() {
      return !!this.$store.state.account.address;
    },
    valid() {
      return this.amount.trim().length > 0;
    },
  },

  methods: {
        getFmtTime(time) {
      	const momentTime = dayjs(time)
      return momentTime.format("D MMM, YYYY HH:mm:ss");
    },

     async afterSubmit(value){
 this.loadingitem = true;

 this.msg = ""
 this.fields = []
 this.value = {}
  if(value == true){
             this.$store.dispatch("updateItem", this.thisitem.id)//.then(result => this.newitem = result)
         this.$router.push("/itemid="+ this.thisitem.id)}
        await this.$store.dispatch("bankBalancesGet");


          this.submitted = false
 
        this.loadingitem = false;  
   
     

    },
   

    async submitLP(itemid) {
      if (!this.flightLP && this.hasAddress) {
        this.flightLP = true;
        this.loadingitem = true;
        let toPay = this.thisitem.estimationprice;
        let deposit = toPay + "token";
       
        const body = { deposit, itemid };
          this.fields = [
          ["buyer", 1, "string", "optional"],

          ["itemid", 2, "string", "optional"],
          ["deposit", 3, "int64", "optional"],
        ];
        this.msg = "MsgPrepayment"
       this.value = {
          buyer: this.$store.state.account.address,
          ...body
        }
  this.submitted = true
      }
    },

    

    async submitSP(itemid) {
      if (!this.flightSP && this.hasAddress) {
        this.flightSP = true;
        this.loadingitem = true;
    
      
        let toPaySP =
          +this.thisitem.estimationprice + +this.thisitem.shippingcost;
       
        let deposit = toPaySP + "token";


        const body = { deposit, itemid };
          this.fields = [
          ["buyer", 1, "string", "optional"],

          ["itemid", 2, "string", "optional"],
          ["deposit", 3, "int64", "optional"],
        ];
        this.msg = "MsgPrepayment"
       this.value = {
          buyer: this.$store.state.account.address,
          ...body
        }
           this.submitted = true
      }
    },



     async submitInterest(itemid) {
      if (!this.flightLP && this.hasAddress) {
    //    this.flightLP = true;
        this.loadingitem = true;
          this.fields = [
          ["estimator", 1, "string", "optional"],

          ["itemid", 2, "string", "optional"],
          ["interested", 3, "bool", "optional"],
        ];
const body = { itemid: itemid,
        interested: false };
this.msg = "MsgUpdateLike"
this.value ={
          estimator: this.$store.state.account.address,
          ...body
        }

  
  this.submitted = true
      }
    },

    async getThisItem() {
      await submit();
      return thisitem();
    },

    getItemPhotos() {
      if (this.showinfo && this.imageurl != "") {
        this.loadingitem = true;
        const id = this.itemid;
 

        const imageRef = databaseRef.ref("ItemPhotoGallery/" + id + "/photos/");
        imageRef.on("value", (snapshot) => {
          const data = snapshot.val();
          if (data != null ) {
            this.photos = data;
            this.loadingitem = false;
          }
        });
        this.loadingitem = false;
      }
    },
  },
};
</script>
