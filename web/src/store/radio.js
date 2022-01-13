import api from '../api'

export default {
  state() {
    return {
      radio: {},
      radioRefreshing: false,
      radioVolumeRefreshing: false,
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
      return state.radioUUID
    },
    radioLoaded(state) {
      return state.radio.uuid == state.radioUUID
    }
  },
  mutations: {
    MERGE_RADIO(state, radio) {
      for (let k in radio) {
        state.radio[k] = radio[k]
      }
    },
    SET_RADIO_POWER(state, power) {
      state.radio.power = power
    },
    SET_RADIO_VOLUME(state, volume) {
      state.radio.volume = volume
    },
    SET_RADIO_REFRESHING(state, radioRefreshing) {
      state.radioRefreshing = radioRefreshing
    },
    SET_RADIO_VOLUME_REFRESHING(state, radioVolumeRefreshing) {
      state.radioVolumeRefreshing = radioVolumeRefreshing
    },
    SET_RADIO_UUID(state, radioUUID) {
      state.radioUUID = radioUUID
    },
    SET_RADIOS(state, radios) {
      state.radios = radios
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
    discoverRadios({ commit, dispatch }) {
      commit("SET_RADIOS_DISCOVERING", true);
      return api.discoverRadios()
        .then(({ ok, error }) => {
          if (ok) {
            dispatch("listRadios");
          } else {
            console.error(error)
          }
        })
        .catch(error => {
          console.error(error)
        })
        .finally(() => {
          commit("SET_RADIOS_DISCOVERING", false)
        })
    },
    listRadios({ commit }) {
      commit("SET_RADIOS_LOADING", true)
      return api.listRadios()
        .then(({ ok, result, error }) => {
          if (ok) {
            commit('SET_RADIOS', result)
          } else {
            console.error(error)
          }
        })
        .catch(error => {
          console.error(error)
        })
        .finally(() => {
          commit("SET_RADIOS_LOADING", false)
        })
    },
    refreshRadio({ commit, state }) {
      commit("SET_RADIO_REFRESHING", true)
      return api.refreshRadio(state.radioUUID)
        .then(({ ok, error }) => {
          if (!ok) {
            console.error(error)
          }
        })
        .catch(error => {
          console.error(error)
        })
        .finally(() => {
          commit("SET_RADIO_REFRESHING", false)
        })
    },
    refreshRadioVolume({ commit, state }) {
      commit("SET_RADIO_VOLUME_REFRESHING", true)
      return api.refreshRadioVolume(state.radioUUID)
        .then(({ ok, error }) => {
          if (!ok) {
            console.error(error)
          }
        })
        .catch(error => {
          console.error(error)
        })
        .finally(() => {
          commit("SET_RADIO_VOLUME_REFRESHING", false)
        })
    },
    toggleRadioPower({ commit, state }) {
      power = !state.radio.power
      return api.patchRadio(state.radioUUID, { power })
        .then(({ ok, error }) => {
          if (ok) {
            commit('SET_RADIO_POWER', power)
          } else {
            console.error(error)
          }
        })
        .catch(error => {
          console.error(error)
        })
    },
    setRadioVolume({ commit }, volume) {
      return api.patchRadio(state.radioUUID, { volume })
        .then(({ ok, error }) => {
          if (ok) {
            commit('SET_RADIO_VOLUME', volume)
          } else {
            console.error(error)
          }
        })
        .catch(error => {
          console.error(error)
        })
    },
    setRadioPreset({ commit }, preset) {
      return api.patchRadio(state.radioUUID, { preset })
        .then(({ ok, error }) => {
          if (ok) {
            commit('MERGE_RADIO', { preset })
          } else {
            console.error(error)
          }
        })
        .catch(error => {
          console.error(error)
        })
    },
    setRadioUUID({ commit, dispatch }, uuid) {
      commit('SET_RADIO_UUID', uuid)
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
        commit("SET_RADIO_WS_CONNECTED", false);
      };

      ws.addEventListener("close", onDisconnect);

      commit("SET_RADIO_WS", ws);
    },
  }
}
