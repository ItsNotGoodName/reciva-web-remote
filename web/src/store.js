import { reactive, readonly } from "vue";

const API_URL = import.meta.env.VITE_API_URL
  ? import.meta.env.VITE_API_URL
  : "";
const WS_URL = import.meta.env.VITE_WS_URL ? import.meta.env.VITE_WS_URL : (() => {
  if (window.location.protocol == "http:") {
    return "ws://" + window.location.host
  }
  return "wss://" + window.location.host
})();

export default {
  state: reactive({
    connecting: false,
    connected: false,
    uuid: "",
    radio: {},
    radios: {},
  }),
  getState() {
    return readonly(this.state);
  },
  updateRadios() {
    return fetch(API_URL + "/v1/radios")
      .then((res) => {
        if (!res.ok) {
          throw Error(res.statusText);
        }
        return res.json();
      })
      .then((data) => {
        let radios = {};
        for (let d in data) {
          radios[data[d].uuid] = data[d].name;
        }
        this.state.radios = radios;
      })
      .catch((error) => {
        console.log(error);
      });
  },
  discoverRadios() {
    fetch(API_URL + "/v1/radios", {
      method: "POST",
    })
      .then(() => {
        this.updateRadios();
      })
      .catch((error) => {
        console.log(error);
      });
  },
  renewRadio(){
    if (!this.state.uuid) return;
    fetch(API_URL + "/v1/radio/" + this.state.radio.uuid + "/renew", {
      method: "POST",
    }).catch((error) => {
      console.log(error);
    });
  },
  refreshRadioVolume(){
    if (!this.state.uuid) return;
    fetch(API_URL + "/v1/radio/" + this.state.radio.uuid + "/volume", {
      method: "POST",
    }).catch((error) => {
      console.log(error);
    });
  },
  toggleRadioPower() {
    if (!this.state.uuid) return;
    fetch(API_URL + "/v1/radio/" + this.state.radio.uuid, {
      method: "PATCH",
      body: JSON.stringify({ power: !this.state.radio.power }),
    }).catch((error) => {
      console.log(error);
    });
  },
  setRadioPreset(preset) {
    if (!this.state.uuid) return;
    fetch(API_URL + "/v1/radio/" + this.state.radio.uuid, {
      method: "PATCH",
      body: JSON.stringify({ preset: preset }),
    }).catch((error) => {
      console.log(error);
    });
  },
  setRadioVolume(volume) {
    if (!this.state.uuid) return;
    fetch(API_URL + "/v1/radio/" + this.state.radio.uuid, {
      method: "PATCH",
      body: JSON.stringify({ volume: volume }),
    }).catch((error) => {
      console.log(error);
    });
  },
  selectRadio(uuid = "") {
    this.state.uuid = uuid;
    if (!this.initWS()) return;
    this.ws.send(uuid);
  },
  initWS() {
    if (this.state.connecting) {
      return false;
    }
    if (this.state.connected) {
      return true;
    }
    this.state.connecting = true;

    // Create websocket
    if (this.state.uuid == undefined || this.state.uuid == "") {
      this.ws = new WebSocket(WS_URL + "/v1/radio/ws");
    } else {
      this.ws = new WebSocket(WS_URL + "/v1/radio/" + this.state.uuid + "/ws");
    }

    // Handle radio state message
    this.ws.addEventListener(
      "message",
      function (event) {
        let radio = JSON.parse(event.data);
        if (radio.uuid != this.state.uuid) return;
        this.state.radio = radio;
      }.bind(this)
    );

    // Handle open
    this.ws.addEventListener(
      "open",
      function () {
        this.state.connected = true;
        this.state.connecting = false;
      }.bind(this)
    );

    // Handle close
    this.ws.addEventListener(
      "close",
      function (event) {
        this.state.connected = false;
        this.state.connecting = false;
      }.bind(this)
    );

    // Handle error
    this.ws.addEventListener(
      "error",
      function (event) {
        this.state.connected = false;
        this.state.connecting = false;
        console.error(event);
      }.bind(this)
    );
  },
};
