import { createApp } from "vue";

import Toast from "vue-toastification";

import App from "./App.vue";
import store from "./store";

import "./scss/index.scss";

createApp(App).use(Toast).use(store).mount("#app");
