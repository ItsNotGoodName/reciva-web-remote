package store

func NewPresets(uris []string) []Preset {
	pts := make([]Preset, 0, len(uris))
	for _, u := range uris {
		pts = append(pts, Preset{URI: u})
	}
	return pts
}
