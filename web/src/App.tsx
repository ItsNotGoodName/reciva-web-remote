import {
  FaBrandsItunesNote,
  FaSolidPlay,
  FaSolidPowerOff,
  FaSolidQuestion,
  FaSolidRadio,
  FaSolidStop,
  FaSolidVolumeHigh,
  FaSolidVolumeLow,
  FaSolidVolumeOff,
} from "solid-icons/fa";
import { IoSearchSharp } from "solid-icons/io";
import { FiRefreshCw } from "solid-icons/fi";
import {
  Switch,
  Match,
  For,
  createSignal,
  type Setter,
  createEffect,
  on,
  Show,
  type JSX,
  splitProps,
  batch,
  type Accessor,
  type Resource,
  createReaction,
  onCleanup,
} from "solid-js";
import { type Component } from "solid-js";
import {
  type StatePreset,
  StateStatus,
  type StateState,
  type ModelRadio,
} from "./api";
import {
  useDiscoverRadios,
  usePatchState,
  useRefreshRadioVolume,
  useRefreshRadioSubscription,
  useRadiosListQuery,
} from "./store";
import { type ClassProps, mergeClass, IOS } from "./utils";
import { useWS } from "./ws";
import {
  DaisyButton,
  DaisyStaticTableCardBody,
  type DaisyStaticTableCardBodyData,
  DaisyTooltip,
  DaisyDropdown,
} from "./Daisy";

const DiscoverButton: Component<
  { discovering: boolean; classButton?: string } & ClassProps
> = (props) => {
  const discoverRadios = useDiscoverRadios();
  const loading = (): boolean => discoverRadios.loading() || props.discovering;

  const discover = () => {
    void discoverRadios.mutate(null);
  };

  return (
    <DaisyTooltip class={props.class} tooltip="Discover">
      <DaisyButton
        class={mergeClass("btn-primary w-14", props.classButton)}
        loading={loading()}
        onClick={discover}
        aria-label="Discover"
      >
        <IoSearchSharp size={20} />
      </DaisyButton>
    </DaisyTooltip>
  );
};

const RadioRefreshSubscriptionButton: Component<
  { radioUUID: Accessor<string>; classButton?: string } & ClassProps
> = (props) => {
  const refreshRadioSubscription = useRefreshRadioSubscription(props.radioUUID);

  const refreshSubscription = () => {
    void refreshRadioSubscription.mutate(null);
  };

  return (
    <DaisyTooltip class={props.class} tooltip="Refresh">
      <DaisyButton
        class={mergeClass("btn-primary w-14", props.classButton)}
        loading={refreshRadioSubscription.loading()}
        onClick={refreshSubscription}
        aria-label="Refresh"
      >
        <FiRefreshCw size={20} />
      </DaisyButton>
    </DaisyTooltip>
  );
};

const RadioPlayerStatusButton: Component<
  {
    status: StateStatus;
    loading?: boolean;
  } & ClassProps
> = (props) => {
  const data = (): { class: string; status: string; element: JSX.Element } => {
    switch (props.status) {
      case StateStatus.StatusConnecting:
        return {
          class: "btn-circle btn-warning animate-spin",
          status: props.status,
          element: <FiRefreshCw size={20} />,
        };
      case StateStatus.StatusPlaying:
        return {
          class: "btn-circle btn-success pl-1",
          status: props.status,
          element: <FaSolidPlay size={20} />,
        };
      case StateStatus.StatusStopped:
        return {
          class: "btn-circle btn-error",
          status: props.status,
          element: <FaSolidStop size={20} />,
        };
      default:
        return {
          class: "btn-circle no-animation btn-info",
          status: "Unknown",
          element: <FaSolidQuestion size={20} />,
        };
    }
  };

  return (
    <DaisyTooltip class={props.class} tooltip={data().status}>
      <DaisyButton
        class={data().class}
        loading={props.loading}
        aria-label={data().status}
      >
        {data().element}
      </DaisyButton>
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
    <DaisyDropdown
      class={mergeClass("no-animation", props.class)}
      buttonProps={{
        class: mergeClass(
          "btn-primary justify-start gap-2 truncate",
          props.classButton
        ),
      }}
      buttonChildren={
        <>
          <span class="badge-info badge badge-lg rounded-md">
            {props.state.preset_number}
          </span>
          {props.state.title_new || props.state.title}
        </>
      }
      dropdownProps={{
        class: mergeClass(
          "card-compact card w-full bg-primary p-2 text-primary-content shadow",
          props.classDropdown
        ),
      }}
      loading={props.loading}
    >
      <DaisyStaticTableCardBody data={data()} title="Stream Information" />
    </DaisyDropdown>
  );
};

const RadioTypeDropdown: Component<
  {
    state: StateState;
    classButton?: string;
    classDropdown?: string;
  } & ClassProps
> = (props) => {
  const data = (): DaisyStaticTableCardBodyData[] => [
    { key: "Name", value: props.state.name },
    { key: "Model Name", value: props.state.model_name },
    { key: "Model Number", value: props.state.model_number },
  ];

  return (
    <DaisyDropdown
      class={props.class}
      buttonProps={{ class: "btn-primary" }}
      buttonChildren={<FaSolidRadio size={20} />}
      dropdownProps={{
        class: mergeClass(
          "card dropdown-content card-compact w-80 bg-primary p-2 text-primary-content shadow",
          props.classDropdown
        ),
      }}
    >
      <DaisyStaticTableCardBody data={data()} title="Radio Information" />
    </DaisyDropdown>
  );
};

const RadioPresetButton: Component<
  {
    selected: boolean;
    preset: StatePreset;
    loading?: boolean;
  } & JSX.HTMLAttributes<HTMLButtonElement>
> = (props) => {
  const [, other] = splitProps(props, [
    "selected",
    "preset",
    "loading",
    "class",
  ]);

  return (
    <DaisyButton
      class={mergeClass("flex gap-2", props.class)}
      classList={{ "btn-primary": props.selected }}
      loading={props.loading}
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
    </DaisyButton>
  );
};

const RadioPresetsList: Component<
  {
    radioUUID: Accessor<string>;
    state: StateState;
  } & ClassProps
> = (props) => {
  const statePatch = usePatchState(props.radioUUID);
  const [lastLoadingNumber, setLoadingNumber] = createSignal(-1);
  const loadingNumber = (): number =>
    statePatch.loading() ? lastLoadingNumber() : -1;

  const setPreset = (preset: number) => {
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
            onclick={[setPreset, preset.number]}
            preset={preset}
          />
        )}
      </For>
    </div>
  );
};

const RadioVolumeButtonGroup: Component<
  {
    radioUUID: Accessor<string>;
    state: StateState;
  } & ClassProps
> = (props) => {
  const statePatch = usePatchState(props.radioUUID);
  const refreshRadioVolume = useRefreshRadioVolume(props.radioUUID);

  const changeVolume = (volumeChange: number) => {
    void statePatch.mutate({ volume: props.state.volume + volumeChange });
  };

  const refreshVolume = () => {
    void refreshRadioVolume.mutate(null);
  };

  return (
    <Show
      when={!props.state.is_muted}
      fallback={
        <DaisyButton
          class={mergeClass("btn-error", props.class)}
          loading={statePatch.loading()}
          aria-label="Volume Muted"
        >
          <FaSolidVolumeOff size={20} />
        </DaisyButton>
      }
    >
      <div class={mergeClass("btn-group flex-nowrap", props.class)}>
        <DaisyButton
          class="btn-info w-14"
          loading={statePatch.loading()}
          aria-label="Lower Volume"
          onClick={[changeVolume, -5]}
        >
          <FaSolidVolumeLow size={20} />
        </DaisyButton>
        <DaisyButton
          class="btn-info w-12 px-0"
          loading={refreshRadioVolume.loading()}
          aria-label={`Volume ${props.state.volume}%`}
          onClick={refreshVolume}
        >
          {props.state.volume}%
        </DaisyButton>
        <DaisyButton
          class="btn-info w-14"
          loading={statePatch.loading()}
          aria-label="Raise Volume"
          onClick={[changeVolume, 5]}
        >
          <FaSolidVolumeHigh size={20} />
        </DaisyButton>
      </div>
    </Show>
  );
};

const RadioPowerButton: Component<
  {
    radioUUID: Accessor<string>;
    state: StateState;
  } & ClassProps
> = (props) => {
  const patchState = usePatchState(props.radioUUID);

  const togglePower = () => {
    void patchState.mutate({
      power: !props.state.power,
    });
  };

  return (
    <DaisyButton
      class={mergeClass("btn w-14", props.class)}
      classList={{
        "btn-success": props.state.power,
        "btn-error": !props.state.power,
      }}
      loading={patchState.loading()}
      onClick={togglePower}
      aria-label={"Power " + (props.state.power ? "On" : "Off")}
    >
      <Switch>
        <Match when={props.state.power}>
          <FaSolidPowerOff size={20} />
        </Match>
        <Match when={!props.state.power}>
          <FaSolidPowerOff size={20} />
        </Match>
      </Switch>
    </DaisyButton>
  );
};

const RadioAudioSourceDropdown: Component<
  {
    radioUUID: Accessor<string>;
    state: StateState;
    classButton?: string;
    classDropdown?: string;
  } & ClassProps
> = (props) => {
  const statePatch = usePatchState(props.radioUUID);
  const [lastLoadingAudioSource, setLastLoadingAudioSource] = createSignal("");
  const loadingAudioSource = (): string =>
    statePatch.loading() ? lastLoadingAudioSource() : "";

  const setAudioSource = (audioSource: string) => {
    batch(() => {
      statePatch.cancel();
      void statePatch.mutate({ audio_source: audioSource });
      setLastLoadingAudioSource(audioSource);
    });
  };

  return (
    <DaisyDropdown
      class={props.class}
      buttonProps={{
        "aria-label": "Audio Source",
        classList: { "btn-secondary": !!props.state.audio_source },
      }}
      buttonChildren={<FaBrandsItunesNote size={20} />}
      dropdownProps={{
        class: mergeClass(
          "menu rounded-box menu-compact w-52 space-y-2 bg-base-200 p-2 shadow",
          props.classDropdown
        ),
      }}
    >
      <span class="mx-auto">Audio Source</span>
      <For each={props.state.audio_sources}>
        {(a) => (
          <DaisyButton
            loading={loadingAudioSource() == a}
            classList={{ "btn-secondary": a == props.state.audio_source }}
            onClick={[setAudioSource, a]}
          >
            {a}
          </DaisyButton>
        )}
      </For>
    </DaisyDropdown>
  );
};

const RadioSelect: Component<
  {
    radioUUID: Accessor<string>;
    setRadioUUID: Setter<string>;
    radios: Resource<ModelRadio[]>;
  } & ClassProps
> = (props) => {
  let select: HTMLSelectElement | undefined;

  // Prevent select.value from defaulting to the first option when radios.data changes
  createEffect(
    on(
      props.radios,
      () => {
        select && (select.value = props.radioUUID());
      },
      { defer: true }
    )
  );

  return (
    <select
      class={mergeClass("select-primary select", props.class)}
      ref={select}
      disabled={props.radios.loading}
      value={props.radioUUID()}
      onChange={(e) => {
        props.setRadioUUID(e.currentTarget.value);
      }}
    >
      <option disabled value="">
        <Switch fallback={<>Select Radio</>}>
          <Match when={props.radios.loading}>Loading...</Match>
          <Match when={!!props.radios.error}>Network Error</Match>
          <Match when={props.radios()?.length == 0}>No Radios Found</Match>
        </Switch>
      </option>
      <Show when={!props.radios.error}>
        <For each={props.radios()}>
          {(radio) => <option value={radio.uuid}>{radio.name}</option>}
        </For>
      </Show>
    </select>
  );
};

const App: Component = () => {
  const [radioUUID, setRadioUUID] = createSignal(
    localStorage.getItem("lastRadioUUID") || ""
  );
  createEffect(() => {
    localStorage.setItem("lastRadioUUID", radioUUID());
  });

  const [data, ws] = useWS(radioUUID);
  const wsReconnecting = () => ws.connecting() && ws.disconnected();

  const { state, discovering } = data;
  const radioLoading = () => state.uuid != radioUUID() || ws.connecting();
  const radioLoaded = () => radioUUID() == state.uuid && ws.connected();

  const radiosListQuery = useRadiosListQuery();

  // Invalidate radios list based on websocket
  const trackRadiosListQuery = createReaction(() => {
    void radiosListQuery[1].refetch("");
  });
  createEffect(() => {
    (discovering() == true || ws.connecting() == true) &&
      trackRadiosListQuery(discovering);
  });
  createEffect(() => {
    ws.connected() == false && trackRadiosListQuery(ws.connected);
  });

  // Reconnect websocket when document is visible
  const onVisibilityChange = () => {
    if (!document.hidden) {
      ws.reconnect();
    }
  };
  document.addEventListener("visibilitychange", onVisibilityChange);
  onCleanup(() => {
    document.removeEventListener("visibilitychange", onVisibilityChange);
  });

  return (
    <div class="h-screen">
      <div class="fixed top-0 z-50 flex w-full gap-2 border-b-2 border-accent border-b-base-300 bg-base-200 p-2">
        <RadioPlayerStatusButton
          class="tooltip-bottom flex"
          status={state.status}
          loading={radioLoading()}
        />
        <RadioPlayerTitleDropdown
          class="dropdown-end flex-1"
          classButton="w-full"
          classDropdown="mt-2"
          state={state}
          loading={radioLoading()}
        />
      </div>
      <div class="container mx-auto px-4 pt-20 pb-36">
        <RadioPresetsList radioUUID={radioUUID} state={state} />
      </div>
      <div class="fixed bottom-0 z-50 w-full space-y-2 ">
        <div class="ml-auto max-w-screen-sm space-y-2 px-2">
          <Show when={ws.connecting() && !wsReconnecting()}>
            <div class="alert shadow-lg">
              <div>
                <span>Connecting to server...</span>
              </div>
            </div>
          </Show>
          <Show when={ws.disconnected()}>
            <div class="alert shadow-lg">
              <div>
                <span>Disconnected from server.</span>
              </div>
              <div class="flex-none">
                <DaisyButton
                  class="btn-primary btn-sm"
                  loading={ws.connecting()}
                  onClick={ws.reconnect}
                >
                  Reconnect
                </DaisyButton>
              </div>
            </div>
          </Show>
        </div>
        <div
          class="flex flex-wrap-reverse gap-2 border-t-2 border-accent border-t-base-300 bg-base-200 p-2"
          classList={{ "pb-5": IOS() }}
        >
          <div class="flex flex-auto gap-2">
            <div class="flex flex-1">
              <DiscoverButton
                classButton="rounded-r-none"
                discovering={discovering()}
              />
              <RadioSelect
                class="w-full min-w-fit flex-1 rounded-l-none"
                radioUUID={radioUUID}
                setRadioUUID={setRadioUUID}
                radios={radiosListQuery[0]}
              />
            </div>
            <Show when={radioLoaded()}>
              <RadioRefreshSubscriptionButton radioUUID={radioUUID} />
            </Show>
          </div>
          <Show when={radioLoaded()}>
            <div class="flex flex-auto gap-2">
              <RadioPowerButton
                class="flex-auto"
                radioUUID={radioUUID}
                state={state}
              />
              <RadioVolumeButtonGroup radioUUID={radioUUID} state={state} />
              <RadioAudioSourceDropdown
                class="dropdown-top dropdown-end"
                classDropdown="mb-2"
                radioUUID={radioUUID}
                state={state}
              />
              <RadioTypeDropdown
                class="dropdown-top dropdown-end"
                classDropdown="mb-2"
                state={state}
              />
            </div>
          </Show>
        </div>
      </div>
    </div>
  );
};

export default App;
