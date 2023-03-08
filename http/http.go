package http

import (
	"fmt"

	"github.com/ItsNotGoodName/reciva-web-remote/internal/hub"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func Start(a API, port int) {
	e := echo.New()
	e.HideBanner = true

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	api := e.Group("/api")

	api.GET("/build", a.GetBuild)
	api.GET("/presets", a.ListPresets)
	api.POST("/presets", a.UpdatePreset)
	api.GET("/presets/*", a.GetPreset)

	api.POST("/radios", a.DiscoverRadios)
	api.GET("/radios", a.ListRadios)

	apiRadios := api.Group("/radios/:uuid")
	apiRadios.Use(a.RadioMiddleware)
	apiRadios.GET("", a.GetRadio)
	apiRadios.POST("/volume", a.RefreshRadioVolume)
	apiRadios.POST("/subscription", a.RefreshRadioSubscription)

	api.GET("/states", a.ListStates)

	apiStates := api.Group("/states/:uuid")
	apiStates.Use(a.RadioMiddleware)
	apiStates.GET("", a.GetState)
	apiStates.PATCH("", a.PatchState)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}

type RadioContext struct {
	echo.Context
	Radio hub.Radio
}
