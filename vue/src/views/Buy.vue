

<template>
  <div>
    <v-app>
      <v-main :class="($vuetify.theme.dark) ? 'grey darken-4' : 'primary lighten-5'" >
        <v-container class="pb-6 pt-0 px-0" >
          
           <v-row class="ma-0 pa-0">
          <v-col cols="2" class="d-none d-sm-flex mx-auto ma-0 pa-0" >
             
                <wallet v-if="this.$store.state.account.address"/>
             
            </v-col >

            <v-col cols="12" sm="8" class="mx-auto pa-0">
              <v-sheet  min-height="70vh" class="rounded-b-xl"  elevation="6">
                <div>
             
                  <item-list-buy />
                </div>
              </v-sheet>
            </v-col>
              <v-col cols="12" sm="2" class="d-none d-sm-flex">
              
            </v-col>
          </v-row>
        </v-container>
      </v-main>
    </v-app>
  </div>
</template>