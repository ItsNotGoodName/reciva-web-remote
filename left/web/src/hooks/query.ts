import { computed, Ref } from "vue"
import { useQuery } from "vue-query";

import { API_URL } from "../constants"
import { KEY_RADIOS, KEY_PRESETS, KEY_PRESET } from "./key";

export function useRadiosQuery() {
  return useQuery(KEY_RADIOS, () =>
    fetch(API_URL + "/api/radios")
      .then((response) => response.json())
      .then((json: APIResponse<State[]>) => {
        if (!json.ok) {
          throw new Error(json.error.message);
        }
        return json.data;
      })
  );
}

export function usePresetsQuery() {
  return useQuery(KEY_PRESETS, () =>
    fetch(API_URL + "/api/presets")
      .then((response) => response.json())
      .then((json: APIResponse<Preset[]>) => {
        if (!json.ok) {
          throw new Error(json.error.message);
        }
        return json.data;
      })
  );
}

export function usePresetQuery(url: Ref<string>) {
  return useQuery([KEY_PRESET, url], () =>
    fetch(API_URL + "/api/preset?url=" + url.value)
      .then((response) => response.json())
      .then((json: APIResponse<Preset>) => {
        if (!json.ok) {
          throw new Error(json.error.message);
        }
        return json.data;
      })
    , { enabled: computed(() => !!url.value) });
}
