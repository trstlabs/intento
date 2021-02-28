import Vue from "vue";
import vuetify from '@/plugins/vuetify' // path to vuetify export
import App from "./App.vue";
import router from "./router";
import store from "./store";
import _ from "lodash";
//import firebase from './firebase';

//Vue.use(Vuetify);

Vue.config.productionTip = false;

Object.defineProperty(Vue.prototype, "$lodash", { value: _ });

const ComponentContext = require.context("./", true, /\.vue$/i, "lazy");
ComponentContext.keys().forEach((componentFilePath) => {
  const componentName = componentFilePath.split("/").pop().split(".")[0];
  Vue.component(componentName, () => ComponentContext(componentFilePath));
});

var firebaseConfig = {
  apiKey: process.env.VUE_APP_API_KEY,
  authDomain: "trustitems-cbb92.firebaseapp.com",
  databaseURL: "https://trustitems-cbb92.firebaseio.com",
  projectId: "trustitems-cbb92",
  storageBucket: "trustitems-cbb92.appspot.com",
  messagingSenderId: "1033037336313",
  appId: "1:1033037336313:web:475aaf4e502bc36adf9c28",
  measurementId: "G-XQ9H3T2QNM",
};



// Initialize Firebase
firebase.initializeApp(firebaseConfig);
firebase.analytics();

new Vue({
  router,
  store,
  vuetify,
  render: (h) => h(App),
  /*created() {
    // Prevent blank screen in Electron builds[deletelater]
    this.$router.push('/')
  }*/

  //created() {



}).$mount("#app"); 