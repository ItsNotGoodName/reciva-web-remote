import { createApp } from 'vue'
import App from './App.vue'
import './index.css'

import { OhVueIcon, addIcons } from "oh-vue-icons";
import { FaPlay, FaStop, FaSync, FaQuestion, FaSearch, FaRedo, FaPowerOff, FaVolumeDown, FaVolumeUp, MdRadio, FaVolumeMute, FaEdit } from "oh-vue-icons/icons";


addIcons(FaPlay, FaStop, FaSync, FaQuestion, FaSearch, FaRedo, FaPowerOff, FaVolumeDown, FaVolumeUp, MdRadio, FaVolumeMute, FaEdit);

createApp(App)
  .component("v-icon", OhVueIcon)
  .mount('#app')

