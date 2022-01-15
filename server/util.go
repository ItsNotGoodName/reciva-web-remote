package server

import (
	"context"
	"log"
	"sort"
	"strings"

	"github.com/ItsNotGoodName/reciva-web-remote/core"
	"github.com/gin-gonic/gin"
)

type RenderJSON struct {
	Ok     bool        `json:"ok"`
	Code   int         `json:"code"`
	Result interface{} `json:"result,omitempty"`
	Error  string      `json:"error,omitempty"`
}

func renderError(c *gin.Context, code int, err error) {
	c.JSON(code, RenderJSON{
		Ok:    false,
		Code:  code,
		Error: err.Error(),
	})
}

func renderJSON(c *gin.Context, code int, result interface{}) {
	c.JSON(code, RenderJSON{
		Ok:     true,
		Code:   code,
		Result: result,
	})
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

func sortPresets(presets []core.Preset) {
	sort.Slice(presets, func(i, j int) bool {
		return strings.Compare(presets[i].URL, presets[j].URL) < 0
	})
}
