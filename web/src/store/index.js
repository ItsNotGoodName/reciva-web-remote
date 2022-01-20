import { createStore } from "vuex";
import { useToast } from "vue-toastification";

import p from "./preset";
import r from "./radio";

const toast = useToast();

export default createStore({
  data() {
    return {
      page: "",
    };
  },
  mutations: {
    SET_PAGE(state, page) {
      state.page = page;
    },
  },
  actions: {
    togglePage({ state, commit }) {
      commit("SET_PAGE", state.page ? "" : "edit");
    },
    _call({ commit }, { promise, loadingMutation }) {
      return new Promise((resolve, reject) => {
        promise
          .then((res) => {
            if (!res.ok) {
              toast.error(res.error);
              reject(res);
            } else {
              resolve(res);
            }
          })
          .catch((error) => {
            toast.error(error.message);
            reject({ error: error.message });
          });

        if (loadingMutation) {
          commit(loadingMutation, true);
          promise.finally(() => {
            commit(loadingMutation, false);
          });
        }
      });
    },
  },
  modules: {
    p,
    r,
  },
});
