<template>
  <div v-if="advanced">
      <div class="mx-4"  v-if="applied[0]" > <p class="caption mb-0 mt-2">Applied filters: </p> <v-chip-group 
    column
        ><v-icon @click="clearList()" small left>
        mdi-close
      </v-icon><div  v-for="filter in applied" :key="filter">
          <v-chip 
           v-if="filter "
          ><v-icon small left>
        mdi-tune
      </v-icon>{{ filter }}
          </v-chip></div>
        </v-chip-group></div>  <v-divider class="ma-4"/>    
    <v-card color="secondary lighten-3" class="pa-2 elevation-5 rounded-lg">

   
   <p class="caption mb-2">Categories: </p>
    <v-chip-group 
    multiple
           show-arrows
          active-class="primary--text"
        ><div  v-for="tag in tags" :key="tag">
          <v-chip v-if="!applied.includes(tag)" color="primary lighten-1 text-capitalize"  @click="updateList(tag)" 
           
           ><v-icon small left>
        mdi-tag-outline
      </v-icon>{{ tag}}
          </v-chip></div>
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
           

        append-icon="$vuetify.icons.custom"
             
            ></v-text-field> </v-col> <v-col> <v-btn icon @click="updatePriceMin"><v-icon>mdi-check</v-icon> </v-btn> <v-btn icon @click="clearList"><v-icon>mdi-close</v-icon> </v-btn> </v-col>
            <v-col>
       <p class="caption">Price maximum: </p>
     
        <v-text-field
              label="Amount"
              type="number"
              v-model="maxPrice"
           

               append-icon="$vuetify.icons.custom"
            ></v-text-field> </v-col> <v-col> <v-btn icon @click="updatePriceMax"><v-icon>mdi-check</v-icon> </v-btn> <v-btn icon @click="clearList"><v-icon>mdi-close</v-icon> </v-btn> </v-col></v-row>
   
    <v-text-field  solo clearable
    prepend-inner-icon="mdi-magnify"
     class="rounded-lg" type="text"
        placeholder="Search title and description..."
        v-model.trim="input"
        v-on:input="search()"
        ref="input"
      background-color="secondary lighten-3" 
       > 

    </v-text-field>
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
      input: '',
      applied: [],
      

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
      
      this.applied.push(tag);
      this.$store.dispatch("tagBuyItemList", tag);},

       updateLocation(tag) {
        
        this.applied.push(tag);
      this.$store.dispatch("locationBuyItemList", tag);},
      
updatePriceMin() {
   //this.$store.dispatch("setBuyItemList");
   this.applied.push(this.minPrice);
      this.$store.dispatch("priceMinBuyItemList", this.minPrice);
      this.minPrice = ''

      },

      updatePriceMax() {
   //this.$store.dispatch("setBuyItemList");
     this.applied.push(this.maxPrice);
      this.$store.dispatch("priceMaxBuyItemList", this.maxPrice);
this.maxPrice = ''
      },
       clearList() {
     this.applied = []
     this.$store.dispatch("setBuyItemList");},

 search() {
   this.applied = []
      let array = this.$store.state.data.item.filter(item => !item.buyer && item.transferable === true)
      
      let rs = array.filter(item => item.description.toLowerCase().includes(this.input) || item.title.toLowerCase().includes(this.input)
      );
      this.applied.push(this.input)
      this.$store.commit("updateBuyItemList", rs);
    },
    },
   

    

   


   
 
};
</script>
