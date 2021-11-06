package api

import (
	"context"
	"database/sql"
	"net/url"

	"github.com/ItsNotGoodName/reciva-web-remote/pkg/radio"
	"github.com/ItsNotGoodName/reciva-web-remote/store"
)

func NewPresetAPI(s *store.Store, h *radio.Hub) *PresetAPI {
	p := PresetAPI{s: s}
	h.PresetMutator = p.PresetMutator
	return &p
}

func (p *PresetAPI) PresetMutator(ctx context.Context, preset *radio.Preset) {
	uri, err := url.ParseRequestURI(preset.URL)
	if err != nil {
		preset.Name = preset.Title
		return
	}
	stream, err := p.GetStreamByURI(ctx, uri.Path)
	if err != nil {
		preset.Name = preset.Title
		return
	}
	preset.Name = stream.Name
}

// GetPresets returns all presets.
func (p *PresetAPI) GetPresets(ctx context.Context) ([]*store.Preset, error) {
	return p.s.GetPresets(ctx)
}

// GetActiveURIS returns active presets as an array of URIS.
func (p *PresetAPI) GetActiveURIS() []string {
	return p.s.Presets
}

// GetActivePresets returns active presets.
func (p *PresetAPI) GetActivePresets(ctx context.Context) ([]*store.Preset, error) {
	var pts []*store.Preset

	for _, uri := range p.s.Presets {
		p, err := p.s.GetPreset(ctx, uri)
		if err != nil {
			return nil, err
		}
		pts = append(pts, p)
	}

	return pts, nil
}

// GetPreset returns preset by uri.
func (p *PresetAPI) GetPreset(ctx context.Context, uri string) (*store.Preset, error) {
	preset, err := p.s.GetPreset(ctx, uri)
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
	URI string `json:"uri"`
}

// UpdatePreset updates preset.
func (p *PresetAPI) UpdatePreset(ctx context.Context, req *UpdatePresetRequest) (*store.Preset, error) {
	preset := &store.Preset{
		URI: req.URI,
		SID: req.SID,
	}
	ok, err := p.s.UpdatePreset(ctx, preset)
	if err != nil {
		return nil, err
	}
	if !ok {
		_, err := p.GetStream(ctx, req.SID)
		if err != nil {
			return nil, err
		}
		return nil, ErrPresetNotFound
	}

	return preset, nil
}

// ClearPresetRequest is request for ClearPreset.
type ClearPresetRequest struct {
	URI string `json:"uri"`
}

// ClearPreset clears preset's SID field.
func (p *PresetAPI) ClearPreset(ctx context.Context, req *ClearPresetRequest) (*store.Preset, error) {
	preset, err := p.GetPreset(ctx, req.URI)
	if err != nil {
		return nil, err
	}
	err = p.s.ClearPreset(ctx, preset)
	if err != nil {
		return nil, err
	}
	return preset, nil
}
