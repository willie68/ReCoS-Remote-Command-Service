import { createApp } from 'vue'
import { createStore } from 'vuex'
import App from './App.vue'

import PrimeVue from 'primevue/config'

import "primevue/resources/themes/saga-blue/theme.css"
import "primevue/resources/primevue.min.css"
import "primeicons/primeicons.css"
import 'primeflex/primeflex.css'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import Toolbar from "primevue/toolbar"
import Button from "primevue/button"
import Password from 'primevue/password'
import SplitButton from 'primevue/splitbutton'
import Panel from 'primevue/panel'
import Menu from 'primevue/menu'
import Splitter from 'primevue/splitter'
import SplitterPanel from 'primevue/splitterpanel'
import Accordion from 'primevue/accordion'
import AccordionTab from 'primevue/accordiontab'
import Listbox from 'primevue/listbox'
import InputNumber from 'primevue/inputnumber'
import Badge from 'primevue/badge'
import BadgeDirective from 'primevue/badgedirective'
import ScrollPanel from 'primevue/scrollpanel'
import Dropdown from 'primevue/dropdown'
import Tooltip from 'primevue/tooltip'
import Fieldset from 'primevue/fieldset'
import Textarea from 'primevue/textarea'
import TabView from 'primevue/tabview'
import TabPanel from 'primevue/tabpanel'
import ColorPicker from 'primevue/colorpicker'
import Checkbox from 'primevue/checkbox'
import ToastService from 'primevue/toastservice'
import Toast from 'primevue/toast'

const store = createStore({
    state () {
      return {
        count: 0,      
        servicePort: 9280,
        baseURL: window.location.protocol + "//localhost:9280/api/v1/",
        password: "",
        authheader: {}
      }
    },
    mutations: {
      increment (state) {
        state.count++
      },
      baseURL (state, baseurl) {
          state.baseURL = baseurl
      },
      password (state, password) {
        state.password = password
        state.authheader = {Authorization: `Basic ${btoa(`admin:${password}`)}`}
      }
    }
  })

const app = createApp(App)

app.directive('badge', BadgeDirective)
app.directive('tooltip', Tooltip)
app.use(store)
app.use(PrimeVue, {ripple: true})
app.use(ToastService);
app.component('Toolbar', Toolbar)
app.component('Button', Button)
app.component('SplitButton', SplitButton)
app.component('Dialog', Dialog)
app.component('InputText', InputText)
app.component('Password', Password)
app.component('Panel', Panel)
app.component('Menu', Menu)
app.component('Splitter', Splitter)
app.component('SplitterPanel', SplitterPanel)
app.component('Accordion', Accordion)
app.component('AccordionTab', AccordionTab)
app.component('Listbox', Listbox)
app.component('InputNumber', InputNumber)
app.component('Badge', Badge)
app.component('BadgeDirective', BadgeDirective)
app.component('ScrollPanel', ScrollPanel)
app.component('Dropdown', Dropdown)
app.component('Fieldset', Fieldset)
app.component('Textarea', Textarea)
app.component('TabView', TabView)
app.component('TabPanel', TabPanel)
app.component('ColorPicker', ColorPicker)
app.component('Checkbox', Checkbox)
app.component('Toast', Toast)

app.mount('#app')