package store

import (
	"github.com/ItsNotGoodName/reciva-web-remote/config"
	"github.com/ItsNotGoodName/reciva-web-remote/pkg/goupnpsub"
)

func NewSettings() *Settings {
	return &Settings{
		Port:    config.DefaultPort,
		CPort:   goupnpsub.DefaultPort,
		Streams: make([]Stream, 0),
		Presets: make([]Preset, 0),
	}
}

func (st *Settings) mergePresets(p []Preset) {
	var newPresets []Preset

	for _, cfgPreset := range p {
		found := false
		for i := range st.Presets {
			if st.Presets[i].URI == cfgPreset.URI {
				newPresets = append(newPresets, st.Presets[i])
				found = true
				break
			}
		}
		if !found {
			newPresets = append(newPresets, cfgPreset)
		}
	}
	st.Presets = newPresets
}
