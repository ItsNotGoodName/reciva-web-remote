package server

import (
	"context"
	"log"

	"github.com/ItsNotGoodName/reciva-web-remote/api"
)

// GetPresetURLS returns all urls for presets
func GetPresetURLS(p *api.PresetAPI) []string {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	presets, err := p.ReadPresets(ctx)
	cancel()
	if err != nil {
		log.Fatal("server.GetPresets:", err.Error())
	}
	urls := make([]string, len(presets))
	for i := range presets {
		urls[i] = presets[i].URL
	}
	return urls
}
