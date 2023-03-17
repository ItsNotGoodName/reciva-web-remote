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
import { API_URL } from "./constants";
import { invalidateWhen, createInvalidateSignal } from "./utils";

const api = new Api({ baseUrl: API_URL });

////////////// Queries

// Build get
export const useBuildGetQuery = () =>
  createResource(() => api.build.buildList());

// Radios list
export const [invalidRadioListQuery, invalidateRadioListQuery] =
  createInvalidateSignal(new Date());
export const useRadiosListQuery = () =>
  createResource(invalidRadioListQuery, () => api.radios.radiosList());

// Presets list
export const [invalidPresetListQuery, invalidatePresetListQuery] =
  createInvalidateSignal(new Date());
export const usePresetListQuery = () =>
  createResource(invalidPresetListQuery, () => api.presets.presetsList());

// Preset get
export const [invalidPresetQuery, invalidatePresetQuery] =
  createInvalidateSignal<string | undefined>(undefined);
export const usePresetQuery = (url: Accessor<string | undefined>) =>
  invalidateWhen(
    () => invalidPresetQuery() === undefined || invalidPresetQuery() == url(),
    createResource<ModelPreset, string>(url, (url: string) =>
      api.presets.presetsDetail(url)
    )
  );

////////////// Mutation

export type MutationReturn<T = void, R = unknown, E = unknown> = {
  mutate: (data: T) => R | Promise<R>;
  loading: Accessor<boolean>;
  error: Accessor<E>;
  cancel: () => void;
};

let lastMutationToken = 0;

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

    token = lastMutationToken++;
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

const cancelMutateWhen = <T, R, U>(
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

// Radios discover
export const useDiscoverRadios = () =>
  createMutation((params) => api.radios.radiosCreate(params));

export const useDeleteRadio = (uuid: Accessor<string>) =>
  cancelMutateWhen(
    uuid,
    createMutation((params) => api.radios.radiosDelete(uuid(), params))
  );

// Radio volume refresh
export const useRefreshRadioVolume = (uuid: Accessor<string>) =>
  cancelMutateWhen(
    uuid,
    createMutation((params) => api.radios.volumeCreate(uuid(), params))
  );

// Radio subscription refresh
export const useRefreshRadioSubscription = (uuid: Accessor<string>) =>
  cancelMutateWhen(
    uuid,
    createMutation((params) => api.radios.subscriptionCreate(uuid(), params))
  );

// State update
export const useUpdateState = (uuid: Accessor<string>) =>
  cancelMutateWhen(
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
      .then(() => invalidatePresetQuery(preset.url))
  );
