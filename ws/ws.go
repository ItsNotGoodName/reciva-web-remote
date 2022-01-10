package ws

import (
	"context"
	"log"
	"time"

	"github.com/ItsNotGoodName/reciva-web-remote/pkg/radio"
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

// Handle handles radio websocket connection.
func Handle(conn *websocket.Conn, hub *radio.Hub, uuid string) {
	log.Println("ws.Handle(INFO): new connection from", conn.RemoteAddr())

	// Start read goroutine
	readC := make(chan string, 1)
	if uuid != "" {
		readC <- uuid
	}
	go handleRead(conn, readC)

	// Start write goroutine
	writeC := make(chan *radio.State, 1)
	writeDoneC := make(chan struct{})
	go handleWrite(conn, writeC, writeDoneC)

	// Start buffering write goroutine
	mergeWriteC := make(chan *radio.State)
	go mergeWrite(mergeWriteC, writeC)

	// Subscribe to hub
	subC := make(chan radio.State, 5)
	sub := radio.NewSub(subC)
	hub.Pub.Subscribe(sub)

	defer func() {
		log.Println("ws.Handle(INFO): closing connection from", conn.RemoteAddr())
		conn.Close()
		hub.Pub.Unsubscribe(sub)
		close(mergeWriteC)
		log.Println("ws.Handle(INFO): connection closed from", conn.RemoteAddr())
	}()

	for {
		log.Println("ws.Handle(INFO): looping in handle from", conn.RemoteAddr())
		select {
		case state, ok := <-subC:
			if !ok {
				return
			}

			if state.UUID != uuid {
				continue
			}

			mergeWriteC <- &state
		case newUUID, ok := <-readC:
			log.Println("ws.Handle(INFO): readC trigger from", conn.RemoteAddr(), newUUID, ok)
			if !ok {
				return
			}

			uuid = newUUID
			state, err := hub.GetRadioState(context.Background(), uuid)
			if err != nil {
				return
			}

			mergeWriteC <- state
		case <-writeDoneC:
			log.Println("ws.Handle(INFO): writeDoneC trigger from", conn.RemoteAddr())
			return
		}
	}

}

// mergeWrite tries to write to outC from inC. If another inC is received before a write to outC, it will merge them before trying to write.
func mergeWrite(inC <-chan *radio.State, outC chan<- *radio.State) {
	var toWrite *radio.State
	var ok bool
	for {
		if toWrite != nil {
			select {
			case outC <- toWrite:
				toWrite = nil
			case state, ok := <-inC:
				if !ok {
					close(outC)
					log.Println("ws.mergeWrite(INFO): closed mergeWrite")
					return
				}
				toWrite.Merge(state)
			}
		} else {
			toWrite, ok = <-inC
			if !ok {
				close(outC)
				log.Println("ws.mergeWrite(INFO): closed mergeWrite")
				return
			}
		}
	}
}

// handleWrite writes state from writeC.
func handleWrite(conn *websocket.Conn, writeC <-chan *radio.State, doneC chan<- struct{}) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		close(doneC)
		log.Println("ws.handleWrite(INFO): closed handleWrite from", conn.RemoteAddr())
	}()

	for {
		select {
		case state, ok := <-writeC:
			// Set 10 second deadline
			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

			if !ok {
				// Tell client the connection is done
				conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// Send state or end on error
			if err := conn.WriteJSON(state); err != nil {
				log.Printf("RadioWS.handleWrite(ERROR): could not write to %s: %s", conn.RemoteAddr(), err)
				return
			}
		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleRead reads uuid fron client and sends it to readC.
func handleRead(conn *websocket.Conn, readC chan string) {
	conn.SetReadLimit(maxMessageSize)
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error { conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("ws.handleRead(ERROR): could not read from %s: %s", conn.RemoteAddr(), err)
			}
			close(readC)
			log.Println("ws.handleRead(INFO): closed handleRead from", conn.RemoteAddr())
			return
		}

		uuid := string(msg)
		select {
		case readC <- uuid:
		default:
			select {
			case <-readC:
			default:
			}
			readC <- uuid
		}
	}
}
