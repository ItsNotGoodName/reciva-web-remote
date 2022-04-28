import { watch, ref, Ref } from "vue";
import { WS_URL } from "../constants"

const subscribe = (ws: WebSocket, uuid: string) => {
  ws.send(JSON.stringify({ type: "state.subscribe", slug: uuid }));
}

const unsubscribe = (ws: WebSocket) => {
  ws.send(JSON.stringify({ type: "state.unsubscribe" }));
}


export function useWS(uuid: Ref<string>) {
  const connecting = ref(true);
  const connected = ref(false);
  const disconnected = ref(false);
  const radio = ref<Radio | undefined>(undefined);

  const connect = () => {
    let ws = new WebSocket(WS_URL + "/api/ws");
    connecting.value = true;

    ws.addEventListener("open", () => {
      connecting.value = false;
      connected.value = true;
      disconnected.value = false;

      if (uuid.value) {
        subscribe(ws, uuid.value);
      }
    });

    ws.addEventListener("message", (event) => {
      let msg = JSON.parse(event.data) as { type: string, slug: Radio };
      if (msg.type == "state.partial" && radio.value) {
        radio.value = { ...radio.value, ...msg.slug };
      } else if (msg.type == "state") {
        radio.value = msg.slug;
      }
    });

    ws.addEventListener("close", () => {
      connecting.value = false
      connected.value = false
      disconnected.value = true;
      radio.value = undefined;

      setTimeout(() => reconnect(), 5000);
    });

    return ws
  }

  let ws = connect()

  const reconnect = () => {
    if (connected.value || connecting.value) {
      return
    }
    ws = connect()
  }

  watch(uuid, () => {
    if (!connected.value) {
      reconnect()
      return
    }

    radio.value = undefined;
    if (uuid.value) {
      subscribe(ws, uuid.value)
    } else {
      unsubscribe(ws)
    }
  })

  return { radio, connecting, connected, disconnected, reconnect }
}