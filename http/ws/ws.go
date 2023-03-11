package ws

import (
	"context"
	"fmt"
	"log"

	"github.com/ItsNotGoodName/reciva-web-remote/internal/hub"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/pubsub"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/radio"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/state"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type CommandSubscribe struct {
	Topics []pubsub.Topic `json:"topics" validate:"required"`
}

type CommandState struct {
	UUID    string `json:"uuid" validate:"required"`
	Partial bool   `json:"partial" validate:"required"`
}

type Command struct {
	Subscribe *CommandSubscribe `json:"subscribe"`
	State     *CommandState     `json:"state"`
}

type Event struct {
	Topic pubsub.Topic `json:"topic" validate:"required"`
	Data  any          `json:"data" validate:"required"`
}

func validateTopic(topic pubsub.Topic) (pubsub.Topic, error) {
	switch pubsub.Topic(topic) {
	case pubsub.DiscoverTopic:
		return pubsub.DiscoverTopic, nil
	case pubsub.StateTopic:
		return pubsub.StateTopic, nil
	// case pubsub.ErrorTopic:
	// 	return pubsub.ErrorTopic, nil
	default:
		return "", fmt.Errorf("invalid topic")
	}
}

func filterValidTopics(topics []pubsub.Topic) []pubsub.Topic {
	pubsubTopics := []pubsub.Topic{}
	for _, topic := range topics {
		pubsubTopic, err := validateTopic(topic)
		if err != nil {
			log.Println("http.wsParseTopics: invalid topic:", topic)
			continue
		}

		pubsubTopics = append(pubsubTopics, pubsubTopic)
	}

	return pubsubTopics
}

func handleRead(conn *websocket.Conn) (<-chan Command, context.CancelFunc) {
	read := make(chan Command)
	readCtx, readCancel := context.WithCancel(context.Background())
	go func() {
		var command Command
		for {
			if err := wsjson.Read(readCtx, conn, &command); err != nil {
				log.Println("http.wsHandleRead: could not read:", err)
				close(read)
				return
			}

			select {
			case read <- command:
			case <-readCtx.Done():
				close(read)
				return
			}
		}
	}()
	return read, readCancel
}

func writeFunc(ctx context.Context, conn *websocket.Conn) func(topic pubsub.Topic, data any) {
	return func(topic pubsub.Topic, data any) {
		if err := wsjson.Write(ctx, conn, Event{Topic: topic, Data: data}); err != nil {
			log.Println("wsWrite:", err)
		}
	}
}

func Handle(ctx context.Context, conn *websocket.Conn, h *hub.Hub) {
	sub, unsub := pubsub.DefaultPub.Subscribe([]pubsub.Topic{})
	read, readCancel := handleRead(conn)
	defer func() {
		unsub()
		readCancel()
	}()
	stateCommand := CommandState{}
	write := writeFunc(ctx, conn)

	for {
		select {
		case command, ok := <-read:
			if !ok {
				return
			}

			if command.State != nil {
				stateCommand = *command.State
			}

			if command.Subscribe != nil {
				unsub()
				topics := filterValidTopics(command.Subscribe.Topics)
				sub, unsub = pubsub.DefaultPub.Subscribe(topics)
				for _, t := range topics {
					if t == pubsub.StateTopic {
						if stateCommand.UUID == "" {
							for _, r := range h.List() {
								if s, err := radio.GetState(ctx, r); err != nil {
									log.Println("ws.Handle:", err)
								} else {
									write(pubsub.StateTopic, s)
								}

							}
						} else {
							if r, err := h.Get(stateCommand.UUID); err != nil {
								log.Println("ws.Handle:", err)
							} else {
								if s, err := radio.GetState(ctx, r); err != nil {
									log.Println("ws.Handle:", err)
								} else {
									write(pubsub.StateTopic, s)
								}
							}

						}
					}
				}
			}
		case msg := <-sub:
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

			write(msg.Topic, data)
		case <-ctx.Done():
			return
		}
	}
}
