package main

import (
	"github.com/ItsNotGoodName/reciva-web-remote/api"
	"github.com/ItsNotGoodName/reciva-web-remote/pkg/goupnpsub"
	"github.com/ItsNotGoodName/reciva-web-remote/pkg/radio"
	"github.com/ItsNotGoodName/reciva-web-remote/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	// Create and start controlpoint
	cp := goupnpsub.NewControlPoint()
	go cp.Start()

	// Create radio hub
	h := radio.NewHub(cp)

	// Create api
	a := api.NewService(h)

	// Get router
	r := gin.Default()

	// Create v1 route
	v1 := r.Group("/v1")

	// Add routes to v1 group
	routes.AddRadioRoutes(v1, a)

	// listen and serve on 0.0.0.0:8080
	r.Run()
}
