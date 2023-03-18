import { type Accessor, createSignal, createEffect, on, batch } from "solid-js";
import { type Store, createStore } from "solid-js/store";
import {
  PubsubTopic,
  type WsCommand,
  type StateState,
  StateStatus,
  type WsEvent,
  type ModelStale,
} from "./api";
import { WS_URL } from "./constants";

const sendCommand = (ws: WebSocket, radioUUID: string) => {
  const topics: Array<PubsubTopic> = [
    PubsubTopic.DiscoverTopic,
    PubsubTopic.StaleTopic,
    PubsubTopic.StateTopic,
  ];

  ws.send(
    JSON.stringify({
      state: { partial: true, uuid: radioUUID },
      subscribe: { topics: topics },
    } as WsCommand)
  );
};

const sendStateCommand = (ws: WebSocket, radioUUID: string) => {
  ws.send(
    JSON.stringify({
      state: { partial: true, uuid: radioUUID },
    } as WsCommand)
  );
};

const DefaultState: StateState = {
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
  stale: Accessor<ModelStale | undefined>;
};

export type WSStatusReturn = {
  connecting: Accessor<boolean>;
  connected: Accessor<boolean>;
  disconnected: Accessor<boolean>;
  reconnect: () => void;
};

export type WSReturn = [WSDataReturn, WSStatusReturn];

export function useWS(radioUUID: Accessor<string>): WSReturn {
  const [connecting, setConnecting] = createSignal(true);
  const [connected, setConnected] = createSignal(false);
  const [disconnected, setDisconnected] = createSignal(false);
  const [shouldReconnect, setShouldReconnect] = createSignal(false);

  const [stale, setStale] = createSignal<ModelStale | undefined>(undefined, {
    equals: false,
  });
  const [discovering, setDiscovering] = createSignal(false);
  const [state, setState] = createStore({ ...DefaultState });

  const connect = () => {
    const ws = new WebSocket(WS_URL);
    batch(() => {
      setConnecting(true);
      setState(DefaultState);
    });
    console.log("WS: Connecting");

    ws.addEventListener("open", () => {
      console.log("WS: Connected");
      batch(() => {
        setConnecting(false);
        setConnected(true);
        setDisconnected(false);
        setShouldReconnect(false);
      });
      sendCommand(ws, radioUUID());
    });

    ws.addEventListener("message", (event) => {
      console.log("WS: Message");
      const msg = JSON.parse(event.data as string) as WsEvent;
      switch (msg.topic) {
        case PubsubTopic.StateTopic:
          const data = msg.data as StateState;
          if (data.uuid == radioUUID()) {
            setState(data);
          }
          break;
        case PubsubTopic.DiscoverTopic:
          setDiscovering(msg.data as boolean);
          break;
        case PubsubTopic.StaleTopic:
          setStale(msg.data as ModelStale);
          break;
      }
    });

    ws.addEventListener("close", () => {
      console.log("WS: Close");
      batch(() => {
        setConnecting(false);
        setConnected(false);
        setDisconnected(true);
        setState(DefaultState);
        setShouldReconnect(true);
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
    on(
      radioUUID,
      () => {
        if (!connected()) {
          reconnect();
          return;
        }
        if (connecting()) {
          return;
        }

        setState(radioUUID() == "" ? DefaultState : { uuid: "" });

        sendStateCommand(ws, radioUUID());
      },
      { defer: true }
    )
  );

  createEffect(
    on(
      shouldReconnect,
      () => {
        if (shouldReconnect()) reconnect();
      },
      { defer: true }
    )
  );

  return [
    {
      state,
      stale,
      discovering,
    },
    {
      connecting,
      connected,
      disconnected,
      reconnect,
    },
  ];
}
