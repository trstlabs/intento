<template>
  <div>
   
     
          
          <div >
            <span class="caption">
              You have
              <span
                class="coin__amount"
               v-for="b in balances" :key="b.denom"
               
                >{{numberFormat( b.amount )}} {{ b.denom }}</span
              >
              on your balance.
            </span>
          </div>
          
        </div>
      
  
 
</template>

<style scoped>






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
   balances() {
      //console.log(this.$store.state.bankBalances)
			return this.$store.getters.bankBalances;
  },
  
},
methods: {
   numberFormat(number) {
			return Intl.NumberFormat().format(number)
		},
}
};
</script>