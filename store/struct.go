package store

type Config struct {
	Presets []Preset `json:"presets"`
}

type Preset struct {
	URL     string `json:"url"`
	NewName string `json:"newName"`
	NewURL  string `json:"newUrl"`
}
type Store struct {
	readonly   bool
	presetOp   chan func(map[string]Preset)
	configOp   chan func(map[string]Preset) map[string]Preset
	configFile string
}
