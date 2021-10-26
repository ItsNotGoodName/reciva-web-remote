import { createStore } from 'vuex'
import api from './api'

export default createStore({
	state() {
		return {
			config: null,
			presets: null,
			streams: null,
			radios: null
		}
	},
	mutations: {
		setConfig(state, config) {
			state.config = config
		},
		setStreams(state, streams) {
			state.streams = streams
		},
		setPresets(state, presets) {
			state.presets = presets
		},
		setRadios(state, radios) {
			let rds = {};
			for (let r in radios) {
				rds[radios[r].uuid] = radios[r].name;
			}
			state.radios = rds;
		},
	},
	actions: {
		loadAll({ dispatch, state }) {
			return dispatch('loadConfig').then(() => {
				dispatch('loadRadios')
				dispatch('loadPresets')
				dispatch('loadStreams')
			})
		},
		loadConfig({ commit }) {
			return api.getConfig()
				.then((config) => {
					commit("setConfig", config)
				})
		},
		loadRadios({ commit }) {
			return api.getRadios()
				.then((radios) => {
					commit("setRadios", radios)
				})
		},
		loadStreams({ commit, state }) {
			if (!state.config.presetsEnabled) {
				return Promise.resolve()
			}
			return api.getStreams()
				.then((streams) => {
					commit("setStreams", streams)
				})
		},
		loadPresets({ commit, state }) {
			if (!state.config.presetsEnabled) {
				return Promise.resolve()
			}
			return api.getPresets()
				.then((presets) => {
					commit("setPresets", presets)
				})
		}
	}
})