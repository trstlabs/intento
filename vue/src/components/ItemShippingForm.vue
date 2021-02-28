<template>
  <div>
    <h2>Complete shipping to buyer</h2>
    <div>
      <input
        type="checkbox"
        id="checkbox"
        v-model="tracking"
        v-bind:value="true"
      />
      <label for="checkbox">I have provided the tracking nr to buyer </label>
      <app-input placeholder="item id" v-model="itemid" />
    </div>
    <app-button @click.native="submit">Accept</app-button>
  </div>
</template>

<script>
export default {
  data() {
    return {
      tracking: false,
      itemid: "",
    };
  },
  methods: {
    async submit() {
      const payload = {
        type: "item/shipping",
        body: {
          tracking: this.tracking,
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