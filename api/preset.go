package api

import (
	"github.com/ItsNotGoodName/reciva-web-remote/store"
)

type PresetAPI struct {
	a *API
	S *store.Store
}

func NewPresetAPI(a *API, s *store.Store) *PresetAPI {
	return &PresetAPI{a: a, S: s}
}

func (p *PresetAPI) GetPresets() *map[string]int {
	pts := p.S.GetPresets()
	ptsMap := make(map[string]int, len(pts))
	for _, v := range pts {
		ptsMap[v.URI] = v.StreamID
	}
	return &ptsMap
}

func (p *PresetAPI) GetStreams() *map[int]string {
	sts := p.S.GetStreams()
	stsMap := make(map[int]string, len(sts))
	for _, v := range sts {
		stsMap[v.SID] = v.Name
	}
	return &stsMap
}

func (p *PresetAPI) DeleteStream(sid int) bool {
	return p.S.DeleteStream(sid) > 0
}

type StreamReq struct {
	Name    *string `json:"name,omitempty"`
	Content *string `json:"content,omitempty"`
}

func (p *PresetAPI) UpdateStream(st *store.Stream, r *StreamReq) bool {
	update := false

	if r.Name != nil && *r.Name != st.Name {
		st.Name = *r.Name
		update = true
	}
	if r.Content != nil && *r.Content != st.Content {
		st.Content = *r.Content
		update = true
	}
	if update {
		if !p.S.UpdateStream(st) {
			return false
		}
	}

	return true
}

type PresetReq struct {
	URI *string `json:"uri,omitempty"`
	SID *int    `json:"sid,omitempty"`
}

func (p *PresetAPI) UpdatePreset(pt *store.Preset, r *PresetReq) bool {
	pt.StreamID = *r.SID
	return p.S.UpdatePreset(pt)
}
