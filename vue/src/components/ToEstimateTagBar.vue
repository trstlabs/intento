<template>
  <div >
   
    <v-chip-group 
    
          
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

import ItemToEstimate from "./ItemToEstimate.vue";
export default {
  //props: ["advanced"],
  components: { ItemToEstimate },
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
      //console.log(this.tag);
      this.$store.dispatch("tagToEstimateList", tag);},



   


   
  },
};
</script>
