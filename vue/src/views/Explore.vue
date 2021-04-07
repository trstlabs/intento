

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
         <v-container class="pa-0" >
          
         
          <v-col cols="12" sm="8" class="d-none d-sm-flex d-md-none mx-auto" >
              <v-sheet min-width="350" class="mx-auto" rounded="lg" elevation="1" >
                <wallet />
              </v-sheet>
            </v-col >
          

          <v-row>
            <v-col cols="12" sm="2" class="d-sm-none d-lg-flex d-md-flex">
              
                <wallet />
              
            </v-col >

            <v-col cols="12" sm="8" class="mx-auto pa-0">
              <v-sheet min-height="70vh" rounded="lg" elevation="6">
                <div>
                     <v-img height="350px" src="img/design/market.png" ></v-img>
                  <leaderboard/>

                <all-items/>
                </div>
              </v-sheet>
            </v-col>
           <v-col cols="12" sm="2" class="d-sm-none d-lg-flex d-md-flex">
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