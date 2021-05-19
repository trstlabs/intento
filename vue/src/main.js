import Vue from "vue";
import vuetify from '@/plugins/vuetify' // path to vuetify export
import App from "./App.vue";
import router from "./router";
import store from "./store";
//import _ from "lodash";
//console.log(require('dotenv').config())
import {firestorePlugin} from "vuefire";

//Vue.use(Vuetify);

Vue.config.productionTip = true;
Vue.config.devtools = false
Vue.config.debug = false;
Vue.config.silent = true;

Vue.use(firestorePlugin);

//Object.defineProperty(Vue.prototype, "$lodash", { value: _ });

const ComponentContext = require.context("./", true, /\.vue$/i, "lazy");
ComponentContext.keys().forEach((componentFilePath) => {
  const componentName = componentFilePath.split("/").pop().split(".")[0];
  Vue.component(componentName, () => ComponentContext(componentFilePath));
});







// Initialize Firebase
//firebase.analytics();

new Vue({
  router,
  store,
  vuetify,
  render: (h) => h(App),



}).$mount("#app"); 