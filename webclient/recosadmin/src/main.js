import { createApp } from 'vue'
import App from './App.vue'
import 'primeflex/primeflex.css';
import "primevue/resources/themes/saga-blue/theme.css";
import "primevue/resources/primevue.min.css";
import "primeicons/primeicons.css";
import PrimeVue from 'primevue/config'
import Dialog from 'primevue/dialog';
import InputText from 'primevue/inputtext';
import Toolbar from "primevue/toolbar";
import Button from "primevue/button";
import Password from 'primevue/password';


const app = createApp(App)

app.use(PrimeVue, {ripple: true})
app.component('Toolbar', Toolbar);
app.component('Button', Button);
app.component('Dialog', Dialog);
app.component('InputText', InputText);
app.component('Password', Password);

app.mount('#app')