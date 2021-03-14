

<template>
  <div>
    <v-app>
      <v-main :class="($vuetify.theme.dark) ? 'grey darken-4' : 'grey lighten-3'" >
        <v-container class="pa-0">
          
         
          <v-col cols="12" sm="8" class=" d-none d-sm-flex d-md-none mx-auto" >
              <v-sheet min-width="350" class="mx-auto" rounded="lg" elevation="1" >
                <wallet v-if="this.$store.state.account.address"/>
              </v-sheet>
            </v-col >
          

          <v-row>
            <v-col cols="12" sm="2" class="d-sm-none d-lg-flex d-md-flex">
             
                <wallet v-if="this.$store.state.account.address"/>
             
            </v-col >

            <v-col cols="12" sm="8" class="mx-auto pa-0 ">
              <v-sheet min-height="70vh" rounded="lg" elevation="6">
                <div>
                  <item-list-buy />
                </div>
              </v-sheet>
            </v-col>
            <v-col cols="12" sm="2" class="d-sm-none d-lg-flex d-md-flex">
              
            </v-col>
          </v-row>
        </v-container>
      </v-main>
    </v-app>
  </div>
</template>