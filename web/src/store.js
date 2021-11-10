import { createStore } from "vuex";

import api from "./api";
import {
  ErrRadioNotSelected,
} from "./constants";
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
    },
    ADD_MESSAGE(state, params) {
      params.id = state.messageID;

      state.messages[state.messageID] = params;
      state.messageID += 1;

      let keys = Object.keys(state.messages);
      if (keys.length > 3) {
        delete state.messages[keys[0]];
      }
    },
  },
  actions: {
    init({ dispatch, commit }) {
      return Promise.all([
        dispatch("loadRadios"),
      ]).finally(() => {
        commit("SET_LOADING", false);
      });
    },
    showEditPage({ commit }) {
      commit("SET_PAGE", "edit");
    },
    showPlayerPage({ commit }) {
      commit("SET_PAGE", "player");
    },
    loadRadios({ commit }) {
      return api.getRadios().then((radios) => {
        commit("SET_RADIOS", radios);
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
    setSelectedRadio({ commit, state, dispatch }, radio) {
      if (radio.uuid == state.radio.uuid) {
        return Promise.resolve();
      }
      commit("SET_SELECTED_RADIO", radio);
      return dispatch("refreshRadioWS");
    },
    discoverRadios({ dispatch, state }) {
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

      // Handle messsage
      ws.addEventListener("message", function (event) {
        let radio = JSON.parse(event.data);
        if (radio.uuid != state.selectedRadio.uuid) return;
        commit("UPDATE_RADIO", radio);
        commit("SET_RADIO_CONNECTED", true);
      });

      ws.addEventListener("open", () => {
      });

      let onDisconnect = function (event) {
        console.log(event);
        commit("SET_RADIO_CONNECTED", false);
        setTimeout(() => {
          dispatch("refreshRadioWS");
        }, 3000);
      };

      // Handle close
      ws.addEventListener("close", onDisconnect);

      commit("SET_RADIO_WS", ws);
    },
  },
});
