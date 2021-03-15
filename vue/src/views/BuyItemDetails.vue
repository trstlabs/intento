



<template>
  <div>
    <v-app>
      <v-main :class="($vuetify.theme.dark) ? 'grey darken-4' : 'grey lighten-3'">
        <v-container class="pa-0">
          
         
          <v-col cols="12" sm="8" class="d-none d-sm-flex d-md-none mx-auto" >
              <v-sheet min-width="350" class="mx-auto" rounded="lg" elevation="1" >
                <wallet />
              </v-sheet>
            </v-col >
          

          <v-row>
            <v-col cols="12" sm="2" class="d-sm-none d-lg-flex d-md-flex">
              
                <wallet />
 
            </v-col >

             <v-col cols="12" sm="8"  class="pa-0 mx-auto">
              <v-sheet min-height="70vh" rounded="lg" elevation="6">
                <div>
                  <!--<faucet/>-->
                  <div>
    
    <div> 
     
      <buy-item-item :itemid="this.$route.params.id" />
      
    </div>
    <v-btn block class="pa-2" to="/buy"> <v-icon left>
        mdi-arrow-left
      </v-icon>All items </v-btn>
  </div>
                </div>
              </v-sheet>
            </v-col>
            <v-col cols="12" sm="2" class="d-sm-none d-lg-flex d-md-flex">
              <!--<v-sheet rounded="lg" min-height="268" elevation="2">
                
              </v-sheet>-->
            </v-col>
          </v-row>
        </v-container>
      </v-main>
    </v-app>
  </div>
</template>

<script>
import BuyItemItem from "../components/BuyItemItem.vue";
import Wallet from "../components/Wallet.vue";
export default {
  components: { BuyItemItem, Wallet },
  data() {
    return { };
  },
  
  

};
</script>
