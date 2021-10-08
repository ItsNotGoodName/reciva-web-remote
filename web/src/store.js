import { reactive, readonly } from "vue";

const API_URL = import.meta.env.VITE_API_URL
  ? import.meta.env.VITE_API_URL
  : "";
const WS_URL = import.meta.env.VITE_WS_URL ? import.meta.env.VITE_WS_URL : "";

export default {
  state: reactive({
    radio: {},
    radios: {},
  }),
  getState() {
    return readonly(this.state);
  },
  refreshRadio() {
    return fetch(API_URL + "/v1/radios")
      .then((res) => {
        if (!res.ok) {
          return;
        }
        return res.json();
      })
      .then((data) => {
        let radios = {};
        for (let d in data) {
          radios[data[d].uuid] = data[d].uuid;
        }
        this.state.radios = radios;
        return radios;
      });
  },
  selectRadio(uuid) {
    if (this.ws != undefined) {
      this.ws.close();
    }
    this.ws = new WebSocket(WS_URL + "/v1/radio/" + uuid + "/ws");
    this.ws.addEventListener(
      "message",
      function (event) {
        console.log(event.data);
        this.state.radio = JSON.parse(event.data);
      }.bind(this)
    );
    this.ws.addEventListener(
      "close",
      function (event) {
        console.log(event);
      }.bind(this)
    );
  },
};
