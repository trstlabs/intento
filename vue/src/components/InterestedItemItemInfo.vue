<template>
  <div class="pa-2 mx-auto"  >
    <v-card elevation="2" rounded="lg">
      <v-progress-linear
        indeterminate
        :active="loadingitem"
      ></v-progress-linear>

      <div class="pa-2 mx-auto">
       
          <p class="pa-2 h3 font-weight-medium "> {{ thisitem.title }} </p>
          
            
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

 <v-chip
      class="ma-1"
      label
      outlined
      medium
    >
    <v-icon left>
        mdi-account-badge
      </v-icon>
      TPP ID: {{ thisitem.id }}
    </v-chip>

<v-chip v-if="thisitem.localpickup != ''"
      class="ma-1"
      label
      outlined
      medium
    >
    <v-icon left>
        mdi-pin
      </v-icon>
      Pickup
    </v-chip>
       
    
          <v-chip v-if="thisitem.shippingcost"
      class="ma-1"
      label
      outlined
      medium
    >
    <v-icon left>
        mdi-package-varient-closed
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
      Shipping cost: {{ thisitem.shippingcost }}<v-icon small right>$vuetify.icons.custom</v-icon>  
    </v-chip>

    <v-chip v-if="thisitem.estimationprice > 0"
      class="ma-1"
      label
      outlined
      medium
    >
    <v-icon left>
        mdi-check-all
      </v-icon>
      Estimation Price: $ {{thisitem.estimationprice}} tokens
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
        this.msg = "MsgCreateBuyer"
       this.value = {
          buyer: this.$store.state.account.address,
          ...body
        }

       this.submitted = false
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
        this.msg = "MsgCreateBuyer"
       this.value = {
          buyer: this.$store.state.account.address,
          ...body
        }
          this.submitted = false
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
this.msg = "MsgUpdateEstimator"
this.value ={
          estimator: this.$store.state.account.address,
          ...body
        }

  

          this.submitted = false
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
