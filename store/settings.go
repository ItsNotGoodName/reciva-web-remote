package store

func NewSettings() *Settings {
	return &Settings{
		Port:    8080,
		CPort:   8058,
		Streams: make([]Stream, 0),
		Presets: make([]Preset, 0),
	}
}
