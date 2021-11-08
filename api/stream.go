package api

import (
	"context"

	"github.com/ItsNotGoodName/reciva-web-remote/store"
)

// ReadStream returns a stream by id.
func (p *PresetAPI) ReadStream(ctx context.Context, id int) (*store.Stream, error) {
	stream, err := p.s.ReadStream(ctx, id)
	if err != nil {
		if err == store.ErrNotFound {
			return nil, ErrStreamNotFound
		}
		return nil, err
	}

	return stream, nil
}

// ReadStreamByURL returns a stream by URL.
func (p *PresetAPI) ReadStreamByURL(ctx context.Context, url string) (*store.Stream, error) {
	preset, err := p.ReadPreset(ctx, url)
	if err != nil {
		return nil, err
	}

	stream, err := p.ReadStream(ctx, preset.SID)
	if err != nil {
		return nil, err
	}

	return stream, nil
}

// ReadStreams returns a list of streams.
func (p *PresetAPI) ReadStreams(ctx context.Context) ([]*store.Stream, error) {
	streams, err := p.s.ReadStreams(ctx)
	if err != nil {
		return nil, err
	}
	if len(streams) == 0 {
		return []*store.Stream{}, nil
	}

	return streams, nil
}

// CreateStreamRequest is a request for CreateStream.
type CreateStreamRequest struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

// CreateStream creates a stream.
func (p *PresetAPI) CreateStream(ctx context.Context, req *CreateStreamRequest) (*store.Stream, error) {
	// Validate request
	if req.Name == "" || len(req.Name) > store.StreamNameMaxLength {
		return nil, ErrStreamNameInvalid
	}
	if len(req.Content) > store.StreamContentMaxLength {
		return nil, ErrStreamContentInvalid
	}

	// Create stream
	stream := &store.Stream{
		Name:    req.Name,
		Content: req.Content,
	}
	err := p.s.CreateStream(ctx, stream)
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

// UpdateStream updates a stream.
func (p *PresetAPI) UpdateStream(ctx context.Context, req *UpdateStreamRequest) (*store.Stream, error) {
	// Validate request
	if req.Name == "" || len(req.Name) > store.StreamNameMaxLength {
		return nil, ErrStreamNameInvalid
	}
	if len(req.Content) > store.StreamContentMaxLength {
		return nil, ErrStreamContentInvalid
	}

	// Update stream
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
