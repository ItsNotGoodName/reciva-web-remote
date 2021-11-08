import { ErrNetwork } from "./constants"

const API_URL = import.meta.env.VITE_API_URL
	? import.meta.env.VITE_API_URL
	: "";
const WS_URL = import.meta.env.VITE_WS_URL ? import.meta.env.VITE_WS_URL : (() => {
	if (window.location.protocol == "http:") {
		return "ws://" + window.location.host
	}
	return "wss://" + window.location.host
})();

const jsonResponse = (req) => {
	return new Promise((resolve, reject) => {
		req
			.then((res) => {
				res.json()
					.then((data) => { res.ok ? resolve(data) : reject(data.err) })
					.catch(() => reject(res.statusText))
			})
			.catch(() => reject(ErrNetwork))
	})
}

const emptyResponse = (req) => {
	return new Promise((resolve, reject) => {
		req
			.then((res) => {
				if (res.ok) {
					return resolve()
				}
				res.json()
					.then((data) => reject(data.err))
					.catch(() => reject(res.statusCode))
			})
			.catch(() => reject(ErrNetwork))
	})
}

export default {
	getConfig() {
		return jsonResponse(fetch(API_URL + "/v1/config"))
	},
	readPresets() {
		return jsonResponse(fetch(API_URL + "/v1/presets"))
	},
	updatePreset(preset) {
		return jsonResponse(fetch(API_URL + "/v1/preset", { method: "POST", body: JSON.stringify(preset) }))
	},
	clearPreset(preset) {
		return jsonResponse(fetch(API_URL + "/v1/preset", { method: "DELETE", body: JSON.stringify(preset) }))
	},
	createStream(stream) {
		return jsonResponse(fetch(API_URL + "/v1/stream/new", { method: "POST", body: JSON.stringify(stream) }))
	},
	readStreams() {
		return jsonResponse(fetch(API_URL + "/v1/streams"))
	},
	readStream(id) {
		return jsonResponse(fetch(API_URL + "/v1/stream/" + id))
	},
	updateStream(stream) {
		return jsonResponse(fetch(API_URL + "/v1/stream/" + stream.id, { method: "POST", body: JSON.stringify(stream) }))
	},
	deleteStream(id) {
		return emptyResponse(fetch(API_URL + "/v1/stream/" + id, { method: "DELETE" }))
	},
	discoverRadios() {
		return emptyResponse(fetch(API_URL + "/v1/radios", { method: "POST" }))
	},
	getRadios() {
		return jsonResponse(fetch(API_URL + "/v1/radios"))
	},
	renewRadio(uuid) {
		return emptyResponse(fetch(API_URL + "/v1/radio/" + uuid + "/renew", { method: "POST" }))
	},
	refreshRadioVolume(uuid) {
		return emptyResponse(fetch(API_URL + "/v1/radio/" + uuid + "/volume", { method: "POST" }))
	},
	updateRadio(uuid, state) {
		return emptyResponse(fetch(API_URL + "/v1/radio/" + uuid, { method: "PATCH", body: JSON.stringify(state) }))
	},
	radioWS(uuid) {
		if (uuid == undefined || uuid == "") {
			return new WebSocket(WS_URL + "/v1/radio/ws")
		}
		return new WebSocket(WS_URL + "/v1/radio/ws?uuid=" + uuid)
	}
}