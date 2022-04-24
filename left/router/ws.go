package router

import (
	"net/http"

	"github.com/ItsNotGoodName/reciva-web-remote/core/radio"
	"github.com/gorilla/websocket"
)

func GetWS(hub radio.HubService, upgrader websocket.Upgrader) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(rw, r, nil)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		conn.Close()
	}
}
