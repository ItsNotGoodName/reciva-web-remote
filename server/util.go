package server

import (
	"context"
	"log"

	"github.com/ItsNotGoodName/reciva-web-remote/core"
	"github.com/gin-gonic/gin"
)

func renderError(c *gin.Context, code int, err error) {
	c.JSON(code, gin.H{"err": err.Error()})
}

// getPresetURLS returns all urls for presets.
func getPresetURLS(p core.PresetStore) []string {
	presets, err := p.ListPresets(context.Background())
	if err != nil {
		log.Fatal("server.getPresetURLS(ERROR):", err)
	}

	urls := make([]string, len(presets))
	for i := range presets {
		urls[i] = presets[i].URL
	}

	return urls
}
