package routes

import (
	"net/http"

	"github.com/ItsNotGoodName/reciva-web-remote/config"
	"github.com/gin-gonic/gin"
)

type ConfigJSON struct {
	PresetsEnabled bool     `json:"presetsEnabled"`
	Presets        []string `json:"presets"`
}

func getConfigHandler(cfg *config.Config) func(c *gin.Context) {
	cfgJSON := ConfigJSON{
		PresetsEnabled: cfg.PresetsEnabled,
		Presets:        cfg.Presets,
	}
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, cfgJSON)
	}
}

func AddConfigRoutes(r *gin.Engine, cfg *config.Config) {
	r.GET("/config.json", getConfigHandler(cfg))
}
