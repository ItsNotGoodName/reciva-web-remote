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

func (p *Preset) Apply(frag *state.Fragment) {
	ctx := context.Background()

	p.fragmentPresets(ctx, frag)

	p.fragmentTitleAndURL(ctx, frag)
}

func (p *Preset) fragmentPresets(ctx context.Context, frag *state.Fragment) {
	if frag.Presets != nil {
		for i := range frag.Presets {
			titleNew := ""
			urlNew := ""

			preset, err := p.store.Get(ctx, frag.Presets[i].URL)
			if err == nil {
				titleNew = preset.TitleNew
				urlNew = preset.URLNew
			}

			frag.Presets[i].TitleNew = titleNew
			frag.Presets[i].URLNew = urlNew
		}
	}
}

func (p *Preset) fragmentTitleAndURL(ctx context.Context, frag *state.Fragment) {
	if frag.URL != nil {
		urlNew := ""
		titleNew := ""

		preset, err := p.store.Get(ctx, *frag.URL)
		if err == nil {
			urlNew = preset.URLNew
			titleNew = preset.TitleNew
		}

		frag.TitleNew = &titleNew
		frag.URLNew = &urlNew
	}
}
