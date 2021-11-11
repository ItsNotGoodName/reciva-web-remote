import { createStore } from "vuex";

import api from "./api";
import p from "./storePreset"

export default createStore({
  state() {
    return {
      page: "player",
      loading: true,
      selectedRadio: null,
      radio: {},
      radioConnected: false,
      radioConnecting: false,
      radioWS: null,
      radios: [],
      message: null
    };
  },
  getters: {
    radioReady(state) {
      return state.selectedRadio && state.radioConnected;
    },
  },
  modules: {
    p
  },
  mutations: {
    SET_PAGE(state, page) {
      state.page = page;
    },
    SET_LOADING(state, loading) {
      state.loading = loading;
    },
    UPDATE_RADIO(state, radio) {
      for (let k in radio) {
        state.radio[k] = radio[k];
      }
    },
    SET_RADIO_POWER(state, power) {
      state.radio.power = power;
    },
    SET_RADIO_VOLUME(state, volume) {
      state.radio.volume = volume;
    },
    SET_SELECTED_RADIO(state, selectedRadio) {
      state.selectedRadio = selectedRadio;
      localStorage.lastRadioUUID = selectedRadio.uuid;
    },
    SET_RADIO_WS(state, radioWS) {
      state.radioWS = radioWS;
    },
    SET_RADIO_CONNECTING(state, radioConnecting) {
      state.radioConnecting = radioConnecting;
      if (radioConnecting) {
        state.radioConnected = false;
      }
    },
    SET_RADIO_CONNECTED(state, radioConnected) {
      state.radioConnected = radioConnected;
      state.radioConnecting = false;
    },
    SET_RADIOS(state, radios) {
      state.radios = radios;
      if (state.selectedRadio) {
        for (let radio of radios) {
          if (radio.uuid == state.selectedRadio.uuid) {
            state.selectedRadio = radio;
            return
          }
        }
        state.selectedRadio = null;
      }
    },
    SET_MESSAGE(state, message) {
      state.message = message;
    }
  },
  actions: {
    init({ dispatch, state }) {
      return dispatch('loadRadios').then(() => {
        if (localStorage.lastRadioUUID) {
          for (let radio of state.radios) {
            if (radio.uuid == localStorage.lastRadioUUID) {
              dispatch('setSelectedRadio', radio);
              return
            }
          }
        }
      })
    },
    showEditPage({ commit }) {
      commit("SET_PAGE", "edit");
    },
    showPlayerPage({ commit }) {
      commit("SET_PAGE", "player");
    },
    loadRadios({ commit }) {
      commit("SET_LOADING", true);
      return api.getRadios().then((radios) => {
        commit("SET_RADIOS", radios);
      }).finally(() => {
        commit("SET_LOADING", false);
      });
    },
    refreshRadio({ dispatch, state }) {
      if (!state.selectedRadio) return Promise.reject(ErrRadioNotSelected);

      return dispatch("refreshRadioWS")
        .then(() => api.renewRadio(state.selectedRadio.uuid))
    },
    playRadioPreset({ state }, num) {
      if (!state.selectedRadio) return Promise.reject(ErrRadioNotSelected);

      return api.updateRadio(state.selectedRadio.uuid, { preset: num });
    },
    toggleRadioPower({ state, commit }) {
      if (!state.selectedRadio) return Promise.reject(ErrRadioNotSelected);

      let newPower = !state.radio.power;
      return api
        .updateRadio(state.selectedRadio.uuid, { power: newPower })
        .then(() => commit("SET_RADIO_POWER", newPower));
    },
    refreshRadioVolume({ state }) {
      if (!state.selectedRadio) return Promise.reject(ErrRadioNotSelected);

      return api.refreshRadioVolume(state.selectedRadio.uuid);
    },
    setRadioVolume({ state }, volume) {
      if (!state.selectedRadio) return Promise.reject(ErrRadioNotSelected);

      return api.updateRadio(state.selectedRadio.uuid, {
        volume: volume,
      });
    },
    setSelectedRadio({ commit, dispatch }, radio) {
      commit("SET_SELECTED_RADIO", radio);
      return dispatch("refreshRadioWS");
    },
    discoverRadios({ dispatch }) {
      return api
        .discoverRadios()
        .then(() => dispatch("loadRadios"))
    },
    refreshRadioWS({ state, commit, dispatch }) {
      // Full state update when radio websocket is connected
      if (state.radioConnected) {
        if (state.selectedRadio) {
          state.radioWS.send(state.selectedRadio.uuid);
        }
        return;
      }

      // Do not create a new websocket if current websocket is connecting or selectedRadio is not set
      if (state.radioConnecting || !state.selectedRadio) {
        return;
      }

      commit("SET_RADIO_CONNECTING", true);
      let ws = api.radioWS(state.selectedRadio.uuid);

      let onMessage = function (event) {
        let radio = JSON.parse(event.data);
        if (radio.uuid != state.selectedRadio.uuid) return;
        commit("UPDATE_RADIO", radio);
      }

      let onFirstMessage = function (event) {
        let radio = JSON.parse(event.data);
        if (radio.uuid != state.selectedRadio.uuid) return;
        commit("UPDATE_RADIO", radio);

        commit("SET_RADIO_CONNECTED", true);
        commit("SET_MESSAGE", null);
        ws.removeEventListener("message", onFirstMessage);
        ws.addEventListener("message", onMessage);
      }

      ws.addEventListener("message", onFirstMessage);

      let onDisconnect = function (event) {
        commit("SET_RADIO_CONNECTED", false);
        commit("SET_MESSAGE", { content: "Disconnected from radio, reconnecting in 3 seconds", severity: "error" });
        setTimeout(() => {
          commit("SET_MESSAGE", { content: "Reconnecting...", severity: "info" });
          dispatch("refreshRadioWS");
        }, 3000);
      };

      // Handle close
      ws.addEventListener("close", onDisconnect);

      commit("SET_RADIO_WS", ws);
    },
  },
});
