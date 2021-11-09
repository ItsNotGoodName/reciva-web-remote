package routes

import (
	"context"
	"log"
	"net/http"
	"net/url"

	"github.com/ItsNotGoodName/reciva-web-remote/api"
	"github.com/gin-gonic/gin"
)

func AddPresetRoutes(r *gin.RouterGroup, p *api.PresetAPI) {
	r.GET("/presets", func(c *gin.Context) {
		presets, err := p.ReadPresets(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}

		c.JSON(http.StatusOK, presets)
	})
}

func newPresetURIHandler(p *api.PresetAPI, url string) func(c *gin.Context) {
	return func(c *gin.Context) {
		// Read the stream
		stream, err := p.ReadPresetByURL(c, url)
		if err != nil {
			code := http.StatusInternalServerError
			c.JSON(code, gin.H{"err": err.Error()})
			return
		}

		c.Writer.WriteString(stream.NewURL)
	}
}

func AddPresetURIRoutes(r *gin.Engine, p *api.PresetAPI) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	presets, err := p.ReadPresets(ctx)
	cancel()

	if err != nil {
		log.Fatal("routes.AddPresetURIRoutes:", err.Error())
	}

	for _, preset := range presets {
		if preset.NewURL == "" {
			continue
		}
		u, _ := url.Parse(preset.URL)
		r.GET(u.Path, newPresetURIHandler(p, preset.URL))
	}
}
