package server

import (
	"net/url"
)

func (s *Server) routes() {
	apiURI := "/v1"

	// Radio api
	{
		g := s.r.Group(apiURI)

		g.GET("/radio/:uuid", ensureUUID(), s.handleGetRadio())
		g.GET("/radio/ws", s.handleGetRadioWS())
		g.GET("/radios", s.handleGetRadios())
		g.POST("/radios", s.handlePostRadios())
		g.Use(ensureRadio(s.h))
		g.PATCH("/radio/:uuid", s.handlePatchRadio())
		g.POST("/radio/:uuid/renew", s.handlePostRadioRenew())
		g.POST("/radio/:uuid/volume", s.handlePostRadioVolume())
	}

	// Preset api
	{
		g := s.r.Group(apiURI)

		g.GET("/preset", s.handleGetPreset())
		g.GET("/presets", s.handleGetPresets())
		g.POST("/preset", s.handlePostPreset())
	}

	// Preset
	{
		urls := GetPresetURLS(s.p)
		for _, rawURL := range urls {
			u, _ := url.Parse(rawURL)
			s.r.GET(u.Path, s.handleGetPresetNewURL(rawURL))
		}
	}
}
