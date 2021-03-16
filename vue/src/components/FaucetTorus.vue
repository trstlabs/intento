
<template>
<div>
  <div class="pa-2 mx-auto"> 
  
  <form @submit.prevent="submit">
    
       <div v-if="status">Status: {{status}}</div>
      <div v-if="serverError">Oops Error!{{serverError}}</div>

     <torus v-if="!this.$store.state.wallet"/>
      <button :disabled="status==='submitting'" type="submit" class="button"></button>
      
     <v-row v-if="!!this.$store.state.wallet" class="justify-center mb-4">
      <vue-recaptcha 
        ref="recaptcha"
        @verify="onCaptchaVerified"
        @expired="onCaptchaExpired"
        @error="onCaptchaError"
        @render="onCaptchaRender"
       
        :sitekey="google">
      </vue-recaptcha></v-row>
      <v-alert type="success" v-if="sucessfulServerResponse">{{sucessfulServerResponse}}</v-alert>
    </form>
  <!--<v-btn class="ma-2" color="primary" block @click="submit()">Receive tokens</v-btn>--> </div>
  
</div>
</template>

<script>
import VueRecaptcha from 'vue-recaptcha'
import axios from 'axios'
export default {

  data () {
    return {
     
      status: null,
     
      sucessfulServerResponse: '',
      serverError: '',
      google: '6LdzO1waAAAAAKDkD1sNFSx552KIrXIr1F_NY_4O'
    }
  },
  methods: {
    submit: function () {

      
      // console.log(this.$refs.recaptcha.execute())
      this.status = 'submitting'
      // this.$refs.recaptcha.reset()
      this.$refs.recaptcha.execute() 
    },
    onCaptchaRender: function(id) {
      //console.log({id})
    },
    onCaptchaError: function(error) {
      console.log({error})
    },
    onCaptchaVerified: async function (recaptchaToken) {
      let accountQuery = await axios.get('https://node.trustpriceprotocol.com/auth/accounts/' + this.$store.state.wallet.address)

      console.log(accountQuery.data.result.value.address)
      if (!accountQuery.data.result.value.address) {
        //console.log("letsgo")
      const self = this
      self.status = 'submitting'
      self.$refs.recaptcha.reset()
      try {
        let response = await axios.post('/.netlify/functions/faucet', {
          recipient:  this.$store.state.wallet.address,
          recaptchaToken: recaptchaToken
        })
        if (response.status === 200) {
          self.sucessfulServerResponse = 'Your cosmos-address is succesfully registered!'
        
        }
         else {
          self.sucessfulServerResponse = response.data
        }
      } catch (err) {
        console.log("ERROR" + err)
        //let foo = getErrorMessage(err)
        //self.serverError = foo === '"read ECONNRESET"' ? 'Opps, we had a connection issue, please try again' : foo
      }
      self.status = ''}
      else {alert("Account already registered on TPP, please log in instead or create a new account")}
    },
    
    onCaptchaExpired: function () {
      self.status = ''
      this.$refs.recaptcha.reset()
    },
    getErrorMessage (err) {
  let responseBody
  responseBody = err.response
  if (!responseBody) {
    responseBody = err
  } else {
    responseBody = err.response.data || responseBody
  }
  return responseBody.message || JSON.stringify(responseBody)
},
  
  },

  
  components: {
    VueRecaptcha
  }
}
</script>