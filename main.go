package main

import (
	"fmt"
	"log"

	"github.com/ItsNotGoodName/reciva-web-remote/api"
	"github.com/ItsNotGoodName/reciva-web-remote/config"
	"github.com/ItsNotGoodName/reciva-web-remote/pkg/goupnpsub"
	"github.com/ItsNotGoodName/reciva-web-remote/pkg/radio"
	"github.com/ItsNotGoodName/reciva-web-remote/routes"
	"github.com/ItsNotGoodName/reciva-web-remote/store"
)

func main() {
	// Create config
	cfg := config.NewConfig()

	// Create store
	s, sErr := store.NewService(cfg)
	if sErr != nil {
		log.Println("main.main: ", sErr)
	}

	// Create and start controlpoint
	cp := goupnpsub.NewControlPointWithPort(cfg.CPort)
	go cp.Start()

	// Create radio hub
	h := radio.NewHub(cp)

	// Create api
	a := api.NewAPI(h)

	// Create router
	r := newRouter()

	// Create websocket upgrader
	u := newUpgrader()

	// Add radio routes to v1 group
	v1 := r.Group("/v1")
	routes.AddRadioRoutes(v1, a, u)

	// Add config routes
	routes.AddConfigRoutes(r, cfg)

	// Check if store has no error
	if sErr == nil {
		// Check if presets are enabled
		if cfg.EnablePresets {
			// Create preset api
			p := api.NewPresetAPI(a, s)
			// Add preset routes
			routes.AddPresetRoutes(r, p)
		} else {
			// Close store
			s.Cancel()
		}
	}

	// Listen and serve
	log.Println("main: listening on port", cfg.Port)
	log.Fatal(r.Run(":" + fmt.Sprint(cfg.Port)))
}
