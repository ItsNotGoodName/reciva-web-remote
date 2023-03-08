package middleware

import (
	"context"

	"github.com/ItsNotGoodName/reciva-web-remote/internal/state"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/store"
)

type StateHook struct {
	store store.Store
}

func NewStateHook(store store.Store) StateHook {
	return StateHook{store: store}
}

func (sh StateHook) OnStart(ctx context.Context, s *state.State, c state.Changed) state.Changed {
	return sh.OnChanged(ctx, s, c)
}

func (sh StateHook) OnChanged(ctx context.Context, s *state.State, c state.Changed) state.Changed {
	return c.Merge(sh.applyPresets(ctx, s, c)).Merge(sh.applyTitleAndURL(ctx, s, c))
}

func (sh StateHook) applyPresets(ctx context.Context, s *state.State, c state.Changed) state.Changed {
	if !c.Is(state.ChangedPresets) {
		return 0
	}

	presets := make([]state.Preset, 0, len(s.Presets))

	for i := range s.Presets {
		presets = append(presets, s.Presets[i])

		p, err := sh.store.GetPreset(ctx, presets[i].URL)
		if err != nil {
			continue
		}

		presets[i].TitleNew = p.TitleNew
		presets[i].URLNew = p.URLNew
	}

	return s.SetPresets(presets)
}

func (sh StateHook) applyTitleAndURL(ctx context.Context, s *state.State, c state.Changed) state.Changed {
	if !c.Is(state.ChangedURL) {
		return 0
	}

	preset, err := sh.store.GetPreset(ctx, s.URL)
	if err != nil {
		return 0
	}

	return s.SetTitleNew(preset.TitleNew).Merge(s.SetURLNew(preset.URLNew))
}
