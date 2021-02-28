<template>
<div class="pa-2 mx-auto">
    <v-card elevation="2" rounded="lg" v-click-outside="clickOutside" >
      <v-progress-linear
        indeterminate
        :active="loadingitem"
      ></v-progress-linear>

            
       <v-row>
         <v-col cols="12" md="8">
            <h4 class="text-capitalize pa-2 text-left">{{ thisitem.title }}</h4>
         

          
          
          </v-col>

          <v-col cols="12">
            <div v-if="imageurl" >
              
                <v-img  class="rounded contain" :src="imageurl"></v-img>
             
            </div>
          </v-col>
        </v-row>
    
        <div>
          <div class="pa-2 mx-auto" elevation="8"  >
            <div>
              <div v-if="photos.photo">
                <v-divider></v-divider>
                <v-carousel
                  cycle
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
             


 <v-card elevation="0" >  <div class="pa-2 overline text-center">Description</div> <v-card-text>
    
     
  <div class="body-1 "> "
           {{thisitem.description }} "
         </div> </v-card-text> </v-card>


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

  
<v-dialog transition="dialog-bottom-transition"
        max-width="300"> <template v-slot:activator="{ on, attrs }">
        <span
          v-bind="attrs"
          v-on="on"
        >
        <v-chip
      class="ma-1"
      label
      outlined
      medium

    >
    <v-icon left small>
        mdi-star-outline
      </v-icon>
     {{thisitem.condition}}/5
    </v-chip> </span> </template> <template v-slot:default="dialog">
          <v-card>
            <v-toolbar 
              color="default"
              
            >Condition (provided by seller)</v-toolbar>
            <v-card-text class="text-left">
           
              <div class="text-p pa-2">
                <v-icon left small>
        mdi-star
      </v-icon><v-icon left small>
        mdi-star-outline
      </v-icon><v-icon left small>
        mdi-star-outline
      </v-icon><v-icon left small>
        mdi-star-outline
      </v-icon><v-icon left small>
        mdi-star-outline
      </v-icon>
                Bad 
                 
              </div>
              <div class="text-p pa-2">
                <v-icon left small>
        mdi-star
      </v-icon><v-icon left small>
        mdi-star
      </v-icon><v-icon left small>
        mdi-star-outline
      </v-icon><v-icon left small>
        mdi-star-outline
      </v-icon><v-icon left small>
        mdi-star-outline
      </v-icon>Fixable 
                 
              </div>
              <div class="text-p pa-2">
                <v-icon left small>
        mdi-star
      </v-icon><v-icon left small>
        mdi-star
      </v-icon><v-icon left small>
        mdi-star
      </v-icon><v-icon left small>
        mdi-star-outline
      </v-icon><v-icon left small>
        mdi-star-outline
      </v-icon>
                Good 
                 
              </div>
              <div class="text-p pa-2">
                <v-icon left small>
        mdi-star
      </v-icon><v-icon left small>
        mdi-star
      </v-icon><v-icon left small>
        mdi-star
      </v-icon><v-icon left small>
        mdi-star
      </v-icon><v-icon left small>
        mdi-star-outline
      </v-icon>
                As New 
                 
              </div><div class="text-p pa-2">
                <v-icon left small>
        mdi-star
      </v-icon><v-icon left small>
        mdi-star
      </v-icon><v-icon left small>
        mdi-star
      </v-icon><v-icon left small>
        mdi-star
      </v-icon><v-icon left small>
        mdi-star
      </v-icon>
                Perfect 
                 
              </div>
            </v-card-text>
            <v-card-actions class="justify-end">
              <v-btn
                text
                @click="dialog.value = false"
              >Close</v-btn>
            </v-card-actions>
            
          </v-card>
        </template>
    </v-dialog >
<v-chip v-if="thisitem.localpickup"
      class="ma-1"
      label
      outlined
      medium
    ><v-icon left> 
        mdi-map-marker-outline
      </v-icon>Local Pickup</v-chip>

      <v-chip v-if="thisitem.shippingcost"
      class="ma-1"
      label
      outlined
      medium
    >
    <v-icon left>
        mdi-package-variant-closed
      </v-icon>
      Shipping Cost: ${{thisitem.shippingcost}} TPP
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
      Price: ${{thisitem.estimationprice}} TPP
    </v-chip>

    <v-chip
      class="ma-1"
      medium label outlined
    >
    <v-icon left>
        mdi-account-outline
      </v-icon>
      Seller: {{ thisitem.creator }}
    </v-chip>
    
    <v-divider class="ma-2"/>
      
  <div class="overline text-center"> Comments </div> 
     <div v-if="thisitem.comments">
<v-chip  v-for="comment in commentlist" v-bind:key="comment" class="ma-2 "
        
    >{{ comment }}
     </v-chip>

     
     </div>
     <div v-if="!thisitem.comments">
<p  class="caption text-center"> No comments to show right now </p> </div>
     
          <v-divider class="ma-4"/>  
          <div v-if="hasAddress" class="ma-4 text-center">
          <wallet-coins /> </div>
             <div class="text-center"> <v-row> <v-col>
            <v-btn block color="primary"
              :disabled="!thisitem.localpickup"
              @click="submitLP(itemid), getThisItem"
            >
              Buy locally for ${{thisitem.estimationprice}} TPP<v-icon right> 
        mdi-map-marker
      </v-icon>
              <div class="button__label" v-if="flight">
                <div class="button__label__icon">
                  <icon-refresh />
                </div>
                Sending transaction...
              </div>
            </v-btn>
            </v-col><v-col>
            <v-btn block color="primary"
              :disabled="thisitem.shippingcost == 0"
              @click="submitSP(itemid), getThisItem"
            >
              Buy for ${{thisitem.estimationprice}} TPP + shipping (${{thisitem.shippingcost}} TPP)<v-icon right> 
        mdi-package-variant-closed
      </v-icon>
              <div class="button__label" v-if="flight">
                <div class="button__label__icon">
                  <icon-refresh />
                </div>
                Sending transaction...
              </div>
            </v-btn> </v-col>
            </v-row> </div> 

            <div v-if="thisitem.buyer != ''">
              <p>Item buyer is {{ thisitem.buyer }}</p>
            </div>
            
            </div>
          </div>
        </div>
       
         <v-card class="pa-2 mx-auto">
            <v-card-title> All Seller items </v-card-title>
     <div v-for="item in SellerItems" v-bind:key="item.id">
     
        
     
       <router-link
                :to="{ name: 'BuyItemDetails', params: { id: item.id } }"
                > {{item.title}}
      {{item.status}}
              </router-link>
  
    </div>    </v-card>
    </v-card>
  </div>
</template>
<script>
import BuyItemDetails from "../views/BuyItemDetails.vue";
export default {
  components: { BuyItemDetails },
  props: ["itemid"],
 
data() {
    return {
      amount: "",
      flight: false,
      flightLP: false,
      flightSP: false,
      showinfo: false,
      imageurl: "",
      loadingitem: true,
      photos: [],
      dialog: false,
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
        this.imageurl = data.photo;
        this.loadingitem = false;
      }
    });
    this.loadingitem = false;
  },
  computed: {
    thisitem() {
  
      return this.$store.getters.getItemByID(this.$route.params.id)  || [];
    },
   
    hasAddress() {
      return !!this.$store.state.account.address;
    },
    valid() {
      return this.amount.trim().length > 0;
    },
    commentlist() {
      return this.thisitem.comments.filter(i => i != "") || [];
      
    },
    SellerItems() { this.$store.dispatch("setSellerItemList", this.thisitem.creator);
    return this.$store.getters.getSellerList  || []
    },
  },

  methods: {

    async submitLP(itemid) {
      if (!this.hasAddress) {alert("Sign in first");};

      if (!this.flightLP && this.hasAddress) {
        this.flightLP = true;
        this.loadingitem = true;
        let toPay = this.thisitem.estimationprice;
        let deposit = toPay + "tpp";
        const type = { type: "buyer" };
        const body = { deposit, itemid };
        await this.$store.dispatch("entitySubmit", { ...type, body });
        await this.$store.dispatch("entityFetch", type);
        await this.$store.dispatch("accountUpdate");
        this.flightLP = false;
        this.loadingitem = false;
      }
    },

    

    async submitSP(itemid) {
       if (!this.hasAddress) {alert("Sign in first");};
      if (!this.flightSP && this.hasAddress) {
        this.flightSP = true;
        this.loadingitem = true;
        console.log("clicked");
         console.log(this.thisitem);
        console.log(this.thisitem.estimationprice);
        console.log(this.thisitem.shippingcost);
        let toPaySP =
          +this.thisitem.estimationprice + +this.thisitem.shippingcost;
        console.log(toPaySP);
        let deposit = toPaySP + "tpp";
        console.log(deposit);
        const type = { type: "buyer" };
        const body = { deposit, itemid };
        await this.$store.dispatch("entitySubmit", { ...type, body });
        await this.$store.dispatch("entityFetch", type);
        await this.$store.dispatch("accountUpdate");

        
        this.flightSP = false;
        this.loadingitem = false;
        this.deposit = "";
        alert("Transaction sent");
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
        const db = firebase.database();

        const imageRef = db.ref("ItemPhotoGallery/" + id);
        imageRef.on("value", (snapshot) => {
          const data = snapshot.val();
          if (data != null && data.photo != null) {
            this.photos = data;
            this.loadingitem = false;
          }
        });
        this.loadingitem = false;
      }
    },
    clickOutside(){
      if(this.showinfo = true ){
      this.showinfo = false};
    },
    
  },
};
</script>


