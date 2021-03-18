<template>
  <div v-if="advanced">
      <v-divider class="ma-4"/>
    <v-card class="pa-2">
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
          append-icon="mdi-earth"
          dense
          v-model="selectedFilter"
          v-on:input="updateLocation(selectedFilter)"
          cache-items
          :items="locations"
          label="Region"
          clearable
          
          outlined
     
          hint="Specify region"
        ></v-select>
         <v-row class="ma-2"> <v-col>
       <p class="caption">Price minimum: </p>
     
        <v-text-field
              label="Amount"
              type="number"
              v-model="minPrice"
           

              prefix="$"
              suffix="TPP"
            ></v-text-field> </v-col> <v-col> <v-btn icon @click="updatePriceMin"><v-icon>mdi-check</v-icon> </v-btn> <v-btn icon @click="clearList"><v-icon>mdi-cancel</v-icon> </v-btn> </v-col>
            <v-col>
       <p class="caption">Price maximum: </p>
     
        <v-text-field
              label="Amount"
              type="number"
              v-model="maxPrice"
           

              prefix="$"
              suffix="TPP"
            ></v-text-field> </v-col> <v-col> <v-btn icon @click="updatePriceMax"><v-icon>mdi-check</v-icon> </v-btn> <v-btn icon @click="clearList"><v-icon>mdi-cancel</v-icon> </v-btn> </v-col></v-row>
    </v-card>
 
 <v-divider class="ma-4"/>
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
      minPrice: 0,
      maxPrice: 0,
      

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
      
updatePriceMin() {
   //this.$store.dispatch("setBuyItemList");
      this.$store.dispatch("priceMinBuyItemList", this.minPrice);

      },

      updatePriceMax() {
   //this.$store.dispatch("setBuyItemList");
      this.$store.dispatch("priceMaxBuyItemList", this.maxPrice);

      },
       clearList() {
     
     this.$store.dispatch("setBuyItemList");},


    },

    

   


   
 
};
</script>
