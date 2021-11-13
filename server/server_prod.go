//go:build prod

package server

import (
	"github.com/ItsNotGoodName/reciva-web-remote/web"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func NewRouter() *gin.Engine {
	// Set release mode
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	web.AddWebRoutes(r)

	return r
}

func NewUpgrader() *websocket.Upgrader {
	return &websocket.Upgrader{}
}
