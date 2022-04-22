package middleware

import (
	"context"

	"github.com/ItsNotGoodName/reciva-web-remote/core/preset"
	"github.com/ItsNotGoodName/reciva-web-remote/core/state"
	"github.com/ItsNotGoodName/reciva-web-remote/pkg/sig"
)

type Preset struct {
	pub   *sig.Pub
	store preset.PresetStore
}

func NewPreset(pub *sig.Pub, presetStore preset.PresetStore) *Preset {
	return &Preset{
		pub:   pub,
		store: presetStore,
	}
}

func (m *Preset) List(ctx context.Context) ([]preset.Preset, error) {
	return m.store.List(ctx)
}

func (m *Preset) Get(ctx context.Context, url string) (*preset.Preset, error) {
	return m.store.Get(ctx, url)
}

func (m *Preset) Update(ctx context.Context, preset *preset.Preset) error {
	err := m.store.Update(ctx, preset)
	if err != nil {
		return err
	}

	m.pub.Publish()
	return nil
}

func (m *Preset) Fragment(frag *state.Fragment) {
	ctx := context.Background()

	if frag.Presets != nil {
		m.fragmentPresets(ctx, frag)
	}

	if frag.URL != nil || frag.Title != nil {
		m.fragmentTitleAndURL(ctx, frag)
	}
}

func (m *Preset) FragmentFromState(s state.State) state.Fragment {
	title := s.Title
	url := s.URL
	var presets []state.Preset
	copy(presets, s.Presets)

	frag := state.NewFragment(s.UUID)
	frag.Title = &title
	frag.URL = &url
	frag.Presets = presets

	ctx := context.Background()
	m.fragmentPresets(ctx, &frag)
	m.fragmentTitleAndURL(ctx, &frag)

	return frag
}

func (m *Preset) fragmentPresets(ctx context.Context, frag *state.Fragment) {
	for i := range frag.Presets {
		newTitle := frag.Presets[i].Title

		preset, err := m.store.Get(ctx, frag.Presets[i].URL)
		if err == nil && preset.NewTitle != "" {
			newTitle = preset.NewTitle
		}

		frag.Presets[i].NewTitle = newTitle
	}
}

func (m *Preset) fragmentTitleAndURL(ctx context.Context, frag *state.Fragment) {
	newURL := *frag.URL
	newTitle := *frag.Title

	preset, err := m.store.Get(ctx, *frag.URL)
	if err == nil {
		newURL = preset.NewURL
		newTitle = preset.NewTitle
	}

	frag.NewTitle = &newTitle
	frag.URL = &newURL
}
