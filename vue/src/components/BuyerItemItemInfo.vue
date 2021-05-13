<template>
  <div>
    <div class="pa-2 mx-auto">
      <v-card elevation="2" rounded="lg">
        <v-progress-linear
          indeterminate
          :active="loadingitem"
        ></v-progress-linear>
        <div class="pa-2 mx-auto">
          <p class="mx-6 overline">{{ thisitem.title }}</p>

          <v-carousel
            height="400"
            hide-delimiter-background
            show-arrows-on-hover
          >
            <v-carousel-item v-for="(photo, i) in photos" :key="i" :src="photo">
            </v-carousel-item>
          </v-carousel>
          <v-divider class="ma-2" />

          <v-row align="start">
            <v-col>
              <v-card elevation="0">
                <div class="pl-4 overline text-center">Description</div>
                <v-card-text>
                  <div class="body-1">
                    {{ thisitem.description }}
                  </div>
                </v-card-text>
              </v-card>

              <v-divider class="ma-2" />
              <v-chip class="ma-1 caption" label outlined medium>
                <v-icon left> mdi-account-badge </v-icon>
                TPP ID: {{ thisitem.id }}
              </v-chip>

              <v-chip
                v-if="thisitem.localpickup != ''"
                class="ma-1 caption"
                label
                target="_blank"
                outlined
                medium
                :href="
                  'https://www.google.com/maps/search/?api=1&query=' +
                  thisitem.localpickup
                "
                ><v-icon left> mdi-map-marker </v-icon>Pickup</v-chip
              >

              <v-chip
                v-if="thisitem.shippingcost > 0"
                class="ma-1 caption"
                label
                outlined
                medium
              >
                <v-icon left> mdi-package-variant-closed </v-icon>
                Shipping Cost: {{ thisitem.shippingcost }} tokens
              </v-chip>

              <v-chip
                v-if="thisitem.estimationprice > 0"
                class="ma-1 caption"
                label
                outlined
                medium
              >
                <v-icon left> mdi-check-all </v-icon>
                Price: {{ thisitem.estimationprice }} tokens
              </v-chip>

              <v-chip
                v-if="thisitem.status"
                class="ma-1 caption"
                label
                outlined
                medium
              >
                <v-icon left> mdi-clock-time-three-outline </v-icon>
                Status: {{ thisitem.status }}
              </v-chip>

              <v-chip
                v-else-if="(thisitem.transferable = true)"
                class="ma-1 caption"
                label
                outlined
                medium
              >
                <v-icon left> mdi-swap-horizontal </v-icon>
                Item Transferable
              </v-chip>

              <v-chip class="ma-1 caption" medium label outlined>
                <v-icon left> mdi-account-outline </v-icon>
                Seller: {{ thisitem.seller }}
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
                <v-stepper-step step="1" complete> Prepayment </v-stepper-step>

                <v-stepper-step :complete="thisitem.status != ''" step="2">
                  Item Transfer
                </v-stepper-step>

                <v-stepper-content step="2">
                  <div v-if="thisitem.tracking === true">
                    <app-text type="p">This item has shipped </app-text>
                    <app-text type="p"
                      >Item has been shipped. Item seller indicated that item is
                      shipped. For more information contact the seller. The
                      protocol has received the request of the seller to arrange
                      tranfer coins.
                    </app-text>
                  </div>

                  <div v-if="thisitem.localpickup == '' && !thisitem.status">
                    <app-text type="p">This item is not shipped yet</app-text>
                    <app-text type="p"
                      >Contact the seller of {{ thisitem.title }}. Item seller
                      will indicate if the item is shipped.
                    </app-text>
                  </div>

                  <div>
                    <div
                      v-if="
                        thisitem.localpickup != '' &&
                        thisitem.status != 'Transferred'
                      "
                    >
                      <app-text class="ma-2" type="p">
                        Arrange a meeting to pick up the item.
                      </app-text>
                      <v-row>
                        <v-btn
                          class="ma-4"
                          color="primary"
                          @click="submitItemTransfer(thisitem.id), getThisItem"
                          ><v-icon left> mdi-checkbox-marked-circle </v-icon>
                          <span v-if="!flightIT"> Complete transfer</span>
                          <div class="button__label" v-else>
                            <div class="button__label__icon"></div>
                            Sending tokens to seller...
                          </div>
                        </v-btn>

                        <v-btn
                          class="ma-4"
                          color="default"
                          :class="[
                            'button',
                            `button__valid__${
                              !!valid && !flightITN && hasAddress
                            }`,
                          ]"
                          @click="submitDeleteBuyer(thisitem.id), getThisItem"
                          ><v-icon left> mdi-cancel </v-icon>
                          <span v-if="!flightITN"> Cancel transfer</span>
                          <div class="button__label" v-else>
                            <div class="button__label__icon"></div>
                            Sending tokens back...
                          </div>
                        </v-btn>
                      </v-row>
                    </div>
                  </div>
                </v-stepper-content>
                <v-stepper-step :complete="thisitem.status != ''" step="3">
                  Complete
                </v-stepper-step>
                <v-stepper-content step="3" class="mx-6 pa-0">
                 
                  <div
                    v-if="
                      thisitem.status === 'Transferred' ||
                      thisitem.status === 'Shipped'
                    "
                  >
                    <p>The transfer is complete. Enjoy your {{thisitem.tags[0]}}!</p>
                    <v-btn block outlined text @click="resell = !resell"> <span v-if="!resell"><v-icon left> mdi-repeat </v-icon> Resell item </span><span v-else> Cancel</span></v-btn>
                    <div class="pa-2 my-4" v-if="resell">
                          <p class="overline"><v-icon left> mdi-repeat </v-icon> Repost</p>
                      <v-textarea
                        class="ma-1"
                        prepend-icon="mdi-text"
                        :rules="rules.noteRules"
                        v-model="fields.note"
                        label="Note (How is the item and why do you resell?)"
                        auto-grow
                      >
                      </v-textarea>

                      <v-row>
                        <v-btn
                          class="pa-2 mt-2"
                          text
                          icon
                          @click="fields.shippingcost = 0"
                        >
                          <v-icon>
                            {{
                              fields.shippingcost === 0
                                ? "mdi-package-variant"
                                : "mdi-package-variant-closed"
                            }}
                          </v-icon>
                        </v-btn>

                        <v-slider
                          class="pa-2 mt-2"
                          hint="Set to 0 tpp no for added cost"
                          thumb-label
                          label="Shipping cost"
                          suffix="tokens"
                          :persistent-hint="fields.shippingcost != 0"
                          placeholder="Added cost"
                          :thumb-size="70"
                          v-model="fields.shippingcost"
                          ><template v-slot:thumb-label="item">
                            {{ item.value
                            }}<v-icon small right>$vuetify.icons.custom</v-icon>
                          </template>
                        </v-slider>
                      </v-row>
                      <v-row>
                        <v-btn
                          class="pa-2 mt-2"
                          text
                          icon
                          @click="fields.discount = 0"
                        >
                          <v-icon>
                            {{
                              fields.discount === 0
                                ? "mdi-brightness-percent-outline"
                                : "mdi-brightness-percent"
                            }}
                          </v-icon>
                        </v-btn>

                        <v-slider
                          class="pa-2 mt-2"
                          hint="Explain discount in the note"
                          thumb-label
                          label="Discount"
                          suffix="tokens"
                          :persistent-hint="fields.discount != 0"
                          placeholder="Discount"
                          :thumb-size="70"
                          v-model="fields.discount"
                          ><template v-slot:thumb-label="item">
                            {{ item.value
                            }}<v-icon small right>$vuetify.icons.custom</v-icon>
                          </template>
                        </v-slider>
                      </v-row>
                      <v-row>
                        <v-col>
                          <v-row>
                            <v-btn
                              class="pa-2"
                              text
                              icon
                              @click="enterlocation = !enterlocation"
                            >
                              <v-icon>
                                {{
                                  enterlocation
                                    ? "mdi-map-marker"
                                    : "mdi-map-marker-off"
                                }}
                              </v-icon>
                            </v-btn>

                            <v-switch
                              class="ml-2"
                              v-model="enterlocation"
                              inset
                              label="Pickup"
                              :persistent-hint="
                                fields.shippingcost != 0 &&
                                enterlocation == tue &&
                                selectedCountries.length > 1
                              "
                              hint="Specify local pickup location in description"
                            ></v-switch> </v-row
                          > <v-text-field
                class="ma-1"
                prepend-icon="mdi-map-marker"
                :rules="rules.pickupRules"
                label="Pickup Location (optional)"
                v-model="fields.localpickup"
                required v-if="enterlocation"
              /></v-col
                        ><v-col>
                          <v-select
                            prepend-icon="mdi-earth"
                            hint="At least one"
                            :persistent-hint="selectedCountries == 0"
                            v-model="selectedCountries"
                            :items="countryCodes"
                            :rules="rules.shippingRules"
                            label="Location"
                            deletable-chips
                            multiple
                            chips
                          >
                          </v-select> </v-col
                      ></v-row>

                      <v-btn  outlined block text @click="submitItemResell"> <v-icon left> mdi-repeat </v-icon>Resell</v-btn>
                    </div>
                  </div> <div v-else>The item status is {{ thisitem.status }}.</div>
                   
                    <div class="pt-4" v-if="makeReview">
                      <p class="overline"><v-icon left> mdi-star </v-icon> Rate</p>
                       <v-rating
                            v-model="fields.rating"
                            
                            color="primary darken-1"
                            background-color="primary lighten-1"
                            
                            
                          ></v-rating>

                      <v-textarea
                        class="ma-1"
                        prepend-icon="mdi-text"
                        :rules="rules.noteRules"
                        v-model="fields.reviewnote"
                        label="Note "
                        auto-grow
                      >
                      </v-textarea>
                      </div> <div class="pt-2"><v-btn  outlined @click="makeReview = !makeReview">  <span v-if="!makeReview"><v-icon left> mdi-star </v-icon> Rate item</span><span v-else> Cancel</span></v-btn>   <v-btn v-if="makeReview" outlined @click="submitItemRating()"> <v-icon left> mdi-star </v-icon>Post Rating</v-btn></div>
                </v-stepper-content>
              </v-stepper>
            </div>
          </div>
        </v-expand-transition></v-card
      >
    </div>
  </div>
</template>

<script>
import ItemListBuyer from "./ItemListBuyer.vue";
import { databaseRef } from "./firebase/db.js";
import {
  SigningStargateClient,
  assertIsBroadcastTxSuccess,
} from "@cosmjs/stargate";
import { Registry } from "@cosmjs/proto-signing/";
import { Type, Field } from "protobufjs";

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
      resell: false,
      photos: [],
      imageurl: "",
      step: 2,
      fields: {
        shippingcost: "0",
        localpickup: "",
        discount: "0",
        note: "",
        reviewnote: "",
        rating: "0",
      },
      makeReview: false,

      rules: {
        pickupRules: [
      
        
          (v) =>
            (v.length <= 25) || "Pickup must be less than 25 characters, enter coordinates instead",
        ],
        shippingRules: [(v) => !!v.length == 1 || "A country is required"],
        noteRules: [
          (v) =>
            (v && v.length <= 80) || "Note must be less than 80 characters",
        ],
      },

      selectedCountries: [],
      countryCodes: ["NL", "BE", "UK", "DE", "US", "CA"],
    };
  },
  beforeCreate() {
    this.loadingitem = true;
  },

  mounted() {
    const id = this.itemid;

    const imageRef = databaseRef.ref("ItemPhotoGallery/" + id + "/photos/");
    imageRef.on("value", (snapshot) => {
      const data = snapshot.val();

      if (data != null ) {
        // console.log(data[0]);
        this.photos = data;
        this.imageurl = data[0];
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
    async submitItemTransfer(itemid) {
      if (this.valid && !this.flightIT && this.hasAddress) {
        this.flightIT = true;

        const body = { itemid };
        const fields = [
          ["buyer", 1, "string", "optional"],
          ["itemid", 2, "string", "optional"],
        ];
        await this.transferSubmit({ body, fields });
        await this.$store.dispatch(
          "setBuyerItemList",
          this.$store.state.account.address
        );
        this.flightIT = false;
      }
    },

    async submitDeleteBuyer(itemid) {
      if (this.valid && !this.flightITN && this.hasAddress) {
        this.flightITN = true;
        const body = { itemid };
        const fields = [
          ["buyer", 1, "string", "optional"],
          ["itemid", 2, "string", "optional"],
        ];
        await this.deleteSubmit({ body, fields });
        await this.$store.dispatch("entityFetch", "buyer");
        await this.$store.dispatch(
          "setBuyerItemList",
          this.$store.state.account.address
        );
        this.flightITN = false;
      }
    },

    async transferSubmit({ body, fields }) {
      const wallet = this.$store.state.wallet;
      const typeUrl = `/${process.env.VUE_APP_PATH}.MsgItemTransfer`;
      let MsgCreate = new Type(`MsgItemTransfer`);
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

      const msg = {
        typeUrl,
        value: {
          buyer: this.$store.state.account.address,
          ...body,
        },
      };
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

    async deleteSubmit({ body, fields }) {
      const wallet = this.$store.state.wallet;
      const typeUrl = `/${process.env.VUE_APP_PATH}.MsgDeleteBuyer`;
      let MsgCreate = new Type(`MsgDeleteBuyer`);
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

      const msg = {
        typeUrl,
        value: {
          buyer: this.$store.state.account.address,
          ...body,
        },
      };
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

    async submitItemResell() {
      if (this.hasAddress) {
        const body = {
          itemid: this.itemid,
          shippingcost: this.fields.shippingcost,
          discount: this.fields.discount,
          localpickup: encodeURI(this.fields.localpickup),
          shippingregion: this.selectedCountries,
          note: this.fields.note,
        };
        const fields = [
          ["seller", 1, "string", "optional"],
          ["itemid", 2, "string", "optional"],
          ["shippingcost", 3, "int64", "optional"],
          ["discount", 4, "int64", "optional"],
          ["localpickup", 5, "string", "optional"],
          ["shippingregion", 6, "string", "repeated"],
          ["note", 7, "string", "optional"],
        ];
        await this.resellSubmit({ body, fields });
        await this.$store.dispatch(
          "setBuyerItemList",
          this.$store.state.account.address
        );
      }
    },

    async resellSubmit({ body, fields }) {
      const wallet = this.$store.state.wallet;
      const typeUrl = `/${process.env.VUE_APP_PATH}.MsgItemResell`;
      let MsgCreate = new Type(`MsgItemResell`);
      const registry = new Registry([[typeUrl, MsgCreate]]);

      fields.forEach((f) => {
        MsgCreate = MsgCreate.add(new Field(f[0], f[1], f[2], f[3]));
      });

      const client = await SigningStargateClient.connectWithSigner(
        process.env.VUE_APP_RPC,
        wallet,
        { registry }
      );

      const msg = {
        typeUrl,
        value: {
          seller: this.$store.state.account.address,
          ...body,
        },
      };
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
      alert("Resell request sent");
    },
    async getThisItem() {
      await submitrevealestimation();
      return this.thisitem();
    },
    createStep() {
      if (this.thisitem.status == "") {
        this.step = 2;
      } else if (this.thisitem.status != "") {
        this.step = 3;
      }
    },

   async submitItemRating() {
       
      if (this.hasAddress) {
     
        const body = {
          itemid: this.itemid,
     rating: this.fields.rating,
          note: this.fields.reviewnote,
        };
        const fields = [
          ["buyer", 1, "string", "optional"],
          ["itemid", 2, "string", "optional"],
          ["rating", 3, "int64", "optional"],
          ["note", 4, "string", "optional"],
        ];
        await this.rateSubmit({ body, fields });
        await this.$store.dispatch(
          "setBuyerItemList",
          this.$store.state.account.address
        );
      }
    },

    async rateSubmit({ body, fields }) {
      const wallet = this.$store.state.wallet;
      const typeUrl = `/${process.env.VUE_APP_PATH}.MsgItemRating`;
      let MsgCreate = new Type(`MsgItemRating`);
      const registry = new Registry([[typeUrl, MsgCreate]]);

      fields.forEach((f) => {
        MsgCreate = MsgCreate.add(new Field(f[0], f[1], f[2], f[3]));
      });

      const client = await SigningStargateClient.connectWithSigner(
        process.env.VUE_APP_RPC,
        wallet,
        { registry }
      );

      const msg = {
        typeUrl,
        value: {
          buyer: this.$store.state.account.address,
          ...body,
        },
      };
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
      alert("Review sent");
    },
    async getThisItem() {
      await submitrevealestimation();
      return this.thisitem();
    },
    createStep() {
      if (this.thisitem.status == "") {
        this.step = 2;
      } else if (this.thisitem.status != "") {
        this.step = 3;
      }
  },
  },
};
</script>

