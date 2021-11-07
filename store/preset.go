package store

import (
	"context"
)

type Preset struct {
	URL string `json:"url"` // URL is the URL of the preset, ex. 'http://example.com/01.m3u'.
	SID int    `json:"sid"` // SID is the stream ID of the preset, 0 means there is no stream associated with it.
}

// AddPreset adds preset with context.
func (s *Store) AddPreset(ctx context.Context, preset *Preset) error {
	_, err := s.db.ExecContext(ctx, "INSERT INTO preset (url, sid) VALUES ($1, $2)", preset.URL, preset.SID)
	return err
}

// GetPresets returns all presets with context.
func (s *Store) GetPresets(ctx context.Context) ([]*Preset, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT url, sid FROM preset")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var presets []*Preset
	for rows.Next() {
		var preset Preset
		if err := rows.Scan(&preset.URL, &preset.SID); err != nil {
			return nil, err
		}
		presets = append(presets, &preset)
	}
	return presets, nil
}

// ClearPreset sets preset's SID to 0 with context.
func (s *Store) ClearPreset(ctx context.Context, preset *Preset) error {
	_, err := s.db.ExecContext(ctx, "UPDATE preset SET sid = 0 WHERE URL = $1", preset.URL)
	preset.SID = 0
	return err
}

// UpdatePreset updates preset's SID with context.
func (s *Store) UpdatePreset(ctx context.Context, preset *Preset) (bool, error) {
	// Update preset's SID
	result, err := s.db.ExecContext(ctx, "UPDATE preset SET sid = $1 WHERE URL = $2", preset.SID, preset.URL)
	if err != nil {
		return false, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rows > 0, nil
}

// GetPreset returns preset by URL with context.
func (s *Store) GetPreset(ctx context.Context, url string) (*Preset, error) {
	var preset Preset
	err := s.db.QueryRowContext(ctx, "SELECT url, sid FROM preset WHERE url = $1", url).Scan(&preset.URL, &preset.SID)
	if err != nil {
		return nil, err
	}
	return &preset, nil
}

// DeleteAllPresets deletes all presets with context.
func (s *Store) DeleteAllPresets(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM preset")
	return err
}
