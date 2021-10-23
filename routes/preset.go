package routes

import (
	"net/http"

	"github.com/ItsNotGoodName/reciva-web-remote/api"
	"github.com/ItsNotGoodName/reciva-web-remote/config"
	"github.com/gin-gonic/gin"
)

func AddPresetAPIRoutes(r *gin.RouterGroup, p *api.PresetAPI) {
}

func AddPresetRoutes(r *gin.Engine, cfg *config.Config, p *api.PresetAPI) {
	for _, v := range cfg.Presets {
		r.GET(v, func(c *gin.Context) {
			preset := p.S.GetPreset(v)
			if preset == nil || preset.StreamID == 0 {
				c.Status(http.StatusNotFound)
				return
			}

			stream := p.S.GetStream(preset.StreamID)
			if stream == nil {
				c.Status(http.StatusNotFound)
				return
			}

			c.Writer.Write([]byte(stream.Content))
		})
	}
}
