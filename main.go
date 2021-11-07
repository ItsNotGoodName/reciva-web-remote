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
	cfg := config.NewConfig(config.WithFlag, config.WithFile)

	// Create and start controlpoint
	cp := upnpsub.NewControlPointWithPort(cfg.CPort)
	go cp.Start()

	// Create radio hub
	h := radio.NewHub(cp)

	// Create router
	r := NewRouter()

	// Create websocket upgrader
	u := NewUpgrader()

	// Add radio routes
	routes.AddRadioRoutes(r.Group(cfg.APIURI), h, u)

	// Add config routes
	routes.AddConfigRoutes(r.Group(cfg.APIURI), cfg)

	// Enable presets based on config
	if cfg.PresetsEnabled {
		// Create store
		if s, err := store.NewStore(cfg); err == nil {
			// Create preset api
			p := api.NewPresetAPI(s, h)
			// Add stream routes
			routes.AddStreamRoutes(r.Group(cfg.APIURI), p)
			// Add preset routes
			routes.AddPresetRoutes(r.Group(cfg.APIURI), p)
			// Add preset stream routes
			routes.AddPresetStreamRoutes(r, p)
		} else {
			log.Println("main:", err)
		}
	}

	// Start hub
	if err := h.Start(); err != nil {
		log.Fatal("main:", err)
	}

	// Listen and serve
	log.Println("main: listening on port", cfg.Port)
	log.Fatal(r.Run(":" + fmt.Sprint(cfg.Port)))
}
