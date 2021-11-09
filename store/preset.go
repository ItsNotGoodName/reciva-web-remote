package store

import (
	"context"
	"encoding/json"
	"os"

	"github.com/ItsNotGoodName/reciva-web-remote/config"
)

// ReadPresetByURL returns a preset by its URL.
func (s *Store) ReadPresetByURL(ctx context.Context, url string) (*config.Preset, error) {
	pChan := make(chan config.Preset)
	errChan := make(chan error)

	s.op <- func(m map[string]config.Preset) {
		preset, ok := m[url]
		if !ok {
			select {
			case errChan <- ErrNotFound:
			case <-ctx.Done():
				errChan <- ctx.Err()
			}
			return
		}
		select {
		case <-ctx.Done():
			errChan <- ctx.Err()
		case pChan <- preset:
		}
	}

	select {
	case err := <-errChan:
		return nil, err
	case preset := <-pChan:
		return &preset, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// ReadPresets returns all presets.
func (s *Store) ReadPresets(ctx context.Context) ([]config.Preset, error) {
	pChan := make(chan []config.Preset)

	s.op <- func(m map[string]config.Preset) {
		presets := make([]config.Preset, 0, len(m))
		for _, preset := range m {
			presets = append(presets, preset)
		}
		select {
		case pChan <- presets:
		case <-ctx.Done():
		}
	}

	select {
	case presets := <-pChan:
		return presets, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// UpdatePreset updates a preset.
func (s *Store) UpdatePreset(ctx context.Context, preset config.Preset) error {
	errChan := make(chan error)

	s.op <- func(m map[string]config.Preset) {
		_, ok := m[preset.URL]
		if !ok {
			select {
			case errChan <- ErrNotFound:
			case <-ctx.Done():
				errChan <- ctx.Err()
			}
			return
		}
		m[preset.URL] = preset
		select {
		case errChan <- nil:
		case <-ctx.Done():
			errChan <- ctx.Err()
		}
	}

	return <-errChan
}

func (s *Store) SaveConfig(ctx context.Context) error {
	errChan := make(chan error)

	s.op <- func(m map[string]config.Preset) {
		err := s.saveConfig(m)
		select {
		case errChan <- err:
		case <-ctx.Done():
			errChan <- ctx.Err()
		}
	}

	return <-errChan
}

func (s *Store) saveConfig(m map[string]config.Preset) error {
	presets := make([]config.Preset, 0, len(m))
	for _, p := range m {
		presets = append(presets, p)
	}

	js := config.ConfigJSON{
		Port:    s.cfg.Port,
		CPort:   s.cfg.CPort,
		Presets: presets,
	}
	b, err := json.MarshalIndent(js, "", "	")
	if err != nil {
		return err
	}

	err = os.WriteFile(s.cfg.ConfigFile, b, 0644)
	if err != nil {
		return err
	}

	return nil
}
