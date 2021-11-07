package api

import (
	"context"
	"database/sql"

	"github.com/ItsNotGoodName/reciva-web-remote/store"
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

// GetStreamByURL returns a stream by URL.
func (p *PresetAPI) GetStreamByURL(ctx context.Context, url string) (*store.Stream, error) {
	preset, err := p.GetPreset(ctx, url)
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

func (r *AddStreamRequest) Validate() error {
	if r.Name == "" || len(r.Name) > store.StreamNameMaxLength {
		return ErrStreamNameInvalid
	}
	if len(r.Content) > store.StreamContentMaxLength {
		return ErrStreamContentInvalid
	}
	return nil
}

// AddStream adds a stream.
func (p *PresetAPI) AddStream(ctx context.Context, req *AddStreamRequest) (*store.Stream, error) {
	stream := &store.Stream{
		Name:    req.Name,
		Content: req.Content,
	}
	err := p.s.AddStream(ctx, stream)
	if err != nil {
		return nil, ErrNameAlreadyExists
	}
	return stream, nil
}

// UpdateStreamRequest is a request for UpdateStream.
type UpdateStreamRequest struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

func (r *UpdateStreamRequest) Validate() error {
	if r.Name == "" || len(r.Name) > store.StreamNameMaxLength {
		return ErrStreamNameInvalid
	}
	if len(r.Content) > store.StreamContentMaxLength {
		return ErrStreamContentInvalid
	}
	return nil
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
	p.h.RefreshPresets()
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
	p.h.RefreshPresets()
	return nil
}
