package store

import (
	"context"

	"github.com/ItsNotGoodName/reciva-web-remote/internal/model"
)

type Store interface {
	ListPresets(ctx context.Context) ([]model.Preset, error)
	GetPreset(ctx context.Context, url string) (*model.Preset, error)
	UpdatePreset(ctx context.Context, p *model.Preset) error
}
