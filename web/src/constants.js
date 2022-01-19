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

export const MAX_MESSAGES = 100;

export const MESSAGE_INFO = "info";
export const MESSAGE_ERROR = "error";
export const MESSAGE_WARNING = "warning";
export const MESSAGE_SUCCESS = "success";
