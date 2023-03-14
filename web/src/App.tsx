import {
  FaBrandsItunesNote,
  FaSolidCircleInfo,
  FaSolidPlay,
  FaSolidPowerOff,
  FaSolidQuestion,
  FaSolidRadio,
  FaSolidStop,
  FaSolidVolumeHigh,
  FaSolidVolumeLow,
  FaSolidVolumeOff,
  FaSolidArrowRotateRight,
  FaSolidArrowsRotate,
  FaSolidMagnifyingGlass,
  FaSolidCircleXmark,
  FaSolidBars,
  FaBrandsGithub,
  FaSolidSpinner,
  FaSolidTag,
  FaSolidHouse,
  FaSolidPen,
} from "solid-icons/fa";
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
  onCleanup,
} from "solid-js";
import { type Component } from "solid-js";
import {
  type StatePreset,
  StateStatus,
  type StateState,
  type ModelRadio,
  type ModelBuild,
  ModelStale,
} from "./api";
import {
  useDiscoverRadios,
  usePatchState,
  useRefreshRadioVolume,
  useRefreshRadioSubscription,
  useRadiosListQuery,
  useBuildGetQuery,
  useDeleteRadio,
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
import { GITHUB_URL, ICON_SIZE, PAGE_EDIT, PAGE_HOME } from "./constants";

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
        <FaSolidMagnifyingGlass size={ICON_SIZE} />
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
        <FaSolidArrowRotateRight size={ICON_SIZE} />
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
          element: <FaSolidArrowsRotate size={ICON_SIZE} />,
        };
      case StateStatus.StatusPlaying:
        return {
          class: "btn-circle btn-success pl-1",
          status: props.status,
          element: <FaSolidPlay size={ICON_SIZE} />,
        };
      case StateStatus.StatusStopped:
        return {
          class: "btn-circle btn-error",
          status: props.status,
          element: <FaSolidStop size={ICON_SIZE} />,
        };
      default:
        return {
          class: "btn-circle btn-info",
          status: "Unknown",
          element: <FaSolidQuestion size={ICON_SIZE} />,
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
      buttonClass="btn-primary w-full justify-start truncate"
      buttonChildren={
        <div class="w-0">
          <span class="badge-info badge badge-lg mr-2 rounded-md">
            {props.state.preset_number}
          </span>
          {props.state.title_new || props.state.title}
        </div>
      }
      dropdownClass="card-compact card w-full bg-primary p-2 text-primary-content shadow my-2"
      loading={props.loading}
    >
      <DaisyStaticTableCardBody data={data()} title="Stream Information" />
    </DaisyDropdown>
  );
};

const RadioTypeDropdown: Component<
  {
    radioUUID: Accessor<string>;
    state: StateState;
  } & ClassProps
> = (props) => {
  const data = (): DaisyStaticTableCardBodyData[] => [
    { key: "Name", value: props.state.name },
    { key: "Model Name", value: props.state.model_name },
    { key: "Model Number", value: props.state.model_number },
  ];

  const deleteRadio = useDeleteRadio(props.radioUUID);

  return (
    <DaisyDropdown
      class={props.class}
      buttonClass="btn-primary"
      buttonChildren={<FaSolidRadio size={ICON_SIZE} />}
      dropdownClass="card-compact card w-80 bg-primary p-2 text-primary-content shadow my-2"
    >
      <DaisyStaticTableCardBody data={data()} title="Radio Information" />
      <div class="mx-2 mb-2 flex">
        <DaisyButton
          class="btn-error btn-sm ml-auto"
          loading={deleteRadio.loading()}
          onClick={() => void deleteRadio.mutate(null)}
        >
          Remove
        </DaisyButton>
      </div>
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
          <FaSolidVolumeOff size={ICON_SIZE} />
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
          <FaSolidVolumeLow size={ICON_SIZE} />
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
          <FaSolidVolumeHigh size={ICON_SIZE} />
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
          <FaSolidPowerOff size={ICON_SIZE} />
        </Match>
        <Match when={!props.state.power}>
          <FaSolidPowerOff size={ICON_SIZE} />
        </Match>
      </Switch>
    </DaisyButton>
  );
};

const RadioAudioSourceDropdown: Component<
  {
    radioUUID: Accessor<string>;
    state: StateState;
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
      aria-label="Audio Source"
      buttonClassList={{ "btn-secondary": !!props.state.audio_source }}
      buttonChildren={<FaBrandsItunesNote size={ICON_SIZE} />}
      dropdownClass="menu rounded-box menu-compact w-52 space-y-2 bg-base-200 p-2 shadow my-2"
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

  // Prevent select.value from defaulting to the first option when props.radios() changes
  createEffect(
    on(
      () => !props.radios.error && props.radios(),
      () => {
        select && (select.value = props.radioUUID());

        if (!props.radios.error) {
          for (const r of props.radios() || []) {
            if (r.uuid == props.radioUUID()) {
              return;
            }
          }
          props.setRadioUUID("");
        }
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
        Select Radio
      </option>
      <Show when={!props.radios.error}>
        <For each={props.radios()}>
          {(radio) => <option value={radio.uuid}>{radio.name}</option>}
        </For>
      </Show>
    </select>
  );
};

const MenuDropdown: Component<
  {
    build: Resource<ModelBuild>;
    page: Accessor<string>;
    setPage: Setter<string>;
  } & ClassProps
> = (props) => {
  return (
    <DaisyDropdown
      class={props.class}
      buttonClassList={{ "btn-success": props.page() != PAGE_HOME }}
      buttonChildren={<FaSolidBars size={ICON_SIZE} />}
    >
      <ul class="menu rounded-box menu-compact mb-2 w-52 min-w-max bg-base-200 p-2 shadow">
        <li>
          <a
            classList={{ active: props.page() == PAGE_HOME }}
            onClick={() => props.setPage(PAGE_HOME)}
          >
            <FaSolidHouse size={ICON_SIZE} />
            Home Page
          </a>
        </li>
        <li>
          <a
            classList={{ active: props.page() == PAGE_EDIT }}
            onClick={() => props.setPage(PAGE_EDIT)}
          >
            <FaSolidPen size={ICON_SIZE} />
            Edit Presets
          </a>
        </li>
        <li>
          <a href={GITHUB_URL}>
            <FaBrandsGithub size={ICON_SIZE} />
            Source Code
          </a>
        </li>
        <Switch>
          <Match when={props.build.loading}>
            <li>
              <a>
                <FaSolidSpinner size={ICON_SIZE} class="animate-spin" /> Version
              </a>
            </li>
          </Match>
          <Match when={!props.build.error}>
            <li>
              <a href={props.build()?.release_url || "#"}>
                <FaSolidTag size={ICON_SIZE} />
                {props.build()?.summary}
              </a>
            </li>
          </Match>
        </Switch>
      </ul>
    </DaisyDropdown>
  );
};

const App: Component = () => {
  const [radioUUID, setRadioUUID] = createSignal(
    localStorage.getItem("lastRadioUUID") || ""
  );
  createEffect(() => {
    localStorage.setItem("lastRadioUUID", radioUUID());
  });

  // WebSocket
  const [{ state, discovering, stale }, ws] = useWS(radioUUID);
  const wsReconnecting = () => ws.connecting() && ws.disconnected();

  const radioLoading = () =>
    (state.uuid != radioUUID() && radioUUID() != "") || ws.connecting();
  const radioLoaded = () =>
    radioUUID() == state.uuid && radioUUID() != "" && ws.connected();

  // Queries
  const buildGetQuery = useBuildGetQuery();
  const radiosListQuery = useRadiosListQuery();

  // Invalidate radios list based on WebSocket
  createEffect(() => {
    stale() == ModelStale.StaleRadios && void radiosListQuery[1].refetch();
  });

  // Reconnect websocket when document is visible
  const onVisibilityChange = () => {
    if (!document.hidden) {
      ws.reconnect();
    }
  };
  document.addEventListener("visibilitychange", onVisibilityChange);
  window.addEventListener("focus", onVisibilityChange);
  onCleanup(() => {
    document.removeEventListener("visibilitychange", onVisibilityChange);
    window.removeEventListener("focus", onVisibilityChange);
  });

  const [page, setPage] = createSignal(PAGE_HOME);

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
          state={state}
          loading={radioLoading()}
        />
      </div>
      <div class="container mx-auto px-4 pt-20 pb-36">
        <Switch>
          <Match when={page() == PAGE_HOME}>
            <RadioPresetsList radioUUID={radioUUID} state={state} />
          </Match>
          <Match when={page() == PAGE_EDIT}>
            <h1>Edit Presets</h1>
          </Match>
        </Switch>
      </div>
      <div class="fixed bottom-0 z-50 w-full space-y-2 ">
        <div class="ml-auto max-w-screen-sm space-y-2 px-2">
          <Show when={!!radiosListQuery[0].error}>
            <div class="alert alert-error shadow-lg">
              <div>
                <FaSolidCircleXmark size={ICON_SIZE} />
                <span>Failed to list radios.</span>
              </div>
            </div>
          </Show>
          <Show when={ws.connecting() && !wsReconnecting()}>
            <div class="alert shadow-lg">
              <div>
                <FaSolidCircleInfo class="fill-info" size={ICON_SIZE} />
                <span>Connecting to server...</span>
              </div>
            </div>
          </Show>
          <Show when={ws.disconnected()}>
            <div class="alert shadow-lg">
              <div>
                <FaSolidCircleInfo class="fill-info" size={ICON_SIZE} />
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
            <MenuDropdown
              class="dropdown-top"
              build={buildGetQuery[0]}
              page={page}
              setPage={setPage}
            />
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
                radioUUID={radioUUID}
                state={state}
              />
              <RadioTypeDropdown
                class="dropdown-top dropdown-end"
                radioUUID={radioUUID}
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
