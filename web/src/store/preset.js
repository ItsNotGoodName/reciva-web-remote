import api from "../api";

export default {
  state: () => ({
    presets: [],
    presetsLoading: false,
  }),
  mutations: {
    SET_PRESETS(state, presets) {
      state.presets = presets;
    },
    SET_PRESETS_LOADING(state, presetsLoading) {
      state.presetsLoading = presetsLoading;
    },
    MERGE_PRESET(state, preset) {
      for (let i = 0; i < state.presets.length; i++) {
        if (state.presets[i].url == preset.url) {
          state.presets[i] = preset;
          return;
        }
      }
      state.presets.push(preset);
    },
  },
  actions: {
    listPresets({ commit }) {
      commit("SET_PRESETS_LOADING", true);
      return api.listPresets()
        .then(({ ok, result, error }) => {
          if (ok) {
            commit("SET_PRESETS", result);
          } else {
            console.error(error);
          }
        })
        .catch(err => {
          console.error(err);
        })
        .finally(() => {
          commit("SET_PRESETS_LOADING", false);
        })
    },
    updatePreset({ commit }, preset) {
      return api.updatePreset(preset)
        .then(({ ok, result, error }) => {
          if (ok) {
            commit("MERGE_PRESET", result);
          } else {
            console.error(error);
          }
        })
        .catch(err => {
          console.error(err);
        })
    },
  },
}