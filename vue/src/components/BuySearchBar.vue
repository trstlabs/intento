<template>
  <div>
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

    <v-text-field
      solo
      clearable
      prepend-inner-icon="mdi-magnify"
      class="rounded-xl"
      type="text"
      placeholder="Search..."
      v-model.trim="input"
      v-on:input="search()"
      ref="input"
       @click:clear="setList()"
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
      this.$store.dispatch("updateBuyItemList", this.input.toLowerCase());
    
    },
   setList() {
      this.$store.dispatch("setBuyItemList");
    },
  },
};
</script>