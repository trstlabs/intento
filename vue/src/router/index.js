import Vue from "vue";
import VueRouter from "vue-router";
import Index from "../views/Index.vue";
import Buy from "../views/Buy.vue";
import Sell from "../views/Sell.vue";
import Estimate from "../views/Estimate.vue";
import Account from "../views/Account.vue";
import BuyItemDetails from "../views/BuyItemDetails.vue";
import Messages from "../views/Messages.vue";

import MissingPage from "../views/MissingPage.vue";

Vue.use(VueRouter);




const routes = [
  {
    path: "/",
    name: "Buy",
    component: Buy,
    meta: {
      title: 'Buy - Marketplace',
      metaTags: [
        {
          name: 'description',
          content: 'Buy and sell items easily.'
        },
        {
          property: 'og:description',
          content: 'Buy and sell items easily.'
        }
      ]
    }
  
    
  },

  {
    path: "/buy/:id",
    name: "BuyItemDetails",
    component: BuyItemDetails
  },
  {
    path: "/sell",
    component: Sell,
    meta: {
      title: 'Sell - Marketplace'},
    
   
    
  
  },
  {
    path: "/earn",
    component: Estimate,
    meta: {
      title: 'Earn - Marketplace'},
  },
  
  {
    path: "/account",
    component: Account,
    meta: {
      title: 'Account - Marketplace'},
  },

  {
    path: "/messages",
    component: Messages,
    meta: {
      title: 'Messages - Marketplace'},
  },
  {
    path: '*',
    name: 'catchAll',
    component: Buy
  }
  /*{
    path: '*',
    component: MissingPage
  }*/
];


const router = new VueRouter({
  mode: "history",
  //base: process.env.BASE_URL,
  routes,
  
});


const DEFAULT_TITLE = 'Marketplace';
router.afterEach((to, from) => {
    // Use next tick to handle router history correctly
    // see: https://github.com/vuejs/vue-router/issues/914#issuecomment-384477609
    Vue.nextTick(() => {
        document.title = to.meta.title || DEFAULT_TITLE;
    });
});

export default router;
