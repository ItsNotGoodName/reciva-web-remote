package store

import (
	"context"
)

// SaveConfig saves to the config file.
func (s *Store) SaveConfig(ctx context.Context) error {
	if s.readonly {
		return ErrReadOnly
	}

	errChan := make(chan error)

	s.configOp <- func(m map[string]Preset) map[string]Preset {
		err := s.saveConfig(s.configFile, m)
		select {
		case errChan <- err:
		case <-ctx.Done():
			errChan <- ctx.Err()
		}
		return m
	}

	return <-errChan
}

// ReadConfig reads the config file.
func (s *Store) ReadConfig(ctx context.Context) error {
	errChan := make(chan error)

	s.configOp <- func(m map[string]Preset) map[string]Preset {
		res, err := readConfig(s.configFile)
		if err != nil {
			select {
			case errChan <- err:
			case <-ctx.Done():
				errChan <- ctx.Err()
			}
			return m
		}

		select {
		case errChan <- nil:
		case <-ctx.Done():
			errChan <- ctx.Err()
		}
		return res
	}

	return <-errChan
}
