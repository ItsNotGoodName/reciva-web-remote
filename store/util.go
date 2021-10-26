package store

func getPresetIndex(pts []Preset, uri string) int {
	for i := range pts {
		if pts[i].URI == uri {
			return i
		}
	}
	return -1
}
