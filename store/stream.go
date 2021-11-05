package store

import "context"

type Stream struct {
	ID      int    `json:"id"`      // ID is the unique ID of the stream.
	Name    string `json:"name"`    // Name of the stream.
	Content string `json:"content"` // Content of the stream.
}

const (
	StreamNameMaxLength    = 64
	StreamContentMaxLength = 1024
)

// AddStream adds stream to the store with context.
func (s *Store) AddStream(ctx context.Context, stream *Stream) error {
	return s.db.QueryRowContext(ctx, "INSERT INTO stream (name, content) VALUES ($1, $2) RETURNING id", stream.Name, stream.Content).Scan(&stream.ID)
}

// GetStreams returns all streams in the store with context.
func (s *Store) GetStreams(ctx context.Context) ([]*Stream, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id, name, content FROM stream")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var streams []*Stream
	for rows.Next() {
		var stream Stream
		if err := rows.Scan(&stream.ID, &stream.Name, &stream.Content); err != nil {
			return nil, err
		}
		streams = append(streams, &stream)
	}

	return streams, nil
}

// GetStream returns stream by ID with context.
func (s *Store) GetStream(ctx context.Context, id int) (*Stream, error) {
	var stream Stream
	err := s.db.QueryRowContext(ctx, "SELECT id, name, content FROM stream WHERE id = $1", id).Scan(&stream.ID, &stream.Name, &stream.Content)
	if err != nil {
		return nil, err
	}

	return &stream, nil
}

// DeleteStream deletes stream with context.
func (s *Store) DeleteStream(ctx context.Context, stream *Stream) (bool, error) {
	txn, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return false, err
	}

	defer txn.Rollback()

	_, err = txn.ExecContext(ctx, "UPDATE preset SET sid = 0 WHERE sid = $1", stream.ID)
	if err != nil {
		return false, err
	}

	result, err := txn.ExecContext(ctx, "DELETE FROM stream WHERE id = $1", stream.ID)
	if err != nil {
		return false, err
	}

	if txn.Commit() != nil {
		return false, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rows > 0, nil
}

// UpdateStream updates stream with context.
func (s *Store) UpdateStream(ctx context.Context, stream *Stream) (bool, error) {
	result, err := s.db.ExecContext(ctx, "UPDATE stream SET name = $1, content = $2 WHERE id = $3", stream.Name, stream.Content, stream.ID)
	if err != nil {
		return false, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rows > 0, nil
}
