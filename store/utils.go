package store

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
