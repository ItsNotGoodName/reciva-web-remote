import api from "./api";

export default {
    state: () => ({
        presets: [],
        presetsLoading: false,
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
        SET_PRESETS_LOADING(state, presetsLoading) {
            state.presetsLoading = presetsLoading;
        },
    },
    actions: {
        readPresets({ commit }) {
            return api.readPresets()
                .then((presets) => {
                    commit("SET_PRESETS", presets);
                })
        },
        updatePreset({ commit }, preset) {
            console.log(preset);
            return api.updatePreset(preset)
                .then((preset) => {
                    commit("ADD_PRESET", preset);
                });
        },
    },
    getters: {

    }
}