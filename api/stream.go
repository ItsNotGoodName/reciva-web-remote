package api

import (
	"context"
	"database/sql"

	"github.com/ItsNotGoodName/reciva-web-remote/store"
	"github.com/mattn/go-sqlite3"
)

// GetStream returns a stream from store by id with context.
func (a *PresetAPI) GetStream(ctx context.Context, id int) (*store.Stream, error) {
	stream, err := a.s.GetStreamByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrStreamNotFound
		}
		return nil, err
	}
	return stream, nil
}

// GetStreams returns a list of streams from store with context.
func (a *PresetAPI) GetStreams(ctx context.Context) ([]*store.Stream, error) {
	streams, err := a.s.GetStreams(ctx)
	if err != nil {

		return nil, err
	}
	if len(streams) == 0 {
		return []*store.Stream{}, nil
	}
	return streams, nil
}

// AddStreamRequest is a request to add a stream to store.
type AddStreamRequest struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

// AddStream adds a stream to store with context.
func (a *PresetAPI) AddStream(ctx context.Context, req *AddStreamRequest) (*store.Stream, error) {
	stream := &store.Stream{
		Name:    req.Name,
		Content: req.Content,
	}
	err := a.s.AddStream(ctx, stream)
	if err != nil {
		if sqlErr, ok := err.(sqlite3.Error); ok {
			if sqlErr.Code == sqlite3.ErrConstraint {
				return nil, ErrNameAlreadyExists
			}
		}
		return nil, err
	}
	return stream, nil
}

// UpdateStreamRequest is a request to update a stream in store.
type UpdateStreamRequest struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

// UpdateStream updates a stream in store with context.
func (a *PresetAPI) UpdateStream(ctx context.Context, req *UpdateStreamRequest) (*store.Stream, error) {
	stream := &store.Stream{
		ID:      req.ID,
		Name:    req.Name,
		Content: req.Content,
	}
	ok, err := a.s.UpdateStream(ctx, stream)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrStreamNotFound
	}
	return stream, nil
}

// DeleteStream deletes a stream from store with context.
func (a *PresetAPI) DeleteStream(ctx context.Context, id int) error {
	ok, err := a.s.DeleteStream(ctx, &store.Stream{ID: id})
	if err != nil {
		return err
	}
	if !ok {
		return ErrStreamNotFound
	}
	return nil
}
