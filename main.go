package main

import (
	"fmt"
	"log"

	"github.com/ItsNotGoodName/go-upnpsub"
	"github.com/ItsNotGoodName/reciva-web-remote/api"
	"github.com/ItsNotGoodName/reciva-web-remote/config"
	"github.com/ItsNotGoodName/reciva-web-remote/pkg/radio"
	"github.com/ItsNotGoodName/reciva-web-remote/routes"
	"github.com/ItsNotGoodName/reciva-web-remote/store"
)

func main() {
	// Create config
	cfg := config.NewConfig(config.WithFlag)

	// Create and start controlpoint
	cp := upnpsub.NewControlPointWithPort(cfg.CPort)
	go cp.Start()

	// Create radio hub
	h := radio.NewHub(cp)

	// Create store
	s, err := store.NewStore(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Create api
	a := api.NewAPI(h, s)

	// Create router
	r := NewRouter()

	// Create websocket upgrader
	u := NewUpgrader()

	// Add radio routes
	routes.AddRadioRoutes(r.Group(cfg.APIURI), a, u)

	// Add config routes
	routes.AddConfigRoutes(r.Group(cfg.APIURI), cfg)

	// Add stream routes
	routes.AddStreamRoutes(r.Group(cfg.APIURI), a)

	// Listen and serve
	log.Println("main: listening on port", cfg.Port)
	log.Fatal(r.Run(":" + fmt.Sprint(cfg.Port)))
}
