import { createStore } from "vuex";

import p from "./preset"
import r from "./radio"

export default createStore({
  state() {
    return {
      page: "play"
    }
  },
  mutations: {
    SET_PAGE(state, page) {
      state.page = page
    }
  },
  actions: {
    togglePage({ commit, state }) {
      commit("SET_PAGE", state.page == "play" ? "preset" : "play")
    }
  },
  modules: {
    p,
    r
  },
});
