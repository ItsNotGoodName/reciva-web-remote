package routes

import (
	"net/http"

	"github.com/ItsNotGoodName/reciva-web-remote/config"
	"github.com/gin-gonic/gin"
)

type ConfigJSON struct {
	PresetsEnabled bool     `json:"presetsEnabled"`
	PresetURLS     []string `json:"presetURLS"`
}

func AddConfigRoutes(r *gin.RouterGroup, cfg *config.Config) {
	cfgJSON := ConfigJSON{PresetsEnabled: cfg.PresetsEnabled, PresetURLS: cfg.URLS}

	r.GET("/config", func(c *gin.Context) {
		c.JSON(http.StatusOK, cfgJSON)
	})
}
