<template>
  <div class="pa-2 mx-lg-auto">
    <div>
      <h2 class="display-1 pa-4 text-center">Browse Items</h2>
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
          <v-card-text>
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
button {
  background: none;
  border: none;
  color: rgba(0, 125, 255);
  padding: 0;
  font-size: inherit;
  font-weight: 800;
  font-family: inherit;
  text-transform: uppercase;
  margin-top: 0.5rem;
  cursor: pointer;
  transition: opacity 0.1s;
  letter-spacing: 0.03em;
  transition: color 0.25s;
  display: inline-flex;
  align-items: center;
}
.item {
  box-shadow: inset 0 0 0 1px rgba(0, 0, 0, 0.1);
  margin-bottom: 1rem;
  padding: 1rem;
  border-radius: 0.5rem;
  overflow: hidden;
}
.item__field {
  display: grid;
  line-height: 1.5;
  grid-template-columns: 15% 1fr;
  grid-template-rows: 1fr;
  word-break: break-all;
}
.item__field__key {
  color: rgba(0, 0, 0, 0.25);
  word-break: keep-all;
  overflow: hidden;
}
button:focus {
  opacity: 0.85;
  outline: none;
}
.button.button__valid__true:active {
  opacity: 0.65;
}
.button__label {
  display: inline-flex;
  align-items: center;
}
.button__label__icon {
  height: 1em;
  width: 1em;
  margin: 0 0.5em 0 0.5em;
  fill: rgba(0, 0, 0, 0.25);
  animation: rotate linear 4s infinite;
}
.button.button__valid__false {
  color: rgba(0, 0, 0, 0.25);
  cursor: not-allowed;
}
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
@keyframes rotate {
  from {
    transform: rotate(0);
  }
  to {
    transform: rotate(-360deg);
  }
}
@media screen and (max-width: 980px) {
  .narrow {
    padding: 0;
  }
}
</style>
