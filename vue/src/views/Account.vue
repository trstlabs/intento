
<!--<template>
  <div>
    <app-layout>
      <div>
        <app-text type="h1">TrustItems (v0.1)</app-text>
        <wallet />
        <item-list-creator />
        <item-list-estimator />
        <item-list-buyer />
      </div>
    </app-layout>
  </div>
</template> -->

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

             <v-col cols="12" sm="8" class="pa-0 mx-auto">
              <v-sheet min-height="70vh" rounded="lg" elevation="6">
                <div>
                   <v-tabs 
      
    fixed-tabs
      :dark="!!$vuetify.theme.dark"
      icons-and-text

  :background-color="($vuetify.theme.dark) ? 'dark' : 'white'"
  >
    <v-tab @click="getItemsFromCreator()">
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
                   
                  <item-list-creator v-if="created"/>
                  
                  <item-list-estimator v-if="estimated" />
                
                  <item-list-buyer v-if="bought"/>
        
                  <item-list-interested v-if="interested"/>
                  
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

   getItemsFromCreator() {
      if (!this.$store.state.account.address) { alert("Sign in first");};
       this.estimated = false
      this.bought = false
      this.interested = false
      this.created = true
      let input = this.$store.state.account.address;
      this.$store.dispatch("setCreatorItemList", input);
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