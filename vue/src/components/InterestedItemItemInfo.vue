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
      Identifier: {{ thisitem.id }}
    </v-chip>

<v-chip v-if="thisitem.localpickup"
      class="ma-1"
      label
      outlined
      medium
    >
    <v-icon left>
        mdi-pin
      </v-icon>
      Local pickup available
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
      Shipping cost: $ {{thisitem.shippingcost}} tokens
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
             
              <v-divider></v-divider>
             

              <v-row> <v-col>
            <v-btn block color="primary"
              v-if="thisitem.localpickup"
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
            <v-btn block color="primary"
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
    </v-card>
  </div>
</template>

<script>
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
   

    async submitLP(itemid) {
      if (!this.flightLP && this.hasAddress) {
        this.flightLP = true;
        this.loadingitem = true;
        let toPay = this.thisitem.estimationprice;
        let deposit = toPay + "token";
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
        let deposit = toPaySP + "token";
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


     async submitInterest(itemid) {
      if (!this.flightLP && this.hasAddress) {
    //    this.flightLP = true;
        this.loadingitem = true;
        
        const type = { type: "estimator/set" };
        const body = { itemid: itemid,
        interested: false };
        await this.$store.dispatch("entitySubmit", { ...type, body });
       // await this.$store.dispatch("entityFetch", type);
       // await this.$store.dispatch("accountUpdate");
     //   this.flightLP = false;
        this.loadingitem = false;
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
  },
};
</script>

<style scoped>

.card__empty {
  margin-bottom: 1rem;
  border: 1px dashed rgba(0, 0, 0, 0.1);
  padding: 1rem;
  display: flex;
  align-items: center;
  justify-content: center;
  box-sizing: border-box;
  border-radius: 8px;

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


<!---
shows item id from buy list
<div id="item-list-buy">
      {{ itemid }}
    </div>
    ---->