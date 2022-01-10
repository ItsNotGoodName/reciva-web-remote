package server

import (
	"net/url"

	"github.com/ItsNotGoodName/reciva-web-remote/core"
	"github.com/ItsNotGoodName/reciva-web-remote/pkg/radio"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func Route(
	engine *gin.Engine,
	upgrader *websocket.Upgrader,
	hub *radio.Hub,
	presetStore core.PresetStore,
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

	{
		r := engine.Group(apiURI)

		r.GET("/preset", handlePresetGet(presetStore))
		r.GET("/presets", handlePresetList(presetStore))
		r.POST("/preset", handlePresetUpdate(presetStore))
	}

	{
		urls := getPresetURLS(presetStore)
		for _, rawURL := range urls {
			u, _ := url.Parse(rawURL)
			engine.GET(u.Path, handlePresetGetNewURL(presetStore, rawURL))
		}
	}
}
