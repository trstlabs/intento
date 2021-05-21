

<template>
  <div>
    <v-app>
      <v-main :class="($vuetify.theme.dark) ? 'grey darken-4' : 'primary lighten-5'" >
        <v-container class="pb-6 pt-0 px-0" >
          
           <v-row class="ma-0 pa-0">
       <!---   <v-col cols="12" class="d-none d-md-flex mx-auto ma-0 pa-0" >
         
                   
<v-navigation-drawer  fixed >
      <v-list
        shaped nav
      >
     
       <div class="text-center">
   <img class="pa-4 "
 
    src="img/brand/icon.png"
    width="77" /></div>
    <wallet v-if="this.$store.state.account.address"/>
        <v-list-item-group
         
          active-class="blue--text text--accent-4"
        >
          <v-list-item to="/">
           <v-list-item-icon><v-icon >
        mdi-shopping
      </v-icon>
          </v-list-item-icon>
            <v-list-item-title >Buy</v-list-item-title>
          </v-list-item>

          <v-list-item to="/sell"> <v-list-item-icon><v-icon >
        mdi-plus-box
      </v-icon>
          </v-list-item-icon>
            <v-list-item-title>Sell</v-list-item-title>
          </v-list-item>

          <v-list-item to="/estimate"> <v-list-item-icon><v-icon   >
        mdi-checkbox-marked
      </v-icon>
          </v-list-item-icon>
            <v-list-item-title>Estimate</v-list-item-title>
          </v-list-item>

          <v-list-item to="/account"> <v-list-item-icon><v-icon  >
        mdi-account-box
      </v-icon>
          </v-list-item-icon>
            <v-list-item-title>Account</v-list-item-title>
          </v-list-item>
          
           <v-list-item to="/messages" > <v-list-item-icon><v-icon>mdi-message-reply
          </v-icon>
          </v-list-item-icon>
            <v-list-item-title>Messages</v-list-item-title>
          </v-list-item>
            <v-list-item to="/explore">
            <v-list-item-title>Explore</v-list-item-title>
          </v-list-item>
       <div  @click="welcome = !welcome" v-if="!this.$store.state.account.address" >
           <v-list-item  @click="welcome = !welcome" inactive>
            <v-list-item-title>Get Started</v-list-item-title>
          </v-list-item></div>
          <v-list-item  target="_blank" href="https://www.trustpriceprotocol.com">
            <v-list-item-title>About TPP</v-list-item-title> 
          </v-list-item>
           <v-list-item  to="/FAQ">
            <v-list-item-title>FAQ</v-list-item-title>
          </v-list-item>
                       <v-list-item inactive
           id="mode-switcher"
          @click="$vuetify.theme.dark = !$vuetify.theme.dark"
        >  <v-list-item-title>Theme</v-list-item-title><v-icon :color="($vuetify.theme.dark) ? 'primary' : 'primary lighten-1'">
            {{ ($vuetify.theme.dark) ? 'mdi-weather-night' : 'mdi-weather-sunny' }}
          </v-icon>
          </v-list-item >
        </v-list-item-group>
        
      </v-list></v-navigation-drawer>

             
            </v-col >-->

            <v-col cols="12" sm="10" class="mx-auto pa-0">
              <v-sheet  min-height="70vh" class="rounded-xl"  elevation="6">
                <div>
             
                  <item-list-buy />
                </div>
              </v-sheet>
            </v-col>
              
          </v-row>
        </v-container>
      </v-main>
    </v-app>
  </div>
</template>