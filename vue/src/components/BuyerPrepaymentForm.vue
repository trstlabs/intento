<template>
  <div>
    <h2>Buy Item</h2>
    <div>
      <app-input type="number" placeholder="Deposit Coins" v-model="deposit" />
      <app-input placeholder="item id" v-model="itemid" />
    </div>
    <app-button @click.native="submit">Provide Prepayment</app-button>
  </div>
</template>

<script>
export default {
  data() {
    return {
      deposit: "",
      itemid: "",
    };
  },
  methods: {
    async submit() {
      const payload = {
        type: "buyer",
        body: {
          deposit: this.deposit,
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