package api

import (
	"github.com/ItsNotGoodName/reciva-web-remote/store"
)

type PresetAPI struct {
	a *API
	s *store.Store
}

func NewPresetAPI(a *API, s *store.Store) *PresetAPI {
	return &PresetAPI{a: a, s: s}
}
