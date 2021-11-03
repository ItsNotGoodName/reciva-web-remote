package api

import "github.com/ItsNotGoodName/reciva-web-remote/store"

func NewPresetAPI(s *store.Store) *PresetAPI {
	return &PresetAPI{s: s}
}
