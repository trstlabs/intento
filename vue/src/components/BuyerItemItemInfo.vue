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
       

          <v-row align="start">
            <v-col>
              <v-card elevation="0">
                <div class="mt-4 overline text-center">Description</div>
                <v-card-text>
                  <div class="body-1">
                    {{ thisitem.description }}
                  </div>
                </v-card-text>
              </v-card>

              <v-divider class="ma-2" />
              <v-chip class="ma-1 caption" outlined>
                <v-icon small left> mdi-account-badge-outline </v-icon>
                TPP ID: {{ thisitem.id }}
              </v-chip>

              <v-chip
                v-if="thisitem.localpickup != ''"
                class="ma-1 caption"
                
                target="_blank"
                outlined
                
                :href="
                  'https://www.google.com/maps/search/?api=1&query=' +
                  thisitem.localpickup
                "
                ><v-icon small left> mdi-map-marker </v-icon>Pickup</v-chip
              >

              <v-chip
                v-if="thisitem.shippingcost > 0"
                class="ma-1 caption"
                
                outlined
                
              >
                <v-icon small left> mdi-package-variant-closed </v-icon>
                Shipping Cost: {{ thisitem.shippingcost }} tokens
              </v-chip>

              <v-chip
                v-if="thisitem.estimationprice > 0"
                class="ma-1 caption"
                
                outlined
                
              >
                <v-icon small left> mdi-check-all </v-icon>
                Price: {{ thisitem.estimationprice }} tokens
              </v-chip>

              <v-chip
                v-if="thisitem.status"
                class="ma-1 caption"
                
                outlined
                
              >
                 <v-icon small left> mdi-tune </v-icon>
                Status: {{ thisitem.status }}
              </v-chip>

              <v-chip
                v-else-if="(thisitem.transferable = true)"
                class="ma-1 caption"
                
                outlined
                
              >
                <v-icon small left> mdi-swap-horizontal </v-icon>
                Item Transferable
              </v-chip>

              <v-chip class="ma-1 caption"   outlined>
                <v-icon small left> mdi-account </v-icon>
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
          <v-btn rounded
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
                    <p >This item has shipped </p>
                    <p 
                      >Item has been shipped. Item seller indicated that item is
                      shipped. For more information contact the seller. The
                      protocol has received the request of the seller to arrange
                      tranfer coins.
                    </p>
                  </div>

                  <div v-if="thisitem.localpickup == '' && !thisitem.status">
                    <p >This item is not shipped yet</p>
                    <p 
                      >Contact the seller of {{ thisitem.title }}. Item seller
                      will indicate if the item is shipped.
                    </p>
                  </div>

                  <div>
                    <div
                      v-if="
                        thisitem.localpickup != '' &&
                        thisitem.status != 'Transferred'
                      "
                    >
                      <p class="ma-2" >
                        Arrange a meeting to pick up the item.
                      </p>
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
                          @click="submitWithdrawal(thisitem.id), getThisItem"
                          ><v-icon left> mdi-close </v-icon>
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
                    <v-btn rounded block outlined text @click="resell = !resell"> <span v-if="!resell"><v-icon left> mdi-repeat </v-icon> Resell item </span><span v-else> Cancel</span></v-btn>
                    <div class="pa-2 my-4" v-if="resell">
                          <p class="overline"><v-icon left> mdi-repeat </v-icon> Repost</p>
                      <v-textarea
                        class="ma-1"
                        prepend-icon="mdi-text"
                        :rules="rules.noteRules"
                        v-model="data.note"
                        label="Note (How is the item and why do you resell?)"
                        auto-grow
                      >
                      </v-textarea>

                      <v-row>
                        <v-btn
                          class="pa-2 mt-2"
                          text
                          icon
                          @click="data.shippingcost = 0"
                        >
                          <v-icon>
                            {{
                              data.shippingcost === 0
                                ? "mdi-package-variant"
                                : "mdi-package-variant-closed"
                            }}
                          </v-icon>
                        </v-btn>

                        <v-slider
                          class="pa-2 mt-2"
                          hint="Set to 0 tpp for no added cost"
                          thumb-label
                          label="Shipping cost"
                          suffix="tokens"
                          :persistent-hint="data.shippingcost != 0"
                          placeholder="Added cost"
                          :thumb-size="70"
                          v-model="data.shippingcost"
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
                          @click="data.discount = 0"
                        >
                          <v-icon>
                            {{
                              data.discount === 0
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
                          :persistent-hint="data.discount != 0"
                          placeholder="Discount"
                          :thumb-size="70"
                          v-model="data.discount"
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
                                data.shippingcost != 0 &&
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
                v-model="data.localpickup"
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

                      <v-btn rounded outlined block text @click="submitItemResell"> <v-icon left> mdi-repeat </v-icon>Resell</v-btn>
                    </div>
                  </div> <div v-else>The item status is {{ thisitem.status }}.</div>
                   
                    <div class="pt-4" v-if="makeReview">
                      <p class="overline"><v-icon left> mdi-star </v-icon> Rate</p>
                       <v-rating
                            v-model="data.rating"
                            
                            color="primary darken-1"
                            background-color="primary lighten-1"
                            
                            
                          ></v-rating>

                      <v-textarea
                        class="ma-1"
                        prepend-icon="mdi-text"
                        :rules="rules.noteRules"
                        v-model="data.reviewnote"
                        label="Note "
                        auto-grow
                      >
                      </v-textarea>
                      </div> <div class="pt-2"><v-btn rounded outlined @click="makeReview = !makeReview">  <span v-if="!makeReview"><v-icon left> mdi-star </v-icon> Rate item</span><span v-else> Cancel</span></v-btn>   <v-btn v-if="makeReview" outlined @click="submitItemRating()"> <v-icon left> mdi-star </v-icon>Post Rating</v-btn></div>
                </v-stepper-content>
              </v-stepper>
            </div>
          </div>
        </v-expand-transition></v-card
      ><sign-tx v-if="submitted" :key="submitted" :fields="fields" :value="value" :msg="msg" @clicked="afterSubmit"></sign-tx>
    </div>
  </div>
</template>

<script>
import ItemListBuyer from "./ItemListBuyer.vue";
import { databaseRef } from "./firebase/db.js";


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
      data: {
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

        fields: [],
      value: {},
      msg: "",
      submitted: false,
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
        this.fields = [
          ["buyer", 1, "string", "optional"],
          ["itemid", 2, "string", "optional"],
        ];
        this.msg = "MsgItemTransfer"
        this.value = {
          buyer: this.$store.state.account.address,
          ...body,
        },
    
     this.submitted = true
      }
    },

     async afterSubmit(value){
 this.loadingitem = true;

 this.msg = ""
 this.fields = []
 this.value = {}
  if(value == true){
         await this.$store.dispatch(
          "setBuyerItemList",
          this.$store.state.account.address
        );
        await this.$store.dispatch("bankBalancesGet");
        await this.$store.dispatch("updateItem", this.thisitem.id)//.then(result => this.newitem = result)
   setTimeout( () => this.$router.push("/itemid="+ this.thisitem.id), 5000)}

          this.submitted = false
              this.flightIT = false;
              this.flightITN = false;
    },


    async submitWithdrawal(itemid) {
      if (this.valid && !this.flightITN && this.hasAddress) {
        this.flightITN = true;
        const body = { itemid };
        this.fields = [
          ["buyer", 1, "string", "optional"],
          ["itemid", 2, "string", "optional"],
        ];
    this.msg ="MsgWithdrawal"
         this.value = {
          buyer: this.$store.state.account.address,
          ...body,
        },
   
      this.submitted = true
       
           }
      },



    async submitItemResell() {
      if (this.hasAddress) {
        const body = {
          itemid: this.itemid,
          shippingcost: this.data.shippingcost,
          discount: this.data.discount,
          localpickup: encodeURI(this.data.localpickup),
          shippingregion: this.selectedCountries,
          note: this.data.note,
        };
        this.msg = "MsgItemResell"

         this.value ={
          seller: this.$store.state.account.address,
          ...body,
        },
        this.fields = [
          ["seller", 1, "string", "optional"],
          ["itemid", 2, "string", "optional"],
          ["shippingcost", 3, "int64", "optional"],
          ["discount", 4, "int64", "optional"],
          ["localpickup", 5, "string", "optional"],
          ["shippingregion", 6, "string", "repeated"],
          ["note", 7, "string", "optional"],
        ];

             this.submitted = true
       
      }
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
     rating: this.data.rating,
          note: this.data.reviewnote,
        };

        this.value = {
          buyer: this.$store.state.account.address,
          ...body,
        },
  
        this.msg = "MsgItemRating"
        this.fields = [
          ["buyer", 1, "string", "optional"],
          ["itemid", 2, "string", "optional"],
          ["rating", 3, "int64", "optional"],
          ["note", 4, "string", "optional"],
        ];
      

             this.submitted = true
      }
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

