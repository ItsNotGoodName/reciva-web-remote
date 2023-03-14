package middleware

import (
	"context"

	"github.com/ItsNotGoodName/reciva-web-remote/internal/model"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/pubsub"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/state"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/store"
)

type Store struct {
	store store.Store
}

func NewStore(store store.Store) store.Store {
	return Store{store: store}
}

func (s Store) ListPresets(ctx context.Context) ([]model.Preset, error) {
	return s.store.ListPresets(ctx)
}
func (s Store) GetPreset(ctx context.Context, url string) (*model.Preset, error) {
	return s.store.GetPreset(ctx, url)
}

func (s Store) UpdatePreset(ctx context.Context, p *model.Preset) error {
	err := s.store.UpdatePreset(ctx, p)
	if err != nil {
		return err
	}

	pubsub.DefaultPub.Publish(pubsub.StateHookStaleTopic, pubsub.StateHookStaleMessage{Changed: state.ChangedAll})
	return nil
}
