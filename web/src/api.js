const ErrNetwork = 'could not contact server'

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
	readPresets() {
		return jsonResponse(fetch(API_URL + "/v1/presets"))
	},
	readPreset(url) {
		return jsonResponse(fetch(API_URL + "/v1/preset?url=" + url))
	},
	updatePreset(preset) {
		return jsonResponse(fetch(API_URL + "/v1/preset", { method: "POST", body: JSON.stringify(preset) }))
	},
	discoverRadios() {
		return emptyResponse(fetch(API_URL + "/v1/radios", { method: "POST" }))
	},
	getRadios() {
		return jsonResponse(fetch(API_URL + "/v1/radios"))
	},
	refreshRadio(uuid) {
		return emptyResponse(fetch(API_URL + "/v1/radio/" + uuid, { method: "POST" }))
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