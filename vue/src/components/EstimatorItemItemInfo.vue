<template>
  <div>
    <div class="pa-2 mx-auto">
      <v-card elevation="2" class="pa-2" rounded="lg">
        <v-progress-linear
          indeterminate
          :active="loadingitem"
        ></v-progress-linear>
        <div class="pa-2 mx-auto" elevation="8">
          <v-row>
           <p class="pa-2 overline"> {{ thisitem.title }} </p>
            <v-spacer /><v-btn
              v-if="
                thisitem.highestestimator != userAddress &&
                thisitem.lowestestimator != userAddress
              "
              text
              @click="removeItem()"
              ><v-icon> mdi-trash-can </v-icon></v-btn
            >
          </v-row>

          <v-divider class="pa-2" />

          <v-row align="start">
            <v-col cols="8">
              <v-chip class="ma-1 caption"  outlined >
                <v-icon left small> mdi-account-badge-outline </v-icon>
                TPP ID: {{ thisitem.id }}
              </v-chip>

              <v-chip
                v-if="thisitem.bestestimator"
                class="ma-1 caption" 
                
                outlined
                
              >
                <v-icon small left> mdi-check-all </v-icon>
                Final Estimation: {{ thisitem.estimationprice }}
                <v-icon small right>$vuetify.icons.custom</v-icon>
              </v-chip>

              <span>
                <div class="pa-2 overline">Description</div>
                <div class="px-2 caption">
                  {{ thisitem.description }}
                </div>
              </span>


              <p
                class="mt-1"
                v-if="thisitem.bestestimator === userAddress"
                
              >
                <v-divider class="ma-4" />
                <v-icon left> mdi-account-check </v-icon>You are the best
                estimator.
                <span class="caption">
                  If the item is transferred, you will be rewarded
                  {{ (thisitem.estimationprice * 0.05).toFixed(0)
                  }}<v-icon small right>$vuetify.icons.custom</v-icon> .</span
                >
              </p>
              <p
                class="mt-1"
                v-else-if="thisitem.lowestestimator === userAddress"
                
              >
                <v-divider class="ma-4" />
                <v-icon left> mdi-account-arrow-left </v-icon>

                You are the lowest estimator.
                <span v-if="!thisitem.estimationprice" class="caption"
                  >If the item owner does not accept the estimation price, you
                  lose {{ thisitem.depositamount
                  }}<v-icon small right>$vuetify.icons.custom</v-icon> .</span
                >
              </p>
              <p
                class="mt-1"
                v-else-if="thisitem.highestestimator === userAddress"
                
              >
                <v-divider class="ma-4" />
                <v-icon left> mdi-account-arrow-right </v-icon>
                You are the highest estimator.<span
                  class="caption"
                  v-if="thisitem.transferable == false"
                >
                  If the seller does not accept the estimation price
                </span>
                <span class="caption" v-else
                  >If the seller does not ship the item, you lose
                  {{ thisitem.depositamount
                  }}<v-icon small right>$vuetify.icons.custom</v-icon> .</span
                >
              </p>
              <p class="mt-1 text-center" v-else type="caption">
                <v-divider class="ma-4" />
                You have estimated this item and you are neither the highest,
                lowest or the best. You may withdrawl your TPP tokens. This is
                done automatically for you when the item transfers.
              </p>
            </v-col>

            <v-col cols="4">
              <div v-if="imageurl" class="d-flex flex-row-reverse text-center">
                <v-avatar class="ma-2 rounded" size="125" tile>
                  <v-img :src="imageurl"></v-img>
                </v-avatar>
              </div>
            </v-col>
          </v-row>
        </div> </v-card
      ><sign-tx
        v-if="submitted"
        :key="submitted"
        :fields="fields"
        :value="value"
        :msg="msg"
        @clicked="afterSubmit"
      ></sign-tx>
    </div>
  </div>
</template>

<script>
import { databaseRef } from "./firebase/db";
import ItemListEstimator from "./ItemListEstimator.vue";


export default {
  props: ["itemid"],
  components: { ItemListEstimator },
  data() {
    return {
      loadingitem: true,

      photos: [],
      imageurl: "",

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
        console.log(data[0]);
        this.photos = data;
        this.imageurl = data[0];
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

       async afterSubmit(value){
 this.loadingitem = true;

 this.msg = ""
 this.fields = []
 this.value = {}
  if(value == true) {
             await this.$store.dispatch("updateItem", this.thisitem.id)//.then(result => this.newitem = result)
       
        await this.$store.dispatch("bankBalancesGet")
alert("Estimation deleted")
  }
          this.submitted = false
       this.flightre = false;
        this.loadingitem = false;  
   
     

    },


    async removeItem() {
      this.loadingitem = true;
      this.flightre = true;
     
      const body = { itemid: this.itemid };
      this.fields = [
        ["estimator", 1, "string", "optional"],
        ["itemid", 2, "string", "optional"],
      ];

      this.msg = "MsgDeleteEstimation"

      this.value = {
          estimator: this.$store.state.account.address,
          ...body,
        },
        


      this.submitted = true
    },

  }
    
};
</script>
