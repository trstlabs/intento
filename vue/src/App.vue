<template>
  <v-app >
    <div>
      <div>
        <!--<router-link name="index" to="/">Home</router-link>-->
        <!--<router-link name="buy" to="/">Buy</router-link>
        <router-link name="sell" to="/sell">Sell</router-link>
        <router-link name="Earn" to="/earn">Earn</router-link>
        <router-link name="Account" to="/account">Account</router-link>-->
      </div>
      <router-view :key="$route.path" />
    </div>
    
    <v-app-bar :color="($vuetify.theme.dark) ? 'grey darken-4' : 'grey lighten-4'"  app dense elevation="2" >
      
   <v-app-bar-nav-icon @click.stop="drawer = !drawer"></v-app-bar-nav-icon>
   <v-col xs="1" class="  mx-auto" >
    

     <v-tooltip bottom>
      <template v-slot:activator="{ on, attrs }">
     <router-link to="/">
     <v-img 
   v-bind="attrs"
          v-on="on"
    src="img/brand/icon.png"
    max-height="32"
    max-width="32"
    contain
  ></v-img>
     </router-link>
     
      </template>
      <span>Trust Price Marketplace</span>
    </v-tooltip>



    </v-col >
<v-col cols="6" xs="10" class="pa-0 d-none d-md-flex">

      <v-tabs show-arrows fixed-tabs  :background-color="($vuetify.theme.dark) ? 'grey darken-4' : 'grey lighten-4'" >
      
       <v-tooltip bottom >
      <template v-slot:activator="{ on, attrs }" >
       

        <v-tab to="/"><v-icon v-bind="attrs"
          v-on="on" large >
        mdi-shopping
      </v-icon></v-tab>
      </template>
      <span>Shop</span>
    </v-tooltip>

    <v-tooltip bottom>
      <template v-slot:activator="{ on, attrs }">
      <v-tab  to="/sell"> <v-icon v-bind="attrs"
          v-on="on" large >
        mdi-plus-box
      </v-icon></v-tab>
      </template>
      <span>Sell</span>
    </v-tooltip>

    <v-tooltip bottom>
      <template class="ma-0 pa-0" v-slot:activator="{ on, attrs }">
      <v-tab  to="/earn"> <v-icon v-bind="attrs"
          v-on="on" large >
        mdi-checkbox-marked
      </v-icon><v-badge 
          color="red"
           :content="messagesToEstimate"
        :value="messagesToEstimate"
        /></v-tab>
      </template>
      <span>Earn</span>
    </v-tooltip>

    <v-tooltip bottom>
      <template v-slot:activator="{ on, attrs }">
      <v-tab  to="/account"> <v-icon v-bind="attrs"
          v-on="on" large >
        mdi-account-box
      </v-icon><v-badge
          color="red"
           :content="messagesAccount"
        :value="messagesAccount"
        />
        </v-tab>
      </template>
      <span>Account</span>
    </v-tooltip>

   
      </v-tabs>

 

   </v-col><v-col class="text-right" > 

      <template><div class="">
      <v-btn rounded elevation="1"  small :color="($vuetify.theme.dark) ? 'primary' : 'primary'"
          to="/messages">
    
          <v-icon>mdi-message-reply
          </v-icon>
        </v-btn></div>
      </template>

  
<!-- <v-tooltip bottom>
      <template v-slot:activator="{ on, attrs }">
      <v-btn 
          icon id="mode-switcher"
          @click="toggledarkmode"
        >
          <v-icon v-bind="attrs"
          v-on="on" :color="($vuetify.theme.dark) ? 'primary' : 'primary lighten-1'">
            {{ ($vuetify.theme.dark) ? 'mdi-weather-night' : 'mdi-weather-sunny' }}
          </v-icon>
        </v-btn>
      </template>
      <span>Switch Theme</span>
    </v-tooltip>-->



    </v-col >
   
    </v-app-bar>
       
    <v-navigation-drawer
      v-model="drawer"
      
      app
      temporary
    > 
      <v-list
        nav
        dense
      >
      <v-alert dense  v-if="!this.$store.state.user && this.$store.state.account.address "
  type="warning" dismissible class="caption"
>Confirm the verification link sent to the email linked to your Google account</v-alert>
       <div class="text-center">
   <img class="pa-4 "
 
    src="img/brand/icon.png"
    width="77" /></div>
      <wallet v-if="this.$store.state.account.address"/>
        <v-list-item-group
         
          active-class="blue--text text--accent-4"
        >
          <v-list-item to="/">
            <v-list-item-title >Buy</v-list-item-title><v-icon >
        mdi-shopping
      </v-icon>
          </v-list-item>

          <v-list-item to="/sell">
            <v-list-item-title>Sell</v-list-item-title><v-icon >
        mdi-plus-box
      </v-icon>
          </v-list-item>

          <v-list-item to="/earn">
            <v-list-item-title>Earn</v-list-item-title><v-icon   >
        mdi-checkbox-marked
      </v-icon>
          </v-list-item>

          <v-list-item to="/account">
            <v-list-item-title>Account</v-list-item-title> <v-icon  >
        mdi-account-box
      </v-icon>
          </v-list-item>
          
           <v-list-item to="/messages">
            <v-list-item-title>Messages</v-list-item-title> <v-icon>mdi-message-reply
          </v-icon>
          </v-list-item>
            <v-list-item to="/explore">
            <v-list-item-title>Explore</v-list-item-title>
          </v-list-item>
        
           <v-list-item v-if="!this.$store.state.account.address" inactive @click="welcome = !welcome">
            <v-list-item-title>Get Started</v-list-item-title>
          </v-list-item>
          <v-list-item  target="_blank" href="https://www.trustpriceprotocol.com">
            <v-list-item-title>About TPP</v-list-item-title>
          </v-list-item>
             <v-list-item inactive
           id="mode-switcher"
          @click="toggledarkmode"
        >  <v-list-item-title>Theme</v-list-item-title><v-icon :color="($vuetify.theme.dark) ? 'primary' : 'primary lighten-1'">
            {{ ($vuetify.theme.dark) ? 'mdi-weather-night' : 'mdi-weather-sunny' }}
          </v-icon>
          </v-list-item >
        </v-list-item-group>
        
      </v-list>
    </v-navigation-drawer>
    
    <welcome v-if="!this.$store.state.account.address && welcome"/> 
    <v-btn
            v-scroll="onScroll"
            v-show="fab"
            fab
            dark
            fixed
            bottom
            right
            color="primary"
            @click="toTop"
          >
            <v-icon>mdi-arrow-up</v-icon>
          </v-btn>
         <v-footer class="text-center caption"  padless>
    <v-col>{{ new Date().getFullYear() }} — <strong>© Trust Price Protocol</strong>
    </v-col>
  </v-footer>
  </v-app>
  
</template>





<script>
import Wallet from './components/Wallet.vue';
export default {
  data() {
    return {

      //dismiss: false,
      //dialog: true,
      fab: false,
      drawer: false,
      welcome: true
     
      
    };
  },
  created() {
    this.$store.dispatch("init");
     
    //this.$store.dispatch("setBuyItemList");
    
    //this.$store.dispatch("initBuyItemList");
  },
  /* messages() {
      return this.$store.getters.getCreatorItemList.length || 0;
    },*/
  
    methods: {
      
        toggledarkmode: function () {
            this.$vuetify.theme.dark = !this.$vuetify.theme.dark;
            localStorage.setItem("dark\_theme", this.$vuetify.theme.dark.toString());
        },

       onScroll (e) {
      if (typeof window === 'undefined') return
      const top = window.pageYOffset ||   e.target.scrollTop || 0
      this.fab = top > 20
    },
    toTop () {
      this.$vuetify.goTo(0)
    }
  
    },
    
    
   computed: {
    messagesAccount() {
       
      return this.$store.getters.getCreatorActionList.length || 0;
    },
    messagesToEstimate() {
      return this.$store.getters.getToEstimateList.length || 0;
    },

  },
    mounted() {
     
    

      //this.messages = this.$store.getters.getCreatorItemList.length;
        const theme = localStorage.getItem("dark\_theme");
        if (theme) {
            if (theme == "true") {
                this.$vuetify.theme.dark = true;
            } else {
                this.$vuetify.theme.dark = false;
            }
        }
    },
};
</script>