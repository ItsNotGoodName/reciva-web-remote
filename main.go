package main

import (
	"net/http"

	"github.com/ItsNotGoodName/reciva-web-remote/api"
	"github.com/ItsNotGoodName/reciva-web-remote/pkg/goupnpsub"
	"github.com/ItsNotGoodName/reciva-web-remote/pkg/radio"
	"github.com/ItsNotGoodName/reciva-web-remote/routes"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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

	// Get websocket upgrader
	upgrader := websocket.Upgrader{}

	// Ignore CORS when not in production
	if gin.Mode() != gin.ReleaseMode {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		r.Use(routes.CORS())
	}

	// Create v1 route
	v1 := r.Group("/v1")

	// Add routes to v1 group
	routes.AddRadioRoutes(v1, a, &upgrader)

	// listen and serve on 0.0.0.0:8080
	r.Run()
}
