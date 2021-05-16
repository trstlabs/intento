<template >
  <div class="pa-2 mx-lg-auto">
    <div class="text-center pa-12" v-if="!showinfo">
      <v-btn :ripple="false" text @click="getItemToEstimate"
        ><v-icon color="primary" left> mdi-refresh </v-icon> Refresh
      </v-btn><v-img class="mx-12" src="img/design/estimate.png" ></v-img>
    </div>
    <v-skeleton-loader background-color="inherit"
      v-if="loadingitem"
      class="mx-auto"
      type="list-item-three-line, image, article"
    ></v-skeleton-loader>

<div v-if="showinfo == true && loadingitem == false" >
    <div elevation="8" v-if="photos[0]"  >
   
      <v-carousel v-if="magnify == false" style="height:100%"
        delimiter-icon="mdi-minus"
        carousel-controls-bg="primary"
        height="300" 
        hide-delimiter-background
        show-arrows-on-hover 
      > 
       <v-carousel-item max-height="300" 
    contain v-for="(photo, i) in photos" :key="i" :src="photo" > 

    <template v-slot:placeholder>
        <v-row
          class="fill-height ma-0"
          align="center"
          justify="center"
        >
          <v-progress-circular
            indeterminate
            color="grey lighten-3"
          ></v-progress-circular>
        </v-row>
      </template>
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

    <template v-slot:placeholder>
        <v-row
          class="fill-height ma-0"
          align="center"
          justify="center"
        >
          <v-progress-circular
            indeterminate
            color="grey lighten-5"
          ></v-progress-circular>
        </v-row>
      </template>
        </v-carousel-item>
      </v-carousel>
    </div>  <v-row class="ml-4 mt-1 mb-1">
      <span  v-for="(photo, index) in photos" :key="index"> <img class="ma-1" @click="show(photo)" height="56"  :src="photo" /></span><v-spacer/><v-btn x-small class="mr-2"
            color="primary"
            icon
            @click="settings = !settings"
          >
          <v-icon> mdi-tune</v-icon>
          </v-btn><v-btn x-small class="mr-4"
            color="primary"
            icon
            @click="magnify = !magnify"
          >
          <v-icon> mdi-crop-free</v-icon>
          </v-btn></v-row> 
<v-dialog
      v-model="fullscreen"
    
    >
     

      <v-card >
        <v-card-title class=" grey lighten-2 ">
         {{item.title}} <v-spacer></v-spacer>
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
    <div v-if="settings"> 
          <v-card class="pa-6 elevation-8 ma-6 rounded-xl" ><v-row  class="mb-2"><v-btn small
        icon
          
            @click="getItemToEstimate()"
          >
          <v-icon> mdi-refresh</v-icon>
          </v-btn> <v-spacer/><v-btn small
        
            icon
            @click="settings = !settings"
          >
          <v-icon> mdi-close</v-icon>
          </v-btn>
          </v-row> <v-select
          append-icon="mdi-tag-outline"
          dense
          v-model="selectedFilter"
          v-on:input="updateList(selectedFilter)"
          cache-items
          :items="tags"
          label="Categories"
          clearable
          solo
          
          :persistent-hint="!selectedFilter"
          hint="Select item category"
        ></v-select>
         <v-select
          append-icon="mdi-tag-outline"
          dense
          v-model="selectedRegion"
          v-on:input="updateRegionList(selectedRegion)"
          cache-items
          :items="locations"
          label="Regions"
          clearable
          solo
          :persistent-hint="!selectedRegion"
          hint="Specify your region"
        ></v-select> </v-card></div>
    <v-card 
      class="pa-2 mt-2"
      elevation="2"
      rounded="xl"
     
    >
      <v-progress-linear
        indeterminate
        :active="loadingitem"
      ></v-progress-linear>

      <div class="pa-2">
        <v-row>
          <v-col class="pa-2">
            <span elevation="0">
              <div class="overline">Title</div>

              <div class="body-1 font-weight-light">
                {{ item.title }}
              </div>
            </span>
          </v-col>
           <v-col class="pa-2 font-weight-light ">
            <span>
              <v-chip-group >
                   <v-chip class="ma-1 text-capitalize"
                  outlined :to="{ name: 'SearchTag', params: { tag: itemtag } }"
                  small
                  v-for="itemtag in item.tags"
                  :key="itemtag"
                >
                  <v-icon small left> mdi-tag-outline </v-icon>
                  {{ itemtag }}
                </v-chip>
               
                
              </v-chip-group>  <v-chip-group>  <v-chip class="ma-1"
                  outlined
                  small
                 :to="{ name: 'BuyItemDetails', params: { id: item.id } }"
                >
                  <v-icon small left> mdi-account-badge </v-icon>
                  TPP ID: {{ item.id }}
                </v-chip>


                 <v-chip v-if="item.flags > 0" class="ma-1 caption"  small outlined>
              <v-icon small left> mdi-shield-alert-outline </v-icon>
              <span v-if="item.flags == 1"> Reported {{ item.flags }} time</span><span v-else>Reported {{ item.flags }}  times</span>
            </v-chip>

              <v-dialog transition="dialog-bottom-transition" max-width="300">
                <template v-slot:activator="{ on, attrs }">
                  <span v-bind="attrs" v-on="on">
                   <v-chip class="ma-1"  outlined small>
                    <v-rating
                      :value="Number(item.condition)"
                      readonly
                      color="gray lighten-1"
                      background-color="gray"
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
              </v-dialog></v-chip-group>
            </span>
          </v-col>
        </v-row>
        <span>
          <div class="overline ">Description</div>

          <div class="caption font-weight-light">
            {{ item.description }}
          </div>
        </span>
      </div>

      <div class="pa-2 mx-auto text-center" elevation="8" v-if="lastitem">
        <v-chip
          v-if="lastitem"
          class="mt-2"
          label
          outlined
          medium
          color="warning"
        >
          <v-icon left> mdi-alarm </v-icon>
          This was the last item, check again later.
        </v-chip>
      </div>

       <v-divider class="mx-4 " />

      <div class="mx-auto">
        <v-row>
          <v-col cols="6" class="mx-auto">
            <v-dialog transition="dialog-bottom-transition" max-width="600">
              <template v-slot:activator="{ on, attrs }">
                    <span v-bind="attrs" v-on="on">
                   <v-icon class="ml-4" small>mdi-information-outline</v-icon>  <p class="text-center caption ">To me, it's worth</p><span v-if="!flight && valid && !timeout" class="caption font-weight-light">Deposit is {{item.depositamount}}<v-icon x-small right>$vuetify.icons.custom</v-icon> 
              
            </span></span>
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
                      Earn a reward when you are the best
                      estimator! For this item, the reward is of {{item.depositamount}} <v-icon x-small right>$vuetify.icons.custom</v-icon>  and is equal to the deposit. However, your deposit is lost when:
                     </div>
                    <div class="caption pa-2">
                      - You are the lowest estimator
                      and the final estimation price is not accepted by the
                      seller.
                    </div>
                    <div class="caption pa-2 mb-2">
                      - You are the highest estimator
                      and the item is not bought by the buyer that provided
                      prepayment.
                    </div>
                     
                     <iframe style="
      

-webkit-mask-image: -webkit-radial-gradient(circle, white 100%, black 100%); /*ios 7 border-radius-bug */
-webkit-transform: rotate(0.000001deg); /*mac os 10.6 safari 5 border-radius-bug */
-webkit-border-radius: 10px; 
-moz-border-radius: 10px;
border-radius: 20px; 
overflow: hidden; 
" width="100%"  height="310" src="https://www.youtube.com/embed/zHXwfePrGvA?start=61&vq=hd1080&autoplay=0&loop=1&modestbranding=1&rel=0&cc_load_policy=1&color=white&mute=1" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>
                               
                  </v-card-text>
                   <v-img class="mx-12" src="img/design/estimate.png" ></v-img><div class="caption text-center pt-4">
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
              v-model="estimation"
              :disabled="lastitem"
            suffix="TPP"
              prepend-icon="$vuetify.icons.custom"
            
            ></v-text-field>
          </v-col>
        </v-row>
      </div>
      <v-divider class="mx-4" />
      <span>
        <div class="pa-2">
          <div>
            <v-chip-group active-class="primary--text " column>
              <v-chip class="font-weight-light "
                :disabled="lastitem"
                small
                v-for="(option, text) in options"
                :key="text"
                @click="updateComment(option.attr)"
              >
                {{ option.name }}
              </v-chip>
            </v-chip-group>
          </div>

          <div class="mx-auto">
            <h3 class="text-left">“</h3>
            <v-text-field
              rounded
              dense
              :disabled="lastitem"
              clearable
              class="caption"
              placeholder="leave a comment (optional)"
              v-model="comment"
            />
          </div>
          <h3 class="text-right">”</h3>
        </div>
      </span>

      <div> 
        <v-btn rounded
          block
          elevation="4"
          color="primary"
          :disabled="!valid || !hasAddress ||  flight || timeout"
          @click="submit(estimation, item.id, interested, comment)"
          ><div v-if="!flight && !valid">
            <v-icon left> mdi-check </v-icon> Estimate item
          </div>
          <div v-if="!flight && valid && !timeout">
            <v-icon > mdi-check-bold </v-icon> Estimate item 
          </div>
          <div v-if="!timeout && !valid && flight">
            <v-icon left> mdi-check </v-icon> Estimate item
          </div>       
             <v-progress-linear v-if="timeout && valid"
      :rotate="360"
      :size="50"
      :width="5"
      :value="timeoutvalue"
      color="white"
    >
      
    </v-progress-linear>
           <div v-if="flight">
              <div class="text-right">Awaiting submission...</div>
            </div>       
        </v-btn>
        <div v-if="timeout && valid">
              <div class="text-right caption">Wait {{ 10-(timeoutvalue/6) }}s</div>
              
            </div>

         
      
      </div> 
    </v-card>
</div>
<sign-tx v-if="submitted" :key="submitted" :fields="fields" :value="value" :msg="msg" @clicked="afterSubmit"></sign-tx>
    <div class="pa-4 mx-auto" v-if="showinfo">
      <v-row class="text-center">
        <v-col class="pa-0">
          <v-tooltip bottom :disabled="!interested" v-if="showinfo">
            <template v-slot:activator="{ on, attrs }">
              <span v-bind="attrs" v-on="on">
                <v-btn
                  class="mx-2"
                  fab
                  dark
                  small
                  color="pink"
                  icon
                  @click="interested = !interested"
                >
                  <v-icon dark> mdi-heart </v-icon>
                </v-btn>
              </span>
            </template>
            <span
              >Find your liked items in the account section when they are
              available.
            </span>
          </v-tooltip> </v-col
        ><v-col class="pa-0">
       

          <v-dialog
            bottom
            :disabled="flag"
            v-if="showinfo"
            v-model="dialog"
            persistent
            max-width="350"
          >
            <template v-slot:activator="{ on, attrs }">
              <v-btn
                class="mx-2"
                fab
                dark
                small
                v-bind="attrs"
                v-on="on"
                color="red"
                icon
              >
                <v-icon dark> mdi-alert-octagon </v-icon>
              </v-btn>
            </template>
            <v-card>
              <v-card-title class="headline"> Report item? <v-icon  @click="reportinfo=!reportinfo" class="ml-4" small>mdi-information-outline</v-icon></v-card-title>
              <v-card-text
                >  <span v-if="reportinfo">
            
          
          If this item is not OK, you can report it here. The protocol
                removes items that are reported
                often. This way they don't make it onto the marketplace. Thanks for keeping TPP safe and helping others.
              <v-divider class="ma-4"/>
             </span>
                 <span class="pt-2 ma-0 subtitle-1"> Please report it if the item is:
                <p class="caption"> Fake </p><p class="caption"> In bad condition </p><p class="caption"> From an untrustworthy seller</p><p class="caption"> Not functioning </p><p class="caption"> Using wrong data </p>
             </span>
           
               
                </v-card-text
              >
              <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn color="primary darken-1" text @click="dialog = false">
                  Discard
                </v-btn>
                <v-btn color="red darken-1" text @click="submitFlag()">
                  Report Item
                </v-btn>
              </v-card-actions>
            </v-card>
          </v-dialog> </v-col
        ><v-col class="pa-0">
          <v-btn
            :disabled="estimation > 1 || !hasAddress || !showinfo"
            icon
            @click="getNewItemByIndex"
            color="primary"
          >
            <v-icon dark> mdi-arrow-right-bold </v-icon>
          </v-btn>
        </v-col>
      </v-row>

     <!-- <div class="pt-12 mx-lg-auto">
      
        <v-select
          append-icon="mdi-tag-outline"
          dense
          v-model="selectedFilter"
          v-on:input="updateList(selectedFilter)"
          cache-items
          :items="tags"
          label="Categories"
          clearable
          rounded
          solo
          persistent-hint
          hint="Specify your expertise"
        ></v-select>
         <v-select
          append-icon="mdi-tag-outline"
          dense
          v-model="selectedRegion"
          v-on:input="updateRegionList(selectedRegion)"
          cache-items
          :items="locations"
          label="Regions"
          clearable
          rounded
          solo
          persistent-hint
          hint="Specify your region"
        ></v-select>
      </div>-->
    </div>
  </div>
</template>

<script>
import ToEstimateTagBar from "./ToEstimateTagBar.vue";
import { databaseRef } from "./firebase/db.js";

export default {
  
  components: { ToEstimateTagBar },
  data() {
    return {
      estimation: "",
      comment: "",

      options: [
        { name: "Great Photos!", attr: "Great Photos!" },
        { name: "Unclear Photos", attr: "I find the photos unclear." },
        {
          name: "Excellent Description",
          attr: "I find the description excellent.",
        },
        { name: "Too Vage", attr: "I find the description too vague." },
        {
          name: "Clear",
          attr:
            "I find the item well described, the buyer will know what to expect.",
        },
        { name: "Looks damaged", attr: "This item appears to be damaged." },
        {
          name: "Repairable",
          attr: "This item seems damaged, but I think it can be repaired.",
        },
        { name: "Used", attr: "This item seems used, so I find it worth less compared similar items" },
        {
          name: "As good as new!",
          attr: "This item appears to look as good as new!",
        },
        { name: "Dirty", attr: "The item looks dirty to me, so I find it worth less compared similar items" },
         {
          name: "Replica",
          attr: "The item looks like a replica to me, please don't sell if this is a replica",
        },

         {
          name: "Serviced",
          attr: "This item seems serviced to me.",
        },
        {
          name: "Can't tell from pics",
          attr: "I can't tell the condition from the pictures, would see it in person.",
        },
         {
          name: "More info",
          attr: "I can't tell the condition from the data, would ask for more info.",
        },
         {
          name: "I'd buy",
          attr: "I like this kind of item, I'd buy it.",
        },
        
      ],

      interested: false,
      flag: false,
      flight: false,
      item: "",
      index: 0,
      showinfo: false,
      lastitem: false,
      loadingitem: false,
      photos: [],
      selectedFilter: "",
      selectedRegion: "",
      timeout: false,
      showphoto: "",
interval: {},
timeoutvalue: 0,
      dialog: false,
      fullscreen: false,
      magnify: false,
      conditionInfo: false,
      settings: false,
      reportinfo: false,


      fields: [],
      value: {},
      msg: "",
      submitted: false,
    };
  },
  beforeDestroy () {
      clearInterval(this.interval)
    },

  mounted() {

    
      
    //console.log(input)
    if (!!this.$store.state.account.address) {
      let input = this.$store.state.account.address;
 this.$store.dispatch("setToEstimateRegions");
     const type = { type: "estimator" };
    this.$store.dispatch("entityFetch",type);
      this.$store.dispatch("setSortedTagList");    
      this.$store.dispatch("setEstimatorItemList", input);
       this.$store.dispatch("setSellerItemList", input);
      this.$store.dispatch("setToEstimateList", this.index);

      this.item = this.items[this.index];
      this.loadItemPhotos();
      this.showinfo = true;
    }
  },

  computed: {
    items() {
      return this.$store.getters.getToEstimateList;
    },
    hasAddress() {
      return !!this.$store.state.account.address;
    },
    valid() {
      return this.estimation.trim().length > 0;
    },
    tags() {
      return this.$store.getters.getTagList;
    },
    locations() {
      return this.$store.getters.getRegionList;
    },
  },

  methods: {
    async submit(estimation, itemid, interested, comment) {
      if (this.valid && !this.flight && this.hasAddress) {
        this.flight = true;
       
 
 
        const body = {
          deposit: this.item.depositamount,
          estimation: estimation,
          itemid: itemid,
          interested: interested,
          comment: comment,
        };
        this.msg = "MsgCreateEstimator"
        this.fields = [
          ["estimator", 1, "string", "optional"],
          ["estimation", 2, "int64", "optional"],
          ["itemid", 3, "string", "optional"],
          ["deposit", 4, "int64", "optional"],
          ["interested", 5, "bool", "optional"],
          ["comment", 6, "string", "optional"],
        ];

        this.value = {
          estimator: this.$store.state.account.address,
          ...body,
        }

        this.submitted = true

      }},

  

      async afterSubmit(value){
 this.loadingitem = true;
 this.submitted = false
 this.msg = ""
 this.fields = []
 this.value = {}


      console.log(value)
      //  await this.estimationSubmit({ ...type,"body, fields });

       if(value == true){
       //this.$store.dispatch("entityFetch", "estimator");
       
             await this.$store.dispatch("bankBalancesGet");
        this.timeout = true
        clearInterval(this.interval)
        this.timeoutvalue = 0
        this.interval = setInterval(() => {
      
        this.timeoutvalue += 6
      }, 600)


         setTimeout(() => this.timeout = false, 6000);
     

       }
        this.getNewItemByIndex();
        //await this.submitRevealEstimation(itemid);
         this.estimation = "";
        this.comment = "";
        //this.flight = false;
        //this.loadingitem = false;
      },

   /* async estimationSubmit({ type, fields, body }) {
      const wallet = this.$store.state.wallet;
      const type2 = type.charAt(0).toUpperCase() + type.slice(1);
      const typeUrl = `/${process.env.VUE_APP_PATH}.MsgCreate${type2}`;
      let MsgCreate = new Type(`MsgCreate${type2}`);
      const registry = new Registry([[typeUrl, MsgCreate]]);
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
          estimator: this.$store.state.account.address,
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
       this.getNewItemByIndex();
      console.log("success!");

      
    },
*/
    async submitFlag() {
      if (!this.flight && this.hasAddress) {
        this.flight = true;
        this.loadingitem = true;
        this.flag = true;
              this.msg = "MsgCreateFlag"

        const body = { flag: true, itemid: this.item.id };
        this.fields = [
          ["estimator", 1, "string", "optional"],
          ["itemid", 2, "string", "optional"],
        
        ];

 this.value = {
          estimator: this.$store.state.account.address,

          ...body,
        },
            this.dialog = false;
        this.flag = false;
        this.submitted = true

    
      }
    },

 /*   async flagSubmit({ type, fields, body }) {
      const wallet = this.$store.state.wallet;
      const type2 = type.charAt(0).toUpperCase() + type.slice(1);
      const typeUrl = `/${process.env.VUE_APP_PATH}.MsgCreateFlag`;
      let MsgCreate = new Type(`MsgCreateFlag`);
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
          estimator: this.$store.state.account.address,

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
      try {
        await this.$store.dispatch("entityFetch", {
          type: type,
        });
      } catch (e) {
        console.log(e);
      }
    },
*/
    async getItemToEstimate() {
      if (!this.hasAddress) {
        alert("Sign in first");
        window.location.reload()
        //return (this.showinfo = false);
      }
      let input = this.$store.state.account.address;
      
      this.$store.dispatch("setSortedTagList");
      this.$store.dispatch("setEstimatorItemList", input);
      this.$store.dispatch("setToEstimateList");

      //let index = 0;
      //this.$store.dispatch("setToEstimateList");
      this.item = this.items[this.index];
      this.lastitem = false;
      this.loadItemPhotos();
      this.showinfo = true;
    },

    async getNewItemByIndex() {
      this.loadingitem = true;
      let oldindex = this.index;
      if (oldindex >= 0 && oldindex < this.items.length - 1) {
        this.index = oldindex + 1;
      }

      console.log(oldindex, this.index);
      this.item = this.items[this.index];
     
    
      if (oldindex === this.index) {
        this.lastitem = true;
      }
      this.loadItemPhotos();
    
      
    },

    show(photo){
      this.showphoto = photo 

      this.fullscreen = true
    },
    loadItemPhotos() {
      this.loadingitem = true;
       if(!this.item){ if(!this.lastitem)
       {
        this.getNewItemByIndex()}else {this.loadingitem = false, alert("No items found to be estimated, try again later.")}

       }else{
      const id = this.item.id;
      //const db = firebase.database();

     const imageRef = databaseRef.ref("ItemPhotoGallery/" + id + "/photos/");
      imageRef.on("value", (snapshot) => {
        const data = snapshot.val();

        if (data != null ) {
        
          // console.log(data[0]);
          this.photos = data;
          //this.photos = { photo: "https://upload.wikimedia.org/wikipedia/commons/thumb/8/85/Points_of_a_horse.jpg/330px-Points_of_a_horse.jpg" };

          this.loadingitem = false;
        }else{ if(this.lastitem){this.loadingitem = false;}else{
          this.photos = []
          this.getNewItemByIndex()}
        }
      });
      //this.loadingitem = false;
      this.interested = false;
      this.flight = false;
        }
    },

 /*   async submitRevealEstimation(itemid) {
      if (this.hasAddress) {
       
        this.getNewItemByIndex();
        const fields = [
          ["creator", 1, "string", "optional"],
          ["itemid", 2, "string", "optional"],
        ];
        // const type = { type: "item" };
        const body = { itemid: itemid };
        this.revealSubmit({ body, fields });
      }
    },
    async revealSubmit({ body, fields }) {
      const wallet = this.$store.state.wallet;
      const typeUrl = `/${process.env.VUE_APP_PATH}.MsgRevealEstimation`;
      let MsgCreate = new Type(`MsgRevealEstimation`);
      const registry = new Registry([[typeUrl, MsgCreate]]);

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
          creator: this.$store.state.account.address,
          ...body,
        },
      };
      const fee = {
        amount: [{ amount: "0", denom: "tpp" }],
        gas: "200000",
      };

      await client.signAndBroadcast(
        this.$store.state.account.address,
        [msg],
        fee
      );
    },
*/
    updateComment(newComment) {
      this.comment = newComment;
    },
    updateList(tag) {
      //console.log(this.tag);
      this.$store.dispatch("tagToEstimateList", tag);
      if (!!this.items[0]) {
        this.item = this.items[0];
        if (!this.items[0]) {
          this.lastitem = true;
        } else {
          this.lastitem = false;
        }
        this.loadItemPhotos();
      } else {
        alert("No items to estimate for " + tag);
        //this.$store.dispatch("setToEstimateList");
        this.getItemToEstimate();
      }
    },

    updateRegionList(region) {
      this.$store.dispatch("regionToEstimateList", region);
      if (!!this.items[0]) {
        this.item = this.items[0];
        if (!this.items[0]) {
          this.lastitem = true;
        } else {
          this.lastitem = false;
        }
        this.loadItemPhotos();
      } else {
        alert("No Items to estimate for " + region);
        //this.$store.dispatch("setToEstimateList");
        this.getItemToEstimate();
      }
    },
  },
};
</script>

