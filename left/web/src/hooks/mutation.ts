import { useMutation } from "vue-query";

import { API_URL } from "../constants"

export function useRadioMutation(uuid: string) {
  return useMutation((req: RadioMutation) => fetch(API_URL + "/api/radio/" + uuid, { body: JSON.stringify(req), method: "PATCH" })
    .then((res) => res.json())
    .then((json: APIResponse<void>) => {
      if (!json.ok) {
        throw new Error(json.error.message);
      }
    }))
};