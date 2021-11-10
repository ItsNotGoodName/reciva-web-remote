import { createApp } from 'vue'
import PrimeVue from 'primevue/config';
import ToastService from 'primevue/toastservice';
import BadgeDirective from 'primevue/badgedirective';

import App from './App.vue'
import store from "./store"

import "primevue/resources/themes/saga-blue/theme.css"
import "primevue/resources/primevue.min.css"
import "primeicons/primeicons.css"
import "/node_modules/primeflex/primeflex.css"

createApp(App)
    .use(PrimeVue)
    .use(ToastService)
    .directive('badge', BadgeDirective)
    .use(store)
    .mount('#app');