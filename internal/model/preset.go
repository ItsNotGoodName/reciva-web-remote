package model

import (
	"fmt"
	"net/url"
	"strings"
)

type Preset struct {
	TitleNew string `json:"title_new" validate:"required"` // TitleNew is the overridden title.
	URL      string `json:"url" validate:"required"`       // URL of the preset.
	URLNew   string `json:"url_new" validate:"required"`   // URLNew is the overridden URL.
}

func (p *Preset) ParseAndValidate() error {
	p.TitleNew = strings.TrimSpace(p.TitleNew)

	urlObj, err := url.ParseRequestURI(p.URL)
	if err != nil {
		return err
	}

	if urlObj.Scheme == "" || urlObj.Host == "" {
		return fmt.Errorf("invalid URL: %s", p.URL)
	}

	return nil
}

func NewPreset(urL, titleNew, urlNew string) (*Preset, error) {
	p := &Preset{
		URL:      urL,
		TitleNew: titleNew,
		URLNew:   urlNew,
	}

	if err := p.ParseAndValidate(); err != nil {
		return nil, err
	}

	return p, nil
}
