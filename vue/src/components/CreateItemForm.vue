<template>
  <div class="pa-2 mx-lg-auto"> 
    <v-dialog v-model="dialog" width="500">
      <template v-slot:activator="{ on, attrs }">
      <span v-if="!fields.title" >
           <h2  v-if="!fields.title" v-bind="attrs"
          v-on="on" class="headline pt-2 font-weight-bold text-center"> Place Item
        </h2>  </span><span  v-else>
          <h2
          v-bind="attrs"
          v-on="on"
          class="headline pt-2 font-weight-bold text-center"
        >
          Place '{{ fields.title }}'
        </h2></span>
        <v-img v-bind="attrs"
          v-on="on" height="300" src="img/design/sell.png"> <v-icon class="ml-4" small>mdi-information-outline</v-icon></v-img>
      </template>

      <v-card class="text-center">
        <v-card-title class="h2 lighten-2"> Info </v-card-title>

        <v-card-text>
          After placing the item, an estimation will be made. After you accept this price, anyone can buy the item!
        </v-card-text>
 <iframe style="
      

-webkit-mask-image: -webkit-radial-gradient(circle, white 100%, black 100%); /*ios 7 border-radius-bug */
-webkit-transform: rotate(0.000001deg); /*mac os 10.6 safari 5 border-radius-bug */
-webkit-border-radius: 10px; 
-moz-border-radius: 10px;
border-radius: 20px; 
overflow: hidden; 
" width="100%"  height="310" src="https://www.youtube.com/embed/zHXwfePrGvA?vq=hd1080&autoplay=0&loop=1&modestbranding=1&rel=0&cc_load_policy=1&color=white&mute=1" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>
        <v-divider></v-divider>

        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="primary" text @click="dialog = false"> Let's go </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    
    <v-stepper class="elevation-0" v-model="e1">
      <v-stepper-header>
        <v-stepper-step :complete="e1 > 1" step="1"> Data </v-stepper-step>

        <v-divider></v-divider>

        <v-stepper-step :complete="e1 > 2" step="2"> Pictures </v-stepper-step>
        <v-divider></v-divider>

        <v-stepper-step :complete="e1 > 3" step="3"> Done! </v-stepper-step>
      </v-stepper-header>

      <v-stepper-items>
    
          <v-stepper-content class="ma-5" step="1">
      
              <v-alert type="warning"
                class="caption text-center"
                v-if="
                  !this.$store.state.user && this.$store.state.account.address
                "
              >
                To post items, Confirm this sign in by clicking the link sent to your
                Google account's email on this device. 
              </v-alert>
              <v-text-field
                class="ma-1"
                prepend-icon="mdi-format-title"
                :rules="rules.titleRules"
                label="Title"
                v-model="fields.title"
                required
              />

              <v-textarea
                class="ma-1"
                prepend-icon="mdi-text"
                :rules="rules.descriptionRules"
                v-model="fields.description"
                label="Description"
                auto-grow
              >
              </v-textarea>

              <v-combobox
                prepend-icon="mdi-tag-outline"
                hint="At least one and at most 5 category tags"
                :persistent-hint="selectedTags == 0 || selectedTags == 5"
                :search-input.sync="search"
                v-model="selectedTags"
                :items="taglist"
                :rules="rules.tagRules"
                label="Categories"
                deletable-chips
                multiple
                chips
              >
                <template v-slot:no-data>
                  <v-list-item>
                    <v-list-item-content>
                      <v-list-item-title>
                        No category tags matching "<strong>{{ search }}</strong
                        >". Press <kbd>enter</kbd> to create a new one
                      </v-list-item-title>
                    </v-list-item-content>
                  </v-list-item>
                </template>
                <template v-slot:selection="{ attrs, item, parent, selected }">
                  <v-chip
                    v-if="selectedTags[0] == item"
                    v-bind="attrs"
                    :input-value="selected"
                    color="primary lighten-2"
                    small
                  >
                    <span class="pr-2">
                      Main: <v-icon x-small>mdi-tag</v-icon> {{ item }}
                    </span>
                    <v-icon small @click="parent.selectItem(item)">
                      mdi-close
                    </v-icon>
                  </v-chip>
                  <v-chip
                    v-else
                    v-bind="attrs"
                    color="secondary darken-1"
                    :input-value="selected"
                    small
                  >
                    <span class="pr-2">
                      <v-icon x-small>mdi-tag</v-icon> {{ item }}
                    </span>
                    <v-icon small @click="parent.selectItem(item)">
                      mdi-close
                    </v-icon>
                  </v-chip>
                </template>
              </v-combobox>

              <div>
                <v-row class="pa-2 mt-4"
                  ><v-btn text icon @click="fields.estimationcount = 3">
                    <v-icon> mdi-check</v-icon></v-btn
                  >
                  <v-slider
                    hint="Lower for faster results, higher for better accuracy"
                    thumb-label
                    :persistent-hint="fields.estimationcount != 3"
                    label="Accuracy"
                    :thumb-size="90"
                    max="12"
                    :rules="rules.estimationcountRules"
                    placeholder="Estimation count"
                    v-model="fields.estimationcount"
                    ><template v-slot:thumb-label="item">
                      {{ item.value }} Estimations
                    </template></v-slider
                  >
                </v-row>
                <v-row class="pa-2">
                  <v-btn text icon @click="fields.condition = 0">
                    <v-icon>{{
                      fields.condition === 0 ? "mdi-star-outline" : "mdi-star"
                    }}</v-icon>
                  </v-btn>
                  <v-slider
                    label="Condition"
                    :hint="
                      'Condition is ' +
                      conditionLabel() +
                      ', explain condition in description'
                    "
                    v-model="fields.condition"
                    :max="4"
                    :persistent-hint="fields.condition != 0"
                    :thumb-size="24"
                    thumb-label
                    ><template v-slot:thumb-label="{ value }">
                      {{ satisfactionEmojis[value] }}
                    </template>
                  </v-slider></v-row
                >

                <v-row class="pa-2 mt-2">
                  <v-btn text icon @click="fields.shippingcost = 0">
                    <v-icon>
                      {{
                        fields.shippingcost === 0
                          ? "mdi-package-variant"
                          : "mdi-package-variant-closed"
                      }}
                    </v-icon>
                  </v-btn>

                  <v-slider
                    hint="Set to 0 tokens no for shipping"
                    thumb-label
                    label="Shipping cost"
                    suffix="tokens"
                    :persistent-hint="fields.shippingcost != 0"
                    placeholder="Shipping cost"
                    :thumb-size="60"
                    thumb-color="primary lighten-1"
                    v-model="fields.shippingcost"
                    ><template v-slot:thumb-label="item">
                      {{ item.value }} <v-icon>$vuetify.icons.custom</v-icon>
                    </template>
                  </v-slider>
                </v-row>
                <v-row>
                  <v-col>
                    <v-row class="ma-0">
                      <v-btn
                        class="mr-auto"
                        text
                        icon
                        @click="fields.localpickup = !fields.localpickup"
                        ><v-icon
                          >{{
                            fields.localpickup
                              ? "mdi-map-marker"
                              : "mdi-map-marker-off"
                          }}
                        </v-icon></v-btn
                      >

                      <v-switch
                        class="mr-auto mt-1"
                        v-model="fields.localpickup"
                        inset
                        label="Local pickup"
                        :persistent-hint="
                          fields.shippingcost != 0 &&
                          fields.localpickup == true &&
                          selectedCountries.length > 1
                        "
                        hint="Specify local pickup location in description"
                      ></v-switch> </v-row></v-col
                  ><v-col>
                    <v-select
                      class="mt-1 pt-0"
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
              </div>
              <div class="text-center pt-6">
                <v-btn
                  color="primary"
                  :disabled="!valid || !!flight || !hasAddress"
                  @click="submit()"
                  ><div v-if="!flight">
                    Next<v-icon> mdi-arrow-right-bold</v-icon>
                  </div>
                  <div v-if="flight">
                    <v-progress-linear
                      indeterminate
                      color="white"
                      class="ma-1"
                    ></v-progress-linear
                    >Creating item ID...
                  </div>
                </v-btn>
              </div>
  
          </v-stepper-content>
       

        <v-stepper-content step="2">
          <div v-if="showpreview">
            <create-item-preview-and-upload
              :thisitem="thisitem"
              v-on:changeStep="updateStepCount($event)"
            />
          </div>
        </v-stepper-content>
        <v-stepper-content step="3">
          <v-alert rounded-lg type="success">
            Submitted, the item will be estimated
          </v-alert>
          <p>
            You can always find your item in your
            <router-link to="/account">account</router-link>. Your item will be
            available to buy after you make it transferable.
          </p>
        </v-stepper-content>
      </v-stepper-items>
    </v-stepper>
  </div>
</template>


<script>
import { assertIsBroadcastTxSuccess } from "@cosmjs/launchpad";
import CreateItemPreviewAndUpload from "./CreateItemPreviewAndUpload.vue";
import { SigningStargateClient } from "@cosmjs/stargate";
import { Registry } from "@cosmjs/proto-signing/";
import { Type, Field } from "protobufjs";

export default {
  components: { CreateItemPreviewAndUpload },
  data() {
    return {
      fields: {
        title: "",
        description: "",
        shippingcost: "0",
        localpickup: false,
        estimationcount: "3",

        condition: "0",
      },

      rules: {
        titleRules: [
          (v) => !!v || "Title is required",
          (v) => (v && v.length > 4) || "Title must be more than 3 characters",
          (v) =>
            (v && v.length <= 80) || "Title must be less than 80 characters",
        ],
        descriptionRules: [
          (v) => !!v || "Description is required",
          (v) =>
            (v && v.length > 4) || "Description must be more than 4 characters",
          (v) =>
            (v && v.length <= 800) ||
            "Description must be less than 800 characters",
        ],
        estimationcountRules: [
          (v) => !!v || "Estimation count is required",
          (v) =>
            (v && v > 2) || "Estimation count must be more than 2 estimators",
          (v) =>
            (v && v < 12) || "Estimation count must be less than 12 estimators",
        ],
        tagRules: [
          (v) => !!v.length == 1 || "Category tag is required",
          (v) => (v && v.length < 6) || "Category tags must be less than 6",
        ],

        shippingRules: [(v) => !!v.length == 1 || "A country is required"],
      },
      itemid: "",
      selectedTags: [],
      selectedCountries: [],
      thisitem: {},
      flight: false,
      showpreview: false,
      e1: 1,
      search: null,
      dialog: false,

      satisfactionEmojis: ["😭", "🙁", "🙂", "😊", "😄"],
      countryCodes: ["NL", "BE", "UK", "DE", "US", "CA"],
    };
  },
  watch: {
    selectedTags(val) {
      if (val.length > 5) {
        this.$nextTick(() => this.selectedTags.pop());
      }
    },
  },

  mounted() {
    //console.log(input)

    this.$store.dispatch("setSortedTagList");
  },

  computed: {
    taglist() {
    
      if (this.selectedTags == 0) {   let list =  [
          "Books",
          "Clothing",
          "Shoes",
          "   Collectible  ",
          "   Consumer Electronic  ",
          "   Home & Garden  ",
          "   Motor",
          "Bike",
          "   Pet supplies ",
          "   Sport",
          "   Toys & Hobbies  ",
          "   Antique ",
          "  Computer",
          "Smartdevice",
          "Smartphone",
          "Sound Device",
          "TV",
          "NFT Art",
          "NFT Collectible",
        ];
        return list
      } else {
        //this.$store.dispatch("setTagList");

        return this.$store.getters.getTagList
      }
    },

    hasAddress() {
      return !!this.$store.state.account.address || alert("Sign in first");
    },

    valid() {
      if (
        this.fields.title.trim().length > 3 &&
        this.fields.description.trim().length > 4 &&
        this.selectedTags.length > 0 &&
        this.selectedCountries.length > 0 &&
        !!this.$store.state.user
      ) {
        return true;
      }
    },
  },

  methods: {
    async submit() {
      if (this.valid && !this.flight && this.hasAddress) {
        this.flight = true;
        const type = { type: "item" };
        const fields = [
          ["creator", 1, "string", "optional"],
          ["title", 2, "string", "optional"],
          ["description", 3, "string", "optional"],
          ["shippingcost", 4, "int64", "optional"],
          ["localpickup", 5, "bool", "optional"],
          ["estimationcount", 6, "int64", "optional"],
          ["tags", 7, "string", "repeated"],
          ["condition", 8, "int64", "optional"],
          ["shippingregion", 9, "string", "repeated"],
          ["depositamount", 10, "int64", "optional"],
        ];
        //const body = [this.$store.state.account.address,"dsaf", "asdf", 33, 1, "sdfsdf", "asdf", 4, "sfda"]
        const body = {
          //creator: this.$store.state.account.address,
          title: this.fields.title,
          description: this.fields.description,
          shippingcost: this.fields.shippingcost,
          localpickup: this.fields.localpickup,
          estimationcount: this.fields.estimationcount,
          tags: this.selectedTags,
          condition: this.fields.condition,
          shippingregion: this.selectedCountries,
          depositamount: this.fields.estimationcount,
        };

        await this.itemSubmit({ ...type, fields, body });

        this.flight = false;
        //this.fields.title = "";
        // this.fields.description = "";
        //this.fields.shippingcost = "";
        // this.fields.localpickup = false;
        //this.fields.estimationcount = "";
        this.itemid = await this.$store.state.newitemID;
        //console.log()
        console.log(this.itemid);
        this.thisitem = await this.$store.getters.getItemByID(this.itemid);
        this.e1 = 2;
        this.showpreview = true;
        //alert("Submitted, find the item in the account section");
      }
    },
    updateStepCount(e1) {
      this.e1 = e1;
    },

    async itemSubmit({ type, fields, body }) {
      const wallet = this.$store.state.wallet;
      const type2 = type.charAt(0).toUpperCase() + type.slice(1);
      const typeUrl = `/${process.env.VUE_APP_PATH}.MsgCreate${type2}`;
      let MsgCreate = new Type(`MsgCreate${type2}`);
      const registry = new Registry([[typeUrl, MsgCreate]]);
      fields.forEach((f) => {
        MsgCreate = MsgCreate.add(new Field(f[0], f[1], f[2], f[3]));
      });

      const [firstAccount] = await wallet.getAccounts();

      const client = await SigningStargateClient.connectWithSigner(
        process.env.VUE_APP_RPC,
        wallet,
        { registry }
      );

      const msg = {
        typeUrl,
        value: {
          creator: this.$store.state.account.address,
          ...body,
        },
      };

      const fee = {
        amount: [{ amount: "0", denom: "tpp" }],
        gas: "200000",
      };
      await this.$store.dispatch("entityFetch", {
        type: type,
      });
      await this.$store.dispatch(
        "setCreatorItemList",
        this.$store.state.account.address
      );
      const selleritems = this.$store.state.creatorItemList || [];

      try {
        const result = await client.signAndBroadcast(
          firstAccount.address,
          [msg],
          fee
        );

        assertIsBroadcastTxSuccess(result);
        await this.$store.dispatch("entityFetch", {
          type: type,
        });
        await this.$store.dispatch(
          "setCreatorItemList",
          this.$store.state.account.address
        );
        let newselleritems = this.$store.state.creatorItemList.map(
          (item) => item.id
        );
        let sorted = newselleritems.sort(
          (selleritems, newselleritems) => newselleritems - selleritems
        );
        console.log(sorted);
        //et len = (selleritems.length)
        // console.log((newselleritems[len].id))
        //this.$store.commit('set', { key: 'newitemID', value: (newselleritems[len].id) })
        this.$store.commit("set", { key: "newitemID", value: sorted[0] });
        await this.$store.dispatch(
          "setSellerItemList",
          this.$store.state.account.address
        );
      } catch (e) {
        console.log(e);
      }
    },

    conditionLabel() {
      if (this.fields.condition === 0) {
        return "'bad'";
      }
      if (this.fields.condition === 1) {
        return "'fixable'";
      }
      if (this.fields.condition === 2) {
        return "'decent'";
      }
      if (this.fields.condition === 3) {
        return "'as new'";
      }
      if (this.fields.condition === 4) {
        return "'perfect'";
      }
    },
  },
};
</script>

