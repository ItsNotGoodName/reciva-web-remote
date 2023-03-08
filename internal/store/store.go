package store

import (
	"context"
	"log"

	"github.com/ItsNotGoodName/reciva-web-remote/internal/model"
)

type Store interface {
	ListPresets(ctx context.Context) ([]model.Preset, error)
	GetPreset(ctx context.Context, url string) (*model.Preset, error)
	UpdatePreset(ctx context.Context, p *model.Preset) error
}

func Must(store Store, err error) Store {
	if err != nil {
		log.Fatalln("store.Must:", err)
	}

	return store
}
