<template>
  <div class="pa-2 mx-lg-auto">

     <v-tooltip right >
      <template v-slot:activator="{ on, attrs }">
        <span
          v-bind="attrs"
          v-on="on"> 
    <p  v-if="creatorItemList.length > 0" class="pa-2 h2 font-weight-medium text-uppercase text-center">
      Your items ({{ creatorItemList.length }}), Actionable ({{creatorActionList.length}})
    </p>

    </span>
    </template>  <span > You have created these items, check out the actions to be done. </span> 
      </v-tooltip>

    <v-btn text
      v-if="creatorItemList.length < 1"
      @click="getItemsFromCreator"
    >
      Display items
    </v-btn>

    <div v-for="item in creatorItemList" v-bind:key="item.id">
      <v-lazy
        v-model="isActive"
        :options="{
          threshold: .5
        }"
       
        transition="fade-transition"
      >
      <creator-item-item-info :itemid="item.id" />
      </v-lazy>
    </div>
    <div class="card__empty" v-if="creatorItemList.length === 0 && dummy">
      <p>No items, place an item first</p>
    </div>
  </div>
</template>

<script>
import AppButton from "./AppButton.vue";
import CreatorItemItemInfo from "./CreatorItemItemInfo.vue";
export default {
  components: { AppButton, CreatorItemItemInfo },
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

  methods: {
    getItemsFromCreator() {
      if (this.$store.state.client == null) { alert("Sign in first");};
      this.dummy = true;
      let input = this.$store.state.account.address;
      this.$store.dispatch("setCreatorItemList", input);
      //this.dummy = false;
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
