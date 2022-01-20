import { useToast } from "vue-toastification";

import api from "../api";

const toast = useToast();

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
    MERGE_PRESETS(state, preset) {
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
    listPresets({ commit, dispatch }) {
      return dispatch("_call", {
        promise: api.listPresets(),
        loadingMutation: "SET_PRESETS_LOADING",
      }).then(({ result }) => {
        commit("SET_PRESETS", result);
      });
    },
    submitPreset({ commit, dispatch, state }) {
      return dispatch("_call", {
        promise: api.updatePreset(state.preset),
        loadingMutation: "SET_PRESET_LOADING",
      }).then(({ result }) => {
        commit("MERGE_PRESETS", result);
        dispatch("hidePreset");
        toast.success("preset updated");
      });
    },
    showPreset({ commit, dispatch, state }, url) {
      if (state.presetLoading) {
        return;
      }

      return dispatch("_call", {
        promise: api.getPreset(url),
        loadingMutation: "SET_PRESET_LOADING",
      }).then(({ result }) => {
        commit("SET_PRESET", result);
        commit("SET_PRESET_VISIBLE", true);
      });
    },
    hidePreset({ commit }) {
      commit("SET_PRESET_VISIBLE", false);
      commit("SET_PRESET", {});
    },
  },
};
