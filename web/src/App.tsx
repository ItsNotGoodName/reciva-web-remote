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
  FaSolidArrowsRotate,
  FaSolidMagnifyingGlass,
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
  batch,
  type Accessor,
  type Resource,
  onCleanup,
  Index,
  createReaction,
  Suspense,
  ErrorBoundary,
} from "solid-js";
import { type Component } from "solid-js";
import {
  StateStatus,
  type StateState,
  type ModelRadio,
  type ModelBuild,
  ModelStale,
} from "./api";
import {
  useDiscoverRadios,
  useUpdateState,
  useRefreshRadioVolume,
  useRefreshRadioSubscription,
  useRadiosListQuery,
  useBuildGetQuery,
  useDeleteRadio,
  invalidatePresetListQuery,
  invalidateRadioListQuery,
} from "./store";
import { type ClassProps, mergeClass, IOS } from "./utils";
import { useWS } from "./ws";
import {
  DaisyButton,
  DaisyTooltip,
  DaisyDropdown,
  DaisyErrorAlert,
} from "./Daisy";
import { GITHUB_URL, ICON_SIZE, PAGE_EDIT, PAGE_HOME } from "./constants";
import { EditPage } from "./pages/Edit";
import { HomePage } from "./pages/Home";

type TableRowData = { key: string; value: JSX.Element };

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
  const data = (): TableRowData[] => [
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
      class={props.class}
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
      <div class="card-body">
        <h2 class="card-title">Stream Information</h2>
        <table class="table-fixed">
          <tbody>
            <Index each={data()}>
              {(data) => (
                <tr>
                  <td class="pb-1 pr-1">
                    <span class="text badge-info badge w-full whitespace-nowrap">
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
    </DaisyDropdown>
  );
};

const RadioTypeDropdown: Component<
  {
    radioUUID: Accessor<string>;
    state: StateState;
  } & ClassProps
> = (props) => {
  const data = (): TableRowData[] => [
    { key: "Name", value: props.state.name },
    { key: "Model Name", value: props.state.model_name },
    { key: "Model Number", value: props.state.model_number },
  ];

  const refreshRadioSubscription = useRefreshRadioSubscription(props.radioUUID);
  const deleteRadio = useDeleteRadio(props.radioUUID);

  return (
    <DaisyDropdown
      class={props.class}
      buttonClass="btn-primary"
      buttonChildren={<FaSolidRadio size={ICON_SIZE} />}
      dropdownClass="card-compact card bg-primary p-2 text-primary-content shadow my-2"
    >
      <div class="card-body">
        <h2 class="card-title">Radio Information</h2>
        <table class="table-fixed">
          <tbody>
            <Index each={data()}>
              {(d) => (
                <tr>
                  <td class="pb-1 pr-1">
                    <span class="text badge-info badge w-full whitespace-nowrap">
                      {d().key}
                    </span>
                  </td>
                  <td>
                    <div class="whitespace-nowrap">{d().value}</div>
                  </td>
                </tr>
              )}
            </Index>
          </tbody>
        </table>
        <div class="card-actions">
          <DaisyButton
            class="btn-info btn-sm ml-auto flex-1"
            loading={refreshRadioSubscription.loading()}
            onClick={() => void refreshRadioSubscription.mutate(null)}
          >
            Refresh
          </DaisyButton>
          <DaisyButton
            class="btn-error btn-sm ml-auto flex-1"
            loading={deleteRadio.loading()}
            onClick={() => void deleteRadio.mutate(null)}
          >
            Remove
          </DaisyButton>
        </div>
      </div>
    </DaisyDropdown>
  );
};

const RadioVolumeButtonGroup: Component<
  {
    radioUUID: Accessor<string>;
    state: StateState;
  } & ClassProps
> = (props) => {
  const updateState = useUpdateState(props.radioUUID);
  const refreshRadioVolume = useRefreshRadioVolume(props.radioUUID);

  const changeVolume = (volumeChange: number) => {
    void updateState.mutate({ volume_delta: volumeChange });
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
          loading={updateState.loading()}
          aria-label="Volume Muted"
        >
          <FaSolidVolumeOff size={ICON_SIZE} />
        </DaisyButton>
      }
    >
      <div class={mergeClass("btn-group flex-nowrap", props.class)}>
        <DaisyButton
          class="btn-info w-14"
          loading={updateState.loading()}
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
          loading={updateState.loading()}
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
  const updateState = useUpdateState(props.radioUUID);

  const togglePower = () => {
    void updateState.mutate({
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
      loading={updateState.loading()}
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
  const updateState = useUpdateState(props.radioUUID);
  const [lastLoadingAudioSource, setLastLoadingAudioSource] = createSignal("");
  const loadingAudioSource = (): string =>
    updateState.loading() ? lastLoadingAudioSource() : "";

  const setAudioSource = (audioSource: string) => {
    batch(() => {
      updateState.cancel();
      void updateState.mutate({ audio_source: audioSource });
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

const RadioListCard: Component<
  {
    radioUUID: Accessor<string>;
    setRadioUUID: Setter<string>;
    radios: Resource<ModelRadio[]>;
    discovering: Accessor<boolean>;
  } & ClassProps
> = (props) => {
  return (
    <div
      class={mergeClass(
        "card w-full bg-base-100 shadow-xl sm:w-96",
        props.class
      )}
    >
      <div class="card-body">
        <h2 class="card-title">
          <Switch fallback={<>Select Radio</>}>
            <Match when={props.discovering()}>Discovering...</Match>
            <Match
              when={
                !props.radios.error &&
                !props.radios.loading &&
                props.radios()?.length == 0
              }
            >
              No Radios Discoverd
            </Match>
          </Switch>
        </h2>
        <Suspense fallback={<>Loading...</>}>
          <Show when={!props.radios.error}>
            <For each={props.radios()}>
              {(radio) => (
                <DaisyButton
                  class="btn-primary"
                  onClick={[props.setRadioUUID, radio.uuid]}
                  value={radio.uuid}
                >
                  {radio.name}
                </DaisyButton>
              )}
            </For>
          </Show>
        </Suspense>
      </div>
    </div>
  );
};

const RadioListSelect: Component<
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

  // Queries
  const buildGetQuery = useBuildGetQuery();
  const radiosListQuery = useRadiosListQuery();

  // invalidations
  createEffect(
    on(
      stale,
      () => {
        // Refetch stale queries
        stale() == ModelStale.StaleRadios &&
          invalidateRadioListQuery(new Date());
        stale() == ModelStale.StalePresets &&
          invalidatePresetListQuery(new Date());
      },
      { defer: true }
    )
  );
  const track = createReaction(() => {
    // Refetch queries
    invalidateRadioListQuery(new Date());
    invalidatePresetListQuery(new Date());
  });
  createEffect(() => {
    ws.disconnected() && track(ws.connected);
  });
  createEffect(
    on(ws.connected, () => {
      // Refetch failed queries
      buildGetQuery[0].error && buildGetQuery[1].refetch();
    })
  );

  // Reconnect websocket when document is visible
  const onFocus = () => {
    if (!document.hidden) {
      ws.reconnect();
    }
  };
  document.addEventListener("visibilitychange", onFocus);
  window.addEventListener("focus", onFocus);
  window.addEventListener("online", onFocus);
  onCleanup(() => {
    document.removeEventListener("visibilitychange", onFocus);
    window.removeEventListener("focus", onFocus);
    window.removeEventListener("online", onFocus);
  });

  const [page, setPage] = createSignal(PAGE_HOME);
  const radioSelected = () => radioUUID() != "";
  const radioLoading = () =>
    (state.uuid != radioUUID() && radioSelected()) || ws.connecting();
  const radioLoaded = () =>
    radioUUID() == state.uuid && radioSelected() && ws.connected();

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
            <Switch>
              <Match when={!radioSelected()}>
                <RadioListCard
                  class="mx-auto"
                  radioUUID={radioUUID}
                  setRadioUUID={setRadioUUID}
                  radios={radiosListQuery[0]}
                  discovering={discovering}
                />
              </Match>
              <Match when={radioLoaded()}>
                <HomePage radioUUID={radioUUID} state={state} />
              </Match>
            </Switch>
          </Match>
          <Match when={page() == PAGE_EDIT}>
            <EditPage />
          </Match>
        </Switch>
      </div>
      <div class="fixed bottom-0 z-50 w-full space-y-2 ">
        <div class="ml-auto max-w-screen-sm space-y-2 px-2">
          <Show when={!!radiosListQuery[0].error}>
            <DaisyErrorAlert>Failed to list radios.</DaisyErrorAlert>
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
              <RadioListSelect
                class="w-full flex-1 rounded-l-none"
                radioUUID={radioUUID}
                setRadioUUID={setRadioUUID}
                radios={radiosListQuery[0]}
              />
            </div>
            <Show when={radioLoaded()}>
              <RadioTypeDropdown
                class="dropdown-top dropdown-end"
                radioUUID={radioUUID}
                state={state}
              />
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
            </div>
          </Show>
        </div>
      </div>
    </div>
  );
};

export default () => (
  <ErrorBoundary
    fallback={(err) => (
      <div class="m-4">
        <DaisyErrorAlert>{err || "Something went wrong."}</DaisyErrorAlert>
      </div>
    )}
  >
    <App />
  </ErrorBoundary>
);
