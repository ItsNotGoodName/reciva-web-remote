package ws

import (
	"context"
	"log"
	"time"

	"github.com/ItsNotGoodName/reciva-web-remote/internal/hub"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/model"
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
	msgC, unsub := pubsub.DefaultPub.SubscribeWithBuffer(lastTopics, 10)
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
		if msg.Topic == pubsub.DiscoverTopic {
			return msg.Topic, msg.Data.(pubsub.DiscoverMessage).Discovering
		}
		if msg.Topic == pubsub.StateTopic {
			data := msg.Data.(pubsub.StateMessage)
			if stateCommand.UUID == "" || !(stateCommand.UUID == "*" || stateCommand.UUID == data.State.UUID) {
				return "", nil
			}
			if stateCommand.Partial && !data.Changed.Is(state.ChangedAll) {
				return msg.Topic, state.GetPartial(&data.State, data.Changed)
			}
			return msg.Topic, data.State
		}
		if msg.Topic == pubsub.StaleTopic {
			return msg.Topic, msg.Data
		}
		if msg.Topic == pubsub.StateHookStaleTopic {
			return pubsub.StaleTopic, model.StalePresets
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
				unsub = pubsub.DefaultPub.Resubscribe(topics, msgC, unsub)

				// Sync
				for _, t := range topics {
					switch t {
					case pubsub.DiscoverTopic:
						if pubsub.DiscoverTopic.In(lastTopics) {
							continue
						}

						if !write(pubsub.DiscoverTopic, d.Discovering()) {
							return
						}
					}
				}

				lastTopics = topics
			}

			// State Command
			if command.State != nil {
				stateCommand = *command.State

				if pubsub.StateTopic.In(lastTopics) {
					// Sync
					if stateCommand.UUID == "*" {
						for _, r := range h.List() {
							if s, err := radio.GetState(ctx, r); err != nil {
								log.Println("ws.Handle:", err)
							} else {
								if !write(pubsub.StateTopic, s) {
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
								if !write(pubsub.StateTopic, s) {
									return
								}
							}
						}
					}
				}
			}
		case msg := <-msgC:
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
