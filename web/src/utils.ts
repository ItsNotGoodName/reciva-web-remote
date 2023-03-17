import {
  type ResourceReturn,
  type Accessor,
  createSignal,
  createEffect,
  on,
  type Signal,
} from "solid-js";

export function createInvalidateSignal<T>(value: T): Signal<T> {
  return createSignal(value, { equals: false });
}

export function invalidateWhen<T, R>(
  dep: Accessor<boolean>,
  query: ResourceReturn<T, R>
): ResourceReturn<T, R> {
  createEffect(
    on(
      dep,
      () => {
        if (dep()) {
          void query[1].refetch();
        }
      },
      { defer: true }
    )
  );
  return query;
}

export type ClassProps = { class?: string };

export function mergeClass(first: string, second?: string) {
  if (second) {
    return first + " " + second;
  }
  return first;
}

export function useDropdown() {
  const [show, setShow] = createSignal(false);
  const toggleShow = (
    e: Event & {
      currentTarget: HTMLElement;
      target: Element;
    }
  ) => {
    if (e.type == "focusout") {
      setShow(false);
      return;
    }

    if (!show()) {
      setShow(true);
    } else {
      setShow(false);
      e.currentTarget.blur();
    }
  };

  return { showDropdown: show, toggleDropdown: toggleShow };
}

// https://stackoverflow.com/a/9039885
export function IOS() {
  return (
    [
      "iPad Simulator",
      "iPhone Simulator",
      "iPod Simulator",
      "iPad",
      "iPhone",
      "iPod",
    ].includes(navigator.platform) ||
    // iPad on iOS 13 detection
    (navigator.userAgent.includes("Mac") && "ontouchend" in document)
  );
}
