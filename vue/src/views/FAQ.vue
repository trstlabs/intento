

<template>
  <div>
    <v-app>
      
      <v-main :class="($vuetify.theme.dark) ? 'grey darken-4' : 'secondary lighten-2'">
       <!-- <v-alert
      dense
      type="warning"
    >
     This App is being developed. Go to <a color="default" href="http://www.trustpriceprotocol.com">main page </a> instead?
    </v-alert> -->
         <v-container class="pb-6 pt-0 px-0" >
          
         
         <v-row class="ma-0 pa-0">
          <v-col cols="2" class="d-none d-sm-flex mx-auto ma-0 pa-0" >
              
               <wallet/>
            
            </v-col >

            <v-col cols="12" sm="8" class="mx-auto pa-0">
              <v-sheet min-height="70vh" class="rounded-b-xl"  elevation="6">
                <div>
                  <v-row class="pt-6 ml-8">
     <v-btn  plain to="/"> <v-icon >
        mdi-arrow-left
      </v-icon> To Shop</v-btn></v-row>   <v-img height="350px" src="img/design/market.png" ></v-img>
                  <info/>


                </div>
              </v-sheet>
            </v-col>
             <v-col cols="12" sm="2" class="d-none d-sm-flex">
              <!--<v-sheet rounded="lg" min-height="268" elevation="2">
                  
                <side-bar-earn />
              </v-sheet>-->
            </v-col>
          </v-row>
        </v-container>
      </v-main>
    </v-app>
  </div>
</template>