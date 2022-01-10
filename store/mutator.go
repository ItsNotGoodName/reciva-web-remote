package store

import (
	"context"

	"github.com/ItsNotGoodName/reciva-web-remote/core"
	"github.com/ItsNotGoodName/reciva-web-remote/pkg/radio"
)

type Mutator struct {
	store core.PresetStore
}

func NewMutator(store core.PresetStore) *Mutator {
	return &Mutator{
		store: store,
	}
}

func (m Mutator) Mutate(ctx context.Context, state *radio.State) *radio.State {
	presets := state.Presets

	for i := range presets {
		preset, err := m.store.GetPreset(ctx, presets[i].URL)
		if err != nil || preset.NewName == "" {
			presets[i].Name = presets[i].Title
		} else {
			presets[i].Name = preset.NewName
		}
	}

	retState := m.MutateNewURL(ctx, state)
	retState.Presets = presets
	return retState
}

func (m Mutator) MutateNewURL(ctx context.Context, state *radio.State) *radio.State {
	preset, err := m.store.GetPreset(ctx, state.URL)
	if err != nil {
		return &radio.State{}
	}

	state.NewURL = &preset.NewURL

	return &radio.State{NewURL: state.NewURL}
}

func (m Mutator) GetTrigger() <-chan struct{} {
	return m.store.PresetChanged()
}
