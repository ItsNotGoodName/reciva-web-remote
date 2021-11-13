//go:build !prod

package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	// Enable CORS on gin
	r.Use(CORS())

	return r
}

func NewUpgrader() *websocket.Upgrader {
	// Ignore origin on websocket
	return &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
}
