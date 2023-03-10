import {
  createResource,
  createSignal,
  type Accessor,
  createEffect,
  on,
} from "solid-js";
import { Api, type HttpPatchState, type StateState } from "./api";
import { type ModelRadio } from "./api";
import { API_URL } from "./constant";
import { createMutation } from "./utils";

const api = new Api({ baseUrl: API_URL });

// Queries

export const radiosQuery = createResource<ModelRadio[], string>(() =>
  api.radios.radiosList().then((res) => res.data)
);
const [stateState, setStaleState] = createSignal<string>("", { equals: false });
export const stateQuery = (uuid: Accessor<string | undefined>) => {
  const query = createResource<StateState, string>(
    () => uuid() || undefined,
    (uuid: string) => api.states.statesDetail(uuid).then((res) => res.data)
  );

  createEffect(
    on(
      stateState,
      () => {
        if (stateState() == uuid()) return query[1].refetch();
      },
      { defer: true }
    )
  );

  return query;
};

// Mutations

export const discover = createMutation(
  () => api.radios.radiosCreate(),
  [radiosQuery]
);
export const statePatch = createMutation(
  (req: HttpPatchState & { uuid: string }) =>
    api.states
      .statesPartialUpdate(req.uuid, req)
      .then(() => setStaleState(req.uuid))
);
