import { createStore } from 'vuex'
import api from './api'

export default createStore({
	state() {
		return {
			config: null,
			radio: {},
			radioUUID: null,
			radioConnected: false,
			radioConnecting: false,
			radioWS: null,
			radios: null,
		}
	},
	mutations: {
		SET_CONFIG(state, config) {
			state.config = config
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
		SET_RADIO_CONNECTING(state, radioConnecting) {
			state.radioConnecting = radioConnecting
			if (radioConnecting) {
				state.radioConnected = false
			}
		},
		SET_RADIO_CONNECTED(state, radioConnected) {
			state.radioConnected = radioConnected
			state.radioConnecting = false
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
		refreshRadio({ dispatch, state }) {
			if (!state.radioUUID) {
				return Promise.resolve()
			}
			dispatch("refreshRadioWS")
			return api.renewRadio(state.radioUUID)
		},
		setRadioUUID({ commit, state, dispatch }, uuid) {
			if (uuid != state.radioUUID) {
				commit("SET_RADIO_UUID", uuid)
				dispatch("refreshRadioWS")
			}
		},
		refreshRadioWS({ state, commit, dispatch }) {
			// Full state update when radio websocket is connected
			if (state.radioConnected) {
				if (state.radioUUID) {
					state.radioWS.send(state.radioUUID)
				}
				return
			}

			// Do not create a new websocket if current websocket is connecting or radioUUID is not set
			if (state.radioConnecting || !state.radioUUID) {
				return
			}

			commit("SET_RADIO_CONNECTING", true)
			let ws = api.radioWS(state.radioUUID)

			// Handle messsage
			ws.addEventListener(
				"message",
				function (event) {
					let radio = JSON.parse(event.data);
					if (radio.uuid != state.radioUUID) return;
					commit("SET_RADIO", radio)
					commit("SET_RADIO_CONNECTED", true)
				}
			);

			let onEnd = function (event) {
				console.log(event)
				commit("SET_RADIO_CONNECTED", false)
				setTimeout(() => {
					dispatch("refreshRadioWS")
				}, 5000)
			}

			// Handle close
			ws.addEventListener("close", onEnd);

			// Handle error
			ws.addEventListener("error", onEnd);

			commit("SET_RADIO_WS", ws)
		}
	}
})