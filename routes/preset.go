package routes

import (
	"net/http"

	"github.com/ItsNotGoodName/reciva-web-remote/api"
	"github.com/gin-gonic/gin"
)

func AddPresetRoutes(r *gin.RouterGroup, p *api.PresetAPI) {
	r.GET("/presets", func(c *gin.Context) {
		// Get active presets
		presets, err := p.GetActivePresets(c)
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
			if err == api.ErrPresetNotFound {
				c.JSON(http.StatusNotFound, gin.H{"err": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
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
			if err == api.ErrPresetNotFound {
				c.JSON(http.StatusNotFound, gin.H{"err": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}

		c.JSON(http.StatusOK, preset)
	})
}

func handlePreset(p *api.PresetAPI, uri string) func(c *gin.Context) {
	return func(c *gin.Context) {
		// Get the stream
		stream, err := p.GetStreamByURI(c, uri)
		if err != nil {
			if err == api.ErrPresetNotFound || err == api.ErrStreamNotFound {
				c.JSON(http.StatusNotFound, gin.H{"err": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		}

		// TODO: sanitize the stream content to prevent XSS
		c.Writer.WriteString(stream.Content)
	}
}

func AddPresetRadioRoutes(r *gin.Engine, p *api.PresetAPI) {
	uris := p.GetActiveURIS()
	for _, uri := range uris {
		r.GET(uri, handlePreset(p, uri))
	}
}
