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

	// Define websocket upgrader
	var upgrader websocket.Upgrader

	// Configure mode based on environment
	if gin.Mode() != gin.ReleaseMode {
		r.Use(routes.CORS())
		upgrader = websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		}
	} else {
		upgrader = websocket.Upgrader{}
	}

	// Create v1 route
	v1 := r.Group("/v1")

	// Add routes to v1 group
	routes.AddRadioRoutes(v1, a, &upgrader)

	// Listen and serve on 0.0.0.0:8080
	r.Run()
}
