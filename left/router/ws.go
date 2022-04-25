package router

import (
	"context"
	"net/http"

	"github.com/ItsNotGoodName/reciva-web-remote/core/app"
	"github.com/ItsNotGoodName/reciva-web-remote/left/ws"
	"github.com/gorilla/websocket"
)

func GetWS(upgrader websocket.Upgrader, handleWS func(*websocket.Conn)) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(rw, r, nil)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		go handleWS(conn)
	}
}

func HandleWS(application *app.App) func(conn *websocket.Conn) {
	return func(conn *websocket.Conn) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		application.Bus(
			ctx,
			ws.HandleRead(ctx, cancel, conn),
			ws.HandleWrite(ctx, cancel, conn),
		)
	}
}
