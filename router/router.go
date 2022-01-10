//go:build !dev

package router

import (
	"github.com/ItsNotGoodName/reciva-web-remote/web"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func NewEngine() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	handleFS(r, web.GetFS())

	return r
}

func NewUpgrader() *websocket.Upgrader {
	return &websocket.Upgrader{}
}
