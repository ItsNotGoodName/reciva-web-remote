import { useQuery } from "vue-query";

import { API_URL } from "../constants"

export function useRadiosQuery() {
  return useQuery("radios", () =>
    fetch(API_URL + "/api/radios")
      .then((response) => response.json())
      .then((json: APIResponse<Radio[]>) => {
        if (!json.ok) {
          throw new Error(json.error.message);
        }
        return json.data;
      })
  );
}
