import {
  Switch,
  Match,
  For,
  createSignal,
  type Setter,
  createEffect,
  on,
  Show,
  type ParentComponent,
  Index,
  type JSX,
  splitProps,
  batch,
  type Accessor,
} from "solid-js";
import { type Component } from "solid-js";
import { type StatePreset, StateStatus, type StateState } from "./api";
import {
  hookWSData,
  useDiscoverRadios,
  radiosListQuery,
  usePatchState,
} from "./store";
import { clickOutside, type ClassProps, mergeClass } from "./utils";
import { useWS } from "./ws";

// Prevent TypeScript from removing clickOutside directive
false && clickOutside;

const DiscoverButton: Component<{ discovering: boolean } & ClassProps> = (
  props
) => {
  console.log("Render: DiscoverButton");
  const discoverRadios = useDiscoverRadios();
  const loading = (): boolean => discoverRadios.loading() || props.discovering;

  const onClick = () => {
    void discoverRadios.mutate(null);
  };

  return (
    <button
      class={mergeClass("btn-primary btn", props.class)}
      classList={{ loading: loading() }}
      onClick={onClick}
      disabled={loading()}
    >
      Discover
    </button>
  );
};

const RadioPlayerStatusButton: Component<
  {
    status: StateStatus;
    loading?: boolean;
  } & ClassProps
> = (props) => {
  console.log("Render: RadioPlayerStatus");
  const data = (): { class: string; status: string } => {
    switch (props.status) {
      case StateStatus.StatusConnecting:
        return { class: "btn btn-circle btn-warning", status: props.status };
      case StateStatus.StatusPlaying:
        return { class: "btn btn-circle btn-success", status: props.status };
      case StateStatus.StatusStopped:
        return { class: "btn btn-circle btn-error", status: props.status };
      default:
        return {
          class: "btn btn-circle no-animation btn-info",
          status: "Unknown",
        };
    }
  };

  return (
    <DaisyTooltip class={props.class} tooltip={data().status}>
      <button
        class={data().class}
        classList={{ loading: props.loading }}
        aria-label={data().status}
      >
        {!props.loading && data().status}
      </button>
    </DaisyTooltip>
  );
};

const RadioPlayerTitleDropdown: Component<
  {
    state: StateState;
    loading?: boolean;
    classButton?: string;
    classDropdown?: string;
  } & ClassProps
> = (props) => {
  console.log("Render: RadioPlayerTitleDropdown");

  const data = (): DaisyStaticTableCardBodyData[] => [
    {
      key: "Metadata",
      value: props.state.metadata,
    },
    {
      key: "Title",
      value: props.state.title,
    },
    {
      key: "New Title",
      value: props.state.title_new,
    },
    {
      key: "Preset Number",
      value: props.state.preset_number,
    },
    {
      key: "URL",
      value: (
        <a class="link-hover link" href={props.state.url}>
          {props.state.url}
        </a>
      ),
    },
    {
      key: "New URL",
      value: (
        <a class="link-hover link" href={props.state.url_new}>
          {props.state.url_new}
        </a>
      ),
    },
  ];

  return (
    <div
      class={mergeClass("dropdown no-animation", props.class)}
      use:clickOutside=""
    >
      <label
        tabindex="0"
        class={mergeClass(
          "btn-primary btn justify-start gap-2 truncate",
          props.classButton
        )}
        classList={{ loading: props.loading }}
      >
        {!props.loading && (
          <>
            <span class="badge-info badge badge-lg rounded-md">
              {props.state.preset_number}
            </span>
            {props.state.title_new || props.state.title}
          </>
        )}
      </label>
      <div
        tabindex="0"
        class={mergeClass(
          "card-compact card dropdown-content w-full bg-primary p-2 text-primary-content shadow",
          props.classDropdown
        )}
      >
        <DaisyStaticTableCardBody data={data()} title="Stream Information" />
      </div>
    </div>
  );
};

const RadioTypeDropdown: Component<
  {
    state: StateState;
    classButton?: string;
    classDropdown?: string;
  } & ClassProps
> = (props) => {
  console.log("Render: RadioTypeDropdown");
  const data = (): DaisyStaticTableCardBodyData[] => [
    { key: "Name", value: props.state.name },
    { key: "Model Name", value: props.state.model_name },
    { key: "Model Number", value: props.state.model_number },
  ];

  return (
    <div class={mergeClass("dropdown", props.class)} use:clickOutside="">
      <label
        tabindex="0"
        class={mergeClass("btn-primary btn", props.classButton)}
      >
        Radio
      </label>
      <div
        tabindex="0"
        class={mergeClass(
          "card-compact card dropdown-content w-80 bg-primary p-2 text-primary-content shadow",
          props.classDropdown
        )}
      >
        <DaisyStaticTableCardBody data={data()} title="Radio Information" />
      </div>
    </div>
  );
};

const RadioPresetButton: Component<
  {
    selected: boolean;
    preset: StatePreset;
    loading?: boolean;
  } & JSX.HTMLAttributes<HTMLButtonElement>
> = (props) => {
  console.log("Render: RadioPresetButton");
  const [, other] = splitProps(props, [
    "selected",
    "preset",
    "loading",
    "class",
  ]);

  return (
    <button
      class={mergeClass("btn flex gap-2", props.class)}
      classList={{ "btn-primary": props.selected, loading: props.loading }}
      {...other}
    >
      <Show when={!props.loading}>
        <span class="badge-info badge badge-lg rounded-md">
          {props.preset.number}
        </span>
        <span class="w-0 flex-grow truncate">
          {props.preset.title_new ? props.preset.title_new : props.preset.title}{" "}
        </span>
      </Show>
    </button>
  );
};

const RadioPresetsList: Component<
  {
    radioUUID: Accessor<string>;
    state: StateState;
  } & ClassProps
> = (props) => {
  console.log("Render: RadioPresetList");
  const statePatch = usePatchState(props.radioUUID);
  const [lastLoadingNumber, setLoadingNumber] = createSignal(-1);
  const loadingNumber = (): number =>
    statePatch.loading() ? lastLoadingNumber() : -1;

  const onClick = (preset: number) => {
    batch(() => {
      statePatch.cancel();
      void statePatch.mutate({ preset: preset });
      setLoadingNumber(preset);
    });
  };

  return (
    <div
      class={mergeClass(
        "grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3",
        props.class
      )}
    >
      <For each={props.state.presets}>
        {(preset) => (
          <RadioPresetButton
            selected={props.state.preset_number == preset.number}
            loading={loadingNumber() == preset.number}
            onclick={[onClick, preset.number]}
            preset={preset}
          />
        )}
      </For>
    </div>
  );
};

const RadioPowerButton: Component<
  {
    radioUUID: Accessor<string>;
    state: StateState;
  } & ClassProps
> = (props) => {
  console.log("Render: RadioPowerButton");
  const patchState = usePatchState(props.radioUUID);

  const toggle = () => {
    void patchState.mutate({
      power: !props.state.power,
    });
  };

  return (
    <button
      class={mergeClass("btn", props.class)}
      classList={{
        "btn-success": props.state.power,
        "btn-error": !props.state.power,
        loading: patchState.loading(),
      }}
      disabled={patchState.loading()}
      onClick={toggle}
    >
      <Switch>
        <Match when={props.state.power}>ON</Match>
        <Match when={!props.state.power}>OFF</Match>
      </Switch>
    </button>
  );
};

const RadioSelect: Component<
  {
    radioUUID: string;
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
        select && (select.value = props.radioUUID);
      },
      { defer: true }
    )
  );

  return (
    <select
      class={mergeClass("select-primary select", props.class)}
      ref={select}
      disabled={radios.loading}
      value={props.radioUUID}
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

const DaisyTooltip: ParentComponent<{ tooltip: string } & ClassProps> = (
  props
) => {
  return (
    <div class={mergeClass("tooltip", props.class)} data-tip={props.tooltip}>
      {props.children}
    </div>
  );
};

type DaisyStaticTableCardBodyData = { key: string; value: JSX.Element };

const DaisyStaticTableCardBody: Component<{
  title: string;
  data: DaisyStaticTableCardBodyData[];
}> = (props) => {
  return (
    <div class="card-body">
      <h3 class="card-title">{props.title}</h3>
      <table class="table-fixed">
        <tbody>
          <Index each={props.data}>
            {(data) => (
              <tr>
                <td>
                  <span class="text badge-info badge mr-2 w-full whitespace-nowrap">
                    {data().key}
                  </span>
                </td>
                <td class="w-full break-all">{data().value}</td>
              </tr>
            )}
          </Index>
        </tbody>
      </table>
    </div>
  );
};

const TopBar: ParentComponent<ClassProps> = (props) => {
  return (
    <div
      class={mergeClass(
        "fixed top-0 z-50 w-full border-b-2 border-accent border-b-base-300 bg-base-200 p-2",
        props.class
      )}
    >
      {props.children}
    </div>
  );
};

const BottomBar: ParentComponent<ClassProps> = (props) => {
  return (
    <div
      class={mergeClass(
        "fixed bottom-0 z-50 w-full border-t-2 border-accent border-t-base-300 bg-base-200 p-2",
        props.class
      )}
    >
      {props.children}
    </div>
  );
};

const App: Component = () => {
  const [radioUUID, setRadioUUID] = createSignal<string>("");
  const radioSelected = () => radioUUID() != "";

  const [data, ws] = useWS(radioUUID);
  hookWSData(data);
  const { state, discovering } = data;
  const loading = () => state.uuid != radioUUID() || ws.connecting();

  return (
    <div class="h-screen">
      <TopBar class="flex gap-2">
        <RadioPlayerStatusButton
          class="tooltip-bottom flex"
          status={state.status}
          loading={loading()}
        />
        <RadioPlayerTitleDropdown
          class="dropdown-end flex-1"
          classButton="w-full"
          classDropdown="mt-2"
          state={state}
          loading={loading()}
        />
      </TopBar>
      <div class="container mx-auto py-20 px-4">
        <RadioPresetsList radioUUID={radioUUID} state={state} />
      </div>
      <BottomBar class="flex gap-2">
        <div class="flex flex-1">
          <DaisyTooltip tooltip="Discover">
            <DiscoverButton
              class="rounded-r-none"
              discovering={discovering()}
            />
          </DaisyTooltip>
          <RadioSelect
            class="w-full flex-1 rounded-l-none"
            radioUUID={radioUUID()}
            setRadioUUID={setRadioUUID}
          />
        </div>
        <Show when={radioSelected()}>
          <RadioPowerButton radioUUID={radioUUID} state={state} />
          <RadioTypeDropdown
            class="dropdown-top dropdown-end"
            classDropdown="mb-2"
            state={state}
          />
        </Show>
      </BottomBar>
    </div>
  );
};

export default App;
