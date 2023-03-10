import {
  Switch,
  Match,
  For,
  Show,
  createSignal,
  Setter,
  Accessor,
  createEffect,
  createResource,
  ErrorBoundary,
  resetErrorBoundaries,
  on,
} from "solid-js";
import type { Component } from "solid-js";
import { Api, HttpPatchState, ModelRadio, StateState } from "./api";
import { API_URL } from "./constant";
import { createMutation } from "./util";

const api = new Api({ baseUrl: API_URL + "/api" }); // TODO: get api path from swagger.json

// Queries

const radiosQuery = createResource<ModelRadio[], string>(() =>
  api.radios.radiosList().then((res) => res.data)
);
const [stateState, setStaleState] = createSignal<string>("", { equals: false });
const stateQuery = (uuid: Accessor<string | undefined>) => {
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

const discoverMutation = createMutation(
  () => api.radios.radiosCreate(),
  [radiosQuery]
);
const stateMutation = createMutation<HttpPatchState & { uuid: string }>((req) =>
  api.states
    .statesPartialUpdate(req.uuid, req)
    .then(() => setStaleState(req.uuid))
);

// Components

const DiscoverButton: Component = () => {
  console.log("Render: DiscoverButton");

  return (
    <button onClick={() => discoverMutation[1].refetch({})}>
      <Switch fallback={<>Discover</>}>
        <Match when={discoverMutation[0].loading}>
          <>Discovering...</>
        </Match>
      </Switch>
    </button>
  );
};

const RadioPowerButton: Component<{
  radioUUID: Accessor<string>;
  state: Accessor<StateState>;
}> = (props) => {
  console.log("Render: RadioPowerButton");

  function toggle() {
    stateMutation[1].refetch({
      uuid: props.radioUUID(),
      power: !props.state().power,
    });
  }

  return (
    <button onClick={toggle}>
      <Switch>
        <Match when={props.state().power}>ON</Match>
        <Match when={!props.state().power}>OFF</Match>
      </Switch>
    </button>
  );
};

const RadioSelect: Component<{
  radioUUID: Accessor<string>;
  setRadioUUID: Setter<string>;
}> = (props) => {
  console.log("Render: RadioSelect");

  const [radios] = radiosQuery;
  let select: HTMLSelectElement | undefined;

  // Prevent select.value from defaulting to first option when radios.data changes
  createEffect(
    on(radios, () => {
      select!.value = props.radioUUID();
    })
  );

  return (
    <select
      ref={select}
      disabled={radios.loading}
      value={props.radioUUID()}
      onChange={(e) => {
        props.setRadioUUID(e.currentTarget.value);
      }}
    >
      <option disabled value="">
        <Switch fallback={<>Select Radio</>}>
          <Match when={radios.loading}>Loading...</Match>
          <Match when={radios.error}>Network Error</Match>
          <Match when={radios() !== undefined && radios()!.length == 0}>
            No Radios Found
          </Match>
        </Switch>
      </option>
      <For each={radios()}>
        {(radio) => <option value={radio.uuid}>{radio.name}</option>}
      </For>
    </select>
  );
};

const App: Component = () => {
  const [radioUUID, setRadioUUID] = createSignal<string>("");
  const [state] = stateQuery(radioUUID);

  createEffect(on(radioUUID, () => resetErrorBoundaries()));

  return (
    <>
      <DiscoverButton />
      <RadioSelect radioUUID={radioUUID} setRadioUUID={setRadioUUID} />
      <ErrorBoundary fallback={<>Could not fetch state</>}>
        <Show when={state() !== undefined}>
          <RadioPowerButton
            radioUUID={radioUUID}
            state={state as Accessor<StateState>}
          />
        </Show>
      </ErrorBoundary>
    </>
  );
};

export default App;
