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
            <h4 class="text-capitalize pa-2 text-left">{{ thisitem.title }}</h4>

            <v-card class="ma-1" elevation="0">
              <v-chip
                class="ma-1 caption"
                label
                color="primary lighten-1"
                medium
              >
                <v-icon left> mdi-check-all </v-icon>
                ${{ thisitem.estimationprice }} TPP
              </v-chip>
              <v-chip
                class="ma-1 caption"
                label
                color="warning lighten-2"
                medium
              >
                <v-icon left> mdi-database-plus </v-icon>
                ${{ (thisitem.estimationprice*0.05).toFixed(0)}} TPP
              </v-chip>

              <p class="ma-1 caption" v-if="thisitem.description.length < 200">
                {{ thisitem.description }}
              </p>
              <p class="ma-1 caption" v-else>
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
            Full Details
          </v-btn>
        </div>
        <v-spacer></v-spacer>

        <v-btn icon @click="(showinfo = !showinfo), getItemPhotos()">
          <v-icon>{{
            showinfo ? "mdi-chevron-up" : "mdi-chevron-down"
          }}</v-icon>
        </v-btn>
      </v-card-actions>
      <v-divider />
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
                  <div class="body-1">" {{ thisitem.description }} "</div>
                </v-card-text>
              </v-card>

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
                v-if="thisitem.shippingcost"
                class="ma-1 caption"
                label
                outlined
                medium
              >
                <v-icon left> mdi-package-variant-closed </v-icon>
                Shipping: ${{ thisitem.shippingcost }} TPP
              </v-chip>

              <v-chip
                v-if="thisitem.bestestimator"
                class="ma-1 caption"
                label
                outlined
                medium
              >
                <v-icon left> mdi-check-all </v-icon>
                Price: ${{ thisitem.estimationprice }} TPP
              </v-chip>

              <v-chip class="ma-1 caption" medium label outlined>
                <v-icon left> mdi-account-outline </v-icon>
                Seller: {{ thisitem.creator }}
              </v-chip>
<!--
              <v-divider class="ma-2" />

              <div class="overline text-center">Comments</div>
              <div v-if="thisitem.comments">
                <v-chip
                  v-for="(listcomment, index) in commentlist"
                  v-bind:key="index"
                  class="ma-2"
                  >{{ listcomment }}
                </v-chip>
              </div>
              <div v-if="!thisitem.comments">
                <p class="caption text-center">No comments to show right now</p>
              </div>

              <v-divider class="ma-4" />
              <div v-if="hasAddress" class="ma-4 text-center">
                <wallet-coins />
              </div>
              <div class="text-center">
                <v-row>
                  <v-col>
                    <v-btn
                      block
                      color="primary"
                      :disabled="!thisitem.localpickup"
                      @click="submitLP(itemid), getThisItem"
                    >
                      Buy locally for ${{ thisitem.estimationprice }} TPP<v-icon
                        right
                      >
                        mdi-map-marker
                      </v-icon>
                      <div class="button__label" v-if="flight">
                        <div class="button__label__icon">
                          <icon-refresh />
                        </div>
                        Sending transaction...
                      </div>
                    </v-btn> </v-col
                  ><v-col>
                    <v-btn
                      block
                      color="primary"
                      :disabled="thisitem.shippingcost == 0"
                      @click="submitSP(itemid), getThisItem"
                    >
                      Buy for ${{ thisitem.estimationprice }} TPP + shipping
                      (${{ thisitem.shippingcost }} TPP)<v-icon right>
                        mdi-package-variant-closed
                      </v-icon>
                      <div class="button__label" v-if="flight">
                        <div class="button__label__icon">
                          <icon-refresh />
                        </div>
                        Sending transaction...
                      </div>
                    </v-btn>
                  </v-col>
                </v-row>
              </div>-->

         
              
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

    const imageRef = databaseRef.ref("ItemPhotoGallery/" + id);
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
  /*
    async submitLP(itemid) {
      if (!this.hasAddress) {
        alert("Sign in first");
      }

      if (!this.flightLP && this.hasAddress) {
        this.flightLP = true;
        this.loadingitem = true;
        let deposit = this.thisitem.estimationprice;
        //let deposit = toPay + "tpp";
        const type = { type: "buyer" };
        const body = { deposit, itemid };
        const fields = [
          ["buyer", 1, "string", "optional"],
          ["itemid", 2, "string", "optional"],
          ["deposit", 3, "int64", "optional"],
        ];
        await this.paySubmit({ body, fields });
        await this.$store.dispatch("entityFetch", type);
        await this.$store.dispatch("accountUpdate");
        this.flightLP = false;
        this.loadingitem = false;
      }
    },

    async submitSP(itemid) {
      if (!this.hasAddress) {
        alert("Sign in first");
      }
      if (!this.flightSP && this.hasAddress) {
        this.flightSP = true;
        this.loadingitem = true;

        let deposit =
          +this.thisitem.estimationprice + +this.thisitem.shippingcost;

        const fields = [
          ["buyer", 1, "string", "optional"],

          ["itemid", 2, "string", "optional"],
          ["deposit", 3, "int64", "optional"],
        ];
        const type = { type: "buyer" };
        const body = { deposit, itemid };
        await this.paySubmit({ body, fields });
        await this.$store.dispatch("entityFetch", type);
        await this.$store.dispatch("accountUpdate");

        this.flightSP = false;
        this.loadingitem = false;
        this.deposit = "";
        alert("Transaction sent");
      }
    },

    async paySubmit({ body, fields }) {
      const wallet = this.$store.state.wallet;
      const typeUrl = `/${process.env.VUE_APP_PATH}.MsgCreateBuyer`;
      let MsgCreate = new Type(`MsgCreateBuyer`);
      const registry = new Registry([[typeUrl, MsgCreate]]);
      console.log(fields);
      fields.forEach((f) => {
        MsgCreate = MsgCreate.add(new Field(f[0], f[1], f[2], f[3]));
      });

      const client = await SigningStargateClient.connectWithSigner(
        process.env.VUE_APP_RPC,
        wallet,
        { registry }
      );
      //console.log("TEST" + client)
      const msg = {
        typeUrl,
        value: {
          buyer: this.$store.state.account.address,
          ...body,
        },
      };

      console.log(msg);
      const fee = {
        amount: [{ amount: "0", denom: "tpp" }],
        gas: "200000",
      };

      const result = await client.signAndBroadcast(
        this.$store.state.account.address,
        [msg],
        fee
      );
      assertIsBroadcastTxSuccess(result);
      alert("Transaction sent");
    },

    

    async getThisItem() {
      await submit();
      return thisitem();
    },*/

    getItemPhotos() {
      if (this.showinfo && this.imageurl != "") {
        this.loadingitem = true;
        const id = this.itemid;

        const imageRef = databaseRef.ref("ItemPhotoGallery/" + id);
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