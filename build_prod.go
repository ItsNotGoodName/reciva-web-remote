//go:build prod

package main

import (
	"embed"

	"github.com/ItsNotGoodName/reciva-web-remote/routes"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	//go:embed web/dist
	dist embed.FS
)

func NewRouter() *gin.Engine {
	// Set release mode
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	// Add web routes
	routes.AddWebRoutes(r, &dist)

	return r
}

func NewUpgrader() *websocket.Upgrader {
	return &websocket.Upgrader{}
}
