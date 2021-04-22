
<template>
<div>
  <div class="pa-2 mx-auto"> 
  
  <form @submit.prevent="submit">
    
       <div v-if="status">Status: {{status}}</div>
      <div v-if="serverError">Oops Error!{{serverError}}</div>

      <v-text-field class="mx-4" v-model="address" required placeholder="cosmos-address" name="address" type="text" />
      <button :disabled="status==='Registering...'" type="submit" class="button"></button>
      
     <v-row v-if="this.$store.state.wallet" class="justify-center mb-4">
      <vue-recaptcha v-if="status == '' || status == 'Registering...' "
        ref="recaptcha"
        @verify="onCaptchaVerified"
        @expired="onCaptchaExpired"
        @error="onCaptchaError"
        @render="onCaptchaRender"
       
        :sitekey="google">
      </vue-recaptcha></v-row>
      <v-alert type="success" v-if="sucessfulServerResponse">{{sucessfulServerResponse}}</v-alert>
    </form>
    <v-btn color="primary" block @click="submit()">Submit</v-btn> </div>
    <v-divider/>
</div>
</template>

<script>
import VueRecaptcha from 'vue-recaptcha'
import axios from 'axios'
export default {

  data () {
    return {
     
      status: '',
      address: null,
      sucessfulServerResponse: '',
      serverError: '',
      google: '6LdzO1waAAAAAKDkD1sNFSx552KIrXIr1F_NY_4O'
    }
  },
  methods: {
   submit: function () {

      
      // console.log(this.$refs.recaptcha.execute())
      this.status = 'Registering...'
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
      let accountQuery = await axios.get('https://node.trustpriceprotocol.com/auth/accounts/' + this.address)

      //console.log(accountQuery.data.result.value.address)
      if (!accountQuery.data.result.value.address) {
        //console.log("letsgo")
 
      this.status = 'Submitting...'
      this.$refs.recaptcha.reset()
      try {
        this.status = 'Getting TPP tokens'
        let response = await axios.post('/.netlify/functions/faucet', {
          recipient:  this.$store.state.wallet.address,
          recaptchaToken: recaptchaToken
        })
        if (response.status === 200) {
           
          this.sucessfulServerResponse = 'Your cosmos-address is succesfully registered!'
          alert('Sign up successfull')
          window.location.reload()
        
        
        }
         else {
          this.sucessfulServerResponse = response.data
        }
      } catch (err) {
        console.log("ERROR" + err)
        //alert("Error receiving TPP tokens on this address")
        window.location.reload()
        //let foo = getErrorMessage(err)
        //this.serverError = foo === '"read ECONNRESET"' ? 'Opps, we had a connection issue, please try again' : foo
      }
      this.status = ''}
      else {alert("Account already registered on TPP, please sign in instead or create a new account")}
    },
    
    onCaptchaExpired: function () {
      this.status = ''
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