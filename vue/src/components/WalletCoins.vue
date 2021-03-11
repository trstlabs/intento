<template>
  <div>
    <div >
     
      </div>
      <div class="text-center"> 

          </div>
          <v-chip outlined label>
          <div class="card__row">
            <span>
              You have
              <span
                class="coin__amount"
               v-for="b in balances" :key="b.denom"
               
                >{{numberFormat( b.amount )}} {{ b.denom }}</span
              >
              on your balance.
            </span>
          </div>
          </v-chip>
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