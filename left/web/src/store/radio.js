import { useToast } from "vue-toastification";

import api from "../api";

const toast = useToast();

export default {
  state() {
    return {
      radioInit: false,
      radio: {},
      radioRefreshing: false,
      radioVolumeRefreshing: false,
      radioVolumeDiff: 0,
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
      return !!state.radioUUID;
    },
    radioReady(state) {
      return state.radio.uuid == state.radioUUID && state.radioWSConnected;
    },
  },
  mutations: {
    SET_RADIO(state, radio) {
      state.radio = radio;
    },
    SET_RADIO_INIT(state, radioInit) {
      state.radioInit = radioInit;
    },
    MERGE_RADIO(state, radio) {
      for (let k in radio) {
        state.radio[k] = radio[k];
      }
    },
    SET_RADIO_POWER(state, power) {
      state.radio.power = power;
    },
    SET_RADIO_PRESET(state, preset) {
      state.radio.preset = preset;
    },
    SET_RADIO_VOLUME_DIFF(state, radioVolumeDiff) {
      state.radioVolumeDiff = radioVolumeDiff;
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
    _callRadio({ dispatch }, opts) {
      return dispatch("_call", opts).catch((res) => {
        if (res.code == 404) {
          dispatch("listRadios");
        }
      });
    },
    refreshRadioAll({ dispatch, getters }) {
      let proms = [dispatch("listRadios"), dispatch("connectRadioWS")];
      if (getters.radioSelected) {
        proms.push(
          dispatch("refreshRadio").then(() => dispatch("refreshRadioVolume"))
        );
      }
      return Promise.all(proms).then(() => {
        toast.success("refreshed");
      });
    },
    initRadio({ dispatch, state, commit }) {
      if (state.radioInit) {
        return;
      }
      commit("SET_RADIO_INIT", true);

      window.addEventListener("focus", () => {
        dispatch("connectRadioWS");
      });

      let lastRadioUUID = localStorage.lastRadioUUID;
      dispatch("connectRadioWS");
      return dispatch("listRadios").then(() => {
        dispatch("setRadioUUID", lastRadioUUID);
      });
    },
    discoverRadios({ dispatch, state }) {
      return dispatch("_call", {
        promise: api.discoverRadios(),
        loadingMutation: "SET_RADIOS_DISCOVERING",
      })
        .then(() => dispatch("listRadios"))
        .then(() =>
          toast.success("discovered " + state.radios.length + " radios")
        );
    },
    listRadios({ commit, dispatch, state }) {
      return dispatch("_call", {
        promise: api.listRadios(),
        loadingMutation: "SET_RADIOS_LOADING",
      }).then(({ result }) => {
        commit("SET_RADIOS", result);
        return dispatch("setRadioUUID", state.radioUUID);
      });
    },

    refreshRadio({ dispatch, state }) {
      return dispatch("_callRadio", {
        promise: api.refreshRadio(state.radioUUID),
        loadingMutation: "SET_RADIO_REFRESHING",
      });
    },
    refreshRadioVolume({ dispatch, state }) {
      return dispatch("_callRadio", {
        promise: api.refreshRadioVolume(state.radioUUID),
        loadingMutation: "SET_RADIO_VOLUME_REFRESHING",
      });
    },
    toggleRadioPower({ commit, dispatch, state }) {
      let power = !state.radio.power;
      return dispatch("_callRadio", {
        promise: api.patchRadio(state.radioUUID, { power }),
      }).then(() => {
        commit("SET_RADIO_POWER", power);
      });
    },
    changeRadioVolume({ state, commit, dispatch }, change) {
      commit("SET_RADIO_VOLUME_DIFF", state.radioVolumeDiff + change);
      return dispatch("_callRadio", {
        promise: api.patchRadio(state.radioUUID, {
          volume: state.radio.volume + state.radioVolumeDiff,
        }),
      }).finally(() => {
        commit("SET_RADIO_VOLUME_DIFF", state.radioVolumeDiff - change);
      });
    },
    setRadioPreset({ commit, state, dispatch }, preset) {
      return dispatch("_callRadio", {
        promise: api.patchRadio(state.radioUUID, { preset }),
      }).then(() => {
        commit("MERGE_RADIO", { preset });
      });
    },
    setRadioUUID({ commit, state, getters }, newUUID) {
      let uuid = "";
      for (let radio of state.radios) {
        if (radio.uuid == newUUID) {
          uuid = radio.uuid;
        }
      }

      commit("SET_RADIO_UUID", uuid);
      commit("SET_RADIO", {});
      if (!(state.radioWSConnected && getters.radioSelected)) {
        return;
      }

      state.radioWS.send(state.radioUUID);
    },
    connectRadioWS({ commit, dispatch, state }) {
      if (state.radioWSConnecting || state.radioWSConnected) {
        return;
      }

      commit("SET_RADIO_WS_CONNECTING", true);

      let toastID = toast("connecting...", { timeout: false });

      let ws = api.getRadioWS();

      ws.addEventListener("open", () => {
        commit("SET_RADIO_WS_CONNECTED", true);
        commit("SET_RADIO_WS_CONNECTING", false);

        toast.clear();

        if (state.radioUUID) {
          ws.send(state.radioUUID);
        }
      });

      ws.addEventListener("message", (event) => {
        let radio = JSON.parse(event.data);
        if (radio.uuid != state.radioUUID) return;
        commit("MERGE_RADIO", radio);
      });

      ws.addEventListener("close", () => {
        commit("SET_RADIO_WS_CONNECTED", false);
        commit("SET_RADIO_WS_CONNECTING", false);

        dispatch("listRadios");

        toast.dismiss(toastID);
        toast.error("disconnected, reconnecting in 10 seconds", {
          timeout: 9000,
        });

        setTimeout(() => dispatch("connectRadioWS"), 10000);
      });

      commit("SET_RADIO_WS", ws);
    },
  },
};
