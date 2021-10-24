package routes

import (
	"net/http"
	"strconv"

	"github.com/ItsNotGoodName/reciva-web-remote/api"
	"github.com/ItsNotGoodName/reciva-web-remote/config"
	"github.com/gin-gonic/gin"
)

type StreamPost struct {
	Name    *string `json:"name,omitempty"`
	Content *string `json:"content,omitempty"`
}

func AddPresetAPIRoutes(r *gin.RouterGroup, p *api.PresetAPI) {
	r.GET("/presets", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	})
	r.POST("/preset", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"uri": "/01.m3u", "sid": 0})
	})
	r.GET("/streams", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	})
	r.GET("/stream/:SID", func(c *gin.Context) {
		sidStr, ok := c.Params.Get("SID")
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		sid, err := strconv.Atoi(sidStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}

		c.JSON(http.StatusOK, gin.H{"sid": sid, "name": "", "content": ""})
	})
	r.DELETE("/stream/:SID", func(c *gin.Context) {
		sidStr, ok := c.Params.Get("SID")
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		sid, err := strconv.Atoi(sidStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}

		deleted := p.S.DeleteStream(sid)
		if deleted < 1 {
			c.JSON(http.StatusNotFound, gin.H{})
			return
		}

		c.JSON(http.StatusOK, gin.H{"deleted": deleted})
	})
	r.PUT("/stream", func(c *gin.Context) {
		// Get JSON
		var streamPost StreamPost
		if err := c.BindJSON(&streamPost); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		// Check JSON
		if streamPost.Name == nil || streamPost.Content == nil {
			c.Status(http.StatusBadRequest)
			return
		}

		// Add stream
		st, err := p.S.AddStream(*streamPost.Name, *streamPost.Content)
		if err != nil {
			c.Status(http.StatusConflict)
			return
		}

		// Return stream
		c.JSON(http.StatusOK, st)
	})
	r.POST("/stream/:SID", func(c *gin.Context) {
		// Get sid
		sidStr, ok := c.Params.Get("SID")
		if !ok {
			c.Status(http.StatusInternalServerError)
			return
		}

		// Convert sid to int
		sid, err := strconv.Atoi(sidStr)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		// Get JSON
		var streamPost StreamPost
		if err := c.BindJSON(&streamPost); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		// Get stream
		st, ok := p.S.GetStream(sid)
		if !ok {
			c.Status(http.StatusNotFound)
			return
		}

		// Update and save stream
		save := false
		if streamPost.Name != nil && *streamPost.Name != st.Name {
			st.Name = *streamPost.Name
			save = true
		}
		if streamPost.Content != nil && *streamPost.Content != st.Content {
			st.Content = *streamPost.Content
			save = true
		}
		if save {
			success := p.S.UpdateStream(st)
			if !success {
				c.Status(http.StatusConflict)
				return
			}
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
