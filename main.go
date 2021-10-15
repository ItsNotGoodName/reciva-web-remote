package main

import (
	"github.com/ItsNotGoodName/reciva-web-remote/api"
	"github.com/ItsNotGoodName/reciva-web-remote/pkg/goupnpsub"
	"github.com/ItsNotGoodName/reciva-web-remote/pkg/radio"
	"github.com/ItsNotGoodName/reciva-web-remote/routes"
)

func main() {
	// Create and start controlpoint
	cp := goupnpsub.NewControlPoint()
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

	// Listen and serve PORT env variable or 8080
	r.Run()
}
