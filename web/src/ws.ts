import { type Accessor, createSignal, createEffect, on } from "solid-js";
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
  //   stateLoading: Accessor<boolean>;
  //   stateSelected: Accessor<boolean>;
};

export type WSStatusReturn = {
  connecting: Accessor<boolean>;
  connected: Accessor<boolean>;
  disconnected: Accessor<boolean>;
  reconnect: () => void;
};

export function useWS(
  stateUUID: Accessor<string>
): [WSDataReturn, WSStatusReturn] {
  const [connecting, setConnecting] = createSignal(true);
  const [connected, setConnected] = createSignal(false);
  const [disconnected, setDisconnected] = createSignal(false);
  const [discovering, setDiscovering] = createSignal(false);
  const [state, setState] = createStore(defaultState);

  //   const stateSelected = () => state.uuid != "";
  //   const stateLoading = () => state.uuid != stateUUID() || connecting();

  const connect = () => {
    const ws = new WebSocket(WS_URL);
    setConnecting(true);
    console.log("WS: Connecting");

    ws.addEventListener("open", () => {
      console.log("WS: Connected");
      setConnecting(false);
      setConnected(true);
      setDisconnected(false);
      subscribe(ws, stateUUID());
    });

    ws.addEventListener("message", (event) => {
      console.log("WS: Message");
      const msg = JSON.parse(event.data as string) as WsEvent;
      if (msg.topic == PubsubTopic.StateTopic) {
        const data = msg.data as StateState;
        if (data.uuid == stateUUID()) {
          setState(
            produce((state) => Object.assign(state, msg.data as StateState))
          );
        }
      } else if (msg.topic == PubsubTopic.DiscoverTopic) {
        setDiscovering(msg.data as boolean);
      }
    });

    ws.addEventListener("close", () => {
      console.log("WS: Close");
      setConnecting(false);
      setConnected(false);
      setDisconnected(true);
      setState(defaultState);
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
    on(stateUUID, () => {
      if (!connected()) {
        reconnect();
        return;
      }
      if (connecting()) {
        return;
      }

      setState(defaultState);
      subscribe(ws, stateUUID());
    })
  );

  return [
    {
      state,
      discovering,
      //   stateLoading,
      //   stateSelected,
    },
    {
      connecting,
      connected,
      disconnected,
      reconnect,
    },
  ];
}
