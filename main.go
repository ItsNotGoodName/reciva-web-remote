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
	cfg := config.NewConfig()

	// Create and start controlpoint
	cp := goupnpsub.NewControlPointWithPort(cfg.CPort)
	go cp.Start()

	// Create radio hub
	h := radio.NewHub(cp)

	// Create api
	a := api.NewService(h)

	// Create router
	r := newRouter()

	// Create websocket upgrader
	u := newUpgrader()

	// Add radio routes to v1 group
	v1 := r.Group("/v1")
	routes.AddRadioRoutes(v1, a, u)

	// Listen and server
	log.Println("main: listening on port", cfg.Port)
	log.Fatal(r.Run(":" + fmt.Sprint(cfg.Port)))
}
