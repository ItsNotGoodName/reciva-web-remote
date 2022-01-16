import api from "../api";

export default {
  state: () => ({
    preset: {},
    presetLoading: false,
    presetVisible: false,
    presets: [],
    presetsLoading: false,
  }),
  mutations: {
    SET_PRESET(state, preset) {
      state.preset = preset;
    },
    SET_PRESET_LOADING(state, loading) {
      state.presetLoading = loading;
    },
    SET_PRESET_VISIBLE(state, presetVisible) {
      state.presetVisible = presetVisible;
    },
    SET_PRESET_NEW_NAME(state, newName) {
      state.preset.newName = newName;
    },
    SET_PRESET_NEW_URL(state, newUrl) {
      state.preset.newUrl = newUrl;
    },
    SET_PRESETS_LOADING(state, presetLoading) {
      state.presetsLoading = presetLoading;
    },
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
    listPresets({ commit, state }) {
      if (state.presetsLoading) {
        return Promise.resolve();
      }

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
    submitPreset({ commit, state }) {
      if (state.presetLoading) {
        return Promise.resolve();
      }

      commit("SET_PRESET_LOADING", true);
      return api.updatePreset(state.preset)
        .then(({ ok, result, error }) => {
          if (ok) {
            commit("MERGE_PRESET", result);
            commit("SET_PRESET_VISIBLE", false);
          } else {
            console.error(error);
          }
        })
        .catch(err => {
          console.error(err);
        })
        .finally(() => {
          commit("SET_PRESET_LOADING", false);
        })
    },
    showPreset({ commit, state }, url) {
      if (state.presetLoading) {
        return Promise.resolve();
      }

      commit("SET_PRESET_LOADING", true);
      return api.getPreset(url)
        .then(({ ok, result, error }) => {
          if (ok) {
            commit("SET_PRESET", result);
            commit("SET_PRESET_VISIBLE", true);
          } else {
            console.error(error);
          }
        })
        .catch(err => {
          console.error(err);
        })
        .finally(() => {
          commit("SET_PRESET_LOADING", false);
        })
    },
    hidePreset({ commit }) {
      commit("SET_PRESET_VISIBLE", false);
      commit("SET_PRESET", {});
    }
  },
}