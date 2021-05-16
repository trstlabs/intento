

<template>
  <div>
    <v-app>
  <v-main :class="($vuetify.theme.dark) ? 'grey darken-4' : 'primary lighten-5'" >
        <v-container class="pb-6 pt-0 px-0" >
          
         
              <v-row class="ma-0 pa-0">
          <v-col cols="2" class="d-none d-sm-flex mx-auto ma-0 pa-0" >
              
                <wallet />
            
            </v-col >

             <v-col cols="12" sm="8" class="pa-0 mx-auto">
              <v-sheet  min-height="70vh" class="rounded-b-xl"  elevation="6">
                <div class="pt-0 mt-0"> 
       
                   <v-tabs 
      
    fixed-tabs
      :dark="!!$vuetify.theme.dark"
      icons-and-text

  :background-color="($vuetify.theme.dark) ? 'dark' : 'light'"
  >
   <v-tab to="/account">
     Transactions<v-icon >
        mdi-cube-send
      </v-icon> 
    </v-tab> 
    <v-tab to="/account=placeditems">
     Created<v-icon >
        mdi-plus-box
      </v-icon> 
    </v-tab> 
    <v-tab to="/account=estimateditems">
      Estimated<v-icon >
        mdi-checkbox-marked
      </v-icon> 
    </v-tab>
    <v-tab to="/account=boughtitems">
      Bought<v-icon >
        mdi-shopping
      </v-icon> 
    </v-tab>
    <v-tab to="/account=likeditems">
      Liked<v-icon >
        mdi-heart
      </v-icon> 
    </v-tab>
  </v-tabs>
                  <!--<faucet/>-->

          <div v-if="this.$route.name == 'account'">
  <v-img src="img/design/buy.png" contain>  <p    class="display-2 pt-4 font-weight-thin gray--text text-center mb-n1">  Transactions</p><p class="overline pt-n10 font-weight-bold gray--text text-center pb-5 "> Browse through history<v-btn text icon @click="setTX"> <v-icon >
        mdi-refresh
      </v-icon></v-btn></p>  </v-img>

            
  
                  <transactions :key="update"/>
                   </div>
                  <item-list-seller v-if="created"/>
                  
                  <item-list-estimator v-if="estimated" />
                
                  <item-list-buyer v-if="bought"/>
        
                  <item-list-interested v-if="interested"/>
                  
                </div>
                
              </v-sheet>
            </v-col> <v-col cols="12" sm="2" class="d-none d-sm-flex">
           
            </v-col>
           </v-row><v-col cols="12" class="d-flex d-sm-none justify-center">
               
                <wallet />
              
            </v-col >
        </v-container>
      </v-main>
    </v-app>
  </div>
</template>

<script>
import Transactions from '../components/Transactions.vue';
import ItemListSeller from '../components/ItemListSeller.vue';
import ItemListEstimator from '../components/ItemListEstimator.vue';
import ItemListBuyer from '../components/ItemListBuyer.vue';
import ItemListInterested from '../components/ItemListInterested.vue';

export default {
  components: { Transactions, ItemListSeller, ItemListEstimator, ItemListBuyer, ItemListInterested },

  data() {
    return {
      created: false,
      estimated: false,
      bought: false,
      interested: false,
      update:true,
    };
  },

  mounted(){
    if (this.$route.params.list == "placeditems"){
this.getItemsFromSeller()

    }else if (this.$route.params.list == "estimateditems"){
this.getItemsFromEstimator()

    }else if (this.$route.params.list == "boughtitems"){
this.getItemsFromBuyer()

    }else if (this.$route.params.list == "likeditems"){
this.getInterestedItems()

    }
  },

  methods: {
    async setTX(){
            this.update = false
      await this.$store.dispatch("setTransactions")
      this.update = true

    },

   getItemsFromSeller() {
      if (this.$store.state.account.address) { 
  
      let input = this.$store.state.account.address;
      this.$store.dispatch("setSellerItemList", input);
};    this.created = true
    },

     getItemsFromEstimator() {
      if (this.$store.state.account.address) {
     
    
      let input = this.$store.state.account.address;
      this.$store.dispatch("setEstimatorItemList", input);
     }   this.estimated = true
    },
     getInterestedItems() {
      if (this.$store.state.account.address) { 

     
      let input = this.$store.state.account.address;
      this.$store.dispatch("setInterestedItemList", input);
      } this.interested = true

    },
     getItemsFromBuyer() {
      if (this.$store.state.account.address) { 
      const type = { type: "buyer" };
      this.$store.dispatch("entityFetch",type);
    
      let input = this.$store.state.account.address;
      this.$store.dispatch("setBuyerItemList", input);
      }  this.bought = true

    },

},
}
</script>