<template>
  <div class="pa-2 mx-lg-auto">
    <v-card class="pa-2 ma-auto" elevation="2" rounded="lg">
      <v-progress-linear
        indeterminate
        :active="loadingitem"
      ></v-progress-linear>

      <v-row>
        <v-col cols="12">
          <div class="subtitle-1  font-weight-medium text-capitalize text-center  mx-auto">{{ thisitem.title }} </div>
        </v-col>

        <v-col cols="12">
          <div v-if="imageurl">
            <v-img class="rounded contain" :src="imageurl"></v-img>
          </div>
        </v-col>
      </v-row>

      <div>
        <div class="pa-2 mx-auto" elevation="8">
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
 <v-chip outlined medium label class="ma-1 caption"
            v-for="tag in thisitem.tags" :key="tag"
          > <v-icon small left>
        mdi-tag-outline
      </v-icon>{{ tag }}</v-chip>
            <v-chip class="ma-1 caption" label outlined medium>
              <v-icon left> mdi-account-badge-outline </v-icon>
              Identifier: {{ thisitem.id }}
            </v-chip>
      

            <v-dialog transition="dialog-bottom-transition" max-width="300">
              <template v-slot:activator="{ on, attrs }">
                <span v-bind="attrs" v-on="on">
                  <v-chip class="ma-1 caption" label outlined medium>
                    <v-rating
      v-model="thisitem.condition.Number"
     readonly
      color="primary"
      background-color="grey lighten-1"
      small dense
    ></v-rating>
                  </v-chip>
                </span>
              </template>
              <template v-slot:default="dialog">
                <v-card>
                  <v-toolbar color="default"
                    >Condition (provided by seller)</v-toolbar
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
              Shipping Cost: ${{ thisitem.shippingcost }} TPP
            </v-chip>
            <v-chip
              outlined
              medium
              label
              class="ma-1 caption"
              v-for="country in thisitem.shippingregion"
              :key="country"
            >
              <v-icon small left> mdi-flag-variant-outline </v-icon
              >{{ country }}</v-chip
            >

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

            <v-chip  @click="createRoom" class="ma-1 caption" medium label outlined>
              <v-icon left> mdi-account-outline </v-icon>
              Seller: {{ thisitem.creator }}
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
            <v-divider class="ma-2" />

            <div class="overline text-center">Comments</div>
            <div v-if="thisitem.comments">
              <v-chip
                v-for="(single, i) in allcomments"
                v-bind:key="i"
                class="ma-2"
                >{{ single }}
              </v-chip>
            </div>
            <div v-if="allcomments.length ==0">
              <p class="caption text-center">No comments to show right now</p>
            </div>

            <v-divider class="ma-4" />
            <div v-if="hasAddress" class="ma-4 text-center">
              <wallet-coins />
            </div>
            <div class="text-center caption pa-2"> You can buy {{thisitem.title}}  for ${{ thisitem.estimationprice }} TPP and ship the item if you live in one of the following locations: "<span
            v-for="loc in thisitem.shippingregion" :key="loc"
          >{{ loc }}</span>" . Additional Shipping cost is ${{thisitem.shippingcost}} TPP. You can arrange a pickup by sending a message to <a @click="createRoom" >{{thisitem.creator}}. </a>   If you buy the item you will receive a cashback reward of ${{ (thisitem.estimationprice*0.05).toFixed(0)}} TPP. With TPP you can withdrawl your payment at any time, up until the item transaction and no transaction costs are applied.</div>
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
                    Buy for ${{ thisitem.estimationprice }} TPP + shipping (${{
                      thisitem.shippingcost
                    }}
                    TPP)<v-icon right> mdi-package-variant-closed </v-icon>
                    <div class="button__label" v-if="flight">
                      <div class="button__label__icon">
                        <icon-refresh />
                      </div>
                      Sending transaction...
                    </div>
                  </v-btn>
                </v-col>
              </v-row>
            </div>

            
          </div>
        </div>
      </div>
      <v-row class="pa-2 mx-auto">

      <v-btn 
        :disabled="!this.$store.state.account.address"
        text 
        @click="createRoom"
      >
          Message Seller</v-btn
      >
      <v-spacer/>
        <v-btn text @click="sellerInfo">Seller Info </v-btn>
         </v-row>
     <div class="pa-2 mx-auto caption">
       <v-card elevation="0" v-if="info">
          <p>This seller has sold {{ sold }} items before</p>
          <!--<p  Of which _ have been transfered by shipping and _ by local pickup.</p>-->
        </v-card>
        <v-card-title> All Seller items </v-card-title>
        <div v-for="item in SellerItems" v-bind:key="item.id">
          <v-card
            elevation="0"
            :to="{ name: 'BuyItemDetails', params: { id: item.id } }"
            ><v-row class="text-left caption ma-2">
              {{ item.title }} <v-spacer /> ${{ item.estimationprice }}TPP
              {{ item.status }}</v-row
            >
          </v-card>
        </div>
      </div>
    </v-card>
  </div>
</template>
<script>
import BuyItemDetails from "../views/BuyItemDetails.vue";
import { usersRef, roomsRef, databaseRef } from "./firebase/db.js";
import { SigningStargateClient, assertIsBroadcastTxSuccess } from "@cosmjs/stargate";
import {  Registry } from '@cosmjs/proto-signing/';
import { Type, Field } from 'protobufjs';

export default {
  components: { BuyItemDetails },
  props: ["itemid"],

  data() {
    return {
      amount: "",
      flight: false,
      flightLP: false,
      flightSP: false,
      info: false,
      imageurl: "",
      loadingitem: true,
      photos: [],
      dialog: false,
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
      return this.$store.getters.getItemByID(this.$route.params.id) || [];
    },

    hasAddress() {
      return !!this.$store.state.account.address;
    },
    valid() {
      return this.amount.trim().length > 0;
    },
    allcomments() {
      return this.thisitem.comments.filter((i) => i != "") || [];
    },
    SellerItems() {
      this.$store.dispatch("setSellerItemList", this.thisitem.creator);
      return this.$store.getters.getSellerList || [];
    },
  },

  methods: {
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
        await this.$store.dispatch("bankBalancesGet");
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
        await this.$store.dispatch("bankBalancesGet");

        this.flightSP = false;
        this.loadingitem = false;
        this.deposit = "";
      }
    },
    async getThisItem() {
      await submit();
      return thisitem();
    },

async paySubmit( { body, fields }) {
      const wallet = this.$store.state.wallet
      const typeUrl = `/${process.env.VUE_APP_PATH}.MsgCreateBuyer`;
      let MsgCreate = new Type(`MsgCreateBuyer`);
      const registry = new Registry([[typeUrl, MsgCreate]]);
      console.log(fields)
      fields.forEach(f => {
        MsgCreate = MsgCreate.add(new Field(f[0], f[1], f[2], f[3]))
      })

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
          ...body
        }
      };

      console.log(msg)
      const fee = {
        amount: [{ amount: '0', denom: 'tpp' }],
        gas: '200000'
      };

      const result = await client.signAndBroadcast(this.$store.state.account.address, [msg], fee);
      assertIsBroadcastTxSuccess(result);
      alert("Transaction sent");

    },
    
    getItemPhotos() {
      if (this.imageurl != "") {
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

    sellerInfo() {
      let rs = this.SellerItems.filter((i) => i.buyer != "");
      this.sold = "no";
      if (rs != "") {
        this.sold = rs.length;
      }

      this.info = !this.info

     
    },
    async createRoom() {

      if (!!this.$store.state.account.address) {
        let user = await usersRef
          .where("username", "==", this.$store.state.account.address)
          .get();
        if (user.docs[0] != null) {
          console.log("User Exists");
          //console.log(user.o_.docs[0].id)
          var userid = user.docs[0].id;
        } else {
          console.log("User does not exist");
          let { id } = await usersRef.add({
            username: this.$store.state.account.address,
          });
          console.log(id);
          await usersRef.doc(id).update({ _id: id });
          var userid = id;
        }

        let creator = await usersRef
          .where("username", "==", this.thisitem.creator)
          .get();
        if (creator.docs[0] != null) {
          console.log("User Exists");
          //console.log(user.o_.docs[0].id)
          var creatorid = creator.docs[0].id;
        } else {
          console.log("User does not exist");
          let { id } = await usersRef.add({ username: this.thisitem.creator });
          console.log(id);
          await usersRef.doc(id).update({ _id: id });
          var creatorid = id;
        }

        await roomsRef.add({
          users: [creatorid, userid],
          lastUpdated: new Date(),
        });
        console.log("asf");
       this.$router.push('/messages')
      }
    },
  },
};
</script>


