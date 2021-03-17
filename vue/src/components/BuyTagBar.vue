<template>
  <div v-if="advanced">
      <v-divider class="pa-2"/>
   <p class="caption mt-2">Categories: </p>
    <v-chip-group 
    
          
          active-class="primary--text"
        >
          <v-chip @click="updateList(tag)" 
            v-for="tag in tags" :key="tag"
          ><v-icon small left>
        mdi-tag-outline
      </v-icon>{{ tag }}
          </v-chip>
        </v-chip-group>
    <p class="caption">Region: </p>
    <v-select
          append-icon="mdi-tag-outline"
          dense
          v-model="selectedFilter"
          v-on:input="updateLocation(selectedFilter)"
          cache-items
          :items="locations"
          label="Region"
          clearable
          rounded
          solo
     
          hint="Specify region"
        ></v-select>
       <p class="caption">Price minimum: </p>
       <v-slider/>
        <p class="caption">Price maximum: </p>
       <v-slider/>
  <v-divider class="pa-2"/>

  </div>
  
</template>


<script>

import ItemListBuy from "./ItemListBuy.vue";
export default {
  props: ["advanced"],
  components: { ItemListBuy },
  data: function () {
    return {
      selectedFilter: "",
      

  };
  },
  /*mounted() {
    
    //console.log(input)
    
      this.$store.dispatch("setTagList");
  
   },*/

  computed: {
    tags() {
      //console.log("computed tags");
      return this.$store.getters.getTagList },
     //return ["asd","sdaf"] },
locations() {
      return this.$store.getters.getLocationList;
    },
  },
  
    

  methods: {

    updateList(tag) {
      this.$store.dispatch("tagBuyItemList", tag);},

       updateLocation(tag) {
      this.$store.dispatch("locationBuyItemList", tag);},
      
    },



   


   
 
};
</script>
