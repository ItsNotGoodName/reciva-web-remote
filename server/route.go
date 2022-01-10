package server

import (
	"github.com/ItsNotGoodName/reciva-web-remote/pkg/radio"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func Route(
	engine *gin.Engine,
	upgrader *websocket.Upgrader,
	hub *radio.Hub,
) {
	apiURI := "/v1"

	{
		r := engine.Group(apiURI)

		r.GET("/radios", handleRadioList(hub))
		r.POST("/radios", handleRadioDiscover(hub))
		r.GET("/radio/ws", handleRadioWS(hub, upgrader))
		r.Use(ensureRadio(hub))
		r.GET("/radio/:uuid", handleRadioGet())
		r.PATCH("/radio/:uuid", handleRadioPatch())
		r.POST("/radio/:uuid", handleRadioRefresh())
		r.POST("/radio/:uuid/volume", handleRadioVolumeRefresh())
	}
}
