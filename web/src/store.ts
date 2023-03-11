import { createEffect, createResource, type Accessor } from "solid-js";
import {
  Api,
  type ModelPreset,
  type HttpPatchState,
  type StateState,
} from "./api";
import { type ModelRadio } from "./api";
import { API_URL } from "./constant";
import { createMutation, once, checkStale, createStaleSignal } from "./utils";
import { type WSDataReturn } from "./ws";

const api = new Api({ baseUrl: API_URL });

// Build get
export const buildGet = once(() => createResource(() => api.build.buildList()));

// Radios list
export const radiosList = createResource<ModelRadio[], string>(() =>
  api.radios.radiosList()
);

// Presets list
export const presetsList = once(() =>
  createResource(() => api.presets.presetsList())
);

// Preset get
const [stalePresets, setStalePresets] = createStaleSignal(undefined);
export const presetResource = (url: Accessor<string | undefined>) =>
  checkStale(
    createResource<ModelPreset, string>(
      () => url() || undefined,
      (url: string) => api.presets.presetsDetail(url)
    ),
    () => !!stalePresets()
  );

// State get
const [staleStateUUID, setStaleStateUUID] = createStaleSignal("");
export const stateResource = (uuid: Accessor<string | undefined>) =>
  checkStale(
    createResource<StateState, string>(
      () => uuid() || undefined,
      (uuid: string) => api.states.statesDetail(uuid)
    ),
    () => staleStateUUID() == uuid(),
    staleStateUUID
  );

// Radios discover
export const radiosDiscover = createMutation(
  () => api.radios.radiosCreate(),
  [radiosList]
);

// Radio volume refresh
export const radioVolumeRefresh = createMutation((uuid: string) =>
  api.radios.volumeCreate(uuid).then(() => setStaleStateUUID(uuid))
);

// Radio subscription refresh
export const radioSubscriptionRefresh = createMutation((uuid: string) =>
  api.radios.subscriptionCreate(uuid)
);

// State update
export const statePatch = createMutation(
  (req: HttpPatchState & { uuid: string }) =>
    api.states
      .statesPartialUpdate(req.uuid, req)
      .then(() => setStaleStateUUID(req.uuid))
);

// Preset update
export const presetUpdate = createMutation((preset: ModelPreset) =>
  api.presets.presetsCreate(preset).then(() => setStalePresets())
);

// Websocket
export const websocketBind = (data: WSDataReturn) => {
  createEffect(() => {
    if (!data.discovering()) {
      void radiosList[1].refetch();
    }
  });
};
