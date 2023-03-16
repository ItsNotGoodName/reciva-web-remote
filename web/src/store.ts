import {
  createEffect,
  createResource,
  createSignal,
  type ResourceReturn,
  type Accessor,
  on,
  type AccessorArray,
  batch,
} from "solid-js";
import {
  Api,
  type HttpPostState,
  type ModelPreset,
  type RequestParams,
} from "./api";
import { type ModelRadio } from "./api";
import { API_URL } from "./constants";
import { staleWhen, createStaleSignal, invalidWhen } from "./utils";

const api = new Api({ baseUrl: API_URL });

export type MutationReturn<T = void, R = unknown, E = unknown> = {
  mutate: (data: T) => R | Promise<R>;
  loading: Accessor<boolean>;
  error: Accessor<E>;
  cancel: () => void;
};

let lastToken = 0;

function createMutation<T, R>(
  mutateFn: (params: RequestParams, data: T) => Promise<R>,
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  invalidateQueries: ResourceReturn<any, any>[] = []
): MutationReturn<T, R> {
  const [loading, setloading] = createSignal(false);
  const [error, setError] = createSignal();
  let token = -1;

  function mutate(data: T): Promise<R> | R {
    batch(() => {
      setloading(true);
      setError();
    });

    token = lastToken++;
    const thisToken = token;

    const pr = mutateFn({ cancelToken: thisToken }, data);

    void pr
      .then(() =>
        Promise.allSettled(
          invalidateQueries.map((query) => void query[1].refetch())
        )
      )
      .catch((e) => {
        if (token == thisToken) {
          batch(() => {
            setloading(false);
            setError(e);
          });
        }
      })
      .then(() => {
        if (token == thisToken) setloading(false);
      });

    return pr;
  }

  function cancel() {
    api.abortRequest(token);
    setloading(false);
  }

  return { loading, mutate, cancel, error };
}

export const cancelWhen = <T, R, U>(
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

////////////// Queries

// Build get
export const useBuildGetQuery = () =>
  createResource(() => api.build.buildList());

// Radios list
export const useRadiosListQuery = () =>
  createResource<ModelRadio[], string>(() => api.radios.radiosList());

// Presets list
export const [invalidPresetListQuery, invalidatePresetListQuery] =
  createStaleSignal(undefined);
export const usePresetListQuery = () =>
  invalidWhen(
    invalidPresetListQuery,
    createResource(() => api.presets.presetsList())
  );

// Preset get
export const [stalePresetQuery, setStalePresetQuery] = createStaleSignal<
  string | undefined
>(undefined);
export const usePresetQuery = (url: Accessor<string | undefined>) =>
  staleWhen(
    createResource<ModelPreset, string>(
      () => url() || undefined,
      (url: string) => api.presets.presetsDetail(url)
    ),
    () => stalePresetQuery() === undefined || stalePresetQuery() == url()
  );

// // State get
// const [staleStateUUID, setStaleStateUUID] = createStaleSignal("");
// export const useStateQuery = (uuid: Accessor<string | undefined>) =>
//   staleWhen(
//     createResource<StateState, string>(
//       () => uuid() || undefined,
//       (uuid: string) => api.states.statesDetail(uuid)
//     ),
//     () => staleStateUUID() == uuid(),
//     staleStateUUID
//   );

////////////// Mutation

// Radios discover
export const useDiscoverRadios = () =>
  createMutation((params) => api.radios.radiosCreate(params));

export const useDeleteRadio = (uuid: Accessor<string>) =>
  cancelWhen(
    uuid,
    createMutation((params) => api.radios.radiosDelete(uuid(), params))
  );

// Radio volume refresh
export const useRefreshRadioVolume = (uuid: Accessor<string>) =>
  cancelWhen(
    uuid,
    createMutation((params) => api.radios.volumeCreate(uuid(), params))
  );

// Radio subscription refresh
export const useRefreshRadioSubscription = (uuid: Accessor<string>) =>
  cancelWhen(
    uuid,
    createMutation((params) => api.radios.subscriptionCreate(uuid(), params))
  );

// State update
export const useUpdateState = (uuid: Accessor<string>) =>
  cancelWhen(
    uuid,
    createMutation((params, req: HttpPostState) =>
      api.states.statesCreate(uuid(), req, params)
    )
  );

// Preset update
export const useUpdatePreset = () =>
  createMutation((params, preset: ModelPreset) =>
    api.presets
      .presetsCreate(preset, params)
      .then(() => setStalePresetQuery(preset.url))
  );
