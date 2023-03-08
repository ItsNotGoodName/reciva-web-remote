package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ItsNotGoodName/go-upnpsub"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/background"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/hub"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/middleware"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/pubsub"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/radio"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/store"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/upnp"
	"github.com/ItsNotGoodName/reciva-web-remote/pkg/interrupt"
)

func main() {
	ctx := interrupt.Context()
	hub := hub.New()
	controlPoint := upnpsub.NewControlPoint()
	store := middleware.NewStore(store.Must(store.NewFile("/tmp/tmp.json")))
	stateHook := middleware.NewStateHook(store)
	discoverer := radio.NewDiscoverer(hub, controlPoint, stateHook)
	done := background.Run(ctx, []background.Background{hub, upnp.NewBackgroundControlPoint(controlPoint), discoverer})
	sub, _ := pubsub.DefaultPub.Subscribe([]pubsub.Topic{pubsub.StateTopic})

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

	go func() {
		time.Sleep(10 * time.Second)
		presets, _ := store.ListPresets(ctx)
		presets[0].TitleNew = "testttt"
		store.UpdatePreset(ctx, &presets[0])
	}()

	<-done
}
