<template>
  <div class="pa-2 mx-lg-auto">

    
    <p  v-if="creatorItemList.length > 0" class="h2 font-weight-medium text-uppercase text-center">
      Total ({{ creatorItemList.length }}), Actionable ({{creatorActionList.length}})<v-btn icon onClick="window.location.reload();" ><v-icon >
        mdi-refresh
      </v-icon></v-btn>
    </p>

 
    

    

    <div v-for="item in creatorItemList" v-bind:key="item.id">
      <v-sheet min-height="250" class="fill-height" color="transparent">
      <v-lazy
        v-model="isActive"
        :options="{
          threshold: .5
        }"
       
        transition="fade-transition"
      >
      <creator-item-item-info :itemid="item.id" />
      </v-lazy> </v-sheet>
    </div>
    <div v-if="creatorItemList.length === 0">
      <p class="caption pa-12 text-center">No items, place an item first<v-btn icon  onClick="window.location.reload();" ><v-icon >
        mdi-refresh
      </v-icon></v-btn></p>
    </div>
    
  </div>
</template>

<script>

import CreatorItemItemInfo from "./CreatorItemItemInfo.vue";
export default {
  components: { CreatorItemItemInfo },
  data() {
    return {
      dummy: false,
      isActive: false, 
    };
  },
 

  computed: {
    creatorItemList() {
      return this.$store.getters.getCreatorItemList || [];
    },
     creatorActionList() {
      return this.$store.getters.getCreatorActionList || [];
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
