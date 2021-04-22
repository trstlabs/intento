<template>
  <div>
    <div v-if="!sent">
       <v-text-field class="mx-4" v-model="email" required placeholder="Email-address" name="email" type="text" />
     <v-btn @click="sent = true, sendMail()"> Send</v-btn>  </div>
    <div v-else>
      <v-alert dense 
  type="success"  class="caption"
>
Confirm this sign-in once, by clicking the link sent to <span v-if="this.email"> {{email}} </span> <span v-else> your Google account's email </span >, on this device.
</v-alert> </div>
  </div>
</template>

<script>

import { auth } from "./firebase/db.js";



export default {

  data() {
    return {

sent: false,
    email: "",


      
    };
  },
  /*created() {
  
        let email = window.localStorage.getItem("emailForSignIn");
        if (email) {
          // User opened the link on a different device. To prevent session fixation
          // attacks, ask the user to provide the associated email again. For example:
          

        // The client SDK will parse the code from the link for you.
        auth
          .signInWithEmailLink(email, window.location.href)
          .then((result) => {
            // Clear email from storage.
            window.localStorage.removeItem('emailForSignIn');
            //window.localStorage.setItem("emailRef", window.location.href);
            this.$store.commit("set", {key: "user", value: result.user } );
            console.log(this.$store.state.user);
            // You can access the new user via result.user
            // Additional user info profile not available via:
            // result.additionalUserInfo.profile == null
            // You can check if the user is new or existing:
            // result.additionalUserInfo.isNewUser
          })
          .catch((error) => {
            console.log(error)
            // Some error occurred, you can inspect the code: error.code
            // Common errors could be invalid email and invalid or expired OTPs.
          });
        }
    
},*/

  computed: {
    
    /*email() {
      //console.log(localStorage.getItem("privkey"));
      return localStorage.getItem("emailForSignIn");
    },*/
  },
  methods: {
    async sendMail() {
      this.loading = true;
      try {
        console.log(this.email)
 let actionCodeSettings = {
          url: "https://marketplace.trustpriceprotocol.com",
          // This must be true.
          handleCodeInApp: true,
          /*iOS: {
    bundleId: 'com.example.ios'
  },
  android: {
    packageName: 'com.example.android',
    installApp: true,
    minimumVersion: '12'
  },*/
          //dynamicLinkDomain: 'marketplace.trustpriceprotocol.com.page.link'
        };

        auth
          .sendSignInLinkToEmail(
            this.email,
            actionCodeSettings
          )
          .then(() => {
            // The link was successfully sent. Inform the user.
            // Save the email locally so you don't need to ask the user for it again
            // if they open the link on the same device.
          
            window.localStorage.setItem("emailForSignIn", this.email);
            //alert("Confirm by clicking the email link on your device")
            // ...
          })
          .catch((error) => {
            var errorCode = error.code;
            var errorMessage = error.message;
            // ...
          });
      } catch (error) {
        this.loading = false;
        console.error(error, "caught");
      }
    },
   
  },
  

};
</script>

