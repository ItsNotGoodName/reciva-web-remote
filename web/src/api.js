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
				if (!res.ok) {
					return reject(res.statusText)
				}
				return res.json();
			})
			.then(data => resolve(data))
			.catch(err => reject(err))
	})
}

const emptyResponse = (req) => {
	return new Promise((resolve, reject) => {
		req
			.then((res) => {
				if (!res.ok) {
					return reject(res.statusText)
				}
				resolve()
			})
			.catch(err => reject(err))
	})
}

export default {
	getConfig() {
		return jsonResponse(fetch(API_URL + "/config.json"))
	},
	getPresets() {
		return jsonResponse(fetch(API_URL + "/v1/presets"))
	},
	updatePreset(preset) {
		return jsonResponse(fetch(API_URL + "/v1/preset", { method: "POST", body: JSON.stringify(preset) }))
	},
	getStreams() {
		return jsonResponse(fetch(API_URL + "/v1/streams")
		)
	},
	getStream(sid) {
		return jsonResponse(fetch(API_URL + "/v1/stream/" + sid))
	},
	newStream(stream) {
		return jsonResponse(fetch(API_URL + "/v1/stream/new", { method: "POST", body: JSON.stringify(stream) }))
	},
	updateStream(stream) {
		return jsonResponse(fetch(API_URL + "/v1/stream/" + stream.sid, { method: "POST", body: JSON.stringify(stream) }))
	},
	deleteStream(stream) {
		return emptyResponse(fetch(API_URL + "/v1/stream/" + stream.sid, { method: "DELETE" }))
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
}