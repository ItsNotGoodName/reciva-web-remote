//go:build dev

package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func NewEngine() *gin.Engine {
	r := gin.Default()

	r.Use(cors())

	return r
}

func NewUpgrader() *websocket.Upgrader {
	return &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
}
