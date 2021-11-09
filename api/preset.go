package api

import (
	"context"

	"github.com/ItsNotGoodName/reciva-web-remote/config"
	"github.com/ItsNotGoodName/reciva-web-remote/pkg/radio"
	"github.com/ItsNotGoodName/reciva-web-remote/store"
)

func NewPresetAPI(s *store.Store, h *radio.Hub) *PresetAPI {
	p := PresetAPI{s: s, h: h}
	h.PresetMutator = p.PresetMutator
	return &p
}

func (p *PresetAPI) PresetMutator(ctx context.Context, preset *radio.Preset) {
	stream, err := p.s.ReadPresetByURL(ctx, preset.URL)
	if err != nil {
		preset.Name = preset.Title
		return
	}
	preset.Name = stream.NewName
}

// ReadPresets returns all presets.
func (p *PresetAPI) ReadPresets(ctx context.Context) ([]config.Preset, error) {
	return p.s.ReadPresets(ctx)
}

// ReadPresetByURL returns a preset by its URL.
func (p *PresetAPI) ReadPresetByURL(ctx context.Context, url string) (*config.Preset, error) {
	return p.s.ReadPresetByURL(ctx, url)
}
