package api

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

func NewRadioWS(conn *websocket.Conn, h *radio.Hub) *RadioWS {
	hc := make(chan radio.State, 2)
	return &RadioWS{
		h:        h,
		conn:     conn,
		hubChan:  &hc,
		readChan: make(chan *radio.State),
		sendChan: make(chan *radio.State),
	}
}

func (rs *RadioWS) Start(uuid string) {
	go rs.balancer(uuid)

	// Register with hub
	rs.h.AddClient(rs.hubChan)

	// Start read handler
	go rs.handleRead(uuid)

	// Start write handler
	go rs.handleWrite()
}

// balancer handles receiving state from radio.Hub and handleRead and sending it to handleWrite.
// It exits when hub closes it's channel.
func (rs *RadioWS) balancer(uuid string) {
	var toSend *radio.State

	// hubChan sends incremental changes while readChan sends full state changes.
	// readChan is used to change the uuid.
	for {
		if toSend == nil {
			select {
			case state, ok := <-*rs.hubChan:
				if !ok {
					close(rs.sendChan)
					return
				}

				if state.UUID != uuid {
					continue
				}

				toSend = &state
			case state := <-rs.readChan:
				uuid = state.UUID
				toSend = state
			}
		}

		select {
		case rs.sendChan <- toSend:
			toSend = nil
		case state, ok := <-*rs.hubChan:
			if !ok {
				close(rs.sendChan)
				return
			}

			if state.UUID != uuid {
				continue
			}

			toSend.Merge(&state)
		case state := <-rs.readChan:
			uuid = state.UUID
			toSend = state
		}
	}
}

func (rs *RadioWS) handleWrite() {
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

			// Send state or end on error
			if err := rs.conn.WriteJSON(state); err != nil {
				log.Println("RadioWS.handleWrite(ERROR):", err)
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

func (rs *RadioWS) handleRead(uuid string) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		rs.h.RemoveClient(rs.hubChan)
		rs.conn.Close()
		cancel()
	}()

	// Send initial state if uuid is set
	if uuid != "" {
		state, err := rs.h.GetRadioState(ctx, uuid)

		if err != nil {
			log.Println("RadioWS.handleRead(ERROR):", err)
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
				log.Println("RadioWS.handleRead(ERROR):", err)
			}
			return
		}

		// Parse uuid and get state from uuid
		state, err := rs.h.GetRadioState(ctx, string(msg))

		// End connection if unable to get state
		if err != nil {
			return
		}

		// Send state to client
		rs.readChan <- state
	}
}
