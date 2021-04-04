



<template>
  <div>
    <v-app>
      <v-main :class="($vuetify.theme.dark) ? 'grey darken-4' : 'secondary lighten-2'">
        <v-container class="pa-0">
          
         
          <v-col cols="12" sm="8" class="d-none d-sm-flex d-md-none mx-auto" >
              <v-sheet min-width="350" class="mx-auto" rounded="lg" elevation="1" >
                <wallet v-if="!!this.$store.state.account.address"/>
              </v-sheet>
            </v-col >
          

          <v-row>
            <v-col cols="12" sm="2" class="d-sm-none d-lg-flex d-md-flex">
              
                <wallet v-if="!!this.$store.state.account.address" />
 
            </v-col >

             <v-col cols="12" sm="8"  class="pa-2 mx-auto">
              <v-sheet min-height="70vh" rounded="lg" elevation="6">
                <div>
               
                  <div>
    
    <div> <v-row class="pa-4">
     <v-btn icon plain to="/buy"> <v-icon >
        mdi-arrow-left
      </v-icon></v-btn><v-spacer/><v-btn icon onClick="window.location.reload();" ><v-icon left>
        mdi-refresh
      </v-icon></v-btn></v-row>
      <buy-item-item :itemid="this.$route.params.id" />
      
    </div>
    <v-btn block plain class="pa-2" to="/buy"> <v-icon left>
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

export default {

  data() {
    return { };
  },
  
  

};
</script>
