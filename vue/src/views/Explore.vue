

<template>
  <div>
    <v-app>
      
  <v-main :class="($vuetify.theme.dark) ? 'grey darken-4' : 'primary lighten-5'" >
       <!-- <v-alert
      dense
      type="warning"
    >
     This App is being developed. Go to <a color="default" href="http://www.trustpriceprotocol.com">main page </a> instead?
    </v-alert> -->
         <v-container class="pb-6 pt-0 px-0" >
          
         
            <v-row class="ma-0 pa-0">
          

            <v-col cols="12" sm="10" class="mx-auto pa-0">
              <v-sheet :class="($vuetify.theme.dark) ? 'grey darken-4' : 'secondary lighten-2' "  min-height="70vh" class="rounded-xl"  elevation="6">
                <div>
                     <v-img height="350px" src="img/design/market.png" ></v-img>
                  <leaderboard/>

                <all-items/>
                </div>
              </v-sheet>
            </v-col>
            
          </v-row>
        </v-container>
      </v-main>
    </v-app>
  </div>
</template>