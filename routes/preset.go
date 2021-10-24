package routes

import (
	"net/http"

	"github.com/ItsNotGoodName/reciva-web-remote/api"
	"github.com/ItsNotGoodName/reciva-web-remote/config"
	"github.com/gin-gonic/gin"
)

func AddPresetAPIRoutes(r *gin.RouterGroup, p *api.PresetAPI) {
	r.GET("/presets", func(c *gin.Context) {
		c.JSON(http.StatusOK, p.S.GetPresets())
	})
	r.POST("/preset", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"uri": "/01.m3u", "sid": 0})
	})
	r.GET("/streams", func(c *gin.Context) {
		c.JSON(http.StatusOK, p.S.GetStreams())
	})
	r.GET("/stream/:SID", ensureSID, func(c *gin.Context) {
		sid := c.GetInt("sid")

		c.JSON(http.StatusOK, gin.H{"sid": sid, "name": "", "content": ""})
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
	r.PUT("/stream", func(c *gin.Context) {
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
		st, err := p.S.AddStream(*streamReq.Name, *streamReq.Content)
		if err != nil {
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

func AddPresetRoutes(r *gin.Engine, cfg *config.Config, p *api.PresetAPI) {
	for _, v := range cfg.Presets {
		r.GET(v, func(c *gin.Context) {
			preset, ok := p.S.GetPreset(v)
			if !ok || preset.StreamID == 0 {
				c.Status(http.StatusNotFound)
				return
			}

			stream, ok := p.S.GetStream(preset.StreamID)
			if !ok {
				c.Status(http.StatusNotFound)
				return
			}

			c.Writer.Write([]byte(stream.Content))
		})
	}
}
