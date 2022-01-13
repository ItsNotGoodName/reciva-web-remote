import { createStore } from "vuex";

import p from "./preset"
import r from "./radio"

export default createStore({
  modules: {
    p,
    r
  },
});
