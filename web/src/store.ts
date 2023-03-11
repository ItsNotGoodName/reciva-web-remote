import { createEffect, createResource, type Accessor } from "solid-js";
import {
  Api,
  type ModelPreset,
  type HttpPatchState,
  type StateState,
} from "./api";
import { type ModelRadio } from "./api";
import { API_URL } from "./constants";
import { createMutation, once, checkStale, createStaleSignal } from "./utils";
import { type WSDataReturn } from "./ws";

const api = new Api({ baseUrl: API_URL });

// Build get
export const buildGetQuery = once(() =>
  createResource(() => api.build.buildList())
);

// Radios list
export const radiosListQuery = createResource<ModelRadio[], string>(() =>
  api.radios.radiosList()
);

// Presets list
export const presetListQuery = once(() =>
  createResource(() => api.presets.presetsList())
);

// Preset get
const [stalePresets, setStalePresets] = createStaleSignal(undefined);
export const usePresetQuery = (url: Accessor<string | undefined>) =>
  checkStale(
    createResource<ModelPreset, string>(
      () => url() || undefined,
      (url: string) => api.presets.presetsDetail(url)
    ),
    () => !!stalePresets()
  );

// State get
const [staleStateUUID, setStaleStateUUID] = createStaleSignal("");
export const useStateQuery = (uuid: Accessor<string | undefined>) =>
  checkStale(
    createResource<StateState, string>(
      () => uuid() || undefined,
      (uuid: string) => api.states.statesDetail(uuid)
    ),
    () => staleStateUUID() == uuid(),
    staleStateUUID
  );

// Radios discover
export const discoverRadios = createMutation(
  () => api.radios.radiosCreate(),
  [radiosListQuery]
);

// Radio volume refresh
export const refreshRadioVolume = createMutation((uuid: string) =>
  api.radios.volumeCreate(uuid).then(() => setStaleStateUUID(uuid))
);

// Radio subscription refresh
export const refreshRadioSubscription = createMutation((uuid: string) =>
  api.radios.subscriptionCreate(uuid)
);

// State update
export const patchState = createMutation(
  (req: HttpPatchState & { uuid: string }) =>
    api.states
      .statesPartialUpdate(req.uuid, req)
      .then(() => setStaleStateUUID(req.uuid))
);

// Preset update
export const updatePreset = createMutation((preset: ModelPreset) =>
  api.presets.presetsCreate(preset).then(() => setStalePresets())
);

// Websocket data hook
export const hookWSData = (data: WSDataReturn) => {
  createEffect(() => {
    if (!data.discovering()) {
      void radiosListQuery[1].refetch();
    }
  });
};
