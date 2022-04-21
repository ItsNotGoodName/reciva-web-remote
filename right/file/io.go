package file

import (
	"encoding/json"
	"os"

	"github.com/ItsNotGoodName/reciva-web-remote/core/preset"
)

type Config struct {
	Presets []Preset `json:"presets"`
}

type Preset struct {
	URL     string `json:"url"`
	NewName string `json:"newName"`
	NewURL  string `json:"newUrl"`
}

func readConfig(file string) (map[string]preset.Preset, error) {
	// Read config file
	b, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	// Unmarshal cfg file
	cfg := Config{}
	if err := json.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}

	// Create map from config
	presets := make(map[string]preset.Preset, len(cfg.Presets))
	for _, p := range cfg.Presets {
		pp, err := preset.NewPreset(p.URL, p.NewName, p.NewURL)
		if err != nil {
			return nil, err
		}

		presets[p.URL] = *pp
	}

	return presets, nil
}

func writeConfig(file string, m map[string]preset.Preset) error {
	// Create presets slice from map
	presets := make([]Preset, 0, len(m))
	for _, p := range m {
		presets = append(presets, Preset{
			URL:     p.URL,
			NewName: p.NewTitle,
			NewURL:  p.NewURL,
		})
	}

	// Create config struct
	cfg := Config{
		Presets: presets,
	}

	// Marshal config struct
	b, err := json.MarshalIndent(cfg, "", "	")
	if err != nil {
		return err
	}

	// Write config file
	return os.WriteFile(file, b, 0644)
}
