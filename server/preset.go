package server

import (
	"net/http"

	"github.com/ItsNotGoodName/reciva-web-remote/core"
	"github.com/gin-gonic/gin"
)

func handlePresetGet(presetStore core.PresetStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get preset
		preset, err := presetStore.GetPreset(c, c.Query("url"))
		if err != nil {
			code := http.StatusInternalServerError
			if err == core.ErrPresetNotFound {
				code = http.StatusNotFound
			}
			renderError(c, code, err)
			return
		}

		c.JSON(http.StatusOK, preset)
	}
}

func handlePresetList(presetStore core.PresetStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		// List presets
		presets, err := presetStore.ListPresets(c)
		if err != nil {
			renderError(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, presets)
	}
}

func handlePresetUpdate(presetStore core.PresetStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse JSON body
		var preset core.Preset
		if err := c.BindJSON(&preset); err != nil {
			renderError(c, http.StatusBadRequest, err)
			return
		}

		// Update preset
		if err := presetStore.UpdatePreset(c, &preset); err != nil {
			code := http.StatusInternalServerError
			if err == core.ErrPresetNotFound {
				code = http.StatusNotFound
			}
			renderError(c, code, err)
			return
		}

		c.JSON(http.StatusOK, preset)
	}
}

func handlePresetGetNewURL(presetStore core.PresetStore, url string) func(c *gin.Context) {
	return func(c *gin.Context) {
		// Get preset
		preset, err := presetStore.GetPreset(c, url)
		if err != nil {
			code := http.StatusInternalServerError
			if err == core.ErrPresetNotFound {
				code = http.StatusNotFound
			}
			renderError(c, code, err)
			return
		}

		c.Writer.WriteString(preset.NewURL)
	}
}
