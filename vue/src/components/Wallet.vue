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
          :class="[
            `button__error__${!!error}`,
            `button__enabled__${!!mnemonicValid}`,
          ]"
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
            
            
            <div class="card__desc">
              {{ address }}
            </div>
          </div>
          <div class="card__row justify-center">
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
  margin: 0.5rem 0;
  
  font-size: 0.875rem;
  font-weight: 400;
  line-height: 1.5;
}
.card__icon {
  width: 1.75rem;
  height: 1.75rem;
  fill: rgba(0, 0, 0, 0.15);
  flex-shrink: 0;
}
.card__desc {
  letter-spacing: 0.02em;
  padding: 0 0.5rem;
  word-break: break-all;
}

.password {
  margin-top: 0.5rem;
  display: flex;
}
.password__input {
  border: none;
  width: 100%;
  padding: 0.75rem;
  box-sizing: border-box;
  font-family: inherit;
  background: rgba(0, 0, 0, 0.03);
  font-size: 0.85rem;
  border-radius: 0.25rem;
 
}
.button {
  margin-left: 1rem;
  background: rgba(0, 0, 0, 0.03);
  padding: 0 1.5rem;
  white-space: nowrap;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 0.85rem;
  color: rgba(0, 0, 0, 0.25);
  text-transform: uppercase;
  letter-spacing: 0.05em;
  border-radius: 0.25rem;
  transition: all 0.1s;
  user-select: none;
}
.button.button__error__true {
  animation: shake 0.82s cubic-bezier(0.36, 0.07, 0.19, 0.97) both;
  background: rgba(255, 228, 228, 0.5);
  color: rgb(255, 0, 0);
}
.button__enabled__false {
  cursor: not-allowed;
}
.button__enabled__true {
  color: rgba(0, 125, 255);
  font-weight: 700;
  cursor: pointer;
}
.button__enabled__true:active {
  color: rgba(0, 125, 255, 0.65);
}
.password__input:focus {
  outline: none;
  background: rgba(0, 0, 0, 0.06);
  color: rgba(0, 0, 0, 0.5);
}
.password__input::placeholder {
  color: rgba(0, 0, 0, 0.35);
  font-weight: 500;
}
.coin__amount {
  text-transform: uppercase;
  font-size: 0.75rem;
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
@keyframes shake {
  10%,
  90% {
    transform: translate3d(-1px, 0, 0);
  }
  20%,
  80% {
    transform: translate3d(2px, 0, 0);
  }
  30%,
  50%,
  70% {
    transform: translate3d(-4px, 0, 0);
  }
  40%,
  60% {
    transform: translate3d(4px, 0, 0);
  }
}
</style>

<script>
import IconUser from "@/components/IconUser.vue";
import * as bip39 from "bip39";


export default {
 
  components: {
    IconUser, 
  },
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