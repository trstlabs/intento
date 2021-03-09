<template>
  <div >
    <div>
    </div>
    <div class="justify-center"><v-col>
   <button> <v-img
    max-height="100"
    max-width="200" 
        
    src="img/google/btn.png"
  @click="login">
  </v-img></button>
  </v-col>
   <v-col class="mx-auto"><button>
  <a  href="https://tor.us/" target="_blank" rel="noopener" > <v-img 
   max-height="60"
    max-width="100" 
         
    src="img/google/directauth.png"
   >
  </v-img></a> </button></v-col>

  </div></div>
   
   
</template>

<script>
import TorusSdk from "@toruslabs/torus-direct-web-sdk";

const GOOGLE = "google";


export default {
  name: "App",
  data() {
    return {
      torusdirectsdk: undefined,
      selectedVerifier: "google",
      loginHint: "",
      verifierMap: {
        [GOOGLE]: {
          name: "Google",
          typeOfLogin: "google",
          clientId: "29876044321-2u525qqtirvtd9v4camu5jpm0egf3al4.apps.googleusercontent.com",
          verifier: "trustpriceprotocol-google-testnet",
        },
       
      },
    };
  },
  computed: {
    loginToConnectionMap() {
      return {
        // [GOOGLE]: { login_hint: 'hello@tor.us', prompt: 'none' }, // This allows seamless login with google
        
      };
    },
  },
  methods: {
    async login(hash, queryParameters) {
      try {
        if (!this.torusdirectsdk) return;
        const jwtParams = this.loginToConnectionMap[this.selectedVerifier] || {};
        const { typeOfLogin, clientId, verifier } = this.verifierMap[this.selectedVerifier];
        console.log(hash, queryParameters, typeOfLogin, clientId, verifier, jwtParams);
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
       // console.log(loginDetails)
       // console.log(loginDetails.privateKey)
       
        this.$store.dispatch("torusSignIn", loginDetails.privateKey)
      } catch (error) {
        console.error(error, "caught");
      }
    },
    console(text) {
      document.querySelector("#console>p").innerHTML = typeof text === "object" ? JSON.stringify(text) : text;
    },
    handleRedirectParameters(hash, queryParameters) {
      const hashParameters = hash.split("&").reduce((result, item) => {
        const [part0, part1] = item.split("=");
        result[part0] = part1;
        return result;
      }, {});
      console.log(hashParameters, queryParameters);
      let instanceParameters = {};
      let error = "";
      if (!queryParameters.preopenInstanceId) {
        if (Object.keys(hashParameters).length > 0 && hashParameters.state) {
          instanceParameters = JSON.parse(atob(decodeURIComponent(decodeURIComponent(hashParameters.state)))) || {};
          error = hashParameters.error_description || hashParameters.error || error;
        } else if (Object.keys(queryParameters).length > 0 && queryParameters.state) {
          instanceParameters = JSON.parse(atob(decodeURIComponent(decodeURIComponent(queryParameters.state)))) || {};
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
      const { error, instanceParameters } = this.handleRedirectParameters(hash, queryParams);
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
        this.selectedVerifier = Object.keys(this.verifierMap).find((x) => this.verifierMap[x].verifier === returnedVerifier);
        this.login(hash, queryParams);
      }
    } catch (error) {
      console.error(error, "mounted caught");
    }
  },
};
</script>

<style>

#console {
  border: 1px solid black;
  height: 80px;
  padding: 2px;
  bottom: 10px;
  position: absolute;
  text-align: left;
  width: calc(100% - 20px);
  border-radius: 5px;
}
#console::before {
  content: "Console :";
  position: absolute;
  top: -20px;
  font-size: 12px;
}
#console > p {
  margin: 0.5em;
  word-wrap: break-word;
}
</style>