<template>
  <div class="pa-2 mx-lg-auto">
    

    <v-tooltip right >
      <template v-slot:activator="{ on, attrs }">
        <span
          v-bind="attrs"
          v-on="on"> 
    <p v-if="estimatorItemList.length > 0" class="h2 font-weight-medium text-uppercase text-center">
      Your Estimated items ({{ estimatorItemList.length }})
    </p>
    </span>
    </template>  <span > Items you estimated. </span> 
      </v-tooltip>

    <v-btn text
      v-if="estimatorItemList.length < 1"
      @click="getItemsFromEstimator"
    >
      Display estimated items
    </v-btn>

    

    <div
      v-for="estimator in estimatorItemList"
      v-bind:key="estimator.itemid"
      
    >
      <estimator-item-item-info :itemid="estimator.itemid" />
    </div>
    <div class="card__empty" v-if="estimatorItemList.length === 0 && dummy">
      <p>No estimations, make an estimation first</p>
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

  methods: {
    getItemsFromEstimator() {
      this.dummy = true;
      let input = this.$store.state.account.address;

      this.$store.dispatch("setEstimatorItemList", input);
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

