<template>
  <div class="pa-2 mx-lg-auto">
   
    <p v-if="buyerItemList.length > 0" class="h2 font-weight-medium text-uppercase text-center">
      Liked Items({{ buyerItemList.length }})<v-btn icon onClick="window.location.reload();" ><v-icon >
        mdi-refresh
      </v-icon></v-btn>
    </p>
   

    <div v-for="item in buyerItemList" v-bind:key="item.id">
     
      <div>
        <div>
          <interested-item-item-info :itemid="item.itemid" />
        </div>
      </div>
    </div>
    <div  v-if="buyerItemList.length === 0">
      
      <p class="caption pa-12 text-center">No items, like and estimate liked items to find items before anyone else <v-btn icon onClick="window.location.reload();" ><v-icon >
        mdi-refresh
      </v-icon></v-btn></p> 
    </div><v-img src="img/design/like.png" ></v-img>
  </div>
</template>

<script>

import InterestedItemItemInfo from "./InterestedItemItemInfo.vue";
export default {
  components: { InterestedItemItemInfo },
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
