package middleware

import (
	"context"

	"github.com/ItsNotGoodName/reciva-web-remote/core/preset"
	"github.com/ItsNotGoodName/reciva-web-remote/core/state"
)

type Preset struct {
	pub   state.MiddlewarePub
	store preset.PresetStore
}

func NewPreset(pub state.MiddlewarePub, presetStore preset.PresetStore) *Preset {
	return &Preset{
		pub:   pub,
		store: presetStore,
	}
}

func (p *Preset) List(ctx context.Context) ([]preset.Preset, error) {
	return p.store.List(ctx)
}

func (p *Preset) Get(ctx context.Context, url string) (*preset.Preset, error) {
	return p.store.Get(ctx, url)
}

func (p *Preset) Update(ctx context.Context, preset *preset.Preset) error {
	err := p.store.Update(ctx, preset)
	if err != nil {
		return err
	}

	p.pub.Publish()
	return nil
}

func (p *Preset) Apply(s *state.State, c state.Changed) state.Changed {
	ctx := context.Background()

	return c.Merge(p.fragmentPresets(ctx, s, c)).Merge(p.fragmentTitleAndURL(ctx, s, c))
}

func (p *Preset) fragmentPresets(ctx context.Context, s *state.State, c state.Changed) state.Changed {
	if !c.Is(state.ChangedPresets) {
		return 0
	}

	presets := s.Presets
	for i := range presets {
		p, err := p.store.Get(ctx, presets[i].URL)
		if err != nil {
			continue
		}

		presets[i].TitleNew = p.TitleNew
		presets[i].URLNew = p.URLNew
	}

	return s.SetPresets(presets)
}

func (p *Preset) fragmentTitleAndURL(ctx context.Context, s *state.State, c state.Changed) state.Changed {
	if !c.Is(state.ChangedURL) {
		return 0
	}

	preset, err := p.store.Get(ctx, s.URL)
	if err != nil {
		return 0
	}

	return s.SetTitleNew(preset.TitleNew).Merge(s.SetURLNew(preset.URLNew))
}
