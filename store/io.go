package store

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/ItsNotGoodName/reciva-web-remote/core"
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

// readPresets reads the config file.
func readPresets(path string) (map[string]core.Preset, error) {
	// Read config file
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Unmarshal config file
	config := core.Config{}
	if err := json.Unmarshal(b, &config); err != nil {
		return nil, err
	}

	// Create map from config
	presets := make(map[string]core.Preset, len(config.Presets))
	for _, p := range config.Presets {
		if err := validateURL(p.URL); err != nil {
			log.Fatal("store.readConfig(ERROR):", err)
		}
		presets[p.URL] = p
	}

	return presets, nil
}

// saveConfig writes the config file.
func saveConfig(configFile string, m map[string]core.Preset) error {
	// Create presets slice from map
	presets := make([]core.Preset, 0, len(m))
	for _, p := range m {
		presets = append(presets, p)
	}

	// Create config struct
	js := core.Config{
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
