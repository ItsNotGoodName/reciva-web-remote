import { createStore } from "vuex";

import api from "./api";
import {
  MsgRadioRefreshed,
  ErrRadioNotSelected,
  MsgDiscoveredRadiosFn,
  MsgConnected,
  MsgDisconnected,
  MsgStreamAdded,
  MsgStreamDeleted,
  MsgStreamUpdated,
  MsgConnecting,
  DefNotificationTimeout,
} from "./constants";

export default createStore({
  state() {
    return {
      config: {
        presetsEnabled: false,
      },
      edit: false,
      isStreamEdit: false,
      notificationID: 0,
      notifications: {},
      presets: [],
      radio: {},
      radioConnected: false,
      radioConnecting: false,
      radioUUID: null,
      radioWS: null,
      radios: [],
      showStream: false,
      streamLoading: false,
      stream: {
        name: "",
        content: "",
        uri: "",
      },
      streams: [],
    };
  },
  mutations: {
    SET_EDIT(state, edit) {
      state.edit = edit;
    },
    SET_CONFIG(state, config) {
      state.config = config;
    },
    UPDATE_RADIO(state, radio) {
      for (let k in radio) {
        state.radio[k] = radio[k];
      }
    },
    SET_RADIO_POWER(state, power) {
      state.radio.power = power;
    },
    SET_RADIO_UUID(state, uuid) {
      state.radioUUID = uuid;
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
      let rds = {};
      for (let r in radios) {
        rds[radios[r].uuid] = radios[r].name;
      }
      state.radios = rds;
    },
    SET_PRESETS(state, presets) {
      state.presets = presets;
    },
    SET_STREAMS(state, streams) {
      state.streams = streams;
    },
    ADD_NOTIFICATION(state, params) {
      params.id = state.notificationID;

      state.notifications[state.notificationID] = params;
      state.notificationID += 1;

      let keys = Object.keys(state.notifications);
      if (keys.length > 3) {
        delete state.notifications[keys[0]];
      }
    },
    CLEAR_NOTIFICATIONS(state) {
      for (let k in state.notifications) {
        delete state.notifications[k];
      }
    },
    CLEAR_NOTIFICATION(state, id) {
      delete state.notifications[id];
    },
    SET_SHOW_STREAM(state, showStream) {
      state.showStream = showStream;
    },
    SET_STREAM(state, stream) {
      state.stream = stream;
    },
    SET_STREAM_NAME(state, name) {
      state.stream.name = name;
    },
    SET_STREAM_CONTENT(state, content) {
      state.stream.content = content;
    },
    SET_STREAM_URI(state, uri) {
      state.stream.uri = uri;
    },
    SET_IS_STREAM_EDIT(state, isStreamEdit) {
      state.isStreamEdit = isStreamEdit;
    },
    SET_STREAM_LOADING(state, streamLoading) {
      state.streamLoading = streamLoading;
    }
  },
  actions: {
    init({ dispatch }) {
      return dispatch("loadConfig").then(() => {
        if (this.state.config.presetsEnabled) {
          return Promise.all([
            dispatch("loadPresets"),
            dispatch("loadStreams"),
            dispatch("loadRadios"),
          ]);
        }
        return dispatch("loadRadios");
      });
    },
    loadConfig({ commit }) {
      return api.getConfig().then((config) => commit("SET_CONFIG", config));
    },
    loadRadios({ commit }) {
      return api.getRadios().then((radios) => {
        commit("SET_RADIOS", radios);
      });
    },
    loadPresets({ commit }) {
      return api.getPresets().then((presets) => {
        commit("SET_PRESETS", presets);
      });
    },
    loadStreams({ dispatch, commit }) {
      return api.getStreams().then((streams) => {
        commit("SET_STREAMS", streams);
      });
    },
    refreshRadio({ dispatch, state }) {
      if (!state.radioUUID) return Promise.reject(ErrRadioNotSelected);

      return dispatch("refreshRadioWS")
        .then(() => api.renewRadio(state.radioUUID))
        .then(() =>
          dispatch("addNotification", {
            type: "success",
            message: MsgRadioRefreshed,
          })
        );
    },
    playRadioPreset({ state }, num) {
      if (!state.radioUUID) return Promise.reject(ErrRadioNotSelected);

      return api.updateRadio(state.radioUUID, { preset: num });
    },
    toggleRadioPower({ state, commit }) {
      if (!state.radioUUID) return Promise.reject(ErrRadioNotSelected);

      let newPower = !state.radio.power;
      return api
        .updateRadio(state.radioUUID, { power: newPower })
        .then(() => commit("SET_RADIO_POWER", newPower));
    },
    refreshRadioVolume({ state }) {
      if (!state.radioUUID) return Promise.reject(ErrRadioNotSelected);

      return api.refreshRadioVolume(state.radioUUID);
    },
    increaseRadioVolume({ state }) {
      if (!state.radioUUID) return Promise.reject(ErrRadioNotSelected);

      return api.updateRadio(state.radioUUID, {
        volume: state.radio.volume + 5,
      });
    },
    decreaseRadioVolume({ state }) {
      if (!state.radioUUID) return Promise.reject(ErrRadioNotSelected);

      return api.updateRadio(state.radioUUID, {
        volume: state.radio.volume - 5,
      });
    },
    setRadioUUID({ commit, state, dispatch }, uuid) {
      if (uuid == state.radioUUID) {
        return Promise.resolve();
      }
      commit("SET_RADIO_UUID", uuid);
      return dispatch("refreshRadioWS");
    },
    discoverRadios({ dispatch, commit, state }) {
      commit("CLEAR_NOTIFICATIONS");
      return api
        .discoverRadios()
        .then(() => dispatch("loadRadios"))
        .then(() =>
          dispatch("addNotification", {
            type: "success",
            message: MsgDiscoveredRadiosFn(Object.keys(state.radios).length),
          })
        );
    },
    refreshRadioWS({ state, commit, dispatch }) {
      // Full state update when radio websocket is connected
      if (state.radioConnected) {
        if (state.radioUUID) {
          state.radioWS.send(state.radioUUID);
        }
        return;
      }

      // Do not create a new websocket if current websocket is connecting or radioUUID is not set
      if (state.radioConnecting || !state.radioUUID) {
        return;
      }

      commit("SET_RADIO_CONNECTING", true);
      let ws = api.radioWS(state.radioUUID);

      // Handle messsage
      ws.addEventListener("message", function (event) {
        let radio = JSON.parse(event.data);
        if (radio.uuid != state.radioUUID) return;
        commit("UPDATE_RADIO", radio);
        commit("SET_RADIO_CONNECTED", true);
      });

      ws.addEventListener("open", () => {
        commit("CLEAR_NOTIFICATIONS");
        dispatch("addNotification", { type: "success", message: MsgConnected });
      });

      let onEnd = function (event) {
        console.log(event);
        commit("SET_RADIO_CONNECTED", false);
        commit("CLEAR_NOTIFICATIONS");
        dispatch("addNotification", {
          type: "error",
          message: MsgDisconnected,
        });
        setTimeout(() => {
          dispatch("addNotification", {
            type: "warning",
            message: MsgConnecting,
          });
          dispatch("refreshRadioWS");
        }, 3000);
      };

      // Handle close
      ws.addEventListener("close", onEnd);

      commit("SET_RADIO_WS", ws);
    },
    addNotification({ commit, state }, params) {
      let id = state.notificationID;
      commit("ADD_NOTIFICATION", params);
      params.timeout != 0 &&
        setTimeout(
          () => {
            commit("CLEAR_NOTIFICATION", id);
          },
          params.timeout ? params.timeout : DefNotificationTimeout
        );
    },
    clearNotification({ commit }, id) {
      commit("CLEAR_NOTIFICATION", id);
    },
    toggleEdit({ state, commit }) {
      commit("SET_EDIT", !state.edit);
    },
    clearStream({ commit }) {
      commit("SET_STREAM", { name: "", content: "", uri: "" });
    },
    hideStream({ dispatch, commit }) {
      commit("SET_SHOW_STREAM", false);
      return dispatch("clearStream");
    },
    addStream({ dispatch, commit }) {
      commit("SET_IS_STREAM_EDIT", false)
      commit("SET_SHOW_STREAM", true);
      return dispatch("clearStream");
    },
    newStream({ state, commit, dispatch }) {
      return api.newStream(state.stream).then((res) => {
        commit("SET_STREAM", res);
        commit("SET_IS_STREAM_EDIT", true)
        commit("SET_SHOW_STREAM", true);
        return dispatch("loadStreams")
      }).then(() => {
        return dispatch("addNotification", {
          type: "success",
          message: MsgStreamAdded,
        })
      });
    },
    updateStream({ dispatch, state, commit }) {
      commit("SET_STREAM_LOADING", true);
      api.updateStream(state.stream)
        .then(() => dispatch("loadStreams"))
        .then(() => {
          dispatch("addNotification", {
            type: "success",
            message: MsgStreamUpdated,
          })
        })
        .finally(() => commit("SET_STREAM_LOADING", false));
    },
    loadStream({ commit }, id) {
      commit("SET_STREAM_LOADING", true);
      return api.getStream(id).then((res) => {
        commit("SET_STREAM", res);
        commit("SET_IS_STREAM_EDIT", true)
        commit("SET_SHOW_STREAM", true);
      }).finally(() => commit("SET_STREAM_LOADING", false));
    },
    deleteStream({ dispatch, state, commit }) {
      commit("SET_STREAM_LOADING", true);
      api.deleteStream(state.stream.id)
        .then(() => dispatch("loadStreams"))
        .then(() => dispatch("hideStream"))
        .then(() => dispatch("addNotification", { type: "success", message: MsgStreamDeleted }))
        .finally(() => {
          commit("SET_STREAM_LOADING", false)
        });
    }
  },
});
