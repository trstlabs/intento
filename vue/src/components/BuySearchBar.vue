<template>
  <div > 
    <!--<v-card outlined elevation="2" class="rounded-lg" :class="($vuetify.theme.dark) ? 'grey darken-3 ' : 'white'">
      <v-btn icon class="hidden-xs-only">
        <v-icon>mdi-magnify</v-icon>
      </v-btn>
      <input 
      
        type="text"
        placeholder="Search..."
        v-model.trim="input"
        v-on:input="search()"
        ref="input"
        @blur="closeSearchBar()"
        style="width: 80%"
      />

      <v-spacer></v-spacer>

      
    </v-card> --->
<!-- <v-btn icon class="hidden-xs-only"> 
        <v-icon>mdi-magnify</v-icon>
      </v-btn>-->
      
    <v-text-field  solo clearable
    prepend-inner-icon="mdi-magnify"
     class="rounded-lg" type="text"
        placeholder="Search..."
        v-model.trim="input"
        v-on:input="search()"
        ref="input"
      @click:clear="search()"
       > 

    </v-text-field>
  </div>
</template>

<script>
export default {
  // name: " listSearch",
  data: () => ({
    input: "",
  }),
  methods: {
    search() {
      this.$store.dispatch("updateBuyItemList", this.input);
      //this.$store.commit("SET_LIST_SEARCH", this.input);
 
   
      //this.$store.commit("SET_LIST_SEARCH", this.input);
    },
    //closeSearchBar() {
      //this.$store.commit("SET_DISPLAY_SEARCH_LIST", false);
    //},
  },
};
</script>

<style scoped>
input {
  border: none;
  font-size: inherit;
  padding: 1rem 1rem;
  margin-bottom: 0rem;
  width: 100%;

  font-family: inherit;
}
input:focus {
  outline: none;
}
input:disabled {
  color: rgba(0, 0, 0, 0.25);
}
</style>