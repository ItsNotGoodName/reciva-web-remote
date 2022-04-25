package ws

import (
	"context"
	"log"
	"time"

	"github.com/ItsNotGoodName/reciva-web-remote/core/app"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

func HandleWrite(ctx context.Context, cancel context.CancelFunc, conn *websocket.Conn) chan<- app.Command {
	writeC := make(chan app.Command)
	go func() {
		ticker := time.NewTicker(pingPeriod)
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
				conn.SetWriteDeadline(time.Now().Add(writeWait))

				// Send msg or end on error
				if err := conn.WriteJSON(msg); err != nil {
					log.Printf("ws.handleWrite: could not write to %s: %s", conn.RemoteAddr(), err)
					return
				}
			case <-ticker.C:
				conn.SetWriteDeadline(time.Now().Add(writeWait))

				// Send ping or end on error
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					log.Printf("ws.handleWrite: could not write ping %s: %s", conn.RemoteAddr(), err)
					return
				}
			}
		}
	}()
	return writeC
}

func HandleRead(ctx context.Context, cancel context.CancelFunc, conn *websocket.Conn) <-chan app.Command {
	readC := make(chan app.Command)
	go func() {
		defer cancel()

		conn.SetReadLimit(maxMessageSize)
		conn.SetReadDeadline(time.Now().Add(pongWait))
		conn.SetPongHandler(func(string) error { conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

		for {
			// Read message or end on error
			var msg app.Command
			err := conn.ReadJSON(&msg)
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("ws.HandleRead: could not read from %s: %s", conn.RemoteAddr(), err)
				}
				return
			}

			// Send message
			select {
			case readC <- msg:
			case <-ctx.Done():
				return
			}
		}
	}()
	return readC
}
