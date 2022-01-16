import api from '../api'
import { MESSAGE_SUCCESS, MESSAGE_ERROR } from '../constants'

export default {
  state() {
    return {
      radio: {},
      radioRefreshing: false,
      radioVolumeRefreshing: false,
      radioVolumeChanging: 0,
      radioUUID: '',
      radios: [],
      radiosDiscovering: false,
      radiosLoading: false,
      radioWS: null,
      radioWSConnecting: false,
      radioWSConnected: false,
    }
  },
  getters: {
    radioSelected(state) {
      return state.radioUUID != ""
    },
    radioLoaded(state) {
      return state.radio.uuid == state.radioUUID && state.radioWSConnected
    }
  },
  mutations: {
    SET_RADIO(state, radio) {
      state.radio = radio
    },
    MERGE_RADIO(state, radio) {
      for (let k in radio) {
        state.radio[k] = radio[k]
      }
    },
    SET_RADIO_POWER(state, power) {
      state.radio.power = power
    },
    SET_RADIO_VOLUME(state, volume) {
      if (volume < 0) {
        state.radio.volume = 0
      } else if (volume > 100) {
        state.radio.volume = 100
      } else {

        state.radio.volume = volume
      }
    },
    CHANGE_RADIO_VOLUME_CHANGING(state, radioVolumeChanging) {
      state.radioVolumeChanging += radioVolumeChanging
    },
    SET_RADIO_VOLUME_REFRESHING(state, radioVolumeRefreshing) {
      state.radioVolumeRefreshing = radioVolumeRefreshing
    },
    SET_RADIO_REFRESHING(state, radioRefreshing) {
      state.radioRefreshing = radioRefreshing
    },
    SET_RADIO_UUID(state, radioUUID) {
      state.radioUUID = radioUUID
      localStorage.lastRadioUUID = radioUUID;
    },
    SET_RADIOS(state, radios) {
      state.radios = radios
      for (let i in radios) {
        if (state.radios[i].uuid == state.radioUUID) {
          return
        }
      }
      state.radioUUID = ''
    },
    SET_RADIOS_DISCOVERING(state, radiosDiscovering) {
      state.radiosDiscovering = radiosDiscovering
    },
    SET_RADIOS_LOADING(state, radiosLoading) {
      state.radiosLoading = radiosLoading
    },
    SET_RADIO_WS_CONNECTING(state, radioWSConnecting) {
      state.radioWSConnecting = radioWSConnecting
    },
    SET_RADIO_WS_CONNECTED(state, radioWSConnected) {
      state.radioWSConnected = radioWSConnected
    },
    SET_RADIO_WS(state, radioWS) {
      state.radioWS = radioWS
    }
  },
  actions: {
    initRadio({ dispatch, state }) {
      return dispatch('listRadios').then(() => {
        if (localStorage.lastRadioUUID) {
          for (let radio of state.radios) {
            if (radio.uuid == localStorage.lastRadioUUID) {
              dispatch('setRadioUUID', radio.uuid);
              return
            }
          }
        }
      })
    },
    discoverRadios({ commit, state, dispatch }) {
      commit("SET_RADIOS_DISCOVERING", true);
      return api.discoverRadios()
        .then(({ ok, error }) => {
          if (ok) {
            return new Promise((resolve) => {
              dispatch("listRadios")
                .then(() => dispatch("addMessage", { type: MESSAGE_SUCCESS, text: "discovered " + state.radios.length + " radios" }))
                .finally(() => resolve())
            })
          } else {
            console.error(error)
            dispatch("addMessage", { type: MESSAGE_ERROR, text: error });
          }
        })
        .catch(error => {
          console.error(error)
          dispatch("addMessage", { type: MESSAGE_ERROR, text: error.message });
        })
        .finally(() => {
          commit("SET_RADIOS_DISCOVERING", false)
        })
    },
    listRadios({ commit, dispatch }) {
      commit("SET_RADIOS_LOADING", true)
      return api.listRadios()
        .then(({ ok, result, error }) => {
          if (ok) {
            commit('SET_RADIOS', result)
          } else {
            console.error(error)
            dispatch("addMessage", { type: MESSAGE_ERROR, text: error });
          }
        })
        .catch(error => {
          console.error(error)
          dispatch("addMessage", { type: MESSAGE_ERROR, text: error.message });
        })
        .finally(() => {
          commit("SET_RADIOS_LOADING", false)
        })
    },
    refreshRadio({ commit, dispatch, state }) {
      commit("SET_RADIO_REFRESHING", true)
      dispatch("refreshRadioWS")
      return api.refreshRadio(state.radioUUID)
        .then(({ ok, error }) => {
          if (!ok) {
            console.error(error)
            dispatch("addMessage", { type: MESSAGE_ERROR, text: error });
            return
          }
          dispatch("addMessage", { type: MESSAGE_SUCCESS, text: "refreshed radio" });
        })
        .catch(error => {
          console.error(error)
          dispatch("addMessage", { type: MESSAGE_ERROR, text: error.message });
        })
        .finally(() => {
          commit("SET_RADIO_REFRESHING", false)
        })
    },
    refreshRadioVolume({ commit, dispatch, state }) {
      commit("SET_RADIO_VOLUME_REFRESHING", true)
      return api.refreshRadioVolume(state.radioUUID)
        .then(({ ok, error }) => {
          if (!ok) {
            console.error(error)
            dispatch("addMessage", { type: MESSAGE_ERROR, text: error });
          }
        })
        .catch(error => {
          console.error(error)
          dispatch("addMessage", { type: MESSAGE_ERROR, text: error.message });
        })
        .finally(() => {
          commit("SET_RADIO_VOLUME_REFRESHING", false)
        })
    },
    toggleRadioPower({ commit, dispatch, state }) {
      let power = !state.radio.power
      return api.patchRadio(state.radioUUID, { power })
        .then(({ ok, error }) => {
          if (ok) {
            commit('SET_RADIO_POWER', power)
          } else {
            console.error(error)
            dispatch("addMessage", { type: MESSAGE_ERROR, text: error });
          }
        })
        .catch(error => {
          console.error(error)
          dispatch("addMessage", { type: MESSAGE_ERROR, text: error.message });
        })
    },
    setRadioVolume({ state, commit, dispatch }, volume) {
      commit("CHANGE_RADIO_VOLUME_CHANGING", 1)
      commit("SET_RADIO_VOLUME", volume)
      return api.patchRadio(state.radioUUID, { volume })
        .then(({ ok, error }) => {
          if (!ok) {
            console.error(error)
            dispatch("addMessage", { type: MESSAGE_ERROR, text: error });
          }
        })
        .catch(error => {
          console.error(error)
          dispatch("addMessage", { type: MESSAGE_ERROR, text: error.message });
        }).finally(() => {
          commit("CHANGE_RADIO_VOLUME_CHANGING", -1)
        })
    },
    setRadioPreset({ commit, state, dispatch }, preset) {
      return api.patchRadio(state.radioUUID, { preset })
        .then(({ ok, error }) => {
          if (ok) {
            commit('MERGE_RADIO', { preset })
          } else {
            console.error(error)
            dispatch("addMessage", { type: MESSAGE_ERROR, text: error });
          }
        })
        .catch(error => {
          console.error(error)
          dispatch("addMessage", { type: MESSAGE_ERROR, text: error.message });
        })
    },
    setRadioUUID({ commit, dispatch }, uuid) {
      commit('SET_RADIO_UUID', uuid)
      commit("SET_RADIO", {})
      dispatch('refreshRadioWS')
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
      }

      let onFirstMessage = function (event) {
        onMessage(event);

        commit("SET_RADIO_WS_CONNECTED", true);
        commit("SET_RADIO_WS_CONNECTING", false);

        ws.removeEventListener("message", onFirstMessage);
        ws.addEventListener("message", onMessage);
      }

      ws.addEventListener("message", onFirstMessage);

      let onDisconnect = function () {
        let wasConnected = state.radioWSConnected;

        commit("SET_RADIO_WS_CONNECTED", false);
        commit("SET_RADIO_WS_CONNECTING", false);

        if (wasConnected) {
          dispatch("refreshRadioWS")
        }
      };

      ws.addEventListener("close", onDisconnect);

      commit("SET_RADIO_WS", ws);
    },
  }
}
