import {
  type JSX,
  type ParentComponent,
  Show,
  splitProps,
  Index,
  type Component,
} from "solid-js";
import { type ClassProps, mergeClass } from "./utils";

export const DaisyButton: ParentComponent<
  { loading?: boolean } & JSX.HTMLAttributes<HTMLButtonElement>
> = (props) => {
  const [, other] = splitProps(props, [
    "class",
    "classList",
    "loading",
    "children",
  ]);

  return (
    <button
      class={mergeClass("btn touch-manipulation", props.class)}
      classList={{ loading: props.loading, ...props.classList }}
      {...other}
    >
      <Show when={!props.loading}>{props.children}</Show>
    </button>
  );
};

export const DaisyDropdownButton: ParentComponent<
  { loading?: boolean } & JSX.HTMLAttributes<HTMLLabelElement>
> = (props) => {
  const [, other] = splitProps(props, [
    "class",
    "classList",
    "loading",
    "children",
  ]);

  return (
    <label
      tabindex="0"
      class={mergeClass("btn touch-manipulation", props.class)}
      classList={{ loading: props.loading, ...props.classList }}
      {...other}
    >
      <Show when={!props.loading}>{props.children}</Show>
    </label>
  );
};

export const DaisyTooltip: ParentComponent<{ tooltip: string } & ClassProps> = (
  props
) => {
  return (
    <div class={mergeClass("tooltip", props.class)} data-tip={props.tooltip}>
      {props.children}
    </div>
  );
};

export type DaisyStaticTableCardBodyData = { key: string; value: JSX.Element };

export const DaisyStaticTableCardBody: Component<{
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
