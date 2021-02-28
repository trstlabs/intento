<template>
  <div>
    <div>
      <h2>Make new Estimation</h2>
    </div>
    <div>
      <app-input placeholder="Estimation" v-model="estimation" />
      <app-input placeholder="ItemID" v-model="itemid" />
    </div>
    <app-button @click.native="submit">Create Estimation</app-button>
  </div>
</template>

<script>
export default {
  data() {
    return {
      estimation: "",
      itemid: "",
    };
  },
  methods: {
    async submit() {
      const payload = {
        type: "estimator",
        body: {
          estimation: this.estimation,
          itemid: this.itemid,
        },
      };
      await this.$store.dispatch("entitySubmit", payload);
      await this.$store.dispatch("entityFetch", payload);
      await this.$store.dispatch("accountUpdate");
      alert("submitted");
    },
  },
};
</script>