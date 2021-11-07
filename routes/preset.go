package routes

import (
	"net/http"
	"net/url"

	"github.com/ItsNotGoodName/reciva-web-remote/api"
	"github.com/gin-gonic/gin"
)

func AddPresetRoutes(r *gin.RouterGroup, p *api.PresetAPI) {
	r.GET("/presets", func(c *gin.Context) {
		// Read active presets
		presets, err := p.ReadActivePresets(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}

		c.JSON(http.StatusOK, presets)
	})

	r.POST("/preset", func(c *gin.Context) {
		// Parse the JSON in the body
		updateReq := api.UpdatePresetRequest{}
		if err := c.BindJSON(&updateReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}

		// Update the preset
		stream, err := p.UpdatePreset(c, &updateReq)
		if err != nil {
			code := http.StatusInternalServerError
			if err == api.ErrPresetNotFound {
				code = http.StatusNotFound
			}
			c.JSON(code, gin.H{"err": err.Error()})
			return
		}

		c.JSON(http.StatusOK, stream)
	})

	r.DELETE("/preset", func(c *gin.Context) {
		// Parse the JSON in the body
		clearReq := api.ClearPresetRequest{}
		if err := c.BindJSON(&clearReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}

		// Delete the preset
		preset, err := p.ClearPreset(c, &clearReq)
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
}

func newPresetStreamHandler(p *api.PresetAPI, url string) func(c *gin.Context) {
	return func(c *gin.Context) {
		// Read the stream
		stream, err := p.ReadStreamByURL(c, url)
		if err != nil {
			code := http.StatusInternalServerError
			if err == api.ErrPresetNotFound || err == api.ErrStreamNotFound {
				code = http.StatusNotFound
			}
			c.JSON(code, gin.H{"err": err.Error()})
			return
		}

		// TODO: sanitize the stream content to prevent XSS
		c.Writer.WriteString(stream.Content)
	}
}

func AddPresetStreamRoutes(r *gin.Engine, p *api.PresetAPI) {
	urls := p.ReadActiveURLS()
	for _, rawURL := range urls {
		u, _ := url.Parse(rawURL)
		uri := u.Path
		r.GET(uri, newPresetStreamHandler(p, rawURL))
	}
}
