<template>
  <div class="pa-0 pb-4 mx-lg-auto">
    <div>
      <p    class="display-2 pt-4 font-weight-thin text-center  pb-5 "> Marketplace</p>
   <div v-if="items[1]">   <!--<v-progress-linear
      :value="onSaleRatio"
 background-color="secondary"
      height="20" class="pb-5"
    ></v-progress-linear>--><v-row class="overline mx-2 pa-2 text-left font-weight-thin">{{items.length}} Items available<v-spacer/>{{totalItems}} items on TPP </v-row></div>
      <div>
        <v-img height="320" src="img/design/market.png "> </v-img>
        <v-container class="mt-n12">
          <v-row>
            <v-col cols="12" sm="7">
              <buy-search-bar />
            </v-col>
            <v-col cols="10" sm="4">
              <buy-options-bar :items="items" />
            </v-col>
            <v-col cols="2" sm="1" class="ml-n1">
              <v-btn fab class="mt-1" small @click="showAdvanced">
                <v-icon> mdi-tune </v-icon>
              </v-btn>
            </v-col>
          </v-row>
          <buy-tag-bar :advanced="advanced" />
        </v-container>

        <div v-if="items[0]">
          <div class="pl-4">
            <v-btn icon x-small to="/faq"
              ><v-icon>mdi-information-outline</v-icon>
            </v-btn>
            <span class="caption" v-if="items[1]">
              {{ items.length }} items available</span
            >
          </div>
          <div v-for="item in items" :key="item.id">
            <div>
              <div>
                <v-sheet class="fill-height" color="transparent"
                  ><v-lazy
                    v-model="isActive"
                    :options="{
                      threshold: 0.5,
                    }"
                    transition="fade-transition"
                  >
                    <buy-item-item-info :itemid="item.id" /></v-lazy
                ></v-sheet>
              </div>
            </div>
          </div>
        </div>
           <v-card @click="clearList()" color="secondary lighten-3"
          class="rounded-lg outlined elevation-1 text-center"
          v-if="items.length < 1"
        ><v-icon>
        mdi-refresh
      </v-icon>
          <v-card-text class="caption">
            No watches to show. Use search / filters to find watches.
          </v-card-text></v-card>
      </div>
    </div>
  </div>
</template>


<script>
import BuyItemItemInfo from "./BuyItemItemInfo.vue";
import BuySearchBar from "./BuySearchBar.vue";
import BuyOptionsBar from "./BuyOptionsBar.vue";
import BuyTagBar from "./BuyTagBar.vue";
export default {
  components: {
    BuyItemItemInfo,
    BuySearchBar,
    BuyOptionsBar,
    BuyTagBar,
  },
  data: function () {
    return {
      advanced: false,
      isActive: false,
    };
  },

  computed: {
    items() {
      //this.$store.dispatch("setBuyItemList");
      return this.$store.getters.getBuyItemList;
    },
    totalItems() {   
       if(this.items){
      return this.$store.state.data.item.length;
       }
    },
  
    onSaleRatio() {
      if(this.items){
      //this.$store.dispatch("setBuyItemList");
      return 100 - (this.totalItems - this.items.length)/this.totalItems*100}else{return 0}
    },
  },

  methods: {
    showAdvanced() {
      this.advanced = !this.advanced;

      this.$store.dispatch("setSortedTagList");
      this.$store.dispatch("setSortedLocationList");
    },
    clearList() {
    this.$router.go()
  }

  },
};
</script>


