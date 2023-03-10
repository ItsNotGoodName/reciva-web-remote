/* eslint-disable @typescript-eslint/no-unsafe-assignment */
export const API_PATH: string = import.meta.env.VITE_API_PATH;
export const API_URL: string = import.meta.env.VITE_API_URL
  ? (import.meta.env.VITE_API_URL as string) + API_PATH
  : API_PATH;
export const WS_URL: string = import.meta.env.VITE_WS_URL
  ? import.meta.env.VITE_WS_URL
  : (() => {
      if (window.location.protocol == "http:") {
        return "ws://" + window.location.host + API_PATH + "/ws";
      }
      return "wss://" + window.location.host + API_PATH + "/ws";
    })();
export const GITHUB_URL: string = import.meta.env.VITE_GITHUB_URL;

export const PAGE_HOME = "";
export const PAGE_EDIT = "edit";
