package ws

import (
	"context"
	"log"
	"time"

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

// HandleWrite writes json to websocket connection.
func HandleWrite(ctx context.Context, conn *websocket.Conn) <-chan interface{} {
	writeC := make(chan interface{})
	go func() {
		ticker := time.NewTicker(pingPeriod)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			case msg := <-writeC:
				// Set 10 second deadline
				conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

				// Send msg or end on error
				if err := conn.WriteJSON(msg); err != nil {
					log.Printf("ws.handleWrite: could not write to %s: %s", conn.RemoteAddr(), err)
					return
				}
			case <-ticker.C:
				conn.SetWriteDeadline(time.Now().Add(writeWait))
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					return
				}
			}
		}
	}()
	return writeC
}

// HandleRead reads string messages from websocket connection.
func HandleRead(ctx context.Context, conn *websocket.Conn) <-chan interface{} {
	readC := make(chan interface{})
	go func() {
		conn.SetReadLimit(maxMessageSize)
		conn.SetReadDeadline(time.Now().Add(pongWait))
		conn.SetPongHandler(func(string) error { conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

		for {
			var msg interface{}
			err := conn.ReadJSON(msg)
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("ws.handleRead: could not read from %s: %s", conn.RemoteAddr(), err)
				}
				close(readC)
				return
			}

			select {
			case readC <- msg:
			case <-ctx.Done():
			}
		}
	}()
	return readC
}
