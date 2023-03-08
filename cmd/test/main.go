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
	hub := hub.New()
	cp := upnpsub.NewControlPoint()
	discoverer := radio.NewDiscoverer(hub, cp, radio.DefaultStateHook{})
	done := background.Run(ctx, []background.Background{hub, upnp.NewBackgroundControlPoint(cp), discoverer})
	sub, _ := pubsub.DefaultPub.Subscribe([]pubsub.Topic{pubsub.DiscoverTopic})

	// Subscription
	go func() {
		for msg := range sub {
			fmt.Println(msg.Data)
		}
	}()

	// Discover radios
	go func() {
		if err := discoverer.Discover(ctx); err != nil {
			log.Fatalln("failed to discover radios:", err)
		}
	}()

	<-done
}
