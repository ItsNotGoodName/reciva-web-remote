import { computed, watch, shallowReactive, ref, Ref } from "vue";

import { WS_URL } from "../constants"

const subscribe = (ws: WebSocket, radioUUID: string) => {
  ws.send(JSON.stringify({ type: "radio.subscribe", slug: radioUUID }));
}

const unsubscribe = (ws: WebSocket) => {
  ws.send(JSON.stringify({ type: "radio.unsubscribe" }));
}

const initialRadio: Radio = {
  audio_source: "",
  audio_sources: [],
  is_muted: false,
  metadata: "",
  model_name: "",
  model_number: "",
  name: "",
  power: false,
  preset_number: 0,
  presets: [],
  status: "",
  title: "",
  title_new: "",
  url: "",
  url_new: "",
  uuid: "",
  volume: 0,
}

export function useWS(radioUUID: Ref<string>) {
  const connecting = ref(true);
  const connected = ref(false);
  const disconnected = ref(false);
  const radio = shallowReactive<Radio>({ ...initialRadio });
  const radioSelected = computed(() => radio.uuid != "")
  const radioLoading = computed(() => (radio.uuid != radioUUID.value) || connecting.value)

  const connect = () => {
    let ws = new WebSocket(WS_URL + "/api/ws");
    connecting.value = true;

    ws.addEventListener("open", () => {
      connecting.value = false;
      connected.value = true;
      disconnected.value = false;

      if (radioUUID.value) {
        subscribe(ws, radioUUID.value);
      }
    });

    ws.addEventListener("message", (event) => {
      let msg = JSON.parse(event.data) as { type: string, slug: any };
      if (msg.type == "radio.partial" || msg.type == "radio") {
        Object.assign(radio, msg.slug);
      }
    });

    ws.addEventListener("close", () => {
      connecting.value = false
      connected.value = false
      disconnected.value = true;
      Object.assign(radio, initialRadio);
    });

    return ws
  }

  let ws = connect()

  const reconnect = () => {
    if (connected.value || connecting.value) {
      return false
    }

    ws = connect()
    return true
  }

  watch(radioUUID, () => {
    if (reconnect()) {
      return
    }

    Object.assign(radio, initialRadio);
    if (radioUUID.value) {
      subscribe(ws, radioUUID.value)
    } else {
      unsubscribe(ws)
    }
  })

  return { radio, radioLoading, radioSelected, connecting, connected, disconnected, reconnect }
}
