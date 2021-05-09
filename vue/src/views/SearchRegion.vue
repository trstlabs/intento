<template>
  <div>
    <v-app>
      <v-main :class="($vuetify.theme.dark) ? 'grey darken-4' : 'secondary lighten-2'">
       <v-container class="pb-6 pt-0 px-0" >
          
         
          <v-row class="ma-0 pa-0">
          <v-col cols="2" class="d-none d-sm-flex mx-auto ma-0 pa-0" >
              
               <wallet v-if="this.$store.state.account.address"/>
            
            </v-col >
             <v-col cols="12" sm="8"  class="pa-0 mx-auto">
              <v-sheet min-height="70vh" rounded="lg" elevation="6">
                <div>
               
                  <div>
    
    <div> <v-row class="pa-4">
     <v-btn icon plain @click="$router.go(-1)"> <v-icon >
        mdi-arrow-left
      </v-icon></v-btn><v-spacer/><v-btn icon onClick="window.location.reload();" ><v-icon left>
        mdi-refresh
      </v-icon></v-btn></v-row>
      <search-list-region :region="this.$route.params.region" />
      
    </div>
    <v-btn block plain class="mx-2" @click="$router.go(-1)"> <v-icon left>
        mdi-arrow-left
      </v-icon>Back </v-btn>
  </div>
                </div>
              </v-sheet>
            </v-col>
              <v-col cols="12" sm="2" class="d-none d-sm-flex">
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
