<template >
  <div class="pa-2 mx-lg-auto">
    <div class="text-center pa-12" v-if="!showinfo">
      <v-btn :ripple="false" text @click="getItemToEstimate"
        ><v-icon color="primary" left> mdi-refresh </v-icon> Refresh
      </v-btn>
    </div>
    <v-skeleton-loader
      v-if="loadingitem"
      class="mx-auto"
      type="list-item-three-line, image, article"
    ></v-skeleton-loader>

<div v-if="showinfo == true && loadingitem == false" >
    <div elevation="8" v-if="photos.photo">
      <v-carousel
        delimiter-icon="mdi-minus"
        carousel-controls-bg="primary"
        contain
        hide-delimiter-background
        show-arrows-on-hover
      >
        <v-carousel-item v-for="(photo, i) in photos" :key="i" :src="photo">
        </v-carousel-item>
      </v-carousel>
    </div>

    <v-card 
      class="pa-2 mt-2"
      elevation="2"
      rounded="lg"
     
    >
      <v-progress-linear
        indeterminate
        :active="loadingitem"
      ></v-progress-linear>

      <div class="pa-2">
        <v-row>
          <v-col class="pa-2">
            <v-card elevation="0">
              <div class="overline">Title</div>

              <div class="body-1">
                {{ item.title }}
              </div>
            </v-card>
          </v-col>
          <v-col class="pa-2">
            <v-card elevation="0">
              <v-chip-group>
                <v-chip
                  outlined
                  small
                  v-for="itemtag in item.tags"
                  :key="itemtag"
                >
                  <v-icon small left> mdi-tag-outline </v-icon>
                  {{ itemtag }}
                </v-chip>
              </v-chip-group>

              <v-dialog transition="dialog-bottom-transition" max-width="300">
                <template v-slot:activator="{ on, attrs }">
                  <span v-bind="attrs" v-on="on">
                    <v-chip class="ma-1" outlined small>
                      <v-icon left small> mdi-star-outline </v-icon>
                      {{ item.condition }}/5
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
            </v-card>
          </v-col>
        </v-row>
        <v-card elevation="0">
          <div class="overline">Description</div>

          <div class="caption">
            {{ item.description }}
          </div>
        </v-card>
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
          <v-col cols="4" class="mx-auto">
            <v-dialog transition="dialog-bottom-transition" max-width="600">
              <template v-slot:activator="{ on, attrs }">
                <span v-bind="attrs" v-on="on"
                
                ><v-icon small>mdi-information-outline</v-icon>
           
                  <p class="text-center caption ">To me, it's worth</p>
                </span>
              </template>
              <template v-slot:default="dialog">
                <v-card>
                  <v-toolbar color="default">Info</v-toolbar>
                  <v-card-text>
                    <div class="text-p pt-4">
                      Earn ~5% of the item value when you are the best
                      estimator. Exept when:
                    </div>
                    <div class="caption pa-2">
                      - Your deposit is lost when you are the lowest estimator
                      and the final estimation price is not accepted by the
                      seller.
                    </div>
                    <div class="caption pa-2">
                      - Your deposit is lost when you are the highest estimator
                      and the item is not bought by the buyer that provided
                      prepayment.
                    </div>
                     <div class="text-p pt-4">
                      Good luck and have fun!
                    </div>
                  </v-card-text>
                  <v-card-actions class="justify-end">
                    <v-btn text @click="dialog.value = false">Close</v-btn>
                  </v-card-actions>
                </v-card>
              </template>
            </v-dialog>
          </v-col>
          <v-col cols="8" class="mx-auto">
            <v-text-field
              label="Amount"
              type="number"
              v-model="estimation"
              :disabled="lastitem"
              prefix="$"
              suffix="TPP"
            ></v-text-field>
          </v-col>
        </v-row>
      </div>
      <v-divider class="mx-4" />
      <v-card elevation="0">
        <div class="pa-2">
          <div>
            <v-chip-group active-class="primary--text" column>
              <v-chip
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
      </v-card>

      <div>
        <v-btn
          block
          elevation="4"
          color="primary"
          :disabled="!valid || !hasAddress || flight || timeout"
          @click="submit(estimation, item.id, interested, comment)"
          ><div v-if="!flight && !valid">
            <v-icon left> mdi-check </v-icon> Estimate item
          </div>
          <div v-if="!flight && valid && !timeout">
            <v-icon left> mdi-check-all </v-icon> Estimate item 
          </div>
          <div v-if="!timeout && !valid && flight">
            <v-icon left> mdi-check </v-icon> Estimate item
          </div>
         
             <v-progress-linear v-if="timeout && valid"
      :rotate="360"
      :size="50"
      :width="5"
      :value="value"
      color="white"
    >
      
    </v-progress-linear>
             <div >
          
         
            <div v-if="flight">
              <div class="text-right">Creating estimation...</div>
            </div>
          </div>
        </v-btn>
        <div v-if="timeout && valid">
              <div class="text-right caption">Wait {{ 10-(value/10) }}s</div>
              
            </div>

             <div v-if="!flight && valid && !timeout">
              <div class="text-right caption">Required deposit is {{item.depositamount}}TPP</div>
              
            </div>
        <!-- tag bar
 <v-chip-group 
    
          
          active-class="primary--text"
        >
          <v-chip @click="updateList(tag)"  outlined
            v-for="tag in tags" :key="tag"
          ><v-icon small left>
        mdi-tag-outline
      </v-icon>{{ tag }}
          </v-chip>
        </v-chip-group>-->
      </div>
    </v-card>
</div>
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
          <!-- <v-tooltip bottom :disabled="flag" v-if="showinfo">
      <template v-slot:activator="{ on, attrs }">
        <span
          v-bind="attrs"
          v-on="on"> 
           <v-btn
      class="mx-2"
      fab
      dark
      small
      color="red"
      :outlined="flag == false"
      @click="submitFlag(true, item.id)"
    >
      <v-icon dark>
        mdi-alert-octagon
      </v-icon>
    </v-btn>
          </span>
    </template>  <span > When this item is not OK, report it. Thank You.</span> 
      </v-tooltip> -->

          <v-dialog
            bottom
            :disabled="flag"
            v-if="showinfo"
            v-model="dialog"
            persistent
            max-width="290"
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
              <v-card-title class="headline"> Report this item? </v-card-title>
              <v-card-text
                >If this item is not OK, you can report it here. TPP protocol
                will automatically remove items that are reported
                often.</v-card-text
              >
              <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn color="red darken-1" text @click="dialog = false">
                  Close
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

      <div class="pt-12 mx-lg-auto">
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
      </div>
    </div>
  </div>
</template>

<script>
import ToEstimateTagBar from "./ToEstimateTagBar.vue";
import { databaseRef } from "./firebase/db.js";

import {
  SigningStargateClient,
  assertIsBroadcastTxSuccess,
} from "@cosmjs/stargate";
import { Registry } from "@cosmjs/proto-signing/";
import { Type, Field } from "protobufjs";

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
        { name: "Looks damaged", attr: "The item appears to be damaged." },
        {
          name: "Repairable",
          attr: "The item seems damaged, but I think it can be repaired.",
        },
        { name: "Used", attr: "The item seems used." },
        {
          name: "As good as new!",
          attr: "The item appears to look as good as new!",
        },
        { name: "Dirty", attr: "The item looks dirty to me." },
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
      timeout: false,
interval: {},
value: 0,
      dialog: false,
      conditionInfo: false,
    };
  },
  beforeDestroy () {
      clearInterval(this.interval)
    },

  mounted() {

    
      
    //console.log(input)
    if (!!this.$store.state.account.address) {
      let input = this.$store.state.account.address;

       const type = { type: "estimator" };
      this.$store.dispatch("entityFetch",type);
      this.$store.dispatch("setSortedTagList");
      this.$store.dispatch("setEstimatorItemList", input);
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
  },

  methods: {
    async submit(estimation, itemid, interested, comment) {
      if (this.valid && !this.flight && this.hasAddress) {
        this.flight = true;
        this.loadingitem = true;
        const type = { type: "estimator" };
        const body = {
          deposit: this.item.depositamount,
          estimation: estimation,
          itemid: itemid,
          interested: interested,
          comment: comment,
        };
        const fields = [
          ["estimator", 1, "string", "optional"],
          ["estimation", 2, "int64", "optional"],
          ["itemid", 3, "string", "optional"],
          ["deposit", 4, "int64", "optional"],
          ["interested", 5, "bool", "optional"],
          ["comment", 6, "string", "optional"],
        ];
        
        await this.estimationSubmit({ ...type, body, fields });
        await this.$store.dispatch("entityFetch", type);
       
             await this.$store.dispatch("bankBalancesGet");
        this.timeout = true
        clearInterval(this.interval)
        this.value = 0
        this.interval = setInterval(() => {
      
        this.value += 10
      }, 1000)


         setTimeout(() => this.timeout = false, 10000);
     

   
        await this.submitRevealEstimation(itemid);
         this.estimation = "";
        this.comment = "";
        //this.flight = false;
        //this.loadingitem = false;
      }
    },

    async estimationSubmit({ type, fields, body }) {
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
      if (!result.data){
        await this.submitRevealEstimation(this.item.id);
      }
      assertIsBroadcastTxSuccess(result);
      
      console.log("success!");

      
    },

    async submitFlag() {
      if (!this.flight && this.hasAddress) {
        this.flight = true;
        this.loadingitem = true;
        this.flag = true;
        const type = { type: "estimator" };
        const body = { flag: true, itemid: this.item.id };
        const fields = [
          ["estimator", 1, "string", "optional"],
          ["itemid", 2, "string", "optional"],
          ["flag", 3, "bool", "optional"],
        ];

        await this.flagSubmit({ ...type, body, fields });

        this.estimation = "";

        this.getNewItemByIndex();
        this.dialog = false;
        this.flag = false;
      }
    },

    async flagSubmit({ type, fields, body }) {
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

    async getItemToEstimate() {
      if (!this.hasAddress) {
        alert("Sign in first");
        window.location.reload()
        //return (this.showinfo = false);
      }
      let input = this.$store.state.account.address;
      
      this.$store.dispatch("setSortedTagList");
      this.$store.dispatch("setEstimatorItemList", input);
      this.$store.dispatch("setToEstimateList", this.index);

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
    loadItemPhotos() {
      this.loadingitem = true;
      const id = this.item.id;
      //const db = firebase.database();

     const imageRef = databaseRef.ref("ItemPhotoGallery/" + id + "/photos/");
      imageRef.on("value", (snapshot) => {
        const data = snapshot.val();

        if (data != null && data.photo != null) {
        
           console.log(data.photo);
          this.photos = data;
          //this.photos = { photo: "https://upload.wikimedia.org/wikipedia/commons/thumb/8/85/Points_of_a_horse.jpg/330px-Points_of_a_horse.jpg" };

          this.loadingitem = false;
        }else{
          this.photos = []
          this.getNewItemByIndex()
        }
      });
      //this.loadingitem = false;
      this.interested = false;
      this.flight = false;
    },

    async submitRevealEstimation(itemid) {
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

    updateComment(newComment) {
      this.comment = newComment;
    },
    updateList(tag) {
      //console.log(this.tag);
      this.$store.dispatch("tagToEstimateList", tag);
      if (!!this.items[0]) {
        this.item = this.items[0];
        if (!this.items[1]) {
          this.lastitem = true;
        } else {
          this.lastitem = false;
        }
        this.loadItemPhotos();
      } else {
        alert("No Items to estimate for:" + tag);
        this.$store.dispatch("setToEstimateList");
        this.getItemToEstimate();
      }
    },
  },
};
</script>

