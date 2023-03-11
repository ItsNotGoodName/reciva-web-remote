import {
  Switch,
  Match,
  For,
  createSignal,
  type Setter,
  type Accessor,
  createEffect,
  on,
  Show,
  type ParentComponent,
} from "solid-js";
import { type Component } from "solid-js";
import { type StateState } from "./api";
import {
  hookWSData,
  discoverRadios,
  radiosListQuery,
  patchState,
} from "./store";
import { type ClassProps, mergeClass } from "./utils";
import { useWS } from "./ws";

const DiscoverButton: Component<
  { discovering: Accessor<boolean> } & ClassProps
> = (props) => {
  console.log("Render: DiscoverButton");

  const loading = () => discoverRadios.loading() || props.discovering();

  return (
    <button
      class={mergeClass(
        "rounded bg-blue-500 py-2 px-4 font-bold text-white hover:bg-blue-700",
        props.class
      )}
      onClick={() => void discoverRadios.mutate(null)}
      disabled={loading()}
    >
      <Switch fallback={<>Discover</>}>
        <Match when={loading()}>
          <>Discovering...</>
        </Match>
      </Switch>
    </button>
  );
};

const RadioPowerButton: Component<
  {
    radioUUID: Accessor<string>;
    state: StateState;
  } & ClassProps
> = (props) => {
  console.log("Render: RadioPowerButton");

  function toggle() {
    void patchState.mutate({
      uuid: props.radioUUID(),
      power: !props.state.power,
    });
  }

  return (
    <button
      class={mergeClass("rounded py-2 px-4 font-bold text-white", props.class)}
      classList={{
        "bg-green-500 hover:bg-green-700": props.state.power,
        "bg-red-500 hover:bg-red-700": !props.state.power,
      }}
      disabled={patchState.loading()}
      onClick={toggle}
    >
      <Switch>
        <Match when={patchState.loading()}>LOADING...</Match>
        <Match when={props.state.power}>ON</Match>
        <Match when={!props.state.power}>OFF</Match>
      </Switch>
    </button>
  );
};

const RadioSelect: Component<
  {
    radioUUID: Accessor<string>;
    setRadioUUID: Setter<string>;
  } & ClassProps
> = (props) => {
  console.log("Render: RadioSelect");

  const [radios] = radiosListQuery;
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
      class={mergeClass("rounded px-4 py-2 hover:bg-gray-300", props.class)}
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

const TopBar: ParentComponent = (props) => {
  return (
    <div class="fixed top-0 z-50 w-full border-b-2 p-2">{props.children}</div>
  );
};

const BottomBar: ParentComponent = (props) => {
  return (
    <div class="fixed bottom-0 z-50 w-full border-t-2 p-2">
      {props.children}
    </div>
  );
};

const App: Component = () => {
  const [radioUUID, setRadioUUID] = createSignal<string>("");
  const showRadioControls = () => radioUUID() != "";

  const [data] = useWS(radioUUID);
  hookWSData(data);
  const { state, discovering } = data;

  return (
    <div class="h-screen">
      <TopBar />
      <div class="py-20"></div>
      <BottomBar>
        <div class="flex gap-2">
          <DiscoverButton discovering={discovering} />
          <RadioSelect
            class="flex-1"
            radioUUID={radioUUID}
            setRadioUUID={setRadioUUID}
          />
          <Show when={showRadioControls()}>
            <RadioPowerButton radioUUID={radioUUID} state={state} />
          </Show>
        </div>
      </BottomBar>
    </div>
  );
};

export default App;
