package store

import (
	"context"
)

// ReadPreset returns a preset by its URL.
func (s *Store) ReadPreset(ctx context.Context, url string) (*Preset, error) {
	pChan := make(chan Preset)
	errChan := make(chan error)

	s.presetOp <- func(m map[string]Preset) {
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
func (s *Store) ReadPresets(ctx context.Context) ([]Preset, error) {
	pChan := make(chan []Preset)

	s.presetOp <- func(m map[string]Preset) {
		presets := make([]Preset, 0, len(m))
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
func (s *Store) UpdatePreset(ctx context.Context, preset *Preset) error {
	if s.readonly {
		return ErrReadOnly
	}

	errChan := make(chan error)

	s.presetOp <- func(m map[string]Preset) {
		_, ok := m[preset.URL]
		if !ok {
			select {
			case errChan <- ErrNotFound:
			case <-ctx.Done():
				errChan <- ctx.Err()
			}
			return
		}
		m[preset.URL] = *preset
		err := s.saveConfig(s.configFile, m)
		select {
		case errChan <- err:
		case <-ctx.Done():
			errChan <- ctx.Err()
		}
	}

	return <-errChan
}
