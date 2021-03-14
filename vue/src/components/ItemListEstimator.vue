<template>
  <div class="pa-2 mx-lg-auto">
    

   
    <p v-if="estimatorItemList.length > 0" class="h2 font-weight-medium text-uppercase text-center">
      Your Estimated items ({{ estimatorItemList.length }})
    </p>
  
   

    <div
      v-for="estimator in estimatorItemList"
      v-bind:key="estimator.itemid"
      
    > <v-lazy
        v-model="isActive"
        :options="{
          threshold: .5
        }"
       
        transition="fade-transition"
      >
      <estimator-item-item-info :itemid="estimator.itemid" /> </v-lazy>
    </div>
    <div v-if="estimatorItemList.length === 0">
      <p class="caption pa-12 text-center">No estimations, make an estimation first</p>
    </div>
  </div>
</template>

<script>
import AppButton from "./AppButton.vue";
import EstimatorItemItemInfo from "./EstimatorItemItemInfo.vue";
export default {
  components: { AppButton, EstimatorItemItemInfo },
  data() {
    return {
      dummy: false,
    };
  },
  computed: {
    estimatorItemList() {
      return this.$store.getters.getEstimatorItemList || [];
    },
  },

 
};
</script>
