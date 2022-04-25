package api

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/ItsNotGoodName/reciva-web-remote/core/app"
	"github.com/gorilla/websocket"
)

func GetWS(upgrader *websocket.Upgrader, handleWS func(*websocket.Conn)) http.HandlerFunc {
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

		application.Bus(
			ctx,
			wsHandleRead(ctx, cancel, conn),
			wsHandleWrite(ctx, cancel, conn),
		)
	}
}

const (
	// Time allowed to write a message to the peer.
	wsWriteWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	wsPongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	wsPingPeriod = (wsPongWait * 9) / 10

	// Maximum message size allowed from peer.
	wsMaxMessageSize = 512
)

func wsHandleWrite(ctx context.Context, cancel context.CancelFunc, conn *websocket.Conn) chan<- app.Command {
	writeC := make(chan app.Command)
	go func() {
		ticker := time.NewTicker(wsPingPeriod)
		defer func() {
			cancel()
			ticker.Stop()
		}()

		for {
			select {
			case <-ctx.Done():
				// Send close message and end
				conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			case msg := <-writeC:
				conn.SetWriteDeadline(time.Now().Add(wsWriteWait))

				// Send msg or end on error
				if err := conn.WriteJSON(msg); err != nil {
					log.Printf("api.wsHandleWrite: could not write to %s: %s", conn.RemoteAddr(), err)
					return
				}
			case <-ticker.C:
				conn.SetWriteDeadline(time.Now().Add(wsWriteWait))

				// Send ping or end on error
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					log.Printf("api.wsHandleWrite: could not write ping %s: %s", conn.RemoteAddr(), err)
					return
				}
			}
		}
	}()
	return writeC
}

func wsHandleRead(ctx context.Context, cancel context.CancelFunc, conn *websocket.Conn) <-chan app.Command {
	readC := make(chan app.Command)
	go func() {
		defer cancel()

		conn.SetReadLimit(wsMaxMessageSize)
		conn.SetReadDeadline(time.Now().Add(wsPongWait))
		conn.SetPongHandler(func(string) error { conn.SetReadDeadline(time.Now().Add(wsPongWait)); return nil })

		for {
			// Read msg or end on error
			var msg app.Command
			err := conn.ReadJSON(&msg)
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("api.wsHandleRead: could not read from %s: %s", conn.RemoteAddr(), err)
				}
				return
			}

			// Send msg
			select {
			case readC <- msg:
			case <-ctx.Done():
				return
			}
		}
	}()
	return readC
}
