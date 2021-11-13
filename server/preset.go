package server

import (
	"net/http"

	"github.com/ItsNotGoodName/reciva-web-remote/api"
	"github.com/ItsNotGoodName/reciva-web-remote/store"
	"github.com/gin-gonic/gin"
)

func (s *Server) handleGetPreset() func(c *gin.Context) {
	return func(c *gin.Context) {
		url := c.Query("url")

		// Read preset
		preset, err := s.p.ReadPreset(c, url)
		if err != nil {
			code := http.StatusInternalServerError
			if err == api.ErrPresetNotFound {
				code = http.StatusNotFound
			}
			c.JSON(code, gin.H{"err": err.Error()})
			return
		}

		c.JSON(http.StatusOK, preset)
	}
}

func (s *Server) handleGetPresets() func(c *gin.Context) {
	return func(c *gin.Context) {
		// Read presets
		presets, err := s.p.ReadPresets(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}

		c.JSON(http.StatusOK, presets)
	}
}

func (s *Server) handlePostPreset() func(c *gin.Context) {
	return func(c *gin.Context) {
		// Parse JSON body
		var preset store.Preset
		if err := c.BindJSON(&preset); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}

		// Update preset
		if err := s.p.UpdatePreset(c, &preset); err != nil {
			code := http.StatusInternalServerError
			if err == api.ErrPresetNotFound {
				code = http.StatusNotFound
			}
			c.JSON(code, gin.H{"err": err.Error()})
			return
		}

		c.JSON(http.StatusOK, preset)
	}
}

func (s *Server) handleGetPresetNewURL(url string) func(c *gin.Context) {
	return func(c *gin.Context) {
		// Read preset
		preset, err := s.p.ReadPreset(c, url)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"err": err.Error()})
			return
		}

		c.Writer.WriteString(preset.NewURL)
	}
}
