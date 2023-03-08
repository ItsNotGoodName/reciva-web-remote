package store

import (
	"context"

	"github.com/ItsNotGoodName/reciva-web-remote/internal"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/model"
)

type Mock struct{}

func NewMock() Mock {
	return Mock{}
}

func (Mock) List(context.Context) ([]model.Preset, error) {
	return []model.Preset{}, nil
}

func (Mock) Get(context.Context, string) (*model.Preset, error) {
	return nil, internal.ErrPresetNotFound
}

func (Mock) Update(context.Context, *model.Preset) error {
	return internal.ErrPresetNotFound
}
