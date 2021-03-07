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
   <v-col xs="1">
    

     <v-tooltip bottom>
      <template v-slot:activator="{ on, attrs }">
      <v-btn
          icon id="mode-switcher"
          @click="toggledarkmode"
        >
          
          <v-img
   v-bind="attrs"
          v-on="on"
    src="img/brand/icon.png"
    max-height="30"
    max-width="30"
    contain
  ></v-img>
        </v-btn>
      </template>
      <span>Trust Price Marketplace</span>
    </v-tooltip>



    </v-col >
<v-col cols="6" xs="10" class="pa-0" >

      <v-tabs grow centered :background-color="($vuetify.theme.dark) ? 'grey darken-4' : 'grey lighten-4'" >
       
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


  

   </v-col><v-col class="text-right"> 

      <template>
      <v-btn icon plain :color="($vuetify.theme.dark) ? 'primary' : 'primary lighten-1'"
          href="/messages">
    
          <v-icon>mdi-message-reply
          </v-icon>
        </v-btn>
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
    <welcome v-if="!this.$store.state.account.address"/> <v-btn
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
          
  </v-app>
  
</template>





<script>
export default {
  data() {
    return {
      //dismiss: false,
      //dialog: true,
      fab: false,
     
      
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