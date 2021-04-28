<template>
  <div class="pa-2 mx-lg-auto">
    <v-card color="secondary lighten-3"  class="pa-2 ma-auto" elevation="2" rounded="lg">
      <v-progress-linear
        indeterminate
        :active="loadingitem"
      ></v-progress-linear>

      <v-row>
        <v-col cols="12">
          <div
            class="subtitle-1 font-weight-medium text-capitalize text-center mx-auto"
          >
            {{ thisitem.title }}
          </div>
        </v-col>
      </v-row>

      <div>
        <div class="pa-2 mx-auto" elevation="8">
          <div>
            <div v-if="photos.photo">
            
              <v-carousel v-if="magnify == false" style="height:100%"
        delimiter-icon="mdi-minus"
        carousel-controls-bg="primary"
        height="300" 
        hide-delimiter-background
        show-arrows-on-hover 
      > 
       <v-carousel-item max-height="300" 
    contain v-for="(photo, i) in photos" :key="i" :src="photo" > 

  
        </v-carousel-item>
      </v-carousel>
      <v-carousel v-if="magnify == true"
        delimiter-icon="mdi-minus"
        carousel-controls-bg="primary"
     contain
        hide-delimiter-background
        show-arrows-on-hover 
      > 
       <v-carousel-item 
    contain v-for="(photo, i) in photos" :key="i" :src="photo" > 

   
        </v-carousel-item>
      </v-carousel>
            </div>
<v-row class="ml-4 mt-1 mb-1">
      <span  v-for="(photo, index) in photos" :key="index"> <img class="ma-1" @click="show(photo)" height="56"  :src="photo" /></span><v-spacer/><v-btn x-small class="mr-4"
            color="primary"
            icon
            @click="magnify = !magnify"
          >
          <v-icon> mdi-crop-free</v-icon>
          </v-btn></v-row> 
<v-dialog
      v-model="fullscreen"
    
    >
     

      <v-card>
        <v-card-title class=" grey lighten-2 ">
         {{thisitem.title}} <v-spacer></v-spacer>
          <v-btn
            color="primary"
            icon
            @click="fullscreen = false"
          >
          <v-icon> mdi-close</v-icon>
          </v-btn>
        </v-card-title>
<v-img :src="showphoto" />
       

      </v-card>
    </v-dialog> 
            <v-card color="secondary lighten-3" elevation="0">
              <div class="pa-2 overline text-center">Description</div>
              <v-card-text>
                <div class="body-1">{{ thisitem.description }}</div>
              </v-card-text>
       
            </v-card>
  
            <v-card v-if="thisitem.note" elevation="0">
                     <v-divider class="mx-4 pa-2" />
              <div class="pl-4 overline text-center">Reseller's Note</div>
              <v-card-text>
                <div class="body-1">
                  {{ thisitem.note }}
                </div>
              </v-card-text></v-card>
<v-divider class="mx-4 pa-2"/>
                <div class="text-center pa-2" v-if="thisitem.estimationprice != 0"><v-row  class="mx-4" >  <v-btn x-small icon @click="iteminfo = !iteminfo">
                <v-icon >mdi-information-outline</v-icon>
              </v-btn>
          <span  v-if="this.$store.state.account.address" >
            <wallet-coins /></span>
            
              <div v-if="iteminfo" class="text-center caption font-weight-light pa-2">
                You can buy {{ thisitem.title }}
                <span v-if="thisitem.shippingcost > 0"
                  >and ship the item if you live in one of the following
                  locations:
                  <span
                    v-for="loc in thisitem.shippingregion"
                    :key="loc"
                    class="font-weight-medium"
                  >
                    {{ loc }}
                  </span>
                  <span v-if="!thisitem.shippingregion[0]"> all locations </span
                  >. Additional cost (e.g. shipping) is ${{
                    thisitem.shippingcost
                  }}
                  TPP.</span
                >
                <span v-if="thisitem.localpickup">
                  and you can arrange a pickup by sending a message to
                  <a @click="createRoom">{{ thisitem.seller }}. </a>
                </span>
                <span v-if="thisitem.discount > 0">
                  Reseller gives a discount of ${{ thisitem.discount}} on
                  the original selling price of ${{
                    thisitem.estimationprice
                  }}TPP.</span
                >
                <span v-if="thisitem.creator == thisitem.seller"
                  >If you buy the item you will receive a cashback reward of ${{
                    (thisitem.estimationprice * 0.05).toFixed(0)
                  }}
                  TPP.
                </span>
                With TPP you can withdrawl your payment at any time, up until
                the item transaction and no transaction costs are applied.
              </div>
            </v-row>

                   <v-img @click="iteminfo = !iteminfo" src="img/design/buy.png" >
           </v-img> 
              <v-row v-if="thisitem.creator == thisitem.seller">
                <v-col  >
                  <v-btn
                    block 
                    color="primary lighten-1"
                    :disabled="!thisitem.localpickup"
                    @click="
                      submit(thisitem.estimationprice),
                        (flightLP = !flightLP),
                        getThisItem
                    "
                    ><div v-if="!flightLP">
                    <v-icon left >
                        mdi-check-all
                      </v-icon>  Buy for {{ thisitem.estimationprice }}<v-icon small right>$vuetify.icons.custom</v-icon> 
                    </div>
                    <div v-if="flightLP">
                      <v-progress-linear
                        indeterminate
                        color="secondary"
                      ></v-progress-linear
                      >Sending transaction...
                    </div>
                  </v-btn> </v-col
                ><v-col>
                  <v-btn
                    block
                    color="primary lighten-1"
                    :disabled="thisitem.shippingcost == 0"
                    @click="
                      submit(
                        Number(thisitem.estimationprice) +
                          Number(thisitem.shippingcost)
                      ),
                        (flightSP = !flightSP),
                        getThisItem
                    "
                    ><div v-if="!flightSP"><v-icon left > mdi-check-all</v-icon
                      ><v-icon left> mdi-plus</v-icon
                      ><v-icon left > mdi-package-variant-closed </v-icon>
                   
                      Buy for {{
                        Number(thisitem.estimationprice) +
                        Number(thisitem.shippingcost)
                      }} <v-icon small >$vuetify.icons.custom</v-icon> </div>
                    <div v-if="flightSP">
                      <v-progress-linear
                        indeterminate
                        color="secondary"
                      ></v-progress-linear
                      >Sending transaction...
                    </div>
                  </v-btn>
                </v-col>
              </v-row>
              <v-row v-else>
                <v-col>
                  <v-btn
                    block
                    color="primary"
                    :disabled="!thisitem.localpickup || thisitem.discount == 0"
                    @click="
                      submit(
                        Number(thisitem.estimationprice) -
                          Number(thisitem.discount)
                      ),
                        (flightLP = !flightLP),
                        getThisItem
                    "
                    ><div v-if="!flightLP">
                      Buy for ${{
                        Number(thisitem.estimationprice) -
                        Number(thisitem.discount)
                      }}TPP <v-icon right> mdi-repeat </v-icon
                      ><v-icon right> mdi-minus </v-icon
                      ><v-icon right> mdi-label-percent</v-icon>
                    </div>
                    <div v-if="flightLP">
                      <v-progress-linear
                        indeterminate
                        color="secondary"
                      ></v-progress-linear
                      >Sending transaction...
                    </div>
                  </v-btn>
                  <v-btn
                    class="mt-2"
                    block
                    color="primary"
                    :disabled="!thisitem.localpickup || thisitem.discount > 0"
                    @click="
                      submit(thisitem.estimationprice),
                        (flightLP = !flightLP),
                        getThisItem
                    "
                    ><div v-if="!flightLP">
                      Buy for ${{ thisitem.estimationprice }}TPP
                      <v-icon right> mdi-repeat </v-icon>
                    </div>
                    <div v-if="flightLP">
                      <v-progress-linear
                        indeterminate
                        color="secondary"
                      ></v-progress-linear
                      >Sending transaction...
                    </div>
                  </v-btn> </v-col
                ><v-col>
                  <v-btn
                    block
                    color="primary"
                    :disabled="
                      thisitem.shippingcost == 0 || thisitem.discount > 0
                    "
                    @click="
                      submit(
                        Number(thisitem.estimationprice) +
                          Number(thisitem.shippingcost)
                      ),
                        (flightSP = !flightSP),
                        getThisItem
                    "
                    ><div v-if="!flightSP">
                      Buy for ${{
                        Number(thisitem.estimationprice) +
                        Number(thisitem.shippingcost)
                      }}TPP <v-icon right> mdi-repeat </v-icon>
                      <v-icon right> mdi-plus </v-icon
                      ><v-icon right> mdi-package-variant-closed </v-icon>
                    </div>
                    <div v-if="flightSP">
                      <v-progress-linear
                        indeterminate
                        color="secondary"
                      ></v-progress-linear
                      >Sending transaction...
                    </div>
                  </v-btn>
                  <v-btn
                    class="mt-2"
                    block
                    color="primary"
                    :disabled="
                      thisitem.shippingcost == 0 || thisitem.discount == 0
                    "
                    @click="
                      submit(
                        Number(thisitem.estimationprice) +
                          Number(thisitem.shippingcost) -
                          Number(thisitem.discount)
                      ),
                        (flightSP = !flightSP),
                        getThisItem
                    "
                    ><div v-if="!flightSP">
                      Buy for ${{
                        Number(thisitem.estimationprice) +
                        Number(thisitem.shippingcost) -
                        Number(thisitem.discount)
                      }}TPP <v-icon right> mdi-repeat </v-icon>
                      <v-icon right> mdi-plus </v-icon
                      ><v-icon right> mdi-package-variant-closed </v-icon
                      ><v-icon right> mdi-minus </v-icon
                      ><v-icon right> mdi-label-percent </v-icon>
                    </div>
                    <div v-if="flightSP">
                      <v-progress-linear
                        indeterminate
                        color="secondary"
                      ></v-progress-linear
                      >Sending transaction...
                    </div>
                  </v-btn>
                </v-col>
                
              </v-row>  <v-divider class="ma-4 pt-2"/>
            </div>
             

            <v-chip class="ma-1 caption" label outlined medium>
              <v-icon left> mdi-account-badge-outline </v-icon>
              Identifier: {{ thisitem.id }}
            </v-chip>
 
       

            <v-dialog transition="dialog-bottom-transition" max-width="300">
              <template v-slot:activator="{ on, attrs }">
                <span v-bind="attrs" v-on="on">
                  <v-chip class="ma-1 caption " label outlined medium>
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
            <span v-if="thisitem.estimationprice > 0"><span v-if="thisitem.creator != thisitem.seller">
              <v-chip
                v-if="
                  thisitem.shippingcost > 0 &&
                  thisitem.localpickup == false &&
                  thisitem.discount == 0
                "
                class="ma-1 caption"
                label
                outlined
              >
                <v-icon left> mdi-repeat </v-icon>
                <v-icon left> mdi-plus </v-icon
                ><v-icon left> mdi-package-variant-closed </v-icon> {{
                  Number(thisitem.estimationprice) +
                  Number(thisitem.shippingcost)
                }}<v-icon small right>$vuetify.icons.custom</v-icon>  
              </v-chip>
              <v-chip
                v-if="
                  thisitem.shippingcost > 0 &&
                  thisitem.localpickup == false &&
                  thisitem.discount > 0
                "
                class="ma-1 caption"
                label
                outlined
              >
                <v-icon left> mdi-repeat </v-icon
                ><v-icon small left> mdi-plus </v-icon
                ><v-icon left> mdi-package-variant-closed </v-icon>
                <v-icon small left> mdi-minus </v-icon
                ><v-icon left> mdi-label-percent</v-icon> {{
                  Number(thisitem.estimationprice) +
                  Number(thisitem.shippingcost) -
                  Number(thisitem.discount)
                }}<v-icon small right>$vuetify.icons.custom</v-icon>  
              </v-chip>
              <v-chip
                v-if="thisitem.discount > 0 && thisitem.localpickup"
                class="ma-1 caption"
                label
                outlined
              >
                <v-icon small> mdi-repeat </v-icon
                ><v-icon small> mdi-minus </v-icon
                ><v-icon small> mdi-label-percent</v-icon> {{
                  thisitem.estimationprice - thisitem.discount
                }}<v-icon small right>$vuetify.icons.custom</v-icon>  
              </v-chip>
              <v-chip class="ma-1 caption" label outlined>
                <v-icon left> mdi-repeat </v-icon>Original price: {{
                  thisitem.estimationprice
                }}<v-icon small right>$vuetify.icons.custom</v-icon>  
              </v-chip>
            </span>
            <span v-else>
              <span v-if="thisitem.localpickup == false">
                <v-chip
                  class="ma-1 caption"
                  label
                 
                  outlined
                >
                  <v-icon left> mdi-check-all </v-icon
                  ><v-icon small left> mdi-plus </v-icon
                  ><v-icon small left> mdi-package-variant-closed </v-icon> {{
                    Number(thisitem.estimationprice) +
                    Number(thisitem.shippingcost)
                  }}<v-icon small right>$vuetify.icons.custom</v-icon>  
                </v-chip></span
              >

              <span v-else>
                <v-chip class="ma-1 caption" label outlined>
                  <v-icon left> mdi-check-all </v-icon>
                  {{ thisitem.estimationprice }}<v-icon small right>$vuetify.icons.custom</v-icon>  
                </v-chip></span
              >
              <v-chip class="ma-1 caption" label outlined>
                <v-icon small left> mdi-hand-heart </v-icon>Cashback: {{
                  (thisitem.estimationprice * 0.05).toFixed(0)
                }}
                <v-icon small right>$vuetify.icons.custom</v-icon>  
              </v-chip>
            </span>
            <v-chip
              v-if="thisitem.localpickup"
              class="ma-1 caption"
              label
              outlined
              ><v-icon left> mdi-map-marker-outline </v-icon>
              Pickup</v-chip
            >

            <v-chip
              v-if="thisitem.shippingcost > 0"
              class="ma-1 caption"
              label
              outlined
            >
              <v-icon left> mdi-package-variant-closed </v-icon>
              Shipping Cost: {{ thisitem.shippingcost}} <v-icon small right>$vuetify.icons.custom</v-icon>  
            </v-chip></span>

            <v-chip
              v-if="thisitem.discount > 0"
              class="ma-1 caption"
              label
              outlined
            >
              <v-icon left> mdi-label-percent </v-icon>
              Discount: {{ thisitem.discount}} <v-icon small right>$vuetify.icons.custom</v-icon>  
            </v-chip>
            <v-chip
              outlined
              label
              class="ma-1 caption"
              v-for="country in thisitem.shippingregion"
              :key="country"
            >
              <v-icon small left> mdi-flag-variant-outline </v-icon
              >{{ country }}</v-chip
            >

            <v-chip @click="createRoom" class="ma-1 caption" label outlined>
              <v-icon left> mdi-account </v-icon>
              Seller: {{ thisitem.seller }}
            </v-chip>
            <v-chip class="ma-1 caption" label outlined v-if="thisitem.creator != thisitem.seller ">
              <v-icon left> mdi-account-outline </v-icon>
              Creator: {{ thisitem.creator }}
            </v-chip>
            <v-chip
              outlined
              label
              class="ma-1 caption"
              v-for="tag in thisitem.tags"
              :key="tag"
            >
              <v-icon small left> mdi-tag-outline </v-icon>{{ tag }}</v-chip
            >
      
 <v-divider class="ma-4" /> 
            <div class="overline text-center">Comments</div>
            <div v-if="thisitem.comments "><div class="font-weight-light" v-for="(single, i) in allcomments"
                v-bind:key="i"><v-icon  small left> mdi-message-text-outline</v-icon> 
            <v-chip
                
                class="ma-2 " color="primary lighten-2"
                > {{ single }}
              </v-chip></div>
            </div>
            <div v-if="allcomments.length == 0">
              <p class="caption text-center">No comments to show right now</p>
            </div>

            <v-divider class="ma-4" />
                
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
        <v-spacer />
        <v-btn text @click="sellerInfo">Seller Info </v-btn>
      </v-row>
      <div class="pa-2 mx-auto caption">
        <span v-if="info">
          <p class="text-center">This seller has sold {{ sold }} items before</p>
          <!--<p  Of which _ have been transfered by shipping and _ by local pickup.</p>-->
        </span>
        <v-card-title v-if="SellerItems[0]" class="overline justify-center">
          All Seller items
        </v-card-title>
        <div v-for="item in SellerItems" v-bind:key="item.id">
          <v-card color="secondary lighten-3" 
            elevation="0"
            :to="{ name: 'BuyItemDetails', params: { id: item.id } }"
          >
            <v-row class="text-left caption ma-2"
              ><span class="font-weight-medium"> {{ item.title }}</span>
              <v-spacer /><v-spacer /> <span> {{ item.status }}</span>
              <span v-if="item.transferable && item.status != ''">
                ${{ item.estimationprice }}TPP </span
              ><span v-if="item.buyer && !item.transferable"> Sold </span>
              <span v-if="!item.estimationprice"> Awaiting estimation </span>
              <span v-if="!item.transferable && item.estimationprice"
                >Not on sale yet</span
              >
              <span v-if="item.thank">Buyer thanked seller</span>
            </v-row>
          </v-card>
        </div>
      </div><v-img src="img/design/transfer.png" ></v-img>
    </v-card>
  </div>
</template>
<script>
import BuyItemDetails from "../views/BuyItemDetails.vue";
import { usersRef, roomsRef, databaseRef } from "./firebase/db.js";
import {
  SigningStargateClient,
  assertIsBroadcastTxSuccess,
} from "@cosmjs/stargate";
import { Registry } from "@cosmjs/proto-signing/";
import { Type, Field } from "protobufjs";

export default {
  components: { BuyItemDetails },
  props: ["itemid"],

  data() {
    return {
      amount: "",
      iteminfo: false,
      flight: false,
      flightLP: false,
      flightSP: false,
      info: false,
      imageurl: "",
      loadingitem: true,
      photos: [],
      dialog: false,
      fullscreen: false,
      magnify: false,
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
        this.photos = data;
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
      this.$store.dispatch("setBuySellerItemList", this.thisitem.seller);
      return this.$store.getters.getBuySellerList || [];
    },
  },

  methods: {
    async submit(deposit) {
      if (!this.hasAddress) {
        alert("Sign in first");
        this.$router.push('/')
        window.location.reload();
      }

      if (this.hasAddress) {
        // this.flightLP = true;
        this.loadingitem = true;
        const type = { type: "buyer" };
        const body = { deposit: deposit, itemid: this.thisitem.id };
        const fields = [
          ["buyer", 1, "string", "optional"],

          ["itemid", 2, "string", "optional"],
          ["deposit", 3, "int64", "optional"],
        ];
        await this.paySubmit({ body, fields });
        await this.$store.dispatch("entityFetch", type);
        await this.$store.dispatch("bankBalancesGet");

        this.loadingitem = false;
      }
      this.flightLP = false;
      this.flightSP = false;
    },

    async getThisItem() {
      await submit();
      return thisitem();
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
      if (!result.data) {
        alert("TX failed");
        window.location.reload();
      }
      assertIsBroadcastTxSuccess(result);
      alert("Transaction sent");
    },

    show(photo){
      this.showphoto = photo 

      this.fullscreen = true
    },

    sellerInfo() {
      let rs = this.SellerItems.filter((i) => i.buyer != "");
      this.sold = "no";
      if (rs != "") {
        this.sold = rs.length;
      }

      this.info = !this.info;
    },
    async createRoom() {
      if (this.$store.state.user.uid) {
        const user = await usersRef
          .where("username", "==", this.thisitem.seller)
          .get();
        console.log(user);

        //let query = await roomsRef.where("users", '', [this.$store.state.user.uid).where("users", "array-contains", user.docs[0].id).get()
        /*await roomsRef.where("users", "==", ["5RlZazMyPgdoHgGfjTud", "B1Xk6qliE2ceNJN6HsoCk2MQO2K2"]).get()
 .then((querySnapshot) => {
    querySnapshot.forEach((doc) => {
      console.log(doc.data())
        console.log(doc.id, ' => ', doc.data());
    });
});*/
        if (user.docs[0]) {
          let query = await roomsRef
            .where("users", "==", [user.docs[0].id, this.$store.state.user.uid])
            .get();
          let otherquery = await roomsRef
            .where("users", "==", [this.$store.state.user.uid, user.docs[0].id])
            .get();
          console.log(query.docs[0]);
          /*

await roomsRef.where("users", "array-contains", this.$store.state.user.uid).get()
   .then((querySnapshot) => {
    querySnapshot.forEach((doc) => {
      console.log(doc.data(users))
        console.log(doc.id, ' => ', doc.data());
    });
});*/

          if (query.docs[0] || otherquery.docs[0]) {
            alert("Seller already added or seller not found");
          } else {
            //await usersRef.doc(id).update({ _id: id });
            await roomsRef.add({
              users: [user.docs[0].id, this.$store.state.user.uid],
              lastUpdated: new Date(),
            });
          }
          this.$router.push("/messages");
        } else {
          alert("Seller DatabaseID not found");
        }
      } else {
        alert("Sign in first (Check your Google email)");
      }
    },
  },
};
</script>


