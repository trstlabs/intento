<template>
  <div class="pa-2 mx-auto">
    <v-card elevation="2" rounded="lg" v-click-outside="clickOutside">
      <v-progress-linear
        indeterminate
        :active="loadingitem"
      ></v-progress-linear>

      <div class="pa-2 mx-auto">
        <v-row>
          <v-col cols="12" md="8">
            <p
              v-if="thisitem.creator != thisitem.seller"
              class="text-capitalize subtitle-2  pa-2 text-left"
            >
              <v-icon small left> mdi-repeat </v-icon>{{ thisitem.title }}
            </p>

            <p v-else class="text-capitalize subtitle-2 pa-2 text-left">
              <v-icon small left>mdi-check-all </v-icon>
              {{ thisitem.title }}
            </p>

            <v-card class="ma-1" elevation="0">
              <p class="ma-1 caption font-weight-light" v-if="thisitem.description.length < 200">
                {{ thisitem.description }}
              </p>
              <p class="ma-1 caption font-weight-light" v-else>
                {{ thisitem.description.substring(0, 148) + ".." }}
              </p>
            </v-card>
          </v-col>

          <v-col cols="12" md="4">
            <div v-if="imageurl">
              <v-img class="rounded contain" :src="imageurl"></v-img>
            </div>
          </v-col>
        </v-row>
      </div>
      <v-card-actions>
        <v-btn 
          color="blue"
          text
          @click="(showinfo = !showinfo), getItemPhotos()"
        >
          Info
        </v-btn>
        <div>
          <v-btn 
            color="blue"
            :to="{ name: 'BuyItemDetails', params: { id: itemid } }"
            text
          >
            Details
          </v-btn>
        </div>
        <div v-if="thisitem.creator != thisitem.seller">
        <!-- <v-chip
            v-if="
              thisitem.shippingcost > 0 &&
              thisitem.localpickup == false &&
              thisitem.discount == 0
            "
            class="ma-1 pl-0 caption"
            label
            color="primary lighten-2"
            small
          >
            <v-chip dark color="primary"  
              ><v-icon small right>$vuetify.icons.custom</v-icon>{{
                Number(thisitem.estimationprice) + Number(thisitem.shippingcost)
              }}</v-chip
            ><v-icon small left> mdi-repeat </v-icon>
            <v-icon small left> mdi-plus </v-icon
            ><v-icon small left> mdi-package-variant-closed </v-icon>
          </v-chip>-->
          
           <span>
   <v-chip  v-if="
               thisitem.shippingcost > 0 &&
              thisitem.localpickup == false &&
              thisitem.discount == 0
            " small dark color="primary" class="pr-0" :to="{ name: 'BuyItemDetails', params: { id: itemid } }"
              >
              <v-hover v-slot="{ hover }" close-delay="300" open-delay="60" >
              <span>
              <span  class="pr-2" v-if="hover" > Buy Now </span><span class="pr-2" v-else><v-icon small right>$vuetify.icons.custom</v-icon>{{
                 Number(thisitem.estimationprice) + Number(thisitem.shippingcost)
              }} </span>
              </span>
</v-hover> 


                  <v-chip label
           
            class="pl-0 caption"
            color="primary lighten-2"
            
          ><v-icon small left> mdi-repeat </v-icon>
            <v-icon small left> mdi-plus </v-icon
            ><v-icon small left> mdi-package-variant-closed </v-icon>
          </v-chip>
          </v-chip
            >
          </span>

<!--

          <v-hover v-slot="{ hover }" close-delay="200"  > <span>
          <v-chip  :to="{ name: 'BuyItemDetails', params: { id: itemid } }"
            v-if="
              thisitem.shippingcost > 0 &&
              thisitem.localpickup == false &&
              thisitem.discount > 0
            "
            class="ma-1 pl-0 caption"
            color="primary lighten-2"
            small
          >
            <v-chip label dark color="primary"
              ><span v-if="hover" > Buy Now </span><span v-else><v-icon small right>$vuetify.icons.custom</v-icon>{{
                Number(thisitem.estimationprice) +
                Number(thisitem.shippingcost) -
                Number(thisitem.discount)
              }}</span></v-chip
            >
            <v-icon small  right> mdi-repeat </v-icon
            ><v-icon small  right> mdi-plus </v-icon
            ><v-icon small  right> mdi-package-variant-closed </v-icon>
            <v-icon small right> mdi-minus </v-icon
            ><v-icon small right> mdi-label-percent</v-icon>
          </v-chip></span></v-hover> 
          -->
 <span>
   <v-chip  v-if="
              thisitem.shippingcost > 0 &&
              thisitem.localpickup == false &&
              thisitem.discount > 0
            " small dark color="primary" class="pr-0" :to="{ name: 'BuyItemDetails', params: { id: itemid } }"
              >
              <v-hover v-slot="{ hover }" close-delay="300" open-delay="60" >
              <span>
              <span  class="pr-2" v-if="hover" > Buy Now </span><span class="pr-2" v-else><v-icon small right>$vuetify.icons.custom</v-icon>{{
                Number(thisitem.estimationprice) +
                Number(thisitem.shippingcost) -
                Number(thisitem.discount)
              }} </span>
              </span>
</v-hover> 


                  <v-chip label
           
            class="pl-0 caption"
            color="primary lighten-2"
            
          ><v-icon small  right> mdi-repeat </v-icon
            ><v-icon small  right> mdi-plus </v-icon
            ><v-icon small  right> mdi-package-variant-closed </v-icon>
            <v-icon small right> mdi-minus </v-icon
            ><v-icon small right> mdi-label-percent</v-icon>
          </v-chip>
          </v-chip
            >
          </span>


        <!--  <v-chip
            v-if="thisitem.discount > 0 && thisitem.localpickup"
            class="ma-1 pl-0 caption"
            
            color="primary lighten-2"
            small
          >
            <v-chip  label dark color="primary">
              <v-icon small right>$vuetify.icons.custom</v-icon>{{ thisitem.estimationprice - thisitem.discount }}</v-chip
            >
            <v-icon small right> mdi-repeat </v-icon>
            <v-icon small right> mdi-minus </v-icon
            ><v-icon small right> mdi-label-percent</v-icon>
          </v-chip>-->
          <span>
   <v-chip  v-if="thisitem.discount > 0 && thisitem.localpickup"
             small dark color="primary" class="pr-0" :to="{ name: 'BuyItemDetails', params: { id: itemid } }"
              >
              <v-hover v-slot="{ hover }" close-delay="300" open-delay="60" >
              <span>
              <span  class="pr-2" v-if="hover" > Buy Now </span><span class="pr-2" v-else><v-icon small right>$vuetify.icons.custom</v-icon>{{
                thisitem.estimationprice - thisitem.discount
              }} </span>
              </span>
</v-hover> 


                  <v-chip label
           
            class="pl-0 caption"
            color="primary lighten-2"
            
          ><v-icon small right> mdi-repeat </v-icon>
            <v-icon small right> mdi-minus </v-icon
            ><v-icon small right> mdi-label-percent</v-icon>
          </v-chip>
          </v-chip
            >
          </span>

        </div>
        <div v-else>

          
          <span v-if="thisitem.localpickup == false">
            <!--<v-chip class="ma-1 caption"  color="primary lighten-1" small>
              <v-chip label dark color="primary">
                <v-icon small right>$vuetify.icons.custom</v-icon>{{
                  Number(thisitem.estimationprice) +
                  Number(thisitem.shippingcost)
                }}</v-chip
              >
              <v-icon right> mdi-check-all </v-icon
              ><v-icon small right> mdi-plus </v-icon
              ><v-icon small right> mdi-package-variant-closed </v-icon>
            </v-chip>-->
          
   <v-chip  
             small dark color="primary" class="pr-0" :to="{ name: 'BuyItemDetails', params: { id: itemid } }"
              >
              <v-hover v-slot="{ hover }" close-delay="300" open-delay="60" >
              <span>
              <span  class="pr-2" v-if="hover" > Buy Now </span><span class="pr-2" v-else><v-icon small right>$vuetify.icons.custom</v-icon>{{
                Number(thisitem.estimationprice) +
                  Number(thisitem.shippingcost)
              }} </span>
              </span>
</v-hover> 


                  <v-chip label
           
            class="pl-0 caption"
            color="primary lighten-1"
            
          > <v-icon right> mdi-check-all </v-icon
              ><v-icon small right> mdi-plus </v-icon
              ><v-icon small right> mdi-package-variant-closed </v-icon>
          </v-chip>
          </v-chip
            >
     
            
            </span
          >

          <span v-else>
            <!--
            <v-chip
              class="ma-1 pl-0 caption"
              
              color="primary lighten-1"
              small
            >
              <v-chip label dark color="primary">
                <v-icon small right>$vuetify.icons.custom</v-icon>{{ thisitem.estimationprice }}</v-chip
              >
              <v-icon right> mdi-check-all </v-icon>
            </v-chip>-->
            <v-chip  
             small dark color="primary lighten-1" class="pr-0" :to="{ name: 'BuyItemDetails', params: { id: itemid } }"
              >
              <v-hover v-slot="{ hover }" close-delay="300" open-delay="60" >
              <span>
              <span  class="pr-2" v-if="hover" > Buy Now </span><span class="pr-3 caption" v-else>{{
                thisitem.estimationprice
              }}<v-icon small right>$vuetify.icons.custom</v-icon></span>
              </span>
</v-hover> 


                  <v-chip label
           
            class="pl-0 caption"
            color="primary"
            
          > <v-icon right> mdi-check-all </v-icon>
          
          </v-chip>
          </v-chip
            >
            
            </span
          >
          <v-chip
            class="ma-1 pl-0 caption"
            
            dark
            color="green lighten-2"
            small
          ><v-chip dark label color="green">
            {{ (thisitem.estimationprice * 0.05).toFixed(0)
            }}<v-icon small right>$vuetify.icons.custom</v-icon></v-chip>
            <v-icon small right> mdi-hand-heart </v-icon>
          </v-chip>
        </div>

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

              <v-card elevation="0">
                <div class="pa-2 overline text-center">Description</div>
                <v-card-text>
                  <div class="caption " >{{ thisitem.description }}</div>
                </v-card-text>
              </v-card>
              <v-chip
                outlined
                medium
                label
                class="ma-1 caption font-weight-light"
                v-for="itemtags in thisitem.tags"
                :key="itemtags"
              >
                <v-icon small left> mdi-tag-outline </v-icon
                >{{ itemtags }}</v-chip
              >
              <v-chip
                outlined
                medium
                label
                class="ma-1 caption"
                v-for="selected in thisitem.shippingregion"
                :key="selected"
              >
                <v-icon small left> mdi-flag-variant-outline </v-icon
                >{{ selected }}</v-chip
              >
              <v-chip class="ma-1 caption" label outlined medium>
                <v-icon left> mdi-account-badge-outline </v-icon>
                Identifier: {{ thisitem.id }}
              </v-chip>

              <v-chip class="ma-1 caption" label outlined medium>
                <v-icon left> mdi-star-outline </v-icon>
                Condition: {{ thisitem.condition }}/5
              </v-chip>

              <v-chip
                v-if="thisitem.localpickup"
                class="ma-1 caption"
                label
                outlined
                medium
                ><v-icon left> mdi-map-marker-outline </v-icon>Local
                Pickup</v-chip
              >

              <v-chip
                v-if="thisitem.shippingcost > 0"
                class="ma-1 caption"
                label
                outlined
                medium
              >
                <v-icon left> mdi-package-variant-closed </v-icon>
                Shipping: {{ thisitem.shippingcost}} <v-icon right small>$vuetify.icons.custom</v-icon> 
              </v-chip>

              <v-chip
                v-if="thisitem.bestestimator"
                class="ma-1 caption"
                label
                outlined
                medium
              >
                <v-icon left> mdi-check-all </v-icon>
                Price: {{ thisitem.estimationprice}} <v-icon right small>$vuetify.icons.custom</v-icon> 
              </v-chip>

              <v-chip class="ma-1 caption" medium label outlined>
                <v-icon left> mdi-account </v-icon>
                Seller: {{ thisitem.seller }}
              </v-chip>
             
            </div>
          </div>
        </div>
      </v-expand-transition>
    </v-card>
  </div>
</template>

<script>
import { databaseRef } from "./firebase/db";
import ItemListBuy from "./ItemListBuy.vue";
import WalletCoins from "./WalletCoins.vue";
import {
  SigningStargateClient,
  assertIsBroadcastTxSuccess,
} from "@cosmjs/stargate";
import { Registry } from "@cosmjs/proto-signing/";
import { Type, Field } from "protobufjs";

export default {
  props: ["itemid"],
  components: { ItemListBuy, WalletCoins },
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

    const imageRef = databaseRef.ref("ItemPhotoGallery/" + id + "/photos/");
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
      return this.$store.getters.getItemByID(this.itemid);
    },

    hasAddress() {
      return !!this.$store.state.account.address;
    },
    valid() {
      return this.amount.trim().length > 0;
    },
    commentlist() {
      return this.thisitem.comments.filter((i) => i != "") || [];
    },
  },

  methods: {
    getItemPhotos() {
      if (this.showinfo && this.imageurl != "") {
        this.loadingitem = true;
        const id = this.itemid;

        const imageRef = databaseRef.ref("ItemPhotoGallery/" + id + "/photos/");
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
    clickOutside() {
      if ((this.showinfo = true)) {
        this.showinfo = false;
      }
    },
  },
};
</script>



<!---
shows item id from buy list
<div id="item-list-buy">
      {{ itemid }}
    </div>
    ---->