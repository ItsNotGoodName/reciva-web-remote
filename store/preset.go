package store

import (
	"context"
)

type Preset struct {
	ID  int    // ID is the unique ID of the preset.
	URI string // URI is the URI of the preset, ex. '/01.m3u'.
	SID int    // SID is the stream ID of the preset, 0 means there is no stream associated with it.
}

// AddPreset adds preset with context.
func (s *Store) AddPreset(ctx context.Context, preset *Preset) error {
	return s.db.QueryRowContext(ctx, "INSERT INTO preset (uri, sid) VALUES ($1, $2) RETURNING id", preset.URI, preset.SID).Scan(&preset.ID)
}

// GetPresets returns all presets with context.
func (s *Store) GetPresets(ctx context.Context) ([]*Preset, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id, uri, sid FROM preset")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var presets []*Preset
	for rows.Next() {
		var preset Preset
		if err := rows.Scan(&preset.ID, &preset.URI, &preset.SID); err != nil {
			return nil, err
		}
		presets = append(presets, &preset)
	}
	return presets, nil
}

// ClearPreset sets preset's SID to 0 with context.
func (s *Store) ClearPreset(ctx context.Context, preset *Preset) error {
	_, err := s.db.ExecContext(ctx, "UPDATE preset SET sid = 0 WHERE id = $1", preset.ID)
	preset.SID = 0
	return err
}

// UpdatePresetSID updates preset's SID with context.
func (s *Store) UpdatePresetSID(ctx context.Context, preset *Preset) error {
	_, err := s.db.ExecContext(ctx, "UPDATE preset SET sid = $1 WHERE id = $2 AND (SELECT COUNT(*) FROM stream WHERE id = $1) > 0 ", preset.SID, preset.ID)
	return err
}

// GetPresetByURI returns preset by uri with context.
func (s *Store) GetPresetByURI(ctx context.Context, uri string) (*Preset, error) {
	var preset Preset
	err := s.db.QueryRowContext(ctx, "SELECT id, uri, sid FROM preset WHERE uri = $1", uri).Scan(&preset.ID, &preset.URI, &preset.SID)
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
