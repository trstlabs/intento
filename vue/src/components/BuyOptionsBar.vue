<template>
  <div >
   
    
    <v-select  class="mx-auto rounded-lg"
             
             v-model="selectedFilter"
             v-on:input="search()"
              :menu-props="{ offsetY: true }"
              solo
              item-text="filter"
              item-value=" abbr"
              hide-details
              :items="filters"
              prepend-inner-icon="mdi-filter"
              label="Sort by"
            ></v-select>

            
  </div>
  
</template>


<script>

import ItemListBuy from "./ItemListBuy.vue";
export default {
  props: ["items"],
  components: { ItemListBuy },
  data: function () {
    return {
      selectedFilter: "",
    filters: [
          { filter: 'All', abbr: 'NE' },
          { filter: 'Local', abbr: 'GA' },
          { filter: 'Shipping', abbr: 'GA' },
          { filter: 'Macbook', abbr: 'transferable' },

 
        ],

  };
  },
  
    

  methods: {


    search() {
      console.log(this.selectedFilter);
      if (this.selectedFilter == "Shipping") {
     let input = this.$store.state.data.item.filter(item => item.buyer === "" && item.transferable === true && item.shippingcost > 0);     
      this.$store.dispatch("filterBuyItemList", input);
      };
      
      if (this.selectedFilter == "Macbook") {this.$store.dispatch("updateBuyItemList", "macbook"); };
      if (this.selectedFilter == "Local") {
        this.$store.dispatch("setLocalBuyItemList"); };
      if (this.selectedFilter == "All") {
        this.$store.dispatch("setBuyItemList");
      
      
      
      
    };},

   
  },
};
</script>
