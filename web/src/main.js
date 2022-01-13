import { createApp } from 'vue'

import App from './App.vue'
import store from "./store"

import "@fortawesome/fontawesome-free/css/all.css";
import "bulma/css/bulma.css";

createApp(App)
  .use(store)
  .mount('#app')