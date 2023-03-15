import {
  type ResourceReturn,
  type Accessor,
  createSignal,
  createEffect,
  on,
  type Signal,
} from "solid-js";

export function createStaleSignal<T>(value: T): Signal<T> {
  return createSignal(value, { equals: false });
}

export function invalidWhen<T, R>(
  deps: Accessor<unknown>,
  query: ResourceReturn<T, R>
): ResourceReturn<T, R> {
  createEffect(
    on(
      deps,
      () => {
        void query[1].refetch();
      },
      { defer: true }
    )
  );
  return query;
}

export function staleWhen<T, R>(
  query: ResourceReturn<T, R>,
  fn: Accessor<boolean>,
  deps: Accessor<unknown> = fn
): ResourceReturn<T, R> {
  createEffect(
    on(
      deps,
      () => {
        if (fn()) {
          void query[1].refetch();
        }
      },
      { defer: true }
    )
  );
  return query;
}

export function once<T>(fn: () => T): () => T {
  let data: T | undefined;
  return function () {
    if (data === undefined) data = fn();

    return data;
  };
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
