<template>
  <div>
    <div></div>
    <div v-if="!this.$store.state.account.address" class="justify-center">
      <v-col>
        <button>
          <v-progress-linear indeterminate v-if="loading"></v-progress-linear>
          <v-img
            max-height="100"
            max-width="200"
            src="img/google/btn.png"
            @click="login"
          >
          </v-img>
        </button>
      </v-col>
      <v-col class="mx-auto"
        ><button>
          <a href="https://tor.us/" target="_blank" rel="noopener">
            <v-img
              max-height="80"
              max-width="120"
              src="img/google/directauth.png"
            >
            </v-img
          ></a></button
      ></v-col>
    </div>
     <div v-else>    <v-alert  
  type="success"  class="caption"
>
Confirm this sign-in once, by clicking the link sent to <span v-if="!!this.email"> {{email}} </span> <span v-else> your Google account's email </span >, on this device.
</v-alert> </div>
  </div>
</template>

<script>
import axios from "axios";
import TorusSdk from "@toruslabs/torus-direct-web-sdk";
import { SigningStargateClient } from "@cosmjs/stargate";
import { DirectSecp256k1Wallet } from "@cosmjs/proto-signing/";
import { fromHex } from "@cosmjs/encoding";
import { auth } from "./firebase/db.js";

const GOOGLE = "google";

export default {
  props: ["privkey"],
  data() {
    return {
      loading: false,
      torusdirectsdk: undefined,
      selectedVerifier: "google",
      loginHint: "",
      verifierMap: {
        [GOOGLE]: {
          name: "Google",
          typeOfLogin: "google",
          clientId:
            "29876044321-2u525qqtirvtd9v4camu5jpm0egf3al4.apps.googleusercontent.com",
          verifier: "trustpriceprotocol-google-testnet",
        },
      },
    };
  },
  created() {
     const email = window.localStorage.getItem("emailForSignIn");
    if (email) {
      //var email = window.localStorage.getItem("emailForSignIn");
      //var emailRef = window.localStorage.getItem("emailRef");
      //if this doesnt work for returning user, put emailRef outside of it and  try signing in right away. First test if auth is able to handle handle it from its localstorage
      //if (emailRef) {
        //auth.signInWithEmailLink(email, emailRef).then((result) => {
          //this.$store.commit("set", {key: "user", value: result.user } );
         // console.log(this.$store.state.user);
       // });

       //https://stackoverflow.com/questions/42878179/how-to-persist-a-firebase-login
       //https://www.youtube.com/watch?v=5VxqV8FhlVg
      auth.onAuthStateChanged(user => {
  if (user){ 
    this.$store.commit("set", {key: "user", value: user } );
            console.log(this.$store.state.user); }} )
      // Confirm the link is a sign-in with email link.
      if (auth.isSignInWithEmailLink(window.location.href)) {
        // Additional state parameters can also be passed via URL.
        // This can be used to continue the user's intended action before triggering
        // the sign-in operation.
        // Get the email if available. This should be available if the user completes
        // the flow on the same device where they started it.

        //var emailRef = window.localStorage.getItem('emailRef');
       
        if (!email) {
          // User opened the link on a different device. To prevent session fixation
          // attacks, ask the user to provide the associated email again. For example:
          email = window.prompt("Please provide your email for sign-in confirmation");
        }

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
      if(!this.privkey){
      this.torusSignIn(this.privkey);}
    }
},

  computed: {
    loginToConnectionMap() {
      return {
        // [GOOGLE]: { login_hint: 'hello@tor.us', prompt: 'none' }, // This allows seamless login with google
      };
    },
    email() {
      //console.log(localStorage.getItem("privkey"));
      return localStorage.getItem("emailForSignIn");
    },
  },
  methods: {
    async login(hash, queryParameters) {
      this.loading = true;
      try {
        if (!this.torusdirectsdk) return;
        const jwtParams = this.loginToConnectionMap[this.selectedVerifier] || {
        
        };
  
        const { typeOfLogin, clientId, verifier } = this.verifierMap[
          this.selectedVerifier
        ];
        console.log(
          hash,
          queryParameters,
          typeOfLogin,
          clientId,
          verifier,
          jwtParams
        );
        const loginDetails = await this.torusdirectsdk.triggerLogin({
          typeOfLogin,
          verifier,
          clientId,
          jwtParams,
          hash,
          queryParameters,
        });

        // const loginDetails = await this.torusdirectsdk.triggerHybridAggregateLogin({
        //   singleLogin: {
        //     typeOfLogin,
        //     verifier,
        //     clientId,
        //     jwtParams,
        //     hash,
        //     queryParameters,
        //   },
        //   aggregateLoginParams: {
        //     aggregateVerifierType: "single_id_verifier",
        //     verifierIdentifier: "tkey-google",
        //     subVerifierDetailsArray: [
        //       {
        //         clientId: "221898609709-obfn3p63741l5333093430j3qeiinaa8.apps.googleusercontent.com",
        //         typeOfLogin: "google",
        //         verifier: "torus",
        //       },
        //     ],
        //   },
        // });

        // AGGREGATE LOGIN
        // const loginDetails = await this.torusdirectsdk.triggerAggregateLogin({
        //   aggregateVerifierType: "single_id_verifier",
        //   verifierIdentifier: "tkey-google",
        //   subVerifierDetailsArray: [
        //     {
        //       clientId: "221898609709-obfn3p63741l5333093430j3qeiinaa8.apps.googleusercontent.com",
        //       typeOfLogin: "google",
        //       verifier: "torus"
        //     }
        //   ]
        // });

        // AGGREGATE LOGIN - AUTH0 (Not working - Sample only)
        // const loginDetails = await torusdirectsdk.triggerAggregateLogin({
        //   aggregateVerifierType: "single_id_verifier",
        //   verifierIdentifier: "google-auth0-gooddollar",
        //   subVerifierDetailsArray: [
        //     {
        //       clientId: config.auth0ClientId,
        //       typeOfLogin: "email_password",
        //       verifier: "auth0",
        //       jwtParams: { domain: config.auth0Domain },
        //     },
        //   ],
        // });
        //this.console(loginDetails);
        // console.log("tasdf" + loginDetails)

        this.torusSignIn(loginDetails.privateKey);

        // console.log(loginDetails.privateKey)

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
            loginDetails.userInfo.email,
            actionCodeSettings
          )
          .then(() => {
            // The link was successfully sent. Inform the user.
            // Save the email locally so you don't need to ask the user for it again
            // if they open the link on the same device.
            console.log(loginDetails.userInfo.email)
            window.localStorage.setItem("emailForSignIn", loginDetails.userInfo.email);
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
    async torusSignIn(details) {
      var uint8array = new TextEncoder().encode(details);
      console.log(details);
      console.log(uint8array);

      const wallet = await DirectSecp256k1Wallet.fromKey(
        fromHex(details),
        "cosmos"
      );
      this.$store.commit("set", { key: "wallet", value: wallet });
      localStorage.setItem("privkey", details);
      //console.log(localStorage.getItem('privkey'))
      const { address } = wallet;

      const url = `${process.env.VUE_APP_API}/auth/accounts/${address}`;
      const acc = (await axios.get(url)).data;
      const account = acc.result.value;

      this.$store.commit("set", { key: "account", value: account });

      //console.log("fdgadagfgfd" + SigningStargateClient.connectWithSigner());
      ////onsole.log(RPC)
      const client = await SigningStargateClient.connectWithSigner(
        process.env.VUE_APP_RPC,
        wallet,
        {}
      );
      this.$store.commit("set", { key: "client", value: client });
      const type = { type: "estimator" };
      await this.$store.dispatch("entityFetch", type);

      this.$store.dispatch("setSellerItemList", account.address);
      this.$store.dispatch("setBuyItemList");
      this.loading = false;
      //console.log(client)
      try {
        await this.$store.dispatch("bankBalancesGet");
      } catch {
        console.log("Error in getting a bank balance.");
      }
    },

    handleRedirectParameters(hash, queryParameters) {
      const hashParameters = hash.split("&").reduce((result, item) => {
        const [part0, part1] = item.split("=");
        result[part0] = part1;
        return result;
      }, {});
      //console.log(hashParameters, queryParameters);
      let instanceParameters = {};
      let error = "";
      if (!queryParameters.preopenInstanceId) {
        if (Object.keys(hashParameters).length > 0 && hashParameters.state) {
          instanceParameters =
            JSON.parse(
              atob(decodeURIComponent(decodeURIComponent(hashParameters.state)))
            ) || {};
          error =
            hashParameters.error_description || hashParameters.error || error;
        } else if (
          Object.keys(queryParameters).length > 0 &&
          queryParameters.state
        ) {
          instanceParameters =
            JSON.parse(
              atob(
                decodeURIComponent(decodeURIComponent(queryParameters.state))
              )
            ) || {};
          if (queryParameters.error) error = queryParameters.error;
        }
      }
      return { error, instanceParameters, hashParameters };
    },
  },
  async mounted() {
    try {
      var url = new URL(location.href);
      const hash = url.hash.substr(1);
      const queryParams = {};
      for (let key of url.searchParams.keys()) {
        queryParams[key] = url.searchParams.get(key);
      }
      const { error, instanceParameters } = this.handleRedirectParameters(
        hash,
        queryParams
      );
      const torusdirectsdk = new TorusSdk({
        baseUrl: `${location.origin}/serviceworker`,
        enableLogging: true,
        network: "testnet", // details for test net
      });

      await torusdirectsdk.init({ skipSw: false });
      this.torusdirectsdk = torusdirectsdk;
      if (hash) {
        if (error) throw new Error(error);
        const { verifier: returnedVerifier } = instanceParameters;
        this.selectedVerifier = Object.keys(this.verifierMap).find(
          (x) => this.verifierMap[x].verifier === returnedVerifier
        );
        this.login(hash, queryParams);
      }
    } catch (error) {
      console.error(error, "mounted caught");
    }
  },
};
</script>

