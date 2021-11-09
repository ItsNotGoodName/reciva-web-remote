package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

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
			// Add preset routes
			routes.AddPresetRoutes(r.Group(cfg.APIURI), p)
			// Add preset routes based on their uri
			routes.AddPresetURIRoutes(r, p)
		} else {
			log.Println("main: presets could not be enabled:", err)
		}
	}

	// Start hub and interrupt handler
	if err := h.Start(); err != nil {
		log.Fatal("main:", err)
	}

	// Start router
	go routes.Start(cfg, r)

	// Listen for interrupt
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	log.Println("main: stopping")

	// Shutdown hub
	if err := h.Stop(); err != nil {
		log.Fatal("main:", err)
	}
}
