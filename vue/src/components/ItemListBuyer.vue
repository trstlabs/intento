<template>
  <div class="pa-2 mx-lg-auto">
   


    <p v-if="buyerItemList.length > 0" class="h2 font-weight-medium text-uppercase text-center">
      Your prepayed items ({{ buyerItemList.length }})<v-btn icon onClick="window.location.reload();"><v-icon >
        mdi-refresh
      </v-icon></v-btn>
    </p>
   

    <div v-for="item in buyerItemList" v-bind:key="item.id">
     
      <div>
        <div>
          <buyer-item-item-info :itemid="item.id" />
        </div>
      </div>
    </div>
    <div v-if="buyerItemList.length === 0 ">
      <p class="caption pa-12 text-center">No items, buy an item first<v-btn icon onClick="window.location.reload();"><v-icon >
        mdi-refresh
      </v-icon></v-btn></p>
    </div>
  </div>
</template>

<script>
import AppButton from "./AppButton.vue";
import BuyerItemItemInfo from "./BuyerItemItemInfo.vue";
export default {
  components: { AppButton, BuyerItemItemInfo },
  data() {
    return {
      dummy: false,
    };
  },

  computed: {
    buyerItemList() {
      return this.$store.getters.getBuyerItemList || [];
    },
  },

 
};
</script> 
