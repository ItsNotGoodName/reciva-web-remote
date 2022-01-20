import { createApp } from "vue";

import "./scss/index.scss";

import Toast from "vue-toastification";

import App from "./App.vue";
import store from "./store";

createApp(App).use(Toast).use(store).mount("#app");
