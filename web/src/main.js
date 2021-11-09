import { createApp } from 'vue'
import PrimeVue from 'primevue/config';
import App from './App.vue'
import store from "./store"

import "primevue/resources/themes/saga-blue/theme.css"
import "primevue/resources/primevue.min.css"
import "primeicons/primeicons.css"
import "/node_modules/primeflex/primeflex.css"

createApp(App)
    .use(PrimeVue)
    .use(store)
    .mount('#app');