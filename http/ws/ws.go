package ws

import (
	"context"
	"log"
	"time"

	"github.com/ItsNotGoodName/reciva-web-remote/internal/hub"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/pubsub"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/radio"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/state"
	"github.com/gorilla/websocket"
)

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

func handleRead(ctx context.Context, cancel context.CancelFunc, conn *websocket.Conn, read chan<- Command) {
	defer cancel()

	conn.SetReadLimit(wsMaxMessageSize)
	conn.SetReadDeadline(time.Now().Add(wsPongWait))
	conn.SetPongHandler(func(string) error { conn.SetReadDeadline(time.Now().Add(wsPongWait)); return nil })

	for {
		// Read command or end on error
		var command Command
		err := conn.ReadJSON(&command)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("ws.handleRead: could not read from %s: %s", conn.RemoteAddr(), err)
			}
			return
		}

		// Send command to handler
		select {
		case read <- command:
		case <-ctx.Done():
			return
		}
	}
}

func Handle(ctx context.Context, conn *websocket.Conn, h *hub.Hub, d *radio.Discoverer) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	readC := make(chan Command)
	go handleRead(ctx, cancel, conn, readC)
	ticker := time.NewTicker(wsPingPeriod)
	defer ticker.Stop()

	lastTopics := []pubsub.Topic{}
	sub, unsub := pubsub.DefaultPub.SubscribeWithBuffer(lastTopics, 25)
	defer func() { unsub() }()

	write := func(topic pubsub.Topic, data any) bool {
		conn.SetWriteDeadline(time.Now().Add(wsWriteWait))

		if err := conn.WriteJSON(Event{Topic: topic, Data: data}); err != nil {
			log.Printf("ws.Handle: could not write to %s: %s", conn.RemoteAddr(), err)
			return false
		}

		return true
	}

	stateCommand := CommandState{}

	parse := func(msg *pubsub.Message) (pubsub.Topic, any) {
		if data, ok := pubsub.ParseDiscover(msg); ok {
			return msg.Topic, data
		}
		if data, ok := pubsub.ParseState(msg); ok {
			if stateCommand.UUID == "" || !(stateCommand.UUID == "*" || stateCommand.UUID == data.State.UUID) {
				return "", nil
			}
			if stateCommand.Partial && !data.Changed.Is(state.ChangedAll) {
				return msg.Topic, state.GetPartial(&data.State, data.Changed)
			}
			return msg.Topic, data.State
		}
		if ok := pubsub.ParseStaleRadios(msg); ok {
			return pubsub.TopicStale, pubsub.StaleRadios
		}
		if _, ok := pubsub.ParseStaleStateHook(msg); ok {
			return pubsub.TopicStale, pubsub.StalePresets
		}
		return "", nil
	}

	for {
		select {
		case <-ctx.Done():
			// Send close message and end
			conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		case command, ok := <-readC:
			if !ok {
				return
			}

			// Subscribe command
			if command.Subscribe != nil {
				// Resubscribe
				topics := uniqueTopics(filterTopics(command.Subscribe.Topics))
				unsub = pubsub.DefaultPub.Resubscribe(topics, sub, unsub)

				// Sync
				for _, t := range topics {
					switch t {
					case pubsub.TopicDiscover:
						if pubsub.TopicDiscover.In(lastTopics) {
							continue
						}

						if !write(pubsub.TopicDiscover, d.Discovering()) {
							return
						}
					}
				}

				lastTopics = topics
			}

			// State Command
			if command.State != nil {
				stateCommand = *command.State

				if pubsub.TopicState.In(lastTopics) {
					// Sync
					if stateCommand.UUID == "*" {
						for _, r := range h.List() {
							if s, err := radio.GetState(ctx, r); err != nil {
								log.Println("ws.Handle:", err)
							} else {
								if !write(pubsub.TopicState, s) {
									return
								}
							}

						}
					} else if stateCommand.UUID != "" {
						if r, err := h.Get(stateCommand.UUID); err != nil {
							log.Println("ws.Handle:", err)
						} else {
							if s, err := radio.GetState(ctx, r); err != nil {
								log.Println("ws.Handle:", err)
							} else {
								if !write(pubsub.TopicState, s) {
									return
								}
							}
						}
					}
				}
			}
		case <-sub.Overflow:
			conn.WriteMessage(websocket.CloseTryAgainLater, []byte{})
			return
		case msg := <-sub.Message:
			// Parse data from pubsub message
			topic, data := parse(&msg)
			if data == nil {
				continue
			}

			// Send event
			if !write(topic, data) {
				return
			}
		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(wsWriteWait))

			// Send ping or end on error
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("ws.Handle: could not write ping %s: %s", conn.RemoteAddr(), err)
				return
			}
		}
	}

}
