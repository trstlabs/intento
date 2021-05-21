<template>
  <div class="pa-0 pb-4 mx-lg-auto">
    <div>
        
  <p    class="display-2 mt-n9 font-weight-thin text-center  "> Marketplace</p><p class="display-1 text-center font-weight-thin">{{this.region.toUpperCase()}}<v-icon x-large left> mdi-flag-variant </v-icon></p>

      <div>  <v-img class="rounded-xl" h height="320" src="img/design/market.png "  >  </v-img>  
        <v-container class="mt-n12">
      
          <v-row>  
            <v-col cols="12" sm="7">
              <buy-search-bar />
            </v-col>
            <v-col cols="10" sm="4">
              <buy-options-bar :items="items" />
            </v-col>
            <v-col cols="2" sm="1" class="ml-n1">
              <v-btn fab class="mt-1" small  @click="showAdvanced">
                <v-icon> mdi-tune </v-icon>
              </v-btn>
            </v-col>   
          </v-row> 
          <buy-tag-bar :advanced="advanced" />
        </v-container>
     
        <div v-if="items[0]">
      <div class="pl-4" >
             <v-btn icon x-small to="/faq"
                ><v-icon  >mdi-information-outline</v-icon>
              </v-btn> <span class="caption" v-if="items[1]"> {{items.length}}  items available in {{region.toUpperCase()}}</span></div>
          <div v-for="item in items" :key="item.id">
            <div>
           <div class="pa-4">
                <v-sheet
                  
                  class="fill-height"
                  color="transparent"
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
       <v-card @click="clearList()" color="secondary lighten-3 "
          class="rounded-xl outlined elevation-1 text-center ma-4 pa-2"
          v-if="items.length < 1"
        >
          <v-card-text class="caption ">
            No items to show. Use search / filters to find items.
          </v-card-text><v-icon>
        mdi-refresh
      </v-icon>
        </v-card>
      </div>
    </div>
  </div>
</template>


<script>
import BuyItemItemInfo from "./BuyItemItemInfo.vue";
import BuySearchBar from "./BuySearchBar.vue";
import BuyOptionsBar from "./BuyOptionsBar.vue";
import BuyTagBar from "./BuyTagBar.vue";
import SearchRegion from "../views/SearchRegion.vue";
export default {
  components: {
    BuyItemItemInfo,
    BuySearchBar,
    BuyOptionsBar,
    BuyTagBar,
    SearchRegion,
  },
  props: ["region"],
  data() {
    return {
      advanced: false,
      isActive: false,
    };
  },

   

  computed: {
    itemsUnfiltered() {
      //this.$store.dispatch("setBuyItemList");
      return this.$store.getters.getBuyItemList;
    },

     items() {

       return this.itemsUnfiltered.filter(item => item.shippingregion.find(region => region.includes(this.region)))
       
    },
  },

  methods: {
    showAdvanced() {
      this.advanced = !this.advanced;

      this.$store.dispatch("setSortedTagList");
      this.$store.dispatch("setSortedLocationList");

    },
     clearList() {
     this.$store.dispatch("setBuyItemList");
  }
  },
};
</script>


