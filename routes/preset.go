package routes

import (
	"net/http"

	"github.com/ItsNotGoodName/reciva-web-remote/api"
	"github.com/ItsNotGoodName/reciva-web-remote/config"
	"github.com/gin-gonic/gin"
)

func AddPresetAPIRoutes(r *gin.RouterGroup, p *api.PresetAPI) {
	r.GET("/presets", func(c *gin.Context) {
		c.JSON(http.StatusOK, p.GetPresets())
	})
	r.POST("/preset", func(c *gin.Context) {
		// Get JSON
		var presetReq api.PresetReq
		if err := c.BindJSON(&presetReq); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		// Check JSON
		if presetReq.SID == nil || presetReq.URI == nil {
			c.Status(http.StatusBadRequest)
			return
		}

		// Get preset
		pt, ok := p.S.GetPreset(*presetReq.URI)
		if !ok {
			c.Status(http.StatusNotFound)
			return
		}

		// Update preset
		if !p.UpdatePreset(pt, &presetReq) {
			c.Status(http.StatusNotFound)
			return
		}

		c.JSON(http.StatusOK, pt)
	})
	r.GET("/streams", func(c *gin.Context) {
		c.JSON(http.StatusOK, p.GetStreams())
	})
	r.GET("/stream/:SID", ensureSID, func(c *gin.Context) {
		sid := c.GetInt("sid")

		// Get stream
		st, ok := p.S.GetStream(sid)
		if !ok {
			c.Status(http.StatusNotFound)
			return
		}

		c.JSON(http.StatusOK, st)
	})
	r.DELETE("/stream/:SID", ensureSID, func(c *gin.Context) {
		sid := c.GetInt("sid")

		// Delete Stream
		if !p.DeleteStream(sid) {
			c.Status(http.StatusNotFound)
			return
		}

		c.Status(http.StatusOK)
	})
	r.POST("/stream/new", func(c *gin.Context) {
		// Get JSON
		var streamReq api.StreamReq
		if err := c.BindJSON(&streamReq); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		// Check JSON
		if streamReq.Name == nil || streamReq.Content == nil {
			c.Status(http.StatusBadRequest)
			return
		}

		// Add stream
		st, ok := p.S.AddStream(*streamReq.Name, *streamReq.Content)
		if !ok {
			c.Status(http.StatusConflict)
			return
		}

		// Return stream
		c.JSON(http.StatusOK, st)
	})
	r.POST("/stream/:SID", ensureSID, func(c *gin.Context) {
		sid := c.GetInt("sid")

		// Get JSON
		var streamReq api.StreamReq
		if err := c.BindJSON(&streamReq); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		// Get stream
		st, ok := p.S.GetStream(sid)
		if !ok {
			c.Status(http.StatusNotFound)
			return
		}

		// Update stream
		if !p.UpdateStream(st, &streamReq) {
			c.Status(http.StatusConflict)
			return
		}

		c.JSON(http.StatusOK, st)
	})
}

func getPresetHandler(p *api.PresetAPI, uri string) func(c *gin.Context) {
	return func(c *gin.Context) {
		// Get preset
		preset, ok := p.S.GetPreset(uri)
		if !ok || preset.SID == 0 {
			c.Status(http.StatusNotFound)
			return
		}

		// Get stream
		stream, ok := p.S.GetStream(preset.SID)
		if !ok {
			c.Status(http.StatusNotFound)
			return
		}

		c.Writer.Write([]byte(stream.Content))
	}
}

func AddPresetRoutes(r *gin.Engine, cfg *config.Config, p *api.PresetAPI) {
	for _, v := range cfg.Presets {
		r.GET(v, getPresetHandler(p, v))
	}
}
