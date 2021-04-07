import Vue from 'vue'
import Vuetify from 'vuetify/lib'
import 'vuetify/dist/vuetify.min.css'
import IconCoin from '../components/IconCoin.vue' 

Vue.use(Vuetify)

const opts = { iconfont: 'md', 
    theme: { 
        themes: {
            light: {
              primary: '#3873F9',
              secondary: '#BDE3F4',
              accent: '#08E4F4',
              error: '#b71c1c',
            },
            dark: {
                primary: '#3873F9',
                secondary: '#3062C6',
                accent: '#08E4F4',
                error: '#b71c1c',
              },
            
     },
    
   },  icons: {
    values: {
      custom: { // name of our custom icon
        component: IconCoin, // our custom component
      },
    },
   },
  };

export default new Vuetify(opts)