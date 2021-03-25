<template>
  <div class="pa-2 mx-auto">
   <div >
       <v-card>
    <v-card-title>
      Current leaderboard </v-card-title>
       <v-card-text class=caption>
      The users that correctly estimated the most (currently transferrable) items. </v-card-text>
   
     <v-simple-table fixed-header
    height="300px">
    <template v-slot:default>
      <thead>
        <tr>
           <th class="text-left">
            Rank
          </th>
          <th class="text-left">
            Expert
          </th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="(item, index) in items"
          :key="index"
        >
        
          <td>{{ index + 1 }}</td>
            <td>{{ item }}</td>
        </tr>
      </tbody>
    </template>
  </v-simple-table>
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
        //filter: {},
        //sortDesc: false,
       // page: 1,
        //itemsPerPage: 4,
        /*sortBy: 'title',
        keys: [
          'title',
          'description',
          'id',
          "estimationprice",
          "seller",
          "buyer",
          "shippingregion",
          "tags",
         
        ],*/
        headers: [
          {
            text: 'Item',
            align: 'start',
            sortable: false,
            value: 'title',
          },
   
         
        //  { text: 'Id ', value: 'id' },
       
         // { text: 'Price (in TPP)', value: 'estimationprice' },
          { text: 'Expert', value: 'bestestimator' },
        
     
        ],
        

    };
  },
   mounted() {
       const type = { type: "item" };
       this.$store.dispatch("entityFetch",type);
      this.loading = false
  },


  computed: {
    items(){
    let rs = this.$store.state.data.item.map(item => item.bestestimator);
     console.log(rs)
    let merged = [].concat.apply([], rs);
      let frequency = {};
      merged.forEach(function (value) { if (value != '') {frequency[value.toLowerCase()] = 0;} });
 console.log(merged)
      let uniques = merged.filter(function (value) {
        return ++frequency[value] == 1;
      });

      let sorted = uniques.sort(function (a, b) {
        return frequency[b] - frequency[a];
      });

    console.log(sorted)
    let toreturn = { text: 'Expert', value: sorted }
    console.log(toreturn)
    return sorted


    },
   
  },

  methods: {
   
   
 
  },
};
</script>

