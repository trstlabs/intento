

<template>
  <div>
    <v-app>
      
 <v-main :class="($vuetify.theme.dark) ? 'grey darken-4' : 'primary lighten-5'" >
     
         <v-container class="pb-6 pt-0 px-0" >
          
         
         <v-row class="ma-0 pa-0">
         

            <v-col cols="12" sm="10" class="mx-auto pa-0">
              <v-sheet min-height="70vh" class="rounded-xl"  elevation="6">
                <div>
                  <v-row class="pt-6 ml-8">
     <v-btn  plain to="/"> <v-icon >
        mdi-arrow-left
      </v-icon> To Shop</v-btn></v-row>   <v-img height="350px" src="img/design/market.png" ></v-img>
                  <info/>


                </div>
              </v-sheet>
            </v-col>
            
          </v-row>
        </v-container>
      </v-main>
    </v-app>
  </div>
</template>