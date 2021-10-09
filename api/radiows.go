package api

import (
	"log"
	"sync"
	"time"

	"github.com/ItsNotGoodName/reciva-web-remote/pkg/radio"
	"github.com/gorilla/websocket"
)

type radioWS struct {
	c         *radio.HubClient
	a         *API
	conn      *websocket.Conn
	h         *radio.Hub
	uuid      string
	uuidMutex sync.Mutex
}

func newRadioWS(conn *websocket.Conn, a *API, h *radio.Hub, uuid string) *radioWS {
	return &radioWS{
		&radio.HubClient{Send: make(chan radio.State)},
		a,
		conn,
		h,
		uuid,
		sync.Mutex{},
	}
}

func (rs *radioWS) start() {
	// Start write handler
	go rs.handleWrite()

	// Send init state if uuid not empty and close ws if send state failed
	if rs.uuid != "" {
		state, ok := rs.a.GetRadioState(rs.uuid)

		if !ok {
			log.Printf("RadioWS.start(ERROR): could not send state with uuid %s", rs.uuid)
			close(rs.c.Send)
			return
		}

		rs.c.Send <- *state
	}

	// Start read handler
	go rs.handleRead()

	// Register client with hub
	rs.h.Register <- rs.c
}

func (rs *radioWS) end() {
	rs.h.Unregister <- rs.c
	rs.conn.Close()
}

func (rs *radioWS) handleWrite() {
	defer rs.end()

	// Loop until channel is closed
	for state := range rs.c.Send {
		// Set 10 second deadline
		rs.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

		// Check if UUID matches
		rs.uuidMutex.Lock()
		if state.UUID != rs.uuid {
			continue
		}
		rs.uuidMutex.Unlock()

		// Send state or end on error
		if err := rs.conn.WriteJSON(state); err != nil {
			log.Println(err)
			return
		}
	}
	// Tell client the connection is done
	rs.conn.WriteMessage(websocket.CloseMessage, []byte{})
}

func (rs *radioWS) handleRead() {
	defer rs.end()

	rs.conn.SetReadLimit(512)
	for {
		_, msg, err := rs.conn.ReadMessage()
		if err != nil {
			if !websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println(err)
			}
			return
		}

		// Parse uuid and get state from uuid
		uuid := string(msg)
		state, ok := rs.a.GetRadioState(uuid)

		// End connection if unable to get state
		if !ok {
			return
		}

		// Update current uuid
		rs.uuidMutex.Lock()
		rs.uuid = uuid
		rs.uuidMutex.Unlock()

		// Send state to client
		rs.c.Send <- *state
	}
}
