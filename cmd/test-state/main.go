package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ItsNotGoodName/go-upnpsub"
	"github.com/ItsNotGoodName/reciva-web-remote/core/radio"
	"github.com/ItsNotGoodName/reciva-web-remote/core/state"
	"github.com/ItsNotGoodName/reciva-web-remote/core/upnp"
)

func main() {
	controlPoint := upnpsub.NewControlPoint()
	go upnpsub.ListenAndServe("", controlPoint)

	pub := state.NewPub()

	sub, err := pub.Subscribe(8)
	if err != nil {
		log.Fatal("could not sub:", err)
	}

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

	runService := radio.NewRunService(pub)
	createService := radio.NewCreateService(controlPoint, runService)

	clients, _, err := upnp.Discover()
	if err != nil {
		log.Fatal("failed to discover radios:", err)
	}
	if len(clients) == 0 {
		log.Fatal("no radios found")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	radio, err := createService.Create(ctx, clients[0])
	if err != nil {
		log.Fatal("failed to create radio:", err)
	}

	time.Sleep(10 * time.Second)

	cancel()
	<-radio.Done()
}
