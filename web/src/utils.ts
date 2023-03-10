import {
  type ResourceReturn,
  type Accessor,
  createSignal,
  createEffect,
  on,
  type Signal,
} from "solid-js";

export type Mutation<T = void, R = unknown> = {
  mutate: (data: T) => R | Promise<R>;
  loading: Accessor<boolean>;
};

export function createStaleSignal<T>(value: T): Signal<T> {
  return createSignal(value, { equals: false });
}

export function checkStale<T, R>(
  query: ResourceReturn<T, R>,
  checkFn: Accessor<boolean>,
  onD: Accessor<unknown> = checkFn
): ResourceReturn<T, R> {
  createEffect(
    on(
      onD,
      () => {
        console.log("stale check");
        if (checkFn()) {
          void query[1].refetch();
          console.log("stale REFETCH");
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

export function createMutation<T, R>(
  mutateFn: (data: T) => R | Promise<R>,
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  invalidateQueries: ResourceReturn<any, any>[] = []
): Mutation<T, R> {
  const [loading, setloading] = createSignal(false);
  let mutationResult: R | Promise<R>;

  function mutate(data: T): Promise<R> | R {
    if (loading()) return mutationResult;

    setloading(true);

    mutationResult = mutateFn(data);

    void Promise.resolve(mutationResult)
      .finally(() => {
        setloading(false);
      })
      .then(() =>
        Promise.allSettled(
          invalidateQueries.map((query) => void query[1].refetch())
        )
      );

    return mutationResult;
  }

  return { loading, mutate };
}
