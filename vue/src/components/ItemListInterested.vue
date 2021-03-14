<template>
  <div class="pa-2 mx-lg-auto">
   
    <p v-if="buyerItemList.length > 0" class="h2 font-weight-medium text-uppercase text-center">
      Liked Items({{ buyerItemList.length }})
    </p>
   

    <div v-for="item in buyerItemList" v-bind:key="item.id">
     
      <div>
        <div>
          <interested-item-item-info :itemid="item.itemid" />
        </div>
      </div>
    </div>
    <div  v-if="buyerItemList.length === 0">
      
      <p class="caption pa-12 text-center">No items, estimate items to find items you are interested in before they are on sale</p> 
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

  
};
</script> 
