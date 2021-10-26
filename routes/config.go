package routes

import (
	"net/http"

	"github.com/ItsNotGoodName/reciva-web-remote/config"
	"github.com/gin-gonic/gin"
)

type ConfigJSON struct {
	EnablePresets bool     `json:"enablePresets"`
	Presets       []string `json:"presets"`
}

func configHandler(cfg *config.Config) func(c *gin.Context) {
	cfgJSON := ConfigJSON{
		EnablePresets: cfg.PresetsEnabled,
		Presets:       cfg.Presets,
	}
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, cfgJSON)
	}
}

func AddConfigRoutes(r *gin.Engine, cfg *config.Config) {
	r.GET("/config.json", configHandler(cfg))
}
