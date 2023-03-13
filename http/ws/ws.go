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
	read := make(chan Command)
	go handleRead(ctx, cancel, conn, read)
	ticker := time.NewTicker(wsPingPeriod)
	sub, unsub := pubsub.DefaultPub.Subscribe([]pubsub.Topic{})
	defer func() {
		cancel()
		ticker.Stop()
		unsub()
	}()

	write := func(topic pubsub.Topic, data any) bool {
		conn.SetWriteDeadline(time.Now().Add(wsWriteWait))

		if err := conn.WriteJSON(Event{Topic: topic, Data: data}); err != nil {
			log.Printf("ws.Handle: could not write to %s: %s", conn.RemoteAddr(), err)
			return false
		}

		return true
	}

	stateCommand := CommandState{}

	for {
		select {
		case <-ctx.Done():
			// Send close message and end
			conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		case command, ok := <-read:
			if !ok {
				return
			}

			if command.State != nil {
				stateCommand = *command.State
			}

			if command.Subscribe != nil {
				topics := filterValidTopics(command.Subscribe.Topics)
				unsub()
				// Subscribe
				sub, unsub = pubsub.DefaultPub.Subscribe(topics)

				// Sync
				for _, t := range topics {
					if t == pubsub.StateTopic {
						if stateCommand.UUID == "" {
							for _, r := range h.List() {
								if s, err := radio.GetState(ctx, r); err != nil {
									log.Println("ws.Handle:", err)
								} else {
									if !write(pubsub.StateTopic, s) {
										return
									}
								}

							}
						} else {
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
					} else if t == pubsub.DiscoverTopic {
						if !write(pubsub.DiscoverTopic, d.Discovering()) {
							return
						}
					}
				}
			}
		case msg := <-sub:
			// Parse data from pubsub message
			data := func() any {
				if msg.Topic == pubsub.DiscoverTopic {
					return msg.Data.(pubsub.DiscoverMessage).Discovering
				}
				if msg.Topic == pubsub.StateTopic {
					data := msg.Data.(pubsub.StateMessage)
					if stateCommand.Partial && !data.Changed.Is(state.ChangedAll) {
						return state.GetPartial(&data.State, data.Changed)
					}
					return data.State
				}
				return nil
			}()
			if data == nil {
				continue
			}

			// Send event
			if !write(msg.Topic, data) {
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
