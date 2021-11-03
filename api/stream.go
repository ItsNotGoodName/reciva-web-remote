package api

import (
	"context"
	"database/sql"

	"github.com/ItsNotGoodName/reciva-web-remote/store"
	"github.com/mattn/go-sqlite3"
)

// GetStream returns a stream by id.
func (p *PresetAPI) GetStream(ctx context.Context, id int) (*store.Stream, error) {
	stream, err := p.s.GetStream(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrStreamNotFound
		}
		return nil, err
	}

	return stream, nil
}

// GetStreamByURI returns a stream by URI.
func (p *PresetAPI) GetStreamByURI(ctx context.Context, uri string) (*store.Stream, error) {
	preset, err := p.GetPreset(ctx, uri)
	if err != nil {
		return nil, err
	}

	stream, err := p.GetStream(ctx, preset.SID)
	if err != nil {
		return nil, err
	}

	return stream, nil
}

// GetStreams returns a list of streams.
func (p *PresetAPI) GetStreams(ctx context.Context) ([]*store.Stream, error) {
	streams, err := p.s.GetStreams(ctx)
	if err != nil {
		return nil, err
	}
	if len(streams) == 0 {
		return []*store.Stream{}, nil
	}

	return streams, nil
}

// AddStreamRequest is a request for AddStream.
type AddStreamRequest struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

// AddStream adds a stream.
func (p *PresetAPI) AddStream(ctx context.Context, req *AddStreamRequest) (*store.Stream, error) {
	stream := &store.Stream{
		Name:    req.Name,
		Content: req.Content,
	}
	err := p.s.AddStream(ctx, stream)
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

// UpdateStreamRequest is a request for UpdateStream.
type UpdateStreamRequest struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

// UpdateStream updates a stream.
func (p *PresetAPI) UpdateStream(ctx context.Context, req *UpdateStreamRequest) (*store.Stream, error) {
	stream := &store.Stream{
		ID:      req.ID,
		Name:    req.Name,
		Content: req.Content,
	}
	ok, err := p.s.UpdateStream(ctx, stream)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrStreamNotFound
	}
	return stream, nil
}

// DeleteStream deletes a stream by ID.
func (p *PresetAPI) DeleteStream(ctx context.Context, id int) error {
	ok, err := p.s.DeleteStream(ctx, &store.Stream{ID: id})
	if err != nil {
		return err
	}
	if !ok {
		return ErrStreamNotFound
	}
	return nil
}
