import { createStore } from 'vuex'
import api from './api'

export default createStore({
	state() {
		return {
			config: null,
			presets: null,
			streams: null,
			radios: null,
			radioWS: null,
			radioUUID: null,
			radio: null
		}
	},
	mutations: {
		SET_CONFIG(state, config) {
			state.config = config
		},
		SET_STREAMS(state, streams) {
			state.streams = streams
		},
		SET_PRESETS(state, presets) {
			state.presets = presets
		},
		SET_RADIO(state, radio) {
			state.radio = radio
		},
		SET_RADIO_UUID(state, uuid) {
			state.radioUUID = uuid
		},
		SET_RADIO_WS(state, radioWS) {
			state.radioWS = radioWS
		},
		SET_RADIOS(state, radios) {
			let rds = {};
			for (let r in radios) {
				rds[radios[r].uuid] = radios[r].name;
			}
			state.radios = rds;
		},
	},
	actions: {
		loadAll({ dispatch }) {
			return dispatch('loadConfig').then(() => {
				dispatch('loadRadios')
				dispatch('loadPresets')
				dispatch('loadStreams')
			})
		},
		loadConfig({ commit }) {
			return api.getConfig()
				.then((config) => {
					commit("SET_CONFIG", config)
				})
		},
		loadRadios({ commit }) {
			return api.getRadios()
				.then((radios) => {
					commit("SET_RADIOS", radios)
				})
		},
		loadStreams({ commit, state }) {
			if (!state.config.presetsEnabled) {
				return Promise.reject()
			}
			return api.getStreams()
				.then((streams) => {
					commit("SET_STREAMS", streams)
				})
		},
		loadPresets({ commit, state }) {
			if (!state.config.presetsEnabled) {
				return Promise.reject()
			}
			return api.getPresets()
				.then((presets) => {
					commit("SET_PRESETS", presets)
				})
		},
		setRadioUUID({ commit, state }, uuid) {
			let ws
			if (!state.radioWS) {
				ws = api.radioWS(state.radioUUID)

				// Handle messsage
				ws.addEventListener(
					"message",
					function (event) {
						let radio = JSON.parse(event.data);
						if (radio.uuid != state.radioUUID) return;
						commit("SET_RADIO", radio)
					}
				);

				// Handle open
				ws.addEventListener(
					"open",
					function (event) {
						console.log(event)
						ws.send(uuid)
					}
				);

				// Handle close
				ws.addEventListener(
					"close",
					function (event) {
						console.log(event);
					}
				);

				// Handle error
				ws.addEventListener(
					"error",
					function (event) {
						console.error(event);
					}
				);

				commit("SET_RADIO_WS", ws)
			} else {
				state.radioWS.send(uuid)
			}

			commit("SET_RADIO_UUID", uuid)
		}
	}
})