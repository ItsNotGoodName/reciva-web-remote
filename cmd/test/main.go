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
	"github.com/ItsNotGoodName/reciva-web-remote/right/mock"
)

func main() {
	// Dependencies
	statePub := pubsub.NewStatePub()
	middlewarePub := pubsub.NewSignalPub()
	runService := radio.NewRunService(
		statePub,
		middleware.NewPreset(middlewarePub, mock.NewPresetStore()),
		middlewarePub,
	)
	controlPoint := upnpsub.NewControlPoint()
	go upnpsub.ListenAndServe("", controlPoint)
	createService := radio.NewCreateService(controlPoint, runService)

	// Discover radios:w
	clients, _, err := upnp.Discover()
	if err != nil {
		log.Fatal("failed to discover radios:", err)
	}
	if len(clients) == 0 {
		log.Fatal("no radios found")
	}

	// Create radio
	radio, err := createService.Create(interrupt.Context(), clients[0])
	if err != nil {
		log.Fatal("failed to create radio:", err)
	}

	// Subscribe to state changes
	sub, unsub := statePub.Subscribe(radio.UUID)
	defer unsub()
	go func() {
		for s := range sub {
			j, err := json.MarshalIndent(s, "", "  ")
			if err != nil {
				log.Fatal("failed to marshal state:", err)
			}

			fmt.Println(string(j))
		}
	}()

	<-radio.Done()
}
