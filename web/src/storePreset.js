import api from "./api";

export default {
    state: () => ({
        presets: [],
        preset: {}
    }),
    mutations: {
        SET_PRESETS(state, presets) {
            state.presets = presets;
        },
        ADD_PRESET(state, preset) {
            for (let i = 0; i < state.presets.length; i++) {
                if (state.presets[i].url == preset.url) {
                    state.presets[i] = preset;
                    return;
                }
            }
            state.presets.push(preset);
        },
        SET_PRESET(state, preset) {
            state.preset = preset;
        }
    },
    actions: {
        readPresets({ commit }) {
            return api.readPresets()
                .then((presets) => {
                    commit("SET_PRESETS", presets);
                })
        },
        updatePreset({ commit }, preset) {
            return api.updatePreset(preset)
                .then((preset) => {
                    commit("ADD_PRESET", preset);
                });
        },
        readPreset({ commit, rootState }) {
            return api.readPreset(rootState.radio.url)
                .then((preset) => {
                    commit("SET_PRESET", preset);
                }).catch(() => {
                    commit("SET_PRESET", {});
                });
        }
    },
    getters: {

    }
}