import { computed, watch, shallowReactive, ref, Ref } from "vue";

import { WS_URL } from "../constants"

const subscribe = (ws: WebSocket, radioUUID: string) => {
  ws.send(JSON.stringify({ type: "state.subscribe", slug: radioUUID }));
}

const unsubscribe = (ws: WebSocket) => {
  ws.send(JSON.stringify({ type: "state.unsubscribe" }));
}

const initialState: State = {
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

export function useWS(stateUUID: Ref<string>) {
  const connecting = ref(true);
  const connected = ref(false);
  const disconnected = ref(false);
  const state = shallowReactive<State>({ ...initialState });
  const stateSelected = computed(() => state.uuid != "")
  const stateLoading = computed(() => (state.uuid != stateUUID.value) || connecting.value)

  const connect = () => {
    let ws = new WebSocket(WS_URL + "/api/ws");
    connecting.value = true;

    ws.addEventListener("open", () => {
      connecting.value = false;
      connected.value = true;
      disconnected.value = false;

      if (stateUUID.value) {
        subscribe(ws, stateUUID.value);
      }
    });

    ws.addEventListener("message", (event) => {
      let msg = JSON.parse(event.data) as { type: string, slug: any };
      if (msg.type == "state.partial" || msg.type == "state") {
        Object.assign(state, msg.slug);
      }
    });

    ws.addEventListener("close", () => {
      connecting.value = false
      connected.value = false
      disconnected.value = true;
      Object.assign(state, initialState);
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

  watch(stateUUID, () => {
    if (reconnect()) {
      return
    }

    Object.assign(state, initialState);
    if (stateUUID.value) {
      subscribe(ws, stateUUID.value)
    } else {
      unsubscribe(ws)
    }
  })

  return { state, stateLoading, stateSelected, connecting, connected, disconnected, reconnect }
}
