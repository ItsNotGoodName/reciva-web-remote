import api from "../api";
import { MESSAGE_SUCCESS } from "../constants";
import { call, radioCall } from "./util";

export default {
  state() {
    return {
      radio: {},
      radioRefreshing: false,
      radioVolumeRefreshing: false,
      radioVolumeChanging: 0,
      radioUUID: "",
      radios: [],
      radiosDiscovering: false,
      radiosLoading: false,
      radioWS: null,
      radioWSConnecting: false,
      radioWSConnected: false,
    };
  },
  getters: {
    radioSelected(state) {
      return state.radioUUID != "";
    },
    radioLoaded(state) {
      return state.radio.uuid == state.radioUUID && state.radioWSConnected;
    },
  },
  mutations: {
    SET_RADIO(state, radio) {
      state.radio = radio;
    },
    MERGE_RADIO(state, radio) {
      for (let k in radio) {
        state.radio[k] = radio[k];
      }
    },
    SET_RADIO_POWER(state, power) {
      state.radio.power = power;
    },
    SET_RADIO_VOLUME(state, volume) {
      if (volume < 0) {
        state.radio.volume = 0;
      } else if (volume > 100) {
        state.radio.volume = 100;
      } else {
        state.radio.volume = volume;
      }
    },
    CHANGE_RADIO_VOLUME_CHANGING(state, radioVolumeChanging) {
      state.radioVolumeChanging += radioVolumeChanging;
    },
    SET_RADIO_VOLUME_REFRESHING(state, radioVolumeRefreshing) {
      state.radioVolumeRefreshing = radioVolumeRefreshing;
    },
    SET_RADIO_REFRESHING(state, radioRefreshing) {
      state.radioRefreshing = radioRefreshing;
    },
    SET_RADIO_UUID(state, radioUUID) {
      state.radioUUID = radioUUID;
      localStorage.lastRadioUUID = radioUUID;
    },
    SET_RADIOS(state, radios) {
      state.radios = radios;
    },
    SET_RADIOS_DISCOVERING(state, radiosDiscovering) {
      state.radiosDiscovering = radiosDiscovering;
    },
    SET_RADIOS_LOADING(state, radiosLoading) {
      state.radiosLoading = radiosLoading;
    },
    SET_RADIO_WS_CONNECTING(state, radioWSConnecting) {
      state.radioWSConnecting = radioWSConnecting;
    },
    SET_RADIO_WS_CONNECTED(state, radioWSConnected) {
      state.radioWSConnected = radioWSConnected;
    },
    SET_RADIO_WS(state, radioWS) {
      state.radioWS = radioWS;
    },
  },
  actions: {
    initRadio({ dispatch }) {
      let lastRadioUUID = localStorage.lastRadioUUID;
      return dispatch("listRadios").then(() => {
        dispatch("setRadioUUID", lastRadioUUID);
      });
    },
    listRadios({ commit, dispatch, state }) {
      return call({
        commit,
        dispatch,
        promise: api.listRadios(),
        loadingMutation: "SET_RADIOS_LOADING",
      }).then(({ result }) => {
        commit("SET_RADIOS", result);
        dispatch("setRadioUUID", state.radioUUID);
      });
    },
    discoverRadios({ commit, dispatch }) {
      return call({
        commit,
        dispatch,
        promise: api.discoverRadios(),
        loadingMutation: "SET_RADIOS_DISCOVERING",
      }).then(() => dispatch("listRadios"));
    },
    refreshRadio({ commit, dispatch, state }) {
      return radioCall({
        commit,
        dispatch,
        promise: api.refreshRadio(state.radioUUID),
        loadingMutation: "SET_RADIO_REFRESHING",
      }).then(() =>
        dispatch("addMessage", {
          type: MESSAGE_SUCCESS,
          text: "refreshed radio",
        })
      );
    },
    refreshRadioVolume({ commit, dispatch, state }) {
      return radioCall({
        commit,
        dispatch,
        promise: api.refreshRadioVolume(state.radioUUID),
        loadingMutation: "SET_RADIO_VOLUME_REFRESHING",
      });
    },
    toggleRadioPower({ commit, dispatch, state }) {
      let power = !state.radio.power;
      return radioCall({
        commit,
        dispatch,
        promise: api.patchRadio(state.radioUUID, { power }),
      }).then(() => {
        commit("SET_RADIO_POWER", power);
      });
    },
    setRadioVolume({ state, commit, dispatch }, volume) {
      commit("CHANGE_RADIO_VOLUME_CHANGING", 1);
      commit("SET_RADIO_VOLUME", volume);
      return radioCall({
        commit,
        dispatch,
        promise: api.patchRadio(state.radioUUID, { volume }),
      }).finally(() => {
        commit("CHANGE_RADIO_VOLUME_CHANGING", -1);
      });
    },
    setRadioPreset({ commit, state, dispatch }, preset) {
      return radioCall({
        commit,
        dispatch,
        promise: api.patchRadio(state.radioUUID, { preset }),
      }).then(() => {
        commit("MERGE_RADIO", { preset });
      });
    },
    setRadioUUID({ commit, dispatch, state }, uuid) {
      for (let radio of state.radios) {
        if (radio.uuid == uuid) {
          commit("SET_RADIO_UUID", uuid);
          commit("SET_RADIO", {});
          dispatch("refreshRadioWS");
          return;
        }
      }

      commit("SET_RADIO_UUID", "");
      commit("SET_RADIO", {});
      dispatch("closeRadioWS");
    },
    closeRadioWS({ state }) {
      if (state.radioWSConnected) {
        state.radioWS.close();
      }
    },
    refreshRadioWS({ commit, state, getters }) {
      if (state.radioWSConnecting || !getters.radioSelected) {
        return;
      }

      if (state.radioWSConnected) {
        state.radioWS.send(state.radioUUID);
        return;
      }

      commit("SET_RADIO_WS_CONNECTING", true);
      let ws = api.getRadioWS(state.radioUUID);

      let onMessage = function (event) {
        let radio = JSON.parse(event.data);
        if (radio.uuid != state.radioUUID) return;
        commit("MERGE_RADIO", radio);
      };

      let onFirstMessage = function (event) {
        onMessage(event);

        commit("SET_RADIO_WS_CONNECTED", true);
        commit("SET_RADIO_WS_CONNECTING", false);

        ws.removeEventListener("message", onFirstMessage);
        ws.addEventListener("message", onMessage);
      };

      ws.addEventListener("message", onFirstMessage);

      let onDisconnect = function () {
        commit("SET_RADIO_WS_CONNECTED", false);
        commit("SET_RADIO_WS_CONNECTING", false);
      };

      ws.addEventListener("close", onDisconnect);

      commit("SET_RADIO_WS", ws);
    },
  },
};
