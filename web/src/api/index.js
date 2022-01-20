import { WS_URL, API_URL } from "../constants";

const jsonResponse = (req) => req.then((res) => res.json());

//const emptyResponse = (req) => {
//return new Promise((resolve, reject) => {
//req
//.then((res) => resolve({ code: res.status, ok: res.ok, error: res.statusText }))
//.catch((err) => reject(err))
//})
//}

export default {
  listPresets() {
    return jsonResponse(fetch(API_URL + "/v1/presets"));
  },
  getPreset(url) {
    return jsonResponse(fetch(API_URL + "/v1/preset?url=" + url));
  },
  updatePreset(preset) {
    return jsonResponse(
      fetch(API_URL + "/v1/preset", {
        method: "POST",
        body: JSON.stringify(preset),
      })
    );
  },
  discoverRadios() {
    return jsonResponse(fetch(API_URL + "/v1/radios", { method: "POST" }));
  },
  listRadios() {
    return jsonResponse(fetch(API_URL + "/v1/radios"));
  },
  refreshRadio(uuid) {
    return jsonResponse(
      fetch(API_URL + "/v1/radio/" + uuid, { method: "POST" })
    );
  },
  refreshRadioVolume(uuid) {
    return jsonResponse(
      fetch(API_URL + "/v1/radio/" + uuid + "/volume", { method: "POST" })
    );
  },
  patchRadio(uuid, state) {
    return jsonResponse(
      fetch(API_URL + "/v1/radio/" + uuid, {
        method: "PATCH",
        body: JSON.stringify(state),
      })
    );
  },
  getRadioWS(uuid) {
    if (uuid) {
      return new WebSocket(WS_URL + "/v1/radio/ws?uuid=" + uuid);
    }
    return new WebSocket(WS_URL + "/v1/radio/ws");
  },
};
