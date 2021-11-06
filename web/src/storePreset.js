import api from "./api";
import {
    ErrStreamChanged,
    ErrStreamLoading,
    MsgStreamAdded,
    MsgStreamDeleted,
    MsgStreamUpdated,
} from "./constants";

export default {
    state: () => ({
        presets: [],
        stream: { id: 0, name: "", content: "" },
        streamChanged: false,
        streamLoading: false,
        streamShow: false,
        streams: {},

    }),
    mutations: {
        SET_PRESETS(state, presets) {
            state.presets = presets;
        },
        ADD_PRESET(state, preset) {
            for (let i = 0; i < state.presets.length; i++) {
                if (state.presets[i].uri === preset.uri) {
                    state.presets[i] = preset;
                    return;
                }
            }
            state.preset.push(preset);
        },
        SET_STREAM(state, stream) {
            state.stream = stream;
            state.streamChanged = false;
        },
        CLEAR_STREAM(state) {
            state.stream = { id: 0, name: "", content: "" };
            state.streamChanged = false;
        },
        SET_STREAM_NAME(state, name) {
            state.stream.name = name;
            state.stream.nameChanged = true;
            state.streamChanged = true;
        },
        SET_STREAM_CONTENT(state, content) {
            state.stream.content = content;
            state.stream.contentChanged = true;
            state.streamChanged = true;
        },
        SET_STREAM_LOADING(state, streamLoading) {
            state.streamLoading = streamLoading;
        },
        SHOW_STREAM(state,) {
            state.streamShow = true;
        },
        HIDE_STREAM(state) {
            state.streamShow = false;
        },
        DELETE_STREAM(state, id) {
            delete state.streams[id];
        },
        ADD_STREAM(state, stream) {
            state.streams[stream.id] = { id: stream.id, name: stream.name, };
        },
        SET_STREAMS(state, streams) {
            let sts = {};
            for (let i = 0; i < streams.length; i++) {
                sts[streams[i].id] = { name: streams[i].name, id: streams[i].id };
            }
            state.streams = sts;
        },
    },
    actions: {
        readPresets({ commit }) {
            return api.readPresets()
                .then((presets) => {
                    commit("SET_PRESETS", presets);
                });
        },
        updatePreset({ commit }, preset) {
            return api.updatePreset(preset)
                .then(() => {
                    commit("ADD_PRESET", preset);
                });
        },
        readStreams({ commit, dispatch, state }) {
            return api.readStreams()
                .then((streams) => {
                    commit("SET_STREAMS", streams);
                    if (state.stream.id && !state.streams[state.stream.id]) {
                        return dispatch("closeStream")
                    }
                });
        },
        closeStream({ commit }) {
            commit("CLEAR_STREAM");
            commit("HIDE_STREAM");
        },
        addStream({ dispatch, state, commit }) {
            if (state.streamChanged)
                return dispatch("addNotification", { type: "warning", message: ErrStreamChanged });

            commit("CLEAR_STREAM")
            commit("SHOW_STREAM");
        },
        createStream({ state, commit, dispatch }) {
            if (state.streamLoading)
                return dispatch("addNotification", { type: "error", message: ErrStreamLoading });

            commit("SET_STREAM_LOADING", true);
            return api.createStream(state.stream)
                .then((stream) => {
                    commit("SET_STREAM", stream);
                    commit("ADD_STREAM", stream);
                    return dispatch("readStreams")
                })
                .then(() => dispatch("addNotification", { type: "success", message: MsgStreamAdded, }))
                .finally(() => commit("SET_STREAM_LOADING", false));
        },
        updateStream({ dispatch, state, commit }) {
            if (state.streamLoading)
                return dispatch("addNotification", { type: "warning", message: ErrStreamLoading });

            commit("SET_STREAM_LOADING", true);
            return api.updateStream(state.stream)
                .then((stream) => {
                    commit("SET_STREAM", stream);
                    return dispatch("readStreams")
                })
                .then(() => dispatch("addNotification", { type: "success", message: MsgStreamUpdated }))
                .finally(() => commit("SET_STREAM_LOADING", false));
        },
        readStream({ commit, state, dispatch }, id) {
            if (state.streamChanged)
                return dispatch("addNotification", { type: "warning", message: ErrStreamChanged });
            if (state.streamLoading)
                return dispatch("addNotification", { type: "warning", message: ErrStreamLoading });

            commit("SET_STREAM_LOADING", true);
            return api.readStream(id)
                .then((stream) => {
                    commit("SHOW_STREAM");
                    commit("SET_STREAM", stream)
                })
                .finally(() => commit("SET_STREAM_LOADING", false));
        },
        deleteStream({ dispatch, state, commit }) {
            if (state.streamLoading)
                return dispatch("addNotification", { type: "warning", message: ErrStreamLoading });

            commit("SET_STREAM_LOADING", true);
            return api.deleteStream(state.stream.id)
                .then(() => {
                    commit("DELETE_STREAM", state.stream.id)
                    return dispatch("closeStream")
                })
                .then(() => dispatch("readStreams"))
                .then(() => dispatch("addNotification", { type: "success", message: MsgStreamDeleted }))
                .finally(() => commit("SET_STREAM_LOADING", false));
        },
    },
    getters: {

    }
}