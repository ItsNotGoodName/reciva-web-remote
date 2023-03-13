import { type Accessor, createSignal, createEffect, on, batch } from "solid-js";
import { type Store, createStore, produce } from "solid-js/store";
import {
  PubsubTopic,
  type WsCommand,
  type StateState,
  StateStatus,
  type WsEvent,
} from "./api";
import { WS_URL } from "./constants";

const subscribe = (ws: WebSocket, radioUUID: string) => {
  const topics: Array<PubsubTopic> = [PubsubTopic.DiscoverTopic];
  if (radioUUID != "") {
    topics.push(PubsubTopic.StateTopic);
  }

  ws.send(
    JSON.stringify({
      state: { partial: true, uuid: radioUUID },
      subscribe: { topics: topics },
    } as WsCommand)
  );
};

const defaultState: StateState = {
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
  status: StateStatus.StatusUnknown,
  title: "",
  title_new: "",
  url: "",
  url_new: "",
  uuid: "",
  volume: 0,
};

export type WSDataReturn = {
  state: Store<StateState>;
  discovering: Accessor<boolean>;
};

export type WSStatusReturn = {
  connecting: Accessor<boolean>;
  connected: Accessor<boolean>;
  disconnected: Accessor<boolean>;
  synced: Accessor<boolean>;
  reconnect: () => void;
};

export type WSReturn = [WSDataReturn, WSStatusReturn];

export function useWS(radioUUID: Accessor<string>): WSReturn {
  const [connecting, setConnecting] = createSignal(true);
  const [connected, setConnected] = createSignal(false);
  const [disconnected, setDisconnected] = createSignal(false);
  const [synced, setSynced] = createSignal(false);

  const [discovering, setDiscovering] = createSignal(false);
  const [state, setState] = createStore(defaultState);

  const connect = () => {
    const ws = new WebSocket(WS_URL);
    batch(() => {
      setConnecting(true);
      setState(defaultState);
    });
    console.log("WS: Connecting");

    ws.addEventListener("open", () => {
      console.log("WS: Connected");
      batch(() => {
        setConnecting(false);
        setConnected(true);
        setDisconnected(false);
      });
      subscribe(ws, radioUUID());
    });

    ws.addEventListener("message", (event) => {
      console.log("WS: Message");
      batch(() => {
        const msg = JSON.parse(event.data as string) as WsEvent;
        switch (msg.topic) {
          case PubsubTopic.StateTopic:
            const data = msg.data as StateState;
            if (data.uuid == radioUUID()) {
              setState(
                produce((state) => Object.assign(state, msg.data as StateState))
              );
            }
            break;
          case PubsubTopic.DiscoverTopic:
            setDiscovering(msg.data as boolean);
            break;
        }
        setSynced(true);
      });
    });

    ws.addEventListener("close", () => {
      console.log("WS: Close");
      batch(() => {
        setConnecting(false);
        setConnected(false);
        setDisconnected(true);
        setSynced(false);
        setState(defaultState);
      });
    });

    return ws;
  };

  let ws = connect();

  const reconnect = () => {
    if (connected() || connecting()) {
      return;
    }

    ws = connect();
  };

  createEffect(
    on(radioUUID, () => {
      if (!connected()) {
        reconnect();
        return;
      }
      if (connecting()) {
        return;
      }

      setState(defaultState);
      subscribe(ws, radioUUID());
    })
  );

  return [
    {
      state,
      discovering,
    },
    {
      connecting,
      connected,
      disconnected,
      reconnect,
      synced,
    },
  ];
}
