export const API_PATH = import.meta.env.VITE_API_PATH;
export const API_URL = import.meta.env.VITE_API_URL
  ? import.meta.env.VITE_API_URL + API_PATH
  : API_PATH;
export const WS_URL = import.meta.env.VITE_WS_URL
  ? import.meta.env.VITE_WS_URL
  : (() => {
      if (window.location.protocol == "http:") {
        return "ws://" + window.location.host + API_PATH + "/ws";
      }
      return "wss://" + window.location.host + API_PATH + "/ws";
    })();
export const GITHUB_URL = import.meta.env.VITE_GITHUB_URL;

export const PAGE_HOME = "";
export const PAGE_EDIT = "edit";
