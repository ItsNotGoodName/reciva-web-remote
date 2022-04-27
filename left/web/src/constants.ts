export const API_URL = import.meta.env.VITE_API_URL
  ? import.meta.env.VITE_API_URL
  : "";
export const WS_URL = import.meta.env.VITE_WS_URL
  ? import.meta.env.VITE_WS_URL
  : (() => {
    if (window.location.protocol == "http:") {
      return "ws://" + window.location.host;
    }
    return "wss://" + window.location.host;
  })();


export const PAGE_HOME = "";
export const PAGE_EDIT = "edit";

export const STATUS_CONNECTING = "Connecting";
export const STATUS_PLAYING = "Playing";
export const STATUS_STOPPED = "Stopped";