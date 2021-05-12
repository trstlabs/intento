



<template>
  <div>
    <v-app>
      <v-main :class="($vuetify.theme.dark) ? 'grey darken-4' : 'secondary lighten-2'">
       <v-container class="pb-6 pt-0 px-0" >
          
         
         <v-row class="ma-0 pa-0">
          <v-col cols="2" class="d-none d-sm-flex mx-auto ma-0 pa-0" >
              
               <wallet v-if="this.$store.state.account.address"/>
 
            </v-col >

             <v-col cols="12" sm="8" class="mx-auto pa-0 mt-n4">
               <v-sheet  min-height="70vh" class="rounded-b-xl"  elevation="6">
              <buy-item-item :itemid="this.$route.params.id" />
  <v-row class="mx-0 mt-2">
     <v-btn icon plain @click="$router.go(-1)"> <v-icon >
        mdi-arrow-left
      </v-icon></v-btn><v-spacer/><v-btn icon @click="$router.go()" ><v-icon left>
        mdi-refresh
      </v-icon></v-btn></v-row>
      
  <v-btn block plain @click="$router.push('/')"> <v-icon left>
        mdi-arrow-left
      </v-icon>To All Items</v-btn>

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
