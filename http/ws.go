package http

import (
	"context"
	"fmt"
	"log"

	"github.com/ItsNotGoodName/reciva-web-remote/internal/pubsub"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/state"
	"github.com/labstack/echo/v4"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type WSCommandSubscribe struct {
	Topics []string `json:"topics"`
}

type WSCommandState struct {
	UUID    string `json:"uuid"`
	Partial bool   `json:"partial"`
}

type WSCommand struct {
	Subscribe *WSCommandSubscribe `json:"subscribe"`
	State     *WSCommandState     `json:"state"`
}

type WSEvent struct {
	Topic pubsub.Topic `json:"topic"`
	Data  any          `json:"data"`
}

func wsParseTopic(topic string) (pubsub.Topic, error) {
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

func wsParseTopics(topics []string) []pubsub.Topic {
	pubsubTopics := []pubsub.Topic{}
	for _, topic := range topics {
		pubsubTopic, err := wsParseTopic(topic)
		if err != nil {
			log.Println("http.wsParseTopics: invalid topic:", topic)
			continue
		}

		pubsubTopics = append(pubsubTopics, pubsubTopic)
	}

	return pubsubTopics
}

func wsHandleRead(conn *websocket.Conn) (<-chan WSCommand, context.CancelFunc) {
	read := make(chan WSCommand)
	readCtx, readCancel := context.WithCancel(context.Background())
	go func() {
		var command WSCommand
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

func wsHandle(ctx context.Context, conn *websocket.Conn) {
	sub, unsub := pubsub.DefaultPub.Subscribe([]pubsub.Topic{})
	read, readCancel := wsHandleRead(conn)
	defer func() {
		unsub()
		readCancel()
	}()
	stateCommand := WSCommandState{}

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
				sub, unsub = pubsub.DefaultPub.Subscribe(wsParseTopics(command.Subscribe.Topics))
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

			if err := wsjson.Write(ctx, conn, WSEvent{Topic: msg.Topic, Data: data}); err != nil {
				log.Println("http.handleWS:", err)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (a API) WS(c echo.Context) error {
	conn, err := websocket.Accept(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer conn.Close(websocket.StatusInternalError, "the sky is falling")

	wsHandle(c.Request().Context(), conn)

	return conn.Close(websocket.StatusNormalClosure, "")
}
