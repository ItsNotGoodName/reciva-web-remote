package store

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
)

func validateURL(u string) error {
	uu, err := url.ParseRequestURI(u)
	if err != nil {
		return err
	}

	if uu.Scheme == "" || uu.Host == "" {
		return fmt.Errorf("invalid URL: %s", u)
	}

	return nil
}

// readConfig reads the config file.
func readConfig(path string) (map[string]Preset, error) {
	// Read config file
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Unmarshal config file
	config := Config{}
	if err := json.Unmarshal(b, &config); err != nil {
		return nil, err
	}

	if len(config.Presets) == 0 {
		return nil, ErrEmptyPresets
	}

	// Create map from config
	presets := make(map[string]Preset, len(config.Presets))
	for _, p := range config.Presets {
		if err := validateURL(p.URL); err != nil {
			log.Fatal("store.readConfig(ERROR):", err)
		}
		presets[p.URL] = p
	}

	return presets, nil
}

// writeConfig writes the config file.
func (s *Store) saveConfig(configFile string, m map[string]Preset) error {
	// Create presets slice from map
	presets := make([]Preset, 0, len(m))
	for _, p := range m {
		presets = append(presets, p)
	}

	// Create config struct
	js := Config{
		Presets: presets,
	}

	// Marshal config struct
	b, err := json.MarshalIndent(js, "", "	")
	if err != nil {
		return err
	}

	// Write config file
	err = os.WriteFile(configFile, b, 0644)
	if err != nil {
		return err
	}

	return nil
}
