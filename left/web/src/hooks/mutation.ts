import { useMutation, useQueryClient } from "vue-query";

import { API_URL } from "../constants"
import { KEY_RADIOS } from "./key";

export function useRadioMutation() {
  return useMutation((req: RadioMutation) => {
    const { uuid, ...jreq } = req
    return fetch(API_URL + "/api/radio/" + uuid, { body: JSON.stringify(jreq), method: "PATCH" })
      .then((res) => res.json())
      .then((json: APIResponse<void>) => {
        if (!json.ok) {
          throw new Error(json.error.message);
        }
      })
  })
};


export function useRadioSubscriptionMutation() {
  return useMutation((uuid: string) => fetch(API_URL + "/api/radio/" + uuid + "/subscription", { method: "POST" })
    .then((res) => res.json())
    .then((json: APIResponse<void>) => {
      if (!json.ok) {
        throw new Error(json.error.message);
      }
    }));
};

export function useRadioVolumeMutation() {
  return useMutation((uuid: string) => fetch(API_URL + "/api/radio/" + uuid + "/volume", { method: "POST" })
    .then((res) => res.json())
    .then((json: APIResponse<void>) => {
      if (!json.ok) {
        throw new Error(json.error.message);
      }
    }));
};

export function useRadiosDiscoverMutation() {
  const queryClient = useQueryClient()
  return useMutation((_: RadiosDiscoverMutation) => fetch(API_URL + "/api/radios", { method: "POST" })
    .then((res) => res.json())
    .then((json: APIResponse<number>) => {
      if (!json.ok) {
        throw new Error(json.error.message);
      }
      return json.data
    }), {
    onSuccess: (_: number) => {
      queryClient.invalidateQueries(KEY_RADIOS)
    }
  })
};