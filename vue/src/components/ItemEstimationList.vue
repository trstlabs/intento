<template>
  <div>
    <app-text type="h2"> Estimate items </app-text>
    <div v-for="item in items" v-bind:key="item.id">
      <app-text type="h2">Item: {{ item.title }}</app-text>
      <app-text type="p">ItemID: {{ item.id }}</app-text>

      <app-text type="p">Description: {{ item.description }}</app-text>
      <app-text type="p"
        >Shipping cost: {{ item.shippingcost }} tokens</app-text
      >
      <app-text type="p"
        >Local pickup available: {{ item.localpickup }}</app-text
      >
      <app-text type="p">Local pickup available: {{ results }}</app-text>
      <app-input placeholder="Estimation" v-model="estimation" />
      <app-button @click.native="submit(estimation, item.id)"
        >Create Estimation</app-button
      >
      <app-text v-if="results" type="p">Already estimated</app-text>
    </div>
  </div>
</template>

<script>
import AppButton from "./AppButton.vue";
export default {
  components: { AppButton },
  data() {
    return {
      estimation: "",
      estimated: false,
    };
  },
  computed: {
    items() {
      return this.$store.state.data.item || [];
    },
    estimators() {
      return this.$store.state.data.estimator;
    },
    loggedin() {
      if (this.$store.state.account.address != null)
        return this.$store.state.account.address;
    },
  },

  methods: {
    results(loggedin) {
      const results = this.estimators.filter((estimator) => {
        return estimator.creator.includes(loggedin);
      });
    },
    //unestimatedlist() {
    //  this.$store.dispatch("ITEM_LIST_UNESTIMATED", this.input);
    //},

    async submit(estimation, itemid) {
      const type = { type: "estimator" };
      const body = { estimation, itemid };
      await this.$store.dispatch("entitySubmit", { ...type, body });
      await this.$store.dispatch("entityFetch", type);
      alert("submitted");
    },
  },
};
</script>