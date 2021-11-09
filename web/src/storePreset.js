import api from "./api";

export default {
    state: () => ({
        presets: [],
        preset: { url: "", newName: "", newUrl: "" },
        presetChanged: false,
        presetLoading: false,
        presetShow: false,
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
            state.presetChanged = false;
        },
        CLEAR_PRESET(state) {
            state.preset = { url: "", newName: "", newUrl: "" };
            state.presetChanged = false;
        },
        SET_PRESET_NEW_NAME(state, newName) {
            state.preset.newName = newName;
            state.preset.newNameChanged = true;
            state.presetChanged = true;
        },
        SET_PRESET_NEW_URL(state, newUrl) {
            state.preset.newUrl = newUrl;
            state.preset.newUrlChanged = true;
            state.presetChanged = true;
        },
        SET_PRESET_LOADING(state, presetLoading) {
            state.presetLoading = presetLoading;
        },
        SHOW_PRESET(state,) {
            state.presetShow = true;
        },
        HIDE_PRESET(state) {
            state.presetShow = false;
        },
    },
    actions: {
        readPresets({ commit }) {
            return api.readPresets()
                .then((presets) => {
                    commit("SET_PRESETS", presets);
                })
        },
        updatePreset({ commit, state }) {
            return api.updatePreset(state.preset)
                .then((preset) => {
                    commit("ADD_PRESET", preset);
                    commit("SET_PRESET", { ...preset });
                });
        },
        closePreset({ commit }) {
            commit("CLEAR_PRESET");
            commit("HIDE_PRESET");
        },
        openPreset({ commit }, url) {
            commit("SET_PRESET_LOADING", true);
            return api.readPreset(url)
                .then((preset) => {
                    commit("SHOW_PRESET");
                    commit("SET_PRESET", preset)
                })
                .finally(() => commit("SET_PRESET_LOADING", false))
        },
    },
    getters: {

    }
}