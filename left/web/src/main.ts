import { createApp } from 'vue'
import App from './App.vue'
import './index.css'

import { VueQueryPlugin } from "vue-query";

import { OhVueIcon, addIcons } from "oh-vue-icons";
import { FaPlay, FaStop, FaSync, FaQuestion, FaSearch, FaRedo, FaPowerOff, FaVolumeDown, FaVolumeUp, MdRadio, FaVolumeMute, FaEdit, FaItunesNote, FaBars, FaGithub, FaHome, FaTag, FaSpinner, FaInfoCircle, FaTimesCircle } from "oh-vue-icons/icons";

addIcons(FaPlay, FaStop, FaSync, FaQuestion, FaSearch, FaRedo, FaPowerOff, FaVolumeDown, FaVolumeUp, MdRadio, FaVolumeMute, FaEdit, FaItunesNote, FaBars, FaGithub, FaHome, FaTag, FaSpinner, FaInfoCircle, FaTimesCircle);

createApp(App)
  .use(VueQueryPlugin)
  .component("v-icon", OhVueIcon)
  .mount('#app')

