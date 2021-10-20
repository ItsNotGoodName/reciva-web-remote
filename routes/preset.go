package routes

import (
	"net/http"

	"github.com/ItsNotGoodName/reciva-web-remote/config"
	"github.com/gin-gonic/gin"
)

func AddPresetRoutes(r *gin.Engine, cfg *config.Config) {
	for p := range cfg.Presets {
		preset := cfg.Presets[p]
		r.GET(preset, func(c *gin.Context) {
			c.String(http.StatusOK, preset)
		})
	}
}
