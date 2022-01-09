package api

import (
	"context"

	"github.com/ItsNotGoodName/reciva-web-remote/pkg/radio"
	"github.com/ItsNotGoodName/reciva-web-remote/store"
)

func NewPresetAPI(s *store.Store, h *radio.Hub) *PresetAPI {
	p := PresetAPI{s: s, h: h}
	// h.PresetMutator = p.PresetMutator
	return &p
}

func (p *PresetAPI) PresetMutator(ctx context.Context, rp *radio.Preset) {
	sp, err := p.s.ReadPreset(ctx, rp.URL)
	if err != nil || sp.NewName == "" {
		rp.Name = rp.Title
		return
	}
	rp.Name = sp.NewName
}

// ReadPresets returns all presets.
func (p *PresetAPI) ReadPresets(ctx context.Context) ([]store.Preset, error) {
	return p.s.ReadPresets(ctx)
}

// ReadPreset returns a preset by its URL.
func (p *PresetAPI) ReadPreset(ctx context.Context, url string) (*store.Preset, error) {
	pt, err := p.s.ReadPreset(ctx, url)
	if err == store.ErrNotFound {
		return nil, ErrPresetNotFound
	}
	return pt, err
}

// UpdatePreset updates a preset.
func (p *PresetAPI) UpdatePreset(ctx context.Context, preset *store.Preset) error {
	// Update the preset
	err := p.s.UpdatePreset(ctx, preset)
	if err == store.ErrNotFound {
		return ErrPresetNotFound
	}
	//p.h.RefreshPresets()

	return err
}
