import { createApp } from 'vue'
import { createStore } from 'vuex'
import App from './App.vue'
import './assets/global.css'

import "primevue/resources/themes/vela-blue/theme.css"
//import 'primevue/resources/themes/saga-blue/theme.css';
import "primevue/resources/primevue.min.css"
import "primeicons/primeicons.css"
import 'primeflex/primeflex.css'

import PrimeVue from 'primevue/config';

import Button from 'primevue/button';
import Dialog from 'primevue/dialog'
import Dropdown from 'primevue/dropdown'
import Menu from 'primevue/menu'
import Sidebar from 'primevue/sidebar'; 
import Toast from 'primevue/toast';
import ToastService from 'primevue/toastservice';
import Tooltip from 'primevue/tooltip'

const store = createStore({
    state () {
      return {
        count: 0,      
        servicePort: 9280,
        baseURL: window.location.protocol + "//localhost:9280/api/v1/",
        packageVersion: process.env.VUE_APP_VERSION || '0',
      }
    },
    mutations: {
      increment (state) {
        state.count++
      },
      baseURL (state, baseurl) {
        state.baseURL = baseurl
        if (!baseurl.endsWith("/")) {
          state.baseURL = state.baseURL + "/"
        }
      },
    }
  })

const app = createApp(App)

app.directive('tooltip', Tooltip)
app.use(store)
app.use(PrimeVue);
app.use(ToastService);


app.component('Button', Button);
app.component('Dialog', Dialog);
app.component('Dropdown', Dropdown);
app.component('Menu', Menu);
app.component('Sidebar', Sidebar);
app.component('Toast', Toast);

app.mount('#app')
