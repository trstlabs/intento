
<!--<template>
  <div>
    <app-layout>
      <div>
        <app-text type="h1">TrustItems (v0.1)</app-text>
        <wallet />
        <item-list-seller />
        <item-list-estimator />
        <item-list-buyer />
      </div>
    </app-layout>
  </div>
</template> -->

<template>
  <div>
    <v-app>
      <v-main :class="($vuetify.theme.dark) ? 'grey darken-4' : 'secondary lighten-2'">
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

  :background-color="($vuetify.theme.dark) ? 'dark' : 'white'"
  >
    <v-tab @click="getItemsFromSeller()">
     Created<v-icon >
        mdi-plus-box
      </v-icon> 
    </v-tab> 
    <v-tab @click="getItemsFromEstimator()">
      Estimated<v-icon >
        mdi-checkbox-marked
      </v-icon> 
    </v-tab>
    <v-tab @click="getItemsFromBuyer()">
      Bought<v-icon >
        mdi-shopping
      </v-icon> 
    </v-tab>
    <v-tab @click="getInterestedItems()">
      Liked<v-icon >
        mdi-heart
      </v-icon> 
    </v-tab>
  </v-tabs>
                  <!--<faucet/>-->
                   
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

export default {

  data() {
    return {
      created: true,
      estimated: false,
      bought: false,
      interested: false,
    };
  },

  

  methods: {

   getItemsFromSeller() {
      if (!this.$store.state.account.address) { alert("Sign in first");};
       this.estimated = false
      this.bought = false
      this.interested = false
      this.created = true
      let input = this.$store.state.account.address;
      this.$store.dispatch("setSellerItemList", input);
      //this.dummy = false;
    },

     getItemsFromEstimator() {
      if (!this.$store.state.account.address) { alert("Sign in first");};
      this.bought = false
      this.interested = false
      this.created = false
       this.estimated = true
      let input = this.$store.state.account.address;
      this.$store.dispatch("setEstimatorItemList", input);
      //this.dummy = false;
    },
     getInterestedItems() {
      if (!this.$store.state.account.address) { alert("Sign in first");};
      this.bought = false
      this.created = false
       this.estimated = false
      this.interested = true
      let input = this.$store.state.account.address;
      this.$store.dispatch("setInterestedItemList", input);
      //this.dummy = false;
    },
     getItemsFromBuyer() {
      if (!this.$store.state.account.address) { alert("Sign in first");};
      const type = { type: "buyer" };
      this.$store.dispatch("entityFetch",type);
      this.interested = false
      this.created = false
       this.estimated = false
      this.bought = true
      let input = this.$store.state.account.address;
      this.$store.dispatch("setBuyerItemList", input);
      //this.dummy = false;
    },

},
}
</script>