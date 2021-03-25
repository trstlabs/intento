<template>
  <div>
    <v-dialog v-model="dialog" width="600" persistent>
      <v-card class="text-center">
        <img class="pa-2" src="img/brand/icon.png" width="77" />
       

        <div v-if="!signin && !signup && !learn" >
           <p class="text-center overline" >
          Continue with an account 
        </p>
       
          <v-card-text class="">
          Welcome to the TPP marketplace. Feel free to use the app and see how it
            behaves. Note: items are examples only and are currently not redeemable for TPP tokens.
            Feedback is always welcome at the
            <a  target="_blank" href="https://www.trustpriceprotocol.com"> main page</a>.
          </v-card-text>
     
             <v-btn class="ma-6" large
              color="primary"
              
              @click="(signin = true), loadContent()"
            >
              Sign In
            </v-btn>


              

          

          <v-card-actions>
            <v-btn color="primary " text @click="(learn = true), loadContent()">
              Learn more
            </v-btn>
            <v-btn
              color="primary lighten-2"
              text
              @click="(dialog = false), loadContent()"
            >
              Look Around
            </v-btn>
            <v-spacer />

            <v-spacer />
            <v-btn
              color="primary "
              text
              @click="(signup = true), loadContent()"
            >
              Sign Up
            </v-btn>
          </v-card-actions>
        </div>
        <div v-if="signin">
           <p class="text-center overline" >
          Continue with an account 
        </p>
          <torus v-if="!wallet" :privkey="this.signkey"/>
           <v-row v-if="!wallet">
              <v-divider class="ma-2" />
              <p class="caption">Or</p>
              <v-divider class="ma-2" />
            </v-row>
            <v-btn text v-if="!wallet" @click="wallet = true">
              Sign in with a mnemonic phrase
            </v-btn>
          <wallet v-if="wallet" @signedIn="updateDialog()" />
          <v-card-actions>
            <v-col>
              <v-btn color="primary" text @click="signin = false, wallet = false"> Back </v-btn>
            </v-col>
            
          </v-card-actions>
        </div>
        

        <div v-if="signup">
          <v-card-text>
            Receive 5 free tokens to get started! Create a new account linked to
            Google using DirectAuth.
          </v-card-text>
          <div v-if="!existing">
            <faucet-torus />
            <v-row>
              <v-divider class="ma-2" />
              <p class="caption">Or</p>
              <v-divider class="ma-2" />
            </v-row>
            <v-btn text @click="existing = true">
              Sign up with an existing cosmos-address.
            </v-btn>
          </div>
          <div v-if="existing">
            <faucet />
          </div>
          <v-card-actions>
            <v-col>
              <v-btn
                block
                color="primary"
                text
                @click="(signup = false), (existing = false)"
              >
                Back
              </v-btn></v-col
            ><v-col>
              <v-btn
                color="primary"
                block
                text
                @click="(signin = true), (signup = false), loadContent()"
              >
                Sign In
              </v-btn></v-col
            >
          </v-card-actions>
        </div>

        <div v-if="learn">
          <v-stepper v-model="e1">
            <v-stepper-header>
              <v-stepper-step :complete="e1 > 1" step="1">
                Shop
              </v-stepper-step>

              <v-divider></v-divider>

              <v-stepper-step :complete="e1 > 2" step="2">
                Sell
              </v-stepper-step>

              <v-divider></v-divider>

              <v-stepper-step step="3"> Earn </v-stepper-step>
            </v-stepper-header>

            <v-stepper-items>
              <v-stepper-content step="1">
                <v-card class="mb-6">
                  <v-card-title 
                    >There is a problem with current online marketplaces. </v-card-title
                  ><v-card-text>
                    <p>
                      Ever wanted to buy something and the item was already
                      granted to another user?
                    </p>
                    <p>
                      Ever paid too much because you had a
                      wrong idea about the item?
                    </p>
                    <p>
                      Trust price protocol is a place where you can trade items hassle free.
                    </p>

                  <v-card class="ma-4">
                     <v-list-item>
      <v-list-item-content>
        <v-list-item-title>When you provide prepayment, you are
                     the only buyer</v-list-item-title><v-list-item-subtitle>
          and you earn a cashback reward of ±5% after the transfer.
        </v-list-item-subtitle>
      </v-list-item-content><v-list-item-icon>
            <v-icon>mdi-plus </v-icon>
          </v-list-item-icon>
    </v-list-item>

    <v-list-item two-line>
      <v-list-item-content>
        <v-list-item-title>Prices are made from independent estimations</v-list-item-title>
        <v-list-item-subtitle>So you
                      know what you buy is right.</v-list-item-subtitle>
      </v-list-item-content><v-list-item-icon>
            <v-icon>mdi-plus </v-icon>
          </v-list-item-icon>
    </v-list-item>

    <v-list-item three-line>
      <v-list-item-content>
        <v-list-item-title>You are in control of your prepayment </v-list-item-title>
        <v-list-item-subtitle>
          and can always get
                      it back until it is transferred.
        </v-list-item-subtitle>
        
        
      </v-list-item-content><v-list-item-icon>
            <v-icon>mdi-plus </v-icon>
          </v-list-item-icon>
    </v-list-item>

                 

                   </v-card>
                
                   

                    <p class="font-weight-medium ma-4">A great way to spend crypto on things you like</p>
                  </v-card-text></v-card
                >
                <v-row class="ma-2">
                  <v-btn text @click="learn = false"> Back </v-btn>
                  <v-spacer />
                  <v-btn color="primary" @click="e1 = 2"> Continue </v-btn>
                </v-row>
              </v-stepper-content>

              <v-stepper-content step="2">
                <v-card class="mb-12"
                  ><v-card-title
                    >An opportunity to enter the crypto universe
                  </v-card-title>
                  <v-card-text>
                    
                    <p>
                      As a seller, you are free to choose 2 options. Ship the
                      item and/or choose “local pickup”.
                    </p>
                    <p>
                      If you choose shipping, you can charge shipping costs and
                      these are separate from the item price.
                    </p>
                    <p>
                      You can set the accuracy of the price. This is set to 3
                      estimations by default. Setting the accuracy higher will
                      result in a better accuracy but this will also take
                      longer. When a final price is made, you may always decline
                      or accept it.
                    </p>
                    <p>
                      You don’t have to worry about pricing, just provide good
                      quality pictures and information and you are done.
                    </p>

                    <p class="font-weight-medium ma-4">
                      Selling items has never been this simple.
                    </p></v-card-text
                  ></v-card
                >
                <v-row class="ma-2">
                  <v-btn text @click="e1 = 1"> Back </v-btn>
                  <v-spacer />
                  <v-btn color="primary" @click="e1 = 3"> Continue </v-btn>
                </v-row>
              </v-stepper-content>

              <v-stepper-content step="3">
                <v-card class="mb-12"
                  ><v-card-title
                    >Trade your time on the internet for crypto
                  </v-card-title>
                  <v-card-text>
                    <p>
                      Do you enjoy looking around the internet and
                      finding/comparing prices? Or looking at things you cannot
                      afford just for fun? Do you want the best value for your
                      money? Then you may enjoy becoming an “estimator”
                    </p>
                    <p>
                      As an “estimator” your job is to moderate items before
                      they enter the marketplace. And you will be rewarded
                      accordingly!
                    </p>
                    <p>
                      Currently, you will earn roughly 5% of the final selling
                      price if you are the estimator closest to the final price. These TPP coins are minted.
                    </p>
                    <p>
                      To cope with bad acting, there is a deposit
                      required for each estimation.
                      These will be returned, except when the following occurs
                      1). The seller did not accept the final price and you are
                      the lowest estimator. 2) The buyer ended up withdrawing
                      prepayment and you are the highest estimator. 
                    </p>

                    <p class="font-weight-medium ma-4">You can trade your earned TPP for
                      the items you like.
                      As a bonus, you can <v-icon small>mdi-heart </v-icon> items and view them once they hit
                      the marketplace. 
                    </p>
                  </v-card-text>
                </v-card>
                <v-row class="ma-2">
                  <v-btn text @click="e1 = 2"> Back </v-btn>
                  <v-spacer />
                  <v-btn
                    text
                    color="primary"
                    @click="(learn = false), (dialog = false)"
                  >
                    Look Around
                  </v-btn>
                  <v-spacer />
                  <v-btn
                    color="primary"
                    @click="(learn = false), (signup = true)"
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
import TorusPlaceholder from './TorusPlaceholder.vue';


/*
import Faucet from "./Faucet.vue";
import FaucetTorus from "./FaucetTorus.vue";
import Torus from "./Torus.vue";
import Wallet from "./Wallet.vue";*/

//const Faucet = () => import("./Faucet.vue");

//import * as bip39 from 'bip39'
export default {
  //components: { Wallet, Faucet, Torus, FaucetTorus },
  components: {TorusPlaceholder },

  data() {
    return {
   
      //dismiss: false,
      login: false,
      dialog: true,
      signin: false,
      signup: false,
      learn: false,
      existing: false,
      e1: 1,
      mnemonic: "",
      wallet: false,
    };
  },
  computed: {
    signkey() {
      //console.log(localStorage.getItem("privkey"));
      return localStorage.getItem("privkey");
    },
  },
 
  methods: {
    // mnemonicGenerate() {
    //	const mnemonic = bip39.generateMnemonic()
    //	this.mnemonic = mnemonic
    //},

    loadContent() {
      this.$store.dispatch("setBuyItemList");

      //this.dialog = false;
    },
    onSignIn() {
      if (this.$store.state.client != null) {
        this.dialog = false;
      } else {
        alert("Sign in unsuccessfull");
      }
    },
    updateDialog() {
      this.dialog = false;
    },
  },
};
</script>
