package routes

import (
	"context"
	"log"
	"net/http"
	"net/url"

	"github.com/ItsNotGoodName/reciva-web-remote/api"
	"github.com/ItsNotGoodName/reciva-web-remote/config"
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

	r.GET("/preset", func(c *gin.Context) {
		url := c.Query("url")

		// Get preset
		preset, err := p.ReadPreset(c, url)
		if err != nil {
			code := http.StatusInternalServerError
			if err == api.ErrPresetNotFound {
				code = http.StatusNotFound
			}
			c.JSON(code, gin.H{"err": err.Error()})
			return
		}

		c.JSON(http.StatusOK, preset)
	})

	r.POST("/preset", func(c *gin.Context) {
		var preset config.Preset
		if err := c.BindJSON(&preset); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}

		if err := p.UpdatePreset(c, &preset); err != nil {
			code := http.StatusInternalServerError
			if err == api.ErrPresetNotFound {
				code = http.StatusNotFound
			}
			c.JSON(code, gin.H{"err": err.Error()})
			return
		}

		c.JSON(http.StatusOK, preset)
	})
}

func newPresetURIHandler(p *api.PresetAPI, url string) func(c *gin.Context) {
	return func(c *gin.Context) {
		// Read the stream
		stream, err := p.ReadPreset(c, url)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
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
		u, _ := url.Parse(preset.URL)
		r.GET(u.Path, newPresetURIHandler(p, preset.URL))
	}
}
