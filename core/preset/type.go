package preset

import (
	"context"
	"fmt"
	"net/url"
)

type (
	Preset struct {
		TitleNew string
		URL      string
		URLNew   string
	}

	PresetStore interface {
		List(ctx context.Context) ([]Preset, error)
		Get(ctx context.Context, url string) (*Preset, error)
		Update(ctx context.Context, preset *Preset) error
	}
)

func ParsePreset(urL, titleNew, urlNew string) (*Preset, error) {
	urlObj, err := url.ParseRequestURI(urL)
	if err != nil {
		return nil, err
	}

	if urlObj.Scheme == "" || urlObj.Host == "" {
		return nil, fmt.Errorf("invalid URL: %s", urL)
	}

	return &Preset{
		URL:      urL,
		TitleNew: titleNew,
		URLNew:   urlNew,
	}, nil
}
