<template>
  <div>
    <v-app>
   <v-main :class="($vuetify.theme.dark) ? 'grey darken-4' : 'primary lighten-5'" >
        <v-container class="pb-6 pt-0 px-0" >
          
         
         <v-row class="ma-0 pa-0"> 
        

            <v-col cols="12" sm="10" class="mx-auto pa-0 ">
              <v-sheet  class="rounded-xl"  elevation="6">
                <div>
                  <create-item-form />
                </div>
              </v-sheet>
            </v-col>
             
          </v-row>
        </v-container>
      </v-main>
    </v-app>
  </div>
</template>