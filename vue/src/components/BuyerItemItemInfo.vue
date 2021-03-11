<template>
  <div>
   <div class="pa-2 mx-auto">
      <v-card elevation="2" rounded="lg">
        <v-progress-linear
          indeterminate
          :active="loadingitem"
        ></v-progress-linear>
        <div class="pa-2 mx-auto">
          <p class="pa-2 h3 font-weight-medium "> {{ thisitem.title }} </p>
          
            
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
            <v-divider class="ma-2" />
         
<v-row align="start">
            <v-col>
               

 
  <v-card elevation="0" >  <div class="pl-4 overline text-center">Description</div> <v-card-text>
    
     
  <div class="body-1 "> "
           {{thisitem.description }} "
         </div> </v-card-text> </v-card>

  <v-divider class="ma-2" />
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

    

<v-chip v-if="thisitem.localpickup"
      class="ma-1"
      label
      outlined
      medium
    ><v-icon left> 
        mdi-map-marker
      </v-icon>Local Pickup</v-chip>

      <v-chip v-if="thisitem.shippingcost > 0"
      class="ma-1"
      label
      outlined
      medium
    >
    <v-icon left>
        mdi-package-variant-closed
      </v-icon>
      Shipping Cost: $ {{thisitem.shippingcost}} tokens
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
      Price: $ {{thisitem.estimationprice}} tokens
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
      Item Transferable
    </v-chip>

<v-chip v-if="thisitem.status"
      class="ma-2"
      label
      outlined
      medium
    >
    <v-icon left>
        mdi-clock-time-three-outline
      </v-icon>
      Status: {{ thisitem.status }}
    </v-chip>

     <v-chip
      class="mt-1"
      medium label outlined
    >
    <v-icon left>
        mdi-account-outline
      </v-icon>
      Seller: {{ thisitem.creator }}
    </v-chip>

    

             
            </v-col>

            <!--<v-col cols="4">
              <div v-if="imageurl" class="d-flex flex-row-reverse">
                <v-avatar class="ma-4" size="125" rounded="lg" tile>
                  <v-img :src="imageurl"></v-img>
                </v-avatar>
              </div>
            </v-col>-->
          </v-row>
          
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
                <v-stepper-step step="1" complete> Provide prepayment Item </v-stepper-step>                
                

                <v-stepper-step :complete="thisitem.status != ''" step="2">
                  Item Transfer
                </v-stepper-step>

                <v-stepper-content step="2">
                  <div v-if="thisitem.tracking === true">
                    
        <app-text type="p">This item has shipped </app-text>
        <app-text type="p"
                    >Item has been shipped. Item seller indicated that item is shipped. For more information contact the seller. The protocol has received the request of the seller to arrange tranfer coins.
                  </app-text>
      </div>

       <div v-if="thisitem.localpickup === false && !thisitem.status">
                    
        <app-text type="p">This item is not shipped yet</app-text>
        <app-text  type="p"
                    >Contact the seller of {{thisitem.title}}. Item seller will indicate if the item is shipped. 
                  </app-text>
      </div>
                  
                  <div>
       <div v-if="thisitem.localpickup === true && thisitem.status != 'Item transferred'">  
         <app-text class="ma-2" type="p"> Arrange a meeting to pick up the item.   </app-text>           
         <v-row>
           
        <v-btn class="ma-4" color="primary"
           

          @click="submitItemTransfer(true, thisitem.id), getThisItem"
        ><v-icon left>
         mdi-checkbox-marked-circle
      </v-icon>
          Complete Transfer
          <div class="button__label" v-if="flightIT">
            <div class="button__label__icon">
              <icon-refresh />
            </div>
            Sending tokens to seller...
          </div>
        </v-btn>
      
      
        <v-btn class="ma-4" color="default"
          :class="[
            'button',
            `button__valid__${!!valid && !flightITN && hasAddress}`,
          ]"
          @click="submitItemTransferN(false, thisitem.id), getThisItem"
        ><v-icon left>
         mdi-cancel
      </v-icon>
          Cancel transfer
          <div class="button__label" v-if="flightITN">
            <div class="button__label__icon">
              <icon-refresh />
            </div>
            Sending tokens back...
          </div>
        </v-btn>
         </v-row>
        </div>
      </div>
                </v-stepper-content>
                <v-stepper-step :complete="thisitem.status != ''" step="3">
                  Item Transferred
                </v-stepper-step>
              </v-stepper>
            </div>
          </div> </v-expand-transition
      ></v-card>
    </div>

  </div>
</template>

<script>
import ItemListBuyer from "./ItemListBuyer.vue";
import {databaseRef} from "./firebase/db.js"
export default {
  props: ["itemid"],
  components: { ItemListBuyer },
  data() {
    return {
      flightIT: false,
      flightITN: false,
      loadingitem: true,
      showactions: false,
      transferbool: false,
      photos: [],
      imageurl: "",
      step: 2,
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
      return this.$store.getters.getItemByID(this.itemid);
    },
    hasAddress() {
      return !!this.$store.state.account.address;
    },
    valid() {
      return this.thisitem.id.trim().length > 0;
    },
  },

  methods: {
   
    async submitItemTransfer(transferable, itemid) {

      if (this.valid && !this.flightIT && this.hasAddress) {
        this.flightIT = true;
        const type = { type: "buyer" };
        const body = { transferable, itemid };
         const fields = [
        ["buyer", 1,'string', "optional"],
         [ "itemid", 2,'string', "optional"] ,                                                    
        ["transferable",3,'bool', "optional"],
      ];
        await this.$store.dispatch("transferSubmit", { ...type, body, fields });
 await this.$store.dispatch("setBuyerItemList", this.$store.state.account.address);
        this.flightIT = false;
    
      }
    },

    async submitItemTransferN(transferable, itemid) {
      if (this.valid && !this.flightITN && this.hasAddress) {
        this.flightITN = true;
        const type = { type: "buyer" };
        const body = { transferable, itemid };
        await this.$store.dispatch("entitySubmit", { ...type, body, fields });
        await this.$store.dispatch("entityFetch", type);
        await this.$store.dispatch("setBuyerItemList", this.$store.state.account.address);
        this.flightITN = false;
     
      }
    },

    async getThisItem() {
      await submitrevealestimation();
      return this.thisitem();
      console.log(this.thisitem);
    },
    createStep() {
      if (this.thisitem.tracking != "") {
        this.step = 2;
      } else if (this.thisitem.transferable === true) {
        this.step = 2;
      } else if (this.thisitem.status != "") {
        this.step = 3;
      }
    },
  },
};
</script>

