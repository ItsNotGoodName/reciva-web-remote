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

func configJSON(cfg *config.Config) *ConfigJSON {
	return &ConfigJSON{
		EnablePresets: cfg.PresetsEnabled,
		Presets:       cfg.Presets,
	}
}

func AddConfigRoutes(r *gin.Engine, cfg *config.Config) {
	r.GET("/config.json", func(c *gin.Context) {
		c.JSON(http.StatusOK, configJSON(cfg))
	})
}
