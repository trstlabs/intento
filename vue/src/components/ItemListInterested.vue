<template>
  <div class="pa-2 mx-lg-auto">
    <v-tooltip right >
      <template v-slot:activator="{ on, attrs }">
        <span
          v-bind="attrs"
          v-on="on"> 
    <p v-if="buyerItemList.length > 0" class="h2 font-weight-medium text-uppercase text-center">
      Liked Items({{ buyerItemList.length }})
    </p>
    </span>
    </template>  <span > Items you liked and estimated and are now for sale. </span> 
      </v-tooltip>


    <v-btn text
      v-if="buyerItemList.length < 1"
      @click="getInterestedItems"
    >
      Liked items
    </v-btn>
     

    <div v-for="item in buyerItemList" v-bind:key="item.id">
     
      <div>
        <div>
          <interested-item-item-info :itemid="item.itemid" />
        </div>
      </div>
    </div>
    <div class="card__empty" v-if="buyerItemList.length === 0 && dummy">
      
      <p class="caption">No items, estimate items to find items you are interested in before they are on sale</p> 
    </div>
  </div>
</template>

<script>
import AppButton from "./AppButton.vue";
import InterestedItemItemInfo from "./InterestedItemItemInfo.vue";
export default {
  components: { AppButton, InterestedItemItemInfo },
  data() {
    return {
      dummy: false,
    };
  },

  computed: {
    buyerItemList() {
      return this.$store.getters.getInterestedItemList || []; 
      
    },
  },

  methods: {
    getInterestedItems() {
      this.dummy = true;
      let input = this.$store.state.account.address
      this.$store.dispatch("setInterestedItemList", input);
      console.log(this.buyerItemList)
      
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