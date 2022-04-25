package mock

import (
	"context"

	"github.com/ItsNotGoodName/reciva-web-remote/core"
	"github.com/ItsNotGoodName/reciva-web-remote/core/preset"
)

type PresetStore struct{}

func NewPresetStore() *PresetStore {
	return &PresetStore{}
}

func (PresetStore) List(context.Context) ([]preset.Preset, error) {
	return []preset.Preset{}, nil
}

func (PresetStore) Get(context.Context, string) (*preset.Preset, error) {
	return nil, core.ErrPresetNotFound
}

func (PresetStore) Update(context.Context, *preset.Preset) error {
	return core.ErrPresetNotFound
}
