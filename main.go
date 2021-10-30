package main

import (
	"fmt"
	"log"

	"github.com/ItsNotGoodName/reciva-web-remote/api"
	"github.com/ItsNotGoodName/reciva-web-remote/config"
	"github.com/ItsNotGoodName/reciva-web-remote/pkg/goupnpsub"
	"github.com/ItsNotGoodName/reciva-web-remote/pkg/radio"
	"github.com/ItsNotGoodName/reciva-web-remote/routes"
)

func main() {
	// Create config
	cfg := config.NewConfig(config.WithFlag)

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

	// Add radio routes
	routes.AddRadioRoutes(r.Group(cfg.APIURI), a, u)

	// Add config routes
	routes.AddConfigRoutes(r.Group(cfg.APIURI), cfg)

	// Listen and serve
	log.Println("main: listening on port", cfg.Port)
	log.Fatal(r.Run(":" + fmt.Sprint(cfg.Port)))
}
