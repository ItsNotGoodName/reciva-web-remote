import {
  QueryClient,
  createQuery,
  QueryClientProvider,
} from "@tanstack/solid-query";
import { Switch, Match, For } from "solid-js";
import type { Component } from "solid-js";
import { Api } from "./api";
import { API_URL } from "./constant";

const queryClient = new QueryClient();
const api = new Api({ baseUrl: API_URL + "/api" }); // TODO: get api path from swagger.json

function Example() {
  const query = createQuery(
    () => ["radios"],
    () => api.radios.radiosList()
  );

  return (
    <div>
      <Switch>
        <Match when={query.isLoading}>
          <p>Loading...</p>
        </Match>
        <Match when={query.isError}>
          <p>Error: unknown error</p>
        </Match>
        <Match when={query.isSuccess}>
          <For each={query.data.data}>{(radio) => <p>{radio.name}</p>}</For>
        </Match>
      </Switch>
    </div>
  );
}

const App: Component = () => {
  return (
    <QueryClientProvider client={queryClient}>
      <Example />
    </QueryClientProvider>
  );
};

export default App;
