import { watch, ref, Ref } from "vue";

import { WS_URL } from "../constants"

const subscribe = (ws: WebSocket, radioUUID: string) => {
  ws.send(JSON.stringify({ type: "state.subscribe", slug: radioUUID }));
}

const unsubscribe = (ws: WebSocket) => {
  ws.send(JSON.stringify({ type: "state.unsubscribe" }));
}


export function useWS(radioUUID: Ref<string>) {
  const connecting = ref(true);
  const connected = ref(false);
  const disconnected = ref(false);
  const radio = ref<Radio | undefined>(undefined);
  let failCount = 0;

  const connect = () => {
    let ws = new WebSocket(WS_URL + "/api/ws");
    connecting.value = true;

    ws.addEventListener("open", () => {
      connecting.value = false;
      connected.value = true;
      disconnected.value = false;
      failCount = 0;

      if (radioUUID.value) {
        subscribe(ws, radioUUID.value);
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
      failCount++;

      if (failCount < 5) {
        setTimeout(reconnect, 2000 * failCount);
      }
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

  watch(radioUUID, () => {
    if (!connected.value) {
      reconnect()
      return
    }

    radio.value = undefined;
    if (radioUUID.value) {
      subscribe(ws, radioUUID.value)
    } else {
      unsubscribe(ws)
    }
  })

  return { radio, connecting, connected, disconnected, reconnect }
}