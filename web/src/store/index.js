import { createStore } from "vuex";

import { MAX_MESSAGES } from "../constants";

import p from "./preset";
import r from "./radio";

const generateID = (function () {
  let id = 0;
  return () => {
    id += 1;
    return id;
  };
})();

export default createStore({
  state() {
    return {
      page: "play",
      messages: [],
      popups: {},
    };
  },
  mutations: {
    SET_PAGE(state, page) {
      state.page = page;
    },
    ADD_MESSAGE(state, message) {
      state.messages.push(message);
      if (state.messages.length > MAX_MESSAGES) {
        state.messages.shift();
      }
    },
    ADD_POPUP(state, popup) {
      state.popups[popup.id] = popup;
    },
    DELETE_POPUP(state, id) {
      delete state.popups[id];
    },
  },
  actions: {
    togglePage({ commit, state }) {
      commit("SET_PAGE", state.page == "play" ? "preset" : "play");
    },
    addMessage({ commit }, { type, text }) {
      let msg = { id: generateID(), type, text };

      commit("ADD_MESSAGE", msg);
      commit("ADD_POPUP", msg);

      setTimeout(() => {
        commit("DELETE_POPUP", msg.id);
      }, 5000);
    },
    deletePopup({ commit }, id) {
      commit("DELETE_POPUP", id);
    },
  },
  modules: {
    p,
    r,
  },
});
