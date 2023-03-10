import { type ResourceReturn, type Accessor, createSignal } from "solid-js";

export type Mutation<T = void, R = unknown> = {
  mutate: (data: T) => R | Promise<R>;
  loading: Accessor<boolean>;
};

export function createMutation<T, R>(
  mutateFn: (data: T) => R | Promise<R>,
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
