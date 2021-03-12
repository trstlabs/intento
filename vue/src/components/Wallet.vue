<template>
  <div>
    <div class="container" >
     
      <div class="overline text-center">{{ address ? "Your Account" : "Sign in" }}</div>
      <div :class="!address && 'card'">
      <div v-if="!address" class="password" >
        <!--<input 
          type="text"
          v-model="password"
          class="password__input"
          placeholder="Password (mnemonic)"
        />-->
        <v-text-field type="password"
          v-model="password"
          
          placeholder="Password (mnemonic)"> </v-text-field>
      </div>
      
  
      <div  v-if="!address" class="password" >
        <v-btn block
          small
          
          @click="signIn"
        >
          Sign in
        </v-btn>
        
      </div>
      <div v-else class="account">
        <div class="card">
          <v-row class="justify-center">
          <v-icon left justify-center large>
        mdi-account-circle
      </v-icon></v-row>
          <div class="card__row">
            
            
            <div class="card__desc caption">
              {{ address }}
            </div>
          </div>
          <div class="card__row caption justify-center">
            <span>
              You have
              <span
                class="coin__amount"
                v-for="b in balances" :key="b.denom"
               
                >{{numberFormat( b.amount )}} {{ b.denom }}</span
              >
              on your balance.
            </span>
          </div> <div  v-if="!!address" >
        <v-btn block text
          small
          
          @click="signOut"
        >
          Sign out
        </v-btn>
        
      </div>
        </div>
       
      </div>
    </div>
   </div>
  </div>
</template>

<style scoped>
.container {
  margin-bottom: 1.5rem;
}
.card {
  background: rgba(0, 0, 0, 0.03);
  border-radius: 0.25rem;
 
  padding: 0.25rem 0.75rem;
  overflow-x: hidden;
}


.card__row {
  display: flex;
  align-items: center;


}

.card__desc {


  word-break: break-all;
}



.coin__amount {
  text-transform: uppercase;

  letter-spacing: 0.02em;
  font-weight: 600;
}
.coin__amount:after {
  content: ",";
  margin-right: 0.25em;
}
.coin__amount:last-child:after {
  content: "";
  margin-right: initial;
}

</style>

<script>

import * as bip39 from "bip39";


export default {
 
 
  data() {
    return {
      password: "",
      error: false,

    
    };
  },
  computed: {
    account() {
      return this.$store.state.account;
    },
    address() {
      const client = this.$store.getters.account
      //console.log(client)
      return client && client.address
    
    },
    mnemonicValid() {
      return bip39.validateMnemonic(this.passwordClean);
    },
    passwordClean() {
      return this.password.trim();
    },

    balances() {
      //console.log(this.$store.state.bankBalances)
			return this.$store.getters.bankBalances;
		}
  },
  methods: {
    async signIn() {
      if (this.mnemonicValid && !this.error) {
        const mnemonic = this.passwordClean;
        await this.$store.dispatch("accountSignIn", { mnemonic })

        this.initConfig();
      }
    },
     async signOut() {
      if (this.address) {
				this.$store.dispatch('accountSignOut')
      }
    },
    numberFormat(number) {
			return Intl.NumberFormat().format(number)
		},
  //set the app according to the logged in user
    async initConfig() {
  
       this.$store.dispatch("setEstimatorItemList", this.address);
       this.$store.dispatch("setToEstimateList", this.address);
       this.$store.dispatch("setCreatorActionList", this.address);
        this.$store.dispatch("setSortedTagList");
         this.$store.dispatch("setCreatorItemList");
          this.$store.dispatch("setBuyerItemList", this.address);

       this.$store.dispatch("setInterestedItemList", this.address);
      this.$emit('signedIn');
      


    }
  },
};
</script>