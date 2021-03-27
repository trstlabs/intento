<template>
  <div >
   
    
    <v-select  class="mx-auto rounded-lg"
             
             v-model="selectedFilter"
             v-on:input="search()"
              :menu-props="{ offsetY: true }"
              solo
              item-text="filter"
              
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
          { filter: 'All',  },
          { filter: 'Pickup', },
          { filter: 'Shipping', },
          { filter: 'Trust Price', },
          { filter: 'Reposted', },
 
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
       if (this.selectedFilter == "Trust Price") {
     let input = this.$store.state.data.item.filter(item => item.buyer === "" && item.transferable === true && item.seller == item.creator);     
      this.$store.dispatch("filterBuyItemList", input);
      };
       if (this.selectedFilter == "Reposted") {
     let input = this.$store.state.data.item.filter(item => item.buyer === "" && item.transferable === true && item.seller != item.creator);     
      this.$store.dispatch("filterBuyItemList", input);
      };
      if (this.selectedFilter == "Pickup") {
        this.$store.dispatch("setLocalBuyItemList"); };
      if (this.selectedFilter == "All") {
        this.$store.dispatch("setBuyItemList");
      
      
      
      
    };},

   
  },
};
</script>
