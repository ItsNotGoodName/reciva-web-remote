package api

import (
	"context"
	"database/sql"

	"github.com/ItsNotGoodName/reciva-web-remote/pkg/radio"
	"github.com/ItsNotGoodName/reciva-web-remote/store"
)

func NewPresetAPI(s *store.Store, h *radio.Hub) *PresetAPI {
	p := PresetAPI{s: s, h: h}
	h.PresetMutator = p.PresetMutator
	return &p
}

func (p *PresetAPI) PresetMutator(ctx context.Context, preset *radio.Preset) {
	stream, err := p.ReadStreamByURL(ctx, preset.URL)
	if err != nil {
		preset.Name = preset.Title
		return
	}
	preset.Name = stream.Name
}

// ReadPresets returns all presets.
func (p *PresetAPI) ReadPresets(ctx context.Context) ([]*store.Preset, error) {
	return p.s.ReadPresets(ctx)
}

// ReadActiveURLS returns active presets as an array of URLS.
func (p *PresetAPI) ReadActiveURLS() []string {
	return p.s.URLS
}

// ReadActivePresets returns active presets.
func (p *PresetAPI) ReadActivePresets(ctx context.Context) ([]*store.Preset, error) {
	var pts []*store.Preset

	for _, url := range p.s.URLS {
		p, err := p.s.ReadPreset(ctx, url)
		if err != nil {
			return nil, err
		}
		pts = append(pts, p)
	}

	return pts, nil
}

// ReadPreset returns preset by url.
func (p *PresetAPI) ReadPreset(ctx context.Context, url string) (*store.Preset, error) {
	preset, err := p.s.ReadPreset(ctx, url)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrPresetNotFound
		}
		return nil, err
	}

	return preset, nil
}

// UpdatePresetRequest is request for UpdatePreset.
type UpdatePresetRequest struct {
	SID int    `json:"sid"`
	URL string `json:"url"`
}

// UpdatePreset updates preset.
func (p *PresetAPI) UpdatePreset(ctx context.Context, req *UpdatePresetRequest) (*store.Preset, error) {
	preset := &store.Preset{
		URL: req.URL,
		SID: req.SID,
	}
	ok, err := p.s.UpdatePreset(ctx, preset)
	if err != nil {
		return nil, err
	}
	if !ok {
		_, err := p.ReadStream(ctx, req.SID)
		if err != nil {
			return nil, err
		}
		return nil, ErrPresetNotFound
	}
	p.h.RefreshPresets()
	return preset, nil
}

// ClearPresetRequest is request for ClearPreset.
type ClearPresetRequest struct {
	URL string `json:"url"`
}

// ClearPreset clears preset's SID field.
func (p *PresetAPI) ClearPreset(ctx context.Context, req *ClearPresetRequest) (*store.Preset, error) {
	preset, err := p.ReadPreset(ctx, req.URL)
	if err != nil {
		return nil, err
	}
	err = p.s.ClearPreset(ctx, preset)
	if err != nil {
		return nil, err
	}
	p.h.RefreshPresets()
	return preset, nil
}
