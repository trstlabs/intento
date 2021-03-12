import Vue from 'vue'
import Vuetify from 'vuetify/lib'
import 'vuetify/dist/vuetify.min.css'

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
    
   }, };

export default new Vuetify(opts)