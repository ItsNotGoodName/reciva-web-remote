import {
  Switch,
  Match,
  For,
  Show,
  createSignal,
  type Setter,
  type Accessor,
  createEffect,
  ErrorBoundary,
  resetErrorBoundaries,
  on,
} from "solid-js";
import { type Component } from "solid-js";
import { type StateState } from "./api";
import { radiosDiscover, radiosList, statePatch, stateResource } from "./store";

const DiscoverButton: Component = () => {
  console.log("Render: DiscoverButton");

  return (
    <button onClick={() => void radiosDiscover.mutate(null)}>
      <Switch fallback={<>Discover</>}>
        <Match when={radiosDiscover.loading()}>
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
    void statePatch.mutate({
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

  const [radios] = radiosList;
  let select: HTMLSelectElement | undefined;

  // Prevent select.value from defaulting to the first option when radios.data changes
  createEffect(
    on(
      radios,
      () => {
        select && (select.value = props.radioUUID());
      },
      { defer: true }
    )
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
          <Match when={radios.error !== undefined}>Network Error</Match>
          <Match when={radios()?.length == 0}>No Radios Found</Match>
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
  const [state] = stateResource(radioUUID);

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
