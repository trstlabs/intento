<template><div >
<v-dialog
      v-model="dialog"
      width="600" persistent
    >
     

      <v-card class="text-center">
        <v-card-title class="justify-center" >
          Welcome! <v-spacer/> <img src="img/brand/icon.png" width="44"/>
        </v-card-title>
       
  <div v-if="!signin && !signup && !learn">
        <v-card-text > 
          This dApp is in development. Feel free to use the app and see how it behaves. Note: items are examples and are not redeemable for tokens. Feedback is always welcome at the <a href="https://www.trustpriceprotocol.com"> main page.</a>
        </v-card-text>
        <torus/>

       

       <!-- <v-card-actions> 
          <v-col>
          <v-btn
            color="primary" block
            text
            @click="signin = true, loadContent()"
          >
            Sign In
          </v-btn> </v-col>
          <v-col>
          <v-btn
            color="primary darken-1" block
            text
            @click="dialog = false, loadContent()"
          >
            Learn more
          </v-btn></v-col><v-col> 
          <v-btn
            color="primary lighten-1" block
            text
            @click="dialog = false, loadContent()"
          >
            Look Around
          </v-btn> 
          </v-col>
          <v-col>
          
          <v-btn block
            color="primary "
            text
            @click="signup = true,loadContent()"
          >
            Sign Up
          </v-btn></v-col>
        </v-card-actions>-->
            <v-card-actions> 
       
          <v-btn
            color="primary" 
            text
            @click="signin = true, loadContent()"
          >
            Sign In
          </v-btn> 
          <v-spacer/>
            <v-btn
            color="primary lighten-1" 
            text
            @click="dialog = false, loadContent()"
          >
            Look Around
          </v-btn> 
          <v-btn
            color="primary lighten-1" 
            text
            @click="learn = true, loadContent()"
          >
            Learn more
          </v-btn>
        
         
          <v-spacer/>
          <v-btn 
            color="primary "
            text
            @click="signup = true,loadContent()"
          >
            Sign Up
          </v-btn>
        </v-card-actions>
  </div>
  <div v-if="signin">
    <torus/>
    <wallet @signedIn="updateDialog()"/>
    <v-card-actions> <v-col>
          <v-btn 
            color="primary"
            text
            @click="signin = false"
          >
            Back
          </v-btn>
         </v-col> 
        </v-card-actions>
  </div>

  <div v-if="signup">
    <v-card-text > 
         Receive 5 free tokens to get started! Create a new account linked to Google using DirectAuth.
        </v-card-text> <div v-if="!existing"> <faucet-torus />  
       <v-row> <v-divider class="ma-2"/> <p class="caption">  Or </p> <v-divider class="ma-2"/> </v-row>
        <v-btn text @click="existing = true"> 
        Sign up with an existing cosmos-address. </v-btn>
</div>
    <faucet v-if="existing"/>

<v-card-actions> <v-col>
          
          
          <v-btn block
            color="primary"
            text
            @click="signup = false, existing = false"
          >
            Back
          </v-btn></v-col><v-col> <v-btn
            color="primary" block
            text
            @click="signin = true, signup = false, loadContent()"
          >
            Sign In
          </v-btn></v-col>
        </v-card-actions>
    
    </div>

    <div v-if="learn">
<v-stepper v-model="e1">
    <v-stepper-header>
      <v-stepper-step
        :complete="e1 > 1"
        step="1"
      >
        Shop
      </v-stepper-step>

      <v-divider></v-divider>

      <v-stepper-step
        :complete="e1 > 2"
        step="2"
      >
        Sell
      </v-stepper-step>

      <v-divider></v-divider>

      <v-stepper-step step="3">
        Earn
      </v-stepper-step>
    </v-stepper-header>

    <v-stepper-items>
      <v-stepper-content step="1">
        <v-card
          class="mb-6"
          
         
        >  <v-card-title>There is a problem with current online marketplaces. </v-card-title><v-card-text > 
        


 <p>Ever wanted to buy something and the item was already granted to another user?  </p>
 <p>Ever paid too much because you bought the item with the wrong idea about the item? </p>
 <p>Trust price protocol is a place where you will have the power you need to buy items more hassle free.  </p>
 <p>* When you provide prepayment, you will from there on be seen as the only buyer </p>
 <p>*Prices are made up from independent estimations. So you know what you buy is the right value. </p>
 <p>*You are in control of your prepayment and can always get it back (until your item is shipped or transferred to you) </p>

 <p>A great way to spend crypto on things you like  </p>
        </v-card-text></v-card>
<v-row class="ma-2">
<v-btn text @click="learn = false">
          Back
        </v-btn>
        <v-spacer />
        <v-btn
          color="primary"
          @click="e1 = 2"
        >
          Continue
        </v-btn>

</v-row>

        
      </v-stepper-content>

      <v-stepper-content step="2">
        <v-card
          class="mb-12"
          
        ><v-card-title>An opportunity to enter the crypto universe </v-card-title>
        <v-card-text> 
 <p>As a seller, you are free to choose 2 options. Ship the item and/or choose “local pickup”.</p>
 <p>If you choose shipping, you may charge shipping costs and these are separate from the item price.</p>
 <p>You can set the accuracy of the price. This is set to 3 estimations by default. Setting the accuracy higher will result in a better accuracy but this will also take longer. When a final price is made, you may always decline or accept it. </p>
 <p>You don’t have to worry about pricing, just provide good quality pictures and information and you are done.</p>

 <p>Selling items has never been this simple. </p></v-card-text></v-card>
<v-row class="ma-2">
<v-btn text @click="e1 = 1">
          Back
        </v-btn>
        <v-spacer />
        <v-btn
          color="primary"
          @click="e1 = 3"
        >
          Continue
        </v-btn>

       
</v-row>
      </v-stepper-content>

      <v-stepper-content step="3">
        <v-card
          class="mb-12"
        
        ><v-card-title>Trade your time on the internet for crypto </v-card-title>
        <v-card-text> <p>Do you enjoy looking around the internet and finding/comparing prices? Or looking at things you cannot afford just for fun? Do you want the best value for your money? Then you may enjoy becoming an “estimator”</p>
 <p>As an “estimator” your job is to moderate items before they enter the marketplace. And you will be rewarded accordingly!</p>
 <p>Currently, you will earn roughly 3% of the final selling price if you are the estimator closest to the final price. </p>
 <p>To cope with bad acting, there is a  very small deposit required for each estimation, of 5 tokens (currently). These will be returned, except when the following occurs 1). The seller did not accept the final price and you are the lowest estimator. 2) The buyer ended up withdrawing prepayment and you are the highest estimator. </p>

 <p>As a bonus, you can like items and view then once they hit the marketplace. And you can trade your earned tokens for the items you like.</p> 
    </v-card-text>    </v-card>
<v-row class="ma-2">
  <v-btn text @click="e1 = 2">
          Back
        </v-btn>
        <v-spacer/>
        <v-btn text
          color="primary"
          @click="learn = false, dialog = false"
        >
          Look Around
        </v-btn>
        <v-spacer/>
        <v-btn
          color="primary"
          @click="learn = false, signup = true"
        >
          Sign Up
        </v-btn>

        </v-row>
      </v-stepper-content>
    </v-stepper-items>
  </v-stepper>
      
    </div>
      </v-card>
    </v-dialog>
    
          </div>
</template>

  <script>
import Faucet from './Faucet.vue';
import FaucetTorus from './FaucetTorus.vue';
import Torus from './Torus.vue';
import Wallet from './Wallet.vue';
//import * as bip39 from 'bip39'
  export default {
  components: { Wallet, Faucet, Torus, FaucetTorus },
  data() {
    return {
      //dismiss: false,
      dialog: true,
      signin: false,
      signup: false,
      learn: false,
      existing: false,
      e1: 1,
      mnemonic: ""
      
    };
  }, 
   methods: {
    // mnemonicGenerate() {
		//	const mnemonic = bip39.generateMnemonic()
		//	this.mnemonic = mnemonic
		//},
       
        loadContent(){

      this.$store.dispatch("setBuyItemList");
 
      //this.dialog = false;
    },
     onSignIn(){
      if (this.$store.state.client != null) { this.dialog = false

      }else{alert("Sign in unsuccessfull")};
    },
    updateDialog(){this.dialog = false}
    },
   
  };

  </script>
<!--
<style scoped>
  .dropdown__textarea {
	display: block;
	height: 8em;
	position: relative;
}
.dropdown__textarea__input {
	border: 0;
	border-radius: 0.5rem;
	padding: 0.5rem;
	box-shadow: inset 0 0 0 1px rgba(0, 0, 0, 0.1);
	outline: none;
	width: 100%;
	height: 100%;
	box-sizing: border-box;
	resize: none;
	letter-spacing: 0.02em;
	font-family: inherit;
	font-size: 0.85rem;
	line-height: 1.4;
	color: rgba(0, 0, 0, 0.85);
}
.dropdown__textarea__icon {
	width: 1rem;
	height: 1rem;
	position: absolute;
	right: 0;
	bottom: 0;
	padding: 1rem;
	fill: rgba(0, 0, 0, 0.5);
	cursor: pointer;
}

</style>-->