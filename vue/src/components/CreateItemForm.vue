<template>
  <div class="pa-2 mx-lg-auto"> 
    <v-dialog v-model="dialog" width="500">
      <template v-slot:activator="{ on, attrs }">
       <v-img rounded v-bind="attrs"
          v-on="on" height="300" src="img/design/sell.png">
      <span v-if="!data.title" >
             <p  v-if="!data.title" v-bind="attrs"
          v-on="on"  class="headline pt-4 font-weight-thin gray--text text-center"> Place An Item</p>
        
        </span><span  v-else>
          <p  v-bind="attrs"
          v-on="on"  class="headline pt-4 font-weight-thin gray--text text-center"> Place {{data.title}}</p>
     </span>
       <v-icon class="ml-4" small>mdi-information-outline</v-icon></v-img>
      </template>

      <v-card class="text-center">
        <v-card-title class="h2 lighten-2"> Info </v-card-title>

        <v-card-text>
          After placing the item, an estimation will be made. After you accept this price, anyone can buy the item!
        </v-card-text>
 <iframe width="100%"  height="281" src="https://www.youtube.com/embed/zHXwfePrGvA?vq=hd1080&autoplay=0&loop=1&modestbranding=1&rel=0&cc_load_policy=1&color=white&mute=1" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>
        <v-divider></v-divider>

        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="primary" text @click="dialog = false"> Let's go </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    
    <v-stepper class="elevation-0" v-model="e1" >
      <v-stepper-header>
        <v-stepper-step :complete="e1 > 1" step="1"> Data </v-stepper-step>

        <v-divider></v-divider>

        <v-stepper-step :complete="e1 > 2" step="2"> Pictures </v-stepper-step>
        <v-divider></v-divider>

        <v-stepper-step :complete="e1 > 3" step="3"> Done! </v-stepper-step>
      </v-stepper-header>

      <v-stepper-items>
    
          <v-stepper-content  class="pa-2 ma-2" step="1" >
      <div  class=" pr-4 ma-2">
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
                v-model="data.title"
                required
              />
<v-alert class="ma-2 caption" dense type="info" dismissible v-if="descrinfo"> Make sure to fully disclose any defects or scratches (and highlight these in the pictures)</v-alert>
              <v-textarea
              @change="descrinfo = !descrinfo"
                class="ma-1"
                prepend-icon="mdi-text"
                :rules="rules.descriptionRules"
                v-model="data.description"
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
                      <span v-if="search">   No category tags matching "<strong>{{ search }}</strong
                        >". Press <kbd>enter</kbd> to create a new one</span>
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
                  ><v-btn text icon @click="data.estimationcount = 3">
                    <v-icon> mdi-check</v-icon></v-btn
                  >
                  <v-slider 
                    hint="Lower for faster results, higher for better accuracy"
                    thumb-label
                    :persistent-hint="data.estimationcount != 3"
                    label="Accuracy"
                    :thumb-size="90"
                    max="12"
                    :rules="rules.estimationcountRules"
                    placeholder="Estimation count"
                    v-model="data.estimationcount"
                    ><template v-slot:thumb-label="item">
                      {{ item.value }} Estimations
                    </template></v-slider
                  >
                </v-row>
                <v-row class="pa-2">
                  <v-btn text icon @click="data.condition = 0">
                    <v-icon>{{
                      data.condition === 0 ? "mdi-star-outline" : "mdi-star"
                    }}</v-icon>
                  </v-btn>
                  <v-slider
                    label="Condition"
                    :hint="
                      'Condition is ' +
                      conditionLabel() +
                      ', please explain condition in description'
                    "
                    v-model="data.condition"
                    :max="4"
                    :persistent-hint="data.condition != 0"
                    :thumb-size="24"
                    thumb-label
                    ><template v-slot:thumb-label="{ value }">
                      {{ satisfactionEmojis[value] }}
                    </template>
                  </v-slider></v-row
                >

                <v-row class="pa-2 mt-2">
                  <v-btn text icon @click="data.shippingcost = 0">
                    <v-icon>
                      {{
                        data.shippingcost === 0
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
                    :persistent-hint="data.shippingcost != 0"
                    placeholder="Shipping cost"
                    :thumb-size="60"
                    thumb-color="primary lighten-1"
                    v-model="data.shippingcost"
                    ><template v-slot:thumb-label="item">
                      {{ item.value }} <v-icon>$vuetify.icons.custom</v-icon>
                    </template>
                  </v-slider>
                </v-row>
                <v-row><v-col>
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
                >
                  <v-col>
                    <v-row class="ma-0">
                      <v-btn
                        class="mr-auto"
                        text
                        icon
                        @click="enterlocation = !enterlocation"
                        ><v-icon
                          >{{
                            enterlocation
                              ? "mdi-map-marker"
                              : "mdi-map-marker-off"
                          }}
                        </v-icon></v-btn
                      >

                      <v-switch
                        class="mr-auto mt-1"
                        v-model="enterlocation"
                        inset
                        label="Pickup"
                        :persistent-hint="
                        
                          enterlocation == true 
                     
                        "
                        hint="Specify location"
                      ></v-switch> </v-row>
                        <v-text-field
                class="ma-1"
                prepend-icon="mdi-map-marker"
                :rules="rules.pickupRules"
                label="Location"
                v-model="data.localpickup"
                required v-if="enterlocation"
              />
              </v-col
                  ></v-row>
              </div><div class="mx-auto text-center" v-if="valid">
                <span class="caption"> Required deposit for price estimators: <v-icon small left>$vuetify.icons.custom</v-icon>{{data.depositamount}} TPP. <v-btn @click="changedeposit = !changedeposit" icon small> <v-icon >
        mdi-pencil
      </v-icon></v-btn></span> 
        <v-row v-if="changedeposit">
          <v-col cols="6" class="mx-auto">
            <v-dialog transition="dialog-bottom-transition" max-width="600">
              <template v-slot:activator="{ on, attrs }">
                <span v-bind="attrs" v-on="on"
                
                ><v-icon class="ml-4 align-center" small>mdi-information-outline</v-icon>
           
                
                </span>
              </template>
              <template v-slot:default="dialog">
                <v-card>
                  <v-toolbar color="default">Info <v-spacer></v-spacer>
          <v-btn
            color="primary"
            icon
            @click="dialog.value = false"
          >
          <v-icon> mdi-close</v-icon>
          </v-btn></v-toolbar>
                  <v-card-text>
                    <div class="text-p pt-4">
                     A deposit is required for all the estimators. You may change this, but don't set it too high because 1) then no one is willing to estimate it 2) when it is higher than the final price, no action is taken. Estimators risk their deposit to make an estimation for you.
                    </div>
                    <div class="caption pa-2">
                      - The lowest estimator loses this
                      when the final estimation price is not accepted by
                      you.
                    </div>
                    <div class="caption pa-2 mb-2">
                      - The highest estimator loses this
                      when the item is not bought by the buyer that provided
                      prepayment.
                    </div>
                     
                     
                                       </v-card-text>
                   <v-img class="mx-12" src="img/design/transfer.png" ></v-img><div class="caption text-center pt-4">
                      Good luck and have fun!
                    </div>
                  <v-card-actions class="justify-end">
                    <v-btn text @click="dialog.value = false">Close</v-btn>
                  </v-card-actions>
                </v-card>
              </template>
            </v-dialog>
          </v-col>
          <v-col cols="6" class="mx-auto">
            <v-text-field
              label="Amount"
              type="number"
              v-model="data.depositamount"
           
            suffix="TPP"
              prepend-icon="$vuetify.icons.custom"
            
            ></v-text-field>
          </v-col>
        </v-row>
      </div>
              <div class="text-center pt-6">
                <v-btn rounded
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
                    >Awaiting transaction creating item ID...
                  </div>
                </v-btn>
              </div>
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
    </v-stepper><sign-tx v-if="submitted" :key="submitted" :fields="fields" :value="value" :msg="msg" @clicked="afterSubmit"></sign-tx>
  </div>
</template>


<script>

import CreateItemPreviewAndUpload from "./CreateItemPreviewAndUpload.vue";


export default {
  components: { CreateItemPreviewAndUpload },
  data() {
    return {
      data: {
        title: "",
        description: "",
        shippingcost: "0",
        localpickup: "",
        estimationcount: "3",
        condition: "0",
depositamount: "3",
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
        pickupRules: [
      
        
          (v) =>
            ( v.length <= 25) || "Pickup must be less than 25 characters, enter coordinates instead",
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
changedeposit: false,
 descrinfo: false,
      satisfactionEmojis: ["😭", "🙁", "🙂", "😊", "😄"],
      countryCodes: ["NL", "BE", "UK", "DE", "US", "CA"],
      enterlocation: false,

        fields: [],
      value: {},
      msg: "",
      submitted: false,
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
          "Watch",
           "Shoes",
          "Clothing",
          "Collectable",
          
        //  "Garden item",
          "Vehicle",
         // "Motor",
          //"Sport",
           "Book",
         // "Antique",
          "Computer",
                "Smartphone",
          "Smart Device",
          "Sound Device",
          "TV",
          "NFT",
          "Other",
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
        this.data.title.trim().length > 3 &&
        this.data.description.trim().length > 4 &&
        this.selectedTags.length > 0 &&
        this.selectedCountries.length > 0 &&
        !!this.$store.state.user
           &&
        (this.data.shippingcost || this.data.localpickup)
      ) {
        return true;
      }
    },
  },

  methods: {
    async submit() {
      if (this.valid && !this.flight && this.hasAddress) {
        this.flight = true;

        this.fields = [
          ["creator", 1, "string", "optional"],
          ["title", 2, "string", "optional"],
          ["description", 3, "string", "optional"],
          ["shippingcost", 4, "int64", "optional"],
          ["localpickup", 5, "string", "optional"],
          ["estimationcount", 6, "int64", "optional"],
          ["tags", 7, "string", "repeated"],
          ["condition", 8, "int64", "optional"],
          ["shippingregion", 9, "string", "repeated"],
          ["depositamount", 10, "int64", "optional"],
        ];
        //const body = [this.$store.state.account.address,"dsaf", "asdf", 33, 1, "sdfsdf", "asdf", 4, "sfda"]
        const body = {
          //creator: this.$store.state.account.address,
          title: this.data.title,
          description: this.data.description,
          shippingcost: this.data.shippingcost,
          localpickup: encodeURI(this.data.localpickup),
          estimationcount: this.data.estimationcount,
          tags: this.selectedTags,
          condition: this.data.condition,
          shippingregion: this.selectedCountries,
          depositamount: this.data.depositamount,
        };

        this.msg = "MsgCreateItem"

        this.value = {
          creator: this.$store.state.account.address,
          ...body,
        }
        
      this.submitted = true

      
    } },
         async afterSubmit(value) {
           
 this.loadingitem = true;

 this.msg = ""
 this.fields = []
 this.value = {}
  if(value == true) {
    const selleritems = this.$store.state.creatorItemList || []

            console.log(selleritems);
         
        const type = { type: "item" }
    await this.$store.dispatch("entityFetch",
          type,
        );
        await this.$store.dispatch(
          "setCreatorItemList",
          this.$store.state.account.address
        );
            console.log("dfsfaegdsfa");
        let newselleritems = this.$store.state.creatorItemList.map(
          (item) => item.id
        );
        let sorted = newselleritems.sort(
          (selleritems, newselleritems) => newselleritems - selleritems
        );
        console.log(sorted);

        this.$store.commit("set", { key: "newitemID", value: sorted[0] });
        await this.$store.dispatch(
          "setSellerItemList",
          this.$store.state.account.address
        );

       

        this.itemid = await this.$store.state.newitemID;

        console.log(this.itemid);
        this.thisitem = await this.$store.getters.getItemByID(this.itemid);
        this.e1 = 2;
        this.showpreview = true;
       
  

         
 
       
  }   this.loadingitem = false;  
  this.submitted = false
    this.flight = false;  
     


    },
    updateStepCount(e1) {
      this.e1 = e1;
    },


    conditionLabel() {
      if (this.data.condition === 0) {
        return "'bad'";
      }
      if (this.data.condition === 1) {
        return "'fixable'";
      }
      if (this.data.condition === 2) {
        return "'decent'";
      }
      if (this.data.condition === 3) {
        return "'as new'";
      }
      if (this.data.condition === 4) {
        return "'perfect'";
      }
    },
  },
};
</script>

