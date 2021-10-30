import { createStore } from 'vuex'
import api from './api'

export default createStore({
	state() {
		return {
			config: null,
			presets: null,
			radio: {},
			radioUUID: null,
			radioWS: null,
			radios: null,
			streams: null,
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
			for (let k in radio) {
				state.radio[k] = radio[k]
			}
		},
		SET_RADIO_POWER(state, power) {
			state.radio.power = power
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
				return Promise.resolve()
			}
			return api.getStreams()
				.then((streams) => {
					commit("SET_STREAMS", streams)
				})
		},
		loadPresets({ commit, state }) {
			if (!state.config.presetsEnabled) {
				return Promise.resolve()
			}
			return api.getPresets()
				.then((presets) => {
					commit("SET_PRESETS", presets)
				})
		},
		refreshRadio({ dispatch, state }) {
			if (!state.radioUUID) {
				return Promise.resolve()
			}
			if (!state.radioWS) {
				dispatch("setRadioUUID", state.radioUUID)
			}
			return api.renewRadio(state.radioUUID)
		},
		setRadioUUID({ commit, state }, uuid) {
			let ws
			if (!state.radioWS) {
				ws = api.radioWS(uuid)

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
					function () {
						if (uuid) {
							ws.send(uuid)
						}
					}
				);

				// Handle close
				ws.addEventListener(
					"close",
					function (event) {
						console.log(event);
						commit("SET_RADIO_WS", null)
					}
				);

				// Handle error
				ws.addEventListener(
					"error",
					function (event) {
						console.error(event);
						commit("SET_RADIO_WS", null)
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