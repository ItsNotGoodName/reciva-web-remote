import {
  type Component,
  type Accessor,
  createSignal,
  batch,
  For,
  type JSX,
  Show,
  splitProps,
} from "solid-js";
import { type StatePreset, type StateState } from "../api";
import { DaisyButton } from "../Daisy";
import { usePatchState } from "../store";
import { type ClassProps, mergeClass } from "../utils";

export const HomePage: Component<
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
