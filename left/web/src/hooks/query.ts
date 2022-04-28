import { useQuery } from "vue-query";

import { API_URL } from "../constants"
import { KEY_SLIM_RADIOS } from "./key";

export function useSlimRadiosQuery() {
  return useQuery(KEY_SLIM_RADIOS, () =>
    fetch(API_URL + "/api/radios/slim")
      .then((response) => response.json())
      .then((json: APIResponse<SlimRadio[]>) => {
        if (!json.ok) {
          throw new Error(json.error.message);
        }
        return json.data;
      })
  );
}
