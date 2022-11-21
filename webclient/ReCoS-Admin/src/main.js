import { createApp } from "vue";
import { createPinia } from "pinia";
import App from "./App.vue";
import "./assets/main.css";
import mitt from 'mitt';

import PrimeVue from "primevue/config";
import "primevue/resources/themes/vela-blue/theme.css";
import "primevue/resources/primevue.min.css";
import "primeicons/primeicons.css";
import "primeflex/primeflex.css";

import Dialog from "primevue/dialog";

const emitter = mitt();
const pinia = createPinia();
const app = createApp(App);

app.config.globalProperties.emitter = emitter;

app.use(PrimeVue);
app.use(pinia);

app.component("Dialog", Dialog);

app.mount("#app");
