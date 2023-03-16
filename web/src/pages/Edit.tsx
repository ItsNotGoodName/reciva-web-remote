import { FaSolidPen } from "solid-icons/fa";
import {
  type Component,
  type Resource,
  type Accessor,
  Switch,
  Match,
  For,
  createSignal,
  batch,
  Show,
  createEffect,
  on,
} from "solid-js";
import { createStore } from "solid-js/store";
import { type ModelPreset } from "../api";
import { DaisyButton, DaisyErrorAlert } from "../Daisy";
import { usePresetListQuery, usePresetQuery, useUpdatePreset } from "../store";
import { type ClassProps, mergeClass } from "../utils";

export const EditPage: Component = () => {
  const [presets] = usePresetListQuery();
  const [presetURL, setPresetURL] = createSignal("");
  const [presetFormOpen, setPresetFormOpen] = createSignal(false);

  const updatePresetUrl = (url: string) => {
    batch(() => {
      setPresetURL(url);
      setPresetFormOpen(true);
    });
    window.scrollTo(0, 0);
  };

  const close = () => {
    setPresetFormOpen(false);
  };

  return (
    <div class="flex flex-wrap-reverse gap-4">
      <PresetsTable
        class="flex-auto md:flex-1"
        presets={presets}
        activeURL={presetURL}
        setActiveURL={updatePresetUrl}
      />
      <Show when={presetFormOpen()}>
        <div class="w-96 flex-1 md:flex-initial">
          <div class="rounded-lg border-2 border-base-200 bg-base-200 p-4">
            <PresetForm presetURL={presetURL} onClose={close} />
          </div>
        </div>
      </Show>
    </div>
  );
};

const PresetsTable: Component<
  {
    presets: Resource<ModelPreset[]>;
    activeURL: Accessor<string>;
    setActiveURL: (url: string) => void;
  } & ClassProps
> = (props) => {
  return (
    <div class={mergeClass("overflow-x-auto", props.class)}>
      <Switch>
        <Match when={!!props.presets.error}>
          <DaisyErrorAlert>Failed to fetch presets.</DaisyErrorAlert>
        </Match>
        <Match when={props.presets() != undefined}>
          <table class="table-compact table w-full">
            <thead>
              <tr>
                <th></th>
                <th>URL</th>
                <th>New Title</th>
              </tr>
            </thead>
            <tbody>
              <For each={props.presets()}>
                {(p) => (
                  <tr classList={{ active: p.url == props.activeURL() }}>
                    <th class="w-0">
                      <DaisyButton
                        class="btn-success btn-sm"
                        aria-label="Edit"
                        onClick={[props.setActiveURL, p.url]}
                      >
                        <FaSolidPen />
                      </DaisyButton>
                    </th>
                    <td class="w-0">{p.url}</td>
                    <td class="w-0">{p.title_new}</td>
                  </tr>
                )}
              </For>
            </tbody>
          </table>
        </Match>
        <Match when={props.presets.loading}>Loading...</Match>
      </Switch>
    </div>
  );
};

const DefaultPresetForm = {
  title_new: "",
  title_new_change: false,
  url: "",
  url_new: "",
  url_new_change: false,
};

const PresetForm: Component<
  { presetURL: Accessor<string>; onClose: () => void } & ClassProps
> = (props) => {
  const [preset, presetQuery] = usePresetQuery(props.presetURL);
  const [presetForm, setPresetForm] = createStore({ ...DefaultPresetForm });
  createEffect(
    on(
      () => !preset.error && preset(),
      () => {
        if (preset.loading) {
          return;
        } else if (preset.error) {
          return;
        }

        const p = preset();
        if (p) {
          setPresetForm({ ...DefaultPresetForm, ...p });
        }
      }
    )
  );

  const updatePreset = useUpdatePreset();
  const submit = (e: Event) => {
    e.preventDefault();
    void updatePreset.mutate({
      title_new: presetForm.title_new,
      url: presetForm.url,
      url_new: presetForm.url_new,
    });
  };
  const reset = () => {
    void presetQuery.refetch();
  };

  const loading = () => preset.loading || updatePreset.loading();

  return (
    <form class={mergeClass("space-y-2", props.class)} onSubmit={submit}>
      <h1 class="text-center text-2xl">Edit Preset</h1>
      <Show when={!!preset.error}>
        <DaisyErrorAlert>Failed to get preset.</DaisyErrorAlert>
      </Show>
      <div class="form-control">
        <label class="label">
          <span class="label-text">URL</span>
        </label>
        <input
          type="text"
          placeholder="URL"
          class="input-bordered input"
          readonly
          value={presetForm.url}
          disabled={loading()}
        />
      </div>
      <div class="form-control">
        <label class="label">
          <span class="label-text">New Title</span>
        </label>
        <input
          classList={{ "input-warning": presetForm.title_new_change }}
          type="text"
          placeholder="New Title"
          class="input-bordered input"
          value={presetForm.title_new}
          onInput={(e) => {
            setPresetForm({
              title_new: e.currentTarget.value,
              title_new_change: true,
            });
          }}
          disabled={loading()}
        />
      </div>
      <div class="form-control">
        <label class="label">
          <span class="label-text">New URL</span>
        </label>
        <textarea
          placeholder="New URL"
          class="textarea-bordered textarea h-24"
          classList={{ "textarea-warning": presetForm.url_new_change }}
          value={presetForm.url_new}
          onInput={(e) => {
            setPresetForm({
              url_new: e.currentTarget.value,
              url_new_change: true,
            });
          }}
          disabled={loading()}
        />
      </div>
      <div class="btn-group flex pt-2">
        <DaisyButton class="flex-1" onClick={props.onClose} type="button">
          Close
        </DaisyButton>
        <DaisyButton
          class="btn-error flex-1"
          onClick={reset}
          type="button"
          loading={loading()}
        >
          Reset
        </DaisyButton>
        <DaisyButton
          class="btn-success flex-1"
          type="submit"
          loading={loading()}
        >
          Save
        </DaisyButton>
      </div>
      <Show when={updatePreset.error()}>
        <DaisyErrorAlert>Failed to update preset</DaisyErrorAlert>
      </Show>
    </form>
  );
};
