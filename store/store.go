package store

import (
	"github.com/ItsNotGoodName/reciva-web-remote/config"
)

type Store struct {
	op  chan func(map[string]config.Preset)
	cfg *config.Config
}

func NewStore(cfg *config.Config) (*Store, error) {
	s := Store{op: make(chan func(map[string]config.Preset)), cfg: cfg}

	// Create presets map
	presets := make(map[string]config.Preset)
	for _, p := range cfg.Presets {
		presets[p.URL] = p
	}
	go s.StoreLoop(presets)

	return &s, nil
}

func (s *Store) StoreLoop(presets map[string]config.Preset) {
	for f := range s.op {
		f(presets)
	}
}
