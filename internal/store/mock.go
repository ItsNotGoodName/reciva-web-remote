package store

import (
	"context"

	"github.com/ItsNotGoodName/reciva-web-remote/core"
	"github.com/ItsNotGoodName/reciva-web-remote/core/preset"
)

type Mock struct{}

func NewMock() Mock {
	return Mock{}
}

func (Mock) List(context.Context) ([]preset.Preset, error) {
	return []preset.Preset{}, nil
}

func (Mock) Get(context.Context, string) (*preset.Preset, error) {
	return nil, core.ErrPresetNotFound
}

func (Mock) Update(context.Context, *preset.Preset) error {
	return core.ErrPresetNotFound
}
