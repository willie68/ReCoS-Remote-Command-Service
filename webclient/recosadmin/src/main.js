import { createApp } from 'vue'
import { createStore } from 'vuex'
import App from './App.vue'
import mitt from 'mitt';

import PrimeVue from 'primevue/config'

import "primevue/resources/themes/vela-blue/theme.css"
import "primevue/resources/primevue.min.css"
import "primeicons/primeicons.css"
import 'primeflex/primeflex.css'
import './assets/style.css'

import Accordion from 'primevue/accordion'
import AccordionTab from 'primevue/accordiontab'
import Button from "primevue/button"
import Badge from 'primevue/badge'
import BadgeDirective from 'primevue/badgedirective'
import Calendar from 'primevue/calendar';
import Card from 'primevue/card';
import Checkbox from 'primevue/checkbox'
import ColorPicker from 'primevue/colorpicker'
import ConfirmationService from 'primevue/confirmationservice';
import ConfirmDialog from 'primevue/confirmdialog';
import Dialog from 'primevue/dialog'
import Dropdown from 'primevue/dropdown'
import Fieldset from 'primevue/fieldset'
import FileUpload from 'primevue/fileupload';
import Listbox from 'primevue/listbox'
import Menu from 'primevue/menu'
import InputNumber from 'primevue/inputnumber'
import InputText from 'primevue/inputtext'
import OrderList from 'primevue/orderlist';
import Panel from 'primevue/panel'
import Password from 'primevue/password'
import PickList from 'primevue/picklist';
import ScrollPanel from 'primevue/scrollpanel'
import SplitButton from 'primevue/splitbutton'
import Splitter from 'primevue/splitter'
import SplitterPanel from 'primevue/splitterpanel'
import TabPanel from 'primevue/tabpanel'
import TabView from 'primevue/tabview'
import Textarea from 'primevue/textarea'
import Toast from 'primevue/toast'
import ToastService from 'primevue/toastservice'
import Toolbar from "primevue/toolbar"
import Tooltip from 'primevue/tooltip'


const store = createStore({
    state () {
      return {
        count: 0,      
        servicePort: 9280,
        baseURL: window.location.protocol + "//localhost:9280/api/v1/",
        password: "",
        authheader: {},
        inconlist: [],
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
      password (state, password) {
        state.password = password
        state.authheader = {Authorization: `Basic ${btoa(`admin:${password}`)}`}
      },
      iconlist (state, iconlist) {
        state.iconlist = iconlist
      },
    }
  })

const emitter = mitt()

const app = createApp(App)

app.config.globalProperties.emitter = emitter;

app.directive('badge', BadgeDirective)
app.directive('tooltip', Tooltip)
app.use(store)
app.use(PrimeVue, {ripple: true})
app.use(ToastService);
app.use(ConfirmationService);

app.component('Accordion', Accordion)
app.component('AccordionTab', AccordionTab)
app.component('Badge', Badge)
app.component('BadgeDirective', BadgeDirective)
app.component('Button', Button)
app.component('Calendar', Calendar)
app.component('Card', Card)
app.component('Checkbox', Checkbox)
app.component('ColorPicker', ColorPicker)
app.component('ConfirmDialog', ConfirmDialog)
app.component('Dialog', Dialog)
app.component('Dropdown', Dropdown)
app.component('Fieldset', Fieldset)
app.component('FileUpload', FileUpload)
app.component('InputNumber', InputNumber)
app.component('InputText', InputText)
app.component('Listbox', Listbox)
app.component('Menu', Menu)
app.component('OrderList', OrderList)
app.component('Panel', Panel)
app.component('Password', Password)
app.component('PickList', PickList)
app.component('ScrollPanel', ScrollPanel)
app.component('SplitButton', SplitButton)
app.component('Splitter', Splitter)
app.component('SplitterPanel', SplitterPanel)
app.component('TabView', TabView)
app.component('TabPanel', TabPanel)
app.component('Textarea', Textarea)
app.component('Toast', Toast)
app.component('Toolbar', Toolbar)

app.mount('#app')