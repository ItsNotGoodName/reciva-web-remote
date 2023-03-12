import {
  createEffect,
  createResource,
  createSignal,
  type ResourceReturn,
  type Accessor,
  on,
  type AccessorArray,
} from "solid-js";
import {
  Api,
  type ModelPreset,
  type HttpPatchState,
  type StateState,
  type RequestParams,
} from "./api";
import { type ModelRadio } from "./api";
import { API_URL } from "./constants";
import { once, checkStale, createStaleSignal } from "./utils";
import { type WSDataReturn } from "./ws";

const api = new Api({ baseUrl: API_URL });

export type MutationReturn<T = void, R = unknown> = {
  mutate: (data: T) => R | Promise<R>;
  loading: Accessor<boolean>;
  cancel: () => void;
};

let lastToken = 0;

function createMutation<T, R>(
  mutateFn: (params: RequestParams, data: T) => Promise<R>,
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  invalidateQueries: ResourceReturn<any, any>[] = []
): MutationReturn<T, R> {
  const [loading, setloading] = createSignal(false);
  let token = -1;

  function mutate(data: T): Promise<R> | R {
    setloading(true);

    token = lastToken++;
    const thisToken = token;

    const pr = mutateFn({ cancelToken: thisToken }, data);

    void pr
      .then(() =>
        Promise.allSettled(
          invalidateQueries.map((query) => void query[1].refetch())
        )
      )
      .finally(() => {
        if (token == thisToken) setloading(false);
      });

    return pr;
  }

  function cancel() {
    api.abortRequest(token);
    setloading(false);
  }

  return { loading, mutate, cancel };
}

export const cancelOn = <T, R, U>(
  deps: Accessor<U> | AccessorArray<U>,
  mutation: MutationReturn<T, R>
): MutationReturn<T, R> => {
  createEffect(
    on(
      deps,
      () => {
        mutation.cancel();
      },
      { defer: true }
    )
  );
  return mutation;
};

// Websocket data hook
export const hookWSData = (data: WSDataReturn) => {
  createEffect(() => {
    if (!data.discovering()) {
      void radiosListQuery[1].refetch();
    }
  });
};

////////////// Queries

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

////////////// Mutation

// Radios discover
export const useDiscoverRadios = () =>
  createMutation(
    (params) => api.radios.radiosCreate(params),
    [radiosListQuery]
  );

// Radio volume refresh
export const useRefreshRadioVolume = (uuid: Accessor<string>) =>
  cancelOn(
    uuid,
    createMutation((params) =>
      api.radios
        .volumeCreate(uuid(), params)
        .then(() => setStaleStateUUID(uuid))
    )
  );

// Radio subscription refresh
export const useRefreshRadioSubscription = (uuid: Accessor<string>) =>
  cancelOn(
    uuid,
    createMutation((params) => api.radios.subscriptionCreate(uuid(), params))
  );

// State update
export const usePatchState = (uuid: Accessor<string>) =>
  cancelOn(
    uuid,
    createMutation((params, req: HttpPatchState) =>
      api.states
        .statesPartialUpdate(uuid(), req, params)
        .then(() => setStaleStateUUID(uuid))
    )
  );

// Preset update
export const useUpdatePreset = () =>
  createMutation((params, preset: ModelPreset) =>
    api.presets.presetsCreate(preset, params).then(() => setStalePresets())
  );
