package api

import (
	"context"
	"log"
	"sync"
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

type radioWS struct {
	readChan  chan *radio.State
	hubChan   *chan radio.State
	a         *API
	sendChan  chan *radio.State
	conn      *websocket.Conn
	uuid      string
	uuidMutex sync.Mutex
}

func newRadioWS(conn *websocket.Conn, a *API, uuid string) *radioWS {
	hc := make(chan radio.State)
	return &radioWS{
		readChan:  make(chan *radio.State),
		hubChan:   &hc,
		a:         a,
		sendChan:  make(chan *radio.State),
		conn:      conn,
		uuid:      uuid,
		uuidMutex: sync.Mutex{},
	}
}

func (rs *radioWS) start() {
	// Aggregate readChan and hubChan into sendChan
	go func() {
		for {
			select {
			case state, ok := <-*rs.hubChan:
				if !ok {
					close(rs.sendChan)
					return
				}
				rs.sendChan <- &state
			case state, ok := <-rs.readChan:
				if !ok {
					close(rs.sendChan)
					return
				}
				rs.sendChan <- state
			}
		}
	}()

	// Register with hub
	rs.a.h.Register <- rs.hubChan

	// Start read handler
	go rs.handleRead()

	// Start write handler
	go rs.handleWrite()
}

func (rs *radioWS) handleWrite() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		rs.conn.Close()
	}()

	for {
		select {
		case state, ok := <-rs.sendChan:
			// Set 10 second deadline
			rs.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

			if !ok {
				// Tell client the connection is done
				rs.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// Check if UUID matches
			rs.uuidMutex.Lock()
			if state.UUID != rs.uuid {
				rs.uuidMutex.Unlock()
				continue
			}
			rs.uuidMutex.Unlock()

			// Send state or end on error
			if err := rs.conn.WriteJSON(state); err != nil {
				log.Println("radioWS.handleWrite:", err)
				return
			}
		case <-ticker.C:
			rs.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := rs.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (rs *radioWS) handleRead() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		rs.a.h.Unregister <- rs.hubChan
		rs.conn.Close()
		cancel()
	}()

	// Send initial state if uuid is set
	if rs.uuid != "" {
		rs.uuidMutex.Lock()
		uuid := rs.uuid
		rs.uuidMutex.Unlock()

		state, ok := rs.a.GetRadioState(ctx, uuid)

		if !ok {
			log.Printf("radioWS.handleRead(ERROR): GetRadioState did not find state with radio uuid %s", uuid)
			return
		}

		rs.readChan <- state
	}

	rs.conn.SetReadLimit(maxMessageSize)
	rs.conn.SetReadDeadline(time.Now().Add(pongWait))
	rs.conn.SetPongHandler(func(string) error { rs.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, msg, err := rs.conn.ReadMessage()
		if err != nil {
			if !websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("radioWS.handleRead:", err)
			}
			return
		}

		// Parse uuid and get state from uuid
		uuid := string(msg)
		state, ok := rs.a.GetRadioState(ctx, uuid)

		// End connection if unable to get state
		if !ok {
			return
		}

		// Update current uuid
		rs.uuidMutex.Lock()
		rs.uuid = uuid
		rs.uuidMutex.Unlock()

		// Send state to client
		rs.readChan <- state
	}
}
