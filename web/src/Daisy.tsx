import { FaSolidCircleXmark } from "solid-icons/fa";
import { type JSX, type ParentComponent, Show, splitProps } from "solid-js";
import { ICON_SIZE } from "./constants";
import { type ClassProps, mergeClass, useDropdown } from "./utils";

export const DaisyErrorAlert: ParentComponent = (props) => {
  return (
    <div class="alert alert-error shadow-lg">
      <div>
        <FaSolidCircleXmark size={ICON_SIZE} />
        {props.children}
      </div>
    </div>
  );
};

export const DaisyButton: ParentComponent<
  { loading?: boolean } & JSX.ButtonHTMLAttributes<HTMLButtonElement>
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

export const DaisyTooltip: ParentComponent<{ tooltip: string } & ClassProps> = (
  props
) => {
  return (
    <div class={mergeClass("tooltip", props.class)} data-tip={props.tooltip}>
      {props.children}
    </div>
  );
};

export const DaisyDropdown: ParentComponent<
  {
    buttonClass?: string;
    buttonClassList?: {
      [k: string]: boolean | undefined;
    };
    buttonChildren?: JSX.Element;
    dropdownClass?: string;
    loading?: boolean;
  } & ClassProps
> = (props) => {
  const { showDropdown, toggleDropdown } = useDropdown();

  return (
    <div
      class={mergeClass("dropdown no-animation", props.class)}
      classList={{ "dropdown-open": showDropdown() }}
      onFocusOut={toggleDropdown}
    >
      <label
        tabindex="0"
        class={mergeClass("btn touch-manipulation", props.buttonClass)}
        classList={{ loading: props.loading, ...props.buttonClassList }}
        onClick={toggleDropdown}
      >
        <Show when={!props.loading}>{props.buttonChildren}</Show>
      </label>
      <div
        tabindex="0"
        class={mergeClass("dropdown-content", props.dropdownClass)}
      >
        {props.children}
      </div>
    </div>
  );
};
