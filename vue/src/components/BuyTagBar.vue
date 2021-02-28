<template>
  <div >
   
    <v-chip-group v-if="advanced"
    
          
          active-class="primary--text"
        >
          <v-chip @click="updateList(tag)"  outlined
            v-for="tag in tags" :key="tag"
          ><v-icon small left>
        mdi-tag-outline
      </v-icon>{{ tag }}
          </v-chip>
        </v-chip-group>
    

  

  </div>
  
</template>


<script>

import ItemListBuy from "./ItemListBuy.vue";
export default {
  props: ["advanced"],
  components: { ItemListBuy },
  data: function () {
    return {
      
      

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

  },
    

  methods: {

    updateList(tag) {
      this.$store.dispatch("tagBuyItemList", tag);},



   


   
  },
};
</script>
