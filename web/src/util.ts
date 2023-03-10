import { ResourceReturn, createResource } from "solid-js";

// TODO: remove this garbage
export function createMutation<T, R = unknown>(
  func: (k: T) => R | Promise<R>,
  invalidateQueries: ResourceReturn<any, any>[] = []
): ResourceReturn<R, T> {
  return createResource<any, object, T>({}, (_, info) => {
    if (typeof info.refetching === "boolean") {
      if (info.refetching == true)
        throw new Error(
          "You need to pass something to the refetch function..."
        );
      return;
    }
    const result = func(info.refetching);

    Promise.resolve(result).then(() =>
      Promise.allSettled(invalidateQueries.map((query) => query[1].refetch()))
    );

    return result;
  });
}
