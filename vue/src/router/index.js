import Vue from "vue";
import VueRouter from "vue-router";
//import Index from "../views/Index.vue";
//import Buy from "../views/Buy.vue";
//import Sell from "../views/Sell.vue";
//import Estimate from "../views/Estimate.vue";
//import Account from "../views/Account.vue";
//import BuyItemDetails from "../views/BuyItemDetails.vue";
//import Messages from "../views/Messages.vue";

//import MissingPage from "../views/MissingPage.vue";

Vue.use(VueRouter);



const routes = [
  {
    path: "/",
    name: "Buy",
    component: () => 
        import(/* webpackChunkName: "buy" */ '@/views/Buy'), 
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
    path: "/itemid=:id",
    name: "BuyItemDetails",
    component: () => 
        import(/* webpackChunkName: "buyitemdetails" */ '@/views/BuyItemDetails'), 
  },
  {
    path: "/sell",
    component: () => 
        import(/* webpackChunkName: "sell" */ '@/views/Sell'), 
    meta: {
      title: 'Sell - Marketplace'},
    
   
    
  
  },
  {
    path: "/earn",
    component: () => 
        import(/* webpackChunkName: "earn" */ '@/views/Estimate'), 
    meta: {
      title: 'Earn - Marketplace'},
  },
  
  {
    path: "/account",
    component: () => 
        import(/* webpackChunkName: "account" */ '@/views/Account'), 
    meta: {
      title: 'Account - Marketplace'},
  },

  {
    path: "/messages",
    //component: Messages,
    component: () => 
        import(/* webpackChunkName: "messages" */ '@/views/Messages'), 
   
    meta: {
      title: 'Messages - Marketplace'},
  },
  {
    path: "/explore",

    component: () => 
        import(/* webpackChunkName: "explore" */ '@/views/Explore'), 
   
    meta: {
      title: 'Explore - Marketplace'},
  },
  {
    path: "/faq",

    component: () => 
        import(/* webpackChunkName: "faq" */ '@/views/FAQ'), 
   
    meta: {
      title: 'FAQ - Marketplace'},
  },
  {
    path: '*',
    //name: 'catchAll',
    component: () => 
        import(/* webpackChunkName: "buy" */ '@/views/Buy'), 
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
  scrollBehavior (to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    } else {
      return { x: 0, y: 0 }
    }
  }
  
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
