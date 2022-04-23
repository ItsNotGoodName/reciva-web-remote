package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/ItsNotGoodName/go-upnpsub"
	"github.com/ItsNotGoodName/reciva-web-remote/core/middleware"
	"github.com/ItsNotGoodName/reciva-web-remote/core/pubsub"
	"github.com/ItsNotGoodName/reciva-web-remote/core/radio"
	"github.com/ItsNotGoodName/reciva-web-remote/core/upnp"
	"github.com/ItsNotGoodName/reciva-web-remote/pkg/interrupt"
	"github.com/ItsNotGoodName/reciva-web-remote/pkg/sig"
	"github.com/ItsNotGoodName/reciva-web-remote/right/mock"
)

func main() {
	controlPoint := upnpsub.NewControlPoint()
	go upnpsub.ListenAndServe("", controlPoint)

	// Subscribe to all radios
	statePub := pubsub.NewStatePub()
	sub := statePub.Subscribe(8, "")
	go func() {
		for s := range sub.Channel() {
			j, err := json.MarshalIndent(s, "", "  ")
			if err != nil {
				log.Fatal("failed to marshal state:", err)
			}

			fmt.Println(string(j))
		}

		log.Println("channel closed")
	}()

	middlewarePub := sig.NewPub()

	runService := radio.NewRunService(
		statePub,
		middleware.NewPreset(middlewarePub, mock.NewPresetStore()),
		middlewarePub,
	)
	createService := radio.NewCreateService(controlPoint, runService)

	clients, _, err := upnp.Discover()
	if err != nil {
		log.Fatal("failed to discover radios:", err)
	}
	if len(clients) == 0 {
		log.Fatal("no radios found")
	}

	radio, err := createService.Create(interrupt.Context(), clients[0])
	if err != nil {
		log.Fatal("failed to create radio:", err)
	}

	<-radio.Done()
}
