import { createApp } from 'vue'
import App from './App.vue'
import 'primeflex/primeflex.css';
import "primevue/resources/themes/saga-blue/theme.css";
import "primevue/resources/primevue.min.css";
import "primeicons/primeicons.css";
import PrimeVue from 'primevue/config'
import Dialog from 'primevue/dialog';
import InputText from 'primevue/inputtext';

const app = createApp(App)

app.use(PrimeVue, {ripple: true})
app.component('Dialog', Dialog);
app.component('InputText', InputText);

app.mount('#app')