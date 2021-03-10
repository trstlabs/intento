<template>
  <div class="pa-2 mx-lg-auto">
    <div>
      <h2 class="headline pa-4 text-center">Browse Items</h2>
      <div>
        <v-container>
          <v-row>
            <v-col cols="12" sm="6">
              <buy-search-bar  />
            </v-col>
            <v-col cols="10" sm="4">
              <buy-options-bar :items="items" />
            </v-col>
           <v-col cols="2" sm="2">
              
            
            <v-btn 
     
      text
      
      small
      
      :outlined="!advanced"
      @click="showAdvanced"
    >
      <v-icon dark>
        mdi-tune
      </v-icon>
    </v-btn>

  </v-col>
          </v-row>
          <buy-tag-bar :advanced="advanced"/>
        </v-container>
        <div v-for="item in items" :key="item.id">
          <div>
            <div >
              <buy-item-item-info :itemid="item.id" />
            </div>
          </div>
        </div>
        <v-card
          class="rounded-lg outlined elevation-1 text-center"
          v-if="items.length < 1"
        >
          <v-card-text class="caption">
            No items to show. Use search / filters to find items or sell an item :)
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
    };
  },
  
  computed: {

   
    items() {
      
      //this.$store.dispatch("setBuyItemList");
      return this.$store.getters.getBuyItemList; console.log(this.$store.getters.getBuyItemList)
    },
  },

  methods: {

    showAdvanced(){
      this.advanced = !this.advanced;

      this.$store.dispatch("setSortedTagList");

    }
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
