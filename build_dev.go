//go:build !prod

package main

import (
	"net/http"

	"github.com/ItsNotGoodName/reciva-web-remote/routes"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func newRouter() *gin.Engine {
	r := gin.Default()

	// Enable CORS on gin
	r.Use(routes.CORS())

	return r
}


func newUpgrader() *websocket.Upgrader {
	// Ignore origin on websocket
	return &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
}