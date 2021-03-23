<template>
  <div class="pa-2 mx-auto">
   <div >
       <v-card>
    <v-card-title>
      All Items
  
      <v-spacer></v-spacer>
      <v-text-field
        v-model="search"
        append-icon="mdi-magnify"
        label="Search"
        single-line
        hide-details
      ></v-text-field>
    </v-card-title>
          <v-card-text class=caption>
      All of the items on the TPP blockchain </v-card-text>
    <v-data-table
      :loading="loading"
    
      :headers="headers"
      :items="items"
      :search="search"
    >
    <template v-slot:[`item.id`]="{ item }">
     <v-btn text :to="{ name: 'BuyItemDetails', params: { id: item.id } }">
        {{ item.id }}
     </v-btn>
    </template>
    
    </v-data-table>
  </v-card>
    </div>

  </div>
</template>

<script>

export default {
  props: [""],
  components: {  },
  data() {
    return {
     itemsPerPageArray: [4, 8, 12],
        search: '',
        loading: true,
      
        headers: [
          {
            text: 'Item',
            align: 'start',
            sortable: false,
            value: 'title',
          },
   
         
          { text: 'ID ', value: 'id' },
       
          { text: 'Price (in TPP)', value: 'estimationprice' },
          { text: 'Creator', value: 'creator' },
        
     
        ],
        items: this.$store.state.data.item

    };
  },
   mounted() {
       const type = { type: "item" };
       this.$store.dispatch("entityFetch",type);
      this.loading = false
  },


  computed: {
   
  },

  methods: {
   
   
 
  },
};
</script>

