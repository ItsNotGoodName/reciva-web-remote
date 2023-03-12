import {
  type ResourceReturn,
  type Accessor,
  createSignal,
  createEffect,
  on,
  type Signal,
  onCleanup,
} from "solid-js";

export function createStaleSignal<T>(value: T): Signal<T> {
  return createSignal(value, { equals: false });
}

export function checkStale<T, R>(
  query: ResourceReturn<T, R>,
  fn: Accessor<boolean>,
  deps: Accessor<unknown> = fn
): ResourceReturn<T, R> {
  createEffect(
    on(
      deps,
      () => {
        console.log("Stale: Check");
        if (fn()) {
          void query[1].refetch();
          console.log("Stale: Refetch");
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

export function clickOutside(el: HTMLInputElement) {
  const onClick = (e: MouseEvent) => {
    !el.contains(e.target as Node) && el.blur();
  };
  document.body.addEventListener("click", onClick);

  onCleanup(() => document.body.removeEventListener("click", onClick));
}

declare module "solid-js" {
  // eslint-disable-next-line @typescript-eslint/no-namespace
  namespace JSX {
    interface Directives {
      // use:clickOutside
      clickOutside: string;
    }
  }
}
