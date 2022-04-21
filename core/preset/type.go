package preset

import (
	"context"
	"fmt"
	"net/url"
)

var (
	ErrPresetNotFound = fmt.Errorf("preset not found")
)

type (
	Preset struct {
		URL      string
		NewTitle string
		NewURL   string
	}

	PresetStore interface {
		List(ctx context.Context) ([]Preset, error)
		Get(ctx context.Context, url string) (*Preset, error)
		Update(ctx context.Context, preset *Preset) error
	}
)

func NewPreset(urL, newTitle, newURL string) (*Preset, error) {
	urlObj, err := url.ParseRequestURI(urL)
	if err != nil {
		return nil, err
	}

	if urlObj.Scheme == "" || urlObj.Host == "" {
		return nil, fmt.Errorf("invalid URL: %s", urL)
	}

	return &Preset{
		URL:      urL,
		NewTitle: newTitle,
		NewURL:   newURL,
	}, nil
}
