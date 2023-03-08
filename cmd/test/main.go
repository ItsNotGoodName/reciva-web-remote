package main

import (
	"fmt"
	"log"

	"github.com/ItsNotGoodName/go-upnpsub"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/background"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/hub"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/pubsub"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/radio"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/upnp"
	"github.com/ItsNotGoodName/reciva-web-remote/pkg/interrupt"
)

func main() {
	ctx := interrupt.Context()

	// Dependencies
	cp := upnpsub.NewControlPoint()
	h := hub.New()
	done := background.Run(ctx, []background.Background{h, upnp.NewBackgroundControlPoint(cp)})

	sub, unsub := pubsub.DefaultPub.Subscribe([]string{pubsub.StateTopic})
	go func() {
		for msg := range sub {
			fmt.Println(msg.Data)
			unsub()
		}
	}()

	// Discover radios
	if err := radio.Discover(ctx, h, cp); err != nil {
		log.Fatalln("failed to discover radios:", err)
	}

	<-done
}
