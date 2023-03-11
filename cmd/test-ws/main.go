package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ItsNotGoodName/reciva-web-remote/http/ws"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/pubsub"
	"github.com/ItsNotGoodName/reciva-web-remote/pkg/interrupt"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func run(ctx context.Context) {

	c, _, err := websocket.Dial(ctx, "ws://localhost:8080/api/ws", nil)
	if err != nil {
		log.Fatalln("could not connect:", err)
	}
	defer c.Close(websocket.StatusInternalError, "the sky is falling")

	command := ws.Command{
		Subscribe: &ws.CommandSubscribe{Topics: []pubsub.Topic{pubsub.StateTopic, pubsub.DiscoverTopic}},
		State:     &ws.CommandState{Partial: true},
	}
	if err := wsjson.Write(ctx, c, &command); err != nil {
		log.Fatalln("could not write:", err)
	}

	for {
		messageType, data, err := c.Read(ctx)
		if err != nil {
			log.Fatalln("could not read:", err)
		}

		fmt.Println(messageType, string(data))
	}
}

func main() {
	ctx := interrupt.Context()

	for i := 0; i < 1000; i++ {
		go run(ctx)
	}

	<-ctx.Done()
}
