<template>
  <div class="pa-0 pb-4 mx-lg-auto">
    <div>
        <p   class="headline pt-4 font-weight-bold text-center ">Shop All </p>

      <div>  <v-img  height="320" src="img/design/market.png " >  </v-img>  
        <v-container class="pt-0">
      
          <v-row>  
            <v-col cols="12" sm="6">
              <buy-search-bar />
            </v-col>
            <v-col cols="10" sm="4">
              <buy-options-bar :items="items" />
            </v-col>
            <v-col cols="2" sm="2">
              <v-btn text small fab @click="showAdvanced">
                <v-icon dark> mdi-tune </v-icon>
              </v-btn>
            </v-col>   
          </v-row> 
          <buy-tag-bar :advanced="advanced" />
        </v-container>
     
        <div v-if="items[0]">
      
             <v-btn icon x-small to="/faq"
                ><v-icon >mdi-information-outline</v-icon>
              </v-btn> <span class="caption" v-if="items[1]"> {{items.length}} items available</span>
          <div v-for="item in items" :key="item.id">
            <div>
              <div>
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
        <v-card
          class="rounded-lg outlined elevation-1 text-center"
          v-if="items.length < 1"
        >
          <v-card-text class="caption">
            No items to show. Use search / filters to find items or sell an item
            :)
          </v-card-text>
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
  },

  methods: {
    showAdvanced() {
      this.advanced = !this.advanced;

      this.$store.dispatch("setSortedTagList");
      this.$store.dispatch("setSortedLocationList");

    },
  },
};
</script>

<style scoped>
.card__empty {
  margin-bottom: 1rem;
  border: 1px dashed rgba(0, 0, 0, 0.1);
  padding: 1rem;
  display: flex;
  align-items: center;
  justify-content: center;
  box-sizing: border-box;
  border-radius: 8px;
  color: rgba(0, 0, 0, 0.25);
  text-align: center;
  min-height: 8rem;
}
</style>
